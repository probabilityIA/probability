package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/orders/domain"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

// Repository implementa el repositorio de órdenes
type Repository struct {
	db db.IDatabase
}

// New crea una nueva instancia del repositorio
func New(database db.IDatabase) domain.IRepository {
	return &Repository{
		db: database,
	}
}

// CreateOrder crea una nueva orden en la base de datos
func (r *Repository) CreateOrder(ctx context.Context, order *models.Order) error {
	return r.db.Conn(ctx).Create(order).Error
}

// GetOrderByID obtiene una orden por su ID
func (r *Repository) GetOrderByID(ctx context.Context, id string) (*models.Order, error) {
	var order models.Order
	err := r.db.Conn(ctx).
		Preload("Business").
		Preload("Integration").
		Preload("PaymentMethod").
		Where("id = ?", id).
		First(&order).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("order not found")
		}
		return nil, err
	}

	return &order, nil
}

// GetOrderByInternalNumber obtiene una orden por su número interno
func (r *Repository) GetOrderByInternalNumber(ctx context.Context, internalNumber string) (*models.Order, error) {
	var order models.Order
	err := r.db.Conn(ctx).
		Preload("Business").
		Preload("Integration").
		Preload("PaymentMethod").
		Where("internal_number = ?", internalNumber).
		First(&order).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("order not found")
		}
		return nil, err
	}

	return &order, nil
}

// ListOrders obtiene una lista paginada de órdenes con filtros
func (r *Repository) ListOrders(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	// Construir query base
	query := r.db.Conn(ctx).Model(&models.Order{})

	// Aplicar filtros
	if businessID, ok := filters["business_id"].(uint); ok && businessID > 0 {
		query = query.Where("business_id = ?", businessID)
	}

	if integrationID, ok := filters["integration_id"].(uint); ok && integrationID > 0 {
		query = query.Where("integration_id = ?", integrationID)
	}

	if integrationType, ok := filters["integration_type"].(string); ok && integrationType != "" {
		query = query.Where("integration_type = ?", integrationType)
	}

	if status, ok := filters["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}

	if customerEmail, ok := filters["customer_email"].(string); ok && customerEmail != "" {
		query = query.Where("customer_email ILIKE ?", "%"+customerEmail+"%")
	}

	if customerPhone, ok := filters["customer_phone"].(string); ok && customerPhone != "" {
		query = query.Where("customer_phone ILIKE ?", "%"+customerPhone+"%")
	}

	if platform, ok := filters["platform"].(string); ok && platform != "" {
		query = query.Where("platform = ?", platform)
	}

	if isPaid, ok := filters["is_paid"].(bool); ok {
		query = query.Where("is_paid = ?", isPaid)
	}

	if warehouseID, ok := filters["warehouse_id"].(uint); ok && warehouseID > 0 {
		query = query.Where("warehouse_id = ?", warehouseID)
	}

	if driverID, ok := filters["driver_id"].(uint); ok && driverID > 0 {
		query = query.Where("driver_id = ?", driverID)
	}

	// Filtros de fecha
	if startDate, ok := filters["start_date"].(string); ok && startDate != "" {
		query = query.Where("created_at >= ?", startDate)
	}

	if endDate, ok := filters["end_date"].(string); ok && endDate != "" {
		query = query.Where("created_at <= ?", endDate)
	}

	// Contar total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Aplicar ordenamiento
	sortBy := "created_at"
	if sort, ok := filters["sort_by"].(string); ok && sort != "" {
		sortBy = sort
	}

	sortOrder := "desc"
	if order, ok := filters["sort_order"].(string); ok && order != "" {
		sortOrder = order
	}

	query = query.Order(fmt.Sprintf("%s %s", sortBy, sortOrder))

	// Aplicar paginación
	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)

	// Precargar relaciones
	query = query.Preload("Business").
		Preload("Integration").
		Preload("PaymentMethod")

	// Ejecutar query
	if err := query.Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// UpdateOrder actualiza una orden existente
func (r *Repository) UpdateOrder(ctx context.Context, order *models.Order) error {
	return r.db.Conn(ctx).Save(order).Error
}

// DeleteOrder elimina (soft delete) una orden
func (r *Repository) DeleteOrder(ctx context.Context, id string) error {
	return r.db.Conn(ctx).Where("id = ?", id).Delete(&models.Order{}).Error
}

// OrderExists verifica si existe una orden con el external_id para una integración
func (r *Repository) OrderExists(ctx context.Context, externalID string, integrationID uint) (bool, error) {
	var count int64
	err := r.db.Conn(ctx).
		Model(&models.Order{}).
		Where("external_id = ? AND integration_id = ?", externalID, integrationID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
