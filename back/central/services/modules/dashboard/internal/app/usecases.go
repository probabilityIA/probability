package app

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/dashboard/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
	"golang.org/x/sync/errgroup"
)

const maxParallelQueries = 6

type UseCase struct {
	repo   domain.IRepository
	cache  domain.IStatsCache
	logger log.ILogger
}

func New(repo domain.IRepository, cache domain.IStatsCache, logger log.ILogger) domain.IUseCase {
	return &UseCase{
		repo:   repo,
		cache:  cache,
		logger: logger,
	}
}

func statsCacheKey(businessID *uint, integrationID *uint, weekStartDate *time.Time, startDate *time.Time, endDate *time.Time) string {
	id := func(v *uint) string {
		if v == nil {
			return "all"
		}
		return fmt.Sprintf("%d", *v)
	}
	day := func(t *time.Time) string {
		if t == nil {
			return "-"
		}
		return t.Format("2006-01-02")
	}
	return strings.Join([]string{
		"dashboard:stats:v1",
		id(businessID),
		id(integrationID),
		day(weekStartDate),
		day(startDate),
		day(endDate),
	}, ":")
}

func (uc *UseCase) GetDashboardStats(ctx context.Context, businessID *uint, integrationID *uint, weekStartDate *time.Time, startDate *time.Time, endDate *time.Time, refresh bool) (*domain.DashboardStats, error) {
	key := statsCacheKey(businessID, integrationID, weekStartDate, startDate, endDate)

	if uc.cache != nil && !refresh {
		if cached, ok := uc.cache.Get(ctx, key); ok {
			return cached, nil
		}
	}

	stats, err := uc.buildStats(ctx, businessID, integrationID, weekStartDate, startDate, endDate)
	if err != nil {
		return nil, err
	}

	if uc.cache != nil {
		uc.cache.Set(ctx, key, stats)
	}

	return stats, nil
}

