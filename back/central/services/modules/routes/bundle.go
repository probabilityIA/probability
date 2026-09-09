package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/routes/internal/app"
	"github.com/secamc93/probability/back/central/services/modules/routes/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/modules/routes/internal/infra/secondary/maps"
	"github.com/secamc93/probability/back/central/services/modules/routes/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

func New(router *gin.RouterGroup, database db.IDatabase, logger log.ILogger, cfg env.IConfig) {
	logger = logger.WithModule("routes")

	repo := repository.New(database)
	optimizer := maps.New(cfg, logger)
	uc := app.New(repo, optimizer)
	h := handlers.New(uc)
	h.RegisterRoutes(router)
}
