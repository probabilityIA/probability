package usecaseorderscore

// Tests unitarios para el paquete usecaseorderscore.
//
// Estrategia:
//   - Testeamos CalculateOrderScore y GetStaticNegativeFactors directamente
//     (son métodos del struct, misma package -> acceso directo).
//   - Para CalculateAndUpdateOrderScore usamos el RepositoryMock de internal/mocks
//     que ya implementa ports.IRepository con testify/mock.
//   - IsCODPayment y RemoveAccents también se testean directamente.

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ─── Helpers ────────────────────────────────────────────────────────────────

// newUC construye una instancia de UseCaseOrderScore con un repo mock.
func newUC(repo *mocks.RepositoryMock) *UseCaseOrderScore {
	return &UseCaseOrderScore{repo: repo}
}

// perfectOrder retorna una orden que cumple todos los criterios
// y no genera ningún factor negativo (score esperado 100).
func perfectOrder() *entities.ProbabilityOrder {
	customerID := uint(42)
	return &entities.ProbabilityOrder{
		CustomerID:         &customerID,
		CustomerEmail:      "juan.perez@ejemplo.com",
		CustomerName:       "Juan Perez",
		Platform:           "Shopify",
		CustomerPhone:      "+573001234567",
		ShippingStreet:     "Calle 123 # 45-67",
		Address2:           "Apto 201",
		CustomerOrderCount: 5,
	}
}

// ─── Tests: CalculateOrderScore ─────────────────────────────────────────────

func TestCalculateOrderScore_OrdenPerfecta_Score100(t *testing.T) {
	uc := newUC(new(mocks.RepositoryMock))
	order := perfectOrder()

	score, factors := uc.CalculateOrderScore(order)

	assert.Equal(t, 100.0, score, "una orden perfecta debe tener score 100")
	assert.Empty(t, factors, "una orden perfecta no debe tener factores negativos")
}

func TestCalculateOrderScore_TodosLosFaltantes_Score30(t *testing.T) {
	uc := newUC(new(mocks.RepositoryMock))

	// Orden con todos los campos inválidos / vacíos
	order := &entities.ProbabilityOrder{
		CustomerEmail:      "correo-invalido",
		CustomerName:       "SoloNombre", // sin apellido (sin espacio)
		Platform:           "",
		CustomerPhone:      "",
		ShippingStreet:     "Cll",  // longitud <= 5
		Address2:           "",
		CustomerOrderCount: 0,
	}

	score, factors := uc.CalculateOrderScore(order)

	// 7 factores × −10 = 30
	assert.Equal(t, 30.0, score)
	assert.Len(t, factors, 7, "deben aparecer los 7 factores negativos")
	assert.Contains(t, factors, "Email válido")
	assert.Contains(t, factors, "Nombre y apellido")
	assert.Contains(t, factors, "Canal de venta")
	assert.Contains(t, factors, "Teléfono")
	assert.Contains(t, factors, "Dirección")
	assert.Contains(t, factors, "Complemento de dirección")
	assert.Contains(t, factors, "Historial de compra")
}

func TestCalculateOrderScore_UnFactorNegativo_Score90(t *testing.T) {
	uc := newUC(new(mocks.RepositoryMock))

	order := perfectOrder()
	order.CustomerPhone = "" // único factor negativo

	score, factors := uc.CalculateOrderScore(order)

	assert.Equal(t, 90.0, score)
	assert.Len(t, factors, 1)
	assert.Contains(t, factors, "Teléfono")
}

func TestCalculateOrderScore_NuncaMenorACero(t *testing.T) {
	uc := newUC(new(mocks.RepositoryMock))

	// Orden terrible + COD → podría generar un score negativo sin el límite
	order := &entities.ProbabilityOrder{
		CustomerEmail:      "",
		CustomerName:       "",
		Platform:           "",
		CustomerPhone:      "",
		ShippingStreet:     "",
		Address2:           "",
		CustomerOrderCount: 0,
	}

	// Agregar pago COD para aplicar reducción del 20 %
	gw := "cod-payment"
	order.Payments = []entities.ProbabilityPayment{
		{Gateway: &gw},
	}

	score, _ := uc.CalculateOrderScore(order)

	assert.GreaterOrEqual(t, score, 0.0, "el score nunca debe ser negativo")
}

