package usecases

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce/internal/mocks"
)

// ─────────────────────────────────────────────────────────────
// SyncOrdersWithParams — camino feliz
// ─────────────────────────────────────────────────────────────

func TestSyncOrdersWithParams_Success(t *testing.T) {
	ctx := context.Background()

	serviceMock := &mocks.IntegrationServiceMock{
		GetIntegrationByIDFn: func(_ context.Context, _ string) (*domain.Integration, error) {
			return &domain.Integration{
				ID:   1,
				Name: "Test WooCommerce",
				Config: map[string]interface{}{
					"store_url": "https://mitienda.com",
				},
			}, nil
		},
	}

	uc := New(&mocks.WooClientMock{}, serviceMock, &mocks.OrderPublisherMock{}, mocks.NewLoggerMock())

	params := map[string]interface{}{
		"created_at_min": time.Now().AddDate(0, 0, -7).Format(time.RFC3339),
	}

	err := uc.SyncOrdersWithParams(ctx, "integration-1", params)

	if err != nil {
		t.Fatalf("esperaba sin error, recibí: %v", err)
	}
}

// ─────────────────────────────────────────────────────────────
// SyncOrdersWithParams — integración no encontrada
// ─────────────────────────────────────────────────────────────

func TestSyncOrdersWithParams_IntegrationNotFound(t *testing.T) {
	ctx := context.Background()

	serviceMock := &mocks.IntegrationServiceMock{
		GetIntegrationByIDFn: func(_ context.Context, _ string) (*domain.Integration, error) {
			return nil, nil
		},
	}

	uc := New(&mocks.WooClientMock{}, serviceMock, &mocks.OrderPublisherMock{}, mocks.NewLoggerMock())

	err := uc.SyncOrdersWithParams(ctx, "nonexistent", map[string]interface{}{})

	if err == nil {
		t.Fatal("esperaba error por integración no encontrada, recibí nil")
	}
	if !errors.Is(err, domain.ErrIntegrationNotFound) {
		t.Errorf("esperaba ErrIntegrationNotFound, recibí: %v", err)
	}
}

// ─────────────────────────────────────────────────────────────
// SyncOrdersWithParams — error al obtener integración
// ─────────────────────────────────────────────────────────────

func TestSyncOrdersWithParams_GetIntegrationError(t *testing.T) {
	ctx := context.Background()
	dbErr := errors.New("database: connection refused")

	serviceMock := &mocks.IntegrationServiceMock{
		GetIntegrationByIDFn: func(_ context.Context, _ string) (*domain.Integration, error) {
			return nil, dbErr
		},
	}

	uc := New(&mocks.WooClientMock{}, serviceMock, &mocks.OrderPublisherMock{}, mocks.NewLoggerMock())

	err := uc.SyncOrdersWithParams(ctx, "integration-1", map[string]interface{}{})

	if err == nil {
		t.Fatal("esperaba error de base de datos, recibí nil")
	}
	if !errors.Is(err, dbErr) {
		t.Errorf("esperaba error wrapping '%v', recibí '%v'", dbErr, err)
	}
}

// ─────────────────────────────────────────────────────────────
// SyncOrdersWithParams — store_url faltante en config
// ─────────────────────────────────────────────────────────────

func TestSyncOrdersWithParams_MissingStoreURL(t *testing.T) {
	ctx := context.Background()

	serviceMock := &mocks.IntegrationServiceMock{
		GetIntegrationByIDFn: func(_ context.Context, _ string) (*domain.Integration, error) {
			return &domain.Integration{
				ID:     1,
				Config: map[string]interface{}{
					// "store_url" ausente
				},
			}, nil
		},
	}

	uc := New(&mocks.WooClientMock{}, serviceMock, &mocks.OrderPublisherMock{}, mocks.NewLoggerMock())

	err := uc.SyncOrdersWithParams(ctx, "integration-1", map[string]interface{}{})

	if err == nil {
		t.Fatal("esperaba error por store_url faltante, recibí nil")
	}
	if !errors.Is(err, domain.ErrMissingStoreURL) {
		t.Errorf("esperaba ErrMissingStoreURL, recibí: %v", err)
	}
}

// ─────────────────────────────────────────────────────────────
// SyncOrdersWithParams — error al descifrar consumer_key
// ─────────────────────────────────────────────────────────────

func TestSyncOrdersWithParams_DecryptConsumerKeyError(t *testing.T) {
	ctx := context.Background()
	decryptErr := errors.New("crypto: decryption failed")

	serviceMock := &mocks.IntegrationServiceMock{
		GetIntegrationByIDFn: func(_ context.Context, _ string) (*domain.Integration, error) {
			return &domain.Integration{
				ID:   1,
				Config: map[string]interface{}{
					"store_url": "https://mitienda.com",
				},
			}, nil
		},
		DecryptCredentialFn: func(_ context.Context, _ string, fieldName string) (string, error) {
			if fieldName == "consumer_key" {
				return "", decryptErr
			}
			return "cs_test", nil
		},
	}

	uc := New(&mocks.WooClientMock{}, serviceMock, &mocks.OrderPublisherMock{}, mocks.NewLoggerMock())

	err := uc.SyncOrdersWithParams(ctx, "integration-1", map[string]interface{}{})

	if err == nil {
		t.Fatal("esperaba error de descifrado, recibí nil")
	}
	if !errors.Is(err, decryptErr) {
		t.Errorf("esperaba error wrapping '%v', recibí '%v'", decryptErr, err)
	}
}

