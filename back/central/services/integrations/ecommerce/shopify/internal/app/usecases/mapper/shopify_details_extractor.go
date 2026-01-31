package mapper

import (
	"encoding/json"

	"gorm.io/datatypes"
)

// ExtractFinancialDetails extrae los detalles financieros del JSON original de Shopify
func ExtractFinancialDetails(rawPayload []byte) (datatypes.JSON, error) {
	var order map[string]interface{}
	if err := json.Unmarshal(rawPayload, &order); err != nil {
		return nil, err
	}

	financialDetails := map[string]interface{}{}

	// Extraer información financiera relevante
	if subtotalPrice, ok := order["subtotal_price"].(string); ok {
		financialDetails["subtotal_price"] = subtotalPrice
	}
	if subtotalPriceSet, ok := order["subtotal_price_set"].(map[string]interface{}); ok {
		financialDetails["subtotal_price_set"] = subtotalPriceSet
	}
	if totalDiscounts, ok := order["total_discounts"].(string); ok {
		financialDetails["total_discounts"] = totalDiscounts
	}
	if totalDiscountsSet, ok := order["total_discounts_set"].(map[string]interface{}); ok {
		financialDetails["total_discounts_set"] = totalDiscountsSet
	}
	if totalTax, ok := order["total_tax"].(string); ok {
		financialDetails["total_tax"] = totalTax
	}
	if totalTaxSet, ok := order["total_tax_set"].(map[string]interface{}); ok {
		financialDetails["total_tax_set"] = totalTaxSet
	}
	if currentSubtotalPrice, ok := order["current_subtotal_price"].(string); ok {
		financialDetails["current_subtotal_price"] = currentSubtotalPrice
	}
	if currentSubtotalPriceSet, ok := order["current_subtotal_price_set"].(map[string]interface{}); ok {
		financialDetails["current_subtotal_price_set"] = currentSubtotalPriceSet
	}
	if currentTotalDiscounts, ok := order["current_total_discounts"].(string); ok {
		financialDetails["current_total_discounts"] = currentTotalDiscounts
	}
	if currentTotalDiscountsSet, ok := order["current_total_discounts_set"].(map[string]interface{}); ok {
		financialDetails["current_total_discounts_set"] = currentTotalDiscountsSet
	}
	if currentTotalTax, ok := order["current_total_tax"].(string); ok {
		financialDetails["current_total_tax"] = currentTotalTax
	}
	if currentTotalTaxSet, ok := order["current_total_tax_set"].(map[string]interface{}); ok {
		financialDetails["current_total_tax_set"] = currentTotalTaxSet
	}
	if taxesIncluded, ok := order["taxes_included"].(bool); ok {
		financialDetails["taxes_included"] = taxesIncluded
	}
	if discountCodes, ok := order["discount_codes"].([]interface{}); ok {
		financialDetails["discount_codes"] = discountCodes
	}
	if discountApplications, ok := order["discount_applications"].([]interface{}); ok {
		financialDetails["discount_applications"] = discountApplications
	}
	if taxLines, ok := order["tax_lines"].([]interface{}); ok {
		financialDetails["tax_lines"] = taxLines
	}
	if totalLineItemsPrice, ok := order["total_line_items_price"].(string); ok {
		financialDetails["total_line_items_price"] = totalLineItemsPrice
	}
	if totalLineItemsPriceSet, ok := order["total_line_items_price_set"].(map[string]interface{}); ok {
		financialDetails["total_line_items_price_set"] = totalLineItemsPriceSet
	}
	if totalOutstanding, ok := order["total_outstanding"].(string); ok {
		financialDetails["total_outstanding"] = totalOutstanding
	}
	if totalTipReceived, ok := order["total_tip_received"].(string); ok {
		financialDetails["total_tip_received"] = totalTipReceived
	}
	if refunds, ok := order["refunds"].([]interface{}); ok {
		financialDetails["refunds"] = refunds
	}

	financialJSON, err := json.Marshal(financialDetails)
	if err != nil {
		return nil, err
	}

	return datatypes.JSON(financialJSON), nil
}

