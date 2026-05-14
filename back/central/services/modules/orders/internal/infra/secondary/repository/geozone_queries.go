package repository

import "context"

func (r *Repository) ResolveOrderGeozone(ctx context.Context, orderID string, businessID uint) error {
	if err := r.db.Conn(ctx).Exec(`
		WITH RECURSIVE src AS (
		    SELECT o.id AS order_id,
		           CASE
		               WHEN o.shipping_lat IS NOT NULL AND o.shipping_lng IS NOT NULL
		                   THEN ST_SetSRID(ST_MakePoint(o.shipping_lng, o.shipping_lat), 4326)
		           END AS p
		    FROM orders o
		    WHERE o.id = ?
		),
		match AS (
		    SELECT g.id, g.type
		    FROM geozones g, src
		    WHERE src.p IS NOT NULL
		      AND g.deleted_at IS NULL
		      AND g.is_active = TRUE
		      AND (g.business_id = 0 OR g.business_id = ?)
		      AND ST_Contains(g.geometry, src.p)
		    ORDER BY CASE g.type
		        WHEN 'barrio' THEN 1
		        WHEN 'neighborhood' THEN 2
		        WHEN 'admin_district' THEN 3
		        WHEN 'locality' THEN 4
		        WHEN 'city' THEN 5
		        WHEN 'state' THEN 6
		        WHEN 'country' THEN 7
		        ELSE 9 END
		    LIMIT 1
		),
		chain AS (
		    SELECT id, parent_id, type, ARRAY[id]::bigint[] AS path
		    FROM geozones WHERE id = (SELECT id FROM match)
		    UNION ALL
		    SELECT g.id, g.parent_id, g.type, c.path || g.id
		    FROM geozones g JOIN chain c ON g.id = c.parent_id
		    WHERE g.deleted_at IS NULL
		),
		levels AS (
		    SELECT
		        MAX(id) FILTER (WHERE type = 'country')         AS country_id,
		        MAX(id) FILTER (WHERE type = 'state')           AS state_id,
		        MAX(id) FILTER (WHERE type = 'city')            AS city_id,
		        MAX(id) FILTER (WHERE type = 'admin_district')  AS admin_district_id,
		        MAX(id) FILTER (WHERE type = 'locality')        AS locality_id,
		        MAX(id) FILTER (WHERE type = 'neighborhood')    AS neighborhood_id,
		        MAX(id) FILTER (WHERE type = 'barrio')          AS barrio_id,
		        (SELECT to_jsonb(path) FROM chain ORDER BY array_length(path, 1) DESC LIMIT 1) AS path_json
		    FROM chain
		)
		UPDATE orders
		SET destination_point = (SELECT p::geography FROM src),
		    destination_geozone_id = (SELECT id FROM match),
		    destination_geozone_path = COALESCE((SELECT path_json FROM levels), destination_geozone_path),
		    geozone_country_id = (SELECT country_id FROM levels),
		    geozone_state_id = (SELECT state_id FROM levels),
		    geozone_city_id = (SELECT city_id FROM levels),
		    geozone_admin_district_id = (SELECT admin_district_id FROM levels),
		    geozone_locality_id = (SELECT locality_id FROM levels),
		    geozone_neighborhood_id = (SELECT neighborhood_id FROM levels),
		    geozone_barrio_id = (SELECT barrio_id FROM levels)
		WHERE id = ? AND (SELECT id FROM match) IS NOT NULL
	`, orderID, businessID, orderID).Error; err != nil {
		return err
	}

	if err := r.db.Conn(ctx).Exec(`
		WITH RECURSIVE target AS (
		    SELECT id,
		           REGEXP_REPLACE(
		             TRIM(unaccent(lower(COALESCE(shipping_city,  '')))),
		             '\s*[,\(]?\s*d\.?\s*c\.?\s*\)?\s*$', '', 'g'
		           ) AS city_norm,
		           REGEXP_REPLACE(
		             TRIM(unaccent(lower(COALESCE(shipping_state, '')))),
		             '\s*[,\(]?\s*d\.?\s*c\.?\s*\)?\s*$', '', 'g'
		           ) AS state_norm
		    FROM orders WHERE id = ? AND geozone_state_id IS NULL
		      AND shipping_city IS NOT NULL AND shipping_city <> ''
		),
		matched_city AS (
		    SELECT t.id AS order_id, g.id AS gid
		    FROM target t
		    JOIN geozones g
		      ON g.deleted_at IS NULL AND g.type = 'city'
		     AND (g.business_id = 0 OR g.business_id = ?)
		     AND REGEXP_REPLACE(unaccent(lower(g.name)), '\s*[,\(]?\s*d\.?\s*c\.?\s*\)?\s*$', '', 'g') = t.city_norm
		    LIMIT 1
		),
		matched_state AS (
		    SELECT t.id AS order_id, g.id AS gid
		    FROM target t
		    JOIN geozones g
		      ON g.deleted_at IS NULL AND g.type = 'state'
		     AND (g.business_id = 0 OR g.business_id = ?)
		     AND REGEXP_REPLACE(unaccent(lower(g.name)), '\s*[,\(]?\s*d\.?\s*c\.?\s*\)?\s*$', '', 'g') IN (t.state_norm, t.city_norm)
		    LIMIT 1
		),
		picked AS (
		    SELECT order_id, gid FROM matched_city
		    UNION ALL
		    SELECT order_id, gid FROM matched_state WHERE NOT EXISTS (SELECT 1 FROM matched_city)
		    LIMIT 1
		),
		chain AS (
		    SELECT g.id, g.parent_id, g.type, ARRAY[g.id]::bigint[] AS path
		    FROM geozones g WHERE g.id = (SELECT gid FROM picked)
		    UNION ALL
		    SELECT g.id, g.parent_id, g.type, c.path || g.id
		    FROM geozones g JOIN chain c ON g.id = c.parent_id
		    WHERE g.deleted_at IS NULL
		),
		levels AS (
		    SELECT
		        MAX(id) FILTER (WHERE type = 'country') AS country_id,
		        MAX(id) FILTER (WHERE type = 'state')   AS state_id,
		        MAX(id) FILTER (WHERE type = 'city')    AS city_id,
		        (SELECT to_jsonb(path) FROM chain ORDER BY array_length(path, 1) DESC LIMIT 1) AS path_json
		    FROM chain
		)
		UPDATE orders
		SET destination_geozone_id = COALESCE(destination_geozone_id, (SELECT gid FROM picked)),
		    destination_geozone_path = COALESCE(destination_geozone_path, (SELECT path_json FROM levels)),
		    geozone_country_id = COALESCE(geozone_country_id, (SELECT country_id FROM levels)),
		    geozone_state_id = COALESCE(geozone_state_id, (SELECT state_id FROM levels)),
		    geozone_city_id = COALESCE(geozone_city_id, (SELECT city_id FROM levels))
		WHERE id = ? AND (SELECT gid FROM picked) IS NOT NULL
	`, orderID, businessID, businessID, orderID).Error; err != nil {
		return err
	}

	return r.resolveOrderBarrio(ctx, orderID, businessID)
}

