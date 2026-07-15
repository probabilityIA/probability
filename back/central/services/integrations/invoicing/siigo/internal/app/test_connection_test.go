package app

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/mocks"
)

func credencialesCompletas() map[string]any {
	return map[string]any{
		"username":   "usuario@empresa.com",
		"access_key": "clave-secreta-123",
		"account_id": "subscription-key-abc",
		"partner_id": "partner-xyz",
	}
}

func configProd() map[string]any {
	return map[string]any{
		"base_url":      "https://api.siigo.com",
		"base_url_test": "http://back-testing:9095",
	}
}

func TestTestConnection_Success(t *testing.T) {
	ctx := context.Background()

	clientMock := &mocks.SiigoClientMock{
		TestAuthenticationFn: func(_ context.Context, _, _, _, _, _ string) error {
			return nil
		},
	}

	uc := New(clientMock, &mocks.IntegrationCoreMock{}, nil, nil, &mocks.LoggerMock{})

	if err := uc.TestConnection(ctx, configProd(), credencialesCompletas()); err != nil {
		t.Errorf("se esperaba exito, se obtuvo error: %v", err)
	}
}

func TestTestConnection_PropagaCredencialesAlCliente(t *testing.T) {
	ctx := context.Background()

	var capturedUsername, capturedAccessKey, capturedAccountID, capturedPartnerID, capturedBaseURL string

	clientMock := &mocks.SiigoClientMock{
		TestAuthenticationFn: func(_ context.Context, username, accessKey, accountID, partnerID, baseURL string) error {
			capturedUsername = username
			capturedAccessKey = accessKey
			capturedAccountID = accountID
			capturedPartnerID = partnerID
			capturedBaseURL = baseURL
			return nil
		},
	}

	credenciales := map[string]any{
		"username":   "test@siigo.com",
		"access_key": "mi-access-key",
		"account_id": "mi-account-id",
		"partner_id": "mi-partner-id",
	}

	uc := New(clientMock, &mocks.IntegrationCoreMock{}, nil, nil, &mocks.LoggerMock{})

	if err := uc.TestConnection(ctx, configProd(), credenciales); err != nil {
		t.Fatalf("se esperaba exito, se obtuvo error: %v", err)
	}
	if capturedUsername != "test@siigo.com" {
		t.Errorf("username esperado %q, se paso %q", "test@siigo.com", capturedUsername)
	}
	if capturedAccessKey != "mi-access-key" {
		t.Errorf("access_key esperado %q, se paso %q", "mi-access-key", capturedAccessKey)
	}
	if capturedAccountID != "mi-account-id" {
		t.Errorf("account_id esperado %q, se paso %q", "mi-account-id", capturedAccountID)
	}
	if capturedPartnerID != "mi-partner-id" {
		t.Errorf("partner_id esperado %q, se paso %q", "mi-partner-id", capturedPartnerID)
	}
	if capturedBaseURL != "https://api.siigo.com" {
		t.Errorf("baseURL esperado %q, se paso %q", "https://api.siigo.com", capturedBaseURL)
	}
}

func TestTestConnection_IsTestingUsaBaseURLTest(t *testing.T) {
	ctx := context.Background()

	var capturedBaseURL string
	clientMock := &mocks.SiigoClientMock{
		TestAuthenticationFn: func(_ context.Context, _, _, _, _, baseURL string) error {
			capturedBaseURL = baseURL
			return nil
		},
	}

	config := map[string]any{
		"base_url":      "https://api.siigo.com",
		"base_url_test": "http://back-testing:9095",
		"is_testing":    true,
	}

	uc := New(clientMock, &mocks.IntegrationCoreMock{}, nil, nil, &mocks.LoggerMock{})

	if err := uc.TestConnection(ctx, config, credencialesCompletas()); err != nil {
		t.Fatalf("se esperaba exito, se obtuvo error: %v", err)
	}
	if capturedBaseURL != "http://back-testing:9095" {
		t.Errorf("con is_testing=true se esperaba base_url_test, se paso %q", capturedBaseURL)
	}
}

func TestTestConnection_AccountIDOpcional(t *testing.T) {
	ctx := context.Background()

	llamado := false
	clientMock := &mocks.SiigoClientMock{
		TestAuthenticationFn: func(_ context.Context, _, _, _, _, _ string) error {
			llamado = true
			return nil
		},
	}

	credenciales := map[string]any{
		"username":   "test@siigo.com",
		"access_key": "mi-access-key",
		"partner_id": "mi-partner-id",
	}

	uc := New(clientMock, &mocks.IntegrationCoreMock{}, nil, nil, &mocks.LoggerMock{})

	if err := uc.TestConnection(ctx, configProd(), credenciales); err != nil {
		t.Errorf("se esperaba exito (account_id es opcional), se obtuvo error: %v", err)
	}
	if !llamado {
		t.Error("se esperaba que el cliente fuera llamado")
	}
}

