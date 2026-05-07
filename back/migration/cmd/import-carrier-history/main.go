package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/secamc93/probability/back/migration/shared/db"
	"github.com/secamc93/probability/back/migration/shared/env"
	"github.com/secamc93/probability/back/migration/shared/log"
)

type bucket struct {
	period      string
	geozoneID   uint
	geozoneType string
	carrier     string
	delivered   int
	returned    int
	cancelled   int
	failed      int
	inTransit   int
}

func bucketKey(period string, geozoneID uint, carrier string) string {
	return fmt.Sprintf("%s|%d|%s", period, geozoneID, carrier)
}

func normalizeCarrier(c string) string {
	c = strings.TrimSpace(c)
	if c == "" {
		return ""
	}
	upper := strings.ToUpper(c)
	var b strings.Builder
	for _, r := range upper {
		if (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func mapStatus(rawStatus string) string {
	switch strings.ToLower(strings.TrimSpace(rawStatus)) {
	case "delivered", "completed":
		return "delivered"
	case "returned":
		return "returned"
	case "cancelled":
		return "cancelled"
	case "shipping_new":
		return "failed"
	case "packing", "ready_to_ship", "sent":
		return "in_transit"
	}
	return ""
}

func main() {
	csvPath := flag.String("csv", "/tmp/carrier_history.csv", "path to csv file")
	defaultPeriod := flag.String("default-period", "2024-12-01", "period for rows without date (YYYY-MM-01)")
	flag.Parse()

	logger := log.New()
	cfg := env.New(logger)
	database := db.New(logger, cfg)
	defer database.Close()
	ctx := context.Background()
	conn := database.Conn(ctx)

	logger.Info(ctx).Str("path", *csvPath).Msg("Loading carrier history CSV")

	type geozoneRow struct {
		ID   uint
		Name string
		Type string
	}
	stateByName := map[string]geozoneRow{}
	cityByName := map[string]geozoneRow{}
	{
		var rows []geozoneRow
		if err := conn.Raw(`
			SELECT id, name, type FROM geozones
			WHERE deleted_at IS NULL AND business_id = 0 AND type IN ('state','city')
		`).Scan(&rows).Error; err != nil {
			logger.Fatal(ctx).Err(err).Msg("failed to load geozones")
		}
		toKey := func(s string) string {
			s = strings.ToLower(strings.TrimSpace(s))
			b := strings.Builder{}
			for _, r := range s {
				switch r {
				case 'á':
					b.WriteRune('a')
				case 'é':
					b.WriteRune('e')
				case 'í':
					b.WriteRune('i')
				case 'ó':
					b.WriteRune('o')
				case 'ú':
					b.WriteRune('u')
				case 'ñ':
					b.WriteRune('n')
				default:
					b.WriteRune(r)
				}
			}
			return b.String()
		}
		for _, r := range rows {
			k := toKey(r.Name)
			if r.Type == "state" {
				stateByName[k] = r
			} else if r.Type == "city" {
				cityByName[k] = r
			}
		}
		logger.Info(ctx).Int("states", len(stateByName)).Int("cities", len(cityByName)).Msg("Loaded geozones")
	}

	toKey := func(s string) string {
		s = strings.ToLower(strings.TrimSpace(s))
		b := strings.Builder{}
		for _, r := range s {
			switch r {
			case 'á':
				b.WriteRune('a')
			case 'é':
				b.WriteRune('e')
			case 'í':
				b.WriteRune('i')
			case 'ó':
				b.WriteRune('o')
			case 'ú':
				b.WriteRune('u')
			case 'ñ':
				b.WriteRune('n')
			default:
				b.WriteRune(r)
			}
		}
		return b.String()
	}

	f, err := os.Open(*csvPath)
	if err != nil {
		logger.Fatal(ctx).Err(err).Msg("open csv")
	}
	defer f.Close()
	rdr := csv.NewReader(f)
	rdr.FieldsPerRecord = -1
	header, err := rdr.Read()
	if err != nil {
		logger.Fatal(ctx).Err(err).Msg("read header")
	}
	col := map[string]int{}
	for i, h := range header {
		col[h] = i
	}

	buckets := map[string]*bucket{}
	skippedNoCarrier := 0
	skippedBadStatus := 0
	skippedNoZone := 0
	totalRows := 0

	for {
		row, err := rdr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Warn(ctx).Err(err).Msg("read row, skipping")
			continue
		}
		totalRows++
		carrier := normalizeCarrier(row[col["carrier"]])
		if carrier == "" {
			skippedNoCarrier++
			continue
		}
		mapped := mapStatus(row[col["status"]])
		if mapped == "" {
			skippedBadStatus++
			continue
		}
		destState := strings.TrimSpace(row[col["dest_state"]])
		destCity := strings.TrimSpace(row[col["dest_city"]])
		var gz geozoneRow
		var found bool
		if destCity != "" {
			if g, ok := cityByName[toKey(destCity)]; ok {
				gz = g
				found = true
			}
		}
		if !found && destState != "" {
			if g, ok := stateByName[toKey(destState)]; ok {
				gz = g
				found = true
			}
		}
		if !found {
			skippedNoZone++
			continue
		}

		period := *defaultPeriod
		if rawDate := strings.TrimSpace(row[col["date"]]); rawDate != "" {
			if t, err := time.Parse("2006-01-02", rawDate); err == nil {
				period = fmt.Sprintf("%04d-%02d-01", t.Year(), int(t.Month()))
			}
		}

		key := bucketKey(period, gz.ID, carrier)
		bk, ok := buckets[key]
		if !ok {
			bk = &bucket{period: period, geozoneID: gz.ID, geozoneType: gz.Type, carrier: carrier}
			buckets[key] = bk
		}
		switch mapped {
		case "delivered":
			bk.delivered++
		case "returned":
			bk.returned++
		case "cancelled":
			bk.cancelled++
		case "failed":
			bk.failed++
		case "in_transit":
			bk.inTransit++
		}
	}

	logger.Info(ctx).
		Int("total_rows", totalRows).
		Int("skipped_no_carrier", skippedNoCarrier).
		Int("skipped_bad_status", skippedBadStatus).
		Int("skipped_no_zone", skippedNoZone).
		Int("buckets", len(buckets)).
		Msg("Aggregation complete")

	upserts := 0
	failed := 0
	for _, b := range buckets {
		total := b.delivered + b.returned + b.cancelled + b.failed + b.inTransit
		if total == 0 {
			continue
		}
		err := conn.Exec(`
			INSERT INTO geozone_monthly_stats
			  (business_id, period, geozone_id, geozone_type, carrier,
			   total_shipments, delivered, cancelled, returned, in_transit, failed, computed_at)
			VALUES (0, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())
			ON CONFLICT (business_id, period, geozone_id, carrier)
			DO UPDATE SET
			  total_shipments = EXCLUDED.total_shipments,
			  delivered = EXCLUDED.delivered,
			  cancelled = EXCLUDED.cancelled,
			  returned = EXCLUDED.returned,
			  in_transit = EXCLUDED.in_transit,
			  failed = EXCLUDED.failed,
			  computed_at = NOW()
		`, b.period, b.geozoneID, b.geozoneType, b.carrier,
			total, b.delivered, b.cancelled, b.returned, b.inTransit, b.failed,
		).Error
		if err != nil {
			failed++
			if failed < 5 {
				logger.Warn(ctx).Err(err).Str("period", b.period).Uint("zone", b.geozoneID).Str("carrier", b.carrier).Msg("upsert failed")
			}
		} else {
			upserts++
		}
	}
	logger.Info(ctx).Int("upserted", upserts).Int("failed", failed).Msg("Done")
	fmt.Printf("imported buckets: %d / %d\n", upserts, len(buckets))
}
