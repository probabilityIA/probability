package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/dtos"
)

// ============================================
// PRODUCTOS (queries replicadas localmente — regla de aislamiento)
// Tabla consultada: products (gestionada por módulo products)
// Solo consultas SELECT de solo lectura
// ============================================

// ListProductsByBusinessID retorna productos del negocio para comparación con proveedor.
// Usa .Table("products") porque el módulo invoicing no importa modelos del módulo products.
func (r *Repository) ListProductsByBusinessID(ctx context.Context, businessID uint) ([]dtos.SystemProduct, error) {
	var rows []struct {
		ID    string  `gorm:"column:id"`
		SKU   string  `gorm:"column:sku"`
		Name  string  `gorm:"column:name"`
		Price float64 `gorm:"column:price"`
	}

	err := r.db.Conn(ctx).
		Table("products").
		Select("id, sku, name, price").
		Where("business_id = ? AND deleted_at IS NULL", businessID).
		Find(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("failed to list products by business_id: %w", err)
	}

	products := make([]dtos.SystemProduct, 0, len(rows))
	for _, row := range rows {
		products = append(products, dtos.SystemProduct{
			ID:    row.ID,
			SKU:   row.SKU,
			Name:  row.Name,
			Price: row.Price,
		})
	}

	return products, nil
}
