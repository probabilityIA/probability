package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateInventoryTraceability(ctx context.Context) error {
	db := r.db.Conn(ctx)

	if err := db.AutoMigrate(
		&models.InventoryState{},
		&models.UnitOfMeasure{},
		&models.ProductUoM{},
		&models.InventoryLot{},
		&models.InventorySerial{},
	); err != nil {
		return fmt.Errorf("failed to auto-migrate traceability tables: %w", err)
	}

	if err := db.AutoMigrate(&models.InventoryLevel{}, &models.StockMovement{}); err != nil {
		return fmt.Errorf("failed to auto-migrate inventory_levels/stock_movements with traceability: %w", err)
	}

	if err := db.Exec(`DROP INDEX IF EXISTS idx_inventory_product_warehouse`).Error; err != nil {
		return fmt.Errorf("failed to drop legacy unique index idx_inventory_product_warehouse: %w", err)
	}

	if err := r.seedInventoryStates(ctx); err != nil {
		return fmt.Errorf("failed to seed inventory states: %w", err)
	}

	if err := r.seedUnitsOfMeasure(ctx); err != nil {
		return fmt.Errorf("failed to seed units of measure: %w", err)
	}

	return nil
}

func (r *Repository) seedInventoryStates(ctx context.Context) error {
	db := r.db.Conn(ctx)
	states := []models.InventoryState{
		{Code: "available", Name: "Disponible", Description: "Stock disponible para venta", IsTerminal: false, IsActive: true},
		{Code: "reserved", Name: "Reservado", Description: "Stock reservado por orden pendiente", IsTerminal: false, IsActive: true},
		{Code: "on_hold", Name: "En espera", Description: "Stock retenido temporalmente", IsTerminal: false, IsActive: true},
		{Code: "damaged", Name: "Daniado", Description: "Stock daniado, no disponible para venta", IsTerminal: true, IsActive: true},
		{Code: "quarantine", Name: "Cuarentena", Description: "Stock en revision", IsTerminal: false, IsActive: true},
		{Code: "expired", Name: "Vencido", Description: "Lote vencido, no disponible", IsTerminal: true, IsActive: true},
		{Code: "in_transit", Name: "En transito", Description: "Stock en transferencia", IsTerminal: false, IsActive: true},
		{Code: "returned", Name: "Devuelto", Description: "Stock devuelto por cliente", IsTerminal: false, IsActive: true},
	}

	for _, s := range states {
		var existing models.InventoryState
		err := db.Where("code = ?", s.Code).First(&existing).Error
		if err == nil {
			continue
		}
		if err := db.Create(&s).Error; err != nil {
			return fmt.Errorf("failed to seed state %s: %w", s.Code, err)
		}
	}
	return nil
}

func (r *Repository) seedUnitsOfMeasure(ctx context.Context) error {
	db := r.db.Conn(ctx)
	units := []models.UnitOfMeasure{
		{Code: "UN", Name: "Unidad", Type: "count", IsActive: true},
		{Code: "KG", Name: "Kilogramo", Type: "weight", IsActive: true},
		{Code: "G", Name: "Gramo", Type: "weight", IsActive: true},
		{Code: "LB", Name: "Libra", Type: "weight", IsActive: true},
		{Code: "L", Name: "Litro", Type: "volume", IsActive: true},
		{Code: "ML", Name: "Mililitro", Type: "volume", IsActive: true},
		{Code: "CM", Name: "Centimetro", Type: "length", IsActive: true},
		{Code: "M", Name: "Metro", Type: "length", IsActive: true},
		{Code: "CAJA", Name: "Caja", Type: "count", IsActive: true},
		{Code: "PALETA", Name: "Paleta", Type: "count", IsActive: true},
	}

	for _, u := range units {
		var existing models.UnitOfMeasure
		err := db.Where("code = ?", u.Code).First(&existing).Error
		if err == nil {
			continue
		}
		if err := db.Create(&u).Error; err != nil {
			return fmt.Errorf("failed to seed uom %s: %w", u.Code, err)
		}
	}
	return nil
}
