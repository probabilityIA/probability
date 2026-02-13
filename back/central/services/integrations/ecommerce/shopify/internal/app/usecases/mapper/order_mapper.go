package mapper

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/domain"
	"gorm.io/datatypes"
)

func MapShopifyOrderToProbability(s *domain.ShopifyOrder) *domain.ProbabilityOrderDTO {
	orderItems := make([]domain.ProbabilityOrderItemDTO, len(s.Items))
	for i, item := range s.Items {
		// Convertir ProductID y VariantID a string si están disponibles
		var productIDStr *string
		if item.ProductID != nil {
			idStr := strconv.FormatInt(*item.ProductID, 10)
			productIDStr = &idStr
		}

		var variantIDStr *string
		if item.VariantID != nil {
			idStr := strconv.FormatInt(*item.VariantID, 10)
			variantIDStr = &idStr
		}

		// Calcular precio total (precio unitario * cantidad - descuento)
		totalPrice := (item.UnitPrice * float64(item.Quantity)) - item.Discount
		if totalPrice < 0 {
			totalPrice = 0
		}

		// Calcular precio total en moneda local (precio unitario * cantidad - descuento)
		totalPricePresentment := (item.UnitPricePresentment * float64(item.Quantity)) - item.DiscountPresentment
		if totalPricePresentment < 0 {
			totalPricePresentment = 0
		}

		// Si el SKU está vacío, generar uno usando ProductID y VariantID
		sku := item.SKU
		if sku == "" {
			if item.VariantID != nil {
				sku = fmt.Sprintf("VAR-%d", *item.VariantID)
			} else if item.ProductID != nil {
				sku = fmt.Sprintf("PROD-%d", *item.ProductID)
			} else {
				sku = fmt.Sprintf("ITEM-%d", i) // Fallback usando índice
			}
		}

		orderItems[i] = domain.ProbabilityOrderItemDTO{
			ProductID:    productIDStr,
			ProductSKU:   sku,
			ProductName:  item.Name,
			ProductTitle: item.Title,
			VariantID:    variantIDStr,
			Quantity:     item.Quantity,
			UnitPrice:    item.UnitPrice,
			TotalPrice:   totalPrice,
			Currency:     s.Currency,
			Discount:     item.Discount,
			Tax:          item.Tax,
			Weight:       item.Weight,
			// Precios en moneda local
			UnitPricePresentment:  item.UnitPricePresentment,
			TotalPricePresentment: totalPricePresentment,
			DiscountPresentment:   item.DiscountPresentment,
			TaxPresentment:        item.TaxPresentment,
		}
	}

	addresses := []domain.ProbabilityAddressDTO{}
	if s.ShippingAddress.Street != "" || s.ShippingAddress.City != "" {
		address := domain.ProbabilityAddressDTO{
			Type:       "shipping",
			Street:     s.ShippingAddress.Street,
			Street2:    s.ShippingAddress.Address2,
			City:       s.ShippingAddress.City,
			State:      s.ShippingAddress.State,
			Country:    s.ShippingAddress.Country,
			PostalCode: s.ShippingAddress.PostalCode,
		}

		// Fallback: Si Address2 está vacío, intentar usar DefaultAddress del cliente
		if address.Street2 == "" && s.Customer.DefaultAddress != nil && s.Customer.DefaultAddress.Address2 != "" {
			fmt.Printf("[Mapper] Usando DefaultAddress.Address2 como fallback para orden %s: %s\n", s.OrderNumber, s.Customer.DefaultAddress.Address2)
			address.Street2 = s.Customer.DefaultAddress.Address2
		}
		if s.ShippingAddress.Coordinates != nil {
			address.Latitude = &s.ShippingAddress.Coordinates.Lat
			address.Longitude = &s.ShippingAddress.Coordinates.Lng
		}
		addresses = append(addresses, address)
	}

	itemsJSON, _ := json.Marshal(orderItems)

	var metadataJSON []byte
	if s.Metadata != nil {
		metadataJSON, _ = json.Marshal(s.Metadata)
	}

	// Extraer y mapear shipments desde fulfillments del raw_data
	shipments := extractShipmentsFromRawData(s.RawData)

	// Calcular montos finales usando presentment money si está disponible (prioridad para mostrar orden en moneda local)
	totalAmount := s.TotalAmount
	currency := s.Currency

	// Si tenemos valores en moneda local (COP, etc), usarlos como principales
	if s.TotalAmountPresentment > 0 && s.CurrencyPresentment != "" {
		totalAmount = s.TotalAmountPresentment
		currency = s.CurrencyPresentment
	}

	subtotal := s.Subtotal
	if s.SubtotalPresentment > 0 {
		subtotal = s.SubtotalPresentment
	}

	tax := s.Tax
	if s.TaxPresentment > 0 {
		tax = s.TaxPresentment
	}

	discount := s.Discount
	if s.DiscountPresentment > 0 {
		discount = s.DiscountPresentment
	}

	shippingCost := s.ShippingCost
	if s.ShippingCostPresentment > 0 {
		shippingCost = s.ShippingCostPresentment
	}

	// Determine COD Total
	var codTotal *float64
	isCOD := false

	// Check metadata for gateway names
	if val, ok := s.Metadata["payment_gateway_names"]; ok {
		if gateways, ok := val.([]string); ok {
			for _, g := range gateways {
				lower := strings.ToLower(g)
				if strings.Contains(lower, "cod") || strings.Contains(lower, "cash") || strings.Contains(lower, "contra") {
					isCOD = true
					break
				}
			}
		}
	}

	// Check tags
	if !isCOD {
		if val, ok := s.Metadata["tags"]; ok {
			if tags, ok := val.(string); ok && tags != "" {
				lower := strings.ToLower(tags)
				if strings.Contains(lower, "cod") || strings.Contains(lower, "contra") {
					isCOD = true
				}
			}
		}
	}

	if isCOD {
		amount := totalAmount
		codTotal = &amount
	}

	probabilityOrder := &domain.ProbabilityOrderDTO{
		BusinessID:      s.BusinessID,
		IntegrationID:   s.IntegrationID,
		IntegrationType: s.IntegrationType,
		Platform:        s.Platform,
		ExternalID:      s.ExternalID,
		OrderNumber:     s.OrderNumber,
		Subtotal:        subtotal,
		Tax:             tax,
		Discount:        discount,
		ShippingCost:    shippingCost,
		TotalAmount:     totalAmount,
		Currency:        currency,
		CodTotal:        codTotal,
		CustomerName:    s.Customer.Name,
		CustomerEmail:   s.Customer.Email,

		CustomerPhone:      s.Customer.Phone,
		CustomerOrderCount: &s.Customer.OrdersCount,
		CustomerTotalSpent: &s.Customer.TotalSpent,
		Status:             s.Status,
		OriginalStatus:     s.OriginalStatus,
		OccurredAt:         s.OccurredAt,
		ImportedAt:         s.ImportedAt,
		Items:              datatypes.JSON(itemsJSON),
		Metadata:           datatypes.JSON(metadataJSON),
		OrderItems:         orderItems,
		Addresses:          addresses,
		Shipments:          shipments,
		OrderStatusURL:     s.OrderStatusURL,
		// Precios en moneda local
		SubtotalPresentment:     s.SubtotalPresentment,
		TaxPresentment:          s.TaxPresentment,
		DiscountPresentment:     s.DiscountPresentment,
		ShippingCostPresentment: s.ShippingCostPresentment,
		TotalAmountPresentment:  s.TotalAmountPresentment,
		CurrencyPresentment:     s.CurrencyPresentment,
		// Facturación - Por defecto todas las órdenes de Shopify son facturables
		// Esto puede ser modificado por filtros y reglas de facturación
		Invoiceable:             true,
	}

	if len(s.RawData) > 0 {
		probabilityOrder.ChannelMetadata = &domain.ProbabilityChannelMetadataDTO{
			ChannelSource: "shopify",
			RawData:       datatypes.JSON(s.RawData),
			Version:       "1.0",
			ReceivedAt:    time.Now(),
			IsLatest:      true,
			SyncStatus:    "synced",
		}
	}

	return probabilityOrder
}

