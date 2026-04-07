package ports

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
)

// IRepository define el contrato del repositorio de configuraciones de notificaciones
type IRepository interface {
	// Create crea una nueva configuración
	Create(ctx context.Context, config *entities.IntegrationNotificationConfig) error

	// Update actualiza una configuración existente
	Update(ctx context.Context, config *entities.IntegrationNotificationConfig) error

	// GetByID obtiene una configuración por su ID
	GetByID(ctx context.Context, id uint) (*entities.IntegrationNotificationConfig, error)

	// List obtiene una lista de configuraciones con filtros opcionales
	List(ctx context.Context, filters dtos.FilterNotificationConfigDTO) ([]entities.IntegrationNotificationConfig, error)

	// Delete elimina una configuración por su ID
	Delete(ctx context.Context, id uint) error

	// GetActiveConfigsByIntegrationAndTrigger obtiene configuraciones activas por integración y trigger
	// Ordenadas por prioridad descendente
	GetActiveConfigsByIntegrationAndTrigger(ctx context.Context, integrationID uint, trigger string) ([]entities.IntegrationNotificationConfig, error)

	// SyncConfigs ejecuta create/update/delete en una transacción atómica
	SyncConfigs(ctx context.Context, businessID uint, integrationID uint,
		toCreate []*entities.IntegrationNotificationConfig,
		toUpdate []*entities.IntegrationNotificationConfig,
		toDeleteIDs []uint,
	) error
}

// IOrderStatusQuerier define el contrato para consultar estados de orden
// Replicado localmente (no compartir repositorios entre módulos)
type IOrderStatusQuerier interface {
	// GetOrderStatusCodesByIDs retorna un map de id->code para los IDs dados
	GetOrderStatusCodesByIDs(ctx context.Context, ids []uint) (map[uint]string, error)
}

// IMessageAuditQuerier define el contrato para consultar logs de auditoría de mensajes
// Consulta las tablas whatsapp_message_logs, whatsapp_conversations y email_logs (replicado localmente)
type IMessageAuditQuerier interface {
	// ListMessageLogs obtiene logs de mensajes WhatsApp con filtros y paginación
	ListMessageLogs(ctx context.Context, filter dtos.MessageAuditFilterDTO) ([]entities.MessageAuditLog, int64, error)

	// GetMessageStats obtiene estadísticas agregadas de mensajes outbound (WhatsApp)
	GetMessageStats(ctx context.Context, businessID uint, dateFrom, dateTo *string) (*entities.MessageAuditStats, error)

	// ListEmailLogs obtiene logs de entregas de email con filtros y paginación
	ListEmailLogs(ctx context.Context, businessID uint, status *string, dateFrom, dateTo *string, page, pageSize int) ([]entities.EmailDeliveryLog, int64, error)

	// ListConversations obtiene conversaciones con resumen para la vista de lista
	ListConversations(ctx context.Context, filter dtos.ConversationListFilterDTO) ([]entities.ConversationSummary, int64, error)

	// GetConversationMessages obtiene los mensajes de una conversación específica para la vista de chat
	// businessID se usa para validar que la conversación pertenece al negocio
	GetConversationMessages(ctx context.Context, conversationID string, businessID uint) (*entities.ConversationSummary, []entities.ConversationMessage, error)
}

// IDeliveryLogRepository persiste logs de entrega de notificaciones (email, SMS, etc.)
// Replicado localmente para que notification_config sea el dueño de la tabla email_logs
type IDeliveryLogRepository interface {
	// CreateEmailLog persiste un log de entrega de email
	CreateEmailLog(ctx context.Context, log *entities.EmailDeliveryLog) error
}

// IWhatsAppPersister persiste eventos de WhatsApp (conversaciones y message logs)
// Los datos vienen desde RabbitMQ (WhatsApp module publica, notification_config consume)
type IWhatsAppPersister interface {
	// CreateConversation persiste una nueva conversación
	CreateConversation(ctx context.Context, conv *entities.WhatsAppConversation) error
	// UpdateConversation actualiza una conversación existente
	UpdateConversation(ctx context.Context, conv *entities.WhatsAppConversation) error
	// ExpireConversation marca una conversación como expirada
	ExpireConversation(ctx context.Context, id string) error
	// CreateMessageLog persiste un nuevo message log
	CreateMessageLog(ctx context.Context, log *entities.WhatsAppMessageLogEntry) error
	// UpdateMessageLogStatus actualiza el estado de un message log
	UpdateMessageLogStatus(ctx context.Context, messageID, status string, deliveredAt, readAt *string) error
}

// INotificationTypeRepository define el contrato del repositorio de tipos de notificaciones
type INotificationTypeRepository interface {
	// GetAll obtiene todos los tipos de notificaciones
	GetAll(ctx context.Context) ([]entities.NotificationType, error)

	// GetByID obtiene un tipo de notificación por su ID
	GetByID(ctx context.Context, id uint) (*entities.NotificationType, error)

	// GetByCode obtiene un tipo de notificación por su código
	GetByCode(ctx context.Context, code string) (*entities.NotificationType, error)

	// Create crea un nuevo tipo de notificación
	Create(ctx context.Context, notificationType *entities.NotificationType) error

	// Update actualiza un tipo de notificación existente
	Update(ctx context.Context, notificationType *entities.NotificationType) error

	// Delete elimina un tipo de notificación por su ID
	Delete(ctx context.Context, id uint) error
}

// INotificationEventTypeRepository define el contrato del repositorio de tipos de eventos de notificación
type INotificationEventTypeRepository interface {
	// GetByNotificationType obtiene todos los tipos de eventos de un tipo de notificación
	GetByNotificationType(ctx context.Context, notificationTypeID uint) ([]entities.NotificationEventType, error)

	// GetAll obtiene todos los tipos de eventos sin filtros
	GetAll(ctx context.Context) ([]entities.NotificationEventType, error)

	// GetByID obtiene un tipo de evento de notificación por su ID
	GetByID(ctx context.Context, id uint) (*entities.NotificationEventType, error)

	// Create crea un nuevo tipo de evento de notificación
	Create(ctx context.Context, eventType *entities.NotificationEventType) error

	// Update actualiza un tipo de evento de notificación existente
	Update(ctx context.Context, eventType *entities.NotificationEventType) error

	// Delete elimina un tipo de evento de notificación por su ID
	Delete(ctx context.Context, id uint) error
}
