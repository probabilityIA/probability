package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const (
	boldIntegrationCode      = "bold_pay"
	boldIntegrationName      = "Bold"
	boldCategoryCode         = "payment"
	boldBaseURL              = "https://integrations.api.bold.co"
	boldBaseURLTest          = "https://integrations.api.bold.co"
	boldImageURL             = "integrations/bold.png"
	boldDescription          = "Pasarela de pagos colombiana: tarjetas, PSE, Nequi, Boton Bancolombia."
	boldCredentialsSchemaRaw = `{
  "type": "object",
  "properties": {
    "api_key": {
      "type": "string",
      "title": "Identity Key (API Key)",
      "description": "Llave de identidad publica de Bold. Se encuentra en panel de Comercios -> Integraciones.",
      "required": true,
      "order": 1,
      "placeholder": "Ingresa tu Identity Key",
      "error_message": "La Identity Key es requerida"
    },
    "secret_key": {
      "type": "string",
      "title": "Secret Key",
      "description": "Llave secreta de Bold. Usada para firmar requests y validar webhooks.",
      "required": true,
      "order": 2,
      "placeholder": "Ingresa tu Secret Key",
      "error_message": "La Secret Key es requerida",
      "format": "password"
    },
    "environment": {
      "type": "string",
      "title": "Ambiente",
      "description": "sandbox para pruebas, production para clientes reales.",
      "required": false,
      "order": 3,
      "enum": ["sandbox", "production"],
      "default": "sandbox"
    }
  },
  "required": ["api_key", "secret_key"]
}`
	boldConfigSchemaRaw = `{
  "type": "object",
  "properties": {
    "webhook_url": {
      "type": "string",
      "title": "Webhook URL",
      "description": "URL publica donde Bold envia notificaciones de pago. Se registra en panel de Comercios.",
      "required": false,
      "order": 1
    }
  }
}`
	boldSetupInstructions = `# Configuracion de Bold

## Paso 1: Obtener credenciales
1. Ingresa al panel de Comercios de Bold (https://comercio.bold.co).
2. Ve a "Integraciones" -> "API y desarrollo".
3. Copia tu **Identity Key** (api_key) y tu **Secret Key**.

## Paso 2: Configurar integracion
1. Pega tu Identity Key y Secret Key en el formulario.
2. Selecciona ambiente: ` + "`sandbox`" + ` para pruebas o ` + "`production`" + ` para uso real.
3. Guarda. La conexion se prueba creando un link de pago de $1000 que puedes anular.

## Paso 3: Webhook (opcional pero recomendado)
1. Registra la URL ` + "`https://<tu-dominio>/webhooks/bold`" + ` en el panel de Bold.
2. Bold enviara eventos ` + "`SALE_APPROVED`" + `, ` + "`SALE_REJECTED`" + `, ` + "`VOID_APPROVED`" + `, ` + "`VOID_REJECTED`" + `.
3. La firma se valida con tu Secret Key (HMAC-SHA256).

## Notas
- Las credenciales se guardan encriptadas (AES-256-GCM).
- Solo super admin puede ver las credenciales descifradas.
- Sandbox y production usan la misma URL base; lo unico que cambia es la api_key.`
)

func (r *Repository) migrateBoldPay(ctx context.Context) error {
	if err := r.db.Conn(ctx).AutoMigrate(&models.BoldWebhookEvent{}, &models.CredentialRevealAudit{}); err != nil {
		return fmt.Errorf("automigrate bold pay tables: %w", err)
	}

	if err := r.db.Conn(ctx).Exec(`SELECT setval('integration_types_id_seq', GREATEST((SELECT COALESCE(MAX(id), 0) FROM integration_types), 1))`).Error; err != nil {
		return fmt.Errorf("resync integration_types sequence: %w", err)
	}

	var category models.IntegrationCategory
	err := r.db.Conn(ctx).Where("code = ?", boldCategoryCode).First(&category).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("integration_category 'payment' not found - run update_integration_categories_es.sql first")
		}
		return fmt.Errorf("query payment category: %w", err)
	}
	categoryID := category.ID

	credentialsSchema := datatypes.JSON([]byte(boldCredentialsSchemaRaw))
	configSchema := datatypes.JSON([]byte(boldConfigSchemaRaw))

	var existing models.IntegrationType
	res := r.db.Conn(ctx).Where("code IN ?", []string{boldIntegrationCode, "bold"}).
		Order("CASE code WHEN 'bold_pay' THEN 0 ELSE 1 END").
		First(&existing)
	if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		return fmt.Errorf("query bold integration_type: %w", res.Error)
	}

	if res.RowsAffected == 0 {
		row := models.IntegrationType{
			Name:              boldIntegrationName,
			Code:              boldIntegrationCode,
			Description:       boldDescription,
			Icon:              "credit-card",
			ImageURL:          boldImageURL,
			IsActive:          true,
			InDevelopment:     false,
			CategoryID:        &categoryID,
			ConfigSchema:      configSchema,
			CredentialsSchema: credentialsSchema,
			SetupInstructions: boldSetupInstructions,
			BaseURL:           boldBaseURL,
			BaseURLTest:       boldBaseURLTest,
		}
		if err := r.db.Conn(ctx).Create(&row).Error; err != nil {
			return fmt.Errorf("create bold integration_type: %w", err)
		}
		return nil
	}

	updates := map[string]any{
		"code":               boldIntegrationCode,
		"name":               boldIntegrationName,
		"description":        boldDescription,
		"icon":               "credit-card",
		"image_url":          boldImageURL,
		"category_id":        categoryID,
		"config_schema":      configSchema,
		"credentials_schema": credentialsSchema,
		"setup_instructions": boldSetupInstructions,
		"base_url":           boldBaseURL,
		"base_url_test":      boldBaseURLTest,
		"is_active":          true,
		"in_development":     false,
	}
	if err := r.db.Conn(ctx).Model(&existing).Updates(updates).Error; err != nil {
		return fmt.Errorf("update bold integration_type: %w", err)
	}
	return nil
}
