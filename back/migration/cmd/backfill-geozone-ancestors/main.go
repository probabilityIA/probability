package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/migration/shared/db"
	"github.com/secamc93/probability/back/migration/shared/env"
	"github.com/secamc93/probability/back/migration/shared/log"
)

const orderBatchSQL = `
WITH RECURSIVE targets AS (
    SELECT id, shipping_lng, shipping_lat, business_id
    FROM orders
    WHERE deleted_at IS NULL
      AND shipping_lat IS NOT NULL
      AND shipping_lng IS NOT NULL
      AND geozone_state_id IS NULL
    LIMIT $1
), src AS (
    SELECT id AS order_id, business_id,
           ST_SetSRID(ST_MakePoint(shipping_lng, shipping_lat), 4326) AS p
    FROM targets
), match AS (
    SELECT DISTINCT ON (s.order_id)
           s.order_id, g.id AS gid, g.type AS gtype
    FROM src s
    JOIN geozones g
      ON g.deleted_at IS NULL
     AND g.is_active
     AND (g.business_id = 0 OR g.business_id = s.business_id)
     AND ST_Contains(g.geometry, s.p)
    ORDER BY s.order_id, CASE g.type
        WHEN 'barrio' THEN 1
        WHEN 'neighborhood' THEN 2
        WHEN 'admin_district' THEN 3
        WHEN 'locality' THEN 4
        WHEN 'city' THEN 5
        WHEN 'state' THEN 6
        WHEN 'country' THEN 7
        ELSE 9 END
), chain AS (
    SELECT m.order_id, g.id, g.type, ARRAY[g.id]::bigint[] AS path
    FROM match m JOIN geozones g ON g.id = m.gid
    UNION ALL
    SELECT c.order_id, p.id, p.type, c.path || p.id
    FROM chain c JOIN geozones p ON p.id = (SELECT parent_id FROM geozones WHERE id = c.id)
    WHERE p.deleted_at IS NULL
), levels AS (
    SELECT order_id,
        MAX(id) FILTER (WHERE type = 'country')         AS country_id,
        MAX(id) FILTER (WHERE type = 'state')           AS state_id,
        MAX(id) FILTER (WHERE type = 'city')            AS city_id,
        MAX(id) FILTER (WHERE type = 'admin_district')  AS admin_district_id,
        MAX(id) FILTER (WHERE type = 'locality')        AS locality_id,
        MAX(id) FILTER (WHERE type = 'neighborhood')    AS neighborhood_id,
        MAX(id) FILTER (WHERE type = 'barrio')          AS barrio_id,
        (
          SELECT to_jsonb(c2.path)
          FROM chain c2
          WHERE c2.order_id = chain.order_id
          ORDER BY array_length(c2.path, 1) DESC
          LIMIT 1
        ) AS path_json
    FROM chain
    GROUP BY order_id
)
UPDATE orders o
SET destination_point = (SELECT p::geography FROM src WHERE src.order_id = o.id),
    destination_geozone_id = m.gid,
    destination_geozone_path = COALESCE(l.path_json, o.destination_geozone_path),
    geozone_country_id = l.country_id,
    geozone_state_id = l.state_id,
    geozone_city_id = l.city_id,
    geozone_admin_district_id = l.admin_district_id,
    geozone_locality_id = l.locality_id,
    geozone_neighborhood_id = l.neighborhood_id,
    geozone_barrio_id = l.barrio_id
FROM match m
JOIN levels l ON l.order_id = m.order_id
WHERE o.id = m.order_id
`

const orderFallbackSQL = `
WITH RECURSIVE candidates AS (
    SELECT o.id, o.business_id, o.shipping_city, o.shipping_state
    FROM orders o
    WHERE o.deleted_at IS NULL
      AND o.geozone_state_id IS NULL
      AND o.shipping_city IS NOT NULL
      AND o.shipping_city <> ''
      AND o.business_id IS NOT NULL
    LIMIT $1
), match AS (
    SELECT DISTINCT ON (c.id) c.id AS order_id, g.id AS gid
    FROM candidates c
    JOIN geozones g
      ON g.deleted_at IS NULL
     AND g.is_active = TRUE
     AND (g.business_id = 0 OR g.business_id = c.business_id)
     AND (
        (g.type = 'city'  AND unaccent(lower(g.name)) = unaccent(lower(c.shipping_city))) OR
        (g.type = 'state' AND c.shipping_state IS NOT NULL AND unaccent(lower(g.name)) = unaccent(lower(c.shipping_state)))
     )
    ORDER BY c.id, CASE g.type WHEN 'city' THEN 1 WHEN 'state' THEN 2 ELSE 9 END
), chain AS (
    SELECT m.order_id, g.id, g.parent_id, g.type, ARRAY[g.id]::bigint[] AS path
    FROM match m JOIN geozones g ON g.id = m.gid
    UNION ALL
    SELECT c.order_id, g.id, g.parent_id, g.type, c.path || g.id
    FROM chain c JOIN geozones g ON g.id = c.parent_id
    WHERE g.deleted_at IS NULL
), levels AS (
    SELECT order_id,
        MAX(id) FILTER (WHERE type = 'country') AS country_id,
        MAX(id) FILTER (WHERE type = 'state')   AS state_id,
        MAX(id) FILTER (WHERE type = 'city')    AS city_id,
        (SELECT to_jsonb(c2.path) FROM chain c2 WHERE c2.order_id = chain.order_id ORDER BY array_length(c2.path,1) DESC LIMIT 1) AS path_json
    FROM chain GROUP BY order_id
)
UPDATE orders o
SET destination_geozone_id = m.gid,
    destination_geozone_path = COALESCE(l.path_json, o.destination_geozone_path),
    geozone_country_id = l.country_id,
    geozone_state_id = l.state_id,
    geozone_city_id = l.city_id
FROM match m
JOIN levels l ON l.order_id = m.order_id
WHERE o.id = m.order_id
`

