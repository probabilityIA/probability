package repository

import (
	"context"
	"fmt"
)

const (
	vig0010OrderID       = "f051894f-9ec9-48b3-b8cb-0b7e8d99247c"
	vig0010CodTotal      = 354900
	vig0010PaymentMethod = 6
)

func (r *Repository) fixVig0010Cod(ctx context.Context) error {
	res := r.db.Conn(ctx).Exec(`
UPDATE orders
SET is_cod = true,
    cod_total = ?,
    cod_includes_shipping = true,
    payment_method_id = ?,
    updated_at = NOW()
WHERE id = ?
  AND deleted_at IS NULL
  AND COALESCE(cod_total, 0) = 0
`, vig0010CodTotal, vig0010PaymentMethod, vig0010OrderID)

	if res.Error != nil {
		return fmt.Errorf("fix VIG-0010 cod: %w", res.Error)
	}

	return nil
}
