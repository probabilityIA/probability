package usecases

import (
	"github.com/secamc93/probability/back/central/services/modules/orders/domain"
)

// UseCases contiene todos los casos de uso del módulo orders
type UseCases struct {
	repo domain.IRepository
}

// New crea una nueva instancia de UseCases
func New(repo domain.IRepository) *UseCases {
	return &UseCases{
		repo: repo,
	}
}
