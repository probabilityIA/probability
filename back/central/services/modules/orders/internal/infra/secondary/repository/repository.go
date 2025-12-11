package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/infra/secondary/repository/mappers"
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
func (r *Repository) CreateOrder(ctx context.Context, order *domain.Order) error {
	dbOrder := mappers.ToDBOrder(order)
	if err := r.db.Conn(ctx).Create(dbOrder).Error; err != nil {
		return err
	}
	// Actualizar el ID del modelo de dominio con el ID generado
	order.ID = dbOrder.ID
	return nil
}

// GetOrderByID obtiene una orden por su ID
func (r *Repository) GetOrderByID(ctx context.Context, id string) (*domain.Order, error) {
	var order models.Order
	err := r.db.Conn(ctx).
		Preload("Business").
		Preload("Integration").
		Preload("PaymentMethod").
		Preload("OrderItems.Product"). // Precargar OrderItems con Product para obtener información del catálogo
		Where("id = ?", id).
		First(&order).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("order not found")
		}
		return nil, err
	}

	return mappers.ToDomainOrder(&order), nil
}

// GetOrderByInternalNumber obtiene una orden por su número interno
func (r *Repository) GetOrderByInternalNumber(ctx context.Context, internalNumber string) (*domain.Order, error) {
	var order models.Order
	err := r.db.Conn(ctx).
		Preload("Business").
		Preload("Integration").
		Preload("PaymentMethod").
		Preload("OrderItems.Product"). // Precargar OrderItems con Product para obtener información del catálogo
		Where("internal_number = ?", internalNumber).
		First(&order).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("order not found")
		}
		return nil, err
	}

	return mappers.ToDomainOrder(&order), nil
}

