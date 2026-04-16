package ports

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/vehicles/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/vehicles/internal/domain/entities"
)

type IRepository interface {
	Create(ctx context.Context, vehicle *entities.Vehicle) (*entities.Vehicle, error)
	GetByID(ctx context.Context, businessID, vehicleID uint) (*entities.Vehicle, error)
	List(ctx context.Context, params dtos.ListVehiclesParams) ([]entities.Vehicle, int64, error)
	Update(ctx context.Context, vehicle *entities.Vehicle) (*entities.Vehicle, error)
	Delete(ctx context.Context, businessID, vehicleID uint) error
	ExistsByLicensePlate(ctx context.Context, businessID uint, plate string, excludeID *uint) (bool, error)
}
