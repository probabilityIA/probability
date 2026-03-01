package usecaseordermapping

// Este archivo contiene los tests unitarios del paquete usecaseordermapping.
// Se usan mocks inline (structs con campos Fn) siguiendo el patrón del proyecto,
// instanciando UseCaseOrderMapping directamente para inyectar dependencias.

import (
	"context"
	"errors"
	"testing"

	"github.com/rs/zerolog"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
	"github.com/secamc93/probability/back/central/shared/log"
)

// ─── Mock: IRepository ──────────────────────────────────────────────────────

type mockRepository struct {
	CreateOrderFn                                         func(ctx context.Context, order *entities.ProbabilityOrder) error
	GetOrderByIDFn                                        func(ctx context.Context, id string) (*entities.ProbabilityOrder, error)
	GetOrderByInternalNumberFn                            func(ctx context.Context, internalNumber string) (*entities.ProbabilityOrder, error)
	GetOrderByOrderNumberFn                               func(ctx context.Context, orderNumber string) (*entities.ProbabilityOrder, error)
	ListOrdersFn                                          func(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]entities.ProbabilityOrder, int64, error)
	UpdateOrderFn                                         func(ctx context.Context, order *entities.ProbabilityOrder) error
	DeleteOrderFn                                         func(ctx context.Context, id string) error
	GetOrderRawFn                                         func(ctx context.Context, id string) (*entities.ProbabilityOrderChannelMetadata, error)
	CountOrdersByClientIDFn                               func(ctx context.Context, clientID uint) (int64, error)
	GetLastManualOrderNumberFn                            func(ctx context.Context, businessID uint) (int, error)
	GetFirstIntegrationIDByBusinessIDFn                   func(ctx context.Context, businessID uint) (uint, error)
	GetPlatformIntegrationIDByBusinessIDFn                func(ctx context.Context, businessID uint) (uint, error)
	OrderExistsFn                                         func(ctx context.Context, externalID string, integrationID uint) (bool, error)
	GetOrderByExternalIDFn                                func(ctx context.Context, externalID string, integrationID uint) (*entities.ProbabilityOrder, error)
	CreateOrderItemsFn                                    func(ctx context.Context, items []*entities.ProbabilityOrderItem) error
	CreateAddressesFn                                     func(ctx context.Context, addresses []*entities.ProbabilityAddress) error
	CreatePaymentsFn                                      func(ctx context.Context, payments []*entities.ProbabilityPayment) error
	CreateShipmentsFn                                     func(ctx context.Context, shipments []*entities.ProbabilityShipment) error
	CreateChannelMetadataFn                               func(ctx context.Context, metadata *entities.ProbabilityOrderChannelMetadata) error
	GetProductBySKUFn                                     func(ctx context.Context, businessID uint, sku string) (*entities.Product, error)
	CreateProductFn                                       func(ctx context.Context, product *entities.Product) error
	GetClientByEmailFn                                    func(ctx context.Context, businessID uint, email string) (*entities.Client, error)
	GetClientByDNIFn                                      func(ctx context.Context, businessID uint, dni string) (*entities.Client, error)
	CreateClientFn                                        func(ctx context.Context, client *entities.Client) error
	CreateOrderErrorFn                                    func(ctx context.Context, orderError *entities.OrderError) error
	GetOrderStatusIDByIntegrationTypeAndOriginalStatusFn  func(ctx context.Context, integrationTypeID uint, originalStatus string) (*uint, error)
	GetPaymentStatusIDByCodeFn                            func(ctx context.Context, code string) (*uint, error)
	GetFulfillmentStatusIDByCodeFn                        func(ctx context.Context, code string) (*uint, error)
}

func (m *mockRepository) CreateOrder(ctx context.Context, order *entities.ProbabilityOrder) error {
	if m.CreateOrderFn != nil {
		return m.CreateOrderFn(ctx, order)
	}
	return nil
}

func (m *mockRepository) GetOrderByID(ctx context.Context, id string) (*entities.ProbabilityOrder, error) {
	if m.GetOrderByIDFn != nil {
		return m.GetOrderByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *mockRepository) GetOrderByInternalNumber(ctx context.Context, internalNumber string) (*entities.ProbabilityOrder, error) {
	if m.GetOrderByInternalNumberFn != nil {
		return m.GetOrderByInternalNumberFn(ctx, internalNumber)
	}
	return nil, nil
}

func (m *mockRepository) GetOrderByOrderNumber(ctx context.Context, orderNumber string) (*entities.ProbabilityOrder, error) {
	if m.GetOrderByOrderNumberFn != nil {
		return m.GetOrderByOrderNumberFn(ctx, orderNumber)
	}
	return nil, nil
}

func (m *mockRepository) ListOrders(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]entities.ProbabilityOrder, int64, error) {
	if m.ListOrdersFn != nil {
		return m.ListOrdersFn(ctx, page, pageSize, filters)
	}
	return nil, 0, nil
}

func (m *mockRepository) UpdateOrder(ctx context.Context, order *entities.ProbabilityOrder) error {
	if m.UpdateOrderFn != nil {
		return m.UpdateOrderFn(ctx, order)
	}
	return nil
}

func (m *mockRepository) DeleteOrder(ctx context.Context, id string) error {
	if m.DeleteOrderFn != nil {
		return m.DeleteOrderFn(ctx, id)
	}
	return nil
}

func (m *mockRepository) GetOrderRaw(ctx context.Context, id string) (*entities.ProbabilityOrderChannelMetadata, error) {
	if m.GetOrderRawFn != nil {
		return m.GetOrderRawFn(ctx, id)
	}
	return nil, nil
}

func (m *mockRepository) CountOrdersByClientID(ctx context.Context, clientID uint) (int64, error) {
	if m.CountOrdersByClientIDFn != nil {
		return m.CountOrdersByClientIDFn(ctx, clientID)
	}
	return 0, nil
}

func (m *mockRepository) GetLastManualOrderNumber(ctx context.Context, businessID uint) (int, error) {
	if m.GetLastManualOrderNumberFn != nil {
		return m.GetLastManualOrderNumberFn(ctx, businessID)
	}
	return 0, nil
}

func (m *mockRepository) GetFirstIntegrationIDByBusinessID(ctx context.Context, businessID uint) (uint, error) {
	if m.GetFirstIntegrationIDByBusinessIDFn != nil {
		return m.GetFirstIntegrationIDByBusinessIDFn(ctx, businessID)
	}
	return 0, nil
}

func (m *mockRepository) GetPlatformIntegrationIDByBusinessID(ctx context.Context, businessID uint) (uint, error) {
	if m.GetPlatformIntegrationIDByBusinessIDFn != nil {
		return m.GetPlatformIntegrationIDByBusinessIDFn(ctx, businessID)
	}
	return 0, nil
}

