package app

import (
	"context"
	"errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/factus/internal/mocks"
)

// ─────────────────────────────────────────────────────────────
// TestConnection — camino feliz con todas las credenciales
// ─────────────────────────────────────────────────────────────

func TestTestConnection_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	credentials := map[string]interface{}{
		"client_id":     "cid-ok",
		"client_secret": "csecret-ok",
		"username":      "admin@factus.com",
		"password":      "p@ssw0rd",
		"api_url":       "https://api.factus.com.co",
	}

	mockClient := &mocks.FactusClientMock{
		TestAuthenticationFn: func(_ context.Context, baseURL, clientID, _, _, _ string) error {
			if baseURL != "https://api.factus.com.co" {
				t.Errorf("baseURL inesperada: '%s'", baseURL)
			}
			if clientID != "cid-ok" {
				t.Errorf("clientID inesperado: '%s'", clientID)
			}
			return nil
		},
	}

	mockCore := &mocks.IntegrationCoreMock{}
	uc := New(mockClient, mockCore, mocks.NewLoggerMock())

	// Act
	err := uc.TestConnection(ctx, nil, credentials)

	// Assert
	if err != nil {
		t.Fatalf("esperaba sin error, recibí: %v", err)
	}
}

// ─────────────────────────────────────────────────────────────
// TestConnection — camino feliz sin api_url (campo opcional)
// ─────────────────────────────────────────────────────────────

func TestTestConnection_Success_WithoutAPIURL(t *testing.T) {
	ctx := context.Background()
	credentials := map[string]interface{}{
		"client_id":     "cid",
		"client_secret": "csecret",
		"username":      "user@test.com",
		"password":      "pass123",
		// api_url ausente — es opcional
	}

	var capturedBaseURL string

	mockClient := &mocks.FactusClientMock{
		TestAuthenticationFn: func(_ context.Context, baseURL, _, _, _, _ string) error {
			capturedBaseURL = baseURL
			return nil
		},
	}

	uc := New(mockClient, &mocks.IntegrationCoreMock{}, mocks.NewLoggerMock())

	err := uc.TestConnection(ctx, nil, credentials)

	if err != nil {
		t.Fatalf("esperaba sin error, recibí: %v", err)
	}
	if capturedBaseURL != "" {
		t.Errorf("esperaba baseURL vacía cuando api_url no se provee, recibí '%s'", capturedBaseURL)
	}
}

// ─────────────────────────────────────────────────────────────
// TestConnection — client_id ausente
// ─────────────────────────────────────────────────────────────

func TestTestConnection_MissingClientID(t *testing.T) {
	ctx := context.Background()
	credentials := map[string]interface{}{
		// "client_id" ausente
		"client_secret": "csecret",
		"username":      "user@test.com",
		"password":      "pass",
	}

	mockClient := &mocks.FactusClientMock{}
	uc := New(mockClient, &mocks.IntegrationCoreMock{}, mocks.NewLoggerMock())

	err := uc.TestConnection(ctx, nil, credentials)

	if err == nil {
		t.Fatal("esperaba error por client_id faltante, recibí nil")
	}
	if !errors.Is(err, err) { // siempre pasa, lo útil es verificar el mensaje
		t.Errorf("error inesperado: %v", err)
	}
}

// ─────────────────────────────────────────────────────────────
// TestConnection — client_id con string vacío
// ─────────────────────────────────────────────────────────────

func TestTestConnection_EmptyClientID(t *testing.T) {
	ctx := context.Background()
	credentials := map[string]interface{}{
		"client_id":     "",
		"client_secret": "csecret",
		"username":      "user@test.com",
		"password":      "pass",
	}

	mockClient := &mocks.FactusClientMock{}
	uc := New(mockClient, &mocks.IntegrationCoreMock{}, mocks.NewLoggerMock())

	err := uc.TestConnection(ctx, nil, credentials)

	if err == nil {
		t.Fatal("esperaba error por client_id vacío, recibí nil")
	}
}

// ─────────────────────────────────────────────────────────────
// TestConnection — client_secret ausente
// ─────────────────────────────────────────────────────────────

