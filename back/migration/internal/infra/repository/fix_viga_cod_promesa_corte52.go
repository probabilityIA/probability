package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type vigaCodPromesaFix struct {
	OrderID     string
	OrderNumber string
	LiquidadoED float64
	Prometido   float64
}

var vigaCodPromesaFixes = []vigaCodPromesaFix{
	{"590aa23b-ca1d-49f4-83c0-d9c1aa7d1b75", "VIG-0048", 84245, 88261.40},
	{"75394fd7-dc0e-4850-9f4a-41f5fb1c51f1", "VIG-0058", 282530, 282775.90},
	{"27888ed9-5b8c-4972-a03d-ccbec14d9727", "VIG-0062", 766478, 767108.00},
	{"7019049d-f84a-4cf4-a97b-7bab7e7290dc", "VIG-0064", 374187, 374505.05},
	{"490ace3f-8f20-4249-b9ec-75dfebba720c", "VIG-0065", 171949, 171948.40},
	{"7e864718-49e6-4443-9436-c129dfb1b41e", "VIG-0067", 255084, 255307.70},
}

func (r *Repository) FixVigaCodPromesaCorte52(ctx context.Context) error {
	applied := 0

	for _, fix := range vigaCodPromesaFixes {
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
			return fmt.Errorf("fix cod promesa viga: %w", err)
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

	fmt.Printf("fix cod promesa viga: %d ordenes devueltas al valor prometido al negocio\n", applied)

	return nil
}
