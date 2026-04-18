package mocks

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
)

func (m *RepositoryMock) CreateLPN(ctx context.Context, lpn *entities.LicensePlate) (*entities.LicensePlate, error) {
	if m.CreateLPNFn != nil {
		return m.CreateLPNFn(ctx, lpn)
	}
	return lpn, nil
}

func (m *RepositoryMock) GetLPNByID(ctx context.Context, businessID, id uint) (*entities.LicensePlate, error) {
	if m.GetLPNByIDFn != nil {
		return m.GetLPNByIDFn(ctx, businessID, id)
	}
	return &entities.LicensePlate{ID: id, BusinessID: businessID}, nil
}

func (m *RepositoryMock) GetLPNByCode(ctx context.Context, businessID uint, code string) (*entities.LicensePlate, error) {
	if m.GetLPNByCodeFn != nil {
		return m.GetLPNByCodeFn(ctx, businessID, code)
	}
	return &entities.LicensePlate{Code: code, BusinessID: businessID}, nil
}

func (m *RepositoryMock) ListLPNs(ctx context.Context, params dtos.ListLPNParams) ([]entities.LicensePlate, int64, error) {
	if m.ListLPNsFn != nil {
		return m.ListLPNsFn(ctx, params)
	}
	return []entities.LicensePlate{}, 0, nil
}

func (m *RepositoryMock) UpdateLPN(ctx context.Context, lpn *entities.LicensePlate) (*entities.LicensePlate, error) {
	if m.UpdateLPNFn != nil {
		return m.UpdateLPNFn(ctx, lpn)
	}
	return lpn, nil
}

func (m *RepositoryMock) DeleteLPN(ctx context.Context, businessID, id uint) error {
	if m.DeleteLPNFn != nil {
		return m.DeleteLPNFn(ctx, businessID, id)
	}
	return nil
}

func (m *RepositoryMock) LPNExistsByCode(ctx context.Context, businessID uint, code string, excludeID *uint) (bool, error) {
	if m.LPNExistsByCodeFn != nil {
		return m.LPNExistsByCodeFn(ctx, businessID, code, excludeID)
	}
	return false, nil
}

func (m *RepositoryMock) AddLPNLine(ctx context.Context, line *entities.LicensePlateLine) (*entities.LicensePlateLine, error) {
	if m.AddLPNLineFn != nil {
		return m.AddLPNLineFn(ctx, line)
	}
	return line, nil
}

func (m *RepositoryMock) ListLPNLines(ctx context.Context, lpnID uint) ([]entities.LicensePlateLine, error) {
	if m.ListLPNLinesFn != nil {
		return m.ListLPNLinesFn(ctx, lpnID)
	}
	return []entities.LicensePlateLine{}, nil
}

func (m *RepositoryMock) DissolveLPN(ctx context.Context, businessID, id uint) error {
	if m.DissolveLPNFn != nil {
		return m.DissolveLPNFn(ctx, businessID, id)
	}
	return nil
}

func (m *RepositoryMock) RecordScanEvent(ctx context.Context, event *entities.ScanEvent) (*entities.ScanEvent, error) {
	if m.RecordScanEventFn != nil {
		return m.RecordScanEventFn(ctx, event)
	}
	return event, nil
}

func (m *RepositoryMock) ResolveScanCode(ctx context.Context, businessID uint, code string) (*entities.ScanResolution, error) {
	if m.ResolveScanCodeFn != nil {
		return m.ResolveScanCodeFn(ctx, businessID, code)
	}
	return nil, nil
}

func (m *RepositoryMock) CreateSyncLog(ctx context.Context, log *entities.InventorySyncLog) (*entities.InventorySyncLog, error) {
	if m.CreateSyncLogFn != nil {
		return m.CreateSyncLogFn(ctx, log)
	}
	return log, nil
}

func (m *RepositoryMock) GetSyncLogByHash(ctx context.Context, businessID uint, direction, hash string) (*entities.InventorySyncLog, error) {
	if m.GetSyncLogByHashFn != nil {
		return m.GetSyncLogByHashFn(ctx, businessID, direction, hash)
	}
	return nil, nil
}

func (m *RepositoryMock) UpdateSyncLogStatus(ctx context.Context, id uint, status, errorMsg string) error {
	if m.UpdateSyncLogStatusFn != nil {
		return m.UpdateSyncLogStatusFn(ctx, id, status, errorMsg)
	}
	return nil
}

func (m *RepositoryMock) ListSyncLogs(ctx context.Context, params dtos.ListSyncLogsParams) ([]entities.InventorySyncLog, int64, error) {
	if m.ListSyncLogsFn != nil {
		return m.ListSyncLogsFn(ctx, params)
	}
	return []entities.InventorySyncLog{}, 0, nil
}
