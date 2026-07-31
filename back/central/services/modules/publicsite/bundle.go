package publicsite

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/pay"
	"github.com/secamc93/probability/back/central/services/modules/publicsite/internal/app"
	"github.com/secamc93/probability/back/central/services/modules/publicsite/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/modules/publicsite/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/storage"
)

func New(router *gin.RouterGroup, database db.IDatabase, logger log.ILogger, environment env.IConfig, payBundle *pay.Bundle, s3 storage.IS3Service) {
	repo := repository.New(database)
	uc := app.New(repo, payBundle, logger)
	h := handlers.New(uc, logger, environment, s3)
	h.RegisterRoutes(router)
}
