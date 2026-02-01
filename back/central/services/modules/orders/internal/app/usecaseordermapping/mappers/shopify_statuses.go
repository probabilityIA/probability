package mappers

// MapShopifyFinancialStatusToPaymentStatus mapea el financial_status de Shopify al código de PaymentStatus de Probability
func MapShopifyFinancialStatusToPaymentStatus(financialStatus string) string {
	switch financialStatus {
	case "pending":
		return "pending"
	case "authorized":
		return "authorized"
	case "paid":
		return "paid"
	case "partially_paid":
		return "partially_paid"
	case "refunded":
		return "refunded"
	case "partially_refunded":
		return "partially_refunded"
	case "voided":
		return "voided"
	case "unpaid":
		return "unpaid"
	default:
		return "pending" // Valor por defecto
	}
}

// MapShopifyFulfillmentStatusToFulfillmentStatus mapea el fulfillment_status de Shopify al código de FulfillmentStatus de Probability
func MapShopifyFulfillmentStatusToFulfillmentStatus(fulfillmentStatus *string) string {
	if fulfillmentStatus == nil || *fulfillmentStatus == "" {
		return "unfulfilled"
	}
	switch *fulfillmentStatus {
	case "unfulfilled":
		return "unfulfilled"
	case "partial":
		return "partial"
	case "fulfilled":
		return "fulfilled"
	case "shipped":
		return "shipped"
	default:
		return "unfulfilled" // Valor por defecto
	}
}
