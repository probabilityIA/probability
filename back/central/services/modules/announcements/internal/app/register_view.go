package app

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/domain/entities"
)

func (uc *UseCase) RegisterView(ctx context.Context, dto dtos.RegisterViewDTO) error {
	view := &entities.AnnouncementView{
		AnnouncementID: dto.AnnouncementID,
		UserID:         dto.UserID,
		BusinessID:     dto.BusinessID,
		Action:         dto.Action,
		LinkID:         dto.LinkID,
		ViewedAt:       time.Now(),
	}

	return uc.repo.RegisterView(ctx, view)
}
