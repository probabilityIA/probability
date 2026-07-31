package repository

import (
	"context"
	"fmt"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"github.com/secamc93/probability/back/migration/shared/models"
)

const (
	tiktokTypeID      = 35
	tiktokTypeCode    = "tiktok"
	tiktokTypeName    = "TikTok"
	tiktokCategoryID  = 1
	tiktokBaseURL     = "https://open.tiktokapis.com"
	tiktokBaseURLTest = "http://back-testing:9101"
	tiktokImageURL    = "integration-types/tiktok.png"
	tiktokDescription = "Integracion con TikTok (Research API): consulta y monitorea las metricas de tu TikTok Shop (rating, resenas, ventas) con snapshots historicos."
	tiktokCredsSchema = `{
  "type": "object",
  "properties": {
    "client_key": {
      "type": "string",
      "label": "Client Key",
      "order": 1,
      "required": true,
      "input_type": "text",
      "description": "Client Key de tu app en TikTok for Developers",
      "placeholder": "awxxxxxxxxxxxxxxxx",
      "help_text": "Lo encuentras en TikTok for Developers > Manage apps > tu app > Credentials.",
      "help_link": "https://developers.tiktok.com/",
      "error_message": "El Client Key es obligatorio"
    },
    "client_secret": {
      "type": "string",
      "label": "Client Secret",
      "order": 2,
      "required": true,
      "input_type": "password",
      "description": "Client Secret de tu app en TikTok for Developers",
      "placeholder": "xxxxxxxxxxxxxxxxxxxxxxxx",
      "help_text": "Lo encuentras en TikTok for Developers > Manage apps > tu app > Credentials.",
      "help_link": "https://developers.tiktok.com/",
      "error_message": "El Client Secret es obligatorio"
    }
  }
}`
	tiktokConfigSchema = `{
  "type": "object",
  "properties": {
    "shop_name": {
      "type": "string",
      "label": "Nombre de la tienda (TikTok Shop)",
      "order": 1,
      "required": true,
      "input_type": "text",
      "description": "Nombre exacto de tu TikTok Shop tal como aparece en la plataforma",
      "placeholder": "mitienda",
      "help_text": "Se usa para consultar las metricas de tu tienda via la Research API.",
      "error_message": "El nombre de la tienda es obligatorio"
    }
  }
}`
	tiktokSetupInstructions = `Pasos para conectar tu TikTok Shop:

1. Crea una app en TikTok for Developers (https://developers.tiktok.com/).
2. Solicita acceso a la Research API con el scope research.data.basic.
3. Una vez aprobada, ve a Manage apps > tu app > Credentials.
4. Copia el Client Key y pegalo en el campo Client Key de Probability.
5. Copia el Client Secret y pegalo en el campo Client Secret de Probability.
6. Escribe el nombre exacto de tu TikTok Shop en el campo de configuracion.
7. Guarda la integracion. Probability validara las credenciales obteniendo un
   token de cliente y consultando la informacion de tu tienda.

Modo de pruebas: al activar "modo pruebas" la integracion no habla con la API real
de TikTok sino con el simulador interno (back-testing), por lo que puedes usar
credenciales ficticias.`
)

func (r *Repository) migrateTiktokIntegrationType(ctx context.Context) error {
	if err := r.db.Conn(ctx).AutoMigrate(&models.TikTokShopSnapshot{}); err != nil {
		return fmt.Errorf("migrateTiktokIntegrationType: creando tabla de snapshots: %w", err)
	}

	var existing models.IntegrationType
	err := r.db.Conn(ctx).
		Where("id = ? OR code = ?", tiktokTypeID, tiktokTypeCode).
		First(&existing).Error

	categoryID := uint(tiktokCategoryID)

	if err == nil {
		updates := map[string]interface{}{
			"name":               tiktokTypeName,
			"code":               tiktokTypeCode,
			"description":        tiktokDescription,
			"image_url":          tiktokImageURL,
			"category_id":        categoryID,
			"base_url":           tiktokBaseURL,
			"base_url_test":      tiktokBaseURLTest,
			"config_schema":      datatypes.JSON(tiktokConfigSchema),
			"credentials_schema": datatypes.JSON(tiktokCredsSchema),
			"setup_instructions": tiktokSetupInstructions,
			"in_development":     false,
		}
		if uerr := r.db.Conn(ctx).Model(&models.IntegrationType{}).
			Where("id = ?", existing.ID).
			Updates(updates).Error; uerr != nil {
			return fmt.Errorf("migrateTiktokIntegrationType: actualizando tipo: %w", uerr)
		}
		return nil
	}

	if err != gorm.ErrRecordNotFound {
		return fmt.Errorf("migrateTiktokIntegrationType: consultando tipo: %w", err)
	}

	integrationType := models.IntegrationType{
		Model:             gorm.Model{ID: tiktokTypeID},
		Name:              tiktokTypeName,
		Code:              tiktokTypeCode,
		Description:       tiktokDescription,
		ImageURL:          tiktokImageURL,
		IsActive:          true,
		InDevelopment:     false,
		CategoryID:        &categoryID,
		ConfigSchema:      datatypes.JSON(tiktokConfigSchema),
		CredentialsSchema: datatypes.JSON(tiktokCredsSchema),
		SetupInstructions: tiktokSetupInstructions,
		BaseURL:           tiktokBaseURL,
		BaseURLTest:       tiktokBaseURLTest,
	}

	if cerr := r.db.Conn(ctx).Create(&integrationType).Error; cerr != nil {
		return fmt.Errorf("migrateTiktokIntegrationType: creando tipo: %w", cerr)
	}

	return nil
}
