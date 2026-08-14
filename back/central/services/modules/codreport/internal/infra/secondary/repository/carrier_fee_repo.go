package repository

import (
	"context"
	"errors"
)

func (r *Repository) UpdateShipmentCarrierFee(ctx context.Context, businessID uint, shipmentID uint, fee float64) error {
	res := r.db.Conn(ctx).Exec(`
UPDATE shipments
SET cod_carrier_fee = ?,
	updated_at = NOW()
WHERE id = ?
  AND deleted_at IS NULL
  AND EXISTS (
	SELECT 1 FROM orders o
	WHERE o.id = shipments.order_id
	  AND o.business_id = ?
	  AND o.deleted_at IS NULL
  )`, fee, shipmentID, businessID)

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("el envio no existe o no pertenece al negocio")
	}
	return nil
}

func (r *Repository) ShipmentCarrierFee(ctx context.Context, businessID uint, shipmentID uint) (float64, string, error) {
	var row struct {
		CodCarrierFee  float64
		TrackingNumber string
		OrderNumber    string
	}

	err := r.db.Conn(ctx).Raw(`
SELECT COALESCE(s.cod_carrier_fee,0) AS cod_carrier_fee,
	COALESCE(s.tracking_number,'') AS tracking_number,
	COALESCE(o.order_number,'') AS order_number
FROM shipments s
JOIN orders o ON o.id = s.order_id AND o.deleted_at IS NULL
WHERE s.id = ? AND s.deleted_at IS NULL AND o.business_id = ?
LIMIT 1`, shipmentID, businessID).Scan(&row).Error

	if err != nil {
		return 0, "", err
	}
	return row.CodCarrierFee, row.OrderNumber, nil
}