const shipmentFallbackSQL = `
WITH RECURSIVE candidates AS (
    SELECT s.id, o.business_id, o.shipping_city, o.shipping_state
    FROM shipments s
    JOIN orders o ON o.id = s.order_id
    WHERE s.deleted_at IS NULL
      AND s.geozone_state_id IS NULL
      AND o.shipping_city IS NOT NULL
      AND o.shipping_city <> ''
      AND o.business_id IS NOT NULL
    LIMIT $1
), match AS (
    SELECT DISTINCT ON (c.id) c.id AS shipment_id, g.id AS gid
    FROM candidates c
    JOIN geozones g
      ON g.deleted_at IS NULL
     AND g.is_active = TRUE
     AND (g.business_id = 0 OR g.business_id = c.business_id)
     AND (
        (g.type = 'city'  AND unaccent(lower(g.name)) = unaccent(lower(c.shipping_city))) OR
        (g.type = 'state' AND c.shipping_state IS NOT NULL AND unaccent(lower(g.name)) = unaccent(lower(c.shipping_state)))
     )
    ORDER BY c.id, CASE g.type WHEN 'city' THEN 1 WHEN 'state' THEN 2 ELSE 9 END
), chain AS (
    SELECT m.shipment_id, g.id, g.parent_id, g.type, ARRAY[g.id]::bigint[] AS path
    FROM match m JOIN geozones g ON g.id = m.gid
    UNION ALL
    SELECT c.shipment_id, g.id, g.parent_id, g.type, c.path || g.id
    FROM chain c JOIN geozones g ON g.id = c.parent_id
    WHERE g.deleted_at IS NULL
), levels AS (
    SELECT shipment_id,
        MAX(id) FILTER (WHERE type = 'country') AS country_id,
        MAX(id) FILTER (WHERE type = 'state')   AS state_id,
        MAX(id) FILTER (WHERE type = 'city')    AS city_id,
        (SELECT to_jsonb(c2.path) FROM chain c2 WHERE c2.shipment_id = chain.shipment_id ORDER BY array_length(c2.path,1) DESC LIMIT 1) AS path_json
    FROM chain GROUP BY shipment_id
)
UPDATE shipments s
SET destination_geozone_id = m.gid,
    destination_geozone_path = COALESCE(l.path_json, s.destination_geozone_path),
    geozone_country_id = l.country_id,
    geozone_state_id = l.state_id,
    geozone_city_id = l.city_id
FROM match m
JOIN levels l ON l.shipment_id = m.shipment_id
WHERE s.id = m.shipment_id
`

