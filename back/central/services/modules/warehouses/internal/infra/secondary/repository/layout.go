package repository

import (
	"context"
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func (r *Repository) GetLayout(ctx context.Context, businessID, warehouseID uint) (*entities.WarehouseLayout, error) {
	var model models.WarehouseLayout
	err := r.db.Conn(ctx).
		Where("warehouse_id = ? AND business_id = ?", warehouseID, businessID).
		First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return layoutModelToEntity(&model), nil
}

func (r *Repository) UpsertLayout(ctx context.Context, dto dtos.SaveLayoutDTO) (*entities.WarehouseLayout, error) {
	nodesJSON := layoutNodesToJSON(dto.Nodes)

	var model models.WarehouseLayout
	err := r.db.Conn(ctx).
		Where("warehouse_id = ? AND business_id = ?", dto.WarehouseID, dto.BusinessID).
		First(&model).Error

	if err == gorm.ErrRecordNotFound {
		model = models.WarehouseLayout{
			WarehouseID:  dto.WarehouseID,
			BusinessID:   dto.BusinessID,
			CanvasWidth:  dto.CanvasWidth,
			CanvasHeight: dto.CanvasHeight,
			GridSize:     dto.GridSize,
			Nodes:        nodesJSON,
		}
		if err := r.db.Conn(ctx).Create(&model).Error; err != nil {
			return nil, err
		}
		return layoutModelToEntity(&model), nil
	}
	if err != nil {
		return nil, err
	}

	updates := map[string]any{
		"canvas_width":  dto.CanvasWidth,
		"canvas_height": dto.CanvasHeight,
		"grid_size":     dto.GridSize,
		"nodes":         nodesJSON,
	}
	if err := r.db.Conn(ctx).Model(&models.WarehouseLayout{}).
		Where("id = ?", model.ID).
		Updates(updates).Error; err != nil {
		return nil, err
	}
	return r.GetLayout(ctx, dto.BusinessID, dto.WarehouseID)
}

func layoutNodesToJSON(nodes []entities.LayoutNode) datatypes.JSON {
	if nodes == nil {
		nodes = []entities.LayoutNode{}
	}
	b, err := json.Marshal(nodes)
	if err != nil {
		return datatypes.JSON([]byte("[]"))
	}
	return datatypes.JSON(b)
}

func layoutModelToEntity(m *models.WarehouseLayout) *entities.WarehouseLayout {
	e := &entities.WarehouseLayout{
		ID:           m.ID,
		WarehouseID:  m.WarehouseID,
		BusinessID:   m.BusinessID,
		CanvasWidth:  m.CanvasWidth,
		CanvasHeight: m.CanvasHeight,
		GridSize:     m.GridSize,
		Nodes:        []entities.LayoutNode{},
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
	if len(m.Nodes) > 0 {
		_ = json.Unmarshal(m.Nodes, &e.Nodes)
	}
	return e
}