// ExtractShippingDetails extrae los detalles de envío del JSON original de Shopify
func ExtractShippingDetails(rawPayload []byte) (datatypes.JSON, error) {
	var order map[string]interface{}
	if err := json.Unmarshal(rawPayload, &order); err != nil {
		return nil, err
	}

	shippingDetails := map[string]interface{}{}

	// Extraer información de envío relevante
	if shippingAddress, ok := order["shipping_address"].(map[string]interface{}); ok {
		shippingDetails["shipping_address"] = shippingAddress
	}
	if shippingLines, ok := order["shipping_lines"].([]interface{}); ok {
		shippingDetails["shipping_lines"] = shippingLines
	}
	if totalShippingPriceSet, ok := order["total_shipping_price_set"].(map[string]interface{}); ok {
		shippingDetails["total_shipping_price_set"] = totalShippingPriceSet
	}
	if totalWeight, ok := order["total_weight"].(float64); ok {
		shippingDetails["total_weight"] = totalWeight
	}
	if fulfillmentStatus, ok := order["fulfillment_status"].(string); ok {
		shippingDetails["fulfillment_status"] = fulfillmentStatus
	}

	shippingJSON, err := json.Marshal(shippingDetails)
	if err != nil {
		return nil, err
	}

	return datatypes.JSON(shippingJSON), nil
}

// ExtractPaymentDetails extrae los detalles de pago del JSON original de Shopify
func ExtractPaymentDetails(rawPayload []byte) (datatypes.JSON, error) {
	var order map[string]interface{}
	if err := json.Unmarshal(rawPayload, &order); err != nil {
		return nil, err
	}

	paymentDetails := map[string]interface{}{}

	// Extraer información de pago relevante
	if financialStatus, ok := order["financial_status"].(string); ok {
		paymentDetails["financial_status"] = financialStatus
	}
	if gateway, ok := order["gateway"].(string); ok {
		paymentDetails["gateway"] = gateway
	}
	if paymentGatewayNames, ok := order["payment_gateway_names"].([]interface{}); ok {
		paymentDetails["payment_gateway_names"] = paymentGatewayNames
	}
	if processingMethod, ok := order["processing_method"].(string); ok {
		paymentDetails["processing_method"] = processingMethod
	}
	if transactions, ok := order["transactions"].([]interface{}); ok {
		paymentDetails["transactions"] = transactions
	}
	if paymentTerms, ok := order["payment_terms"].(map[string]interface{}); ok {
		paymentDetails["payment_terms"] = paymentTerms
	}
	if billingAddress, ok := order["billing_address"].(map[string]interface{}); ok {
		paymentDetails["billing_address"] = billingAddress
	}

	paymentJSON, err := json.Marshal(paymentDetails)
	if err != nil {
		return nil, err
	}

	return datatypes.JSON(paymentJSON), nil
}

// ExtractFulfillmentDetails extrae los detalles de fulfillment del JSON original de Shopify
func ExtractFulfillmentDetails(rawPayload []byte) (datatypes.JSON, error) {
	var order map[string]interface{}
	if err := json.Unmarshal(rawPayload, &order); err != nil {
		return nil, err
	}

	fulfillmentDetails := map[string]interface{}{}

	// Extraer información de fulfillment relevante
	if fulfillments, ok := order["fulfillments"].([]interface{}); ok {
		fulfillmentDetails["fulfillments"] = fulfillments
	}
	if fulfillmentStatus, ok := order["fulfillment_status"].(string); ok {
		fulfillmentDetails["fulfillment_status"] = fulfillmentStatus
	}
	if locationID, ok := order["location_id"].(float64); ok {
		fulfillmentDetails["location_id"] = locationID
	}
	if lineItems, ok := order["line_items"].([]interface{}); ok {
		// Extraer información de fulfillment de cada line item
		lineItemsWithFulfillment := []map[string]interface{}{}
		for _, item := range lineItems {
			if itemMap, ok := item.(map[string]interface{}); ok {
				itemFulfillment := map[string]interface{}{}
				if fulfillableQuantity, ok := itemMap["fulfillable_quantity"].(float64); ok {
					itemFulfillment["fulfillable_quantity"] = fulfillableQuantity
				}
				if fulfillmentService, ok := itemMap["fulfillment_service"].(string); ok {
					itemFulfillment["fulfillment_service"] = fulfillmentService
				}
				if fulfillmentStatus, ok := itemMap["fulfillment_status"].(string); ok {
					itemFulfillment["fulfillment_status"] = fulfillmentStatus
				}
				if requiresShipping, ok := itemMap["requires_shipping"].(bool); ok {
					itemFulfillment["requires_shipping"] = requiresShipping
				}
				lineItemsWithFulfillment = append(lineItemsWithFulfillment, itemFulfillment)
			}
		}
		if len(lineItemsWithFulfillment) > 0 {
			fulfillmentDetails["line_items_fulfillment"] = lineItemsWithFulfillment
		}
	}

	fulfillmentJSON, err := json.Marshal(fulfillmentDetails)
	if err != nil {
		return nil, err
	}

	return datatypes.JSON(fulfillmentJSON), nil
}
