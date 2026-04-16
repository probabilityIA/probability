package entities

type ScoreBreakdown struct {
	FinalScore      float64          `json:"final_score"`
	Categories      []CategoryResult `json:"categories"`
	NegativeFactors []string         `json:"negative_factors"`
}

type CategoryResult struct {
	Name          string   `json:"name"`
	Weight        float64  `json:"weight"`
	RawScore      float64  `json:"raw_score"`
	WeightedScore float64  `json:"weighted_score"`
	Factors       []string `json:"factors"`
}
