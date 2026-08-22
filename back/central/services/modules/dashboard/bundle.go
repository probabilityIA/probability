package dashboard

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/dashboard/internal/app"
	"github.com/secamc93/probability/back/central/services/modules/dashboard/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/modules/dashboard/internal/infra/secondary/cache"
	"github.com/secamc93/probability/back/central/services/modules/dashboard/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/redis"
)

func New(router *gin.RouterGroup, database db.IDatabase, redisClient redis.IRedis, logger log.ILogger) {
	logger = logger.WithModule("Dashboard")

	repo := repository.New(database, logger)
	statsCache := cache.New(redisClient, logger)

	uc := app.New(repo, statsCache, logger)

	h := handlers.New(uc, logger)

	h.RegisterRoutes(router)
}
