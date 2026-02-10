package services

import "context"

// NoopChatRepository is a no-op implementation of ChatRepository.
type NoopChatRepository struct{}

func (n *NoopChatRepository) FindChatsByPhones(_ context.Context, _ []string, _ int, _ int) ([]Chat, error) {
	return nil, nil
}

func (n *NoopChatRepository) CreateChat(_ context.Context, _ string, _ string, _ string, _ int, _ int) (string, error) {
	return "000000000000000000000000", nil
}

func (n *NoopChatRepository) UpdateChatLeadID(_ context.Context, _ string, _ int) error {
	return nil
}

// NoopWhatsAppValidator always returns true.
type NoopWhatsAppValidator struct{}

func (n *NoopWhatsAppValidator) ValidatePhone(_ context.Context, _ string, _ int) (bool, error) {
	return true, nil
}

// NoopEventEmitter does nothing.
type NoopEventEmitter struct{}

func (n *NoopEventEmitter) Emit(_ context.Context, _ string, _ interface{}) error {
	return nil
}

// NoopCacheClearer does nothing.
type NoopCacheClearer struct{}

func (n *NoopCacheClearer) ClearLeadCache(_ context.Context, _ int) error {
	return nil
}
