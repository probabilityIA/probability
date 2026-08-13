package repository

import (
	"context"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/products/internal/domain"
)

const channelsQuery = `
SELECT
    pbi.integration_id,
    COALESCE(i.name, '')             AS integration_name,
    COALESCE(it.code, '')            AS integration_code,
    COALESCE(i.category, '')         AS integration_category,
    COALESCE(it.image_url, '')       AS integration_image_url,
    COALESCE(i.is_active, false)     AS is_active,
    pbi.external_product_id,
    pbi.external_variant_id,
    pbi.external_sku,
    pbi.external_barcode,
    pbi.external_logistic_type       AS logistic_type,
    pbi.last_pushed_qty,
    pbi.created_at                   AS linked_at,
    pbi.updated_at,
    inv.started_at                   AS inventory_sync_at,
    COALESCE(inv.status, '')         AS inventory_sync_status,
    cmp.started_at                   AS compare_at,
    COALESCE(cmp.status, '')         AS compare_status
FROM product_business_integrations pbi
JOIN integrations i ON i.id = pbi.integration_id AND i.deleted_at IS NULL
LEFT JOIN integration_types it ON it.id = i.integration_type_id AND it.deleted_at IS NULL
LEFT JOIN LATERAL (
    SELECT r.started_at, r.status
    FROM integration_sync_runs r
    WHERE r.integration_id = pbi.integration_id
      AND r.business_id = pbi.business_id
      AND r.kind = 'inventory'
      AND r.deleted_at IS NULL
    ORDER BY r.started_at DESC
    LIMIT 1
) inv ON true
LEFT JOIN LATERAL (
    SELECT r.started_at, r.status
    FROM integration_sync_runs r
    WHERE r.integration_id = pbi.integration_id
      AND r.business_id = pbi.business_id
      AND r.kind = 'products'
      AND r.deleted_at IS NULL
    ORDER BY r.started_at DESC
    LIMIT 1
) cmp ON true
WHERE pbi.product_id = ? AND pbi.business_id = ? AND pbi.deleted_at IS NULL
ORDER BY i.name`

func (r *Repository) GetProductChannels(ctx context.Context, businessID uint, productID string) ([]domain.ProductChannel, error) {
	var rows []struct {
		IntegrationID       uint
		IntegrationName     string
		IntegrationCode     string
		IntegrationCategory string
		IntegrationImageURL string
		IsActive            bool
		ExternalProductID   string
		ExternalVariantID   *string
		ExternalSKU         *string
		ExternalBarcode     *string
		LogisticType        *string
		LastPushedQty       *int
		LinkedAt            time.Time
		UpdatedAt           time.Time
		InventorySyncAt     *time.Time
		InventorySyncStatus string
		CompareAt           *time.Time
		CompareStatus       string
	}

	if err := r.db.Conn(ctx).Raw(channelsQuery, productID, businessID).Scan(&rows).Error; err != nil {
		return nil, err
	}

	channels := make([]domain.ProductChannel, 0, len(rows))
	for _, row := range rows {
		channels = append(channels, domain.ProductChannel{
			IntegrationID:       row.IntegrationID,
			IntegrationName:     row.IntegrationName,
			IntegrationCode:     row.IntegrationCode,
			IntegrationCategory: row.IntegrationCategory,
			IntegrationImageURL: row.IntegrationImageURL,
			IsActive:            row.IsActive,
			ExternalProductID:   row.ExternalProductID,
			ExternalVariantID:   row.ExternalVariantID,
			ExternalSKU:         row.ExternalSKU,
			ExternalBarcode:     row.ExternalBarcode,
			LogisticType:        row.LogisticType,
			LastPushedQty:       row.LastPushedQty,
			LinkedAt:            row.LinkedAt,
			UpdatedAt:           row.UpdatedAt,
			InventorySyncAt:     row.InventorySyncAt,
			InventorySyncStatus: row.InventorySyncStatus,
			CompareAt:           row.CompareAt,
			CompareStatus:       row.CompareStatus,
		})
	}
	return channels, nil
}

const stockQuery = `
SELECT
    il.warehouse_id,
    COALESCE(w.name, '')   AS warehouse_name,
    COALESCE(w.code, '')   AS warehouse_code,
    SUM(il.quantity)       AS quantity,
    SUM(il.reserved_qty)   AS reserved_qty,
    SUM(il.available_qty)  AS available_qty,
    MAX(il.updated_at)     AS updated_at
FROM inventory_levels il
LEFT JOIN warehouses w ON w.id = il.warehouse_id AND w.deleted_at IS NULL
WHERE il.product_id = ? AND il.business_id = ? AND il.deleted_at IS NULL
GROUP BY il.warehouse_id, w.name, w.code
ORDER BY COALESCE(w.name, '')`

func (r *Repository) GetProductStock(ctx context.Context, businessID uint, productID string) ([]domain.ProductStockLevel, error) {
	var levels []domain.ProductStockLevel
	if err := r.db.Conn(ctx).Raw(stockQuery, productID, businessID).Scan(&levels).Error; err != nil {
		return nil, err
	}
	return levels, nil
}

const movementsQuery = `
SELECT
    sm.created_at,
    COALESCE(mt.name, '')      AS type_name,
    COALESCE(mt.direction, '') AS direction,
    COALESCE(sm.reason, '')    AS reason,
    sm.quantity,
    sm.previous_qty,
    sm.new_qty,
    COALESCE(w.name, '')       AS warehouse_name,
    COALESCE(i.name, '')       AS integration_name
FROM stock_movements sm
LEFT JOIN stock_movement_types mt ON mt.id = sm.movement_type_id AND mt.deleted_at IS NULL
LEFT JOIN warehouses w ON w.id = sm.warehouse_id AND w.deleted_at IS NULL
LEFT JOIN integrations i ON i.id = sm.integration_id AND i.deleted_at IS NULL
WHERE sm.product_id = ? AND sm.business_id = ? AND sm.deleted_at IS NULL
ORDER BY sm.created_at DESC
LIMIT ?`

func (r *Repository) GetProductMovements(ctx context.Context, businessID uint, productID string, limit int) ([]domain.ProductMovement, error) {
	if limit <= 0 || limit > 50 {
		limit = 10
	}
	var movements []domain.ProductMovement
	if err := r.db.Conn(ctx).Raw(movementsQuery, productID, businessID, limit).Scan(&movements).Error; err != nil {
		return nil, err
	}
	return movements, nil
}

const siblingsQuery = `
SELECT COUNT(*)
FROM products p
WHERE p.family_id = ? AND p.business_id = ? AND p.id <> ? AND p.deleted_at IS NULL`

func (r *Repository) CountFamilySiblings(ctx context.Context, businessID uint, familyID uint, productID string) (int64, error) {
	var total int64
	if err := r.db.Conn(ctx).Raw(siblingsQuery, familyID, businessID, productID).Scan(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func (r *Repository) GetProductFamilyByName(ctx context.Context, businessID uint, name string) (*domain.ProductFamily, error) {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return nil, nil
	}

	var fila struct {
		ID   uint
		Name string
	}
	err := r.db.Conn(ctx).Table("product_families").
		Select("id, name").
		Where("business_id = ? AND deleted_at IS NULL AND LOWER(TRIM(name)) = LOWER(?)", businessID, trimmed).
		Limit(1).Scan(&fila).Error
	if err != nil {
		return nil, err
	}
	if fila.ID == 0 {
		return nil, nil
	}
	return &domain.ProductFamily{ID: fila.ID, BusinessID: businessID, Name: fila.Name}, nil
}
