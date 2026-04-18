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
	var existing models.StockMovementType
	if err := db.Where("code = ?", "count_adjustment").First(&existing).Error; err == nil {
		return nil
	}
	mt := models.StockMovementType{
		Code:        "count_adjustment",
		Name:        "Ajuste por conteo ciclico",
		Description: "Movimiento generado por aprobacion de discrepancia en conteo ciclico",
		Direction:   "neutral",
		IsActive:    true,
	}
	return db.Create(&mt).Error
}
