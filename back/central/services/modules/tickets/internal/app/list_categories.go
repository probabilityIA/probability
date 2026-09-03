package app

import "context"

func (uc *UseCase) ListCategories(ctx context.Context) ([]string, error) {
	return uc.repo.ListCategories(ctx)
}
