package main

import (
	"context"
	"io/ioutil"

	"github.com/secamc93/probability/back/migration/shared/db"
	"github.com/secamc93/probability/back/migration/shared/env"
	"github.com/secamc93/probability/back/migration/shared/log"
)

type customConfig struct {
	env.IConfig
}

func (c *customConfig) Get(key string) string {
	if key == "DB_NAME" {
		return "probability"
	}
	return c.IConfig.Get(key)
}

func main() {
	// 1. Init Logger
	logger := log.New()

	// 2. Init Config
	baseCfg := env.New(logger)
	cfg := &customConfig{IConfig: baseCfg}

	// 3. Init DB
	database := db.New(logger, cfg)
	defer database.Close()

	// 4. List tables for debugging
	var tables []string
	database.Conn(context.Background()).Raw("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'").Scan(&tables)
	logger.Info().Interface("tables", tables).Msg("Existing tables in 'public' schema of 'probability' DB")

	if len(tables) == 0 {
		logger.Fatal(context.Background()).Msg("Still no tables found in 'probability' database. Check RDS configuration.")
	}

	// 5. Read SQL file
	sqlPath := "shared/sql/add_shipment_permissions_to_business_role.sql"
	sqlContent, err := ioutil.ReadFile(sqlPath)
	if err != nil {
		logger.Fatal(context.Background()).Err(err).Msg("Failed to read SQL file")
	}

	// 6. Execute SQL
	logger.Info().Msg("Executing SQL migration script...")

	// Usar GORM Exec para ejecutar el SQL
	if err := database.Conn(context.Background()).Exec(string(sqlContent)).Error; err != nil {
		logger.Fatal(context.Background()).Err(err).Msg("Failed to execute SQL")
	}

	logger.Info().Msg("SQL script executed successfully! Permissions granted to 'Administrador' role.")
}
