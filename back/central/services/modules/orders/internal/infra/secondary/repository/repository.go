package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/infra/secondary/repository/mappers"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

// Repository implementa el repositorio de órdenes
type Repository struct {
	db           db.IDatabase
	imageURLBase string
}

// New crea una nueva instancia del repositorio
func New(database db.IDatabase, config env.IConfig) ports.IRepository {
	imageURLBase := config.Get("URL_BASE_DOMAIN_S3")
	return &Repository{
		db:           database,
		imageURLBase: imageURLBase,
	}
}

// CreateOrder crea una nueva orden en la base de datos
func (r *Repository) CreateOrder(ctx context.Context, order *entities.ProbabilityOrder) error {
	// Validaciones críticas antes de insertar
	if order.ExternalID == "" {
		return fmt.Errorf("error: intentando insertar orden sin external_id - OrderNumber: %s", order.OrderNumber)
	}
	if order.IntegrationID == 0 {
		return fmt.Errorf("error: intentando insertar orden sin integration_id - ExternalID: %s", order.ExternalID)
	}
	if order.BusinessID == nil || *order.BusinessID == 0 {
		return fmt.Errorf("error: intentando insertar orden sin business_id - ExternalID: %s", order.ExternalID)
	}

	dbOrder := mappers.ToDBOrder(order)
	if err := r.db.Conn(ctx).Create(dbOrder).Error; err != nil {
		// Detectar error de clave duplicada para external_id + integration_id
		errMsg := err.Error()
		if strings.Contains(errMsg, "duplicate key value violates unique constraint") &&
			(strings.Contains(errMsg, "idx_integration_external_id") || strings.Contains(errMsg, "SQLSTATE 23505")) {
			return domainerrors.ErrOrderAlreadyExists
		}
		return err
	}
	// Actualizar el ID del modelo de dominio con el ID generado
	order.ID = dbOrder.ID
	return nil
}

// GetFirstIntegrationIDByBusinessID obtiene la primera integración disponible para un negocio
func (r *Repository) GetFirstIntegrationIDByBusinessID(ctx context.Context, businessID uint) (uint, error) {
	var integration models.Integration
	err := r.db.Conn(ctx).
		Where("business_id = ?", businessID).
		Or("business_id IS NULL"). // Algunas integraciones pueden ser globales
		Order("business_id DESC, is_default DESC, id ASC").
		First(&integration).Error

	if err != nil {
		return 0, fmt.Errorf("error finding integration for business %d: %w", businessID, err)
	}

	return integration.ID, nil
}

// GetOrderByID obtiene una orden por su ID
func (r *Repository) GetOrderByID(ctx context.Context, id string) (*entities.ProbabilityOrder, error) {
	var order models.Order
	err := r.db.Conn(ctx).
		Preload("Business").
		Preload("Integration.IntegrationType"). // Precargar Integration con IntegrationType para obtener el logo
		Preload("PaymentMethod").
		Preload("OrderStatus").        // Precargar OrderStatus para obtener información del estado de Probability
		Preload("PaymentStatus").      // Precargar PaymentStatus
		Preload("FulfillmentStatus").  // Precargar FulfillmentStatus
		Preload("OrderItems.Product"). // Precargar OrderItems con Product para obtener información del catálogo
		Preload("ChannelMetadata").    // Precargar ChannelMetadata para acceso a RawData en scoring
		Where("id = ?", id).
		First(&order).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("order not found")
		}
		return nil, err
	}

	return mappers.ToDomainOrder(&order, r.imageURLBase), nil
}

// GetOrderByInternalNumber obtiene una orden por su número interno
func (r *Repository) GetOrderByInternalNumber(ctx context.Context, internalNumber string) (*entities.ProbabilityOrder, error) {
	var order models.Order
	err := r.db.Conn(ctx).
		Preload("Business").
		Preload("Integration.IntegrationType"). // Precargar Integration con IntegrationType para obtener el logo
		Preload("PaymentMethod").
		Preload("OrderStatus").        // Precargar OrderStatus para obtener información del estado de Probability
		Preload("PaymentStatus").      // Precargar PaymentStatus
		Preload("FulfillmentStatus").  // Precargar FulfillmentStatus
		Preload("OrderItems.Product"). // Precargar OrderItems con Product para obtener información del catálogo
		Preload("ChannelMetadata").    // Precargar ChannelMetadata para acceso a RawData en scoring
		Where("internal_number = ?", internalNumber).
		First(&order).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("order not found")
		}
		return nil, err
	}

	return mappers.ToDomainOrder(&order, r.imageURLBase), nil
}

