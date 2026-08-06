package mocks

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/shipping_margins/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/shipping_margins/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/shipping_margins/internal/domain/ports"
)

type RepositoryMock struct {
	CreateFn                  func(ctx context.Context, m *entities.ShippingMargin) (*entities.ShippingMargin, error)
	GetByIDFn                 func(ctx context.Context, businessID, id uint) (*entities.ShippingMargin, error)
	GetByBusinessAndCarrierFn func(ctx context.Context, businessID uint, carrierCode string) (*entities.ShippingMargin, error)
	ListFn                    func(ctx context.Context, params dtos.ListShippingMarginsParams) ([]entities.ShippingMargin, int64, error)
	UpdateFn                  func(ctx context.Context, m *entities.ShippingMargin) (*entities.ShippingMargin, error)
	ExistsByCarrierFn         func(ctx context.Context, businessID uint, carrierCode string, excludeID *uint) (bool, error)
	ProfitReportFn            func(ctx context.Context, params dtos.ProfitReportParams) (*dtos.ProfitReportResponse, error)
	ProfitReportDetailFn      func(ctx context.Context, params dtos.ProfitReportDetailParams) (*dtos.ProfitReportDetailResponse, error)

	CreatedCarriers []string
}

var _ ports.IRepository = (*RepositoryMock)(nil)

func (m *RepositoryMock) Create(ctx context.Context, sm *entities.ShippingMargin) (*entities.ShippingMargin, error) {
	m.CreatedCarriers = append(m.CreatedCarriers, sm.CarrierCode)
	if m.CreateFn != nil {
		return m.CreateFn(ctx, sm)
	}
	return sm, nil
}

func (m *RepositoryMock) GetByID(ctx context.Context, businessID, id uint) (*entities.ShippingMargin, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, businessID, id)
	}
	return &entities.ShippingMargin{ID: id, BusinessID: businessID}, nil
}

func (m *RepositoryMock) GetByBusinessAndCarrier(ctx context.Context, businessID uint, carrierCode string) (*entities.ShippingMargin, error) {
	if m.GetByBusinessAndCarrierFn != nil {
		return m.GetByBusinessAndCarrierFn(ctx, businessID, carrierCode)
	}
	return nil, nil
}

func (m *RepositoryMock) List(ctx context.Context, params dtos.ListShippingMarginsParams) ([]entities.ShippingMargin, int64, error) {
	if m.ListFn != nil {
		return m.ListFn(ctx, params)
	}
	return nil, 0, nil
}

func (m *RepositoryMock) Update(ctx context.Context, sm *entities.ShippingMargin) (*entities.ShippingMargin, error) {
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, sm)
	}
	return sm, nil
}

func (m *RepositoryMock) ExistsByCarrier(ctx context.Context, businessID uint, carrierCode string, excludeID *uint) (bool, error) {
	if m.ExistsByCarrierFn != nil {
		return m.ExistsByCarrierFn(ctx, businessID, carrierCode, excludeID)
	}
	return false, nil
}

func (m *RepositoryMock) ProfitReport(ctx context.Context, params dtos.ProfitReportParams) (*dtos.ProfitReportResponse, error) {
	if m.ProfitReportFn != nil {
		return m.ProfitReportFn(ctx, params)
	}
	return &dtos.ProfitReportResponse{}, nil
}

func (m *RepositoryMock) ProfitReportDetail(ctx context.Context, params dtos.ProfitReportDetailParams) (*dtos.ProfitReportDetailResponse, error) {
	if m.ProfitReportDetailFn != nil {
		return m.ProfitReportDetailFn(ctx, params)
	}
	return &dtos.ProfitReportDetailResponse{}, nil
}

type CacheCall struct {
	BusinessID  uint
	CarrierCode string
}

type CacheWriterMock struct {
	UpsertFn    func(ctx context.Context, m *entities.ShippingMargin) error
	DeleteFn    func(ctx context.Context, businessID uint, carrierCode string) error
	UpsertCalls []CacheCall
	DeleteCalls []CacheCall
}

var _ ports.ICacheWriter = (*CacheWriterMock)(nil)

func (m *CacheWriterMock) Upsert(ctx context.Context, sm *entities.ShippingMargin) error {
	m.UpsertCalls = append(m.UpsertCalls, CacheCall{BusinessID: sm.BusinessID, CarrierCode: sm.CarrierCode})
	if m.UpsertFn != nil {
		return m.UpsertFn(ctx, sm)
	}
	return nil
}

func (m *CacheWriterMock) Delete(ctx context.Context, businessID uint, carrierCode string) error {
	m.DeleteCalls = append(m.DeleteCalls, CacheCall{BusinessID: businessID, CarrierCode: carrierCode})
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, businessID, carrierCode)
	}
	return nil
}
