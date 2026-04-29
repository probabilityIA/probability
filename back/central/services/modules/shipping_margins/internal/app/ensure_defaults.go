package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/shipping_margins/internal/domain/entities"
)

func (uc *UseCase) EnsureDefaultsForBusiness(ctx context.Context, businessID uint) error {
	if businessID == 0 {
		return nil
	}
	for _, c := range entities.DefaultCarriers {
		exists, err := uc.repo.ExistsByCarrier(ctx, businessID, c.Code, nil)
		if err != nil {
			return err
		}
		if exists {
			continue
		}
		m := &entities.ShippingMargin{
			BusinessID:      businessID,
			CarrierCode:     c.Code,
			CarrierName:     c.Name,
			MarginAmount:    0,
			InsuranceMargin: 0,
			IsActive:        true,
		}
		created, err := uc.repo.Create(ctx, m)
		if err != nil {
			return err
		}
		if uc.cache != nil {
			_ = uc.cache.Upsert(ctx, created)
		}
	}
	return nil
}