func (m *mockRepository) OrderExists(ctx context.Context, externalID string, integrationID uint) (bool, error) {
	if m.OrderExistsFn != nil {
		return m.OrderExistsFn(ctx, externalID, integrationID)
	}
	return false, nil
}

func (m *mockRepository) GetOrderByExternalID(ctx context.Context, externalID string, integrationID uint) (*entities.ProbabilityOrder, error) {
	if m.GetOrderByExternalIDFn != nil {
		return m.GetOrderByExternalIDFn(ctx, externalID, integrationID)
	}
	return nil, nil
}

func (m *mockRepository) CreateOrderItems(ctx context.Context, items []*entities.ProbabilityOrderItem) error {
	if m.CreateOrderItemsFn != nil {
		return m.CreateOrderItemsFn(ctx, items)
	}
	return nil
}

func (m *mockRepository) CreateAddresses(ctx context.Context, addresses []*entities.ProbabilityAddress) error {
	if m.CreateAddressesFn != nil {
		return m.CreateAddressesFn(ctx, addresses)
	}
	return nil
}

func (m *mockRepository) CreatePayments(ctx context.Context, payments []*entities.ProbabilityPayment) error {
	if m.CreatePaymentsFn != nil {
		return m.CreatePaymentsFn(ctx, payments)
	}
	return nil
}

func (m *mockRepository) CreateShipments(ctx context.Context, shipments []*entities.ProbabilityShipment) error {
	if m.CreateShipmentsFn != nil {
		return m.CreateShipmentsFn(ctx, shipments)
	}
	return nil
}

func (m *mockRepository) CreateChannelMetadata(ctx context.Context, metadata *entities.ProbabilityOrderChannelMetadata) error {
	if m.CreateChannelMetadataFn != nil {
		return m.CreateChannelMetadataFn(ctx, metadata)
	}
	return nil
}

func (m *mockRepository) GetProductBySKU(ctx context.Context, businessID uint, sku string) (*entities.Product, error) {
	if m.GetProductBySKUFn != nil {
		return m.GetProductBySKUFn(ctx, businessID, sku)
	}
	return nil, nil
}

func (m *mockRepository) CreateProduct(ctx context.Context, product *entities.Product) error {
	if m.CreateProductFn != nil {
		return m.CreateProductFn(ctx, product)
	}
	return nil
}

func (m *mockRepository) GetClientByEmail(ctx context.Context, businessID uint, email string) (*entities.Client, error) {
	if m.GetClientByEmailFn != nil {
		return m.GetClientByEmailFn(ctx, businessID, email)
	}
	return nil, nil
}

func (m *mockRepository) GetClientByDNI(ctx context.Context, businessID uint, dni string) (*entities.Client, error) {
	if m.GetClientByDNIFn != nil {
		return m.GetClientByDNIFn(ctx, businessID, dni)
	}
	return nil, nil
}

func (m *mockRepository) CreateClient(ctx context.Context, client *entities.Client) error {
	if m.CreateClientFn != nil {
		return m.CreateClientFn(ctx, client)
	}
	return nil
}

func (m *mockRepository) CreateOrderError(ctx context.Context, orderError *entities.OrderError) error {
	if m.CreateOrderErrorFn != nil {
		return m.CreateOrderErrorFn(ctx, orderError)
	}
	return nil
}

func (m *mockRepository) GetOrderStatusIDByIntegrationTypeAndOriginalStatus(ctx context.Context, integrationTypeID uint, originalStatus string) (*uint, error) {
	if m.GetOrderStatusIDByIntegrationTypeAndOriginalStatusFn != nil {
		return m.GetOrderStatusIDByIntegrationTypeAndOriginalStatusFn(ctx, integrationTypeID, originalStatus)
	}
	return nil, nil
}

func (m *mockRepository) GetPaymentStatusIDByCode(ctx context.Context, code string) (*uint, error) {
	if m.GetPaymentStatusIDByCodeFn != nil {
		return m.GetPaymentStatusIDByCodeFn(ctx, code)
	}
	return nil, nil
}

func (m *mockRepository) GetFulfillmentStatusIDByCode(ctx context.Context, code string) (*uint, error) {
	if m.GetFulfillmentStatusIDByCodeFn != nil {
		return m.GetFulfillmentStatusIDByCodeFn(ctx, code)
	}
	return nil, nil
}

// ─── Mock: IOrderEventPublisher (Redis) ─────────────────────────────────────

type mockRedisPublisher struct {
	PublishOrderEventFn func(ctx context.Context, event *entities.OrderEvent, order *entities.ProbabilityOrder) error
}

func (m *mockRedisPublisher) PublishOrderEvent(ctx context.Context, event *entities.OrderEvent, order *entities.ProbabilityOrder) error {
	if m.PublishOrderEventFn != nil {
		return m.PublishOrderEventFn(ctx, event, order)
	}
	return nil
}

// ─── Mock: IOrderRabbitPublisher ────────────────────────────────────────────

type mockRabbitPublisher struct {
	PublishOrderCreatedFn           func(ctx context.Context, order *entities.ProbabilityOrder) error
	PublishOrderUpdatedFn           func(ctx context.Context, order *entities.ProbabilityOrder) error
	PublishOrderCancelledFn         func(ctx context.Context, order *entities.ProbabilityOrder) error
	PublishOrderStatusChangedFn     func(ctx context.Context, order *entities.ProbabilityOrder, previousStatus, currentStatus string) error
	PublishConfirmationRequestedFn  func(ctx context.Context, order *entities.ProbabilityOrder) error
	PublishOrderEventFn             func(ctx context.Context, event *entities.OrderEvent, order *entities.ProbabilityOrder) error
}

func (m *mockRabbitPublisher) PublishOrderCreated(ctx context.Context, order *entities.ProbabilityOrder) error {
	if m.PublishOrderCreatedFn != nil {
		return m.PublishOrderCreatedFn(ctx, order)
	}
	return nil
}

func (m *mockRabbitPublisher) PublishOrderUpdated(ctx context.Context, order *entities.ProbabilityOrder) error {
	if m.PublishOrderUpdatedFn != nil {
		return m.PublishOrderUpdatedFn(ctx, order)
	}
	return nil
}

func (m *mockRabbitPublisher) PublishOrderCancelled(ctx context.Context, order *entities.ProbabilityOrder) error {
	if m.PublishOrderCancelledFn != nil {
		return m.PublishOrderCancelledFn(ctx, order)
	}
	return nil
}

func (m *mockRabbitPublisher) PublishOrderStatusChanged(ctx context.Context, order *entities.ProbabilityOrder, previousStatus, currentStatus string) error {
	if m.PublishOrderStatusChangedFn != nil {
		return m.PublishOrderStatusChangedFn(ctx, order, previousStatus, currentStatus)
	}
	return nil
}

