package push

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/secamc93/probability/back/central/services/modules/push/internal/app"
	"github.com/secamc93/probability/back/central/services/modules/push/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/modules/push/internal/infra/primary/queue/consumer"
	"github.com/secamc93/probability/back/central/services/modules/push/internal/infra/secondary/fcm"
	"github.com/secamc93/probability/back/central/services/modules/push/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type Bundle struct {
	UseCase app.IUseCase
}

func New(router *gin.RouterGroup, database db.IDatabase, logger log.ILogger, cfg env.IConfig, queue rabbitmq.IQueue) *Bundle {
	logger = logger.WithModule("push")

	repo := repository.New(database)
	sender := fcm.New(cfg, logger)
	useCase := app.New(repo, sender, logger)

	handler := handlers.New(useCase, logger)
	handler.RegisterRoutes(router)

	if queue != nil {
		pushConsumer := consumer.New(queue, useCase, logger)
		if err := pushConsumer.Start(context.Background()); err != nil {
			logger.Error(context.Background()).Err(err).Msg("no se pudo iniciar el consumidor de push")
		}
	}

	return &Bundle{UseCase: useCase}
}
