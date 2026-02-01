package ports

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/domain/entities"
)

// ═══════════════════════════════════════════════════════════════
// REPOSITORIOS (Secondary Ports - Driven Adapters)
// ═══════════════════════════════════════════════════════════════

// IProviderRepository define las operaciones de persistencia para proveedores
type IProviderRepository interface {
	Create(ctx context.Context, provider *entities.Provider) error
	GetByID(ctx context.Context, id uint) (*entities.Provider, error)
	GetByBusinessAndType(ctx context.Context, businessID uint, providerTypeCode string) (*entities.Provider, error)
	GetDefaultByBusiness(ctx context.Context, businessID uint) (*entities.Provider, error)
	List(ctx context.Context, filters *dtos.ProviderFiltersDTO) ([]*entities.Provider, error)
	Update(ctx context.Context, provider *entities.Provider) error
	Delete(ctx context.Context, id uint) error
}

// IProviderTypeRepository define las operaciones de persistencia para tipos de proveedor
type IProviderTypeRepository interface {
	GetByCode(ctx context.Context, code string) (*entities.ProviderType, error)
	List(ctx context.Context) ([]*entities.ProviderType, error)
	GetActive(ctx context.Context) ([]*entities.ProviderType, error)
}

// ═══════════════════════════════════════════════════════════════
// CLIENTE DE SOFTPYMES (Secondary Port - Driven Adapter)
// ═══════════════════════════════════════════════════════════════

// ISoftpymesClient define las operaciones con la API de Softpymes
type ISoftpymesClient interface {
	// TestAuthentication verifica que las credenciales sean válidas
	// referer: Identificación de la instancia del cliente (requerido por API)
	TestAuthentication(ctx context.Context, apiKey, apiSecret, referer string) error

	// CreateInvoice crea una factura en Softpymes
	CreateInvoice(ctx context.Context, invoiceData map[string]interface{}) error

	// CreateCreditNote crea una nota crédito en Softpymes
	CreateCreditNote(ctx context.Context, creditNoteData map[string]interface{}) error
}