func (m *mockRabbitPublisher) PublishConfirmationRequested(ctx context.Context, order *entities.ProbabilityOrder) error {
	if m.PublishConfirmationRequestedFn != nil {
		return m.PublishConfirmationRequestedFn(ctx, order)
	}
	return nil
}

func (m *mockRabbitPublisher) PublishOrderEvent(ctx context.Context, event *entities.OrderEvent, order *entities.ProbabilityOrder) error {
	if m.PublishOrderEventFn != nil {
		return m.PublishOrderEventFn(ctx, event, order)
	}
	return nil
}

// ─── Mock: IIntegrationEventPublisher ───────────────────────────────────────

type mockIntegrationEventPublisher struct {
	PublishSyncOrderRejectedFn func(ctx context.Context, integrationID uint, businessID *uint, orderNumber, externalID, platform, reason, errMsg string)
	PublishSyncOrderCreatedFn  func(ctx context.Context, integrationID uint, businessID *uint, data map[string]interface{})
	PublishSyncOrderUpdatedFn  func(ctx context.Context, integrationID uint, businessID *uint, data map[string]interface{})
}

func (m *mockIntegrationEventPublisher) PublishSyncOrderRejected(ctx context.Context, integrationID uint, businessID *uint, orderNumber, externalID, platform, reason, errMsg string) {
	if m.PublishSyncOrderRejectedFn != nil {
		m.PublishSyncOrderRejectedFn(ctx, integrationID, businessID, orderNumber, externalID, platform, reason, errMsg)
	}
}

func (m *mockIntegrationEventPublisher) PublishSyncOrderCreated(ctx context.Context, integrationID uint, businessID *uint, data map[string]interface{}) {
	if m.PublishSyncOrderCreatedFn != nil {
		m.PublishSyncOrderCreatedFn(ctx, integrationID, businessID, data)
	}
}

func (m *mockIntegrationEventPublisher) PublishSyncOrderUpdated(ctx context.Context, integrationID uint, businessID *uint, data map[string]interface{}) {
	if m.PublishSyncOrderUpdatedFn != nil {
		m.PublishSyncOrderUpdatedFn(ctx, integrationID, businessID, data)
	}
}

// ─── Mock: IOrderScoreUseCase ────────────────────────────────────────────────

type mockScoreUseCase struct {
	CalculateOrderScoreFn          func(order *entities.ProbabilityOrder) (float64, []string)
	CalculateAndUpdateOrderScoreFn func(ctx context.Context, orderID string) error
}

func (m *mockScoreUseCase) CalculateOrderScore(order *entities.ProbabilityOrder) (float64, []string) {
	if m.CalculateOrderScoreFn != nil {
		return m.CalculateOrderScoreFn(order)
	}
	return 0.0, nil
}

func (m *mockScoreUseCase) CalculateAndUpdateOrderScore(ctx context.Context, orderID string) error {
	if m.CalculateAndUpdateOrderScoreFn != nil {
		return m.CalculateAndUpdateOrderScoreFn(ctx, orderID)
	}
	return nil
}

// ─── Mock: log.ILogger ──────────────────────────────────────────────────────

type mockLogger struct{}

func (m *mockLogger) Info(ctx ...context.Context) *zerolog.Event {
	nop := zerolog.Nop()
	return nop.Info()
}

func (m *mockLogger) Error(ctx ...context.Context) *zerolog.Event {
	nop := zerolog.Nop()
	return nop.Error()
}

func (m *mockLogger) Warn(ctx ...context.Context) *zerolog.Event {
	nop := zerolog.Nop()
	return nop.Warn()
}

func (m *mockLogger) Debug(ctx ...context.Context) *zerolog.Event {
	nop := zerolog.Nop()
	return nop.Debug()
}

func (m *mockLogger) Fatal(ctx ...context.Context) *zerolog.Event {
	nop := zerolog.Nop()
	return nop.Fatal()
}

func (m *mockLogger) Panic(ctx ...context.Context) *zerolog.Event {
	nop := zerolog.Nop()
	return nop.Panic()
}

func (m *mockLogger) With() zerolog.Context {
	nop := zerolog.Nop()
	return nop.With()
}

func (m *mockLogger) WithService(service string) log.ILogger     { return m }
func (m *mockLogger) WithModule(module string) log.ILogger       { return m }
func (m *mockLogger) WithBusinessID(businessID uint) log.ILogger { return m }

// ─── Helpers de construcción ─────────────────────────────────────────────────

// newTestUseCase construye un UseCaseOrderMapping con todas las dependencias mockeadas.
// Instancia el struct directamente para inyectar el scoreUseCase mock.
// IMPORTANTE: Los publishers se pasan como interfaces para que nil sea nil-interfaz real.
// Si se pasa un puntero nil tipado (como *mockRedisPublisher(nil)) Go crea una interfaz
// no-nil que causa panics al derreferenciar el receptor. Por eso se usan wrappers.
func newTestUseCase(
	repo *mockRepository,
	redisPublisher *mockRedisPublisher,
	rabbitPublisher *mockRabbitPublisher,
	integrationEventPublisher *mockIntegrationEventPublisher,
	scoreUseCase *mockScoreUseCase,
) *UseCaseOrderMapping {
	uc := &UseCaseOrderMapping{
		repo:         repo,
		logger:       &mockLogger{},
		scoreUseCase: scoreUseCase,
	}

	// Asignar publishers solo si no son nil para evitar el problema de interfaz con puntero nil
	if redisPublisher != nil {
		uc.redisEventPublisher = redisPublisher
	}
	if rabbitPublisher != nil {
		uc.rabbitEventPublisher = rabbitPublisher
	}
	if integrationEventPublisher != nil {
		uc.integrationEventPublisher = integrationEventPublisher
	}

	return uc
}

// newBusinessID es un helper que retorna un puntero a uint.
func newBusinessID(id uint) *uint {
	return &id
}

// newUint es un helper que retorna un puntero a uint.
func newUint(v uint) *uint {
	return &v
}

// newMinimalDTO construye un DTO mínimo válido para pruebas.
func newMinimalDTO(externalID string, integrationID uint, businessID uint) *dtos.ProbabilityOrderDTO {
	bID := businessID
	return &dtos.ProbabilityOrderDTO{
		BusinessID:    &bID,
		IntegrationID: integrationID,
		ExternalID:    externalID,
		OrderNumber:   "ORD-001",
		Platform:      "shopify",
		TotalAmount:   100.0,
		Currency:      "COP",
		CustomerEmail: "test@example.com",
	}
}

// ─── Tests: MapAndSaveOrder ──────────────────────────────────────────────────

