-- Geozone analytics: denormalize ancestor ids onto orders/shipments
-- and create monthly aggregated stats table for delivery probability.

-- ORDERS: destination point + deepest geozone id + denormalized ancestors
ALTER TABLE orders
    ADD COLUMN IF NOT EXISTS destination_point geography(Point, 4326),
    ADD COLUMN IF NOT EXISTS destination_geozone_id BIGINT
        REFERENCES geozones(id) ON UPDATE CASCADE ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS destination_geozone_path JSONB,
    ADD COLUMN IF NOT EXISTS geozone_country_id BIGINT,
    ADD COLUMN IF NOT EXISTS geozone_state_id BIGINT,
    ADD COLUMN IF NOT EXISTS geozone_city_id BIGINT,
    ADD COLUMN IF NOT EXISTS geozone_admin_district_id BIGINT,
    ADD COLUMN IF NOT EXISTS geozone_locality_id BIGINT,
    ADD COLUMN IF NOT EXISTS geozone_neighborhood_id BIGINT,
    ADD COLUMN IF NOT EXISTS geozone_barrio_id BIGINT;

CREATE INDEX IF NOT EXISTS idx_orders_destination_geozone_id ON orders (destination_geozone_id);
CREATE INDEX IF NOT EXISTS idx_orders_destination_point ON orders USING GIST (destination_point);
CREATE INDEX IF NOT EXISTS idx_orders_geozone_state_id ON orders (geozone_state_id);
CREATE INDEX IF NOT EXISTS idx_orders_geozone_city_id ON orders (geozone_city_id);
CREATE INDEX IF NOT EXISTS idx_orders_geozone_locality_id ON orders (geozone_locality_id);
CREATE INDEX IF NOT EXISTS idx_orders_geozone_neighborhood_id ON orders (geozone_neighborhood_id);
CREATE INDEX IF NOT EXISTS idx_orders_geozone_barrio_id ON orders (geozone_barrio_id);

-- SHIPMENTS: denormalized ancestors (id columns) + carrier index for stats
ALTER TABLE shipments
    ADD COLUMN IF NOT EXISTS geozone_country_id BIGINT,
    ADD COLUMN IF NOT EXISTS geozone_state_id BIGINT,
    ADD COLUMN IF NOT EXISTS geozone_city_id BIGINT,
    ADD COLUMN IF NOT EXISTS geozone_admin_district_id BIGINT,
    ADD COLUMN IF NOT EXISTS geozone_locality_id BIGINT,
    ADD COLUMN IF NOT EXISTS geozone_neighborhood_id BIGINT,
    ADD COLUMN IF NOT EXISTS geozone_barrio_id BIGINT;

CREATE INDEX IF NOT EXISTS idx_shipments_geozone_state_id ON shipments (geozone_state_id);
CREATE INDEX IF NOT EXISTS idx_shipments_geozone_city_id ON shipments (geozone_city_id);
CREATE INDEX IF NOT EXISTS idx_shipments_geozone_locality_id ON shipments (geozone_locality_id);
CREATE INDEX IF NOT EXISTS idx_shipments_geozone_neighborhood_id ON shipments (geozone_neighborhood_id);
CREATE INDEX IF NOT EXISTS idx_shipments_geozone_barrio_id ON shipments (geozone_barrio_id);
CREATE INDEX IF NOT EXISTS idx_shipments_carrier_status ON shipments (carrier, status);

-- Monthly aggregated stats per (business, geozone, carrier, period).
-- A row with carrier_id = 0 represents the all-carriers aggregate for that geozone.
CREATE TABLE IF NOT EXISTS geozone_monthly_stats (
    id BIGSERIAL PRIMARY KEY,
    business_id BIGINT NOT NULL,
    period DATE NOT NULL,
    geozone_id BIGINT NOT NULL REFERENCES geozones(id) ON UPDATE CASCADE ON DELETE CASCADE,
    geozone_type VARCHAR(32) NOT NULL,
    carrier VARCHAR(128) NOT NULL DEFAULT '',
    total_shipments INT NOT NULL DEFAULT 0,
    delivered INT NOT NULL DEFAULT 0,
    cancelled INT NOT NULL DEFAULT 0,
    returned INT NOT NULL DEFAULT 0,
    in_transit INT NOT NULL DEFAULT 0,
    failed INT NOT NULL DEFAULT 0,
    total_attempts INT NOT NULL DEFAULT 0,
    avg_delivery_days NUMERIC(6,2),
    computed_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_geozone_monthly_stats UNIQUE (business_id, period, geozone_id, carrier)
);

CREATE INDEX IF NOT EXISTS idx_gms_lookup
    ON geozone_monthly_stats (business_id, geozone_id, carrier, period DESC);

CREATE INDEX IF NOT EXISTS idx_gms_type_period
    ON geozone_monthly_stats (business_id, geozone_type, period DESC);
