package models

import (
	"time"

	"gorm.io/gorm"
)

// LeadImportStatus represents the status of a lead import job
type LeadImportStatus string

const (
	LeadImportStatusFailed     LeadImportStatus = "FAILED"
	LeadImportStatusProcessing LeadImportStatus = "PROCESSING"
	LeadImportStatusFinished   LeadImportStatus = "FINISHED"
)

// LeadImport represents a bulk lead import job
type LeadImport struct {
	ID            int              `json:"id" gorm:"primaryKey;autoIncrement"`
	Name          string           `json:"name" gorm:"type:varchar(255);not null"`
	Status        LeadImportStatus `json:"status" gorm:"type:varchar(20);default:PROCESSING;not null"`
	TotalCreated  int              `json:"total_created" gorm:"not null;default:0"`
	TotalExisting int              `json:"total_existing" gorm:"not null;default:0"`
	TotalErrors   int              `json:"total_errors" gorm:"not null;default:0"`
	IsDeleted     bool             `json:"is_deleted" gorm:"default:false;not null"`
	CreatorID     int              `json:"creator_id" gorm:"not null"`
	CompanyID     int              `json:"company_id" gorm:"not null"`
	SourceID      int              `json:"source_id" gorm:"not null"`
	AccountID     int              `json:"account_id" gorm:"not null"`
	CreatedAt     time.Time        `json:"created_at" gorm:"not null"`
	UpdatedAt     time.Time        `json:"updated_at" gorm:"not null"`
}

// TableName overrides the table name to match the Sequelize model (amigocare schema)
func (LeadImport) TableName() string {
	return "amigocare.lead_imports"
}

// BeforeCreate sets timestamps (optional, GORM does this if using gorm.Model)
func (li *LeadImport) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	if li.CreatedAt.IsZero() {
		li.CreatedAt = now
	}
	if li.UpdatedAt.IsZero() {
		li.UpdatedAt = now
	}
	return nil
}

// BeforeUpdate sets UpdatedAt
func (li *LeadImport) BeforeUpdate(tx *gorm.DB) error {
	li.UpdatedAt = time.Now()
	return nil
}
