package publicsite

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/publicsite/internal/app"
	"github.com/secamc93/probability/back/central/services/modules/publicsite/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/modules/publicsite/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

func New(router *gin.RouterGroup, database db.IDatabase, logger log.ILogger, environment env.IConfig) {
	repo := repository.New(database)
	uc := app.New(repo, logger)
	h := handlers.New(uc, logger, environment)
	h.RegisterRoutes(router)
}
