package repository

import (
	"context"
	"fmt"
)

var vigaCodRevertConfirmadas = []vigaCodPromesaFix{
	{"f051894f-9ec9-48b3-b8cb-0b7e8d99247c", "VIG-0010", 354820.90, 354900},
	{"fb99a9f5-7764-43cc-904e-9e7042af1519", "VIG-0033", 182351.40, 179073},
	{"85c54e23-92bf-46d9-a81d-b349dbaef632", "VIG-0035", 77274.10, 75879},
	{"80857e42-579d-4602-8d30-ff7a6e2ace5b", "VIG-0045", 137666.85, 137254},
	{"cbdc2366-6395-4e7c-9c10-5fe45debbffe", "VIG-0047", 95825.05, 93848},
	{"6885c6d0-b1e2-4946-af44-d842ecd1ca7d", "VIG-0049", 157273.35, 154168},
	{"5630a16d-885a-46d0-b2f6-43cc1edf1d6d", "VIG-0054", 227771.10, 226909},
}

func (r *Repository) FixVigaCodRevertConfirmadas(ctx context.Context) error {
	applied := 0

	for _, fix := range vigaCodRevertConfirmadas {
		res := r.db.Conn(ctx).Exec(`
UPDATE orders
SET cod_total = ?,
    updated_at = NOW()
WHERE id = ?
  AND business_id = ?
  AND deleted_at IS NULL
  AND cod_total = ?
  AND EXISTS (
      SELECT 1
      FROM cod_payment_cut_order cpo
      JOIN cod_payment_cut c ON c.id = cpo.cod_payment_cut_id
      WHERE cpo.order_id = orders.id
        AND cpo.deleted_at IS NULL
        AND c.deleted_at IS NULL
        AND c.status = 'confirmed'
  )
`, fix.Prometido, fix.OrderID, vigaCodPayoutBusinessID, fix.LiquidadoED)

		if res.Error != nil {
			return fmt.Errorf("orden %s: %w", fix.OrderNumber, res.Error)
		}
		applied += int(res.RowsAffected)
	}

	if applied > 0 {
		fmt.Printf("revert viga cod confirmadas: %d ordenes devueltas al valor del corte cerrado\n", applied)
	}

	return nil
}
