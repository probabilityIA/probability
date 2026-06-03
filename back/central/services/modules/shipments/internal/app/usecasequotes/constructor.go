package usecasequotes

import (
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

type UseCaseQuotes struct {
	repo domain.IRepository
}

func New(repo domain.IRepository) *UseCaseQuotes {
	return &UseCaseQuotes{repo: repo}
}
