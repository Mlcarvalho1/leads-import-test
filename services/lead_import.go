package services

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"leads-import/models"

	"gorm.io/gorm"
)

const ChunkSize = 100

type LeadImportService struct {
	DB       *gorm.DB
	Chats    ChatRepository
	WhatsApp WhatsAppValidator
	Events   EventEmitter
	Cache    CacheClearer
}

type StartImportInput struct {
	Request   models.ImportRequest
	Rows      []models.ParsedRow
	CompanyID int
	UserID    int
	Token     string
}

func (s *LeadImportService) StartImport(input StartImportInput) (int, error) {
	// 1. Validate source_id exists
	var source models.LeadSource
	if err := s.DB.Where("id = ? AND is_deleted = false", input.Request.SourceID).First(&source).Error; err != nil {
		return 0, fmt.Errorf("invalid source_id: source not found")
	}

	// 2. Validate account_id exists and belongs to company
	var account models.MessagingAccount
	if err := s.DB.Where("id = ? AND company_id = ? AND is_deleted = false", input.Request.AccountID, input.CompanyID).First(&account).Error; err != nil {
		return 0, fmt.Errorf("invalid account_id: account not found or does not belong to company")
	}

	// 3. Validate import name uniqueness
	var existingCount int64
	s.DB.Table("amigocare.lead_imports").
		Where("name = ? AND account_id = ? AND company_id = ? AND is_deleted = false", input.Request.Name, input.Request.AccountID, input.CompanyID).
		Count(&existingCount)
	if existingCount > 0 {
		return 0, fmt.Errorf("import name already exists for this account")
	}

	// 4. Validate tag_ids if provided
	if len(input.Request.TagIDs) > 5 {
		return 0, fmt.Errorf("max 5 tag_ids allowed")
	}
	if len(input.Request.TagIDs) > 0 {
		var tagCount int64
		s.DB.Model(&models.Tag{}).
			Where("id IN ? AND company_id = ? AND is_deleted = false", input.Request.TagIDs, input.CompanyID).
			Count(&tagCount)
		if int(tagCount) != len(input.Request.TagIDs) {
			return 0, fmt.Errorf("one or more tag_ids are invalid")
		}
	}

	// 5. Check rate limit
	if err := CheckRateLimit(s.DB, input.CompanyID, input.Request.AccountID); err != nil {
		return 0, err
	}

	// 6. Insert lead_imports record
	importRecord := models.LeadImport{
		Name:      input.Request.Name,
		Status:    models.LeadImportStatusProcessing,
		CompanyID: input.CompanyID,
		CreatorID: input.UserID,
		SourceID:  input.Request.SourceID,
		AccountID: input.Request.AccountID,
	}
	if err := s.DB.Create(&importRecord).Error; err != nil {
		return 0, fmt.Errorf("failed to create import record: %w", err)
	}

	// 7. Launch async processing
	go s.processImport(importRecord.ID, input)

	return importRecord.ID, nil
}

