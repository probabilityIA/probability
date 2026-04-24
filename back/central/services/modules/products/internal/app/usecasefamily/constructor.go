package usecasefamily

import "github.com/secamc93/probability/back/central/services/modules/products/internal/domain"

// UseCaseFamily contiene los casos de uso CRUD de familias de producto.
type UseCaseFamily struct {
	repo domain.IRepository
}

// New crea una nueva instancia de UseCaseFamily.
func New(repo domain.IRepository) *UseCaseFamily {
	return &UseCaseFamily{repo: repo}
}
