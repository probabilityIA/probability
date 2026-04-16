package ports

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/entities"
)

type IRepository interface {
	Create(ctx context.Context, client *entities.Client) (*entities.Client, error)
	GetByID(ctx context.Context, businessID, clientID uint) (*entities.Client, error)
	List(ctx context.Context, params dtos.ListClientsParams) ([]entities.Client, int64, error)
	Update(ctx context.Context, client *entities.Client) (*entities.Client, error)
	Delete(ctx context.Context, businessID, clientID uint) error

	ExistsByEmail(ctx context.Context, businessID uint, email string, excludeID *uint) (bool, error)
	ExistsByDni(ctx context.Context, businessID uint, dni string, excludeID *uint) (bool, error)

	GetCustomerSummary(ctx context.Context, businessID, customerID uint) (*entities.CustomerSummary, error)
	ListCustomerAddresses(ctx context.Context, params dtos.ListCustomerAddressesParams) ([]entities.CustomerAddress, int64, error)
	ListCustomerProducts(ctx context.Context, params dtos.ListCustomerProductsParams) ([]entities.CustomerProductHistory, int64, error)
	ListCustomerOrderItems(ctx context.Context, params dtos.ListCustomerOrderItemsParams) ([]entities.CustomerOrderItem, int64, error)
	UpsertCustomerSummary(ctx context.Context, summary *entities.CustomerSummary) error
	UpsertCustomerAddress(ctx context.Context, address *entities.CustomerAddress) error
	UpsertCustomerProductHistory(ctx context.Context, product *entities.CustomerProductHistory) error
	UpsertCustomerOrderItem(ctx context.Context, item *entities.CustomerOrderItem) error
	UpdateOrderItemsStatus(ctx context.Context, orderID string, status string) error
	FindClientByPhone(ctx context.Context, businessID uint, phone string) (*entities.Client, error)
	FindClientByDNI(ctx context.Context, businessID uint, dni string) (*entities.Client, error)
	FindClientByEmail(ctx context.Context, businessID uint, email string) (*entities.Client, error)
	UpdateClientFields(ctx context.Context, clientID uint, updates map[string]any) error
}
