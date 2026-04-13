package announcements

import (
	"github.com/gin-gonic/gin"

	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/app"
	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/infra/secondary/repository"
	storageadapter "github.com/secamc93/probability/back/central/services/modules/announcements/internal/infra/secondary/storage"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/storage"
)

func New(router *gin.RouterGroup, database db.IDatabase, logger log.ILogger, s3 storage.IS3Service) {
	logger = logger.WithModule("announcements")

	repo := repository.New(database)
	storageService := storageadapter.New(s3)
	useCase := app.New(repo, storageService, logger)
	handler := handlers.New(useCase, logger)
	handler.RegisterRoutes(router)
}
