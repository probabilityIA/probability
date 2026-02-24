package app

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/mocks"
)

// ─────────────────────────────────────────────────────────────────────────────
// Helpers para construir maps de credenciales y configuración en tests
// ─────────────────────────────────────────────────────────────────────────────

func validCredentials() map[string]interface{} {
	return map[string]interface{}{
		"api_key":    "test-api-key-12345",
		"api_secret": "test-api-secret-67890",
	}
}

func validConfig() map[string]interface{} {
	return map[string]interface{}{
		"referer": "https://mi-tienda.softpymes.com",
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Tests: caso feliz
// ─────────────────────────────────────────────────────────────────────────────

// TestTestConnection_Success verifica que cuando el cliente autentica correctamente
// el use case retorna nil.
func TestTestConnection_Success(t *testing.T) {
	// Arrange
	mockClient := &mocks.SoftpymesClientMock{
		TestAuthenticationFn: func(_ context.Context, apiKey, apiSecret, referer, baseURL string) error {
			return nil
		},
	}
	uc := buildUseCase(mockClient)
	ctx := context.Background()

	// Act
	err := uc.TestConnection(ctx, validConfig(), validCredentials())

	// Assert
	if err != nil {
		t.Errorf("TestConnection() esperaba nil, obtuvo: %v", err)
	}
}

// TestTestConnection_PassesCorrectCredentialsToClient verifica que el use case
// pasa exactamente las credenciales recibidas al cliente.
func TestTestConnection_PassesCorrectCredentialsToClient(t *testing.T) {
	// Arrange
	var capturedKey, capturedSecret, capturedReferer string
	mockClient := &mocks.SoftpymesClientMock{
		TestAuthenticationFn: func(_ context.Context, apiKey, apiSecret, referer, _ string) error {
			capturedKey = apiKey
			capturedSecret = apiSecret
			capturedReferer = referer
			return nil
		},
	}
	uc := buildUseCase(mockClient)
	ctx := context.Background()

	expectedKey := "my-api-key"
	expectedSecret := "my-api-secret"
	expectedReferer := "https://cliente.softpymes.com"

	creds := map[string]interface{}{
		"api_key":    expectedKey,
		"api_secret": expectedSecret,
	}
	cfg := map[string]interface{}{
		"referer": expectedReferer,
	}

	// Act
	_ = uc.TestConnection(ctx, cfg, creds)

	// Assert
	if capturedKey != expectedKey {
		t.Errorf("api_key propagada = %q, esperaba %q", capturedKey, expectedKey)
	}
	if capturedSecret != expectedSecret {
		t.Errorf("api_secret propagada = %q, esperaba %q", capturedSecret, expectedSecret)
	}
	if capturedReferer != expectedReferer {
		t.Errorf("referer propagado = %q, esperaba %q", capturedReferer, expectedReferer)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Tests: validación de campos requeridos en credentials
// ─────────────────────────────────────────────────────────────────────────────

// TestTestConnection_MissingApiKey verifica que se retorna error cuando
// api_key no está presente en credentials.
func TestTestConnection_MissingApiKey(t *testing.T) {
	// Arrange
	mockClient := &mocks.SoftpymesClientMock{}
	uc := buildUseCase(mockClient)
	ctx := context.Background()

	creds := map[string]interface{}{
		"api_secret": "some-secret",
		// api_key ausente
	}

	// Act
	err := uc.TestConnection(ctx, validConfig(), creds)

	// Assert
	if err == nil {
		t.Fatal("TestConnection() debería retornar error cuando api_key está ausente")
	}
	expected := "api_key is required in credentials"
	if err.Error() != expected {
		t.Errorf("error = %q, esperaba %q", err.Error(), expected)
	}
}

// TestTestConnection_EmptyApiKey verifica que se retorna error cuando
// api_key está presente pero es un string vacío.
func TestTestConnection_EmptyApiKey(t *testing.T) {
	// Arrange
	mockClient := &mocks.SoftpymesClientMock{}
	uc := buildUseCase(mockClient)
	ctx := context.Background()

	creds := map[string]interface{}{
		"api_key":    "",
		"api_secret": "some-secret",
	}

	// Act
	err := uc.TestConnection(ctx, validConfig(), creds)

	// Assert
	if err == nil {
		t.Fatal("TestConnection() debería retornar error cuando api_key está vacío")
	}
	expected := "api_key is required in credentials"
	if err.Error() != expected {
		t.Errorf("error = %q, esperaba %q", err.Error(), expected)
	}
}

// TestTestConnection_ApiKeyWrongType verifica que se retorna error cuando
// api_key tiene un tipo incorrecto (no string).
func TestTestConnection_ApiKeyWrongType(t *testing.T) {
	// Arrange
	mockClient := &mocks.SoftpymesClientMock{}
	uc := buildUseCase(mockClient)
	ctx := context.Background()

	creds := map[string]interface{}{
		"api_key":    12345, // tipo incorrecto
		"api_secret": "some-secret",
	}

	// Act
	err := uc.TestConnection(ctx, validConfig(), creds)

	// Assert
	if err == nil {
		t.Fatal("TestConnection() debería retornar error cuando api_key tiene tipo incorrecto")
	}
}

// TestTestConnection_MissingApiSecret verifica que se retorna error cuando
// api_secret no está presente en credentials.
func TestTestConnection_MissingApiSecret(t *testing.T) {
	// Arrange
	mockClient := &mocks.SoftpymesClientMock{}
	uc := buildUseCase(mockClient)
	ctx := context.Background()

	creds := map[string]interface{}{
		"api_key": "some-key",
		// api_secret ausente
	}

	// Act
	err := uc.TestConnection(ctx, validConfig(), creds)

	// Assert
	if err == nil {
		t.Fatal("TestConnection() debería retornar error cuando api_secret está ausente")
	}
	expected := "api_secret is required in credentials"
	if err.Error() != expected {
		t.Errorf("error = %q, esperaba %q", err.Error(), expected)
	}
}

// TestTestConnection_EmptyApiSecret verifica que se retorna error cuando
// api_secret está presente pero es un string vacío.
func TestTestConnection_EmptyApiSecret(t *testing.T) {
	// Arrange
	mockClient := &mocks.SoftpymesClientMock{}
	uc := buildUseCase(mockClient)
	ctx := context.Background()

	creds := map[string]interface{}{
		"api_key":    "some-key",
		"api_secret": "",
	}

	// Act
	err := uc.TestConnection(ctx, validConfig(), creds)

	// Assert
	if err == nil {
		t.Fatal("TestConnection() debería retornar error cuando api_secret está vacío")
	}
	expected := "api_secret is required in credentials"
	if err.Error() != expected {
		t.Errorf("error = %q, esperaba %q", err.Error(), expected)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Tests: validación de campos requeridos en config
// ─────────────────────────────────────────────────────────────────────────────

// TestTestConnection_MissingReferer verifica que se retorna error cuando
// referer no está presente en config.
func TestTestConnection_MissingReferer(t *testing.T) {
	// Arrange
	mockClient := &mocks.SoftpymesClientMock{}
	uc := buildUseCase(mockClient)
	ctx := context.Background()

	cfg := map[string]interface{}{
		// referer ausente
	}

	// Act
	err := uc.TestConnection(ctx, cfg, validCredentials())

	// Assert
	if err == nil {
		t.Fatal("TestConnection() debería retornar error cuando referer está ausente")
	}
	expected := "referer is required in config (identificación de instancia del cliente)"
	if err.Error() != expected {
		t.Errorf("error = %q, esperaba %q", err.Error(), expected)
	}
}

// TestTestConnection_EmptyReferer verifica que se retorna error cuando
// referer está presente pero es un string vacío.
func TestTestConnection_EmptyReferer(t *testing.T) {
	// Arrange
	mockClient := &mocks.SoftpymesClientMock{}
	uc := buildUseCase(mockClient)
	ctx := context.Background()

	cfg := map[string]interface{}{
		"referer": "",
	}

	// Act
	err := uc.TestConnection(ctx, cfg, validCredentials())

	// Assert
	if err == nil {
		t.Fatal("TestConnection() debería retornar error cuando referer está vacío")
	}
	expected := "referer is required in config (identificación de instancia del cliente)"
	if err.Error() != expected {
		t.Errorf("error = %q, esperaba %q", err.Error(), expected)
	}
}

// TestTestConnection_RefererWrongType verifica que se retorna error cuando
// referer tiene un tipo incorrecto (no string).
func TestTestConnection_RefererWrongType(t *testing.T) {
	// Arrange
	mockClient := &mocks.SoftpymesClientMock{}
	uc := buildUseCase(mockClient)
	ctx := context.Background()

	cfg := map[string]interface{}{
		"referer": true, // tipo incorrecto
	}

	// Act
	err := uc.TestConnection(ctx, cfg, validCredentials())

	// Assert
	if err == nil {
		t.Fatal("TestConnection() debería retornar error cuando referer tiene tipo incorrecto")
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Tests: errores del cliente
// ─────────────────────────────────────────────────────────────────────────────

// TestTestConnection_ClientAuthenticationError verifica que cuando el cliente
// retorna error, el use case propaga el error envuelto correctamente.
func TestTestConnection_ClientAuthenticationError(t *testing.T) {
	// Arrange
	clientErr := errors.New("invalid credentials")
	mockClient := &mocks.SoftpymesClientMock{
		TestAuthenticationFn: func(_ context.Context, _, _, _, _ string) error {
			return clientErr
		},
	}
	uc := buildUseCase(mockClient)
	ctx := context.Background()

	// Act
	err := uc.TestConnection(ctx, validConfig(), validCredentials())

	// Assert
	if err == nil {
		t.Fatal("TestConnection() debería retornar error cuando el cliente falla")
	}

	// Verificar que el error original está envuelto con errors.Is
	if !errors.Is(err, clientErr) {
		t.Errorf("error no envuelve el error original del cliente. got: %v", err)
	}

	// Verificar que el mensaje contiene el contexto esperado
	expectedFragment := "failed to connect to Softpymes"
	if !containsString(err.Error(), expectedFragment) {
		t.Errorf("error.Error() = %q, debería contener %q", err.Error(), expectedFragment)
	}
}

// TestTestConnection_ClientNetworkError verifica el comportamiento con un error
// de red genérico del cliente.
func TestTestConnection_ClientNetworkError(t *testing.T) {
	// Arrange
	networkErr := fmt.Errorf("connection timeout: dial tcp: i/o timeout")
	mockClient := &mocks.SoftpymesClientMock{
		TestAuthenticationFn: func(_ context.Context, _, _, _, _ string) error {
			return networkErr
		},
	}
	uc := buildUseCase(mockClient)
	ctx := context.Background()

	// Act
	err := uc.TestConnection(ctx, validConfig(), validCredentials())

	// Assert
	if err == nil {
		t.Fatal("TestConnection() debería retornar error ante un error de red")
	}
	if !errors.Is(err, networkErr) {
		t.Errorf("error no envuelve el error de red original. got: %v", err)
	}
}

// TestTestConnection_ValidationOrderApiKeyFirst verifica que la validación de
// api_key ocurre antes que la de api_secret (orden de verificaciones en código).
func TestTestConnection_ValidationOrderApiKeyFirst(t *testing.T) {
	// Arrange
	mockClient := &mocks.SoftpymesClientMock{}
	uc := buildUseCase(mockClient)
	ctx := context.Background()

	creds := map[string]interface{}{
		// Ambas ausentes: se debe reportar api_key primero
	}

	// Act
	err := uc.TestConnection(ctx, validConfig(), creds)

	// Assert
	if err == nil {
		t.Fatal("TestConnection() debería retornar error")
	}
	expected := "api_key is required in credentials"
	if err.Error() != expected {
		t.Errorf("error = %q, esperaba %q (api_key debe validarse primero)", err.Error(), expected)
	}
}

// TestTestConnection_TableDriven ejecuta un conjunto de escenarios usando
// table-driven tests para cubrir múltiples combinaciones de entradas inválidas.
func TestTestConnection_TableDriven(t *testing.T) {
	tests := []struct {
		name        string
		config      map[string]interface{}
		credentials map[string]interface{}
		clientErr   error
		wantErr     bool
		wantErrMsg  string
	}{
		{
			name:        "credenciales y config válidos, cliente OK",
			config:      validConfig(),
			credentials: validCredentials(),
			clientErr:   nil,
			wantErr:     false,
		},
		{
			name:   "api_key vacío",
			config: validConfig(),
			credentials: map[string]interface{}{
				"api_key":    "",
				"api_secret": "secret",
			},
			wantErr:    true,
			wantErrMsg: "api_key is required in credentials",
		},
		{
			name:   "api_secret vacío",
			config: validConfig(),
			credentials: map[string]interface{}{
				"api_key":    "key",
				"api_secret": "",
			},
			wantErr:    true,
			wantErrMsg: "api_secret is required in credentials",
		},
		{
			name: "referer vacío",
			config: map[string]interface{}{
				"referer": "",
			},
			credentials: validCredentials(),
			wantErr:     true,
			wantErrMsg:  "referer is required in config (identificación de instancia del cliente)",
		},
		{
			name:        "cliente retorna error de autenticación",
			config:      validConfig(),
			credentials: validCredentials(),
			clientErr:   errors.New("401 unauthorized"),
			wantErr:     true,
			wantErrMsg:  "failed to connect to Softpymes",
		},
		{
			name:        "maps vacíos (sin ningún campo)",
			config:      map[string]interface{}{},
			credentials: map[string]interface{}{},
			wantErr:     true,
			wantErrMsg:  "api_key is required in credentials",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockClient := &mocks.SoftpymesClientMock{
				TestAuthenticationFn: func(_ context.Context, _, _, _, _ string) error {
					return tt.clientErr
				},
			}
			uc := buildUseCase(mockClient)
			ctx := context.Background()

			// Act
			err := uc.TestConnection(ctx, tt.config, tt.credentials)

			// Assert
			if tt.wantErr {
				if err == nil {
					t.Fatalf("esperaba error, obtuvo nil")
				}
				if tt.wantErrMsg != "" && !containsString(err.Error(), tt.wantErrMsg) {
					t.Errorf("error.Error() = %q, debería contener %q", err.Error(), tt.wantErrMsg)
				}
			} else {
				if err != nil {
					t.Errorf("no esperaba error, obtuvo: %v", err)
				}
			}
		})
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────────────────────────────────────

// containsString reporta si s contiene substr.
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		}())
}
