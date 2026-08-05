package mocks

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/domain/ports"
)

type RepositoryMock struct {
	CreateFn           func(ctx context.Context, dto dtos.CreateGeozoneDTO) (*entities.Geozone, error)
	GetByIDFn          func(ctx context.Context, id uint, includeGeom bool) (*entities.Geozone, error)
	GetByCodeFn        func(ctx context.Context, businessID uint, geozoneType, code string) (*entities.Geozone, error)
	ListFn             func(ctx context.Context, params dtos.ListGeozonesParams) ([]entities.Geozone, int64, error)
	LookupContainingFn func(ctx context.Context, params dtos.LookupParams) ([]entities.Geozone, error)
	DeleteFn           func(ctx context.Context, id uint) error
	GetForDisplayFn    func(ctx context.Context, params dtos.DisplayParams) ([]dtos.DisplayFeature, error)
	ResolveAncestorsFn func(ctx context.Context, lat, lng float64, businessID uint) (*entities.GeozoneAncestors, error)

	CreatedDTOs   []dtos.CreateGeozoneDTO
	DisplayParams []dtos.DisplayParams
	nextID        uint
}

var _ ports.IRepository = (*RepositoryMock)(nil)

func (m *RepositoryMock) Create(ctx context.Context, dto dtos.CreateGeozoneDTO) (*entities.Geozone, error) {
	m.CreatedDTOs = append(m.CreatedDTOs, dto)
	if m.CreateFn != nil {
		return m.CreateFn(ctx, dto)
	}
	m.nextID++
	return &entities.Geozone{ID: m.nextID, Name: dto.Name, Type: dto.Type}, nil
}

func (m *RepositoryMock) GetByID(ctx context.Context, id uint, includeGeom bool) (*entities.Geozone, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, id, includeGeom)
	}
	return &entities.Geozone{ID: id}, nil
}

func (m *RepositoryMock) GetByCode(ctx context.Context, businessID uint, geozoneType, code string) (*entities.Geozone, error) {
	if m.GetByCodeFn != nil {
		return m.GetByCodeFn(ctx, businessID, geozoneType, code)
	}
	return nil, nil
}

func (m *RepositoryMock) List(ctx context.Context, params dtos.ListGeozonesParams) ([]entities.Geozone, int64, error) {
	if m.ListFn != nil {
		return m.ListFn(ctx, params)
	}
	return nil, 0, nil
}

func (m *RepositoryMock) LookupContaining(ctx context.Context, params dtos.LookupParams) ([]entities.Geozone, error) {
	if m.LookupContainingFn != nil {
		return m.LookupContainingFn(ctx, params)
	}
	return nil, nil
}

func (m *RepositoryMock) Delete(ctx context.Context, id uint) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, id)
	}
	return nil
}

func (m *RepositoryMock) GetForDisplay(ctx context.Context, params dtos.DisplayParams) ([]dtos.DisplayFeature, error) {
	m.DisplayParams = append(m.DisplayParams, params)
	if m.GetForDisplayFn != nil {
		return m.GetForDisplayFn(ctx, params)
	}
	return nil, nil
}

func (m *RepositoryMock) ResolveAncestors(ctx context.Context, lat, lng float64, businessID uint) (*entities.GeozoneAncestors, error) {
	if m.ResolveAncestorsFn != nil {
		return m.ResolveAncestorsFn(ctx, lat, lng, businessID)
	}
	return nil, nil
}

type ResolverMock struct {
	ResolveFn func(ctx context.Context, lat, lng float64, businessID uint) (*entities.GeozoneAncestors, error)
	Calls     int
}

var _ ports.IResolver = (*ResolverMock)(nil)

func (m *ResolverMock) Resolve(ctx context.Context, lat, lng float64, businessID uint) (*entities.GeozoneAncestors, error) {
	m.Calls++
	if m.ResolveFn != nil {
		return m.ResolveFn(ctx, lat, lng, businessID)
	}
	return nil, nil
}

