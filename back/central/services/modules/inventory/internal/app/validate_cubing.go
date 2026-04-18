package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/app/request"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/app/response"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/errors"
)

func (uc *useCase) ValidateCubing(ctx context.Context, dto request.ValidateCubingDTO) (*response.CubingCheckResult, error) {
	if dto.Quantity <= 0 {
		return nil, domainerrors.ErrInvalidQuantity
	}

	_, _, _, err := uc.repo.GetProductByID(ctx, dto.ProductID, dto.BusinessID)
	if err != nil {
		return nil, domainerrors.ErrProductNotFound
	}

	dims, err := uc.repo.GetProductDimensions(ctx, dto.ProductID, dto.BusinessID)
	if err != nil {
		return nil, err
	}

	capacity, err := uc.repo.GetLocationCapacity(ctx, dto.LocationID)
	if err != nil {
		return nil, err
	}

	weightKgPerUnit := convertWeightToKg(dims.Weight, dims.WeightU)
	volumeCm3PerUnit := convertVolumeToCm3(dims.Length, dims.Width, dims.Height, dims.DimU)

	qty := float64(dto.Quantity)
	weightNeeded := weightKgPerUnit * qty
	volumeNeeded := volumeCm3PerUnit * qty

	occupied, err := uc.repo.GetLocationOccupiedQty(ctx, dto.LocationID)
	if err != nil {
		return nil, err
	}

	result := &response.CubingCheckResult{
		Fits:            true,
		WeightNeededKg:  weightNeeded,
		VolumeNeededCm3: volumeNeeded,
		OccupiedQty:     occupied,
	}

	if capacity.MaxWeightKg != nil {
		result.WeightMaxKg = *capacity.MaxWeightKg
		if weightNeeded > *capacity.MaxWeightKg {
			result.Fits = false
			result.Reason = "weight exceeds location capacity"
		}
	}
	if capacity.MaxVolumeCm3 != nil {
		result.VolumeMaxCm3 = *capacity.MaxVolumeCm3
		if volumeNeeded > *capacity.MaxVolumeCm3 {
			result.Fits = false
			if result.Reason == "" {
				result.Reason = "volume exceeds location capacity"
			} else {
				result.Reason = "weight and volume exceed location capacity"
			}
		}
	}

	return result, nil
}

func convertWeightToKg(weight float64, unit string) float64 {
	switch unit {
	case "g":
		return weight / 1000
	case "lb":
		return weight * 0.453592
	case "oz":
		return weight * 0.0283495
	case "kg", "":
		return weight
	default:
		return weight
	}
}

func convertVolumeToCm3(length, width, height float64, unit string) float64 {
	factor := 1.0
	switch unit {
	case "mm":
		factor = 0.1
	case "m":
		factor = 100
	case "in":
		factor = 2.54
	case "cm", "":
		factor = 1
	}
	lc := length * factor
	wc := width * factor
	hc := height * factor
	return lc * wc * hc
}
