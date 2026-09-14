package authz

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/authz/internal/app"
	"github.com/secamc93/probability/back/central/services/auth/authz/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/auth/authz/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/auth/authz/internal/infra/secondary/cache"
	"github.com/secamc93/probability/back/central/services/auth/authz/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/redis"
)

type IModuleAccess = ports.IModuleAccess

type Bundle struct {
	UseCase app.IUseCase
}

func New(router *gin.RouterGroup, database db.IDatabase, redisClient redis.IRedis, logger log.ILogger, moduleAccess IModuleAccess) *Bundle {
	moduleLogger := logger.WithModule("authz")

	repo := repository.New(database)
	useCase := app.New(repo, moduleAccess, cache.New(redisClient), moduleLogger)
	handlers.New(useCase, moduleLogger).RegisterRoutes(router)

	return &Bundle{UseCase: useCase}
}
