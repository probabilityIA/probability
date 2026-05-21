package app

import (
	"context"
	"sort"

	"github.com/secamc93/probability/back/central/services/modules/codreport/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/codreport/internal/domain/entities"
)

func (uc *UseCase) CarrierConfigs(ctx context.Context, businessID uint) ([]entities.CarrierConfig, error) {
	configs, err := uc.repo.CarrierConfigs(ctx, businessID)
	if err != nil {
		return nil, err
	}
	existing := map[string]bool{}
	for i := range configs {
		existing[configs[i].CarrierName] = true
	}

	discovered, err := uc.repo.DiscoveredCarriers(ctx, businessID)
	if err == nil {
		for _, name := range discovered {
			if name == "" || existing[name] {
				continue
			}
			configs = append(configs, entities.CarrierConfig{
				BusinessID:         businessID,
				CarrierName:        name,
				DiscountPercentage: 0,
				IsActive:           true,
			})
		}
	}

	sort.Slice(configs, func(i, j int) bool {
		return configs[i].CarrierName < configs[j].CarrierName
	})
	return configs, nil
}

func (uc *UseCase) SaveCarrierConfig(ctx context.Context, d dtos.SaveCarrierConfigDTO) (*entities.CarrierConfig, error) {
	return uc.repo.SaveCarrierConfig(ctx, d)
}