func TestCalculateOrderScore_AplicaReduccionCOD(t *testing.T) {
	uc := newUC(new(mocks.RepositoryMock))

	// Orden perfecta (100) con COD → 100 × 0.8 = 80
	order := perfectOrder()
	gw := "cod"
	order.Payments = []entities.ProbabilityPayment{
		{Gateway: &gw},
	}

	score, factors := uc.CalculateOrderScore(order)

	assert.Equal(t, 80.0, score)
	assert.Contains(t, factors, "Pago Contra Entrega")
}

func TestCalculateOrderScore_Redondeo2Decimales(t *testing.T) {
	uc := newUC(new(mocks.RepositoryMock))

	// 3 factores negativos → 70, con COD → 70 × 0.8 = 56.0 (exacto)
	order := perfectOrder()
	order.CustomerPhone = ""
	order.Address2 = ""
	order.CustomerOrderCount = 0

	gw := "cash on delivery"
	order.Payments = []entities.ProbabilityPayment{
		{Gateway: &gw},
	}

	score, _ := uc.CalculateOrderScore(order)

	// Verificar que el resultado tiene a lo sumo 2 decimales
	rounded := float64(int(score*100)) / 100
	assert.Equal(t, rounded, score)
}

// ─── Tests: GetStaticNegativeFactors ─────────────────────────────────────────

func TestGetStaticNegativeFactors_EmailInvalido(t *testing.T) {
	uc := newUC(new(mocks.RepositoryMock))

	casos := []struct {
		email   string
		esValido bool
	}{
		{"juan@ejemplo.com", true},
		{"", false},
		{"sin-arroba", false},
		{"@nodomain", false},
		{"user@.com", false},
	}

	for _, tc := range casos {
		order := perfectOrder()
		order.CustomerEmail = tc.email
		factors := uc.GetStaticNegativeFactors(order)

		if tc.esValido {
			assert.NotContains(t, factors, "Email válido", "email '%s' debería ser válido", tc.email)
		} else {
			assert.Contains(t, factors, "Email válido", "email '%s' debería ser inválido", tc.email)
		}
	}
}

func TestGetStaticNegativeFactors_NombreApellido(t *testing.T) {
	uc := newUC(new(mocks.RepositoryMock))

	cases := []struct {
		name      string
		hayFactor bool
	}{
		{"Juan Perez", false},        // nombre y apellido
		{"Maria De Los Angeles", false}, // múltiples palabras
		{"Solo", true},              // sin espacio
		{"", true},                  // vacío
		{"   ", true},               // solo espacios
	}

	for _, tc := range cases {
		order := perfectOrder()
		order.CustomerName = tc.name
		factors := uc.GetStaticNegativeFactors(order)

		if tc.hayFactor {
			assert.Contains(t, factors, "Nombre y apellido", "nombre '%s' debería generar factor", tc.name)
		} else {
			assert.NotContains(t, factors, "Nombre y apellido", "nombre '%s' no debería generar factor", tc.name)
		}
	}
}

func TestGetStaticNegativeFactors_DireccionCorta(t *testing.T) {
	uc := newUC(new(mocks.RepositoryMock))

	// La condicion del codigo es: len(ShippingStreet) <= 5 → genera factor "Direccion"
	// Longitud 6 o mas → NO genera factor
	cases := []struct {
		street    string
		hayFactor bool
	}{
		{"Calle 123 # 45", false},     // 14 chars → no genera factor
		{"Avenida 45 # 10", false},    // 15 chars → no genera factor
		{"123456", false},             // 6 chars → no genera factor (primer valor sin factor)
		{"12345", true},               // 5 chars <= 5 → genera factor
		{"abc", true},                 // 3 chars <= 5 → genera factor
		{"", true},                    // 0 chars <= 5 → genera factor
	}

	for _, tc := range cases {
		order := perfectOrder()
		order.ShippingStreet = tc.street
		factors := uc.GetStaticNegativeFactors(order)

		if tc.hayFactor {
			assert.Contains(t, factors, "Dirección",
				"calle '%s' (len %d) debería generar factor Dirección", tc.street, len(tc.street))
		} else {
			assert.NotContains(t, factors, "Dirección",
				"calle '%s' (len %d) NO debería generar factor Dirección", tc.street, len(tc.street))
		}
	}
}

