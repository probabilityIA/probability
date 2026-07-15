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
	Matched           int
	OnlyInProbability []ProductBrief
	OnlyInSiigo       []ProductBrief
	ProbabilityNoSKU  int
	SiigoNoSKU        int
}
