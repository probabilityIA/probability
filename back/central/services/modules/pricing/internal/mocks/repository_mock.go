package mocks

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/pricing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/pricing/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/pricing/internal/domain/ports"
)

type RepositoryMock struct {
	CreateClientGroupFn    func(ctx context.Context, group *entities.ClientGroup) (*entities.ClientGroup, error)
	UpdateClientGroupFn    func(ctx context.Context, group *entities.ClientGroup) (*entities.ClientGroup, error)
	GetClientGroupFn       func(ctx context.Context, businessID, groupID uint) (*entities.ClientGroup, error)
	ListClientGroupsFn     func(ctx context.Context, params dtos.ListClientGroupsParams) ([]entities.ClientGroup, int64, error)
	DeleteClientGroupFn    func(ctx context.Context, businessID, groupID uint) error
	GroupNameExistsFn      func(ctx context.Context, businessID, groupID uint, name string) (bool, error)
	ListGroupMembersFn     func(ctx context.Context, params dtos.ListGroupMembersParams) ([]entities.ClientSummary, int64, error)
	ListAvailableClientsFn func(ctx context.Context, params dtos.ListAvailableClientsParams) ([]entities.ClientSummary, int64, error)
	AddGroupMembersFn      func(ctx context.Context, dto dtos.AddGroupMembersDTO) error
	RemoveGroupMemberFn    func(ctx context.Context, businessID, groupID, clientID uint) error
	GetClientGroupIDFn     func(ctx context.Context, businessID, clientID uint) (*uint, error)
	ListCatalogPricesFn    func(ctx context.Context, params dtos.ListCatalogPricesParams) ([]entities.CatalogPriceRow, int64, error)
	SaveCatalogPricesFn    func(ctx context.Context, dto dtos.SaveCatalogPricesDTO) error
	GetEffectivePriceFn    func(ctx context.Context, params dtos.EffectivePriceParams) (*entities.EffectivePrice, error)
}

var _ ports.IRepository = (*RepositoryMock)(nil)

func (m *RepositoryMock) CreateClientGroup(ctx context.Context, group *entities.ClientGroup) (*entities.ClientGroup, error) {
	if m.CreateClientGroupFn != nil {
		return m.CreateClientGroupFn(ctx, group)
	}
	return group, nil
}

func (m *RepositoryMock) UpdateClientGroup(ctx context.Context, group *entities.ClientGroup) (*entities.ClientGroup, error) {
	if m.UpdateClientGroupFn != nil {
		return m.UpdateClientGroupFn(ctx, group)
	}
	return group, nil
}

func (m *RepositoryMock) GetClientGroup(ctx context.Context, businessID, groupID uint) (*entities.ClientGroup, error) {
	if m.GetClientGroupFn != nil {
		return m.GetClientGroupFn(ctx, businessID, groupID)
	}
	return nil, nil
}

func (m *RepositoryMock) ListClientGroups(ctx context.Context, params dtos.ListClientGroupsParams) ([]entities.ClientGroup, int64, error) {
	if m.ListClientGroupsFn != nil {
		return m.ListClientGroupsFn(ctx, params)
	}
	return nil, 0, nil
}

func (m *RepositoryMock) DeleteClientGroup(ctx context.Context, businessID, groupID uint) error {
	if m.DeleteClientGroupFn != nil {
		return m.DeleteClientGroupFn(ctx, businessID, groupID)
	}
	return nil
}

func (m *RepositoryMock) GroupNameExists(ctx context.Context, businessID, groupID uint, name string) (bool, error) {
	if m.GroupNameExistsFn != nil {
		return m.GroupNameExistsFn(ctx, businessID, groupID, name)
	}
	return false, nil
}

func (m *RepositoryMock) ListGroupMembers(ctx context.Context, params dtos.ListGroupMembersParams) ([]entities.ClientSummary, int64, error) {
	if m.ListGroupMembersFn != nil {
		return m.ListGroupMembersFn(ctx, params)
	}
	return nil, 0, nil
}

func (m *RepositoryMock) ListAvailableClients(ctx context.Context, params dtos.ListAvailableClientsParams) ([]entities.ClientSummary, int64, error) {
	if m.ListAvailableClientsFn != nil {
		return m.ListAvailableClientsFn(ctx, params)
	}
	return nil, 0, nil
}

func (m *RepositoryMock) AddGroupMembers(ctx context.Context, dto dtos.AddGroupMembersDTO) error {
	if m.AddGroupMembersFn != nil {
		return m.AddGroupMembersFn(ctx, dto)
	}
	return nil
}

func (m *RepositoryMock) RemoveGroupMember(ctx context.Context, businessID, groupID, clientID uint) error {
	if m.RemoveGroupMemberFn != nil {
		return m.RemoveGroupMemberFn(ctx, businessID, groupID, clientID)
	}
	return nil
}

func (m *RepositoryMock) GetClientGroupID(ctx context.Context, businessID, clientID uint) (*uint, error) {
	if m.GetClientGroupIDFn != nil {
		return m.GetClientGroupIDFn(ctx, businessID, clientID)
	}
	return nil, nil
}

func (m *RepositoryMock) ListCatalogPrices(ctx context.Context, params dtos.ListCatalogPricesParams) ([]entities.CatalogPriceRow, int64, error) {
	if m.ListCatalogPricesFn != nil {
		return m.ListCatalogPricesFn(ctx, params)
	}
	return nil, 0, nil
}

func (m *RepositoryMock) SaveCatalogPrices(ctx context.Context, dto dtos.SaveCatalogPricesDTO) error {
	if m.SaveCatalogPricesFn != nil {
		return m.SaveCatalogPricesFn(ctx, dto)
	}
	return nil
}

func (m *RepositoryMock) GetEffectivePrice(ctx context.Context, params dtos.EffectivePriceParams) (*entities.EffectivePrice, error) {
	if m.GetEffectivePriceFn != nil {
		return m.GetEffectivePriceFn(ctx, params)
	}
	return nil, nil
}
