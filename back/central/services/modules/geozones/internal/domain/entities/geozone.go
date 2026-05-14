package entities

import (
	"encoding/json"
	"time"
)

type Geozone struct {
	ID         uint
	BusinessID uint
	ParentID   *uint
	Type       string
	Code       *string
	Name       string
	Geometry   json.RawMessage
	Centroid   json.RawMessage
	Properties json.RawMessage
	IsActive   bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