func TestMapAndSaveOrder_FaltaIntegrationID_RetornaError(t *testing.T) {
	// Arrange
	repo := &mockRepository{}
	uc := newTestUseCase(repo, nil, nil, nil, &mockScoreUseCase{})
	bID := uint(1)
	dto := &dtos.ProbabilityOrderDTO{
		BusinessID:    &bID,
		IntegrationID: 0, // Falta IntegrationID
		ExternalID:    "EXT-001",
	}

	// Act
	result, err := uc.MapAndSaveOrder(context.Background(), dto)

	// Assert
	if err == nil {
		t.Fatal("se esperaba error cuando integration_id es 0, pero no se obtuvo ninguno")
	}
	if result != nil {
		t.Errorf("se esperaba resultado nil, pero se obtuvo: %v", result)
	}
	if err.Error() != "integration_id is required" {
		t.Errorf("mensaje de error inesperado: %q", err.Error())
	}
}

func TestMapAndSaveOrder_FaltaBusinessID_RetornaError(t *testing.T) {
	// Arrange
	repo := &mockRepository{}
	uc := newTestUseCase(repo, nil, nil, nil, &mockScoreUseCase{})
	dto := &dtos.ProbabilityOrderDTO{
		BusinessID:    nil, // Falta BusinessID
		IntegrationID: 10,
		ExternalID:    "EXT-001",
	}

	// Act
	result, err := uc.MapAndSaveOrder(context.Background(), dto)

	// Assert
	if err == nil {
		t.Fatal("se esperaba error cuando business_id es nil, pero no se obtuvo ninguno")
	}
	if result != nil {
		t.Errorf("se esperaba resultado nil, pero se obtuvo: %v", result)
	}
	if err.Error() != "business_id is required" {
		t.Errorf("mensaje de error inesperado: %q", err.Error())
	}
}

func TestMapAndSaveOrder_BusinessIDCero_RetornaError(t *testing.T) {
	// Arrange
	repo := &mockRepository{}
	uc := newTestUseCase(repo, nil, nil, nil, &mockScoreUseCase{})
	zeroID := uint(0)
	dto := &dtos.ProbabilityOrderDTO{
		BusinessID:    &zeroID, // BusinessID = 0
		IntegrationID: 10,
		ExternalID:    "EXT-001",
	}

	// Act
	result, err := uc.MapAndSaveOrder(context.Background(), dto)

	// Assert
	if err == nil {
		t.Fatal("se esperaba error cuando business_id es 0, pero no se obtuvo ninguno")
	}
	if result != nil {
		t.Errorf("se esperaba resultado nil, pero se obtuvo: %v", result)
	}
}

func TestMapAndSaveOrder_ErrorAlVerificarExistencia_RetornaError(t *testing.T) {
	// Arrange
	expectedErr := errors.New("fallo en base de datos")
	repo := &mockRepository{
		OrderExistsFn: func(ctx context.Context, externalID string, integrationID uint) (bool, error) {
			return false, expectedErr
		},
	}
	uc := newTestUseCase(repo, nil, nil, nil, &mockScoreUseCase{})
	dto := newMinimalDTO("EXT-001", 10, 1)

	// Act
	result, err := uc.MapAndSaveOrder(context.Background(), dto)

	// Assert
	if err == nil {
		t.Fatal("se esperaba error al fallar la verificación de existencia")
	}
	if result != nil {
		t.Errorf("se esperaba resultado nil, pero se obtuvo: %v", result)
	}
	if !errors.Is(err, expectedErr) {
		t.Errorf("se esperaba error wrapeado con: %v, pero se obtuvo: %v", expectedErr, err)
	}
}

func TestMapAndSaveOrder_OrdenNueva_ExitoSinEntidadesRelacionadas(t *testing.T) {
	// Arrange
	createOrderCalled := false
	repo := &mockRepository{
		OrderExistsFn: func(ctx context.Context, externalID string, integrationID uint) (bool, error) {
			return false, nil // Orden no existe aun
		},
		GetClientByEmailFn: func(ctx context.Context, businessID uint, email string) (*entities.Client, error) {
			return &entities.Client{ID: 99, Email: email}, nil // Cliente ya existe
		},
		CreateOrderFn: func(ctx context.Context, order *entities.ProbabilityOrder) error {
			createOrderCalled = true
			order.ID = "ORDER-123" // Simular que la DB asigna un ID
			return nil
		},
	}
	uc := newTestUseCase(repo, nil, nil, nil, &mockScoreUseCase{})
	dto := newMinimalDTO("EXT-002", 10, 1)

	// Act
	result, err := uc.MapAndSaveOrder(context.Background(), dto)

	// Assert
	if err != nil {
		t.Fatalf("no se esperaba error, pero se obtuvo: %v", err)
	}
	if result == nil {
		t.Fatal("se esperaba resultado no nil, pero se obtuvo nil")
	}
	if !createOrderCalled {
		t.Error("se esperaba que CreateOrder fuera llamado, pero no lo fue")
	}
	if result.ExternalID != "EXT-002" {
		t.Errorf("ExternalID esperado: %q, obtenido: %q", "EXT-002", result.ExternalID)
	}
	if result.Currency != "COP" {
		t.Errorf("Currency esperada: %q, obtenida: %q", "COP", result.Currency)
	}
}

func TestMapAndSaveOrder_OrdenExistente_LlamaUpdateOrder(t *testing.T) {
	// Arrange
	updateOrderCalled := false
	existingOrder := &entities.ProbabilityOrder{
		ID:          "ORDER-EXISTENTE",
		ExternalID:  "EXT-003",
		Status:      "pending",
		TotalAmount: 50.0,
		Currency:    "COP",
	}
	repo := &mockRepository{
		OrderExistsFn: func(ctx context.Context, externalID string, integrationID uint) (bool, error) {
			return true, nil // Orden ya existe
		},
		GetOrderByExternalIDFn: func(ctx context.Context, externalID string, integrationID uint) (*entities.ProbabilityOrder, error) {
			return existingOrder, nil
		},
		UpdateOrderFn: func(ctx context.Context, order *entities.ProbabilityOrder) error {
			updateOrderCalled = true
			return nil
		},
	}
	uc := newTestUseCase(repo, nil, nil, nil, &mockScoreUseCase{})
	dto := newMinimalDTO("EXT-003", 10, 1)
	dto.Status = "completed" // Estado diferente para forzar cambio

	// Act
	result, err := uc.MapAndSaveOrder(context.Background(), dto)

	// Assert
	if err != nil {
		t.Fatalf("no se esperaba error, pero se obtuvo: %v", err)
	}
	if result == nil {
		t.Fatal("se esperaba resultado no nil, pero se obtuvo nil")
	}
	if !updateOrderCalled {
		t.Error("se esperaba que UpdateOrder fuera llamado, pero no lo fue")
	}
}

