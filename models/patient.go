package models

import "time"

type Patient struct {
	ID                          int        `json:"id" gorm:"primaryKey"`
	CompanyID                   int        `json:"company_id" gorm:"not null"`
	ContactCellphone            string     `json:"contact_cellphone" gorm:"type:varchar(25)"`
	ContactCellphoneDialCode    string     `json:"contact_cellphone_dial_code" gorm:"type:varchar(25)"`
	ContactCellphoneCountryCode string     `json:"contact_cellphone_country_code" gorm:"type:varchar(25)"`
	DeletedAt                   *time.Time `json:"deleted_at"`
}

func (Patient) TableName() string {
	return "amigocare.patients"
}
