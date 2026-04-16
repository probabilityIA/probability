package app

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/entities"
	"github.com/secamc93/probability/back/central/shared/log"
)

func (uc *UseCase) ProcessOrderEvent(ctx context.Context, event dtos.OrderEventDTO) error {
	ctx = log.WithFunctionCtx(ctx, "ProcessOrderEvent")

	customerID, err := uc.resolveCustomerID(ctx, event)
	if err != nil || customerID == 0 {
		uc.log.Warn(ctx).
			Str("order_id", event.OrderID).
			Str("phone", event.CustomerPhone).
			Msg("could not resolve customer_id, skipping")
		return nil
	}

	switch event.EventType {
	case "order.created":
		return uc.handleOrderCreated(ctx, customerID, event)
	case "order.status_changed":
		return uc.handleStatusChanged(ctx, customerID, event)
	case "order.updated":
		return uc.handleOrderUpdated(ctx, customerID, event)
	}

	return nil
}

func (uc *UseCase) resolveCustomerID(ctx context.Context, event dtos.OrderEventDTO) (uint, error) {
	if event.CustomerID != nil && *event.CustomerID > 0 {
		return *event.CustomerID, nil
	}

	if event.CustomerDNI != "" {
		client, err := uc.repo.FindClientByDNI(ctx, event.BusinessID, event.CustomerDNI)
		if err != nil {
			return 0, err
		}
		if client != nil {
			return client.ID, nil
		}
	}

	if event.CustomerPhone != "" {
		client, err := uc.repo.FindClientByPhone(ctx, event.BusinessID, event.CustomerPhone)
		if err != nil {
			return 0, err
		}
		if client != nil {
			return client.ID, nil
		}
	}

	if event.CustomerEmail != "" {
		client, err := uc.repo.FindClientByEmail(ctx, event.BusinessID, event.CustomerEmail)
		if err != nil {
			return 0, err
		}
		if client != nil {
			return client.ID, nil
		}
	}

	return 0, nil
}

func (uc *UseCase) syncClientData(ctx context.Context, customerID uint, event dtos.OrderEventDTO) {
	client, err := uc.repo.GetByID(ctx, event.BusinessID, customerID)
	if err != nil || client == nil {
		return
	}

	updates := map[string]any{}

	if event.CustomerName != "" && event.CustomerName != client.Name {
		updates["name"] = event.CustomerName
	}

	if event.CustomerEmail != "" {
		if client.Email == nil || *client.Email != event.CustomerEmail {
			updates["email"] = event.CustomerEmail
		}
	}

	if event.CustomerPhone != "" && event.CustomerPhone != client.Phone {
		updates["phone"] = event.CustomerPhone
	}

	if event.CustomerDNI != "" {
		if client.Dni == nil || *client.Dni != event.CustomerDNI {
			updates["dni"] = event.CustomerDNI
		}
	}

	if len(updates) == 0 {
		return
	}

	if err := uc.repo.UpdateClientFields(ctx, customerID, updates); err != nil {
		uc.log.Error(ctx).Err(err).Uint("customer_id", customerID).Msg("failed to sync client data from order")
	}
}

