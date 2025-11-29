package server

import (
	"central/cmd/internal/routes"
	"central/services/auth"
	"central/services/auth/business"
	"central/services/auth/middleware"
	"central/services/horizontalproperty"
	"central/services/restaurants/customer"
	"central/services/restaurants/reserve"
	"central/services/restaurants/rooms"
	"central/services/restaurants/tables"
	"central/shared/db"
	"central/shared/email"
	"central/shared/env"
	"central/shared/log"
	"central/shared/storage"
	"context"
	"fmt"
)

func Init(ctx context.Context) error {
	logger := log.New()
	environment := env.New(logger)

	database := db.New(logger, environment)
	s3 := storage.New(environment, logger)
	email := email.New(environment, logger)

	middleware.InitFromEnv(environment, logger)
	r := routes.BuildRouter(ctx, logger, environment)

	routes.SetupSwagger(r, environment, logger)
	jwtService := middleware.GetJWTService()

	v1Group := r.Group("/api/v1")

	auth.New(database, environment, logger, s3, v1Group, jwtService)
	customer.New(database, environment, logger, v1Group)
	business.New(database, environment, logger, s3, v1Group)
	horizontalproperty.New(database, logger, s3, environment, v1Group)
	reserve.New(database, environment, logger, email, v1Group)
	rooms.New(database, environment, logger, v1Group)
	tables.New(database, environment, logger, v1Group)

	LogStartupInfo(ctx, logger, environment)

	port := environment.Get("HTTP_PORT")

	addr := fmt.Sprintf(":%s", port)
	return r.Run(addr)
}
