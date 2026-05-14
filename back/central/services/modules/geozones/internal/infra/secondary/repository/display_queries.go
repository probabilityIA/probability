package repository

import (
	"context"
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/domain/dtos"
)

func (r *Repository) GetForDisplay(ctx context.Context, params dtos.DisplayParams) ([]dtos.DisplayFeature, error) {
	geomCol := "ST_AsGeoJSON(geometry)"
	args := []any{}

	if params.Tolerance > 0 {
		geomCol = "ST_AsGeoJSON(ST_SimplifyPreserveTopology(geometry::geometry, ?))"
		args = append(args, params.Tolerance)
	}

	where := "deleted_at IS NULL AND business_id = 0 AND is_active = TRUE"
	if params.Type != "" {
		where += " AND type = ?"
		args = append(args, params.Type)
	}
	if params.Bbox != nil {
		where += " AND geometry && ST_MakeEnvelope(?, ?, ?, ?, 4326)"
		args = append(args, params.Bbox.MinLng, params.Bbox.MinLat, params.Bbox.MaxLng, params.Bbox.MaxLat)
	}
	if params.ParentID != nil {
		where += " AND parent_id = ?"
		args = append(args, *params.ParentID)
	}

	type row struct {
		ID   uint
		Type string
		Code *string
		Name string
		Geom string
	}

	query := `SELECT id, type, code, name, ` + geomCol + ` AS geom FROM geozones WHERE ` + where + ` ORDER BY type, name`

	var rows []row
	if err := r.db.Conn(ctx).Raw(query, args...).Scan(&rows).Error; err != nil {
		return nil, err
	}

	out := make([]dtos.DisplayFeature, 0, len(rows))
	for i := range rows {
		out = append(out, dtos.DisplayFeature{
			Type:     "Feature",
			Geometry: json.RawMessage(rows[i].Geom),
			Properties: dtos.DisplayFeatureProperties{
				ID:   rows[i].ID,
				Type: rows[i].Type,
				Code: rows[i].Code,
				Name: rows[i].Name,
			},
		})
	}
	return out, nil
}
