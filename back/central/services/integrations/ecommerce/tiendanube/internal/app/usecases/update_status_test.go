package usecases

import (
	"context"
	"testing"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiendanube/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

type patchRegistrado struct {
	FulfillmentID string
	Status        string
	Tracking      *domain.TrackingInfo
}

type fakeFulfillmentClient struct {
	domain.ITiendanubeClient
	fulfillments []domain.FulfillmentOrder
	patches      []patchRegistrado
	eventos      []domain.TrackingEvent
	cancelada    string
	listarErr    error
}

func (f *fakeFulfillmentClient) ListFulfillmentOrders(ctx context.Context, cred domain.Credential, orderID string) ([]domain.FulfillmentOrder, error) {
	if f.listarErr != nil {
		return nil, f.listarErr
	}
	return f.fulfillments, nil
}

func (f *fakeFulfillmentClient) UpdateFulfillmentOrder(ctx context.Context, cred domain.Credential, orderID, fulfillmentOrderID, status string, tracking *domain.TrackingInfo) error {
	f.patches = append(f.patches, patchRegistrado{FulfillmentID: fulfillmentOrderID, Status: status, Tracking: tracking})
	return nil
}

func (f *fakeFulfillmentClient) CreateTrackingEvent(ctx context.Context, cred domain.Credential, orderID, fulfillmentOrderID string, event domain.TrackingEvent) error {
	f.eventos = append(f.eventos, event)
	return nil
}

func (f *fakeFulfillmentClient) CancelOrder(ctx context.Context, cred domain.Credential, orderID string) error {
	f.cancelada = orderID
	return nil
}

func integracionConSync(activo bool) *domain.Integration {
	return &domain.Integration{
		ID:      77,
		StoreID: "8126740",
		BaseURL: "https://api.tiendanube.com/v1",
		Config: map[string]interface{}{
			"store_id":                     "8126740",
			domain.ConfigStatusSyncEnabled: activo,
		},
	}
}

func nuevoCasoEstados(t *testing.T, activo bool, estadoInicial string) (*fakeFulfillmentClient, ITiendanubeUseCase) {
	t.Helper()
	cliente := &fakeFulfillmentClient{
		fulfillments: []domain.FulfillmentOrder{{ID: "fo-1", Number: 1, Status: estadoInicial}},
	}
	uc := New(cliente, &fakeService{integration: integracionConSync(activo)}, nil, &fakeRepo{}, nil, log.New())
	return cliente, uc
}

func actualizacion(status string) domain.OrderStatusUpdate {
	return domain.OrderStatusUpdate{
		ExternalOrderID: "2051703380",
		Status:          status,
		TrackingNumber:  "RM123456789CO",
		TrackingURL:     "https://rastreo.example/RM123456789CO",
		Carrier:         "Interrapidisimo",
		City:            "Bogota",
		HappenedAt:      time.Date(2026, 8, 22, 15, 0, 0, 0, time.UTC),
	}
}

func TestPasosHasta(t *testing.T) {
	casos := []struct {
		nombre   string
		actual   string
		objetivo string
		esperado []string
	}{
		{"de cero a entregado", domain.FulfillmentUnpacked, domain.FulfillmentDelivered, []string{domain.FulfillmentPacked, domain.FulfillmentDispatched, domain.FulfillmentDelivered}},
		{"de empacado a despachado", domain.FulfillmentPacked, domain.FulfillmentDispatched, []string{domain.FulfillmentDispatched}},
		{"no retrocede", domain.FulfillmentDispatched, domain.FulfillmentPacked, nil},
		{"mismo estado no repite", domain.FulfillmentPacked, domain.FulfillmentPacked, nil},
		{"estado desconocido arranca en unpacked", "loquesea", domain.FulfillmentPacked, []string{domain.FulfillmentPacked}},
		{"minusculas del canal", "packed", domain.FulfillmentDispatched, []string{domain.FulfillmentDispatched}},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			got := pasosHasta(caso.actual, caso.objetivo)
			if len(got) != len(caso.esperado) {
				t.Fatalf("pasosHasta = %v, se esperaba %v", got, caso.esperado)
			}
			for i := range got {
				if got[i] != caso.esperado[i] {
					t.Fatalf("pasosHasta = %v, se esperaba %v", got, caso.esperado)
				}
			}
		})
	}
}

func TestUpdateOrderStatusOmiteSiElSyncEstaApagado(t *testing.T) {
	cliente, uc := nuevoCasoEstados(t, false, domain.FulfillmentUnpacked)

	if err := uc.UpdateOrderStatus(context.Background(), testIntegrationID, actualizacion("in_transit")); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if len(cliente.patches) != 0 || len(cliente.eventos) != 0 {
		t.Fatalf("con el sync apagado no debia tocarse el canal")
	}
}

func TestUpdateOrderStatusOmiteEstadoSinHomologacion(t *testing.T) {
	cliente, uc := nuevoCasoEstados(t, true, domain.FulfillmentUnpacked)

	if err := uc.UpdateOrderStatus(context.Background(), testIntegrationID, actualizacion("on_hold")); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if len(cliente.patches) != 0 || len(cliente.eventos) != 0 {
		t.Fatalf("un estado sin homologacion no debia tocar el canal")
	}
}

