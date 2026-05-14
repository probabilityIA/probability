-- Refresh of geozone_carrier_stats: full upsert from live shipments + historical monthly stats.
-- Idempotent and safe to run concurrently with reads.
--
-- Rows are written for every (geozone_level, geozone_id, carrier_key) combination plus
-- a pseudo "__global__" level (geozone_id = 0) holding all-zone aggregates per carrier_key.
-- Stale rows (not touched in this refresh) are pruned at the end.

WITH
-- 1) Carriers to exclude (placeholders / internals).
excluded_carriers AS (
    SELECT UNNEST(ARRAY[
        'MANUAL','TEST','GIFTCARD','ENVIOCLICK'
    ])::text AS carrier_key
),

-- 2) Live shipments aggregated per (level, geozone_id, carrier_key).
live_per_level AS (
    SELECT 'barrio'::text AS geozone_level, geozone_barrio_id::bigint AS geozone_id, carrier_key,
           MAX(carrier) AS carrier_display,
           COUNT(*) FILTER (WHERE status IN ('delivered','failed','returned')) AS total,
           COUNT(*) FILTER (WHERE status = 'delivered') AS delivered,
           COUNT(*) FILTER (WHERE status = 'cancelled') AS cancelled,
           COUNT(*) FILTER (WHERE status = 'returned')  AS returned,
           COUNT(*) FILTER (WHERE status = 'failed')    AS failed,
           COUNT(*) FILTER (WHERE status NOT IN ('delivered','cancelled','returned','failed')) AS in_transit
    FROM shipments
    WHERE deleted_at IS NULL AND geozone_barrio_id IS NOT NULL
      AND carrier_key <> '' AND carrier_key NOT IN (SELECT carrier_key FROM excluded_carriers)
    GROUP BY geozone_barrio_id, carrier_key
    UNION ALL
    SELECT 'neighborhood', geozone_neighborhood_id::bigint, carrier_key,
           MAX(carrier),
           COUNT(*) FILTER (WHERE status IN ('delivered','failed','returned')),
           COUNT(*) FILTER (WHERE status = 'delivered'),
           COUNT(*) FILTER (WHERE status = 'cancelled'),
           COUNT(*) FILTER (WHERE status = 'returned'),
           COUNT(*) FILTER (WHERE status = 'failed'),
           COUNT(*) FILTER (WHERE status NOT IN ('delivered','cancelled','returned','failed'))
    FROM shipments
    WHERE deleted_at IS NULL AND geozone_neighborhood_id IS NOT NULL
      AND carrier_key <> '' AND carrier_key NOT IN (SELECT carrier_key FROM excluded_carriers)
    GROUP BY geozone_neighborhood_id, carrier_key
    UNION ALL
    SELECT 'admin_district', geozone_admin_district_id::bigint, carrier_key,
           MAX(carrier),
           COUNT(*) FILTER (WHERE status IN ('delivered','failed','returned')),
           COUNT(*) FILTER (WHERE status = 'delivered'),
           COUNT(*) FILTER (WHERE status = 'cancelled'),
           COUNT(*) FILTER (WHERE status = 'returned'),
           COUNT(*) FILTER (WHERE status = 'failed'),
           COUNT(*) FILTER (WHERE status NOT IN ('delivered','cancelled','returned','failed'))
    FROM shipments
    WHERE deleted_at IS NULL AND geozone_admin_district_id IS NOT NULL
      AND carrier_key <> '' AND carrier_key NOT IN (SELECT carrier_key FROM excluded_carriers)
    GROUP BY geozone_admin_district_id, carrier_key
    UNION ALL
    SELECT 'locality', geozone_locality_id::bigint, carrier_key,
           MAX(carrier),
           COUNT(*) FILTER (WHERE status IN ('delivered','failed','returned')),
           COUNT(*) FILTER (WHERE status = 'delivered'),
           COUNT(*) FILTER (WHERE status = 'cancelled'),
           COUNT(*) FILTER (WHERE status = 'returned'),
           COUNT(*) FILTER (WHERE status = 'failed'),
           COUNT(*) FILTER (WHERE status NOT IN ('delivered','cancelled','returned','failed'))
    FROM shipments
    WHERE deleted_at IS NULL AND geozone_locality_id IS NOT NULL
      AND carrier_key <> '' AND carrier_key NOT IN (SELECT carrier_key FROM excluded_carriers)
    GROUP BY geozone_locality_id, carrier_key
    UNION ALL
    SELECT 'city', geozone_city_id::bigint, carrier_key,
           MAX(carrier),
           COUNT(*) FILTER (WHERE status IN ('delivered','failed','returned')),
           COUNT(*) FILTER (WHERE status = 'delivered'),
           COUNT(*) FILTER (WHERE status = 'cancelled'),
           COUNT(*) FILTER (WHERE status = 'returned'),
           COUNT(*) FILTER (WHERE status = 'failed'),
           COUNT(*) FILTER (WHERE status NOT IN ('delivered','cancelled','returned','failed'))
    FROM shipments
    WHERE deleted_at IS NULL AND geozone_city_id IS NOT NULL
      AND carrier_key <> '' AND carrier_key NOT IN (SELECT carrier_key FROM excluded_carriers)
    GROUP BY geozone_city_id, carrier_key
    UNION ALL
    SELECT 'state', geozone_state_id::bigint, carrier_key,
           MAX(carrier),
           COUNT(*) FILTER (WHERE status IN ('delivered','failed','returned')),
           COUNT(*) FILTER (WHERE status = 'delivered'),
           COUNT(*) FILTER (WHERE status = 'cancelled'),
           COUNT(*) FILTER (WHERE status = 'returned'),
           COUNT(*) FILTER (WHERE status = 'failed'),
           COUNT(*) FILTER (WHERE status NOT IN ('delivered','cancelled','returned','failed'))
    FROM shipments
    WHERE deleted_at IS NULL AND geozone_state_id IS NOT NULL
      AND carrier_key <> '' AND carrier_key NOT IN (SELECT carrier_key FROM excluded_carriers)
    GROUP BY geozone_state_id, carrier_key
    UNION ALL
    SELECT 'country', geozone_country_id::bigint, carrier_key,
           MAX(carrier),
           COUNT(*) FILTER (WHERE status IN ('delivered','failed','returned')),
           COUNT(*) FILTER (WHERE status = 'delivered'),
           COUNT(*) FILTER (WHERE status = 'cancelled'),
           COUNT(*) FILTER (WHERE status = 'returned'),
           COUNT(*) FILTER (WHERE status = 'failed'),
           COUNT(*) FILTER (WHERE status NOT IN ('delivered','cancelled','returned','failed'))
    FROM shipments
    WHERE deleted_at IS NULL AND geozone_country_id IS NOT NULL
      AND carrier_key <> '' AND carrier_key NOT IN (SELECT carrier_key FROM excluded_carriers)
    GROUP BY geozone_country_id, carrier_key
),