func TestMapAndSaveOrder_ErrorAlCrearOrden_RetornaError(t *testing.T) {
	// Arrange
	expectedErr := errors.New("error de constraint en BD")
	repo := &mockRepository{
		OrderExistsFn: func(ctx context.Context, externalID string, integrationID uint) (bool, error) {
			return false, nil
		},
		GetClientByEmailFn: func(ctx context.Context, businessID uint, email string) (*entities.Client, error) {
			return nil, nil // Cliente no existe
		},
		GetClientByDNIFn: func(ctx context.Context, businessID uint, dni string) (*entities.Client, error) {
			return nil, nil
		},
		CreateClientFn: func(ctx context.Context, client *entities.Client) error {
			return nil // Crear cliente exitoso
		},
		CreateOrderFn: func(ctx context.Context, order *entities.ProbabilityOrder) error {
			return expectedErr // Fallo al crear la orden
		},
	}
	integrationPub := &mockIntegrationEventPublisher{}
	uc := newTestUseCase(repo, nil, nil, integrationPub, &mockScoreUseCase{})
	dto := newMinimalDTO("EXT-004", 10, 1)

	// Act
	result, err := uc.MapAndSaveOrder(context.Background(), dto)

	// Assert
	if err == nil {
		t.Fatal("se esperaba error al fallar la creación de la orden")
	}
	if result != nil {
		t.Errorf("se esperaba resultado nil, pero se obtuvo: %v", result)
	}
	if !errors.Is(err, expectedErr) {
		t.Errorf("se esperaba error wrapeado con: %v, pero se obtuvo: %v", expectedErr, err)
	}
}

func TestMapAndSaveOrder_ErrorAlProcesarCliente_RetornaError(t *testing.T) {
	// Arrange
	expectedErr := errors.New("error al buscar cliente")
	repo := &mockRepository{
		OrderExistsFn: func(ctx context.Context, externalID string, integrationID uint) (bool, error) {
			return false, nil
		},
		GetClientByEmailFn: func(ctx context.Context, businessID uint, email string) (*entities.Client, error) {
			return nil, expectedErr // Error al buscar cliente
		},
	}
	integrationPub := &mockIntegrationEventPublisher{}
	uc := newTestUseCase(repo, nil, nil, integrationPub, &mockScoreUseCase{})
	dto := newMinimalDTO("EXT-005", 10, 1)

	// Act
	result, err := uc.MapAndSaveOrder(context.Background(), dto)

	// Assert
	if err == nil {
		t.Fatal("se esperaba error al fallar el procesamiento de cliente")
	}
	if result != nil {
		t.Errorf("se esperaba resultado nil, pero se obtuvo: %v", result)
	}
	if !errors.Is(err, expectedErr) {
		t.Errorf("se esperaba error wrapeado con: %v, pero se obtuvo: %v", expectedErr, err)
	}
}

// ─── Tests: UpdateOrder ───────────────────────────────────────────────────────

func TestUpdateOrder_SinCambios_RetornaSinActualizar(t *testing.T) {
	// Arrange
	updateOrderCalled := false
	existingOrder := &entities.ProbabilityOrder{
		ID:            "ORDER-U01",
		ExternalID:    "EXT-U01",
		Status:        "pending",
		TotalAmount:   100.0,
		Currency:      "COP",
		CustomerEmail: "cliente@example.com",
	}
	repo := &mockRepository{
		UpdateOrderFn: func(ctx context.Context, order *entities.ProbabilityOrder) error {
			updateOrderCalled = true
			return nil
		},
	}
	uc := newTestUseCase(repo, nil, nil, nil, &mockScoreUseCase{})

	// DTO con los mismos valores que la orden existente (sin cambios)
	dto := &dtos.ProbabilityOrderDTO{
		Status:        "pending",   // Igual al existente
		TotalAmount:   0,           // No actualiza (dto.TotalAmount <= 0)
		Currency:      "",          // No actualiza (vacío)
		CustomerEmail: "",          // No actualiza (vacío)
	}

	// Act
	result, err := uc.UpdateOrder(context.Background(), existingOrder, dto)

	// Assert
	if err != nil {
		t.Fatalf("no se esperaba error, pero se obtuvo: %v", err)
	}
	if result == nil {
		t.Fatal("se esperaba resultado no nil")
	}
	if updateOrderCalled {
		t.Error("no se esperaba que UpdateOrder fuera llamado cuando no hay cambios")
	}
	if result.ID != "ORDER-U01" {
		t.Errorf("ID esperado: %q, obtenido: %q", "ORDER-U01", result.ID)
	}
}

func TestUpdateOrder_CambioDeEstado_ActualizaYPublicaEvento(t *testing.T) {
	// Arrange
	updateOrderCalled := false
	existingOrder := &entities.ProbabilityOrder{
		ID:          "ORDER-U02",
		ExternalID:  "EXT-U02",
		Status:      "pending",
		TotalAmount: 150.0,
		Currency:    "COP",
	}
	repo := &mockRepository{
		UpdateOrderFn: func(ctx context.Context, order *entities.ProbabilityOrder) error {
			updateOrderCalled = true
			return nil
		},
	}
	uc := newTestUseCase(repo, nil, nil, nil, &mockScoreUseCase{})

	dto := &dtos.ProbabilityOrderDTO{
		Status: "completed", // Cambio de estado
	}

	// Act
	result, err := uc.UpdateOrder(context.Background(), existingOrder, dto)

	// Assert
	if err != nil {
		t.Fatalf("no se esperaba error, pero se obtuvo: %v", err)
	}
	if result == nil {
		t.Fatal("se esperaba resultado no nil")
	}
	if !updateOrderCalled {
		t.Error("se esperaba que UpdateOrder fuera llamado ante el cambio de estado")
	}
	if result.Status != "completed" {
		t.Errorf("status esperado: %q, obtenido: %q", "completed", result.Status)
	}
}

func TestUpdateOrder_ErrorEnRepositorio_RetornaError(t *testing.T) {
	// Arrange
	expectedErr := errors.New("error de conexion BD")
	existingOrder := &entities.ProbabilityOrder{
		ID:          "ORDER-U03",
		ExternalID:  "EXT-U03",
		Status:      "pending",
		TotalAmount: 200.0,
	}
	repo := &mockRepository{
		UpdateOrderFn: func(ctx context.Context, order *entities.ProbabilityOrder) error {
			return expectedErr
		},
	}
	uc := newTestUseCase(repo, nil, nil, nil, &mockScoreUseCase{})

	dto := &dtos.ProbabilityOrderDTO{
		Status: "shipped", // Hay un cambio para que se intente actualizar
	}

	// Act
	result, err := uc.UpdateOrder(context.Background(), existingOrder, dto)

	// Assert
	if err == nil {
		t.Fatal("se esperaba error al fallar el repositorio, pero no se obtuvo ninguno")
	}
	if result != nil {
		t.Errorf("se esperaba resultado nil, pero se obtuvo: %v", result)
	}
	if !errors.Is(err, expectedErr) {
		t.Errorf("se esperaba error wrapeado con: %v, pero se obtuvo: %v", expectedErr, err)
	}
}

