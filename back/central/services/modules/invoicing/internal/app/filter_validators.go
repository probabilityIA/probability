package app

import (
<<<<<<< HEAD
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/ports"
=======
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/errors"
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
)

// FilterValidator define la interfaz para validadores de filtros
type FilterValidator interface {
<<<<<<< HEAD
	Validate(order *ports.OrderData) error
=======
	Validate(order *dtos.OrderData) error
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
}

// ═══════════════════════════════════════════════════════════════
// VALIDADORES DE MONTO
// ═══════════════════════════════════════════════════════════════

type MinAmountValidator struct {
	MinAmount float64
}

<<<<<<< HEAD
func (v *MinAmountValidator) Validate(order *ports.OrderData) error {
=======
func (v *MinAmountValidator) Validate(order *dtos.OrderData) error {
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
	if order.TotalAmount < v.MinAmount {
		return errors.ErrOrderBelowMinAmount
	}
	return nil
}

type MaxAmountValidator struct {
	MaxAmount float64
}

<<<<<<< HEAD
func (v *MaxAmountValidator) Validate(order *ports.OrderData) error {
=======
func (v *MaxAmountValidator) Validate(order *dtos.OrderData) error {
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
	if order.TotalAmount > v.MaxAmount {
		return errors.ErrOrderAboveMaxAmount
	}
	return nil
}

// ═══════════════════════════════════════════════════════════════
// VALIDADORES DE PAGO
// ═══════════════════════════════════════════════════════════════

type PaymentStatusValidator struct {
	RequiredStatus string
}

<<<<<<< HEAD
func (v *PaymentStatusValidator) Validate(order *ports.OrderData) error {
=======
func (v *PaymentStatusValidator) Validate(order *dtos.OrderData) error {
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
	if v.RequiredStatus == "paid" && !order.IsPaid {
		return errors.ErrOrderNotPaid
	}
	return nil
}

type PaymentMethodsValidator struct {
	AllowedMethods []uint
}

<<<<<<< HEAD
func (v *PaymentMethodsValidator) Validate(order *ports.OrderData) error {
=======
func (v *PaymentMethodsValidator) Validate(order *dtos.OrderData) error {
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
	if len(v.AllowedMethods) == 0 {
		return nil // Sin restricción
	}

	for _, methodID := range v.AllowedMethods {
		if methodID == order.PaymentMethodID {
			return nil // Permitido
		}
	}
	return errors.ErrPaymentMethodNotAllowed
}

// ═══════════════════════════════════════════════════════════════
// VALIDADORES DE ORDEN
// ═══════════════════════════════════════════════════════════════

type OrderTypesValidator struct {
	AllowedTypes []string
}

<<<<<<< HEAD
func (v *OrderTypesValidator) Validate(order *ports.OrderData) error {
=======
func (v *OrderTypesValidator) Validate(order *dtos.OrderData) error {
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
	if len(v.AllowedTypes) == 0 {
		return nil
	}

	for _, allowedType := range v.AllowedTypes {
		if allowedType == order.OrderTypeName {
			return nil
		}
	}
	return errors.ErrOrderTypeNotAllowed
}

type ExcludeStatusesValidator struct {
	ExcludedStatuses []string
}

<<<<<<< HEAD
func (v *ExcludeStatusesValidator) Validate(order *ports.OrderData) error {
=======
func (v *ExcludeStatusesValidator) Validate(order *dtos.OrderData) error {
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
	for _, excludedStatus := range v.ExcludedStatuses {
		if excludedStatus == order.Status {
			return errors.ErrOrderStatusExcluded
		}
	}
	return nil
}

// ═══════════════════════════════════════════════════════════════
// VALIDADORES DE PRODUCTOS
// ═══════════════════════════════════════════════════════════════

type ExcludeProductsValidator struct {
	ExcludedSKUs []string
}

<<<<<<< HEAD
func (v *ExcludeProductsValidator) Validate(order *ports.OrderData) error {
=======
func (v *ExcludeProductsValidator) Validate(order *dtos.OrderData) error {
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
	for _, item := range order.Items {
		for _, excludedSKU := range v.ExcludedSKUs {
			if item.SKU == excludedSKU {
				return errors.ErrProductExcluded
			}
		}
	}
	return nil
}

type IncludeProductsOnlyValidator struct {
	AllowedSKUs []string
}

<<<<<<< HEAD
func (v *IncludeProductsOnlyValidator) Validate(order *ports.OrderData) error {
=======
func (v *IncludeProductsOnlyValidator) Validate(order *dtos.OrderData) error {
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
	if len(v.AllowedSKUs) == 0 {
		return nil
	}

	for _, item := range order.Items {
		found := false
		for _, allowedSKU := range v.AllowedSKUs {
			if item.SKU == allowedSKU {
				found = true
				break
			}
		}
		if !found {
			return errors.ErrProductNotAllowed
		}
	}
	return nil
}

type ItemsCountValidator struct {
	MinCount *int
	MaxCount *int
}

<<<<<<< HEAD
func (v *ItemsCountValidator) Validate(order *ports.OrderData) error {
=======
func (v *ItemsCountValidator) Validate(order *dtos.OrderData) error {
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
	itemCount := len(order.Items)

	if v.MinCount != nil && itemCount < *v.MinCount {
		return errors.ErrMinItemsNotMet
	}

	if v.MaxCount != nil && itemCount > *v.MaxCount {
		return errors.ErrMaxItemsExceeded
	}

	return nil
}

// ═══════════════════════════════════════════════════════════════
// VALIDADORES DE CLIENTE
// ═══════════════════════════════════════════════════════════════

type CustomerTypesValidator struct {
	AllowedTypes []string
}

<<<<<<< HEAD
func (v *CustomerTypesValidator) Validate(order *ports.OrderData) error {
=======
func (v *CustomerTypesValidator) Validate(order *dtos.OrderData) error {
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
	if len(v.AllowedTypes) == 0 || order.CustomerType == nil {
		return nil
	}

	for _, allowedType := range v.AllowedTypes {
		if allowedType == *order.CustomerType {
			return nil
		}
	}
	return errors.ErrCustomerTypeNotAllowed
}

type ExcludeCustomersValidator struct {
	ExcludedCustomerIDs []string
}

<<<<<<< HEAD
func (v *ExcludeCustomersValidator) Validate(order *ports.OrderData) error {
=======
func (v *ExcludeCustomersValidator) Validate(order *dtos.OrderData) error {
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
	if order.CustomerID == nil {
		return nil
	}

	for _, excludedID := range v.ExcludedCustomerIDs {
		if excludedID == *order.CustomerID {
			return errors.ErrCustomerExcluded
		}
	}
	return nil
}

// ═══════════════════════════════════════════════════════════════
// VALIDADORES DE UBICACIÓN
// ═══════════════════════════════════════════════════════════════

type ShippingRegionsValidator struct {
	AllowedRegions []string
}

<<<<<<< HEAD
func (v *ShippingRegionsValidator) Validate(order *ports.OrderData) error {
=======
func (v *ShippingRegionsValidator) Validate(order *dtos.OrderData) error {
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
	if len(v.AllowedRegions) == 0 || order.ShippingState == nil {
		return nil
	}

	for _, allowedRegion := range v.AllowedRegions {
		if allowedRegion == *order.ShippingState {
			return nil
		}
	}
	return errors.ErrShippingRegionNotAllowed
}

// ═══════════════════════════════════════════════════════════════
// VALIDADORES DE FECHA
// ═══════════════════════════════════════════════════════════════

type DateRangeValidator struct {
	StartDate *string
	EndDate   *string
}

<<<<<<< HEAD
func (v *DateRangeValidator) Validate(order *ports.OrderData) error {
=======
func (v *DateRangeValidator) Validate(order *dtos.OrderData) error {
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
	// Si no hay restricciones de fecha, pasar
	if v.StartDate == nil && v.EndDate == nil {
		return nil
	}

	// TODO: Implementar validación de fechas cuando se necesite
	// Por ahora, esta validación siempre pasa
	return nil
}
