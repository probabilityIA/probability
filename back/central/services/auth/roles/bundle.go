package roles

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/roles/internal/app"
	rolehandler "github.com/secamc93/probability/back/central/services/auth/roles/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/auth/roles/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
)

// New inicializa el módulo de roles
func New(
	router *gin.RouterGroup,
	db db.IDatabase,
	logger log.ILogger,
) {
	// 1. Inicializar Repositorio
	repo := repository.New(db, logger)

	// 2. Inicializar Caso de Uso
	roleUC := app.New(repo, logger)

	// 3. Inicializar Handler
	roleH := rolehandler.New(roleUC, logger)

	// 4. Registrar Rutas
	roleH.RegisterRoutes(router, roleH, logger)
}