// ─────────────────────────────────────────────────────────────
// SyncOrdersWithParams — error al descifrar consumer_secret
// ─────────────────────────────────────────────────────────────

func TestSyncOrdersWithParams_DecryptConsumerSecretError(t *testing.T) {
	ctx := context.Background()
	decryptErr := errors.New("crypto: decryption failed")

	serviceMock := &mocks.IntegrationServiceMock{
		GetIntegrationByIDFn: func(_ context.Context, _ string) (*domain.Integration, error) {
			return &domain.Integration{
				ID:   1,
				Config: map[string]interface{}{
					"store_url": "https://mitienda.com",
				},
			}, nil
		},
		DecryptCredentialFn: func(_ context.Context, _ string, fieldName string) (string, error) {
			if fieldName == "consumer_secret" {
				return "", decryptErr
			}
			return "ck_test", nil
		},
	}

	uc := New(&mocks.WooClientMock{}, serviceMock, &mocks.OrderPublisherMock{}, mocks.NewLoggerMock())

	err := uc.SyncOrdersWithParams(ctx, "integration-1", map[string]interface{}{})

	if err == nil {
		t.Fatal("esperaba error de descifrado, recibí nil")
	}
	if !errors.Is(err, decryptErr) {
		t.Errorf("esperaba error wrapping '%v', recibí '%v'", decryptErr, err)
	}
}

// ─────────────────────────────────────────────────────────────
// buildQueryParams — verificar parseo de parámetros
// ─────────────────────────────────────────────────────────────

func TestBuildQueryParams_WithAllParams(t *testing.T) {
	minDate := "2026-01-01T00:00:00Z"
	maxDate := "2026-01-31T23:59:59Z"
	params := map[string]interface{}{
		"created_at_min": minDate,
		"created_at_max": maxDate,
		"status":         "processing",
	}

	qp := buildQueryParams(params)

	if qp.After == nil {
		t.Fatal("After no debe ser nil")
	}
	if qp.Before == nil {
		t.Fatal("Before no debe ser nil")
	}
	if qp.Status != "processing" {
		t.Errorf("Status esperado 'processing', recibí '%s'", qp.Status)
	}
	if qp.OrderBy != "date" {
		t.Errorf("OrderBy por defecto esperado 'date', recibí '%s'", qp.OrderBy)
	}
	if qp.Order != "desc" {
		t.Errorf("Order por defecto esperado 'desc', recibí '%s'", qp.Order)
	}
}

func TestBuildQueryParams_EmptyMap(t *testing.T) {
	qp := buildQueryParams(map[string]interface{}{})

	if qp.After != nil {
		t.Errorf("After debe ser nil para mapa vacío")
	}
	if qp.Before != nil {
		t.Errorf("Before debe ser nil para mapa vacío")
	}
	if qp.Status != "" {
		t.Errorf("Status debe ser vacío para mapa vacío")
	}
}

func TestBuildQueryParams_NotAMap(t *testing.T) {
	qp := buildQueryParams("not a map")

	if qp.OrderBy != "date" {
		t.Errorf("OrderBy por defecto esperado 'date', recibí '%s'", qp.OrderBy)
	}
}

func TestBuildQueryParams_InvalidDateFormat(t *testing.T) {
	params := map[string]interface{}{
		"created_at_min": "invalid-date",
	}

	qp := buildQueryParams(params)

	if qp.After != nil {
		t.Errorf("After debe ser nil con fecha inválida")
	}
}

// ─────────────────────────────────────────────────────────────
// SyncOrders — verifica que use los últimos 30 días
// ─────────────────────────────────────────────────────────────

func TestSyncOrders_Uses30DaysDefault(t *testing.T) {
	ctx := context.Background()

	serviceMock := &mocks.IntegrationServiceMock{
		GetIntegrationByIDFn: func(_ context.Context, _ string) (*domain.Integration, error) {
			return &domain.Integration{
				ID:   1,
				Config: map[string]interface{}{
					"store_url": "https://mitienda.com",
				},
			}, nil
		},
	}

	uc := New(&mocks.WooClientMock{}, serviceMock, &mocks.OrderPublisherMock{}, mocks.NewLoggerMock())

	err := uc.SyncOrders(ctx, "integration-1")

	if err != nil {
		t.Fatalf("esperaba sin error, recibí: %v", err)
	}
}
