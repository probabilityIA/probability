package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/dashboard/internal/domain"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
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
func (r *Repository) GetTotalOrders(ctx context.Context, businessID *uint, integrationID *uint, startDate *time.Time, endDate *time.Time) (int64, error) {
	var count int64
	query := r.db.Conn(ctx).Model(&models.Order{})

	// Aplicar filtro por business_id si está especificado y no es super user
	if businessID != nil && *businessID > 0 {
		query = query.Where("business_id = ?", *businessID)
	}

	// Aplicar filtro por integration_id si está especificado
	if integrationID != nil && *integrationID > 0 {
		query = query.Where("integration_id = ?", *integrationID)
	}

	// Aplicar filtro por rango de fechas
	if startDate != nil {
		query = query.Where("orders.created_at >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("orders.created_at < ?", *endDate)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// GetOrdersToday obtiene el total de órdenes creadas hoy
func (r *Repository) GetOrdersToday(ctx context.Context, businessID *uint, integrationID *uint, startDate *time.Time, endDate *time.Time) (int64, error) {
	var count int64
	query := r.db.Conn(ctx).Model(&models.Order{}).
		Where("DATE(created_at AT TIME ZONE 'America/Bogota') = DATE(NOW() AT TIME ZONE 'America/Bogota')")

	// Aplicar filtro por business_id si está especificado
	if businessID != nil && *businessID > 0 {
		query = query.Where("business_id = ?", *businessID)
	}

	// Aplicar filtro por integration_id si está especificado
	if integrationID != nil && *integrationID > 0 {
		query = query.Where("integration_id = ?", *integrationID)
	}

	// Aplicar filtro por rango de fechas
	if startDate != nil {
		query = query.Where("orders.created_at >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("orders.created_at < ?", *endDate)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// GetOrdersByIntegrationType obtiene el conteo de órdenes agrupado por tipo de integración
func (r *Repository) GetOrdersByIntegrationType(ctx context.Context, businessID *uint, integrationID *uint, startDate *time.Time, endDate *time.Time) ([]domain.OrderCountByIntegrationType, error) {
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

	// Aplicar filtro por integration_id si está especificado
	if integrationID != nil && *integrationID > 0 {
		query = query.Where("orders.integration_id = ?", *integrationID)
	}

	// Aplicar filtro por rango de fechas
	if startDate != nil {
		query = query.Where("orders.created_at >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("orders.created_at < ?", *endDate)
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
func (r *Repository) GetTopCustomers(ctx context.Context, businessID *uint, integrationID *uint, limit int, startDate *time.Time, endDate *time.Time) ([]domain.TopCustomer, error) {
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

	// Aplicar filtro por integration_id si está especificado
	if integrationID != nil && *integrationID > 0 {
		query = query.Where("orders.integration_id = ?", *integrationID)
	}

	// Aplicar filtro por rango de fechas
	if startDate != nil {
		query = query.Where("orders.created_at >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("orders.created_at < ?", *endDate)
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

// GetOrdersByLocation obtiene el conteo de órdenes agrupado por ubicación (estado/departamento)
func (r *Repository) GetOrdersByLocation(ctx context.Context, businessID *uint, integrationID *uint, limit int, startDate *time.Time, endDate *time.Time) ([]domain.OrderCountByLocation, error) {
	type Result struct {
		City       string `gorm:"column:city"`
		State      string `gorm:"column:state"`
		OrderCount int64  `gorm:"column:order_count"`
	}

	var results []Result
	query := r.db.Conn(ctx).
		Model(&models.Order{}).
		Select("orders.shipping_state as city, '' as state, COUNT(*) as order_count").
		Where("orders.shipping_state IS NOT NULL AND orders.shipping_state != ''").
		Group("orders.shipping_state").
		Order("order_count DESC").
		Limit(limit)

	// Aplicar filtro por business_id si está especificado y no es super user
	if businessID != nil && *businessID > 0 {
		query = query.Where("orders.business_id = ?", *businessID)
	}

	// Aplicar filtro por integration_id si está especificado
	if integrationID != nil && *integrationID > 0 {
		query = query.Where("orders.integration_id = ?", *integrationID)
	}

	// Aplicar filtro por rango de fechas
	if startDate != nil {
		query = query.Where("orders.created_at >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("orders.created_at < ?", *endDate)
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

// GetTopDrivers obtiene los top N transportadores por número de órdenes
func (r *Repository) GetTopDrivers(ctx context.Context, businessID *uint, integrationID *uint, limit int, startDate *time.Time, endDate *time.Time) ([]domain.TopDriver, error) {
	type Result struct {
		DriverName string `gorm:"column:driver_name"`
		DriverID   *uint  `gorm:"column:driver_id"`
		OrderCount int64  `gorm:"column:order_count"`
	}

	var results []Result
	query := r.db.Conn(ctx).
		Model(&models.Order{}).
		Select("orders.driver_name, orders.driver_id, COUNT(*) as order_count").
		Where("orders.driver_name != '' AND orders.driver_id IS NOT NULL").
		Group("orders.driver_id, orders.driver_name").
		Order("order_count DESC").
		Limit(limit)

	// Aplicar filtro por business_id si está especificado y no es super user
	if businessID != nil && *businessID > 0 {
		query = query.Where("orders.business_id = ?", *businessID)
	}

	// Aplicar filtro por integration_id si está especificado
	if integrationID != nil && *integrationID > 0 {
		query = query.Where("orders.integration_id = ?", *integrationID)
	}

	// Aplicar filtro por rango de fechas
	if startDate != nil {
		query = query.Where("orders.created_at >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("orders.created_at < ?", *endDate)
	}

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	// Mapear resultados
	drivers := make([]domain.TopDriver, len(results))
	for i, result := range results {
		drivers[i] = domain.TopDriver{
			DriverName: result.DriverName,
			DriverID:   result.DriverID,
			OrderCount: result.OrderCount,
		}
	}

	return drivers, nil
}

// GetDriversByLocation obtiene transportadores agrupados por ubicación
func (r *Repository) GetDriversByLocation(ctx context.Context, businessID *uint, integrationID *uint, limit int, startDate *time.Time, endDate *time.Time) ([]domain.DriverByLocation, error) {
	type Result struct {
		DriverName string `gorm:"column:driver_name"`
		City       string `gorm:"column:city"`
		State      string `gorm:"column:state"`
		OrderCount int64  `gorm:"column:order_count"`
	}

	var results []Result
	query := r.db.Conn(ctx).
		Model(&models.Order{}).
		Select("orders.driver_name, orders.shipping_city as city, orders.shipping_state as state, COUNT(*) as order_count").
		Where("orders.driver_name != '' AND orders.driver_id IS NOT NULL AND orders.shipping_city != ''").
		Group("orders.driver_id, orders.driver_name, orders.shipping_city, orders.shipping_state").
		Order("order_count DESC").
		Limit(limit)

	// Aplicar filtro por business_id si está especificado y no es super user
	if businessID != nil && *businessID > 0 {
		query = query.Where("orders.business_id = ?", *businessID)
	}

	// Aplicar filtro por integration_id si está especificado
	if integrationID != nil && *integrationID > 0 {
		query = query.Where("orders.integration_id = ?", *integrationID)
	}

	// Aplicar filtro por rango de fechas
	if startDate != nil {
		query = query.Where("orders.created_at >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("orders.created_at < ?", *endDate)
	}

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	// Mapear resultados
	drivers := make([]domain.DriverByLocation, len(results))
	for i, result := range results {
		drivers[i] = domain.DriverByLocation{
			DriverName: result.DriverName,
			City:       result.City,
			State:      result.State,
			OrderCount: result.OrderCount,
		}
	}

	return drivers, nil
}

// GetTopProducts obtiene los top N productos por número de órdenes
func (r *Repository) GetTopProducts(ctx context.Context, businessID *uint, integrationID *uint, limit int, startDate *time.Time, endDate *time.Time) ([]domain.TopProduct, error) {
	type Result struct {
		ProductName string  `gorm:"column:product_name"`
		ProductID   string  `gorm:"column:product_id"`
		SKU         string  `gorm:"column:sku"`
		OrderCount  int64   `gorm:"column:order_count"`
		TotalSold   float64 `gorm:"column:total_sold"`
	}

	var results []Result
	query := r.db.Conn(ctx).
		Model(&models.OrderItem{}).
		Select("products.name as product_name, order_items.product_id, products.sku, COUNT(DISTINCT order_items.order_id) as order_count, SUM(order_items.total_price) as total_sold").
		Joins("JOIN orders ON orders.id = order_items.order_id").
		Joins("LEFT JOIN products ON products.id = order_items.product_id").
		Where("order_items.product_id IS NOT NULL AND products.id IS NOT NULL").
		Group("order_items.product_id, products.name, products.sku").
		Order("order_count DESC, total_sold DESC").
		Limit(limit)

	// Aplicar filtro por business_id si está especificado y no es super user
	if businessID != nil && *businessID > 0 {
		query = query.Where("orders.business_id = ?", *businessID)
	}

	// Aplicar filtro por integration_id si está especificado
	if integrationID != nil && *integrationID > 0 {
		query = query.Where("orders.integration_id = ?", *integrationID)
	}

	// Aplicar filtro por rango de fechas
	if startDate != nil {
		query = query.Where("orders.created_at >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("orders.created_at < ?", *endDate)
	}

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	// Mapear resultados
	products := make([]domain.TopProduct, len(results))
	for i, result := range results {
		products[i] = domain.TopProduct{
			ProductName: result.ProductName,
			ProductID:   result.ProductID,
			SKU:         result.SKU,
			OrderCount:  result.OrderCount,
			TotalSold:   result.TotalSold,
		}
	}

	return products, nil
}

// GetProductsByCategory obtiene productos agrupados por categoría
func (r *Repository) GetProductsByCategory(ctx context.Context, businessID *uint, integrationID *uint, startDate *time.Time, endDate *time.Time) ([]domain.ProductByCategory, error) {
	type Result struct {
		Category string `gorm:"column:category"`
		Count    int64  `gorm:"column:count"`
	}

	var results []Result
	var query *gorm.DB

	// Si hay filtro de integración, debemos obtener productos DE LAS ÓRDENES de esa integración (ventas)
	// Si NO hay filtro de integración, obtenemos del catálogo (inventario/existencia)
	if integrationID != nil && *integrationID > 0 {
		query = r.db.Conn(ctx).
			Model(&models.OrderItem{}).
			Select("products.category, COUNT(DISTINCT order_items.product_id) as count").
			Joins("JOIN orders ON orders.id = order_items.order_id").
			Joins("LEFT JOIN products ON products.id = order_items.product_id").
			Where("order_items.product_id IS NOT NULL AND products.id IS NOT NULL AND products.category != ''").
			Group("products.category").
			Order("count DESC")

		query = query.Where("orders.integration_id = ?", *integrationID)
	} else {
		query = r.db.Conn(ctx).
			Model(&models.Product{}).
			Select("products.category, COUNT(DISTINCT products.id) as count").
			Where("products.category != ''").
			Group("products.category").
			Order("count DESC")
	}

	// Aplicar filtro por business_id si está especificado y no es super user
	if businessID != nil && *businessID > 0 {
		query = query.Where("products.business_id = ?", *businessID)
	}

	// Aplicar filtro por rango de fechas (solo para órdenes)
	if startDate != nil || endDate != nil {
		if integrationID != nil && *integrationID > 0 {
			if startDate != nil {
				query = query.Where("orders.created_at >= ?", *startDate)
			}
			if endDate != nil {
				query = query.Where("orders.created_at < ?", *endDate)
			}
		}
	}

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	// Mapear resultados
	categories := make([]domain.ProductByCategory, len(results))
	for i, result := range results {
		categories[i] = domain.ProductByCategory{
			Category: result.Category,
			Count:    result.Count,
		}
	}

	return categories, nil
}

// GetProductsByBrand obtiene productos agrupados por marca
func (r *Repository) GetProductsByBrand(ctx context.Context, businessID *uint, integrationID *uint, startDate *time.Time, endDate *time.Time) ([]domain.ProductByBrand, error) {
	type Result struct {
		Brand string `gorm:"column:brand"`
		Count int64  `gorm:"column:count"`
	}

	var results []Result
	var query *gorm.DB

	// Misma lógica que Category: Si hay integrationID, filtrar por ventas. Si no, catálogo.
	if integrationID != nil && *integrationID > 0 {
		query = r.db.Conn(ctx).
			Model(&models.OrderItem{}).
			Select("products.brand, COUNT(DISTINCT order_items.product_id) as count").
			Joins("JOIN orders ON orders.id = order_items.order_id").
			Joins("LEFT JOIN products ON products.id = order_items.product_id").
			Where("order_items.product_id IS NOT NULL AND products.id IS NOT NULL AND products.brand != ''").
			Group("products.brand").
			Order("count DESC")

		query = query.Where("orders.integration_id = ?", *integrationID)
	} else {
		query = r.db.Conn(ctx).
			Model(&models.Product{}).
			Select("products.brand, COUNT(DISTINCT products.id) as count").
			Where("products.brand != ''").
			Group("products.brand").
			Order("count DESC")
	}

	// Aplicar filtro por business_id si está especificado y no es super user
	if businessID != nil && *businessID > 0 {
		query = query.Where("products.business_id = ?", *businessID)
	}

	// Aplicar filtro por rango de fechas (solo para órdenes)
	if startDate != nil || endDate != nil {
		if integrationID != nil && *integrationID > 0 {
			if startDate != nil {
				query = query.Where("orders.created_at >= ?", *startDate)
			}
			if endDate != nil {
				query = query.Where("orders.created_at < ?", *endDate)
			}
		}
	}

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	// Mapear resultados
	brands := make([]domain.ProductByBrand, len(results))
	for i, result := range results {
		brands[i] = domain.ProductByBrand{
			Brand: result.Brand,
			Count: result.Count,
		}
	}

	return brands, nil
}

// GetShipmentsByStatus obtiene envíos agrupados por estado
func (r *Repository) GetShipmentsByStatus(ctx context.Context, businessID *uint, integrationID *uint, startDate *time.Time, endDate *time.Time) ([]domain.ShipmentsByStatus, error) {
	type Result struct {
		Status string `gorm:"column:status"`
		Count  int64  `gorm:"column:count"`
	}

	var results []Result
	query := r.db.Conn(ctx).
		Model(&models.Shipment{}).
		Select("shipments.status, COUNT(*) as count").
		Joins("JOIN orders ON orders.id = shipments.order_id").
		Where("orders.deleted_at IS NULL").
		Group("shipments.status").
		Order("count DESC")

	// Aplicar filtro por business_id si está especificado y no es super user
	if businessID != nil && *businessID > 0 {
		query = query.Where("orders.business_id = ?", *businessID)
	}

	// Aplicar filtro por integration_id si está especificado
	if integrationID != nil && *integrationID > 0 {
		query = query.Where("orders.integration_id = ?", *integrationID)
	}

	// Aplicar filtro por rango de fechas
	if startDate != nil {
		query = query.Where("orders.created_at >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("orders.created_at < ?", *endDate)
	}

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	// Mapear resultados
	statuses := make([]domain.ShipmentsByStatus, len(results))
	for i, result := range results {
		statuses[i] = domain.ShipmentsByStatus{
			Status: result.Status,
			Count:  result.Count,
		}
	}

	return statuses, nil
}

// GetShipmentsByStatusFiltered obtiene envíos agrupados por estado (solo: pending, in_transit, delivered)
// Solo cuenta el shipment más reciente por orden para evitar duplicados
func (r *Repository) GetShipmentsByStatusFiltered(ctx context.Context, businessID *uint, integrationID *uint, startDate *time.Time, endDate *time.Time) ([]domain.ShipmentsByStatus, error) {
	type Result struct {
		Status string `gorm:"column:status"`
		Count  int64  `gorm:"column:count"`
	}

	var results []Result

	// Build the appropriate Raw SQL query based on filters
	var query *gorm.DB
	baseSQLPart := `
		SELECT s.status, COUNT(*) as count
		FROM (
			SELECT DISTINCT ON (order_id) id, order_id, status
			FROM shipments
			ORDER BY order_id, created_at DESC
		) s
		JOIN orders o ON o.id = s.order_id
		WHERE s.status IN (?, ?, ?)
		AND o.deleted_at IS NULL
	`
	orderBySQLPart := `
		GROUP BY s.status
		ORDER BY CASE
			WHEN s.status = 'pending' THEN 1
			WHEN s.status = 'in_transit' THEN 2
			WHEN s.status = 'delivered' THEN 3
			ELSE 4
		END
	`

	// Construir filtros de fecha
	dateFilterPart := ""
	var dateParams []interface{}
	if startDate != nil {
		dateFilterPart += " AND o.created_at >= ?"
		dateParams = append(dateParams, *startDate)
	}
	if endDate != nil {
		dateFilterPart += " AND o.created_at < ?"
		dateParams = append(dateParams, *endDate)
	}

	// Handle both business_id and integration_id filters
	if businessID != nil && *businessID > 0 && integrationID != nil && *integrationID > 0 {
		allParams := []interface{}{"pending", "in_transit", "delivered", *businessID, *integrationID}
		allParams = append(allParams, dateParams...)
		query = r.db.Conn(ctx).Raw(
			baseSQLPart+`AND o.business_id = ? AND o.integration_id = ?`+dateFilterPart+orderBySQLPart,
			allParams...,
		)
	} else if businessID != nil && *businessID > 0 {
		allParams := []interface{}{"pending", "in_transit", "delivered", *businessID}
		allParams = append(allParams, dateParams...)
		query = r.db.Conn(ctx).Raw(
			baseSQLPart+`AND o.business_id = ?`+dateFilterPart+orderBySQLPart,
			allParams...,
		)
	} else if integrationID != nil && *integrationID > 0 {
		allParams := []interface{}{"pending", "in_transit", "delivered", *integrationID}
		allParams = append(allParams, dateParams...)
		query = r.db.Conn(ctx).Raw(
			baseSQLPart+`AND o.integration_id = ?`+dateFilterPart+orderBySQLPart,
			allParams...,
		)
	} else {
		allParams := []interface{}{"pending", "in_transit", "delivered"}
		allParams = append(allParams, dateParams...)
		query = r.db.Conn(ctx).Raw(
			baseSQLPart+dateFilterPart+orderBySQLPart,
			allParams...,
		)
	}

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	// Mapear resultados
	statuses := make([]domain.ShipmentsByStatus, len(results))
	for i, result := range results {
		statuses[i] = domain.ShipmentsByStatus{
			Status: result.Status,
			Count:  result.Count,
		}
	}

	return statuses, nil
}

// GetShipmentsByCarrier obtiene envíos agrupados por transportista
// Solo cuenta el shipment más reciente por orden para evitar duplicados
func (r *Repository) GetShipmentsByCarrier(ctx context.Context, businessID *uint, integrationID *uint, startDate *time.Time, endDate *time.Time) ([]domain.ShipmentsByCarrier, error) {
	type Result struct {
		Carrier string `gorm:"column:carrier"`
		Count   int64  `gorm:"column:count"`
	}

	var results []Result

	baseSQLPart := `
		SELECT TRIM(LOWER(s.carrier)) as carrier, COUNT(*) as count
		FROM (
			SELECT DISTINCT ON (order_id) id, order_id, carrier
			FROM shipments
			ORDER BY order_id, created_at DESC
		) s
		JOIN orders o ON o.id = s.order_id
		WHERE s.carrier IS NOT NULL AND s.carrier != ''
		AND o.deleted_at IS NULL
	`
	orderBySQLPart := `
		GROUP BY TRIM(LOWER(s.carrier))
		ORDER BY count DESC
	`

	// Construir filtros de fecha
	dateFilterPart := ""
	var dateParams []interface{}
	if startDate != nil {
		dateFilterPart += " AND o.created_at >= ?"
		dateParams = append(dateParams, *startDate)
	}
	if endDate != nil {
		dateFilterPart += " AND o.created_at < ?"
		dateParams = append(dateParams, *endDate)
	}

	var query *gorm.DB
	if businessID != nil && *businessID > 0 && integrationID != nil && *integrationID > 0 {
		allParams := []interface{}{*businessID, *integrationID}
		allParams = append(allParams, dateParams...)
		query = r.db.Conn(ctx).Raw(
			baseSQLPart+`AND o.business_id = ? AND o.integration_id = ?`+dateFilterPart+orderBySQLPart,
			allParams...,
		)
	} else if businessID != nil && *businessID > 0 {
		allParams := []interface{}{*businessID}
		allParams = append(allParams, dateParams...)
		query = r.db.Conn(ctx).Raw(
			baseSQLPart+`AND o.business_id = ?`+dateFilterPart+orderBySQLPart,
			allParams...,
		)
	} else if integrationID != nil && *integrationID > 0 {
		allParams := []interface{}{*integrationID}
		allParams = append(allParams, dateParams...)
		query = r.db.Conn(ctx).Raw(
			baseSQLPart+`AND o.integration_id = ?`+dateFilterPart+orderBySQLPart,
			allParams...,
		)
	} else {
		query = r.db.Conn(ctx).Raw(baseSQLPart + dateFilterPart + orderBySQLPart, dateParams...)
	}

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	// Mapear resultados
	carriers := make([]domain.ShipmentsByCarrier, len(results))
	for i, result := range results {
		carriers[i] = domain.ShipmentsByCarrier{
			Carrier: result.Carrier,
			Count:   result.Count,
		}
	}

	return carriers, nil
}

// GetShipmentsByCarrierToday obtiene envíos agrupados por transportista del día actual
// Solo cuenta el shipment más reciente por orden para evitar duplicados
func (r *Repository) GetShipmentsByCarrierToday(ctx context.Context, businessID *uint, integrationID *uint, startDate *time.Time, endDate *time.Time) ([]domain.ShipmentsByCarrier, error) {
	type Result struct {
		Carrier string `gorm:"column:carrier"`
		Count   int64  `gorm:"column:count"`
	}

	var results []Result

	baseSQLPart := `
		SELECT TRIM(LOWER(s.carrier)) as carrier, COUNT(*) as count
		FROM (
			SELECT DISTINCT ON (order_id) id, order_id, carrier
			FROM shipments
			WHERE DATE(created_at AT TIME ZONE 'America/Bogota') = DATE(NOW() AT TIME ZONE 'America/Bogota')
			ORDER BY order_id, created_at DESC
		) s
		JOIN orders o ON o.id = s.order_id
		WHERE s.carrier IS NOT NULL AND s.carrier != ''
		AND o.deleted_at IS NULL
	`
	orderBySQLPart := `
		GROUP BY TRIM(LOWER(s.carrier))
		ORDER BY count DESC
	`

	// Construir filtros de fecha (nota: GetShipmentsByCarrierToday ya filtra por hoy, pero se respeta el filtro si se proporciona)
	dateFilterPart := ""
	var dateParams []interface{}
	if startDate != nil {
		dateFilterPart += " AND o.created_at >= ?"
		dateParams = append(dateParams, *startDate)
	}
	if endDate != nil {
		dateFilterPart += " AND o.created_at < ?"
		dateParams = append(dateParams, *endDate)
	}

	var query *gorm.DB
	if businessID != nil && *businessID > 0 && integrationID != nil && *integrationID > 0 {
		allParams := []interface{}{*businessID, *integrationID}
		allParams = append(allParams, dateParams...)
		query = r.db.Conn(ctx).Raw(
			baseSQLPart+`AND o.business_id = ? AND o.integration_id = ?`+dateFilterPart+orderBySQLPart,
			allParams...,
		)
	} else if businessID != nil && *businessID > 0 {
		allParams := []interface{}{*businessID}
		allParams = append(allParams, dateParams...)
		query = r.db.Conn(ctx).Raw(
			baseSQLPart+`AND o.business_id = ?`+dateFilterPart+orderBySQLPart,
			allParams...,
		)
	} else if integrationID != nil && *integrationID > 0 {
		allParams := []interface{}{*integrationID}
		allParams = append(allParams, dateParams...)
		query = r.db.Conn(ctx).Raw(
			baseSQLPart+`AND o.integration_id = ?`+dateFilterPart+orderBySQLPart,
			allParams...,
		)
	} else {
		query = r.db.Conn(ctx).Raw(baseSQLPart + dateFilterPart + orderBySQLPart, dateParams...)
	}

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	// Mapear resultados
	carriers := make([]domain.ShipmentsByCarrier, len(results))
	for i, result := range results {
		carriers[i] = domain.ShipmentsByCarrier{
			Carrier: result.Carrier,
			Count:   result.Count,
		}
	}

	return carriers, nil
}

// GetShipmentsByWarehouse obtiene envíos agrupados por almacén
// Solo cuenta el shipment más reciente por orden para evitar duplicados
func (r *Repository) GetShipmentsByWarehouse(ctx context.Context, businessID *uint, integrationID *uint, limit int, startDate *time.Time, endDate *time.Time) ([]domain.ShipmentsByWarehouse, error) {
	type Result struct {
		WarehouseName string `gorm:"column:warehouse_name"`
		WarehouseID   *uint  `gorm:"column:warehouse_id"`
		Count         int64  `gorm:"column:count"`
	}

	var results []Result

	baseSQLPart := `
		SELECT s.warehouse_name, s.warehouse_id, COUNT(*) as count
		FROM (
			SELECT DISTINCT ON (order_id) id, order_id, warehouse_name, warehouse_id
			FROM shipments
			ORDER BY order_id, created_at DESC
		) s
		JOIN orders o ON o.id = s.order_id
		WHERE s.warehouse_name != '' AND s.warehouse_id IS NOT NULL
		AND o.deleted_at IS NULL
	`
	orderBySQLPart := `
		GROUP BY s.warehouse_id, s.warehouse_name
		ORDER BY count DESC
		LIMIT ?
	`

	// Construir filtros de fecha
	dateFilterPart := ""
	var dateParams []interface{}
	if startDate != nil {
		dateFilterPart += " AND o.created_at >= ?"
		dateParams = append(dateParams, *startDate)
	}
	if endDate != nil {
		dateFilterPart += " AND o.created_at < ?"
		dateParams = append(dateParams, *endDate)
	}

	var query *gorm.DB
	if businessID != nil && *businessID > 0 && integrationID != nil && *integrationID > 0 {
		allParams := []interface{}{*businessID, *integrationID}
		allParams = append(allParams, dateParams...)
		allParams = append(allParams, limit)
		query = r.db.Conn(ctx).Raw(
			baseSQLPart+`AND o.business_id = ? AND o.integration_id = ?`+dateFilterPart+orderBySQLPart,
			allParams...,
		)
	} else if businessID != nil && *businessID > 0 {
		allParams := []interface{}{*businessID}
		allParams = append(allParams, dateParams...)
		allParams = append(allParams, limit)
		query = r.db.Conn(ctx).Raw(
			baseSQLPart+`AND o.business_id = ?`+dateFilterPart+orderBySQLPart,
			allParams...,
		)
	} else if integrationID != nil && *integrationID > 0 {
		allParams := []interface{}{*integrationID}
		allParams = append(allParams, dateParams...)
		allParams = append(allParams, limit)
		query = r.db.Conn(ctx).Raw(
			baseSQLPart+`AND o.integration_id = ?`+dateFilterPart+orderBySQLPart,
			allParams...,
		)
	} else {
		allParams := append(dateParams, limit)
		query = r.db.Conn(ctx).Raw(baseSQLPart + dateFilterPart + orderBySQLPart, allParams...)
	}

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	// Mapear resultados
	warehouses := make([]domain.ShipmentsByWarehouse, len(results))
	for i, result := range results {
		warehouses[i] = domain.ShipmentsByWarehouse{
			WarehouseName: result.WarehouseName,
			WarehouseID:   result.WarehouseID,
			Count:         result.Count,
		}
	}

	return warehouses, nil
}

// GetShipmentsByDayOfWeek obtiene órdenes agrupadas por día para una semana específica (lunes a domingo)
func (r *Repository) GetShipmentsByDayOfWeek(ctx context.Context, businessID *uint, integrationID *uint, startDate *time.Time) ([]domain.ShipmentsByDayOfWeek, error) {
	// Si no se proporciona startDate, usar el lunes de la semana actual (en zona horaria America/Bogota)
	if startDate == nil {
		loc, _ := time.LoadLocation("America/Bogota")
		now := time.Now().In(loc)
		// Calcular el lunes de la semana actual
		daysToMonday := int(now.Weekday()) - 1
		if daysToMonday < 0 {
			daysToMonday = 6 // Si es domingo, ir al lunes de la semana anterior
		}
		monday := now.AddDate(0, 0, -daysToMonday)
		// Normalizar a inicio del día
		monday = time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, monday.Location())
		startDate = &monday
	}

	// Calcular el fin de la semana (domingo)
	endDate := startDate.AddDate(0, 0, 7)

	// Obtener todas las órdenes de la semana
	var orders []models.Order
	query := r.db.Conn(ctx).
		Model(&models.Order{}).
		Where("orders.created_at >= ?", startDate).
		Where("orders.created_at < ?", endDate).
		Where("orders.deleted_at IS NULL")

	// Aplicar filtro por business_id si está especificado
	if businessID != nil && *businessID > 0 {
		query = query.Where("orders.business_id = ?", *businessID)
	}

	// Aplicar filtro por integration_id si está especificado
	if integrationID != nil && *integrationID > 0 {
		query = query.Where("orders.integration_id = ?", *integrationID)
	}

	if err := query.Find(&orders).Error; err != nil {
		r.logger.Error().Err(err).Msg("Error al obtener órdenes para agrupar por día")
		return nil, err
	}

	// Agrupar órdenes por fecha en el código Go (convertir a zona horaria America/Bogota)
	loc, _ := time.LoadLocation("America/Bogota")
	countMap := make(map[string]int64)
	for _, order := range orders {
		// Convertir created_at a fecha YYYY-MM-DD en zona horaria America/Bogota
		dateStr := order.CreatedAt.In(loc).Format("2006-01-02")
		countMap[dateStr]++
	}

	// Generar 7 días de la semana con sus conteos y porcentajes
	dayNames := []string{"Lunes", "Martes", "Miércoles", "Jueves", "Viernes", "Sábado", "Domingo"}
	ordersByDay := make([]domain.ShipmentsByDayOfWeek, 7)

	for i := 0; i < 7; i++ {
		currentDate := startDate.AddDate(0, 0, i)
		dateStr := currentDate.Format("2006-01-02")
		currentCount := countMap[dateStr]

		// Calcular porcentaje vs día anterior
		var percentageVsPrevious *float64
		if i > 0 {
			previousDateStr := startDate.AddDate(0, 0, i-1).Format("2006-01-02")
			previousCount := countMap[previousDateStr]

			if previousCount > 0 {
				// Calcular porcentaje: ((actual - anterior) / anterior) * 100
				percentage := float64(currentCount-previousCount) / float64(previousCount) * 100
				percentageVsPrevious = &percentage
			} else if currentCount > 0 {
				// Si anterior fue 0 pero actual > 0, es +100%
				percentage := 100.0
				percentageVsPrevious = &percentage
			}
		}

		ordersByDay[i] = domain.ShipmentsByDayOfWeek{
			Date:                 dateStr,
			DayName:              dayNames[i],
			Count:                currentCount,
			PercentageVsPrevious: percentageVsPrevious,
		}
	}

	return ordersByDay, nil
}

// GetOrdersByDepartment obtiene TODAS las órdenes agrupadas por departamento
func (r *Repository) GetOrdersByDepartment(ctx context.Context, businessID *uint, integrationID *uint, startDate *time.Time, endDate *time.Time) ([]domain.OrdersByDepartment, error) {
	type Result struct {
		Department string `gorm:"column:department"`
		Count      int64  `gorm:"column:count"`
	}

	var results []Result
	query := r.db.Conn(ctx).
		Model(&models.Order{}).
		Select("orders.shipping_state as department, COUNT(*) as count").
		Where("orders.shipping_state IS NOT NULL AND orders.shipping_state != ''").
		Group("orders.shipping_state").
		Order("count DESC")

	// Aplicar filtro por business_id si está especificado
	if businessID != nil && *businessID > 0 {
		query = query.Where("orders.business_id = ?", *businessID)
	}

	// Aplicar filtro por integration_id si está especificado
	if integrationID != nil && *integrationID > 0 {
		query = query.Where("orders.integration_id = ?", *integrationID)
	}

	// Aplicar filtro por rango de fechas
	if startDate != nil {
		query = query.Where("orders.created_at >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("orders.created_at < ?", *endDate)
	}

	if err := query.Scan(&results).Error; err != nil {
		r.logger.Error().Err(err).Msg("Error al obtener órdenes por departamento")
		return nil, err
	}

	// Mapear resultados
	departments := make([]domain.OrdersByDepartment, len(results))
	for i, result := range results {
		departments[i] = domain.OrdersByDepartment{
			Department: result.Department,
			Count:      result.Count,
		}
	}

	return departments, nil
}

// GetOrdersByBusiness obtiene órdenes agrupadas por business (solo para super admin)
func (r *Repository) GetOrdersByBusiness(ctx context.Context, limit int, startDate *time.Time, endDate *time.Time) ([]domain.OrdersByBusiness, error) {
	// Usar una subconsulta para obtener el nombre del business
	// Primero obtenemos los business_ids y sus conteos, luego obtenemos los nombres
	query := r.db.Conn(ctx).
		Model(&models.Order{}).
		Select("orders.business_id, COUNT(*) as order_count").
		Where("orders.business_id IS NOT NULL").
		Group("orders.business_id").
		Order("order_count DESC").
		Limit(limit)

	// Aplicar filtro por rango de fechas
	if startDate != nil {
		query = query.Where("orders.created_at >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("orders.created_at < ?", *endDate)
	}

	// Obtener los business_ids y conteos
	type BusinessCount struct {
		BusinessID uint  `gorm:"column:business_id"`
		OrderCount int64 `gorm:"column:order_count"`
	}
	var businessCounts []BusinessCount
	if err := query.Scan(&businessCounts).Error; err != nil {
		return nil, err
	}

	// Si no hay resultados, retornar vacío
	if len(businessCounts) == 0 {
		return []domain.OrdersByBusiness{}, nil
	}

	// Obtener los business_ids únicos
	businessIDs := make([]uint, len(businessCounts))
	businessCountMap := make(map[uint]int64)
	for i, bc := range businessCounts {
		businessIDs[i] = bc.BusinessID
		businessCountMap[bc.BusinessID] = bc.OrderCount
	}

	// Obtener los nombres de los businesses
	var businesses []models.Business
	if err := r.db.Conn(ctx).
		Model(&models.Business{}).
		Where("id IN ?", businessIDs).
		Find(&businesses).Error; err != nil {
		// Si no se pueden obtener los nombres, retornar solo los IDs
		r.logger.Warn().Err(err).Msg("No se pudieron obtener los nombres de los businesses, retornando solo IDs")
		resultsList := make([]domain.OrdersByBusiness, len(businessCounts))
		for i, bc := range businessCounts {
			resultsList[i] = domain.OrdersByBusiness{
				BusinessID:   bc.BusinessID,
				BusinessName: fmt.Sprintf("Business #%d", bc.BusinessID),
				OrderCount:   bc.OrderCount,
			}
		}
		return resultsList, nil
	}

	// Crear un mapa de ID a nombre
	businessNameMap := make(map[uint]string)
	for _, b := range businesses {
		businessNameMap[b.ID] = b.Name
	}

	// Mapear resultados con los nombres
	resultsList := make([]domain.OrdersByBusiness, len(businessCounts))
	for i, bc := range businessCounts {
		name := businessNameMap[bc.BusinessID]
		if name == "" {
			name = fmt.Sprintf("Business #%d", bc.BusinessID)
		}
		resultsList[i] = domain.OrdersByBusiness{
			BusinessID:   bc.BusinessID,
			BusinessName: name,
			OrderCount:   bc.OrderCount,
		}
	}

	return resultsList, nil
}

// GetOrdersByWeek obtiene órdenes agrupadas por semana (últimas 12 semanas)
// Semana comienza el lunes (ISO 8601)
func (r *Repository) GetOrdersByWeek(ctx context.Context, businessID *uint, integrationID *uint, startDate *time.Time, endDate *time.Time) ([]domain.OrdersByWeek, error) {
	type Result struct {
		StartDate string `gorm:"column:start_date"`
		EndDate   string `gorm:"column:end_date"`
		Count     int64  `gorm:"column:count"`
	}

	var results []Result

	baseSQLPart := `
		SELECT
			TO_CHAR(DATE_TRUNC('week', orders.created_at AT TIME ZONE 'America/Bogota'), 'YYYY-MM-DD') as start_date,
			TO_CHAR(DATE_TRUNC('week', orders.created_at AT TIME ZONE 'America/Bogota') + INTERVAL '6 days', 'YYYY-MM-DD') as end_date,
			COUNT(*) as count
		FROM orders
		WHERE orders.deleted_at IS NULL
	`
	groupByPart := `
		GROUP BY DATE_TRUNC('week', orders.created_at AT TIME ZONE 'America/Bogota')
		ORDER BY DATE_TRUNC('week', orders.created_at AT TIME ZONE 'America/Bogota') DESC
		LIMIT 12
	`

	// Construir filtros de fecha
	dateFilterPart := ""
	var dateParams []interface{}
	if startDate != nil {
		dateFilterPart += " AND orders.created_at >= ?"
		dateParams = append(dateParams, *startDate)
	}
	if endDate != nil {
		dateFilterPart += " AND orders.created_at < ?"
		dateParams = append(dateParams, *endDate)
	}

	var query *gorm.DB

	// Construcción del query con filtros opcionales
	if businessID != nil && *businessID > 0 && integrationID != nil && *integrationID > 0 {
		allParams := []interface{}{*businessID, *integrationID}
		allParams = append(allParams, dateParams...)
		query = r.db.Conn(ctx).Raw(
			baseSQLPart+`AND orders.business_id = ? AND orders.integration_id = ? `+dateFilterPart+groupByPart,
			allParams...,
		)
	} else if businessID != nil && *businessID > 0 {
		allParams := []interface{}{*businessID}
		allParams = append(allParams, dateParams...)
		query = r.db.Conn(ctx).Raw(
			baseSQLPart+`AND orders.business_id = ? `+dateFilterPart+groupByPart,
			allParams...,
		)
	} else if integrationID != nil && *integrationID > 0 {
		allParams := []interface{}{*integrationID}
		allParams = append(allParams, dateParams...)
		query = r.db.Conn(ctx).Raw(
			baseSQLPart+`AND orders.integration_id = ? `+dateFilterPart+groupByPart,
			allParams...,
		)
	} else {
		query = r.db.Conn(ctx).Raw(baseSQLPart + dateFilterPart + groupByPart, dateParams...)
	}

	if err := query.Scan(&results).Error; err != nil {
		r.logger.Error().Err(err).Msg("Error al obtener órdenes por semana")
		return nil, err
	}

	// Invertir orden (queremos de más antigua a más reciente)
	for i := len(results) / 2; i >= 0; i-- {
		opp := len(results) - 1 - i
		results[i], results[opp] = results[opp], results[i]
	}

	// Mapear resultados
	ordersByWeek := make([]domain.OrdersByWeek, len(results))
	for i, result := range results {
		weekNumber := i + 1
		weekLabel := fmt.Sprintf("Sem %d - %s a %s", weekNumber, result.StartDate, result.EndDate)

		ordersByWeek[i] = domain.OrdersByWeek{
			Week:       weekLabel,
			WeekNumber: weekNumber,
			StartDate:  result.StartDate,
			EndDate:    result.EndDate,
			Count:      result.Count,
		}
	}

	return ordersByWeek, nil
}

// GetOrdersByMonth obtiene órdenes agrupadas por mes del año actual con porcentaje mes-a-mes
func (r *Repository) GetOrdersByMonth(ctx context.Context, businessID *uint, integrationID *uint, startDate *time.Time, endDate *time.Time) ([]domain.OrdersByMonth, error) {
	type Result struct {
		Month int   `gorm:"column:month"`
		Year  int   `gorm:"column:year"`
		Count int64 `gorm:"column:count"`
	}

	var results []Result

	baseSQLPart := `
		SELECT
			EXTRACT(MONTH FROM orders.created_at)::int as month,
			EXTRACT(YEAR FROM orders.created_at)::int as year,
			COUNT(*) as count
		FROM orders
		WHERE EXTRACT(YEAR FROM orders.created_at) = EXTRACT(YEAR FROM NOW())
	`
	groupByPart := `
		GROUP BY EXTRACT(MONTH FROM orders.created_at), EXTRACT(YEAR FROM orders.created_at)
		ORDER BY year ASC, month ASC
	`

	// Construir filtros de fecha
	dateFilterPart := ""
	var dateParams []interface{}
	if startDate != nil {
		dateFilterPart += " AND orders.created_at >= ?"
		dateParams = append(dateParams, *startDate)
	}
	if endDate != nil {
		dateFilterPart += " AND orders.created_at < ?"
		dateParams = append(dateParams, *endDate)
	}

	var query *gorm.DB

	// Construcción del query con filtros opcionales
	if businessID != nil && *businessID > 0 && integrationID != nil && *integrationID > 0 {
		allParams := []interface{}{*businessID, *integrationID}
		allParams = append(allParams, dateParams...)
		query = r.db.Conn(ctx).Raw(
			baseSQLPart+`AND orders.business_id = ? AND orders.integration_id = ? `+dateFilterPart+groupByPart,
			allParams...,
		)
	} else if businessID != nil && *businessID > 0 {
		allParams := []interface{}{*businessID}
		allParams = append(allParams, dateParams...)
		query = r.db.Conn(ctx).Raw(
			baseSQLPart+`AND orders.business_id = ? `+dateFilterPart+groupByPart,
			allParams...,
		)
	} else if integrationID != nil && *integrationID > 0 {
		allParams := []interface{}{*integrationID}
		allParams = append(allParams, dateParams...)
		query = r.db.Conn(ctx).Raw(
			baseSQLPart+`AND orders.integration_id = ? `+dateFilterPart+groupByPart,
			allParams...,
		)
	} else {
		query = r.db.Conn(ctx).Raw(baseSQLPart + dateFilterPart + groupByPart, dateParams...)
	}

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	monthNames := []string{"Enero", "Febrero", "Marzo", "Abril", "Mayo", "Junio", "Julio", "Agosto", "Septiembre", "Octubre", "Noviembre", "Diciembre"}

	// Crear mapa de mes -> count para acceso rápido
	monthCounts := make(map[int]int64)
	for _, result := range results {
		monthCounts[result.Month] = result.Count
	}

	ordersByMonth := make([]domain.OrdersByMonth, len(results))
	now := time.Now()
	currentMonth := int(now.Month())
	currentDay := now.Day()

	for i, result := range results {
		monthName := monthNames[result.Month-1]
		percentage := float64(0)

		// Si es el primer mes (Enero) o no hay datos del mes anterior, mostrar 0%
		if result.Month == 1 {
			percentage = 0
		} else {
			prevMonth := result.Month - 1
			prevCount := monthCounts[prevMonth]

			if prevCount == 0 {
				percentage = 0
			} else {
				currentCount := result.Count

				// Si es el mes actual (incompleto), normalizar por días disponibles
				if result.Month == currentMonth {
					// Obtener último día del mes anterior
					firstOfCurrentMonth := time.Date(result.Year, time.Month(result.Month), 1, 0, 0, 0, 0, now.Location())
					lastOfPrevMonth := firstOfCurrentMonth.AddDate(0, 0, -1)
					daysInPrevMonth := lastOfPrevMonth.Day()

					// Proyectar órdenes de los primeros 'currentDay' del mes anterior
					projectedPrevMonthCount := (float64(prevCount) / float64(daysInPrevMonth)) * float64(currentDay)

					// Calcular porcentaje comparando con la proyección
					if projectedPrevMonthCount > 0 {
						percentage = ((float64(currentCount) - projectedPrevMonthCount) / projectedPrevMonthCount) * 100
					}
				} else {
					// Meses completos: comparación directa
					percentage = ((float64(currentCount) - float64(prevCount)) / float64(prevCount)) * 100
				}
			}
		}

		ordersByMonth[i] = domain.OrdersByMonth{
			Month:       monthName,
			MonthNumber: int(result.Month),
			Year:        result.Year,
			Count:       result.Count,
			Percentage:  percentage,
		}
	}

	return ordersByMonth, nil
}

// GetTopSellingDays obtiene los TOP N días de mayor demanda (fechas específicas)
func (r *Repository) GetTopSellingDays(ctx context.Context, businessID *uint, integrationID *uint, limit int) ([]domain.TopSellingDay, error) {
	type Result struct {
		Date  time.Time `gorm:"column:date"`
		Count int64     `gorm:"column:count"`
	}

	var results []Result

	query := r.db.Conn(ctx).
		Model(&models.Order{}).
		Select("DATE(orders.created_at AT TIME ZONE 'America/Bogota') as date, COUNT(*) as count").
		Where("orders.deleted_at IS NULL").
		Group("DATE(orders.created_at AT TIME ZONE 'America/Bogota')").
		Order("count DESC").
		Limit(limit)

	// Aplicar filtro por business_id si está especificado
	if businessID != nil && *businessID > 0 {
		query = query.Where("orders.business_id = ?", *businessID)
	}

	// Aplicar filtro por integration_id si está especificado
	if integrationID != nil && *integrationID > 0 {
		query = query.Where("orders.integration_id = ?", *integrationID)
	}

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	// Mapear resultados a TopSellingDay
	topDays := make([]domain.TopSellingDay, len(results))
	dayNames := []string{"Domingo", "Lunes", "Martes", "Miércoles", "Jueves", "Viernes", "Sábado"}
	monthNames := []string{"ene", "feb", "mar", "abr", "may", "jun", "jul", "ago", "sep", "oct", "nov", "dic"}

	for i, result := range results {
		dayName := dayNames[result.Date.Weekday()]
		monthNum := result.Date.Month() - 1
		monthShort := monthNames[monthNum]
		formatted := fmt.Sprintf("%s %d %s", dayName, result.Date.Day(), monthShort)

		topDays[i] = domain.TopSellingDay{
			Date:      result.Date.Format("2006-01-02"),
			DayName:   dayName,
			Formatted: formatted,
			Total:     result.Count,
		}
	}

	return topDays, nil
}
