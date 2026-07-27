package syncruns

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/syncruns/internal/app"
	"github.com/secamc93/probability/back/central/services/integrations/syncruns/internal/infra/primary/handlers"
	syncqueue "github.com/secamc93/probability/back/central/services/integrations/syncruns/internal/infra/primary/queue"
	"github.com/secamc93/probability/back/central/services/integrations/syncruns/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

func New(router *gin.RouterGroup, database db.IDatabase, logger log.ILogger, rabbitMQ rabbitmq.IQueue) {
	repo := repository.New(database, logger)
	useCase := app.New(repo, logger)

	handler := handlers.New(useCase, logger)
	handler.RegisterRoutes(router)

	syncqueue.New(rabbitMQ, useCase, logger).Start(context.Background())
}
