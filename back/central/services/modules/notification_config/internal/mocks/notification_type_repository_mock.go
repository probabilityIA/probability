package mocks

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
)

// NotificationTypeRepositoryMock - Mock del repositorio de tipos de notificaciones
type NotificationTypeRepositoryMock struct {
	GetAllFn   func(ctx context.Context) ([]entities.NotificationType, error)
	GetByIDFn  func(ctx context.Context, id uint) (*entities.NotificationType, error)
	GetByCodeFn func(ctx context.Context, code string) (*entities.NotificationType, error)
	CreateFn   func(ctx context.Context, notificationType *entities.NotificationType) error
	UpdateFn   func(ctx context.Context, notificationType *entities.NotificationType) error
	DeleteFn   func(ctx context.Context, id uint) error
}

func (m *NotificationTypeRepositoryMock) GetAll(ctx context.Context) ([]entities.NotificationType, error) {
	if m.GetAllFn != nil {
		return m.GetAllFn(ctx)
	}
	return []entities.NotificationType{}, nil
}

func (m *NotificationTypeRepositoryMock) GetByID(ctx context.Context, id uint) (*entities.NotificationType, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *NotificationTypeRepositoryMock) GetByCode(ctx context.Context, code string) (*entities.NotificationType, error) {
	if m.GetByCodeFn != nil {
		return m.GetByCodeFn(ctx, code)
	}
	return nil, nil
}

func (m *NotificationTypeRepositoryMock) Create(ctx context.Context, notificationType *entities.NotificationType) error {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, notificationType)
	}
	return nil
}

func (m *NotificationTypeRepositoryMock) Update(ctx context.Context, notificationType *entities.NotificationType) error {
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, notificationType)
	}
	return nil
}

func (m *NotificationTypeRepositoryMock) Delete(ctx context.Context, id uint) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, id)
	}
	return nil
}
