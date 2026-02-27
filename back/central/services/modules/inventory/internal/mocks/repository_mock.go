package mocks

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/ports"
)

// RepositoryMock implementa ports.IRepository para tests
type RepositoryMock struct {
	GetProductInventoryFn                              func(ctx context.Context, params dtos.GetProductInventoryParams) ([]entities.InventoryLevel, error)
	ListWarehouseInventoryFn                           func(ctx context.Context, params dtos.ListWarehouseInventoryParams) ([]entities.InventoryLevel, int64, error)
	GetOrCreateLevelFn                                 func(ctx context.Context, productID string, warehouseID uint, locationID *uint, businessID uint) (*entities.InventoryLevel, error)
	UpdateLevelFn                                      func(ctx context.Context, level *entities.InventoryLevel) error
	CreateMovementFn                                   func(ctx context.Context, movement *entities.StockMovement) (*entities.StockMovement, error)
	ListMovementsFn                                    func(ctx context.Context, params dtos.ListMovementsParams) ([]entities.StockMovement, int64, error)
	ListMovementTypesFn                                func(ctx context.Context, params dtos.ListStockMovementTypesParams) ([]entities.StockMovementType, int64, error)
	GetMovementTypeByIDFn                              func(ctx context.Context, id uint) (*entities.StockMovementType, error)
	GetMovementTypeIDByCodeFn                          func(ctx context.Context, code string) (uint, error)
	CreateMovementTypeFn                               func(ctx context.Context, movType *entities.StockMovementType) (*entities.StockMovementType, error)
	UpdateMovementTypeFn                               func(ctx context.Context, movType *entities.StockMovementType) error
	DeleteMovementTypeFn                               func(ctx context.Context, id uint) error
	GetProductByIDFn                                   func(ctx context.Context, productID string, businessID uint) (string, string, bool, error)
	UpdateProductStockQuantityFn                       func(ctx context.Context, productID string, totalQuantity int) error
	WarehouseExistsFn                                  func(ctx context.Context, warehouseID uint, businessID uint) (bool, error)
	GetProductIntegrationsFn                           func(ctx context.Context, productID string, businessID uint) ([]ports.ProductIntegrationInfo, error)
	AdjustStockTxFn                                    func(ctx context.Context, params dtos.AdjustStockTxParams) (*dtos.AdjustStockTxResult, error)
	TransferStockTxFn                                  func(ctx context.Context, params dtos.TransferStockTxParams) (*dtos.TransferStockTxResult, error)
	ReserveStockTxFn                                   func(ctx context.Context, params dtos.ReserveStockTxParams) (*dtos.ReserveStockTxResult, error)
	ConfirmSaleTxFn                                    func(ctx context.Context, params dtos.ConfirmSaleTxParams) error
	ReleaseStockTxFn                                   func(ctx context.Context, params dtos.ReleaseTxParams) error
	ReturnStockTxFn                                    func(ctx context.Context, params dtos.ReturnStockTxParams) error
	GetDefaultWarehouseIDFn                            func(ctx context.Context, businessID uint) (uint, error)
}

func (m *RepositoryMock) GetProductInventory(ctx context.Context, params dtos.GetProductInventoryParams) ([]entities.InventoryLevel, error) {
	if m.GetProductInventoryFn != nil {
		return m.GetProductInventoryFn(ctx, params)
	}
	return nil, nil
}

func (m *RepositoryMock) ListWarehouseInventory(ctx context.Context, params dtos.ListWarehouseInventoryParams) ([]entities.InventoryLevel, int64, error) {
	if m.ListWarehouseInventoryFn != nil {
		return m.ListWarehouseInventoryFn(ctx, params)
	}
	return nil, 0, nil
}

func (m *RepositoryMock) GetOrCreateLevel(ctx context.Context, productID string, warehouseID uint, locationID *uint, businessID uint) (*entities.InventoryLevel, error) {
	if m.GetOrCreateLevelFn != nil {
		return m.GetOrCreateLevelFn(ctx, productID, warehouseID, locationID, businessID)
	}
	return nil, nil
}

func (m *RepositoryMock) UpdateLevel(ctx context.Context, level *entities.InventoryLevel) error {
	if m.UpdateLevelFn != nil {
		return m.UpdateLevelFn(ctx, level)
	}
	return nil
}

func (m *RepositoryMock) CreateMovement(ctx context.Context, movement *entities.StockMovement) (*entities.StockMovement, error) {
	if m.CreateMovementFn != nil {
		return m.CreateMovementFn(ctx, movement)
	}
	return movement, nil
}

func (m *RepositoryMock) ListMovements(ctx context.Context, params dtos.ListMovementsParams) ([]entities.StockMovement, int64, error) {
	if m.ListMovementsFn != nil {
		return m.ListMovementsFn(ctx, params)
	}
	return nil, 0, nil
}

