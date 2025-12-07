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
)

func Init(ctx context.Context) error {
	logger := log.New()
	environment := env.New(logger)

	database := db.New(logger, environment)
	// s3 := storage.New(environment, logger)
	_ = email.New(environment, logger)

	// Initialize RabbitMQ (opcional - si falla, se registra warning y continúa)
	var rabbitMQ rabbitmq.IQueue
	rabbitMQInstance, err := rabbitmq.New(logger, environment)
	if err != nil {
		logger.Warn().
			Err(err).
			Msg("Failed to connect to RabbitMQ, queue features will not be available")
	} else {
		rabbitMQ = rabbitMQInstance
		logger.Info().Msg("RabbitMQ connected successfully")
	}

	middleware.InitFromEnv(environment, logger)
	r := routes.BuildRouter(ctx, logger, environment)

	routes.SetupSwagger(r, environment, logger)
	// jwtService := middleware.GetJWTService()

	v1Group := r.Group("/api/v1")

	// Initialize Auth Modules
	auth.New(v1Group, database, logger, environment)

	// Initialize Integrations Module (coordina core, WhatsApp, Shopify, etc.)
	integrations.New(v1Group, database, logger, environment, rabbitMQ)

	// Initialize Order Module
	modules.New(v1Group, database, logger, environment, rabbitMQ)

	LogStartupInfo(ctx, logger, environment)

	port := environment.Get("HTTP_PORT")

	addr := fmt.Sprintf(":%s", port)
	return r.Run(addr)
}
