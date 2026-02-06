package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/ports"
	"github.com/secamc93/probability/back/migration/shared/models"
)

type OrderRepository struct {
	*Repository
}

func NewOrderRepository(base *Repository) ports.IOrderRepository {
	return &OrderRepository{Repository: base}
}

// OrderItemJSON representa un item de orden parseado desde JSON
type OrderItemJSON struct {
	ProductID   *string  `json:"product_id"`
	SKU         string   `json:"sku"`
	Name        string   `json:"name"`
	Description *string  `json:"description"`
	Quantity    int      `json:"quantity"`
	UnitPrice   float64  `json:"unit_price"`
	TotalPrice  float64  `json:"total_price"`
	Tax         float64  `json:"tax"`
	TaxRate     *float64 `json:"tax_rate"`
	Discount    float64  `json:"discount"`
}

// GetByID obtiene una orden por su ID y la mapea a OrderData
func (r *OrderRepository) GetByID(ctx context.Context, orderID string) (*ports.OrderData, error) {
	var order models.Order

	// Consultar orden (sin preloads innecesarios)
	err := r.db.Conn(ctx).
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
func (r *OrderRepository) UpdateInvoiceInfo(ctx context.Context, orderID string, invoiceID string, invoiceURL string) error {
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
func (r *OrderRepository) mapToOrderData(order *models.Order) *ports.OrderData {
	if order == nil {
		return nil
	}

	orderData := &ports.OrderData{
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

	// Parsear items desde JSON
	items, err := r.parseOrderItems(order.Items)
	if err != nil {
		r.log.Warn(context.Background()).
			Err(err).
			Str("order_id", order.ID).
			Msg("Failed to parse order items from JSON, using empty items")
		items = []ports.OrderItemData{}
	}
	orderData.Items = items

	return orderData
}

// parseOrderItems parsea el campo Items (JSON) a OrderItemData
func (r *OrderRepository) parseOrderItems(itemsJSON []byte) ([]ports.OrderItemData, error) {
	if len(itemsJSON) == 0 {
		return []ports.OrderItemData{}, nil
	}

	var jsonItems []OrderItemJSON
	if err := json.Unmarshal(itemsJSON, &jsonItems); err != nil {
		return nil, fmt.Errorf("failed to unmarshal items JSON: %w", err)
	}

	items := make([]ports.OrderItemData, len(jsonItems))
	for i, jsonItem := range jsonItems {
		items[i] = ports.OrderItemData{
			ProductID:   jsonItem.ProductID,
			SKU:         jsonItem.SKU,
			Name:        jsonItem.Name,
			Description: jsonItem.Description,
			Quantity:    jsonItem.Quantity,
			UnitPrice:   jsonItem.UnitPrice,
			TotalPrice:  jsonItem.TotalPrice,
			Tax:         jsonItem.Tax,
			TaxRate:     jsonItem.TaxRate,
			Discount:    jsonItem.Discount,
		}
	}

	return items, nil
}
