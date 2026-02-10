package models

import "time"

type Tag struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"type:varchar(255)"`
	IsDeleted   bool      `json:"is_deleted" gorm:"default:false"`
	CompanyID   int       `json:"company_id" gorm:"not null"`
	CreatorID   int       `json:"creator_id" gorm:"not null"`
	DestroyerID *int      `json:"destroyer_id"`
	CreatedAt   time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"not null"`
}

func (Tag) TableName() string {
	return "amigocare.tags"
}
