package tours

import (
	"github.com/gin-gonic/gin"

	"github.com/secamc93/probability/back/central/services/modules/tours/internal/app"
	"github.com/secamc93/probability/back/central/services/modules/tours/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/modules/tours/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
)

type Bundle struct {
	UseCase app.IUseCase
}

func New(router *gin.RouterGroup, database db.IDatabase, logger log.ILogger) *Bundle {
	logger = logger.WithModule("tours")

	repo := repository.New(database)
	useCase := app.New(repo, logger)
	handler := handlers.New(useCase, logger)
	handler.RegisterRoutes(router)

	return &Bundle{UseCase: useCase}
}
