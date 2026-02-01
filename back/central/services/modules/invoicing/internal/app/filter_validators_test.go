package app

import (
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/ports"
	"github.com/stretchr/testify/assert"
)

// ═══════════════════════════════════════════════════════════════
// TESTS DE VALIDADORES DE MONTO
// ═══════════════════════════════════════════════════════════════

func TestMinAmountValidator(t *testing.T) {
	validator := &MinAmountValidator{MinAmount: 100000}

	t.Run("Orden por encima del mínimo - debe pasar", func(t *testing.T) {
		order := &ports.OrderData{TotalAmount: 150000}
		err := validator.Validate(order)
		assert.Nil(t, err)
	})

	t.Run("Orden exactamente en el mínimo - debe pasar", func(t *testing.T) {
		order := &ports.OrderData{TotalAmount: 100000}
		err := validator.Validate(order)
		assert.Nil(t, err)
	})

	t.Run("Orden por debajo del mínimo - debe fallar", func(t *testing.T) {
		order := &ports.OrderData{TotalAmount: 50000}
		err := validator.Validate(order)
		assert.Equal(t, errors.ErrOrderBelowMinAmount, err)
	})
}

func TestMaxAmountValidator(t *testing.T) {
	validator := &MaxAmountValidator{MaxAmount: 5000000}

	t.Run("Orden por debajo del máximo - debe pasar", func(t *testing.T) {
		order := &ports.OrderData{TotalAmount: 150000}
		err := validator.Validate(order)
		assert.Nil(t, err)
	})

	t.Run("Orden exactamente en el máximo - debe pasar", func(t *testing.T) {
		order := &ports.OrderData{TotalAmount: 5000000}
		err := validator.Validate(order)
		assert.Nil(t, err)
	})

	t.Run("Orden por encima del máximo - debe fallar", func(t *testing.T) {
		order := &ports.OrderData{TotalAmount: 6000000}
		err := validator.Validate(order)
		assert.Equal(t, errors.ErrOrderAboveMaxAmount, err)
	})
}

// ═══════════════════════════════════════════════════════════════
// TESTS DE VALIDADORES DE PAGO
// ═══════════════════════════════════════════════════════════════

func TestPaymentStatusValidator(t *testing.T) {
	validator := &PaymentStatusValidator{RequiredStatus: "paid"}

	t.Run("Orden pagada - debe pasar", func(t *testing.T) {
		order := &ports.OrderData{IsPaid: true}
		err := validator.Validate(order)
		assert.Nil(t, err)
	})

	t.Run("Orden no pagada - debe fallar", func(t *testing.T) {
		order := &ports.OrderData{IsPaid: false}
		err := validator.Validate(order)
		assert.Equal(t, errors.ErrOrderNotPaid, err)
	})
}

func TestPaymentMethodsValidator(t *testing.T) {
	validator := &PaymentMethodsValidator{AllowedMethods: []uint{1, 3, 5}}

	t.Run("Método permitido - debe pasar", func(t *testing.T) {
		order := &ports.OrderData{PaymentMethodID: 3}
		err := validator.Validate(order)
		assert.Nil(t, err)
	})

	t.Run("Método no permitido - debe fallar", func(t *testing.T) {
		order := &ports.OrderData{PaymentMethodID: 2}
		err := validator.Validate(order)
		assert.Equal(t, errors.ErrPaymentMethodNotAllowed, err)
	})

	t.Run("Sin restricciones (lista vacía) - siempre pasa", func(t *testing.T) {
		validatorSinRestricciones := &PaymentMethodsValidator{AllowedMethods: []uint{}}
		order := &ports.OrderData{PaymentMethodID: 999}
		err := validatorSinRestricciones.Validate(order)
		assert.Nil(t, err)
	})
}

// ═══════════════════════════════════════════════════════════════
// TESTS DE VALIDADORES DE ORDEN
// ═══════════════════════════════════════════════════════════════

func TestOrderTypesValidator(t *testing.T) {
	validator := &OrderTypesValidator{AllowedTypes: []string{"delivery", "pickup"}}

	t.Run("Tipo permitido - debe pasar", func(t *testing.T) {
		order := &ports.OrderData{OrderTypeName: "delivery"}
		err := validator.Validate(order)
		assert.Nil(t, err)
	})

	t.Run("Tipo no permitido - debe fallar", func(t *testing.T) {
		order := &ports.OrderData{OrderTypeName: "dine_in"}
		err := validator.Validate(order)
		assert.Equal(t, errors.ErrOrderTypeNotAllowed, err)
	})

	t.Run("Sin restricciones - siempre pasa", func(t *testing.T) {
		validatorSinRestricciones := &OrderTypesValidator{AllowedTypes: []string{}}
		order := &ports.OrderData{OrderTypeName: "cualquier_tipo"}
		err := validatorSinRestricciones.Validate(order)
		assert.Nil(t, err)
	})
}

