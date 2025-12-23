package domain

import (
	"context"
)

// IRepository define la interfaz del repositorio para obtener estadísticas
type IRepository interface {
	// GetTotalOrders obtiene el total de órdenes
	// Si businessID es nil o 0, retorna todas las órdenes (super user)
	// Si businessID > 0, filtra por ese negocio
	GetTotalOrders(ctx context.Context, businessID *uint, integrationID *uint) (int64, error)

	// GetOrdersByIntegrationType obtiene el conteo de órdenes agrupado por tipo de integración
	GetOrdersByIntegrationType(ctx context.Context, businessID *uint, integrationID *uint) ([]OrderCountByIntegrationType, error)

	// GetTopCustomers obtiene los top N clientes por número de órdenes
	GetTopCustomers(ctx context.Context, businessID *uint, integrationID *uint, limit int) ([]TopCustomer, error)

	// GetOrdersByLocation obtiene el conteo de órdenes agrupado por ubicación (ciudad/estado)
	GetOrdersByLocation(ctx context.Context, businessID *uint, integrationID *uint, limit int) ([]OrderCountByLocation, error)

	// Transportadores
	GetTopDrivers(ctx context.Context, businessID *uint, integrationID *uint, limit int) ([]TopDriver, error)
	GetDriversByLocation(ctx context.Context, businessID *uint, integrationID *uint, limit int) ([]DriverByLocation, error)

	// Productos
	GetTopProducts(ctx context.Context, businessID *uint, integrationID *uint, limit int) ([]TopProduct, error)
	GetProductsByCategory(ctx context.Context, businessID *uint, integrationID *uint) ([]ProductByCategory, error)
	GetProductsByBrand(ctx context.Context, businessID *uint, integrationID *uint) ([]ProductByBrand, error)

	// Envíos
	GetShipmentsByStatus(ctx context.Context, businessID *uint, integrationID *uint) ([]ShipmentsByStatus, error)
	GetShipmentsByCarrier(ctx context.Context, businessID *uint, integrationID *uint) ([]ShipmentsByCarrier, error)
	GetShipmentsByWarehouse(ctx context.Context, businessID *uint, integrationID *uint, limit int) ([]ShipmentsByWarehouse, error)

	// Businesses (solo super admin, businessID debe ser nil)
	GetOrdersByBusiness(ctx context.Context, limit int) ([]OrdersByBusiness, error)
}

// IUseCase define la interfaz del caso de uso para obtener estadísticas del dashboard
type IUseCase interface {
	GetDashboardStats(ctx context.Context, businessID *uint, integrationID *uint) (*DashboardStats, error)
}
