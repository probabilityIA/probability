package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/domain"
)

func TestGetOrder_Success(t *testing.T) {
	ctx := context.Background()
	businessID := uint(42)

	integrationSvc := &mockIntegrationService{
		GetIntegrationByIDFn: func(ctx context.Context, id string) (*domain.Integration, error) {
			return &domain.Integration{
				ID:         7,
				BusinessID: &businessID,
				Name:       "mi-tienda.myshopify.com",
				Config:     map[string]interface{}{"store_name": "mi-tienda.myshopify.com"},
			}, nil
		},
	}

	shopifyClient := &mockShopifyClient{
		GetOrderFn: func(ctx context.Context, storeName, accessToken, orderID string) (*domain.ShopifyOrder, error) {
			return &domain.ShopifyOrder{
				ExternalID:  "12345",
				OrderNumber: "#1001",
				TotalAmount: 150000,
				Currency:    "COP",
				Customer:    domain.ShopifyCustomer{Name: "Test Customer"},
				Metadata:    map[string]interface{}{},
			}, nil
		},
	}

	publisher := &mockOrderPublisher{}
	uc := newTestUseCase(integrationSvc, shopifyClient, publisher, &mockSyncEventPublisher{})

	err := uc.GetOrder(ctx, "7", "12345")

	if err != nil {
		t.Fatalf("se esperaba nil, se obtuvo: %v", err)
	}
	if len(publisher.PublishedOrders) != 1 {
		t.Fatalf("se esperaba 1 orden publicada, se publicaron %d", len(publisher.PublishedOrders))
	}

	published := publisher.PublishedOrders[0]
	if published.IntegrationID != 7 {
		t.Errorf("IntegrationID incorrecto: got %d, want 7", published.IntegrationID)
	}
	if published.IntegrationType != "shopify" {
		t.Errorf("IntegrationType incorrecto: got %q, want %q", published.IntegrationType, "shopify")
	}
	if published.BusinessID == nil || *published.BusinessID != businessID {
		t.Errorf("BusinessID incorrecto: got %v, want %d", published.BusinessID, businessID)
	}
	if published.Invoiceable != true {
		t.Error("Invoiceable debe ser true para ordenes de Shopify")
	}
}

func TestGetOrder_IntegrationNotFound(t *testing.T) {
	ctx := context.Background()
	expectedErr := errors.New("integration not found")

	integrationSvc := &mockIntegrationService{
		GetIntegrationByIDFn: func(ctx context.Context, id string) (*domain.Integration, error) {
			return nil, expectedErr
		},
	}

	publisher := &mockOrderPublisher{}
	uc := newTestUseCase(integrationSvc, &mockShopifyClient{}, publisher, &mockSyncEventPublisher{})

	err := uc.GetOrder(ctx, "999", "12345")

	if err == nil {
		t.Fatal("se esperaba error, se obtuvo nil")
	}
	if !errors.Is(err, expectedErr) {
		t.Errorf("error incorrecto: got %v, want %v", err, expectedErr)
	}
	if len(publisher.PublishedOrders) != 0 {
		t.Error("no deberia publicarse ninguna orden cuando falla la integracion")
	}
}

func TestGetOrder_StoreNameMissing(t *testing.T) {
	ctx := context.Background()
	businessID := uint(1)

	integrationSvc := &mockIntegrationService{
		GetIntegrationByIDFn: func(ctx context.Context, id string) (*domain.Integration, error) {
			return &domain.Integration{
				ID:         1,
				BusinessID: &businessID,
				Name:       "", // nombre vacio
				Config:     map[string]interface{}{},
			}, nil
		},
	}

	uc := newTestUseCase(integrationSvc, &mockShopifyClient{}, &mockOrderPublisher{}, &mockSyncEventPublisher{})

	err := uc.GetOrder(ctx, "1", "12345")

	if err == nil {
		t.Fatal("se esperaba error por store_name faltante, se obtuvo nil")
	}
}

func TestGetOrder_DecryptCredentialError(t *testing.T) {
	ctx := context.Background()
	decryptErr := errors.New("decryption failed")

	integrationSvc := &mockIntegrationService{
		GetIntegrationByIDFn: func(ctx context.Context, id string) (*domain.Integration, error) {
			return newIntegration(1, "mi-tienda.myshopify.com"), nil
		},
		DecryptCredentialFn: func(ctx context.Context, id string, fieldName string) (string, error) {
			return "", decryptErr
		},
	}

	uc := newTestUseCase(integrationSvc, &mockShopifyClient{}, &mockOrderPublisher{}, &mockSyncEventPublisher{})

	err := uc.GetOrder(ctx, "1", "12345")

	if err == nil {
		t.Fatal("se esperaba error de descifrado, se obtuvo nil")
	}
	if !errors.Is(err, decryptErr) {
		t.Errorf("error incorrecto: got %v, want %v", err, decryptErr)
	}
}

