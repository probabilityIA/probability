package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/entities"
	errs "github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/errors"
	"github.com/secamc93/probability/back/central/shared/moduleregistry"
)

func (uc *UseCase) GrantOverride(ctx context.Context, dto dtos.GrantOverrideDTO) error {
	if !moduleregistry.IsValid(dto.ModuleCode) {
		return errs.ErrInvalidModuleCode
	}

	override := &entities.BusinessModuleOverride{
		BusinessID:      dto.BusinessID,
		ModuleCode:      dto.ModuleCode,
		GrantedByUserID: dto.GrantedByUserID,
		Notes:           dto.Notes,
	}

	return uc.repo.CreateOverride(ctx, override)
}
