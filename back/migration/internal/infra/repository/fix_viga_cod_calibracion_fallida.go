package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

const vigaCodCalibracionBusinessID = 46

type vigaCodCalibracionFix struct {
	OrderID     string
	OrderNumber string
	Tracking    string
	CurrentCod  float64
	RealPayout  float64
	RealFee     float64
}

var vigaCodCalibracionFixes = []vigaCodCalibracionFix{
	{"63ab9177-d47b-4cff-83ef-5343abc9a064", "VIG-0071", "240058568295", 111113, 109255, 7974},
	{"952b0c53-ac89-4f7a-a3a7-6db8b1f519aa", "VIG-0072", "240058578276", 90254, 89343, 6714},
	{"d74b0eaa-1ae8-4b0a-8d4a-e99ffa507cfc", "VIG-0069", "240058579495", 135929, 132594, 9451},
}

func (r *Repository) fixVigaCodCalibracionFallida(ctx context.Context) error {
	applied := 0

	for _, fix := range vigaCodCalibracionFixes {
		err := r.db.Conn(ctx).Transaction(func(tx *gorm.DB) error {
			res := tx.Exec(`
UPDATE orders
SET cod_total = ?,
    updated_at = NOW()
WHERE id = ?
  AND business_id = ?
  AND deleted_at IS NULL
  AND cod_total = ?
  AND NOT EXISTS (
      SELECT 1 FROM cod_payment_cut_order cpo
      WHERE cpo.order_id = orders.id AND cpo.deleted_at IS NULL
  )
`, fix.RealPayout, fix.OrderID, vigaCodCalibracionBusinessID, fix.CurrentCod)

			if res.Error != nil {
				return fmt.Errorf("orden %s: %w", fix.OrderNumber, res.Error)
			}
			if res.RowsAffected == 0 {
				return nil
			}

			feeRes := tx.Exec(`
UPDATE shipments
SET cod_carrier_fee = ?,
    updated_at = NOW()
WHERE order_id = ?
  AND tracking_number = ?
  AND deleted_at IS NULL
`, fix.RealFee, fix.OrderID, fix.Tracking)

			if feeRes.Error != nil {
				return fmt.Errorf("envio %s: %w", fix.Tracking, feeRes.Error)
			}

			applied++
			return nil
		})
		if err != nil {
			return fmt.Errorf("fix cod calibracion fallida viga: %w", err)
		}
	}

	if applied > 0 {
		fmt.Printf("fix cod calibracion fallida viga: %d ordenes ajustadas al valor que liquida EnvioClick\n", applied)
	}

	return nil
}