func TestTestConnection_APIURLOverrideEnCredenciales(t *testing.T) {
	ctx := context.Background()

	var capturedBaseURL string
	clientMock := &mocks.SiigoClientMock{
		TestAuthenticationFn: func(_ context.Context, _, _, _, _, baseURL string) error {
			capturedBaseURL = baseURL
			return nil
		},
	}

	credenciales := map[string]any{
		"username":   "test@siigo.com",
		"access_key": "mi-access-key",
		"partner_id": "mi-partner-id",
		"api_url":    "https://override.siigo.com",
	}

	uc := New(clientMock, &mocks.IntegrationCoreMock{}, nil, nil, &mocks.LoggerMock{})

	if err := uc.TestConnection(ctx, configProd(), credenciales); err != nil {
		t.Fatalf("se esperaba exito, se obtuvo error: %v", err)
	}
	if capturedBaseURL != "https://override.siigo.com" {
		t.Errorf("api_url en credentials debe sobrescribir base_url, se paso %q", capturedBaseURL)
	}
}

func TestTestConnection_SinURLConfiguradaFalla(t *testing.T) {
	ctx := context.Background()

	clientMock := &mocks.SiigoClientMock{
		TestAuthenticationFn: func(_ context.Context, _, _, _, _, _ string) error {
			t.Error("el cliente no debe ser llamado si no hay URL")
			return nil
		},
	}

	uc := New(clientMock, &mocks.IntegrationCoreMock{}, nil, nil, &mocks.LoggerMock{})

	err := uc.TestConnection(ctx, nil, credencialesCompletas())
	if err == nil {
		t.Fatal("se esperaba error cuando no hay URL configurada")
	}
	if !strings.Contains(err.Error(), "URL de Siigo no configurada") {
		t.Errorf("mensaje de error inesperado: %v", err)
	}
}

func TestTestConnection_ClienteRetornaError(t *testing.T) {
	ctx := context.Background()
	errEsperado := errors.New("siigo authentication failed: invalid credentials")

	clientMock := &mocks.SiigoClientMock{
		TestAuthenticationFn: func(_ context.Context, _, _, _, _, _ string) error {
			return errEsperado
		},
	}

	uc := New(clientMock, &mocks.IntegrationCoreMock{}, nil, nil, &mocks.LoggerMock{})

	err := uc.TestConnection(ctx, configProd(), credencialesCompletas())
	if err == nil {
		t.Fatal("se esperaba un error del cliente, se obtuvo nil")
	}
	if !errors.Is(err, errEsperado) {
		t.Errorf("error esperado %v, se obtuvo %v", errEsperado, err)
	}
}

func TestTestConnection_ValidacionDeCamposRequeridos(t *testing.T) {
	tests := []struct {
		nombre       string
		credenciales map[string]any
		errorEspMsg  string
	}{
		{
			nombre: "falta username",
			credenciales: map[string]any{
				"access_key": "clave",
				"partner_id": "partner",
			},
			errorEspMsg: "username",
		},
		{
			nombre: "username vacio",
			credenciales: map[string]any{
				"username":   "",
				"access_key": "clave",
				"partner_id": "partner",
			},
			errorEspMsg: "username",
		},
		{
			nombre: "falta access_key",
			credenciales: map[string]any{
				"username":   "user@test.com",
				"partner_id": "partner",
			},
			errorEspMsg: "access_key",
		},
		{
			nombre: "access_key vacio",
			credenciales: map[string]any{
				"username":   "user@test.com",
				"access_key": "",
				"partner_id": "partner",
			},
			errorEspMsg: "access_key",
		},
		{
			nombre: "falta partner_id",
			credenciales: map[string]any{
				"username":   "user@test.com",
				"access_key": "clave",
			},
			errorEspMsg: "partner_id",
		},
		{
			nombre: "partner_id vacio",
			credenciales: map[string]any{
				"username":   "user@test.com",
				"access_key": "clave",
				"partner_id": "",
			},
			errorEspMsg: "partner_id",
		},
		{
			nombre:       "mapa de credenciales vacio",
			credenciales: map[string]any{},
			errorEspMsg:  "username",
		},
	}

	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			ctx := context.Background()
			clienteLlamado := false
			mockConCaptura := &mocks.SiigoClientMock{
				TestAuthenticationFn: func(_ context.Context, _, _, _, _, _ string) error {
					clienteLlamado = true
					return nil
				},
			}

			uc := New(mockConCaptura, &mocks.IntegrationCoreMock{}, nil, nil, &mocks.LoggerMock{})

			err := uc.TestConnection(ctx, configProd(), tt.credenciales)
			if err == nil {
				t.Fatalf("[%s] se esperaba error de validacion, se obtuvo nil", tt.nombre)
			}
			if tt.errorEspMsg != "" && !strings.Contains(err.Error(), tt.errorEspMsg) {
				t.Errorf("[%s] error esperado conteniendo %q, se obtuvo: %q", tt.nombre, tt.errorEspMsg, err.Error())
			}
			if clienteLlamado {
				t.Errorf("[%s] el cliente no deberia ser llamado cuando la validacion falla", tt.nombre)
			}
		})
	}
}

func TestTestConnection_TipoIncorrectoEnCredenciales(t *testing.T) {
	ctx := context.Background()
	uc := New(&mocks.SiigoClientMock{}, &mocks.IntegrationCoreMock{}, nil, nil, &mocks.LoggerMock{})

	credenciales := map[string]any{
		"username":   12345,
		"access_key": "clave",
		"partner_id": "partner",
	}

	if err := uc.TestConnection(ctx, configProd(), credenciales); err == nil {
		t.Fatal("se esperaba error por tipo incorrecto en username, se obtuvo nil")
	}
}
