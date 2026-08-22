package domain

import "context"

type MatrixColumn struct {
	IntegrationID uint
	Name          string
	Code          string
	ImageURL      string
	IsSales       bool
}

type MatrixCell struct {
	IntegrationID uint
	Present       bool
	SKU           string
	Barcode       string
	ExternalID    string
	VariantID     string
	SKUMatches    bool
	SKUKnown      bool
}

type MatrixRow struct {
	ProductID string
	SKU       string
	Barcode   string
	Name      string
	ImageURL  string
	Cells     []MatrixCell
}

const (
	MatrixSearchAll     = "all"
	MatrixSearchSKU     = "sku"
	MatrixSearchName    = "name"
	MatrixSearchBarcode = "barcode"
)

type MatrixQuery struct {
	BusinessID uint
	Search     string
	SearchBy   string
	PresentIn  []uint
	MissingIn  []uint
	Page       int
	PageSize   int
	All        bool
}

func (q MatrixQuery) Offset() int {
	if q.Page < 1 {
		return 0
	}
	return (q.Page - 1) * q.PageSize
}

func (q *MatrixQuery) Normalize() {
	switch q.SearchBy {
	case MatrixSearchSKU, MatrixSearchName, MatrixSearchBarcode:
	default:
		q.SearchBy = MatrixSearchAll
	}
	if q.Page < 1 {
		q.Page = 1
	}
	if q.PageSize < 1 {
		q.PageSize = 10
	}
	if q.PageSize > 100 {
		q.PageSize = 100
	}
}

type MatrixPage struct {
	Columns    []MatrixColumn
	Rows       []MatrixRow
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
}

type IMatrixRepository interface {
	MatchMatrix(ctx context.Context, query MatrixQuery) (*MatrixPage, error)
}
