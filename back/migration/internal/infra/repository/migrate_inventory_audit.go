package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateInventoryAudit(ctx context.Context) error {
	db := r.db.Conn(ctx)
	if err := db.AutoMigrate(
		&models.CycleCountPlan{},
		&models.CycleCountTask{},
		&models.CycleCountLine{},
		&models.InventoryDiscrepancy{},
	); err != nil {
		return fmt.Errorf("failed to auto-migrate inventory audit tables: %w", err)
	}

	if err := r.seedCountAdjustmentMovementType(ctx); err != nil {
		return fmt.Errorf("failed to seed count_adjustment movement type: %w", err)
	}

	return nil
}

func (r *Repository) seedCountAdjustmentMovementType(ctx context.Context) error {
	db := r.db.Conn(ctx)
	if err := db.Exec(`SELECT setval(pg_get_serial_sequence('stock_movement_types', 'id'), GREATEST(COALESCE((SELECT MAX(id) FROM stock_movement_types), 1), 1))`).Error; err != nil {
		return err
	}
	var existing models.StockMovementType
	if err := db.Unscoped().Where("code = ?", "count_adjustment").First(&existing).Error; err == nil {
		return nil
	}
	sql := `INSERT INTO stock_movement_types (created_at, updated_at, code, name, description, is_active, direction)
		VALUES (NOW(), NOW(), 'count_adjustment', 'Ajuste por conteo ciclico', 'Movimiento generado por aprobacion de discrepancia en conteo ciclico', true, 'neutral')
		ON CONFLICT (code) DO NOTHING`
	return db.Exec(sql).Error
}