-- 3) Historical aggregates from geozone_monthly_stats.
historical_per_level AS (
    SELECT geozone_type AS geozone_level, geozone_id::bigint, carrier_key,
           MAX(carrier) AS carrier_display,
           SUM(delivered + failed + returned)::bigint AS total,
           SUM(delivered)::bigint  AS delivered,
           SUM(cancelled)::bigint  AS cancelled,
           SUM(returned)::bigint   AS returned,
           SUM(failed)::bigint     AS failed,
           SUM(in_transit)::bigint AS in_transit
    FROM geozone_monthly_stats
    WHERE carrier_key <> '' AND carrier_key NOT IN (SELECT carrier_key FROM excluded_carriers)
    GROUP BY geozone_type, geozone_id, carrier_key
),

-- 4) Per-zone combined (live + historical).
combined_zones AS (
    SELECT geozone_level, geozone_id, carrier_key,
           MAX(carrier_display) AS carrier_display,
           SUM(total)::bigint     AS total,
           SUM(delivered)::bigint AS delivered,
           SUM(cancelled)::bigint AS cancelled,
           SUM(returned)::bigint  AS returned,
           SUM(failed)::bigint    AS failed,
           SUM(in_transit)::bigint AS in_transit
    FROM (
        SELECT * FROM live_per_level
        UNION ALL
        SELECT * FROM historical_per_level
    ) u
    GROUP BY geozone_level, geozone_id, carrier_key
),

