package main

import (
	"context"
	"os"

	"github.com/secamc93/probability/back/migration/internal/infra/repository"
	"github.com/secamc93/probability/back/migration/internal/seeds"
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

	// 6. Run Seeds
	logger.Info().Msg("Starting database seeding...")
	gormDB := database.Conn(context.Background())

	// Seed payment methods
	if err := seeds.SeedPaymentMethods(gormDB); err != nil {
		logger.Error(context.Background()).Err(err).Msg("Failed to seed payment methods")
	} else {
		logger.Info().Msg("Payment methods seeded successfully")
	}

	// Seed Shopify mappings
	if err := seeds.SeedShopifyMappings(gormDB); err != nil {
		logger.Error(context.Background()).Err(err).Msg("Failed to seed Shopify mappings")
	} else {
		logger.Info().Msg("Shopify mappings seeded successfully")
	}

	// Seed WhatsApp mappings
	if err := seeds.SeedWhatsAppMappings(gormDB); err != nil {
		logger.Error(context.Background()).Err(err).Msg("Failed to seed WhatsApp mappings")
	} else {
		logger.Info().Msg("WhatsApp mappings seeded successfully")
	}

	logger.Info().Msg("Database seeding completed successfully")
}
