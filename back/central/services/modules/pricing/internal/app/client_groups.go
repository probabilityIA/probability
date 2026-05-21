package app

import (
	"context"
	"strings"

	"github.com/secamc93/probability/back/central/services/modules/pricing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/pricing/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/pricing/internal/domain/errors"
)

func (uc *UseCase) CreateClientGroup(ctx context.Context, dto dtos.SaveClientGroupDTO) (*entities.ClientGroup, error) {
	name := strings.TrimSpace(dto.Name)
	if name == "" {
		return nil, domainerrors.ErrGroupNameRequired
	}

	exists, err := uc.repo.GroupNameExists(ctx, dto.BusinessID, 0, name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domainerrors.ErrGroupNameDuplicate
	}

	return uc.repo.CreateClientGroup(ctx, &entities.ClientGroup{
		BusinessID:  dto.BusinessID,
		Name:        name,
		Description: strings.TrimSpace(dto.Description),
		IsActive:    dto.IsActive,
	})
}

func (uc *UseCase) UpdateClientGroup(ctx context.Context, dto dtos.SaveClientGroupDTO) (*entities.ClientGroup, error) {
	name := strings.TrimSpace(dto.Name)
	if name == "" {
		return nil, domainerrors.ErrGroupNameRequired
	}

	exists, err := uc.repo.GroupNameExists(ctx, dto.BusinessID, dto.ID, name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domainerrors.ErrGroupNameDuplicate
	}

	return uc.repo.UpdateClientGroup(ctx, &entities.ClientGroup{
		ID:          dto.ID,
		BusinessID:  dto.BusinessID,
		Name:        name,
		Description: strings.TrimSpace(dto.Description),
		IsActive:    dto.IsActive,
	})
}

func (uc *UseCase) GetClientGroup(ctx context.Context, businessID, groupID uint) (*entities.ClientGroup, error) {
	return uc.repo.GetClientGroup(ctx, businessID, groupID)
}

func (uc *UseCase) ListClientGroups(ctx context.Context, params dtos.ListClientGroupsParams) ([]entities.ClientGroup, int64, error) {
	return uc.repo.ListClientGroups(ctx, params)
}

func (uc *UseCase) DeleteClientGroup(ctx context.Context, businessID, groupID uint) error {
	return uc.repo.DeleteClientGroup(ctx, businessID, groupID)
}
