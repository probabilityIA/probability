package usecases

// Tests de los use cases de eventos de orden:
// ProcessOrderPaid, ProcessOrderFulfilled, ProcessOrderCancelled,
// ProcessOrderUpdated, ProcessOrderPartiallyFulfilled.
//
// Todos tienen la misma estructura: reciben shopDomain + *ShopifyOrder,
// buscan la integracion por external_id y publican la orden transformada.

import (
	"context"
	"errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/domain"
)

// buildOrderForEvent construye una orden Shopify minima para tests de eventos.
func buildOrderForEvent(externalID, orderNumber string) *domain.ShopifyOrder {
	return &domain.ShopifyOrder{
		ExternalID:  externalID,
		OrderNumber: orderNumber,
		TotalAmount: 200000,
		Currency:    "COP",
		Metadata:    map[string]interface{}{},
	}
}

// buildIntegrationSvc construye un mockIntegrationService que retorna la integracion esperada.
func buildIntegrationSvc(t *testing.T, expectedDomain string, integrationID uint) *mockIntegrationService {
	businessID := uint(11)
	return &mockIntegrationService{
		GetIntegrationByExternalIDFn: func(ctx context.Context, externalID string, integrationType int) (*domain.Integration, error) {
			if externalID != expectedDomain {
				t.Errorf("externalID incorrecto: got %q, want %q", externalID, expectedDomain)
			}
			return &domain.Integration{ID: integrationID, BusinessID: &businessID}, nil
		},
	}
}

// ──────────────────────────────────────────────────────────────────────────────
// ProcessOrderPaid
// ──────────────────────────────────────────────────────────────────────────────

func TestProcessOrderPaid_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	shopDomain := "tienda.myshopify.com"
	publisher := &mockOrderPublisher{}
	uc := newTestUseCase(buildIntegrationSvc(t, shopDomain, 3), &mockShopifyClient{}, publisher)
	order := buildOrderForEvent("order-paid-1", "#2001")

	// Act
	err := uc.ProcessOrderPaid(ctx, shopDomain, order)

	// Assert
	if err != nil {
		t.Fatalf("ProcessOrderPaid: se esperaba nil, se obtuvo: %v", err)
	}
	if len(publisher.PublishedOrders) != 1 {
		t.Fatalf("se esperaba 1 orden publicada, se publicaron %d", len(publisher.PublishedOrders))
	}
	if publisher.PublishedOrders[0].IntegrationType != "shopify" {
		t.Errorf("IntegrationType incorrecto: got %q", publisher.PublishedOrders[0].IntegrationType)
	}
}

func TestProcessOrderPaid_NilOrder(t *testing.T) {
	ctx := context.Background()
	uc := newTestUseCase(&mockIntegrationService{}, &mockShopifyClient{}, &mockOrderPublisher{})

	err := uc.ProcessOrderPaid(ctx, "tienda.myshopify.com", nil)

	if err == nil {
		t.Fatal("ProcessOrderPaid con nil order deberia retornar error")
	}
}

