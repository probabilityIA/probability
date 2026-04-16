package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/domain"
)

func TestSyncOrders_Success(t *testing.T) {
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
		GetOrdersFn: func(ctx context.Context, storeName, accessToken string, params *domain.GetOrdersParams) ([]domain.ShopifyOrder, string, error) {
			return []domain.ShopifyOrder{
				{ExternalID: "o1", Metadata: map[string]interface{}{}},
			}, "", nil
		},
	}

	publisher := &mockOrderPublisher{}
	syncPub := &mockSyncEventPublisher{}
	uc := newTestUseCase(integrationSvc, shopifyClient, publisher, syncPub)

	err := uc.SyncOrders(ctx, "7")

	if err != nil {
		t.Fatalf("se esperaba nil, se obtuvo: %v", err)
	}
	if len(publisher.PublishedOrders) != 1 {
		t.Errorf("se esperaba 1 orden publicada, se publicaron %d", len(publisher.PublishedOrders))
	}

	// Verificar SSE events: started + completed
	if len(syncPub.Events) < 2 {
		t.Fatalf("se esperaban al menos 2 SSE events (started+completed), se obtuvieron %d", len(syncPub.Events))
	}
	if syncPub.Events[0].EventType != "integration.sync.started" {
		t.Errorf("primer evento debe ser started: got %q", syncPub.Events[0].EventType)
	}
	if syncPub.Events[len(syncPub.Events)-1].EventType != "integration.sync.completed" {
		t.Errorf("ultimo evento debe ser completed: got %q", syncPub.Events[len(syncPub.Events)-1].EventType)
	}
}

func TestSyncOrdersWithParams_Success_WithOrders(t *testing.T) {
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
		GetOrdersFn: func(ctx context.Context, storeName, accessToken string, params *domain.GetOrdersParams) ([]domain.ShopifyOrder, string, error) {
			return []domain.ShopifyOrder{
				{ExternalID: "o1", Metadata: map[string]interface{}{}},
				{ExternalID: "o2", Metadata: map[string]interface{}{}},
				{ExternalID: "o3", Metadata: map[string]interface{}{}},
			}, "", nil
		},
	}

	publisher := &mockOrderPublisher{}
	syncPub := &mockSyncEventPublisher{}
	uc := newTestUseCase(integrationSvc, shopifyClient, publisher, syncPub)

	params := &domain.SyncOrdersParams{
		Status:          domain.OrderStatusAny,
		FinancialStatus: "paid",
	}

	err := uc.SyncOrdersWithParams(ctx, "7", params)

	if err != nil {
		t.Fatalf("se esperaba nil, se obtuvo: %v", err)
	}
	if len(publisher.PublishedOrders) != 3 {
		t.Errorf("se esperaban 3 ordenes publicadas, se publicaron %d", len(publisher.PublishedOrders))
	}

	// Verificar total_fetched en completed event
	completedEvent := syncPub.Events[len(syncPub.Events)-1]
	if totalFetched, ok := completedEvent.Data["total_fetched"].(int); ok {
		if totalFetched != 3 {
			t.Errorf("total_fetched incorrecto: got %d, want 3", totalFetched)
		}
	}
}

func TestSyncOrdersWithParams_GetIntegrationError(t *testing.T) {
	ctx := context.Background()
	expectedErr := errors.New("integration not found")

	integrationSvc := &mockIntegrationService{
		GetIntegrationByIDFn: func(ctx context.Context, id string) (*domain.Integration, error) {
			return nil, expectedErr
		},
	}

	syncPub := &mockSyncEventPublisher{}
	uc := newTestUseCase(integrationSvc, &mockShopifyClient{}, &mockOrderPublisher{}, syncPub)

	params := &domain.SyncOrdersParams{Status: domain.OrderStatusAny}

	err := uc.SyncOrdersWithParams(ctx, "7", params)

	if err == nil {
		t.Fatal("se esperaba error, se obtuvo nil")
	}
	if !errors.Is(err, expectedErr) {
		t.Errorf("error incorrecto: got %v, want %v", err, expectedErr)
	}
	// No debe haber SSE events cuando falla antes de parsear el integration_id
	if len(syncPub.Events) != 0 {
		t.Errorf("no deberian publicarse SSE events al fallar GetIntegration: got %d", len(syncPub.Events))
	}
}

func TestSyncOrdersWithParams_DecryptCredentialError(t *testing.T) {
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

	uc := newTestUseCase(integrationSvc, &mockShopifyClient{}, &mockOrderPublisher{}, &mockSyncEventPublisher{})

	params := &domain.SyncOrdersParams{Status: domain.OrderStatusAny}

	err := uc.SyncOrdersWithParams(ctx, "1", params)

	if err == nil {
		t.Fatal("se esperaba error de descifrado, se obtuvo nil")
	}
	if !errors.Is(err, decryptErr) {
		t.Errorf("error incorrecto: got %v, want %v", err, decryptErr)
	}
}

func TestSyncOrdersWithParams_InvalidIntegrationID(t *testing.T) {
	ctx := context.Background()

	integrationSvc := &mockIntegrationService{
		GetIntegrationByIDFn: func(ctx context.Context, id string) (*domain.Integration, error) {
			return newIntegration(1, "mi-tienda.myshopify.com"), nil
		},
	}

	uc := newTestUseCase(integrationSvc, &mockShopifyClient{}, &mockOrderPublisher{}, &mockSyncEventPublisher{})

	params := &domain.SyncOrdersParams{Status: domain.OrderStatusAny}

	err := uc.SyncOrdersWithParams(ctx, "abc", params) // ID no numerico

	if err == nil {
		t.Fatal("se esperaba error por integration_id invalido, se obtuvo nil")
	}
	if !errors.Is(err, domain.ErrInvalidIntegrationID) {
		t.Errorf("error debe envolver ErrInvalidIntegrationID: got %v", err)
	}
}

