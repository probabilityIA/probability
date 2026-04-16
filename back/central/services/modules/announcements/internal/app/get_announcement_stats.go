package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/domain/entities"
)

func (uc *UseCase) GetAnnouncementStats(ctx context.Context, announcementID uint) (*entities.AnnouncementStats, error) {
	return uc.repo.GetStats(ctx, announcementID)
}
