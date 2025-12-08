package usecaseproduct

import (
	"github.com/secamc93/probability/back/central/services/modules/products/internal/domain"
)

// UseCaseProduct contiene los casos de uso CRUD básicos de productos
type UseCaseProduct struct {
	repo domain.IRepository
}

// New crea una nueva instancia de UseCaseProduct
func New(repo domain.IRepository) *UseCaseProduct {
	return &UseCaseProduct{
		repo: repo,
	}
}

