package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/domain/entities"
)

// ListProviders lista proveedores según filtros
func (uc *useCase) ListProviders(ctx context.Context, filters *dtos.ProviderFiltersDTO) ([]*entities.Provider, error) {
	return uc.providerRepo.List(ctx, filters)
}
