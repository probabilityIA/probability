package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/domain/entities"
)

func (uc *UseCase) GetAnnouncement(ctx context.Context, id uint) (*entities.Announcement, error) {
	return uc.repo.GetByID(ctx, id)
}
