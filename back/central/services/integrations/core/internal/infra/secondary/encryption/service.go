package encryption

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"

	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

type encryptionService struct {
	key []byte
	log log.ILogger
}

// newEncryptionService crea una nueva instancia del servicio de encriptación (privado)
func newEncryptionService(config env.IConfig, logger log.ILogger) *encryptionService {
	encryptionKey := config.Get("ENCRYPTION_KEY")
	if encryptionKey == "" {
		logger.Fatal(context.Background()).
			Msg("ENCRYPTION_KEY no está configurada - es requerida para encriptar credenciales")
	}

	// La clave debe tener exactamente 32 bytes para AES-256
	key := []byte(encryptionKey)
	if len(key) != 32 {
		logger.Fatal(context.Background()).
			Int("key_length", len(key)).
			Msg("ENCRYPTION_KEY debe tener exactamente 32 bytes (256 bits) para AES-256")
	}

	return &encryptionService{
		key: key,
		log: logger,
	}
}

// EncryptCredentials encripta un mapa de credenciales
func (s *encryptionService) EncryptCredentials(ctx context.Context, credentials map[string]interface{}) ([]byte, error) {
	// Convertir el mapa a JSON
	jsonData, err := json.Marshal(credentials)
	if err != nil {
		s.log.Error(ctx).Err(err).Msg("Error al serializar credenciales a JSON")
		return nil, fmt.Errorf("error al serializar credenciales: %w", err)
	}

	// Encriptar el JSON
	encrypted, err := s.encrypt(jsonData)
	if err != nil {
		s.log.Error(ctx).Err(err).Msg("Error al encriptar credenciales")
		return nil, fmt.Errorf("error al encriptar credenciales: %w", err)
	}

	return encrypted, nil
}

// DecryptCredentials desencripta credenciales
func (s *encryptionService) DecryptCredentials(ctx context.Context, encryptedData []byte) (map[string]interface{}, error) {
	// Desencriptar
	decrypted, err := s.decrypt(encryptedData)
	if err != nil {
		s.log.Error(ctx).Err(err).Msg("Error al desencriptar credenciales")
		return nil, fmt.Errorf("error al desencriptar credenciales: %w", err)
	}

	// Convertir JSON a mapa
	var credentials map[string]interface{}
	if err := json.Unmarshal(decrypted, &credentials); err != nil {
		s.log.Error(ctx).Err(err).Msg("Error al deserializar credenciales desde JSON")
		return nil, fmt.Errorf("error al deserializar credenciales: %w", err)
	}

	return credentials, nil
}

// EncryptValue encripta un valor individual
func (s *encryptionService) EncryptValue(ctx context.Context, value string) (string, error) {
	encrypted, err := s.encrypt([]byte(value))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// DecryptValue desencripta un valor individual
func (s *encryptionService) DecryptValue(ctx context.Context, encryptedValue string) (string, error) {
	encrypted, err := base64.StdEncoding.DecodeString(encryptedValue)
	if err != nil {
		return "", fmt.Errorf("error al decodificar valor encriptado: %w", err)
	}

	decrypted, err := s.decrypt(encrypted)
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}

// encrypt encripta datos usando AES-256-GCM
func (s *encryptionService) encrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(s.key)
	if err != nil {
		return nil, fmt.Errorf("error al crear cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("error al crear GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("error al generar nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// decrypt desencripta datos usando AES-256-GCM
func (s *encryptionService) decrypt(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(s.key)
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

	return plaintext, nil
}
