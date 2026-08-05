package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

const vigaCodPayoutBusinessID = 46

type vigaCodPayoutFix struct {
	OrderID     string
	OrderNumber string
	Tracking    string
	CurrentCod  float64
	RealPayout  float64
	RealFee     float64
}

var vigaCodPayoutFixes = []vigaCodPayoutFix{
	{"b92b31be-c436-4e66-89f3-47ef8d3dd036", "VIG-0022", "84151679648", 73050.40, 73200, 6116},
	{"29f6862e-be60-4551-9d8c-aadba1adc991", "VIG-0025", "240056906880", 106163.60, 105864, 7760},
	{"fb99a9f5-7764-43cc-904e-9e7042af1519", "VIG-0033", "240057385411", 182351.40, 179073, 12391},
	{"85c54e23-92bf-46d9-a81d-b349dbaef632", "VIG-0035", "240057388565", 77274.10, 75879, 5863},
	{"80857e42-579d-4602-8d30-ff7a6e2ace5b", "VIG-0045", "240057615406", 137666.85, 137254, 9745},
	{"cbdc2366-6395-4e7c-9c10-5fe45debbffe", "VIG-0047", "240057763304", 95825.05, 93848, 6999},
	{"590aa23b-ca1d-49f4-83c0-d9c1aa7d1b75", "VIG-0048", "84151693358", 88261.40, 84245, 6116},
	{"6885c6d0-b1e2-4946-af44-d842ecd1ca7d", "VIG-0049", "240057810509", 157273.35, 154168, 10816},
	{"5630a16d-885a-46d0-b2f6-43cc1edf1d6d", "VIG-0054", "240057899857", 227771.10, 226909, 15417},
	{"20a63d7d-8367-4d00-9524-354790ce7889", "VIG-0057", "240057983211", 113214.75, 112651, 8189},
	{"75394fd7-dc0e-4850-9f4a-41f5fb1c51f1", "VIG-0058", "84151697392", 282775.90, 282530, 7912},
	{"ae77a290-3c60-43e9-bc8d-11a99f1d93f7", "VIG-0060", "240058130173", 198291.25, 197359, 13548},
}

func (r *Repository) fixVigaCodRealPayout(ctx context.Context) error {
	applied := 0

	for _, fix := range vigaCodPayoutFixes {
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
			return fmt.Errorf("fix cod real payout viga: %w", err)
		}
	}

	if applied > 0 {
		fmt.Printf("fix cod real payout viga: %d ordenes ajustadas al valor liquidado por EnvioClick\n", applied)
	}

	return nil
}
