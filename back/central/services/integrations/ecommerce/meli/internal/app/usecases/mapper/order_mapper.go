package mapper

import (
	"fmt"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/canonical"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
)

// MapMeliOrderToProbability convierte una orden de MercadoLibre al DTO canónico de Probability.
// shippingDetail puede ser nil si no se pudo obtener el detalle del envío.
func MapMeliOrderToProbability(order *domain.MeliOrder, shippingDetail *domain.MeliShippingDetail, rawJSON []byte) *canonical.ProbabilityOrderDTO {
	now := time.Now()

	// Calcular shipping cost
	shippingCost := 0.0
	if shippingDetail != nil && shippingDetail.ShippingOption != nil {
		shippingCost = shippingDetail.ShippingOption.Cost
	}

	// Subtotal = TotalAmount - ShippingCost + CouponAmount
	// MeLi: TotalAmount ya incluye descuentos pero no shipping
	subtotal := order.TotalAmount + order.CouponAmount

	// Customer
	customerName := strings.TrimSpace(fmt.Sprintf("%s %s", order.Buyer.FirstName, order.Buyer.LastName))
	if customerName == "" {
		customerName = order.Buyer.Nickname
	}

	customerPhone := ""
	if order.Buyer.Phone.Number != "" {
		if order.Buyer.Phone.AreaCode != "" {
			customerPhone = order.Buyer.Phone.AreaCode + order.Buyer.Phone.Number
		} else {
			customerPhone = order.Buyer.Phone.Number
		}
	}

	customerDNI := ""
	if order.Buyer.BillingInfo != nil {
		customerDNI = order.Buyer.BillingInfo.DocNumber
	}

	// Coupon
	var coupon *string
	if order.CouponID != nil && *order.CouponID != "" {
		coupon = order.CouponID
	}

	// Status mapping
	status := mapMeliOrderStatus(order.Status)

	dto := &canonical.ProbabilityOrderDTO{
		IntegrationType: "mercado_libre",
		Platform:        "mercadolibre",
		ExternalID:      fmt.Sprintf("%d", order.ID),
		OrderNumber:     fmt.Sprintf("%d", order.ID),
		Subtotal:        subtotal,
		Tax:             0, // MeLi no separa impuestos
		Discount:        order.CouponAmount,
		ShippingCost:    shippingCost,
		TotalAmount:     order.TotalAmount,
		Currency:        order.CurrencyID,
		CustomerName:    customerName,
		CustomerEmail:   order.Buyer.Email,
		CustomerPhone:   customerPhone,
		CustomerDNI:     customerDNI,
		Status:          status,
		OriginalStatus:  order.Status,
		Coupon:          coupon,
		OccurredAt:      order.DateCreated,
		ImportedAt:      now,
	}

	// Order items
	dto.OrderItems = make([]canonical.ProbabilityOrderItemDTO, 0, len(order.OrderItems))
	for _, item := range order.OrderItems {
		productID := item.Item.ID

		var variantID *string
		if item.Item.VariationID != nil {
			vid := fmt.Sprintf("%d", *item.Item.VariationID)
			variantID = &vid
		}

		sku := ""
		if item.Item.SellerSKU != nil {
			sku = *item.Item.SellerSKU
		}

		dto.OrderItems = append(dto.OrderItems, canonical.ProbabilityOrderItemDTO{
			ProductID:    &productID,
			ProductSKU:   sku,
			ProductName:  item.Item.Title,
			ProductTitle: item.Item.Title,
			VariantID:    variantID,
			Quantity:     item.Quantity,
			UnitPrice:    item.UnitPrice,
			TotalPrice:   item.UnitPrice * float64(item.Quantity),
			Currency:     item.Currency,
		})
	}

	// Addresses (from shipping detail)
	if shippingDetail != nil && shippingDetail.ReceiverAddress != nil {
		addr := shippingDetail.ReceiverAddress
		street := addr.StreetName
		if addr.StreetNumber != "" {
			street = addr.StreetName + " " + addr.StreetNumber
		}

		dto.Addresses = append(dto.Addresses, canonical.ProbabilityAddressDTO{
			Type:       "shipping",
			Street:     street,
			City:       addr.City.Name,
			State:      addr.State.Name,
			Country:    addr.Country.Name,
			PostalCode: addr.ZipCode,
			Latitude:   addr.Latitude,
			Longitude:  addr.Longitude,
		})

		if addr.Comment != "" {
			instructions := addr.Comment
			dto.Addresses[len(dto.Addresses)-1].Instructions = &instructions
		}
	}

	// Payments
	dto.Payments = make([]canonical.ProbabilityPaymentDTO, 0, len(order.Payments))
	for _, p := range order.Payments {
		gateway := p.PaymentMethodID
		paymentStatus := mapMeliPaymentStatus(p.Status)

		dto.Payments = append(dto.Payments, canonical.ProbabilityPaymentDTO{
			Amount:   p.TransactionAmount,
			Currency: p.CurrencyID,
			Status:   paymentStatus,
			PaidAt:   p.DateApproved,
			Gateway:  &gateway,
		})
	}

	// Shipment
	if shippingDetail != nil {
		carrier := "mercadoenvios"
		shStatus := mapMeliShippingStatus(shippingDetail.Status)

		shipment := canonical.ProbabilityShipmentDTO{
			Carrier:      &carrier,
			Status:       shStatus,
			ShippingCost: &shippingCost,
		}

		if shippingDetail.ShippingOption != nil {
			optName := shippingDetail.ShippingOption.Name
			shipment.CarrierCode = &optName
			if shippingDetail.ShippingOption.EstimatedDeliveryTime != nil {
				shipment.EstimatedDelivery = shippingDetail.ShippingOption.EstimatedDeliveryTime.Date
			}
		}

		dto.Shipments = append(dto.Shipments, shipment)
	}

	// Channel metadata
	if rawJSON != nil {
		dto.ChannelMetadata = &canonical.ProbabilityChannelMetadataDTO{
			ChannelSource: "mercadolibre",
			RawData:       rawJSON,
			Version:       "v2",
			ReceivedAt:    now,
			IsLatest:      true,
			SyncStatus:    "synced",
		}
	}

	return dto
}

// mapMeliOrderStatus mapea el estado de orden MeLi al estado canónico de Probability.
func mapMeliOrderStatus(meliStatus string) string {
	switch meliStatus {
	case "confirmed":
		return "pending"
	case "payment_required":
		return "pending"
	case "payment_in_process":
		return "pending"
	case "paid":
		return "paid"
	case "partially_paid":
		return "paid"
	case "cancelled":
		return "cancelled"
	default:
		return meliStatus
	}
}

// mapMeliPaymentStatus mapea el estado de pago MeLi al estado canónico.
func mapMeliPaymentStatus(meliStatus string) string {
	switch meliStatus {
	case "approved":
		return "paid"
	case "pending", "in_process", "in_mediation":
		return "pending"
	case "rejected":
		return "failed"
	case "refunded":
		return "refunded"
	case "cancelled":
		return "cancelled"
	default:
		return meliStatus
	}
}

// mapMeliShippingStatus mapea el estado de envío MeLi al estado canónico.
func mapMeliShippingStatus(meliStatus string) string {
	switch meliStatus {
	case "ready_to_ship":
		return "pending"
	case "shipped":
		return "shipped"
	case "delivered":
		return "delivered"
	case "not_delivered":
		return "failed"
	case "cancelled":
		return "cancelled"
	default:
		return "pending"
	}
}