func TestSyncOrdersWithParams_GetOrdersError_PublishesFailedEvent(t *testing.T) {
	ctx := context.Background()
	businessID := uint(42)
	getOrdersErr := errors.New("shopify api 500")

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
		GetOrdersFn: func(ctx context.Context, storeName, accessToken string, params *domain.GetOrdersParams) ([]domain.ShopifyOrder, string, error) {
			return nil, "", getOrdersErr
		},
	}

	syncPub := &mockSyncEventPublisher{}
	uc := newTestUseCase(integrationSvc, shopifyClient, &mockOrderPublisher{}, syncPub)

	params := &domain.SyncOrdersParams{Status: domain.OrderStatusAny}

	err := uc.SyncOrdersWithParams(ctx, "7", params)

	if err == nil {
		t.Fatal("se esperaba error de GetOrders, se obtuvo nil")
	}

	// Verificar SSE events: started + failed
	hasStarted := false
	hasFailed := false
	for _, e := range syncPub.Events {
		if e.EventType == "integration.sync.started" {
			hasStarted = true
		}
		if e.EventType == "integration.sync.failed" {
			hasFailed = true
		}
	}
	if !hasStarted {
		t.Error("deberia haber un evento started")
	}
	if !hasFailed {
		t.Error("deberia haber un evento failed cuando GetOrders falla")
	}
}

func TestSyncOrdersWithParams_ParamFiltersApplied(t *testing.T) {
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

	var capturedParams *domain.GetOrdersParams
	shopifyClient := &mockShopifyClient{
		GetOrdersFn: func(ctx context.Context, storeName, accessToken string, params *domain.GetOrdersParams) ([]domain.ShopifyOrder, string, error) {
			capturedParams = params
			return []domain.ShopifyOrder{}, "", nil
		},
	}

	uc := newTestUseCase(integrationSvc, shopifyClient, &mockOrderPublisher{}, &mockSyncEventPublisher{})

	params := &domain.SyncOrdersParams{
		Status:            "open",
		FinancialStatus:   "paid",
		FulfillmentStatus: "shipped",
	}

	err := uc.SyncOrdersWithParams(ctx, "7", params)

	if err != nil {
		t.Fatalf("se esperaba nil, se obtuvo: %v", err)
	}
	if capturedParams == nil {
		t.Fatal("params no fueron capturados")
	}
	if capturedParams.Status != "open" {
		t.Errorf("Status: got %q, want %q", capturedParams.Status, "open")
	}
	if capturedParams.FinancialStatus != "paid" {
		t.Errorf("FinancialStatus: got %q, want %q", capturedParams.FinancialStatus, "paid")
	}
	if capturedParams.FulfillmentStatus != "shipped" {
		t.Errorf("FulfillmentStatus: got %q, want %q", capturedParams.FulfillmentStatus, "shipped")
	}
}

func TestSyncOrdersWithParams_EmptyFilterDefaults(t *testing.T) {
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

	var capturedParams *domain.GetOrdersParams
	shopifyClient := &mockShopifyClient{
		GetOrdersFn: func(ctx context.Context, storeName, accessToken string, params *domain.GetOrdersParams) ([]domain.ShopifyOrder, string, error) {
			capturedParams = params
			return []domain.ShopifyOrder{}, "", nil
		},
	}

	uc := newTestUseCase(integrationSvc, shopifyClient, &mockOrderPublisher{}, &mockSyncEventPublisher{})

	// Todos los filtros vacios
	params := &domain.SyncOrdersParams{}

	err := uc.SyncOrdersWithParams(ctx, "7", params)

	if err != nil {
		t.Fatalf("se esperaba nil, se obtuvo: %v", err)
	}
	if capturedParams == nil {
		t.Fatal("params no fueron capturados")
	}
	// Status vacio debe defaultear a "any"
	if capturedParams.Status != domain.OrderStatusAny {
		t.Errorf("Status deberia defaultear a %q: got %q", domain.OrderStatusAny, capturedParams.Status)
	}
	if capturedParams.FinancialStatus != domain.FinancialStatusAny {
		t.Errorf("FinancialStatus deberia defaultear a %q: got %q", domain.FinancialStatusAny, capturedParams.FinancialStatus)
	}
}

func TestSyncOrdersWithParams_SSEEventContainsCorrectIDs(t *testing.T) {
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
		GetOrdersFn: func(ctx context.Context, storeName, accessToken string, params *domain.GetOrdersParams) ([]domain.ShopifyOrder, string, error) {
			return []domain.ShopifyOrder{}, "", nil
		},
	}

	syncPub := &mockSyncEventPublisher{}
	uc := newTestUseCase(integrationSvc, shopifyClient, &mockOrderPublisher{}, syncPub)

	params := &domain.SyncOrdersParams{Status: domain.OrderStatusAny}

	err := uc.SyncOrdersWithParams(ctx, "7", params)

	if err != nil {
		t.Fatalf("se esperaba nil, se obtuvo: %v", err)
	}

	for _, event := range syncPub.Events {
		if event.IntegrationID != 7 {
			t.Errorf("SSE event IntegrationID: got %d, want 7", event.IntegrationID)
		}
		if event.BusinessID == nil || *event.BusinessID != businessID {
			t.Errorf("SSE event BusinessID: got %v, want %d", event.BusinessID, businessID)
		}
	}
}