const shipmentBatchSQL = `
WITH RECURSIVE targets AS (
    SELECT s.id, o.shipping_lng, o.shipping_lat, o.business_id
    FROM shipments s
    JOIN orders o ON o.id = s.order_id
    WHERE s.deleted_at IS NULL
      AND o.shipping_lat IS NOT NULL
      AND o.shipping_lng IS NOT NULL
      AND s.geozone_state_id IS NULL
    LIMIT $1
), src AS (
    SELECT id AS shipment_id, business_id,
           ST_SetSRID(ST_MakePoint(shipping_lng, shipping_lat), 4326) AS p
    FROM targets
), match AS (
    SELECT DISTINCT ON (s.shipment_id)
           s.shipment_id, g.id AS gid, g.type AS gtype
    FROM src s
    JOIN geozones g
      ON g.deleted_at IS NULL
     AND g.is_active
     AND (g.business_id = 0 OR g.business_id = s.business_id)
     AND ST_Contains(g.geometry, s.p)
    ORDER BY s.shipment_id, CASE g.type
        WHEN 'barrio' THEN 1
        WHEN 'neighborhood' THEN 2
        WHEN 'admin_district' THEN 3
        WHEN 'locality' THEN 4
        WHEN 'city' THEN 5
        WHEN 'state' THEN 6
        WHEN 'country' THEN 7
        ELSE 9 END
), chain AS (
    SELECT m.shipment_id, g.id, g.type, ARRAY[g.id]::bigint[] AS path
    FROM match m JOIN geozones g ON g.id = m.gid
    UNION ALL
    SELECT c.shipment_id, p.id, p.type, c.path || p.id
    FROM chain c JOIN geozones p ON p.id = (SELECT parent_id FROM geozones WHERE id = c.id)
    WHERE p.deleted_at IS NULL
), levels AS (
    SELECT shipment_id,
        MAX(id) FILTER (WHERE type = 'country')         AS country_id,
        MAX(id) FILTER (WHERE type = 'state')           AS state_id,
        MAX(id) FILTER (WHERE type = 'city')            AS city_id,
        MAX(id) FILTER (WHERE type = 'admin_district')  AS admin_district_id,
        MAX(id) FILTER (WHERE type = 'locality')        AS locality_id,
        MAX(id) FILTER (WHERE type = 'neighborhood')    AS neighborhood_id,
        MAX(id) FILTER (WHERE type = 'barrio')          AS barrio_id,
        (
          SELECT to_jsonb(c2.path)
          FROM chain c2
          WHERE c2.shipment_id = chain.shipment_id
          ORDER BY array_length(c2.path, 1) DESC
          LIMIT 1
        ) AS path_json
    FROM chain
    GROUP BY shipment_id
)
UPDATE shipments s
SET destination_point = (SELECT p::geography FROM src WHERE src.shipment_id = s.id),
    destination_geozone_id = m.gid,
    destination_geozone_path = COALESCE(l.path_json, s.destination_geozone_path),
    geozone_country_id = l.country_id,
    geozone_state_id = l.state_id,
    geozone_city_id = l.city_id,
    geozone_admin_district_id = l.admin_district_id,
    geozone_locality_id = l.locality_id,
    geozone_neighborhood_id = l.neighborhood_id,
    geozone_barrio_id = l.barrio_id
FROM match m
JOIN levels l ON l.shipment_id = m.shipment_id
WHERE s.id = m.shipment_id
`

func main() {
	target := flag.String("target", "all", "what to backfill: orders, shipments, all")
	batch := flag.Int("batch", 500, "rows per batch")
	flag.Parse()

	logger := log.New()
	cfg := env.New(logger)
	database := db.New(logger, cfg)
	defer database.Close()

	ctx := context.Background()
	conn := database.Conn(ctx)

	if *target == "orders" || *target == "all" {
		processed := 0
		for {
			start := time.Now()
			res := conn.Exec(orderBatchSQL, *batch)
			if res.Error != nil {
				logger.Fatal(ctx).Err(res.Error).Msg("orders backfill failed")
			}
			rows := res.RowsAffected
			processed += int(rows)
			logger.Info(ctx).Int64("rows", rows).Int("total", processed).Dur("dur", time.Since(start)).Msg("orders batch")
			if rows == 0 {
				break
			}
		}
		fmt.Printf("orders backfill complete (point-in-polygon): %d rows\n", processed)

		fallbackProcessed := 0
		for {
			start := time.Now()
			res := conn.Exec(orderFallbackSQL, *batch)
			if res.Error != nil {
				logger.Fatal(ctx).Err(res.Error).Msg("orders fallback backfill failed")
			}
			rows := res.RowsAffected
			fallbackProcessed += int(rows)
			logger.Info(ctx).Int64("rows", rows).Int("total", fallbackProcessed).Dur("dur", time.Since(start)).Msg("orders fallback batch")
			if rows == 0 {
				break
			}
		}
		fmt.Printf("orders fallback by name complete: %d rows\n", fallbackProcessed)
	}

	if *target == "shipments" || *target == "all" {
		processed := 0
		for {
			start := time.Now()
			res := conn.Exec(shipmentBatchSQL, *batch)
			if res.Error != nil {
				logger.Fatal(ctx).Err(res.Error).Msg("shipments backfill failed")
			}
			rows := res.RowsAffected
			processed += int(rows)
			logger.Info(ctx).Int64("rows", rows).Int("total", processed).Dur("dur", time.Since(start)).Msg("shipments batch")
			if rows == 0 {
				break
			}
		}
		fmt.Printf("shipments backfill complete (point-in-polygon): %d rows\n", processed)

		fallbackProcessed := 0
		for {
			start := time.Now()
			res := conn.Exec(shipmentFallbackSQL, *batch)
			if res.Error != nil {
				logger.Fatal(ctx).Err(res.Error).Msg("shipments fallback backfill failed")
			}
			rows := res.RowsAffected
			fallbackProcessed += int(rows)
			logger.Info(ctx).Int64("rows", rows).Int("total", fallbackProcessed).Dur("dur", time.Since(start)).Msg("shipments fallback batch")
			if rows == 0 {
				break
			}
		}
		fmt.Printf("shipments fallback by name complete: %d rows\n", fallbackProcessed)
	}
}
