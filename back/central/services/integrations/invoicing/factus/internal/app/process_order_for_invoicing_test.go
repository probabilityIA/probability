package app

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/factus/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/factus/internal/mocks"
)

// ─────────────────────────────────────────────────────────────
// Helpers para construir fixtures reutilizables
// ─────────────────────────────────────────────────────────────

func buildProcessInvoiceRequest() *dtos.ProcessInvoiceRequest {
	return &dtos.ProcessInvoiceRequest{
		InvoiceID:     42,
		Operation:     "create",
		CorrelationID: "corr-abc-123",
		IntegrationID: 7,
		OrderID:       "order-999",
		Total:         120000,
		Subtotal:      100000,
		Tax:           20000,
		Discount:      0,
		ShippingCost:  5000,
		Currency:      "COP",
		Customer: dtos.CustomerData{
			Name:  "Juan Perez",
			Email: "juan@example.com",
			DNI:   "123456789",
		},
		Items: []dtos.ItemData{
			{
				SKU:        "SKU-01",
				Name:       "Producto A",
				Quantity:   2,
				UnitPrice:  50000,
				TotalPrice: 100000,
			},
		},
		Config: map[string]interface{}{
			"document_id": float64(1),
		},
	}
}

func buildPublicIntegration(config map[string]interface{}) *core.PublicIntegration {
	if config == nil {
		config = map[string]interface{}{}
	}
	return &core.PublicIntegration{
		ID:              7,
		IntegrationType: 7,
		Config:          config,
	}
}

func buildSuccessInvoiceResult() *dtos.CreateInvoiceResult {
	return &dtos.CreateInvoiceResult{
		InvoiceNumber: "SETP990000001",
		ExternalID:    "ext-001",
		CUFE:          "cufe-hash-abc",
		QRCode:        "qr-data",
		Total:         "120000.00",
		IssuedAt:      "2026-02-23T10:00:00Z",
		AuditData: &dtos.AuditData{
			RequestURL:     "https://api.factus.com.co/v1/bills/validate",
			ResponseStatus: 200,
			ResponseBody:   `{"status":"ok"}`,
		},
	}
}

// ─────────────────────────────────────────────────────────────
// CreateInvoice — camino feliz
// ─────────────────────────────────────────────────────────────

func TestCreateInvoice_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	req := buildProcessInvoiceRequest()

	integrationConfig := map[string]interface{}{"tax_id": float64(1)}
	expectedResult := buildSuccessInvoiceResult()

	mockCore := &mocks.IntegrationCoreMock{
		GetIntegrationByIDFn: func(_ context.Context, _ string) (*core.PublicIntegration, error) {
			return buildPublicIntegration(integrationConfig), nil
		},
		DecryptCredentialFn: func(_ context.Context, _ string, field string) (string, error) {
			m := map[string]string{
				"client_id":     "real-client-id",
				"client_secret": "real-client-secret",
				"username":      "factus@biz.com",
				"password":      "s3cr3t",
				"api_url":       "https://api.factus.com.co",
			}
			return m[field], nil
		},
	}

	mockClient := &mocks.FactusClientMock{
		CreateInvoiceFn: func(_ context.Context, clientReq *dtos.CreateInvoiceRequest) (*dtos.CreateInvoiceResult, error) {
			// Verificar que las credenciales se propagaron correctamente
			if clientReq.Credentials.ClientID != "real-client-id" {
				t.Errorf("esperaba client_id 'real-client-id', recibí '%s'", clientReq.Credentials.ClientID)
			}
			if clientReq.Credentials.Username != "factus@biz.com" {
				t.Errorf("esperaba username 'factus@biz.com', recibí '%s'", clientReq.Credentials.Username)
			}
			return expectedResult, nil
		},
	}

	uc := New(mockClient, mockCore, mocks.NewLoggerMock())

	// Act
	result, err := uc.CreateInvoice(ctx, req)

	// Assert
	if err != nil {
		t.Fatalf("esperaba sin error, recibí: %v", err)
	}
	if result == nil {
		t.Fatal("esperaba resultado no nulo")
	}
	if result.InvoiceNumber != expectedResult.InvoiceNumber {
		t.Errorf("InvoiceNumber: esperaba '%s', recibí '%s'", expectedResult.InvoiceNumber, result.InvoiceNumber)
	}
	if result.CUFE != expectedResult.CUFE {
		t.Errorf("CUFE: esperaba '%s', recibí '%s'", expectedResult.CUFE, result.CUFE)
	}
	if result.AuditData == nil {
		t.Error("esperaba AuditData en el resultado")
	}
}

