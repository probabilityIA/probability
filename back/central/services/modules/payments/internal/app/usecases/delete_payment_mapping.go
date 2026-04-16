package usecases

import "context"

// DeletePaymentMapping elimina un mapeo
func (uc *UseCase) DeletePaymentMapping(ctx context.Context, id uint) error {
	return uc.repo.DeletePaymentMapping(ctx, id)
}