func TestProcessOrderPaid_IntegrationError(t *testing.T) {
	ctx := context.Background()
	expectedErr := errors.New("db error")

	integrationSvc := &mockIntegrationService{
		GetIntegrationByExternalIDFn: func(ctx context.Context, externalID string, integrationType int) (*domain.Integration, error) {
			return nil, expectedErr
		},
	}

	publisher := &mockOrderPublisher{}
	uc := newTestUseCase(integrationSvc, &mockShopifyClient{}, publisher)

	err := uc.ProcessOrderPaid(ctx, "tienda.myshopify.com", buildOrderForEvent("o1", "#1"))

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

func TestProcessOrderPaid_PublisherError(t *testing.T) {
	ctx := context.Background()
	publishErr := errors.New("queue unavailable")

	publisher := &mockOrderPublisher{
		PublishFn: func(ctx context.Context, order *domain.ProbabilityOrderDTO) error {
			return publishErr
		},
	}

	uc := newTestUseCase(buildIntegrationSvc(t, "tienda.myshopify.com", 1), &mockShopifyClient{}, publisher)

	err := uc.ProcessOrderPaid(ctx, "tienda.myshopify.com", buildOrderForEvent("o1", "#1"))

	if err == nil {
		t.Fatal("se esperaba error del publisher, se obtuvo nil")
	}
	if !errors.Is(err, publishErr) {
		t.Errorf("error incorrecto: got %v, want %v", err, publishErr)
	}
}

// ──────────────────────────────────────────────────────────────────────────────
// ProcessOrderFulfilled
// ──────────────────────────────────────────────────────────────────────────────

func TestProcessOrderFulfilled_Success(t *testing.T) {
	ctx := context.Background()
	shopDomain := "tienda.myshopify.com"
	publisher := &mockOrderPublisher{}
	uc := newTestUseCase(buildIntegrationSvc(t, shopDomain, 3), &mockShopifyClient{}, publisher)

	err := uc.ProcessOrderFulfilled(ctx, shopDomain, buildOrderForEvent("order-fulfilled-1", "#2002"))

	if err != nil {
		t.Fatalf("ProcessOrderFulfilled: se esperaba nil, se obtuvo: %v", err)
	}
	if len(publisher.PublishedOrders) != 1 {
		t.Fatalf("se esperaba 1 orden publicada, se publicaron %d", len(publisher.PublishedOrders))
	}
}

func TestProcessOrderFulfilled_NilOrder(t *testing.T) {
	ctx := context.Background()
	uc := newTestUseCase(&mockIntegrationService{}, &mockShopifyClient{}, &mockOrderPublisher{})

	err := uc.ProcessOrderFulfilled(ctx, "tienda.myshopify.com", nil)

	if err == nil {
		t.Fatal("ProcessOrderFulfilled con nil order deberia retornar error")
	}
}

func TestProcessOrderFulfilled_PublisherError(t *testing.T) {
	ctx := context.Background()
	publishErr := errors.New("queue full")

	publisher := &mockOrderPublisher{
		PublishFn: func(ctx context.Context, order *domain.ProbabilityOrderDTO) error {
			return publishErr
		},
	}

	uc := newTestUseCase(buildIntegrationSvc(t, "tienda.myshopify.com", 1), &mockShopifyClient{}, publisher)

	err := uc.ProcessOrderFulfilled(ctx, "tienda.myshopify.com", buildOrderForEvent("o1", "#1"))

	if err == nil {
		t.Fatal("se esperaba error del publisher, se obtuvo nil")
	}
	if !errors.Is(err, publishErr) {
		t.Errorf("error incorrecto: got %v, want %v", err, publishErr)
	}
}

// ──────────────────────────────────────────────────────────────────────────────
// ProcessOrderCancelled
// ──────────────────────────────────────────────────────────────────────────────

func TestProcessOrderCancelled_Success(t *testing.T) {
	ctx := context.Background()
	shopDomain := "tienda.myshopify.com"
	publisher := &mockOrderPublisher{}
	uc := newTestUseCase(buildIntegrationSvc(t, shopDomain, 3), &mockShopifyClient{}, publisher)

	err := uc.ProcessOrderCancelled(ctx, shopDomain, buildOrderForEvent("order-cancelled-1", "#2003"))

	if err != nil {
		t.Fatalf("ProcessOrderCancelled: se esperaba nil, se obtuvo: %v", err)
	}
	if len(publisher.PublishedOrders) != 1 {
		t.Fatalf("se esperaba 1 orden publicada, se publicaron %d", len(publisher.PublishedOrders))
	}
}

func TestProcessOrderCancelled_NilOrder(t *testing.T) {
	ctx := context.Background()
	uc := newTestUseCase(&mockIntegrationService{}, &mockShopifyClient{}, &mockOrderPublisher{})

	err := uc.ProcessOrderCancelled(ctx, "tienda.myshopify.com", nil)

	if err == nil {
		t.Fatal("ProcessOrderCancelled con nil order deberia retornar error")
	}
}

func TestProcessOrderCancelled_IntegrationError(t *testing.T) {
	ctx := context.Background()
	expectedErr := errors.New("network error")

	integrationSvc := &mockIntegrationService{
		GetIntegrationByExternalIDFn: func(ctx context.Context, externalID string, integrationType int) (*domain.Integration, error) {
			return nil, expectedErr
		},
	}

	uc := newTestUseCase(integrationSvc, &mockShopifyClient{}, &mockOrderPublisher{})

	err := uc.ProcessOrderCancelled(ctx, "tienda.myshopify.com", buildOrderForEvent("o1", "#1"))

	if err == nil {
		t.Fatal("se esperaba error, se obtuvo nil")
	}
	if !errors.Is(err, expectedErr) {
		t.Errorf("error incorrecto: got %v, want %v", err, expectedErr)
	}
}

// ──────────────────────────────────────────────────────────────────────────────
// ProcessOrderUpdated
// ──────────────────────────────────────────────────────────────────────────────

func TestProcessOrderUpdated_Success(t *testing.T) {
	ctx := context.Background()
	shopDomain := "tienda.myshopify.com"
	publisher := &mockOrderPublisher{}
	uc := newTestUseCase(buildIntegrationSvc(t, shopDomain, 3), &mockShopifyClient{}, publisher)

	err := uc.ProcessOrderUpdated(ctx, shopDomain, buildOrderForEvent("order-updated-1", "#2004"))

	if err != nil {
		t.Fatalf("ProcessOrderUpdated: se esperaba nil, se obtuvo: %v", err)
	}
	if len(publisher.PublishedOrders) != 1 {
		t.Fatalf("se esperaba 1 orden publicada, se publicaron %d", len(publisher.PublishedOrders))
	}
}

func TestProcessOrderUpdated_NilOrder(t *testing.T) {
	ctx := context.Background()
	uc := newTestUseCase(&mockIntegrationService{}, &mockShopifyClient{}, &mockOrderPublisher{})

	err := uc.ProcessOrderUpdated(ctx, "tienda.myshopify.com", nil)

	if err == nil {
		t.Fatal("ProcessOrderUpdated con nil order deberia retornar error")
	}
}

func TestProcessOrderUpdated_BusinessIDAndIntegrationTypeSet(t *testing.T) {
	// Verifica que se asignan correctamente los campos de integracion a la orden
	ctx := context.Background()
	shopDomain := "tienda.myshopify.com"
	businessID := uint(55)
	integrationID := uint(9)

	integrationSvc := &mockIntegrationService{
		GetIntegrationByExternalIDFn: func(ctx context.Context, externalID string, integrationType int) (*domain.Integration, error) {
			return &domain.Integration{ID: integrationID, BusinessID: &businessID}, nil
		},
	}

	publisher := &mockOrderPublisher{}
	uc := newTestUseCase(integrationSvc, &mockShopifyClient{}, publisher)

	order := buildOrderForEvent("o1", "#1")

	err := uc.ProcessOrderUpdated(ctx, shopDomain, order)

	if err != nil {
		t.Fatalf("se esperaba nil, se obtuvo: %v", err)
	}
	if order.BusinessID == nil || *order.BusinessID != businessID {
		t.Errorf("BusinessID no asignado correctamente: got %v, want %d", order.BusinessID, businessID)
	}
	if order.IntegrationID != integrationID {
		t.Errorf("IntegrationID no asignado correctamente: got %d, want %d", order.IntegrationID, integrationID)
	}
	if order.IntegrationType != "shopify" {
		t.Errorf("IntegrationType incorrecto: got %q, want %q", order.IntegrationType, "shopify")
	}
}

// ──────────────────────────────────────────────────────────────────────────────
// ProcessOrderPartiallyFulfilled
// ──────────────────────────────────────────────────────────────────────────────

func TestProcessOrderPartiallyFulfilled_Success(t *testing.T) {
	ctx := context.Background()
	shopDomain := "tienda.myshopify.com"
	publisher := &mockOrderPublisher{}
	uc := newTestUseCase(buildIntegrationSvc(t, shopDomain, 3), &mockShopifyClient{}, publisher)

	err := uc.ProcessOrderPartiallyFulfilled(ctx, shopDomain, buildOrderForEvent("order-partial-1", "#2005"))

	if err != nil {
		t.Fatalf("ProcessOrderPartiallyFulfilled: se esperaba nil, se obtuvo: %v", err)
	}
	if len(publisher.PublishedOrders) != 1 {
		t.Fatalf("se esperaba 1 orden publicada, se publicaron %d", len(publisher.PublishedOrders))
	}
}

func TestProcessOrderPartiallyFulfilled_NilOrder(t *testing.T) {
	ctx := context.Background()
	uc := newTestUseCase(&mockIntegrationService{}, &mockShopifyClient{}, &mockOrderPublisher{})

	err := uc.ProcessOrderPartiallyFulfilled(ctx, "tienda.myshopify.com", nil)

	if err == nil {
		t.Fatal("ProcessOrderPartiallyFulfilled con nil order deberia retornar error")
	}
}

func TestProcessOrderPartiallyFulfilled_PublisherError(t *testing.T) {
	ctx := context.Background()
	publishErr := errors.New("broker error")

	publisher := &mockOrderPublisher{
		PublishFn: func(ctx context.Context, order *domain.ProbabilityOrderDTO) error {
			return publishErr
		},
	}

	uc := newTestUseCase(buildIntegrationSvc(t, "tienda.myshopify.com", 1), &mockShopifyClient{}, publisher)

	err := uc.ProcessOrderPartiallyFulfilled(ctx, "tienda.myshopify.com", buildOrderForEvent("o1", "#1"))

	if err == nil {
		t.Fatal("se esperaba error del publisher, se obtuvo nil")
	}
	if !errors.Is(err, publishErr) {
		t.Errorf("error incorrecto: got %v, want %v", err, publishErr)
	}
}
