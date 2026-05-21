package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/pricing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/pricing/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/pricing/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

type IUseCase interface {
	CreateClientGroup(ctx context.Context, dto dtos.SaveClientGroupDTO) (*entities.ClientGroup, error)
	UpdateClientGroup(ctx context.Context, dto dtos.SaveClientGroupDTO) (*entities.ClientGroup, error)
	GetClientGroup(ctx context.Context, businessID, groupID uint) (*entities.ClientGroup, error)
	ListClientGroups(ctx context.Context, params dtos.ListClientGroupsParams) ([]entities.ClientGroup, int64, error)
	DeleteClientGroup(ctx context.Context, businessID, groupID uint) error

	ListGroupMembers(ctx context.Context, params dtos.ListGroupMembersParams) ([]entities.ClientSummary, int64, error)
	ListAvailableClients(ctx context.Context, params dtos.ListAvailableClientsParams) ([]entities.ClientSummary, int64, error)
	AddGroupMembers(ctx context.Context, dto dtos.AddGroupMembersDTO) error
	RemoveGroupMember(ctx context.Context, businessID, groupID, clientID uint) error

	ListCatalogPrices(ctx context.Context, params dtos.ListCatalogPricesParams) ([]entities.CatalogPriceRow, int64, error)
	SaveCatalogPrices(ctx context.Context, dto dtos.SaveCatalogPricesDTO) error
	GetEffectivePrice(ctx context.Context, params dtos.EffectivePriceParams) (*entities.EffectivePrice, error)
}

type UseCase struct {
	repo ports.IRepository
	log  log.ILogger
}

func New(repo ports.IRepository, logger log.ILogger) IUseCase {
	return &UseCase{repo: repo, log: logger}
}
