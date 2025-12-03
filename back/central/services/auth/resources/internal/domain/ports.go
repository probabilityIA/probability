package domain

import (
	"context"
)

// IRepository define las operaciones del repositorio de recursos
// Esta interfaz contiene solo los métodos que se usan en los casos de uso del módulo resources
type IRepository interface {
	// Métodos de recursos
	GetResources(ctx context.Context, filters ResourceFilters) ([]Resource, int64, error)
	GetResourceByID(ctx context.Context, id uint) (*Resource, error)
	GetResourceByName(ctx context.Context, name string) (*Resource, error)
	CreateResource(ctx context.Context, resource Resource) (uint, error)
	UpdateResource(ctx context.Context, id uint, resource Resource) (string, error)
	DeleteResource(ctx context.Context, id uint) (string, error)
}
