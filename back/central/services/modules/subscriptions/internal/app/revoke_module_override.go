package app

import "context"

func (uc *UseCase) RevokeOverride(ctx context.Context, businessID uint, moduleCode string) error {
	return uc.repo.DeleteOverride(ctx, businessID, moduleCode)
}
