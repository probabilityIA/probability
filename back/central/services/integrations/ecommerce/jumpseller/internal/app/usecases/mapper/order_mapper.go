package mapper

import (
	"fmt"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/canonical"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/jumpseller/internal/domain"
)

const (
	paymentMethodCreditCard   uint = 1
	paymentMethodPaypal       uint = 3
	paymentMethodBankTransfer uint = 4
	paymentMethodCash         uint = 5
	paymentMethodCOD          uint = 6
	paymentMethodMercadoPago  uint = 7
	paymentMethodStripe       uint = 8
)

func MapJumpsellerOrderToProbability(order *domain.JumpsellerOrder, rawJSON []byte) *canonical.ProbabilityOrderDTO {
	now := time.Now()

	totalTax := order.Tax + order.ShippingTax
	discount := order.Discount + order.ShippingDiscount

	customerName := strings.TrimSpace(fmt.Sprintf("%s %s", order.BillingAddress.Name, order.BillingAddress.Surname))
	if customerName == "" {
		customerName = strings.TrimSpace(order.Customer.Name)
	}

	var notes *string
	if order.AdditionalInfo != "" {
		info := order.AdditionalInfo
		notes = &info
	}

	dto := &canonical.ProbabilityOrderDTO{
		IntegrationType: "jumpseller",
		Platform:        "jumpseller",
		ExternalID:      fmt.Sprintf("%d", order.ID),
		OrderNumber:     fmt.Sprintf("%d", order.ID),
		Subtotal:        order.Subtotal,
		Tax:             totalTax,
		Discount:        discount,
		ShippingCost:    order.Shipping,
		TotalAmount:     order.Total,
		Currency:        order.Currency,
		CustomerName:    customerName,
		CustomerEmail:   order.Customer.Email,
		CustomerPhone:   order.Customer.Phone,
		CustomerDNI:     order.BillingAddress.TaxID,
		Status:          mapJumpsellerStatus(order.Status),
		OriginalStatus:  order.Status,
		Notes:           notes,
		OccurredAt:      order.CreatedAt,
		ImportedAt:      now,
	}

	dto.OrderItems = make([]canonical.ProbabilityOrderItemDTO, 0, len(order.Products))
	for _, item := range order.Products {
		productID := fmt.Sprintf("%d", item.ID)

		var variantID *string
		if item.VariantID > 0 {
			vid := fmt.Sprintf("%d", item.VariantID)
			variantID = &vid
		}

		var weight *float64
		if item.Weight > 0 {
			w := item.Weight
			weight = &w
		}

		dto.OrderItems = append(dto.OrderItems, canonical.ProbabilityOrderItemDTO{
			ProductID:    &productID,
			ProductSKU:   item.SKU,
			ProductName:  item.Name,
			ProductTitle: item.Name,
			VariantID:    variantID,
			Quantity:     item.Qty,
			UnitPrice:    item.Price,
			TotalPrice:   item.Price*float64(item.Qty) - item.Discount,
			Currency:     order.Currency,
			Discount:     item.Discount,
			Tax:          item.Tax,
			Weight:       weight,
		})
	}

	dto.Addresses = make([]canonical.ProbabilityAddressDTO, 0, 2)
	dto.Addresses = append(dto.Addresses, mapAddress(order.BillingAddress, "billing"))
	if order.ShippingRequired {
		dto.Addresses = append(dto.Addresses, mapAddress(order.ShippingAddress, "shipping"))
	}

	if order.PaymentMethodName != "" || order.PaymentMethodType != "" {
		paymentStatus := "pending"
		var paidAt *time.Time
		if order.Status == domain.StatusPaid {
			paymentStatus = "completed"
			paid := order.CreatedAt
			paidAt = &paid
		}

		gateway := order.PaymentMethodName
		paymentMethodID := mapJumpsellerPaymentMethod(order.PaymentMethodName, order.PaymentMethodType)

		dto.Payments = append(dto.Payments, canonical.ProbabilityPaymentDTO{
			PaymentMethodID: paymentMethodID,
			Amount:          order.Total,
			Currency:        order.Currency,
			Status:          paymentStatus,
			PaidAt:          paidAt,
			Gateway:         &gateway,
		})

		if paymentMethodID == paymentMethodCOD {
			codTotal := order.Total
			dto.CodTotal = &codTotal
		}
	}

	if order.ShippingRequired || order.ShippingMethod != "" {
		carrier := order.ShippingMethod
		shippingCost := order.Shipping

		shipment := canonical.ProbabilityShipmentDTO{
			Carrier:      &carrier,
			Status:       mapShipmentStatus(order.ShipmentStatus, order.Status),
			ShippingCost: &shippingCost,
		}

		if order.TrackingNumber != "" {
			tracking := order.TrackingNumber
			shipment.TrackingNumber = &tracking
		}
		if order.TrackingCompany != "" {
			company := order.TrackingCompany
			shipment.CarrierCode = &company
		}
		if order.TrackingURL != "" {
			trackingURL := order.TrackingURL
			shipment.TrackingURL = &trackingURL
		}

		dto.Shipments = append(dto.Shipments, shipment)
	}

	if rawJSON != nil {
		dto.ChannelMetadata = &canonical.ProbabilityChannelMetadataDTO{
			ChannelSource: "jumpseller",
			RawData:       rawJSON,
			Version:       "v1",
			ReceivedAt:    now,
			IsLatest:      true,
			SyncStatus:    "synced",
		}
	}

	dto.Invoiceable = strings.EqualFold(order.Currency, "COP")

	return dto
}

