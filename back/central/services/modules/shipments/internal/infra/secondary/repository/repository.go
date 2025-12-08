package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/infra/secondary/repository/mappers"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

// Repository implementa el repositorio de envíos
type Repository struct {
	db db.IDatabase
}

// New crea una nueva instancia del repositorio
func New(database db.IDatabase) domain.IRepository {
	return &Repository{
		db: database,
	}
}

// CreateShipment crea un nuevo envío en la base de datos
func (r *Repository) CreateShipment(ctx context.Context, shipment *domain.Shipment) error {
	dbShipment := mappers.ToDBShipment(shipment)
	if err := r.db.Conn(ctx).Create(dbShipment).Error; err != nil {
		return err
	}
	// Actualizar el ID del modelo de dominio con el ID generado
	shipment.ID = dbShipment.ID
	return nil
}

// GetShipmentByID obtiene un envío por su ID
func (r *Repository) GetShipmentByID(ctx context.Context, id uint) (*domain.Shipment, error) {
	var shipment models.Shipment
	err := r.db.Conn(ctx).
		Preload("Order").
		Preload("ShippingAddress").
		Where("id = ?", id).
		First(&shipment).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrShipmentNotFound
		}
		return nil, err
	}

	return mappers.ToDomainShipment(&shipment), nil
}

// GetShipmentByTrackingNumber obtiene un envío por su número de tracking
func (r *Repository) GetShipmentByTrackingNumber(ctx context.Context, trackingNumber string) (*domain.Shipment, error) {
	var shipment models.Shipment
	err := r.db.Conn(ctx).
		Preload("Order").
		Preload("ShippingAddress").
		Where("tracking_number = ?", trackingNumber).
		First(&shipment).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrShipmentNotFound
		}
		return nil, err
	}

	return mappers.ToDomainShipment(&shipment), nil
}

// GetShipmentsByOrderID obtiene todos los envíos de una orden
func (r *Repository) GetShipmentsByOrderID(ctx context.Context, orderID string) ([]domain.Shipment, error) {
	var shipments []models.Shipment
	err := r.db.Conn(ctx).
		Preload("Order").
		Preload("ShippingAddress").
		Where("order_id = ?", orderID).
		Order("created_at DESC").
		Find(&shipments).Error

	if err != nil {
		return nil, err
	}

	domainShipments := make([]domain.Shipment, len(shipments))
	for i, shipment := range shipments {
		domainShipments[i] = *mappers.ToDomainShipment(&shipment)
	}

	return domainShipments, nil
}

