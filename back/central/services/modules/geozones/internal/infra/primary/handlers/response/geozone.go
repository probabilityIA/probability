package response

import (
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/domain/entities"
)

type GeozoneResponse struct {
	ID         uint            `json:"id"`
	BusinessID uint            `json:"business_id"`
	ParentID   *uint           `json:"parent_id"`
	Type       string          `json:"type"`
	Code       *string         `json:"code"`
	Name       string          `json:"name"`
	Geometry   json.RawMessage `json:"geometry,omitempty"`
	Centroid   json.RawMessage `json:"centroid,omitempty"`
	Properties json.RawMessage `json:"properties"`
	IsActive   bool            `json:"is_active"`
}

func FromEntity(g *entities.Geozone) GeozoneResponse {
	return GeozoneResponse{
		ID:         g.ID,
		BusinessID: g.BusinessID,
		ParentID:   g.ParentID,
		Type:       g.Type,
		Code:       g.Code,
		Name:       g.Name,
		Geometry:   g.Geometry,
		Centroid:   g.Centroid,
		Properties: g.Properties,
		IsActive:   g.IsActive,
	}
}

type ListResponse struct {
	Data       []GeozoneResponse `json:"data"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	TotalPages int               `json:"total_pages"`
}

type BulkImportResponse struct {
	Created int      `json:"created"`
	Skipped int      `json:"skipped"`
	Errors  []string `json:"errors,omitempty"`
}
