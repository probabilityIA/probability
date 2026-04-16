package resources

import (
	"github.com/secamc93/probability/back/central/services/auth/resources/internal/app"
	"github.com/secamc93/probability/back/central/services/auth/resources/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/auth/resources/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"

	"github.com/gin-gonic/gin"
)

// New inicializa y registra todos los componentes del módulo de recursos
func New(db db.IDatabase, logger log.ILogger, v1Group *gin.RouterGroup) {
	Repository := repository.New(db, logger)
	UseCase := app.New(Repository, logger)
	Handlers := handlers.New(UseCase, logger)
	handlers.RegisterRoutes(v1Group, Handlers, logger)
}
