package repository

import (
	"context"
	"fmt"
)

const (
	vigaCodMarginBusinessID = 46
	vigaCodMarginAmount     = 500
)

func (r *Repository) seedVigaCodMarginAmount(ctx context.Context) error {
	res := r.db.Conn(ctx).Exec(`
UPDATE shipping_margin
SET cod_margin_amount = ?,
    cod_margin_percent = 0,
    updated_at = NOW()
WHERE business_id = ?
  AND deleted_at IS NULL
  AND cod_margin_percent > 0
  AND cod_margin_amount = 0
`, vigaCodMarginAmount, vigaCodMarginBusinessID)

	if res.Error != nil {
		return fmt.Errorf("seed cod margin amount viga: %w", res.Error)
	}

	if res.RowsAffected > 0 {
		fmt.Printf("seed cod margin amount viga: %d transportadoras pasadas a monto fijo de %d\n", res.RowsAffected, vigaCodMarginAmount)
	}

	return nil
}
