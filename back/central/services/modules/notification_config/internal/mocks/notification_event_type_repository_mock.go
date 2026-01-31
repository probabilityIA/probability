package mocks

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
)

// NotificationEventTypeRepositoryMock - Mock del repositorio de tipos de eventos de notificación
type NotificationEventTypeRepositoryMock struct {
	GetByNotificationTypeFn func(ctx context.Context, notificationTypeID uint) ([]entities.NotificationEventType, error)
	GetByIDFn               func(ctx context.Context, id uint) (*entities.NotificationEventType, error)
	CreateFn                func(ctx context.Context, eventType *entities.NotificationEventType) error
	UpdateFn                func(ctx context.Context, eventType *entities.NotificationEventType) error
	DeleteFn                func(ctx context.Context, id uint) error
}

func (m *NotificationEventTypeRepositoryMock) GetByNotificationType(ctx context.Context, notificationTypeID uint) ([]entities.NotificationEventType, error) {
	if m.GetByNotificationTypeFn != nil {
		return m.GetByNotificationTypeFn(ctx, notificationTypeID)
	}
	return []entities.NotificationEventType{}, nil
}

func (m *NotificationEventTypeRepositoryMock) GetByID(ctx context.Context, id uint) (*entities.NotificationEventType, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *NotificationEventTypeRepositoryMock) Create(ctx context.Context, eventType *entities.NotificationEventType) error {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, eventType)
	}
	return nil
}

func (m *NotificationEventTypeRepositoryMock) Update(ctx context.Context, eventType *entities.NotificationEventType) error {
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, eventType)
	}
	return nil
}

func (m *NotificationEventTypeRepositoryMock) Delete(ctx context.Context, id uint) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, id)
	}
	return nil
}
