package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/dtos"
	"github.com/secamc93/probability/back/migration/shared/models"
)

func detectIsCOD(order *models.Order) bool {
	if order.CodTotal != nil && *order.CodTotal > 0 {
		return true
	}
	if len(order.PaymentDetails) == 0 {
		return false
	}
	var details struct {
		Gateway             string   `json:"gateway"`
		PaymentGatewayNames []string `json:"payment_gateway_names"`
	}
	if err := json.Unmarshal(order.PaymentDetails, &details); err != nil {
		return false
	}
	keywords := []string{"cod", "cash", "contra"}
	gw := strings.ToLower(details.Gateway)
	for _, kw := range keywords {
		if gw != "" && strings.Contains(gw, kw) {
			return true
		}
	}
	for _, name := range details.PaymentGatewayNames {
		n := strings.ToLower(name)
		for _, kw := range keywords {
			if strings.Contains(n, kw) {
				return true
			}
		}
	}
	return false
}

func (r *Repository) CountOrdersOutsideBusiness(ctx context.Context, orderIDs []string, businessID uint) (int64, error) {
	if len(orderIDs) == 0 || businessID == 0 {
		return 0, nil
	}

	var count int64
	err := r.db.Conn(ctx).Model(&models.Order{}).
		Where("id IN ?", orderIDs).
		Where("business_id IS DISTINCT FROM ?", businessID).
		Count(&count).Error
	if err != nil {
		r.log.Error(ctx).Err(err).Uint("business_id", businessID).Msg("Failed to validate order ownership")
		return 0, fmt.Errorf("failed to validate order ownership: %w", err)
	}

	return count, nil
}

