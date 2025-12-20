package paymentstatus

import (
	"github.com/gin-gonic/gin"

	"github.com/secamc93/probability/back/central/services/modules/paymentstatus/app"
	"github.com/secamc93/probability/back/central/services/modules/paymentstatus/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/modules/paymentstatus/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

// New inicializa el módulo de payment statuses
func New(router *gin.RouterGroup, database db.IDatabase, logger log.ILogger, environment env.IConfig) {
	logger = logger.WithModule("Payment Status")

	// 1. Repositorio
	repo := repository.New(database, logger)

	// 2. Casos de uso
	uc := app.New(repo, logger)

	// 3. Handlers
	h := handlers.New(uc, logger, environment)

	// 4. Rutas
	h.RegisterRoutes(router)
}
