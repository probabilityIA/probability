package mobile

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/mobile/internal/app"
	"github.com/secamc93/probability/back/central/services/modules/mobile/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/modules/mobile/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
)

func New(router *gin.RouterGroup, database db.IDatabase) {
	repo := repository.New(database)
	uc := app.New(repo)
	h := handlers.New(uc)
	h.RegisterRoutes(router)
}
