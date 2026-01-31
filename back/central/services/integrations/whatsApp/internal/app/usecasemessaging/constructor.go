package usecasemessaging

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/whatsApp/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/whatsApp/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/integrations/whatsApp/internal/domain/ports"
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

// useCase contiene todas las dependencias compartidas
type Usecases struct {
	whatsApp         ports.IWhatsApp
	conversationRepo ports.IConversationRepository
	messageLogRepo   ports.IMessageLogRepository
	integrationRepo  ports.IIntegrationRepository
	publisher        ports.IEventPublisher
	log              log.ILogger
	config           env.IConfig
}

// New crea la instancia única de use case con todas las dependencias
func New(
	whatsApp ports.IWhatsApp,
	conversationRepo ports.IConversationRepository,
	messageLogRepo ports.IMessageLogRepository,
	integrationRepo ports.IIntegrationRepository,
	publisher ports.IEventPublisher,
	logger log.ILogger,
	config env.IConfig,
) IUseCase {
	return &Usecases{
		whatsApp:         whatsApp,
		conversationRepo: conversationRepo,
		messageLogRepo:   messageLogRepo,
		integrationRepo:  integrationRepo,
		publisher:        publisher,
		log:              logger,
		config:           config,
	}
}
