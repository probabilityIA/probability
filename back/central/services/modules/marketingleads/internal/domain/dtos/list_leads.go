package dtos

type ListLeadsParams struct {
	Page     int
	PageSize int
}

type ListLeadsResult struct {
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
}
