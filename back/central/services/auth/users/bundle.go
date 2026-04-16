package users

import (
	"github.com/gin-gonic/gin"
	app "github.com/secamc93/probability/back/central/services/auth/users/internal/app"
	"github.com/secamc93/probability/back/central/services/auth/users/internal/infra/primary/handlers"

	"github.com/secamc93/probability/back/central/services/auth/users/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/storage"
)

// New inicializa el módulo de users
func New(
	router *gin.RouterGroup,
	db db.IDatabase,
	logger log.ILogger,
	cfg env.IConfig,
	s3 storage.IS3Service,
) {
	// 1. Inicializar Repositorio
	repo := repository.New(db, logger)

	// 2. Inicializar Caso de Uso
	userUC := app.New(repo, logger, s3, cfg)

	// 3. Inicializar Handler
	userH := handlers.New(userUC, logger)

	// 4. Registrar Rutas
	userH.RegisterRoutes(router, userH, logger)
}
