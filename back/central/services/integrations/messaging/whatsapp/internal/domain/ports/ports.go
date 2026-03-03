package ports

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/entities"
)

type IWhatsApp interface {
	SendMessage(ctx context.Context, phoneNumberID uint, msg entities.TemplateMessage, accessToken string) (string, error)
}

// ============================================
// CACHE INTERFACES (reemplazan repositorios DB)
// ============================================

// IConversationCache define operaciones de cache para conversaciones WhatsApp.
// Las conversaciones se mantienen en Redis como fuente primaria de estado
// en tiempo real, y se persisten asincrónicamente en DB via RabbitMQ.
type IConversationCache interface {
	// GetByID obtiene una conversación del cache por su ID
	GetByID(ctx context.Context, id string) (*entities.Conversation, error)

	// GetByPhoneAndOrder obtiene una conversación por teléfono + número de orden
	GetByPhoneAndOrder(ctx context.Context, phoneNumber, orderNumber string) (*entities.Conversation, error)

	// GetActiveByPhone obtiene la conversación activa de un teléfono
	GetActiveByPhone(ctx context.Context, phoneNumber string) (*entities.Conversation, error)

	// Save guarda una conversación en cache (clave principal + índices).
	// Si el estado es terminal, elimina el índice active.
	Save(ctx context.Context, conversation *entities.Conversation) error

	// Expire marca una conversación como expirada y limpia índices activos
	Expire(ctx context.Context, id string) error
}

// ICredentialsCache lee credenciales de WhatsApp desde Redis.
// Las claves son calentadas por integrations/core al startup.
type ICredentialsCache interface {
	// GetWhatsAppConfig obtiene credenciales de WhatsApp para un business.
	// Lee de: integration:idx:biz:{businessID}:type:2 → integration:creds:{id} + integration:meta:{id}
	GetWhatsAppConfig(ctx context.Context, businessID uint) (*WhatsAppConfig, error)

	// GetWhatsAppDefaultConfig obtiene credenciales globales de plataforma.
	// Lee de: integration:platform_creds:2
	GetWhatsAppDefaultConfig(ctx context.Context) (*WhatsAppConfig, error)
}

// ============================================
// PERSISTENCE PUBLISHER (async DB via RabbitMQ)
// ============================================

// IPersistencePublisher publica eventos para persistencia asíncrona en DB.
// notification_config consume estos eventos y los persiste.
type IPersistencePublisher interface {
	// Conversation events
	PublishConversationCreated(ctx context.Context, conversation *entities.Conversation) error
	PublishConversationUpdated(ctx context.Context, conversation *entities.Conversation) error
	PublishConversationExpired(ctx context.Context, conversationID string) error

	// MessageLog events
	PublishMessageLogCreated(ctx context.Context, messageLog *entities.MessageLog) error
	PublishMessageStatusUpdated(ctx context.Context, messageID string, status entities.MessageStatus, timestamps map[string]time.Time) error
}

// ============================================
// EVENT PUBLISHER (business events via RabbitMQ)
// ============================================

// IEventPublisher define el contrato para publicar eventos en RabbitMQ
type IEventPublisher interface {
	PublishOrderConfirmed(ctx context.Context, orderNumber, phoneNumber string, businessID uint) error
	PublishOrderCancelled(ctx context.Context, orderNumber, reason, phoneNumber string, businessID uint) error
	PublishNoveltyRequested(ctx context.Context, orderNumber, noveltyType, phoneNumber string, businessID uint) error
	PublishHandoffRequested(ctx context.Context, orderNumber, phoneNumber string, businessID uint, conversationID string) error
}

// WhatsAppConfig contiene la configuración de WhatsApp
type WhatsAppConfig struct {
	PhoneNumberID uint
	AccessToken   string
	IntegrationID uint
	WhatsAppURL   string
}
