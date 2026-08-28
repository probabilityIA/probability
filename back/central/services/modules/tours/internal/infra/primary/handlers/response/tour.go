package response

type TourProgress struct {
	TourKey     string `json:"tour_key"`
	Version     int    `json:"version"`
	Status      string `json:"status"`
	StepIndex   int    `json:"step_index"`
	CompletedAt string `json:"completed_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}
