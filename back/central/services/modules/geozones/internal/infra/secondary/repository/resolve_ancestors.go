package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/domain/entities"
)

func (r *Repository) ResolveAncestors(ctx context.Context, lat, lng float64, businessID uint) (*entities.GeozoneAncestors, error) {
	type row struct {
		ID   uint
		Type string
	}
	var rows []row
	err := r.db.Conn(ctx).Raw(`
		WITH RECURSIVE src AS (
		    SELECT ST_SetSRID(ST_MakePoint(?, ?), 4326) AS p
		),
		match AS (
		    SELECT g.id, g.type
		    FROM geozones g, src
		    WHERE g.deleted_at IS NULL
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
		    SELECT id, parent_id, type FROM geozones WHERE id = (SELECT id FROM match)
		    UNION ALL
		    SELECT g.id, g.parent_id, g.type
		    FROM geozones g JOIN chain c ON g.id = c.parent_id
		    WHERE g.deleted_at IS NULL
		)
		SELECT id, type FROM chain
	`, lng, lat, businessID).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return &entities.GeozoneAncestors{}, nil
	}
	out := &entities.GeozoneAncestors{Path: make([]uint, 0, len(rows))}
	deepestPriority := 99
	priority := func(t string) int {
		switch t {
		case "barrio":
			return 1
		case "neighborhood":
			return 2
		case "admin_district":
			return 3
		case "locality":
			return 4
		case "city":
			return 5
		case "state":
			return 6
		case "country":
			return 7
		}
		return 99
	}
	for _, r := range rows {
		id := r.ID
		out.Path = append(out.Path, id)
		switch r.Type {
		case "country":
			out.CountryID = &id
		case "state":
			out.StateID = &id
		case "city":
			out.CityID = &id
		case "admin_district":
			out.AdminDistrictID = &id
		case "locality":
			out.LocalityID = &id
		case "neighborhood":
			out.NeighborhoodID = &id
		case "barrio":
			out.BarrioID = &id
		}
		if p := priority(r.Type); p < deepestPriority {
			deepestPriority = p
			d := id
			out.DeepestID = &d
		}
	}
	return out, nil
}
