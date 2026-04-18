package response

type CubingCheckResult struct {
	Fits            bool    `json:"fits"`
	Reason          string  `json:"reason,omitempty"`
	WeightNeededKg  float64 `json:"weight_needed_kg"`
	WeightMaxKg     float64 `json:"weight_max_kg"`
	VolumeNeededCm3 float64 `json:"volume_needed_cm3"`
	VolumeMaxCm3    float64 `json:"volume_max_cm3"`
	OccupiedQty     int     `json:"occupied_qty"`
}
