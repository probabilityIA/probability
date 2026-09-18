package repository

import (
	"context"
	"strconv"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm/clause"
)

func (r *Repository) ResolveOrderWarehouse(ctx context.Context, businessID, integrationID uint) (*entities.WarehouseRef, error) {
	preferred := r.integrationSingleWarehouseID(ctx, integrationID)

	var row struct {
		ID   uint
		Name string
	}
	query := r.db.Conn(ctx).
		Model(&models.Warehouse{}).
		Select("id, name").
		Where("business_id = ? AND is_active = true", businessID)
	if preferred > 0 {
		query = query.Order(clause.OrderBy{Expression: clause.Expr{SQL: "(id = ?) DESC", Vars: []interface{}{preferred}}})
	}
	err := query.
		Order("is_default DESC").
		Order("id ASC").
		Limit(1).
		Scan(&row).Error
	if err != nil {
		return nil, err
	}
	if row.ID == 0 {
		return nil, nil
	}
	return &entities.WarehouseRef{ID: row.ID, Name: row.Name}, nil
}

func (r *Repository) integrationSingleWarehouseID(ctx context.Context, integrationID uint) uint {
	if integrationID == 0 {
		return 0
	}
	var result struct {
		Mode   string
		Single string
	}
	err := r.db.Conn(ctx).
		Model(&models.Integration{}).
		Select("COALESCE(config->>'inventory_warehouse_mode', '') AS mode, COALESCE(config->>'inventory_single_warehouse_id', '') AS single").
		Where("id = ?", integrationID).
		Limit(1).
		Scan(&result).Error
	if err != nil || result.Mode != "single" {
		return 0
	}
	id, err := strconv.ParseUint(result.Single, 10, 64)
	if err != nil {
		return 0
	}
	return uint(id)
}