func mapAddress(address domain.Address, addressType string) canonical.ProbabilityAddressDTO {
	street := strings.TrimSpace(address.Address)
	if address.StreetNumber != "" {
		street = strings.TrimSpace(street + " " + address.StreetNumber)
	}

	return canonical.ProbabilityAddressDTO{
		Type:       addressType,
		FirstName:  address.Name,
		LastName:   address.Surname,
		Street:     street,
		City:       address.City,
		State:      address.Region,
		Country:    address.Country,
		PostalCode: address.Postal,
		Latitude:   address.Latitude,
		Longitude:  address.Longitude,
	}
}

func mapJumpsellerStatus(status string) string {
	switch status {
	case domain.StatusPendingPayment:
		return "pending"
	case domain.StatusPaid:
		return "paid"
	case domain.StatusCanceled:
		return "cancelled"
	case domain.StatusAbandoned:
		return "abandoned"
	default:
		return strings.ToLower(strings.ReplaceAll(status, " ", "_"))
	}
}

func mapShipmentStatus(shipmentStatus, orderStatus string) string {
	switch strings.ToLower(strings.ReplaceAll(shipmentStatus, " ", "_")) {
	case domain.ShipmentDelivered:
		return "delivered"
	case domain.ShipmentInTransit, "shipped":
		return "in_transit"
	case domain.ShipmentFailed:
		return "failed"
	case domain.ShipmentRequested:
		return "pending"
	}

	if orderStatus == domain.StatusCanceled {
		return "cancelled"
	}
	return "pending"
}

func mapJumpsellerPaymentMethod(name, methodType string) uint {
	m := strings.ToLower(name + " " + methodType)
	switch {
	case strings.Contains(m, "contra") || strings.Contains(m, "cash on delivery") || strings.Contains(m, "cod"):
		return paymentMethodCOD
	case strings.Contains(m, "transfer") || strings.Contains(m, "transferencia"):
		return paymentMethodBankTransfer
	case strings.Contains(m, "efectivo") || strings.Contains(m, "cash"):
		return paymentMethodCash
	case strings.Contains(m, "paypal"):
		return paymentMethodPaypal
	case strings.Contains(m, "stripe"):
		return paymentMethodStripe
	case strings.Contains(m, "mercado"):
		return paymentMethodMercadoPago
	default:
		return paymentMethodCreditCard
	}
}