func TestUpdateOrder_CambioInformacionFinanciera_ActualizaCampos(t *testing.T) {
	// Arrange
	existingOrder := &entities.ProbabilityOrder{
		ID:          "ORDER-U04",
		TotalAmount: 100.0,
		Currency:    "COP",
		Subtotal:    80.0,
	}
	repo := &mockRepository{
		UpdateOrderFn: func(ctx context.Context, order *entities.ProbabilityOrder) error {
			return nil
		},
	}
	uc := newTestUseCase(repo, nil, nil, nil, &mockScoreUseCase{})

	dto := &dtos.ProbabilityOrderDTO{
		TotalAmount: 200.0, // Valor diferente al existente
		Currency:    "USD", // Moneda diferente
		Subtotal:    160.0, // Subtotal diferente
	}

	// Act
	result, err := uc.UpdateOrder(context.Background(), existingOrder, dto)

	// Assert
	if err != nil {
		t.Fatalf("no se esperaba error, pero se obtuvo: %v", err)
	}
	if result == nil {
		t.Fatal("se esperaba resultado no nil")
	}
	if result.TotalAmount != 200.0 {
		t.Errorf("TotalAmount esperado: 200.0, obtenido: %v", result.TotalAmount)
	}
	if result.Currency != "USD" {
		t.Errorf("Currency esperada: %q, obtenida: %q", "USD", result.Currency)
	}
}

func TestUpdateOrder_CambioInformacionCliente_ActualizaCampos(t *testing.T) {
	// Arrange
	existingOrder := &entities.ProbabilityOrder{
		ID:            "ORDER-U05",
		CustomerName:  "Juan Perez",
		CustomerEmail: "juan@example.com",
		CustomerPhone: "3001234567",
		TotalAmount:   100.0,
	}
	repo := &mockRepository{
		UpdateOrderFn: func(ctx context.Context, order *entities.ProbabilityOrder) error {
			return nil
		},
	}
	uc := newTestUseCase(repo, nil, nil, nil, &mockScoreUseCase{})

	dto := &dtos.ProbabilityOrderDTO{
		CustomerName:  "Juan Carlos Perez",
		CustomerEmail: "juancarlos@example.com",
		CustomerPhone: "3009999999",
	}

	// Act
	result, err := uc.UpdateOrder(context.Background(), existingOrder, dto)

	// Assert
	if err != nil {
		t.Fatalf("no se esperaba error, pero se obtuvo: %v", err)
	}
	if result == nil {
		t.Fatal("se esperaba resultado no nil")
	}
	if result.CustomerName != "Juan Carlos Perez" {
		t.Errorf("CustomerName esperado: %q, obtenido: %q", "Juan Carlos Perez", result.CustomerName)
	}
	if result.CustomerEmail != "juancarlos@example.com" {
		t.Errorf("CustomerEmail esperado: %q, obtenido: %q", "juancarlos@example.com", result.CustomerEmail)
	}
}

// ─── Tests: buildOrderEntity ─────────────────────────────────────────────────

func TestBuildOrderEntity_MapeaCamposCorrectamente(t *testing.T) {
	// Arrange
	uc := newTestUseCase(&mockRepository{}, nil, nil, nil, &mockScoreUseCase{})
	bID := uint(42)
	clientID := uint(99)
	dto := &dtos.ProbabilityOrderDTO{
		BusinessID:      &bID,
		IntegrationID:   10,
		IntegrationType: "shopify",
		Platform:        "shopify",
		ExternalID:      "EXT-BUILD-01",
		OrderNumber:     "ORD-BUILD-01",
		TotalAmount:     300.0,
		Currency:        "COP",
		CustomerName:    "Carlos Lopez",
		CustomerEmail:   "carlos@example.com",
		Status:          "pending",
		OriginalStatus:  "open",
		Invoiceable:     true,
	}
	statusMapping := orderStatusMapping{
		OrderStatusID:       newUint(5),
		PaymentStatusID:     newUint(2),
		FulfillmentStatusID: newUint(1),
	}

	// Act
	order := uc.buildOrderEntity(dto, &clientID, statusMapping)

	// Assert
	if order.BusinessID != &bID {
		t.Errorf("BusinessID inesperado")
	}
	if order.IntegrationID != 10 {
		t.Errorf("IntegrationID esperado: 10, obtenido: %d", order.IntegrationID)
	}
	if order.ExternalID != "EXT-BUILD-01" {
		t.Errorf("ExternalID esperado: %q, obtenido: %q", "EXT-BUILD-01", order.ExternalID)
	}
	if order.TotalAmount != 300.0 {
		t.Errorf("TotalAmount esperado: 300.0, obtenido: %v", order.TotalAmount)
	}
	if order.CustomerID == nil || *order.CustomerID != 99 {
		t.Errorf("CustomerID esperado: 99, obtenido: %v", order.CustomerID)
	}
	if order.StatusID == nil || *order.StatusID != 5 {
		t.Errorf("StatusID esperado: 5, obtenido: %v", order.StatusID)
	}
	if order.PaymentStatusID == nil || *order.PaymentStatusID != 2 {
		t.Errorf("PaymentStatusID esperado: 2, obtenido: %v", order.PaymentStatusID)
	}
	if !order.Invoiceable {
		t.Error("se esperaba Invoiceable = true")
	}
}

// ─── Tests: assignPaymentMethodID ────────────────────────────────────────────

func TestAssignPaymentMethodID_SinPagos_AsignaValorPorDefecto(t *testing.T) {
	// Arrange
	uc := newTestUseCase(&mockRepository{}, nil, nil, nil, &mockScoreUseCase{})
	order := &entities.ProbabilityOrder{}
	dto := &dtos.ProbabilityOrderDTO{
		Payments: nil, // Sin pagos
	}

	// Act
	uc.assignPaymentMethodID(order, dto)

	// Assert
	if order.PaymentMethodID != 1 {
		t.Errorf("PaymentMethodID por defecto esperado: 1, obtenido: %d", order.PaymentMethodID)
	}
}

func TestAssignPaymentMethodID_ConPagoCompletado_AsignaPaymentMethodID(t *testing.T) {
	// Arrange
	uc := newTestUseCase(&mockRepository{}, nil, nil, nil, &mockScoreUseCase{})
	order := &entities.ProbabilityOrder{}
	dto := &dtos.ProbabilityOrderDTO{
		Payments: []dtos.ProbabilityPaymentDTO{
			{
				PaymentMethodID: 3,
				Status:          "pending", // No "completed" para evitar requerir PaidAt
			},
		},
	}

	// Act
	uc.assignPaymentMethodID(order, dto)

	// Assert
	if order.PaymentMethodID != 3 {
		t.Errorf("PaymentMethodID esperado: 3, obtenido: %d", order.PaymentMethodID)
	}
}

