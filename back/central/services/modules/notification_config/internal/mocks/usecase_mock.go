package mocks

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
)

// UseCaseMock - Mock del caso de uso
type UseCaseMock struct {
	// Notification Configs
	CreateFn             func(ctx context.Context, dto dtos.CreateNotificationConfigDTO) (*dtos.NotificationConfigResponseDTO, error)
	UpdateFn             func(ctx context.Context, id uint, dto dtos.UpdateNotificationConfigDTO) (*dtos.NotificationConfigResponseDTO, error)
	GetByIDFn            func(ctx context.Context, id uint) (*dtos.NotificationConfigResponseDTO, error)
	ListFn               func(ctx context.Context, filters dtos.FilterNotificationConfigDTO) ([]dtos.NotificationConfigResponseDTO, error)
	DeleteFn             func(ctx context.Context, id uint) error
	ValidateConditionsFn func(config *entities.IntegrationNotificationConfig, orderStatusID uint, paymentMethodID uint) bool

	// Notification Types
	GetNotificationTypesFn      func(ctx context.Context) ([]entities.NotificationType, error)
	GetNotificationTypeByIDFn   func(ctx context.Context, id uint) (*entities.NotificationType, error)
	GetNotificationTypeByCodeFn func(ctx context.Context, code string) (*entities.NotificationType, error)
	CreateNotificationTypeFn    func(ctx context.Context, notificationType *entities.NotificationType) error
	UpdateNotificationTypeFn    func(ctx context.Context, notificationType *entities.NotificationType) error
	DeleteNotificationTypeFn    func(ctx context.Context, id uint) error

	// Sync
	SyncByIntegrationFn func(ctx context.Context, dto dtos.SyncNotificationConfigsDTO) (*dtos.SyncNotificationConfigsResponseDTO, error)

	// Notification Event Types
	GetEventTypesByNotificationTypeFn func(ctx context.Context, notificationTypeID uint) ([]entities.NotificationEventType, error)
	ListAllEventTypesFn               func(ctx context.Context) ([]entities.NotificationEventType, error)
	GetNotificationEventTypeByIDFn    func(ctx context.Context, id uint) (*entities.NotificationEventType, error)
	CreateNotificationEventTypeFn     func(ctx context.Context, eventType *entities.NotificationEventType) error
	UpdateNotificationEventTypeFn     func(ctx context.Context, eventType *entities.NotificationEventType) error
	DeleteNotificationEventTypeFn     func(ctx context.Context, id uint) error

	// Message Audit
	ListMessageAuditFn      func(ctx context.Context, filter dtos.MessageAuditFilterDTO) (*dtos.PaginatedMessageAuditResponseDTO, error)
	GetMessageAuditStatsFn  func(ctx context.Context, businessID uint, dateFrom, dateTo *string) (*dtos.MessageAuditStatsResponseDTO, error)
}

// Notification Configs
func (m *UseCaseMock) Create(ctx context.Context, dto dtos.CreateNotificationConfigDTO) (*dtos.NotificationConfigResponseDTO, error) {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, dto)
	}
	return nil, nil
}

func (m *UseCaseMock) Update(ctx context.Context, id uint, dto dtos.UpdateNotificationConfigDTO) (*dtos.NotificationConfigResponseDTO, error) {
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, id, dto)
	}
	return nil, nil
}

func (m *UseCaseMock) GetByID(ctx context.Context, id uint) (*dtos.NotificationConfigResponseDTO, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *UseCaseMock) List(ctx context.Context, filters dtos.FilterNotificationConfigDTO) ([]dtos.NotificationConfigResponseDTO, error) {
	if m.ListFn != nil {
		return m.ListFn(ctx, filters)
	}
	return []dtos.NotificationConfigResponseDTO{}, nil
}

func (m *UseCaseMock) Delete(ctx context.Context, id uint) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, id)
	}
	return nil
}

func (m *UseCaseMock) ValidateConditions(config *entities.IntegrationNotificationConfig, orderStatusID uint, paymentMethodID uint) bool {
	if m.ValidateConditionsFn != nil {
		return m.ValidateConditionsFn(config, orderStatusID, paymentMethodID)
	}
	return false
}

