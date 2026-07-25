package app

import (
	"context"
	"strings"

	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/errors"
)

func (uc *UseCase) ListServices(ctx context.Context) ([]entities.Service, error) {
	return uc.repo.ListServices(ctx)
}

func (uc *UseCase) validateService(dto *dtos.SaveServiceDTO) error {
	dto.Code = strings.ToUpper(strings.TrimSpace(dto.Code))
	dto.Name = strings.TrimSpace(dto.Name)
	dto.UnspscCode = strings.TrimSpace(dto.UnspscCode)
	if dto.Code == "" {
		return domainerrors.ErrCodeRequired
	}
	if dto.Name == "" {
		return domainerrors.ErrNameRequired
	}
	if dto.UnitPrice < 0 {
		return domainerrors.ErrInvalidAmount
	}
	return nil
}

func (uc *UseCase) CreateService(ctx context.Context, dto dtos.SaveServiceDTO) (*entities.Service, error) {
	if err := uc.validateService(&dto); err != nil {
		return nil, err
	}
	exists, err := uc.repo.ServiceCodeExists(ctx, dto.Code, nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domainerrors.ErrDuplicateCode
	}
	dto.IsActive = true
	return uc.repo.CreateService(ctx, dto)
}

func (uc *UseCase) UpdateService(ctx context.Context, dto dtos.SaveServiceDTO) (*entities.Service, error) {
	if err := uc.validateService(&dto); err != nil {
		return nil, err
	}
	if _, err := uc.repo.GetServiceByID(ctx, dto.ID); err != nil {
		return nil, err
	}
	exists, err := uc.repo.ServiceCodeExists(ctx, dto.Code, &dto.ID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domainerrors.ErrDuplicateCode
	}
	return uc.repo.UpdateService(ctx, dto)
}

func (uc *UseCase) DeleteService(ctx context.Context, id uint) error {
	if _, err := uc.repo.GetServiceByID(ctx, id); err != nil {
		return err
	}
	return uc.repo.DeleteService(ctx, id)
}
