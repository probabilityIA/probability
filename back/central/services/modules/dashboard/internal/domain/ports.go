package domain

import (
	"context"
)

// IRepository define la interfaz del repositorio para obtener estadísticas
type IRepository interface {
	// GetTotalOrders obtiene el total de órdenes
	// Si businessID es nil o 0, retorna todas las órdenes (super user)
	// Si businessID > 0, filtra por ese negocio
	GetTotalOrders(ctx context.Context, businessID *uint) (int64, error)

	// GetOrdersByIntegrationType obtiene el conteo de órdenes agrupado por tipo de integración
	GetOrdersByIntegrationType(ctx context.Context, businessID *uint) ([]OrderCountByIntegrationType, error)

	// GetTopCustomers obtiene los top N clientes por número de órdenes
	GetTopCustomers(ctx context.Context, businessID *uint, limit int) ([]TopCustomer, error)

	// GetOrdersByLocation obtiene el conteo de órdenes agrupado por ubicación (ciudad/estado)
	GetOrdersByLocation(ctx context.Context, businessID *uint, limit int) ([]OrderCountByLocation, error)
}

// IUseCase define la interfaz del caso de uso para obtener estadísticas del dashboard
type IUseCase interface {
	GetDashboardStats(ctx context.Context, businessID *uint) (*DashboardStats, error)
}
