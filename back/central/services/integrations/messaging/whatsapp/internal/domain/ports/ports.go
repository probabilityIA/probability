package ports

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/entities"
)

type IWhatsApp interface {
	SendMessage(ctx context.Context, phoneNumberID uint, msg entities.TemplateMessage, accessToken string) (string, error)
}

// IConversationRepository define el contrato para el repositorio de conversaciones
type IConversationRepository interface {
	Create(ctx context.Context, conversation *entities.Conversation) error
	GetByID(ctx context.Context, id string) (*entities.Conversation, error)
	GetByPhoneAndOrder(ctx context.Context, phoneNumber, orderNumber string) (*entities.Conversation, error)
	GetActiveByPhone(ctx context.Context, phoneNumber string) (*entities.Conversation, error)
	Update(ctx context.Context, conversation *entities.Conversation) error
	Expire(ctx context.Context, id string) error
	Delete(ctx context.Context, id string) error
}

// IMessageLogRepository define el contrato para el repositorio de logs de mensajes
type IMessageLogRepository interface {
	Create(ctx context.Context, messageLog *entities.MessageLog) error
	GetByID(ctx context.Context, id string) (*entities.MessageLog, error)
	GetByMessageID(ctx context.Context, messageID string) (*entities.MessageLog, error)
	GetByConversation(ctx context.Context, conversationID string) ([]entities.MessageLog, error)
	UpdateStatus(ctx context.Context, messageID string, status entities.MessageStatus, timestamps map[string]time.Time) error
	Delete(ctx context.Context, id string) error
}

// IEventPublisher define el contrato para publicar eventos en RabbitMQ
type IEventPublisher interface {
	PublishOrderConfirmed(ctx context.Context, orderNumber, phoneNumber string, businessID uint) error
	PublishOrderCancelled(ctx context.Context, orderNumber, reason, phoneNumber string, businessID uint) error
	PublishNoveltyRequested(ctx context.Context, orderNumber, noveltyType, phoneNumber string, businessID uint) error
	PublishHandoffRequested(ctx context.Context, orderNumber, phoneNumber string, businessID uint, conversationID string) error
}

// IIntegrationRepository define el contrato para obtener configuraciones de integraciones
type IIntegrationRepository interface {
	// GetWhatsAppConfig obtiene la configuración de WhatsApp para un business
	// Retorna phone_number_id y access_token desencriptado
	GetWhatsAppConfig(ctx context.Context, businessID uint) (*WhatsAppConfig, error)

	// GetWhatsAppDefaultConfig obtiene las credenciales globales de WhatsApp
	// desde el tipo de integración (platform_credentials_encrypted).
	// Usado para alertas de plataforma que no pertenecen a ningún business.
	GetWhatsAppDefaultConfig(ctx context.Context) (*WhatsAppConfig, error)
}

// WhatsAppConfig contiene la configuración de WhatsApp obtenida desde la base de datos
type WhatsAppConfig struct {
	PhoneNumberID uint
	AccessToken   string
	IntegrationID uint
}
