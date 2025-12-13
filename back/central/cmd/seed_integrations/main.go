package main

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/datatypes"
)

func main() {
	ctx := context.Background()
	logger := log.New()
	environment := env.New(logger)
	database := db.New(logger, environment)

	// Seed Integration Types
	types := []models.IntegrationType{
		{
			Name:        "WhatsApp",
			Code:        "whatsapp",
			Category:    "internal",
			Description: "Integración con WhatsApp Business API",
			IsActive:    true,
			ConfigSchema: datatypes.JSON([]byte(`{
				"type": "object",
				"properties": {
					"phone_number_id": {"type": "string"},
					"webhook_url": {"type": "string"}
				},
				"required": ["phone_number_id"]
			}`)),
			CredentialsSchema: datatypes.JSON([]byte(`{
				"type": "object",
				"properties": {
					"access_token": {"type": "string"}
				},
				"required": ["access_token"]
			}`)),
		},
		{
			Name:        "Shopify",
			Code:        "shopify",
			Category:    "external",
			Description: "Integración con tiendas Shopify",
			IsActive:    true,
			ConfigSchema: datatypes.JSON([]byte(`{
				"type": "object",
				"properties": {
					"store_name": {"type": "string"},
					"api_version": {"type": "string"}
				},
				"required": ["store_name"]
			}`)),
			CredentialsSchema: datatypes.JSON([]byte(`{
				"type": "object",
				"properties": {
					"access_token": {"type": "string"},
					"api_key": {"type": "string"}
				},
				"required": ["access_token"]
			}`)),
		},
		{
			Name:              "Mercado Libre",
			Code:              "mercado_libre",
			Category:          "external",
			Description:       "Integración con Mercado Libre",
			IsActive:          true,
			ConfigSchema:      datatypes.JSON([]byte(`{}`)),
			CredentialsSchema: datatypes.JSON([]byte(`{}`)),
		},
	}

	for _, t := range types {
		var existing models.IntegrationType
		if err := database.Conn(ctx).Where("code = ?", t.Code).First(&existing).Error; err != nil {
			// Create new
			t.CreatedAt = time.Now()
			t.UpdatedAt = time.Now()
			if err := database.Conn(ctx).Create(&t).Error; err != nil {
				logger.Error(ctx).Err(err).Str("code", t.Code).Msg("Failed to create integration type")
			} else {
				logger.Info(ctx).Str("code", t.Code).Msg("Created integration type")
			}
		} else {
			// Update existing to ensure category is set
			existing.Category = t.Category
			existing.ConfigSchema = t.ConfigSchema
			existing.CredentialsSchema = t.CredentialsSchema
			if err := database.Conn(ctx).Save(&existing).Error; err != nil {
				logger.Error(ctx).Err(err).Str("code", t.Code).Msg("Failed to update integration type")
			} else {
				logger.Info(ctx).Str("code", t.Code).Msg("Updated integration type")
			}
		}
	}
}
