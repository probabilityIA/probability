package mocks

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
)

func (m *RepositoryMock) CreateLot(ctx context.Context, lot *entities.InventoryLot) (*entities.InventoryLot, error) {
	if m.CreateLotFn != nil {
		return m.CreateLotFn(ctx, lot)
	}
	return lot, nil
}

func (m *RepositoryMock) GetLotByID(ctx context.Context, businessID, lotID uint) (*entities.InventoryLot, error) {
	if m.GetLotByIDFn != nil {
		return m.GetLotByIDFn(ctx, businessID, lotID)
	}
	return &entities.InventoryLot{ID: lotID, BusinessID: businessID}, nil
}

func (m *RepositoryMock) ListLots(ctx context.Context, params dtos.ListLotsParams) ([]entities.InventoryLot, int64, error) {
	if m.ListLotsFn != nil {
		return m.ListLotsFn(ctx, params)
	}
	return []entities.InventoryLot{}, 0, nil
}

func (m *RepositoryMock) UpdateLot(ctx context.Context, lot *entities.InventoryLot) (*entities.InventoryLot, error) {
	if m.UpdateLotFn != nil {
		return m.UpdateLotFn(ctx, lot)
	}
	return lot, nil
}

func (m *RepositoryMock) DeleteLot(ctx context.Context, businessID, lotID uint) error {
	if m.DeleteLotFn != nil {
		return m.DeleteLotFn(ctx, businessID, lotID)
	}
	return nil
}

func (m *RepositoryMock) LotExistsByCode(ctx context.Context, businessID uint, productID, code string, excludeID *uint) (bool, error) {
	if m.LotExistsByCodeFn != nil {
		return m.LotExistsByCodeFn(ctx, businessID, productID, code, excludeID)
	}
	return false, nil
}

func (m *RepositoryMock) CreateSerial(ctx context.Context, serial *entities.InventorySerial) (*entities.InventorySerial, error) {
	if m.CreateSerialFn != nil {
		return m.CreateSerialFn(ctx, serial)
	}
	return serial, nil
}

func (m *RepositoryMock) GetSerialByID(ctx context.Context, businessID, serialID uint) (*entities.InventorySerial, error) {
	if m.GetSerialByIDFn != nil {
		return m.GetSerialByIDFn(ctx, businessID, serialID)
	}
	return &entities.InventorySerial{ID: serialID, BusinessID: businessID}, nil
}

func (m *RepositoryMock) ListSerials(ctx context.Context, params dtos.ListSerialsParams) ([]entities.InventorySerial, int64, error) {
	if m.ListSerialsFn != nil {
		return m.ListSerialsFn(ctx, params)
	}
	return []entities.InventorySerial{}, 0, nil
}

func (m *RepositoryMock) UpdateSerial(ctx context.Context, serial *entities.InventorySerial) (*entities.InventorySerial, error) {
	if m.UpdateSerialFn != nil {
		return m.UpdateSerialFn(ctx, serial)
	}
	return serial, nil
}

func (m *RepositoryMock) SerialExists(ctx context.Context, businessID uint, productID, serial string, excludeID *uint) (bool, error) {
	if m.SerialExistsFn != nil {
		return m.SerialExistsFn(ctx, businessID, productID, serial, excludeID)
	}
	return false, nil
}

func (m *RepositoryMock) ListInventoryStates(ctx context.Context) ([]entities.InventoryState, error) {
	if m.ListInventoryStatesFn != nil {
		return m.ListInventoryStatesFn(ctx)
	}
	return []entities.InventoryState{}, nil
}

func (m *RepositoryMock) GetInventoryStateByCode(ctx context.Context, code string) (*entities.InventoryState, error) {
	if m.GetInventoryStateByCodeFn != nil {
		return m.GetInventoryStateByCodeFn(ctx, code)
	}
	return &entities.InventoryState{Code: code}, nil
}

func (m *RepositoryMock) ListUoMs(ctx context.Context) ([]entities.UnitOfMeasure, error) {
	if m.ListUoMsFn != nil {
		return m.ListUoMsFn(ctx)
	}
	return []entities.UnitOfMeasure{}, nil
}

func (m *RepositoryMock) GetUoMByCode(ctx context.Context, code string) (*entities.UnitOfMeasure, error) {
	if m.GetUoMByCodeFn != nil {
		return m.GetUoMByCodeFn(ctx, code)
	}
	return &entities.UnitOfMeasure{Code: code}, nil
}

func (m *RepositoryMock) GetUoMByID(ctx context.Context, uomID uint) (*entities.UnitOfMeasure, error) {
	if m.GetUoMByIDFn != nil {
		return m.GetUoMByIDFn(ctx, uomID)
	}
	return &entities.UnitOfMeasure{ID: uomID}, nil
}

func (m *RepositoryMock) CreateProductUoM(ctx context.Context, pu *entities.ProductUoM) (*entities.ProductUoM, error) {
	if m.CreateProductUoMFn != nil {
		return m.CreateProductUoMFn(ctx, pu)
	}
	return pu, nil
}

func (m *RepositoryMock) ListProductUoMs(ctx context.Context, params dtos.ListProductUoMParams) ([]entities.ProductUoM, error) {
	if m.ListProductUoMsFn != nil {
		return m.ListProductUoMsFn(ctx, params)
	}
	return []entities.ProductUoM{}, nil
}

func (m *RepositoryMock) DeleteProductUoM(ctx context.Context, businessID, id uint) error {
	if m.DeleteProductUoMFn != nil {
		return m.DeleteProductUoMFn(ctx, businessID, id)
	}
	return nil
}

func (m *RepositoryMock) GetBaseProductUoM(ctx context.Context, businessID uint, productID string) (*entities.ProductUoM, error) {
	if m.GetBaseProductUoMFn != nil {
		return m.GetBaseProductUoMFn(ctx, businessID, productID)
	}
	return &entities.ProductUoM{}, nil
}

func (m *RepositoryMock) ChangeStateTx(ctx context.Context, params dtos.ChangeInventoryStateTxParams) (*entities.StockMovement, error) {
	if m.ChangeStateTxFn != nil {
		return m.ChangeStateTxFn(ctx, params)
	}
	return &entities.StockMovement{}, nil
}

func (m *RepositoryMock) ListLotsForReserve(ctx context.Context, productID string, warehouseID, businessID uint, strategy string) ([]entities.InventoryLot, error) {
	if m.ListLotsForReserveFn != nil {
		return m.ListLotsForReserveFn(ctx, productID, warehouseID, businessID, strategy)
	}
	return []entities.InventoryLot{}, nil
}
