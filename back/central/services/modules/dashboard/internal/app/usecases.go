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

	stats := &domain.DashboardStats{
		TotalOrders:             totalOrders,
		OrdersByIntegrationType: ordersByIntegrationType,
		TopCustomers:            topCustomers,
		OrdersByLocation:        ordersByLocation,
	}

	return stats, nil
}