// ─────────────────────────────────────────────────────────────
// CreateInvoice — combina config de integración con config del request
// ─────────────────────────────────────────────────────────────

func TestCreateInvoice_ConfigMerge_RequestOverridesIntegration(t *testing.T) {
	// Arrange
	ctx := context.Background()
	req := buildProcessInvoiceRequest()
	// El request tiene document_id = 1, la integración tiene document_id = 99
	// El request debe ganar (prioridad más alta)
	req.Config = map[string]interface{}{"document_id": float64(1), "extra": "from_request"}

	integrationConfig := map[string]interface{}{
		"document_id":       float64(99), // debe ser sobreescrito por el request
		"integration_field": "base_value",
	}

	var capturedConfig map[string]interface{}

	mockCore := &mocks.IntegrationCoreMock{
		GetIntegrationByIDFn: func(_ context.Context, _ string) (*core.PublicIntegration, error) {
			return buildPublicIntegration(integrationConfig), nil
		},
	}

	mockClient := &mocks.FactusClientMock{
		CreateInvoiceFn: func(_ context.Context, clientReq *dtos.CreateInvoiceRequest) (*dtos.CreateInvoiceResult, error) {
			capturedConfig = clientReq.Config
			return buildSuccessInvoiceResult(), nil
		},
	}

	uc := New(mockClient, mockCore, mocks.NewLoggerMock())

	// Act
	_, err := uc.CreateInvoice(ctx, req)

	// Assert
	if err != nil {
		t.Fatalf("esperaba sin error, recibí: %v", err)
	}
	if capturedConfig == nil {
		t.Fatal("el config enviado al cliente no debe ser nulo")
	}
	// El campo del request debe prevalecer sobre el de la integración
	if capturedConfig["document_id"] != float64(1) {
		t.Errorf("document_id: esperaba 1, recibí %v", capturedConfig["document_id"])
	}
	// El campo exclusivo de la integración debe estar presente
	if capturedConfig["integration_field"] != "base_value" {
		t.Errorf("integration_field: esperaba 'base_value', recibí %v", capturedConfig["integration_field"])
	}
	// El campo exclusivo del request también debe estar presente
	if capturedConfig["extra"] != "from_request" {
		t.Errorf("extra: esperaba 'from_request', recibí %v", capturedConfig["extra"])
	}
}

// ─────────────────────────────────────────────────────────────
// CreateInvoice — error al obtener integración
// ─────────────────────────────────────────────────────────────

func TestCreateInvoice_GetIntegrationError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	req := buildProcessInvoiceRequest()
	expectedErr := errors.New("integration not found")

	mockCore := &mocks.IntegrationCoreMock{
		GetIntegrationByIDFn: func(_ context.Context, _ string) (*core.PublicIntegration, error) {
			return nil, expectedErr
		},
	}
	mockClient := &mocks.FactusClientMock{}

	uc := New(mockClient, mockCore, mocks.NewLoggerMock())

	// Act
	result, err := uc.CreateInvoice(ctx, req)

	// Assert
	if err == nil {
		t.Fatal("esperaba error, recibí nil")
	}
	if !errors.Is(err, expectedErr) {
		t.Errorf("esperaba error '%v', recibí '%v'", expectedErr, err)
	}
	if result != nil {
		t.Errorf("esperaba resultado nulo ante error de integración, recibí %v", result)
	}
}

// ─────────────────────────────────────────────────────────────
// CreateInvoice — error al descifrar client_id
// ─────────────────────────────────────────────────────────────

