package dashboard

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/dashboard/internal/app"
	"github.com/secamc93/probability/back/central/services/modules/dashboard/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/modules/dashboard/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
)

// New inicializa el módulo de dashboard
func New(router *gin.RouterGroup, database db.IDatabase, logger log.ILogger) {
	logger = logger.WithModule("Dashboard")

	// 1. Repositorio
	repo := repository.New(database, logger)

	// 2. Casos de uso
	uc := app.New(repo, logger)

	// 3. Handlers
	h := handlers.New(uc, logger)

	// 4. Rutas
	h.RegisterRoutes(router)
}
