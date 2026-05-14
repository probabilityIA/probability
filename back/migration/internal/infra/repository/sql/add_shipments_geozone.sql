ALTER TABLE shipments
    ADD COLUMN IF NOT EXISTS destination_geozone_id BIGINT
        REFERENCES geozones(id) ON UPDATE CASCADE ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS destination_point geography(Point, 4326),
    ADD COLUMN IF NOT EXISTS destination_geozone_path JSONB;

CREATE INDEX IF NOT EXISTS idx_shipments_destination_geozone_id
    ON shipments (destination_geozone_id);

CREATE INDEX IF NOT EXISTS idx_shipments_destination_point
    ON shipments USING GIST (destination_point);

CREATE INDEX IF NOT EXISTS idx_shipments_destination_geozone_path
    ON shipments USING GIN (destination_geozone_path);

UPDATE shipments s
SET destination_point = ST_SetSRID(ST_MakePoint(o.shipping_lng, o.shipping_lat), 4326)::geography
FROM orders o
WHERE o.id = s.order_id
  AND s.deleted_at IS NULL
  AND s.destination_point IS NULL
  AND o.shipping_lat IS NOT NULL
  AND o.shipping_lng IS NOT NULL;
