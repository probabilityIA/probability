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

type HumanSession struct {
	ConversationID string
	BusinessID     uint
	PhoneNumber    string
}

type IConversationCache interface {
	GetByID(ctx context.Context, id string) (*entities.Conversation, error)
	GetByPhoneAndOrder(ctx context.Context, phoneNumber, orderNumber string) (*entities.Conversation, error)
	GetActiveByPhone(ctx context.Context, phoneNumber string) (*entities.Conversation, error)
	Save(ctx context.Context, conversation *entities.Conversation) error
	Expire(ctx context.Context, id string) error
	ActivateHumanSession(ctx context.Context, phoneNumber, conversationID string, businessID uint) error
	GetHumanSession(ctx context.Context, phoneNumber string) (*HumanSession, error)
	SetAIPaused(ctx context.Context, phoneNumber, conversationID string, businessID uint) error
	IsAIPaused(ctx context.Context, phoneNumber string) bool
	ClearAIPaused(ctx context.Context, phoneNumber string) error
}

type ICredentialsCache interface {
	GetWhatsAppConfig(ctx context.Context, businessID uint) (*WhatsAppConfig, error)
	GetWhatsAppConfigByIntegration(ctx context.Context, integrationID uint) (*WhatsAppConfig, error)
	GetWhatsAppDefaultConfig(ctx context.Context) (*WhatsAppConfig, error)
	ResolveByPhoneNumberID(ctx context.Context, phoneNumberID string) (*IntegrationOwner, error)
}

type IntegrationOwner struct {
	IntegrationID uint
	BusinessID    uint
}

type IPersistencePublisher interface {
	PublishConversationCreated(ctx context.Context, conversation *entities.Conversation) error
	PublishConversationUpdated(ctx context.Context, conversation *entities.Conversation) error
	PublishConversationExpired(ctx context.Context, conversationID string) error
	PublishMessageLogCreated(ctx context.Context, messageLog *entities.MessageLog) error
	PublishMessageStatusUpdated(ctx context.Context, messageID string, status entities.MessageStatus, timestamps map[string]time.Time) error
}

type IEventPublisher interface {
	PublishOrderConfirmed(ctx context.Context, orderNumber, phoneNumber string, businessID uint) error
	PublishOrderCancelled(ctx context.Context, orderNumber, reason, phoneNumber string, businessID uint) error
	PublishNoveltyRequested(ctx context.Context, orderNumber, noveltyType, phoneNumber string, businessID uint) error
	PublishHandoffRequested(ctx context.Context, orderNumber, phoneNumber string, businessID uint, conversationID string) error
}

type WhatsAppConfig struct {
	PhoneNumberID uint
	AccessToken   string
	IntegrationID uint
	WhatsAppURL   string
	WABAID        string
	OwnNumber     bool
}

type IAIForwarder interface {
	ForwardToAI(ctx context.Context, phoneNumber, messageText, messageID, messageType string) error
}

type IPlatformCredentialsGetter interface {
	GetCachedPlatformCredentials(ctx context.Context, integrationTypeID uint) (map[string]any, error)
	GetIntegrationIDByBusinessAndType(ctx context.Context, businessID, integrationTypeID uint) (uint, error)
	GetIntegrationConfigAndCredentials(ctx context.Context, integrationID uint) (map[string]any, map[string]any, error)
	FindIntegrationByConfigValue(ctx context.Context, integrationTypeID uint, field, value string) (uint, uint, error)
}

type ISSEEventPublisher interface {
	PublishMessageReceived(ctx context.Context, businessID uint, conversationID, phoneNumber, messageID, content string) error
	PublishConversationStarted(ctx context.Context, businessID uint, conversationID, phoneNumber string) error
	PublishMessageStatusUpdated(ctx context.Context, businessID uint, messageID, status string) error
}

type TemplateStatus struct {
	Name        string
	Language    string
	Status      string
	Category    string
	MetaID      string
	Reason      string
	UpdatedAt   time.Time
	Provisioned bool
}

type WABATemplatesSnapshot struct {
	IntegrationID uint
	BusinessID    uint
	WABAID        string
	Templates     []TemplateStatus
	RefreshedAt   time.Time
}

type ITemplateAPI interface {
	ListTemplates(ctx context.Context, wabaID, accessToken string) ([]TemplateDefinitionRemote, error)
	CreateTemplate(ctx context.Context, wabaID, accessToken string, template TemplateDefinitionRemote) (string, error)
}

type TemplateDefinitionRemote struct {
	ID              string
	Name            string
	Language        string
	Category        string
	Status          string
	RejectedReason  string
	ParameterFormat string
	Components      []map[string]any
}

type ITemplateStatusCache interface {
	Get(ctx context.Context, integrationID uint) (*WABATemplatesSnapshot, error)
	Save(ctx context.Context, snapshot *WABATemplatesSnapshot) error
	UpdateStatusByWABA(ctx context.Context, wabaID, name, language, status, reason string) error
	IndexWABA(ctx context.Context, wabaID string, integrationID uint) error
	IntegrationByWABA(ctx context.Context, wabaID string) (uint, error)
}