func TestGetStaticNegativeFactors_ComplementoDireccion_DesdeAddresses(t *testing.T) {
	uc := newUC(new(mocks.RepositoryMock))

	// Address2 vacío en el campo plano, pero existe en Addresses
	order := perfectOrder()
	order.Address2 = ""
	order.Addresses = []entities.ProbabilityAddress{
		{Type: "shipping", Street2: "Apto 301"},
	}

	factors := uc.GetStaticNegativeFactors(order)

	assert.NotContains(t, factors, "Complemento de dirección",
		"debe encontrar Address2 desde el slice Addresses")
}

func TestGetStaticNegativeFactors_ComplementoDireccion_DesdeChannelMetadata(t *testing.T) {
	uc := newUC(new(mocks.RepositoryMock))

	rawData, _ := json.Marshal(map[string]interface{}{
		"shipping_address": map[string]interface{}{
			"address2": "Piso 5",
		},
	})

	order := perfectOrder()
	order.Address2 = ""
	order.Addresses = nil
	order.ChannelMetadata = []entities.ProbabilityOrderChannelMetadata{
		{RawData: rawData},
	}

	factors := uc.GetStaticNegativeFactors(order)

	assert.NotContains(t, factors, "Complemento de dirección",
		"debe encontrar Address2 desde ChannelMetadata RawData")
}

func TestGetStaticNegativeFactors_HistorialCompra_Cero(t *testing.T) {
	uc := newUC(new(mocks.RepositoryMock))

	order := perfectOrder()
	order.CustomerOrderCount = 0

	factors := uc.GetStaticNegativeFactors(order)

	assert.Contains(t, factors, "Historial de compra")
}

func TestGetStaticNegativeFactors_HistorialCompra_Mayor(t *testing.T) {
	uc := newUC(new(mocks.RepositoryMock))

	order := perfectOrder()
	order.CustomerOrderCount = 3

	factors := uc.GetStaticNegativeFactors(order)

	assert.NotContains(t, factors, "Historial de compra")
}

// ─── Tests: IsCODPayment ──────────────────────────────────────────────────────

func TestIsCODPayment_GatewayConCod(t *testing.T) {
	uc := newUC(new(mocks.RepositoryMock))

	keywords := []string{"cod", "cash", "contra"}
	for _, kw := range keywords {
		gw := kw
		order := &entities.ProbabilityOrder{
			Payments: []entities.ProbabilityPayment{
				{Gateway: &gw},
			},
		}
		assert.True(t, uc.IsCODPayment(order), "gateway '%s' debe ser COD", kw)
	}
}

func TestIsCODPayment_GatewayNormal_NoEsCOD(t *testing.T) {
	uc := newUC(new(mocks.RepositoryMock))

	gw := "credit_card"
	order := &entities.ProbabilityOrder{
		Payments: []entities.ProbabilityPayment{
			{Gateway: &gw},
		},
	}

	assert.False(t, uc.IsCODPayment(order))
}

func TestIsCODPayment_CodTotalPositivo(t *testing.T) {
	uc := newUC(new(mocks.RepositoryMock))

	total := 50000.0
	order := &entities.ProbabilityOrder{
		CodTotal: &total,
	}

	assert.True(t, uc.IsCODPayment(order))
}

func TestIsCODPayment_PaymentDetails_Gateway(t *testing.T) {
	uc := newUC(new(mocks.RepositoryMock))

	details, _ := json.Marshal(map[string]interface{}{
		"gateway": "contraentrega",
	})

	order := &entities.ProbabilityOrder{
		PaymentDetails: details,
	}

	assert.True(t, uc.IsCODPayment(order))
}

func TestIsCODPayment_PaymentDetails_GatewayNames(t *testing.T) {
	uc := newUC(new(mocks.RepositoryMock))

	details, _ := json.Marshal(map[string]interface{}{
		"payment_gateway_names": []string{"credit_card", "cash_on_delivery"},
	})

	order := &entities.ProbabilityOrder{
		PaymentDetails: details,
	}

	assert.True(t, uc.IsCODPayment(order))
}

func TestIsCODPayment_Metadata_Tags(t *testing.T) {
	uc := newUC(new(mocks.RepositoryMock))

	meta, _ := json.Marshal(map[string]interface{}{
		"tags": "vip, contra entrega, express",
	})

	order := &entities.ProbabilityOrder{
		Metadata: meta,
	}

	assert.True(t, uc.IsCODPayment(order))
}

