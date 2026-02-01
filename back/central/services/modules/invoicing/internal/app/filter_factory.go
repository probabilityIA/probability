package app

import (
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
)

// CreateValidators crea todos los validadores según la configuración
func CreateValidators(config *entities.FilterConfig) []FilterValidator {
	var validators []FilterValidator

	// Monto
	if config.MinAmount != nil {
		validators = append(validators, &MinAmountValidator{MinAmount: *config.MinAmount})
	}
	if config.MaxAmount != nil {
		validators = append(validators, &MaxAmountValidator{MaxAmount: *config.MaxAmount})
	}

	// Pago
	if config.PaymentStatus != nil {
		validators = append(validators, &PaymentStatusValidator{RequiredStatus: *config.PaymentStatus})
	}
	if len(config.PaymentMethods) > 0 {
		validators = append(validators, &PaymentMethodsValidator{AllowedMethods: config.PaymentMethods})
	}

	// Orden
	if len(config.OrderTypes) > 0 {
		validators = append(validators, &OrderTypesValidator{AllowedTypes: config.OrderTypes})
	}
	if len(config.ExcludeStatuses) > 0 {
		validators = append(validators, &ExcludeStatusesValidator{ExcludedStatuses: config.ExcludeStatuses})
	}

	// Productos
	if len(config.ExcludeProducts) > 0 {
		validators = append(validators, &ExcludeProductsValidator{ExcludedSKUs: config.ExcludeProducts})
	}
	if len(config.IncludeProductsOnly) > 0 {
		validators = append(validators, &IncludeProductsOnlyValidator{AllowedSKUs: config.IncludeProductsOnly})
	}
	if config.MinItemsCount != nil || config.MaxItemsCount != nil {
		validators = append(validators, &ItemsCountValidator{
			MinCount: config.MinItemsCount,
			MaxCount: config.MaxItemsCount,
		})
	}

	// Cliente
	if len(config.CustomerTypes) > 0 {
		validators = append(validators, &CustomerTypesValidator{AllowedTypes: config.CustomerTypes})
	}
	if len(config.ExcludeCustomerIDs) > 0 {
		validators = append(validators, &ExcludeCustomersValidator{ExcludedCustomerIDs: config.ExcludeCustomerIDs})
	}

	// Ubicación
	if len(config.ShippingRegions) > 0 {
		validators = append(validators, &ShippingRegionsValidator{AllowedRegions: config.ShippingRegions})
	}

	// Fecha
	if config.DateRange != nil {
		var startDate, endDate *string
		if config.DateRange.StartDate != nil {
			start := config.DateRange.StartDate.Format("2006-01-02")
			startDate = &start
		}
		if config.DateRange.EndDate != nil {
			end := config.DateRange.EndDate.Format("2006-01-02")
			endDate = &end
		}
		validators = append(validators, &DateRangeValidator{
			StartDate: startDate,
			EndDate:   endDate,
		})
	}

	return validators
}
