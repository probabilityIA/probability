package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

var vigaCodPromesaResto = []vigaCodPromesaFix{
	{"b92b31be-c436-4e66-89f3-47ef8d3dd036", "VIG-0022", 73200, 73050.40},
	{"29f6862e-be60-4551-9d8c-aadba1adc991", "VIG-0025", 105864, 106163.60},
	{"37c134e4-4d27-4333-a454-2b8b2ac1d405", "VIG-0055", 130860, 131619.15},
	{"01606b9f-8652-42f6-997b-49a7b65c33b7", "VIG-0056", 112651, 113214.75},
	{"20a63d7d-8367-4d00-9524-354790ce7889", "VIG-0057", 112651, 113214.75},
	{"ae77a290-3c60-43e9-bc8d-11a99f1d93f7", "VIG-0060", 197359, 198291.25},
	{"e83c791a-cde8-44e8-a92e-555506cb40d8", "VIG-0061", 158535, 159300.90},
	{"8566068c-d0dc-40f6-8535-3a482f04d130", "VIG-0063", 112651, 113214.75},
	{"04d9262a-e2b4-4085-b27a-afe663de05d4", "VIG-0066", 64701, 65050.75},
	{"3345d1bb-709d-448c-bd2f-c1ae3ce2b347", "VIG-0068", 158535, 159300.90},
	{"d74b0eaa-1ae8-4b0a-8d4a-e99ffa507cfc", "VIG-0069", 132594, 135929.00},
	{"63ab9177-d47b-4cff-83ef-5343abc9a064", "VIG-0071", 109255, 111113.00},
	{"952b0c53-ac89-4f7a-a3a7-6db8b1f519aa", "VIG-0072", 89343, 90254.00},
}

func (r *Repository) FixVigaCodPromesaResto(ctx context.Context) error {
	applied := 0

	for _, fix := range vigaCodPromesaResto {
		err := r.db.Conn(ctx).Transaction(func(tx *gorm.DB) error {
			res := tx.Exec(`
UPDATE orders
SET cod_total = ?,
    updated_at = NOW()
WHERE id = ?
  AND business_id = ?
  AND deleted_at IS NULL
  AND cod_total = ?
`, fix.Prometido, fix.OrderID, vigaCodPayoutBusinessID, fix.LiquidadoED)

			if res.Error != nil {
				return fmt.Errorf("orden %s: %w", fix.OrderNumber, res.Error)
			}
			if res.RowsAffected == 0 {
				return nil
			}

			cutRes := tx.Exec(`
UPDATE cod_payment_cut_order cpo
SET cod_amount = ?,
    updated_at = NOW()
FROM cod_payment_cut c
WHERE cpo.cod_payment_cut_id = c.id
  AND cpo.order_id = ?
  AND cpo.deleted_at IS NULL
  AND c.deleted_at IS NULL
  AND c.status = 'draft'
`, fix.Prometido, fix.OrderID)

			if cutRes.Error != nil {
				return fmt.Errorf("corte de %s: %w", fix.OrderNumber, cutRes.Error)
			}

			applied++
			return nil
		})
		if err != nil {
			return fmt.Errorf("fix cod promesa resto viga: %w", err)
		}
	}

	if applied == 0 {
		return nil
	}

	totales := r.db.Conn(ctx).Exec(`
UPDATE cod_payment_cut c
SET total_collected = t.suma,
    total_net = t.suma - COALESCE(c.total_discount, 0),
    orders_count = t.cantidad,
    updated_at = NOW()
FROM (
    SELECT cpo.cod_payment_cut_id AS cut_id,
           SUM(cpo.cod_amount) AS suma,
           COUNT(*) AS cantidad
    FROM cod_payment_cut_order cpo
    WHERE cpo.deleted_at IS NULL
    GROUP BY cpo.cod_payment_cut_id
) t
WHERE c.id = t.cut_id
  AND c.business_id = ?
  AND c.status = 'draft'
  AND c.deleted_at IS NULL
`, vigaCodPayoutBusinessID)

	if totales.Error != nil {
		return fmt.Errorf("totales del corte viga: %w", totales.Error)
	}

	fmt.Printf("fix cod promesa resto viga: %d ordenes devueltas al valor prometido al negocio\n", applied)

	return nil
}
