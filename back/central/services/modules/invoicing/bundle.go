package invoicing

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/app"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/primary/queue/consumer"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/secondary/providers/softpymes"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/secondary/queue"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// New inicializa el módulo de facturación
func New(
	router *gin.RouterGroup,
	database db.IDatabase,
	logger log.ILogger,
	config env.IConfig,
	rabbitMQ rabbitmq.IQueue,
) {
	ctx := context.Background()
	moduleLogger := logger.WithModule("invoicing")

	moduleLogger.Info(ctx).Msg("Initializing invoicing module")

	// ═══════════════════════════════════════════════════════════════
	// 1. INFRAESTRUCTURA SECUNDARIA (Adaptadores de salida)
	// ═══════════════════════════════════════════════════════════════

	// Repositorios (GORM)
	repos := repository.New(database, moduleLogger)

	// Cliente de Softpymes
	softpymesBaseURL := config.Get("SOFTPYMES_API_URL")
	if softpymesBaseURL == "" {
		softpymesBaseURL = "https://api-integracion.softpymes.com.co/app/integration/"
	}
	softpymesClient := softpymes.New(softpymesBaseURL, moduleLogger)

	// Event publisher (RabbitMQ)
	eventPublisher := queue.NewEventPublisher(rabbitMQ, moduleLogger)

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
		repos.Invoice,
		repos.InvoiceItem,
		repos.Provider,
		repos.ProviderType,
		repos.Config,
		repos.SyncLog,
		repos.CreditNote,
		nil,             // Order repository - se debe inyectar desde orders module
		softpymesClient, // Provider client (Softpymes)
		nil,             // Encryption - TODO: agregar cuando esté disponible
		eventPublisher,  // Event publisher (RabbitMQ)
		moduleLogger,
	)

	// ═══════════════════════════════════════════════════════════════
	// 3. INFRAESTRUCTURA PRIMARIA (Adaptadores de entrada)
	// ═══════════════════════════════════════════════════════════════

	// HTTP Handlers
	handler := handlers.New(useCase, moduleLogger)
	handler.RegisterRoutes(router)

	// Consumers (RabbitMQ)
	if rabbitMQ != nil {
		consumers := consumer.NewConsumers(
			rabbitMQ,
			useCase,
			repos.SyncLog,
			moduleLogger,
		)

		// Iniciar Order Consumer (escucha events de órdenes)
		go func() {
			if err := consumers.Order.Start(ctx); err != nil {
				moduleLogger.Error(ctx).Err(err).Msg("Failed to start order consumer")
			}
		}()

		// Iniciar Retry Consumer (cron de reintentos cada 5 minutos)
		go consumers.Retry.Start(ctx)

		moduleLogger.Info(ctx).Msg("Consumers started successfully")
	} else {
		moduleLogger.Warn(ctx).Msg("RabbitMQ not available - consumers not started")
	}

	moduleLogger.Info(ctx).Msg("Invoicing module initialized successfully")
}
