package usecaseupdateorder

import (
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func fPtr(v float64) *float64 {
	return &v
}

func direccionEnvio() dtos.ProbabilityAddressDTO {
	return dtos.ProbabilityAddressDTO{
		Type:       "shipping",
		Street:     "Calle 1 # 2-3",
		Street2:    "Apto 101",
		City:       "Bogota",
		State:      "Cundinamarca",
		Country:    "Colombia",
		PostalCode: "110111",
		Latitude:   fPtr(4.65),
		Longitude:  fPtr(-74.05),
	}
}

func TestUpdateShippingAddress_SinDirecciones_NoCambiaNada(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)

	changed := uc.updateShippingAddress(&entities.ProbabilityOrder{}, &dtos.ProbabilityOrderDTO{})

	assert.False(t, changed)
}

func TestUpdateShippingAddress_ActualizaTodosLosCampos(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)
	order := &entities.ProbabilityOrder{}

	changed := uc.updateShippingAddress(order, &dtos.ProbabilityOrderDTO{
		Addresses: []dtos.ProbabilityAddressDTO{direccionEnvio()},
	})

	assert.True(t, changed)
	assert.Equal(t, "Calle 1 # 2-3", order.ShippingStreet)
	assert.Equal(t, "Apto 101", order.Address2)
	assert.Equal(t, "Bogota", order.ShippingCity)
	assert.Equal(t, "Cundinamarca", order.ShippingState)
	assert.Equal(t, "Colombia", order.ShippingCountry)
	assert.Equal(t, "110111", order.ShippingPostalCode)
	require.NotNil(t, order.ShippingLat)
	require.NotNil(t, order.ShippingLng)
	assert.InDelta(t, 4.65, *order.ShippingLat, 0.001)
	assert.InDelta(t, -74.05, *order.ShippingLng, 0.001)
}

func TestUpdateShippingAddress_IgnoraDireccionDeFacturacion(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)
	order := &entities.ProbabilityOrder{}

	changed := uc.updateShippingAddress(order, &dtos.ProbabilityOrderDTO{
		Addresses: []dtos.ProbabilityAddressDTO{{
			Type:   "billing",
			Street: "Otra calle",
			City:   "Medellin",
		}},
	})

	assert.False(t, changed)
	assert.Equal(t, "", order.ShippingStreet)
	assert.Equal(t, "", order.ShippingCity)
}

func TestUpdateShippingAddress_TomaLaPrimeraDeEnvio(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)
	order := &entities.ProbabilityOrder{}

	uc.updateShippingAddress(order, &dtos.ProbabilityOrderDTO{
		Addresses: []dtos.ProbabilityAddressDTO{
			{Type: "billing", City: "Cali"},
			{Type: "shipping", City: "Bogota"},
			{Type: "shipping", City: "Medellin"},
		},
	})

	assert.Equal(t, "Bogota", order.ShippingCity)
}

func TestUpdateShippingAddress_MismosValores_NoMarcaCambio(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)
	order := &entities.ProbabilityOrder{
		ShippingStreet:     "Calle 1 # 2-3",
		Address2:           "Apto 101",
		ShippingCity:       "Bogota",
		ShippingState:      "Cundinamarca",
		ShippingCountry:    "Colombia",
		ShippingPostalCode: "110111",
		ShippingLat:        fPtr(4.65),
		ShippingLng:        fPtr(-74.05),
	}

	changed := uc.updateShippingAddress(order, &dtos.ProbabilityOrderDTO{
		Addresses: []dtos.ProbabilityAddressDTO{direccionEnvio()},
	})

	assert.False(t, changed)
}

func TestUpdateShippingAddress_CamposVacios_NoBorranLosExistentes(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)
	order := &entities.ProbabilityOrder{
		ShippingStreet:  "Calle 1 # 2-3",
		ShippingCity:    "Bogota",
		ShippingCountry: "Colombia",
		ShippingLat:     fPtr(4.65),
	}

	changed := uc.updateShippingAddress(order, &dtos.ProbabilityOrderDTO{
		Addresses: []dtos.ProbabilityAddressDTO{{Type: "shipping"}},
	})

	assert.False(t, changed)
	assert.Equal(t, "Calle 1 # 2-3", order.ShippingStreet)
	assert.Equal(t, "Bogota", order.ShippingCity)
	require.NotNil(t, order.ShippingLat)
}

func TestUpdateShippingAddress_SoloCoordenadas(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)
	order := &entities.ProbabilityOrder{ShippingLat: fPtr(1.0), ShippingLng: fPtr(2.0)}

	changed := uc.updateShippingAddress(order, &dtos.ProbabilityOrderDTO{
		Addresses: []dtos.ProbabilityAddressDTO{{
			Type:      "shipping",
			Latitude:  fPtr(4.65),
			Longitude: fPtr(-74.05),
		}},
	})

	assert.True(t, changed)
	assert.InDelta(t, 4.65, *order.ShippingLat, 0.001)
	assert.InDelta(t, -74.05, *order.ShippingLng, 0.001)
}

