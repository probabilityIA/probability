package repository

import (
	"context"
	"fmt"
)

const fixCodTotalCheckoutWooSQL = `
WITH lineas AS (
	SELECT o.id, sl
	FROM orders o
	CROSS JOIN LATERAL jsonb_array_elements(
		CASE WHEN jsonb_typeof(o.shipping_details->'shipping_lines') = 'array'
			THEN o.shipping_details->'shipping_lines' ELSE '[]'::jsonb END
	) sl
	WHERE o.platform = 'woocommerce'
	  AND o.is_cod
	  AND o.cod_includes_shipping
	  AND o.cod_total > 0
	  AND o.cod_checkout_carrier_fee = 0
	  AND o.deleted_at IS NULL
),
meta AS (
	SELECT l.id, md
	FROM lineas l
	CROSS JOIN LATERAL jsonb_array_elements(
		CASE WHEN jsonb_typeof(l.sl->'meta_data') = 'array'
			THEN l.sl->'meta_data' ELSE '[]'::jsonb END
	) md
),
candidatas AS (
	SELECT o.id,
		o.cod_total,
		o.total_amount - COALESCE(o.discount, 0) AS productos,
		(SELECT COALESCE(SUM((l.sl->>'total')::numeric), 0)
		   FROM lineas l
		  WHERE l.id = o.id
		    AND l.sl->>'total' ~ '^[0-9]+(\.[0-9]+){0,1}$') AS envio,
		(SELECT MAX((m.md->>'value')::numeric)
		   FROM meta m
		  WHERE m.id = o.id
		    AND m.md->>'key' = 'cod_carrier_fee'
		    AND m.md->>'value' ~ '^[0-9]+(\.[0-9]+){0,1}$') AS fee
	FROM orders o
	WHERE o.id IN (SELECT id FROM meta WHERE md->>'key' = 'quote_id' AND COALESCE(md->>'value', '') <> '')
	  AND o.id IN (SELECT id FROM meta WHERE md->>'key' = 'cod' AND md->>'value' = '1')
	  AND NOT EXISTS (
		SELECT 1 FROM cod_payment_cut_order cpo
		WHERE cpo.order_id = o.id AND cpo.deleted_at IS NULL
	  )
)
UPDATE orders o
SET cod_checkout_carrier_fee = c.fee,
    cod_total = CASE
        WHEN ABS(c.cod_total - (c.productos + c.envio)) <= 1 THEN c.cod_total - c.fee
        ELSE c.cod_total
    END,
    updated_at = NOW()
FROM candidatas c
WHERE o.id = c.id
  AND c.fee > 0
  AND (
      ABS(c.cod_total - (c.productos + c.envio)) <= 1
      OR ABS(c.cod_total - (c.productos + c.envio - c.fee)) <= 1
  )
`

func (r *Repository) FixCodTotalCheckoutWoo(ctx context.Context) error {
	res := r.db.Conn(ctx).Exec(fixCodTotalCheckoutWooSQL)
	if res.Error != nil {
		return fmt.Errorf("fix cod_total checkout woocommerce: %w", res.Error)
	}
	if res.RowsAffected > 0 {
		fmt.Printf("fix cod_total checkout woocommerce: %d ordenes separadas en cod_total neto + comision del checkout\n", res.RowsAffected)
	}
	return nil
}
