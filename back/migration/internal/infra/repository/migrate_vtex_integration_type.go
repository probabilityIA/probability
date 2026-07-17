package repository

import (
	"context"
	"fmt"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"github.com/secamc93/probability/back/migration/shared/models"
)

const (
	vtexTypeID      = 16
	vtexTypeCode    = "vtex"
	vtexTypeName    = "VTEX"
	vtexCategoryID  = 1
	vtexBaseURL     = "https://{accountName}.vtexcommercestable.com.br"
	vtexImageURL    = "integration-types/1771905400_vtex.png"
	vtexDescription = "Integracion con tiendas VTEX: importa ordenes en tiempo real, empuja stock y sincroniza el catalogo."
	vtexCredsSchema = `{
  "type": "object",
  "properties": {
    "app_key": {
      "type": "string",
      "label": "App Key",
      "order": 1,
      "required": true,
      "input_type": "password",
      "description": "Clave de aplicacion de la API de VTEX",
      "placeholder": "vtexappkey-mitienda-XXXXXX",
      "help_text": "La encuentras en el Admin de VTEX: Configuracion de la cuenta > Claves de API.",
      "help_link": "https://developers.vtex.com/docs/guides/authentication",
      "error_message": "El App Key es obligatorio"
    },
    "app_token": {
      "type": "string",
      "label": "App Token",
      "order": 2,
      "required": true,
      "input_type": "password",
      "description": "Token de aplicacion de la API de VTEX",
      "placeholder": "xxxxxxxxxxxxxxxxxxxxxxxx",
      "help_text": "Solo se muestra una vez, al crear la clave en Configuracion de la cuenta > Claves de API.",
      "help_link": "https://developers.vtex.com/docs/guides/authentication",
      "error_message": "El App Token es obligatorio"
    }
  }
}`
	vtexConfigSchema = `{
  "type": "object",
  "properties": {
    "account_name": {
      "type": "string",
      "label": "Nombre de la cuenta VTEX",
      "order": 1,
      "required": true,
      "input_type": "text",
      "description": "El nombre de tu cuenta VTEX, tal como aparece en la URL del Admin",
      "placeholder": "mitienda",
      "help_text": "Si tu Admin es mitienda.myvtex.com, el nombre de la cuenta es mitienda."
    },
    "is_seller": {
      "type": "boolean",
      "label": "La cuenta es un seller",
      "order": 2,
      "required": false,
      "input_type": "checkbox",
      "description": "Activalo si vendes como seller en un marketplace VTEX y no como tienda propia"
    },
    "inventory_sync_enabled": {
      "type": "boolean",
      "label": "Sincronizar inventario hacia VTEX",
      "order": 3,
      "required": false,
      "input_type": "checkbox",
      "description": "Empuja el stock de Probability a los SKUs asociados de VTEX"
    },
    "status_sync_enabled": {
      "type": "boolean",
      "label": "Sincronizar estados hacia VTEX",
      "order": 4,
      "required": false,
      "input_type": "checkbox",
      "description": "Actualiza el estado y el tracking de la orden en VTEX"
    }
  }
}`
	vtexSetupInstructions = `Pasos para conectar tu tienda VTEX:

1. Ingresa al Admin de VTEX de tu tienda.
2. Ve a Configuracion de la cuenta > Claves de API.
3. Crea una nueva clave de aplicacion.
4. Copia el "App Key" y pegalo en el campo App Key de Probability.
5. Copia el "App Token" y pegalo en el campo App Token de Probability.
   Ojo: el App Token solo se muestra una vez, al momento de crear la clave.
6. Asigna a la clave los roles de OMS (ordenes), Catalogo y Logistica (inventario).
7. En Probability escribe el nombre de tu cuenta VTEX (si tu Admin es
   mitienda.myvtex.com, el nombre de la cuenta es mitienda).
8. Usa "Probar Conexion" para validar las credenciales y guarda la integracion.

Webhooks: VTEX admite una sola configuracion de hook por cuenta
(POST /api/orders/hook/config). Al registrarla se reemplaza cualquier hook
anterior de esa cuenta, incluido el de otra herramienta.`
)

func (r *Repository) migrateVtexIntegrationType(ctx context.Context) error {
	var existing models.IntegrationType
	err := r.db.Conn(ctx).
		Where("id = ? OR code = ?", vtexTypeID, vtexTypeCode).
		First(&existing).Error

	categoryID := uint(vtexCategoryID)

	if err == nil {
		updates := map[string]interface{}{
			"name":               vtexTypeName,
			"code":               vtexTypeCode,
			"description":        vtexDescription,
			"category_id":        categoryID,
			"base_url":           vtexBaseURL,
			"config_schema":      datatypes.JSON(vtexConfigSchema),
			"credentials_schema": datatypes.JSON(vtexCredsSchema),
			"setup_instructions": vtexSetupInstructions,
		}
		if uerr := r.db.Conn(ctx).Model(&models.IntegrationType{}).
			Where("id = ?", existing.ID).
			Updates(updates).Error; uerr != nil {
			return fmt.Errorf("migrateVtexIntegrationType: actualizando tipo: %w", uerr)
		}
		return nil
	}

	if err != gorm.ErrRecordNotFound {
		return fmt.Errorf("migrateVtexIntegrationType: consultando tipo: %w", err)
	}

	integrationType := models.IntegrationType{
		Model:             gorm.Model{ID: vtexTypeID},
		Name:              vtexTypeName,
		Code:              vtexTypeCode,
		Description:       vtexDescription,
		ImageURL:          vtexImageURL,
		IsActive:          true,
		InDevelopment:     true,
		CategoryID:        &categoryID,
		ConfigSchema:      datatypes.JSON(vtexConfigSchema),
		CredentialsSchema: datatypes.JSON(vtexCredsSchema),
		SetupInstructions: vtexSetupInstructions,
		BaseURL:           vtexBaseURL,
	}

	if cerr := r.db.Conn(ctx).Create(&integrationType).Error; cerr != nil {
		return fmt.Errorf("migrateVtexIntegrationType: creando tipo: %w", cerr)
	}

	return nil
}
