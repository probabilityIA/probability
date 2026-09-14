package server

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/cmd/internal/routes"
	"github.com/secamc93/probability/back/central/services/auth"
	"github.com/secamc93/probability/back/central/services/auth/authz"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/events"
	"github.com/secamc93/probability/back/central/services/integrations"
	"github.com/secamc93/probability/back/central/services/modules"
	"github.com/secamc93/probability/back/central/shared/bedrock"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/email"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
	"github.com/secamc93/probability/back/central/shared/redis"
	"github.com/secamc93/probability/back/central/shared/storage"
)

func Init(ctx context.Context) error {
	logger := log.New()
	environment := env.New(logger)

	database := db.New(logger, environment)
	emailService := email.New(environment, logger)

	s3Service := storage.New(environment, logger)

	queueRegistry := NewQueueRegistry()
	rabbitMQ, err := rabbitmq.New(logger, environment)
	if err != nil {
		logger.Error(ctx).
			Err(err).
			Msg("Failed to initialize RabbitMQ - consumers will be disabled")
		rabbitMQ = rabbitmq.NewNoop()
	} else {
		if rmq, ok := rabbitMQ.(interface {
			SetQueueRegistry(rabbitmq.QueueRegistryCallback)
		}); ok {
			rmq.SetQueueRegistry(queueRegistry.Register)
		}
	}

	redisRegistry := NewRedisRegistry()
	redisClient := redis.New(logger, environment)
	if redisClient != nil {
		if rc, ok := redisClient.(interface {
			SetCacheRegistry(redis.CacheRegistryCallback)
			SetChannelRegistry(redis.ChannelRegistryCallback)
		}); ok {
			rc.SetCacheRegistry(redisRegistry.RegisterCachePrefix)
			rc.SetChannelRegistry(redisRegistry.RegisterChannel)
		}
	}

	bedrockClient := bedrock.New(logger, environment)

	middleware.InitFromEnv(environment, logger)
	r := routes.BuildRouter(ctx, logger, environment)

	v1Group := r.Group("/api/v1")

	authBundle := auth.New(v1Group, database, logger, environment, s3Service, rabbitMQ)

	events.New(v1Group, logger, rabbitMQ, redisClient)

	integrationCore, dianEmitter := integrations.New(v1Group, database, logger, environment, rabbitMQ, s3Service, redisClient, emailService)

	modulesBundle := modules.New(v1Group, database, logger, environment, rabbitMQ, redisClient, s3Service, bedrockClient, integrationCore, dianEmitter)

	authz.New(v1Group, database, redisClient, logger, modulesBundle.Subscriptions.UseCase)

	authBundle.Demo.SetOnBusinessCreated(modulesBundle.Subscriptions.UseCase.AssignTrialSubscription)
	authBundle.Business.SetOnBusinessCreated(modulesBundle.Subscriptions.UseCase.AssignTrialSubscription)

	LogStartupInfo(ctx, logger, environment, queueRegistry, redisRegistry)

	port := environment.Get("HTTP_PORT")

	addr := fmt.Sprintf(":%s", port)
	return r.Run(addr)
}
