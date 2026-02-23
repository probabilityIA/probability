package mocks

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/entities"
)

// ConversationRepositoryMock implementa ports.IConversationRepository para tests unitarios
type ConversationRepositoryMock struct {
	CreateFn              func(ctx context.Context, conversation *entities.Conversation) error
	GetByIDFn             func(ctx context.Context, id string) (*entities.Conversation, error)
	GetByPhoneAndOrderFn  func(ctx context.Context, phoneNumber, orderNumber string) (*entities.Conversation, error)
	GetActiveByPhoneFn    func(ctx context.Context, phoneNumber string) (*entities.Conversation, error)
	UpdateFn              func(ctx context.Context, conversation *entities.Conversation) error
	ExpireFn              func(ctx context.Context, id string) error
	DeleteFn              func(ctx context.Context, id string) error
}

func (m *ConversationRepositoryMock) Create(ctx context.Context, conversation *entities.Conversation) error {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, conversation)
	}
	return nil
}

func (m *ConversationRepositoryMock) GetByID(ctx context.Context, id string) (*entities.Conversation, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *ConversationRepositoryMock) GetByPhoneAndOrder(ctx context.Context, phoneNumber, orderNumber string) (*entities.Conversation, error) {
	if m.GetByPhoneAndOrderFn != nil {
		return m.GetByPhoneAndOrderFn(ctx, phoneNumber, orderNumber)
	}
	return nil, nil
}

func (m *ConversationRepositoryMock) GetActiveByPhone(ctx context.Context, phoneNumber string) (*entities.Conversation, error) {
	if m.GetActiveByPhoneFn != nil {
		return m.GetActiveByPhoneFn(ctx, phoneNumber)
	}
	return nil, nil
}

func (m *ConversationRepositoryMock) Update(ctx context.Context, conversation *entities.Conversation) error {
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, conversation)
	}
	return nil
}

func (m *ConversationRepositoryMock) Expire(ctx context.Context, id string) error {
	if m.ExpireFn != nil {
		return m.ExpireFn(ctx, id)
	}
	return nil
}

func (m *ConversationRepositoryMock) Delete(ctx context.Context, id string) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, id)
	}
	return nil
}
