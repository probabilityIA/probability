CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE IF NOT EXISTS geozones (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    business_id BIGINT NOT NULL DEFAULT 0,
    parent_id BIGINT REFERENCES geozones(id) ON UPDATE CASCADE ON DELETE SET NULL,
    type VARCHAR(32) NOT NULL,
    code VARCHAR(64),
    name VARCHAR(255) NOT NULL,
    geometry geometry(MultiPolygon, 4326) NOT NULL,
    centroid geography(Point, 4326),
    properties JSONB NOT NULL DEFAULT '{}'::jsonb,
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE INDEX IF NOT EXISTS idx_geozones_geometry ON geozones USING GIST (geometry);
CREATE INDEX IF NOT EXISTS idx_geozones_centroid ON geozones USING GIST (centroid);
CREATE INDEX IF NOT EXISTS idx_geozones_business_id ON geozones (business_id);
CREATE INDEX IF NOT EXISTS idx_geozones_parent_id ON geozones (parent_id);
CREATE INDEX IF NOT EXISTS idx_geozones_type ON geozones (type);
CREATE INDEX IF NOT EXISTS idx_geozones_code ON geozones (code);
CREATE INDEX IF NOT EXISTS idx_geozones_deleted_at ON geozones (deleted_at);

CREATE UNIQUE INDEX IF NOT EXISTS uq_geozones_business_type_code
    ON geozones (business_id, type, code)
    WHERE code IS NOT NULL AND deleted_at IS NULL;
