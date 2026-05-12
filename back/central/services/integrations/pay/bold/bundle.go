package bold

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/pay/bold/internal/app"
	"github.com/secamc93/probability/back/central/services/integrations/pay/bold/internal/infra/primary/consumer"
	"github.com/secamc93/probability/back/central/services/integrations/pay/bold/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/integrations/pay/bold/internal/infra/secondary/client"
	boldqueue "github.com/secamc93/probability/back/central/services/integrations/pay/bold/internal/infra/secondary/queue"
	"github.com/secamc93/probability/back/central/services/integrations/pay/bold/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

const rawWebhookRetentionDays = 15

func New(
	router *gin.RouterGroup,
	coreSvc core.IIntegrationCore,
	logger log.ILogger,
	database db.IDatabase,
	rabbit rabbitmq.IQueue,
) {
	logger = logger.WithModule("bold")
	ctx := context.Background()

	boldClient := client.New(logger)
	integrationRepo := repository.New(coreSvc, logger)
	rawWebhookLog := repository.NewRawWebhookLogRepository(database, logger)
	if intType, err := coreSvc.GetIntegrationTypeByCode(ctx, "bold_pay"); err == nil && intType != nil {
		if setter, ok := rawWebhookLog.(interface{ SetIntegrationTypeID(uint) }); ok {
			setter.SetIntegrationTypeID(intType.ID)
		}
	}
	responsePublisher := boldqueue.New(rabbit, logger)
	webhookPublisher := boldqueue.NewWebhookPublisher(rabbit, logger)

	useCase := app.New(boldClient, integrationRepo, responsePublisher, logger)
	webhookUseCase := app.NewWebhookUseCase(integrationRepo, webhookPublisher, logger)

	webhookHandlers := handlers.NewWebhookHandlers(webhookUseCase, rawWebhookLog, logger)
	webhookHandlers.RegisterRoutes(router)

	go startRawWebhookRetention(ctx, rawWebhookLog, logger)

	if rabbit != nil {
		boldConsumer := consumer.New(rabbit, useCase, logger)
		go func() {
			if err := boldConsumer.Start(ctx); err != nil {
				logger.Error(ctx).Err(err).Msg("Bold consumer failed")
			}
		}()
		logger.Info(ctx).Msg("Bold consumer started")
	} else {
		logger.Warn(ctx).Msg("RabbitMQ not available - Bold consumer disabled")
	}
}

func startRawWebhookRetention(ctx context.Context, repo interface {
	DeleteOlderThan(ctx context.Context, days int) (int64, error)
}, logger log.ILogger) {
	if repo == nil {
		return
	}
	run := func() {
		deleted, err := repo.DeleteOlderThan(ctx, rawWebhookRetentionDays)
		if err != nil {
			logger.Warn(ctx).Err(err).Msg("bold raw webhook retention cleanup failed")
			return
		}
		if deleted > 0 {
			logger.Info(ctx).Int64("deleted", deleted).Int("retention_days", rawWebhookRetentionDays).Msg("bold raw webhook logs purged")
		}
	}
	run()
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			run()
		}
	}
}
