package models

type ParsedRow struct {
	Name        string
	Phone       string
	CPF         string
	Email       string
	TagNames    []string
	DialCode    string
	CountryCode string
}
