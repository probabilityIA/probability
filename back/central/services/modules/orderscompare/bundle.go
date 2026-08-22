package orderscompare

import (
	"github.com/gin-gonic/gin"
	integrationcore "github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/modules/orderscompare/internal/app"
	"github.com/secamc93/probability/back/central/services/modules/orderscompare/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/modules/orderscompare/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/modules/orderscompare/internal/infra/secondary/channels"
	"github.com/secamc93/probability/back/central/services/modules/orderscompare/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
)

type Bundle struct {
	UseCase ports.IUseCase
}

func New(router *gin.RouterGroup, database db.IDatabase, logger log.ILogger, integrationCore integrationcore.IIntegrationCore) *Bundle {
	logger = logger.WithModule("orderscompare")

	repo := repository.New(database, logger)
	registry := channels.New(integrationCore)
	useCase := app.New(repo, registry, logger)

	handlers.New(useCase, logger).RegisterRoutes(router)

	return &Bundle{UseCase: useCase}
}
