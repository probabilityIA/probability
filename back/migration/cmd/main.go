package main

import (
	"context"
	"os"

	"github.com/secamc93/probability/back/migration/internal/infra/repository"
	"github.com/secamc93/probability/back/migration/shared/db"
	"github.com/secamc93/probability/back/migration/shared/env"
	"github.com/secamc93/probability/back/migration/shared/log"
)

func main() {
	// 1. Init Logger
	logger := log.New()

	// 2. Init Config
	cfg := env.New(logger)

	// 3. Init DB
	database := db.New(logger, cfg)
	defer database.Close()

	// 4. Init Repository
	repo := repository.New(database, cfg)

	// 5. Run Migrations
	logger.Info().Msg("Starting database migration...")
	if err := repo.Migrate(context.Background()); err != nil {
		logger.Fatal(context.Background()).Err(err).Msg("Migration failed")
		os.Exit(1)
	}
	logger.Info().Msg("Database migration completed successfully")

}
