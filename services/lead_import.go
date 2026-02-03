package services

import (
	"encoding/csv"
	"io"
	"strings"
	"your-app/database"
	"your-app/models"
)

// ImportLeadsFromCSV reads a CSV from r and bulk-inserts leads.
// Expected CSV columns: name, email, phone (header optional; order: name, email, phone)
func ImportLeadsFromCSV(r io.Reader) (imported int, err error) {
	reader := csv.NewReader(r)
	reader.TrimLeadingSpace = true
	reader.FieldsPerRecord = -1

	rows, err := reader.ReadAll()
	if err != nil {
		return 0, err
	}

	if len(rows) == 0 {
		return 0, nil
	}

	// Skip header if it looks like a header (first cell is "name" or "nome" etc)
	start := 0
	if len(rows) > 0 {
		first := strings.ToLower(strings.TrimSpace(rows[0][0]))
		if first == "name" || first == "nome" || first == "email" {
			start = 1
		}
	}

	var leads []models.Lead
	for i := start; i < len(rows); i++ {
		row := rows[i]
		if len(row) < 2 {
			continue
		}
		name := strings.TrimSpace(row[0])
		email := strings.TrimSpace(row[1])
		phone := ""
		if len(row) > 2 {
			phone = strings.TrimSpace(row[2])
		}
		if name == "" && email == "" {
			continue
		}
		leads = append(leads, models.Lead{
			Name:   name,
			Email:  email,
			Phone:  phone,
			Source: "csv_import",
		})
	}

	if len(leads) == 0 {
		return 0, nil
	}

	if err := database.GetDB().CreateInBatches(leads, 100).Error; err != nil {
		return 0, err
	}
	return len(leads), nil
}