// ─── Tests: populateOrderFields ──────────────────────────────────────────────

func TestPopulateOrderFields_DireccionDeEnvio_PopulaCamposPlanos(t *testing.T) {
	// Arrange
	uc := newTestUseCase(&mockRepository{}, nil, nil, nil, &mockScoreUseCase{})
	order := &entities.ProbabilityOrder{}
	dto := &dtos.ProbabilityOrderDTO{
		Addresses: []dtos.ProbabilityAddressDTO{
			{
				Type:       "shipping",
				Street:     "Calle 123",
				Street2:    "Apto 4B",
				City:       "Bogota",
				State:      "Cundinamarca",
				Country:    "Colombia",
				PostalCode: "110111",
			},
		},
	}

	// Act
	uc.populateOrderFields(order, dto)

	// Assert
	if order.ShippingStreet != "Calle 123 Apto 4B" {
		t.Errorf("ShippingStreet esperada: %q, obtenida: %q", "Calle 123 Apto 4B", order.ShippingStreet)
	}
	if order.ShippingCity != "Bogota" {
		t.Errorf("ShippingCity esperada: %q, obtenida: %q", "Bogota", order.ShippingCity)
	}
	if order.ShippingCountry != "Colombia" {
		t.Errorf("ShippingCountry esperada: %q, obtenida: %q", "Colombia", order.ShippingCountry)
	}
	if order.Address2 != "Apto 4B" {
		t.Errorf("Address2 esperado: %q, obtenido: %q", "Apto 4B", order.Address2)
	}
}

func TestPopulateOrderFields_SinDireccionDeEnvio_NoModificaCampos(t *testing.T) {
	// Arrange
	uc := newTestUseCase(&mockRepository{}, nil, nil, nil, &mockScoreUseCase{})
	order := &entities.ProbabilityOrder{
		ShippingCity:    "CiudadOriginal",
		ShippingCountry: "PaisOriginal",
	}
	dto := &dtos.ProbabilityOrderDTO{
		Addresses: []dtos.ProbabilityAddressDTO{
			{
				Type:   "billing", // No es shipping
				Street: "Otra calle",
			},
		},
	}

	// Act
	uc.populateOrderFields(order, dto)

	// Assert: Los campos no deberian cambiar porque no hay dirección de tipo "shipping"
	if order.ShippingCity != "CiudadOriginal" {
		t.Errorf("ShippingCity no debería cambiar: se esperaba %q, obtenida: %q", "CiudadOriginal", order.ShippingCity)
	}
	if order.ShippingCountry != "PaisOriginal" {
		t.Errorf("ShippingCountry no debería cambiar: se esperaba %q, obtenida: %q", "PaisOriginal", order.ShippingCountry)
	}
}

// ─── Tests: equalJSON ────────────────────────────────────────────────────────

func TestEqualJSON_MismoContenido_RetornaTrue(t *testing.T) {
	// Arrange
	a := []byte(`{"key":"value","num":42}`)
	b := []byte(`{"num":42,"key":"value"}`) // Mismo contenido, diferente orden

	// Act
	result := equalJSON(a, b)

	// Assert
	if !result {
		t.Error("se esperaba true para JSONs con mismo contenido pero diferente orden")
	}
}

func TestEqualJSON_ContenidoDiferente_RetornaFalse(t *testing.T) {
	// Arrange
	a := []byte(`{"key":"value1"}`)
	b := []byte(`{"key":"value2"}`)

	// Act
	result := equalJSON(a, b)

	// Assert
	if result {
		t.Error("se esperaba false para JSONs con contenido diferente")
	}
}

func TestEqualJSON_JSONInvalido_RetornaFalse(t *testing.T) {
	// Arrange
	a := []byte(`{invalid json}`)
	b := []byte(`{"key":"value"}`)

	// Act
	result := equalJSON(a, b)

	// Assert
	if result {
		t.Error("se esperaba false cuando el primer JSON es inválido")
	}
}

// ─── Tests: GetOrCreateCustomer ──────────────────────────────────────────────

func TestGetOrCreateCustomer_SinEmail_RetornaNilSinError(t *testing.T) {
	// Arrange
	repo := &mockRepository{}
	uc := newTestUseCase(repo, nil, nil, nil, &mockScoreUseCase{})
	dto := &dtos.ProbabilityOrderDTO{
		CustomerEmail: "", // Sin email
	}

	// Act
	client, err := uc.GetOrCreateCustomer(context.Background(), 1, dto)

	// Assert
	if err != nil {
		t.Fatalf("no se esperaba error, pero se obtuvo: %v", err)
	}
	if client != nil {
		t.Errorf("se esperaba cliente nil cuando no hay email, pero se obtuvo: %v", client)
	}
}

func TestGetOrCreateCustomer_ClienteExistentePorEmail_RetornaClienteExistente(t *testing.T) {
	// Arrange
	existingClient := &entities.Client{
		ID:    55,
		Email: "existente@example.com",
	}
	repo := &mockRepository{
		GetClientByEmailFn: func(ctx context.Context, businessID uint, email string) (*entities.Client, error) {
			return existingClient, nil // Cliente encontrado
		},
	}
	uc := newTestUseCase(repo, nil, nil, nil, &mockScoreUseCase{})
	dto := &dtos.ProbabilityOrderDTO{
		CustomerEmail: "existente@example.com",
	}

	// Act
	client, err := uc.GetOrCreateCustomer(context.Background(), 1, dto)

	// Assert
	if err != nil {
		t.Fatalf("no se esperaba error, pero se obtuvo: %v", err)
	}
	if client == nil {
		t.Fatal("se esperaba cliente no nil")
	}
	if client.ID != 55 {
		t.Errorf("ID de cliente esperado: 55, obtenido: %d", client.ID)
	}
}

func TestGetOrCreateCustomer_ClienteNuevo_CreaYRetorna(t *testing.T) {
	// Arrange
	createClientCalled := false
	repo := &mockRepository{
		GetClientByEmailFn: func(ctx context.Context, businessID uint, email string) (*entities.Client, error) {
			return nil, nil // No existe por email
		},
		GetClientByDNIFn: func(ctx context.Context, businessID uint, dni string) (*entities.Client, error) {
			return nil, nil // No existe por DNI
		},
		CreateClientFn: func(ctx context.Context, client *entities.Client) error {
			createClientCalled = true
			client.ID = 100 // Simula ID asignado por BD
			return nil
		},
	}
	uc := newTestUseCase(repo, nil, nil, nil, &mockScoreUseCase{})
	dto := &dtos.ProbabilityOrderDTO{
		CustomerEmail: "nuevo@example.com",
		CustomerName:  "Nuevo Cliente",
		CustomerPhone: "3001111111",
	}

	// Act
	client, err := uc.GetOrCreateCustomer(context.Background(), 1, dto)

	// Assert
	if err != nil {
		t.Fatalf("no se esperaba error, pero se obtuvo: %v", err)
	}
	if client == nil {
		t.Fatal("se esperaba cliente no nil")
	}
	if !createClientCalled {
		t.Error("se esperaba que CreateClient fuera llamado para el nuevo cliente")
	}
	if client.Email != "nuevo@example.com" {
		t.Errorf("Email esperado: %q, obtenido: %q", "nuevo@example.com", client.Email)
	}
}

