package app

import (
	"context"
	"errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/mocks"
)

// credencialesCompletas retorna un mapa con las cuatro credenciales requeridas por Siigo.
func credencialesCompletas() map[string]interface{} {
	return map[string]interface{}{
		"username":   "usuario@empresa.com",
		"access_key": "clave-secreta-123",
		"account_id": "subscription-key-abc",
		"partner_id": "partner-xyz",
	}
}

// TestTestConnection_Success verifica el flujo feliz: credenciales completas
// y el cliente de Siigo confirma la autenticación sin error.
func TestTestConnection_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()

	clientMock := &mocks.SiigoClientMock{
		TestAuthenticationFn: func(
			ctx context.Context,
			username, accessKey, accountID, partnerID, baseURL string,
		) error {
			return nil // Autenticación exitosa
		},
	}

	uc := New(clientMock, &mocks.IntegrationCoreMock{}, &mocks.LoggerMock{})

	// Act
	err := uc.TestConnection(ctx, nil, credencialesCompletas())

	// Assert
	if err != nil {
		t.Errorf("se esperaba éxito, se obtuvo error: %v", err)
	}
}

// TestTestConnection_PropagaCredencialesAlCliente verifica que los valores del mapa
// de credenciales se pasen correctamente al cliente de Siigo.
func TestTestConnection_PropagaCredencialesAlCliente(t *testing.T) {
	// Arrange
	ctx := context.Background()

	var capturedUsername, capturedAccessKey, capturedAccountID, capturedPartnerID, capturedBaseURL string

	clientMock := &mocks.SiigoClientMock{
		TestAuthenticationFn: func(
			ctx context.Context,
			username, accessKey, accountID, partnerID, baseURL string,
		) error {
			capturedUsername = username
			capturedAccessKey = accessKey
			capturedAccountID = accountID
			capturedPartnerID = partnerID
			capturedBaseURL = baseURL
			return nil
		},
	}

	credenciales := map[string]interface{}{
		"username":   "test@siigo.com",
		"access_key": "mi-access-key",
		"account_id": "mi-account-id",
		"partner_id": "mi-partner-id",
		"api_url":    "https://api.siigo.com",
	}

	uc := New(clientMock, &mocks.IntegrationCoreMock{}, &mocks.LoggerMock{})

	// Act
	err := uc.TestConnection(ctx, nil, credenciales)

	// Assert
	if err != nil {
		t.Fatalf("se esperaba éxito, se obtuvo error: %v", err)
	}
	if capturedUsername != "test@siigo.com" {
		t.Errorf("username esperado %q, se pasó %q", "test@siigo.com", capturedUsername)
	}
	if capturedAccessKey != "mi-access-key" {
		t.Errorf("access_key esperado %q, se pasó %q", "mi-access-key", capturedAccessKey)
	}
	if capturedAccountID != "mi-account-id" {
		t.Errorf("account_id esperado %q, se pasó %q", "mi-account-id", capturedAccountID)
	}
	if capturedPartnerID != "mi-partner-id" {
		t.Errorf("partner_id esperado %q, se pasó %q", "mi-partner-id", capturedPartnerID)
	}
	if capturedBaseURL != "https://api.siigo.com" {
		t.Errorf("api_url esperado %q, se pasó %q", "https://api.siigo.com", capturedBaseURL)
	}
}

