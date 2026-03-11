package usecasecreateorder

import (
	"context"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
)

// saveRelatedEntities guarda todas las entidades relacionadas (items, addresses, payments, shipments, metadata)
func (uc *UseCaseCreateOrder) saveRelatedEntities(ctx context.Context, order *entities.ProbabilityOrder, dto *dtos.ProbabilityOrderDTO) error {
	if err := uc.saveOrderItems(ctx, order, dto); err != nil {
		return err
	}

	if err := uc.saveAddresses(ctx, order, dto); err != nil {
		return err
	}

	if err := uc.savePayments(ctx, order, dto); err != nil {
		return err
	}

	if err := uc.saveShipments(ctx, order, dto); err != nil {
		return err
	}

	if err := uc.saveChannelMetadata(ctx, order, dto); err != nil {
		return err
	}

	return nil
}

// saveOrderItems guarda los items de la orden
func (uc *UseCaseCreateOrder) saveOrderItems(ctx context.Context, order *entities.ProbabilityOrder, dto *dtos.ProbabilityOrderDTO) error {
	if len(dto.OrderItems) == 0 {
		return nil
	}

	orderItems := make([]*entities.ProbabilityOrderItem, len(dto.OrderItems))
	for i, itemDTO := range dto.OrderItems {
		// Validar/Crear Producto
		product, err := uc.GetOrCreateProduct(ctx, *dto.BusinessID, itemDTO)
		if err != nil {
			return fmt.Errorf("error processing product for item %s: %w", itemDTO.ProductSKU, err)
		}

		var productID *string
		if product != nil {
			productID = &product.ID
		}

		orderItems[i] = &entities.ProbabilityOrderItem{
			OrderID:          order.ID,
			ProductID:        productID,
			ProductSKU:       itemDTO.ProductSKU,
			ProductName:      itemDTO.ProductName,
			ProductTitle:     itemDTO.ProductTitle,
			VariantID:        itemDTO.VariantID,
			Quantity:         itemDTO.Quantity,
			UnitPrice:        itemDTO.UnitPrice,
			TotalPrice:       itemDTO.TotalPrice,
			Currency:         itemDTO.Currency,
			Discount:         itemDTO.Discount,
			DiscountPercent:  itemDTO.DiscountPercent,
			Tax:              itemDTO.Tax,
			TaxRate:          itemDTO.TaxRate,
			UnitPriceBase:            itemDTO.UnitPriceBase,
			UnitPriceBasePresentment: itemDTO.UnitPriceBasePresentment,
			ImageURL:         itemDTO.ImageURL,
			ProductURL:       itemDTO.ProductURL,
			Weight:           itemDTO.Weight,
			RequiresShipping: true,
			IsGiftCard:       false,
			Metadata:         itemDTO.Metadata,
			// Precios en moneda local
			UnitPricePresentment:  itemDTO.UnitPricePresentment,
			TotalPricePresentment: itemDTO.TotalPricePresentment,
			DiscountPresentment:   itemDTO.DiscountPresentment,
			TaxPresentment:        itemDTO.TaxPresentment,
		}
	}

	return uc.repo.CreateOrderItems(ctx, orderItems)
}

// saveAddresses guarda las direcciones de la orden
func (uc *UseCaseCreateOrder) saveAddresses(ctx context.Context, order *entities.ProbabilityOrder, dto *dtos.ProbabilityOrderDTO) error {
	if len(dto.Addresses) == 0 {
		return nil
	}

	addresses := make([]*entities.ProbabilityAddress, len(dto.Addresses))
	for i, addrDTO := range dto.Addresses {
		addresses[i] = &entities.ProbabilityAddress{
			Type:         addrDTO.Type,
			OrderID:      order.ID,
			FirstName:    addrDTO.FirstName,
			LastName:     addrDTO.LastName,
			Company:      addrDTO.Company,
			Phone:        addrDTO.Phone,
			Street:       addrDTO.Street,
			Street2:      addrDTO.Street2,
			City:         addrDTO.City,
			State:        addrDTO.State,
			Country:      addrDTO.Country,
			PostalCode:   addrDTO.PostalCode,
			Latitude:     addrDTO.Latitude,
			Longitude:    addrDTO.Longitude,
			Instructions: addrDTO.Instructions,
			IsDefault:    false,
			Metadata:     addrDTO.Metadata,
		}
	}

	return uc.repo.CreateAddresses(ctx, addresses)
}

