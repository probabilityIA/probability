package publisher

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/shopify/internal/domain"
	ordersdomain "github.com/secamc93/probability/back/central/services/modules/orders/domain"
	"github.com/secamc93/probability/back/central/shared/log"
	"gorm.io/datatypes"
)

type Publisher struct {
	logger       log.ILogger
	orderUseCase ordersdomain.IOrderMappingUseCase
}

func New(logger log.ILogger, orderUseCase ordersdomain.IOrderMappingUseCase) domain.OrderPublisher {
	return &Publisher{
		logger:       logger,
		orderUseCase: orderUseCase,
	}
}

// Publish publica una orden al sistema (guardando directamente en BD)
func (p *Publisher) Publish(ctx context.Context, order *domain.UnifiedOrder) error {
	// 1. Mapear UnifiedOrder a CanonicalOrderDTO
	// Marshal metadata to JSON for RawData
	rawDataJSON, _ := json.Marshal(order.Metadata)

	// Mapear a CanonicalOrderDTO
	dto := &ordersdomain.CanonicalOrderDTO{
		BusinessID:      order.BusinessID,
		IntegrationID:   order.IntegrationID,
		IntegrationType: "shopify",
		Platform:        "shopify",
		ExternalID:      order.ExternalID,
		OrderNumber:     order.OrderNumber,
		InternalNumber:  order.OrderNumber, // Usamos el mismo por defecto, puede cambiar luego

		// Información financiera
		// NOTA: UnifiedOrder actual no tiene desglose de impuestos/descuentos/envío expuesto en el struct principal
		// Se deberían extraer de Metadata o calcular de items si fuera necesario. Por ahora van en 0.
		Subtotal:     order.TotalAmount, // Asumimos subtotal = total por ahora si no hay desglose
		Tax:          0,
		Discount:     0,
		ShippingCost: 0,
		TotalAmount:  order.TotalAmount,
		Currency:     order.Currency,

		// Cliente
		CustomerName:  order.Customer.Name,
		CustomerEmail: order.Customer.Email,
		CustomerPhone: order.Customer.Phone,
		// CustomerDNI:   order.Customer.Note, // Shopify no tiene campo DNI nativo accesible fácil aqui sin meta

		// Estado
		Status:         p.mapShopifyStatus(order.Status),
		OriginalStatus: order.OriginalStatus, // Use OriginalStatus from UnifiedOrder

		// Fechas
		OccurredAt: order.OccurredAt,
		ImportedAt: time.Now(),

		// Coord
		// ShippingLat/Lng will be set if coords exist below or via pointer logic if needed

		// Metadata del canal (Raw Data)
		ChannelMetadata: &ordersdomain.CanonicalChannelMetadataDTO{
			ChannelSource: "shopify",
			RawData:       datatypes.JSON(rawDataJSON),
			Version:       "1.0",
			ReceivedAt:    time.Now(),
			SyncStatus:    "synced",
			IsLatest:      true,
		},
	}

	// Mapear Dirección de Envío (Shipping Address)
	shippingAddress := ordersdomain.CanonicalAddressDTO{
		Type:       "shipping",
		FirstName:  order.Customer.Name, // Using customer name as recipient by default
		Street:     order.ShippingAddress.Street,
		Street2:    order.ShippingAddress.Address2,
		City:       order.ShippingAddress.City,
		State:      order.ShippingAddress.State,
		Country:    order.ShippingAddress.Country,
		PostalCode: order.ShippingAddress.PostalCode,
	}

	// Map coordinates if available
	if order.ShippingAddress.Coordinates != nil {
		shippingAddress.Latitude = &order.ShippingAddress.Coordinates.Lat
		shippingAddress.Longitude = &order.ShippingAddress.Coordinates.Lng
	}
	dto.Addresses = append(dto.Addresses, shippingAddress)

	// Agregar Dirección de Origen (Hardcoded: Bogotá, per user request)
	originAddress := ordersdomain.CanonicalAddressDTO{
		Type:    "origin",
		City:    "Bogotá",
		State:   "Cundinamarca",
		Country: "Colombia",
		Street:  "Bodega Central", // Default placeholder
	}
	dto.Addresses = append(dto.Addresses, originAddress)

	// Mapear Items
	for _, item := range order.Items {
		canonicalItem := ordersdomain.CanonicalOrderItemDTO{
			ProductSKU:   item.SKU,
			ProductName:  item.Name,
			Quantity:     item.Quantity,
			UnitPrice:    item.UnitPrice,
			TotalPrice:   item.UnitPrice * float64(item.Quantity),
			Currency:     order.Currency,
			ProductTitle: item.Name,
		}
		dto.OrderItems = append(dto.OrderItems, canonicalItem)
	}

	// 2. Guardar usando el caso de uso de Orders
	_, err := p.orderUseCase.MapAndSaveOrder(ctx, dto)
	if err != nil {
		p.logger.Error().
			Err(err).
			Str("order_number", order.OrderNumber).
			Msg("Failed to save order directly via use case")
		return fmt.Errorf("failed to save order: %w", err)
	}

	p.logger.Info().
		Str("order_number", order.OrderNumber).
		Msg("Order saved successfully via direct use case call")

	return nil
}

func (p *Publisher) mapShopifyStatus(shopifyStatus string) string {
	switch shopifyStatus {
	case "paid":
		return "confirmed"
	case "pending":
		return "pending"
	case "refunded":
		return "cancelled"
	case "voided":
		return "cancelled"
	default:
		return "pending"
	}
}
