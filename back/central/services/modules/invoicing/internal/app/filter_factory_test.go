package app

import (
	"fmt"
	"testing"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func f64p(v float64) *float64 {
	return &v
}

func intp(v int) *int {
	return &v
}

func strp(v string) *string {
	return &v
}

func boolp(v bool) *bool {
	return &v
}

func tipos(validators []FilterValidator) []string {
	out := make([]string, 0, len(validators))
	for _, v := range validators {
		out = append(out, fmt.Sprintf("%T", v))
	}
	return out
}

func TestCreateValidators_ConfigVacia_SoloValidaMoneda(t *testing.T) {
	got := CreateValidators(&entities.FilterConfig{}, "")

	require.Len(t, got, 1)
	assert.IsType(t, &CurrencyCOPValidator{}, got[0])
}

func TestCreateValidators_MonedaSiempreEsElPrimero(t *testing.T) {
	got := CreateValidators(&entities.FilterConfig{
		MinAmount:     f64p(1000),
		MaxAmount:     f64p(9000),
		PaymentStatus: strp("paid"),
	}, "")

	require.NotEmpty(t, got)
	assert.IsType(t, &CurrencyCOPValidator{}, got[0])
}

func TestCreateValidators_MontoMinimoYMaximo(t *testing.T) {
	got := CreateValidators(&entities.FilterConfig{
		MinAmount: f64p(10000),
		MaxAmount: f64p(500000),
	}, "")

	require.Len(t, got, 3)

	minV, ok := got[1].(*MinAmountValidator)
	require.True(t, ok)
	assert.InDelta(t, 10000.0, minV.MinAmount, 0.001)

	maxV, ok := got[2].(*MaxAmountValidator)
	require.True(t, ok)
	assert.InDelta(t, 500000.0, maxV.MaxAmount, 0.001)
}

func TestCreateValidators_SoloMontoMinimo(t *testing.T) {
	got := CreateValidators(&entities.FilterConfig{MinAmount: f64p(10000)}, "")

	require.Len(t, got, 2)
	assert.IsType(t, &MinAmountValidator{}, got[1])
}

func TestCreateValidators_EstadoDePagoConCOD(t *testing.T) {
	got := CreateValidators(&entities.FilterConfig{
		PaymentStatus: strp("paid"),
		InvoiceCOD:    boolp(true),
	}, "")

	require.Len(t, got, 2)
	v, ok := got[1].(*PaymentStatusValidator)
	require.True(t, ok)
	assert.Equal(t, "paid", v.RequiredStatus)
	assert.True(t, v.AllowCOD)
}

func TestCreateValidators_EstadoDePagoSinInvoiceCOD_NoPermiteCOD(t *testing.T) {
	got := CreateValidators(&entities.FilterConfig{PaymentStatus: strp("paid")}, "")

	v, ok := got[1].(*PaymentStatusValidator)
	require.True(t, ok)
	assert.False(t, v.AllowCOD, "sin invoice_cod explicito el COD no debe facturarse")
}

func TestCreateValidators_InvoiceCODFalseExplicito(t *testing.T) {
	got := CreateValidators(&entities.FilterConfig{
		PaymentStatus: strp("paid"),
		InvoiceCOD:    boolp(false),
	}, "")

	v, ok := got[1].(*PaymentStatusValidator)
	require.True(t, ok)
	assert.False(t, v.AllowCOD)
}

func TestCreateValidators_InvoiceCODSinPaymentStatus_NoCreaValidador(t *testing.T) {
	got := CreateValidators(&entities.FilterConfig{InvoiceCOD: boolp(true)}, "")

	require.Len(t, got, 1)
	assert.NotContains(t, tipos(got), "*app.PaymentStatusValidator")
}

func TestCreateValidators_MetodosDePago(t *testing.T) {
	got := CreateValidators(&entities.FilterConfig{PaymentMethods: []uint{1, 3, 7}}, "")

	require.Len(t, got, 2)
	v, ok := got[1].(*PaymentMethodsValidator)
	require.True(t, ok)
	assert.Equal(t, []uint{1, 3, 7}, v.AllowedMethods)
}

func TestCreateValidators_ListasVacias_NoCreanValidadores(t *testing.T) {
	got := CreateValidators(&entities.FilterConfig{
		PaymentMethods:      []uint{},
		OrderTypes:          []string{},
		ExcludeStatuses:     []string{},
		ExcludeProducts:     []string{},
		IncludeProductsOnly: []string{},
		CustomerTypes:       []string{},
		ExcludeCustomerIDs:  []string{},
		ShippingRegions:     []string{},
	}, "")

	assert.Len(t, got, 1)
}

func TestCreateValidators_FiltrosDeOrden(t *testing.T) {
	got := CreateValidators(&entities.FilterConfig{
		OrderTypes:      []string{"delivery", "pickup"},
		ExcludeStatuses: []string{"cancelled", "refunded"},
	}, "")

	require.Len(t, got, 3)

	tipoV, ok := got[1].(*OrderTypesValidator)
	require.True(t, ok)
	assert.Equal(t, []string{"delivery", "pickup"}, tipoV.AllowedTypes)

	excV, ok := got[2].(*ExcludeStatusesValidator)
	require.True(t, ok)
	assert.Equal(t, []string{"cancelled", "refunded"}, excV.ExcludedStatuses)
}

func TestCreateValidators_FiltrosDeProducto(t *testing.T) {
	got := CreateValidators(&entities.FilterConfig{
		ExcludeProducts:     []string{"SKU-1"},
		IncludeProductsOnly: []string{"SKU-2", "SKU-3"},
	}, "")

	require.Len(t, got, 3)

	excV, ok := got[1].(*ExcludeProductsValidator)
	require.True(t, ok)
	assert.Equal(t, []string{"SKU-1"}, excV.ExcludedSKUs)

	incV, ok := got[2].(*IncludeProductsOnlyValidator)
	require.True(t, ok)
	assert.Len(t, incV.AllowedSKUs, 2)
}

func TestCreateValidators_ConteoDeItems(t *testing.T) {
	cases := []struct {
		name   string
		config *entities.FilterConfig
		want   bool
	}{
		{name: "solo minimo", config: &entities.FilterConfig{MinItemsCount: intp(1)}, want: true},
		{name: "solo maximo", config: &entities.FilterConfig{MaxItemsCount: intp(50)}, want: true},
		{name: "ambos", config: &entities.FilterConfig{MinItemsCount: intp(1), MaxItemsCount: intp(50)}, want: true},
		{name: "ninguno", config: &entities.FilterConfig{}, want: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := CreateValidators(tc.config, "")
			assert.Equal(t, tc.want, len(got) == 2)
		})
	}
}

