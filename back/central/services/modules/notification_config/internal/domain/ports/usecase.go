package ports

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
)

type IUseCase interface {
	Create(ctx context.Context, dto dtos.CreateNotificationConfigDTO) (*dtos.NotificationConfigResponseDTO, error)

	Update(ctx context.Context, id uint, dto dtos.UpdateNotificationConfigDTO) (*dtos.NotificationConfigResponseDTO, error)

	GetByID(ctx context.Context, id uint) (*dtos.NotificationConfigResponseDTO, error)

	List(ctx context.Context, filters dtos.FilterNotificationConfigDTO) ([]dtos.NotificationConfigResponseDTO, error)

	Delete(ctx context.Context, id uint) error

	ValidateConditions(config *entities.IntegrationNotificationConfig, orderStatusID uint, paymentMethodID uint) bool

	SyncByIntegration(ctx context.Context, dto dtos.SyncNotificationConfigsDTO) (*dtos.SyncNotificationConfigsResponseDTO, error)

	GetNotificationTypes(ctx context.Context) ([]entities.NotificationType, error)

	GetNotificationTypeByID(ctx context.Context, id uint) (*entities.NotificationType, error)

	GetNotificationTypeByCode(ctx context.Context, code string) (*entities.NotificationType, error)

	CreateNotificationType(ctx context.Context, notificationType *entities.NotificationType) error

	UpdateNotificationType(ctx context.Context, notificationType *entities.NotificationType) error

	DeleteNotificationType(ctx context.Context, id uint) error

	GetEventTypesByNotificationType(ctx context.Context, notificationTypeID uint) ([]entities.NotificationEventType, error)

	ListAllEventTypes(ctx context.Context) ([]entities.NotificationEventType, error)

	GetNotificationEventTypeByID(ctx context.Context, id uint) (*entities.NotificationEventType, error)

	CreateNotificationEventType(ctx context.Context, eventType *entities.NotificationEventType) error

	UpdateNotificationEventType(ctx context.Context, eventType *entities.NotificationEventType) error

	DeleteNotificationEventType(ctx context.Context, id uint) error

	ListMessageAudit(ctx context.Context, filter dtos.MessageAuditFilterDTO) (*dtos.PaginatedMessageAuditResponseDTO, error)

	GetMessageAuditStats(ctx context.Context, businessID uint, dateFrom, dateTo *string) (*dtos.MessageAuditStatsResponseDTO, error)

	ListConversations(ctx context.Context, filter dtos.ConversationListFilterDTO) (*dtos.PaginatedConversationListResponseDTO, error)

	GetConversationMessages(ctx context.Context, conversationID string, businessID uint) (*dtos.ConversationDetailResponseDTO, error)
	MarkConversationRead(ctx context.Context, conversationID string, businessID uint, userID *uint) error
}
