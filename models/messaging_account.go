package models

type MessagingAccount struct {
	ID        int  `json:"id" gorm:"primaryKey"`
	CompanyID int  `json:"company_id" gorm:"not null"`
	IsDeleted bool `json:"is_deleted" gorm:"default:false"`
}

func (MessagingAccount) TableName() string {
	return "amigocare.messaging_accounts"
}
