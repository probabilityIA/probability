package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/dtos"
	"github.com/secamc93/probability/back/migration/shared/models"
)

// GetByID obtiene una orden por su ID y la mapea a OrderData
func (r *Repository) GetOrderByID(ctx context.Context, orderID string) (*dtos.OrderData, error) {
	var order models.Order

	// Consultar orden CON Preload de OrderItems para obtener los product_id correctos
	// IMPORTANTE: OrderItems tiene los product_id de la tabla products (PRD_xxx)
	// mientras que el JSONB Items tiene los IDs externos de las plataformas (Shopify, etc)
	err := r.db.Conn(ctx).
		Preload("OrderItems.Product").  // Preload items y productos
		Where("id = ?", orderID).
		First(&order).Error

	if err != nil {
		r.log.Error(ctx).Err(err).Str("order_id", orderID).Msg("Failed to get order by ID")
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	// Mapear a OrderData
	orderData := r.mapToOrderData(&order)
	return orderData, nil
}

// UpdateInvoiceInfo actualiza la información de factura en una orden
func (r *Repository) UpdateOrderInvoiceInfo(ctx context.Context, orderID string, invoiceID string, invoiceURL string) error {
	result := r.db.Conn(ctx).Model(&models.Order{}).
		Where("id = ?", orderID).
		Updates(map[string]interface{}{
			"invoice_id":  invoiceID,
			"invoice_url": invoiceURL,
		})

	if result.Error != nil {
		r.log.Error(ctx).
			Err(result.Error).
			Str("order_id", orderID).
			Str("invoice_id", invoiceID).
			Msg("Failed to update invoice info")
		return fmt.Errorf("failed to update invoice info: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("order not found: %s", orderID)
	}

	return nil
}

// mapToOrderData convierte un modelo GORM Order a OrderData del dominio
func (r *Repository) mapToOrderData(order *models.Order) *dtos.OrderData {
	if order == nil {
		return nil
	}

	orderData := &dtos.OrderData{
		ID:              order.ID,
		BusinessID:      uint(0),
		IntegrationID:   order.IntegrationID,
		OrderNumber:     order.OrderNumber,
		TotalAmount:     order.TotalAmount,
		Subtotal:        order.Subtotal,
		Tax:             order.Tax,
		Discount:        order.Discount,
		ShippingCost:    order.ShippingCost,
		Currency:        order.Currency,
		CustomerName:    order.CustomerName,
		CustomerEmail:   order.CustomerEmail,
		CustomerPhone:   order.CustomerPhone,
		CustomerDNI:     order.CustomerDNI,
		IsPaid:          order.IsPaid,
		PaymentMethodID: order.PaymentMethodID,
		Invoiceable:     order.Invoiceable,
		Status:          order.Status,
		CreatedAt:       order.CreatedAt,
	}

	// BusinessID puede ser nil en el modelo, asignar valor
	if order.BusinessID != nil {
		orderData.BusinessID = *order.BusinessID
	}

	// CustomerID
	if order.CustomerID != nil {
		customerIDStr := fmt.Sprintf("%d", *order.CustomerID)
		orderData.CustomerID = &customerIDStr
	}

	// Shipping address
	if order.ShippingCity != "" {
		orderData.ShippingCity = &order.ShippingCity
	}
	if order.ShippingState != "" {
		orderData.ShippingState = &order.ShippingState
	}
	if order.ShippingCountry != "" {
		orderData.ShippingCountry = &order.ShippingCountry
	}

	// Order Type - usar campos directos del modelo
	if order.OrderTypeID != nil {
		orderData.OrderTypeID = *order.OrderTypeID
	}
	if order.OrderTypeName != "" {
		orderData.OrderTypeName = order.OrderTypeName
	}

	// Mapear items desde la tabla normalizada order_items (NO desde JSONB)
	// IMPORTANTE: order_items contiene los product_id correctos de la tabla products
	// mientras que el JSONB Items contiene IDs externos de las plataformas
	items := make([]dtos.OrderItemData, 0, len(order.OrderItems))
	for _, orderItem := range order.OrderItems {
		item := dtos.OrderItemData{
			ProductID:   orderItem.ProductID,  // ✅ ID correcto de tabla products (PRD_xxx)
			SKU:         "",                    // Se obtendrá del Product si existe
			Name:        "",                    // Se obtendrá del Product si existe
			Description: nil,
			Quantity:    orderItem.Quantity,
			UnitPrice:   orderItem.UnitPrice,
			TotalPrice:  orderItem.TotalPrice,
			Tax:         orderItem.Tax,
			TaxRate:     orderItem.TaxRate,
			Discount:    orderItem.Discount,
		}

		// Obtener información del producto desde la relación (si existe y no está soft-deleted)
		if orderItem.Product.ID != "" {
			item.SKU = orderItem.Product.SKU
			item.Name = orderItem.Product.Name
			if orderItem.Product.Description != "" {
				desc := orderItem.Product.Description
				item.Description = &desc
			}
		}

		items = append(items, item)
	}
	orderData.Items = items

	return orderData
}

// GetInvoiceableOrders obtiene órdenes facturables paginadas
// Filtra por:
// - business_id (multi-tenant isolation)
//   - Si businessID = 0 (super admin): retorna órdenes de TODOS los businesses
//   - Si businessID != 0 (usuario normal): retorna solo órdenes de ese business
// - invoiceable = true
// - invoice_id IS NULL (no facturadas previamente)
func (r *Repository) GetInvoiceableOrders(ctx context.Context, businessID uint, page, pageSize int) ([]*dtos.OrderData, int64, error) {
	var orders []models.Order
	var total int64

	// Validar parámetros de paginación
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// Construir query base
	countQuery := r.db.Conn(ctx).Model(&models.Order{}).
		Where("invoiceable = ?", true).
		Where("invoice_id IS NULL")

	// Super admin (businessID = 0): ver todas las órdenes de todos los businesses
	// Usuario normal: solo su business
	if businessID != 0 {
		countQuery = countQuery.Where("business_id = ?", businessID)
		r.log.Debug(ctx).
			Uint("business_id", businessID).
			Msg("Filtering by specific business (normal user)")
	} else {
		r.log.Debug(ctx).
			Msg("Super admin mode: listing orders from ALL businesses")
	}

	if err := countQuery.Count(&total).Error; err != nil {
		r.log.Error(ctx).
			Err(err).
			Uint("business_id", businessID).
			Msg("Failed to count invoiceable orders")
		return nil, 0, fmt.Errorf("failed to count invoiceable orders: %w", err)
	}

	// Si no hay órdenes, retornar inmediatamente
	if total == 0 {
		return []*dtos.OrderData{}, 0, nil
	}

	// Construir query de resultados
	offset := (page - 1) * pageSize
	resultsQuery := r.db.Conn(ctx).
		Where("invoiceable = ?", true).
		Where("invoice_id IS NULL").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize)

	// Aplicar filtro de business si no es super admin
	if businessID != 0 {
		resultsQuery = resultsQuery.Where("business_id = ?", businessID)
	}

	err := resultsQuery.Find(&orders).Error

	if err != nil {
		r.log.Error(ctx).
			Err(err).
			Uint("business_id", businessID).
			Msg("Failed to get invoiceable orders")
		return nil, 0, fmt.Errorf("failed to get invoiceable orders: %w", err)
	}

	// Mapear a OrderData
	orderDataList := make([]*dtos.OrderData, len(orders))
	for i, order := range orders {
		orderDataList[i] = r.mapToOrderData(&order)
	}

	r.log.Info(ctx).
		Uint("business_id", businessID).
		Bool("is_super_admin", businessID == 0).
		Int64("total", total).
		Int("page", page).
		Int("page_size", pageSize).
		Int("returned", len(orderDataList)).
		Msg("Retrieved invoiceable orders")

	return orderDataList, total, nil
}
