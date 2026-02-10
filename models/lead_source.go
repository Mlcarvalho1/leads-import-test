package models

type LeadSource struct {
	ID        int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string `json:"name" gorm:"type:varchar(255)"`
	IsDeleted bool   `json:"is_deleted" gorm:"default:false"`
}

func (LeadSource) TableName() string {
	return "amigocare.lead_sources"
}
