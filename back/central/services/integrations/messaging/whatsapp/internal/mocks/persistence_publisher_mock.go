package mocks

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/entities"
)

// PersistencePublisherMock implementa ports.IPersistencePublisher para tests unitarios
type PersistencePublisherMock struct {
	PublishConversationCreatedFn func(ctx context.Context, conversation *entities.Conversation) error
	PublishConversationUpdatedFn func(ctx context.Context, conversation *entities.Conversation) error
	PublishConversationExpiredFn func(ctx context.Context, conversationID string) error
	PublishMessageLogCreatedFn   func(ctx context.Context, messageLog *entities.MessageLog) error
	PublishMessageStatusUpdatedFn func(ctx context.Context, messageID string, status entities.MessageStatus, timestamps map[string]time.Time) error
}

func (m *PersistencePublisherMock) PublishConversationCreated(ctx context.Context, conversation *entities.Conversation) error {
	if m.PublishConversationCreatedFn != nil {
		return m.PublishConversationCreatedFn(ctx, conversation)
	}
	return nil
}

func (m *PersistencePublisherMock) PublishConversationUpdated(ctx context.Context, conversation *entities.Conversation) error {
	if m.PublishConversationUpdatedFn != nil {
		return m.PublishConversationUpdatedFn(ctx, conversation)
	}
	return nil
}

func (m *PersistencePublisherMock) PublishConversationExpired(ctx context.Context, conversationID string) error {
	if m.PublishConversationExpiredFn != nil {
		return m.PublishConversationExpiredFn(ctx, conversationID)
	}
	return nil
}

func (m *PersistencePublisherMock) PublishMessageLogCreated(ctx context.Context, messageLog *entities.MessageLog) error {
	if m.PublishMessageLogCreatedFn != nil {
		return m.PublishMessageLogCreatedFn(ctx, messageLog)
	}
	return nil
}

func (m *PersistencePublisherMock) PublishMessageStatusUpdated(ctx context.Context, messageID string, status entities.MessageStatus, timestamps map[string]time.Time) error {
	if m.PublishMessageStatusUpdatedFn != nil {
		return m.PublishMessageStatusUpdatedFn(ctx, messageID, status, timestamps)
	}
	return nil
}
