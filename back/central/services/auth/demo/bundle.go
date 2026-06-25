package demo

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/demo/internal/app"
	"github.com/secamc93/probability/back/central/services/auth/demo/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/auth/demo/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/email"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

func New(router *gin.RouterGroup, database db.IDatabase, logger log.ILogger, cfg env.IConfig) {
	repo := repository.New(database, logger, cfg.Get("ENCRYPTION_KEY"))
	emailService := email.New(cfg, logger)
	useCase := app.New(repo, emailService, logger, cfg)
	handler := handlers.New(useCase, logger)
	handler.RegisterRoutes(router)
}