func (uc *UseCase) handleOrderCreated(ctx context.Context, customerID uint, event dtos.OrderEventDTO) error {
	uc.syncClientData(ctx, customerID, event)

	paidCount := 0
	if event.IsPaid {
		paidCount = 1
	}
	inProgress := 1
	if isTerminalStatus(event.Status) {
		inProgress = 0
	}
	delivered := 0
	if event.Status == "delivered" {
		delivered = 1
		inProgress = 0
	}
	cancelled := 0
	if isCancelledStatus(event.Status) {
		cancelled = 1
		inProgress = 0
	}

	summary := &entities.CustomerSummary{
		CustomerID:        customerID,
		BusinessID:        event.BusinessID,
		TotalOrders:       1,
		DeliveredOrders:   delivered,
		CancelledOrders:   cancelled,
		InProgressOrders:  inProgress,
		TotalSpent:        event.TotalAmount,
		TotalPaidOrders:   paidCount,
		AvgDeliveryScore:  event.DeliveryProbability,
		FirstOrderAt:      &event.OrderedAt,
		LastOrderAt:       &event.OrderedAt,
		PreferredPlatform: event.Platform,
	}

	if err := uc.repo.UpsertCustomerSummary(ctx, summary); err != nil {
		uc.log.Error(ctx).Err(err).Uint("customer_id", customerID).Msg("failed to upsert summary")
	}

	if event.ShippingCity != "" || event.ShippingStreet != "" {
		addr := &entities.CustomerAddress{
			CustomerID: customerID,
			BusinessID: event.BusinessID,
			Street:     event.ShippingStreet,
			City:       event.ShippingCity,
			State:      event.ShippingState,
			Country:    event.ShippingCountry,
			PostalCode: event.ShippingPostalCode,
			Latitude:   event.ShippingLat,
			Longitude:  event.ShippingLng,
			TimesUsed:  1,
			LastUsedAt: event.OrderedAt,
		}
		if err := uc.repo.UpsertCustomerAddress(ctx, addr); err != nil {
			uc.log.Error(ctx).Err(err).Uint("customer_id", customerID).Msg("failed to upsert address")
		}
	}

	for _, item := range event.Items {
		productID := ""
		if item.ProductID != nil {
			productID = *item.ProductID
		}
		if productID == "" {
			continue
		}

		product := &entities.CustomerProductHistory{
			CustomerID:     customerID,
			BusinessID:     event.BusinessID,
			ProductID:      productID,
			ProductName:    item.ProductName,
			ProductSKU:     item.ProductSKU,
			ProductImage:   item.ProductImage,
			TimesOrdered:   1,
			TotalQuantity:  item.Quantity,
			TotalSpent:     item.TotalPrice,
			FirstOrderedAt: event.OrderedAt,
			LastOrderedAt:  event.OrderedAt,
		}
		if err := uc.repo.UpsertCustomerProductHistory(ctx, product); err != nil {
			uc.log.Error(ctx).Err(err).Str("product_id", productID).Msg("failed to upsert product history")
		}

		orderItem := &entities.CustomerOrderItem{
			CustomerID:   customerID,
			BusinessID:   event.BusinessID,
			OrderID:      event.OrderID,
			OrderNumber:  event.OrderNumber,
			ProductID:    item.ProductID,
			ProductName:  item.ProductName,
			ProductSKU:   item.ProductSKU,
			ProductImage: item.ProductImage,
			Quantity:     item.Quantity,
			UnitPrice:    item.UnitPrice,
			TotalPrice:   item.TotalPrice,
			OrderStatus:  event.Status,
			OrderedAt:    event.OrderedAt,
		}
		if err := uc.repo.UpsertCustomerOrderItem(ctx, orderItem); err != nil {
			uc.log.Error(ctx).Err(err).Str("order_id", event.OrderID).Msg("failed to upsert order item")
		}
	}

	return nil
}

func (uc *UseCase) handleStatusChanged(ctx context.Context, customerID uint, event dtos.OrderEventDTO) error {
	prev := event.PreviousStatus
	curr := event.CurrentStatus

	summary := &entities.CustomerSummary{
		CustomerID: customerID,
		BusinessID: event.BusinessID,
	}

	if curr == "delivered" && prev != "delivered" {
		summary.DeliveredOrders = 1
		if !isTerminalStatus(prev) {
			summary.InProgressOrders = -1
		}
	}

	if isCancelledStatus(curr) && !isCancelledStatus(prev) {
		summary.CancelledOrders = 1
		if !isTerminalStatus(prev) {
			summary.InProgressOrders = -1
		}
		summary.TotalSpent = -event.TotalAmount
	}

	if event.IsPaid {
		summary.TotalPaidOrders = 1
	}

	if event.DeliveryProbability > 0 {
		summary.AvgDeliveryScore = event.DeliveryProbability
	}

	now := time.Now()
	summary.LastOrderAt = &now

	if err := uc.repo.UpsertCustomerSummary(ctx, summary); err != nil {
		uc.log.Error(ctx).Err(err).Uint("customer_id", customerID).Msg("failed to update summary on status change")
	}

	if curr != "" {
		if err := uc.repo.UpdateOrderItemsStatus(ctx, event.OrderID, curr); err != nil {
			uc.log.Error(ctx).Err(err).Str("order_id", event.OrderID).Msg("failed to update order items status")
		}
	}

	return nil
}

func (uc *UseCase) handleOrderUpdated(ctx context.Context, customerID uint, event dtos.OrderEventDTO) error {
	if event.Status != "" {
		if err := uc.repo.UpdateOrderItemsStatus(ctx, event.OrderID, event.Status); err != nil {
			uc.log.Error(ctx).Err(err).Str("order_id", event.OrderID).Msg("failed to update order items status")
		}
	}
	return nil
}

func isTerminalStatus(status string) bool {
	switch status {
	case "delivered", "cancelled", "voided", "refunded", "partially_refunded":
		return true
	}
	return false
}

func isCancelledStatus(status string) bool {
	switch status {
	case "cancelled", "voided", "refunded":
		return true
	}
	return false
}
