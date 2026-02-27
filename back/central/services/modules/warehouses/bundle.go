package warehouses

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/app"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
)

// New inicializa el módulo de warehouses
func New(router *gin.RouterGroup, database db.IDatabase) {
	// 1. Init Repository
	repo := repository.New(database)

	// 2. Init Use Cases
	uc := app.New(repo)

	// 3. Init Handlers
	h := handlers.New(uc)

	// 4. Register Routes
	h.RegisterRoutes(router)
}