-- 5) Global aggregates (all zones combined) per carrier_key.
live_global AS (
    SELECT carrier_key,
           MAX(carrier) AS carrier_display,
           COUNT(*) FILTER (WHERE status IN ('delivered','failed','returned')) AS total,
           COUNT(*) FILTER (WHERE status = 'delivered') AS delivered,
           COUNT(*) FILTER (WHERE status = 'cancelled') AS cancelled,
           COUNT(*) FILTER (WHERE status = 'returned')  AS returned,
           COUNT(*) FILTER (WHERE status = 'failed')    AS failed,
           COUNT(*) FILTER (WHERE status NOT IN ('delivered','cancelled','returned','failed')) AS in_transit
    FROM shipments
    WHERE deleted_at IS NULL
      AND carrier_key <> '' AND carrier_key NOT IN (SELECT carrier_key FROM excluded_carriers)
    GROUP BY carrier_key
),
historical_global AS (
    SELECT carrier_key,
           MAX(carrier) AS carrier_display,
           SUM(delivered + failed + returned)::bigint AS total,
           SUM(delivered)::bigint  AS delivered,
           SUM(cancelled)::bigint  AS cancelled,
           SUM(returned)::bigint   AS returned,
           SUM(failed)::bigint     AS failed,
           SUM(in_transit)::bigint AS in_transit
    FROM geozone_monthly_stats
    WHERE carrier_key <> '' AND carrier_key NOT IN (SELECT carrier_key FROM excluded_carriers)
    GROUP BY carrier_key
),
combined_global AS (
    SELECT '__global__'::text AS geozone_level, 0::bigint AS geozone_id, carrier_key,
           MAX(carrier_display) AS carrier_display,
           SUM(total)::bigint     AS total,
           SUM(delivered)::bigint AS delivered,
           SUM(cancelled)::bigint AS cancelled,
           SUM(returned)::bigint  AS returned,
           SUM(failed)::bigint    AS failed,
           SUM(in_transit)::bigint AS in_transit
    FROM (
        SELECT * FROM live_global
        UNION ALL
        SELECT * FROM historical_global
    ) u
    GROUP BY carrier_key
),

-- 6) All rows to write.
all_rows AS (
    SELECT * FROM combined_zones
    UNION ALL
    SELECT * FROM combined_global
),

-- 7) Upsert.
upserted AS (
    INSERT INTO geozone_carrier_stats AS gcs (
        geozone_level, geozone_id, carrier_key, carrier_display,
        total, delivered, cancelled, returned, failed, in_transit,
        sample_sufficient, last_refreshed_at
    )
    SELECT
        geozone_level, geozone_id, carrier_key,
        COALESCE(carrier_display, ''),
        total, delivered, cancelled, returned, failed, in_transit,
        (total >= 5),
        NOW()
    FROM all_rows
    ON CONFLICT (geozone_level, geozone_id, carrier_key) DO UPDATE
    SET carrier_display   = EXCLUDED.carrier_display,
        total             = EXCLUDED.total,
        delivered         = EXCLUDED.delivered,
        cancelled         = EXCLUDED.cancelled,
        returned          = EXCLUDED.returned,
        failed            = EXCLUDED.failed,
        in_transit        = EXCLUDED.in_transit,
        sample_sufficient = EXCLUDED.sample_sufficient,
        last_refreshed_at = NOW()
    RETURNING last_refreshed_at
)

-- 8) Delete rows that were not touched in this refresh (no longer have data).
DELETE FROM geozone_carrier_stats
WHERE last_refreshed_at < (SELECT MIN(last_refreshed_at) FROM upserted);
