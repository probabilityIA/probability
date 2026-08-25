package shippingconfig

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/shippingconfig/internal/app"
	"github.com/secamc93/probability/back/central/services/modules/shippingconfig/internal/domain"
	"github.com/secamc93/probability/back/central/services/modules/shippingconfig/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/modules/shippingconfig/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
)

func New(router *gin.RouterGroup, database db.IDatabase) domain.IUseCase {
	repo := repository.New(database)
	uc := app.New(repo)
	h := handlers.New(uc)
	h.RegisterRoutes(router)
	return uc
}
