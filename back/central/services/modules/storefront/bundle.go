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
)

// New inicializa el modulo de storefront
func New(router *gin.RouterGroup, database db.IDatabase, logger log.ILogger, rabbitMQ rabbitmq.IQueue, environment env.IConfig) {
	// 1. Init Repository
	repo := repository.New(database)

	// 2. Init Publisher
	publisher := queue.NewStorefrontPublisher(rabbitMQ, logger)

	// 3. Init Use Cases
	uc := app.New(repo, logger, publisher)

	// 4. Init Handlers
	h := handlers.New(uc, logger, environment)

	// 5. Register Routes
	h.RegisterRoutes(router)
}
