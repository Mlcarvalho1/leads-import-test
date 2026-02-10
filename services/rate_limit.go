package services

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

func CheckRateLimit(db *gorm.DB, companyID int, accountID int) error {
	oneHourAgo := time.Now().Add(-1 * time.Hour)

	var importCount int64
	err := db.Table("amigocare.lead_imports").
		Where("company_id = ? AND account_id = ? AND created_at > ? AND is_deleted = false", companyID, accountID, oneHourAgo).
		Count(&importCount).Error
	if err != nil {
		return fmt.Errorf("failed to check rate limit: %w", err)
	}

	if importCount >= 5 {
		return fmt.Errorf("rate limit exceeded: max 5 imports per hour per account")
	}

	var totalCreated int64
	err = db.Table("amigocare.lead_imports").
		Where("company_id = ? AND account_id = ? AND created_at > ? AND is_deleted = false", companyID, accountID, oneHourAgo).
		Select("COALESCE(SUM(total_created), 0)").
		Scan(&totalCreated).Error
	if err != nil {
		return fmt.Errorf("failed to check rate limit: %w", err)
	}

	if totalCreated >= 5000 {
		return fmt.Errorf("rate limit exceeded: max 5000 leads per hour per account")
	}

	return nil
}
