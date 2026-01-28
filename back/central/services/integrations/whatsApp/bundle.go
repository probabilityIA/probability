package whatsApp

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/whatsApp/internal/app"
	"github.com/secamc93/probability/back/central/services/integrations/whatsApp/internal/app/usecasetestconnection"
	"github.com/secamc93/probability/back/central/services/integrations/whatsApp/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/whatsApp/internal/infra/primary/handlers"
	primaryqueue "github.com/secamc93/probability/back/central/services/integrations/whatsApp/internal/infra/primary/queue"
	"github.com/secamc93/probability/back/central/services/integrations/whatsApp/internal/infra/secondary/client"
	"github.com/secamc93/probability/back/central/services/integrations/whatsApp/internal/infra/secondary/queue"
	"github.com/secamc93/probability/back/central/services/integrations/whatsApp/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// IWhatsAppBundle define la interfaz del bundle de WhatsApp
type IWhatsAppBundle interface {
	// SendMessage envía un mensaje de WhatsApp con el número de orden y número de teléfono
	SendMessage(ctx context.Context, orderNumber, phoneNumber string) (string, error)
	// RegisterRoutes registra las rutas HTTP del módulo
	RegisterRoutes(router *gin.RouterGroup)
}

type bundle struct {
	wa                  domain.IWhatsApp
	usecase             app.IUseCaseSendMessage
	sendTemplateUseCase app.ISendTemplateMessageUseCase
	handleWebhookUseCase app.IHandleWebhookUseCase
	testUsecase         usecasetestconnection.ITestConnectionUseCase
	templateHandler     *handlers.TemplateHandler
	webhookHandler      *handlers.WebhookHandler
}

// New crea una nueva instancia del bundle de WhatsApp con todas sus dependencias
func New(config env.IConfig, logger log.ILogger, database db.IDatabase, rabbit rabbitmq.IQueue) core.IIntegrationContract {
	logger = logger.WithModule("whatsapp")

	// 1. Capa de infraestructura secundaria (adaptadores de salida)
	// Cliente HTTP de WhatsApp
	wa := client.New(config)

	// Repositorios
	conversationRepo := repository.NewConversationRepository(database, logger)
	messageLogRepo := repository.NewMessageLogRepository(database, logger)

	// Publisher de eventos RabbitMQ
	publisher := queue.NewWebhookPublisher(rabbit, logger)

	// 2. Capa de aplicación (casos de uso)
	// Usecase legacy (mantener compatibilidad)
	usecase := app.New(wa, logger, config)

	// Nuevos usecases
	sendTemplateUseCase := app.NewSendTemplateMessage(
		wa,
		conversationRepo,
		messageLogRepo,
		logger,
		config,
	)

	conversationManager := app.NewConversationManager(
		conversationRepo,
		logger,
	)

	handleWebhookUseCase := app.NewHandleWebhook(
		conversationRepo,
		messageLogRepo,
		sendTemplateUseCase,
		publisher,
		conversationManager,
		logger,
	)

	testUsecase := usecasetestconnection.New(config, logger)

	// 3. Capa de infraestructura primaria (adaptadores de entrada)
	templateHandler := handlers.NewTemplateHandler(sendTemplateUseCase, logger)
	webhookHandler := handlers.NewWebhookHandler(handleWebhookUseCase, logger, config)

	// 4. Inicializar consumidor de órdenes (si RabbitMQ está disponible)
	if rabbit != nil {
		orderConsumer := primaryqueue.NewOrderConfirmationConsumer(
			rabbit,
			sendTemplateUseCase,
			logger,
		)

		// Arrancar consumidor en goroutine
		go func() {
			if err := orderConsumer.Start(context.Background()); err != nil {
				logger.Error().Err(err).Msg("Error starting order confirmation consumer")
			}
		}()

		logger.Info().Msg("Order confirmation consumer initialized")
	}

	return &bundle{
		wa:                   wa,
		usecase:              usecase,
		sendTemplateUseCase:  sendTemplateUseCase,
		handleWebhookUseCase: handleWebhookUseCase,
		testUsecase:          testUsecase,
		templateHandler:      templateHandler,
		webhookHandler:       webhookHandler,
	}
}

// RegisterRoutes registra todas las rutas HTTP del módulo
func (b *bundle) RegisterRoutes(router *gin.RouterGroup) {
	whatsapp := router.Group("/whatsapp")
	{
		// Endpoint para envío de plantillas
		whatsapp.POST("/send-template", b.templateHandler.SendTemplate)

		// Endpoints de webhook
		whatsapp.GET("/webhook", b.webhookHandler.VerifyWebhook)
		whatsapp.POST("/webhook", b.webhookHandler.ReceiveWebhook)
	}
}

// SendMessage expone el método simplificado para enviar mensajes (legacy)
func (b *bundle) SendMessage(ctx context.Context, orderNumber, phoneNumber string) (string, error) {
	req := domain.SendMessageRequest{
		OrderNumber: orderNumber,
		PhoneNumber: phoneNumber,
	}
	return b.usecase.SendMessage(ctx, req)
}

// TestConnection prueba la conexión enviando un mensaje de prueba
func (b *bundle) TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error {
	// Factory para crear clientes de WhatsApp con configuración dinámica
	clientFactory := func(cfg env.IConfig) domain.IWhatsApp {
		return client.New(cfg)
	}

	// Delegar al caso de uso pasando los mapas directamente
	return b.testUsecase.TestConnection(ctx, config, credentials, clientFactory)
}

// SyncOrdersByIntegrationID no está soportado para WhatsApp
func (b *bundle) SyncOrdersByIntegrationID(ctx context.Context, integrationID string) error {
	return fmt.Errorf("order synchronization is not supported for WhatsApp integration")
}

// SyncOrdersByIntegrationIDWithParams no está soportado para WhatsApp
func (b *bundle) SyncOrdersByIntegrationIDWithParams(ctx context.Context, integrationID string, params interface{}) error {
	return fmt.Errorf("order synchronization is not supported for WhatsApp integration")
}

// SyncOrdersByBusiness no está soportado para WhatsApp
func (b *bundle) SyncOrdersByBusiness(ctx context.Context, businessID uint) error {
	return fmt.Errorf("order synchronization is not supported for WhatsApp integration")
}

// GetWebhookURL retorna la URL del webhook de WhatsApp
func (b *bundle) GetWebhookURL(ctx context.Context, baseURL string, integrationID uint) (*core.WebhookInfo, error) {
	webhookURL := fmt.Sprintf("%s/api/integrations/whatsapp/webhook", baseURL)

	return &core.WebhookInfo{
		URL:    webhookURL,
		Events: []string{"messages", "message_template_status_update"},
		Method: "POST",
	}, nil
}
