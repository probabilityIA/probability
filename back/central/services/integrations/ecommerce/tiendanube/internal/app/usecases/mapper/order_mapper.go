package mapper

import (
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/canonical"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiendanube/internal/domain"
)

const (
	paymentMethodCreditCard   uint = 1
	paymentMethodBankTransfer uint = 4
	paymentMethodCash         uint = 5
	paymentMethodCOD          uint = 6
	paymentMethodMercadoPago  uint = 7
)

func parseTime(raw string) time.Time {
	if strings.TrimSpace(raw) == "" {
		return time.Time{}
	}
	for _, layout := range []string{time.RFC3339, "2006-01-02T15:04:05-0700", "2006-01-02 15:04:05"} {
		if parsed, err := time.Parse(layout, raw); err == nil {
			return parsed
		}
	}
	return time.Time{}
}

func MapTiendanubeOrderToProbability(order *domain.TiendanubeOrder, rawJSON []byte) *canonical.ProbabilityOrderDTO {
	now := time.Now()

	customerName := strings.TrimSpace(order.ContactName)
	if customerName == "" {
		customerName = strings.TrimSpace(order.ShippingAddress.FirstName + " " + order.ShippingAddress.LastName)
	}
	if customerName == "" {
		customerName = strings.TrimSpace(order.BillingAddress.Name)
	}

	var notes *string
	if strings.TrimSpace(order.Note) != "" {
		note := order.Note
		notes = &note
	}

	occurredAt := parseTime(order.CreatedAt)
	if occurredAt.IsZero() {
		occurredAt = now
	}

	dto := &canonical.ProbabilityOrderDTO{
		IntegrationType: "tiendanube",
		Platform:        "tiendanube",
		ExternalID:      order.ID,
		OrderNumber:     order.Number,
		Subtotal:        order.Subtotal,
		Discount:        order.Discount,
		ShippingCost:    order.ShippingCost,
		TotalAmount:     order.Total,
		Currency:        order.Currency,
		CustomerName:    customerName,
		CustomerEmail:   order.ContactEmail,
		CustomerPhone:   order.ContactPhone,
		CustomerDNI:     order.ContactIdentification,
		Status:          MapOrderStatus(order.Status, order.PaymentStatus, order.ShippingStatus),
		OriginalStatus:  order.Status,
		Notes:           notes,
		OccurredAt:      occurredAt,
		ImportedAt:      now,
	}

	dto.OrderItems = make([]canonical.ProbabilityOrderItemDTO, 0, len(order.Items))
	for _, item := range order.Items {
		productID := item.ProductID
		var variantID *string
		if item.VariantID != "" && item.VariantID != "0" {
			vid := item.VariantID
			variantID = &vid
		}
		var weight *float64
		if item.Weight > 0 {
			w := item.Weight
			weight = &w
		}
		var imageURL *string
		if item.ImageURL != "" {
			img := item.ImageURL
			imageURL = &img
		}

		dto.OrderItems = append(dto.OrderItems, canonical.ProbabilityOrderItemDTO{
			ProductID:    &productID,
			ProductSKU:   item.SKU,
			ProductName:  item.Name,
			ProductTitle: item.Name,
			VariantID:    variantID,
			Quantity:     item.Quantity,
			UnitPrice:    item.Price,
			TotalPrice:   item.Price * float64(item.Quantity),
			Currency:     order.Currency,
			ImageURL:     imageURL,
			Weight:       weight,
		})
	}

	dto.Addresses = make([]canonical.ProbabilityAddressDTO, 0, 2)
	if hasAddress(order.BillingAddress) {
		dto.Addresses = append(dto.Addresses, mapAddress(order.BillingAddress, "billing"))
	}
	if hasAddress(order.ShippingAddress) {
		dto.Addresses = append(dto.Addresses, mapAddress(order.ShippingAddress, "shipping"))
	}

	if order.Gateway != "" || order.GatewayName != "" {
		paymentStatus := mapPaymentStatus(order.PaymentStatus)
		var paidAt *time.Time
		if paid := parseTime(order.PaidAt); !paid.IsZero() {
			paidAt = &paid
		}

		gateway := order.GatewayName
		if gateway == "" {
			gateway = order.Gateway
		}
		paymentMethodID := MapPaymentMethod(order.Gateway, order.GatewayName)

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

	if order.ShippingOption != "" || order.TrackingNumber != "" {
		carrier := order.ShippingOption
		shippingCost := order.ShippingCost

		shipment := canonical.ProbabilityShipmentDTO{
			Carrier:      &carrier,
			Status:       MapShipmentStatus(order.ShippingStatus),
			ShippingCost: &shippingCost,
		}
		if order.TrackingNumber != "" {
			tracking := order.TrackingNumber
			shipment.TrackingNumber = &tracking
		}
		if order.TrackingURL != "" {
			trackingURL := order.TrackingURL
			shipment.TrackingURL = &trackingURL
		}

		dto.Shipments = append(dto.Shipments, shipment)
	}

	if rawJSON != nil {
		dto.ChannelMetadata = &canonical.ProbabilityChannelMetadataDTO{
			ChannelSource: "tiendanube",
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

func hasAddress(address domain.TiendanubeAddress) bool {
	return strings.TrimSpace(address.Street) != "" ||
		strings.TrimSpace(address.City) != "" ||
		strings.TrimSpace(address.Locality) != ""
}

func mapAddress(address domain.TiendanubeAddress, addressType string) canonical.ProbabilityAddressDTO {
	street := strings.TrimSpace(address.Street)
	if address.Number != "" {
		street = strings.TrimSpace(street + " " + address.Number)
	}

	firstName := address.FirstName
	lastName := address.LastName
	if firstName == "" && lastName == "" && address.Name != "" {
		parts := strings.SplitN(strings.TrimSpace(address.Name), " ", 2)
		firstName = parts[0]
		if len(parts) > 1 {
			lastName = parts[1]
		}
	}

	city := address.City
	if city == "" {
		city = address.Locality
	}

	return canonical.ProbabilityAddressDTO{
		Type:       addressType,
		FirstName:  firstName,
		LastName:   lastName,
		Phone:      address.Phone,
		Street:     street,
		Street2:    address.Floor,
		City:       city,
		State:      address.Province,
		Country:    address.Country,
		PostalCode: address.Zipcode,
	}
}
