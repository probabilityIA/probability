package ports

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/pricing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/pricing/internal/domain/entities"
)

type IRepository interface {
	CreateClientGroup(ctx context.Context, group *entities.ClientGroup) (*entities.ClientGroup, error)
	UpdateClientGroup(ctx context.Context, group *entities.ClientGroup) (*entities.ClientGroup, error)
	GetClientGroup(ctx context.Context, businessID, groupID uint) (*entities.ClientGroup, error)
	ListClientGroups(ctx context.Context, params dtos.ListClientGroupsParams) ([]entities.ClientGroup, int64, error)
	DeleteClientGroup(ctx context.Context, businessID, groupID uint) error
	GroupNameExists(ctx context.Context, businessID, groupID uint, name string) (bool, error)

	ListGroupMembers(ctx context.Context, params dtos.ListGroupMembersParams) ([]entities.ClientSummary, int64, error)
	ListAvailableClients(ctx context.Context, params dtos.ListAvailableClientsParams) ([]entities.ClientSummary, int64, error)
	AddGroupMembers(ctx context.Context, dto dtos.AddGroupMembersDTO) error
	RemoveGroupMember(ctx context.Context, businessID, groupID, clientID uint) error
	GetClientGroupID(ctx context.Context, businessID, clientID uint) (*uint, error)

	ListCatalogPrices(ctx context.Context, params dtos.ListCatalogPricesParams) ([]entities.CatalogPriceRow, int64, error)
	SaveCatalogPrices(ctx context.Context, dto dtos.SaveCatalogPricesDTO) error
	GetEffectivePrice(ctx context.Context, params dtos.EffectivePriceParams) (*entities.EffectivePrice, error)
}
