package invoicing

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/app"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/primary/queue/consumer"
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
	integrationCore core.IIntegrationCore,
) {
	ctx := context.Background()
	moduleLogger := logger.WithModule("invoicing")

	// ═══════════════════════════════════════════════════════════════
	// 1. INFRAESTRUCTURA SECUNDARIA (Adaptadores de salida)
	// ═══════════════════════════════════════════════════════════════

	// Repositorios (GORM) - solo para invoices, configs, sync logs
	repos := repository.New(database, moduleLogger)

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
		repos.Config,
		repos.SyncLog,
		repos.CreditNote,
		repos.Order,       // Order repository - implementado localmente
		integrationCore,   // Integration Core (reemplaza provider repos y client)
		nil,               // Encryption - TODO: agregar cuando esté disponible
		eventPublisher,    // Event publisher (RabbitMQ)
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
	} else {
		moduleLogger.Warn(ctx).Msg("RabbitMQ not available - consumers not started")
	}
}
