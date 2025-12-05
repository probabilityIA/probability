package orders

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/orders/app/usecases"
	"github.com/secamc93/probability/back/central/services/modules/orders/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/modules/orders/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

// New inicializa el módulo de orders
func New(router *gin.RouterGroup, database db.IDatabase, logger log.ILogger, environment env.IConfig) {
	// 1. Init Repositories
	repo := repository.New(database)

	// 2. Init Use Cases
	uc := usecases.New(repo)

	// 3. Init Handlers
	h := handlers.New(uc)

	// 4. Register Routes
	h.RegisterRoutes(router)
}
