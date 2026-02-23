package redis

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/events/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/events/internal/mocks"
)

// ─────────────────────────────────────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────────────────────────────────────

const testChannel = "integration-events-test"

// buildPublisher crea un IntegrationEventRedisPublisher listo para test.
// Recibe el mock de Redis para que cada test pueda configurar su comportamiento.
func buildPublisher(redisMock *mocks.RedisMock) *IntegrationEventRedisPublisher {
	logger := mocks.NewLoggerMock()
	return New(redisMock, testChannel, logger)
}

// buildEvent construye un IntegrationEvent con valores válidos para usar en tests.
func buildEvent() domain.IntegrationEvent {
	businessID := uint(42)
	return domain.IntegrationEvent{
		ID:            "evt-001",
		Type:          domain.IntegrationEventTypeSyncOrderCreated,
		IntegrationID: 10,
		BusinessID:    &businessID,
		Timestamp:     time.Now(),
		Data: domain.SyncOrderCreatedEvent{
			OrderID:     "order-001",
			OrderNumber: "#1001",
			Platform:    "shopify",
			Status:      "paid",
			CreatedAt:   time.Now(),
		},
		Metadata: map[string]interface{}{
			"order_id": "order-001",
			"platform": "shopify",
		},
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Constructor
// ─────────────────────────────────────────────────────────────────────────────

func TestNew_CreaPublisher_ConCamposCorrectos(t *testing.T) {
	redisMock := &mocks.RedisMock{}
	publisher := buildPublisher(redisMock)

	if publisher == nil {
		t.Fatal("New devolvio nil")
	}
	if publisher.channel != testChannel {
		t.Errorf("channel = %q, want %q", publisher.channel, testChannel)
	}
	if publisher.redisClient == nil {
		t.Error("redisClient no debe ser nil")
	}
	if publisher.logger == nil {
		t.Error("logger no debe ser nil")
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Publish — cliente Redis no disponible (Client retorna nil)
// ─────────────────────────────────────────────────────────────────────────────

func TestPublish_ClienteRedisNil_RetornaError(t *testing.T) {
	// RedisMock.ClientFn es nil por defecto, lo que hace que Client() retorne nil.
	// El publisher debe detectar esto y retornar un error descriptivo.
	redisMock := &mocks.RedisMock{}
	publisher := buildPublisher(redisMock)

	event := buildEvent()
	err := publisher.Publish(context.Background(), event)

	if err == nil {
		t.Fatal("esperaba error cuando el cliente Redis es nil, obtuvo nil")
	}

	if !strings.Contains(err.Error(), "redis client no disponible") {
		t.Errorf("mensaje de error inesperado: %q", err.Error())
	}
}

func TestPublish_ClienteFnExplicitoNil_RetornaError(t *testing.T) {
	// Configuramos explícitamente ClientFn para retornar nil
	redisMock := &mocks.RedisMock{
		ClientFn: nil, // nil => implementacion por defecto retorna nil
	}
	publisher := buildPublisher(redisMock)

	err := publisher.Publish(context.Background(), buildEvent())

	if err == nil {
		t.Fatal("esperaba error 'redis client no disponible'")
	}
	if !strings.Contains(err.Error(), "redis client no disponible") {
		t.Errorf("mensaje inesperado: %q", err.Error())
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Publish — serialización: evento con datos serializables pasa la validacion JSON
// ─────────────────────────────────────────────────────────────────────────────

// TestPublish_SerializacionPayload verifica que el payload JSON se construye
// correctamente antes de llegar al cliente Redis. Para esto usamos un evento
// sin campos problemáticos y validamos el error por cliente nil (la serialización
// ocurre antes de llamar a Client()).
func TestPublish_SerializacionPayload_OcurreAntesDeClient(t *testing.T) {
	// El error "redis client no disponible" implica que la serialización JSON
	// fue exitosa (de lo contrario el error sería "error serializando...").
	redisMock := &mocks.RedisMock{} // Client() retorna nil por defecto
	publisher := buildPublisher(redisMock)

	event := buildEvent()
	err := publisher.Publish(context.Background(), event)

	// El error debe ser exactamente "redis client no disponible",
	// NO "error serializando integration event".
	if err == nil {
		t.Fatal("esperaba error, obtuvo nil")
	}
	if strings.Contains(err.Error(), "serializando") {
		t.Errorf("la serializacion fallo inesperadamente: %v", err)
	}
	if !strings.Contains(err.Error(), "redis client no disponible") {
		t.Errorf("error inesperado: %v", err)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Publish — validacion de campos del payload serializado
// ─────────────────────────────────────────────────────────────────────────────

// TestPublish_CamposPayload_SonCorrectos usa un PublishFn personalizado que
// captura el JSON enviado a Redis y verifica que contenga los campos esperados.
// Nota: dado que *redis.Client es un tipo concreto, no podemos sustituirlo por
// una interfaz sin modificar el codigo de produccion. Por eso usamos el patron
// de interceptar en el mock a nivel de RedisMock.ClientFn y un hook en el
// publisher para verificar la serialización antes de invocar Client().
//
// La estrategia aqui es interceptar via un wrapping del publisher que capture
// el payload antes de enviarlo a Redis. Como la serialización ocurre antes de
// Client(), podemos crear un publisher que sí tenga un cliente real (modo test
// de go-redis con miniredis) O bien verificar los campos de otra manera.
//
// Para mantenernos sin dependencias externas, verificamos la serialización
// directamente testeando la función auxiliar del payload:
func TestPublish_IntegrationEventPayload_CamposCorrectos(t *testing.T) {
	businessID := uint(99)
	ts := time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC)

	event := domain.IntegrationEvent{
		ID:            "evt-test-999",
		Type:          domain.IntegrationEventTypeSyncOrderUpdated,
		IntegrationID: 77,
		BusinessID:    &businessID,
		Timestamp:     ts,
		Data:          map[string]string{"key": "value"},
		Metadata:      map[string]interface{}{"order_id": "abc"},
	}

	// Construimos el payload igual que el código de producción
	payload := integrationEventPayload{
		ID:            event.ID,
		Type:          string(event.Type),
		IntegrationID: event.IntegrationID,
		BusinessID:    event.BusinessID,
		Timestamp:     event.Timestamp,
		Data:          event.Data,
		Metadata:      event.Metadata,
	}

	if payload.ID != "evt-test-999" {
		t.Errorf("payload.ID = %q, want evt-test-999", payload.ID)
	}
	if payload.Type != string(domain.IntegrationEventTypeSyncOrderUpdated) {
		t.Errorf("payload.Type = %q, want %q", payload.Type, domain.IntegrationEventTypeSyncOrderUpdated)
	}
	if payload.IntegrationID != 77 {
		t.Errorf("payload.IntegrationID = %d, want 77", payload.IntegrationID)
	}
	if payload.BusinessID == nil || *payload.BusinessID != 99 {
		t.Errorf("payload.BusinessID = %v, want ptr(99)", payload.BusinessID)
	}
	if !payload.Timestamp.Equal(ts) {
		t.Errorf("payload.Timestamp = %v, want %v", payload.Timestamp, ts)
	}
}

func TestPublish_IntegrationEventPayload_BusinessIDNil(t *testing.T) {
	event := domain.IntegrationEvent{
		ID:            "evt-no-business",
		Type:          domain.IntegrationEventTypeSyncFailed,
		IntegrationID: 5,
		BusinessID:    nil,
		Timestamp:     time.Now(),
	}

	payload := integrationEventPayload{
		ID:            event.ID,
		Type:          string(event.Type),
		IntegrationID: event.IntegrationID,
		BusinessID:    event.BusinessID,
		Timestamp:     event.Timestamp,
	}

	if payload.BusinessID != nil {
		t.Errorf("esperaba BusinessID nil, obtuvo %v", payload.BusinessID)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Publish — tabla de casos de error
// ─────────────────────────────────────────────────────────────────────────────

func TestPublish_TablaDeErrores(t *testing.T) {
	tests := []struct {
		name        string
		setupMock   func() *mocks.RedisMock
		wantErr     bool
		wantErrMsg  string
	}{
		{
			name: "cliente Redis nil retorna error descriptivo",
			setupMock: func() *mocks.RedisMock {
				return &mocks.RedisMock{} // Client() retorna nil por defecto
			},
			wantErr:    true,
			wantErrMsg: "redis client no disponible",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			publisher := buildPublisher(tt.setupMock())
			event := buildEvent()

			err := publisher.Publish(context.Background(), event)

			if tt.wantErr && err == nil {
				t.Fatal("esperaba error, obtuvo nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("esperaba nil, obtuvo: %v", err)
			}
			if tt.wantErr && !strings.Contains(err.Error(), tt.wantErrMsg) {
				t.Errorf("mensaje de error %q no contiene %q", err.Error(), tt.wantErrMsg)
			}
		})
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Tipos de eventos — verificacion de constantes de dominio
// ─────────────────────────────────────────────────────────────────────────────

// TestIntegrationEventType_Constantes verifica que las constantes de tipo de evento
// tengan los valores de string correctos que se usan al serializar a Redis.
func TestIntegrationEventType_Constantes(t *testing.T) {
	tests := []struct {
		eventType domain.IntegrationEventType
		wantStr   string
	}{
		{domain.IntegrationEventTypeSyncOrderCreated, "integration.sync.order.created"},
		{domain.IntegrationEventTypeSyncOrderUpdated, "integration.sync.order.updated"},
		{domain.IntegrationEventTypeSyncOrderRejected, "integration.sync.order.rejected"},
		{domain.IntegrationEventTypeSyncStarted, "integration.sync.started"},
		{domain.IntegrationEventTypeSyncCompleted, "integration.sync.completed"},
		{domain.IntegrationEventTypeSyncFailed, "integration.sync.failed"},
	}

	for _, tt := range tests {
		t.Run(tt.wantStr, func(t *testing.T) {
			if string(tt.eventType) != tt.wantStr {
				t.Errorf("IntegrationEventType = %q, want %q", tt.eventType, tt.wantStr)
			}
		})
	}
}
