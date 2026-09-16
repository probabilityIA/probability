package dtos

import "time"

const (
	ReviewAll           = "all"
	ReviewNoDestination = "no_destination"
	ReviewNegative      = "negative"
	ReviewPositive      = "positive"
	ReviewErrors        = "errors"
	ReviewNotClicked    = "not_clicked"
)

type ReviewFilter struct {
	BusinessID *uint
	Kind       string
	Search     string
	From       *time.Time
	To         *time.Time
	Page       int
	PageSize   int
}

type PaginatedResponse[T any] struct {
	Data       []T
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
}
