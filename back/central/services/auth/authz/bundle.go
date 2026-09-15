package authz

import (
	"context"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/authz/internal/app"
	"github.com/secamc93/probability/back/central/services/auth/authz/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/auth/authz/internal/infra/primary/handlers"
	guardmw "github.com/secamc93/probability/back/central/services/auth/authz/internal/infra/primary/middleware"
	"github.com/secamc93/probability/back/central/services/auth/authz/internal/infra/secondary/cache"
	"github.com/secamc93/probability/back/central/services/auth/authz/internal/infra/secondary/repository"
	sharedauthz "github.com/secamc93/probability/back/central/shared/authz"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/redis"
)

const apiPrefix = "/api/v1"

type IModuleAccess = ports.IModuleAccess

type Bundle struct {
	UseCase app.IUseCase
	guard   *guardmw.Guard
	log     log.ILogger
}

func New(router *gin.RouterGroup, database db.IDatabase, redisClient redis.IRedis, logger log.ILogger) *Bundle {
	moduleLogger := logger.WithModule("authz")

	repo := repository.New(database)
	useCase := app.New(repo, nil, cache.New(redisClient), moduleLogger)
	handlers.New(useCase, moduleLogger).RegisterRoutes(router)

	return &Bundle{
		UseCase: useCase,
		guard:   guardmw.New(useCase, sharedauthz.RoutePolicies, moduleLogger, apiPrefix),
		log:     moduleLogger,
	}
}

func (b *Bundle) Middleware() gin.HandlerFunc {
	return b.guard.Handler()
}

func (b *Bundle) SetModuleAccess(modules IModuleAccess) {
	b.UseCase.SetModuleAccess(modules)
}

func (b *Bundle) VisibleNavigation(ctx context.Context, userID, tokenBusinessID, requestedBusinessID uint) ([]sharedauthz.NavEntry, error) {
	access, err := b.UseCase.GetEffectiveAccess(ctx, userID, tokenBusinessID, requestedBusinessID)
	if err != nil {
		return nil, err
	}

	visible := make(map[string]bool)
	for _, item := range b.UseCase.Navigation(access) {
		visible[item.Key] = true
	}

	entries := make([]sharedauthz.NavEntry, 0, len(visible))
	for _, entry := range sharedauthz.Navigation {
		if visible[entry.Key] {
			entries = append(entries, entry)
		}
	}
	return entries, nil
}

func (b *Bundle) ReportCoverage(ctx context.Context, routes gin.RoutesInfo) {
	missing := guardmw.Coverage(sharedauthz.RoutePolicies, routes, apiPrefix)
	if len(missing) > 0 {
		b.log.Warn(ctx).
			Int("undeclared_routes", len(missing)).
			Str("sample", strings.Join(firstN(missing, 10), " | ")).
			Msg("[authz] rutas sin politica declarada")
	} else {
		b.log.Info(ctx).Int("routes", len(routes)).Msg("[authz] todas las rutas tienen politica declarada")
	}

	if path := os.Getenv("AUTHZ_ROUTES_DUMP"); path != "" {
		if err := guardmw.DumpRoutes(sharedauthz.RoutePolicies, routes, apiPrefix, path); err != nil {
			b.log.Warn(ctx).Err(err).Msg("[authz] no se pudo escribir el volcado de rutas")
		}
	}
}

func firstN(list []string, n int) []string {
	if len(list) < n {
		return list
	}
	return list[:n]
}