// extractShipmentsFromRawData extrae los shipments desde los fulfillments del JSON raw_data
func extractShipmentsFromRawData(rawData []byte) []domain.ProbabilityShipmentDTO {
	if len(rawData) == 0 {
		return nil
	}

	var orderMap map[string]interface{}
	if err := json.Unmarshal(rawData, &orderMap); err != nil {
		return nil
	}

	fulfillments, ok := orderMap["fulfillments"].([]interface{})
	if !ok || len(fulfillments) == 0 {
		return nil
	}

	shipments := make([]domain.ProbabilityShipmentDTO, 0, len(fulfillments))

	for _, fulfillmentRaw := range fulfillments {
		fulfillment, ok := fulfillmentRaw.(map[string]interface{})
		if !ok {
			continue
		}

		// Mapear campos básicos del fulfillment
		trackingNumber := getStringPtr(fulfillment, "tracking_number")
		trackingURL := getStringPtr(fulfillment, "tracking_url")
		trackingCompany := getStringPtr(fulfillment, "tracking_company")
		shipmentStatus := getStringPtr(fulfillment, "shipment_status")
		service := getStringPtr(fulfillment, "service")

		// Si no hay tracking_number pero hay tracking_numbers array, usar el primero
		if trackingNumber == nil {
			if trackingNumbers, ok := fulfillment["tracking_numbers"].([]interface{}); ok && len(trackingNumbers) > 0 {
				if firstNum, ok := trackingNumbers[0].(string); ok {
					trackingNumber = &firstNum
				}
			}
		}

		// Si no hay tracking_url pero hay tracking_urls array, usar el primero
		if trackingURL == nil {
			if trackingURLs, ok := fulfillment["tracking_urls"].([]interface{}); ok && len(trackingURLs) > 0 {
				if firstURL, ok := trackingURLs[0].(string); ok {
					trackingURL = &firstURL
				}
			}
		}

		// Determinar status del shipment
		status := "pending"
		if shipmentStatus != nil {
			switch *shipmentStatus {
			case "confirmed", "success":
				status = "in_transit"
			case "delivered":
				status = "delivered"
			case "failure", "cancelled":
				status = "failed"
			default:
				status = "pending"
			}
		} else if fulfillmentStatus, ok := fulfillment["status"].(string); ok {
			switch fulfillmentStatus {
			case "success", "confirmed":
				status = "in_transit"
			case "failure", "cancelled":
				status = "failed"
			default:
				status = "pending"
			}
		}

		// Parsear fechas
		var shippedAt *time.Time
		if createdAtStr, ok := fulfillment["created_at"].(string); ok {
			if parsed, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
				shippedAt = &parsed
			}
		}

		var deliveredAt *time.Time
		if updatedAtStr, ok := fulfillment["updated_at"].(string); ok && status == "delivered" {
			if parsed, err := time.Parse(time.RFC3339, updatedAtStr); err == nil {
				deliveredAt = &parsed
			}
		}

		// Serializar metadata del fulfillment
		metadataJSON, _ := json.Marshal(fulfillment)

		shipment := domain.ProbabilityShipmentDTO{
			TrackingNumber: trackingNumber,
			TrackingURL:    trackingURL,
			Carrier:        trackingCompany,
			Status:         status,
			ShippedAt:      shippedAt,
			DeliveredAt:    deliveredAt,
			Metadata:       metadataJSON,
		}

		// Si hay service, usarlo como carrier si no hay tracking_company
		if shipment.Carrier == nil && service != nil {
			shipment.Carrier = service
		}

		shipments = append(shipments, shipment)
	}

	return shipments
}