type ProbabilityRepositoryMock struct {
	AncestorsByOrderIDFn    func(ctx context.Context, orderID string, businessID uint) (*entities.GeozoneAncestors, error)
	ProbabilityByOrderFn    func(ctx context.Context, ancestors *entities.GeozoneAncestors) ([]dtos.ProbabilityResult, error)
	ProbabilityForCarrierFn func(ctx context.Context, ancestors *entities.GeozoneAncestors, carrierKey string) (*dtos.ProbabilityResult, error)
	RefreshAggregatesFn     func(ctx context.Context) error

	CarrierKeys []string
}

var _ ports.IProbabilityRepository = (*ProbabilityRepositoryMock)(nil)

func (m *ProbabilityRepositoryMock) AncestorsByOrderID(ctx context.Context, orderID string, businessID uint) (*entities.GeozoneAncestors, error) {
	if m.AncestorsByOrderIDFn != nil {
		return m.AncestorsByOrderIDFn(ctx, orderID, businessID)
	}
	return &entities.GeozoneAncestors{}, nil
}

func (m *ProbabilityRepositoryMock) ProbabilityByOrder(ctx context.Context, ancestors *entities.GeozoneAncestors) ([]dtos.ProbabilityResult, error) {
	if m.ProbabilityByOrderFn != nil {
		return m.ProbabilityByOrderFn(ctx, ancestors)
	}
	return nil, nil
}

func (m *ProbabilityRepositoryMock) ProbabilityForCarrier(ctx context.Context, ancestors *entities.GeozoneAncestors, carrierKey string) (*dtos.ProbabilityResult, error) {
	m.CarrierKeys = append(m.CarrierKeys, carrierKey)
	if m.ProbabilityForCarrierFn != nil {
		return m.ProbabilityForCarrierFn(ctx, ancestors, carrierKey)
	}
	return nil, nil
}

func (m *ProbabilityRepositoryMock) RefreshAggregates(ctx context.Context) error {
	if m.RefreshAggregatesFn != nil {
		return m.RefreshAggregatesFn(ctx)
	}
	return nil
}

type DisplayCacheMock struct {
	GetFn      func(ctx context.Context, key string) ([]byte, bool)
	SetFn      func(ctx context.Context, key string, value []byte) error
	FlushAllFn func(ctx context.Context) error

	GetKeys []string
	SetKeys []string
	Flushed int
}

var _ ports.IDisplayCache = (*DisplayCacheMock)(nil)

func (m *DisplayCacheMock) Get(ctx context.Context, key string) ([]byte, bool) {
	m.GetKeys = append(m.GetKeys, key)
	if m.GetFn != nil {
		return m.GetFn(ctx, key)
	}
	return nil, false
}

func (m *DisplayCacheMock) Set(ctx context.Context, key string, value []byte) error {
	m.SetKeys = append(m.SetKeys, key)
	if m.SetFn != nil {
		return m.SetFn(ctx, key, value)
	}
	return nil
}

func (m *DisplayCacheMock) FlushAll(ctx context.Context) error {
	m.Flushed++
	if m.FlushAllFn != nil {
		return m.FlushAllFn(ctx)
	}
	return nil
}

type ProbabilityCacheMock struct {
	GetByOrderFn      func(ctx context.Context, businessID uint, orderID string) ([]dtos.ProbabilityResult, bool)
	SetByOrderFn      func(ctx context.Context, businessID uint, orderID string, results []dtos.ProbabilityResult) error
	InvalidateOrderFn func(ctx context.Context, businessID uint, orderID string) error

	SetCalls int
}

var _ ports.IProbabilityCache = (*ProbabilityCacheMock)(nil)

func (m *ProbabilityCacheMock) GetByOrder(ctx context.Context, businessID uint, orderID string) ([]dtos.ProbabilityResult, bool) {
	if m.GetByOrderFn != nil {
		return m.GetByOrderFn(ctx, businessID, orderID)
	}
	return nil, false
}

func (m *ProbabilityCacheMock) SetByOrder(ctx context.Context, businessID uint, orderID string, results []dtos.ProbabilityResult) error {
	m.SetCalls++
	if m.SetByOrderFn != nil {
		return m.SetByOrderFn(ctx, businessID, orderID, results)
	}
	return nil
}

func (m *ProbabilityCacheMock) InvalidateOrder(ctx context.Context, businessID uint, orderID string) error {
	if m.InvalidateOrderFn != nil {
		return m.InvalidateOrderFn(ctx, businessID, orderID)
	}
	return nil
}
