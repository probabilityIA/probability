package domain

import (
	"context"
	"time"
)

// IRepository define la interfaz del repositorio para obtener estadísticas
type IRepository interface {
	// GetTotalOrders obtiene el total de órdenes
	// Si businessID es nil o 0, retorna todas las órdenes (super user)
	// Si businessID > 0, filtra por ese negocio
	GetTotalOrders(ctx context.Context, businessID *uint, integrationID *uint, startDate *time.Time, endDate *time.Time) (int64, error)

	// GetOrdersToday obtiene el total de órdenes creadas hoy
	GetOrdersToday(ctx context.Context, businessID *uint, integrationID *uint, startDate *time.Time, endDate *time.Time) (int64, error)

	// GetOrdersByIntegrationType obtiene el conteo de órdenes agrupado por tipo de integración
	GetOrdersByIntegrationType(ctx context.Context, businessID *uint, integrationID *uint, startDate *time.Time, endDate *time.Time) ([]OrderCountByIntegrationType, error)

	// GetTopCustomers obtiene los top N clientes por número de órdenes
	GetTopCustomers(ctx context.Context, businessID *uint, integrationID *uint, limit int, startDate *time.Time, endDate *time.Time) ([]TopCustomer, error)

	// GetOrdersByLocation obtiene el conteo de órdenes agrupado por ubicación (ciudad/estado)
	GetOrdersByLocation(ctx context.Context, businessID *uint, integrationID *uint, limit int, startDate *time.Time, endDate *time.Time) ([]OrderCountByLocation, error)

	// Transportadores
	GetTopDrivers(ctx context.Context, businessID *uint, integrationID *uint, limit int, startDate *time.Time, endDate *time.Time) ([]TopDriver, error)
	GetDriversByLocation(ctx context.Context, businessID *uint, integrationID *uint, limit int, startDate *time.Time, endDate *time.Time) ([]DriverByLocation, error)

	// Productos
	GetTopProducts(ctx context.Context, businessID *uint, integrationID *uint, limit int, startDate *time.Time, endDate *time.Time) ([]TopProduct, error)
	GetProductsByCategory(ctx context.Context, businessID *uint, integrationID *uint, startDate *time.Time, endDate *time.Time) ([]ProductByCategory, error)
	GetProductsByBrand(ctx context.Context, businessID *uint, integrationID *uint, startDate *time.Time, endDate *time.Time) ([]ProductByBrand, error)

	// Envíos
	GetShipmentsByStatus(ctx context.Context, businessID *uint, integrationID *uint, startDate *time.Time, endDate *time.Time) ([]ShipmentsByStatus, error)
	GetShipmentsByStatusFiltered(ctx context.Context, businessID *uint, integrationID *uint, startDate *time.Time, endDate *time.Time) ([]ShipmentsByStatus, error)
	GetShipmentsByCarrier(ctx context.Context, businessID *uint, integrationID *uint, startDate *time.Time, endDate *time.Time) ([]ShipmentsByCarrier, error)
	GetShipmentsByCarrierToday(ctx context.Context, businessID *uint, integrationID *uint, startDate *time.Time, endDate *time.Time) ([]ShipmentsByCarrier, error)
	GetShipmentsByWarehouse(ctx context.Context, businessID *uint, integrationID *uint, limit int, startDate *time.Time, endDate *time.Time) ([]ShipmentsByWarehouse, error)
	GetShipmentsByDayOfWeek(ctx context.Context, businessID *uint, integrationID *uint, weekStartDate *time.Time) ([]ShipmentsByDayOfWeek, error)
	GetOrdersByDepartment(ctx context.Context, businessID *uint, integrationID *uint, startDate *time.Time, endDate *time.Time) ([]OrdersByDepartment, error)

	// Businesses (solo super admin, businessID debe ser nil)
	GetOrdersByBusiness(ctx context.Context, limit int, startDate *time.Time, endDate *time.Time) ([]OrdersByBusiness, error)

	// Órdenes por mes
	GetOrdersByMonth(ctx context.Context, businessID *uint, integrationID *uint, startDate *time.Time, endDate *time.Time) ([]OrdersByMonth, error)

	// Órdenes por semana (últimas 12 semanas)
	GetOrdersByWeek(ctx context.Context, businessID *uint, integrationID *uint, startDate *time.Time, endDate *time.Time) ([]OrdersByWeek, error)

	// TOP 5 días de mayor demanda (fechas específicas con más órdenes en toda la historia)
	GetTopSellingDays(ctx context.Context, businessID *uint, integrationID *uint, limit int) ([]TopSellingDay, error)
}

// IUseCase define la interfaz del caso de uso para obtener estadísticas del dashboard
type IUseCase interface {
	GetDashboardStats(ctx context.Context, businessID *uint, integrationID *uint, weekStartDate *time.Time, startDate *time.Time, endDate *time.Time) (*DashboardStats, error)
	GetTopSellingDays(ctx context.Context, businessID *uint, integrationID *uint, limit int) ([]TopSellingDay, error)
}