func TestCreateValidators_ConteoDeItems_CopiaLosLimites(t *testing.T) {
	got := CreateValidators(&entities.FilterConfig{
		MinItemsCount: intp(2),
		MaxItemsCount: intp(20),
	}, "")

	require.Len(t, got, 2)
	v, ok := got[1].(*ItemsCountValidator)
	require.True(t, ok)
	require.NotNil(t, v.MinCount)
	require.NotNil(t, v.MaxCount)
	assert.Equal(t, 2, *v.MinCount)
	assert.Equal(t, 20, *v.MaxCount)
}

func TestCreateValidators_FiltrosDeCliente(t *testing.T) {
	got := CreateValidators(&entities.FilterConfig{
		CustomerTypes:      []string{"natural"},
		ExcludeCustomerIDs: []string{"cust-1", "cust-2"},
	}, "")

	require.Len(t, got, 3)

	tipoV, ok := got[1].(*CustomerTypesValidator)
	require.True(t, ok)
	assert.Equal(t, []string{"natural"}, tipoV.AllowedTypes)

	excV, ok := got[2].(*ExcludeCustomersValidator)
	require.True(t, ok)
	assert.Len(t, excV.ExcludedCustomerIDs, 2)
}

func TestCreateValidators_RegionesDeEnvio(t *testing.T) {
	got := CreateValidators(&entities.FilterConfig{ShippingRegions: []string{"Bogota", "Medellin"}}, "")

	require.Len(t, got, 2)
	v, ok := got[1].(*ShippingRegionsValidator)
	require.True(t, ok)
	assert.Equal(t, []string{"Bogota", "Medellin"}, v.AllowedRegions)
}

