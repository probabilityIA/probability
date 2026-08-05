package usecaseupdateorder

import (
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateFinancialFields_ActualizaTodosLosMontos(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)
	order := &entities.ProbabilityOrder{}

	changed := uc.updateFinancialFields(order, &dtos.ProbabilityOrderDTO{
		Subtotal:     100000,
		Tax:          19000,
		Discount:     5000,
		ShippingCost: 12000,
		TotalAmount:  126000,
		Currency:     "COP",
		CodTotal:     fPtr(126000),
		IsCod:        true,
	})

	assert.True(t, changed)
	assert.InDelta(t, 100000.0, order.Subtotal, 0.001)
	assert.InDelta(t, 19000.0, order.Tax, 0.001)
	assert.InDelta(t, 5000.0, order.Discount, 0.001)
	assert.InDelta(t, 12000.0, order.ShippingCost, 0.001)
	assert.InDelta(t, 126000.0, order.TotalAmount, 0.001)
	assert.Equal(t, "COP", order.Currency)
	require.NotNil(t, order.CodTotal)
	assert.InDelta(t, 126000.0, *order.CodTotal, 0.001)
	assert.True(t, order.IsCod)
}

func TestUpdateFinancialFields_MismosValores_NoMarcaCambio(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)
	order := &entities.ProbabilityOrder{
		Subtotal:     100000,
		Tax:          19000,
		Discount:     5000,
		ShippingCost: 12000,
		TotalAmount:  126000,
		Currency:     "COP",
		CodTotal:     fPtr(126000),
		IsCod:        true,
	}

	changed := uc.updateFinancialFields(order, &dtos.ProbabilityOrderDTO{
		Subtotal:     100000,
		Tax:          19000,
		Discount:     5000,
		ShippingCost: 12000,
		TotalAmount:  126000,
		Currency:     "COP",
		CodTotal:     fPtr(126000),
		IsCod:        true,
	})

	assert.False(t, changed)
}

func TestUpdateFinancialFields_SubtotalCero_NoBorraElExistente(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)
	order := &entities.ProbabilityOrder{Subtotal: 100000, TotalAmount: 126000}

	uc.updateFinancialFields(order, &dtos.ProbabilityOrderDTO{Subtotal: 0, TotalAmount: 0})

	assert.InDelta(t, 100000.0, order.Subtotal, 0.001, "un subtotal cero no debe borrar el monto real")
	assert.InDelta(t, 126000.0, order.TotalAmount, 0.001)
}

func TestUpdateFinancialFields_ImpuestoCeroSiSeAplica(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)
	order := &entities.ProbabilityOrder{Tax: 19000, Discount: 5000, ShippingCost: 12000}

	changed := uc.updateFinancialFields(order, &dtos.ProbabilityOrderDTO{})

	assert.True(t, changed)
	assert.InDelta(t, 0.0, order.Tax, 0.001, "impuesto, descuento y envio si aceptan cero explicito")
	assert.InDelta(t, 0.0, order.Discount, 0.001)
	assert.InDelta(t, 0.0, order.ShippingCost, 0.001)
}

func TestUpdateFinancialFields_MonedaVacia_NoBorraLaExistente(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)
	order := &entities.ProbabilityOrder{Currency: "COP"}

	uc.updateFinancialFields(order, &dtos.ProbabilityOrderDTO{Currency: ""})

	assert.Equal(t, "COP", order.Currency)
}

func TestUpdateFinancialFields_CodTotalNulo_NoBorraElExistente(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)
	order := &entities.ProbabilityOrder{CodTotal: fPtr(150000), IsCod: true}

	uc.updateFinancialFields(order, &dtos.ProbabilityOrderDTO{IsCod: true})

	require.NotNil(t, order.CodTotal)
	assert.InDelta(t, 150000.0, *order.CodTotal, 0.001)
}

func TestUpdateFinancialFields_IsCodSeSincronizaSiempre(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)
	order := &entities.ProbabilityOrder{IsCod: true}

	changed := uc.updateFinancialFields(order, &dtos.ProbabilityOrderDTO{IsCod: false})

	assert.True(t, changed)
	assert.False(t, order.IsCod, "is_cod no tiene puntero, un false explicito si desmarca la contraentrega")
}

