package dtos

type ListReferralsParams struct {
	Page     int
	PageSize int
}

type ListReferralsResult struct {
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
}
