package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/shipping_margins/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/shipping_margins/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/shipping_margins/internal/domain/ports"
)

type IUseCase interface {
	Create(ctx context.Context, dto dtos.CreateShippingMarginDTO) (*entities.ShippingMargin, error)
	Get(ctx context.Context, businessID, id uint) (*entities.ShippingMargin, error)
	List(ctx context.Context, params dtos.ListShippingMarginsParams) ([]entities.ShippingMargin, int64, error)
	Update(ctx context.Context, dto dtos.UpdateShippingMarginDTO) (*entities.ShippingMargin, error)
	Delete(ctx context.Context, businessID, id uint) error
}

type UseCase struct {
	repo  ports.IRepository
	cache ports.ICacheWriter
}

func New(repo ports.IRepository, cache ports.ICacheWriter) IUseCase {
	return &UseCase{repo: repo, cache: cache}
}