func TestUpdateOrderStatusEnTransitoAvanzaYRegistraEvento(t *testing.T) {
	cliente, uc := nuevoCasoEstados(t, true, domain.FulfillmentUnpacked)

	if err := uc.UpdateOrderStatus(context.Background(), testIntegrationID, actualizacion("in_transit")); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	if len(cliente.patches) != 2 {
		t.Fatalf("se esperaban 2 patches, hubo %d: %+v", len(cliente.patches), cliente.patches)
	}
	if cliente.patches[0].Status != domain.FulfillmentPacked || cliente.patches[1].Status != domain.FulfillmentDispatched {
		t.Fatalf("la progresion no fue PACKED -> DISPATCHED: %+v", cliente.patches)
	}
	if cliente.patches[0].Tracking != nil {
		t.Fatal("el tracking no debia ir en el paso PACKED")
	}
	if cliente.patches[1].Tracking == nil || cliente.patches[1].Tracking.Code != "RM123456789CO" {
		t.Fatalf("el tracking debia viajar en DISPATCHED: %+v", cliente.patches[1].Tracking)
	}
	if !cliente.patches[1].Tracking.NotifyCustomer {
		t.Fatal("se debia pedir notificar al cliente")
	}

	if len(cliente.eventos) != 1 || cliente.eventos[0].Status != domain.TrackingInTransit {
		t.Fatalf("se esperaba un evento in_transit: %+v", cliente.eventos)
	}
	if cliente.eventos[0].Address != "Bogota" {
		t.Fatalf("el evento debia llevar la ciudad: %+v", cliente.eventos[0])
	}
	if cliente.eventos[0].HappenedAt.IsZero() {
		t.Fatal("el evento debia llevar la fecha del cambio de estado")
	}
}

func TestUpdateOrderStatusEntregadoDesdeDespachado(t *testing.T) {
	cliente, uc := nuevoCasoEstados(t, true, domain.FulfillmentDispatched)

	if err := uc.UpdateOrderStatus(context.Background(), testIntegrationID, actualizacion("delivered")); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	if len(cliente.patches) != 1 || cliente.patches[0].Status != domain.FulfillmentDelivered {
		t.Fatalf("se esperaba un unico patch a DELIVERED: %+v", cliente.patches)
	}
	if len(cliente.eventos) != 1 || cliente.eventos[0].Status != domain.TrackingDelivered {
		t.Fatalf("se esperaba el evento delivered: %+v", cliente.eventos)
	}
}

func TestUpdateOrderStatusNovedadNoRetrocedeYSoloRegistraEvento(t *testing.T) {
	cliente, uc := nuevoCasoEstados(t, true, domain.FulfillmentDispatched)

	if err := uc.UpdateOrderStatus(context.Background(), testIntegrationID, actualizacion("delivery_novelty")); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	if len(cliente.patches) != 0 {
		t.Fatalf("ya estaba despachado, no debia haber patches: %+v", cliente.patches)
	}
	if len(cliente.eventos) != 1 || cliente.eventos[0].Status != domain.TrackingAttemptFailed {
		t.Fatalf("se esperaba el evento delivery_attempt_failed: %+v", cliente.eventos)
	}
}

func TestUpdateOrderStatusCancelada(t *testing.T) {
	cliente, uc := nuevoCasoEstados(t, true, domain.FulfillmentUnpacked)

	if err := uc.UpdateOrderStatus(context.Background(), testIntegrationID, actualizacion("cancelled")); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	if cliente.cancelada != "2051703380" {
		t.Fatalf("se esperaba cancelar la orden, cancelada=%q", cliente.cancelada)
	}
	if len(cliente.patches) != 0 || len(cliente.eventos) != 0 {
		t.Fatal("cancelar no debia tocar fulfillment ni eventos")
	}
}

func TestUpdateOrderStatusSinGuiaNoMandaTracking(t *testing.T) {
	cliente, uc := nuevoCasoEstados(t, true, domain.FulfillmentPacked)

	update := actualizacion("in_transit")
	update.TrackingNumber = ""
	update.TrackingURL = ""

	if err := uc.UpdateOrderStatus(context.Background(), testIntegrationID, update); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	if len(cliente.patches) != 1 || cliente.patches[0].Tracking != nil {
		t.Fatalf("sin guia no debia viajar tracking_info: %+v", cliente.patches)
	}
}

func TestUpdateOrderStatusSinFulfillmentOrders(t *testing.T) {
	cliente := &fakeFulfillmentClient{}
	uc := New(cliente, &fakeService{integration: integracionConSync(true)}, nil, &fakeRepo{}, nil, log.New())

	if err := uc.UpdateOrderStatus(context.Background(), testIntegrationID, actualizacion("in_transit")); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if len(cliente.patches) != 0 || len(cliente.eventos) != 0 {
		t.Fatal("sin fulfillment orders no hay nada que actualizar")
	}
}

func TestUpdateOrderStatusSinIDExterno(t *testing.T) {
	_, uc := nuevoCasoEstados(t, true, domain.FulfillmentUnpacked)

	update := actualizacion("in_transit")
	update.ExternalOrderID = "  "

	if err := uc.UpdateOrderStatus(context.Background(), testIntegrationID, update); err != domain.ErrMissingExternalOrder {
		t.Fatalf("se esperaba ErrMissingExternalOrder, hubo %v", err)
	}
}
