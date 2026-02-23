package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/domain"
)

func TestVerifyWebhooksByURL_Success_MatchingWebhooks(t *testing.T) {
	// Arrange
	ctx := context.Background()
	integrationID := "integration-1"
	baseURL := "https://api.miempresa.com"

	// Construimos la URL que el use case generara internamente:
	// fmt.Sprintf("%s/integrations/shopify/webhook", baseURL)
	expectedAddress := "https://api.miempresa.com/integrations/shopify/webhook"

	integrationSvc := &mockIntegrationService{
		GetIntegrationByIDFn: func(ctx context.Context, id string) (*domain.Integration, error) {
			return newIntegration(1, "mi-tienda.myshopify.com"), nil
		},
		DecryptCredentialFn: func(ctx context.Context, id string, fieldName string) (string, error) {
			return "shpat_test_token", nil
		},
	}

	shopifyClient := &mockShopifyClient{
		ListWebhooksFn: func(ctx context.Context, storeName, accessToken string) ([]domain.WebhookInfo, error) {
			// Shopify retorna webhooks: dos coinciden con nuestra URL, uno no
			return []domain.WebhookInfo{
				{ID: "wh-1", Address: expectedAddress, Topic: "orders/create"},
				{ID: "wh-2", Address: expectedAddress, Topic: "orders/paid"},
				{ID: "wh-ext", Address: "https://otro-servicio.com/webhooks", Topic: "orders/create"},
			}, nil
		},
	}

	uc := newTestUseCase(integrationSvc, shopifyClient, &mockOrderPublisher{})

	// Act
	matching, err := uc.VerifyWebhooksByURL(ctx, integrationID, baseURL)

	// Assert
	if err != nil {
		t.Fatalf("se esperaba nil, se obtuvo: %v", err)
	}
	if len(matching) != 2 {
		t.Fatalf("se esperaban 2 webhooks coincidentes, se encontraron %d", len(matching))
	}
	for _, wh := range matching {
		if wh.ID == "wh-ext" {
			t.Error("el webhook externo no debia estar en los resultados")
		}
	}
}

func TestVerifyWebhooksByURL_NoMatches(t *testing.T) {
	// Arrange
	ctx := context.Background()
	baseURL := "https://api.miempresa.com"

	integrationSvc := &mockIntegrationService{
		GetIntegrationByIDFn: func(ctx context.Context, id string) (*domain.Integration, error) {
			return newIntegration(1, "mi-tienda.myshopify.com"), nil
		},
		DecryptCredentialFn: func(ctx context.Context, id string, fieldName string) (string, error) {
			return "shpat_test_token", nil
		},
	}

	shopifyClient := &mockShopifyClient{
		ListWebhooksFn: func(ctx context.Context, storeName, accessToken string) ([]domain.WebhookInfo, error) {
			// Ningun webhook coincide con nuestra URL
			return []domain.WebhookInfo{
				{ID: "wh-ext", Address: "https://otro-servicio.com/webhooks", Topic: "orders/create"},
			}, nil
		},
	}

	uc := newTestUseCase(integrationSvc, shopifyClient, &mockOrderPublisher{})

	// Act
	matching, err := uc.VerifyWebhooksByURL(ctx, "integration-1", baseURL)

	// Assert
	if err != nil {
		t.Fatalf("se esperaba nil, se obtuvo: %v", err)
	}
	if len(matching) != 0 {
		t.Errorf("se esperaban 0 webhooks coincidentes, se encontraron %d", len(matching))
	}
}

func TestVerifyWebhooksByURL_URLNormalization(t *testing.T) {
	// Arrange: comprueba que URLs con trailing slash se normalizan correctamente
	ctx := context.Background()
	baseURL := "https://api.miempresa.com"

	// El address tiene trailing slash, debe normalizarse y coincidir igualmente
	expectedAddressWithSlash := "https://api.miempresa.com/integrations/shopify/webhook/"

	integrationSvc := &mockIntegrationService{
		GetIntegrationByIDFn: func(ctx context.Context, id string) (*domain.Integration, error) {
			return newIntegration(1, "mi-tienda.myshopify.com"), nil
		},
		DecryptCredentialFn: func(ctx context.Context, id string, fieldName string) (string, error) {
			return "shpat_test_token", nil
		},
	}

	shopifyClient := &mockShopifyClient{
		ListWebhooksFn: func(ctx context.Context, storeName, accessToken string) ([]domain.WebhookInfo, error) {
			return []domain.WebhookInfo{
				{ID: "wh-1", Address: expectedAddressWithSlash, Topic: "orders/create"},
			}, nil
		},
	}

	uc := newTestUseCase(integrationSvc, shopifyClient, &mockOrderPublisher{})

	// Act
	matching, err := uc.VerifyWebhooksByURL(ctx, "integration-1", baseURL)

	// Assert
	if err != nil {
		t.Fatalf("se esperaba nil, se obtuvo: %v", err)
	}
	if len(matching) != 1 {
		t.Errorf("URL con trailing slash deberia normalizarse y coincidir, se encontraron %d matches", len(matching))
	}
}

func TestVerifyWebhooksByURL_IntegrationServiceError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	expectedErr := errors.New("integration not found")

	integrationSvc := &mockIntegrationService{
		GetIntegrationByIDFn: func(ctx context.Context, id string) (*domain.Integration, error) {
			return nil, expectedErr
		},
	}

	uc := newTestUseCase(integrationSvc, &mockShopifyClient{}, &mockOrderPublisher{})

	// Act
	matching, err := uc.VerifyWebhooksByURL(ctx, "integration-1", "https://api.miempresa.com")

	// Assert
	if err == nil {
		t.Fatal("se esperaba error, se obtuvo nil")
	}
	if !errors.Is(err, expectedErr) {
		t.Errorf("error incorrecto: got %v, want %v", err, expectedErr)
	}
	if matching != nil {
		t.Error("se esperaba nil como resultado")
	}
}

func TestVerifyWebhooksByURL_ShopifyClientError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	clientErr := errors.New("shopify api 500")

	integrationSvc := &mockIntegrationService{
		GetIntegrationByIDFn: func(ctx context.Context, id string) (*domain.Integration, error) {
			return newIntegration(1, "mi-tienda.myshopify.com"), nil
		},
		DecryptCredentialFn: func(ctx context.Context, id string, fieldName string) (string, error) {
			return "shpat_test_token", nil
		},
	}

	shopifyClient := &mockShopifyClient{
		ListWebhooksFn: func(ctx context.Context, storeName, accessToken string) ([]domain.WebhookInfo, error) {
			return nil, clientErr
		},
	}

	uc := newTestUseCase(integrationSvc, shopifyClient, &mockOrderPublisher{})

	// Act
	matching, err := uc.VerifyWebhooksByURL(ctx, "integration-1", "https://api.miempresa.com")

	// Assert
	if err == nil {
		t.Fatal("se esperaba error del cliente Shopify, se obtuvo nil")
	}
	if !errors.Is(err, clientErr) {
		t.Errorf("error incorrecto: got %v, want %v", err, clientErr)
	}
	if matching != nil {
		t.Error("se esperaba nil como resultado")
	}
}
