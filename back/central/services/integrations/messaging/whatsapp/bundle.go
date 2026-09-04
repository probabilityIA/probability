package whatsApp

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/app/usecasemessaging"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/app/usecasetestconnection"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/primary/queue/consumerai"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/primary/queue/consumeralert"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/primary/queue/consumerauthotp"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/primary/queue/consumerorder"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/primary/queue/consumershipment"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/primary/queue/consumersubscriptionalert"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/primary/queue/consumerwalletalert"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/primary/queue/consumerwebhook"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/secondary/cache"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/secondary/client"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/secondary/queue"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
	redisclient "github.com/secamc93/probability/back/central/shared/redis"
)

type IWhatsAppBundle interface {
	SendMessage(ctx context.Context, orderNumber, phoneNumber string) (string, error)
	SendDiagnosticResultTemplate(ctx context.Context, phone, name, surveyTitle, level string) (string, error)
	RegisterRoutes(router *gin.RouterGroup)
	SetPlatformCredsGetter(getter ports.IPlatformCredentialsGetter)
}

const diagnosticResultTemplateName = "resultado_diagnostico"

func buildDiagnosticoResultTemplateMessage(phoneNumber, name, surveyTitle, level string) entities.TemplateMessage {
	return entities.TemplateMessage{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               phoneNumber,
		Type:             "template",
		Template: entities.TemplateData{
			Name:     diagnosticResultTemplateName,
			Language: entities.TemplateLanguage{Code: "es"},
			Components: []entities.TemplateComponent{
				{
					Type: "body",
					Parameters: []entities.TemplateParameter{
						{Type: "text", Text: name},
						{Type: "text", Text: surveyTitle},
						{Type: "text", Text: level},
					},
				},
			},
		},
	}
}

type bundle struct {
	core.BaseIntegration
	wa          ports.IWhatsApp
	useCase     usecasemessaging.IUseCase
	testUsecase usecasetestconnection.ITestConnectionUseCase
	handler     handlers.IHandler
	credsCache  cache.ICredentialsCacheMutable
}

func New(config env.IConfig, logger log.ILogger, rabbit rabbitmq.IQueue, redisClient redisclient.IRedis) core.IIntegrationContract {
	logger = logger.WithModule("whatsapp")

	convCache, credsCache := cache.New(redisClient, logger)

	persistPub := queue.NewPersistencePublisher(rabbit, logger)

	whatsappURL := config.Get("WHATSAPP_URL")
	logger.Info().Str("whatsapp_url", whatsappURL).Msg("WhatsApp URL loaded from .env")

	wa := client.New(whatsappURL, logger)

	publisher := queue.NewWebhookPublisher(rabbit, logger)

	var aiForwarder ports.IAIForwarder
	if rabbit != nil {
		aiForwarder = queue.NewAIForwarder(rabbit, redisClient, logger)
	}

	var ssePublisher ports.ISSEEventPublisher
	if rabbit != nil {
		ssePublisher = queue.NewSSEPublisher(rabbit, logger)
	} else {
		ssePublisher = queue.NewNoopSSEPublisher()
	}

	clientFactory := func(baseURL string) ports.IWhatsApp {
		return client.New(baseURL, logger)
	}

	useCase := usecasemessaging.New(
		wa,
		convCache,
		credsCache,
		persistPub,
		publisher,
		logger,
		config,
		aiForwarder,
		ssePublisher,
		clientFactory,
	)

	testUsecase := usecasetestconnection.New(config, logger)

	handler := handlers.New(useCase, logger, config, rabbit)

	if rabbit != nil {
		webhookConsumer := consumerwebhook.New(rabbit, useCase, logger)
		webhookConsumer.Start(context.Background())

		orderConsumer := consumerorder.New(
			rabbit,
			useCase,
			logger,
		)

		go func() {
			if err := orderConsumer.Start(context.Background()); err != nil {
				logger.Error().Err(err).Msg("Error starting order confirmation consumer")
			}
		}()

		alertConsumer := consumeralert.New(rabbit, wa, credsCache, logger, config)
		go func() {
			if err := alertConsumer.Start(context.Background()); err != nil {
				logger.Error().Err(err).Msg("Error starting monitoring alert consumer")
			}
		}()

		walletAlertConsumer := consumerwalletalert.New(rabbit, useCase, logger)
		go func() {
			if err := walletAlertConsumer.Start(context.Background()); err != nil {
				logger.Error().Err(err).Msg("Error starting wallet balance alert consumer")
			}
		}()

		subscriptionAlertConsumer := consumersubscriptionalert.New(rabbit, useCase, logger)
		go func() {
			if err := subscriptionAlertConsumer.Start(context.Background()); err != nil {
				logger.Error().Err(err).Msg("Error starting subscription payment window consumer")
			}
		}()

		authOTPConsumer := consumerauthotp.New(rabbit, wa, credsCache, logger)
		go func() {
			if err := authOTPConsumer.Start(context.Background()); err != nil {
				logger.Error().Err(err).Msg("Error starting auth OTP consumer")
			}
		}()

		shipmentConsumer := consumershipment.New(rabbit, useCase, logger)
		go func() {
			if err := shipmentConsumer.Start(context.Background()); err != nil {
				logger.Error().Err(err).Msg("Error starting shipment guide notification consumer")
			}
		}()

		aiResponseConsumer := consumerai.New(rabbit, wa, credsCache, logger)
		go func() {
			if err := aiResponseConsumer.Start(context.Background()); err != nil {
				logger.Error().Err(err).Msg("Error starting AI response consumer")
			}
		}()
	}

	return &bundle{
		wa:          wa,
		useCase:     useCase,
		testUsecase: testUsecase,
		handler:     handler,
		credsCache:  credsCache,
	}
}

