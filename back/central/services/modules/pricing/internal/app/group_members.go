package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/pricing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/pricing/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/pricing/internal/domain/errors"
)

func (uc *UseCase) ListGroupMembers(ctx context.Context, params dtos.ListGroupMembersParams) ([]entities.ClientSummary, int64, error) {
	return uc.repo.ListGroupMembers(ctx, params)
}

func (uc *UseCase) ListAvailableClients(ctx context.Context, params dtos.ListAvailableClientsParams) ([]entities.ClientSummary, int64, error) {
	return uc.repo.ListAvailableClients(ctx, params)
}

func (uc *UseCase) AddGroupMembers(ctx context.Context, dto dtos.AddGroupMembersDTO) error {
	if len(dto.ClientIDs) == 0 {
		return domainerrors.ErrNoClients
	}
	return uc.repo.AddGroupMembers(ctx, dto)
}

func (uc *UseCase) RemoveGroupMember(ctx context.Context, businessID, groupID, clientID uint) error {
	return uc.repo.RemoveGroupMember(ctx, businessID, groupID, clientID)
}
