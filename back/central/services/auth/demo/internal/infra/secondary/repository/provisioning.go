package repository

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"

	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/datatypes"
)

func parseEncryptionKey(key string) []byte {
	if decoded, err := base64.StdEncoding.DecodeString(key); err == nil && len(decoded) == 32 {
		return decoded
	}
	return []byte(key)
}

func (r *Repository) encryptCredentials(creds map[string]any) (datatypes.JSON, error) {
	if len(creds) == 0 {
		return datatypes.JSON([]byte("{}")), nil
	}
	plain, err := json.Marshal(creds)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(r.encKey)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	ciphertext := gcm.Seal(nonce, nonce, plain, nil)
	encoded := base64.StdEncoding.EncodeToString(ciphertext)
	wrapper, err := json.Marshal(map[string]string{"encrypted": encoded})
	if err != nil {
		return nil, err
	}
	return datatypes.JSON(wrapper), nil
}

func (r *Repository) ProvisionDemoIntegrations(ctx context.Context, businessID, userID uint) error {
	db := r.database.Conn(ctx)

	var existingWh models.Warehouse
	if err := db.Where("business_id = ?", businessID).First(&existingWh).Error; err != nil {
		warehouse := models.Warehouse{
			BusinessID:    businessID,
			Name:          "Bodega Principal",
			Code:          "BOD-DEMO",
			Address:       "Calle 100 #15-20",
			City:          "Bogota",
			State:         "Cundinamarca",
			Country:       "CO",
			ZipCode:       "110111",
			Phone:         "3001234567",
			ContactName:   "Demo",
			IsActive:      true,
			IsDefault:     true,
			IsFulfillment: true,
			StructureType: "simple",
			Company:       "Demo",
			FirstName:     "Demo",
			LastName:      "Bodega",
			CityDaneCode:  "11001000",
			PostalCode:    "110111",
			Street:        "Calle 100 #15-20",
		}
		if err := db.Create(&warehouse).Error; err != nil {
			return fmt.Errorf("error creando bodega demo: %w", err)
		}
	}

	platformCreds, _ := r.encryptCredentials(nil)
	siigoCreds, err := r.encryptCredentials(map[string]any{
		"username":   "demo@siigo.com",
		"access_key": "demo",
		"account_id": "demo",
		"partner_id": "demo",
	})
	if err != nil {
		return fmt.Errorf("error cifrando creds siigo: %w", err)
	}
	envioclickCreds, err := r.encryptCredentials(map[string]any{"api_key": "demo"})
	if err != nil {
		return fmt.Errorf("error cifrando creds envioclick: %w", err)
	}

	integrations := []models.Integration{
		{
			Name:              "Probability",
			Code:              fmt.Sprintf("platform_demo_%d", businessID),
			Category:          "platform",
			IntegrationTypeID: 6,
			BusinessID:        &businessID,
			IsActive:          true,
			IsDefault:         true,
			IsTesting:         false,
			Config:            datatypes.JSON([]byte("{}")),
			Credentials:       platformCreds,
			Description:       "Integracion interna de la plataforma (ordenes manuales) - demo",
		},
		{
			Name:              "Siigo (Demo)",
			Code:              fmt.Sprintf("siigo_test_%d", businessID),
			Category:          "invoicing",
			IntegrationTypeID: 8,
			BusinessID:        &businessID,
			IsActive:          true,
			IsDefault:         true,
			IsTesting:         true,
			Config:            datatypes.JSON([]byte("{}")),
			Credentials:       siigoCreds,
			Description:       "Facturacion electronica en modo prueba (mock) - demo",
		},
		{
			Name:              "EnvioClick (Demo)",
			Code:              fmt.Sprintf("envioclick_test_%d", businessID),
			Category:          "shipping",
			IntegrationTypeID: 12,
			BusinessID:        &businessID,
			IsActive:          true,
			IsDefault:         true,
			IsTesting:         true,
			Config:            datatypes.JSON([]byte(`{"auto_generate_guide_enabled":true}`)),
			Credentials:       envioclickCreds,
			Description:       "Generacion de guias en modo prueba (mock) - demo",
		},
		{
			Name:              "Inventario",
			Code:              fmt.Sprintf("inventory_business_%d", businessID),
			Category:          "internal",
			IntegrationTypeID: 32,
			BusinessID:        &businessID,
			IsActive:          true,
			IsDefault:         true,
			IsTesting:         false,
			Config:            datatypes.JSON([]byte("{}")),
			Credentials:       platformCreds,
			Description:       "Modulo interno de inventario - demo",
		},
	}

	var siigoIntegrationID uint
	for i := range integrations {
		integrations[i].CreatedByID = userID
		var existing models.Integration
		err := db.Where("code = ?", integrations[i].Code).First(&existing).Error
		if err == nil {
			if err := db.Model(&models.Integration{}).Where("id = ?", existing.ID).Updates(map[string]any{
				"credentials": integrations[i].Credentials,
				"config":      integrations[i].Config,
				"is_testing":  integrations[i].IsTesting,
				"is_active":   true,
			}).Error; err != nil {
				return fmt.Errorf("error actualizando integracion %s: %w", integrations[i].Code, err)
			}
			if integrations[i].IntegrationTypeID == 8 {
				siigoIntegrationID = existing.ID
			}
			continue
		}
		if err := db.Create(&integrations[i]).Error; err != nil {
			return fmt.Errorf("error creando integracion %s: %w", integrations[i].Code, err)
		}
		if integrations[i].IntegrationTypeID == 8 {
			siigoIntegrationID = integrations[i].ID
		}
	}

	if siigoIntegrationID > 0 {
		var existingCfg models.InvoicingConfig
		err := db.Where("business_id = ?", businessID).First(&existingCfg).Error
		if err != nil {
			cfg := models.InvoicingConfig{
				BusinessID:             businessID,
				InvoicingIntegrationID: &siigoIntegrationID,
				Enabled:                true,
				AutoInvoice:            true,
				CreatedByID:            userID,
			}
			if err := db.Create(&cfg).Error; err != nil {
				return fmt.Errorf("error creando invoicing_config: %w", err)
			}
		}
	}

	r.logger.Info().Uint("business_id", businessID).Msg("Integraciones demo (test) aprovisionadas")
	return nil
}
