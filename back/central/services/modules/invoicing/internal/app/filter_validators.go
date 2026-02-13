package app

import (
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/errors"
)

// FilterValidator define la interfaz para validadores de filtros
type FilterValidator interface {
	Validate(order *dtos.OrderData) error
}

// ═══════════════════════════════════════════════════════════════
// VALIDADORES DE MONTO
// ═══════════════════════════════════════════════════════════════

type MinAmountValidator struct {
	MinAmount float64
}

func (v *MinAmountValidator) Validate(order *dtos.OrderData) error {
	if order.TotalAmount < v.MinAmount {
		return errors.ErrOrderBelowMinAmount
	}
	return nil
}

type MaxAmountValidator struct {
	MaxAmount float64
}

func (v *MaxAmountValidator) Validate(order *dtos.OrderData) error {
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

func (v *PaymentStatusValidator) Validate(order *dtos.OrderData) error {
	if v.RequiredStatus == "paid" && !order.IsPaid {
		return errors.ErrOrderNotPaid
	}
	return nil
}

type PaymentMethodsValidator struct {
	AllowedMethods []uint
}

func (v *PaymentMethodsValidator) Validate(order *dtos.OrderData) error {
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

func (v *OrderTypesValidator) Validate(order *dtos.OrderData) error {
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

func (v *ExcludeStatusesValidator) Validate(order *dtos.OrderData) error {
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

func (v *ExcludeProductsValidator) Validate(order *dtos.OrderData) error {
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

func (v *IncludeProductsOnlyValidator) Validate(order *dtos.OrderData) error {
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

func (v *ItemsCountValidator) Validate(order *dtos.OrderData) error {
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

func (v *CustomerTypesValidator) Validate(order *dtos.OrderData) error {
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

func (v *ExcludeCustomersValidator) Validate(order *dtos.OrderData) error {
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

func (v *ShippingRegionsValidator) Validate(order *dtos.OrderData) error {
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

func (v *DateRangeValidator) Validate(order *dtos.OrderData) error {
	// Si no hay restricciones de fecha, pasar
	if v.StartDate == nil && v.EndDate == nil {
		return nil
	}

	// TODO: Implementar validación de fechas cuando se necesite
	// Por ahora, esta validación siempre pasa
	return nil
}