func TestExcludeStatusesValidator(t *testing.T) {
	validator := &ExcludeStatusesValidator{ExcludedStatuses: []string{"cancelled", "refunded"}}

	t.Run("Estado permitido - debe pasar", func(t *testing.T) {
		order := &ports.OrderData{Status: "confirmed"}
		err := validator.Validate(order)
		assert.Nil(t, err)
	})

	t.Run("Estado excluido - debe fallar", func(t *testing.T) {
		order := &ports.OrderData{Status: "cancelled"}
		err := validator.Validate(order)
		assert.Equal(t, errors.ErrOrderStatusExcluded, err)
	})
}

// ═══════════════════════════════════════════════════════════════
// TESTS DE VALIDADORES DE PRODUCTOS
// ═══════════════════════════════════════════════════════════════

func TestExcludeProductsValidator(t *testing.T) {
	validator := &ExcludeProductsValidator{ExcludedSKUs: []string{"GIFT-CARD-001", "SKU-123"}}

	t.Run("Sin productos excluidos - debe pasar", func(t *testing.T) {
		order := &ports.OrderData{
			Items: []ports.OrderItemData{
				{SKU: "PROD-A"},
				{SKU: "PROD-B"},
			},
		}
		err := validator.Validate(order)
		assert.Nil(t, err)
	})

	t.Run("Con producto excluido - debe fallar", func(t *testing.T) {
		order := &ports.OrderData{
			Items: []ports.OrderItemData{
				{SKU: "PROD-A"},
				{SKU: "GIFT-CARD-001"},
			},
		}
		err := validator.Validate(order)
		assert.Equal(t, errors.ErrProductExcluded, err)
	})
}

func TestIncludeProductsOnlyValidator(t *testing.T) {
	validator := &IncludeProductsOnlyValidator{AllowedSKUs: []string{"PROD-A", "PROD-B"}}

	t.Run("Solo productos permitidos - debe pasar", func(t *testing.T) {
		order := &ports.OrderData{
			Items: []ports.OrderItemData{
				{SKU: "PROD-A"},
				{SKU: "PROD-B"},
			},
		}
		err := validator.Validate(order)
		assert.Nil(t, err)
	})

	t.Run("Productos fuera de la lista - debe fallar", func(t *testing.T) {
		order := &ports.OrderData{
			Items: []ports.OrderItemData{
				{SKU: "PROD-A"},
				{SKU: "PROD-C"}, // No está en AllowedSKUs
			},
		}
		err := validator.Validate(order)
		assert.Equal(t, errors.ErrProductNotAllowed, err)
	})

	t.Run("Sin restricciones - siempre pasa", func(t *testing.T) {
		validatorSinRestricciones := &IncludeProductsOnlyValidator{AllowedSKUs: []string{}}
		order := &ports.OrderData{
			Items: []ports.OrderItemData{
				{SKU: "CUALQUIER-SKU"},
			},
		}
		err := validatorSinRestricciones.Validate(order)
		assert.Nil(t, err)
	})
}

func TestItemsCountValidator(t *testing.T) {
	minCount := 2
	maxCount := 10

	t.Run("Dentro del rango - debe pasar", func(t *testing.T) {
		validator := &ItemsCountValidator{MinCount: &minCount, MaxCount: &maxCount}
		order := &ports.OrderData{
			Items: []ports.OrderItemData{
				{SKU: "A"}, {SKU: "B"}, {SKU: "C"},
			},
		}
		err := validator.Validate(order)
		assert.Nil(t, err)
	})

	t.Run("Por debajo del mínimo - debe fallar", func(t *testing.T) {
		validator := &ItemsCountValidator{MinCount: &minCount}
		order := &ports.OrderData{
			Items: []ports.OrderItemData{{SKU: "A"}},
		}
		err := validator.Validate(order)
		assert.Equal(t, errors.ErrMinItemsNotMet, err)
	})

	t.Run("Por encima del máximo - debe fallar", func(t *testing.T) {
		validator := &ItemsCountValidator{MaxCount: &maxCount}
		order := &ports.OrderData{
			Items: make([]ports.OrderItemData, 15),
		}
		err := validator.Validate(order)
		assert.Equal(t, errors.ErrMaxItemsExceeded, err)
	})

	t.Run("Sin restricciones - siempre pasa", func(t *testing.T) {
		validator := &ItemsCountValidator{}
		order := &ports.OrderData{
			Items: make([]ports.OrderItemData, 100),
		}
		err := validator.Validate(order)
		assert.Nil(t, err)
	})
}

