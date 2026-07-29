package marketingleads

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/marketingleads/internal/app"
	"github.com/secamc93/probability/back/central/services/modules/marketingleads/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/modules/marketingleads/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/modules/marketingleads/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
)

type IWhatsAppSender = ports.IWhatsAppSender

func New(router *gin.RouterGroup, database db.IDatabase, logger log.ILogger, wa IWhatsAppSender) {
	repo := repository.New(database)
	uc := app.New(repo, wa, logger)
	h := handlers.New(uc, logger)
	h.RegisterRoutes(router)
}
