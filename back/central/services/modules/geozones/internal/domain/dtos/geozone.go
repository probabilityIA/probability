package dtos

import "encoding/json"

const (
	TypeCountry       = "country"
	TypeState         = "state"
	TypeCity          = "city"
	TypeAdminDistrict = "admin_district"
	TypeLocality      = "locality"
	TypeNeighborhood  = "neighborhood"
	TypeBarrio        = "barrio"
	TypeCustom        = "custom"
)

type ListGeozonesParams struct {
	BusinessID  uint
	Type        string
	ParentID    *uint
	Code        string
	Search      string
	IncludeGeom bool
	Page        int
	PageSize    int
}

func (p ListGeozonesParams) Offset() int {
	return (p.Page - 1) * p.PageSize
}

type CreateGeozoneDTO struct {
	BusinessID uint
	ParentID   *uint
	Type       string
	Code       *string
	Name       string
	Geometry   json.RawMessage
	Properties json.RawMessage
}

type LookupParams struct {
	BusinessID uint
	Lat        float64
	Lng        float64
	Type       string
}

type BulkImportFeature struct {
	Type       string          `json:"type"`
	Code       *string         `json:"code"`
	Name       string          `json:"name"`
	ParentCode *string         `json:"parent_code"`
	Geometry   json.RawMessage `json:"geometry"`
	Properties json.RawMessage `json:"properties"`
}

type BulkImportDTO struct {
	BusinessID uint
	Features   []BulkImportFeature
}

type BulkImportResult struct {
	Created int
	Skipped int
	Errors  []string
}
