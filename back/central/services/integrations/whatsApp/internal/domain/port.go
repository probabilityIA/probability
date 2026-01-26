package domain

import (
	"context"
	"time"
)

type IWhatsApp interface {
	SendMessage(ctx context.Context, phoneNumberID uint, msg TemplateMessage) (string, error)
}

// IConversationRepository define el contrato para el repositorio de conversaciones
type IConversationRepository interface {
	Create(ctx context.Context, conversation *Conversation) error
	GetByID(ctx context.Context, id string) (*Conversation, error)
	GetByPhoneAndOrder(ctx context.Context, phoneNumber, orderNumber string) (*Conversation, error)
	GetActiveByPhone(ctx context.Context, phoneNumber string) (*Conversation, error)
	Update(ctx context.Context, conversation *Conversation) error
	Expire(ctx context.Context, id string) error
	Delete(ctx context.Context, id string) error
}

// IMessageLogRepository define el contrato para el repositorio de logs de mensajes
type IMessageLogRepository interface {
	Create(ctx context.Context, messageLog *MessageLog) error
	GetByID(ctx context.Context, id string) (*MessageLog, error)
	GetByMessageID(ctx context.Context, messageID string) (*MessageLog, error)
	GetByConversation(ctx context.Context, conversationID string) ([]MessageLog, error)
	UpdateStatus(ctx context.Context, messageID string, status MessageStatus, timestamps map[string]time.Time) error
	Delete(ctx context.Context, id string) error
}

// IEventPublisher define el contrato para publicar eventos en RabbitMQ
type IEventPublisher interface {
	PublishOrderConfirmed(ctx context.Context, orderNumber, phoneNumber string, businessID uint) error
	PublishOrderCancelled(ctx context.Context, orderNumber, reason, phoneNumber string, businessID uint) error
	PublishNoveltyRequested(ctx context.Context, orderNumber, noveltyType, phoneNumber string, businessID uint) error
	PublishHandoffRequested(ctx context.Context, orderNumber, phoneNumber string, businessID uint, conversationID string) error
}