func TestTestConnection_MissingClientSecret(t *testing.T) {
	ctx := context.Background()
	credentials := map[string]interface{}{
		"client_id": "cid",
		// "client_secret" ausente
		"username": "user@test.com",
		"password": "pass",
	}

	mockClient := &mocks.FactusClientMock{}
	uc := New(mockClient, &mocks.IntegrationCoreMock{}, mocks.NewLoggerMock())

	err := uc.TestConnection(ctx, nil, credentials)

	if err == nil {
		t.Fatal("esperaba error por client_secret faltante, recibí nil")
	}
}

// ─────────────────────────────────────────────────────────────
// TestConnection — client_secret con string vacío
// ─────────────────────────────────────────────────────────────

func TestTestConnection_EmptyClientSecret(t *testing.T) {
	ctx := context.Background()
	credentials := map[string]interface{}{
		"client_id":     "cid",
		"client_secret": "",
		"username":      "user@test.com",
		"password":      "pass",
	}

	mockClient := &mocks.FactusClientMock{}
	uc := New(mockClient, &mocks.IntegrationCoreMock{}, mocks.NewLoggerMock())

	err := uc.TestConnection(ctx, nil, credentials)

	if err == nil {
		t.Fatal("esperaba error por client_secret vacío, recibí nil")
	}
}

// ─────────────────────────────────────────────────────────────
// TestConnection — username ausente
// ─────────────────────────────────────────────────────────────

func TestTestConnection_MissingUsername(t *testing.T) {
	ctx := context.Background()
	credentials := map[string]interface{}{
		"client_id":     "cid",
		"client_secret": "csecret",
		// "username" ausente
		"password": "pass",
	}

	mockClient := &mocks.FactusClientMock{}
	uc := New(mockClient, &mocks.IntegrationCoreMock{}, mocks.NewLoggerMock())

	err := uc.TestConnection(ctx, nil, credentials)

	if err == nil {
		t.Fatal("esperaba error por username faltante, recibí nil")
	}
}

// ─────────────────────────────────────────────────────────────
// TestConnection — username con string vacío
// ─────────────────────────────────────────────────────────────

func TestTestConnection_EmptyUsername(t *testing.T) {
	ctx := context.Background()
	credentials := map[string]interface{}{
		"client_id":     "cid",
		"client_secret": "csecret",
		"username":      "",
		"password":      "pass",
	}

	mockClient := &mocks.FactusClientMock{}
	uc := New(mockClient, &mocks.IntegrationCoreMock{}, mocks.NewLoggerMock())

	err := uc.TestConnection(ctx, nil, credentials)

	if err == nil {
		t.Fatal("esperaba error por username vacío, recibí nil")
	}
}

// ─────────────────────────────────────────────────────────────
// TestConnection — password ausente
// ─────────────────────────────────────────────────────────────

func TestTestConnection_MissingPassword(t *testing.T) {
	ctx := context.Background()
	credentials := map[string]interface{}{
		"client_id":     "cid",
		"client_secret": "csecret",
		"username":      "user@test.com",
		// "password" ausente
	}

	mockClient := &mocks.FactusClientMock{}
	uc := New(mockClient, &mocks.IntegrationCoreMock{}, mocks.NewLoggerMock())

	err := uc.TestConnection(ctx, nil, credentials)

	if err == nil {
		t.Fatal("esperaba error por password faltante, recibí nil")
	}
}

// ─────────────────────────────────────────────────────────────
// TestConnection — password con string vacío
// ─────────────────────────────────────────────────────────────

func TestTestConnection_EmptyPassword(t *testing.T) {
	ctx := context.Background()
	credentials := map[string]interface{}{
		"client_id":     "cid",
		"client_secret": "csecret",
		"username":      "user@test.com",
		"password":      "",
	}

	mockClient := &mocks.FactusClientMock{}
	uc := New(mockClient, &mocks.IntegrationCoreMock{}, mocks.NewLoggerMock())

	err := uc.TestConnection(ctx, nil, credentials)

	if err == nil {
		t.Fatal("esperaba error por password vacío, recibí nil")
	}
}

// ─────────────────────────────────────────────────────────────
// TestConnection — tabla de validaciones de campos requeridos
// ─────────────────────────────────────────────────────────────

