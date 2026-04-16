package ports

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/drivers/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/drivers/internal/domain/entities"
)

type IRepository interface {
	Create(ctx context.Context, driver *entities.Driver) (*entities.Driver, error)
	GetByID(ctx context.Context, businessID, driverID uint) (*entities.Driver, error)
	List(ctx context.Context, params dtos.ListDriversParams) ([]entities.Driver, int64, error)
	Update(ctx context.Context, driver *entities.Driver) (*entities.Driver, error)
	Delete(ctx context.Context, businessID, driverID uint) error
	ExistsByIdentification(ctx context.Context, businessID uint, identification string, excludeID *uint) (bool, error)
}
