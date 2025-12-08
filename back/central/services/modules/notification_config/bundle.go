package notification_config

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/app"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
)

func New(router *gin.RouterGroup, database db.IDatabase) {
	// 1. Init Repository
	repo := repository.New(database)

	// 2. Init Use Case
	useCase := app.New(repo)

	// 3. Init Handler
	handler := handlers.New(useCase)

	// 4. Register Routes
	handler.RegisterRoutes(router)
}
