package repository

import (
	"context"
	"strings"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

func (r *Repository) ResolveShipmentGeozone(ctx context.Context, shipmentID uint, businessID uint) error {
	return r.db.Conn(ctx).Exec(`
		WITH RECURSIVE src AS (
		    SELECT s.id AS shipment_id,
		           CASE
		               WHEN o.shipping_lat IS NOT NULL AND o.shipping_lng IS NOT NULL
		                   THEN ST_SetSRID(ST_MakePoint(o.shipping_lng, o.shipping_lat), 4326)
		           END AS p
		    FROM shipments s
		    LEFT JOIN orders o ON o.id = s.order_id
		    WHERE s.id = ?
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
		UPDATE shipments
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
	`, shipmentID, businessID, shipmentID).Error
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
