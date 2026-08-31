package request

type SaveProgress struct {
	TourKey   string `json:"tour_key" binding:"required"`
	Version   int    `json:"version" binding:"required"`
	Status    string `json:"status" binding:"required"`
	StepIndex int    `json:"step_index"`
}

type SkipAllTour struct {
	TourKey string `json:"tour_key" binding:"required"`
	Version int    `json:"version" binding:"required"`
}

type SkipAll struct {
	Tours []SkipAllTour `json:"tours" binding:"required,min=1"`
}