func (b *bundle) RegisterRoutes(router *gin.RouterGroup) {
	b.handler.RegisterRoutes(router)
}

func (b *bundle) SetPlatformCredsGetter(getter ports.IPlatformCredentialsGetter) {
	b.handler.SetPlatformCredsGetter(getter)
	if b.credsCache != nil {
		b.credsCache.SetResolver(getter)
	}
}

func (b *bundle) SendMessage(ctx context.Context, orderNumber, phoneNumber string) (string, error) {
	req := dtos.SendMessageRequest{
		OrderNumber: orderNumber,
		PhoneNumber: phoneNumber,
	}
	return b.useCase.SendMessage(ctx, req)
}

func (b *bundle) SendDiagnosticResultTemplate(ctx context.Context, phone, name, surveyTitle, level string) (string, error) {
	config, err := b.credsCache.GetWhatsAppDefaultConfig(ctx)
	if err != nil {
		return "", fmt.Errorf("credenciales de plataforma WhatsApp no disponibles: %w", err)
	}

	msg := buildDiagnosticoResultTemplateMessage(phone, name, surveyTitle, level)
	return b.wa.SendMessage(ctx, config.PhoneNumberID, msg, config.AccessToken)
}

func (b *bundle) TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error {
	clientFactory := func(baseURL string, logger log.ILogger) ports.IWhatsApp {
		return client.New(baseURL, logger)
	}

	return b.testUsecase.TestConnection(ctx, config, credentials, clientFactory)
}

func (b *bundle) GetWebhookURL(ctx context.Context, baseURL string, integrationID uint) (*core.WebhookInfo, error) {
	webhookURL := fmt.Sprintf("%s/integrations/whatsapp/webhook", baseURL)

	return &core.WebhookInfo{
		URL:         webhookURL,
		Method:      "POST",
		Description: "Configura este webhook en Meta Business Manager (WhatsApp > Configuration > Webhook) para recibir eventos de mensajes, estados de entrega (sent, delivered, read, failed) y respuestas de botones en tiempo real. Asegúrate de suscribirte a los campos 'messages' y 'message_template_status_update'.",
		Events: []string{
			"messages",
			"message_template_status_update",
		},
	}, nil
}
