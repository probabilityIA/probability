package ai

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/ai/internal/app/usecases"
	"github.com/secamc93/probability/back/central/services/modules/ai/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/modules/ai/internal/infra/secondary/openrouter"
	"github.com/secamc93/probability/back/central/shared/log"
)

func New(router *gin.RouterGroup, logger log.ILogger) {
	// dependencies
	client := openrouter.New(logger)
	useCase := usecases.NewGetRecommendationUseCase(client)
	handler := handlers.NewRecommendationHandler(useCase, logger)

	// routes
	group := router.Group("/ai")
	{
		group.GET("/recommendation", handler.GetRecommendation)
	}
}