// GetOrderByOrderNumber obtiene una orden por su order_number
func (r *Repository) GetOrderByOrderNumber(ctx context.Context, orderNumber string) (*entities.ProbabilityOrder, error) {
	var order models.Order
	err := r.db.Conn(ctx).
		Preload("Business").
		Preload("Integration.IntegrationType"). // Precargar Integration con IntegrationType para obtener el logo
		Preload("PaymentMethod").
		Preload("OrderStatus").        // Precargar OrderStatus para obtener información del estado de Probability
		Preload("PaymentStatus").      // Precargar PaymentStatus
		Preload("FulfillmentStatus").  // Precargar FulfillmentStatus
		Preload("OrderItems.Product"). // Precargar OrderItems con Product para obtener información del catálogo
		Preload("ChannelMetadata").    // Precargar ChannelMetadata para acceso a RawData en scoring
		Where("order_number = ?", orderNumber).
		First(&order).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("order not found")
		}
		return nil, err
	}

	return mappers.ToDomainOrder(&order, r.imageURLBase), nil
}

// ListOrders obtiene una lista paginada de órdenes con filtros
func (r *Repository) ListOrders(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]entities.ProbabilityOrder, int64, error) {
	var dbOrders []models.Order
	var total int64

	query := r.db.Conn(ctx).Model(&models.Order{})

	// Aplicar filtros
	if businessID, ok := filters["business_id"].(uint); ok && businessID > 0 {
		query = query.Where("business_id = ?", businessID)
	}

	if integrationID, ok := filters["integration_id"].(uint); ok && integrationID > 0 {
		query = query.Where("integration_id = ?", integrationID)
	}

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

	// Filtro por status_id (estado de Probability)
	if statusID, ok := filters["status_id"].(uint); ok && statusID > 0 {
		query = query.Where("status_id = ?", statusID)
	}
	// Mantener compatibilidad con filtro antiguo por status (string) si se necesita
	if status, ok := filters["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}

	if platform, ok := filters["platform"].(string); ok && platform != "" {
		query = query.Where("platform = ?", platform)
	}

	if isPaid, ok := filters["is_paid"].(bool); ok {
		query = query.Where("is_paid = ?", isPaid)
	}

	if isCOD, ok := filters["is_cod"].(bool); ok {
		if isCOD {
			query = query.Where(`
				cod_total > 0 OR 
				payment_details->>'gateway' ILIKE '%cod%' OR 
				payment_details->>'gateway' ILIKE '%cash%' OR 
				payment_details->>'gateway' ILIKE '%contra%' OR
				EXISTS (
					SELECT 1 FROM json_array_elements_text(CAST(payment_details->'payment_gateway_names' AS json)) as gateway
					WHERE gateway ILIKE '%cod%' OR gateway ILIKE '%cash%' OR gateway ILIKE '%contra%'
				)
			`)
		} else {
			query = query.Where(`
				(cod_total IS NULL OR cod_total = 0) AND 
				(payment_details->>'gateway' IS NULL OR (
					payment_details->>'gateway' NOT ILIKE '%cod%' AND 
					payment_details->>'gateway' NOT ILIKE '%cash%' AND 
					payment_details->>'gateway' NOT ILIKE '%contra%'
				)) AND
				NOT EXISTS (
					SELECT 1 FROM json_array_elements_text(CAST(payment_details->'payment_gateway_names' AS json)) as gateway
					WHERE gateway ILIKE '%cod%' OR gateway ILIKE '%cash%' OR gateway ILIKE '%contra%'
				)
			`)
		}
	}

	if paymentStatusID, ok := filters["payment_status_id"].(uint); ok && paymentStatusID > 0 {
		query = query.Where("payment_status_id = ?", paymentStatusID)
	}

	if fulfillmentStatusID, ok := filters["fulfillment_status_id"].(uint); ok && fulfillmentStatusID > 0 {
		query = query.Where("fulfillment_status_id = ?", fulfillmentStatusID)
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
		Preload("Integration.IntegrationType"). // Precargar Integration con IntegrationType para obtener el logo
		Preload("PaymentMethod").
		Preload("OrderStatus").       // Precargar OrderStatus para obtener información del estado de Probability
		Preload("PaymentStatus").     // Precargar PaymentStatus
		Preload("FulfillmentStatus"). // Precargar FulfillmentStatus
		Preload("OrderItems.Product") // Precargar OrderItems con Product para obtener información del catálogo

	// Paginación
	offset = (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&dbOrders).Error; err != nil {
		return nil, 0, err
	}

	// Mapear a dominio
	orders := make([]entities.ProbabilityOrder, len(dbOrders))
	for i, dbOrder := range dbOrders {
		orders[i] = *mappers.ToDomainOrder(&dbOrder, r.imageURLBase)
	}

	return orders, total, nil
}

// GetOrderRaw obtiene los metadatos crudos de una orden
func (r *Repository) GetOrderRaw(ctx context.Context, id string) (*entities.ProbabilityOrderChannelMetadata, error) {
	var dbMetadata models.OrderChannelMetadata
	if err := r.db.Conn(ctx).Where("order_id = ?", id).First(&dbMetadata).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("raw data not found for this order")
		}
		return nil, err
	}
	return mappers.ToDomainChannelMetadata(&dbMetadata), nil
}

