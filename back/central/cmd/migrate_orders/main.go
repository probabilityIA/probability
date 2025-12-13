package main

import (
	"context"
	"log"

	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/env"
	logger "github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/migration/shared/models"
)

func main() {
	ctx := context.Background()
	l := logger.New()
	environment := env.New(l)
	database := db.New(l, environment)

	log.Println("Starting migration for Order model...")

	// AutoMigrate will create the missing 'order_status_url' column
	if err := database.Conn(ctx).AutoMigrate(&models.Order{}); err != nil {
		l.Fatal(ctx).Err(err).Msg("Failed to auto-migrate Order model")
	}

	log.Println("Migration completed successfully.")
}
