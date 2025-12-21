package repository

import (
	"context"
	"fmt"

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

// GetTopDrivers obtiene los top N transportadores por número de órdenes
func (r *Repository) GetTopDrivers(ctx context.Context, businessID *uint, limit int) ([]domain.TopDriver, error) {
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
func (r *Repository) GetDriversByLocation(ctx context.Context, businessID *uint, limit int) ([]domain.DriverByLocation, error) {
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
func (r *Repository) GetTopProducts(ctx context.Context, businessID *uint, limit int) ([]domain.TopProduct, error) {
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
func (r *Repository) GetProductsByCategory(ctx context.Context, businessID *uint) ([]domain.ProductByCategory, error) {
	type Result struct {
		Category string `gorm:"column:category"`
		Count    int64  `gorm:"column:count"`
	}

	var results []Result
	query := r.db.Conn(ctx).
		Model(&models.Product{}).
		Select("products.category, COUNT(DISTINCT products.id) as count").
		Where("products.category != ''").
		Group("products.category").
		Order("count DESC")

	// Aplicar filtro por business_id si está especificado y no es super user
	if businessID != nil && *businessID > 0 {
		query = query.Where("products.business_id = ?", *businessID)
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
func (r *Repository) GetProductsByBrand(ctx context.Context, businessID *uint) ([]domain.ProductByBrand, error) {
	type Result struct {
		Brand string `gorm:"column:brand"`
		Count int64  `gorm:"column:count"`
	}

	var results []Result
	query := r.db.Conn(ctx).
		Model(&models.Product{}).
		Select("products.brand, COUNT(DISTINCT products.id) as count").
		Where("products.brand != ''").
		Group("products.brand").
		Order("count DESC")

	// Aplicar filtro por business_id si está especificado y no es super user
	if businessID != nil && *businessID > 0 {
		query = query.Where("products.business_id = ?", *businessID)
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
func (r *Repository) GetShipmentsByStatus(ctx context.Context, businessID *uint) ([]domain.ShipmentsByStatus, error) {
	type Result struct {
		Status string `gorm:"column:status"`
		Count  int64  `gorm:"column:count"`
	}

	var results []Result
	query := r.db.Conn(ctx).
		Model(&models.Shipment{}).
		Select("shipments.status, COUNT(*) as count").
		Joins("JOIN orders ON orders.id = shipments.order_id").
		Group("shipments.status").
		Order("count DESC")

	// Aplicar filtro por business_id si está especificado y no es super user
	if businessID != nil && *businessID > 0 {
		query = query.Where("orders.business_id = ?", *businessID)
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
func (r *Repository) GetShipmentsByCarrier(ctx context.Context, businessID *uint) ([]domain.ShipmentsByCarrier, error) {
	type Result struct {
		Carrier string `gorm:"column:carrier"`
		Count   int64  `gorm:"column:count"`
	}

	var results []Result
	query := r.db.Conn(ctx).
		Model(&models.Shipment{}).
		Select("COALESCE(shipments.carrier, 'Sin transportista') as carrier, COUNT(*) as count").
		Joins("JOIN orders ON orders.id = shipments.order_id").
		Group("shipments.carrier").
		Order("count DESC")

	// Aplicar filtro por business_id si está especificado y no es super user
	if businessID != nil && *businessID > 0 {
		query = query.Where("orders.business_id = ?", *businessID)
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
func (r *Repository) GetShipmentsByWarehouse(ctx context.Context, businessID *uint, limit int) ([]domain.ShipmentsByWarehouse, error) {
	type Result struct {
		WarehouseName string `gorm:"column:warehouse_name"`
		WarehouseID   *uint  `gorm:"column:warehouse_id"`
		Count         int64  `gorm:"column:count"`
	}

	var results []Result
	query := r.db.Conn(ctx).
		Model(&models.Shipment{}).
		Select("shipments.warehouse_name, shipments.warehouse_id, COUNT(*) as count").
		Joins("JOIN orders ON orders.id = shipments.order_id").
		Where("shipments.warehouse_name != '' AND shipments.warehouse_id IS NOT NULL").
		Group("shipments.warehouse_id, shipments.warehouse_name").
		Order("count DESC").
		Limit(limit)

	// Aplicar filtro por business_id si está especificado y no es super user
	if businessID != nil && *businessID > 0 {
		query = query.Where("orders.business_id = ?", *businessID)
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

// GetOrdersByBusiness obtiene órdenes agrupadas por business (solo para super admin)
func (r *Repository) GetOrdersByBusiness(ctx context.Context, limit int) ([]domain.OrdersByBusiness, error) {
	// Usar una subconsulta para obtener el nombre del business
	// Primero obtenemos los business_ids y sus conteos, luego obtenemos los nombres
	query := r.db.Conn(ctx).
		Model(&models.Order{}).
		Select("orders.business_id, COUNT(*) as order_count").
		Where("orders.business_id IS NOT NULL").
		Group("orders.business_id").
		Order("order_count DESC").
		Limit(limit)

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