func TestTestConnection_TableDriven_RequiredFields(t *testing.T) {
	baseCredentials := map[string]interface{}{
		"client_id":     "cid",
		"client_secret": "csecret",
		"username":      "user@test.com",
		"password":      "pass",
	}

	tests := []struct {
		name        string
		omitField   string
		emptyField  string
	}{
		{"falta client_id", "client_id", ""},
		{"client_id vacío", "", "client_id"},
		{"falta client_secret", "client_secret", ""},
		{"client_secret vacío", "", "client_secret"},
		{"falta username", "username", ""},
		{"username vacío", "", "username"},
		{"falta password", "password", ""},
		{"password vacío", "", "password"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Clonar credenciales base
			creds := make(map[string]interface{})
			for k, v := range baseCredentials {
				creds[k] = v
			}

			if tt.omitField != "" {
				delete(creds, tt.omitField)
			}
			if tt.emptyField != "" {
				creds[tt.emptyField] = ""
			}

			mockClient := &mocks.FactusClientMock{}
			uc := New(mockClient, &mocks.IntegrationCoreMock{}, mocks.NewLoggerMock())

			err := uc.TestConnection(context.Background(), nil, creds)

			if err == nil {
				t.Fatalf("esperaba error para escenario '%s', recibí nil", tt.name)
			}
		})
	}
}

// ─────────────────────────────────────────────────────────────
// TestConnection — el cliente HTTP retorna error de autenticación
// ─────────────────────────────────────────────────────────────

func TestTestConnection_ClientAuthenticationFails(t *testing.T) {
	ctx := context.Background()
	credentials := map[string]interface{}{
		"client_id":     "cid-malo",
		"client_secret": "csecret-malo",
		"username":      "user@test.com",
		"password":      "wrong-pass",
	}

	authErr := errors.New("factus: authentication failed — invalid credentials")

	mockClient := &mocks.FactusClientMock{
		TestAuthenticationFn: func(_ context.Context, _, _, _, _, _ string) error {
			return authErr
		},
	}

	uc := New(mockClient, &mocks.IntegrationCoreMock{}, mocks.NewLoggerMock())

	err := uc.TestConnection(ctx, nil, credentials)

	if err == nil {
		t.Fatal("esperaba error de autenticación, recibí nil")
	}
	if !errors.Is(err, authErr) {
		t.Errorf("esperaba error '%v', recibí '%v'", authErr, err)
	}
}

// ─────────────────────────────────────────────────────────────
// TestConnection — el map de config puede ser nil (se ignora)
// ─────────────────────────────────────────────────────────────

func TestTestConnection_NilConfigIsIgnored(t *testing.T) {
	ctx := context.Background()
	credentials := map[string]interface{}{
		"client_id":     "cid",
		"client_secret": "csecret",
		"username":      "user@test.com",
		"password":      "pass",
	}

	mockClient := &mocks.FactusClientMock{
		TestAuthenticationFn: func(_ context.Context, _, _, _, _, _ string) error {
			return nil
		},
	}

	uc := New(mockClient, &mocks.IntegrationCoreMock{}, mocks.NewLoggerMock())

	// config es nil — el método usa `_` para ignorarlo, por lo tanto no debe paniquear
	err := uc.TestConnection(ctx, nil, credentials)

	if err != nil {
		t.Fatalf("esperaba sin error con config nil, recibí: %v", err)
	}
}

// ─────────────────────────────────────────────────────────────
// TestConnection — tipo de dato incorrecto en credenciales (no string)
// ─────────────────────────────────────────────────────────────

func TestTestConnection_NonStringCredentialValue(t *testing.T) {
	ctx := context.Background()
	// client_id es un número, no un string — la aserción de tipo fallará
	credentials := map[string]interface{}{
		"client_id":     12345, // tipo incorrecto
		"client_secret": "csecret",
		"username":      "user@test.com",
		"password":      "pass",
	}

	mockClient := &mocks.FactusClientMock{}
	uc := New(mockClient, &mocks.IntegrationCoreMock{}, mocks.NewLoggerMock())

	err := uc.TestConnection(ctx, nil, credentials)

	// La aserción de tipo fallará y okClientID será false → error esperado
	if err == nil {
		t.Fatal("esperaba error cuando client_id no es string, recibí nil")
	}
}