func (r *Repository) GetOrderByID(ctx context.Context, orderID string) (*dtos.OrderData, error) {
	var order models.Order

	err := r.db.Conn(ctx).
		Preload("OrderItems.Product").
		Where("id = ?", orderID).
		First(&order).Error

	if err != nil {
		r.log.Error(ctx).Err(err).Str("order_id", orderID).Msg("Failed to get order by ID")
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	orderData := r.mapToOrderData(&order)
	return orderData, nil
}

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

func (r *Repository) mapToOrderData(order *models.Order) *dtos.OrderData {
	if order == nil {
		return nil
	}

	orderData := &dtos.OrderData{
		ID:               order.ID,
		BusinessID:       uint(0),
		IntegrationID:    order.IntegrationID,
		OrderNumber:      order.OrderNumber,
		TotalAmount:      order.TotalAmount,
		Subtotal:         order.Subtotal,
		Tax:              order.Tax,
		Discount:         order.Discount,
		ShippingCost:     order.ShippingCost,
		ShippingDiscount: order.ShippingDiscount,
		FreeShipping:     order.FreeShipping,
		Currency:         order.Currency,
		CustomerName:     order.CustomerName,
		CustomerEmail:    order.CustomerEmail,
		CustomerPhone:    order.CustomerPhone,
		CustomerDNI:      order.CustomerDNI,
		IsPaid:           order.IsPaid,
		IsCOD:            detectIsCOD(order),
		PaymentMethodID:  order.PaymentMethodID,
		Invoiceable:      order.Invoiceable,
		IsTest:           order.IsTest,
		Status:           order.Status,
		CreatedAt:        order.CreatedAt,
	}

	if order.BusinessID != nil {
		orderData.BusinessID = *order.BusinessID
	}

	if order.CustomerID != nil {
		customerIDStr := fmt.Sprintf("%d", *order.CustomerID)
		orderData.CustomerID = &customerIDStr
	}

	if order.ShippingCity != "" {
		orderData.ShippingCity = &order.ShippingCity
	}
	if order.ShippingState != "" {
		orderData.ShippingState = &order.ShippingState
	}
	if order.ShippingCountry != "" {
		orderData.ShippingCountry = &order.ShippingCountry
	}

	if order.OrderTypeID != nil {
		orderData.OrderTypeID = *order.OrderTypeID
	}
	if order.OrderTypeName != "" {
		orderData.OrderTypeName = order.OrderTypeName
	}

	items := make([]dtos.OrderItemData, 0, len(order.OrderItems))
	for _, orderItem := range order.OrderItems {
		item := dtos.OrderItemData{
			ProductID:                orderItem.ProductID,
			SKU:                      "",
			Name:                     "",
			Description:              nil,
			Quantity:                 orderItem.Quantity,
			UnitPrice:                orderItem.UnitPrice,
			UnitPriceBase:            orderItem.UnitPriceBase,
			TotalPrice:               orderItem.TotalPrice,
			Tax:                      orderItem.Tax,
			TaxRate:                  orderItem.TaxRate,
			Discount:                 orderItem.Discount,
			DiscountPercent:          orderItem.DiscountPercent,
			UnitPricePresentment:     orderItem.UnitPricePresentment,
			UnitPriceBasePresentment: orderItem.UnitPriceBasePresentment,
			TotalPricePresentment:    orderItem.TotalPricePresentment,
			DiscountPresentment:      orderItem.DiscountPresentment,
			TaxPresentment:           orderItem.TaxPresentment,
		}

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

func (r *Repository) GetInvoiceableOrders(ctx context.Context, filter dtos.InvoiceableOrdersFilter) ([]*dtos.OrderData, int64, error) {
	filter.Sanitize()
	var orders []models.Order
	var total int64

	if filter.BusinessID == 0 {
		return nil, 0, fmt.Errorf("business_id es requerido para listar ordenes facturables")
	}

	noInvoiceSubquery := "NOT EXISTS (SELECT 1 FROM invoices i WHERE i.order_id = orders.id AND i.deleted_at IS NULL AND i.status <> 'cancelled')"
	base := r.db.Conn(ctx).Model(&models.Order{}).
		Where("business_id = ?", filter.BusinessID).
		Where("invoiceable = ?", true).
		Where("invoice_id IS NULL").
		Where(noInvoiceSubquery)

	if filter.StartDate != nil {
		base = base.Where("created_at >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		base = base.Where("created_at <= ?", *filter.EndDate)
	}
	if filter.OrderNumber != "" {
		base = base.Where("order_number ILIKE ?", "%"+filter.OrderNumber+"%")
	}
	if filter.CustomerName != "" {
		base = base.Where("customer_name ILIKE ?", "%"+filter.CustomerName+"%")
	}
	if filter.CustomerEmail != "" {
		base = base.Where("customer_email ILIKE ?", "%"+filter.CustomerEmail+"%")
	}
	if filter.PaymentStatusID > 0 {
		base = base.Where("payment_status_id = ?", filter.PaymentStatusID)
	}
	if filter.FulfillmentStatusID > 0 {
		base = base.Where("fulfillment_status_id = ?", filter.FulfillmentStatusID)
	}

	if err := base.Count(&total).Error; err != nil {
		r.log.Error(ctx).Err(err).Msg("Failed to count invoiceable orders")
		return nil, 0, fmt.Errorf("failed to count invoiceable orders: %w", err)
	}
	if total == 0 {
		return []*dtos.OrderData{}, 0, nil
	}

	offset := (filter.Page - 1) * filter.PageSize
	if err := base.
		Order(fmt.Sprintf("%s %s", filter.SortColumn(), filter.SortOrder)).
		Offset(offset).
		Limit(filter.PageSize).
		Find(&orders).Error; err != nil {
		r.log.Error(ctx).Err(err).Msg("Failed to get invoiceable orders")
		return nil, 0, fmt.Errorf("failed to get invoiceable orders: %w", err)
	}

	orderDataList := make([]*dtos.OrderData, len(orders))
	for i, order := range orders {
		orderDataList[i] = r.mapToOrderData(&order)
	}
	return orderDataList, total, nil
}
