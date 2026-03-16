package usecasemessaging

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

// IUseCase define todos los métodos de los casos de uso de WhatsApp
type IUseCase interface {
	// SendMessage (legacy)
	SendMessage(ctx context.Context, req dtos.SendMessageRequest) (string, error)

	// SendTemplate
	SendTemplate(ctx context.Context, templateName, phoneNumber string, variables map[string]string, orderNumber string, businessID uint) (string, error)
	SendTemplateWithConversation(ctx context.Context, templateName, phoneNumber string, variables map[string]string, conversationID string) (string, error)

	// HandleWebhook
	HandleIncomingMessage(ctx context.Context, whPayload dtos.WebhookPayloadDTO) error
	HandleMessageStatus(ctx context.Context, whPayload dtos.WebhookPayloadDTO) error

	// ConversationManager
	TransitionState(ctx context.Context, conversation *entities.Conversation, userResponse string) (*dtos.StateTransitionDTO, error)
	GetInitialState() entities.ConversationState
	IsTerminalState(state entities.ConversationState) bool
}

// WhatsAppClientFactory crea un client HTTP de WhatsApp con la URL dada
type WhatsAppClientFactory func(baseURL string) ports.IWhatsApp

// usecases contiene todas las dependencias compartidas
type usecases struct {
	whatsApp         ports.IWhatsApp
	clientFactory    WhatsAppClientFactory
	conversationCache ports.IConversationCache
	credentialsCache  ports.ICredentialsCache
	persistPublisher  ports.IPersistencePublisher
	publisher        ports.IEventPublisher
	log              log.ILogger
	config           env.IConfig
}

// New crea la instancia única de use case con todas las dependencias
func New(
	whatsApp ports.IWhatsApp,
	conversationCache ports.IConversationCache,
	credentialsCache ports.ICredentialsCache,
	persistPublisher ports.IPersistencePublisher,
	publisher ports.IEventPublisher,
	logger log.ILogger,
	config env.IConfig,
	clientFactory ...WhatsAppClientFactory,
) IUseCase {
	uc := &usecases{
		whatsApp:          whatsApp,
		conversationCache: conversationCache,
		credentialsCache:  credentialsCache,
		persistPublisher:  persistPublisher,
		publisher:         publisher,
		log:               logger,
		config:            config,
	}
	if len(clientFactory) > 0 {
		uc.clientFactory = clientFactory[0]
	}
	return uc
}
