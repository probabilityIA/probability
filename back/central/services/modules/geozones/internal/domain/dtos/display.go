package dtos

import "encoding/json"

type DisplayFeatureProperties struct {
	ID   uint    `json:"id"`
	Type string  `json:"type"`
	Code *string `json:"code,omitempty"`
	Name string  `json:"name"`
}

type DisplayFeature struct {
	Type       string                   `json:"type"`
	Geometry   json.RawMessage          `json:"geometry"`
	Properties DisplayFeatureProperties `json:"properties"`
}

type DisplayFeatureCollection struct {
	Type     string           `json:"type"`
	Features []DisplayFeature `json:"features"`
}

type DisplayParams struct {
	Type      string
	Tolerance float64
	Bbox      *Bbox
	ParentID  *uint
}

type Bbox struct {
	MinLng float64
	MinLat float64
	MaxLng float64
	MaxLat float64
}
