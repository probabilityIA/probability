package subscriptions

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/app"
	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/infra/primary/worker"
	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
)

type Bundle struct {
	UseCase app.IUseCase
}

func New(router *gin.RouterGroup, database db.IDatabase, logger log.ILogger, wallet ports.IWalletDebiter, announcements ports.IAnnouncementsGateway) *Bundle {
	moduleLogger := logger.WithModule("subscriptions")

	repo := repository.New(database)
	useCase := app.New(repo, wallet, announcements, moduleLogger)
	handler := handlers.New(useCase)
	handler.RegisterRoutes(router)

	expiryWorker := worker.New(useCase, moduleLogger)
	go expiryWorker.Start(context.Background())

	return &Bundle{UseCase: useCase}
}
