package usecaseupdateorder

import (
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
	"github.com/stretchr/testify/assert"
)

func intPtr(v int) *int {
	return &v
}

func TestUpdateCustomerFields_SinDatos_NoCambiaNada(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)

	changed := uc.updateCustomerFields(&entities.ProbabilityOrder{}, &dtos.ProbabilityOrderDTO{})

	assert.False(t, changed)
}

func TestUpdateCustomerFields_ActualizaTodo(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)
	order := &entities.ProbabilityOrder{}

	changed := uc.updateCustomerFields(order, &dtos.ProbabilityOrderDTO{
		CustomerName:       "Ana Gomez",
		CustomerEmail:      "ana@test.com",
		CustomerPhone:      "3001234567",
		CustomerDNI:        "1020304050",
		CustomerOrderCount: intPtr(5),
		CustomerTotalSpent: sPtr("750000"),
	})

	assert.True(t, changed)
	assert.Equal(t, "Ana Gomez", order.CustomerName)
	assert.Equal(t, "ana@test.com", order.CustomerEmail)
	assert.Equal(t, "3001234567", order.CustomerPhone)
	assert.Equal(t, "1020304050", order.CustomerDNI)
	assert.Equal(t, 5, order.CustomerOrderCount)
	assert.Equal(t, "750000", order.CustomerTotalSpent)
}

func TestUpdateCustomerFields_MismosValores_NoMarcaCambio(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)
	order := &entities.ProbabilityOrder{
		CustomerName:       "Ana Gomez",
		CustomerEmail:      "ana@test.com",
		CustomerPhone:      "3001234567",
		CustomerDNI:        "1020304050",
		CustomerOrderCount: 5,
		CustomerTotalSpent: "750000",
	}

	changed := uc.updateCustomerFields(order, &dtos.ProbabilityOrderDTO{
		CustomerName:       "Ana Gomez",
		CustomerEmail:      "ana@test.com",
		CustomerPhone:      "3001234567",
		CustomerDNI:        "1020304050",
		CustomerOrderCount: intPtr(5),
		CustomerTotalSpent: sPtr("750000"),
	})

	assert.False(t, changed)
}

func TestUpdateCustomerFields_VaciosNoBorranLosExistentes(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)
	order := &entities.ProbabilityOrder{
		CustomerName:  "Ana Gomez",
		CustomerEmail: "ana@test.com",
		CustomerPhone: "3001234567",
		CustomerDNI:   "1020304050",
	}

	changed := uc.updateCustomerFields(order, &dtos.ProbabilityOrderDTO{})

	assert.False(t, changed)
	assert.Equal(t, "Ana Gomez", order.CustomerName)
	assert.Equal(t, "ana@test.com", order.CustomerEmail)
	assert.Equal(t, "3001234567", order.CustomerPhone)
	assert.Equal(t, "1020304050", order.CustomerDNI)
}

func TestUpdateCustomerFields_ContadorCeroExplicito_SiCambia(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)
	order := &entities.ProbabilityOrder{CustomerOrderCount: 5, CustomerTotalSpent: "750000"}

	changed := uc.updateCustomerFields(order, &dtos.ProbabilityOrderDTO{
		CustomerOrderCount: intPtr(0),
		CustomerTotalSpent: sPtr("0"),
	})

	assert.True(t, changed)
	assert.Equal(t, 0, order.CustomerOrderCount)
	assert.Equal(t, "0", order.CustomerTotalSpent)
}

func TestUpdateCustomerFields_SoloUnCampo(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)
	order := &entities.ProbabilityOrder{CustomerName: "Ana", CustomerEmail: "ana@test.com"}

	changed := uc.updateCustomerFields(order, &dtos.ProbabilityOrderDTO{CustomerPhone: "3009999999"})

	assert.True(t, changed)
	assert.Equal(t, "Ana", order.CustomerName)
	assert.Equal(t, "3009999999", order.CustomerPhone)
}

func TestUpdateShippingFields_CombinaTrackingYDireccion(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)
	order := &entities.ProbabilityOrder{}

	changed := uc.updateShippingFields(order, &dtos.ProbabilityOrderDTO{
		Shipments: []dtos.ProbabilityShipmentDTO{{TrackingNumber: sPtr("034057376067")}},
		Addresses: []dtos.ProbabilityAddressDTO{{Type: "shipping", City: "Bogota"}},
	})

	assert.True(t, changed)
	assert.Equal(t, "034057376067", *order.TrackingNumber)
	assert.Equal(t, "Bogota", order.ShippingCity)
}

func TestUpdateShippingFields_SoloTracking(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)
	order := &entities.ProbabilityOrder{}

	changed := uc.updateShippingFields(order, &dtos.ProbabilityOrderDTO{
		Shipments: []dtos.ProbabilityShipmentDTO{{TrackingNumber: sPtr("034057376067")}},
	})

	assert.True(t, changed)
	assert.Equal(t, "", order.ShippingCity)
}

func TestUpdateShippingFields_SoloDireccion(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)
	order := &entities.ProbabilityOrder{}

	changed := uc.updateShippingFields(order, &dtos.ProbabilityOrderDTO{
		Addresses: []dtos.ProbabilityAddressDTO{{Type: "shipping", City: "Cali"}},
	})

	assert.True(t, changed)
	assert.Nil(t, order.TrackingNumber)
	assert.Equal(t, "Cali", order.ShippingCity)
}

func TestUpdateShippingFields_SinNada_NoMarcaCambio(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)

	changed := uc.updateShippingFields(&entities.ProbabilityOrder{}, &dtos.ProbabilityOrderDTO{})

	assert.False(t, changed)
}
