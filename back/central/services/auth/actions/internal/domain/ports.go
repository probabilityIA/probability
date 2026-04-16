package domain

import (
	"context"
)

// IRepository define las operaciones del repositorio de actions
// Esta interfaz contiene solo los métodos que se usan en los casos de uso del módulo actions
type IRepository interface {
	// Métodos de actions
	GetActions(ctx context.Context, page, pageSize int, name string) ([]Action, int64, error)
	GetActionByID(ctx context.Context, id uint) (*Action, error)
	GetActionByName(ctx context.Context, name string) (*Action, error)
	CreateAction(ctx context.Context, action Action) (uint, error)
	UpdateAction(ctx context.Context, id uint, action Action) (string, error)
	DeleteAction(ctx context.Context, id uint) (string, error)
}
