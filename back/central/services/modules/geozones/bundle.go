package geozones

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/app"
	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/infra/secondary/cache"
	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/redis"
)

type Bundle struct {
	UseCase app.IUseCase
}

func New(router *gin.RouterGroup, database db.IDatabase, logger log.ILogger, rdb redis.IRedis) *Bundle {
	repo := repository.New(database)
	displayCache := cache.New(rdb, logger)
	uc := app.New(repo, displayCache)
	h := handlers.New(uc)
	h.RegisterRoutes(router)
	return &Bundle{UseCase: uc}
}
