package actions

import (
	"github.com/secamc93/probability/back/central/services/auth/actions/internal/app"
	"github.com/secamc93/probability/back/central/services/auth/actions/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/auth/actions/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"

	"github.com/gin-gonic/gin"
)

// New inicializa y registra todos los componentes del módulo de actions
func New(database db.IDatabase, logger log.ILogger, v1Group *gin.RouterGroup) {
	// Repositorio
	repo := repository.New(database, logger)

	// Casos de uso
	useCase := app.New(repo, logger)

	// Handlers HTTP
	handlers := handlers.New(useCase, logger)

	// Rutas /actions
	handlers.RegisterRoutes(v1Group, handlers, logger)
}
