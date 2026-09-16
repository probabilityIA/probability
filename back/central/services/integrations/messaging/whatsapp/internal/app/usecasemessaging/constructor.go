package usecasemessaging

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

type IUseCase interface {
	SystemFlows() []entities.SystemFlow
	SendManualMedia(ctx context.Context, conversationID string, phoneNumber string, businessID uint, media entities.OutboundMedia, caption string, sentBy string) (string, error)
	SendMessage(ctx context.Context, req dtos.SendMessageRequest) (string, error)

	SendTemplate(ctx context.Context, templateName, phoneNumber string, variables map[string]string, orderNumber string, businessID uint) (string, error)
	SendTemplateWithHeaderImage(ctx context.Context, templateName, phoneNumber string, variables map[string]string, headerImageURL string, orderNumber string, businessID uint) (string, error)
	SendCustomTemplate(ctx context.Context, templateName, language, phoneNumber string, parameters []string, headerImageURL string, businessID uint) (string, error)
	SendPlatformTemplate(ctx context.Context, templateName, phoneNumber string, variables map[string]string, orderNumber string, businessID uint) (string, error)
	SendTemplateWithConversation(ctx context.Context, templateName, phoneNumber string, variables map[string]string, conversationID string) (string, error)

	SendManualReply(ctx context.Context, conversationID, phoneNumber string, businessID uint, text, sentBy string) (string, error)

	PauseAI(ctx context.Context, conversationID, phoneNumber string, businessID uint) error

	ResumeAI(ctx context.Context, conversationID, phoneNumber string, businessID uint) error

	HandleIncomingMessage(ctx context.Context, whPayload dtos.WebhookPayloadDTO) error
	HandleMessageStatus(ctx context.Context, whPayload dtos.WebhookPayloadDTO) error

	TransitionState(ctx context.Context, conversation *entities.Conversation, userResponse string) (*dtos.StateTransitionDTO, error)
	GetInitialState() entities.ConversationState
	IsTerminalState(state entities.ConversationState) bool
}

type IUseCaseMutable interface {
	SetMediaDependencies(factory MediaAPIFactory, storage ports.IChatMediaStorage)
	IUseCase
	SetButtonReplyPublisher(publisher ports.IButtonReplyPublisher)
}

type WhatsAppClientFactory func(baseURL string) ports.IWhatsApp

type usecases struct {
	mediaFactory      MediaAPIFactory
	chatMedia         ports.IChatMediaStorage
	whatsApp          ports.IWhatsApp
	clientFactory     WhatsAppClientFactory
	conversationCache ports.IConversationCache
	credentialsCache  ports.ICredentialsCache
	persistPublisher  ports.IPersistencePublisher
	publisher         ports.IEventPublisher
	ssePublisher      ports.ISSEEventPublisher
	aiForwarder       ports.IAIForwarder
	buttonReplyPub    ports.IButtonReplyPublisher
	log               log.ILogger
	config            env.IConfig
}

func (u *usecases) SetButtonReplyPublisher(publisher ports.IButtonReplyPublisher) {
	u.buttonReplyPub = publisher
}

func New(
	whatsApp ports.IWhatsApp,
	conversationCache ports.IConversationCache,
	credentialsCache ports.ICredentialsCache,
	persistPublisher ports.IPersistencePublisher,
	publisher ports.IEventPublisher,
	logger log.ILogger,
	config env.IConfig,
	aiForwarder ports.IAIForwarder,
	ssePublisher ports.ISSEEventPublisher,
	clientFactory ...WhatsAppClientFactory,
) IUseCaseMutable {
	uc := &usecases{
		whatsApp:          whatsApp,
		conversationCache: conversationCache,
		credentialsCache:  credentialsCache,
		persistPublisher:  persistPublisher,
		publisher:         publisher,
		ssePublisher:      ssePublisher,
		aiForwarder:       aiForwarder,
		log:               logger,
		config:            config,
	}
	if len(clientFactory) > 0 {
		uc.clientFactory = clientFactory[0]
	}
	return uc
}
