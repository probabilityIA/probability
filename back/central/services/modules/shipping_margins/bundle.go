package shipping_margins

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/shipping_margins/internal/app"
	"github.com/secamc93/probability/back/central/services/modules/shipping_margins/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/modules/shipping_margins/internal/infra/secondary/cache"
	"github.com/secamc93/probability/back/central/services/modules/shipping_margins/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/redis"
)

type Bundle struct {
	UseCase app.IUseCase
}

func New(router *gin.RouterGroup, database db.IDatabase, logger log.ILogger, redisClient redis.IRedis) *Bundle {
	repo := repository.New(database)
	cacheWriter := cache.New(redisClient, logger)
	uc := app.New(repo, cacheWriter)
	h := handlers.New(uc)
	h.RegisterRoutes(router)
	return &Bundle{UseCase: uc}
}
