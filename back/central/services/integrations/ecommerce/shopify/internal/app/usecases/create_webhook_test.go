package usecases

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/domain"
)

func TestCreateWebhook_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	integrationID := "integration-1"
	baseURL := "https://api.miempresa.com"

	expectedWebhookURL := "https://api.miempresa.com/api/v1/integrations/shopify/webhook"

	var configUpdateCalled bool
	var capturedConfig map[string]interface{}

	integrationSvc := &mockIntegrationService{
		GetIntegrationByIDFn: func(ctx context.Context, id string) (*domain.Integration, error) {
			return newIntegration(1, "mi-tienda.myshopify.com"), nil
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

	createdWebhookIDs := []string{}
	shopifyClient := &mockShopifyClient{
		// VerifyWebhooksByURL llama internamente a ListWebhooks
		ListWebhooksFn: func(ctx context.Context, storeName, accessToken string) ([]domain.WebhookInfo, error) {
			return []domain.WebhookInfo{}, nil // sin webhooks previos
		},
		CreateWebhookFn: func(ctx context.Context, storeName, accessToken, webhookURL, event string) (string, error) {
			// Verificar que la URL construida es la correcta
			if webhookURL != expectedWebhookURL {
				t.Errorf("webhookURL incorrecto: got %q, want %q", webhookURL, expectedWebhookURL)
			}
			id := "wh-" + event
			createdWebhookIDs = append(createdWebhookIDs, id)
			return id, nil
		},
	}

	uc := newTestUseCase(integrationSvc, shopifyClient, &mockOrderPublisher{})

	// Act
	result, err := uc.CreateWebhook(ctx, integrationID, baseURL)

	// Assert
	if err != nil {
		t.Fatalf("se esperaba nil, se obtuvo: %v", err)
	}
	if result == nil {
		t.Fatal("el resultado no debe ser nil")
	}
	if result.WebhookURL != expectedWebhookURL {
		t.Errorf("WebhookURL incorrecto: got %q, want %q", result.WebhookURL, expectedWebhookURL)
	}
	// Se esperan 6 eventos registrados
	if len(result.CreatedWebhooks) != 6 {
		t.Errorf("se esperaban 6 webhooks creados, se crearon %d", len(result.CreatedWebhooks))
	}
	if !configUpdateCalled {
		t.Fatal("UpdateIntegrationConfig no fue llamado")
	}
	if capturedConfig["webhook_url"] != expectedWebhookURL {
		t.Errorf("webhook_url en config incorrecto: got %v, want %q", capturedConfig["webhook_url"], expectedWebhookURL)
	}
}

func TestCreateWebhook_LocalhostBlocked(t *testing.T) {
	// Arrange: los webhooks no se pueden crear en localhost
	ctx := context.Background()

	integrationSvc := &mockIntegrationService{
		GetIntegrationByIDFn: func(ctx context.Context, id string) (*domain.Integration, error) {
			return newIntegration(1, "mi-tienda.myshopify.com"), nil
		},
		DecryptCredentialFn: func(ctx context.Context, id string, fieldName string) (string, error) {
			return "shpat_test_token", nil
		},
	}

	uc := newTestUseCase(integrationSvc, &mockShopifyClient{}, &mockOrderPublisher{})

	tests := []struct {
		name    string
		baseURL string
	}{
		{"localhost", "http://localhost:3000"},
		{"127.0.0.1", "http://127.0.0.1:8080"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result, err := uc.CreateWebhook(ctx, "integration-1", tt.baseURL)

			// Assert
			if err == nil {
				t.Fatalf("[%s] se esperaba error por localhost, se obtuvo nil", tt.name)
			}
			// El resultado debe incluir la URL aunque haya error
			if result == nil {
				t.Fatalf("[%s] el result no debe ser nil aunque haya error", tt.name)
			}
			if !strings.Contains(err.Error(), "localhost") && !strings.Contains(err.Error(), "pruebas") {
				t.Errorf("[%s] el error deberia mencionar el entorno de pruebas: %v", tt.name, err)
			}
		})
	}
}

