package mocks

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/entities"
)

// ConversationCacheMock implementa ports.IConversationCache para tests unitarios
type ConversationCacheMock struct {
	GetByIDFn           func(ctx context.Context, id string) (*entities.Conversation, error)
	GetByPhoneAndOrderFn func(ctx context.Context, phoneNumber, orderNumber string) (*entities.Conversation, error)
	GetActiveByPhoneFn  func(ctx context.Context, phoneNumber string) (*entities.Conversation, error)
	SaveFn              func(ctx context.Context, conversation *entities.Conversation) error
	ExpireFn            func(ctx context.Context, id string) error
}

func (m *ConversationCacheMock) GetByID(ctx context.Context, id string) (*entities.Conversation, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *ConversationCacheMock) GetByPhoneAndOrder(ctx context.Context, phoneNumber, orderNumber string) (*entities.Conversation, error) {
	if m.GetByPhoneAndOrderFn != nil {
		return m.GetByPhoneAndOrderFn(ctx, phoneNumber, orderNumber)
	}
	return nil, nil
}

func (m *ConversationCacheMock) GetActiveByPhone(ctx context.Context, phoneNumber string) (*entities.Conversation, error) {
	if m.GetActiveByPhoneFn != nil {
		return m.GetActiveByPhoneFn(ctx, phoneNumber)
	}
	return nil, nil
}

func (m *ConversationCacheMock) Save(ctx context.Context, conversation *entities.Conversation) error {
	if m.SaveFn != nil {
		return m.SaveFn(ctx, conversation)
	}
	return nil
}

func (m *ConversationCacheMock) Expire(ctx context.Context, id string) error {
	if m.ExpireFn != nil {
		return m.ExpireFn(ctx, id)
	}
	return nil
}