func TestGetOrder_ShopifyClientError(t *testing.T) {
	ctx := context.Background()
	clientErr := errors.New("shopify api 500")

	integrationSvc := &mockIntegrationService{
		GetIntegrationByIDFn: func(ctx context.Context, id string) (*domain.Integration, error) {
			return newIntegration(1, "mi-tienda.myshopify.com"), nil
		},
	}

	shopifyClient := &mockShopifyClient{
		GetOrderFn: func(ctx context.Context, storeName, accessToken, orderID string) (*domain.ShopifyOrder, error) {
			return nil, clientErr
		},
	}

	uc := newTestUseCase(integrationSvc, shopifyClient, &mockOrderPublisher{}, &mockSyncEventPublisher{})

	err := uc.GetOrder(ctx, "1", "12345")

	if err == nil {
		t.Fatal("se esperaba error del cliente Shopify, se obtuvo nil")
	}
	if !errors.Is(err, clientErr) {
		t.Errorf("error incorrecto: got %v, want %v", err, clientErr)
	}
}

func TestGetOrder_PublishError(t *testing.T) {
	ctx := context.Background()
	businessID := uint(1)

	integrationSvc := &mockIntegrationService{
		GetIntegrationByIDFn: func(ctx context.Context, id string) (*domain.Integration, error) {
			return &domain.Integration{
				ID:         1,
				BusinessID: &businessID,
				Name:       "mi-tienda.myshopify.com",
				Config:     map[string]interface{}{"store_name": "mi-tienda.myshopify.com"},
			}, nil
		},
	}

	shopifyClient := &mockShopifyClient{
		GetOrderFn: func(ctx context.Context, storeName, accessToken, orderID string) (*domain.ShopifyOrder, error) {
			return &domain.ShopifyOrder{
				ExternalID: "12345",
				Metadata:   map[string]interface{}{},
			}, nil
		},
	}

	publisher := &mockOrderPublisher{
		PublishFn: func(ctx context.Context, order *domain.ProbabilityOrderDTO) error {
			return errors.New("rabbitmq connection lost")
		},
	}

	uc := newTestUseCase(integrationSvc, shopifyClient, publisher, &mockSyncEventPublisher{})

	err := uc.GetOrder(ctx, "1", "12345")

	if err == nil {
		t.Fatal("se esperaba error del publisher, se obtuvo nil")
	}
	if !errors.Is(err, domain.ErrPublishFailed) {
		t.Errorf("error debe envolver ErrPublishFailed: got %v", err)
	}
}

func TestGetOrder_OrderFieldsAssigned(t *testing.T) {
	ctx := context.Background()
	businessID := uint(55)

	integrationSvc := &mockIntegrationService{
		GetIntegrationByIDFn: func(ctx context.Context, id string) (*domain.Integration, error) {
			return &domain.Integration{
				ID:         9,
				BusinessID: &businessID,
				Name:       "mi-tienda.myshopify.com",
				Config:     map[string]interface{}{"store_name": "mi-tienda.myshopify.com"},
			}, nil
		},
	}

	shopifyClient := &mockShopifyClient{
		GetOrderFn: func(ctx context.Context, storeName, accessToken, orderID string) (*domain.ShopifyOrder, error) {
			return &domain.ShopifyOrder{
				ExternalID:  "order-99",
				OrderNumber: "#99",
				TotalAmount: 50000,
				Currency:    "COP",
				Metadata:    map[string]interface{}{},
			}, nil
		},
	}

	publisher := &mockOrderPublisher{}
	uc := newTestUseCase(integrationSvc, shopifyClient, publisher, &mockSyncEventPublisher{})

	err := uc.GetOrder(ctx, "9", "order-99")

	if err != nil {
		t.Fatalf("se esperaba nil, se obtuvo: %v", err)
	}

	published := publisher.PublishedOrders[0]
	if published.IntegrationID != 9 {
		t.Errorf("IntegrationID: got %d, want 9", published.IntegrationID)
	}
	if published.BusinessID == nil || *published.BusinessID != businessID {
		t.Errorf("BusinessID: got %v, want %d", published.BusinessID, businessID)
	}
	if published.IntegrationType != "shopify" {
		t.Errorf("IntegrationType: got %q, want %q", published.IntegrationType, "shopify")
	}
}