func TestCreateValidators_RangoDeFechas_FormateaAISO(t *testing.T) {
	inicio := time.Date(2026, 8, 1, 13, 45, 0, 0, time.UTC)
	fin := time.Date(2026, 8, 31, 23, 0, 0, 0, time.UTC)

	got := CreateValidators(&entities.FilterConfig{
		DateRange: &entities.DateRangeFilter{StartDate: &inicio, EndDate: &fin},
	}, "")

	require.Len(t, got, 2)
	v, ok := got[1].(*DateRangeValidator)
	require.True(t, ok)
	require.NotNil(t, v.StartDate)
	require.NotNil(t, v.EndDate)
	assert.Equal(t, "2026-08-01", *v.StartDate)
	assert.Equal(t, "2026-08-31", *v.EndDate)
}

func TestCreateValidators_RangoDeFechas_SoloInicio(t *testing.T) {
	inicio := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)

	got := CreateValidators(&entities.FilterConfig{
		DateRange: &entities.DateRangeFilter{StartDate: &inicio},
	}, "")

	require.Len(t, got, 2)
	v, ok := got[1].(*DateRangeValidator)
	require.True(t, ok)
	require.NotNil(t, v.StartDate)
	assert.Nil(t, v.EndDate)
}

func TestCreateValidators_RangoDeFechasVacio_CreaValidadorSinLimites(t *testing.T) {
	got := CreateValidators(&entities.FilterConfig{DateRange: &entities.DateRangeFilter{}}, "")

	require.Len(t, got, 2)
	v, ok := got[1].(*DateRangeValidator)
	require.True(t, ok)
	assert.Nil(t, v.StartDate)
	assert.Nil(t, v.EndDate)
}

func TestCreateValidators_ConfiguracionCompleta_CreaTodos(t *testing.T) {
	inicio := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	fin := time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)

	got := CreateValidators(&entities.FilterConfig{
		MinAmount:           f64p(1000),
		MaxAmount:           f64p(999999),
		PaymentStatus:       strp("paid"),
		InvoiceCOD:          boolp(true),
		PaymentMethods:      []uint{1},
		OrderTypes:          []string{"delivery"},
		ExcludeStatuses:     []string{"cancelled"},
		ExcludeProducts:     []string{"SKU-X"},
		IncludeProductsOnly: []string{"SKU-Y"},
		MinItemsCount:       intp(1),
		MaxItemsCount:       intp(10),
		CustomerTypes:       []string{"juridica"},
		ExcludeCustomerIDs:  []string{"c-1"},
		ShippingRegions:     []string{"Bogota"},
		DateRange:           &entities.DateRangeFilter{StartDate: &inicio, EndDate: &fin},
	}, "")

	assert.Len(t, got, 14)

	nombres := tipos(got)
	esperados := []string{
		"*app.CurrencyCOPValidator",
		"*app.MinAmountValidator",
		"*app.MaxAmountValidator",
		"*app.PaymentStatusValidator",
		"*app.PaymentMethodsValidator",
		"*app.OrderTypesValidator",
		"*app.ExcludeStatusesValidator",
		"*app.ExcludeProductsValidator",
		"*app.IncludeProductsOnlyValidator",
		"*app.ItemsCountValidator",
		"*app.CustomerTypesValidator",
		"*app.ExcludeCustomersValidator",
		"*app.ShippingRegionsValidator",
		"*app.DateRangeValidator",
	}
	assert.Equal(t, esperados, nombres)
}