// ListOrders obtiene una lista paginada de órdenes con filtros
func (r *Repository) ListOrders(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]domain.Order, int64, error) {
	var dbOrders []models.Order
	var total int64

	query := r.db.Conn(ctx).Model(&models.Order{})

	// Aplicar filtros
	if customerEmail, ok := filters["customer_email"].(string); ok && customerEmail != "" {
		query = query.Where("customer_email ILIKE ?", "%"+customerEmail+"%")
	}

	if customerPhone, ok := filters["customer_phone"].(string); ok && customerPhone != "" {
		query = query.Where("customer_phone ILIKE ?", "%"+customerPhone+"%")
	}

	if orderNumber, ok := filters["order_number"].(string); ok && orderNumber != "" {
		query = query.Where("order_number ILIKE ?", "%"+orderNumber+"%")
	}

	if internalNumber, ok := filters["internal_number"].(string); ok && internalNumber != "" {
		query = query.Where("internal_number ILIKE ?", "%"+internalNumber+"%")
	}

	if status, ok := filters["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
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
		Preload("PaymentMethod").
		Preload("OrderItems.Product") // Precargar OrderItems con Product para obtener información del catálogo

	// Paginación
	offset = (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&dbOrders).Error; err != nil {
		return nil, 0, err
	}

	// Mapear a dominio
	orders := make([]domain.Order, len(dbOrders))
	for i, dbOrder := range dbOrders {
		orders[i] = *mappers.ToDomainOrder(&dbOrder)
	}

	return orders, total, nil
}

// GetOrderRaw obtiene los metadatos crudos de una orden
func (r *Repository) GetOrderRaw(ctx context.Context, id string) (*domain.OrderChannelMetadata, error) {
	var dbMetadata models.OrderChannelMetadata
	if err := r.db.Conn(ctx).Where("order_id = ?", id).First(&dbMetadata).Error; err != nil {
		return nil, err
	}
	return mappers.ToDomainChannelMetadata(&dbMetadata), nil
}

// UpdateOrder actualiza una orden existente
func (r *Repository) UpdateOrder(ctx context.Context, order *domain.Order) error {
	dbOrder := mappers.ToDBOrder(order)
	return r.db.Conn(ctx).Save(dbOrder).Error
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

// ───────────────────────────────────────────
//
//	MÉTODOS PARA TABLAS RELACIONADAS
//
// ───────────────────────────────────────────

// CreateOrderItems crea múltiples items de orden
func (r *Repository) CreateOrderItems(ctx context.Context, items []*domain.OrderItem) error {
	if len(items) == 0 {
		return nil
	}

	// Convertir []*domain.OrderItem a []domain.OrderItem para usar el mapper
	domainItems := make([]domain.OrderItem, len(items))
	for i, item := range items {
		domainItems[i] = *item
	}

	// Usar el mapper para convertir a modelos de BD
	dbItems := mappers.ToDBOrderItems(domainItems)

	// Convertir a slice de punteros para CreateInBatches
	dbItemsPtrs := make([]*models.OrderItem, len(dbItems))
	for i := range dbItems {
		dbItemsPtrs[i] = &dbItems[i]
	}

	return r.db.Conn(ctx).CreateInBatches(dbItemsPtrs, 100).Error
}

// CreateAddresses crea múltiples direcciones
func (r *Repository) CreateAddresses(ctx context.Context, addresses []*domain.Address) error {
	if len(addresses) == 0 {
		return nil
	}

	dbAddresses := make([]*models.Address, len(addresses))
	for i, addr := range addresses {
		dbAddr := &models.Address{
			Model: gorm.Model{
				ID:        addr.ID,
				CreatedAt: addr.CreatedAt,
				UpdatedAt: addr.UpdatedAt,
				DeletedAt: gorm.DeletedAt{},
			},
			Type:         addr.Type,
			OrderID:      addr.OrderID,
			FirstName:    addr.FirstName,
			LastName:     addr.LastName,
			Company:      addr.Company,
			Phone:        addr.Phone,
			Street:       addr.Street,
			Street2:      addr.Street2,
			City:         addr.City,
			State:        addr.State,
			Country:      addr.Country,
			PostalCode:   addr.PostalCode,
			Latitude:     addr.Latitude,
			Longitude:    addr.Longitude,
			Instructions: addr.Instructions,
			IsDefault:    addr.IsDefault,
			Metadata:     addr.Metadata,
		}
		if addr.DeletedAt != nil {
			dbAddr.DeletedAt = gorm.DeletedAt{Time: *addr.DeletedAt, Valid: true}
		}
		dbAddresses[i] = dbAddr
	}

	return r.db.Conn(ctx).CreateInBatches(dbAddresses, 100).Error
}

// CreatePayments crea múltiples pagos
func (r *Repository) CreatePayments(ctx context.Context, payments []*domain.Payment) error {
	if len(payments) == 0 {
		return nil
	}

	dbPayments := make([]*models.Payment, len(payments))
	for i, p := range payments {
		dbPay := &models.Payment{
			Model: gorm.Model{
				ID:        p.ID,
				CreatedAt: p.CreatedAt,
				UpdatedAt: p.UpdatedAt,
				DeletedAt: gorm.DeletedAt{},
			},
			OrderID:          p.OrderID,
			PaymentMethodID:  p.PaymentMethodID,
			Amount:           p.Amount,
			Currency:         p.Currency,
			ExchangeRate:     p.ExchangeRate,
			Status:           p.Status,
			PaidAt:           p.PaidAt,
			ProcessedAt:      p.ProcessedAt,
			TransactionID:    p.TransactionID,
			PaymentReference: p.PaymentReference,
			Gateway:          p.Gateway,
			RefundAmount:     p.RefundAmount,
			RefundedAt:       p.RefundedAt,
			FailureReason:    p.FailureReason,
			Metadata:         p.Metadata,
		}
		if p.DeletedAt != nil {
			dbPay.DeletedAt = gorm.DeletedAt{Time: *p.DeletedAt, Valid: true}
		}
		dbPayments[i] = dbPay
	}

	return r.db.Conn(ctx).CreateInBatches(dbPayments, 100).Error
}

// CreateShipments crea múltiples envíos
func (r *Repository) CreateShipments(ctx context.Context, shipments []*domain.Shipment) error {
	if len(shipments) == 0 {
		return nil
	}

	dbShipments := make([]*models.Shipment, len(shipments))
	for i, s := range shipments {
		dbShip := &models.Shipment{
			Model: gorm.Model{
				ID:        s.ID,
				CreatedAt: s.CreatedAt,
				UpdatedAt: s.UpdatedAt,
				DeletedAt: gorm.DeletedAt{},
			},
			OrderID:           s.OrderID,
			TrackingNumber:    s.TrackingNumber,
			TrackingURL:       s.TrackingURL,
			Carrier:           s.Carrier,
			CarrierCode:       s.CarrierCode,
			GuideID:           s.GuideID,
			GuideURL:          s.GuideURL,
			Status:            s.Status,
			ShippedAt:         s.ShippedAt,
			DeliveredAt:       s.DeliveredAt,
			ShippingAddressID: s.ShippingAddressID,
			ShippingCost:      s.ShippingCost,
			InsuranceCost:     s.InsuranceCost,
			TotalCost:         s.TotalCost,
			Weight:            s.Weight,
			Height:            s.Height,
			Width:             s.Width,
			Length:            s.Length,
			WarehouseID:       s.WarehouseID,
			WarehouseName:     s.WarehouseName,
			DriverID:          s.DriverID,
			DriverName:        s.DriverName,
			IsLastMile:        s.IsLastMile,
			EstimatedDelivery: s.EstimatedDelivery,
			DeliveryNotes:     s.DeliveryNotes,
			Metadata:          s.Metadata,
		}
		if s.DeletedAt != nil {
			dbShip.DeletedAt = gorm.DeletedAt{Time: *s.DeletedAt, Valid: true}
		}
		dbShipments[i] = dbShip
	}

	return r.db.Conn(ctx).CreateInBatches(dbShipments, 100).Error
}

// CreateChannelMetadata crea metadata del canal
func (r *Repository) CreateChannelMetadata(ctx context.Context, metadata *domain.OrderChannelMetadata) error {
	if metadata == nil {
		return nil
	}
	dbMetadata := mappers.ToDBChannelMetadata(metadata)
	return r.db.Conn(ctx).Create(dbMetadata).Error
}

// ───────────────────────────────────────────
//
//	MÉTODOS DE CATÁLOGO (VALIDACIÓN)
//
// ───────────────────────────────────────────

// GetProductBySKU busca un producto por SKU y BusinessID
func (r *Repository) GetProductBySKU(ctx context.Context, businessID uint, sku string) (*domain.Product, error) {
	var product models.Product
	err := r.db.Conn(ctx).
		Where("business_id = ? AND sku = ?", businessID, sku).
		First(&product).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Retornar nil si no existe, no error
		}
		return nil, err
	}
	return mappers.ToDomainProduct(&product), nil
}

