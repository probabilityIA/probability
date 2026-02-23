package app

import (
	"context"
	"testing"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/events/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/events/internal/mocks"
)

// ─────────────────────────────────────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────────────────────────────────────

func newUintPtr(v uint) *uint { return &v }

func newFloat64Ptr(v float64) *float64 { return &v }

// buildService crea un IntegrationEventService listo para usar en tests,
// junto con el mock del publisher para inspeccionar los eventos publicados.
func buildService() (*mocks.IntegrationEventPublisherMock, domain.IIntegrationEventService) {
	pub := &mocks.IntegrationEventPublisherMock{}
	svc := NewIntegrationEventService(pub)
	return pub, svc
}

// ─────────────────────────────────────────────────────────────────────────────
// PublishSyncOrderCreated
// ─────────────────────────────────────────────────────────────────────────────

func TestIntegrationEventService_PublishSyncOrderCreated_PublicaEvento(t *testing.T) {
	pub, svc := buildService()
	ctx := context.Background()

	businessID := newUintPtr(42)
	total := newFloat64Ptr(150.75)
	createdAt := time.Now().Truncate(time.Second)

	data := domain.SyncOrderCreatedEvent{
		OrderID:       "order-001",
		OrderNumber:   "#1001",
		ExternalID:    "ext-001",
		Platform:      "shopify",
		CustomerEmail: "cliente@example.com",
		TotalAmount:   total,
		Currency:      "COP",
		Status:        "paid",
		CreatedAt:     createdAt,
	}

	err := svc.PublishSyncOrderCreated(ctx, 10, businessID, data)

	// Assert: no debe retornar error
	if err != nil {
		t.Fatalf("esperaba nil, obtuvo error: %v", err)
	}

	// Assert: el publisher debe haber recibido exactamente un evento
	events := pub.PublishedEvents()
	if len(events) != 1 {
		t.Fatalf("esperaba 1 evento publicado, obtuvo %d", len(events))
	}

	event := events[0]

	if event.Type != domain.IntegrationEventTypeSyncOrderCreated {
		t.Errorf("tipo de evento incorrecto: got %q, want %q", event.Type, domain.IntegrationEventTypeSyncOrderCreated)
	}
	if event.IntegrationID != 10 {
		t.Errorf("integration_id incorrecto: got %d, want 10", event.IntegrationID)
	}
	if event.BusinessID == nil || *event.BusinessID != 42 {
		t.Errorf("business_id incorrecto: got %v, want ptr(42)", event.BusinessID)
	}
	if event.ID == "" {
		t.Error("el evento debe tener un ID no vacío")
	}
	if event.Timestamp.IsZero() {
		t.Error("el evento debe tener un timestamp no cero")
	}

	// Verificar metadata esperada
	assertMetadataStr(t, event.Metadata, "order_id", data.OrderID)
	assertMetadataStr(t, event.Metadata, "order_number", data.OrderNumber)
	assertMetadataStr(t, event.Metadata, "external_id", data.ExternalID)
	assertMetadataStr(t, event.Metadata, "platform", data.Platform)
	assertMetadataStr(t, event.Metadata, "status", data.Status)
}

