package sprints

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/sprints/internal/app"
	"github.com/secamc93/probability/back/central/services/modules/sprints/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/modules/sprints/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
)

func New(router *gin.RouterGroup, database db.IDatabase, logger log.ILogger) {
	logger = logger.WithModule("sprints")
	repo := repository.New(database)
	uc := app.New(repo, logger)
	h := handlers.New(uc, logger)
	h.RegisterRoutes(router)
}
