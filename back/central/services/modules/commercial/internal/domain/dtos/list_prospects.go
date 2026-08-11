package dtos

const (
	DefaultPage     = 1
	DefaultPageSize = 10
	MaxPageSize     = 100
)

type ListProspectsFilters struct {
	Platform string
	Status   string
	City     string
	Search   string
	MinScore *float64
	MaxScore *float64
	HasCOD   *bool
	Seen     *bool
	Page     int
	PageSize int
}

type ListResult struct {
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
}

func NormalizePagination(page, pageSize int) (int, int) {
	if page < 1 {
		page = DefaultPage
	}
	if pageSize < 1 {
		pageSize = DefaultPageSize
	}
	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}
	return page, pageSize
}

func BuildListResult(total int64, page, pageSize int) ListResult {
	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}
	if totalPages < 1 {
		totalPages = 1
	}
	return ListResult{
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}
