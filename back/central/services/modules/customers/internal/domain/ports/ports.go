package ports

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/entities"
)

// IRepository define los métodos del repositorio del módulo clients
type IRepository interface {
	Create(ctx context.Context, client *entities.Client) (*entities.Client, error)
	GetByID(ctx context.Context, businessID, clientID uint) (*entities.Client, error)
	List(ctx context.Context, params dtos.ListClientsParams) ([]entities.Client, int64, error)
	Update(ctx context.Context, client *entities.Client) (*entities.Client, error)
	Delete(ctx context.Context, businessID, clientID uint) error

	// Validaciones de unicidad
	ExistsByEmail(ctx context.Context, businessID uint, email string, excludeID *uint) (bool, error)
	ExistsByDni(ctx context.Context, businessID uint, dni string, excludeID *uint) (bool, error)

	// Stats de órdenes (replica query de orders - sin compartir repositorio)
	GetOrderStats(ctx context.Context, clientID uint) (orderCount int64, totalSpent float64, lastOrderAt *time.Time, err error)
}