func TestUpdateAdditionalFields_SinDatos_NoCambiaNada(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)
	order := &entities.ProbabilityOrder{}

	changed := uc.updateAdditionalFields(order, &dtos.ProbabilityOrderDTO{})

	assert.False(t, changed)
}

func TestUpdateAdditionalFields_ActualizaTodo(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)
	order := &entities.ProbabilityOrder{}
	aprobado := true

	changed := uc.updateAdditionalFields(order, &dtos.ProbabilityOrderDTO{
		OrderTypeID:     uintPtr(2),
		OrderTypeName:   "delivery",
		Notes:           sPtr("dejar en porteria"),
		Coupon:          sPtr("DESC10"),
		Approved:        &aprobado,
		UserID:          uintPtr(42),
		UserName:        "Ana",
		Invoiceable:     true,
		InvoiceURL:      sPtr("https://s3/factura.pdf"),
		InvoiceID:       sPtr("FV-100"),
		InvoiceProvider: sPtr("siigo"),
		OrderStatusURL:  "https://shop/status",
	})

	assert.True(t, changed)
	assert.Equal(t, uint(2), *order.OrderTypeID)
	assert.Equal(t, "delivery", order.OrderTypeName)
	assert.Equal(t, "dejar en porteria", *order.Notes)
	assert.Equal(t, "DESC10", *order.Coupon)
	assert.True(t, *order.Approved)
	assert.Equal(t, uint(42), *order.UserID)
	assert.Equal(t, "Ana", order.UserName)
	assert.True(t, order.Invoiceable)
	assert.Equal(t, "https://s3/factura.pdf", *order.InvoiceURL)
	assert.Equal(t, "FV-100", *order.InvoiceID)
	assert.Equal(t, "siigo", *order.InvoiceProvider)
	assert.Equal(t, "https://shop/status", order.OrderStatusURL)
}

func TestUpdateAdditionalFields_MismosValores_NoMarcaCambio(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)
	aprobado := true
	order := &entities.ProbabilityOrder{
		OrderTypeID:     uintPtr(2),
		OrderTypeName:   "delivery",
		Notes:           sPtr("nota"),
		Coupon:          sPtr("DESC10"),
		Approved:        &aprobado,
		UserID:          uintPtr(42),
		UserName:        "Ana",
		Invoiceable:     true,
		InvoiceURL:      sPtr("https://s3/factura.pdf"),
		InvoiceID:       sPtr("FV-100"),
		InvoiceProvider: sPtr("siigo"),
		OrderStatusURL:  "https://shop/status",
	}

	changed := uc.updateAdditionalFields(order, &dtos.ProbabilityOrderDTO{
		OrderTypeID:     uintPtr(2),
		OrderTypeName:   "delivery",
		Notes:           sPtr("nota"),
		Coupon:          sPtr("DESC10"),
		Approved:        &aprobado,
		UserID:          uintPtr(42),
		UserName:        "Ana",
		Invoiceable:     true,
		InvoiceURL:      sPtr("https://s3/factura.pdf"),
		InvoiceID:       sPtr("FV-100"),
		InvoiceProvider: sPtr("siigo"),
		OrderStatusURL:  "https://shop/status",
	})

	assert.False(t, changed)
}

func TestUpdateAdditionalFields_InvoiceableSeSincronizaSiempre(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)
	order := &entities.ProbabilityOrder{Invoiceable: true}

	changed := uc.updateAdditionalFields(order, &dtos.ProbabilityOrderDTO{Invoiceable: false})

	assert.True(t, changed)
	assert.False(t, order.Invoiceable, "invoiceable no tiene puntero, un false explicito si desmarca la orden")
}

func TestUpdateAdditionalFields_NulosNoBorranLosExistentes(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)
	order := &entities.ProbabilityOrder{
		Notes:         sPtr("nota original"),
		Coupon:        sPtr("DESC10"),
		OrderTypeName: "delivery",
		UserName:      "Ana",
		InvoiceID:     sPtr("FV-100"),
	}

	uc.updateAdditionalFields(order, &dtos.ProbabilityOrderDTO{})

	assert.Equal(t, "nota original", *order.Notes)
	assert.Equal(t, "DESC10", *order.Coupon)
	assert.Equal(t, "delivery", order.OrderTypeName)
	assert.Equal(t, "Ana", order.UserName)
	assert.Equal(t, "FV-100", *order.InvoiceID)
}

func TestUpdateAdditionalFields_ApprovedFalseExplicito_SiCambia(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)
	aprobado := true
	rechazado := false
	order := &entities.ProbabilityOrder{Approved: &aprobado}

	changed := uc.updateAdditionalFields(order, &dtos.ProbabilityOrderDTO{Approved: &rechazado})

	assert.True(t, changed)
	assert.False(t, *order.Approved)
}
