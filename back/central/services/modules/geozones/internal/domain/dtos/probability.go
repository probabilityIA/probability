package dtos

type ProbabilityRequest struct {
	BusinessID uint
	Lat        *float64
	Lng        *float64
	OrderID    string
	Carrier    string
}

type ProbabilityLevelStats struct {
	GeozoneID    uint   `json:"geozone_id"`
	GeozoneType  string `json:"geozone_type"`
	GeozoneName  string `json:"geozone_name,omitempty"`
	Total        int64  `json:"total"`
	Delivered    int64  `json:"delivered"`
	Cancelled    int64  `json:"cancelled"`
	Returned     int64  `json:"returned"`
	InTransit    int64  `json:"in_transit"`
}

type ProbabilityResult struct {
	Found        bool                  `json:"found"`
	DeliveryRate *float64              `json:"delivery_rate,omitempty"`
	Level        string                `json:"level,omitempty"`
	Carrier      string                `json:"carrier,omitempty"`
	Stats        *ProbabilityLevelStats `json:"stats,omitempty"`
}
