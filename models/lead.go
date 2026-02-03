package models

import "gorm.io/gorm"

// Lead represents an imported lead
type Lead struct {
	gorm.Model
	Name   string `json:"name" gorm:"index"`
	Email  string `json:"email" gorm:"index"`
	Phone  string `json:"phone"`
	Source string `json:"source"` // e.g. "csv_import"
}
