package repository

import (
	"context"
	"fmt"
)

type cotizacionFallidaSinMarcar struct {
	quoteID int
	message string
}

var cotizacionesFallidasSinMarcar = []cotizacionFallidaSinMarcar{
	{
		quoteID: 6542,
		message: "Falta un correo electrónico o no tiene un formato válido. Revisa el correo del destinatario en la orden y el de tu dirección de origen, y vuelve a generar la guía.",
	},
	{
		quoteID: 6441,
		message: "Un teléfono no tiene el formato correcto (entre 7 y 10 dígitos). Revisa el teléfono del destinatario en la orden y el de tu dirección de origen, y vuelve a generar la guía.",
	},
}

func (r *Repository) fixVig0095CotizacionFallida(ctx context.Context) error {
	for _, c := range cotizacionesFallidasSinMarcar {
		res := r.db.Conn(ctx).Exec(`
UPDATE shipping_quotes q
SET status = 'failed',
    error_message = ?,
    updated_at = NOW()
FROM shipments s
WHERE q.id = ?
  AND q.deleted_at IS NULL
  AND q.status = 'guide_generated'
  AND s.order_id = q.order_uuid
  AND s.deleted_at IS NULL
  AND s.status = 'failed'
  AND COALESCE(s.tracking_number, '') = ''
  AND COALESCE(s.guide_url, '') = ''
`, c.message, c.quoteID)

		if res.Error != nil {
			return fmt.Errorf("fix cotizacion fallida %d: %w", c.quoteID, res.Error)
		}
	}

	return nil
}
