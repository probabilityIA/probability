package repository

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"
	"fmt"

	meliErrors "github.com/secamc93/probability/back/central/services/integrations/pay/melipago/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/integrations/pay/melipago/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
	"gorm.io/gorm"
)

// IntegrationRepository obtiene credenciales de MercadoPago desde integration_types
type IntegrationRepository struct {
	db            db.IDatabase
	log           log.ILogger
	encryptionKey []byte
}

// New crea una nueva instancia del repositorio de integración MercadoPago
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
		log:           logger.WithModule("melipago.integration_repository"),
		encryptionKey: key,
	}
}

// GetMeliPagoConfig obtiene las credenciales de MercadoPago desde platform_credentials_encrypted
// Credenciales JSON esperadas: {"access_token": "...", "environment": "sandbox|production"}
func (r *IntegrationRepository) GetMeliPagoConfig(ctx context.Context) (*ports.MeliPagoConfig, error) {
	type IntegrationTypeRow struct {
		ID                           uint   `json:"id"`
		PlatformCredentialsEncrypted []byte `json:"platform_credentials_encrypted"`
	}

	var row IntegrationTypeRow
	err := r.db.Conn(ctx).
		Table("integration_types").
		Select("id, platform_credentials_encrypted").
		Where("code = ?", "melipago_pay").
		Where("deleted_at IS NULL").
		First(&row).Error

	if err == gorm.ErrRecordNotFound {
		return nil, meliErrors.ErrMeliPagoConfigNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("error querying melipago integration type: %w", err)
	}

	if len(row.PlatformCredentialsEncrypted) == 0 {
		return nil, fmt.Errorf("melipago platform credentials not configured")
	}

	credentials, err := r.decryptCredentials(ctx, row.PlatformCredentialsEncrypted)
	if err != nil {
		return nil, fmt.Errorf("error decrypting melipago credentials: %w", err)
	}

	accessToken, _ := credentials["access_token"].(string)
	environment, _ := credentials["environment"].(string)

	if accessToken == "" {
		return nil, meliErrors.ErrInvalidCredentials
	}
	if environment == "" {
		environment = "sandbox"
	}

	r.log.Info(ctx).
		Uint("integration_type_id", row.ID).
		Str("environment", environment).
		Msg("MercadoPago config retrieved")

	return &ports.MeliPagoConfig{
		AccessToken: accessToken,
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
