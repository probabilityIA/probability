package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/ports"
)

// IUseCase define los casos de uso del módulo clients
type IUseCase interface {
	CreateClient(ctx context.Context, dto dtos.CreateClientDTO) (*entities.Client, error)
	GetClient(ctx context.Context, businessID, clientID uint) (*entities.Client, error)
	ListClients(ctx context.Context, params dtos.ListClientsParams) ([]entities.Client, int64, error)
	UpdateClient(ctx context.Context, dto dtos.UpdateClientDTO) (*entities.Client, error)
	DeleteClient(ctx context.Context, businessID, clientID uint) error
}

// UseCase implementa IUseCase
type UseCase struct {
	repo ports.IRepository
}

// New crea una nueva instancia del use case
func New(repo ports.IRepository) IUseCase {
	return &UseCase{repo: repo}
}
