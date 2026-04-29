package usecaseshipment

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

type UseCaseShipment struct {
	repo         domain.IRepository
	marginReader domain.IShippingMarginReader
}

func New(repo domain.IRepository, marginReader domain.IShippingMarginReader) *UseCaseShipment {
	uc := &UseCaseShipment{
		repo:         repo,
		marginReader: marginReader,
	}

	uc.repo.EnsureAllBusinessesActive(context.Background())

	return uc
}

