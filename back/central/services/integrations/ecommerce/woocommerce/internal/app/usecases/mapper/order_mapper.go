package mapper

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/canonical"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce/internal/domain"
)

// MapWooOrderToProbability convierte una orden WooCommerce a el DTO canónico de Probability.
func MapWooOrderToProbability(order *domain.WooCommerceOrder, rawJSON []byte) *canonical.ProbabilityOrderDTO {
	now := time.Now()
	totalAmount := parseFloat(order.Total)
	totalTax := parseFloat(order.TotalTax)
	discount := parseFloat(order.DiscountTotal)
	shippingCost := parseFloat(order.ShippingTotal)
	subtotal := totalAmount - totalTax - shippingCost + discount

	// Customer name from billing
	customerName := strings.TrimSpace(fmt.Sprintf("%s %s", order.Billing.FirstName, order.Billing.LastName))

	// Notes
	var notes *string
	if order.CustomerNote != "" {
		notes = &order.CustomerNote
	}

	// Coupon
	var coupon *string
	if len(order.CouponLines) > 0 {
		codes := make([]string, len(order.CouponLines))
		for i, cl := range order.CouponLines {
			codes[i] = cl.Code
		}
		joined := strings.Join(codes, ", ")
		coupon = &joined
	}

	// Map status
	status := mapWooStatus(order.Status)

	dto := &canonical.ProbabilityOrderDTO{
		IntegrationType: "woocommerce",
		Platform:        "woocommerce",
		ExternalID:      fmt.Sprintf("%d", order.ID),
		OrderNumber:     order.Number,
		Subtotal:        subtotal,
		Tax:             totalTax,
		Discount:        discount,
		ShippingCost:    shippingCost,
		TotalAmount:     totalAmount,
		Currency:        order.Currency,
		CustomerName:    customerName,
		CustomerEmail:   order.Billing.Email,
		CustomerPhone:   order.Billing.Phone,
		Status:          status,
		OriginalStatus:  order.Status,
		Notes:           notes,
		Coupon:          coupon,
		OccurredAt:      order.DateCreated,
		ImportedAt:      now,
	}

	// Order items
	dto.OrderItems = make([]canonical.ProbabilityOrderItemDTO, 0, len(order.LineItems))
	for _, item := range order.LineItems {
		productID := fmt.Sprintf("%d", item.ProductID)
		unitPrice := item.Price
		totalPrice := parseFloat(item.Total)
		tax := parseFloat(item.TotalTax)
		sku := item.SKU

		var variantID *string
		if item.VariationID > 0 {
			vid := fmt.Sprintf("%d", item.VariationID)
			variantID = &vid
		}

		var imageURL *string
		if item.ImageURL != "" {
			imageURL = &item.ImageURL
		}

		dto.OrderItems = append(dto.OrderItems, canonical.ProbabilityOrderItemDTO{
			ProductID:    &productID,
			ProductSKU:   sku,
			ProductName:  item.Name,
			ProductTitle: item.Name,
			VariantID:    variantID,
			Quantity:     item.Quantity,
			UnitPrice:    unitPrice,
			TotalPrice:   totalPrice,
			Currency:     order.Currency,
			Tax:          tax,
			ImageURL:     imageURL,
		})
	}

	// Addresses
	dto.Addresses = make([]canonical.ProbabilityAddressDTO, 0, 2)

	// Billing address
	dto.Addresses = append(dto.Addresses, canonical.ProbabilityAddressDTO{
		Type:       "billing",
		FirstName:  order.Billing.FirstName,
		LastName:   order.Billing.LastName,
		Company:    order.Billing.Company,
		Phone:      order.Billing.Phone,
		Street:     order.Billing.Address1,
		Street2:    order.Billing.Address2,
		City:       order.Billing.City,
		State:      order.Billing.State,
		Country:    order.Billing.Country,
		PostalCode: order.Billing.Postcode,
	})

	// Shipping address
	dto.Addresses = append(dto.Addresses, canonical.ProbabilityAddressDTO{
		Type:       "shipping",
		FirstName:  order.Shipping.FirstName,
		LastName:   order.Shipping.LastName,
		Company:    order.Shipping.Company,
		Phone:      order.Shipping.Phone,
		Street:     order.Shipping.Address1,
		Street2:    order.Shipping.Address2,
		City:       order.Shipping.City,
		State:      order.Shipping.State,
		Country:    order.Shipping.Country,
		PostalCode: order.Shipping.Postcode,
	})

	// Payment
	if order.PaymentMethod != "" {
		paymentStatus := "pending"
		var paidAt *time.Time
		if order.DatePaid != nil {
			paidAt = order.DatePaid
			paymentStatus = "paid"
		}

		gateway := order.PaymentMethod
		dto.Payments = append(dto.Payments, canonical.ProbabilityPaymentDTO{
			Amount:   totalAmount,
			Currency: order.Currency,
			Status:   paymentStatus,
			PaidAt:   paidAt,
			Gateway:  &gateway,
		})
	}

	// Shipments from shipping lines
	for _, sl := range order.ShippingLines {
		carrier := sl.MethodTitle
		carrierCode := sl.MethodID
		shCost := parseFloat(sl.Total)

		dto.Shipments = append(dto.Shipments, canonical.ProbabilityShipmentDTO{
			Carrier:      &carrier,
			CarrierCode:  &carrierCode,
			Status:       mapShipmentStatus(order.Status),
			ShippingCost: &shCost,
		})
	}

	// Channel metadata with raw data
	if rawJSON != nil {
		dto.ChannelMetadata = &canonical.ProbabilityChannelMetadataDTO{
			ChannelSource: "woocommerce",
			RawData:       rawJSON,
			Version:       "v3",
			ReceivedAt:    now,
			IsLatest:      true,
			SyncStatus:    "synced",
		}
	}

	return dto
}

// mapWooStatus mapea el estado de WooCommerce al estado canónico de Probability.
func mapWooStatus(wooStatus string) string {
	switch wooStatus {
	case "pending":
		return "pending"
	case "processing":
		return "paid"
	case "on-hold":
		return "on_hold"
	case "completed":
		return "fulfilled"
	case "cancelled":
		return "cancelled"
	case "refunded":
		return "refunded"
	case "failed":
		return "failed"
	case "trash":
		return "deleted"
	default:
		return wooStatus
	}
}

// mapShipmentStatus mapea el estado de la orden WooCommerce a un estado de envío.
func mapShipmentStatus(wooStatus string) string {
	switch wooStatus {
	case "completed":
		return "delivered"
	case "processing":
		return "pending"
	case "on-hold":
		return "pending"
	case "cancelled", "refunded", "failed":
		return "cancelled"
	default:
		return "pending"
	}
}

func parseFloat(s string) float64 {
	v, _ := strconv.ParseFloat(s, 64)
	return v
}
