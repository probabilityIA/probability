package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/domain/entities"
)

func (uc *UseCase) ListCategories(ctx context.Context) ([]entities.AnnouncementCategory, error) {
	return uc.repo.ListCategories(ctx)
}