// CreateProduct crea un nuevo producto
func (r *Repository) CreateProduct(ctx context.Context, product *domain.Product) error {
	dbProduct := mappers.ToDBProduct(product)
	if err := r.db.Conn(ctx).Create(dbProduct).Error; err != nil {
		return err
	}
	product.ID = dbProduct.ID
	return nil
}

// GetClientByEmail busca un cliente por Email y BusinessID
func (r *Repository) GetClientByEmail(ctx context.Context, businessID uint, email string) (*domain.Client, error) {
	var client models.Client
	err := r.db.Conn(ctx).
		Where("business_id = ? AND email = ?", businessID, email).
		First(&client).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Retornar nil si no existe
		}
		return nil, err
	}
	return mappers.ToDomainClient(&client), nil
}

// GetClientByDNI busca un cliente por DNI y BusinessID
func (r *Repository) GetClientByDNI(ctx context.Context, businessID uint, dni string) (*domain.Client, error) {
	if dni == "" {
		return nil, nil // No buscar si el DNI está vacío
	}

	var client models.Client
	err := r.db.Conn(ctx).
		Where("business_id = ? AND dni = ?", businessID, dni).
		First(&client).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Retornar nil si no existe
		}
		return nil, err
	}
	return mappers.ToDomainClient(&client), nil
}

// CreateClient crea un nuevo cliente
func (r *Repository) CreateClient(ctx context.Context, client *domain.Client) error {
	dbClient := mappers.ToDBClient(client)
	if err := r.db.Conn(ctx).Create(dbClient).Error; err != nil {
		return err
	}
	client.ID = dbClient.ID
	return nil
}
