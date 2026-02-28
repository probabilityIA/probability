package ports

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
)

// IUseCase define el contrato de los casos de uso
type IUseCase interface {
	// ========== Notification Configs ==========
	// Create crea una nueva configuración
	Create(ctx context.Context, dto dtos.CreateNotificationConfigDTO) (*dtos.NotificationConfigResponseDTO, error)

	// Update actualiza una configuración existente
	Update(ctx context.Context, id uint, dto dtos.UpdateNotificationConfigDTO) (*dtos.NotificationConfigResponseDTO, error)

	// GetByID obtiene una configuración por su ID
	GetByID(ctx context.Context, id uint) (*dtos.NotificationConfigResponseDTO, error)

	// List obtiene una lista de configuraciones con filtros
	List(ctx context.Context, filters dtos.FilterNotificationConfigDTO) ([]dtos.NotificationConfigResponseDTO, error)

	// Delete elimina una configuración
	Delete(ctx context.Context, id uint) error

	// ValidateConditions valida si una orden cumple las condiciones de una configuración
	// NUEVA ESTRUCTURA: Usa orderStatusID (uint) en lugar de orderStatus (string)
	ValidateConditions(config *entities.IntegrationNotificationConfig, orderStatusID uint, paymentMethodID uint) bool

	// ========== Sync ==========
	// SyncByIntegration sincroniza las reglas de notificación para una integración
	// Crea, actualiza y elimina configs en una sola operación transaccional
	SyncByIntegration(ctx context.Context, dto dtos.SyncNotificationConfigsDTO) (*dtos.SyncNotificationConfigsResponseDTO, error)

	// ========== Notification Types ==========
	// GetNotificationTypes obtiene todos los tipos de notificaciones
	GetNotificationTypes(ctx context.Context) ([]entities.NotificationType, error)

	// GetNotificationTypeByID obtiene un tipo de notificación por su ID
	GetNotificationTypeByID(ctx context.Context, id uint) (*entities.NotificationType, error)

	// GetNotificationTypeByCode obtiene un tipo de notificación por su código
	GetNotificationTypeByCode(ctx context.Context, code string) (*entities.NotificationType, error)

	// CreateNotificationType crea un nuevo tipo de notificación
	CreateNotificationType(ctx context.Context, notificationType *entities.NotificationType) error

	// UpdateNotificationType actualiza un tipo de notificación existente
	UpdateNotificationType(ctx context.Context, notificationType *entities.NotificationType) error

	// DeleteNotificationType elimina un tipo de notificación por su ID
	DeleteNotificationType(ctx context.Context, id uint) error

	// ========== Notification Event Types ==========
	// GetEventTypesByNotificationType obtiene todos los tipos de eventos de un tipo de notificación
	GetEventTypesByNotificationType(ctx context.Context, notificationTypeID uint) ([]entities.NotificationEventType, error)

	// ListAllEventTypes obtiene todos los tipos de eventos de notificación sin filtros
	ListAllEventTypes(ctx context.Context) ([]entities.NotificationEventType, error)

	// GetNotificationEventTypeByID obtiene un tipo de evento de notificación por su ID
	GetNotificationEventTypeByID(ctx context.Context, id uint) (*entities.NotificationEventType, error)

	// CreateNotificationEventType crea un nuevo tipo de evento de notificación
	CreateNotificationEventType(ctx context.Context, eventType *entities.NotificationEventType) error

	// UpdateNotificationEventType actualiza un tipo de evento de notificación existente
	UpdateNotificationEventType(ctx context.Context, eventType *entities.NotificationEventType) error

	// DeleteNotificationEventType elimina un tipo de evento de notificación por su ID
	DeleteNotificationEventType(ctx context.Context, id uint) error
}
