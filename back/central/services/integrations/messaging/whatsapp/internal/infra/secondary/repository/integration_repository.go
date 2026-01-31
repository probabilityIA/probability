package repository

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
	"gorm.io/gorm"
)

// IntegrationRepository implementa ports.IIntegrationRepository
type IntegrationRepository struct {
	db            db.IDatabase
	log           log.ILogger
	encryptionKey []byte
}

// GetWhatsAppConfig obtiene phone_number_id y access_token desde la base de datos
func (r *IntegrationRepository) GetWhatsAppConfig(ctx context.Context, businessID uint) (*ports.WhatsAppConfig, error) {
	// Estructura para almacenar la integración
	type Integration struct {
		ID          uint            `json:"id"`
		Config      json.RawMessage `json:"config"`
		Credentials json.RawMessage `json:"credentials"`
	}

	var integration Integration

	// Primero intentar obtener la integración del business específico
	err := r.db.Conn(ctx).
		Table("integrations").
		Select("id, config, credentials").
		Where("integration_type_id = ?", 2).
		Where("business_id = ?", businessID).
		First(&integration).Error

	if err == gorm.ErrRecordNotFound {
		// Si no existe, usar la integración global (business_id IS NULL)
		r.log.Info(ctx).
			Uint("business_id", businessID).
			Msg("[Integration Repository] - no se encontró integración específica, usando global")

		err = r.db.Conn(ctx).
			Table("integrations").
			Select("id, config, credentials").
			Where("integration_type_id = ?", 2).
			Where("business_id IS NULL").
			First(&integration).Error

		if err != nil {
			r.log.Error(ctx).Err(err).Msg("[Integration Repository] - no se encontró integración de WhatsApp")
			return nil, fmt.Errorf("no se encontró integración de WhatsApp")
		}
	} else if err != nil {
		r.log.Error(ctx).Err(err).Msg("[Integration Repository] - error consultando integración")
		return nil, fmt.Errorf("error consultando integración: %w", err)
	}

	// 1. Parsear el config JSON para obtener phone_number_id
	var config map[string]interface{}
	if err := json.Unmarshal(integration.Config, &config); err != nil {
		r.log.Error(ctx).Err(err).Msg("[Integration Repository] - error parseando config")
		return nil, fmt.Errorf("error parseando config: %w", err)
	}

	// Extraer phone_number_id
	phoneNumberIDValue, exists := config["phone_number_id"]
	if !exists {
		r.log.Error(ctx).Msg("[Integration Repository] - phone_number_id no encontrado en config")
		return nil, fmt.Errorf("phone_number_id no encontrado en configuración")
	}

	phoneNumberIDStr, ok := phoneNumberIDValue.(string)
	if !ok {
		r.log.Error(ctx).Msg("[Integration Repository] - phone_number_id no es string")
		return nil, fmt.Errorf("phone_number_id debe ser string")
	}

	phoneNumberID, err := strconv.ParseUint(phoneNumberIDStr, 10, 64)
	if err != nil {
		r.log.Error(ctx).Err(err).Str("phone_number_id", phoneNumberIDStr).Msg("[Integration Repository] - error parseando phone_number_id")
		return nil, fmt.Errorf("error parseando phone_number_id: %w", err)
	}

	// 2. Parsear credentials para obtener el wrapper encriptado
	var credentialsWrapper map[string]interface{}
	if err := json.Unmarshal(integration.Credentials, &credentialsWrapper); err != nil {
		r.log.Error(ctx).Err(err).Msg("[Integration Repository] - error parseando credentials wrapper")
		return nil, fmt.Errorf("error parseando credentials: %w", err)
	}

	// Extraer el valor encriptado del wrapper
	encryptedValue, exists := credentialsWrapper["encrypted"]
	if !exists {
		r.log.Error(ctx).Msg("[Integration Repository] - credentials no tienen campo 'encrypted'")
		return nil, fmt.Errorf("credentials inválidas: falta campo 'encrypted'")
	}

	encryptedStr, ok := encryptedValue.(string)
	if !ok {
		r.log.Error(ctx).Msg("[Integration Repository] - encrypted no es string")
		return nil, fmt.Errorf("credentials inválidas: 'encrypted' debe ser string")
	}

	// Decodificar base64
	encryptedBytes, err := base64.StdEncoding.DecodeString(encryptedStr)
	if err != nil {
		r.log.Error(ctx).Err(err).Msg("[Integration Repository] - error decodificando credentials desde base64")
		return nil, fmt.Errorf("error decodificando credentials: %w", err)
	}

	// Desencriptar usando AES-GCM
	decryptedCredentials, err := r.decryptCredentials(ctx, encryptedBytes)
	if err != nil {
		r.log.Error(ctx).Err(err).Msg("[Integration Repository] - error desencriptando credentials")
		return nil, fmt.Errorf("error desencriptando credentials: %w", err)
	}

	// Extraer access_token de las credenciales desencriptadas
	accessTokenValue, exists := decryptedCredentials["access_token"]
	if !exists {
		r.log.Error(ctx).Msg("[Integration Repository] - access_token no encontrado en credentials")
		return nil, fmt.Errorf("access_token no encontrado en credentials")
	}

	accessToken, ok := accessTokenValue.(string)
	if !ok {
		r.log.Error(ctx).Msg("[Integration Repository] - access_token no es string")
		return nil, fmt.Errorf("access_token debe ser string")
	}

	r.log.Info(ctx).
		Uint("integration_id", integration.ID).
		Uint("phone_number_id", uint(phoneNumberID)).
		Msg("[Integration Repository] - Configuración de WhatsApp obtenida desde DB")

	return &ports.WhatsAppConfig{
		PhoneNumberID: uint(phoneNumberID),
		AccessToken:   accessToken,
		IntegrationID: integration.ID,
	}, nil
}

// decryptCredentials desencripta credenciales usando AES-256-GCM
func (r *IntegrationRepository) decryptCredentials(ctx context.Context, ciphertext []byte) (map[string]interface{}, error) {
	block, err := aes.NewCipher(r.encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("error al crear cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("error al crear GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext demasiado corto")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("error al desencriptar: %w", err)
	}

	// Convertir JSON a mapa
	var credentials map[string]interface{}
	if err := json.Unmarshal(plaintext, &credentials); err != nil {
		return nil, fmt.Errorf("error al deserializar credenciales: %w", err)
	}

	return credentials, nil
}