// TestTestConnection_APIURLEsOpcional verifica que el campo api_url sea opcional:
// si no está presente, el use case igual llama al cliente (con string vacío).
func TestTestConnection_APIURLEsOpcional(t *testing.T) {
	// Arrange
	ctx := context.Background()

	llamadoAlCliente := false
	clientMock := &mocks.SiigoClientMock{
		TestAuthenticationFn: func(
			ctx context.Context,
			username, accessKey, accountID, partnerID, baseURL string,
		) error {
			llamadoAlCliente = true
			if baseURL != "" {
				t.Errorf("se esperaba baseURL vacío cuando api_url no está en credenciales, se obtuvo %q", baseURL)
			}
			return nil
		},
	}

	// Credenciales sin api_url
	credenciales := map[string]interface{}{
		"username":   "test@siigo.com",
		"access_key": "mi-access-key",
		"account_id": "mi-account-id",
		"partner_id": "mi-partner-id",
	}

	uc := New(clientMock, &mocks.IntegrationCoreMock{}, &mocks.LoggerMock{})

	// Act
	err := uc.TestConnection(ctx, nil, credenciales)

	// Assert
	if err != nil {
		t.Errorf("se esperaba éxito (api_url es opcional), se obtuvo error: %v", err)
	}
	if !llamadoAlCliente {
		t.Error("se esperaba que el cliente fuera llamado, pero no lo fue")
	}
}

// TestTestConnection_ConfigEsIgnorado verifica que el parámetro config (primer argumento)
// sea ignorado por el use case (la firma del contrato lo requiere pero Siigo no lo usa).
func TestTestConnection_ConfigEsIgnorado(t *testing.T) {
	// Arrange
	ctx := context.Background()

	clientMock := &mocks.SiigoClientMock{
		TestAuthenticationFn: func(_ context.Context, _, _, _, _, _ string) error {
			return nil
		},
	}

	config := map[string]interface{}{
		"alguna_config": "valor_que_siigo_ignora",
	}

	uc := New(clientMock, &mocks.IntegrationCoreMock{}, &mocks.LoggerMock{})

	// Act
	err := uc.TestConnection(ctx, config, credencialesCompletas())

	// Assert — debe tener éxito aunque config tenga valores arbitrarios
	if err != nil {
		t.Errorf("se esperaba éxito ignorando config, se obtuvo error: %v", err)
	}
}

// TestTestConnection_ClienteRetornaError verifica que el error del cliente de Siigo
// se propague sin modificaciones al llamador.
func TestTestConnection_ClienteRetornaError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	errEsperado := errors.New("siigo authentication failed: invalid credentials")

	clientMock := &mocks.SiigoClientMock{
		TestAuthenticationFn: func(
			_ context.Context, _, _, _, _, _ string,
		) error {
			return errEsperado
		},
	}

	uc := New(clientMock, &mocks.IntegrationCoreMock{}, &mocks.LoggerMock{})

	// Act
	err := uc.TestConnection(ctx, nil, credencialesCompletas())

	// Assert
	if err == nil {
		t.Fatal("se esperaba un error del cliente, se obtuvo nil")
	}
	if !errors.Is(err, errEsperado) {
		t.Errorf("error esperado %v, se obtuvo %v", errEsperado, err)
	}
}

