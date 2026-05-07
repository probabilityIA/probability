package repository

import (
	"context"
	"errors"

	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/domain/ports"
	"gorm.io/gorm"
)

func (r *Repository) AncestorsByOrderID(ctx context.Context, orderID string, businessID uint) (*entities.GeozoneAncestors, error) {
	type row struct {
		DestinationGeozoneID    *uint
		GeozoneCountryID        *uint
		GeozoneStateID          *uint
		GeozoneCityID           *uint
		GeozoneAdminDistrictID  *uint
		GeozoneLocalityID       *uint
		GeozoneNeighborhoodID   *uint
		GeozoneBarrioID         *uint
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

var allowedLevelColumns = map[string]bool{
	"geozone_country_id":         true,
	"geozone_state_id":           true,
	"geozone_city_id":            true,
	"geozone_admin_district_id":  true,
	"geozone_locality_id":        true,
	"geozone_neighborhood_id":    true,
	"geozone_barrio_id":          true,
}

func (r *Repository) AggregateAtLevel(ctx context.Context, businessID uint, levelColumn string, geozoneID uint, carrier string) (ports.LevelAggregate, error) {
	if !allowedLevelColumns[levelColumn] {
		return ports.LevelAggregate{}, errors.New("invalid level column")
	}
	type agg struct {
		Total     int64
		Delivered int64
		Cancelled int64
		Returned  int64
		InTransit int64
	}
	var out agg
	_ = businessID
	args := []any{geozoneID}
	liveCarrierFilter := ""
	if carrier != "" {
		liveCarrierFilter = " AND UPPER(REGEXP_REPLACE(s.carrier, '[^a-zA-Z0-9]', '', 'g')) = UPPER(REGEXP_REPLACE(?, '[^a-zA-Z0-9]', '', 'g'))"
		args = append(args, carrier)
	}
	args = append(args, geozoneID)
	monthlyCarrierFilter := ""
	if carrier != "" {
		monthlyCarrierFilter = " AND UPPER(REGEXP_REPLACE(m.carrier, '[^a-zA-Z0-9]', '', 'g')) = UPPER(REGEXP_REPLACE(?, '[^a-zA-Z0-9]', '', 'g'))"
		args = append(args, carrier)
	}
	q := `
		WITH live AS (
		    SELECT
		      COUNT(*) FILTER (WHERE s.status IN ('delivered','failed','returned')) AS total,
		      COUNT(*) FILTER (WHERE s.status = 'delivered') AS delivered,
		      COUNT(*) FILTER (WHERE s.status = 'cancelled') AS cancelled,
		      COUNT(*) FILTER (WHERE s.status = 'returned')  AS returned,
		      COUNT(*) FILTER (WHERE s.status NOT IN ('delivered','cancelled','returned','failed')) AS in_transit
		    FROM shipments s
		    WHERE s.deleted_at IS NULL
		      AND s.` + levelColumn + ` = ?` + liveCarrierFilter + `
		), historical AS (
		    SELECT
		      COALESCE(SUM(m.delivered + m.failed + m.returned), 0) AS total,
		      COALESCE(SUM(m.delivered), 0) AS delivered,
		      COALESCE(SUM(m.cancelled), 0) AS cancelled,
		      COALESCE(SUM(m.returned),  0) AS returned,
		      COALESCE(SUM(m.in_transit),0) AS in_transit
		    FROM geozone_monthly_stats m
		    WHERE m.geozone_id = ?` + monthlyCarrierFilter + `
		)
		SELECT
		  l.total + h.total AS total,
		  l.delivered + h.delivered AS delivered,
		  l.cancelled + h.cancelled AS cancelled,
		  l.returned  + h.returned  AS returned,
		  l.in_transit + h.in_transit AS in_transit
		FROM live l, historical h
		LIMIT 1
	`
	if err := r.db.Conn(ctx).Raw(q, args...).Scan(&out).Error; err != nil {
		return ports.LevelAggregate{}, err
	}
	return ports.LevelAggregate{
		Total:     out.Total,
		Delivered: out.Delivered,
		Cancelled: out.Cancelled,
		Returned:  out.Returned,
		InTransit: out.InTransit,
	}, nil
}

func (r *Repository) CarriersForBusiness(ctx context.Context, businessID uint) ([]string, error) {
	_ = businessID
	type row struct {
		Carrier string
		Total   int64
	}
	var rows []row
	err := r.db.Conn(ctx).Raw(`
		WITH live_carriers AS (
		    SELECT
		        UPPER(REGEXP_REPLACE(s.carrier, '[^a-zA-Z0-9]', '', 'g')) AS key,
		        s.carrier AS original,
		        1 AS sort_order
		    FROM shipments s
		    WHERE s.deleted_at IS NULL
		      AND s.carrier IS NOT NULL AND s.carrier <> ''
		), historical_carriers AS (
		    SELECT
		        UPPER(REGEXP_REPLACE(m.carrier, '[^a-zA-Z0-9]', '', 'g')) AS key,
		        m.carrier AS original,
		        2 AS sort_order
		    FROM geozone_monthly_stats m
		    WHERE m.carrier IS NOT NULL AND m.carrier <> ''
		), all_carriers AS (
		    SELECT key, original, sort_order FROM live_carriers
		    UNION ALL
		    SELECT key, original, sort_order FROM historical_carriers
		), filtered AS (
		    SELECT * FROM all_carriers
		    WHERE LOWER(original) NOT IN ('manual','test','gift_card','giftcard','envioclick','envio click','envio_click')
		)
		SELECT DISTINCT ON (key)
		    original AS carrier,
		    (SELECT COUNT(*) FROM filtered f2 WHERE f2.key = filtered.key) AS total
		FROM filtered
		ORDER BY key, sort_order
	`).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, len(rows))
	for _, x := range rows {
		if x.Total < 3 {
			continue
		}
		out = append(out, x.Carrier)
	}
	return out, nil
}

func (r *Repository) GlobalCarrierStats(ctx context.Context, carrier string) (int64, int64, error) {
	type row struct {
		Total     int64
		Delivered int64
	}
	var out row
	err := r.db.Conn(ctx).Raw(`
		WITH live AS (
		    SELECT
		      COUNT(*) FILTER (WHERE s.status IN ('delivered','failed','returned')) AS total,
		      COUNT(*) FILTER (WHERE s.status = 'delivered') AS delivered
		    FROM shipments s
		    WHERE s.deleted_at IS NULL
		      AND UPPER(REGEXP_REPLACE(s.carrier, '[^a-zA-Z0-9]', '', 'g')) = UPPER(REGEXP_REPLACE(?, '[^a-zA-Z0-9]', '', 'g'))
		), historical AS (
		    SELECT
		      COALESCE(SUM(m.delivered + m.failed + m.returned), 0) AS total,
		      COALESCE(SUM(m.delivered), 0) AS delivered
		    FROM geozone_monthly_stats m
		    WHERE UPPER(REGEXP_REPLACE(m.carrier, '[^a-zA-Z0-9]', '', 'g')) = UPPER(REGEXP_REPLACE(?, '[^a-zA-Z0-9]', '', 'g'))
		)
		SELECT l.total + h.total AS total, l.delivered + h.delivered AS delivered
		FROM live l, historical h
	`, carrier, carrier).Scan(&out).Error
	if err != nil {
		return 0, 0, err
	}
	return out.Delivered, out.Total, nil
}

func (r *Repository) GeozoneNameAndType(ctx context.Context, geozoneID uint) (string, string, error) {
	type row struct {
		Name string
		Type string
	}
	var rec row
	if err := r.db.Conn(ctx).Raw(
		`SELECT name, type FROM geozones WHERE id = ? LIMIT 1`, geozoneID,
	).Scan(&rec).Error; err != nil {
		return "", "", err
	}
	return rec.Name, rec.Type, nil
}