func TestCreateInvoice_DecryptClientIDError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	req := buildProcessInvoiceRequest()
	decryptErr := errors.New("encryption key not found")

	mockCore := &mocks.IntegrationCoreMock{
		GetIntegrationByIDFn: func(_ context.Context, _ string) (*core.PublicIntegration, error) {
			return buildPublicIntegration(nil), nil
		},
		DecryptCredentialFn: func(_ context.Context, _ string, field string) (string, error) {
			if field == "client_id" {
				return "", decryptErr
			}
			return "valor", nil
		},
	}
	mockClient := &mocks.FactusClientMock{}

	uc := New(mockClient, mockCore, mocks.NewLoggerMock())

	// Act
	result, err := uc.CreateInvoice(ctx, req)

	// Assert
	if err == nil {
		t.Fatal("esperaba error, recibí nil")
	}
	if !errors.Is(err, decryptErr) {
		t.Errorf("esperaba error de descifrado, recibí '%v'", err)
	}
	if result != nil {
		t.Errorf("esperaba resultado nulo, recibí %v", result)
	}
}

// ─────────────────────────────────────────────────────────────
// CreateInvoice — error al descifrar client_secret
// ─────────────────────────────────────────────────────────────

func TestCreateInvoice_DecryptClientSecretError(t *testing.T) {
	ctx := context.Background()
	req := buildProcessInvoiceRequest()
	decryptErr := errors.New("could not decrypt client_secret")

	mockCore := &mocks.IntegrationCoreMock{
		GetIntegrationByIDFn: func(_ context.Context, _ string) (*core.PublicIntegration, error) {
			return buildPublicIntegration(nil), nil
		},
		DecryptCredentialFn: func(_ context.Context, _ string, field string) (string, error) {
			if field == "client_secret" {
				return "", decryptErr
			}
			return "valor", nil
		},
	}
	mockClient := &mocks.FactusClientMock{}

	uc := New(mockClient, mockCore, mocks.NewLoggerMock())

	result, err := uc.CreateInvoice(ctx, req)

	if err == nil {
		t.Fatal("esperaba error, recibí nil")
	}
	if !errors.Is(err, decryptErr) {
		t.Errorf("esperaba error '%v', recibí '%v'", decryptErr, err)
	}
	if result != nil {
		t.Errorf("esperaba resultado nulo, recibí %v", result)
	}
}

// ─────────────────────────────────────────────────────────────
// CreateInvoice — error al descifrar username
// ─────────────────────────────────────────────────────────────

func TestCreateInvoice_DecryptUsernameError(t *testing.T) {
	ctx := context.Background()
	req := buildProcessInvoiceRequest()
	decryptErr := errors.New("could not decrypt username")

	mockCore := &mocks.IntegrationCoreMock{
		GetIntegrationByIDFn: func(_ context.Context, _ string) (*core.PublicIntegration, error) {
			return buildPublicIntegration(nil), nil
		},
		DecryptCredentialFn: func(_ context.Context, _ string, field string) (string, error) {
			if field == "username" {
				return "", decryptErr
			}
			return "valor", nil
		},
	}
	mockClient := &mocks.FactusClientMock{}

	uc := New(mockClient, mockCore, mocks.NewLoggerMock())

	result, err := uc.CreateInvoice(ctx, req)

	if err == nil {
		t.Fatal("esperaba error, recibí nil")
	}
	if !errors.Is(err, decryptErr) {
		t.Errorf("esperaba error '%v', recibí '%v'", decryptErr, err)
	}
	if result != nil {
		t.Errorf("esperaba resultado nulo, recibí %v", result)
	}
}

// ─────────────────────────────────────────────────────────────
// CreateInvoice — error al descifrar password
// ─────────────────────────────────────────────────────────────

func TestCreateInvoice_DecryptPasswordError(t *testing.T) {
	ctx := context.Background()
	req := buildProcessInvoiceRequest()
	decryptErr := errors.New("could not decrypt password")

	mockCore := &mocks.IntegrationCoreMock{
		GetIntegrationByIDFn: func(_ context.Context, _ string) (*core.PublicIntegration, error) {
			return buildPublicIntegration(nil), nil
		},
		DecryptCredentialFn: func(_ context.Context, _ string, field string) (string, error) {
			if field == "password" {
				return "", decryptErr
			}
			return "valor", nil
		},
	}
	mockClient := &mocks.FactusClientMock{}

	uc := New(mockClient, mockCore, mocks.NewLoggerMock())

	result, err := uc.CreateInvoice(ctx, req)

	if err == nil {
		t.Fatal("esperaba error, recibí nil")
	}
	if !errors.Is(err, decryptErr) {
		t.Errorf("esperaba error '%v', recibí '%v'", decryptErr, err)
	}
	if result != nil {
		t.Errorf("esperaba resultado nulo, recibí %v", result)
	}
}

