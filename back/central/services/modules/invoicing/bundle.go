package invoicing

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/app"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/primary/queue/consumer"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/secondary/queue"
	invoicingRedis "github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/secondary/redis"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
	"github.com/secamc93/probability/back/central/shared/redis"
)

// New inicializa el módulo de facturación
func New(
	router *gin.RouterGroup,
	database db.IDatabase,
	logger log.ILogger,
	config env.IConfig,
	rabbitMQ rabbitmq.IQueue,
	redisClient redis.IRedis,
	integrationCore core.IIntegrationCore,
) {
	ctx := context.Background()
	moduleLogger := logger.WithModule("invoicing")

	// ═══════════════════════════════════════════════════════════════
	// 1. INFRAESTRUCTURA SECUNDARIA (Adaptadores de salida)
	// ═══════════════════════════════════════════════════════════════

	// Servicio de caché para configuraciones de facturación
	configCache := invoicingRedis.NewConfigCache(redisClient, config, moduleLogger)

	// Repositorio único (GORM) - implementa TODAS las interfaces
	// Usa el servicio de caché para configuraciones de facturación
	repo := repository.New(database, configCache, moduleLogger)

	// Event publisher (RabbitMQ)
	eventPublisher := queue.New(rabbitMQ, moduleLogger)

	// SSE publisher (Redis Pub/Sub) para notificaciones en tiempo real
	sseChannel := config.Get("REDIS_INVOICE_EVENTS_CHANNEL")
	if sseChannel == "" {
		sseChannel = "probability:invoicing:events"
	}
	var ssePublisher = invoicingRedis.NewNoopSSEPublisher()
	if redisClient != nil {
		ssePublisher = invoicingRedis.NewSSEPublisher(redisClient, moduleLogger, sseChannel)

		// ═══════════════════════════════════════════════════════════════
		// REGISTRAR PREFIJOS DE CACHÉ Y CANALES PARA STARTUP LOGS
		// ═══════════════════════════════════════════════════════════════
		redisClient.RegisterCachePrefix("probability:invoicing:config:*")
		redisClient.RegisterChannel(sseChannel)
	} else {
		moduleLogger.Warn(ctx).Msg("Redis no disponible - SSE deshabilitado")
	}

	// Encryption service (para credenciales)
	encryptionKey := config.Get("ENCRYPTION_KEY")
	if encryptionKey == "" {
		moduleLogger.Warn(ctx).Msg("ENCRYPTION_KEY not set - using default (INSECURE)")
		encryptionKey = "default-encryption-key-change-me-in-production"
	}
	// TODO: Crear encryption service cuando esté disponible
	// encryptionService := encryption.New(encryptionKey, moduleLogger)

	// ═══════════════════════════════════════════════════════════════
	// 2. CAPA DE APLICACIÓN (Casos de uso)
	// ═══════════════════════════════════════════════════════════════

	useCase := app.New(
		repo,            // IRepository único (implementa TODAS las interfaces)
		integrationCore, // Integration Core (reemplaza provider repos y client)
		nil,             // Encryption - TODO: agregar cuando esté disponible
		eventPublisher,  // Event publisher (RabbitMQ)
		ssePublisher,    // SSE publisher (Redis Pub/Sub)
		moduleLogger,
	)

	// ═══════════════════════════════════════════════════════════════
	// 3. INFRAESTRUCTURA PRIMARIA (Adaptadores de entrada)
	// ═══════════════════════════════════════════════════════════════

	// HTTP Handlers
	handler := handlers.New(useCase, repo, moduleLogger)
	handler.RegisterRoutes(router)

	// Consumers (RabbitMQ)
	if rabbitMQ != nil {
		consumers := consumer.NewConsumers(
			rabbitMQ,
			useCase,
			repo, // IRepository único
			ssePublisher,
			moduleLogger,
		)

		// Iniciar Order Consumer (escucha events de órdenes)
		go func() {
			if err := consumers.Order.Start(ctx); err != nil {
				moduleLogger.Error(ctx).Err(err).Msg("Error al iniciar consumer de órdenes")
			}
		}()

		// Iniciar Retry Consumer (cron de reintentos cada 5 minutos)
		go consumers.Retry.Start(ctx)

		// Iniciar Bulk Invoice Consumer (procesa facturas masivas)
		go func() {
			if err := consumers.BulkInvoice.Start(ctx); err != nil {
				moduleLogger.Error(ctx).Err(err).Msg("Error al iniciar consumer de facturación masiva")
			}
		}()

		moduleLogger.Info(ctx).Msg("All consumers started successfully")
	} else {
		moduleLogger.Warn(ctx).Msg("RabbitMQ no disponible - consumers de facturación deshabilitados")
	}
}
