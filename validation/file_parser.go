package validation

import (
	"encoding/csv"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"

	"leads-import/models"

	"github.com/xuri/excelize/v2"
)

type RowError struct {
	Row     int    `json:"row"`
	Column  string `json:"column"`
	Message string `json:"message"`
}

func ParseFile(file multipart.File, filename string) ([]models.ParsedRow, []RowError, error) {
	ext := strings.ToLower(filepath.Ext(filename))

	var rawRows [][]string
	var err error

	switch ext {
	case ".csv":
		rawRows, err = parseCSV(file)
	case ".xlsx":
		rawRows, err = parseExcel(file, filename)
	default:
		return nil, nil, fmt.Errorf("unsupported file format: %s (use .csv or .xlsx)", ext)
	}

	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse file: %w", err)
	}

	if len(rawRows) == 0 {
		return nil, nil, fmt.Errorf("file is empty")
	}

	// Validate header
	header := rawRows[0]
	if err := validateHeader(header); err != nil {
		return nil, nil, err
	}

	dataRows := rawRows[1:]
	if len(dataRows) == 0 {
		return nil, nil, fmt.Errorf("file must have at least 1 data row")
	}
	if len(dataRows) > 5000 {
		return nil, nil, fmt.Errorf("file must have at most 5000 data rows, got %d", len(dataRows))
	}

	var parsed []models.ParsedRow
	var errors []RowError

	for i, row := range dataRows {
		rowNum := i + 2 // 1-indexed, skip header

		// Pad row to 5 columns
		for len(row) < 5 {
			row = append(row, "")
		}

		name := strings.TrimSpace(row[0])
		phone := strings.TrimSpace(row[1])
		cpf := strings.TrimSpace(row[2])
		email := strings.TrimSpace(row[3])
		tagsRaw := strings.TrimSpace(row[4])

		// Name: required, max 255
		if name == "" {
			errors = append(errors, RowError{Row: rowNum, Column: "name", Message: "name is required"})
			continue
		}
		if len(name) > 255 {
			errors = append(errors, RowError{Row: rowNum, Column: "name", Message: "name must be at most 255 characters"})
			continue
		}

		// Phone: required, must be valid
		if phone == "" {
			errors = append(errors, RowError{Row: rowNum, Column: "phone", Message: "phone is required"})
			continue
		}
		phoneInfo, err := ParsePhone(phone)
		if err != nil {
			errors = append(errors, RowError{Row: rowNum, Column: "phone", Message: err.Error()})
			continue
		}

		// CPF: optional, validate if present
		validCPF := ""
		if cpf != "" {
			validCPF, err = ValidateCPF(cpf)
			if err != nil {
				errors = append(errors, RowError{Row: rowNum, Column: "cpf", Message: err.Error()})
				continue
			}
		}

		// Email: optional, validate if present
		if email != "" {
			if err := ValidateEmail(email); err != nil {
				errors = append(errors, RowError{Row: rowNum, Column: "email", Message: err.Error()})
				continue
			}
		}

		// Tags: optional, max 5, max 255 chars total
		var tagNames []string
		if tagsRaw != "" {
			if len(tagsRaw) > 255 {
				errors = append(errors, RowError{Row: rowNum, Column: "tags", Message: "tags must be at most 255 characters"})
				continue
			}
			parts := strings.Split(tagsRaw, ",")
			for _, p := range parts {
				t := strings.TrimSpace(p)
				if t != "" {
					tagNames = append(tagNames, t)
				}
			}
			if len(tagNames) > 5 {
				errors = append(errors, RowError{Row: rowNum, Column: "tags", Message: "max 5 tags per row"})
				continue
			}
		}

		parsed = append(parsed, models.ParsedRow{
			Name:        name,
			Phone:       phoneInfo.National,
			CPF:         validCPF,
			Email:       email,
			TagNames:    tagNames,
			DialCode:    phoneInfo.DialCode,
			CountryCode: phoneInfo.CountryCode,
		})
	}

	return parsed, errors, nil
}

func validateHeader(header []string) error {
	expected := []string{"name", "phone", "cpf", "email", "tags"}
	if len(header) != 5 {
		return fmt.Errorf("file must have exactly 5 columns (name, phone, cpf, email, tags), got %d", len(header))
	}
	for i, col := range header {
		if strings.TrimSpace(strings.ToLower(col)) != expected[i] {
			return fmt.Errorf("column %d must be '%s', got '%s'", i+1, expected[i], strings.TrimSpace(col))
		}
	}
	return nil
}

func parseCSV(r io.Reader) ([][]string, error) {
	reader := csv.NewReader(r)
	reader.TrimLeadingSpace = true
	reader.FieldsPerRecord = -1
	return reader.ReadAll()
}

func parseExcel(file multipart.File, filename string) ([][]string, error) {
	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, fmt.Errorf("failed to open excel file: %w", err)
	}
	defer f.Close()

	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		return nil, fmt.Errorf("excel file has no sheets")
	}

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("failed to read excel rows: %w", err)
	}

	return rows, nil
}
