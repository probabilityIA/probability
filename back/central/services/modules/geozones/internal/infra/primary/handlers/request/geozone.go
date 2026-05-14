package request

import "encoding/json"

type CreateGeozoneRequest struct {
	ParentID   *uint           `json:"parent_id"`
	Type       string          `json:"type" binding:"required"`
	Code       *string         `json:"code"`
	Name       string          `json:"name" binding:"required"`
	Geometry   json.RawMessage `json:"geometry" binding:"required"`
	Properties json.RawMessage `json:"properties"`
}

type BulkFeatureProperties struct {
	Type       string  `json:"type"`
	Code       *string `json:"code"`
	Name       string  `json:"name"`
	ParentCode *string `json:"parent_code"`
}

type BulkFeature struct {
	Type       string                `json:"type"`
	Geometry   json.RawMessage       `json:"geometry"`
	Properties BulkFeatureProperties `json:"properties"`
}

type BulkImportRequest struct {
	Type     string        `json:"type"`
	Features []BulkFeature `json:"features"`
}
