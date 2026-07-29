package repository

import (
	"context"
	"fmt"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"github.com/secamc93/probability/back/migration/shared/models"
)

const (
	shipitTypeID      = 34
	shipitTypeCode    = "shipit"
	shipitTypeName    = "Shipit"
	shipitCategoryID  = 5
	shipitBaseURL     = "https://api.shipit.cl"
	shipitBaseURLTest = "http://back-testing:9100"
	shipitImageURL    = "integration-types/1785346702_shipit.png"
	shipitIcon        = "truck"
	shipitDescription = "Integracion con Shipit (Chile): genera envios, cotiza tarifas entre couriers y sincroniza el estado de las guias."
	shipitCredsSchema = `{
  "type": "object",
  "properties": {
    "email": {
      "type": "string",
      "label": "Email de la cuenta",
      "order": 1,
      "required": true,
      "input_type": "text",
      "description": "Email de la cuenta Shipit (header X-Shipit-Email)",
      "placeholder": "cuenta@dominio.cl",
      "help_text": "Es el email con el que inicias sesion en app.shipit.cl",
      "help_link": "https://app.shipit.cl/settings?tab=api",
      "error_message": "El email es obligatorio"
    },
    "access_token": {
      "type": "string",
      "label": "Access Token",
      "order": 2,
      "required": true,
      "input_type": "password",
      "description": "Token de acceso de la API de Shipit (header X-Shipit-Access-Token)",
      "placeholder": "xxxxxxxxxxxxxxxxxxxx",
      "help_text": "Lo encuentras en Shipit: Configuracion > API.",
      "help_link": "https://app.shipit.cl/settings?tab=api",
      "error_message": "El access token es obligatorio"
    }
  }
}`
	shipitConfigSchema = `{
  "type": "object",
  "properties": {
    "base_url": {
      "type": "string",
      "label": "URL base de la API",
      "order": 1,
      "required": false,
      "input_type": "text",
      "description": "URL base de la API de Shipit (no modificar salvo pruebas)",
      "placeholder": "https://api.shipit.cl"
    },
    "base_url_test": {
      "type": "string",
      "label": "URL base de pruebas",
      "order": 2,
      "required": false,
      "input_type": "text",
      "description": "Si se define, la integracion habla con este simulador en lugar de la API real"
    }
  }
}`
	shipitSetupInstructions = `Pasos para conectar tu cuenta Shipit:

1. Inicia sesion en https://app.shipit.cl
2. Ve a Configuracion > API (https://app.shipit.cl/settings?tab=api).
3. Copia el email de tu cuenta y pegalo en el campo Email de Probability.
4. Copia el Access Token y pegalo en el campo Access Token de Probability.
5. Guarda la integracion.

Notas:
- Shipit opera en Chile: las direcciones de destino deben incluir una comuna
  chilena valida (se resuelve automaticamente por el nombre de la ciudad).
- Para recibir actualizaciones automaticas de estado, configura el webhook de
  packages en Shipit apuntando a: {API_URL}/api/v1/webhooks/shipit
- Modo de pruebas: al activar "modo pruebas" la integracion habla con el
  simulador interno (back-testing) y ademas marca los envios como sandbox.`
)

func (r *Repository) migrateShipitIntegrationType(ctx context.Context) error {
	var existing models.IntegrationType
	err := r.db.Conn(ctx).
		Where("id = ? OR code = ?", shipitTypeID, shipitTypeCode).
		First(&existing).Error

	categoryID := uint(shipitCategoryID)

	if err == nil {
		updates := map[string]interface{}{
			"name":               shipitTypeName,
			"code":               shipitTypeCode,
			"description":        shipitDescription,
			"icon":               shipitIcon,
			"image_url":          shipitImageURL,
			"category_id":        categoryID,
			"base_url":           shipitBaseURL,
			"base_url_test":      shipitBaseURLTest,
			"config_schema":      datatypes.JSON(shipitConfigSchema),
			"credentials_schema": datatypes.JSON(shipitCredsSchema),
			"setup_instructions": shipitSetupInstructions,
		}
		if uerr := r.db.Conn(ctx).Model(&models.IntegrationType{}).
			Where("id = ?", existing.ID).
			Updates(updates).Error; uerr != nil {
			return fmt.Errorf("migrateShipitIntegrationType: actualizando tipo: %w", uerr)
		}
		return nil
	}

	if err != gorm.ErrRecordNotFound {
		return fmt.Errorf("migrateShipitIntegrationType: consultando tipo: %w", err)
	}

	integrationType := models.IntegrationType{
		Model:             gorm.Model{ID: shipitTypeID},
		Name:              shipitTypeName,
		Code:              shipitTypeCode,
		Description:       shipitDescription,
		Icon:              shipitIcon,
		ImageURL:          shipitImageURL,
		IsActive:          true,
		InDevelopment:     false,
		CategoryID:        &categoryID,
		ConfigSchema:      datatypes.JSON(shipitConfigSchema),
		CredentialsSchema: datatypes.JSON(shipitCredsSchema),
		SetupInstructions: shipitSetupInstructions,
		BaseURL:           shipitBaseURL,
		BaseURLTest:       shipitBaseURLTest,
	}

	if cerr := r.db.Conn(ctx).Create(&integrationType).Error; cerr != nil {
		return fmt.Errorf("migrateShipitIntegrationType: creando tipo: %w", cerr)
	}

	return nil
}
