package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce/internal/mocks"
)

// ─────────────────────────────────────────────────────────────
// TestConnection — camino feliz con todas las credenciales
// ─────────────────────────────────────────────────────────────

func TestTestConnection_Success(t *testing.T) {
	ctx := context.Background()
	config := map[string]interface{}{
		"store_url": "https://mitienda.com",
	}
	credentials := map[string]interface{}{
		"consumer_key":    "ck_xxxxxxxxxxxxxxxxxxxx",
		"consumer_secret": "cs_xxxxxxxxxxxxxxxxxxxx",
	}

	mockClient := &mocks.WooClientMock{
		TestConnectionFn: func(_ context.Context, storeURL, consumerKey, consumerSecret string) error {
			if storeURL != "https://mitienda.com" {
				t.Errorf("storeURL inesperada: '%s'", storeURL)
			}
			if consumerKey != "ck_xxxxxxxxxxxxxxxxxxxx" {
				t.Errorf("consumerKey inesperado: '%s'", consumerKey)
			}
			if consumerSecret != "cs_xxxxxxxxxxxxxxxxxxxx" {
				t.Errorf("consumerSecret inesperado: '%s'", consumerSecret)
			}
			return nil
		},
	}

	uc := New(mockClient, &mocks.IntegrationServiceMock{}, &mocks.OrderPublisherMock{}, mocks.NewLoggerMock())

	err := uc.TestConnection(ctx, config, credentials)

	if err != nil {
		t.Fatalf("esperaba sin error, recibí: %v", err)
	}
}

// ─────────────────────────────────────────────────────────────
// TestConnection — store_url ausente
// ─────────────────────────────────────────────────────────────

func TestTestConnection_MissingStoreURL(t *testing.T) {
	ctx := context.Background()
	config := map[string]interface{}{
		// "store_url" ausente
	}
	credentials := map[string]interface{}{
		"consumer_key":    "ck_test",
		"consumer_secret": "cs_test",
	}

	uc := New(&mocks.WooClientMock{}, &mocks.IntegrationServiceMock{}, &mocks.OrderPublisherMock{}, mocks.NewLoggerMock())

	err := uc.TestConnection(ctx, config, credentials)

	if err == nil {
		t.Fatal("esperaba error por store_url faltante, recibí nil")
	}
	if !errors.Is(err, domain.ErrMissingStoreURL) {
		t.Errorf("esperaba ErrMissingStoreURL, recibí: %v", err)
	}
}

// ─────────────────────────────────────────────────────────────
// TestConnection — store_url vacío
// ─────────────────────────────────────────────────────────────

func TestTestConnection_EmptyStoreURL(t *testing.T) {
	ctx := context.Background()
	config := map[string]interface{}{
		"store_url": "",
	}
	credentials := map[string]interface{}{
		"consumer_key":    "ck_test",
		"consumer_secret": "cs_test",
	}

	uc := New(&mocks.WooClientMock{}, &mocks.IntegrationServiceMock{}, &mocks.OrderPublisherMock{}, mocks.NewLoggerMock())

	err := uc.TestConnection(ctx, config, credentials)

	if err == nil {
		t.Fatal("esperaba error por store_url vacío, recibí nil")
	}
}

// ─────────────────────────────────────────────────────────────
// TestConnection — consumer_key ausente
// ─────────────────────────────────────────────────────────────

func TestTestConnection_MissingConsumerKey(t *testing.T) {
	ctx := context.Background()
	config := map[string]interface{}{
		"store_url": "https://mitienda.com",
	}
	credentials := map[string]interface{}{
		// "consumer_key" ausente
		"consumer_secret": "cs_test",
	}

	uc := New(&mocks.WooClientMock{}, &mocks.IntegrationServiceMock{}, &mocks.OrderPublisherMock{}, mocks.NewLoggerMock())

	err := uc.TestConnection(ctx, config, credentials)

	if err == nil {
		t.Fatal("esperaba error por consumer_key faltante, recibí nil")
	}
	if !errors.Is(err, domain.ErrMissingConsumerKey) {
		t.Errorf("esperaba ErrMissingConsumerKey, recibí: %v", err)
	}
}

// ─────────────────────────────────────────────────────────────
// TestConnection — consumer_key vacío
// ─────────────────────────────────────────────────────────────

func TestTestConnection_EmptyConsumerKey(t *testing.T) {
	ctx := context.Background()
	config := map[string]interface{}{
		"store_url": "https://mitienda.com",
	}
	credentials := map[string]interface{}{
		"consumer_key":    "",
		"consumer_secret": "cs_test",
	}

	uc := New(&mocks.WooClientMock{}, &mocks.IntegrationServiceMock{}, &mocks.OrderPublisherMock{}, mocks.NewLoggerMock())

	err := uc.TestConnection(ctx, config, credentials)

	if err == nil {
		t.Fatal("esperaba error por consumer_key vacío, recibí nil")
	}
}

// ─────────────────────────────────────────────────────────────
// TestConnection — consumer_secret ausente
// ─────────────────────────────────────────────────────────────

func TestTestConnection_MissingConsumerSecret(t *testing.T) {
	ctx := context.Background()
	config := map[string]interface{}{
		"store_url": "https://mitienda.com",
	}
	credentials := map[string]interface{}{
		"consumer_key": "ck_test",
		// "consumer_secret" ausente
	}

	uc := New(&mocks.WooClientMock{}, &mocks.IntegrationServiceMock{}, &mocks.OrderPublisherMock{}, mocks.NewLoggerMock())

	err := uc.TestConnection(ctx, config, credentials)

	if err == nil {
		t.Fatal("esperaba error por consumer_secret faltante, recibí nil")
	}
	if !errors.Is(err, domain.ErrMissingConsumerSecret) {
		t.Errorf("esperaba ErrMissingConsumerSecret, recibí: %v", err)
	}
}

