package services

import (
	"encoding/csv"
	"io"
	"strings"

	"your-app/database"
	"your-app/models"
)

const BatchSize = 100

func isHeaderRow(row []string) bool {
	if len(row) == 0 {
		return false
	}
	first := strings.ToLower(strings.TrimSpace(row[0]))
	return first == "name" || first == "nome" || first == "email"
}

func ImportLeadsFromCSV(r io.Reader) (int, error) {
	reader := csv.NewReader(r)
	reader.TrimLeadingSpace = true
	reader.FieldsPerRecord = -1

	db := database.GetDB()

	imported := 0
	batch := make([]models.Lead, 0, BatchSize)

	isFirstRow := true

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return imported, err
		}

		if isFirstRow && isHeaderRow(row) {
			isFirstRow = false
			continue
		}
		isFirstRow = false

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

		batch = append(batch, models.Lead{
			Name:   name,
			Email:  email,
			Phone:  phone,
			Source: "csv_import",
		})

		if len(batch) >= BatchSize {
			if err := db.CreateInBatches(batch, BatchSize).Error; err != nil {
				return imported, err
			}
			imported += len(batch)
			batch = batch[:0] // reutiliza memória
		}
	}

	// flush final
	if len(batch) > 0 {
		if err := db.CreateInBatches(batch, BatchSize).Error; err != nil {
			return imported, err
		}
		imported += len(batch)
	}

	return imported, nil
}
