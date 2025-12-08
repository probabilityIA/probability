package test

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/test/internal/app/generator"
	"github.com/secamc93/probability/back/central/services/integrations/test/internal/app/usecases"
	"github.com/secamc93/probability/back/central/services/integrations/test/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/integrations/test/internal/infra/primary/worker"
	"github.com/secamc93/probability/back/central/services/integrations/test/internal/infra/secondary/queue"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// New inicializa el módulo de test para generar órdenes aleatorias
func New(router *gin.RouterGroup, logger log.ILogger, rabbitMQ rabbitmq.IQueue) {
	// 1. Init Generator
	orderGenerator := generator.New()

	// 2. Init Publisher for RabbitMQ
	if rabbitMQ == nil {
		logger.Warn().Msg("RabbitMQ not available, test module will not be able to publish orders")
		return
	}

	orderPublisher := queue.New(rabbitMQ, logger)

	// 3. Init Use Cases
	uc := usecases.New(orderGenerator, orderPublisher)

	// 4. Init Handlers
	h := handlers.New(uc)

	// 5. Register Routes
	h.RegisterRoutes(router)

	// 6. Init and Start Scheduler (genera órdenes cada 30 segundos de las tres plataformas)
	businessID := uint(7)
	schedulerConfig := &worker.SchedulerConfig{
		Interval:        30 * time.Second,
		OrdersPerBatch:  1,
		IntegrationID:   1, // Se ignorará, se usa por plataforma
		BusinessID:      &businessID,
		Platform:        "", // Se rotará entre las tres plataformas
		Status:          "pending",
		IncludePayment:  true,
		IncludeShipment: true,
	}

	orderScheduler := worker.NewOrderScheduler(uc, logger, schedulerConfig)
	orderScheduler.Start(context.Background())

	logger.Info().
		Dur("interval", schedulerConfig.Interval).
		Int("orders_per_batch", schedulerConfig.OrdersPerBatch).
		Msg("Test module initialized - automatic order generation started")
}
