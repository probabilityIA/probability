package mocks

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/app/request"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/app/response"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
)

func (m *UseCaseMock) CreateLot(ctx context.Context, dto request.CreateLotDTO) (*entities.InventoryLot, error) {
	if m.CreateLotFn != nil {
		return m.CreateLotFn(ctx, dto)
	}
	return &entities.InventoryLot{}, nil
}

func (m *UseCaseMock) GetLot(ctx context.Context, businessID, lotID uint) (*entities.InventoryLot, error) {
	if m.GetLotFn != nil {
		return m.GetLotFn(ctx, businessID, lotID)
	}
	return &entities.InventoryLot{ID: lotID, BusinessID: businessID}, nil
}

func (m *UseCaseMock) ListLots(ctx context.Context, params dtos.ListLotsParams) ([]entities.InventoryLot, int64, error) {
	if m.ListLotsFn != nil {
		return m.ListLotsFn(ctx, params)
	}
	return []entities.InventoryLot{}, 0, nil
}

func (m *UseCaseMock) UpdateLot(ctx context.Context, dto request.UpdateLotDTO) (*entities.InventoryLot, error) {
	if m.UpdateLotFn != nil {
		return m.UpdateLotFn(ctx, dto)
	}
	return &entities.InventoryLot{}, nil
}

func (m *UseCaseMock) DeleteLot(ctx context.Context, businessID, lotID uint) error {
	if m.DeleteLotFn != nil {
		return m.DeleteLotFn(ctx, businessID, lotID)
	}
	return nil
}

func (m *UseCaseMock) CreateSerial(ctx context.Context, dto request.CreateSerialDTO) (*entities.InventorySerial, error) {
	if m.CreateSerialFn != nil {
		return m.CreateSerialFn(ctx, dto)
	}
	return &entities.InventorySerial{}, nil
}

func (m *UseCaseMock) GetSerial(ctx context.Context, businessID, serialID uint) (*entities.InventorySerial, error) {
	if m.GetSerialFn != nil {
		return m.GetSerialFn(ctx, businessID, serialID)
	}
	return &entities.InventorySerial{ID: serialID, BusinessID: businessID}, nil
}

func (m *UseCaseMock) ListSerials(ctx context.Context, params dtos.ListSerialsParams) ([]entities.InventorySerial, int64, error) {
	if m.ListSerialsFn != nil {
		return m.ListSerialsFn(ctx, params)
	}
	return []entities.InventorySerial{}, 0, nil
}

func (m *UseCaseMock) UpdateSerial(ctx context.Context, dto request.UpdateSerialDTO) (*entities.InventorySerial, error) {
	if m.UpdateSerialFn != nil {
		return m.UpdateSerialFn(ctx, dto)
	}
	return &entities.InventorySerial{}, nil
}

func (m *UseCaseMock) ListInventoryStates(ctx context.Context) ([]entities.InventoryState, error) {
	if m.ListInventoryStatesFn != nil {
		return m.ListInventoryStatesFn(ctx)
	}
	return []entities.InventoryState{}, nil
}

func (m *UseCaseMock) ChangeInventoryState(ctx context.Context, dto request.ChangeInventoryStateDTO) (*entities.StockMovement, error) {
	if m.ChangeInventoryStateFn != nil {
		return m.ChangeInventoryStateFn(ctx, dto)
	}
	return &entities.StockMovement{}, nil
}

func (m *UseCaseMock) ListUoMs(ctx context.Context) ([]entities.UnitOfMeasure, error) {
	if m.ListUoMsFn != nil {
		return m.ListUoMsFn(ctx)
	}
	return []entities.UnitOfMeasure{}, nil
}

func (m *UseCaseMock) ListProductUoMs(ctx context.Context, businessID uint, productID string) ([]entities.ProductUoM, error) {
	if m.ListProductUoMsFn != nil {
		return m.ListProductUoMsFn(ctx, businessID, productID)
	}
	return []entities.ProductUoM{}, nil
}

func (m *UseCaseMock) CreateProductUoM(ctx context.Context, dto request.CreateProductUoMDTO) (*entities.ProductUoM, error) {
	if m.CreateProductUoMFn != nil {
		return m.CreateProductUoMFn(ctx, dto)
	}
	return &entities.ProductUoM{}, nil
}

func (m *UseCaseMock) DeleteProductUoM(ctx context.Context, businessID, id uint) error {
	if m.DeleteProductUoMFn != nil {
		return m.DeleteProductUoMFn(ctx, businessID, id)
	}
	return nil
}

func (m *UseCaseMock) ConvertUoM(ctx context.Context, dto request.ConvertUoMDTO) (*response.ConvertUoMResult, error) {
	if m.ConvertUoMFn != nil {
		return m.ConvertUoMFn(ctx, dto)
	}
	return &response.ConvertUoMResult{}, nil
}
