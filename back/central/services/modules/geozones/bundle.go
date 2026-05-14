package geozones

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/app"
	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/infra/primary/handlers"
	probqueue "github.com/secamc93/probability/back/central/services/modules/geozones/internal/infra/primary/queue"
	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/infra/secondary/cache"
	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/infra/secondary/metrics"
	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
	"github.com/secamc93/probability/back/central/shared/redis"
)

const aggregateRefreshInterval = 6 * time.Hour

type Bundle struct {
	UseCase          app.IUseCase
	Resolver         ports.IResolver
	ProbabilityRepo  ports.IProbabilityRepository
	ProbabilityCache ports.IProbabilityCache
}

func New(router *gin.RouterGroup, database db.IDatabase, logger log.ILogger, rdb redis.IRedis, queue rabbitmq.IQueue) *Bundle {
	repoStruct := repository.NewStruct(database)
	displayCache := cache.New(rdb, logger)
	probabilityCache := cache.NewProbabilityCache(rdb, logger)
	uc := app.New(repoStruct, displayCache)
	resolver := app.NewResolver(repoStruct)
	probabilityUC := app.NewProbability(repoStruct, repoStruct, resolver, probabilityCache)
	h := handlers.New(uc, probabilityUC)
	h.RegisterRoutes(router)

	ctx := context.Background()
	go startAggregateRefresh(ctx, repoStruct, logger)

	if queue != nil {
		consumer := probqueue.NewProbabilityConsumer(queue, probabilityCache, logger)
		consumer.Start(ctx)
	}

	return &Bundle{
		UseCase:          uc,
		Resolver:         resolver,
		ProbabilityRepo:  repoStruct,
		ProbabilityCache: probabilityCache,
	}
}

func startAggregateRefresh(ctx context.Context, repo ports.IProbabilityRepository, logger log.ILogger) {
	run := func() {
		start := time.Now()
		if err := repo.RefreshAggregates(ctx); err != nil {
			logger.Warn(ctx).Err(err).Msg("geozone_carrier_stats refresh failed")
			return
		}
		elapsed := time.Since(start)
		metrics.AggregateRefreshDuration.Observe(elapsed.Seconds())
		metrics.AggregateRefreshLastSuccess.SetToCurrentTime()
		logger.Info(ctx).Dur("duration", elapsed).Msg("geozone_carrier_stats refresh completed")
	}
	run()
	ticker := time.NewTicker(aggregateRefreshInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			run()
		}
	}
}
