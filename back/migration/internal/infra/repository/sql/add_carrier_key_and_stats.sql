-- Geozone probability optimization: carrier_key normalization + aggregated stats table.
--
-- Three parts:
--   1) carrier_key as generated columns on shipments and geozone_monthly_stats
--      to make REGEXP_REPLACE/UPPER lookups sargable via plain b-tree indexes.
--   2) Partial indexes on (carrier_key, geozone_<level>_id) WHERE deleted_at IS NULL
--      to accelerate live aggregations.
--   3) geozone_carrier_stats table that materializes (live + historical) totals
--      per (geozone_level, geozone_id, carrier_key) so the by-carrier endpoint
--      can be answered with a single indexed SELECT.

-- --------------------------------------------------------------------------
-- 1) carrier_key generated columns
-- --------------------------------------------------------------------------

ALTER TABLE shipments
    ADD COLUMN IF NOT EXISTS carrier_key VARCHAR(64)
        GENERATED ALWAYS AS (
            UPPER(REGEXP_REPLACE(COALESCE(carrier, ''), '[^a-zA-Z0-9]', '', 'g'))
        ) STORED;

ALTER TABLE geozone_monthly_stats
    ADD COLUMN IF NOT EXISTS carrier_key VARCHAR(64)
        GENERATED ALWAYS AS (
            UPPER(REGEXP_REPLACE(COALESCE(carrier, ''), '[^a-zA-Z0-9]', '', 'g'))
        ) STORED;

-- --------------------------------------------------------------------------
-- 2) Indexes on carrier_key + geozone levels
-- --------------------------------------------------------------------------

CREATE INDEX IF NOT EXISTS idx_shipments_carrier_key
    ON shipments (carrier_key)
    WHERE deleted_at IS NULL AND carrier_key <> '';

CREATE INDEX IF NOT EXISTS idx_shipments_ckey_country
    ON shipments (carrier_key, geozone_country_id)
    WHERE deleted_at IS NULL AND geozone_country_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_shipments_ckey_state
    ON shipments (carrier_key, geozone_state_id)
    WHERE deleted_at IS NULL AND geozone_state_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_shipments_ckey_city
    ON shipments (carrier_key, geozone_city_id)
    WHERE deleted_at IS NULL AND geozone_city_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_shipments_ckey_admin
    ON shipments (carrier_key, geozone_admin_district_id)
    WHERE deleted_at IS NULL AND geozone_admin_district_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_shipments_ckey_locality
    ON shipments (carrier_key, geozone_locality_id)
    WHERE deleted_at IS NULL AND geozone_locality_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_shipments_ckey_neighborhood
    ON shipments (carrier_key, geozone_neighborhood_id)
    WHERE deleted_at IS NULL AND geozone_neighborhood_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_shipments_ckey_barrio
    ON shipments (carrier_key, geozone_barrio_id)
    WHERE deleted_at IS NULL AND geozone_barrio_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_gms_carrier_key
    ON geozone_monthly_stats (carrier_key, geozone_id);

-- --------------------------------------------------------------------------
-- 3) geozone_carrier_stats: precomputed (live + historical) by zone+carrier
-- --------------------------------------------------------------------------

CREATE TABLE IF NOT EXISTS geozone_carrier_stats (
    geozone_level     VARCHAR(20) NOT NULL,
    geozone_id        BIGINT      NOT NULL,
    carrier_key       VARCHAR(64) NOT NULL,
    carrier_display   VARCHAR(128) NOT NULL DEFAULT '',
    total             BIGINT      NOT NULL DEFAULT 0,
    delivered         BIGINT      NOT NULL DEFAULT 0,
    cancelled         BIGINT      NOT NULL DEFAULT 0,
    returned          BIGINT      NOT NULL DEFAULT 0,
    failed            BIGINT      NOT NULL DEFAULT 0,
    in_transit        BIGINT      NOT NULL DEFAULT 0,
    sample_sufficient BOOLEAN     NOT NULL DEFAULT FALSE,
    last_refreshed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (geozone_level, geozone_id, carrier_key)
);

CREATE INDEX IF NOT EXISTS idx_gcs_lookup
    ON geozone_carrier_stats (geozone_id, geozone_level)
    WHERE sample_sufficient = TRUE;

CREATE INDEX IF NOT EXISTS idx_gcs_carrier
    ON geozone_carrier_stats (carrier_key, geozone_level);

CREATE INDEX IF NOT EXISTS idx_gcs_global
    ON geozone_carrier_stats (carrier_key)
    WHERE geozone_level = '__global__';

-- --------------------------------------------------------------------------
-- 4) Indice compuesto parcial para listado paginado de ordenes
--    (business_id, created_at DESC) WHERE deleted_at IS NULL
--    Acelera el endpoint GET /api/v1/orders (de ~600ms a ~10ms).
-- --------------------------------------------------------------------------

CREATE INDEX IF NOT EXISTS idx_orders_business_created_active
    ON orders (business_id, created_at DESC)
    WHERE deleted_at IS NULL;
