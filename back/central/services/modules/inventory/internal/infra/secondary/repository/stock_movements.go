package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

func (r *Repository) CreateMovement(ctx context.Context, movement *entities.StockMovement) (*entities.StockMovement, error) {
	model := &models.StockMovement{
		ProductID:      movement.ProductID,
		WarehouseID:    movement.WarehouseID,
		LocationID:     movement.LocationID,
		BusinessID:     movement.BusinessID,
		MovementTypeID: movement.MovementTypeID,
		Reason:         movement.Reason,
		Quantity:       movement.Quantity,
		PreviousQty:    movement.PreviousQty,
		NewQty:         movement.NewQty,
		ReferenceType:  movement.ReferenceType,
		ReferenceID:    movement.ReferenceID,
		IntegrationID:  movement.IntegrationID,
		Notes:          movement.Notes,
		CreatedByID:    movement.CreatedByID,
	}

	if err := r.db.Conn(ctx).Create(model).Error; err != nil {
		return nil, err
	}

	movement.ID = model.ID
	movement.CreatedAt = model.CreatedAt
	return movement, nil
}

// createMovementTx crea un StockMovement dentro de una transacción activa
func (r *Repository) createMovementTx(tx *gorm.DB, movement *models.StockMovement) error {
	return tx.Create(movement).Error
}

func (r *Repository) ListMovements(ctx context.Context, params dtos.ListMovementsParams) ([]entities.StockMovement, int64, error) {
	var modelsList []models.StockMovement
	var total int64

	query := r.db.Conn(ctx).Model(&models.StockMovement{}).
		Where("stock_movements.business_id = ?", params.BusinessID)

	if params.ProductID != "" {
		query = query.Where("stock_movements.product_id = ?", params.ProductID)
	}
	if params.WarehouseID != nil {
		query = query.Where("stock_movements.warehouse_id = ?", *params.WarehouseID)
	}
	if params.Type != "" {
		// Filtrar por code del tipo de movimiento
		query = query.Joins("INNER JOIN stock_movement_types smt ON smt.id = stock_movements.movement_type_id").
			Where("smt.code = ?", params.Type)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(params.Offset()).Limit(params.PageSize).Order("stock_movements.created_at DESC").Find(&modelsList).Error; err != nil {
		return nil, 0, err
	}

	movements := make([]entities.StockMovement, len(modelsList))
	for i, m := range modelsList {
		e := stockMovementModelToEntity(&m)

		// Enriquecer con nombre del producto y bodega
		var prod struct {
			Name string
			SKU  string
		}
		r.db.Conn(ctx).Table("products").Select("name, sku").Where("id = ? AND deleted_at IS NULL", m.ProductID).Scan(&prod)
		e.ProductName = prod.Name
		e.ProductSKU = prod.SKU

		var whName string
		r.db.Conn(ctx).Table("warehouses").Select("name").Where("id = ? AND deleted_at IS NULL", m.WarehouseID).Scan(&whName)
		e.WarehouseName = whName

		// Enriquecer con datos del tipo de movimiento
		var mt struct {
			Code string
			Name string
		}
		r.db.Conn(ctx).Table("stock_movement_types").Select("code, name").Where("id = ?", m.MovementTypeID).Scan(&mt)
		e.MovementTypeCode = mt.Code
		e.MovementTypeName = mt.Name

		movements[i] = *e
	}
	return movements, total, nil
}