func (r *Repository) resolveOrderBarrio(ctx context.Context, orderID string, businessID uint) error {
	return r.db.Conn(ctx).Exec(`
		WITH RECURSIVE target AS (
		    SELECT id, geozone_city_id,
		           TRIM(unaccent(lower(COALESCE(NULLIF(shipping_neighborhood, ''),
		                                        SPLIT_PART(shipping_street, ' | ', 3))))) AS barrio_norm
		    FROM orders
		    WHERE id = ?
		      AND geozone_city_id IS NOT NULL
		      AND geozone_barrio_id IS NULL
		),
		candidates AS (
		    SELECT g.id, g.parent_id, g.type, g.name
		    FROM target t
		    JOIN geozones g
		      ON g.deleted_at IS NULL
		     AND g.type IN ('barrio','neighborhood')
		     AND (g.business_id = 0 OR g.business_id = ?)
		     AND unaccent(lower(g.name)) = t.barrio_norm
		    WHERE t.barrio_norm <> ''
		),
		anc AS (
		    SELECT c.id AS leaf_id, c.id AS cur_id, c.parent_id, 0 AS depth
		    FROM candidates c
		    UNION ALL
		    SELECT a.leaf_id, g.id, g.parent_id, a.depth + 1
		    FROM anc a
		    JOIN geozones g ON g.id = a.parent_id AND g.deleted_at IS NULL
		    WHERE a.depth < 8
		),
		matched AS (
		    SELECT DISTINCT ON (t.id)
		           t.id AS order_id, a.leaf_id AS barrio_id
		    FROM target t
		    JOIN anc a ON a.cur_id = t.geozone_city_id
		    ORDER BY t.id, a.depth ASC
		),
		chain AS (
		    SELECT m.order_id, g.id, g.parent_id, g.type, ARRAY[g.id]::bigint[] AS path
		    FROM matched m
		    JOIN geozones g ON g.id = m.barrio_id
		    UNION ALL
		    SELECT c.order_id, g.id, g.parent_id, g.type, c.path || g.id
		    FROM chain c
		    JOIN geozones g ON g.id = c.parent_id AND g.deleted_at IS NULL
		),
		levels AS (
		    SELECT c.order_id,
		           MAX(c.id) FILTER (WHERE c.type = 'neighborhood')   AS neighborhood_id,
		           MAX(c.id) FILTER (WHERE c.type = 'admin_district') AS admin_district_id,
		           MAX(c.id) FILTER (WHERE c.type = 'locality')       AS locality_id,
		           MAX(c.id) FILTER (WHERE c.type = 'barrio')         AS barrio_id,
		           (SELECT to_jsonb(path) FROM chain c2
		            WHERE c2.order_id = c.order_id
		            ORDER BY array_length(path,1) DESC LIMIT 1) AS path_json
		    FROM chain c
		    GROUP BY c.order_id
		)
		UPDATE orders o
		SET destination_geozone_id   = COALESCE(l.barrio_id, o.destination_geozone_id),
		    destination_geozone_path = COALESCE(l.path_json, o.destination_geozone_path),
		    geozone_barrio_id        = COALESCE(o.geozone_barrio_id,        l.barrio_id),
		    geozone_neighborhood_id  = COALESCE(o.geozone_neighborhood_id,  l.neighborhood_id),
		    geozone_admin_district_id = COALESCE(o.geozone_admin_district_id, l.admin_district_id),
		    geozone_locality_id      = COALESCE(o.geozone_locality_id,      l.locality_id)
		FROM levels l
		WHERE o.id = l.order_id AND l.barrio_id IS NOT NULL
	`, orderID, businessID).Error
}
