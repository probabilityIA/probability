package ports

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/entities"
)

type IWhatsApp interface {
	SendMessage(ctx context.Context, phoneNumberID uint, msg entities.TemplateMessage, accessToken string) (string, error)
	SendTextMessage(ctx context.Context, phoneNumberID uint, toPhone, text, accessToken string) (string, error)
}

// ============================================
// CACHE INTERFACES (reemplazan repositorios DB)
// ============================================

// HumanSession representa una sesión de atención humana activa.
// Se crea cuando un agente humano responde manualmente a un cliente.
// Permite que las respuestas del cliente lleguen al dashboard en lugar del bot AI.
type HumanSession struct {
	ConversationID string
	BusinessID     uint
	PhoneNumber    string
}

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

	// ActivateHumanSession crea o renueva una sesión de atención humana para un teléfono.
	// TTL: 24h (alineado con la ventana de servicio de WhatsApp).
	// Cada llamada renueva el TTL automáticamente.
	ActivateHumanSession(ctx context.Context, phoneNumber, conversationID string, businessID uint) error

	// GetHumanSession retorna la sesión humana activa para un teléfono.
	// Retorna nil, nil si no existe (no es un error).
	GetHumanSession(ctx context.Context, phoneNumber string) (*HumanSession, error)

	// SetAIPaused pausa el bot AI para un teléfono. TTL: 24h.
	SetAIPaused(ctx context.Context, phoneNumber, conversationID string, businessID uint) error

	// IsAIPaused retorna true si el AI está pausado para este teléfono.
	IsAIPaused(ctx context.Context, phoneNumber string) bool

	// ClearAIPaused reactiva el AI para un teléfono.
	ClearAIPaused(ctx context.Context, phoneNumber string) error
}

// ICredentialsCache lee credenciales de WhatsApp desde Redis.
// Las claves son calentadas por integrations/core al startup.
type ICredentialsCache interface {
	// GetWhatsAppConfig obtiene credenciales de WhatsApp para un business.
	// Lee de: integration:idx:biz:{businessID}:type:2 -> integration:creds:{id} + integration:meta:{id}
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

// ============================================
// AI FORWARDER (forward messages to AI Sales agent)
// ============================================

// IAIForwarder reenvia mensajes sin conversacion activa al agente de ventas AI.
// La implementacion lee platform_creds de Redis para verificar ai_sales_enabled
// y obtener el business_id demo antes de publicar a whatsapp.ai.incoming.
type IAIForwarder interface {
	ForwardToAI(ctx context.Context, phoneNumber, messageText, messageID, messageType string) error
}

// IPlatformCredentialsGetter obtiene credenciales de plataforma cacheadas por tipo de integración.
// Implementado por integrations/core — evita que WhatsApp dependa de Redis directamente.
type IPlatformCredentialsGetter interface {
	GetCachedPlatformCredentials(ctx context.Context, integrationTypeID uint) (map[string]any, error)
}

// ISSEEventPublisher publica eventos SSE al exchange de eventos (para notificar al frontend en tiempo real).
// Implementado por infra/secondary/queue/sse_publisher.go usando rabbitmq.PublishEvent.
type ISSEEventPublisher interface {
	PublishMessageReceived(ctx context.Context, businessID uint, conversationID, phoneNumber, messageID, content string) error
	PublishConversationStarted(ctx context.Context, businessID uint, conversationID, phoneNumber string) error
	PublishMessageStatusUpdated(ctx context.Context, businessID uint, messageID, status string) error
}
