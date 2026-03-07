package products

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/products/internal/app/usecases"
	"github.com/secamc93/probability/back/central/services/modules/products/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/modules/products/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/storage"
)

// New inicializa el módulo de products
func New(router *gin.RouterGroup, database db.IDatabase, logger log.ILogger, environment env.IConfig, s3 storage.IS3Service) {
	// 1. Init Repositories
	repo := repository.New(database)

	// 2. Init Use Cases
	uc := usecases.New(repo)

	// 3. Init Handlers
	h := handlers.New(uc, logger, s3, environment)

	// 4. Register Routes
	h.RegisterRoutes(router)
}
