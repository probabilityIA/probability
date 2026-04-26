package tickets

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/app"
	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/infra/secondary/repository"
	storageadapter "github.com/secamc93/probability/back/central/services/modules/tickets/internal/infra/secondary/storage"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/storage"
)

func New(router *gin.RouterGroup, database db.IDatabase, logger log.ILogger, s3 storage.IS3Service) {
	logger = logger.WithModule("tickets")
	repo := repository.New(database)
	storageSvc := storageadapter.New(s3)
	uc := app.New(repo, storageSvc, logger)
	h := handlers.New(uc, logger)
	h.RegisterRoutes(router)
}
