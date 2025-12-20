package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/dashboard/internal/domain"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/migration/shared/models"
)

// Repository implementa domain.IRepository
type Repository struct {
	db     db.IDatabase
	logger log.ILogger
}

// New crea una nueva instancia del repositorio
func New(database db.IDatabase, logger log.ILogger) domain.IRepository {
	return &Repository{
		db:     database,
		logger: logger,
	}
}

// GetTotalOrders obtiene el total de órdenes
func (r *Repository) GetTotalOrders(ctx context.Context, businessID *uint) (int64, error) {
	var count int64
	query := r.db.Conn(ctx).Model(&models.Order{})

	// Aplicar filtro por business_id si está especificado y no es super user
	if businessID != nil && *businessID > 0 {
		query = query.Where("business_id = ?", *businessID)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// GetOrdersByIntegrationType obtiene el conteo de órdenes agrupado por tipo de integración
func (r *Repository) GetOrdersByIntegrationType(ctx context.Context, businessID *uint) ([]domain.OrderCountByIntegrationType, error) {
	type Result struct {
		IntegrationType string `gorm:"column:integration_type"`
		Count           int64  `gorm:"column:count"`
	}

	var results []Result
	query := r.db.Conn(ctx).
		Model(&models.Order{}).
		Select("orders.integration_type, COUNT(*) as count").
		Group("orders.integration_type").
		Order("count DESC")

	// Aplicar filtro por business_id si está especificado y no es super user
	if businessID != nil && *businessID > 0 {
		query = query.Where("orders.business_id = ?", *businessID)
	}

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	// Mapear resultados
	stats := make([]domain.OrderCountByIntegrationType, len(results))
	for i, result := range results {
		stats[i] = domain.OrderCountByIntegrationType{
			IntegrationType: result.IntegrationType,
			Count:           result.Count,
		}
	}

	return stats, nil
}

// GetTopCustomers obtiene los top N clientes por número de órdenes
func (r *Repository) GetTopCustomers(ctx context.Context, businessID *uint, limit int) ([]domain.TopCustomer, error) {
	type Result struct {
		CustomerName  string `gorm:"column:customer_name"`
		CustomerEmail string `gorm:"column:customer_email"`
		OrderCount    int64  `gorm:"column:order_count"`
	}

	var results []Result
	query := r.db.Conn(ctx).
		Model(&models.Order{}).
		Select("orders.customer_name, orders.customer_email, COUNT(*) as order_count").
		Where("orders.customer_email != ''").
		Group("orders.customer_name, orders.customer_email").
		Order("order_count DESC").
		Limit(limit)

	// Aplicar filtro por business_id si está especificado y no es super user
	if businessID != nil && *businessID > 0 {
		query = query.Where("orders.business_id = ?", *businessID)
	}

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	// Mapear resultados
	customers := make([]domain.TopCustomer, len(results))
	for i, result := range results {
		customers[i] = domain.TopCustomer{
			CustomerName:  result.CustomerName,
			CustomerEmail: result.CustomerEmail,
			OrderCount:    result.OrderCount,
		}
	}

	return customers, nil
}

// GetOrdersByLocation obtiene el conteo de órdenes agrupado por ubicación
func (r *Repository) GetOrdersByLocation(ctx context.Context, businessID *uint, limit int) ([]domain.OrderCountByLocation, error) {
	type Result struct {
		City       string `gorm:"column:city"`
		State      string `gorm:"column:state"`
		OrderCount int64  `gorm:"column:order_count"`
	}

	var results []Result
	query := r.db.Conn(ctx).
		Model(&models.Order{}).
		Select("orders.shipping_city as city, orders.shipping_state as state, COUNT(*) as order_count").
		Where("orders.shipping_city != ''").
		Group("orders.shipping_city, orders.shipping_state").
		Order("order_count DESC").
		Limit(limit)

	// Aplicar filtro por business_id si está especificado y no es super user
	if businessID != nil && *businessID > 0 {
		query = query.Where("orders.business_id = ?", *businessID)
	}

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	// Mapear resultados
	locations := make([]domain.OrderCountByLocation, len(results))
	for i, result := range results {
		locations[i] = domain.OrderCountByLocation{
			City:       result.City,
			State:      result.State,
			OrderCount: result.OrderCount,
		}
	}

	return locations, nil
}
