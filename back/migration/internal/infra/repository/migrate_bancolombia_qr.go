package repository

import (
	"context"
	"fmt"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"github.com/secamc93/probability/back/migration/shared/models"
)

const (
	bancolombiaTypeCode     = "bancolombia_qr"
	bancolombiaTypeName     = "Bancolombia QR"
	bancolombiaCategoryCode = "payment"
	bancolombiaBaseURL      = "https://api.bancolombia.com"
	bancolombiaBaseURLTest  = "https://sb-api.bancolombia.com"
	bancolombiaImageURL     = "integration-types/bancolombia.png"
	bancolombiaDescription  = "Recibe pagos con codigo QR Bancolombia: genera un QR por transaccion y recibe la confirmacion de la transferencia via webhook."
	bancolombiaCredsSchema  = `{
  "type": "object",
  "properties": {
    "client_id": {
      "type": "string",
      "label": "Client ID (produccion)",
      "order": 1,
      "required": false,
      "input_type": "text",
      "description": "Client ID de la app registrada en el API Market de Bancolombia (ambiente productivo)",
      "help_link": "https://developer.bancolombia.com/"
    },
    "client_secret": {
      "type": "string",
      "label": "Client Secret (produccion)",
      "order": 2,
      "required": false,
      "input_type": "password",
      "description": "Client Secret de la app en el API Market (ambiente productivo)"
    },
    "webhook_secret": {
      "type": "string",
      "label": "Webhook Secret (produccion)",
      "order": 3,
      "required": false,
      "input_type": "password",
      "description": "Secreto para validar la firma HMAC de las notificaciones webhook. Si se deja vacio, el webhook se acepta sin validar firma."
    },
    "test_client_id": {
      "type": "string",
      "label": "Client ID (sandbox)",
      "order": 4,
      "required": false,
      "input_type": "text",
      "description": "Client ID de la app en el sandbox del API Market"
    },
    "test_client_secret": {
      "type": "string",
      "label": "Client Secret (sandbox)",
      "order": 5,
      "required": false,
      "input_type": "password",
      "description": "Client Secret de la app en el sandbox"
    },
    "test_webhook_secret": {
      "type": "string",
      "label": "Webhook Secret (sandbox)",
      "order": 6,
      "required": false,
      "input_type": "password",
      "description": "Secreto para validar la firma del webhook en sandbox"
    },
    "merchant_id": {
      "type": "string",
      "label": "Merchant ID",
      "order": 7,
      "required": false,
      "input_type": "text",
      "description": "Identificador del comercio asignado por Bancolombia para la generacion de QR"
    },
    "merchant_name": {
      "type": "string",
      "label": "Nombre del comercio",
      "order": 8,
      "required": false,
      "input_type": "text",
      "description": "Nombre del comercio que aparece en el QR"
    },
    "scope": {
      "type": "string",
      "label": "OAuth Scope",
      "order": 9,
      "required": false,
      "input_type": "text",
      "description": "Scope OAuth del producto QR. Default: QR-Management:write:app. Confirmar en la seccion de seguridad del API en el API Market."
    },
    "token_path": {
      "type": "string",
      "label": "Token Path",
      "order": 10,
      "required": false,
      "input_type": "text",
      "description": "Path del endpoint de token OAuth. Default: /security/oauth-provider/oauth2/token"
    },
    "qr_generate_path": {
      "type": "string",
      "label": "QR Generate Path",
      "order": 11,
      "required": false,
      "input_type": "text",
      "description": "Path del endpoint de generacion de QR. Default: /v1/qr-management/codes. Confirmar contra la especificacion del API Market."
    }
  }
}`
	bancolombiaConfigSchema = `{
  "type": "object",
  "properties": {}
}`
	bancolombiaSetupInstructions = `Pasos para conectar pagos con QR Bancolombia:

1. Crea una cuenta en el portal de desarrolladores de Bancolombia
   (https://developer.bancolombia.com/) y registra una aplicacion.
2. Suscribe la aplicacion al producto de APIs "QR Code" (sandbox para pruebas).
3. Copia el Client ID y Client Secret de la app y registralos como credenciales
   de plataforma de esta integracion (campos de sandbox para pruebas).
4. Configura el Merchant ID y el nombre del comercio entregados por Bancolombia.
5. Registra la URL del webhook de notificacion de transferencias en el portal:
   produccion /api/v1/webhooks/bancolombia y sandbox /api/v1/webhooks/bancolombia/test.
6. Si Bancolombia entrega un secreto de firma para las notificaciones,
   registralo en Webhook Secret para validar la firma HMAC de cada evento.

Notas:
- La autenticacion usa OAuth2 client_credentials contra el API Market
  (token valido por 20 minutos, se renueva automaticamente).
- Los paths del API de QR se pueden ajustar en las credenciales de plataforma
  (token_path, qr_generate_path, scope) sin recompilar, por si la
  especificacion del API Market difiere de los valores por defecto.`
)

func (r *Repository) migrateBancolombiaQR(ctx context.Context) error {
	if err := r.db.Conn(ctx).Exec(
		`SELECT setval('integration_types_id_seq', GREATEST((SELECT COALESCE(MAX(id), 0) FROM integration_types), 1))`,
	).Error; err != nil {
		return fmt.Errorf("migrateBancolombiaQR: resync secuencia: %w", err)
	}

	var category models.IntegrationCategory
	if err := r.db.Conn(ctx).Where("code = ?", bancolombiaCategoryCode).First(&category).Error; err != nil {
		return fmt.Errorf("migrateBancolombiaQR: categoria payment no encontrada: %w", err)
	}
	categoryID := category.ID

	var existing models.IntegrationType
	err := r.db.Conn(ctx).Where("code = ?", bancolombiaTypeCode).First(&existing).Error

	if err == nil {
		updates := map[string]interface{}{
			"name":               bancolombiaTypeName,
			"description":        bancolombiaDescription,
			"image_url":          bancolombiaImageURL,
			"category_id":        categoryID,
			"base_url":           bancolombiaBaseURL,
			"base_url_test":      bancolombiaBaseURLTest,
			"config_schema":      datatypes.JSON(bancolombiaConfigSchema),
			"credentials_schema": datatypes.JSON(bancolombiaCredsSchema),
			"setup_instructions": bancolombiaSetupInstructions,
		}
		if uerr := r.db.Conn(ctx).Model(&models.IntegrationType{}).
			Where("id = ?", existing.ID).
			Updates(updates).Error; uerr != nil {
			return fmt.Errorf("migrateBancolombiaQR: actualizando tipo: %w", uerr)
		}
		return nil
	}

	if err != gorm.ErrRecordNotFound {
		return fmt.Errorf("migrateBancolombiaQR: consultando tipo: %w", err)
	}

	integrationType := models.IntegrationType{
		Name:              bancolombiaTypeName,
		Code:              bancolombiaTypeCode,
		Description:       bancolombiaDescription,
		Icon:              "qr-code",
		ImageURL:          bancolombiaImageURL,
		IsActive:          true,
		InDevelopment:     true,
		CategoryID:        &categoryID,
		ConfigSchema:      datatypes.JSON(bancolombiaConfigSchema),
		CredentialsSchema: datatypes.JSON(bancolombiaCredsSchema),
		SetupInstructions: bancolombiaSetupInstructions,
		BaseURL:           bancolombiaBaseURL,
		BaseURLTest:       bancolombiaBaseURLTest,
	}

	if cerr := r.db.Conn(ctx).Create(&integrationType).Error; cerr != nil {
		return fmt.Errorf("migrateBancolombiaQR: creando tipo: %w", cerr)
	}

	return nil
}
