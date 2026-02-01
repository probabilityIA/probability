package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/domain/entities"
)

// GetProvider obtiene un proveedor por ID
func (uc *useCase) GetProvider(ctx context.Context, id uint) (*entities.Provider, error) {
	return uc.providerRepo.GetByID(ctx, id)
}
