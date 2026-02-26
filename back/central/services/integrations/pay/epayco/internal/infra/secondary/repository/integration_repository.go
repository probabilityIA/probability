package repository

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"
	"fmt"

	epaycoErrors "github.com/secamc93/probability/back/central/services/integrations/pay/epayco/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/integrations/pay/epayco/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
	"gorm.io/gorm"
)

// IntegrationRepository obtiene credenciales de ePayco desde integration_types
type IntegrationRepository struct {
	db            db.IDatabase
	log           log.ILogger
	encryptionKey []byte
}

// New crea una nueva instancia del repositorio de integración ePayco
func New(database db.IDatabase, logger log.ILogger, encryptionKey string) ports.IIntegrationRepository {
	key := []byte(encryptionKey)
	if len(key) > 32 {
		key = key[:32]
	} else {
		for len(key) < 32 {
			key = append(key, 0)
		}
	}
	return &IntegrationRepository{
		db:            database,
		log:           logger.WithModule("epayco.integration_repository"),
		encryptionKey: key,
	}
}

// GetEPaycoConfig obtiene las credenciales de ePayco desde platform_credentials_encrypted
// Credenciales JSON esperadas: {"customer_id": "...", "key": "...", "environment": "test|production"}
func (r *IntegrationRepository) GetEPaycoConfig(ctx context.Context) (*ports.EPaycoConfig, error) {
	type IntegrationTypeRow struct {
		ID                           uint   `json:"id"`
		PlatformCredentialsEncrypted []byte `json:"platform_credentials_encrypted"`
	}

	var row IntegrationTypeRow
	err := r.db.Conn(ctx).
		Table("integration_types").
		Select("id, platform_credentials_encrypted").
		Where("code = ?", "epayco_pay").
		Where("deleted_at IS NULL").
		First(&row).Error

	if err == gorm.ErrRecordNotFound {
		return nil, epaycoErrors.ErrEPaycoConfigNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("error querying epayco integration type: %w", err)
	}

	if len(row.PlatformCredentialsEncrypted) == 0 {
		return nil, fmt.Errorf("epayco platform credentials not configured")
	}

	credentials, err := r.decryptCredentials(ctx, row.PlatformCredentialsEncrypted)
	if err != nil {
		return nil, fmt.Errorf("error decrypting epayco credentials: %w", err)
	}

	customerID, _ := credentials["customer_id"].(string)
	key, _ := credentials["key"].(string)
	environment, _ := credentials["environment"].(string)

	if customerID == "" || key == "" {
		return nil, epaycoErrors.ErrInvalidCredentials
	}
	if environment == "" {
		environment = "test"
	}

	r.log.Info(ctx).
		Uint("integration_type_id", row.ID).
		Str("environment", environment).
		Msg("ePayco config retrieved")

	return &ports.EPaycoConfig{
		CustomerID:  customerID,
		Key:         key,
		Environment: environment,
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