func TestUpdateStructuredData_SinDatos_NoCambiaNada(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)

	changed := uc.updateStructuredData(&entities.ProbabilityOrder{}, &dtos.ProbabilityOrderDTO{})

	assert.False(t, changed)
}

func TestUpdateStructuredData_ActualizaLosCincoJSONB(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)
	order := &entities.ProbabilityOrder{}

	changed := uc.updateStructuredData(order, &dtos.ProbabilityOrderDTO{
		Metadata:           []byte(`{"canal":"shopify"}`),
		FinancialDetails:   []byte(`{"financial_status":"paid"}`),
		ShippingDetails:    []byte(`{"carrier":"Servientrega"}`),
		PaymentDetails:     []byte(`{"gateway":"bold"}`),
		FulfillmentDetails: []byte(`{"fulfillment_status":"shipped"}`),
	})

	assert.True(t, changed)
	assert.JSONEq(t, `{"canal":"shopify"}`, string(order.Metadata))
	assert.JSONEq(t, `{"financial_status":"paid"}`, string(order.FinancialDetails))
	assert.JSONEq(t, `{"carrier":"Servientrega"}`, string(order.ShippingDetails))
	assert.JSONEq(t, `{"gateway":"bold"}`, string(order.PaymentDetails))
	assert.JSONEq(t, `{"fulfillment_status":"shipped"}`, string(order.FulfillmentDetails))
}

func TestUpdateStructuredData_MismoJSONConOtroOrden_NoMarcaCambio(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)
	order := &entities.ProbabilityOrder{
		Metadata:           []byte(`{"a":1,"b":2}`),
		FinancialDetails:   []byte(`{"x":true}`),
		ShippingDetails:    []byte(`{"y":"z"}`),
		PaymentDetails:     []byte(`{"p":1}`),
		FulfillmentDetails: []byte(`{"f":null}`),
	}

	changed := uc.updateStructuredData(order, &dtos.ProbabilityOrderDTO{
		Metadata:           []byte(`{"b":2,"a":1}`),
		FinancialDetails:   []byte(`{"x":true}`),
		ShippingDetails:    []byte(`{"y":"z"}`),
		PaymentDetails:     []byte(`{"p":1}`),
		FulfillmentDetails: []byte(`{"f":null}`),
	})

	assert.False(t, changed, "el mismo JSON con las claves en otro orden no cuenta como cambio")
}

func TestUpdateStructuredData_JSONDistinto_MarcaCambio(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)
	order := &entities.ProbabilityOrder{Metadata: []byte(`{"canal":"shopify"}`)}

	changed := uc.updateStructuredData(order, &dtos.ProbabilityOrderDTO{
		Metadata: []byte(`{"canal":"woocommerce"}`),
	})

	assert.True(t, changed)
	assert.JSONEq(t, `{"canal":"woocommerce"}`, string(order.Metadata))
}

func TestUpdateStructuredData_CampoVacio_NoBorraElExistente(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)
	order := &entities.ProbabilityOrder{
		Metadata:       []byte(`{"canal":"shopify"}`),
		PaymentDetails: []byte(`{"gateway":"bold"}`),
	}

	changed := uc.updateStructuredData(order, &dtos.ProbabilityOrderDTO{})

	assert.False(t, changed)
	assert.JSONEq(t, `{"canal":"shopify"}`, string(order.Metadata))
	assert.JSONEq(t, `{"gateway":"bold"}`, string(order.PaymentDetails))
}

func TestUpdateStructuredData_SoloUnoCambia(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)
	order := &entities.ProbabilityOrder{
		Metadata:        []byte(`{"a":1}`),
		ShippingDetails: []byte(`{"carrier":"Servientrega"}`),
	}

	changed := uc.updateStructuredData(order, &dtos.ProbabilityOrderDTO{
		Metadata:        []byte(`{"a":1}`),
		ShippingDetails: []byte(`{"carrier":"Interrapidisimo"}`),
	})

	assert.True(t, changed)
	assert.JSONEq(t, `{"a":1}`, string(order.Metadata))
	assert.JSONEq(t, `{"carrier":"Interrapidisimo"}`, string(order.ShippingDetails))
}
