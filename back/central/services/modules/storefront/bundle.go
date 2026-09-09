package storefront

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/app"
	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/infra/secondary/queue"
	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
	"github.com/secamc93/probability/back/central/shared/storage"
)

func New(router *gin.RouterGroup, database db.IDatabase, logger log.ILogger, rabbitMQ rabbitmq.IQueue, environment env.IConfig, s3 storage.IS3Service) {
	repo := repository.New(database)
	publisher := queue.NewStorefrontPublisher(rabbitMQ, logger)
	uc := app.New(repo, logger, publisher)
	h := handlers.New(uc, logger, environment, s3)
	h.RegisterRoutes(router)
}
