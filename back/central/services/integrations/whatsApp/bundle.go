package whatsApp

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/whatsApp/internal/app/usecasemessaging"
	"github.com/secamc93/probability/back/central/services/integrations/whatsApp/internal/app/usecasetestconnection"
	"github.com/secamc93/probability/back/central/services/integrations/whatsApp/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/whatsApp/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/integrations/whatsApp/internal/infra/primary/consumer/consumerevent"
	"github.com/secamc93/probability/back/central/services/integrations/whatsApp/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/integrations/whatsApp/internal/infra/primary/queue/consumerorder"
	"github.com/secamc93/probability/back/central/services/integrations/whatsApp/internal/infra/secondary/client"
	"github.com/secamc93/probability/back/central/services/integrations/whatsApp/internal/infra/secondary/queue"
	"github.com/secamc93/probability/back/central/services/integrations/whatsApp/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/services/modules"
	"github.com/secamc93/probability/back/central/shared/db"
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
}

type bundle struct {
	wa          ports.IWhatsApp
	useCase     usecasemessaging.IUseCase
	testUsecase usecasetestconnection.ITestConnectionUseCase
	handler     handlers.IHandler
}

// New crea una nueva instancia del bundle de WhatsApp con todas sus dependencias
func New(config env.IConfig, logger log.ILogger, database db.IDatabase, rabbit rabbitmq.IQueue, redisClient redisclient.IRedis, moduleBundles *modules.ModuleBundles) core.IIntegrationContract {
	logger = logger.WithModule("whatsapp")

	// 1. Capa de infraestructura secundaria (adaptadores de salida)
	// Cliente HTTP de WhatsApp
	wa := client.New(config)

	// Preparar encryption key para IntegrationRepository
	encryptionKeyStr := config.Get("ENCRYPTION_KEY")
	var encryptionKey []byte
	decoded, err := base64.StdEncoding.DecodeString(encryptionKeyStr)
	if err == nil && len(decoded) == 32 {
		encryptionKey = decoded
	} else {
		encryptionKey = []byte(encryptionKeyStr)
	}

	// Repositorios (constructor consolidado)
	conversationRepo, messageLogRepo, integrationRepo := repository.New(database, logger, encryptionKey)

	// Publisher de eventos RabbitMQ
	publisher := queue.NewWebhookPublisher(rabbit, logger)

	// 2. Capa de aplicación (casos de uso)
	// Casos de uso (constructor consolidado)
	useCase := usecasemessaging.New(
		wa,
		conversationRepo,
		messageLogRepo,
		integrationRepo,
		publisher,
		logger,
		config,
	)

	// Test usecase (subdirectorio separado)
	testUsecase := usecasetestconnection.New(config, logger)

	// 3. Capa de infraestructura primaria (adaptadores de entrada)
	handler := handlers.New(useCase, logger, config)

	// 4. Inicializar consumidor de órdenes (si RabbitMQ está disponible)
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

		logger.Info().Msg("Order confirmation consumer initialized")
	}

	// 5. Inicializar consumidor de eventos Redis → RabbitMQ (si Redis y RabbitMQ están disponibles)
	if redisClient != nil && rabbit != nil {
		// Crear repositorios de consultas
		orderQueries := repository.NewOrderQueries(database, logger)
		integrationQueries := repository.NewIntegrationQueries(database, logger)

		// Obtener canal de eventos de Redis desde configuración
		redisChannel := config.Get("REDIS_ORDER_EVENTS_CHANNEL")
		if redisChannel == "" {
			redisChannel = "probability:orders:events" // Valor por defecto
		}

		// Crear repositorio de notification_config (consulta directa a BD)
		notificationConfigRepo := repository.NewNotificationConfigRepository(database, logger)

		// Crear consumer de eventos
		orderEventConsumer := consumerevent.New(
			redisClient,
			rabbit,
			notificationConfigRepo, // ← Pasa repositorio
			integrationQueries,
			orderQueries,
			logger,
			redisChannel,
		)

		// Arrancar consumer en goroutine
		go func() {
			if err := orderEventConsumer.Start(context.Background()); err != nil {
				logger.Error().Err(err).Msg("Error starting WhatsApp order event consumer")
			}
		}()

		logger.Info().Msg("WhatsApp order event consumer initialized")
	}

	return &bundle{
		wa:          wa,
		useCase:     useCase,
		testUsecase: testUsecase,
		handler:     handler,
	}
}

// RegisterRoutes registra todas las rutas HTTP del módulo
func (b *bundle) RegisterRoutes(router *gin.RouterGroup) {
	// Delegar al handler
	b.handler.RegisterRoutes(router)
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
	// Factory para crear clientes de WhatsApp con configuración dinámica
	clientFactory := func(cfg env.IConfig) ports.IWhatsApp {
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
	// Construir la URL del webhook
	// El webhook se recibe en: /integrations/whatsapp/webhook
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
