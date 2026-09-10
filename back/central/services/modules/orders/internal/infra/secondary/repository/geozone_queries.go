package repository

import "context"

func (r *Repository) ResolveOrderGeozone(ctx context.Context, orderID string, businessID uint) error {
	if err := r.db.Conn(ctx).Exec(`
		WITH RECURSIVE ord AS (
		    SELECT o.id, o.shipping_lat, o.shipping_lng,
		           REGEXP_REPLACE(
		             TRIM(unaccent(lower(COALESCE(o.shipping_city,  '')))),
		             '\s*[,\(]{0,1}\s*d\.{0,1}\s*c\.{0,1}\s*\){0,1}\s*$', '', 'g'
		           ) AS city_norm,
		           REGEXP_REPLACE(
		             TRIM(unaccent(lower(COALESCE(o.shipping_state, '')))),
		             '\s*[,\(]{0,1}\s*d\.{0,1}\s*c\.{0,1}\s*\){0,1}\s*$', '', 'g'
		           ) AS state_norm
		    FROM orders o WHERE o.id = ?
		),
		src AS (
		    SELECT CASE
		               WHEN shipping_lat IS NOT NULL AND shipping_lng IS NOT NULL
		                   THEN ST_SetSRID(ST_MakePoint(shipping_lng, shipping_lat), 4326)
		           END AS p
		    FROM ord
		),
		point_match AS (
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
		point_chain AS (
		    SELECT id, parent_id, type, ARRAY[id]::bigint[] AS path
		    FROM geozones WHERE id = (SELECT id FROM point_match)
		    UNION ALL
		    SELECT g.id, g.parent_id, g.type, c.path || g.id
		    FROM geozones g JOIN point_chain c ON g.id = c.parent_id
		    WHERE g.deleted_at IS NULL
		),
		point_levels AS (
		    SELECT
		        MAX(id) FILTER (WHERE type = 'country')         AS country_id,
		        MAX(id) FILTER (WHERE type = 'state')           AS state_id,
		        MAX(id) FILTER (WHERE type = 'city')            AS city_id,
		        MAX(id) FILTER (WHERE type = 'admin_district')  AS admin_district_id,
		        MAX(id) FILTER (WHERE type = 'locality')        AS locality_id,
		        MAX(id) FILTER (WHERE type = 'neighborhood')    AS neighborhood_id,
		        MAX(id) FILTER (WHERE type = 'barrio')          AS barrio_id,
		        (SELECT to_jsonb(path) FROM point_chain ORDER BY array_length(path, 1) DESC LIMIT 1) AS path_json
		    FROM point_chain
		),
		text_city AS (
		    SELECT g.id AS gid
		    FROM ord, geozones g
		    LEFT JOIN geozones gs ON gs.id = g.parent_id AND gs.type = 'state'
		    WHERE ord.city_norm <> ''
		      AND g.deleted_at IS NULL AND g.type = 'city'
		      AND (g.business_id = 0 OR g.business_id = ?)
		      AND REGEXP_REPLACE(unaccent(lower(g.name)), '\s*[,\(]{0,1}\s*d\.{0,1}\s*c\.{0,1}\s*\){0,1}\s*$', '', 'g') = ord.city_norm
		    ORDER BY CASE
		        WHEN ord.state_norm <> '' AND gs.id IS NOT NULL
		             AND REGEXP_REPLACE(unaccent(lower(gs.name)), '\s*[,\(]{0,1}\s*d\.{0,1}\s*c\.{0,1}\s*\){0,1}\s*$', '', 'g') = ord.state_norm
		        THEN 0 ELSE 1 END
		    LIMIT 1
		),
		text_state AS (
		    SELECT g.id AS gid
		    FROM ord, geozones g
		    WHERE NOT EXISTS (SELECT 1 FROM text_city)
		      AND (ord.state_norm <> '' OR ord.city_norm <> '')
		      AND g.deleted_at IS NULL AND g.type = 'state'
		      AND (g.business_id = 0 OR g.business_id = ?)
		      AND REGEXP_REPLACE(unaccent(lower(g.name)), '\s*[,\(]{0,1}\s*d\.{0,1}\s*c\.{0,1}\s*\){0,1}\s*$', '', 'g') IN (ord.state_norm, ord.city_norm)
		    LIMIT 1
		),
		text_picked AS (
		    SELECT gid FROM text_city
		    UNION ALL
		    SELECT gid FROM text_state WHERE NOT EXISTS (SELECT 1 FROM text_city)
		    LIMIT 1
		),
		text_chain AS (
		    SELECT id, parent_id, type, ARRAY[id]::bigint[] AS path
		    FROM geozones WHERE id = (SELECT gid FROM text_picked)
		    UNION ALL
		    SELECT g.id, g.parent_id, g.type, c.path || g.id
		    FROM geozones g JOIN text_chain c ON g.id = c.parent_id
		    WHERE g.deleted_at IS NULL
		),
		text_levels AS (
		    SELECT
		        MAX(id) FILTER (WHERE type = 'country') AS country_id,
		        MAX(id) FILTER (WHERE type = 'state')   AS state_id,
		        MAX(id) FILTER (WHERE type = 'city')    AS city_id,
		        (SELECT to_jsonb(path) FROM text_chain ORDER BY array_length(path, 1) DESC LIMIT 1) AS path_json
		    FROM text_chain
		),
		decision AS (
		    SELECT CASE
		        WHEN (SELECT city_id FROM point_levels) IS NOT NULL
		             AND (SELECT city_id FROM text_levels) IS NOT NULL
		             AND (SELECT city_id FROM point_levels) <> (SELECT city_id FROM text_levels)
		            THEN 'text_wins'
		        WHEN (SELECT id FROM point_match) IS NOT NULL THEN 'point'
		        WHEN (SELECT gid FROM text_picked) IS NOT NULL THEN 'text'
		        ELSE 'none'
		    END AS mode
		)
		UPDATE orders o
		SET destination_point = CASE WHEN (SELECT mode FROM decision) = 'text_wins'
		        THEN NULL ELSE (SELECT p::geography FROM src) END,
		    destination_geozone_id = CASE (SELECT mode FROM decision)
		        WHEN 'point' THEN (SELECT id FROM point_match)
		        WHEN 'text' THEN (SELECT gid FROM text_picked)
		        WHEN 'text_wins' THEN (SELECT gid FROM text_picked)
		        ELSE o.destination_geozone_id END,
		    destination_geozone_path = CASE (SELECT mode FROM decision)
		        WHEN 'point' THEN (SELECT path_json FROM point_levels)
		        WHEN 'text' THEN (SELECT path_json FROM text_levels)
		        WHEN 'text_wins' THEN (SELECT path_json FROM text_levels)
		        ELSE o.destination_geozone_path END,
		    geozone_country_id = CASE (SELECT mode FROM decision)
		        WHEN 'point' THEN (SELECT country_id FROM point_levels)
		        WHEN 'text' THEN (SELECT country_id FROM text_levels)
		        WHEN 'text_wins' THEN (SELECT country_id FROM text_levels)
		        ELSE o.geozone_country_id END,
		    geozone_state_id = CASE (SELECT mode FROM decision)
		        WHEN 'point' THEN (SELECT state_id FROM point_levels)
		        WHEN 'text' THEN (SELECT state_id FROM text_levels)
		        WHEN 'text_wins' THEN (SELECT state_id FROM text_levels)
		        ELSE o.geozone_state_id END,
		    geozone_city_id = CASE (SELECT mode FROM decision)
		        WHEN 'point' THEN (SELECT city_id FROM point_levels)
		        WHEN 'text' THEN (SELECT city_id FROM text_levels)
		        WHEN 'text_wins' THEN (SELECT city_id FROM text_levels)
		        ELSE o.geozone_city_id END,
		    geozone_admin_district_id = CASE (SELECT mode FROM decision)
		        WHEN 'point' THEN (SELECT admin_district_id FROM point_levels)
		        WHEN 'text_wins' THEN NULL
		        ELSE o.geozone_admin_district_id END,
		    geozone_locality_id = CASE (SELECT mode FROM decision)
		        WHEN 'point' THEN (SELECT locality_id FROM point_levels)
		        WHEN 'text_wins' THEN NULL
		        ELSE o.geozone_locality_id END,
		    geozone_neighborhood_id = CASE (SELECT mode FROM decision)
		        WHEN 'point' THEN (SELECT neighborhood_id FROM point_levels)
		        WHEN 'text_wins' THEN NULL
		        ELSE o.geozone_neighborhood_id END,
		    geozone_barrio_id = CASE (SELECT mode FROM decision)
		        WHEN 'point' THEN (SELECT barrio_id FROM point_levels)
		        WHEN 'text_wins' THEN NULL
		        ELSE o.geozone_barrio_id END
		WHERE o.id = ? AND (SELECT mode FROM decision) <> 'none'
	`, orderID, businessID, businessID, businessID, orderID).Error; err != nil {
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
