package usecaseupdateorder

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
	CreateOrderFn                                        func(ctx context.Context, order *entities.ProbabilityOrder) error
	GetOrderByIDFn                                       func(ctx context.Context, id string) (*entities.ProbabilityOrder, error)
	GetOrderByInternalNumberFn                           func(ctx context.Context, internalNumber string) (*entities.ProbabilityOrder, error)
	GetOrderByOrderNumberFn                              func(ctx context.Context, orderNumber string) (*entities.ProbabilityOrder, error)
	ListOrdersFn                                         func(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]entities.ProbabilityOrder, int64, error)
	UpdateOrderFn                                        func(ctx context.Context, order *entities.ProbabilityOrder) error
	DeleteOrderFn                                        func(ctx context.Context, id string) error
	GetOrderRawFn                                        func(ctx context.Context, id string) (*entities.ProbabilityOrderChannelMetadata, error)
	CountOrdersByClientIDFn                              func(ctx context.Context, clientID uint) (int64, error)
	GetLastManualOrderNumberFn                           func(ctx context.Context, businessID uint) (int, error)
	GetFirstIntegrationIDByBusinessIDFn                  func(ctx context.Context, businessID uint) (uint, error)
	GetPlatformIntegrationIDByBusinessIDFn               func(ctx context.Context, businessID uint) (uint, error)
	OrderExistsFn                                        func(ctx context.Context, externalID string, integrationID uint) (bool, error)
	GetOrderByExternalIDFn                               func(ctx context.Context, externalID string, integrationID uint) (*entities.ProbabilityOrder, error)
	CreateOrderItemsFn                                   func(ctx context.Context, items []*entities.ProbabilityOrderItem) error
	CreateAddressesFn                                    func(ctx context.Context, addresses []*entities.ProbabilityAddress) error
	CreatePaymentsFn                                     func(ctx context.Context, payments []*entities.ProbabilityPayment) error
	CreateShipmentsFn                                    func(ctx context.Context, shipments []*entities.ProbabilityShipment) error
	CreateChannelMetadataFn                              func(ctx context.Context, metadata *entities.ProbabilityOrderChannelMetadata) error
	GetProductBySKUFn                                    func(ctx context.Context, businessID uint, sku string) (*entities.Product, error)
	CreateProductFn                                      func(ctx context.Context, product *entities.Product) error
	GetClientByEmailFn                                   func(ctx context.Context, businessID uint, email string) (*entities.Client, error)
	GetClientByDNIFn                                     func(ctx context.Context, businessID uint, dni string) (*entities.Client, error)
	CreateClientFn                                       func(ctx context.Context, client *entities.Client) error
	CreateOrderErrorFn                                   func(ctx context.Context, orderError *entities.OrderError) error
	GetOrderStatusIDByIntegrationTypeAndOriginalStatusFn func(ctx context.Context, integrationTypeID uint, originalStatus string) (*uint, error)
	GetOrderStatusIDByCodeFn                             func(ctx context.Context, code string) (*uint, error)
	GetPaymentStatusIDByCodeFn                           func(ctx context.Context, code string) (*uint, error)
	GetFulfillmentStatusIDByCodeFn                       func(ctx context.Context, code string) (*uint, error)
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
func (m *mockRepository) GetOrderStatusIDByCode(ctx context.Context, code string) (*uint, error) {
	if m.GetOrderStatusIDByCodeFn != nil {
		return m.GetOrderStatusIDByCodeFn(ctx, code)
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

func (m *mockRepository) UpdateProductPrice(ctx context.Context, productID string, price float64) error {
	return nil
}

// ─── Mock: IOrderRabbitPublisher ────────────────────────────────────────────

type mockRabbitPublisher struct {
	PublishOrderCreatedFn          func(ctx context.Context, order *entities.ProbabilityOrder) error
	PublishOrderUpdatedFn          func(ctx context.Context, order *entities.ProbabilityOrder) error
	PublishOrderCancelledFn        func(ctx context.Context, order *entities.ProbabilityOrder) error
	PublishOrderStatusChangedFn    func(ctx context.Context, order *entities.ProbabilityOrder, previousStatus, currentStatus string) error
	PublishConfirmationRequestedFn func(ctx context.Context, order *entities.ProbabilityOrder) error
	PublishOrderEventFn            func(ctx context.Context, event *entities.OrderEvent, order *entities.ProbabilityOrder) error
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
	PublishSyncOrderCreatedFn  func(ctx context.Context, integrationID uint, businessID *uint, data map[string]interface{})
	PublishSyncOrderUpdatedFn  func(ctx context.Context, integrationID uint, businessID *uint, data map[string]interface{})
	PublishSyncOrderRejectedFn func(ctx context.Context, integrationID uint, businessID *uint, data map[string]interface{})
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
func (m *mockIntegrationEventPublisher) PublishSyncOrderRejected(ctx context.Context, integrationID uint, businessID *uint, data map[string]interface{}) {
	if m.PublishSyncOrderRejectedFn != nil {
		m.PublishSyncOrderRejectedFn(ctx, integrationID, businessID, data)
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

// ─── Helpers ────────────────────────────────────────────────────────────────

func newTestUpdateUseCase(
	repo *mockRepository,
	rabbitPublisher *mockRabbitPublisher,
	integrationEventPublisher *mockIntegrationEventPublisher,
	scoreUseCase *mockScoreUseCase,
) *UseCaseUpdateOrder {
	uc := &UseCaseUpdateOrder{
		repo:         repo,
		logger:       &mockLogger{},
		scoreUseCase: scoreUseCase,
	}
	if rabbitPublisher != nil {
		uc.rabbitEventPublisher = rabbitPublisher
	}
	if integrationEventPublisher != nil {
		uc.integrationEventPublisher = integrationEventPublisher
	}
	return uc
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
	uc := newTestUpdateUseCase(repo, nil, nil, &mockScoreUseCase{})

	// DTO con los mismos valores que la orden existente (sin cambios)
	dto := &dtos.ProbabilityOrderDTO{
		Status:        "pending", // Igual al existente
		TotalAmount:   0,         // No actualiza (dto.TotalAmount <= 0)
		Currency:      "",        // No actualiza (vacío)
		CustomerEmail: "",        // No actualiza (vacío)
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
	uc := newTestUpdateUseCase(repo, nil, nil, &mockScoreUseCase{})

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
	uc := newTestUpdateUseCase(repo, nil, nil, &mockScoreUseCase{})

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
	uc := newTestUpdateUseCase(repo, nil, nil, &mockScoreUseCase{})

	dto := &dtos.ProbabilityOrderDTO{
		TotalAmount: 200.0, // Valor diferente al existente
		Currency:    "USD",  // Moneda diferente
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
	uc := newTestUpdateUseCase(repo, nil, nil, &mockScoreUseCase{})

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