// Notification Types
func (m *UseCaseMock) GetNotificationTypes(ctx context.Context) ([]entities.NotificationType, error) {
	if m.GetNotificationTypesFn != nil {
		return m.GetNotificationTypesFn(ctx)
	}
	return []entities.NotificationType{}, nil
}

func (m *UseCaseMock) GetNotificationTypeByID(ctx context.Context, id uint) (*entities.NotificationType, error) {
	if m.GetNotificationTypeByIDFn != nil {
		return m.GetNotificationTypeByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *UseCaseMock) GetNotificationTypeByCode(ctx context.Context, code string) (*entities.NotificationType, error) {
	if m.GetNotificationTypeByCodeFn != nil {
		return m.GetNotificationTypeByCodeFn(ctx, code)
	}
	return nil, nil
}

func (m *UseCaseMock) CreateNotificationType(ctx context.Context, notificationType *entities.NotificationType) error {
	if m.CreateNotificationTypeFn != nil {
		return m.CreateNotificationTypeFn(ctx, notificationType)
	}
	return nil
}

func (m *UseCaseMock) UpdateNotificationType(ctx context.Context, notificationType *entities.NotificationType) error {
	if m.UpdateNotificationTypeFn != nil {
		return m.UpdateNotificationTypeFn(ctx, notificationType)
	}
	return nil
}

func (m *UseCaseMock) DeleteNotificationType(ctx context.Context, id uint) error {
	if m.DeleteNotificationTypeFn != nil {
		return m.DeleteNotificationTypeFn(ctx, id)
	}
	return nil
}

// Notification Event Types
func (m *UseCaseMock) GetEventTypesByNotificationType(ctx context.Context, notificationTypeID uint) ([]entities.NotificationEventType, error) {
	if m.GetEventTypesByNotificationTypeFn != nil {
		return m.GetEventTypesByNotificationTypeFn(ctx, notificationTypeID)
	}
	return []entities.NotificationEventType{}, nil
}

func (m *UseCaseMock) ListAllEventTypes(ctx context.Context) ([]entities.NotificationEventType, error) {
	if m.ListAllEventTypesFn != nil {
		return m.ListAllEventTypesFn(ctx)
	}
	return []entities.NotificationEventType{}, nil
}

func (m *UseCaseMock) GetNotificationEventTypeByID(ctx context.Context, id uint) (*entities.NotificationEventType, error) {
	if m.GetNotificationEventTypeByIDFn != nil {
		return m.GetNotificationEventTypeByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *UseCaseMock) CreateNotificationEventType(ctx context.Context, eventType *entities.NotificationEventType) error {
	if m.CreateNotificationEventTypeFn != nil {
		return m.CreateNotificationEventTypeFn(ctx, eventType)
	}
	return nil
}

func (m *UseCaseMock) UpdateNotificationEventType(ctx context.Context, eventType *entities.NotificationEventType) error {
	if m.UpdateNotificationEventTypeFn != nil {
		return m.UpdateNotificationEventTypeFn(ctx, eventType)
	}
	return nil
}

func (m *UseCaseMock) DeleteNotificationEventType(ctx context.Context, id uint) error {
	if m.DeleteNotificationEventTypeFn != nil {
		return m.DeleteNotificationEventTypeFn(ctx, id)
	}
	return nil
}

// Sync
func (m *UseCaseMock) SyncByIntegration(ctx context.Context, dto dtos.SyncNotificationConfigsDTO) (*dtos.SyncNotificationConfigsResponseDTO, error) {
	if m.SyncByIntegrationFn != nil {
		return m.SyncByIntegrationFn(ctx, dto)
	}
	return nil, nil
}

// Message Audit
func (m *UseCaseMock) ListMessageAudit(ctx context.Context, filter dtos.MessageAuditFilterDTO) (*dtos.PaginatedMessageAuditResponseDTO, error) {
	if m.ListMessageAuditFn != nil {
		return m.ListMessageAuditFn(ctx, filter)
	}
	return nil, nil
}

func (m *UseCaseMock) GetMessageAuditStats(ctx context.Context, businessID uint, dateFrom, dateTo *string) (*dtos.MessageAuditStatsResponseDTO, error) {
	if m.GetMessageAuditStatsFn != nil {
		return m.GetMessageAuditStatsFn(ctx, businessID, dateFrom, dateTo)
	}
	return nil, nil
}
