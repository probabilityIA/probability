package login

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/login/internal/app"
	authhandler "github.com/secamc93/probability/back/central/services/auth/login/internal/infra/primary/handlers"
	otpqueue "github.com/secamc93/probability/back/central/services/auth/login/internal/infra/secondary/queue"
	"github.com/secamc93/probability/back/central/services/auth/login/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/email"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/jwt"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// New inicializa el módulo de login
func New(
	router *gin.RouterGroup,
	db db.IDatabase,
	logger log.ILogger,
	cfg env.IConfig,
	queue rabbitmq.IQueue,
) {
	// 1. Inicializar Repositorio
	repo := repository.New(db, logger)

	// 2. Inicializar Servicio JWT
	jwtService := jwt.New(cfg.Get("JWT_SECRET"))

	emailService := email.New(cfg, logger)

	otpPublisher := otpqueue.New(queue, logger)

	// 3. Inicializar Caso de Uso
	authUC := app.New(repo, jwtService, emailService, otpPublisher, logger, cfg)

	// 4. Inicializar Handler
	authH := authhandler.New(authUC, logger)

	// 5. Registrar Rutas
	authH.RegisterRoutes(router, authH, logger)
}
