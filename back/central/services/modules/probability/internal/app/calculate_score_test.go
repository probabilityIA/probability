package app

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"testing"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/probability/internal/domain/entities"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)


type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) GetOrderForScoring(ctx context.Context, id string) (*entities.ScoreOrder, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.ScoreOrder), args.Error(1)
}

func (m *mockRepository) CountOrdersByClientID(ctx context.Context, clientID uint) (int64, error) {
	args := m.Called(ctx, clientID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockRepository) UpdateOrderScore(ctx context.Context, orderID string, score float64, factors []byte, breakdown []byte) error {
	args := m.Called(ctx, orderID, score, factors, breakdown)
	return args.Error(0)
}

func (m *mockRepository) GetCustomerOrderHistory(ctx context.Context, customerID uint, excludeOrderID string) (*entities.CustomerHistory, error) {
	args := m.Called(ctx, customerID, excludeOrderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.CustomerHistory), args.Error(1)
}

func (m *mockRepository) GetCustomerDeliveryHistory(ctx context.Context, customerID uint) (*entities.DeliveryHistory, error) {
	args := m.Called(ctx, customerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.DeliveryHistory), args.Error(1)
}

func (m *mockRepository) GetOrderItemCount(ctx context.Context, orderID string) (int, error) {
	args := m.Called(ctx, orderID)
	return args.Get(0).(int), args.Error(1)
}

func (m *mockRepository) GetPaymentMethodCategory(ctx context.Context, paymentMethodID uint) (string, error) {
	args := m.Called(ctx, paymentMethodID)
	return args.Get(0).(string), args.Error(1)
}

type mockPublisher struct {
	mock.Mock
}

func (m *mockPublisher) PublishScoreCalculated(ctx context.Context, orderID, orderNumber string, businessID, integrationID uint) error {
	args := m.Called(ctx, orderID, orderNumber, businessID, integrationID)
	return args.Error(0)
}


func newUC(repo *mockRepository, pub *mockPublisher) *UseCaseScore {
	return &UseCaseScore{repo: repo, publisher: pub, log: log.New()}
}

func newPureUC() *UseCaseScore {
	return &UseCaseScore{}
}

func perfectOrder() *entities.ScoreOrder {
	customerID := uint(42)
	now := time.Now()
	firstOrder := now.AddDate(0, -6, 0)
	lastOrder := now.AddDate(0, 0, -3)
	return &entities.ScoreOrder{
		CustomerID:         &customerID,
		CustomerEmail:      "juan.perez@ejemplo.com",
		CustomerName:       "Juan Perez",
		Platform:           "Shopify",
		CustomerPhone:      "+573001234567",
		ShippingStreet:     "Calle 123 # 45-67",
		Address2:           "Apto 201",
		CustomerOrderCount: 5,
		TotalAmount:        150000,
		IsPaid:             true,
		OrderItemCount:     2,
		CustomerHistory: &entities.CustomerHistory{
			TotalOrders:       10,
			TotalSpent:        1500000,
			AvgOrderValue:     150000,
			FirstOrderDate:    &firstOrder,
			LastOrderDate:     &lastOrder,
			NoveltyCount:      0,
			CODOrderCount:     0,
			DistinctAddresses: 2,
			FailedPayments:    0,
		},
		DeliveryHistory: &entities.DeliveryHistory{
			TotalShipments:  10,
			FailedShipments: 0,
		},
	}
}

func setupEnrichmentMocks(repoMock *mockRepository, ctx context.Context, customerID uint, orderID string) {
	repoMock.On("GetCustomerOrderHistory", ctx, customerID, orderID).Return(&entities.CustomerHistory{}, nil).Maybe()
	repoMock.On("GetCustomerDeliveryHistory", ctx, customerID).Return(&entities.DeliveryHistory{}, nil).Maybe()
	repoMock.On("GetOrderItemCount", ctx, orderID).Return(0, nil).Maybe()
	repoMock.On("GetPaymentMethodCategory", ctx, mock.AnythingOfType("uint")).Return("", nil).Maybe()
}


func TestCalculateOrderScore_OrdenPerfecta_Score100(t *testing.T) {
	uc := newPureUC()
	order := perfectOrder()
	score, factors, breakdown := uc.CalculateOrderScore(order)
	assert.Equal(t, 100.0, score, "una orden perfecta debe tener score 100")
	assert.Empty(t, factors, "una orden perfecta no debe tener factores negativos")
	assert.NotNil(t, breakdown, "breakdown debe estar presente")
	assert.Equal(t, score, breakdown.FinalScore)
	assert.Len(t, breakdown.Categories, 5, "debe haber 5 categorias")
}

func TestCalculateOrderScore_TodosLosFaltantes(t *testing.T) {
	uc := newPureUC()
	order := &entities.ScoreOrder{
		CustomerEmail: "correo-invalido", CustomerName: "SoloNombre",
		Platform: "", CustomerPhone: "", ShippingStreet: "Cll",
		Address2: "", CustomerOrderCount: 0,
	}
	score, factors, _ := uc.CalculateOrderScore(order)
	assert.Less(t, score, 55.0, "score debe ser bajo con todos los datos faltantes")
	assert.GreaterOrEqual(t, score, 0.0)
	assert.Contains(t, factors, "Email válido")
	assert.Contains(t, factors, "Nombre y apellido")
	assert.Contains(t, factors, "Canal de venta")
	assert.Contains(t, factors, "Teléfono")
	assert.Contains(t, factors, "Dirección")
	assert.Contains(t, factors, "Complemento de dirección")
	assert.Contains(t, factors, "Historial de compra")
}

func TestCalculateOrderScore_UnFactorNegativo(t *testing.T) {
	uc := newPureUC()
	order := perfectOrder()
	order.CustomerPhone = ""
	score, factors, _ := uc.CalculateOrderScore(order)
	assert.Less(t, score, 100.0)
	assert.Greater(t, score, 80.0)
	assert.Contains(t, factors, "Teléfono")
}

func TestCalculateOrderScore_NuncaMenorACero(t *testing.T) {
	uc := newPureUC()
	order := &entities.ScoreOrder{}
	gw := "cod-payment"
	order.Payments = []entities.ScorePayment{{Gateway: &gw}}
	score, _, _ := uc.CalculateOrderScore(order)
	assert.GreaterOrEqual(t, score, 0.0)
}

func TestCalculateOrderScore_AplicaReduccionCOD(t *testing.T) {
	uc := newPureUC()
	order := perfectOrder()
	gw := "cod"
	order.Payments = []entities.ScorePayment{{Gateway: &gw}}
	score, factors, _ := uc.CalculateOrderScore(order)
	assert.Less(t, score, 100.0)
	assert.Contains(t, factors, "Pago Contra Entrega")
}

func TestCalculateOrderScore_Redondeo2Decimales(t *testing.T) {
	uc := newPureUC()
	order := perfectOrder()
	order.CustomerPhone = ""
	order.Address2 = ""
	order.CustomerOrderCount = 0
	gw := "cash on delivery"
	order.Payments = []entities.ScorePayment{{Gateway: &gw}}
	score, _, _ := uc.CalculateOrderScore(order)
	rounded := math.Round(score*100) / 100
	assert.Equal(t, rounded, score)
}


func TestGetStaticNegativeFactors_EmailInvalido(t *testing.T) {
	uc := newPureUC()
	casos := []struct {
		email    string
		esValido bool
	}{
		{"juan@ejemplo.com", true}, {"", false}, {"sin-arroba", false},
		{"@nodomain", false}, {"user@.com", false},
	}
	for _, tc := range casos {
		order := perfectOrder()
		order.CustomerEmail = tc.email
		factors := uc.GetStaticNegativeFactors(order)
		if tc.esValido {
			assert.NotContains(t, factors, "Email válido")
		} else {
			assert.Contains(t, factors, "Email válido")
		}
	}
}

func TestGetStaticNegativeFactors_NombreApellido(t *testing.T) {
	uc := newPureUC()
	cases := []struct {
		name      string
		hayFactor bool
	}{
		{"Juan Perez", false}, {"Maria De Los Angeles", false},
		{"Solo", true}, {"", true}, {"   ", true},
	}
	for _, tc := range cases {
		order := perfectOrder()
		order.CustomerName = tc.name
		factors := uc.GetStaticNegativeFactors(order)
		if tc.hayFactor {
			assert.Contains(t, factors, "Nombre y apellido")
		} else {
			assert.NotContains(t, factors, "Nombre y apellido")
		}
	}
}

func TestGetStaticNegativeFactors_DireccionCorta(t *testing.T) {
	uc := newPureUC()
	cases := []struct {
		street    string
		hayFactor bool
	}{
		{"Calle 123 # 45", false}, {"123456", false},
		{"12345", true}, {"abc", true}, {"", true},
	}
	for _, tc := range cases {
		order := perfectOrder()
		order.ShippingStreet = tc.street
		factors := uc.GetStaticNegativeFactors(order)
		if tc.hayFactor {
			assert.Contains(t, factors, "Dirección")
		} else {
			assert.NotContains(t, factors, "Dirección")
		}
	}
}

func TestGetStaticNegativeFactors_ComplementoDireccion_DesdeAddresses(t *testing.T) {
	uc := newPureUC()
	order := perfectOrder()
	order.Address2 = ""
	order.Addresses = []entities.ScoreAddress{{Type: "shipping", Street2: "Apto 301"}}
	factors := uc.GetStaticNegativeFactors(order)
	assert.NotContains(t, factors, "Complemento de dirección")
}

func TestGetStaticNegativeFactors_ComplementoDireccion_DesdeChannelMetadata(t *testing.T) {
	uc := newPureUC()
	rawData, _ := json.Marshal(map[string]interface{}{
		"shipping_address": map[string]interface{}{"address2": "Piso 5"},
	})
	order := perfectOrder()
	order.Address2 = ""
	order.Addresses = nil
	order.ChannelMetadata = []entities.ScoreChannelMetadata{{RawData: rawData}}
	factors := uc.GetStaticNegativeFactors(order)
	assert.NotContains(t, factors, "Complemento de dirección")
}

func TestGetStaticNegativeFactors_HistorialCompra_Cero(t *testing.T) {
	uc := newPureUC()
	order := perfectOrder()
	order.CustomerOrderCount = 0
	factors := uc.GetStaticNegativeFactors(order)
	assert.Contains(t, factors, "Historial de compra")
}

func TestGetStaticNegativeFactors_HistorialCompra_Mayor(t *testing.T) {
	uc := newPureUC()
	order := perfectOrder()
	order.CustomerOrderCount = 3
	factors := uc.GetStaticNegativeFactors(order)
	assert.NotContains(t, factors, "Historial de compra")
}


func TestIsCODPayment_GatewayConCod(t *testing.T) {
	uc := newPureUC()
	for _, kw := range []string{"cod", "cash", "contra"} {
		gw := kw
		order := &entities.ScoreOrder{Payments: []entities.ScorePayment{{Gateway: &gw}}}
		assert.True(t, uc.IsCODPayment(order))
	}
}

func TestIsCODPayment_GatewayNormal_NoEsCOD(t *testing.T) {
	uc := newPureUC()
	gw := "credit_card"
	order := &entities.ScoreOrder{Payments: []entities.ScorePayment{{Gateway: &gw}}}
	assert.False(t, uc.IsCODPayment(order))
}

func TestIsCODPayment_CodTotalPositivo(t *testing.T) {
	uc := newPureUC()
	total := 50000.0
	assert.True(t, uc.IsCODPayment(&entities.ScoreOrder{CodTotal: &total}))
}

func TestIsCODPayment_PaymentDetails_Gateway(t *testing.T) {
	uc := newPureUC()
	details, _ := json.Marshal(map[string]interface{}{"gateway": "contraentrega"})
	assert.True(t, uc.IsCODPayment(&entities.ScoreOrder{PaymentDetails: details}))
}

func TestIsCODPayment_PaymentDetails_GatewayNames(t *testing.T) {
	uc := newPureUC()
	details, _ := json.Marshal(map[string]interface{}{
		"payment_gateway_names": []string{"credit_card", "cash_on_delivery"},
	})
	assert.True(t, uc.IsCODPayment(&entities.ScoreOrder{PaymentDetails: details}))
}

func TestIsCODPayment_Metadata_Tags(t *testing.T) {
	uc := newPureUC()
	meta, _ := json.Marshal(map[string]interface{}{"tags": "vip, contra entrega, express"})
	assert.True(t, uc.IsCODPayment(&entities.ScoreOrder{Metadata: meta}))
}

func TestIsCODPayment_SinIndicadores_False(t *testing.T) {
	uc := newPureUC()
	gw := "pse"
	assert.False(t, uc.IsCODPayment(&entities.ScoreOrder{Payments: []entities.ScorePayment{{Gateway: &gw}}}))
}


func TestRemoveAccents(t *testing.T) {
	uc := newPureUC()
	cases := []struct{ input, expected string }{
		{"café", "cafe"}, {"Ñoño", "Nono"}, {"José María", "Jose Maria"},
		{"sin acentos", "sin acentos"}, {"", ""},
	}
	for _, tc := range cases {
		assert.Equal(t, tc.expected, uc.RemoveAccents(tc.input))
	}
}


func TestCalculateAndUpdateOrderScore_Exitoso(t *testing.T) {
	repoMock := new(mockRepository)
	pubMock := new(mockPublisher)
	uc := newUC(repoMock, pubMock)
	ctx := context.Background()
	orderID := "order-abc-123"
	customerID := uint(10)

	orden := &entities.ScoreOrder{
		ID: orderID, CustomerID: &customerID,
		CustomerEmail: "test@example.com", CustomerName: "Ana Gomez",
		Platform: "Shopify", CustomerPhone: "+573009876543",
		ShippingStreet: "Carrera 50 # 12-34", Address2: "Oficina 5",
		CustomerOrderCount: 2, IntegrationID: 1, OrderNumber: "ORD-001",
	}

	repoMock.On("GetOrderForScoring", ctx, orderID).Return(orden, nil)
	setupEnrichmentMocks(repoMock, ctx, customerID, orderID)
	repoMock.On("UpdateOrderScore", ctx, orderID, mock.AnythingOfType("float64"), mock.AnythingOfType("[]uint8"), mock.AnythingOfType("[]uint8")).Return(nil)
	pubMock.On("PublishScoreCalculated", ctx, orderID, "ORD-001", mock.AnythingOfType("uint"), mock.AnythingOfType("uint")).Return(nil)

	err := uc.CalculateAndUpdateOrderScore(ctx, orderID)
	assert.NoError(t, err)
	repoMock.AssertExpectations(t)
	pubMock.AssertExpectations(t)
}

func TestCalculateAndUpdateOrderScore_ErrorAlObtenerOrden(t *testing.T) {
	repoMock := new(mockRepository)
	pubMock := new(mockPublisher)
	uc := newUC(repoMock, pubMock)
	ctx := context.Background()
	orderID := "order-no-existe"

	repoMock.On("GetOrderForScoring", ctx, orderID).Return(nil, errors.New("record not found"))
	err := uc.CalculateAndUpdateOrderScore(ctx, orderID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get order for scoring")
	repoMock.AssertNotCalled(t, "UpdateOrderScore", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestCalculateAndUpdateOrderScore_ErrorAlActualizar(t *testing.T) {
	repoMock := new(mockRepository)
	pubMock := new(mockPublisher)
	uc := newUC(repoMock, pubMock)
	ctx := context.Background()
	orderID := "order-update-fail"
	customerID := uint(5)

	orden := &entities.ScoreOrder{
		ID: orderID, CustomerID: &customerID, IntegrationID: 3,
		CustomerEmail: "a@b.com", CustomerName: "Luis Torres", Platform: "Woo",
		CustomerPhone: "+1", ShippingStreet: "Calle larga 999", Address2: "Apto 1",
		CustomerOrderCount: 2, OrderNumber: "ORD-FAIL",
	}

	repoMock.On("GetOrderForScoring", ctx, orderID).Return(orden, nil)
	setupEnrichmentMocks(repoMock, ctx, customerID, orderID)
	repoMock.On("UpdateOrderScore", ctx, orderID, mock.AnythingOfType("float64"), mock.AnythingOfType("[]uint8"), mock.AnythingOfType("[]uint8")).Return(errors.New("db write error"))

	err := uc.CalculateAndUpdateOrderScore(ctx, orderID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update order score")
}

func TestCalculateAndUpdateOrderScore_ConsultaHistorialCliente_Local(t *testing.T) {
	repoMock := new(mockRepository)
	pubMock := new(mockPublisher)
	uc := newUC(repoMock, pubMock)
	ctx := context.Background()
	orderID := "local-order-1"
	customerID := uint(7)

	orden := &entities.ScoreOrder{
		ID: orderID, CustomerID: &customerID, IntegrationID: 0,
		CustomerEmail: "local@test.com", CustomerName: "Pedro Ramirez", Platform: "manual",
		CustomerPhone: "+1", ShippingStreet: "Diagonal 10 # 20", Address2: "Casa",
		CustomerOrderCount: 0, OrderNumber: "ORD-LOCAL",
	}

	repoMock.On("GetOrderForScoring", ctx, orderID).Return(orden, nil)
	repoMock.On("CountOrdersByClientID", ctx, customerID).Return(int64(3), nil)
	setupEnrichmentMocks(repoMock, ctx, customerID, orderID)
	repoMock.On("UpdateOrderScore", ctx, orderID, mock.AnythingOfType("float64"), mock.AnythingOfType("[]uint8"), mock.AnythingOfType("[]uint8")).Return(nil)
	pubMock.On("PublishScoreCalculated", ctx, orderID, "ORD-LOCAL", mock.AnythingOfType("uint"), mock.AnythingOfType("uint")).Return(nil)

	err := uc.CalculateAndUpdateOrderScore(ctx, orderID)
	assert.NoError(t, err)
	assert.Equal(t, 3, orden.CustomerOrderCount)
}

func TestCalculateAndUpdateOrderScore_IntegracionCount0DBMayor1(t *testing.T) {
	repoMock := new(mockRepository)
	pubMock := new(mockPublisher)
	uc := newUC(repoMock, pubMock)
	ctx := context.Background()
	orderID := "shopify-order-old"
	customerID := uint(20)

	orden := &entities.ScoreOrder{
		ID: orderID, CustomerID: &customerID, IntegrationID: 5,
		CustomerEmail: "shopify@test.com", CustomerName: "Carla Reyes", Platform: "Shopify",
		CustomerPhone: "+57", ShippingStreet: "Transversal 80 # 5", Address2: "Bl 3",
		CustomerOrderCount: 0, OrderNumber: "ORD-SHOPIFY",
	}

	repoMock.On("GetOrderForScoring", ctx, orderID).Return(orden, nil)
	repoMock.On("CountOrdersByClientID", ctx, customerID).Return(int64(5), nil)
	setupEnrichmentMocks(repoMock, ctx, customerID, orderID)
	repoMock.On("UpdateOrderScore", ctx, orderID, mock.AnythingOfType("float64"), mock.AnythingOfType("[]uint8"), mock.AnythingOfType("[]uint8")).Return(nil)
	pubMock.On("PublishScoreCalculated", ctx, orderID, "ORD-SHOPIFY", mock.AnythingOfType("uint"), mock.AnythingOfType("uint")).Return(nil)

	err := uc.CalculateAndUpdateOrderScore(ctx, orderID)
	assert.NoError(t, err)
	assert.Equal(t, 5, orden.CustomerOrderCount)
}

func TestCalculateAndUpdateOrderScore_NegativeFactors_OrdenPerfecta(t *testing.T) {
	repoMock := new(mockRepository)
	pubMock := new(mockPublisher)
	uc := newUC(repoMock, pubMock)
	ctx := context.Background()
	orderID := "perfect-order-id"
	customerID := uint(99)

	orden := perfectOrder()
	orden.ID = orderID
	orden.CustomerID = &customerID
	orden.IntegrationID = 1
	orden.OrderNumber = "ORD-PERFECT"

	repoMock.On("GetOrderForScoring", ctx, orderID).Return(orden, nil)
	setupEnrichmentMocks(repoMock, ctx, customerID, orderID)
	repoMock.On("UpdateOrderScore", ctx, orderID, mock.AnythingOfType("float64"), []byte("[]"), mock.AnythingOfType("[]uint8")).Return(nil)
	pubMock.On("PublishScoreCalculated", ctx, orderID, "ORD-PERFECT", mock.AnythingOfType("uint"), mock.AnythingOfType("uint")).Return(nil)

	err := uc.CalculateAndUpdateOrderScore(ctx, orderID)
	assert.NoError(t, err)
	repoMock.AssertExpectations(t)
}