// ─────────────────────────────────────────────────────────────
// CreateInvoice — el cliente HTTP retorna error pero sí retorna AuditData
// ─────────────────────────────────────────────────────────────

func TestCreateInvoice_ClientError_PropagatesAuditData(t *testing.T) {
	// Arrange
	ctx := context.Background()
	req := buildProcessInvoiceRequest()
	apiErr := errors.New("factus API returned 422 Unprocessable Entity")

	auditData := &dtos.AuditData{
		RequestURL:     "https://api.factus.com.co/v1/bills/validate",
		ResponseStatus: 422,
		ResponseBody:   `{"errors":["campo requerido faltante"]}`,
	}

	mockCore := &mocks.IntegrationCoreMock{
		GetIntegrationByIDFn: func(_ context.Context, _ string) (*core.PublicIntegration, error) {
			return buildPublicIntegration(nil), nil
		},
	}

	mockClient := &mocks.FactusClientMock{
		CreateInvoiceFn: func(_ context.Context, _ *dtos.CreateInvoiceRequest) (*dtos.CreateInvoiceResult, error) {
			// El cliente retorna resultado parcial (con AuditData) Y error, como documenta el contrato
			return &dtos.CreateInvoiceResult{AuditData: auditData}, apiErr
		},
	}

	uc := New(mockClient, mockCore, mocks.NewLoggerMock())

	// Act
	result, err := uc.CreateInvoice(ctx, req)

	// Assert
	if err == nil {
		t.Fatal("esperaba error, recibí nil")
	}
	if !errors.Is(err, apiErr) {
		t.Errorf("esperaba error '%v', recibí '%v'", apiErr, err)
	}
	// El contrato del use case garantiza que SIEMPRE retorna resultado (incluso en error)
	// para propagar el AuditData hacia el consumer
	if result == nil {
		t.Fatal("esperaba resultado no nulo incluso en caso de error (para propagar AuditData)")
	}
	if result.AuditData == nil {
		t.Error("esperaba AuditData en resultado de error")
	}
	if result.AuditData.ResponseStatus != 422 {
		t.Errorf("AuditData.ResponseStatus: esperaba 422, recibí %d", result.AuditData.ResponseStatus)
	}
}

// ─────────────────────────────────────────────────────────────
// CreateInvoice — el cliente retorna error SIN AuditData
// ─────────────────────────────────────────────────────────────

func TestCreateInvoice_ClientError_NilResult_ReturnsEmptyResult(t *testing.T) {
	// Arrange
	ctx := context.Background()
	req := buildProcessInvoiceRequest()
	apiErr := errors.New("network timeout")

	mockCore := &mocks.IntegrationCoreMock{
		GetIntegrationByIDFn: func(_ context.Context, _ string) (*core.PublicIntegration, error) {
			return buildPublicIntegration(nil), nil
		},
	}

	mockClient := &mocks.FactusClientMock{
		CreateInvoiceFn: func(_ context.Context, _ *dtos.CreateInvoiceRequest) (*dtos.CreateInvoiceResult, error) {
			// El cliente retorna nil en el resultado (caso de timeout / error de red)
			return nil, apiErr
		},
	}

	uc := New(mockClient, mockCore, mocks.NewLoggerMock())

	// Act
	result, err := uc.CreateInvoice(ctx, req)

	// Assert
	if err == nil {
		t.Fatal("esperaba error, recibí nil")
	}
	if !errors.Is(err, apiErr) {
		t.Errorf("esperaba error '%v', recibí '%v'", apiErr, err)
	}
	// El use case debe retornar un result vacío (no nil) para que el consumer no panic
	if result == nil {
		t.Fatal("el use case debe retornar resultado vacío (no nil) incluso cuando el cliente retorna nil")
	}
}

// ─────────────────────────────────────────────────────────────
// CreateInvoice — el mensaje de error debe incluir el integrationID
// ─────────────────────────────────────────────────────────────

