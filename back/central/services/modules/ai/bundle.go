package ai

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/ai/internal/app"
	"github.com/secamc93/probability/back/central/services/modules/ai/internal/infra/primary/handlers"
	bedrockadapter "github.com/secamc93/probability/back/central/services/modules/ai/internal/infra/secondary/bedrock"
	"github.com/secamc93/probability/back/central/services/modules/ai/internal/infra/secondary/cache"
	"github.com/secamc93/probability/back/central/services/modules/ai/internal/infra/secondary/navigation"
	"github.com/secamc93/probability/back/central/services/modules/ai/internal/infra/secondary/openrouter"
	"github.com/secamc93/probability/back/central/shared/bedrock"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/redis"
)

type VisibleNavigationFunc = navigation.VisibleNavigationFunc

func New(router *gin.RouterGroup, logger log.ILogger, cfg env.IConfig, bedrockClient bedrock.IBedrock, redisClient redis.IRedis, visibleNavigation VisibleNavigationFunc) {
	moduleLogger := logger.WithModule("ai")

	useCase := app.New(
		openrouter.New(moduleLogger),
		bedrockadapter.New(bedrockClient, cfg),
		navigation.New(visibleNavigation),
		cache.New(redisClient),
		moduleLogger,
	)

	handlers.New(useCase, moduleLogger).RegisterRoutes(router)
}
