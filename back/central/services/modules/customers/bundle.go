package customers

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/customers/internal/app"
	"github.com/secamc93/probability/back/central/services/modules/customers/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/modules/customers/internal/infra/primary/queue"
	"github.com/secamc93/probability/back/central/services/modules/customers/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

func New(router *gin.RouterGroup, database db.IDatabase, logger log.ILogger, rabbitMQ rabbitmq.IQueue) {
	repo := repository.New(database)
	uc := app.New(repo, logger)
	h := handlers.New(uc)
	h.RegisterRoutes(router)

	consumer := queue.NewOrderConsumer(rabbitMQ, uc, logger)
	consumer.Start(context.Background())
}