func TestCreateInvoice_ErrorMessage_ContainsIntegrationID(t *testing.T) {
	ctx := context.Background()
	req := buildProcessInvoiceRequest()
	req.IntegrationID = 99

	mockCore := &mocks.IntegrationCoreMock{
		GetIntegrationByIDFn: func(_ context.Context, _ string) (*core.PublicIntegration, error) {
			return nil, errors.New("not found")
		},
	}
	mockClient := &mocks.FactusClientMock{}

	uc := New(mockClient, mockCore, mocks.NewLoggerMock())

	_, err := uc.CreateInvoice(ctx, req)

	if err == nil {
		t.Fatal("esperaba error, recibí nil")
	}
	if !strings.Contains(err.Error(), "99") {
		t.Errorf("el mensaje de error debe contener el integration_id '99', recibí: %s", err.Error())
	}
}

// ─────────────────────────────────────────────────────────────
// CreateInvoice — api_url es opcional (puede estar vacía)
// ─────────────────────────────────────────────────────────────

func TestCreateInvoice_ApiURLIsOptional(t *testing.T) {
	// Arrange
	ctx := context.Background()
	req := buildProcessInvoiceRequest()

	var capturedBaseURL string

	mockCore := &mocks.IntegrationCoreMock{
		GetIntegrationByIDFn: func(_ context.Context, _ string) (*core.PublicIntegration, error) {
			return buildPublicIntegration(nil), nil
		},
		DecryptCredentialFn: func(_ context.Context, _ string, field string) (string, error) {
			// api_url no existe en las credenciales → error ignorado por el use case
			if field == "api_url" {
				return "", errors.New("field not found")
			}
			defaults := map[string]string{
				"client_id":     "cid",
				"client_secret": "csecret",
				"username":      "user@test.com",
				"password":      "pass",
			}
			return defaults[field], nil
		},
	}

	mockClient := &mocks.FactusClientMock{
		CreateInvoiceFn: func(_ context.Context, clientReq *dtos.CreateInvoiceRequest) (*dtos.CreateInvoiceResult, error) {
			capturedBaseURL = clientReq.Credentials.BaseURL
			return buildSuccessInvoiceResult(), nil
		},
	}

	uc := New(mockClient, mockCore, mocks.NewLoggerMock())

	// Act
	result, err := uc.CreateInvoice(ctx, req)

	// Assert
	if err != nil {
		t.Fatalf("el error de api_url no debe propagarse, recibí: %v", err)
	}
	if result == nil {
		t.Fatal("esperaba resultado no nulo")
	}
	// La URL base debe quedar vacía cuando no hay credencial
	if capturedBaseURL != "" {
		t.Errorf("esperaba BaseURL vacía, recibí '%s'", capturedBaseURL)
	}
}

// ─────────────────────────────────────────────────────────────
// CreateInvoice — tabla de escenarios (table-driven)
// ─────────────────────────────────────────────────────────────

func TestCreateInvoice_TableDriven_DecryptErrors(t *testing.T) {
	credentialFields := []string{"client_id", "client_secret", "username", "password"}

	for _, failField := range credentialFields {
		failField := failField // captura para closure
		t.Run("falla descifrado de "+failField, func(t *testing.T) {
			ctx := context.Background()
			req := buildProcessInvoiceRequest()
			decryptErr := errors.New("error descifrado " + failField)

			mockCore := &mocks.IntegrationCoreMock{
				GetIntegrationByIDFn: func(_ context.Context, _ string) (*core.PublicIntegration, error) {
					return buildPublicIntegration(nil), nil
				},
				DecryptCredentialFn: func(_ context.Context, _ string, field string) (string, error) {
					if field == failField {
						return "", decryptErr
					}
					return "valor-ok", nil
				},
			}
			mockClient := &mocks.FactusClientMock{}

			uc := New(mockClient, mockCore, mocks.NewLoggerMock())

			result, err := uc.CreateInvoice(ctx, req)

			if err == nil {
				t.Fatalf("esperaba error al fallar descifrado de '%s'", failField)
			}
			if !errors.Is(err, decryptErr) {
				t.Errorf("esperaba error '%v', recibí '%v'", decryptErr, err)
			}
			if result != nil {
				t.Errorf("esperaba resultado nulo, recibí %v", result)
			}
		})
	}
}
