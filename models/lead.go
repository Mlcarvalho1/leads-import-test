package models

import "time"

type Lead struct {
	ID                          int        `json:"id" gorm:"primaryKey;autoIncrement"`
	Name                        *string    `json:"name" gorm:"type:varchar(255)"`
	Email                       *string    `json:"email" gorm:"type:varchar(255)"`
	CPF                         *string    `json:"cpf" gorm:"type:varchar(11)"`
	ContactCellphone            string     `json:"contact_cellphone" gorm:"type:varchar(25);not null"`
	ContactCellphoneDialCode    string     `json:"contact_cellphone_dial_code" gorm:"type:varchar(25);default:'55'"`
	ContactCellphoneCountryCode string     `json:"contact_cellphone_country_code" gorm:"type:varchar(25);default:'BR'"`
	SourceID                    int        `json:"source_id" gorm:"not null"`
	ChannelID                   int        `json:"channel_id" gorm:"not null"`
	ChatID                      *string    `json:"chat_id"`
	ImportID                    int        `json:"import_id" gorm:"not null"`
	CompanyID                   int        `json:"company_id" gorm:"not null"`
	AmigocareMessagingAccountID int        `json:"amigocare_messaging_account_id" gorm:"not null"`
	CreatorID                   int        `json:"creator_id" gorm:"not null"`
	PatientID                   *int       `json:"patient_id"`
	AttendanceID                *int       `json:"attendance_id"`
	ConvertedBy                 *int       `json:"converted_by"`
	ConvertedAt                 *time.Time `json:"converted_at"`
	ConvertedFrom               *string    `json:"converted_from"`
	IsDeleted                   bool       `json:"is_deleted" gorm:"default:false;not null"`
	CreatedAt                   time.Time  `json:"created_at" gorm:"not null"`
	UpdatedAt                   time.Time  `json:"updated_at" gorm:"not null"`
}

func (Lead) TableName() string {
	return "amigocare.amigocare_leads"
}
