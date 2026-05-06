package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/secamc93/probability/back/migration/shared/db"
	"github.com/secamc93/probability/back/migration/shared/env"
	"github.com/secamc93/probability/back/migration/shared/log"
)

const (
	dptoURL = "https://raw.githubusercontent.com/macortesgu/MGN_2021_geojson/main/MGN2021_DPTO_web.geo.json"
	mpioURL = "https://raw.githubusercontent.com/macortesgu/MGN_2021_geojson/main/MGN2021_MPIO_web.geo.json"
)

type feature struct {
	Geometry   json.RawMessage        `json:"geometry"`
	Properties map[string]interface{} `json:"properties"`
}

type featureCollection struct {
	Features []feature `json:"features"`
}

func main() {
	logger := log.New()
	cfg := env.New(logger)
	database := db.New(logger, cfg)
	defer database.Close()

	ctx := context.Background()

	dpto := mustLoad(dptoURL, "/tmp/dane/MGN2021_DPTO.geo.json")
	mpio := mustLoad(mpioURL, "/tmp/dane/MGN2021_MPIO.geo.json")

	logger.Info().Int("departamentos", len(dpto.Features)).Int("municipios", len(mpio.Features)).Msg("Loaded DANE data")

	countryID, err := upsertCountry(ctx, database)
	if err != nil {
		logger.Fatal(ctx).Err(err).Msg("upsert country failed")
	}
	logger.Info().Uint("country_id", countryID).Msg("Country Colombia ready")

	dptoIDByCode := make(map[string]uint, len(dpto.Features))
	created, skipped := 0, 0
	for _, f := range dpto.Features {
		code := str(f.Properties, "DPTO_CCDGO")
		name := str(f.Properties, "DPTO_CNMBR")
		if code == "" || name == "" {
			skipped++
			continue
		}
		id, ok, err := upsertGeozone(ctx, database, "state", code, name, &countryID, f.Geometry)
		if err != nil {
			logger.Error(ctx).Str("code", code).Str("name", name).Err(err).Msg("upsert dpto")
			skipped++
			continue
		}
		dptoIDByCode[code] = id
		if ok {
			created++
		}
	}
	logger.Info().Int("created", created).Int("skipped", skipped).Int("total", len(dpto.Features)).Msg("Departments done")

	created, skipped = 0, 0
	for _, f := range mpio.Features {
		dptoCode := str(f.Properties, "DPTO_CCDGO")
		code := str(f.Properties, "MPIO_CDPMP")
		name := str(f.Properties, "MPIO_CNMBR")
		if code == "" || name == "" {
			skipped++
			continue
		}
		parentID, ok := dptoIDByCode[dptoCode]
		if !ok {
			skipped++
			continue
		}
		_, ok, err := upsertGeozone(ctx, database, "city", code, name, &parentID, f.Geometry)
		if err != nil {
			logger.Error(ctx).Str("code", code).Str("name", name).Err(err).Msg("upsert mpio")
			skipped++
			continue
		}
		if ok {
			created++
		}
	}
	logger.Info().Int("created", created).Int("skipped", skipped).Int("total", len(mpio.Features)).Msg("Municipalities done")
}

func upsertCountry(ctx context.Context, database db.IDatabase) (uint, error) {
	var existing uint
	err := database.Conn(ctx).Raw(`SELECT id FROM geozones WHERE business_id = 0 AND type = 'country' AND code = 'CO' AND deleted_at IS NULL LIMIT 1`).Scan(&existing).Error
	if err != nil {
		return 0, err
	}
	if existing > 0 {
		return existing, nil
	}
	var id uint
	err = database.Conn(ctx).Raw(`
		INSERT INTO geozones (business_id, type, code, name, geometry, centroid, properties, is_active)
		VALUES (
			0, 'country', 'CO', 'Colombia',
			ST_Multi(ST_SetSRID(ST_GeomFromText('POLYGON((-79 -5, -66 -5, -66 13, -79 13, -79 -5))'), 4326))::geometry(MultiPolygon, 4326),
			ST_PointOnSurface(ST_SetSRID(ST_MakePoint(-74.297, 4.571), 4326))::geography,
			'{"source":"placeholder","note":"replaced after departments loaded"}'::jsonb,
			TRUE
		)
		RETURNING id
	`).Scan(&id).Error
	return id, err
}

func upsertGeozone(ctx context.Context, database db.IDatabase, gtype, code, name string, parentID *uint, geom json.RawMessage) (uint, bool, error) {
	var existing uint
	err := database.Conn(ctx).Raw(`SELECT id FROM geozones WHERE business_id = 0 AND type = ? AND code = ? AND deleted_at IS NULL LIMIT 1`, gtype, code).Scan(&existing).Error
	if err != nil {
		return 0, false, err
	}
	if existing > 0 {
		return existing, false, nil
	}
	var id uint
	props := fmt.Sprintf(`{"source":"DANE-MGN-2021","dane_code":"%s"}`, code)
	err = database.Conn(ctx).Raw(`
		INSERT INTO geozones (business_id, parent_id, type, code, name, geometry, centroid, properties, is_active)
		VALUES (
			0, ?, ?, ?, ?,
			ST_Multi(ST_SetSRID(ST_GeomFromGeoJSON(?), 4326))::geometry(MultiPolygon, 4326),
			ST_PointOnSurface(ST_SetSRID(ST_GeomFromGeoJSON(?), 4326))::geography,
			?::jsonb,
			TRUE
		)
		RETURNING id
	`, parentID, gtype, code, name, string(geom), string(geom), props).Scan(&id).Error
	if err != nil {
		return 0, false, err
	}
	return id, true, nil
}

func mustLoad(url, cachePath string) *featureCollection {
	if data, err := os.ReadFile(cachePath); err == nil {
		return parse(data)
	}
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	_ = os.MkdirAll("/tmp/dane", 0o755)
	_ = os.WriteFile(cachePath, data, 0o644)
	return parse(data)
}

func parse(data []byte) *featureCollection {
	var fc featureCollection
	if err := json.Unmarshal(data, &fc); err != nil {
		panic(err)
	}
	return &fc
}

func str(m map[string]interface{}, k string) string {
	v, ok := m[k]
	if !ok || v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", v)
}