func TestIsCODPayment_SinIndicadores_False(t *testing.T) {
	uc := newUC(new(mocks.RepositoryMock))

	gw := "pse"
	order := &entities.ProbabilityOrder{
		Payments: []entities.ProbabilityPayment{
			{Gateway: &gw},
		},
	}

	assert.False(t, uc.IsCODPayment(order))
}

// ─── Tests: RemoveAccents ────────────────────────────────────────────────────

func TestRemoveAccents(t *testing.T) {
	uc := newUC(new(mocks.RepositoryMock))

	cases := []struct {
		input    string
		expected string
	}{
		{"café", "cafe"},
		{"Ñoño", "Nono"},
		{"José María", "Jose Maria"},
		{"sin acentos", "sin acentos"},
		{"", ""},
	}

	for _, tc := range cases {
		result := uc.RemoveAccents(tc.input)
		assert.Equal(t, tc.expected, result, "RemoveAccents('%s')", tc.input)
	}
}

// ─── Tests: CalculateAndUpdateOrderScore ─────────────────────────────────────

func TestCalculateAndUpdateOrderScore_Exitoso(t *testing.T) {
	repoMock := new(mocks.RepositoryMock)
	uc := newUC(repoMock)

	ctx := context.Background()
	orderID := "order-abc-123"
	customerID := uint(10)

	orden := &entities.ProbabilityOrder{
		ID:                 orderID,
		CustomerID:         &customerID,
		CustomerEmail:      "test@example.com",
		CustomerName:       "Ana Gomez",
		Platform:           "Shopify",
		CustomerPhone:      "+573009876543",
		ShippingStreet:     "Carrera 50 # 12-34",
		Address2:           "Oficina 5",
		CustomerOrderCount: 2,
		IntegrationID:      1,
	}

	// GetOrderByID retorna la orden
	repoMock.On("GetOrderByID", ctx, orderID).Return(orden, nil)

	// UpdateOrder debe ser llamado con la orden actualizada
	repoMock.On("UpdateOrder", ctx, mock.AnythingOfType("*entities.ProbabilityOrder")).Return(nil)

	err := uc.CalculateAndUpdateOrderScore(ctx, orderID)

	assert.NoError(t, err)
	assert.NotNil(t, orden.DeliveryProbability, "DeliveryProbability debe quedar asignado")
	repoMock.AssertExpectations(t)
}