// TestTestConnection_ValidacionDeCamposRequeridos usa table-driven tests para
// verificar que cada campo requerido dispara el error correcto cuando falta o está vacío.
func TestTestConnection_ValidacionDeCamposRequeridos(t *testing.T) {
	tests := []struct {
		nombre       string
		credenciales map[string]interface{}
		errorEspMsg  string // Fragmento del mensaje de error esperado
	}{
		{
			nombre: "falta username",
			credenciales: map[string]interface{}{
				"access_key": "clave",
				"account_id": "cuenta",
				"partner_id": "partner",
			},
			errorEspMsg: "username",
		},
		{
			nombre: "username vacio",
			credenciales: map[string]interface{}{
				"username":   "",
				"access_key": "clave",
				"account_id": "cuenta",
				"partner_id": "partner",
			},
			errorEspMsg: "username",
		},
		{
			nombre: "falta access_key",
			credenciales: map[string]interface{}{
				"username":   "user@test.com",
				"account_id": "cuenta",
				"partner_id": "partner",
			},
			errorEspMsg: "access_key",
		},
		{
			nombre: "access_key vacio",
			credenciales: map[string]interface{}{
				"username":   "user@test.com",
				"access_key": "",
				"account_id": "cuenta",
				"partner_id": "partner",
			},
			errorEspMsg: "access_key",
		},
		{
			nombre: "falta account_id",
			credenciales: map[string]interface{}{
				"username":   "user@test.com",
				"access_key": "clave",
				"partner_id": "partner",
			},
			errorEspMsg: "account_id",
		},
		{
			nombre: "account_id vacio",
			credenciales: map[string]interface{}{
				"username":   "user@test.com",
				"access_key": "clave",
				"account_id": "",
				"partner_id": "partner",
			},
			errorEspMsg: "account_id",
		},
		{
			nombre: "falta partner_id",
			credenciales: map[string]interface{}{
				"username":   "user@test.com",
				"access_key": "clave",
				"account_id": "cuenta",
			},
			errorEspMsg: "partner_id",
		},
		{
			nombre: "partner_id vacio",
			credenciales: map[string]interface{}{
				"username":   "user@test.com",
				"access_key": "clave",
				"account_id": "cuenta",
				"partner_id": "",
			},
			errorEspMsg: "partner_id",
		},
		{
			nombre:       "mapa de credenciales vacio",
			credenciales: map[string]interface{}{},
			errorEspMsg:  "username",
		},
	}

	// El cliente NO debe ser llamado cuando hay error de validación
	clientMock := &mocks.SiigoClientMock{
		TestAuthenticationFn: func(_ context.Context, _, _, _, _, _ string) error {
			return nil // Si se llega aquí, el test fallará por su propia aserción
		},
	}

	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			// Arrange
			ctx := context.Background()
			clienteLlamado := false
			mockConCaptura := &mocks.SiigoClientMock{
				TestAuthenticationFn: func(_ context.Context, _, _, _, _, _ string) error {
					clienteLlamado = true
					return clientMock.TestAuthenticationFn(ctx, "", "", "", "", "")
				},
			}

			uc := New(mockConCaptura, &mocks.IntegrationCoreMock{}, &mocks.LoggerMock{})

			// Act
			err := uc.TestConnection(ctx, nil, tt.credenciales)

			// Assert: debe retornar error de validación
			if err == nil {
				t.Fatalf("[%s] se esperaba error de validación, se obtuvo nil", tt.nombre)
			}

			// El mensaje de error debe mencionar el campo faltante
			if tt.errorEspMsg != "" {
				encontrado := false
				msg := err.Error()
				for _, palabra := range []string{tt.errorEspMsg} {
					if contains(msg, palabra) {
						encontrado = true
						break
					}
				}
				if !encontrado {
					t.Errorf("[%s] error esperado conteniendo %q, se obtuvo: %q", tt.nombre, tt.errorEspMsg, msg)
				}
			}

			// El cliente NO debe ser invocado si la validación falla
			if clienteLlamado {
				t.Errorf("[%s] el cliente no debería ser llamado cuando la validación falla", tt.nombre)
			}
		})
	}
}

// TestTestConnection_TipoIncorrectoEnCredenciales verifica que si un campo existe
// en el mapa pero no es string (tipo incorrecto), se trate como ausente.
func TestTestConnection_TipoIncorrectoEnCredenciales(t *testing.T) {
	// Arrange
	ctx := context.Background()
	uc := New(&mocks.SiigoClientMock{}, &mocks.IntegrationCoreMock{}, &mocks.LoggerMock{})

	credenciales := map[string]interface{}{
		"username":   12345, // int en lugar de string
		"access_key": "clave",
		"account_id": "cuenta",
		"partner_id": "partner",
	}

	// Act
	err := uc.TestConnection(ctx, nil, credenciales)

	// Assert: la aserción de tipo falla → username queda vacío → error de validación
	if err == nil {
		t.Fatal("se esperaba error por tipo incorrecto en username, se obtuvo nil")
	}
}

// contains es un helper local para verificar substring sin importar strings en los tests.
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && indexOf(s, substr) >= 0))
}

// indexOf encuentra el índice de substr en s, o -1 si no está.
func indexOf(s, substr string) int {
	if len(substr) == 0 {
		return 0
	}
	if len(substr) > len(s) {
		return -1
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
