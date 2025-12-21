package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/dashboard/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

// IUseCase define la interfaz del caso de uso
type IUseCase interface {
	GetDashboardStats(ctx context.Context, businessID *uint) (*domain.DashboardStats, error)
}

// UseCase implementa la lógica de negocio para el dashboard
type UseCase struct {
	repo   domain.IRepository
	logger log.ILogger
}

// New crea una nueva instancia del caso de uso
func New(repo domain.IRepository, logger log.ILogger) IUseCase {
	return &UseCase{
		repo:   repo,
		logger: logger,
	}
}

// GetDashboardStats obtiene todas las estadísticas del dashboard
func (uc *UseCase) GetDashboardStats(ctx context.Context, businessID *uint) (*domain.DashboardStats, error) {
	// Obtener total de órdenes
	totalOrders, err := uc.repo.GetTotalOrders(ctx, businessID)
	if err != nil {
		uc.logger.Error().Err(err).Msg("Error al obtener total de órdenes")
		return nil, err
	}

	// Obtener órdenes por tipo de integración
	ordersByIntegrationType, err := uc.repo.GetOrdersByIntegrationType(ctx, businessID)
	if err != nil {
		uc.logger.Error().Err(err).Msg("Error al obtener órdenes por tipo de integración")
		return nil, err
	}

	// Obtener top clientes (top 10)
	topCustomers, err := uc.repo.GetTopCustomers(ctx, businessID, 10)
	if err != nil {
		uc.logger.Error().Err(err).Msg("Error al obtener top clientes")
		return nil, err
	}

	// Obtener órdenes por ubicación (top 10 ciudades)
	ordersByLocation, err := uc.repo.GetOrdersByLocation(ctx, businessID, 10)
	if err != nil {
		uc.logger.Error().Err(err).Msg("Error al obtener órdenes por ubicación")
		return nil, err
	}

	// Obtener top transportadores (top 10)
	topDrivers, err := uc.repo.GetTopDrivers(ctx, businessID, 10)
	if err != nil {
		uc.logger.Error().Err(err).Msg("Error al obtener top transportadores")
		return nil, err
	}

	// Obtener transportadores por ubicación (top 10)
	driversByLocation, err := uc.repo.GetDriversByLocation(ctx, businessID, 10)
	if err != nil {
		uc.logger.Error().Err(err).Msg("Error al obtener transportadores por ubicación")
		return nil, err
	}

	// Obtener top productos (top 10)
	topProducts, err := uc.repo.GetTopProducts(ctx, businessID, 10)
	if err != nil {
		uc.logger.Error().Err(err).Msg("Error al obtener top productos")
		return nil, err
	}

	// Obtener productos por categoría
	productsByCategory, err := uc.repo.GetProductsByCategory(ctx, businessID)
	if err != nil {
		uc.logger.Error().Err(err).Msg("Error al obtener productos por categoría")
		return nil, err
	}

	// Obtener productos por marca
	productsByBrand, err := uc.repo.GetProductsByBrand(ctx, businessID)
	if err != nil {
		uc.logger.Error().Err(err).Msg("Error al obtener productos por marca")
		return nil, err
	}

	// Obtener envíos por estado
	shipmentsByStatus, err := uc.repo.GetShipmentsByStatus(ctx, businessID)
	if err != nil {
		uc.logger.Error().Err(err).Msg("Error al obtener envíos por estado")
		return nil, err
	}

	// Obtener envíos por transportista
	shipmentsByCarrier, err := uc.repo.GetShipmentsByCarrier(ctx, businessID)
	if err != nil {
		uc.logger.Error().Err(err).Msg("Error al obtener envíos por transportista")
		return nil, err
	}

	// Obtener envíos por almacén (top 10)
	shipmentsByWarehouse, err := uc.repo.GetShipmentsByWarehouse(ctx, businessID, 10)
	if err != nil {
		uc.logger.Error().Err(err).Msg("Error al obtener envíos por almacén")
		return nil, err
	}

	stats := &domain.DashboardStats{
		TotalOrders:             totalOrders,
		OrdersByIntegrationType: ordersByIntegrationType,
		TopCustomers:            topCustomers,
		OrdersByLocation:        ordersByLocation,
		TopDrivers:              topDrivers,
		DriversByLocation:       driversByLocation,
		TopProducts:             topProducts,
		ProductsByCategory:      productsByCategory,
		ProductsByBrand:         productsByBrand,
		ShipmentsByStatus:       shipmentsByStatus,
		ShipmentsByCarrier:      shipmentsByCarrier,
		ShipmentsByWarehouse:    shipmentsByWarehouse,
	}

	// Obtener estadísticas de businesses solo si NO hay filtro de business aplicado (businessID == nil)
	// Esto significa que el super admin está viendo todos los businesses
	// Si está filtrando por un business específico, no mostrar esta gráfica
	if businessID == nil {
		ordersByBusiness, err := uc.repo.GetOrdersByBusiness(ctx, 10)
		if err != nil {
			uc.logger.Error().Err(err).Msg("Error al obtener órdenes por business")
			// No retornar error, solo loguear y continuar
		} else {
			stats.OrdersByBusiness = ordersByBusiness
		}
	} else {
		// Si hay un filtro de business aplicado, limpiar la lista de businesses
		stats.OrdersByBusiness = []domain.OrdersByBusiness{}
	}

	return stats, nil
}
