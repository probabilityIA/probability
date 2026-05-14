-- Agrega shipping_neighborhood a orders: barrio explicito del cliente para
-- poder resolver hasta nivel 'barrio' en la jerarquia de geozonas.
ALTER TABLE orders ADD COLUMN IF NOT EXISTS shipping_neighborhood VARCHAR(120);

CREATE INDEX IF NOT EXISTS idx_orders_shipping_neighborhood
    ON orders (LOWER(shipping_neighborhood))
    WHERE shipping_neighborhood IS NOT NULL AND shipping_neighborhood <> '';
