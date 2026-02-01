package app

import (
	"context"
	"errors"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
)

// ═══════════════════════════════════════════════════════════════
// MÉTODOS DEPRECADOS - Migrados a integrations/core
// ═══════════════════════════════════════════════════════════════
// NOTA: Estos métodos están deprecados y retornan errores.
// La gestión de proveedores de facturación ahora se realiza a través
// del módulo integrations/core que proporciona un catálogo centralizado
// de todas las integraciones de facturación.
//
// Para gestionar proveedores de facturación, usar:
// - integrations/core endpoints: GET /integrations?category=invoicing
// ═══════════════════════════════════════════════════════════════

var (
	// ErrProviderManagementDeprecated indica que la gestión de proveedores está deprecada
	ErrProviderManagementDeprecated = errors.New("gestión de proveedores deprecada, usar integrations/core")
)

// CreateProvider está deprecado y retorna error
// Usar integrations/core para crear integraciones de facturación
func (uc *useCase) CreateProvider(ctx context.Context, dto *dtos.CreateProviderDTO) (*entities.InvoicingProvider, error) {
	uc.log.Warn(ctx).Msg("CreateProvider is deprecated, use integrations/core instead")
	return nil, ErrProviderManagementDeprecated
}

// UpdateProvider está deprecado y retorna error
// Usar integrations/core para actualizar integraciones de facturación
func (uc *useCase) UpdateProvider(ctx context.Context, id uint, dto *dtos.UpdateProviderDTO) error {
	uc.log.Warn(ctx).Msg("UpdateProvider is deprecated, use integrations/core instead")
	return ErrProviderManagementDeprecated
}

// GetProvider está deprecado y retorna error
// Usar integrations/core para obtener integraciones de facturación
func (uc *useCase) GetProvider(ctx context.Context, id uint) (*entities.InvoicingProvider, error) {
	uc.log.Warn(ctx).Msg("GetProvider is deprecated, use integrations/core instead")
	return nil, ErrProviderManagementDeprecated
}

// ListProviders está deprecado y retorna error
// Usar integrations/core para listar integraciones de facturación
func (uc *useCase) ListProviders(ctx context.Context, businessID uint) ([]*entities.InvoicingProvider, error) {
	uc.log.Warn(ctx).Msg("ListProviders is deprecated, use integrations/core instead")
	return nil, ErrProviderManagementDeprecated
}

// TestProviderConnection está deprecado y retorna error
// Usar integrations/core para validar credenciales de integraciones
func (uc *useCase) TestProviderConnection(ctx context.Context, id uint) error {
	uc.log.Warn(ctx).Msg("TestProviderConnection is deprecated, use integrations/core instead")
	return ErrProviderManagementDeprecated
}

// ListProviderTypes está deprecado y retorna error
// Usar integrations/core para listar tipos de integraciones de facturación
func (uc *useCase) ListProviderTypes(ctx context.Context) ([]*entities.InvoicingProviderType, error) {
	uc.log.Warn(ctx).Msg("ListProviderTypes is deprecated, use integrations/core instead")
	return nil, ErrProviderManagementDeprecated
}
