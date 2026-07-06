package demo

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/demo/internal/app"
	"github.com/secamc93/probability/back/central/services/auth/demo/internal/infra/primary/handlers"
	otpqueue "github.com/secamc93/probability/back/central/services/auth/demo/internal/infra/secondary/queue"
	"github.com/secamc93/probability/back/central/services/auth/demo/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/email"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

func New(router *gin.RouterGroup, database db.IDatabase, logger log.ILogger, cfg env.IConfig, queue rabbitmq.IQueue) {
	repo := repository.New(database, logger, cfg.Get("ENCRYPTION_KEY"))
	emailService := email.New(cfg, logger)
	otpPublisher := otpqueue.New(queue, logger)
	useCase := app.New(repo, emailService, otpPublisher, logger, cfg)
	handler := handlers.New(useCase, logger)
	handler.RegisterRoutes(router)
}
