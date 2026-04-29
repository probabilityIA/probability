package repository

import (
	"context"
	"fmt"
	"strings"
)

func quote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}

func (r *Repository) migrateSeedDefaultShippingMargins(ctx context.Context) error {
	db := r.db.Conn(ctx)

	carriers := []struct {
		Code string
		Name string
	}{
		{"servientrega", "Servientrega"},
		{"interrapidisimo", "Interrapidisimo"},
		{"coordinadora", "Coordinadora"},
		{"envia", "Envia"},
		{"tcc", "TCC"},
		{"deprisa", "Deprisa"},
		{"99minutos", "99Minutos"},
		{"mipaquete", "MiPaquete"},
		{"enviame", "Enviame"},
	}

	for _, c := range carriers {
		sql := fmt.Sprintf(`
INSERT INTO shipping_margin (created_at, updated_at, business_id, carrier_code, carrier_name, margin_amount, insurance_margin, is_active)
SELECT NOW(), NOW(), b.id, %s, %s, 0, 0, true
FROM business b
WHERE b.deleted_at IS NULL
  AND EXISTS (
    SELECT 1
    FROM integrations i
    JOIN integration_types it ON it.id = i.integration_type_id
    JOIN integration_categories ic ON ic.id = it.category_id
    WHERE i.business_id = b.id
      AND i.deleted_at IS NULL
      AND ic.code = 'shipping'
  )
  AND NOT EXISTS (
    SELECT 1 FROM shipping_margin sm
    WHERE sm.business_id = b.id AND sm.carrier_code = %s AND sm.deleted_at IS NULL
  )
`, quote(c.Code), quote(c.Name), quote(c.Code))
		if err := db.Exec(sql).Error; err != nil {
			return fmt.Errorf("seed margin %s: %w", c.Code, err)
		}
	}

	return nil
}
