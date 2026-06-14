package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/errors"
)

func (u *UseCase) GetLayout(ctx context.Context, businessID, warehouseID uint) (*entities.WarehouseLayout, error) {
	exists, err := u.repo.WarehouseExists(ctx, businessID, warehouseID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, domainerrors.ErrWarehouseNotFound
	}
	layout, err := u.repo.GetLayout(ctx, businessID, warehouseID)
	if err != nil {
		return nil, err
	}
	if layout == nil {
		return &entities.WarehouseLayout{
			WarehouseID:  warehouseID,
			BusinessID:   businessID,
			CanvasWidth:  1200,
			CanvasHeight: 800,
			GridSize:     20,
			Scale:        40,
			Nodes:        []entities.LayoutNode{},
		}, nil
	}
	return layout, nil
}

func (u *UseCase) SaveLayout(ctx context.Context, dto dtos.SaveLayoutDTO) (*entities.WarehouseLayout, error) {
	exists, err := u.repo.WarehouseExists(ctx, dto.BusinessID, dto.WarehouseID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, domainerrors.ErrWarehouseNotFound
	}
	if dto.CanvasWidth <= 0 {
		dto.CanvasWidth = 1200
	}
	if dto.CanvasHeight <= 0 {
		dto.CanvasHeight = 800
	}
	if dto.GridSize <= 0 {
		dto.GridSize = 20
	}
	if dto.Scale <= 0 {
		dto.Scale = 40
	}
	return u.repo.UpsertLayout(ctx, dto)
}
