package encryption

import (
	"context"
	"io"
	"testing"

	"github.com/rs/zerolog"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/stretchr/testify/assert"
)

// ============================================
// Logger mock sin dependencias externas
// ============================================

type mockEncLogger struct {
	l zerolog.Logger
}

func newMockEncLogger() *mockEncLogger {
	return &mockEncLogger{l: zerolog.New(io.Discard)}
}

func (m *mockEncLogger) Info(ctx ...context.Context) *zerolog.Event  { return m.l.Info() }
func (m *mockEncLogger) Error(ctx ...context.Context) *zerolog.Event { return m.l.Error() }
func (m *mockEncLogger) Debug(ctx ...context.Context) *zerolog.Event { return m.l.Debug() }
func (m *mockEncLogger) Warn(ctx ...context.Context) *zerolog.Event  { return m.l.Warn() }
func (m *mockEncLogger) Fatal(ctx ...context.Context) *zerolog.Event { return m.l.Fatal() }
func (m *mockEncLogger) Panic(ctx ...context.Context) *zerolog.Event { return m.l.Panic() }
func (m *mockEncLogger) With() zerolog.Context                       { return m.l.With() }
func (m *mockEncLogger) WithService(s string) log.ILogger            { return m }
func (m *mockEncLogger) WithModule(s string) log.ILogger             { return m }
func (m *mockEncLogger) WithBusinessID(id uint) log.ILogger          { return m }

// testKey es una clave AES-256 válida (32 bytes exactos para tests)
const testKey = "01234567890123456789012345678901"

// newTestEncryptionService crea el servicio de encriptación directamente con clave conocida
func newTestEncryptionService() *encryptionService {
	return &encryptionService{
		key: []byte(testKey),
		log: newMockEncLogger(),
	}
}

// ============================================
// encrypt / decrypt (métodos internos)
// ============================================

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	// Arrange
	svc := newTestEncryptionService()
	plaintext := []byte("dato secreto de prueba")

	// Act — encriptar
	ciphertext, err := svc.encrypt(plaintext)
	assert.NoError(t, err)
	assert.NotNil(t, ciphertext)
	assert.NotEqual(t, plaintext, ciphertext)

	// Act — desencriptar
	decrypted, err := svc.decrypt(ciphertext)
	assert.NoError(t, err)

	// Assert — debe ser igual al original
	assert.Equal(t, plaintext, decrypted)
}

func TestEncrypt_GeneraNonceDiferenteCadaVez(t *testing.T) {
	// Arrange
	svc := newTestEncryptionService()
	plaintext := []byte("mismo dato")

	// Act — encriptar dos veces el mismo plaintext
	ciphertext1, err1 := svc.encrypt(plaintext)
	ciphertext2, err2 := svc.encrypt(plaintext)

	// Assert — cada cifrado produce un resultado diferente (nonce aleatorio)
	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NotEqual(t, ciphertext1, ciphertext2)
}

func TestDecrypt_CiphertextDemasiadoCorto(t *testing.T) {
	// Arrange
	svc := newTestEncryptionService()

	// Un ciphertext de 2 bytes es menor que el tamaño del nonce GCM (12 bytes)
	ciphertextCorto := []byte{0x01, 0x02}

	// Act
	resultado, err := svc.decrypt(ciphertextCorto)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "demasiado corto")
	assert.Nil(t, resultado)
}

func TestDecrypt_CiphertextManipulado(t *testing.T) {
	// Arrange
	svc := newTestEncryptionService()
	plaintext := []byte("datos válidos")

	ciphertext, err := svc.encrypt(plaintext)
	assert.NoError(t, err)

	// Corromper el ciphertext
	ciphertext[len(ciphertext)-1] ^= 0xFF

	// Act
	resultado, err := svc.decrypt(ciphertext)

	// Assert — debe fallar la verificación de autenticidad GCM
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "desencriptar")
	assert.Nil(t, resultado)
}

// ============================================
// EncryptCredentials / DecryptCredentials
// ============================================

func TestEncryptCredentials_MapABytes(t *testing.T) {
	// Arrange
	svc := newTestEncryptionService()
	ctx := context.Background()

	credentials := map[string]interface{}{
		"api_key":      "key_123",
		"shop_domain":  "mi-tienda.myshopify.com",
		"access_token": "shpat_abc123",
	}

	// Act
	encrypted, err := svc.EncryptCredentials(ctx, credentials)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, encrypted)
	assert.Greater(t, len(encrypted), 0)
}

