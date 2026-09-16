package ai

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/ai/internal/app"
	"github.com/secamc93/probability/back/central/services/modules/ai/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/modules/ai/internal/infra/primary/queue/consumer"
	bedrockadapter "github.com/secamc93/probability/back/central/services/modules/ai/internal/infra/secondary/bedrock"
	"github.com/secamc93/probability/back/central/services/modules/ai/internal/infra/secondary/businessdata"
	"github.com/secamc93/probability/back/central/services/modules/ai/internal/infra/secondary/cache"
	"github.com/secamc93/probability/back/central/services/modules/ai/internal/infra/secondary/navigation"
	"github.com/secamc93/probability/back/central/services/modules/ai/internal/infra/secondary/openrouter"
	"github.com/secamc93/probability/back/central/services/modules/ai/internal/infra/secondary/queue"
	"github.com/secamc93/probability/back/central/services/modules/ai/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/bedrock"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
	"github.com/secamc93/probability/back/central/shared/redis"
)

type VisibleNavigationFunc = navigation.VisibleNavigationFunc

type Dependencies struct {
	Config            env.IConfig
	Database          db.IDatabase
	RabbitMQ          rabbitmq.IQueue
	Bedrock           bedrock.IBedrock
	Redis             redis.IRedis
	VisibleNavigation VisibleNavigationFunc
}

func New(router *gin.RouterGroup, logger log.ILogger, deps Dependencies) {
	moduleLogger := logger.WithModule("ai")
	ctx := context.Background()

	useCase := app.New(
		openrouter.New(moduleLogger),
		bedrockadapter.New(deps.Bedrock, deps.Config),
		navigation.New(deps.VisibleNavigation),
		cache.New(deps.Redis),
		queue.New(deps.RabbitMQ),
		repository.New(deps.Database),
		businessdata.New(deps.Database),
		repository.NewAlerts(deps.Database),
		moduleLogger,
	)

	handlers.New(useCase, moduleLogger).RegisterRoutes(router)

	if deps.RabbitMQ != nil {
		queueConsumer := consumer.New(deps.RabbitMQ, useCase, moduleLogger)
		go func() {
			if err := queueConsumer.Start(ctx); err != nil {
				moduleLogger.Error(ctx).Err(err).Msg("[ai.assistant] no se pudo iniciar el consumidor de conversaciones")
			}
		}()
		go func() {
			if err := queueConsumer.StartAlerts(ctx); err != nil {
				moduleLogger.Error(ctx).Err(err).Msg("[ai.assistant] no se pudo iniciar el consumidor de alertas")
			}
		}()
	}

	if deps.Database != nil {
		go startConversationRetention(ctx, useCase, moduleLogger)
	}
}

func startConversationRetention(ctx context.Context, useCase app.IUseCase, logger log.ILogger) {
	run := func() {
		deleted, err := useCase.PurgeExpiredConversations(ctx)
		if err != nil {
			logger.Warn(ctx).Err(err).Msg("[ai.assistant] fallo la limpieza de conversaciones vencidas")
			return
		}
		if deleted > 0 {
			logger.Info(ctx).Int64("deleted", deleted).Msg("[ai.assistant] mensajes de mas de un ano eliminados")
		}
		purged, err := useCase.PurgeExpiredAlerts(ctx)
		if err != nil {
			logger.Warn(ctx).Err(err).Msg("[ai.assistant] fallo la limpieza de alertas vencidas")
			return
		}
		if purged > 0 {
			logger.Info(ctx).Int64("deleted", purged).Msg("[ai.assistant] alertas de mas de 90 dias eliminadas")
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
