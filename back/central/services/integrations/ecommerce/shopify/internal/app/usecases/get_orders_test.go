package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/domain"
)

func TestGetOrders_Success_SinglePage(t *testing.T) {
	ctx := context.Background()
	businessID := uint(42)

	integration := &domain.Integration{
		ID:         7,
		BusinessID: &businessID,
		Name:       "mi-tienda.myshopify.com",
		IsActive:   true,
	}

	shopifyClient := &mockShopifyClient{
		GetOrdersFn: func(ctx context.Context, storeName, accessToken string, params *domain.GetOrdersParams) ([]domain.ShopifyOrder, string, error) {
			return []domain.ShopifyOrder{
				{ExternalID: "o1", Metadata: map[string]interface{}{}},
				{ExternalID: "o2", Metadata: map[string]interface{}{}},
				{ExternalID: "o3", Metadata: map[string]interface{}{}},
			}, "", nil // sin next page
		},
	}

	publisher := &mockOrderPublisher{}
	uc := newTestUseCase(&mockIntegrationService{}, shopifyClient, publisher, &mockSyncEventPublisher{})

	total, err := uc.GetOrders(ctx, integration, "mi-tienda.myshopify.com", "token", &domain.GetOrdersParams{Limit: 250})

	if err != nil {
		t.Fatalf("se esperaba nil, se obtuvo: %v", err)
	}
	if total != 3 {
		t.Errorf("total incorrecto: got %d, want 3", total)
	}
	if len(publisher.PublishedOrders) != 3 {
		t.Errorf("se esperaban 3 ordenes publicadas, se publicaron %d", len(publisher.PublishedOrders))
	}
}

func TestGetOrders_Success_MultiPage(t *testing.T) {
	ctx := context.Background()
	businessID := uint(42)

	integration := &domain.Integration{
		ID:         7,
		BusinessID: &businessID,
		Name:       "mi-tienda.myshopify.com",
		IsActive:   true,
	}

	shopifyClient := &mockShopifyClient{
		GetOrdersFn: func(ctx context.Context, storeName, accessToken string, params *domain.GetOrdersParams) ([]domain.ShopifyOrder, string, error) {
			return []domain.ShopifyOrder{
				{ExternalID: "o1", Metadata: map[string]interface{}{}},
				{ExternalID: "o2", Metadata: map[string]interface{}{}},
			}, "https://shopify.com/next-page", nil
		},
		GetOrdersByURLFn: func(ctx context.Context, nextPageURL, accessToken string) ([]domain.ShopifyOrder, string, error) {
			return []domain.ShopifyOrder{
				{ExternalID: "o3", Metadata: map[string]interface{}{}},
			}, "", nil
		},
	}

	publisher := &mockOrderPublisher{}
	uc := newTestUseCase(&mockIntegrationService{}, shopifyClient, publisher, &mockSyncEventPublisher{})

	total, err := uc.GetOrders(ctx, integration, "mi-tienda.myshopify.com", "token", &domain.GetOrdersParams{Limit: 250})

	if err != nil {
		t.Fatalf("se esperaba nil, se obtuvo: %v", err)
	}
	if total != 3 {
		t.Errorf("total incorrecto: got %d, want 3", total)
	}
}

func TestGetOrders_EmptyResults(t *testing.T) {
	ctx := context.Background()
	businessID := uint(42)

	integration := &domain.Integration{
		ID:         7,
		BusinessID: &businessID,
		Name:       "mi-tienda.myshopify.com",
		IsActive:   true,
	}

	shopifyClient := &mockShopifyClient{
		GetOrdersFn: func(ctx context.Context, storeName, accessToken string, params *domain.GetOrdersParams) ([]domain.ShopifyOrder, string, error) {
			return []domain.ShopifyOrder{}, "", nil
		},
	}

	publisher := &mockOrderPublisher{}
	uc := newTestUseCase(&mockIntegrationService{}, shopifyClient, publisher, &mockSyncEventPublisher{})

	total, err := uc.GetOrders(ctx, integration, "mi-tienda.myshopify.com", "token", &domain.GetOrdersParams{Limit: 250})

	if err != nil {
		t.Fatalf("se esperaba nil, se obtuvo: %v", err)
	}
	if total != 0 {
		t.Errorf("total incorrecto: got %d, want 0", total)
	}
	if len(publisher.PublishedOrders) != 0 {
		t.Errorf("no deberian publicarse ordenes, se publicaron %d", len(publisher.PublishedOrders))
	}
}

