package mocks

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/entities"
)

// MessageLogRepositoryMock implementa ports.IMessageLogRepository para tests unitarios
type MessageLogRepositoryMock struct {
	CreateFn           func(ctx context.Context, messageLog *entities.MessageLog) error
	GetByIDFn          func(ctx context.Context, id string) (*entities.MessageLog, error)
	GetByMessageIDFn   func(ctx context.Context, messageID string) (*entities.MessageLog, error)
	GetByConversationFn func(ctx context.Context, conversationID string) ([]entities.MessageLog, error)
	UpdateStatusFn     func(ctx context.Context, messageID string, status entities.MessageStatus, timestamps map[string]time.Time) error
	DeleteFn           func(ctx context.Context, id string) error
}

func (m *MessageLogRepositoryMock) Create(ctx context.Context, messageLog *entities.MessageLog) error {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, messageLog)
	}
	return nil
}

func (m *MessageLogRepositoryMock) GetByID(ctx context.Context, id string) (*entities.MessageLog, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *MessageLogRepositoryMock) GetByMessageID(ctx context.Context, messageID string) (*entities.MessageLog, error) {
	if m.GetByMessageIDFn != nil {
		return m.GetByMessageIDFn(ctx, messageID)
	}
	return nil, nil
}

func (m *MessageLogRepositoryMock) GetByConversation(ctx context.Context, conversationID string) ([]entities.MessageLog, error) {
	if m.GetByConversationFn != nil {
		return m.GetByConversationFn(ctx, conversationID)
	}
	return nil, nil
}

func (m *MessageLogRepositoryMock) UpdateStatus(ctx context.Context, messageID string, status entities.MessageStatus, timestamps map[string]time.Time) error {
	if m.UpdateStatusFn != nil {
		return m.UpdateStatusFn(ctx, messageID, status, timestamps)
	}
	return nil
}

func (m *MessageLogRepositoryMock) Delete(ctx context.Context, id string) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, id)
	}
	return nil
}