// ═══════════════════════════════════════════════════════════════
// TESTS DE VALIDADORES DE CLIENTE
// ═══════════════════════════════════════════════════════════════

func TestCustomerTypesValidator(t *testing.T) {
	validator := &CustomerTypesValidator{AllowedTypes: []string{"natural", "juridica"}}

	t.Run("Tipo permitido - debe pasar", func(t *testing.T) {
		customerType := "natural"
		order := &ports.OrderData{CustomerType: &customerType}
		err := validator.Validate(order)
		assert.Nil(t, err)
	})

	t.Run("Tipo no permitido - debe fallar", func(t *testing.T) {
		customerType := "otro"
		order := &ports.OrderData{CustomerType: &customerType}
		err := validator.Validate(order)
		assert.Equal(t, errors.ErrCustomerTypeNotAllowed, err)
	})

	t.Run("CustomerType nil - siempre pasa", func(t *testing.T) {
		order := &ports.OrderData{CustomerType: nil}
		err := validator.Validate(order)
		assert.Nil(t, err)
	})
}

func TestExcludeCustomersValidator(t *testing.T) {
	validator := &ExcludeCustomersValidator{ExcludedCustomerIDs: []string{"123", "456"}}

	t.Run("Cliente no excluido - debe pasar", func(t *testing.T) {
		customerID := "789"
		order := &ports.OrderData{CustomerID: &customerID}
		err := validator.Validate(order)
		assert.Nil(t, err)
	})

	t.Run("Cliente excluido - debe fallar", func(t *testing.T) {
		customerID := "123"
		order := &ports.OrderData{CustomerID: &customerID}
		err := validator.Validate(order)
		assert.Equal(t, errors.ErrCustomerExcluded, err)
	})

	t.Run("CustomerID nil - siempre pasa", func(t *testing.T) {
		order := &ports.OrderData{CustomerID: nil}
		err := validator.Validate(order)
		assert.Nil(t, err)
	})
}

// ═══════════════════════════════════════════════════════════════
// TESTS DE VALIDADORES DE UBICACIÓN
// ═══════════════════════════════════════════════════════════════

func TestShippingRegionsValidator(t *testing.T) {
	validator := &ShippingRegionsValidator{AllowedRegions: []string{"Bogotá", "Medellín", "Cali"}}

	t.Run("Región permitida - debe pasar", func(t *testing.T) {
		shippingState := "Bogotá"
		order := &ports.OrderData{ShippingState: &shippingState}
		err := validator.Validate(order)
		assert.Nil(t, err)
	})

	t.Run("Región no permitida - debe fallar", func(t *testing.T) {
		shippingState := "Barranquilla"
		order := &ports.OrderData{ShippingState: &shippingState}
		err := validator.Validate(order)
		assert.Equal(t, errors.ErrShippingRegionNotAllowed, err)
	})

	t.Run("ShippingState nil - siempre pasa", func(t *testing.T) {
		order := &ports.OrderData{ShippingState: nil}
		err := validator.Validate(order)
		assert.Nil(t, err)
	})

	t.Run("Sin restricciones - siempre pasa", func(t *testing.T) {
		validatorSinRestricciones := &ShippingRegionsValidator{AllowedRegions: []string{}}
		shippingState := "Cualquier Ciudad"
		order := &ports.OrderData{ShippingState: &shippingState}
		err := validatorSinRestricciones.Validate(order)
		assert.Nil(t, err)
	})
}

// ═══════════════════════════════════════════════════════════════
// TESTS DE VALIDADORES DE FECHA
// ═══════════════════════════════════════════════════════════════

func TestDateRangeValidator(t *testing.T) {
	t.Run("Sin restricciones - siempre pasa", func(t *testing.T) {
		validator := &DateRangeValidator{}
		order := &ports.OrderData{}
		err := validator.Validate(order)
		assert.Nil(t, err)
	})

	// TODO: Implementar tests completos cuando DateRangeValidator esté implementado
}
