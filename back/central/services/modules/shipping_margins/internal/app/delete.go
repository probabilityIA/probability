package app

import "context"

func (uc *UseCase) Delete(ctx context.Context, businessID, id uint) error {
	existing, err := uc.repo.GetByID(ctx, businessID, id)
	if err != nil {
		return err
	}
	if err := uc.repo.Delete(ctx, businessID, id); err != nil {
		return err
	}
	if uc.cache != nil {
		_ = uc.cache.Delete(ctx, existing.BusinessID, existing.CarrierCode)
	}
	return nil
}
