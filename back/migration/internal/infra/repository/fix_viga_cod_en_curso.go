package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

var vigaCodEnCursoFixes = []vigaCodPayoutFix{
	{"37c134e4-4d27-4333-a454-2b8b2ac1d405", "VIG-0055", "240057917044", 131619.15, 130860, 9341},
	{"01606b9f-8652-42f6-997b-49a7b65c33b7", "VIG-0056", "240057968236", 113214.75, 112651, 8189},
	{"e83c791a-cde8-44e8-a92e-555506cb40d8", "VIG-0061", "240058281660", 159300.90, 158535, 11092},
	{"27888ed9-5b8c-4972-a03d-ccbec14d9727", "VIG-0062", "84151700852", 767108.00, 766478, 19710},
	{"8566068c-d0dc-40f6-8535-3a482f04d130", "VIG-0063", "240058307393", 113214.75, 112651, 8189},
	{"7019049d-f84a-4cf4-a97b-7bab7e7290dc", "VIG-0064", "84151702342", 374505.05, 374187, 10146},
	{"490ace3f-8f20-4249-b9ec-75dfebba720c", "VIG-0065", "84151702488", 171948.40, 171949, 6116},
	{"04d9262a-e2b4-4085-b27a-afe663de05d4", "VIG-0066", "240058455557", 65050.75, 64701, 5155},
	{"7e864718-49e6-4443-9436-c129dfb1b41e", "VIG-0067", "84151703577", 255307.70, 255084, 7242},
	{"3345d1bb-709d-448c-bd2f-c1ae3ce2b347", "VIG-0068", "240058483021", 159300.90, 158535, 11092},
}

func (r *Repository) fixVigaCodEnCurso(ctx context.Context) error {
	applied := 0

	for _, fix := range vigaCodEnCursoFixes {
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
`, fix.RealPayout, fix.OrderID, vigaCodPayoutBusinessID, fix.CurrentCod)

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
			return fmt.Errorf("fix cod en curso viga: %w", err)
		}
	}

	if applied > 0 {
		fmt.Printf("fix cod en curso viga: %d ordenes ajustadas al valor liquidado por EnvioClick\n", applied)
	}

	return nil
}
