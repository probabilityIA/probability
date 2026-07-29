package dtos

import "github.com/secamc93/probability/back/central/shared/productmatch"

type ProductForSync struct {
	ID      string
	SKU     string
	Barcode string
	Name    string
}

func (p ProductForSync) MatchItem() productmatch.Item {
	return productmatch.Item{
		SKU:     p.SKU,
		Barcode: p.Barcode,
		Name:    p.Name,
	}
}

type ProductBrief struct {
	SKU          string
	Name         string
	MatchedBy    string
	MatchedValue string
}

type ReconcileResult struct {
	Matched              int
	MatchedNotAssociated []ProductBrief
	OnlyInProbability    []ProductBrief
	OnlyInSiigo          []ProductBrief
	ProbabilityNoSKU     int
	SiigoNoSKU           int
	MatchRules           []productmatch.Rule
}
