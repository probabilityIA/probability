package login

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/login/internal/app"
	authhandler "github.com/secamc93/probability/back/central/services/auth/login/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/auth/login/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/jwt"
	"github.com/secamc93/probability/back/central/shared/log"
)

// New inicializa el módulo de login
func New(
	router *gin.RouterGroup,
	db db.IDatabase,
	logger log.ILogger,
	cfg env.IConfig,
) {
	// 1. Inicializar Repositorio
	repo := repository.New(db, logger)

	// 2. Inicializar Servicio JWT
	jwtService := jwt.New(cfg.Get("JWT_SECRET"))

	// 3. Inicializar Caso de Uso
	// Al usar type aliases en el dominio, jwtService satisface la interfaz domain.IJWTService
	authUC := app.New(repo, jwtService, logger, cfg)

	// 4. Inicializar Handler
	authH := authhandler.New(authUC, logger)

	// 5. Registrar Rutas
	authH.RegisterRoutes(router, authH, logger)
}
