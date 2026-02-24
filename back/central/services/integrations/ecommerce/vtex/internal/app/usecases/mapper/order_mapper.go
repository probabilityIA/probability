package mapper

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/datatypes"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/canonical"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/vtex/internal/domain"
)

// MapVTEXOrderToProbability convierte una orden de VTEX al DTO canónico de Probability.
// VTEX maneja valores en centavos (int). Dividir por 100 para obtener el valor real.
func MapVTEXOrderToProbability(order *domain.VTEXOrder, rawJSON []byte) *canonical.ProbabilityOrderDTO {
	now := time.Now()

	// Valores monetarios: VTEX usa centavos
	totalAmount := centavosToFloat(order.Value)
	subtotal := centavosToFloat(order.TotalItems)
	discount := centavosToFloat(order.TotalDiscount)
	shippingCost := centavosToFloat(order.TotalFreight)

	// Si TotalDiscount es negativo en VTEX, convertir a positivo para canonical
	if discount < 0 {
		discount = -discount
	}

	// Currency: extraer del primer item de Totals o usar default
	currency := extractCurrency(order)

	// Customer
	customerName := ""
	customerEmail := ""
	customerPhone := ""
	customerDNI := ""
	if order.ClientProfileData != nil {
		cp := order.ClientProfileData
		if cp.IsCorporate && cp.CorporateName != "" {
			customerName = cp.CorporateName
		} else {
			customerName = strings.TrimSpace(fmt.Sprintf("%s %s", cp.FirstName, cp.LastName))
		}
		customerEmail = cp.Email
		customerPhone = cp.Phone
		customerDNI = cp.Document
	}

	// Status mapping
	status := mapVTEXOrderStatus(order.Status)

	dto := &canonical.ProbabilityOrderDTO{
		IntegrationType: "vtex",
		Platform:        "vtex",
		ExternalID:      order.OrderID,
		OrderNumber:     order.Sequence,
		Subtotal:        subtotal,
		Tax:             extractTax(order),
		Discount:        discount,
		ShippingCost:    shippingCost,
		TotalAmount:     totalAmount,
		Currency:        currency,
		CustomerName:    customerName,
		CustomerEmail:   customerEmail,
		CustomerPhone:   customerPhone,
		CustomerDNI:     customerDNI,
		Status:          status,
		OriginalStatus:  order.Status,
		OccurredAt:      order.CreationDate,
		ImportedAt:      now,
	}

	// Order items
	dto.OrderItems = make([]canonical.ProbabilityOrderItemDTO, 0, len(order.Items))
	for _, item := range order.Items {
		productID := item.ProductID
		skuID := item.ID

		unitPrice := centavosToFloat(item.SellingPrice)
		totalPrice := unitPrice * float64(item.Quantity)
		itemTax := centavosToFloat(item.Tax)

		var imageURL *string
		if item.ImageURL != "" {
			imageURL = &item.ImageURL
		}

		var detailURL *string
		if item.DetailURL != "" {
			detailURL = &item.DetailURL
		}

		dto.OrderItems = append(dto.OrderItems, canonical.ProbabilityOrderItemDTO{
			ProductID:    &productID,
			ProductSKU:   item.RefID,
			ProductName:  item.Name,
			ProductTitle: item.SKUName,
			VariantID:    &skuID,
			Quantity:     item.Quantity,
			UnitPrice:    unitPrice,
			TotalPrice:   totalPrice,
			Currency:     currency,
			Tax:          itemTax,
			ImageURL:     imageURL,
			ProductURL:   detailURL,
		})
	}

	// Addresses (from shipping data)
	if order.ShippingData != nil && order.ShippingData.Address != nil {
		addr := order.ShippingData.Address
		street := addr.Street
		if addr.Number != "" {
			street = addr.Street + " " + addr.Number
		}

		street2 := ""
		if addr.Complement != "" {
			street2 = addr.Complement
		}
		if addr.Neighborhood != "" {
			if street2 != "" {
				street2 += ", "
			}
			street2 += addr.Neighborhood
		}

		// Split receiver name into first and last
		firstName, lastName := splitName(addr.ReceiverName)

		var lat, lng *float64
		if len(addr.GeoCoordinates) >= 2 {
			// VTEX geoCoordinates: [longitude, latitude]
			lng = &addr.GeoCoordinates[0]
			lat = &addr.GeoCoordinates[1]
		}

		var instructions *string
		if addr.Reference != "" {
			instructions = &addr.Reference
		}

		dto.Addresses = append(dto.Addresses, canonical.ProbabilityAddressDTO{
			Type:         "shipping",
			FirstName:    firstName,
			LastName:     lastName,
			Street:       street,
			Street2:      street2,
			City:         addr.City,
			State:        addr.State,
			Country:      addr.Country,
			PostalCode:   addr.PostalCode,
			Latitude:     lat,
			Longitude:    lng,
			Instructions: instructions,
		})
	}

	// Payments
	if order.PaymentData != nil {
		for _, tx := range order.PaymentData.Transactions {
			for _, p := range tx.Payments {
				paymentStatus := mapVTEXPaymentStatus(order.Status)
				gateway := p.PaymentSystemName

				var transactionID *string
				if p.TID != "" {
					transactionID = &p.TID
				}

				dto.Payments = append(dto.Payments, canonical.ProbabilityPaymentDTO{
					Amount:        centavosToFloat(p.Value),
					Currency:      currency,
					Status:        paymentStatus,
					Gateway:       &gateway,
					TransactionID: transactionID,
				})
			}
		}
	}

	// Shipments (from packages or logistics info)
	if order.PackageAttachment != nil {
		for _, pkg := range order.PackageAttachment.Packages {
			shipment := canonical.ProbabilityShipmentDTO{
				Status: mapVTEXShippingStatus(order.Status),
			}

			if pkg.TrackingNumber != "" {
				shipment.TrackingNumber = &pkg.TrackingNumber
			}
			if pkg.TrackingURL != "" {
				shipment.TrackingURL = &pkg.TrackingURL
			}
			if pkg.Courier != "" {
				shipment.Carrier = &pkg.Courier
			}
			if pkg.InvoiceURL != "" {
				invoiceURL := pkg.InvoiceURL
				dto.InvoiceURL = &invoiceURL
			}

			shippingCostVal := centavosToFloat(order.TotalFreight)
			shipment.ShippingCost = &shippingCostVal

			// Estimated delivery from logistics info
			if order.ShippingData != nil {
				for _, li := range order.ShippingData.LogisticsInfo {
					if li.ShippingEstimateDate != nil {
						shipment.EstimatedDelivery = li.ShippingEstimateDate
						break
					}
					if li.DeliveryCompany != "" && shipment.Carrier == nil {
						carrier := li.DeliveryCompany
						shipment.Carrier = &carrier
					}
				}
			}

			dto.Shipments = append(dto.Shipments, shipment)
		}
	}

	// If no packages yet but we have logistics info, create a pending shipment
	if len(dto.Shipments) == 0 && order.ShippingData != nil && len(order.ShippingData.LogisticsInfo) > 0 {
		li := order.ShippingData.LogisticsInfo[0]
		shipment := canonical.ProbabilityShipmentDTO{
			Status: "pending",
		}
		if li.DeliveryCompany != "" {
			shipment.Carrier = &li.DeliveryCompany
		}
		shippingCostVal := centavosToFloat(order.TotalFreight)
		shipment.ShippingCost = &shippingCostVal
		if li.ShippingEstimateDate != nil {
			shipment.EstimatedDelivery = li.ShippingEstimateDate
		}
		dto.Shipments = append(dto.Shipments, shipment)
	}

	// Channel metadata
	if rawJSON != nil {
		dto.ChannelMetadata = &canonical.ProbabilityChannelMetadataDTO{
			ChannelSource: "vtex",
			RawData:       datatypes.JSON(rawJSON),
			Version:       "v1",
			ReceivedAt:    now,
			IsLatest:      true,
			SyncStatus:    "synced",
		}
	}

	return dto
}

