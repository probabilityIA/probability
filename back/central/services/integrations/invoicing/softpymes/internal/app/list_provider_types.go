package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/domain/entities"
)

// ListProviderTypes lista los tipos de proveedores de facturación disponibles
func (uc *useCase) ListProviderTypes(ctx context.Context) ([]*entities.ProviderType, error) {
	return uc.providerTypeRepo.GetActive(ctx)
}
