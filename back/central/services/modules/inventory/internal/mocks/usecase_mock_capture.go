package mocks

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/app/request"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/app/response"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
)

func (m *UseCaseMock) CreateLPN(ctx context.Context, dto request.CreateLPNDTO) (*entities.LicensePlate, error) {
	if m.CreateLPNFn != nil {
		return m.CreateLPNFn(ctx, dto)
	}
	return &entities.LicensePlate{}, nil
}

func (m *UseCaseMock) GetLPN(ctx context.Context, businessID, id uint) (*entities.LicensePlate, error) {
	if m.GetLPNFn != nil {
		return m.GetLPNFn(ctx, businessID, id)
	}
	return &entities.LicensePlate{ID: id, BusinessID: businessID}, nil
}

func (m *UseCaseMock) ListLPNs(ctx context.Context, params dtos.ListLPNParams) ([]entities.LicensePlate, int64, error) {
	if m.ListLPNsFn != nil {
		return m.ListLPNsFn(ctx, params)
	}
	return []entities.LicensePlate{}, 0, nil
}

func (m *UseCaseMock) UpdateLPN(ctx context.Context, dto request.UpdateLPNDTO) (*entities.LicensePlate, error) {
	if m.UpdateLPNFn != nil {
		return m.UpdateLPNFn(ctx, dto)
	}
	return &entities.LicensePlate{}, nil
}

func (m *UseCaseMock) DeleteLPN(ctx context.Context, businessID, id uint) error {
	if m.DeleteLPNFn != nil {
		return m.DeleteLPNFn(ctx, businessID, id)
	}
	return nil
}

func (m *UseCaseMock) AddToLPN(ctx context.Context, dto request.AddToLPNDTO) (*entities.LicensePlateLine, error) {
	if m.AddToLPNFn != nil {
		return m.AddToLPNFn(ctx, dto)
	}
	return &entities.LicensePlateLine{}, nil
}

func (m *UseCaseMock) MoveLPN(ctx context.Context, dto request.MoveLPNDTO) (*entities.LicensePlate, error) {
	if m.MoveLPNFn != nil {
		return m.MoveLPNFn(ctx, dto)
	}
	return &entities.LicensePlate{}, nil
}

func (m *UseCaseMock) DissolveLPN(ctx context.Context, dto request.DissolveLPNDTO) error {
	if m.DissolveLPNFn != nil {
		return m.DissolveLPNFn(ctx, dto)
	}
	return nil
}

func (m *UseCaseMock) MergeLPN(ctx context.Context, dto request.MergeLPNDTO) (*entities.LicensePlate, error) {
	if m.MergeLPNFn != nil {
		return m.MergeLPNFn(ctx, dto)
	}
	return &entities.LicensePlate{}, nil
}

func (m *UseCaseMock) Scan(ctx context.Context, dto request.ScanDTO) (*response.ScanResponse, error) {
	if m.ScanFn != nil {
		return m.ScanFn(ctx, dto)
	}
	return &response.ScanResponse{}, nil
}

func (m *UseCaseMock) InboundSync(ctx context.Context, dto request.InboundSyncDTO) (*response.InboundSyncResult, error) {
	if m.InboundSyncFn != nil {
		return m.InboundSyncFn(ctx, dto)
	}
	return &response.InboundSyncResult{}, nil
}

func (m *UseCaseMock) ListSyncLogs(ctx context.Context, params dtos.ListSyncLogsParams) ([]entities.InventorySyncLog, int64, error) {
	if m.ListSyncLogsFn != nil {
		return m.ListSyncLogsFn(ctx, params)
	}
	return []entities.InventorySyncLog{}, 0, nil
}
