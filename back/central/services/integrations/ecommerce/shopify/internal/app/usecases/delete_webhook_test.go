package usecases

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/domain"
)

func TestDeleteWebhook_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	integrationID := "integration-5"
	webhookID := "wh-999"
	businessID := uint(1)

	var configUpdateCalled bool
	var capturedConfig map[string]interface{}

	integrationSvc := &mockIntegrationService{
		GetIntegrationByIDFn: func(ctx context.Context, id string) (*domain.Integration, error) {
			return &domain.Integration{
				ID:         5,
				BusinessID: &businessID,
				Config: map[string]interface{}{
					"store_name":  "mi-tienda.myshopify.com",
					"webhook_ids": []interface{}{"wh-111", "wh-999", "wh-222"},
				},
			}, nil
		},
		DecryptCredentialFn: func(ctx context.Context, id string, fieldName string) (string, error) {
			return "shpat_test_token", nil
		},
		UpdateIntegrationConfigFn: func(ctx context.Context, id string, config map[string]interface{}) error {
			configUpdateCalled = true
			capturedConfig = config
			return nil
		},
	}

	var deletedWebhookID string
	shopifyClient := &mockShopifyClient{
		DeleteWebhookFn: func(ctx context.Context, storeName, accessToken, whID string) error {
			deletedWebhookID = whID
			return nil
		},
	}

	uc := newTestUseCase(integrationSvc, shopifyClient, &mockOrderPublisher{})

	// Act
	err := uc.DeleteWebhook(ctx, integrationID, webhookID)

	// Assert
	if err != nil {
		t.Fatalf("se esperaba nil, se obtuvo: %v", err)
	}
	if deletedWebhookID != webhookID {
		t.Errorf("webhookID eliminado incorrecto: got %q, want %q", deletedWebhookID, webhookID)
	}
	if !configUpdateCalled {
		t.Fatal("UpdateIntegrationConfig no fue llamado")
	}
	// Verificar que el webhook eliminado no esta en la nueva lista
	if ids, ok := capturedConfig["webhook_ids"].([]interface{}); ok {
		for _, id := range ids {
			if fmt.Sprintf("%v", id) == webhookID {
				t.Errorf("el webhook %q no deberia estar en la lista actualizada", webhookID)
			}
		}
	}
}

func TestDeleteWebhook_LastWebhookCleared(t *testing.T) {
	// Arrange: cuando solo queda un webhook y se elimina, webhook_configured debe quedar en false
	ctx := context.Background()
	businessID := uint(1)
	webhookID := "wh-solo"

	var capturedConfig map[string]interface{}

	integrationSvc := &mockIntegrationService{
		GetIntegrationByIDFn: func(ctx context.Context, id string) (*domain.Integration, error) {
			return &domain.Integration{
				ID:         5,
				BusinessID: &businessID,
				Config: map[string]interface{}{
					"store_name":  "mi-tienda.myshopify.com",
					"webhook_ids": []interface{}{webhookID}, // solo uno
				},
			}, nil
		},
		DecryptCredentialFn: func(ctx context.Context, id string, fieldName string) (string, error) {
			return "shpat_test_token", nil
		},
		UpdateIntegrationConfigFn: func(ctx context.Context, id string, config map[string]interface{}) error {
			capturedConfig = config
			return nil
		},
	}

	shopifyClient := &mockShopifyClient{
		DeleteWebhookFn: func(ctx context.Context, storeName, accessToken, whID string) error {
			return nil
		},
	}

	uc := newTestUseCase(integrationSvc, shopifyClient, &mockOrderPublisher{})

	// Act
	err := uc.DeleteWebhook(ctx, "integration-5", webhookID)

	// Assert
	if err != nil {
		t.Fatalf("se esperaba nil, se obtuvo: %v", err)
	}

	if configured, ok := capturedConfig["webhook_configured"].(bool); !ok || configured {
		t.Error("webhook_configured debe ser false cuando se elimina el ultimo webhook")
	}
}

func TestDeleteWebhook_IntegrationServiceError(t *testing.T) {
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
	err := uc.DeleteWebhook(ctx, "integration-1", "wh-1")

	// Assert
	if err == nil {
		t.Fatal("se esperaba error, se obtuvo nil")
	}
	if !errors.Is(err, expectedErr) {
		t.Errorf("error incorrecto: got %v, want %v", err, expectedErr)
	}
}

func TestDeleteWebhook_DecryptCredentialError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	decryptErr := errors.New("decrypt failed")

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
	err := uc.DeleteWebhook(ctx, "integration-1", "wh-1")

	// Assert
	if err == nil {
		t.Fatal("se esperaba error de descifrado, se obtuvo nil")
	}
	if !errors.Is(err, decryptErr) {
		t.Errorf("error incorrecto: got %v, want %v", err, decryptErr)
	}
}

func TestDeleteWebhook_StoreNameMissing(t *testing.T) {
	// Arrange
	ctx := context.Background()
	businessID := uint(1)

	integrationSvc := &mockIntegrationService{
		GetIntegrationByIDFn: func(ctx context.Context, id string) (*domain.Integration, error) {
			return &domain.Integration{
				ID:         1,
				BusinessID: &businessID,
				Config:     map[string]interface{}{}, // sin store_name
			}, nil
		},
		DecryptCredentialFn: func(ctx context.Context, id string, fieldName string) (string, error) {
			return "shpat_test_token", nil
		},
	}

	uc := newTestUseCase(integrationSvc, &mockShopifyClient{}, &mockOrderPublisher{})

	// Act
	err := uc.DeleteWebhook(ctx, "integration-1", "wh-1")

	// Assert
	if err == nil {
		t.Fatal("se esperaba error por store_name faltante, se obtuvo nil")
	}
}

func TestDeleteWebhook_ShopifyClientError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	clientErr := errors.New("webhook not found in shopify")

	integrationSvc := &mockIntegrationService{
		GetIntegrationByIDFn: func(ctx context.Context, id string) (*domain.Integration, error) {
			return newIntegration(1, "mi-tienda.myshopify.com"), nil
		},
		DecryptCredentialFn: func(ctx context.Context, id string, fieldName string) (string, error) {
			return "shpat_test_token", nil
		},
	}

	shopifyClient := &mockShopifyClient{
		DeleteWebhookFn: func(ctx context.Context, storeName, accessToken, webhookID string) error {
			return clientErr
		},
	}

	uc := newTestUseCase(integrationSvc, shopifyClient, &mockOrderPublisher{})

	// Act
	err := uc.DeleteWebhook(ctx, "integration-1", "wh-1")

	// Assert
	if err == nil {
		t.Fatal("se esperaba error del cliente Shopify, se obtuvo nil")
	}
	if !errors.Is(err, clientErr) {
		t.Errorf("error incorrecto: got %v, want %v", err, clientErr)
	}
}
