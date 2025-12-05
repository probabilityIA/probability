package domain

import (
	"context"

	"github.com/secamc93/probability/back/migration/shared/models"
)

// ───────────────────────────────────────────
//
//	REPOSITORY INTERFACE
//
// ───────────────────────────────────────────

// IRepository define todos los métodos de repositorio del módulo orders
type IRepository interface {
	// CRUD Operations
	CreateOrder(ctx context.Context, order *models.Order) error
	GetOrderByID(ctx context.Context, id string) (*models.Order, error)
	GetOrderByInternalNumber(ctx context.Context, internalNumber string) (*models.Order, error)
	ListOrders(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]models.Order, int64, error)
	UpdateOrder(ctx context.Context, order *models.Order) error
	DeleteOrder(ctx context.Context, id string) error

	// Validation
	OrderExists(ctx context.Context, externalID string, integrationID uint) (bool, error)
}
