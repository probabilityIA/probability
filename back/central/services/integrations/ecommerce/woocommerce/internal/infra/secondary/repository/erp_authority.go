package repository

import "context"

const erpFeedsInventoryQuery = `
SELECT COUNT(*)
FROM integrations i
JOIN integration_types it ON it.id = i.integration_type_id AND it.deleted_at IS NULL
WHERE i.business_id = ?
  AND i.deleted_at IS NULL
  AND i.is_active = true
  AND it.code = 'siigo'
  AND COALESCE((i.config->>'inventory_sync_enabled')::boolean, false) = true`

func (r *ProductRepository) ERPFeedsInventory(ctx context.Context, businessID uint) (bool, error) {
	var total int64
	if err := r.db.Conn(ctx).Raw(erpFeedsInventoryQuery, businessID).Scan(&total).Error; err != nil {
		return false, err
	}
	return total > 0, nil
}