// ListShipments obtiene una lista paginada de envíos con filtros optimizados
func (r *Repository) ListShipments(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]domain.Shipment, int64, error) {
	var shipments []models.Shipment
	var total int64

	// Construir query base
	query := r.db.Conn(ctx).Model(&models.Shipment{})

	// Filtro por order_id
	if orderID, ok := filters["order_id"].(string); ok && orderID != "" {
		query = query.Where("shipments.order_id = ?", orderID)
	}

	// Filtro por múltiples order_ids
	if orderIDs, ok := filters["order_ids"].([]string); ok && len(orderIDs) > 0 {
		query = query.Where("shipments.order_id IN ?", orderIDs)
	}

	// Filtro por tracking_number (búsqueda parcial)
	if trackingNumber, ok := filters["tracking_number"].(string); ok && trackingNumber != "" {
		query = query.Where("shipments.tracking_number ILIKE ?", "%"+trackingNumber+"%")
	}

	// Filtro por múltiples tracking_numbers
	if trackingNumbers, ok := filters["tracking_numbers"].([]string); ok && len(trackingNumbers) > 0 {
		query = query.Where("shipments.tracking_number IN ?", trackingNumbers)
	}

	// Filtro por carrier
	if carrier, ok := filters["carrier"].(string); ok && carrier != "" {
		query = query.Where("shipments.carrier ILIKE ?", "%"+carrier+"%")
	}

	// Filtro por carrier_code
	if carrierCode, ok := filters["carrier_code"].(string); ok && carrierCode != "" {
		query = query.Where("shipments.carrier_code = ?", carrierCode)
	}

	// Filtro por status
	if status, ok := filters["status"].(string); ok && status != "" {
		query = query.Where("shipments.status = ?", status)
	}

	// Filtro por múltiples statuses
	if statuses, ok := filters["statuses"].([]string); ok && len(statuses) > 0 {
		query = query.Where("shipments.status IN ?", statuses)
	}

	// Filtro por guide_id
	if guideID, ok := filters["guide_id"].(string); ok && guideID != "" {
		query = query.Where("shipments.guide_id = ?", guideID)
	}

	// Filtro por warehouse_id
	if warehouseID, ok := filters["warehouse_id"].(uint); ok && warehouseID > 0 {
		query = query.Where("shipments.warehouse_id = ?", warehouseID)
	}

	// Filtro por driver_id
	if driverID, ok := filters["driver_id"].(uint); ok && driverID > 0 {
		query = query.Where("shipments.driver_id = ?", driverID)
	}

	// Filtro por is_last_mile
	if isLastMile, ok := filters["is_last_mile"].(bool); ok {
		query = query.Where("shipments.is_last_mile = ?", isLastMile)
	}

	// Filtros de fecha - shipped_at
	if shippedAfter, ok := filters["shipped_after"].(string); ok && shippedAfter != "" {
		query = query.Where("shipments.shipped_at >= ?", shippedAfter)
	}

	if shippedBefore, ok := filters["shipped_before"].(string); ok && shippedBefore != "" {
		query = query.Where("shipments.shipped_at <= ?", shippedBefore)
	}

	// Filtros de fecha - delivered_at
	if deliveredAfter, ok := filters["delivered_after"].(string); ok && deliveredAfter != "" {
		query = query.Where("shipments.delivered_at >= ?", deliveredAfter)
	}

	if deliveredBefore, ok := filters["delivered_before"].(string); ok && deliveredBefore != "" {
		query = query.Where("shipments.delivered_at <= ?", deliveredBefore)
	}

	// Filtros de fecha - created_at (compatibilidad)
	if startDate, ok := filters["start_date"].(string); ok && startDate != "" {
		query = query.Where("shipments.created_at >= ?", startDate)
	}

	if endDate, ok := filters["end_date"].(string); ok && endDate != "" {
		query = query.Where("shipments.created_at <= ?", endDate)
	}

	if createdAfter, ok := filters["created_after"].(string); ok && createdAfter != "" {
		query = query.Where("shipments.created_at >= ?", createdAfter)
	}

	if createdBefore, ok := filters["created_before"].(string); ok && createdBefore != "" {
		query = query.Where("shipments.created_at <= ?", createdBefore)
	}

	// Filtros de fecha - updated_at
	if updatedAfter, ok := filters["updated_after"].(string); ok && updatedAfter != "" {
		query = query.Where("shipments.updated_at >= ?", updatedAfter)
	}

	if updatedBefore, ok := filters["updated_before"].(string); ok && updatedBefore != "" {
		query = query.Where("shipments.updated_at <= ?", updatedBefore)
	}

	// Filtro por business_id (a través de JOIN con Order -> Business)
	if businessID, ok := filters["business_id"].(uint); ok && businessID > 0 {
		query = query.
			Joins("INNER JOIN orders ON shipments.order_id = orders.id").
			Where("orders.business_id = ?", businessID)
	}

	// Filtro por integration_id (a través de JOIN con Order -> Integration)
	if integrationID, ok := filters["integration_id"].(uint); ok && integrationID > 0 {
		query = query.
			Joins("INNER JOIN orders ON shipments.order_id = orders.id").
			Where("orders.integration_id = ?", integrationID)
	}

	// Filtro por integration_type (a través de JOIN con Order -> Integration -> IntegrationType)
	if integrationType, ok := filters["integration_type"].(string); ok && integrationType != "" {
		query = query.
			Joins("INNER JOIN orders ON shipments.order_id = orders.id").
			Joins("INNER JOIN integrations ON orders.integration_id = integrations.id").
			Joins("INNER JOIN integration_types ON integrations.integration_type_id = integration_types.id").
			Where("integration_types.code = ?", integrationType)
	}

	// Usar DISTINCT si hay JOINs para evitar duplicados
	hasJoins := filters["business_id"] != nil || filters["integration_id"] != nil || filters["integration_type"] != nil
	if hasJoins {
		query = query.Distinct("shipments.id")
	}

	// Contar total (antes de aplicar paginación y ordenamiento)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Aplicar ordenamiento
	sortBy := "shipments.created_at"
	if sort, ok := filters["sort_by"].(string); ok && sort != "" {
		// Mapear campos de ordenamiento
		sortFieldMap := map[string]string{
			"id":                "shipments.id",
			"order_id":          "shipments.order_id",
			"tracking_number":   "shipments.tracking_number",
			"status":            "shipments.status",
			"carrier":           "shipments.carrier",
			"shipped_at":        "shipments.shipped_at",
			"delivered_at":      "shipments.delivered_at",
			"created_at":        "shipments.created_at",
			"updated_at":        "shipments.updated_at",
			"warehouse_id":      "shipments.warehouse_id",
			"driver_id":         "shipments.driver_id",
		}
		if mappedField, exists := sortFieldMap[sort]; exists {
			sortBy = mappedField
		}
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
	query = query.Preload("Order").Preload("ShippingAddress")

	// Ejecutar query
	if err := query.Find(&shipments).Error; err != nil {
		return nil, 0, err
	}

	// Convertir a dominio
	domainShipments := make([]domain.Shipment, len(shipments))
	for i, shipment := range shipments {
		domainShipments[i] = *mappers.ToDomainShipment(&shipment)
	}

	return domainShipments, total, nil
}

// UpdateShipment actualiza un envío existente
func (r *Repository) UpdateShipment(ctx context.Context, shipment *domain.Shipment) error {
	dbShipment := mappers.ToDBShipment(shipment)
	return r.db.Conn(ctx).Save(dbShipment).Error
}

// DeleteShipment elimina (soft delete) un envío
func (r *Repository) DeleteShipment(ctx context.Context, id uint) error {
	return r.db.Conn(ctx).Where("id = ?", id).Delete(&models.Shipment{}).Error
}

// ShipmentExists verifica si existe un envío con el tracking number para una orden
func (r *Repository) ShipmentExists(ctx context.Context, orderID string, trackingNumber string) (bool, error) {
	var count int64
	err := r.db.Conn(ctx).
		Model(&models.Shipment{}).
		Where("order_id = ? AND tracking_number = ?", orderID, trackingNumber).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

