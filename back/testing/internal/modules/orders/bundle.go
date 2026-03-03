package orders

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/testing/internal/modules/orders/internal/app"
	"github.com/secamc93/probability/back/testing/internal/modules/orders/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/testing/internal/modules/orders/internal/infra/secondary/client"
	"github.com/secamc93/probability/back/testing/internal/modules/orders/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/testing/internal/shared/db"
	"github.com/secamc93/probability/back/testing/shared/log"
)

func New(router *gin.RouterGroup, database db.IDatabase, centralAPIURL string, logger log.ILogger) {
	repo := repository.New(database)
	centralClient := client.New(centralAPIURL)
	useCase := app.New(repo, centralClient, logger)
	handler := handlers.New(useCase, logger)

	handler.RegisterRoutes(router)

	logger.Info().Msg("Orders testing module initialized")
}
