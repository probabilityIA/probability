package server

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/cmd/internal/routes"
	"github.com/secamc93/probability/back/central/services/auth"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/integrations"
	"github.com/secamc93/probability/back/central/services/modules"
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
	_ = email.New(environment, logger)

	// Initialize S3
	s3Service := storage.New(environment, logger)

	// Initialize RabbitMQ
	rabbitMQ, _ := rabbitmq.New(logger, environment)

	// Initialize Redis
	redisClient := redis.New(logger, environment)

	middleware.InitFromEnv(environment, logger)
	r := routes.BuildRouter(ctx, logger, environment)

	routes.SetupSwagger(r, environment, logger)
	// jwtService := middleware.GetJWTService()

	v1Group := r.Group("/api/v1")

	// Initialize Auth Modules
	auth.New(v1Group, database, logger, environment, s3Service)

	// Initialize Integrations Module (coordina core, WhatsApp, Shopify, etc.)
	integrations.New(v1Group, database, logger, environment, rabbitMQ)

	// Initialize Order Module
	modules.New(v1Group, database, logger, environment, rabbitMQ, redisClient)

	LogStartupInfo(ctx, logger, environment)

	port := environment.Get("HTTP_PORT")

	addr := fmt.Sprintf(":%s", port)
	return r.Run(addr)
}
