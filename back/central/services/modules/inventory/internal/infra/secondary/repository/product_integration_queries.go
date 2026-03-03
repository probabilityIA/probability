package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/ports"
	"gorm.io/gorm"
)

// ============================================
// MÉTODOS DE CONSULTA A TABLAS DE OTROS MÓDULOS
// (Replicados localmente - no compartir repos)
// ============================================

// GetProductByID obtiene datos básicos de un producto.
// Tabla consultada: products (gestionada por módulo products)
func (r *Repository) GetProductByID(ctx context.Context, productID string, businessID uint) (string, string, bool, error) {
	var result struct {
		Name           string
		SKU            string
		TrackInventory bool
	}

	err := r.db.Conn(ctx).
		Table("products").
		Select("name, sku, track_inventory").
		Where("id = ? AND business_id = ? AND deleted_at IS NULL", productID, businessID).
		Scan(&result).Error

	if err != nil {
		return "", "", false, err
	}
	if result.Name == "" {
		return "", "", false, gorm.ErrRecordNotFound
	}

	return result.Name, result.SKU, result.TrackInventory, nil
}

// UpdateProductStockQuantity actualiza el campo StockQuantity del producto con el total de inventario.
// Tabla consultada: products (gestionada por módulo products)
func (r *Repository) UpdateProductStockQuantity(ctx context.Context, productID string, totalQuantity int) error {
	return r.db.Conn(ctx).
		Table("products").
		Where("id = ? AND deleted_at IS NULL", productID).
		Update("stock_quantity", totalQuantity).Error
}

// WarehouseExists verifica si una bodega existe y pertenece al negocio.
// Tabla consultada: warehouses (gestionada por módulo warehouses)
func (r *Repository) WarehouseExists(ctx context.Context, warehouseID uint, businessID uint) (bool, error) {
	var count int64
	err := r.db.Conn(ctx).
		Table("warehouses").
		Where("id = ? AND business_id = ? AND deleted_at IS NULL", warehouseID, businessID).
		Count(&count).Error
	return count > 0, err
}

// GetDefaultWarehouseID obtiene la bodega por defecto del negocio.
// Busca primero una bodega marcada como is_default=true, si no hay, usa la primera activa.
// Tabla consultada: warehouses (gestionada por módulo warehouses)
func (r *Repository) GetDefaultWarehouseID(ctx context.Context, businessID uint) (uint, error) {
	var result struct {
		ID uint
	}

	// Intentar bodega por defecto
	err := r.db.Conn(ctx).
		Table("warehouses").
		Select("id").
		Where("business_id = ? AND is_default = true AND is_active = true AND deleted_at IS NULL", businessID).
		Limit(1).
		Scan(&result).Error

	if err != nil {
		return 0, err
	}
	if result.ID > 0 {
		return result.ID, nil
	}

	// Fallback: primera bodega activa
	err = r.db.Conn(ctx).
		Table("warehouses").
		Select("id").
		Where("business_id = ? AND is_active = true AND deleted_at IS NULL", businessID).
		Order("id ASC").
		Limit(1).
		Scan(&result).Error

	if err != nil {
		return 0, err
	}
	if result.ID > 0 {
		return result.ID, nil
	}

	return 0, gorm.ErrRecordNotFound
}

// GetProductIntegrations obtiene las integraciones vinculadas a un producto.
// Tabla consultada: product_business_integrations + integration_types (replicadas)
func (r *Repository) GetProductIntegrations(ctx context.Context, productID string, businessID uint) ([]ports.ProductIntegrationInfo, error) {
	var results []struct {
		IntegrationID       uint
		ExternalProductID   string
		IntegrationTypeCode string
	}

	err := r.db.Conn(ctx).
		Table("product_business_integrations pbi").
		Select("pbi.integration_id, pbi.external_product_id, it.code AS integration_type_code").
		Joins("INNER JOIN integrations i ON i.id = pbi.integration_id AND i.deleted_at IS NULL AND i.is_active = true").
		Joins("INNER JOIN integration_types it ON it.id = i.integration_type_id AND it.deleted_at IS NULL").
		Where("pbi.product_id = ? AND pbi.business_id = ? AND pbi.deleted_at IS NULL", productID, businessID).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	infos := make([]ports.ProductIntegrationInfo, len(results))
	for i, r := range results {
		infos[i] = ports.ProductIntegrationInfo{
			IntegrationID:       r.IntegrationID,
			ExternalProductID:   r.ExternalProductID,
			IntegrationTypeCode: r.IntegrationTypeCode,
		}
	}
	return infos, nil
}
