package dtos

type ListClientsParams struct {
	BusinessID uint
	Search     string
	Email      string
	Dni        string
	Name       string
	Page       int
	PageSize   int
}

func (p ListClientsParams) Offset() int {
	if p.Page < 1 {
		p.Page = 1
	}
	return (p.Page - 1) * p.PageSize
}

type CreateClientDTO struct {
	BusinessID uint
	Name       string
	Email      *string
	Phone      string
	Dni        *string
	Address    *string
	City       *string
	Notes      *string
}

type UpdateClientDTO struct {
	ID         uint
	BusinessID uint
	Name       string
	Email      *string
	Phone      string
	Dni        *string
	Address    *string
	City       *string
	Notes      *string
}

type BulkClientRowResult struct {
	Row     int
	Name    string
	Success bool
	Error   string
}

type BulkClientResultDTO struct {
	TotalRows    int
	SuccessCount int
	FailedCount  int
	Results      []BulkClientRowResult
}
