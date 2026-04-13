package app

import (
	"context"

	domainerrors "github.com/secamc93/probability/back/central/services/modules/announcements/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/domain/entities"
)

func (uc *UseCase) ChangeStatus(ctx context.Context, dto dtos.ChangeStatusDTO) error {
	if !isValidStatus(dto.Status) {
		return domainerrors.ErrInvalidStatus
	}

	_, err := uc.repo.GetByID(ctx, dto.ID)
	if err != nil {
		return err
	}

	return uc.repo.ChangeStatus(ctx, dto.ID, dto.Status)
}

func isValidStatus(s entities.AnnouncementStatus) bool {
	switch s {
	case entities.StatusDraft, entities.StatusScheduled, entities.StatusActive, entities.StatusInactive:
		return true
	}
	return false
}
