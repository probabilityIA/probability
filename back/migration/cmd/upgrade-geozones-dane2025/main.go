package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

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
	if len(os.Args) < 3 {
		fmt.Println("uso: upgrade-geozones-dane2025 <dpto.geojson> <mpio.geojson>")
		os.Exit(2)
	}
	dptoPath := os.Args[1]
	mpioPath := os.Args[2]

	logger := log.New()
	cfg := env.New(logger)
	database := db.New(logger, cfg)
	defer database.Close()
	ctx := context.Background()

	dpto := mustLoad(dptoPath)
	mpio := mustLoad(mpioPath)
	logger.Info().Int("dpto", len(dpto.Features)).Int("mpio", len(mpio.Features)).Msg("Loaded DANE 2025")

	upd, miss := 0, 0
	for _, f := range dpto.Features {
		code := str(f.Properties, "dpto_ccdgo", "DPTO_CCDGO")
		if code == "" {
			miss++
			continue
		}
		if err := update(ctx, database, "state", code, f.Geometry); err != nil {
			logger.Error(ctx).Str("code", code).Err(err).Msg("update dpto")
			miss++
			continue
		}
		upd++
	}
	logger.Info().Int("updated", upd).Int("missed", miss).Msg("Departments done")

	upd, miss = 0, 0
	for _, f := range mpio.Features {
		code := str(f.Properties, "mpio_cdpmp", "MPIO_CDPMP")
		if code == "" {
			miss++
			continue
		}
		if err := update(ctx, database, "city", code, f.Geometry); err != nil {
			logger.Error(ctx).Str("code", code).Err(err).Msg("update mpio")
			miss++
			continue
		}
		upd++
	}
	logger.Info().Int("updated", upd).Int("missed", miss).Msg("Municipalities done")

	logger.Info().Msg("Recomputing country geometry from union of states...")
	res := database.Conn(ctx).Exec(`
		UPDATE geozones
		SET geometry = (SELECT ST_Multi(ST_Union(geometry))::geometry(MultiPolygon, 4326)
		                FROM geozones WHERE business_id = 0 AND type = 'state' AND deleted_at IS NULL),
		    centroid = (SELECT ST_PointOnSurface(ST_Union(geometry))::geography
		                FROM geozones WHERE business_id = 0 AND type = 'state' AND deleted_at IS NULL),
		    properties = '{"source":"DANE-MGN-2025"}'::jsonb,
		    updated_at = NOW()
		WHERE business_id = 0 AND type = 'country' AND code = 'CO'
	`)
	if res.Error != nil {
		logger.Error(ctx).Err(res.Error).Msg("recompute country")
	} else {
		logger.Info().Int64("rows", res.RowsAffected).Msg("Country recomputed")
	}
}

func update(ctx context.Context, database db.IDatabase, gtype, code string, geom json.RawMessage) error {
	props := fmt.Sprintf(`{"source":"DANE-MGN-2025","dane_code":"%s"}`, code)
	res := database.Conn(ctx).Exec(`
		UPDATE geozones
		SET geometry = ST_Multi(ST_SetSRID(ST_GeomFromGeoJSON(?), 4326))::geometry(MultiPolygon, 4326),
		    centroid = ST_PointOnSurface(ST_SetSRID(ST_GeomFromGeoJSON(?), 4326))::geography,
		    properties = ?::jsonb,
		    updated_at = NOW()
		WHERE business_id = 0 AND type = ? AND code = ? AND deleted_at IS NULL
	`, string(geom), string(geom), props, gtype, code)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("no row matched")
	}
	return nil
}

func mustLoad(path string) *featureCollection {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var fc featureCollection
	if err := json.Unmarshal(data, &fc); err != nil {
		panic(err)
	}
	return &fc
}

func str(m map[string]interface{}, keys ...string) string {
	for _, k := range keys {
		if v, ok := m[k]; ok && v != nil {
			if s, ok := v.(string); ok && strings.TrimSpace(s) != "" {
				return s
			}
		}
	}
	return ""
}
