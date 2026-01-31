package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
)

// ListProviderTypes lista todos los tipos de proveedores disponibles
func (uc *useCase) ListProviderTypes(ctx context.Context) ([]*entities.InvoicingProviderType, error) {
	return uc.providerTypeRepo.GetActive(ctx)
}
