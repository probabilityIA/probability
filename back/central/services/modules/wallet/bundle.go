package wallet

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/wallet/app/usecases"
	"github.com/secamc93/probability/back/central/services/modules/wallet/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/modules/wallet/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/services/modules/wallet/infra/secondary/services"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

func New(router *gin.RouterGroup, database db.IDatabase, logger log.ILogger, environment env.IConfig) {
	// 1. Init Dependencies
	repo := repository.New(database)
	nequi := services.NewNequiService(environment, logger)

	// 2. Init Use Cases
	uc := usecases.New(repo, nequi)

	// 3. Init Handlers
	h := handlers.New(uc)

	// 4. Register Routes
	h.RegisterRoutes(router)
}