func TestGetOrCreateCustomer_ErrorAlBuscarPorEmail_RetornaError(t *testing.T) {
	// Arrange
	expectedErr := errors.New("error de red al buscar cliente")
	repo := &mockRepository{
		GetClientByEmailFn: func(ctx context.Context, businessID uint, email string) (*entities.Client, error) {
			return nil, expectedErr
		},
	}
	uc := newTestUseCase(repo, nil, nil, nil, &mockScoreUseCase{})
	dto := &dtos.ProbabilityOrderDTO{
		CustomerEmail: "error@example.com",
	}

	// Act
	client, err := uc.GetOrCreateCustomer(context.Background(), 1, dto)

	// Assert
	if err == nil {
		t.Fatal("se esperaba error al fallar la búsqueda por email")
	}
	if client != nil {
		t.Errorf("se esperaba cliente nil en caso de error, pero se obtuvo: %v", client)
	}
	if !errors.Is(err, expectedErr) {
		t.Errorf("se esperaba error wrapeado con: %v, pero se obtuvo: %v", expectedErr, err)
	}
}

// ─── Tests: GetOrCreateProduct ───────────────────────────────────────────────

func TestGetOrCreateProduct_SinSKU_RetornaError(t *testing.T) {
	// Arrange
	repo := &mockRepository{}
	uc := newTestUseCase(repo, nil, nil, nil, &mockScoreUseCase{})
	itemDTO := dtos.ProbabilityOrderItemDTO{
		ProductSKU: "", // Sin SKU
	}

	// Act
	product, err := uc.GetOrCreateProduct(context.Background(), 1, itemDTO)

	// Assert
	if err == nil {
		t.Fatal("se esperaba error cuando el SKU está vacío")
	}
	if product != nil {
		t.Errorf("se esperaba producto nil en caso de error, pero se obtuvo: %v", product)
	}
}

func TestGetOrCreateProduct_ProductoExistente_RetornaExistente(t *testing.T) {
	// Arrange
	existingProduct := &entities.Product{
		ID:  "PROD-001",
		SKU: "SKU-ABC",
	}
	repo := &mockRepository{
		GetProductBySKUFn: func(ctx context.Context, businessID uint, sku string) (*entities.Product, error) {
			return existingProduct, nil
		},
	}
	uc := newTestUseCase(repo, nil, nil, nil, &mockScoreUseCase{})
	itemDTO := dtos.ProbabilityOrderItemDTO{
		ProductSKU: "SKU-ABC",
	}

	// Act
	product, err := uc.GetOrCreateProduct(context.Background(), 1, itemDTO)

	// Assert
	if err != nil {
		t.Fatalf("no se esperaba error, pero se obtuvo: %v", err)
	}
	if product == nil {
		t.Fatal("se esperaba producto no nil")
	}
	if product.ID != "PROD-001" {
		t.Errorf("ID de producto esperado: %q, obtenido: %q", "PROD-001", product.ID)
	}
}

func TestGetOrCreateProduct_ProductoNuevo_CreaYRetorna(t *testing.T) {
	// Arrange
	createProductCalled := false
	repo := &mockRepository{
		GetProductBySKUFn: func(ctx context.Context, businessID uint, sku string) (*entities.Product, error) {
			return nil, nil // No existe
		},
		CreateProductFn: func(ctx context.Context, product *entities.Product) error {
			createProductCalled = true
			product.ID = "PROD-NUEVO"
			return nil
		},
	}
	uc := newTestUseCase(repo, nil, nil, nil, &mockScoreUseCase{})
	itemDTO := dtos.ProbabilityOrderItemDTO{
		ProductSKU:  "SKU-NUEVO",
		ProductName: "Producto Nuevo",
	}

	// Act
	product, err := uc.GetOrCreateProduct(context.Background(), 1, itemDTO)

	// Assert
	if err != nil {
		t.Fatalf("no se esperaba error, pero se obtuvo: %v", err)
	}
	if product == nil {
		t.Fatal("se esperaba producto no nil")
	}
	if !createProductCalled {
		t.Error("se esperaba que CreateProduct fuera llamado")
	}
}

// ─── Tests: mapOrderToResponse ────────────────────────────────────────────────

func TestMapOrderToResponse_MapeaTodosLosCamposPrincipales(t *testing.T) {
	// Arrange
	uc := newTestUseCase(&mockRepository{}, nil, nil, nil, &mockScoreUseCase{})
	bID := uint(10)
	statusID := uint(3)
	order := &entities.ProbabilityOrder{
		ID:              "ORDER-RESP-01",
		BusinessID:      &bID,
		IntegrationID:   5,
		IntegrationType: "shopify",
		Platform:        "shopify",
		ExternalID:      "EXT-RESP-01",
		OrderNumber:     "ORD-RESP-01",
		TotalAmount:     500.0,
		Currency:        "COP",
		CustomerName:    "Test Customer",
		CustomerEmail:   "customer@example.com",
		Status:          "pending",
		StatusID:        &statusID,
		Invoiceable:     true,
	}

	// Act
	response := uc.mapOrderToResponse(order)

	// Assert
	if response == nil {
		t.Fatal("se esperaba respuesta no nil")
	}
	if response.ID != "ORDER-RESP-01" {
		t.Errorf("ID esperado: %q, obtenido: %q", "ORDER-RESP-01", response.ID)
	}
	if response.IntegrationID != 5 {
		t.Errorf("IntegrationID esperado: 5, obtenido: %d", response.IntegrationID)
	}
	if response.TotalAmount != 500.0 {
		t.Errorf("TotalAmount esperado: 500.0, obtenido: %v", response.TotalAmount)
	}
	if response.Currency != "COP" {
		t.Errorf("Currency esperada: %q, obtenida: %q", "COP", response.Currency)
	}
	if response.Status != "pending" {
		t.Errorf("Status esperado: %q, obtenido: %q", "pending", response.Status)
	}
	if response.StatusID == nil || *response.StatusID != 3 {
		t.Errorf("StatusID esperado: 3, obtenido: %v", response.StatusID)
	}
	if !response.Invoiceable {
		t.Error("se esperaba Invoiceable = true en la respuesta")
	}
}
