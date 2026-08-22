package mapper

import "strings"

func MapOrderStatus(status, paymentStatus, shippingStatus string) string {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "cancelled", "canceled":
		return "cancelled"
	case "closed":
		return "completed"
	}

	switch strings.ToLower(strings.TrimSpace(paymentStatus)) {
	case "paid":
		if strings.EqualFold(shippingStatus, "delivered") {
			return "completed"
		}
		return "paid"
	case "refunded", "partially_refunded":
		return "refunded"
	case "voided":
		return "cancelled"
	case "abandoned":
		return "abandoned"
	case "authorized", "partially_paid":
		return "pending"
	}

	return "pending"
}

func mapPaymentStatus(paymentStatus string) string {
	switch strings.ToLower(strings.TrimSpace(paymentStatus)) {
	case "paid":
		return "completed"
	case "refunded", "partially_refunded":
		return "refunded"
	case "voided", "abandoned":
		return "cancelled"
	default:
		return "pending"
	}
}

func MapShipmentStatus(shippingStatus string) string {
	switch strings.ToLower(strings.TrimSpace(shippingStatus)) {
	case "delivered":
		return "delivered"
	case "shipped", "fulfilled", "partially_fulfilled":
		return "in_transit"
	case "unpacked", "unshipped", "partially_packed":
		return "pending"
	default:
		return "pending"
	}
}

func MapPaymentMethod(gateway, gatewayName string) uint {
	m := strings.ToLower(gateway + " " + gatewayName)
	switch {
	case strings.Contains(m, "contra") || strings.Contains(m, "cash on delivery") || strings.Contains(m, "cod"):
		return paymentMethodCOD
	case strings.Contains(m, "transfer") || strings.Contains(m, "transferencia") || strings.Contains(m, "deposit"):
		return paymentMethodBankTransfer
	case strings.Contains(m, "efectivo") || strings.Contains(m, "cash"):
		return paymentMethodCash
	case strings.Contains(m, "mercado"):
		return paymentMethodMercadoPago
	default:
		return paymentMethodCreditCard
	}
}
