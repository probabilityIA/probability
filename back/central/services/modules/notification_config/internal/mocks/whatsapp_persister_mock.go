package mocks

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
)

// WhatsAppPersisterMock - Mock de IWhatsAppPersister para testing
type WhatsAppPersisterMock struct {
	CreateConversationFn     func(ctx context.Context, conv *entities.WhatsAppConversation) error
	UpdateConversationFn     func(ctx context.Context, conv *entities.WhatsAppConversation) error
	ExpireConversationFn     func(ctx context.Context, id string) error
	CreateMessageLogFn       func(ctx context.Context, log *entities.WhatsAppMessageLogEntry) error
	UpdateMessageLogStatusFn func(ctx context.Context, messageID, status string, deliveredAt, readAt *string) error
}

func (m *WhatsAppPersisterMock) CreateConversation(ctx context.Context, conv *entities.WhatsAppConversation) error {
	if m.CreateConversationFn != nil {
		return m.CreateConversationFn(ctx, conv)
	}
	return nil
}

func (m *WhatsAppPersisterMock) UpdateConversation(ctx context.Context, conv *entities.WhatsAppConversation) error {
	if m.UpdateConversationFn != nil {
		return m.UpdateConversationFn(ctx, conv)
	}
	return nil
}

func (m *WhatsAppPersisterMock) ExpireConversation(ctx context.Context, id string) error {
	if m.ExpireConversationFn != nil {
		return m.ExpireConversationFn(ctx, id)
	}
	return nil
}

func (m *WhatsAppPersisterMock) CreateMessageLog(ctx context.Context, log *entities.WhatsAppMessageLogEntry) error {
	if m.CreateMessageLogFn != nil {
		return m.CreateMessageLogFn(ctx, log)
	}
	return nil
}

func (m *WhatsAppPersisterMock) UpdateMessageLogStatus(ctx context.Context, messageID, status string, deliveredAt, readAt *string) error {
	if m.UpdateMessageLogStatusFn != nil {
		return m.UpdateMessageLogStatusFn(ctx, messageID, status, deliveredAt, readAt)
	}
	return nil
}