// UpdateOrder actualiza una orden existente
func (r *Repository) UpdateOrder(ctx context.Context, order *entities.ProbabilityOrder) error {
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

// GetOrderByExternalID obtiene una orden por external_id e integration_id
func (r *Repository) GetOrderByExternalID(ctx context.Context, externalID string, integrationID uint) (*entities.ProbabilityOrder, error) {
	var order models.Order
	err := r.db.Conn(ctx).
		Preload("Business").
		Preload("Integration.IntegrationType"). // Precargar Integration con IntegrationType para obtener el logo
		Preload("PaymentMethod").
		Preload("OrderStatus").       // Precargar OrderStatus para obtener información del estado de Probability
		Preload("PaymentStatus").     // Precargar PaymentStatus
		Preload("FulfillmentStatus"). // Precargar FulfillmentStatus
		Preload("OrderItems.Product").
		Where("external_id = ? AND integration_id = ?", externalID, integrationID).
		First(&order).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("order not found")
		}
		return nil, err
	}

	return mappers.ToDomainOrder(&order, r.imageURLBase), nil
}

// ───────────────────────────────────────────
//
//	MÉTODOS PARA TABLAS RELACIONADAS
//
// ───────────────────────────────────────────

// CreateOrderItems crea múltiples items de orden
func (r *Repository) CreateOrderItems(ctx context.Context, items []*entities.ProbabilityOrderItem) error {
	if len(items) == 0 {
		return nil
	}

	// Convertir []*entities.ProbabilityOrderItem a []entities.ProbabilityOrderItem para usar el mapper
	domainItems := make([]entities.ProbabilityOrderItem, len(items))
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
func (r *Repository) CreateAddresses(ctx context.Context, addresses []*entities.ProbabilityAddress) error {
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
func (r *Repository) CreatePayments(ctx context.Context, payments []*entities.ProbabilityPayment) error {
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
func (r *Repository) CreateShipments(ctx context.Context, shipments []*entities.ProbabilityShipment) error {
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
func (r *Repository) CreateChannelMetadata(ctx context.Context, metadata *entities.ProbabilityOrderChannelMetadata) error {
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
func (r *Repository) GetProductBySKU(ctx context.Context, businessID uint, sku string) (*entities.Product, error) {
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
func (r *Repository) CreateProduct(ctx context.Context, product *entities.Product) error {
	dbProduct := mappers.ToDBProduct(product)
	if err := r.db.Conn(ctx).Create(dbProduct).Error; err != nil {
		return err
	}
	product.ID = dbProduct.ID
	return nil
}

// GetClientByEmail busca un cliente por Email y BusinessID
func (r *Repository) GetClientByEmail(ctx context.Context, businessID uint, email string) (*entities.Client, error) {
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
func (r *Repository) GetClientByDNI(ctx context.Context, businessID uint, dni string) (*entities.Client, error) {
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
func (r *Repository) CreateClient(ctx context.Context, client *entities.Client) error {
	dbClient := mappers.ToDBClient(client)
	if err := r.db.Conn(ctx).Create(dbClient).Error; err != nil {
		return err
	}
	client.ID = dbClient.ID
	return nil
}

// CountOrdersByClientID cuenta las órdenes de un cliente
func (r *Repository) CountOrdersByClientID(ctx context.Context, clientID uint) (int64, error) {
	var count int64
	err := r.db.Conn(ctx).
		Model(&models.Order{}).
		Where("customer_id = ?", clientID).
		Count(&count).Error
	return count, err
}

// GetLastManualOrderNumber obtiene el último número de secuencia para órdenes manuales
func (r *Repository) GetLastManualOrderNumber(ctx context.Context, businessID uint) (int, error) {
	var lastOrder models.Order
	err := r.db.Conn(ctx).
		Where("business_id = ? AND platform = 'manual' AND order_number LIKE 'prob-%'", businessID).
		Order("order_number DESC").
		First(&lastOrder).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, err
	}

	// Parsear el número de 'prob-XXXX'
	var num int
	_, err = fmt.Sscanf(lastOrder.OrderNumber, "prob-%d", &num)
	if err != nil {
		return 0, nil // Si no se puede parsear, empezamos de nuevo
	}
	return num, nil
}

// CreateOrderError guarda un error ocurrido durante el procesamiento de una orden
func (r *Repository) CreateOrderError(ctx context.Context, orderError *entities.OrderError) error {
	if orderError == nil {
		return fmt.Errorf("orderError cannot be nil")
	}

	dbError := &models.OrderError{
		ExternalID:      orderError.ExternalID,
		IntegrationID:   orderError.IntegrationID,
		BusinessID:      orderError.BusinessID,
		IntegrationType: orderError.IntegrationType,
		Platform:        orderError.Platform,
		ErrorType:       orderError.ErrorType,
		ErrorMessage:    orderError.ErrorMessage,
		ErrorStack:      orderError.ErrorStack,
		RawData:         orderError.RawData,
		Status:          orderError.Status,
		ResolvedAt:      orderError.ResolvedAt,
		ResolvedBy:      orderError.ResolvedBy,
		Resolution:      orderError.Resolution,
	}

	if orderError.Status == "" {
		dbError.Status = "new"
	}

	return r.db.Conn(ctx).Create(dbError).Error
}