// getStringPtr obtiene un puntero a string desde un map, retorna nil si no existe
func getStringPtr(m map[string]interface{}, key string) *string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok && str != "" {
			return &str
		}
	}
	return nil
}

// EnrichOrderWithDetails extrae y agrega PaymentDetails y FulfillmentDetails al ProbabilityOrderDTO
// desde el rawPayload (JSON original de Shopify)
func EnrichOrderWithDetails(probabilityOrder *domain.ProbabilityOrderDTO, rawPayload []byte) {
	if len(rawPayload) == 0 {
		return
	}

	// Extraer FinancialDetails
	if financialDetails, err := ExtractFinancialDetails(rawPayload); err == nil {
		probabilityOrder.FinancialDetails = financialDetails
	}

	// Extraer ShippingDetails
	if shippingDetails, err := ExtractShippingDetails(rawPayload); err == nil {
		probabilityOrder.ShippingDetails = shippingDetails
	}

	// Extraer PaymentDetails (incluye financial_status)
	if paymentDetails, err := ExtractPaymentDetails(rawPayload); err == nil {
		probabilityOrder.PaymentDetails = paymentDetails
	}

	// Extraer FulfillmentDetails (incluye fulfillment_status)
	if fulfillmentDetails, err := ExtractFulfillmentDetails(rawPayload); err == nil {
		probabilityOrder.FulfillmentDetails = fulfillmentDetails
	}
}
