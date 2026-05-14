package geozones

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/app"
	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/infra/secondary/cache"
	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/redis"
)

type Bundle struct {
	UseCase  app.IUseCase
	Resolver ports.IResolver
}

func New(router *gin.RouterGroup, database db.IDatabase, logger log.ILogger, rdb redis.IRedis) *Bundle {
	repoStruct := repository.NewStruct(database)
	displayCache := cache.New(rdb, logger)
	uc := app.New(repoStruct, displayCache)
	resolver := app.NewResolver(repoStruct)
	probabilityUC := app.NewProbability(repoStruct, repoStruct, resolver)
	h := handlers.New(uc, probabilityUC)
	h.RegisterRoutes(router)
	return &Bundle{UseCase: uc, Resolver: resolver}
}