// ─────────────────────────────────────────────────────────────
// TestConnection — consumer_secret vacío
// ─────────────────────────────────────────────────────────────

func TestTestConnection_EmptyConsumerSecret(t *testing.T) {
	ctx := context.Background()
	config := map[string]interface{}{
		"store_url": "https://mitienda.com",
	}
	credentials := map[string]interface{}{
		"consumer_key":    "ck_test",
		"consumer_secret": "",
	}

	uc := New(&mocks.WooClientMock{}, &mocks.IntegrationServiceMock{}, &mocks.OrderPublisherMock{}, mocks.NewLoggerMock())

	err := uc.TestConnection(ctx, config, credentials)

	if err == nil {
		t.Fatal("esperaba error por consumer_secret vacío, recibí nil")
	}
}

// ─────────────────────────────────────────────────────────────
// TestConnection — tabla de validaciones de campos requeridos
// ─────────────────────────────────────────────────────────────

func TestTestConnection_TableDriven_RequiredFields(t *testing.T) {
	baseConfig := map[string]interface{}{
		"store_url": "https://mitienda.com",
	}
	baseCredentials := map[string]interface{}{
		"consumer_key":    "ck_test",
		"consumer_secret": "cs_test",
	}

	tests := []struct {
		name       string
		configMod  func(map[string]interface{})
		credsMod   func(map[string]interface{})
	}{
		{"falta store_url", func(c map[string]interface{}) { delete(c, "store_url") }, nil},
		{"store_url vacío", func(c map[string]interface{}) { c["store_url"] = "" }, nil},
		{"falta consumer_key", nil, func(c map[string]interface{}) { delete(c, "consumer_key") }},
		{"consumer_key vacío", nil, func(c map[string]interface{}) { c["consumer_key"] = "" }},
		{"falta consumer_secret", nil, func(c map[string]interface{}) { delete(c, "consumer_secret") }},
		{"consumer_secret vacío", nil, func(c map[string]interface{}) { c["consumer_secret"] = "" }},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			config := make(map[string]interface{})
			for k, v := range baseConfig {
				config[k] = v
			}
			creds := make(map[string]interface{})
			for k, v := range baseCredentials {
				creds[k] = v
			}

			if tt.configMod != nil {
				tt.configMod(config)
			}
			if tt.credsMod != nil {
				tt.credsMod(creds)
			}

			uc := New(&mocks.WooClientMock{}, &mocks.IntegrationServiceMock{}, &mocks.OrderPublisherMock{}, mocks.NewLoggerMock())

			err := uc.TestConnection(context.Background(), config, creds)

			if err == nil {
				t.Fatalf("esperaba error para escenario '%s', recibí nil", tt.name)
			}
		})
	}
}

// ─────────────────────────────────────────────────────────────
// TestConnection — el cliente HTTP retorna error
// ─────────────────────────────────────────────────────────────

func TestTestConnection_ClientFails(t *testing.T) {
	ctx := context.Background()
	config := map[string]interface{}{
		"store_url": "https://mitienda.com",
	}
	credentials := map[string]interface{}{
		"consumer_key":    "ck_bad",
		"consumer_secret": "cs_bad",
	}

	clientErr := errors.New("woocommerce: unauthorized — invalid consumer key/secret")

	mockClient := &mocks.WooClientMock{
		TestConnectionFn: func(_ context.Context, _, _, _ string) error {
			return clientErr
		},
	}

	uc := New(mockClient, &mocks.IntegrationServiceMock{}, &mocks.OrderPublisherMock{}, mocks.NewLoggerMock())

	err := uc.TestConnection(ctx, config, credentials)

	if err == nil {
		t.Fatal("esperaba error de conexión, recibí nil")
	}
	if !errors.Is(err, clientErr) {
		t.Errorf("esperaba error wrapping '%v', recibí '%v'", clientErr, err)
	}
}

// ─────────────────────────────────────────────────────────────
// TestConnection — tipo de dato incorrecto en config (no string)
// ─────────────────────────────────────────────────────────────

func TestTestConnection_NonStringStoreURL(t *testing.T) {
	ctx := context.Background()
	config := map[string]interface{}{
		"store_url": 12345, // tipo incorrecto
	}
	credentials := map[string]interface{}{
		"consumer_key":    "ck_test",
		"consumer_secret": "cs_test",
	}

	uc := New(&mocks.WooClientMock{}, &mocks.IntegrationServiceMock{}, &mocks.OrderPublisherMock{}, mocks.NewLoggerMock())

	err := uc.TestConnection(ctx, config, credentials)

	if err == nil {
		t.Fatal("esperaba error cuando store_url no es string, recibí nil")
	}
}

// ─────────────────────────────────────────────────────────────
// TestConnection — tipo de dato incorrecto en credentials (no string)
// ─────────────────────────────────────────────────────────────

func TestTestConnection_NonStringConsumerKey(t *testing.T) {
	ctx := context.Background()
	config := map[string]interface{}{
		"store_url": "https://mitienda.com",
	}
	credentials := map[string]interface{}{
		"consumer_key":    42, // tipo incorrecto
		"consumer_secret": "cs_test",
	}

	uc := New(&mocks.WooClientMock{}, &mocks.IntegrationServiceMock{}, &mocks.OrderPublisherMock{}, mocks.NewLoggerMock())

	err := uc.TestConnection(ctx, config, credentials)

	if err == nil {
		t.Fatal("esperaba error cuando consumer_key no es string, recibí nil")
	}
}
