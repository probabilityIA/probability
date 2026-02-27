package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *Repository) GetProductInventory(ctx context.Context, params dtos.GetProductInventoryParams) ([]entities.InventoryLevel, error) {
	// Read-through: intentar cache primero
	if r.cache != nil {
		cached, _ := r.cache.GetProductLevels(ctx, params.ProductID, params.BusinessID)
		if cached != nil {
			return cached, nil
		}
	}

	var modelsList []models.InventoryLevel

	err := r.db.Conn(ctx).
		Where("product_id = ? AND business_id = ?", params.ProductID, params.BusinessID).
		Find(&modelsList).Error
	if err != nil {
		return nil, err
	}

	levels := make([]entities.InventoryLevel, len(modelsList))
	for i, m := range modelsList {
		e := inventoryLevelModelToEntity(&m)

		// Enriquecer con nombre de bodega
		var wh struct {
			Name string
			Code string
		}
		r.db.Conn(ctx).Table("warehouses").Select("name, code").Where("id = ? AND deleted_at IS NULL", m.WarehouseID).Scan(&wh)
		e.WarehouseName = wh.Name
		e.WarehouseCode = wh.Code

		levels[i] = *e
	}

	// Write-through: guardar en cache (background)
	if r.cache != nil && len(levels) > 0 {
		go r.cache.SetProductLevels(context.Background(), params.ProductID, params.BusinessID, levels)
	}

	return levels, nil
}

func (r *Repository) ListWarehouseInventory(ctx context.Context, params dtos.ListWarehouseInventoryParams) ([]entities.InventoryLevel, int64, error) {
	var modelsList []models.InventoryLevel
	var total int64

	query := r.db.Conn(ctx).Model(&models.InventoryLevel{}).
		Where("warehouse_id = ? AND business_id = ?", params.WarehouseID, params.BusinessID)

	if params.Search != "" {
		like := "%" + params.Search + "%"
		query = query.Where("product_id IN (SELECT id FROM products WHERE (name ILIKE ? OR sku ILIKE ?) AND deleted_at IS NULL)", like, like)
	}

	if params.LowStock {
		query = query.Where("min_stock IS NOT NULL AND quantity <= min_stock")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(params.Offset()).Limit(params.PageSize).Order("updated_at DESC").Find(&modelsList).Error; err != nil {
		return nil, 0, err
	}

	levels := make([]entities.InventoryLevel, len(modelsList))
	for i, m := range modelsList {
		e := inventoryLevelModelToEntity(&m)

		// Enriquecer con nombre del producto
		var prod struct {
			Name string
			SKU  string
		}
		r.db.Conn(ctx).Table("products").Select("name, sku").Where("id = ? AND deleted_at IS NULL", m.ProductID).Scan(&prod)
		e.ProductName = prod.Name
		e.ProductSKU = prod.SKU

		levels[i] = *e
	}
	return levels, total, nil
}

func (r *Repository) GetOrCreateLevel(ctx context.Context, productID string, warehouseID uint, locationID *uint, businessID uint) (*entities.InventoryLevel, error) {
	var model models.InventoryLevel

	query := r.db.Conn(ctx).Where("product_id = ? AND warehouse_id = ? AND business_id = ?", productID, warehouseID, businessID)
	if locationID != nil {
		query = query.Where("location_id = ?", *locationID)
	} else {
		query = query.Where("location_id IS NULL")
	}

	err := query.First(&model).Error
	if err == nil {
		return inventoryLevelModelToEntity(&model), nil
	}
	if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// Crear nuevo nivel
	model = models.InventoryLevel{
		ProductID:    productID,
		WarehouseID:  warehouseID,
		LocationID:   locationID,
		BusinessID:   businessID,
		Quantity:     0,
		ReservedQty:  0,
		AvailableQty: 0,
	}

	if err := r.db.Conn(ctx).Create(&model).Error; err != nil {
		return nil, err
	}

	return inventoryLevelModelToEntity(&model), nil
}

func (r *Repository) UpdateLevel(ctx context.Context, level *entities.InventoryLevel) error {
	err := r.db.Conn(ctx).Model(&models.InventoryLevel{}).
		Where("id = ?", level.ID).
		Updates(map[string]interface{}{
			"quantity":      level.Quantity,
			"reserved_qty":  level.ReservedQty,
			"available_qty": level.AvailableQty,
			"min_stock":     level.MinStock,
			"max_stock":     level.MaxStock,
			"reorder_point": level.ReorderPoint,
		}).Error
	if err != nil {
		return err
	}

	// Write-through: actualizar cache (background)
	if r.cache != nil {
		go r.cache.SetLevel(context.Background(), level.ProductID, level.WarehouseID, level)
		go r.cache.InvalidateProduct(context.Background(), level.ProductID, level.BusinessID)
	}

	return nil
}

// ============================================
// MÉTODOS TRANSACCIONALES INTERNOS
// Reciben *gorm.DB (la transacción activa)
// ============================================

// getOrCreateLevelTx obtiene o crea un InventoryLevel dentro de una transacción.
// Usa SELECT FOR UPDATE para bloqueo pesimista de la fila.
// Si la fila no existe, la crea; si falla por unique constraint (concurrencia),
// reintenta el SELECT FOR UPDATE.
func (r *Repository) getOrCreateLevelTx(tx *gorm.DB, productID string, warehouseID uint, locationID *uint, businessID uint) (*models.InventoryLevel, error) {
	var model models.InventoryLevel

	query := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("product_id = ? AND warehouse_id = ? AND business_id = ?", productID, warehouseID, businessID)

	if locationID != nil {
		query = query.Where("location_id = ?", *locationID)
	} else {
		query = query.Where("location_id IS NULL")
	}

	err := query.First(&model).Error
	if err == nil {
		return &model, nil // Encontrado y bloqueado
	}
	if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// No existe: crear
	model = models.InventoryLevel{
		ProductID:    productID,
		WarehouseID:  warehouseID,
		LocationID:   locationID,
		BusinessID:   businessID,
		Quantity:     0,
		ReservedQty:  0,
		AvailableQty: 0,
	}

	if createErr := tx.Create(&model).Error; createErr != nil {
		// Si falla por unique constraint (otra goroutine creó primero),
		// reintentar SELECT FOR UPDATE
		retryQuery := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("product_id = ? AND warehouse_id = ? AND business_id = ?", productID, warehouseID, businessID)
		if locationID != nil {
			retryQuery = retryQuery.Where("location_id = ?", *locationID)
		} else {
			retryQuery = retryQuery.Where("location_id IS NULL")
		}
		if retryErr := retryQuery.First(&model).Error; retryErr != nil {
			return nil, createErr // Retornar el error original del Create
		}
		return &model, nil
	}

	return &model, nil
}

// updateLevelTx actualiza un InventoryLevel dentro de una transacción activa
func (r *Repository) updateLevelTx(tx *gorm.DB, level *models.InventoryLevel) error {
	return tx.Model(level).
		Updates(map[string]interface{}{
			"quantity":      level.Quantity,
			"reserved_qty":  level.ReservedQty,
			"available_qty": level.AvailableQty,
		}).Error
}
