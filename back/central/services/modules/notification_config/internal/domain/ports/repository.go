package ports

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
)

type IRepository interface {
	Create(ctx context.Context, config *entities.IntegrationNotificationConfig) error

	Update(ctx context.Context, config *entities.IntegrationNotificationConfig) error

	GetByID(ctx context.Context, id uint) (*entities.IntegrationNotificationConfig, error)

	List(ctx context.Context, filters dtos.FilterNotificationConfigDTO) ([]entities.IntegrationNotificationConfig, error)

	Delete(ctx context.Context, id uint) error

	GetActiveConfigsByIntegrationAndTrigger(ctx context.Context, integrationID uint, trigger string) ([]entities.IntegrationNotificationConfig, error)

	SyncConfigs(ctx context.Context, businessID uint, integrationID uint,
		toCreate []*entities.IntegrationNotificationConfig,
		toUpdate []*entities.IntegrationNotificationConfig,
		toDeleteIDs []uint,
	) error
}

type IDefaultRulesQuerier interface {
	PlatformIntegrationID(ctx context.Context, businessID uint) (uint, error)
	OrderStatusIDsByCodes(ctx context.Context, codes []string) ([]uint, error)
}

type IOrderStatusQuerier interface {
	GetOrderStatusCodesByIDs(ctx context.Context, ids []uint) (map[uint]string, error)
}

type IMessageAuditQuerier interface {
	ListMessageLogs(ctx context.Context, filter dtos.MessageAuditFilterDTO) ([]entities.MessageAuditLog, int64, error)

	GetMessageStats(ctx context.Context, businessID uint, dateFrom, dateTo *string) (*entities.MessageAuditStats, error)

	ListEmailLogs(ctx context.Context, businessID uint, status *string, dateFrom, dateTo *string, page, pageSize int) ([]entities.EmailDeliveryLog, int64, error)

	ListConversations(ctx context.Context, filter dtos.ConversationListFilterDTO) ([]entities.ConversationSummary, int64, error)

	GetConversationMessages(ctx context.Context, conversationID string, businessID uint) (*entities.ConversationSummary, []entities.ConversationMessage, error)
	CountUnreadConversations(ctx context.Context, filter dtos.ConversationListFilterDTO) (int64, error)
	MarkConversationRead(ctx context.Context, conversationID string, businessID uint, userID *uint) error
}

type IDeliveryLogRepository interface {
	CreateEmailLog(ctx context.Context, log *entities.EmailDeliveryLog) error
}

type IWhatsAppPersister interface {
	CreateConversation(ctx context.Context, conv *entities.WhatsAppConversation) error
	UpdateConversation(ctx context.Context, conv *entities.WhatsAppConversation) error
	ExpireConversation(ctx context.Context, id string) error
	CreateMessageLog(ctx context.Context, log *entities.WhatsAppMessageLogEntry) error
	UpdateMessageLogStatus(ctx context.Context, messageID, status string, deliveredAt, readAt *string) error
}

type INotificationTypeRepository interface {
	GetAll(ctx context.Context) ([]entities.NotificationType, error)

	GetByID(ctx context.Context, id uint) (*entities.NotificationType, error)

	GetByCode(ctx context.Context, code string) (*entities.NotificationType, error)

	Create(ctx context.Context, notificationType *entities.NotificationType) error

	Update(ctx context.Context, notificationType *entities.NotificationType) error

	Delete(ctx context.Context, id uint) error
}

type INotificationEventTypeRepository interface {
	GetByNotificationType(ctx context.Context, notificationTypeID uint) ([]entities.NotificationEventType, error)

	GetAll(ctx context.Context) ([]entities.NotificationEventType, error)

	GetByID(ctx context.Context, id uint) (*entities.NotificationEventType, error)

	Create(ctx context.Context, eventType *entities.NotificationEventType) error

	Update(ctx context.Context, eventType *entities.NotificationEventType) error

	Delete(ctx context.Context, id uint) error
}

type IChatMediaSigner interface {
	PresignGet(ctx context.Context, key, filename string, ttl time.Duration) (string, error)
}

type IChatRetentionRepository interface {
	PurgeChatHistory(ctx context.Context, cutoff time.Time) (entities.ChatPurgeResult, error)
}
