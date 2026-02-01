package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

// IUseCase define todos los casos de uso del módulo de gestión de proveedores Softpymes
type IUseCase interface {
	// Proveedores
	CreateProvider(ctx context.Context, dto *dtos.CreateProviderDTO) (*entities.Provider, error)
	GetProvider(ctx context.Context, id uint) (*entities.Provider, error)
	ListProviders(ctx context.Context, filters *dtos.ProviderFiltersDTO) ([]*entities.Provider, error)
	UpdateProvider(ctx context.Context, id uint, dto *dtos.UpdateProviderDTO) (*entities.Provider, error)
	DeleteProvider(ctx context.Context, id uint) error
	TestProviderConnection(ctx context.Context, id uint) error

	// Tipos de proveedores
	ListProviderTypes(ctx context.Context) ([]*entities.ProviderType, error)
}

// useCase implementa todos los casos de uso del módulo
type useCase struct {
	// Repositorios
	providerRepo     ports.IProviderRepository
	providerTypeRepo ports.IProviderTypeRepository

	// Cliente de Softpymes
	softpymesClient ports.ISoftpymesClient

	// Logger
	log log.ILogger
}

// New crea una nueva instancia del use case
func New(
	providerRepo ports.IProviderRepository,
	providerTypeRepo ports.IProviderTypeRepository,
	softpymesClient ports.ISoftpymesClient,
	logger log.ILogger,
) IUseCase {
	return &useCase{
		providerRepo:     providerRepo,
		providerTypeRepo: providerTypeRepo,
		softpymesClient:  softpymesClient,
		log:              logger.WithModule("softpymes.usecase"),
	}
}
