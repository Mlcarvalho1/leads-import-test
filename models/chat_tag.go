package models

import "time"

type ChatTag struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	ChatID      string    `json:"chat_id" gorm:"not null"`
	TagID       int       `json:"tag_id" gorm:"not null"`
	LeadID      int       `json:"lead_id" gorm:"not null"`
	CompanyID   int       `json:"company_id" gorm:"not null"`
	CreatorID   int       `json:"creator_id" gorm:"not null"`
	DestroyerID *int      `json:"destroyer_id"`
	IsDeleted   bool      `json:"is_deleted" gorm:"default:false;not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"not null"`
}

func (ChatTag) TableName() string {
	return "amigocare.chat_tags"
}
