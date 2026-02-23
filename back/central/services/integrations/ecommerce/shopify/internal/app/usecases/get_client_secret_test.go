package usecases

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/domain"
)

func TestGetClientSecretByShopDomain_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	shopDomain := "mi-tienda.myshopify.com"
	integrationID := uint(7)
	expectedSecret := "super-secret-key-abc"

	integrationSvc := &mockIntegrationService{
		GetIntegrationByExternalIDFn: func(ctx context.Context, externalID string, integrationType int) (*domain.Integration, error) {
			if externalID != shopDomain {
				t.Errorf("externalID inesperado: got %q, want %q", externalID, shopDomain)
			}
			if integrationType != domain.IntegrationTypeID {
				t.Errorf("integrationType inesperado: got %d, want %d", integrationType, domain.IntegrationTypeID)
			}
			return &domain.Integration{ID: integrationID}, nil
		},
		DecryptCredentialFn: func(ctx context.Context, id string, fieldName string) (string, error) {
			expectedIDStr := fmt.Sprintf("%d", integrationID)
			if id != expectedIDStr {
				t.Errorf("integrationID incorrecto: got %q, want %q", id, expectedIDStr)
			}
			if fieldName != "client_secret" {
				t.Errorf("fieldName incorrecto: got %q, want %q", fieldName, "client_secret")
			}
			return expectedSecret, nil
		},
	}

	uc := newTestUseCase(integrationSvc, &mockShopifyClient{}, &mockOrderPublisher{})

	// Act
	secret, err := uc.GetClientSecretByShopDomain(ctx, shopDomain)

	// Assert
	if err != nil {
		t.Fatalf("se esperaba nil, se obtuvo: %v", err)
	}
	if secret != expectedSecret {
		t.Errorf("secret incorrecto: got %q, want %q", secret, expectedSecret)
	}
}

func TestGetClientSecretByShopDomain_IntegrationNotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()

	integrationSvc := &mockIntegrationService{
		GetIntegrationByExternalIDFn: func(ctx context.Context, externalID string, integrationType int) (*domain.Integration, error) {
			return nil, nil // Retorna nil sin error (integracion no encontrada)
		},
	}

	uc := newTestUseCase(integrationSvc, &mockShopifyClient{}, &mockOrderPublisher{})

	// Act
	secret, err := uc.GetClientSecretByShopDomain(ctx, "dominio-desconocido.myshopify.com")

	// Assert
	if err == nil {
		t.Fatal("se esperaba error por integracion no encontrada, se obtuvo nil")
	}
	if secret != "" {
		t.Errorf("se esperaba secret vacio, se obtuvo %q", secret)
	}
}

func TestGetClientSecretByShopDomain_IntegrationServiceError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	expectedErr := errors.New("db connection error")

	integrationSvc := &mockIntegrationService{
		GetIntegrationByExternalIDFn: func(ctx context.Context, externalID string, integrationType int) (*domain.Integration, error) {
			return nil, expectedErr
		},
	}

	uc := newTestUseCase(integrationSvc, &mockShopifyClient{}, &mockOrderPublisher{})

	// Act
	secret, err := uc.GetClientSecretByShopDomain(ctx, "mi-tienda.myshopify.com")

	// Assert
	if err == nil {
		t.Fatal("se esperaba error, se obtuvo nil")
	}
	if !errors.Is(err, expectedErr) {
		t.Errorf("error incorrecto: got %v, want %v", err, expectedErr)
	}
	if secret != "" {
		t.Errorf("se esperaba secret vacio, se obtuvo %q", secret)
	}
}

func TestGetClientSecretByShopDomain_DecryptError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	decryptErr := errors.New("decryption failed: invalid key")

	integrationSvc := &mockIntegrationService{
		GetIntegrationByExternalIDFn: func(ctx context.Context, externalID string, integrationType int) (*domain.Integration, error) {
			return &domain.Integration{ID: 5}, nil
		},
		DecryptCredentialFn: func(ctx context.Context, id string, fieldName string) (string, error) {
			return "", decryptErr
		},
	}

	uc := newTestUseCase(integrationSvc, &mockShopifyClient{}, &mockOrderPublisher{})

	// Act
	secret, err := uc.GetClientSecretByShopDomain(ctx, "mi-tienda.myshopify.com")

	// Assert
	if err == nil {
		t.Fatal("se esperaba error de descifrado, se obtuvo nil")
	}
	if !errors.Is(err, decryptErr) {
		t.Errorf("error incorrecto: got %v, want %v", err, decryptErr)
	}
	if secret != "" {
		t.Errorf("se esperaba secret vacio, se obtuvo %q", secret)
	}
}
