package main

import (
	"context"

	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

func main() {
	ctx := context.Background()
	logger := log.New()
	environment := env.New(logger)
	database := db.New(logger, environment)

	logger.Info(ctx).Msg("Starting Order Cleanup Process...")

	// Execute Truncate/Delete commands in order of dependencies
	tables := []string{
		"order_items",
		"addresses",
		"payments",
		"shipments",
		"order_channel_metadata",
		"orders",
		"client", // Clean clients too as per user request
	}

	for _, table := range tables {
		logger.Info(ctx).Str("table", table).Msg("Deleting records from table")
		if err := database.Conn(ctx).Exec("TRUNCATE TABLE " + table + " RESTART IDENTITY CASCADE").Error; err != nil {
			// Fallback to DELETE if TRUNCATE fails (e.g. permissions)
			logger.Warn(ctx).Err(err).Msg("TRUNCATE failed, attempting DELETE FROM")
			if err := database.Conn(ctx).Exec("DELETE FROM " + table).Error; err != nil {
				logger.Fatal(ctx).Err(err).Str("table", table).Msg("Failed to clean table")
			}
		}
	}

	logger.Info(ctx).Msg("✅ Successfully cleaned all order data!")
}
