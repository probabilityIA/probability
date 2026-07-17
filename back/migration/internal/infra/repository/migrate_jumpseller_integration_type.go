package repository

import (
	"context"
	"fmt"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"github.com/secamc93/probability/back/migration/shared/models"
)

const (
	jumpsellerTypeID      = 33
	jumpsellerTypeCode    = "jumpseller"
	jumpsellerTypeName    = "Jumpseller"
	jumpsellerCategoryID  = 1
	jumpsellerBaseURL     = "https://api.jumpseller.com/v1"
	jumpsellerBaseURLTest = "http://back-testing:9097"
	jumpsellerImageURL    = "integration-types/1784241131_jumpseller.png"
	jumpsellerDescription = "Integracion con tiendas Jumpseller: importa ordenes en tiempo real, empuja stock y sincroniza estados."
	jumpsellerCredsSchema = `{
  "type": "object",
  "properties": {
    "api_key": {
      "type": "string",
      "label": "Login",
      "order": 1,
      "required": true,
      "input_type": "text",
      "description": "Login de la API de Jumpseller",
      "placeholder": "tu-login",
      "help_text": "Lo encuentras en Jumpseller: Cuenta > API. Es el campo Login.",
      "help_link": "https://jumpseller.com/support/api/",
      "error_message": "El Login es obligatorio"
    },
    "api_secret": {
      "type": "string",
      "label": "Auth Token",
      "order": 2,
      "required": true,
      "input_type": "password",
      "description": "Auth Token de la API de Jumpseller",
      "placeholder": "xxxxxxxxxxxxxxxxxxxxxxxx",
      "help_text": "Lo encuentras en Jumpseller: Cuenta > API. Es el campo Auth Token.",
      "help_link": "https://jumpseller.com/support/api/",
      "error_message": "El Auth Token es obligatorio"
    }
  }
}`
	jumpsellerConfigSchema = `{
  "type": "object",
  "properties": {
    "inventory_sync_enabled": {
      "type": "boolean",
      "label": "Sincronizar inventario hacia Jumpseller",
      "order": 1,
      "required": false,
      "input_type": "checkbox",
      "description": "Empuja el stock de Probability a los productos asociados de Jumpseller"
    },
    "status_sync_enabled": {
      "type": "boolean",
      "label": "Sincronizar estados hacia Jumpseller",
      "order": 2,
      "required": false,
      "input_type": "checkbox",
      "description": "Actualiza el estado y el tracking de la orden en Jumpseller"
    }
  }
}`
	jumpsellerSetupInstructions = `Pasos para conectar tu tienda Jumpseller:

1. Inicia sesion en el panel de administracion de tu tienda Jumpseller.
2. Ve a Cuenta > API (https://jumpseller.com/admin/account/api).
3. Activa el acceso a la API si aun no lo esta.
4. Copia el campo "Login" y pegalo en el campo Login de Probability.
5. Copia el campo "Auth Token" y pegalo en el campo Auth Token de Probability.
6. Guarda la integracion. Probability validara las credenciales contra tu tienda
   y registrara automaticamente los webhooks de ordenes (order_created, order_paid,
   order_shipped y order_canceled).

Modo de pruebas: al activar "modo pruebas" la integracion no habla con la API real
de Jumpseller sino con el simulador interno (back-testing), por lo que puedes usar
credenciales ficticias.`
)

func (r *Repository) migrateJumpsellerIntegrationType(ctx context.Context) error {
	var existing models.IntegrationType
	err := r.db.Conn(ctx).
		Where("id = ? OR code = ?", jumpsellerTypeID, jumpsellerTypeCode).
		First(&existing).Error

	categoryID := uint(jumpsellerCategoryID)

	if err == nil {
		updates := map[string]interface{}{
			"name":               jumpsellerTypeName,
			"code":               jumpsellerTypeCode,
			"description":        jumpsellerDescription,
			"image_url":          jumpsellerImageURL,
			"category_id":        categoryID,
			"base_url":           jumpsellerBaseURL,
			"base_url_test":      jumpsellerBaseURLTest,
			"config_schema":      datatypes.JSON(jumpsellerConfigSchema),
			"credentials_schema": datatypes.JSON(jumpsellerCredsSchema),
			"setup_instructions": jumpsellerSetupInstructions,
		}
		if uerr := r.db.Conn(ctx).Model(&models.IntegrationType{}).
			Where("id = ?", existing.ID).
			Updates(updates).Error; uerr != nil {
			return fmt.Errorf("migrateJumpsellerIntegrationType: actualizando tipo: %w", uerr)
		}
		return nil
	}

	if err != gorm.ErrRecordNotFound {
		return fmt.Errorf("migrateJumpsellerIntegrationType: consultando tipo: %w", err)
	}

	integrationType := models.IntegrationType{
		Model:             gorm.Model{ID: jumpsellerTypeID},
		Name:              jumpsellerTypeName,
		Code:              jumpsellerTypeCode,
		Description:       jumpsellerDescription,
		ImageURL:          jumpsellerImageURL,
		IsActive:          true,
		InDevelopment:     true,
		CategoryID:        &categoryID,
		ConfigSchema:      datatypes.JSON(jumpsellerConfigSchema),
		CredentialsSchema: datatypes.JSON(jumpsellerCredsSchema),
		SetupInstructions: jumpsellerSetupInstructions,
		BaseURL:           jumpsellerBaseURL,
		BaseURLTest:       jumpsellerBaseURLTest,
	}

	if cerr := r.db.Conn(ctx).Create(&integrationType).Error; cerr != nil {
		return fmt.Errorf("migrateJumpsellerIntegrationType: creando tipo: %w", cerr)
	}

	return nil
}
