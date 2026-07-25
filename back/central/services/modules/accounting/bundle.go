package accounting

import (
	"context"

	"github.com/gin-gonic/gin"
	integrationsCore "github.com/secamc93/probability/back/central/services/integrations/core"
	factusinv "github.com/secamc93/probability/back/central/services/integrations/invoicing/factus"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/app"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/infra/primary/worker"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/infra/secondary/dian"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/email"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

type Bundle struct {
	UseCase app.IUseCase
}

func New(router *gin.RouterGroup, database db.IDatabase, logger log.ILogger, cfg env.IConfig, integrationCore integrationsCore.IIntegrationCore, dianEmitter *factusinv.PlatformEmitter) *Bundle {
	repo := repository.New(database)
	emailService := email.New(cfg, logger)
	uc := app.New(repo, emailService, dian.NewFactusEmitter(dianEmitter), dian.NewIntegrationAdapter(integrationCore), logger)
	h := handlers.New(uc)
	h.RegisterRoutes(router)

	syncWorker := worker.NewSyncWorker(uc, logger)
	go syncWorker.Start(context.Background())

	return &Bundle{UseCase: uc}
}