func TestEncryptDecryptCredentials_RoundTrip(t *testing.T) {
	// Arrange
	svc := newTestEncryptionService()
	ctx := context.Background()

	original := map[string]interface{}{
		"api_key":     "key_abc",
		"api_secret":  "secret_xyz",
		"shop_domain": "test.myshopify.com",
	}

	// Act — encriptar y desencriptar
	encrypted, err := svc.EncryptCredentials(ctx, original)
	assert.NoError(t, err)

	decrypted, err := svc.DecryptCredentials(ctx, encrypted)
	assert.NoError(t, err)

	// Assert — debe recuperar los valores originales
	assert.Equal(t, original["api_key"], decrypted["api_key"])
	assert.Equal(t, original["api_secret"], decrypted["api_secret"])
	assert.Equal(t, original["shop_domain"], decrypted["shop_domain"])
}

func TestDecryptCredentials_DatosCorruptos(t *testing.T) {
	// Arrange
	svc := newTestEncryptionService()
	ctx := context.Background()

	// Datos que no son un ciphertext válido
	datosCorruptos := []byte("esto no es un ciphertext valido 0123456789abc")

	// Act
	resultado, err := svc.DecryptCredentials(ctx, datosCorruptos)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resultado)
}

func TestDecryptCredentials_JSONInvalidoTrasDesencriptar(t *testing.T) {
	// Arrange
	svc := newTestEncryptionService()
	ctx := context.Background()

	// Encriptar datos que no son JSON válido (texto plano)
	plainNonJSON := []byte("no soy JSON")
	ciphertext, err := svc.encrypt(plainNonJSON)
	assert.NoError(t, err)

	// Act — intentar desencriptar como credenciales
	resultado, err := svc.DecryptCredentials(ctx, ciphertext)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "deserializar")
	assert.Nil(t, resultado)
}

// ============================================
// EncryptValue / DecryptValue
// ============================================

func TestEncryptValue_RetornaStringBase64(t *testing.T) {
	// Arrange
	svc := newTestEncryptionService()
	ctx := context.Background()

	// Act
	encrypted, err := svc.EncryptValue(ctx, "valor secreto")

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, encrypted)
}

func TestEncryptDecryptValue_RoundTrip(t *testing.T) {
	// Arrange
	svc := newTestEncryptionService()
	ctx := context.Background()
	original := "token_secreto_12345"

	// Act — encriptar y desencriptar
	encrypted, err := svc.EncryptValue(ctx, original)
	assert.NoError(t, err)

	decrypted, err := svc.DecryptValue(ctx, encrypted)
	assert.NoError(t, err)

	// Assert
	assert.Equal(t, original, decrypted)
}

func TestDecryptValue_StringBase64Invalido(t *testing.T) {
	// Arrange
	svc := newTestEncryptionService()
	ctx := context.Background()

	// Act — base64 inválido
	resultado, err := svc.DecryptValue(ctx, "!!!no es base64!!!")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "decodificar")
	assert.Equal(t, "", resultado)
}

func TestEncryptValue_CadenaVacia(t *testing.T) {
	// Arrange
	svc := newTestEncryptionService()
	ctx := context.Background()

	// Act — encriptar cadena vacía
	encrypted, err := svc.EncryptValue(ctx, "")
	assert.NoError(t, err)

	decrypted, err := svc.DecryptValue(ctx, encrypted)
	assert.NoError(t, err)

	// Assert — debe recuperar la cadena vacía
	assert.Equal(t, "", decrypted)
}

func TestEncryptValue_StringLargo(t *testing.T) {
	// Arrange
	svc := newTestEncryptionService()
	ctx := context.Background()

	// Crear un string de 1000 caracteres
	largo := ""
	for i := 0; i < 1000; i++ {
		largo += "a"
	}

	// Act
	encrypted, err := svc.EncryptValue(ctx, largo)
	assert.NoError(t, err)

	decrypted, err := svc.DecryptValue(ctx, encrypted)
	assert.NoError(t, err)

	// Assert
	assert.Equal(t, largo, decrypted)
	assert.Len(t, decrypted, 1000)
}
