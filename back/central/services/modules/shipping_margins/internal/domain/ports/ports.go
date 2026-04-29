package ports

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/shipping_margins/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/shipping_margins/internal/domain/entities"
)

type IRepository interface {
	Create(ctx context.Context, m *entities.ShippingMargin) (*entities.ShippingMargin, error)
	GetByID(ctx context.Context, businessID, id uint) (*entities.ShippingMargin, error)
	GetByBusinessAndCarrier(ctx context.Context, businessID uint, carrierCode string) (*entities.ShippingMargin, error)
	List(ctx context.Context, params dtos.ListShippingMarginsParams) ([]entities.ShippingMargin, int64, error)
	Update(ctx context.Context, m *entities.ShippingMargin) (*entities.ShippingMargin, error)
	Delete(ctx context.Context, businessID, id uint) error
	ExistsByCarrier(ctx context.Context, businessID uint, carrierCode string, excludeID *uint) (bool, error)
}

type ICacheWriter interface {
	Upsert(ctx context.Context, m *entities.ShippingMargin) error
	Delete(ctx context.Context, businessID uint, carrierCode string) error
}
