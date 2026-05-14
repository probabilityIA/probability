package repository

import (
	"context"
	_ "embed"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/infra/secondary/metrics"
	"gorm.io/gorm"
)

//go:embed aggregate_refresh.sql
var aggregateRefreshSQL string

var nonAlnum = regexp.MustCompile(`[^a-zA-Z0-9]`)

func NormalizeCarrierKey(carrier string) string {
	return strings.ToUpper(nonAlnum.ReplaceAllString(carrier, ""))
}

func (r *Repository) RefreshAggregates(ctx context.Context) error {
	return r.db.Conn(ctx).Exec(aggregateRefreshSQL).Error
}

func (r *Repository) AncestorsByOrderID(ctx context.Context, orderID string, businessID uint) (*entities.GeozoneAncestors, error) {
	type row struct {
		DestinationGeozoneID   *uint
		GeozoneCountryID       *uint
		GeozoneStateID         *uint
		GeozoneCityID          *uint
		GeozoneAdminDistrictID *uint
		GeozoneLocalityID      *uint
		GeozoneNeighborhoodID  *uint
		GeozoneBarrioID        *uint
	}
	var rec row
	err := r.db.Conn(ctx).Raw(`
		SELECT destination_geozone_id,
		       geozone_country_id, geozone_state_id, geozone_city_id,
		       geozone_admin_district_id, geozone_locality_id,
		       geozone_neighborhood_id, geozone_barrio_id
		FROM orders
		WHERE id = ? AND (business_id = ? OR ? = 0)
		LIMIT 1
	`, orderID, businessID, businessID).Scan(&rec).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	if rec.DestinationGeozoneID == nil && rec.GeozoneCountryID == nil && rec.GeozoneStateID == nil {
		return nil, nil
	}
	return &entities.GeozoneAncestors{
		DeepestID:       rec.DestinationGeozoneID,
		CountryID:       rec.GeozoneCountryID,
		StateID:         rec.GeozoneStateID,
		CityID:          rec.GeozoneCityID,
		AdminDistrictID: rec.GeozoneAdminDistrictID,
		LocalityID:      rec.GeozoneLocalityID,
		NeighborhoodID:  rec.GeozoneNeighborhoodID,
		BarrioID:        rec.GeozoneBarrioID,
	}, nil
}

type aggregateRow struct {
	CarrierKey      string
	CarrierDisplay  string
	GeozoneLevel    *string
	GeozoneID       *uint64
	GeozoneName     *string
	Total           int64
	Delivered       int64
	Cancelled       int64
	Returned        int64
	InTransit       int64
	GlobalTotal     int64
	GlobalDelivered int64
}

func (r *Repository) ProbabilityByOrder(ctx context.Context, ancestors *entities.GeozoneAncestors) ([]dtos.ProbabilityResult, error) {
	if ancestors == nil {
		return nil, nil
	}
	type levelRef struct {
		level string
		id    *uint
	}
	refs := []levelRef{
		{"barrio", ancestors.BarrioID},
		{"neighborhood", ancestors.NeighborhoodID},
		{"admin_district", ancestors.AdminDistrictID},
		{"locality", ancestors.LocalityID},
		{"city", ancestors.CityID},
		{"state", ancestors.StateID},
		{"country", ancestors.CountryID},
	}
	levels := make([]string, 0, len(refs))
	ids := make([]uint64, 0, len(refs))
	for _, ref := range refs {
		if ref.id == nil || *ref.id == 0 {
			continue
		}
		levels = append(levels, ref.level)
		ids = append(ids, uint64(*ref.id))
	}
	q := `
WITH order_zones(geozone_level, geozone_id) AS (
    SELECT UNNEST(?::text[]), UNNEST(?::bigint[])
),
zone_stats AS (
    SELECT gcs.*
    FROM geozone_carrier_stats gcs
    JOIN order_zones oz USING (geozone_level, geozone_id)
    WHERE gcs.total > 0
),
deepest AS (
    SELECT DISTINCT ON (carrier_key) zs.*
    FROM zone_stats zs
    ORDER BY carrier_key,
             array_position(ARRAY['barrio','neighborhood','admin_district','locality','city','state','country'], geozone_level)
)
SELECT
    g.carrier_key,
    g.carrier_display,
    d.geozone_level,
    d.geozone_id,
    gz.name AS geozone_name,
    COALESCE(d.total, 0)      AS total,
    COALESCE(d.delivered, 0)  AS delivered,
    COALESCE(d.cancelled, 0)  AS cancelled,
    COALESCE(d.returned, 0)   AS returned,
    COALESCE(d.in_transit, 0) AS in_transit,
    g.total     AS global_total,
    g.delivered AS global_delivered
FROM geozone_carrier_stats g
LEFT JOIN deepest d ON d.carrier_key = g.carrier_key
LEFT JOIN geozones gz ON gz.id = d.geozone_id
WHERE g.geozone_level = '__global__' AND g.total > 0
ORDER BY g.carrier_display
`
	var rows []aggregateRow
	if len(levels) == 0 {
		levels = []string{""}
		ids = []uint64{0}
	}
	start := time.Now()
	defer func() { metrics.ProbabilityQueryDuration.Observe(time.Since(start).Seconds()) }()
	if err := r.db.Conn(ctx).Raw(q, levels, ids).Scan(&rows).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	out := make([]dtos.ProbabilityResult, 0, len(rows))
	for _, row := range rows {
		res := dtos.ProbabilityResult{
			Carrier: row.CarrierDisplay,
		}
		if row.GlobalTotal > 0 {
			gr := float64(row.GlobalDelivered) / float64(row.GlobalTotal)
			res.GlobalRate = &gr
			res.GlobalTotal = row.GlobalTotal
		}
		if row.GeozoneLevel != nil && row.GeozoneID != nil && row.Total > 0 {
			rate := float64(row.Delivered) / float64(row.Total)
			res.Found = true
			res.DeliveryRate = &rate
			res.Level = *row.GeozoneLevel
			if row.Total < 5 {
				res.IsEstimated = true
				res.EstimateSource = "zone_low_sample"
			}
			name := ""
			if row.GeozoneName != nil {
				name = *row.GeozoneName
			}
			res.Stats = &dtos.ProbabilityLevelStats{
				GeozoneID:   uint(*row.GeozoneID),
				GeozoneType: *row.GeozoneLevel,
				GeozoneName: name,
				Total:       row.Total,
				Delivered:   row.Delivered,
				Cancelled:   row.Cancelled,
				Returned:    row.Returned,
				InTransit:   row.InTransit,
			}
		}
		out = append(out, res)
	}
	return out, nil
}

func (r *Repository) ProbabilityForCarrier(ctx context.Context, ancestors *entities.GeozoneAncestors, carrierKey string) (*dtos.ProbabilityResult, error) {
	if ancestors == nil {
		return nil, nil
	}
	all, err := r.ProbabilityByOrder(ctx, ancestors)
	if err != nil {
		return nil, err
	}
	for i := range all {
		if NormalizeCarrierKey(all[i].Carrier) == carrierKey {
			return &all[i], nil
		}
	}
	return nil, nil
}

var _ ports.IProbabilityRepository = (*Repository)(nil)
