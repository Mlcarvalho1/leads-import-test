package models

type LeadChannel struct {
	ID        int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string `json:"name" gorm:"type:varchar(100)"`
	IsDeleted bool   `json:"is_deleted" gorm:"default:false"`
}

func (LeadChannel) TableName() string {
	return "amigocare.lead_channels"
}