func TestIntegrationEventService_PublishSyncOrderCreated_BusinessIDNil(t *testing.T) {
	pub, svc := buildService()
	ctx := context.Background()

	data := domain.SyncOrderCreatedEvent{
		OrderID:     "order-002",
		OrderNumber: "#1002",
		CreatedAt:   time.Now(),
	}

	err := svc.PublishSyncOrderCreated(ctx, 5, nil, data)

	if err != nil {
		t.Fatalf("esperaba nil, obtuvo: %v", err)
	}

	events := pub.PublishedEvents()
	if len(events) != 1 {
		t.Fatalf("esperaba 1 evento, obtuvo %d", len(events))
	}
	if events[0].BusinessID != nil {
		t.Errorf("business_id deberia ser nil, obtuvo %v", events[0].BusinessID)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// PublishSyncOrderUpdated
// ─────────────────────────────────────────────────────────────────────────────

func TestIntegrationEventService_PublishSyncOrderUpdated_PublicaEvento(t *testing.T) {
	pub, svc := buildService()
	ctx := context.Background()

	businessID := newUintPtr(7)
	data := domain.SyncOrderUpdatedEvent{
		OrderID:     "order-003",
		OrderNumber: "#1003",
		ExternalID:  "ext-003",
		Platform:    "shopify",
		Status:      "fulfilled",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := svc.PublishSyncOrderUpdated(ctx, 20, businessID, data)

	if err != nil {
		t.Fatalf("esperaba nil, obtuvo: %v", err)
	}

	events := pub.PublishedEvents()
	if len(events) != 1 {
		t.Fatalf("esperaba 1 evento, obtuvo %d", len(events))
	}

	event := events[0]
	if event.Type != domain.IntegrationEventTypeSyncOrderUpdated {
		t.Errorf("tipo incorrecto: got %q, want %q", event.Type, domain.IntegrationEventTypeSyncOrderUpdated)
	}
	if event.IntegrationID != 20 {
		t.Errorf("integration_id incorrecto: got %d, want 20", event.IntegrationID)
	}

	assertMetadataStr(t, event.Metadata, "order_id", data.OrderID)
	assertMetadataStr(t, event.Metadata, "platform", data.Platform)
	assertMetadataStr(t, event.Metadata, "status", data.Status)
}

// ─────────────────────────────────────────────────────────────────────────────
// PublishSyncOrderRejected
// ─────────────────────────────────────────────────────────────────────────────

func TestIntegrationEventService_PublishSyncOrderRejected_PublicaEvento(t *testing.T) {
	pub, svc := buildService()
	ctx := context.Background()

	businessID := newUintPtr(99)
	data := domain.SyncOrderRejectedEvent{
		OrderID:     "order-004",
		OrderNumber: "#1004",
		ExternalID:  "ext-004",
		Platform:    "shopify",
		Reason:      "producto no encontrado",
		Error:       "SKU inexistente",
		RejectedAt:  time.Now(),
	}

	err := svc.PublishSyncOrderRejected(ctx, 30, businessID, data)

	if err != nil {
		t.Fatalf("esperaba nil, obtuvo: %v", err)
	}

	events := pub.PublishedEvents()
	if len(events) != 1 {
		t.Fatalf("esperaba 1 evento, obtuvo %d", len(events))
	}

	event := events[0]
	if event.Type != domain.IntegrationEventTypeSyncOrderRejected {
		t.Errorf("tipo incorrecto: got %q, want %q", event.Type, domain.IntegrationEventTypeSyncOrderRejected)
	}

	assertMetadataStr(t, event.Metadata, "reason", data.Reason)
	assertMetadataStr(t, event.Metadata, "error", data.Error)
	assertMetadataStr(t, event.Metadata, "order_number", data.OrderNumber)
}

// ─────────────────────────────────────────────────────────────────────────────
// PublishSyncStarted
// ─────────────────────────────────────────────────────────────────────────────

func TestIntegrationEventService_PublishSyncStarted_PublicaEvento(t *testing.T) {
	pub, svc := buildService()
	ctx := context.Background()

	businessID := newUintPtr(1)
	data := domain.SyncStartedEvent{
		IntegrationID:   50,
		IntegrationType: "shopify",
		StartedAt:       time.Now(),
	}

	err := svc.PublishSyncStarted(ctx, 50, businessID, data)

	if err != nil {
		t.Fatalf("esperaba nil, obtuvo: %v", err)
	}

	events := pub.PublishedEvents()
	if len(events) != 1 {
		t.Fatalf("esperaba 1 evento, obtuvo %d", len(events))
	}

	event := events[0]
	if event.Type != domain.IntegrationEventTypeSyncStarted {
		t.Errorf("tipo incorrecto: got %q, want %q", event.Type, domain.IntegrationEventTypeSyncStarted)
	}

	assertMetadataStr(t, event.Metadata, "integration_type", data.IntegrationType)
}

// ─────────────────────────────────────────────────────────────────────────────
// PublishSyncCompleted
// ─────────────────────────────────────────────────────────────────────────────

func TestIntegrationEventService_PublishSyncCompleted_PublicaEvento(t *testing.T) {
	pub, svc := buildService()
	ctx := context.Background()

	businessID := newUintPtr(2)
	data := domain.SyncCompletedEvent{
		IntegrationID:   60,
		IntegrationType: "shopify",
		TotalOrders:     100,
		CreatedOrders:   80,
		UpdatedOrders:   15,
		RejectedOrders:  5,
		Duration:        2 * time.Minute,
		CompletedAt:     time.Now(),
	}

	err := svc.PublishSyncCompleted(ctx, 60, businessID, data)

	if err != nil {
		t.Fatalf("esperaba nil, obtuvo: %v", err)
	}

	events := pub.PublishedEvents()
	if len(events) != 1 {
		t.Fatalf("esperaba 1 evento, obtuvo %d", len(events))
	}

	event := events[0]
	if event.Type != domain.IntegrationEventTypeSyncCompleted {
		t.Errorf("tipo incorrecto: got %q, want %q", event.Type, domain.IntegrationEventTypeSyncCompleted)
	}

	// Verificar contadores en metadata
	assertMetadataInt(t, event.Metadata, "total_orders", data.TotalOrders)
	assertMetadataInt(t, event.Metadata, "created_orders", data.CreatedOrders)
	assertMetadataInt(t, event.Metadata, "updated_orders", data.UpdatedOrders)
	assertMetadataInt(t, event.Metadata, "rejected_orders", data.RejectedOrders)
	assertMetadataStr(t, event.Metadata, "integration_type", data.IntegrationType)
}

// ─────────────────────────────────────────────────────────────────────────────
// PublishSyncFailed
// ─────────────────────────────────────────────────────────────────────────────

func TestIntegrationEventService_PublishSyncFailed_PublicaEvento(t *testing.T) {
	pub, svc := buildService()
	ctx := context.Background()

	businessID := newUintPtr(3)
	data := domain.SyncFailedEvent{
		IntegrationID:   70,
		IntegrationType: "shopify",
		Error:           "timeout al conectar con Shopify",
		FailedAt:        time.Now(),
	}

	err := svc.PublishSyncFailed(ctx, 70, businessID, data)

	if err != nil {
		t.Fatalf("esperaba nil, obtuvo: %v", err)
	}

	events := pub.PublishedEvents()
	if len(events) != 1 {
		t.Fatalf("esperaba 1 evento, obtuvo %d", len(events))
	}

	event := events[0]
	if event.Type != domain.IntegrationEventTypeSyncFailed {
		t.Errorf("tipo incorrecto: got %q, want %q", event.Type, domain.IntegrationEventTypeSyncFailed)
	}

	assertMetadataStr(t, event.Metadata, "integration_type", data.IntegrationType)
	assertMetadataStr(t, event.Metadata, "error", data.Error)
}

// ─────────────────────────────────────────────────────────────────────────────
// Publicaciones multiples e independencia entre llamadas
// ─────────────────────────────────────────────────────────────────────────────

func TestIntegrationEventService_VariosEventos_SeAcumulanEnOrden(t *testing.T) {
	pub, svc := buildService()
	ctx := context.Background()
	businessID := newUintPtr(1)

	// Publicar tres eventos distintos
	_ = svc.PublishSyncStarted(ctx, 1, businessID, domain.SyncStartedEvent{IntegrationType: "shopify"})
	_ = svc.PublishSyncOrderCreated(ctx, 1, businessID, domain.SyncOrderCreatedEvent{OrderID: "o1", CreatedAt: time.Now()})
	_ = svc.PublishSyncOrderRejected(ctx, 1, businessID, domain.SyncOrderRejectedEvent{OrderID: "o2", RejectedAt: time.Now()})

	events := pub.PublishedEvents()
	if len(events) != 3 {
		t.Fatalf("esperaba 3 eventos, obtuvo %d", len(events))
	}

	expectedTypes := []domain.IntegrationEventType{
		domain.IntegrationEventTypeSyncStarted,
		domain.IntegrationEventTypeSyncOrderCreated,
		domain.IntegrationEventTypeSyncOrderRejected,
	}
	for i, want := range expectedTypes {
		if events[i].Type != want {
			t.Errorf("evento[%d]: tipo %q, want %q", i, events[i].Type, want)
		}
	}
}

func TestIntegrationEventService_CadaEventoTieneIDUnico(t *testing.T) {
	pub, svc := buildService()
	ctx := context.Background()
	businessID := newUintPtr(1)

	for i := 0; i < 5; i++ {
		_ = svc.PublishSyncStarted(ctx, uint(i), businessID, domain.SyncStartedEvent{IntegrationType: "shopify"})
	}

	events := pub.PublishedEvents()
	if len(events) != 5 {
		t.Fatalf("esperaba 5 eventos, obtuvo %d", len(events))
	}

	seen := make(map[string]bool)
	for _, e := range events {
		if e.ID == "" {
			t.Error("evento tiene ID vacio")
		}
		if seen[e.ID] {
			t.Errorf("ID duplicado detectado: %q", e.ID)
		}
		seen[e.ID] = true
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Stubs: SetSyncState, GetSyncState, DeleteSyncState, IncrementSyncCounter
// ─────────────────────────────────────────────────────────────────────────────

func TestIntegrationEventService_SetSyncState_RetornaNil(t *testing.T) {
	_, svc := buildService()
	state := domain.SyncState{IntegrationID: 1, Status: domain.SyncStatusInProgress}

	err := svc.SetSyncState(context.Background(), 1, state)
	if err != nil {
		t.Errorf("esperaba nil, obtuvo: %v", err)
	}
}

func TestIntegrationEventService_GetSyncState_RetornaNilNil(t *testing.T) {
	_, svc := buildService()

	result, err := svc.GetSyncState(context.Background(), 1)
	if err != nil {
		t.Errorf("esperaba nil error, obtuvo: %v", err)
	}
	if result != nil {
		t.Errorf("esperaba nil result, obtuvo: %v", result)
	}
}

func TestIntegrationEventService_DeleteSyncState_RetornaNil(t *testing.T) {
	_, svc := buildService()

	err := svc.DeleteSyncState(context.Background(), 1)
	if err != nil {
		t.Errorf("esperaba nil, obtuvo: %v", err)
	}
}

func TestIntegrationEventService_IncrementSyncCounter_RetornaNil(t *testing.T) {
	_, svc := buildService()

	err := svc.IncrementSyncCounter(context.Background(), 1, "created")
	if err != nil {
		t.Errorf("esperaba nil, obtuvo: %v", err)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Helpers de aserción de metadata
// ─────────────────────────────────────────────────────────────────────────────

// assertMetadataStr verifica que metadata[key] sea el string esperado.
func assertMetadataStr(t *testing.T, metadata map[string]interface{}, key, want string) {
	t.Helper()
	val, ok := metadata[key]
	if !ok {
		t.Errorf("metadata[%q] no encontrado", key)
		return
	}
	got, ok := val.(string)
	if !ok {
		t.Errorf("metadata[%q] tipo incorrecto: %T", key, val)
		return
	}
	if got != want {
		t.Errorf("metadata[%q] = %q, want %q", key, got, want)
	}
}

// assertMetadataInt verifica que metadata[key] sea el int esperado.
func assertMetadataInt(t *testing.T, metadata map[string]interface{}, key string, want int) {
	t.Helper()
	val, ok := metadata[key]
	if !ok {
		t.Errorf("metadata[%q] no encontrado", key)
		return
	}
	got, ok := val.(int)
	if !ok {
		t.Errorf("metadata[%q] tipo incorrecto: %T (want int)", key, val)
		return
	}
	if got != want {
		t.Errorf("metadata[%q] = %d, want %d", key, got, want)
	}
}
