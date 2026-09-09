package whatsApp

import (
	"context"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/app/usecaseconnection"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/app/usecaseembedded"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/app/usecasemessaging"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/app/usecasenumbers"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/app/usecasetemplates"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/app/usecasetestconnection"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/primary/queue/consumerai"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/primary/queue/consumeralert"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/primary/queue/consumerauthotp"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/primary/queue/consumercampaign"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/primary/queue/consumerorder"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/secondary/mapimage"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/primary/queue/consumerscheduled"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/primary/queue/consumershipment"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/primary/queue/consumersubscriptionalert"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/primary/queue/consumertemplates"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/primary/queue/consumerwalletalert"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/primary/queue/consumerwebhook"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/secondary/cache"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/secondary/client"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/secondary/queue"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
	redisclient "github.com/secamc93/probability/back/central/shared/redis"
	"github.com/secamc93/probability/back/central/shared/storage"
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
	wa               ports.IWhatsApp
	useCase          usecasemessaging.IUseCase
	testUsecase      usecasetestconnection.ITestConnectionUseCase
	templatesUseCase usecasetemplates.IUseCase
	connectionUC     usecaseconnection.IUseCaseMutable
	numbersUC        usecasenumbers.IUseCaseMutable
	embeddedUC       usecaseembedded.IUseCaseMutable
	handler          handlers.IHandler
	credsCache       cache.ICredentialsCacheMutable
}

func New(config env.IConfig, logger log.ILogger, rabbit rabbitmq.IQueue, redisClient redisclient.IRedis, s3 storage.IS3Service) core.IIntegrationContract {
	logger = logger.WithModule("whatsapp")

	convCache, credsCache, templatesCache := cache.New(redisClient, logger)

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

	templatesAPIFactory := func(baseURL string) ports.ITemplateAPI {
		if baseURL == "" {
			baseURL = whatsappURL
		}
		return client.NewTemplatesClient(baseURL, logger)
	}

	var templateStatusPublisher usecasetemplates.IStatusPublisher
	if rabbit != nil {
		templateStatusPublisher = queue.NewTemplateStatusPublisher(rabbit, logger)
	}
	templatesUseCase := usecasetemplates.New(credsCache, templatesCache, templatesAPIFactory, templateStatusPublisher, logger)
	connectionUseCase := usecaseconnection.New(credsCache, templatesAPIFactory, logger)

	phoneNumbersAPIFactory := func(baseURL string) ports.IPhoneNumbersAPI {
		if baseURL == "" {
			baseURL = whatsappURL
		}
		return client.NewPhoneNumbersClient(baseURL, logger)
	}

	numbersUseCase := usecasenumbers.New(credsCache, phoneNumbersAPIFactory, logger)

	embeddedSignupAPIFactory := func(baseURL string) ports.IEmbeddedSignupAPI {
		if baseURL == "" {
			baseURL = whatsappURL
		}
		return client.NewEmbeddedSignupClient(baseURL, logger)
	}

	embeddedUseCase := usecaseembedded.New(credsCache, embeddedSignupAPIFactory, phoneNumbersAPIFactory, logger)

	handler := handlers.New(useCase, templatesUseCase, connectionUseCase, numbersUseCase, embeddedUseCase, logger, config, rabbit)

	if rabbit != nil {
		webhookConsumer := consumerwebhook.New(rabbit, useCase, templatesUseCase, logger)
		webhookConsumer.Start(context.Background())

		mapGenerator := mapimage.New(config, s3, logger)
		orderConsumer := consumerorder.New(
			rabbit,
			useCase,
			mapGenerator,
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

		templateSubmitConsumer := consumertemplates.New(rabbit, templatesUseCase, logger)
		go func() {
			if err := templateSubmitConsumer.Start(context.Background()); err != nil {
				logger.Error().Err(err).Msg("Error starting whatsapp template submit consumer")
			}
		}()

		scheduledSendConsumer := consumerscheduled.New(rabbit, useCase, logger)
		go func() {
			if err := scheduledSendConsumer.Start(context.Background()); err != nil {
				logger.Error().Err(err).Msg("Error starting scheduled send consumer")
			}
		}()

		campaignSendConsumer := consumercampaign.New(rabbit, useCase, logger)
		go func() {
			if err := campaignSendConsumer.Start(context.Background()); err != nil {
				logger.Error().Err(err).Msg("Error starting campaign send consumer")
			}
		}()
	}

	return &bundle{
		wa:               wa,
		useCase:          useCase,
		testUsecase:      testUsecase,
		templatesUseCase: templatesUseCase,
		connectionUC:     connectionUseCase,
		numbersUC:        numbersUseCase,
		embeddedUC:       embeddedUseCase,
		handler:          handler,
		credsCache:       credsCache,
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
	if b.connectionUC != nil {
		b.connectionUC.SetResolver(getter)
	}
	if b.numbersUC != nil {
		b.numbersUC.SetResolver(getter)
	}
	if b.embeddedUC != nil {
		b.embeddedUC.SetResolver(getter)
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

	config, credentials = b.fillWithPlatformCredentials(ctx, config, credentials)

	return b.testUsecase.TestConnection(ctx, config, credentials, clientFactory)
}

func (b *bundle) fillWithPlatformCredentials(
	ctx context.Context,
	config map[string]interface{},
	credentials map[string]interface{},
) (map[string]interface{}, map[string]interface{}) {
	if config == nil {
		config = map[string]interface{}{}
	}
	if credentials == nil {
		credentials = map[string]interface{}{}
	}

	token, _ := credentials["access_token"].(string)
	phoneID, _ := config["phone_number_id"].(string)
	if token != "" && phoneID != "" {
		return config, credentials
	}

	platform, err := b.credsCache.GetWhatsAppDefaultConfig(ctx)
	if err != nil {
		return config, credentials
	}

	if token == "" {
		credentials["access_token"] = platform.AccessToken
	}
	if phoneID == "" && platform.PhoneNumberID != 0 {
		config["phone_number_id"] = strconv.FormatUint(uint64(platform.PhoneNumberID), 10)
	}

	return config, credentials
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
