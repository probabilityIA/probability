package dtos

import "github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/entities"

type SaveLayoutDTO struct {
	WarehouseID  uint
	BusinessID   uint
	CanvasWidth  float64
	CanvasHeight float64
	GridSize     float64
	Scale        float64
	Nodes        []entities.LayoutNode
}
