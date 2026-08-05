package repository

import (
	"context"
	"fmt"
)

const defaultCodMarginAmount = 500

func (r *Repository) seedCodMarginAmount(ctx context.Context) error {
	res := r.db.Conn(ctx).Exec(`
UPDATE shipping_margin
SET cod_margin_amount = ?,
    cod_margin_percent = 0,
    updated_at = NOW()
WHERE deleted_at IS NULL
  AND cod_margin_percent > 0
  AND cod_margin_amount = 0
`, defaultCodMarginAmount)

	if res.Error != nil {
		return fmt.Errorf("seed cod margin amount: %w", res.Error)
	}

	if res.RowsAffected > 0 {
		fmt.Printf("seed cod margin amount: %d margenes pasados de porcentaje a monto fijo de %d\n", res.RowsAffected, defaultCodMarginAmount)
	}

	return nil
}
