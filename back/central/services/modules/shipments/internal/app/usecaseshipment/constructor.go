package usecaseshipment

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

type UseCaseShipment struct {
	repo         domain.IRepository
	marginReader domain.IShippingMarginReader
	pdfUploader  domain.IPDFUploader
}

func New(repo domain.IRepository, marginReader domain.IShippingMarginReader, pdfUploader domain.IPDFUploader) *UseCaseShipment {
	uc := &UseCaseShipment{
		repo:         repo,
		marginReader: marginReader,
		pdfUploader:  pdfUploader,
	}

	uc.repo.EnsureAllBusinessesActive(context.Background())

	return uc
}