func TestGetOrders_BusinessIDNil_ReturnsError(t *testing.T) {
	ctx := context.Background()

	integration := &domain.Integration{
		ID:         7,
		BusinessID: nil, // sin business_id
		Name:       "mi-tienda.myshopify.com",
		IsActive:   true,
	}

	shopifyClient := &mockShopifyClient{
		GetOrdersFn: func(ctx context.Context, storeName, accessToken string, params *domain.GetOrdersParams) ([]domain.ShopifyOrder, string, error) {
			return []domain.ShopifyOrder{
				{ExternalID: "o1", Metadata: map[string]interface{}{}},
			}, "", nil
		},
	}

	uc := newTestUseCase(&mockIntegrationService{}, shopifyClient, &mockOrderPublisher{}, &mockSyncEventPublisher{})

	_, err := uc.GetOrders(ctx, integration, "mi-tienda.myshopify.com", "token", &domain.GetOrdersParams{Limit: 250})

	if err == nil {
		t.Fatal("se esperaba error por BusinessID nil, se obtuvo nil")
	}
	if !errors.Is(err, domain.ErrBusinessIDMissing) {
		t.Errorf("error debe envolver ErrBusinessIDMissing: got %v", err)
	}
}

func TestGetOrders_ShopifyClientError(t *testing.T) {
	ctx := context.Background()
	businessID := uint(42)
	clientErr := errors.New("shopify api 429 rate limited")

	integration := &domain.Integration{
		ID:         7,
		BusinessID: &businessID,
		Name:       "mi-tienda.myshopify.com",
		IsActive:   true,
	}

	shopifyClient := &mockShopifyClient{
		GetOrdersFn: func(ctx context.Context, storeName, accessToken string, params *domain.GetOrdersParams) ([]domain.ShopifyOrder, string, error) {
			return nil, "", clientErr
		},
	}

	uc := newTestUseCase(&mockIntegrationService{}, shopifyClient, &mockOrderPublisher{}, &mockSyncEventPublisher{})

	_, err := uc.GetOrders(ctx, integration, "mi-tienda.myshopify.com", "token", &domain.GetOrdersParams{Limit: 250})

	if err == nil {
		t.Fatal("se esperaba error del client, se obtuvo nil")
	}
	if !errors.Is(err, clientErr) {
		t.Errorf("error incorrecto: got %v, want %v", err, clientErr)
	}
}

func TestGetOrders_ShopifyClientErrorOnSecondPage(t *testing.T) {
	ctx := context.Background()
	businessID := uint(42)
	clientErr := errors.New("shopify api timeout")

	integration := &domain.Integration{
		ID:         7,
		BusinessID: &businessID,
		Name:       "mi-tienda.myshopify.com",
		IsActive:   true,
	}

	shopifyClient := &mockShopifyClient{
		GetOrdersFn: func(ctx context.Context, storeName, accessToken string, params *domain.GetOrdersParams) ([]domain.ShopifyOrder, string, error) {
			return []domain.ShopifyOrder{
				{ExternalID: "o1", Metadata: map[string]interface{}{}},
			}, "https://shopify.com/next-page", nil
		},
		GetOrdersByURLFn: func(ctx context.Context, nextPageURL, accessToken string) ([]domain.ShopifyOrder, string, error) {
			return nil, "", clientErr
		},
	}

	publisher := &mockOrderPublisher{}
	uc := newTestUseCase(&mockIntegrationService{}, shopifyClient, publisher, &mockSyncEventPublisher{})

	total, err := uc.GetOrders(ctx, integration, "mi-tienda.myshopify.com", "token", &domain.GetOrdersParams{Limit: 250})

	if err == nil {
		t.Fatal("se esperaba error en segunda pagina, se obtuvo nil")
	}
	// La primera pagina se publico correctamente
	if total != 1 {
		t.Errorf("total deberia incluir la primera pagina: got %d, want 1", total)
	}
	if len(publisher.PublishedOrders) != 1 {
		t.Errorf("se esperaba 1 orden publicada de la primera pagina, se publicaron %d", len(publisher.PublishedOrders))
	}
}