func (s *LeadImportService) processImport(importID int, input StartImportInput) {
	ctx := context.Background()

	totalCreated := 0
	totalExisting := 0
	totalErrors := 0
	finalStatus := models.LeadImportStatusFinished

	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic in processImport: %v", r)
			finalStatus = models.LeadImportStatusFailed
		}
		s.DB.Table("amigocare.lead_imports").Where("id = ?", importID).Updates(map[string]interface{}{
			"status":         string(finalStatus),
			"total_created":  totalCreated,
			"total_existing": totalExisting,
			"total_errors":   totalErrors,
			"updated_at":     time.Now(),
		})

		_ = s.Cache.ClearLeadCache(ctx, input.CompanyID)
		_ = s.Events.Emit(ctx, "lead:import-finished", map[string]interface{}{
			"import_id":  importID,
			"company_id": input.CompanyID,
		})
	}()

	// 1. Filter duplicates
	phones := make([]string, len(input.Rows))
	for i, row := range input.Rows {
		phones[i] = row.Phone
	}

	// Check existing leads by phone
	var existingLeadPhones []string
	s.DB.Model(&models.Lead{}).
		Where("contact_cellphone IN ? AND company_id = ? AND amigocare_messaging_account_id = ? AND is_deleted = false",
			phones, input.CompanyID, input.Request.AccountID).
		Pluck("contact_cellphone", &existingLeadPhones)

	// Check existing patients by phone
	var existingPatientPhones []string
	s.DB.Model(&models.Patient{}).
		Where("contact_cellphone IN ? AND company_id = ? AND deleted_at IS NULL",
			phones, input.CompanyID).
		Pluck("contact_cellphone", &existingPatientPhones)

	// Check existing chats in MongoDB
	existingChats, err := s.Chats.FindChatsByPhones(ctx, phones, input.Request.AccountID, input.CompanyID)
	if err != nil {
		log.Printf("failed to check existing chats: %v", err)
		finalStatus = models.LeadImportStatusFailed
		return
	}
	existingChatPhones := make(map[string]bool)
	for _, chat := range existingChats {
		existingChatPhones[chat.Phone] = true
	}

	// Build duplicate set
	duplicatePhones := make(map[string]bool)
	for _, p := range existingLeadPhones {
		duplicatePhones[p] = true
	}
	for _, p := range existingPatientPhones {
		duplicatePhones[p] = true
	}
	for p := range existingChatPhones {
		duplicatePhones[p] = true
	}

	// Separate duplicates from non-duplicates
	var nonDuplicates []models.ParsedRow
	for _, row := range input.Rows {
		if duplicatePhones[row.Phone] {
			totalExisting++
		} else {
			nonDuplicates = append(nonDuplicates, row)
		}
	}

	if len(nonDuplicates) == 0 {
		return
	}

	// 2. Resolve tags
	allTagNames := make(map[string]bool)
	for _, row := range nonDuplicates {
		for _, t := range row.TagNames {
			allTagNames[t] = true
		}
	}

	tagNameToID := make(map[string]int)
	for name := range allTagNames {
		var tag models.Tag
		err := s.DB.Where("LOWER(name) = LOWER(?) AND company_id = ? AND is_deleted = false", name, input.CompanyID).First(&tag).Error
		if err != nil {
			tag = models.Tag{
				Name:      name,
				CompanyID: input.CompanyID,
				CreatorID: input.UserID,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			if err := s.DB.Create(&tag).Error; err != nil {
				log.Printf("failed to create tag '%s': %v", name, err)
				continue
			}
		}
		tagNameToID[strings.ToLower(name)] = tag.ID
	}

	// Add tag_ids from request
	var requestTags []models.Tag
	if len(input.Request.TagIDs) > 0 {
		s.DB.Where("id IN ? AND is_deleted = false", input.Request.TagIDs).Find(&requestTags)
	}
	requestTagIDs := make([]int, 0)
	for _, t := range requestTags {
		requestTagIDs = append(requestTagIDs, t.ID)
	}

	// 3. Get IMPORT channel ID
	var importChannel models.LeadChannel
	if err := s.DB.Where("LOWER(name) = 'import' AND is_deleted = false").First(&importChannel).Error; err != nil {
		log.Printf("failed to find IMPORT channel: %v", err)
		finalStatus = models.LeadImportStatusFailed
		return
	}

	// 4. Process non-duplicates in chunks
	for i := 0; i < len(nonDuplicates); i += ChunkSize {
		end := i + ChunkSize
		if end > len(nonDuplicates) {
			end = len(nonDuplicates)
		}
		chunk := nonDuplicates[i:end]

		for _, row := range chunk {
			valid, err := s.WhatsApp.ValidatePhone(ctx, row.Phone, input.Request.AccountID)
			if err != nil {
				log.Printf("WhatsApp validation error for %s: %v", row.Phone, err)
				totalErrors++
				continue
			}
			if !valid {
				totalErrors++
				continue
			}

			chatID, err := s.Chats.CreateChat(ctx, row.Phone, row.DialCode, row.CountryCode, input.Request.AccountID, input.CompanyID)
			if err != nil {
				log.Printf("failed to create chat for %s: %v", row.Phone, err)
				totalErrors++
				continue
			}

			var namePtr *string
			if row.Name != "" {
				n := row.Name
				namePtr = &n
			}
			var emailPtr *string
			if row.Email != "" {
				e := row.Email
				emailPtr = &e
			}
			var cpfPtr *string
			if row.CPF != "" {
				c := row.CPF
				cpfPtr = &c
			}

			lead := models.Lead{
				Name:                        namePtr,
				Email:                       emailPtr,
				CPF:                         cpfPtr,
				ContactCellphone:            row.Phone,
				ContactCellphoneDialCode:    row.DialCode,
				ContactCellphoneCountryCode: row.CountryCode,
				SourceID:                    input.Request.SourceID,
				ChannelID:                   importChannel.ID,
				ChatID:                      &chatID,
				ImportID:                    importID,
				CompanyID:                   input.CompanyID,
				AmigocareMessagingAccountID: input.Request.AccountID,
				CreatorID:                   input.UserID,
				IsDeleted:                   false,
				CreatedAt:                   time.Now(),
				UpdatedAt:                   time.Now(),
			}

			if err := s.DB.Create(&lead).Error; err != nil {
				log.Printf("failed to create lead for %s: %v", row.Phone, err)
				totalErrors++
				continue
			}

			if err := s.Chats.UpdateChatLeadID(ctx, chatID, lead.ID); err != nil {
				log.Printf("failed to update chat lead ID: %v", err)
			}

			// Create chat_tags
			rowTagIDs := make([]int, 0)
			for _, tagName := range row.TagNames {
				if id, ok := tagNameToID[strings.ToLower(tagName)]; ok {
					rowTagIDs = append(rowTagIDs, id)
				}
			}
			rowTagIDs = append(rowTagIDs, requestTagIDs...)

			seen := make(map[int]bool)
			for _, id := range rowTagIDs {
				if seen[id] {
					continue
				}
				seen[id] = true
				chatTag := models.ChatTag{
					ChatID:    chatID,
					TagID:     id,
					LeadID:    lead.ID,
					CompanyID: input.CompanyID,
					CreatorID: input.UserID,
					IsDeleted: false,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				if err := s.DB.Create(&chatTag).Error; err != nil {
					log.Printf("failed to create chat_tag: %v", err)
				}
			}

			totalCreated++
		}
	}
}
