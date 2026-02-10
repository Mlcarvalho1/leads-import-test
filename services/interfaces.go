package services

import "context"

type Chat struct {
	ID        string
	Phone     string
	AccountID int
	CompanyID int
	LeadID    *int
}

type ChatRepository interface {
	FindChatsByPhones(ctx context.Context, phones []string, accountID int, companyID int) ([]Chat, error)
	CreateChat(ctx context.Context, phone string, dialCode string, countryCode string, accountID int, companyID int) (string, error)
	UpdateChatLeadID(ctx context.Context, chatID string, leadID int) error
}

type WhatsAppValidator interface {
	ValidatePhone(ctx context.Context, phone string, accountID int) (bool, error)
}

type EventEmitter interface {
	Emit(ctx context.Context, event string, data interface{}) error
}

type CacheClearer interface {
	ClearLeadCache(ctx context.Context, companyID int) error
}
