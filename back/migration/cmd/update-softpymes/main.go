package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/secamc93/probability/back/migration/shared/db"
	"github.com/secamc93/probability/back/migration/shared/env"
	applog "github.com/secamc93/probability/back/migration/shared/log"
)

func main() {
	// Cargar .env
	if err := godotenv.Load("../../../.env"); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Inicializar logger
	logger := applog.New()

	// Inicializar configuración
	cfg := env.New(logger)

	// Conectar a base de datos
	database := db.New(logger, cfg)

	ctx := context.Background()

	// Ejecutar UPDATE
	sql := `UPDATE integration_types
SET
    credentials_schema = '{"type": "object", "properties": {"api_key": {"type": "string", "title": "API Key", "description": "Clave de API proporcionada por Softpymes", "required": true, "order": 1, "placeholder": "Ingresa tu API Key de Softpymes", "error_message": "La API Key es requerida"}, "api_secret": {"type": "string", "title": "API Secret", "description": "Secreto de API proporcionado por Softpymes", "required": true, "order": 2, "placeholder": "Ingresa tu API Secret de Softpymes", "error_message": "El API Secret es requerido", "format": "password"}}, "required": ["api_key", "api_secret"]}'::jsonb,
    config_schema = '{"type": "object", "properties": {}}'::jsonb,
    setup_instructions = 'Configura tu integración de Softpymes con tu API Key y API Secret'
WHERE code = 'softpymes'`

	result := database.Conn(ctx).Exec(sql)
	if result.Error != nil {
		log.Fatal("Error executing UPDATE:", result.Error)
	}

	fmt.Println("✅ Softpymes schemas updated successfully!")
	fmt.Println("📊 Rows affected:", result.RowsAffected)

	// Verificar el cambio
	var count int64
	database.Conn(ctx).Raw("SELECT COUNT(*) FROM integration_types WHERE code = 'softpymes' AND credentials_schema IS NOT NULL").Scan(&count)
	fmt.Println("✅ Verification: integration_types with credentials_schema:", count)

	os.Exit(0)
}