// savePayments guarda los pagos de la orden
func (uc *UseCaseCreateOrder) savePayments(ctx context.Context, order *entities.ProbabilityOrder, dto *dtos.ProbabilityOrderDTO) error {
	if len(dto.Payments) == 0 {
		return nil
	}

	payments := make([]*entities.ProbabilityPayment, len(dto.Payments))
	for i, payDTO := range dto.Payments {
		payments[i] = &entities.ProbabilityPayment{
			OrderID:          order.ID,
			PaymentMethodID:  payDTO.PaymentMethodID,
			Amount:           payDTO.Amount,
			Currency:         payDTO.Currency,
			ExchangeRate:     payDTO.ExchangeRate,
			Status:           payDTO.Status,
			PaidAt:           payDTO.PaidAt,
			ProcessedAt:      payDTO.ProcessedAt,
			TransactionID:    payDTO.TransactionID,
			PaymentReference: payDTO.PaymentReference,
			Gateway:          payDTO.Gateway,
			RefundAmount:     payDTO.RefundAmount,
			RefundedAt:       payDTO.RefundedAt,
			FailureReason:    payDTO.FailureReason,
			Metadata:         payDTO.Metadata,
		}
	}

	return uc.repo.CreatePayments(ctx, payments)
}

// saveShipments guarda los envíos de la orden
func (uc *UseCaseCreateOrder) saveShipments(ctx context.Context, order *entities.ProbabilityOrder, dto *dtos.ProbabilityOrderDTO) error {
	if len(dto.Shipments) == 0 {
		return nil
	}

	shipments := make([]*entities.ProbabilityShipment, len(dto.Shipments))
	for i, shipDTO := range dto.Shipments {
		shipments[i] = &entities.ProbabilityShipment{
			OrderID:           &order.ID,
			TrackingNumber:    shipDTO.TrackingNumber,
			TrackingURL:       shipDTO.TrackingURL,
			Carrier:           shipDTO.Carrier,
			CarrierCode:       shipDTO.CarrierCode,
			GuideID:           shipDTO.GuideID,
			GuideURL:          shipDTO.GuideURL,
			Status:            shipDTO.Status,
			ShippedAt:         shipDTO.ShippedAt,
			DeliveredAt:       shipDTO.DeliveredAt,
			ShippingAddressID: shipDTO.ShippingAddressID,
			ShippingCost:      shipDTO.ShippingCost,
			InsuranceCost:     shipDTO.InsuranceCost,
			TotalCost:         shipDTO.TotalCost,
			Weight:            shipDTO.Weight,
			Height:            shipDTO.Height,
			Width:             shipDTO.Width,
			Length:            shipDTO.Length,
			WarehouseID:       shipDTO.WarehouseID,
			WarehouseName:     shipDTO.WarehouseName,
			DriverID:          shipDTO.DriverID,
			DriverName:        shipDTO.DriverName,
			IsLastMile:        shipDTO.IsLastMile,
			EstimatedDelivery: shipDTO.EstimatedDelivery,
			DeliveryNotes:     shipDTO.DeliveryNotes,
			Metadata:          shipDTO.Metadata,
		}
	}

	return uc.repo.CreateShipments(ctx, shipments)
}

// saveChannelMetadata guarda los metadatos del canal
func (uc *UseCaseCreateOrder) saveChannelMetadata(ctx context.Context, order *entities.ProbabilityOrder, dto *dtos.ProbabilityOrderDTO) error {
	if dto.ChannelMetadata == nil {
		return nil
	}

	metadata := &entities.ProbabilityOrderChannelMetadata{
		OrderID:       order.ID,
		ChannelSource: dto.ChannelMetadata.ChannelSource,
		IntegrationID: dto.IntegrationID,
		RawData:       dto.ChannelMetadata.RawData,
		Version:       dto.ChannelMetadata.Version,
		ReceivedAt:    dto.ChannelMetadata.ReceivedAt,
		ProcessedAt:   dto.ChannelMetadata.ProcessedAt,
		IsLatest:      dto.ChannelMetadata.IsLatest,
		LastSyncedAt:  dto.ChannelMetadata.LastSyncedAt,
		SyncStatus:    dto.ChannelMetadata.SyncStatus,
	}

	if metadata.ReceivedAt.IsZero() {
		metadata.ReceivedAt = time.Now()
	}
	if metadata.SyncStatus == "" {
		metadata.SyncStatus = "pending"
	}

	return uc.repo.CreateChannelMetadata(ctx, metadata)
}
