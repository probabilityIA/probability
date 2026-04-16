package permissions

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/permissions/internal/app"
	permissionhandler "github.com/secamc93/probability/back/central/services/auth/permissions/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/auth/permissions/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
)

// New inicializa el módulo de permissions
func New(
	router *gin.RouterGroup,
	db db.IDatabase,
	logger log.ILogger,
) {
	// 1. Inicializar Repositorio
	repo := repository.New(db, logger)

	// 2. Inicializar Caso de Uso
	permissionUC := app.New(repo, logger)

	// 3. Inicializar Handler
	permissionH := permissionhandler.New(permissionUC, logger)

	// 4. Registrar Rutas
	permissionH.RegisterRoutes(router, permissionH, logger)
}