func (uc *UseCase) buildStats(ctx context.Context, businessID *uint, integrationID *uint, weekStartDate *time.Time, startDate *time.Time, endDate *time.Time) (*domain.DashboardStats, error) {
	stats := &domain.DashboardStats{
		OrdersByBusiness: []domain.OrdersByBusiness{},
	}

	g, gctx := errgroup.WithContext(ctx)
	g.SetLimit(maxParallelQueries)

	g.Go(func() error {
		v, err := uc.repo.GetTotalOrders(gctx, businessID, integrationID, startDate, endDate)
		if err != nil {
			uc.logger.Error(gctx).Err(err).Msg("Error al obtener total de ordenes")
			return err
		}
		stats.TotalOrders = v
		return nil
	})

	g.Go(func() error {
		v, err := uc.repo.GetOrdersToday(gctx, businessID, integrationID, startDate, endDate)
		if err != nil {
			uc.logger.Error(gctx).Err(err).Msg("Error al obtener ordenes de hoy")
			v = 0
		}
		stats.OrdersToday = v
		return nil
	})

	g.Go(func() error {
		v, err := uc.repo.GetOrdersByIntegrationType(gctx, businessID, integrationID, startDate, endDate)
		if err != nil {
			uc.logger.Error(gctx).Err(err).Msg("Error al obtener ordenes por tipo de integracion")
			return err
		}
		stats.OrdersByIntegrationType = v
		return nil
	})

	g.Go(func() error {
		v, err := uc.repo.GetTopCustomers(gctx, businessID, integrationID, 10, startDate, endDate)
		if err != nil {
			uc.logger.Error(gctx).Err(err).Msg("Error al obtener top clientes")
			return err
		}
		stats.TopCustomers = v
		return nil
	})

	g.Go(func() error {
		v, err := uc.repo.GetOrdersByLocation(gctx, businessID, integrationID, 10, startDate, endDate)
		if err != nil {
			uc.logger.Error(gctx).Err(err).Msg("Error al obtener ordenes por ubicacion")
			return err
		}
		stats.OrdersByLocation = v
		return nil
	})

	g.Go(func() error {
		v, err := uc.repo.GetTopDrivers(gctx, businessID, integrationID, 10, startDate, endDate)
		if err != nil {
			uc.logger.Error(gctx).Err(err).Msg("Error al obtener top transportadores")
			return err
		}
		stats.TopDrivers = v
		return nil
	})

	g.Go(func() error {
		v, err := uc.repo.GetDriversByLocation(gctx, businessID, integrationID, 10, startDate, endDate)
		if err != nil {
			uc.logger.Error(gctx).Err(err).Msg("Error al obtener transportadores por ubicacion")
			return err
		}
		stats.DriversByLocation = v
		return nil
	})

	g.Go(func() error {
		v, err := uc.repo.GetTopProducts(gctx, businessID, integrationID, 10, startDate, endDate)
		if err != nil {
			uc.logger.Error(gctx).Err(err).Msg("Error al obtener top productos")
			return err
		}
		stats.TopProducts = v
		return nil
	})

	g.Go(func() error {
		v, err := uc.repo.GetProductsByCategory(gctx, businessID, integrationID, startDate, endDate)
		if err != nil {
			uc.logger.Error(gctx).Err(err).Msg("Error al obtener productos por categoria")
			return err
		}
		stats.ProductsByCategory = v
		return nil
	})

	g.Go(func() error {
		v, err := uc.repo.GetProductsByBrand(gctx, businessID, integrationID, startDate, endDate)
		if err != nil {
			uc.logger.Error(gctx).Err(err).Msg("Error al obtener productos por marca")
			return err
		}
		stats.ProductsByBrand = v
		return nil
	})

	g.Go(func() error {
		v, err := uc.repo.GetShipmentsByStatusFiltered(gctx, businessID, integrationID, startDate, endDate)
		if err != nil {
			uc.logger.Error(gctx).Err(err).Msg("Error al obtener envios por estado")
			return err
		}
		stats.ShipmentsByStatus = v
		return nil
	})

	g.Go(func() error {
		v, err := uc.repo.GetShipmentsByCarrier(gctx, businessID, integrationID, startDate, endDate)
		if err != nil {
			uc.logger.Error(gctx).Err(err).Msg("Error al obtener envios por transportadora")
			return err
		}
		stats.ShipmentsByCarrier = v
		return nil
	})

	g.Go(func() error {
		v, err := uc.repo.GetShipmentsByCarrierToday(gctx, businessID, integrationID, startDate, endDate)
		if err != nil {
			uc.logger.Error(gctx).Err(err).Msg("Error al obtener envios por transportadora de hoy")
			v = []domain.ShipmentsByCarrier{}
		}
		stats.ShipmentsByCarrierToday = v
		return nil
	})

	g.Go(func() error {
		v, err := uc.repo.GetShipmentsByWarehouse(gctx, businessID, integrationID, 10, startDate, endDate)
		if err != nil {
			uc.logger.Error(gctx).Err(err).Msg("Error al obtener envios por bodega")
			return err
		}
		stats.ShipmentsByWarehouse = v
		return nil
	})

	g.Go(func() error {
		v, err := uc.repo.GetShipmentsByDayOfWeek(gctx, businessID, integrationID, weekStartDate)
		if err != nil {
			uc.logger.Error(gctx).Err(err).Msg("Error al obtener envios por dia de la semana")
			v = []domain.ShipmentsByDayOfWeek{}
		}
		stats.ShipmentsByDayOfWeek = v
		return nil
	})

	g.Go(func() error {
		v, err := uc.repo.GetOrdersByDepartment(gctx, businessID, integrationID, startDate, endDate)
		if err != nil {
			uc.logger.Error(gctx).Err(err).Msg("Error al obtener ordenes por departamento")
			v = []domain.OrdersByDepartment{}
		}
		stats.OrdersByDepartment = v
		return nil
	})

	g.Go(func() error {
		v, err := uc.repo.GetOrdersByMonth(gctx, businessID, integrationID, startDate, endDate)
		if err != nil {
			uc.logger.Error(gctx).Err(err).Msg("Error al obtener ordenes por mes")
			v = []domain.OrdersByMonth{}
		}
		stats.OrdersByMonth = v
		return nil
	})

	g.Go(func() error {
		v, err := uc.repo.GetOrdersByWeek(gctx, businessID, integrationID, startDate, endDate)
		if err != nil {
			uc.logger.Error(gctx).Err(err).Msg("Error al obtener ordenes por semana")
			v = []domain.OrdersByWeek{}
		}
		stats.OrdersByWeek = v
		return nil
	})

	if businessID == nil {
		g.Go(func() error {
			v, err := uc.repo.GetOrdersByBusiness(gctx, 10, startDate, endDate)
			if err != nil {
				uc.logger.Error(gctx).Err(err).Msg("Error al obtener ordenes por business")
				return nil
			}
			stats.OrdersByBusiness = v
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return stats, nil
}

func (uc *UseCase) GetTopSellingDays(ctx context.Context, businessID *uint, integrationID *uint, limit int) ([]domain.TopSellingDay, error) {
	topDays, err := uc.repo.GetTopSellingDays(ctx, businessID, integrationID, limit)
	if err != nil {
		uc.logger.Error(ctx).Err(err).Msg("Error al obtener TOP dias de mayor demanda")
		return nil, err
	}

	return topDays, nil
}
