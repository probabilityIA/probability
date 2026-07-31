package websiteconfig

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/websiteconfig/internal/app"
	"github.com/secamc93/probability/back/central/services/modules/websiteconfig/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/modules/websiteconfig/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/storage"
)

func New(router *gin.RouterGroup, database db.IDatabase, logger log.ILogger, s3 storage.IS3Service, environment env.IConfig) {
	repo := repository.New(database)
	uc := app.New(repo, logger)
	h := handlers.New(uc, logger, s3, environment)
	h.RegisterRoutes(router)
}
