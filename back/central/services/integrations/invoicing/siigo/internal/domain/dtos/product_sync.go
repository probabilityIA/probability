package dtos

type ProductForSync struct {
	ID   string
	SKU  string
	Name string
}

type ProductBrief struct {
	SKU  string
	Name string
}

type ReconcileResult struct {
	Matched              int
	MatchedNotAssociated []ProductBrief
	OnlyInProbability    []ProductBrief
	OnlyInSiigo          []ProductBrief
	ProbabilityNoSKU     int
	SiigoNoSKU           int
}