// centavosToFloat convierte centavos (int) a float64 (dividir por 100).
func centavosToFloat(centavos int) float64 {
	return float64(centavos) / 100.0
}

// extractCurrency extrae el código de moneda de la orden.
func extractCurrency(order *domain.VTEXOrder) string {
	// Intentar desde Totals
	for _, t := range order.Totals {
		if t.ID == "Items" || t.ID == "Shipping" {
			// Los totals no tienen currency directamente, pero el orden general sí
			break
		}
	}

	// Usar la currency del primer order summary si está disponible
	if order.Currency != "" {
		return order.Currency
	}

	// Default para VTEX LATAM
	return "BRL"
}

// extractTax extrae el total de impuestos de la orden.
func extractTax(order *domain.VTEXOrder) float64 {
	// Buscar en Totals el ID "Tax"
	for _, t := range order.Totals {
		if t.ID == "Tax" {
			return centavosToFloat(t.Value)
		}
	}
	// Sumar tax de items individuales
	totalTax := 0
	for _, item := range order.Items {
		totalTax += item.Tax * item.Quantity
	}
	return centavosToFloat(totalTax)
}

// splitName divide un nombre completo en firstName y lastName.
func splitName(fullName string) (string, string) {
	parts := strings.Fields(strings.TrimSpace(fullName))
	if len(parts) == 0 {
		return "", ""
	}
	if len(parts) == 1 {
		return parts[0], ""
	}
	return parts[0], strings.Join(parts[1:], " ")
}

// mapVTEXOrderStatus mapea el estado de orden VTEX al estado canónico de Probability.
// Referencia: https://help.vtex.com/en/tutorial/order-flow-and-status
func mapVTEXOrderStatus(vtexStatus string) string {
	switch vtexStatus {
	case "waiting-for-sellers-confirmation", "order-created":
		return "pending"
	case "payment-pending", "waiting-for-authorization", "approve-payment":
		return "pending"
	case "payment-approved", "authorize-fulfillment":
		return "paid"
	case "window-to-cancel":
		return "paid"
	case "ready-for-handling", "start-handling", "handling":
		return "processing"
	case "waiting-for-mkt-authorization", "waiting-ffmt-authorization":
		return "processing"
	case "invoice", "invoiced":
		return "invoiced"
	case "replaced", "cancellation-requested", "cancel":
		return "cancelled"
	case "canceled":
		return "cancelled"
	default:
		return vtexStatus
	}
}

// mapVTEXPaymentStatus mapea el estado de la orden VTEX a un estado de pago canónico.
func mapVTEXPaymentStatus(vtexOrderStatus string) string {
	switch vtexOrderStatus {
	case "payment-approved", "authorize-fulfillment", "window-to-cancel",
		"ready-for-handling", "start-handling", "handling",
		"invoice", "invoiced":
		return "paid"
	case "payment-pending", "waiting-for-authorization", "approve-payment":
		return "pending"
	case "canceled", "cancel", "cancellation-requested":
		return "cancelled"
	default:
		return "pending"
	}
}

// mapVTEXShippingStatus mapea el estado de la orden VTEX a un estado de envío canónico.
func mapVTEXShippingStatus(vtexOrderStatus string) string {
	switch vtexOrderStatus {
	case "invoiced", "invoice":
		return "shipped"
	case "ready-for-handling", "start-handling", "handling":
		return "pending"
	case "canceled", "cancel":
		return "cancelled"
	default:
		return "pending"
	}
}
