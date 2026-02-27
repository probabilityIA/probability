package mocks

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
)

// UseCaseMock implementa app.IUseCase para tests de handlers
type UseCaseMock struct {
	GetProductInventoryFn    func(ctx context.Context, params dtos.GetProductInventoryParams) ([]entities.InventoryLevel, error)
	ListWarehouseInventoryFn func(ctx context.Context, params dtos.ListWarehouseInventoryParams) ([]entities.InventoryLevel, int64, error)
	AdjustStockFn            func(ctx context.Context, dto dtos.AdjustStockDTO) (*entities.StockMovement, error)
	TransferStockFn          func(ctx context.Context, dto dtos.TransferStockDTO) error
	ListMovementsFn          func(ctx context.Context, params dtos.ListMovementsParams) ([]entities.StockMovement, int64, error)
	ListMovementTypesFn      func(ctx context.Context, params dtos.ListStockMovementTypesParams) ([]entities.StockMovementType, int64, error)
	CreateMovementTypeFn     func(ctx context.Context, dto dtos.CreateStockMovementTypeDTO) (*entities.StockMovementType, error)
	UpdateMovementTypeFn     func(ctx context.Context, dto dtos.UpdateStockMovementTypeDTO) (*entities.StockMovementType, error)
	DeleteMovementTypeFn     func(ctx context.Context, id uint) error
	ReserveStockForOrderFn   func(ctx context.Context, orderID string, businessID uint, warehouseID *uint, items []dtos.OrderInventoryItem) (*dtos.OrderStockResult, error)
	ConfirmSaleForOrderFn    func(ctx context.Context, orderID string, businessID uint, warehouseID *uint, items []dtos.OrderInventoryItem) (*dtos.OrderStockResult, error)
	ReleaseStockForOrderFn   func(ctx context.Context, orderID string, businessID uint, warehouseID *uint, items []dtos.OrderInventoryItem) (*dtos.OrderStockResult, error)
	ReturnStockForOrderFn    func(ctx context.Context, orderID string, businessID uint, warehouseID *uint, items []dtos.OrderInventoryItem) (*dtos.OrderStockResult, error)
}

func (m *UseCaseMock) GetProductInventory(ctx context.Context, params dtos.GetProductInventoryParams) ([]entities.InventoryLevel, error) {
	if m.GetProductInventoryFn != nil {
		return m.GetProductInventoryFn(ctx, params)
	}
	return nil, nil
}

func (m *UseCaseMock) ListWarehouseInventory(ctx context.Context, params dtos.ListWarehouseInventoryParams) ([]entities.InventoryLevel, int64, error) {
	if m.ListWarehouseInventoryFn != nil {
		return m.ListWarehouseInventoryFn(ctx, params)
	}
	return nil, 0, nil
}

func (m *UseCaseMock) AdjustStock(ctx context.Context, dto dtos.AdjustStockDTO) (*entities.StockMovement, error) {
	if m.AdjustStockFn != nil {
		return m.AdjustStockFn(ctx, dto)
	}
	return &entities.StockMovement{ID: 1}, nil
}

func (m *UseCaseMock) TransferStock(ctx context.Context, dto dtos.TransferStockDTO) error {
	if m.TransferStockFn != nil {
		return m.TransferStockFn(ctx, dto)
	}
	return nil
}

func (m *UseCaseMock) ListMovements(ctx context.Context, params dtos.ListMovementsParams) ([]entities.StockMovement, int64, error) {
	if m.ListMovementsFn != nil {
		return m.ListMovementsFn(ctx, params)
	}
	return nil, 0, nil
}

func (m *UseCaseMock) ListMovementTypes(ctx context.Context, params dtos.ListStockMovementTypesParams) ([]entities.StockMovementType, int64, error) {
	if m.ListMovementTypesFn != nil {
		return m.ListMovementTypesFn(ctx, params)
	}
	return nil, 0, nil
}

func (m *UseCaseMock) CreateMovementType(ctx context.Context, dto dtos.CreateStockMovementTypeDTO) (*entities.StockMovementType, error) {
	if m.CreateMovementTypeFn != nil {
		return m.CreateMovementTypeFn(ctx, dto)
	}
	return &entities.StockMovementType{ID: 1, Code: dto.Code, Name: dto.Name}, nil
}

func (m *UseCaseMock) UpdateMovementType(ctx context.Context, dto dtos.UpdateStockMovementTypeDTO) (*entities.StockMovementType, error) {
	if m.UpdateMovementTypeFn != nil {
		return m.UpdateMovementTypeFn(ctx, dto)
	}
	return &entities.StockMovementType{ID: dto.ID, Name: dto.Name}, nil
}

func (m *UseCaseMock) DeleteMovementType(ctx context.Context, id uint) error {
	if m.DeleteMovementTypeFn != nil {
		return m.DeleteMovementTypeFn(ctx, id)
	}
	return nil
}

func (m *UseCaseMock) ReserveStockForOrder(ctx context.Context, orderID string, businessID uint, warehouseID *uint, items []dtos.OrderInventoryItem) (*dtos.OrderStockResult, error) {
	if m.ReserveStockForOrderFn != nil {
		return m.ReserveStockForOrderFn(ctx, orderID, businessID, warehouseID, items)
	}
	return &dtos.OrderStockResult{OrderID: orderID, Success: true}, nil
}

func (m *UseCaseMock) ConfirmSaleForOrder(ctx context.Context, orderID string, businessID uint, warehouseID *uint, items []dtos.OrderInventoryItem) (*dtos.OrderStockResult, error) {
	if m.ConfirmSaleForOrderFn != nil {
		return m.ConfirmSaleForOrderFn(ctx, orderID, businessID, warehouseID, items)
	}
	return &dtos.OrderStockResult{OrderID: orderID, Success: true}, nil
}

func (m *UseCaseMock) ReleaseStockForOrder(ctx context.Context, orderID string, businessID uint, warehouseID *uint, items []dtos.OrderInventoryItem) (*dtos.OrderStockResult, error) {
	if m.ReleaseStockForOrderFn != nil {
		return m.ReleaseStockForOrderFn(ctx, orderID, businessID, warehouseID, items)
	}
	return &dtos.OrderStockResult{OrderID: orderID, Success: true}, nil
}

func (m *UseCaseMock) ReturnStockForOrder(ctx context.Context, orderID string, businessID uint, warehouseID *uint, items []dtos.OrderInventoryItem) (*dtos.OrderStockResult, error) {
	if m.ReturnStockForOrderFn != nil {
		return m.ReturnStockForOrderFn(ctx, orderID, businessID, warehouseID, items)
	}
	return &dtos.OrderStockResult{OrderID: orderID, Success: true}, nil
}