func TestCreateWebhook_IntegrationServiceError(t *testing.T) {
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
	result, err := uc.CreateWebhook(ctx, "integration-1", "https://api.miempresa.com")

	// Assert
	if err == nil {
		t.Fatal("se esperaba error, se obtuvo nil")
	}
	if !errors.Is(err, expectedErr) {
		t.Errorf("error incorrecto: got %v, want %v", err, expectedErr)
	}
	if result != nil {
		t.Error("se esperaba nil como resultado")
	}
}

func TestCreateWebhook_StoreNameMissing(t *testing.T) {
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
	result, err := uc.CreateWebhook(ctx, "integration-1", "https://api.miempresa.com")

	// Assert
	if err == nil {
		t.Fatal("se esperaba error por store_name faltante, se obtuvo nil")
	}
	if result != nil {
		t.Error("se esperaba nil como resultado")
	}
}

func TestCreateWebhook_AllEventsFailToCreate_ReturnsError(t *testing.T) {
	// Arrange: si no se puede crear ningun webhook, el use case retorna error
	ctx := context.Background()
	createErr := errors.New("shopify api 403 forbidden")

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
			return []domain.WebhookInfo{}, nil
		},
		CreateWebhookFn: func(ctx context.Context, storeName, accessToken, webhookURL, event string) (string, error) {
			return "", createErr // Todos los eventos fallan
		},
	}

	uc := newTestUseCase(integrationSvc, shopifyClient, &mockOrderPublisher{})

	// Act
	result, err := uc.CreateWebhook(ctx, "integration-1", "https://api.miempresa.com")

	// Assert
	if err == nil {
		t.Fatal("se esperaba error cuando ningun webhook se puede crear, se obtuvo nil")
	}
	if result == nil {
		t.Fatal("el resultado no debe ser nil aunque haya error")
	}
	if len(result.CreatedWebhooks) != 0 {
		t.Errorf("no deberia haber webhooks creados, se encontraron %d", len(result.CreatedWebhooks))
	}
}

func TestCreateWebhook_DeletesExistingWebhooksBeforeCreating(t *testing.T) {
	// Arrange: si existen webhooks previos con la misma URL, deben eliminarse antes de crear nuevos
	ctx := context.Background()
	baseURL := "https://api.miempresa.com"

	existingWebhook := domain.WebhookInfo{
		ID:      "old-wh-1",
		Address: baseURL + "/integrations/shopify/webhook",
		Topic:   "orders/create",
	}

	var deletedIDs []string

	integrationSvc := &mockIntegrationService{
		GetIntegrationByIDFn: func(ctx context.Context, id string) (*domain.Integration, error) {
			return newIntegration(1, "mi-tienda.myshopify.com"), nil
		},
		DecryptCredentialFn: func(ctx context.Context, id string, fieldName string) (string, error) {
			return "shpat_test_token", nil
		},
		UpdateIntegrationConfigFn: func(ctx context.Context, id string, config map[string]interface{}) error {
			return nil
		},
	}

	shopifyClient := &mockShopifyClient{
		ListWebhooksFn: func(ctx context.Context, storeName, accessToken string) ([]domain.WebhookInfo, error) {
			return []domain.WebhookInfo{existingWebhook}, nil
		},
		DeleteWebhookFn: func(ctx context.Context, storeName, accessToken, webhookID string) error {
			deletedIDs = append(deletedIDs, webhookID)
			return nil
		},
		CreateWebhookFn: func(ctx context.Context, storeName, accessToken, webhookURL, event string) (string, error) {
			return "new-wh-" + event, nil
		},
	}

	uc := newTestUseCase(integrationSvc, shopifyClient, &mockOrderPublisher{})

	// Act
	result, err := uc.CreateWebhook(ctx, "integration-1", baseURL)

	// Assert
	if err != nil {
		t.Fatalf("se esperaba nil, se obtuvo: %v", err)
	}
	if len(deletedIDs) == 0 {
		t.Error("debia eliminarse el webhook anterior antes de crear nuevos")
	}
	if result.DeletedWebhooks == nil || len(result.DeletedWebhooks) == 0 {
		t.Error("DeletedWebhooks debe contener el webhook eliminado")
	}
}
