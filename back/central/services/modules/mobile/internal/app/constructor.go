package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/mobile/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/mobile/internal/domain/ports"
)

type IUseCase interface {
	GetOrderFull(ctx context.Context, businessID uint, orderID string) (*entities.OrderFull, error)
}

type UseCase struct {
	repo ports.IRepository
}

func New(repo ports.IRepository) IUseCase {
	return &UseCase{repo: repo}
}
