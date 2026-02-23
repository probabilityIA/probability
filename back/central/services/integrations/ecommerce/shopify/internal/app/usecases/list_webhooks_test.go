package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/domain"
)

func TestListWebhooks_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	integrationID := "integration-1"

	expectedWebhooks := []domain.WebhookInfo{
		{ID: "wh-1", Address: "https://api.ejemplo.com/shopify/webhook", Topic: "orders/create"},
		{ID: "wh-2", Address: "https://api.ejemplo.com/shopify/webhook", Topic: "orders/paid"},
	}

	integrationSvc := &mockIntegrationService{
		GetIntegrationByIDFn: func(ctx context.Context, id string) (*domain.Integration, error) {
			if id != integrationID {
				t.Errorf("integrationID inesperado: got %q, want %q", id, integrationID)
			}
			return newIntegration(1, "mi-tienda.myshopify.com"), nil
		},
		DecryptCredentialFn: func(ctx context.Context, id string, fieldName string) (string, error) {
			return "shpat_test_token", nil
		},
	}

	shopifyClient := &mockShopifyClient{
		ListWebhooksFn: func(ctx context.Context, storeName, accessToken string) ([]domain.WebhookInfo, error) {
			if storeName != "mi-tienda.myshopify.com" {
				t.Errorf("storeName inesperado: got %q", storeName)
			}
			return expectedWebhooks, nil
		},
	}

	uc := newTestUseCase(integrationSvc, shopifyClient, &mockOrderPublisher{})

	// Act
	webhooks, err := uc.ListWebhooks(ctx, integrationID)

	// Assert
	if err != nil {
		t.Fatalf("se esperaba nil, se obtuvo: %v", err)
	}
	if len(webhooks) != len(expectedWebhooks) {
		t.Fatalf("cantidad de webhooks incorrecta: got %d, want %d", len(webhooks), len(expectedWebhooks))
	}
	for i, wh := range webhooks {
		if wh.ID != expectedWebhooks[i].ID {
			t.Errorf("webhook[%d].ID incorrecto: got %q, want %q", i, wh.ID, expectedWebhooks[i].ID)
		}
	}
}

func TestListWebhooks_IntegrationNotFound(t *testing.T) {
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
	webhooks, err := uc.ListWebhooks(ctx, "integration-inexistente")

	// Assert
	if err == nil {
		t.Fatal("se esperaba error, se obtuvo nil")
	}
	if !errors.Is(err, expectedErr) {
		t.Errorf("error incorrecto: got %v, want %v", err, expectedErr)
	}
	if webhooks != nil {
		t.Errorf("se esperaba nil como resultado, se obtuvo %v", webhooks)
	}
}

func TestListWebhooks_DecryptCredentialError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	decryptErr := errors.New("cannot decrypt access_token")

	integrationSvc := &mockIntegrationService{
		GetIntegrationByIDFn: func(ctx context.Context, id string) (*domain.Integration, error) {
			return newIntegration(1, "mi-tienda.myshopify.com"), nil
		},
		DecryptCredentialFn: func(ctx context.Context, id string, fieldName string) (string, error) {
			return "", decryptErr
		},
	}

	uc := newTestUseCase(integrationSvc, &mockShopifyClient{}, &mockOrderPublisher{})

	// Act
	webhooks, err := uc.ListWebhooks(ctx, "integration-1")

	// Assert
	if err == nil {
		t.Fatal("se esperaba error al descifrar credencial, se obtuvo nil")
	}
	if !errors.Is(err, decryptErr) {
		t.Errorf("error incorrecto: got %v, want %v", err, decryptErr)
	}
	if webhooks != nil {
		t.Error("se esperaba nil como resultado de webhooks")
	}
}

func TestListWebhooks_StoreNameMissingInConfig(t *testing.T) {
	// Arrange
	ctx := context.Background()
	businessID := uint(1)

	integrationSvc := &mockIntegrationService{
		GetIntegrationByIDFn: func(ctx context.Context, id string) (*domain.Integration, error) {
			// Retorna integracion sin store_name en el config
			return &domain.Integration{
				ID:         1,
				BusinessID: &businessID,
				Config:     map[string]interface{}{}, // config vacio, sin store_name
			}, nil
		},
		DecryptCredentialFn: func(ctx context.Context, id string, fieldName string) (string, error) {
			return "shpat_test_token", nil
		},
	}

	uc := newTestUseCase(integrationSvc, &mockShopifyClient{}, &mockOrderPublisher{})

	// Act
	webhooks, err := uc.ListWebhooks(ctx, "integration-1")

	// Assert
	if err == nil {
		t.Fatal("se esperaba error por store_name faltante, se obtuvo nil")
	}
	if webhooks != nil {
		t.Error("se esperaba nil como resultado de webhooks")
	}
}

func TestListWebhooks_ClientError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	clientErr := errors.New("shopify api error 503")

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
	webhooks, err := uc.ListWebhooks(ctx, "integration-1")

	// Assert
	if err == nil {
		t.Fatal("se esperaba error del cliente Shopify, se obtuvo nil")
	}
	if !errors.Is(err, clientErr) {
		t.Errorf("error incorrecto: got %v, want %v", err, clientErr)
	}
	if webhooks != nil {
		t.Error("se esperaba nil como resultado de webhooks")
	}
}

func TestListWebhooks_EmptyResult(t *testing.T) {
	// Arrange
	ctx := context.Background()

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
			return []domain.WebhookInfo{}, nil // Lista vacia, es valido
		},
	}

	uc := newTestUseCase(integrationSvc, shopifyClient, &mockOrderPublisher{})

	// Act
	webhooks, err := uc.ListWebhooks(ctx, "integration-1")

	// Assert
	if err != nil {
		t.Fatalf("se esperaba nil, se obtuvo: %v", err)
	}
	if len(webhooks) != 0 {
		t.Errorf("se esperaba lista vacia, se obtuvo %d webhooks", len(webhooks))
	}
}
