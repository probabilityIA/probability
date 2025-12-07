package server

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/cmd/internal/routes"
	business "github.com/secamc93/probability/back/central/services/auth/bussines"
	"github.com/secamc93/probability/back/central/services/auth/login"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/auth/permissions"
	"github.com/secamc93/probability/back/central/services/auth/roles"
	"github.com/secamc93/probability/back/central/services/auth/users"
	"github.com/secamc93/probability/back/central/services/integrations"
	"github.com/secamc93/probability/back/central/services/modules"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/email"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	sharedQueue "github.com/secamc93/probability/back/central/shared/queue"
)

func Init(ctx context.Context) error {
	logger := log.New()
	environment := env.New(logger)

	database := db.New(logger, environment)
	// s3 := storage.New(environment, logger)
	_ = email.New(environment, logger)

	// Initialize RabbitMQ (opcional - si falla, se registra warning y continúa)
	var rabbitMQ sharedQueue.IQueue
	rabbitMQInstance, err := sharedQueue.New(logger, environment)
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
	login.New(v1Group, database, logger, environment)
	permissions.New(v1Group, database, logger)
	roles.New(v1Group, database, logger)
	users.New(v1Group, database, logger, environment, nil)
	business.New(v1Group, database, logger, environment, nil)

	// Initialize Integrations Module (coordina core, WhatsApp, Shopify, etc.)
	integrations.New(v1Group, database, logger, environment)

	// Initialize Order Module
	modules.New(v1Group, database, logger, environment, rabbitMQ)

	LogStartupInfo(ctx, logger, environment)

	port := environment.Get("HTTP_PORT")

	addr := fmt.Sprintf(":%s", port)
	return r.Run(addr)
}
