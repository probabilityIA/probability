package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/secamc93/probability/back/migration/shared/db"
	"github.com/secamc93/probability/back/migration/shared/env"
	"github.com/secamc93/probability/back/migration/shared/log"
)

type feature struct {
	Geometry   json.RawMessage        `json:"geometry"`
	Properties map[string]interface{} `json:"properties"`
}

type featureCollection struct {
	Features []feature `json:"features"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("uso: seed-corregimientos <zona_urbana.geojson>")
		os.Exit(2)
	}
	logger := log.New()
	cfg := env.New(logger)
	database := db.New(logger, cfg)
	defer database.Close()
	ctx := context.Background()

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		logger.Fatal(ctx).Err(err).Msg("read file")
	}
	var fc featureCollection
	if err := json.Unmarshal(data, &fc); err != nil {
		logger.Fatal(ctx).Err(err).Msg("parse json")
	}
	logger.Info().Int("features", len(fc.Features)).Msg("loaded")

	cityIDByCode := make(map[string]uint)
	type cityRow struct {
		ID   uint
		Code string
	}
	var cities []cityRow
	if err := database.Conn(ctx).Raw(`SELECT id, code FROM geozones WHERE business_id = 0 AND type = 'city' AND deleted_at IS NULL`).Scan(&cities).Error; err != nil {
		logger.Fatal(ctx).Err(err).Msg("load cities")
	}
	for _, c := range cities {
		cityIDByCode[c.Code] = c.ID
	}
	logger.Info().Int("cities_loaded", len(cityIDByCode)).Msg("city lookup ready")

	created, skipped, errs := 0, 0, 0
	for i, f := range fc.Features {
		code := str(f.Properties, "zu_cdivi")
		mpioCode := str(f.Properties, "mpio_cdpmp")
		name := str(f.Properties, "zu_cnmbre")

		if code == "" || name == "" || mpioCode == "" {
			skipped++
			continue
		}
		if code == mpioCode+"000" {
			skipped++
			continue
		}

		parentID, ok := cityIDByCode[mpioCode]
		if !ok {
			skipped++
			continue
		}

		var existing uint
		if err := database.Conn(ctx).Raw(`SELECT id FROM geozones WHERE business_id = 0 AND type = 'locality' AND code = ? AND deleted_at IS NULL LIMIT 1`, code).Scan(&existing).Error; err != nil {
			errs++
			continue
		}
		if existing > 0 {
			skipped++
			continue
		}

		props := fmt.Sprintf(`{"source":"DANE-MGN-2025-corregimiento","dane_code":"%s","mpio_code":"%s"}`, code, mpioCode)
		err := database.Conn(ctx).Exec(`
			INSERT INTO geozones (business_id, parent_id, type, code, name, geometry, centroid, properties, is_active)
			VALUES (
				0, ?, 'locality', ?, ?,
				ST_Multi(ST_SetSRID(ST_GeomFromGeoJSON(?), 4326))::geometry(MultiPolygon, 4326),
				ST_PointOnSurface(ST_SetSRID(ST_GeomFromGeoJSON(?), 4326))::geography,
				?::jsonb,
				TRUE
			)
		`, parentID, code, name, string(f.Geometry), string(f.Geometry), props).Error
		if err != nil {
			errs++
			if errs < 10 {
				logger.Error(ctx).Str("code", code).Str("name", name).Err(err).Msg("insert failed")
			}
			continue
		}
		created++
		if (i+1)%500 == 0 {
			logger.Info().Int("processed", i+1).Int("created", created).Msg("progress")
		}
	}

	logger.Info().Int("created", created).Int("skipped", skipped).Int("errors", errs).Msg("DONE")
}

func str(m map[string]interface{}, k string) string {
	if v, ok := m[k]; ok && v != nil {
		if s, ok := v.(string); ok {
			return s
		}
		return fmt.Sprintf("%v", v)
	}
	return ""
}
