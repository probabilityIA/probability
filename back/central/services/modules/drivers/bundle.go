package drivers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/drivers/internal/app"
	"github.com/secamc93/probability/back/central/services/modules/drivers/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/modules/drivers/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
)

func New(router *gin.RouterGroup, database db.IDatabase) {
	repo := repository.New(database)
	uc := app.New(repo)
	h := handlers.New(uc)
	h.RegisterRoutes(router)
}