func TestCalculateAndUpdateOrderScore_ErrorAlObtenerOrden(t *testing.T) {
	repoMock := new(mocks.RepositoryMock)
	uc := newUC(repoMock)

	ctx := context.Background()
	orderID := "order-no-existe"
	dbErr := errors.New("record not found")

	repoMock.On("GetOrderByID", ctx, orderID).Return(nil, dbErr)

	err := uc.CalculateAndUpdateOrderScore(ctx, orderID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get order")
	repoMock.AssertExpectations(t)
	repoMock.AssertNotCalled(t, "UpdateOrder", mock.Anything, mock.Anything)
}

func TestCalculateAndUpdateOrderScore_ErrorAlActualizar(t *testing.T) {
	repoMock := new(mocks.RepositoryMock)
	uc := newUC(repoMock)

	ctx := context.Background()
	orderID := "order-update-fail"
	customerID := uint(5)

	// CustomerOrderCount > 0 e IntegrationID > 0:
	// la condicion shouldCheckDB = (IntegrationID==0 || CustomerOrderCount==0) es FALSE
	// por lo tanto NO se llama CountOrdersByClientID y el test es predecible.
	orden := &entities.ProbabilityOrder{
		ID:                 orderID,
		CustomerID:         &customerID,
		IntegrationID:      3,
		CustomerEmail:      "a@b.com",
		CustomerName:       "Luis Torres",
		Platform:           "Woo",
		CustomerPhone:      "+1",
		ShippingStreet:     "Calle larga 999",
		Address2:           "Apto 1",
		CustomerOrderCount: 2, // > 0 y IntegrationID != 0 → no consulta DB
	}

	repoMock.On("GetOrderByID", ctx, orderID).Return(orden, nil)
	updateErr := errors.New("db write error")
	repoMock.On("UpdateOrder", ctx, mock.AnythingOfType("*entities.ProbabilityOrder")).Return(updateErr)

	err := uc.CalculateAndUpdateOrderScore(ctx, orderID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update order with score")
	repoMock.AssertExpectations(t)
}

func TestCalculateAndUpdateOrderScore_ConsultaHistorialCliente_Local(t *testing.T) {
	// Para órdenes locales (IntegrationID == 0), siempre consulta la DB
	repoMock := new(mocks.RepositoryMock)
	uc := newUC(repoMock)

	ctx := context.Background()
	orderID := "local-order-1"
	customerID := uint(7)

	orden := &entities.ProbabilityOrder{
		ID:                 orderID,
		CustomerID:         &customerID,
		IntegrationID:      0, // local
		CustomerEmail:      "local@test.com",
		CustomerName:       "Pedro Ramirez",
		Platform:           "manual",
		CustomerPhone:      "+1",
		ShippingStreet:     "Diagonal 10 # 20",
		Address2:           "Casa",
		CustomerOrderCount: 0, // se actualizará desde la DB
	}

	repoMock.On("GetOrderByID", ctx, orderID).Return(orden, nil)
	// La DB reporta 3 órdenes previas
	repoMock.On("CountOrdersByClientID", ctx, customerID).Return(int64(3), nil)
	repoMock.On("UpdateOrder", ctx, mock.AnythingOfType("*entities.ProbabilityOrder")).Return(nil)

	err := uc.CalculateAndUpdateOrderScore(ctx, orderID)

	assert.NoError(t, err)
	// CustomerOrderCount debe haberse actualizado con el valor de la DB
	assert.Equal(t, 3, orden.CustomerOrderCount)
	repoMock.AssertExpectations(t)
}

func TestCalculateAndUpdateOrderScore_ConsultaHistorialCliente_IntegracionCount0DBMayor1(t *testing.T) {
	// Integración con CustomerOrderCount == 0 y DB > 1 → debe recuperar historial
	repoMock := new(mocks.RepositoryMock)
	uc := newUC(repoMock)

	ctx := context.Background()
	orderID := "shopify-order-old"
	customerID := uint(20)

	orden := &entities.ProbabilityOrder{
		ID:                 orderID,
		CustomerID:         &customerID,
		IntegrationID:      5, // integración
		CustomerEmail:      "shopify@test.com",
		CustomerName:       "Carla Reyes",
		Platform:           "Shopify",
		CustomerPhone:      "+57",
		ShippingStreet:     "Transversal 80 # 5",
		Address2:           "Bl 3",
		CustomerOrderCount: 0, // dato perdido
	}

	repoMock.On("GetOrderByID", ctx, orderID).Return(orden, nil)
	// DB tiene 5 órdenes → cliente recurrente
	repoMock.On("CountOrdersByClientID", ctx, customerID).Return(int64(5), nil)
	repoMock.On("UpdateOrder", ctx, mock.AnythingOfType("*entities.ProbabilityOrder")).Return(nil)

	err := uc.CalculateAndUpdateOrderScore(ctx, orderID)

	assert.NoError(t, err)
	assert.Equal(t, 5, orden.CustomerOrderCount, "historial debe recuperarse de la DB cuando DB > 1")
	repoMock.AssertExpectations(t)
}

func TestCalculateAndUpdateOrderScore_NegativeFactors_OrdenPerfecta(t *testing.T) {
	// Verifica que NegativeFactors quede como "[]" cuando no hay factores
	repoMock := new(mocks.RepositoryMock)
	uc := newUC(repoMock)

	ctx := context.Background()
	orderID := "perfect-order-id"
	customerID := uint(99)

	orden := perfectOrder()
	orden.ID = orderID
	orden.CustomerID = &customerID
	orden.IntegrationID = 1 // integración con count > 0

	repoMock.On("GetOrderByID", ctx, orderID).Return(orden, nil)
	repoMock.On("UpdateOrder", ctx, mock.AnythingOfType("*entities.ProbabilityOrder")).Return(nil)

	err := uc.CalculateAndUpdateOrderScore(ctx, orderID)

	assert.NoError(t, err)
	assert.JSONEq(t, "[]", string(orden.NegativeFactors))
	repoMock.AssertExpectations(t)
}
