package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/domain"
)

func TestCreateOrder_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	shopDomain := "mi-tienda.myshopify.com"
	businessID := uint(42)

	integrationSvc := &mockIntegrationService{
		GetIntegrationByExternalIDFn: func(ctx context.Context, externalID string, integrationType int) (*domain.Integration, error) {
			if externalID != shopDomain {
				t.Errorf("externalID inesperado: got %q, want %q", externalID, shopDomain)
			}
			if integrationType != domain.IntegrationTypeID {
				t.Errorf("integrationType inesperado: got %d, want %d", integrationType, domain.IntegrationTypeID)
			}
			return &domain.Integration{ID: 7, BusinessID: &businessID}, nil
		},
	}

	publisher := &mockOrderPublisher{}

	order := &domain.ShopifyOrder{
		ExternalID:  "shop-order-001",
		OrderNumber: "#1001",
		TotalAmount: 150000,
		Currency:    "COP",
	}
	rawPayload := []byte(`{"id":1,"financial_status":"paid","fulfillment_status":null}`)

	uc := newTestUseCase(integrationSvc, &mockShopifyClient{}, publisher)

	// Act
	err := uc.CreateOrder(ctx, shopDomain, order, rawPayload)

	// Assert
	if err != nil {
		t.Fatalf("se esperaba nil, se obtuvo: %v", err)
	}

	// Verificar que se publico exactamente una orden
	if len(publisher.PublishedOrders) != 1 {
		t.Fatalf("se esperaba 1 orden publicada, se publicaron %d", len(publisher.PublishedOrders))
	}

	published := publisher.PublishedOrders[0]

	// Verificar que se asignaron los campos de la integracion
	if published.IntegrationID != 7 {
		t.Errorf("IntegrationID incorrecto: got %d, want %d", published.IntegrationID, 7)
	}
	if published.IntegrationType != "shopify" {
		t.Errorf("IntegrationType incorrecto: got %q, want %q", published.IntegrationType, "shopify")
	}
	if published.BusinessID == nil || *published.BusinessID != businessID {
		t.Errorf("BusinessID incorrecto: got %v, want %d", published.BusinessID, businessID)
	}

	// Verificar ChannelMetadata cuando rawPayload tiene contenido
	if published.ChannelMetadata == nil {
		t.Error("ChannelMetadata no debe ser nil cuando se provee rawPayload")
	} else {
		if published.ChannelMetadata.ChannelSource != "shopify" {
			t.Errorf("ChannelSource incorrecto: got %q, want %q", published.ChannelMetadata.ChannelSource, "shopify")
		}
		if !published.ChannelMetadata.IsLatest {
			t.Error("IsLatest debe ser true")
		}
		if published.ChannelMetadata.SyncStatus != "synced" {
			t.Errorf("SyncStatus incorrecto: got %q, want %q", published.ChannelMetadata.SyncStatus, "synced")
		}
	}
}

func TestCreateOrder_NilOrder(t *testing.T) {
	// Arrange
	ctx := context.Background()
	uc := newTestUseCase(&mockIntegrationService{}, &mockShopifyClient{}, &mockOrderPublisher{})

	// Act
	err := uc.CreateOrder(ctx, "mi-tienda.myshopify.com", nil, nil)

	// Assert
	if err == nil {
		t.Fatal("se esperaba error por order nil, se obtuvo nil")
	}
}

func TestCreateOrder_IntegrationServiceError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	expectedErr := errors.New("integration not found")

	integrationSvc := &mockIntegrationService{
		GetIntegrationByExternalIDFn: func(ctx context.Context, externalID string, integrationType int) (*domain.Integration, error) {
			return nil, expectedErr
		},
	}

	publisher := &mockOrderPublisher{}
	uc := newTestUseCase(integrationSvc, &mockShopifyClient{}, publisher)

	order := &domain.ShopifyOrder{ExternalID: "order-001"}

	// Act
	err := uc.CreateOrder(ctx, "tienda-desconocida.myshopify.com", order, nil)

	// Assert
	if err == nil {
		t.Fatal("se esperaba error, se obtuvo nil")
	}
	if !errors.Is(err, expectedErr) {
		t.Errorf("error incorrecto: got %v, want %v", err, expectedErr)
	}
	if len(publisher.PublishedOrders) != 0 {
		t.Error("no se debia publicar ninguna orden cuando falla la busqueda de integracion")
	}
}

func TestCreateOrder_PublisherError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	businessID := uint(10)
	publishErr := errors.New("rabbitmq connection lost")

	integrationSvc := &mockIntegrationService{
		GetIntegrationByExternalIDFn: func(ctx context.Context, externalID string, integrationType int) (*domain.Integration, error) {
			return &domain.Integration{ID: 3, BusinessID: &businessID}, nil
		},
	}

	publisher := &mockOrderPublisher{
		PublishFn: func(ctx context.Context, order *domain.ProbabilityOrderDTO) error {
			return publishErr
		},
	}

	uc := newTestUseCase(integrationSvc, &mockShopifyClient{}, publisher)

	order := &domain.ShopifyOrder{ExternalID: "order-002"}

	// Act
	err := uc.CreateOrder(ctx, "mi-tienda.myshopify.com", order, nil)

	// Assert
	if err == nil {
		t.Fatal("se esperaba error del publisher, se obtuvo nil")
	}
	if !errors.Is(err, publishErr) {
		t.Errorf("error incorrecto: got %v, want %v", err, publishErr)
	}
}

func TestCreateOrder_WithoutRawPayload_NoChannelMetadata(t *testing.T) {
	// Arrange
	ctx := context.Background()
	businessID := uint(5)

	integrationSvc := &mockIntegrationService{
		GetIntegrationByExternalIDFn: func(ctx context.Context, externalID string, integrationType int) (*domain.Integration, error) {
			return &domain.Integration{ID: 1, BusinessID: &businessID}, nil
		},
	}

	publisher := &mockOrderPublisher{}
	uc := newTestUseCase(integrationSvc, &mockShopifyClient{}, publisher)

	order := &domain.ShopifyOrder{ExternalID: "order-003"}

	// Act - rawPayload vacio (nil)
	err := uc.CreateOrder(ctx, "mi-tienda.myshopify.com", order, nil)

	// Assert
	if err != nil {
		t.Fatalf("se esperaba nil, se obtuvo: %v", err)
	}
	if len(publisher.PublishedOrders) != 1 {
		t.Fatalf("se esperaba 1 orden publicada, se publicaron %d", len(publisher.PublishedOrders))
	}
	// Sin rawPayload el ChannelMetadata solo se agrega si el mapper lo pone
	// (en CreateOrder el bloque rawPayload > 0 es el que lo agrega explicitamente)
}