func (m *RepositoryMock) ListMovementTypes(ctx context.Context, params dtos.ListStockMovementTypesParams) ([]entities.StockMovementType, int64, error) {
	if m.ListMovementTypesFn != nil {
		return m.ListMovementTypesFn(ctx, params)
	}
	return nil, 0, nil
}

func (m *RepositoryMock) GetMovementTypeByID(ctx context.Context, id uint) (*entities.StockMovementType, error) {
	if m.GetMovementTypeByIDFn != nil {
		return m.GetMovementTypeByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *RepositoryMock) GetMovementTypeIDByCode(ctx context.Context, code string) (uint, error) {
	if m.GetMovementTypeIDByCodeFn != nil {
		return m.GetMovementTypeIDByCodeFn(ctx, code)
	}
	return 1, nil
}

func (m *RepositoryMock) CreateMovementType(ctx context.Context, movType *entities.StockMovementType) (*entities.StockMovementType, error) {
	if m.CreateMovementTypeFn != nil {
		return m.CreateMovementTypeFn(ctx, movType)
	}
	return movType, nil
}

func (m *RepositoryMock) UpdateMovementType(ctx context.Context, movType *entities.StockMovementType) error {
	if m.UpdateMovementTypeFn != nil {
		return m.UpdateMovementTypeFn(ctx, movType)
	}
	return nil
}

func (m *RepositoryMock) DeleteMovementType(ctx context.Context, id uint) error {
	if m.DeleteMovementTypeFn != nil {
		return m.DeleteMovementTypeFn(ctx, id)
	}
	return nil
}

func (m *RepositoryMock) GetProductByID(ctx context.Context, productID string, businessID uint) (string, string, bool, error) {
	if m.GetProductByIDFn != nil {
		return m.GetProductByIDFn(ctx, productID, businessID)
	}
	return "Producto Test", "SKU-001", true, nil
}

func (m *RepositoryMock) UpdateProductStockQuantity(ctx context.Context, productID string, totalQuantity int) error {
	if m.UpdateProductStockQuantityFn != nil {
		return m.UpdateProductStockQuantityFn(ctx, productID, totalQuantity)
	}
	return nil
}

func (m *RepositoryMock) WarehouseExists(ctx context.Context, warehouseID uint, businessID uint) (bool, error) {
	if m.WarehouseExistsFn != nil {
		return m.WarehouseExistsFn(ctx, warehouseID, businessID)
	}
	return true, nil
}

func (m *RepositoryMock) GetProductIntegrations(ctx context.Context, productID string, businessID uint) ([]ports.ProductIntegrationInfo, error) {
	if m.GetProductIntegrationsFn != nil {
		return m.GetProductIntegrationsFn(ctx, productID, businessID)
	}
	return nil, nil
}

func (m *RepositoryMock) AdjustStockTx(ctx context.Context, params dtos.AdjustStockTxParams) (*dtos.AdjustStockTxResult, error) {
	if m.AdjustStockTxFn != nil {
		return m.AdjustStockTxFn(ctx, params)
	}
	return &dtos.AdjustStockTxResult{
		Movement:    &entities.StockMovement{ID: 1},
		NewQuantity: params.Quantity,
	}, nil
}

func (m *RepositoryMock) TransferStockTx(ctx context.Context, params dtos.TransferStockTxParams) (*dtos.TransferStockTxResult, error) {
	if m.TransferStockTxFn != nil {
		return m.TransferStockTxFn(ctx, params)
	}
	return &dtos.TransferStockTxResult{
		FromNewQty: 0,
		ToNewQty:   params.Quantity,
	}, nil
}

func (m *RepositoryMock) ReserveStockTx(ctx context.Context, params dtos.ReserveStockTxParams) (*dtos.ReserveStockTxResult, error) {
	if m.ReserveStockTxFn != nil {
		return m.ReserveStockTxFn(ctx, params)
	}
	return &dtos.ReserveStockTxResult{
		Reserved:   params.Quantity,
		Sufficient: true,
	}, nil
}

func (m *RepositoryMock) ConfirmSaleTx(ctx context.Context, params dtos.ConfirmSaleTxParams) error {
	if m.ConfirmSaleTxFn != nil {
		return m.ConfirmSaleTxFn(ctx, params)
	}
	return nil
}

func (m *RepositoryMock) ReleaseStockTx(ctx context.Context, params dtos.ReleaseTxParams) error {
	if m.ReleaseStockTxFn != nil {
		return m.ReleaseStockTxFn(ctx, params)
	}
	return nil
}

func (m *RepositoryMock) ReturnStockTx(ctx context.Context, params dtos.ReturnStockTxParams) error {
	if m.ReturnStockTxFn != nil {
		return m.ReturnStockTxFn(ctx, params)
	}
	return nil
}

func (m *RepositoryMock) GetDefaultWarehouseID(ctx context.Context, businessID uint) (uint, error) {
	if m.GetDefaultWarehouseIDFn != nil {
		return m.GetDefaultWarehouseIDFn(ctx, businessID)
	}
	return 1, nil
}
