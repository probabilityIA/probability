package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/entities"
)

func (uc *UseCase) ListAdminBusinesses(ctx context.Context, page, pageSize int, search, statusFilter string) ([]entities.AdminBusinessRow, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return uc.repo.ListBusinessesForAdmin(ctx, page, pageSize, search, statusFilter)
}

func (uc *UseCase) GetAdminKPIs(ctx context.Context) (entities.AdminKPIs, error) {
	return uc.repo.GetAdminKPIs(ctx)
}