func TestGetOrders_PublishErrorContinues(t *testing.T) {
	ctx := context.Background()
	businessID := uint(42)

	integration := &domain.Integration{
		ID:         7,
		BusinessID: &businessID,
		Name:       "mi-tienda.myshopify.com",
		IsActive:   true,
	}

	shopifyClient := &mockShopifyClient{
		GetOrdersFn: func(ctx context.Context, storeName, accessToken string, params *domain.GetOrdersParams) ([]domain.ShopifyOrder, string, error) {
			return []domain.ShopifyOrder{
				{ExternalID: "o1", Metadata: map[string]interface{}{}},
				{ExternalID: "o2", Metadata: map[string]interface{}{}},
				{ExternalID: "o3", Metadata: map[string]interface{}{}},
			}, "", nil
		},
	}

	publishCount := 0
	publisher := &mockOrderPublisher{
		PublishFn: func(ctx context.Context, order *domain.ProbabilityOrderDTO) error {
			publishCount++
			if publishCount == 2 {
				return errors.New("rabbitmq error") // Falla la segunda
			}
			return nil
		},
	}

	uc := newTestUseCase(&mockIntegrationService{}, shopifyClient, publisher, &mockSyncEventPublisher{})

	total, err := uc.GetOrders(ctx, integration, "mi-tienda.myshopify.com", "token", &domain.GetOrdersParams{Limit: 250})

	// No debe retornar error: los errores de publish se logean pero se continua
	if err != nil {
		t.Fatalf("se esperaba nil (continue on publish error), se obtuvo: %v", err)
	}
	// Solo 2 se publicaron exitosamente (la 2da fallo)
	if total != 2 {
		t.Errorf("total publicado incorrecto: got %d, want 2", total)
	}
}

func TestGetOrders_OrderFieldsAssigned(t *testing.T) {
	ctx := context.Background()
	businessID := uint(55)

	integration := &domain.Integration{
		ID:         9,
		BusinessID: &businessID,
		Name:       "mi-tienda.myshopify.com",
		IsActive:   true,
	}

	shopifyClient := &mockShopifyClient{
		GetOrdersFn: func(ctx context.Context, storeName, accessToken string, params *domain.GetOrdersParams) ([]domain.ShopifyOrder, string, error) {
			return []domain.ShopifyOrder{
				{ExternalID: "o1", OrderNumber: "#100", Metadata: map[string]interface{}{}},
			}, "", nil
		},
	}

	publisher := &mockOrderPublisher{}
	uc := newTestUseCase(&mockIntegrationService{}, shopifyClient, publisher, &mockSyncEventPublisher{})

	_, err := uc.GetOrders(ctx, integration, "mi-tienda.myshopify.com", "token", &domain.GetOrdersParams{Limit: 250})

	if err != nil {
		t.Fatalf("se esperaba nil, se obtuvo: %v", err)
	}

	published := publisher.PublishedOrders[0]
	if published.IntegrationID != 9 {
		t.Errorf("IntegrationID: got %d, want 9", published.IntegrationID)
	}
	if published.IntegrationType != "shopify" {
		t.Errorf("IntegrationType: got %q, want %q", published.IntegrationType, "shopify")
	}
	if published.BusinessID == nil || *published.BusinessID != businessID {
		t.Errorf("BusinessID: got %v, want %d", published.BusinessID, businessID)
	}
}
