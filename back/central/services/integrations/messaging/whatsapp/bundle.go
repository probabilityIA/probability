package whatsApp

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/app/usecasemessaging"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/app/usecasetestconnection"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/primary/queue/consumerai"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/primary/queue/consumeralert"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/primary/queue/consumerorder"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/primary/queue/consumershipment"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/secondary/cache"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/secondary/client"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/secondary/queue"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
	redisclient "github.com/secamc93/probability/back/central/shared/redis"
)

// IWhatsAppBundle define la interfaz del bundle de WhatsApp
type IWhatsAppBundle interface {
	// SendMessage envía un mensaje de WhatsApp con el número de orden y número de teléfono
	SendMessage(ctx context.Context, orderNumber, phoneNumber string) (string, error)
	// RegisterRoutes registra las rutas HTTP del módulo
	RegisterRoutes(router *gin.RouterGroup)
	// SetPlatformCredsGetter inyecta el getter de credenciales de plataforma (de core) después de la construcción
	SetPlatformCredsGetter(getter ports.IPlatformCredentialsGetter)
}

type bundle struct {
	core.BaseIntegration
	wa          ports.IWhatsApp
	useCase     usecasemessaging.IUseCase
	testUsecase usecasetestconnection.ITestConnectionUseCase
	handler     handlers.IHandler
	credsCache  cache.ICredentialsCacheMutable
}

// New crea una nueva instancia del bundle de WhatsApp con todas sus dependencias.
// WhatsApp es stateless (sin DB directa): lee credenciales de Redis cache
// y persiste conversaciones/logs asincrónicamente via RabbitMQ.
func New(config env.IConfig, logger log.ILogger, rabbit rabbitmq.IQueue, redisClient redisclient.IRedis) core.IIntegrationContract {
	logger = logger.WithModule("whatsapp")

	// 1. Capa de infraestructura secundaria (adaptadores de salida)
	// Cache de conversaciones + credenciales (Redis)
	convCache, credsCache := cache.New(redisClient, logger)

	// Publisher de persistencia (async DB via RabbitMQ)
	persistPub := queue.NewPersistencePublisher(rabbit, logger)

	// WhatsApp URL desde .env
	whatsappURL := config.Get("WHATSAPP_URL")
	logger.Info().Str("whatsapp_url", whatsappURL).Msg("WhatsApp URL loaded from .env")

	// Cliente HTTP de WhatsApp
	wa := client.New(whatsappURL, logger)

	// Publisher de eventos de negocio (RabbitMQ)
	publisher := queue.NewWebhookPublisher(rabbit, logger)

	// AI Forwarder: publica mensajes sin conversación activa a whatsapp.ai.incoming
	var aiForwarder ports.IAIForwarder
	if rabbit != nil {
		aiForwarder = queue.NewAIForwarder(rabbit, redisClient, logger)
	}

	// SSE Publisher: notifica al frontend vía RabbitMQ -> events module -> SSE endpoint
	var ssePublisher ports.ISSEEventPublisher
	if rabbit != nil {
		ssePublisher = queue.NewSSEPublisher(rabbit, logger)
	} else {
		ssePublisher = queue.NewNoopSSEPublisher()
	}

	// 2. Capa de aplicación (casos de uso)
	// Factory para crear clients con URL dinámica (de platform_creds, no de .env)
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

	// Test usecase (subdirectorio separado)
	testUsecase := usecasetestconnection.New(config, logger)

	// 3. Capa de infraestructura primaria (adaptadores de entrada)
	handler := handlers.New(useCase, logger, config)

	// 4. Inicializar consumidores (si RabbitMQ está disponible)
	if rabbit != nil {
		orderConsumer := consumerorder.New(
			rabbit,
			useCase,
			logger,
		)

		// Arrancar consumidor en goroutine
		go func() {
			if err := orderConsumer.Start(context.Background()); err != nil {
				logger.Error().Err(err).Msg("Error starting order confirmation consumer")
			}
		}()

		// Inicializar consumidor de alertas de monitoreo
		alertConsumer := consumeralert.New(rabbit, wa, credsCache, logger)
		go func() {
			if err := alertConsumer.Start(context.Background()); err != nil {
				logger.Error().Err(err).Msg("Error starting monitoring alert consumer")
			}
		}()

		// Inicializar consumidor de notificaciones de guía de envío
		shipmentConsumer := consumershipment.New(rabbit, useCase, logger)
		go func() {
			if err := shipmentConsumer.Start(context.Background()); err != nil {
				logger.Error().Err(err).Msg("Error starting shipment guide notification consumer")
			}
		}()

		// Inicializar consumidor de respuestas AI (whatsapp.ai.response -> envío de texto libre)
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

// RegisterRoutes registra todas las rutas HTTP del módulo
func (b *bundle) RegisterRoutes(router *gin.RouterGroup) {
	b.handler.RegisterRoutes(router)
}

func (b *bundle) SetPlatformCredsGetter(getter ports.IPlatformCredentialsGetter) {
	b.handler.SetPlatformCredsGetter(getter)
	if b.credsCache != nil {
		b.credsCache.SetResolver(getter)
	}
}

// SendMessage expone el método simplificado para enviar mensajes (legacy)
func (b *bundle) SendMessage(ctx context.Context, orderNumber, phoneNumber string) (string, error) {
	req := dtos.SendMessageRequest{
		OrderNumber: orderNumber,
		PhoneNumber: phoneNumber,
	}
	return b.useCase.SendMessage(ctx, req)
}

// TestConnection prueba la conexión enviando un mensaje de prueba
func (b *bundle) TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error {
	// Factory para crear clientes de WhatsApp con URL dinámica
	clientFactory := func(baseURL string, logger log.ILogger) ports.IWhatsApp {
		return client.New(baseURL, logger)
	}

	// Delegar al caso de uso pasando los mapas directamente
	return b.testUsecase.TestConnection(ctx, config, credentials, clientFactory)
}

// GetWebhookURL retorna la URL del webhook de WhatsApp (implementa IWebhookProvider)
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
