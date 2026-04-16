package repository

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/pay/nequi/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
	"gorm.io/gorm"
)

// IntegrationRepository obtiene credenciales de Nequi desde integration_types
type IntegrationRepository struct {
	db            db.IDatabase
	log           log.ILogger
	encryptionKey []byte
}

// New crea una nueva instancia del repositorio de integración Nequi
func New(database db.IDatabase, logger log.ILogger, encryptionKey string) ports.IIntegrationRepository {
	key := []byte(encryptionKey)
	// Ajustar key a 32 bytes (AES-256)
	if len(key) > 32 {
		key = key[:32]
	} else {
		for len(key) < 32 {
			key = append(key, 0)
		}
	}
	return &IntegrationRepository{
		db:            database,
		log:           logger.WithModule("nequi.integration_repository"),
		encryptionKey: key,
	}
}

// GetNequiConfig obtiene las credenciales de Nequi desde platform_credentials_encrypted
func (r *IntegrationRepository) GetNequiConfig(ctx context.Context) (*ports.NequiConfig, error) {
	type IntegrationTypeRow struct {
		ID                           uint   `json:"id"`
		PlatformCredentialsEncrypted []byte `json:"platform_credentials_encrypted"`
	}

	var row IntegrationTypeRow
	err := r.db.Conn(ctx).
		Table("integration_types").
		Select("id, platform_credentials_encrypted").
		Where("code = ?", "nequi_pay").
		Where("deleted_at IS NULL").
		First(&row).Error

	if err == gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("integration type 'nequi_pay' not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error querying nequi integration type: %w", err)
	}

	if len(row.PlatformCredentialsEncrypted) == 0 {
		return nil, fmt.Errorf("nequi platform credentials not configured")
	}

	credentials, err := r.decryptCredentials(ctx, row.PlatformCredentialsEncrypted)
	if err != nil {
		return nil, fmt.Errorf("error decrypting nequi credentials: %w", err)
	}

	apiKey, _ := credentials["api_key"].(string)
	environment, _ := credentials["environment"].(string)
	phoneCode, _ := credentials["phone_code"].(string)

	if apiKey == "" {
		return nil, fmt.Errorf("api_key not found in nequi credentials")
	}
	if environment == "" {
		environment = "sandbox"
	}
	if phoneCode == "" {
		phoneCode = "NIT_1"
	}

	r.log.Info(ctx).
		Uint("integration_type_id", row.ID).
		Str("environment", environment).
		Msg("Nequi config retrieved")

	return &ports.NequiConfig{
		APIKey:      apiKey,
		Environment: environment,
		PhoneCode:   phoneCode,
	}, nil
}

func (r *IntegrationRepository) decryptCredentials(ctx context.Context, ciphertext []byte) (map[string]interface{}, error) {
	block, err := aes.NewCipher(r.encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("error creating cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("error creating GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("error decrypting: %w", err)
	}

	var credentials map[string]interface{}
	if err := json.Unmarshal(plaintext, &credentials); err != nil {
		return nil, fmt.Errorf("error deserializing credentials: %w", err)
	}

	return credentials, nil
}
