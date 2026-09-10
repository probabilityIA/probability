package repository

import (
	"context"
	"strings"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

func (r *Repository) ResolveShipmentGeozone(ctx context.Context, shipmentID uint, businessID uint) error {
	return r.db.Conn(ctx).Exec(`
		WITH RECURSIVE ord AS (
		    SELECT o.shipping_lat, o.shipping_lng,
		           REGEXP_REPLACE(
		             TRIM(unaccent(lower(COALESCE(o.shipping_city,  '')))),
		             '\s*[,\(]{0,1}\s*d\.{0,1}\s*c\.{0,1}\s*\){0,1}\s*$', '', 'g'
		           ) AS city_norm,
		           REGEXP_REPLACE(
		             TRIM(unaccent(lower(COALESCE(o.shipping_state, '')))),
		             '\s*[,\(]{0,1}\s*d\.{0,1}\s*c\.{0,1}\s*\){0,1}\s*$', '', 'g'
		           ) AS state_norm
		    FROM shipments s
		    LEFT JOIN orders o ON o.id = s.order_id
		    WHERE s.id = ?
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
		UPDATE shipments s
		SET destination_point = CASE WHEN (SELECT mode FROM decision) = 'text_wins'
		        THEN NULL ELSE (SELECT p::geography FROM src) END,
		    destination_geozone_id = CASE (SELECT mode FROM decision)
		        WHEN 'point' THEN (SELECT id FROM point_match)
		        WHEN 'text' THEN (SELECT gid FROM text_picked)
		        WHEN 'text_wins' THEN (SELECT gid FROM text_picked)
		        ELSE s.destination_geozone_id END,
		    destination_geozone_path = CASE (SELECT mode FROM decision)
		        WHEN 'point' THEN (SELECT path_json FROM point_levels)
		        WHEN 'text' THEN (SELECT path_json FROM text_levels)
		        WHEN 'text_wins' THEN (SELECT path_json FROM text_levels)
		        ELSE s.destination_geozone_path END,
		    geozone_country_id = CASE (SELECT mode FROM decision)
		        WHEN 'point' THEN (SELECT country_id FROM point_levels)
		        WHEN 'text' THEN (SELECT country_id FROM text_levels)
		        WHEN 'text_wins' THEN (SELECT country_id FROM text_levels)
		        ELSE s.geozone_country_id END,
		    geozone_state_id = CASE (SELECT mode FROM decision)
		        WHEN 'point' THEN (SELECT state_id FROM point_levels)
		        WHEN 'text' THEN (SELECT state_id FROM text_levels)
		        WHEN 'text_wins' THEN (SELECT state_id FROM text_levels)
		        ELSE s.geozone_state_id END,
		    geozone_city_id = CASE (SELECT mode FROM decision)
		        WHEN 'point' THEN (SELECT city_id FROM point_levels)
		        WHEN 'text' THEN (SELECT city_id FROM text_levels)
		        WHEN 'text_wins' THEN (SELECT city_id FROM text_levels)
		        ELSE s.geozone_city_id END,
		    geozone_admin_district_id = CASE (SELECT mode FROM decision)
		        WHEN 'point' THEN (SELECT admin_district_id FROM point_levels)
		        WHEN 'text_wins' THEN NULL
		        ELSE s.geozone_admin_district_id END,
		    geozone_locality_id = CASE (SELECT mode FROM decision)
		        WHEN 'point' THEN (SELECT locality_id FROM point_levels)
		        WHEN 'text_wins' THEN NULL
		        ELSE s.geozone_locality_id END,
		    geozone_neighborhood_id = CASE (SELECT mode FROM decision)
		        WHEN 'point' THEN (SELECT neighborhood_id FROM point_levels)
		        WHEN 'text_wins' THEN NULL
		        ELSE s.geozone_neighborhood_id END,
		    geozone_barrio_id = CASE (SELECT mode FROM decision)
		        WHEN 'point' THEN (SELECT barrio_id FROM point_levels)
		        WHEN 'text_wins' THEN NULL
		        ELSE s.geozone_barrio_id END
		WHERE s.id = ? AND (SELECT mode FROM decision) <> 'none'
	`, shipmentID, businessID, businessID, businessID, shipmentID).Error
}

func (r *Repository) GetShipmentStatsByGeozone(ctx context.Context, filter domain.ShipmentStatsFilter) ([]domain.ShipmentStatsByGeozone, error) {
	where := []string{"s.deleted_at IS NULL", "s.destination_geozone_id IS NOT NULL"}
	args := []any{}

	where = append(where, "EXISTS (SELECT 1 FROM orders o WHERE o.id = s.order_id AND o.business_id = ?)")
	args = append(args, filter.BusinessID)

	if filter.Carrier != "" {
		where = append(where, "s.carrier = ?")
		args = append(args, filter.Carrier)
	}
	if filter.From != nil {
		where = append(where, "s.created_at >= ?")
		args = append(args, *filter.From)
	}
	if filter.To != nil {
		where = append(where, "s.created_at <= ?")
		args = append(args, *filter.To)
	}

	typeFilter := ""
	if filter.Type != "" {
		typeFilter = " AND g.type = ?"
		args = append(args, filter.Type)
	}

	limit := filter.Limit
	if limit <= 0 || limit > 1000 {
		limit = 100
	}
	args = append(args, limit)

	query := `
		SELECT g.id, g.type, g.code, g.name, g.parent_id,
		       COUNT(*) AS total,
		       COUNT(*) FILTER (WHERE s.status = 'delivered') AS delivered,
		       COUNT(*) FILTER (WHERE s.status = 'cancelled') AS cancelled,
		       COUNT(*) FILTER (WHERE s.status NOT IN ('delivered','cancelled')) AS in_transit,
		       CASE WHEN COUNT(*) > 0
		            THEN ROUND(100.0 * COUNT(*) FILTER (WHERE s.status = 'delivered') / COUNT(*), 2)
		            ELSE 0 END AS success_rate
		FROM shipments s
		JOIN geozones g ON g.id = s.destination_geozone_id
		WHERE ` + strings.Join(where, " AND ") + typeFilter + `
		GROUP BY g.id, g.type, g.code, g.name, g.parent_id
		ORDER BY total DESC
		LIMIT ?
	`

	type row struct {
		ID          uint
		Type        string
		Code        *string
		Name        string
		ParentID    *uint
		Total       int64
		Delivered   int64
		Cancelled   int64
		InTransit   int64
		SuccessRate float64
	}
	var rows []row
	if err := r.db.Conn(ctx).Raw(query, args...).Scan(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.ShipmentStatsByGeozone, len(rows))
	for i, x := range rows {
		out[i] = domain.ShipmentStatsByGeozone{
			GeozoneID:   x.ID,
			Type:        x.Type,
			Code:        x.Code,
			Name:        x.Name,
			ParentID:    x.ParentID,
			Total:       x.Total,
			Delivered:   x.Delivered,
			Cancelled:   x.Cancelled,
			InTransit:   x.InTransit,
			SuccessRate: x.SuccessRate,
		}
	}
	return out, nil
}
