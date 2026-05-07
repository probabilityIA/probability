package repository

import "context"

func (r *Repository) ResolveOrderGeozone(ctx context.Context, orderID string, businessID uint) error {
	return r.db.Conn(ctx).Exec(`
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
		WHERE id = ?
	`, orderID, businessID, orderID).Error
}
