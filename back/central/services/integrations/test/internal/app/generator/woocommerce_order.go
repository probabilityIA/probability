package generator

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/test/internal/domain"
	"gorm.io/datatypes"
)

// GenerateWooCommerceOrderJSON genera un JSON de orden de WooCommerce simulado
func (g *OrderGenerator) GenerateWooCommerceOrderJSON(r *rand.Rand) (map[string]interface{}, error) {
	// Seleccionar cliente aleatorio
	customer := domain.FakeCustomers[r.Intn(len(domain.FakeCustomers))]
	
	// Seleccionar dirección aleatoria
	address := domain.FakeAddresses[r.Intn(len(domain.FakeAddresses))]

	// Generar 1-3 items aleatorios
	numItems := r.Intn(3) + 1
	lineItems := make([]map[string]interface{}, numItems)
	subtotal := 0.0
	totalTax := 0.0

	for i := 0; i < numItems; i++ {
		product := domain.FakeProducts[r.Intn(len(domain.FakeProducts))]
		quantity := r.Intn(2) + 1
		price := product.Price
		itemSubtotal := price * float64(quantity)
		itemTax := itemSubtotal * 0.19
		totalTax += itemTax
		subtotal += itemSubtotal

		lineItems[i] = map[string]interface{}{
			"id":           315 + i,
			"name":         product.Name,
			"product_id":   99 + i,
			"variation_id": 0,
			"quantity":     quantity,
			"tax_class":    "",
			"subtotal":     fmt.Sprintf("%.2f", itemSubtotal),
			"subtotal_tax": fmt.Sprintf("%.2f", itemTax),
			"total":        fmt.Sprintf("%.2f", itemSubtotal),
			"total_tax":    fmt.Sprintf("%.2f", itemTax),
			"taxes": []map[string]interface{}{
				{
					"id":       75 + i,
					"total":    fmt.Sprintf("%.2f", itemTax),
					"subtotal": fmt.Sprintf("%.2f", itemTax),
				},
			},
			"meta_data": []interface{}{},
			"sku":        product.SKU,
			"price":      price,
		}
	}

	// Calcular shipping
	shippingTotal := float64(r.Intn(5000) + 5000) // Entre 5000 y 10000 COP
	total := subtotal + totalTax + shippingTotal

	// Generar IDs
	orderID := 727 + r.Intn(10000)
	orderKey := fmt.Sprintf("wc_order_%x", r.Int63())

	now := time.Now()
	dateCreated := now.Add(-time.Duration(r.Intn(24)) * time.Hour).Format("2006-01-02T15:04:05")
	dateModified := now.Add(-time.Duration(r.Intn(5)) * time.Hour).Format("2006-01-02T15:04:05")
	datePaid := dateCreated

	// Estado aleatorio
	statuses := []string{"processing", "pending", "on-hold", "completed", "cancelled", "refunded", "failed"}
	status := statuses[r.Intn(len(statuses))]

	// Métodos de pago
	paymentMethods := []map[string]string{
		{"id": "bacs", "title": "Transferencia Directa"},
		{"id": "cod", "title": "Pago Contra Entrega"},
		{"id": "cheque", "title": "Cheque"},
		{"id": "paypal", "title": "PayPal"},
		{"id": "stripe", "title": "Tarjeta de Crédito"},
	}
	paymentMethod := paymentMethods[r.Intn(len(paymentMethods))]

	// Construir JSON de WooCommerce
	wooOrder := map[string]interface{}{
		"id":                  orderID,
		"parent_id":           0,
		"number":              fmt.Sprintf("%d", orderID),
		"order_key":           orderKey,
		"created_via":         "checkout",
		"version":             "5.9.0",
		"status":              status,
		"currency":            "COP",
		"date_created":        dateCreated,
		"date_modified":       dateModified,
		"discount_total":      "0.00",
		"discount_tax":        "0.00",
		"shipping_total":      fmt.Sprintf("%.2f", shippingTotal),
		"shipping_tax":        "0.00",
		"cart_tax":            fmt.Sprintf("%.2f", totalTax),
		"total":               fmt.Sprintf("%.2f", total),
		"total_tax":           fmt.Sprintf("%.2f", totalTax),
		"prices_include_tax":  false,
		"customer_id":         r.Intn(100) + 1,
		"customer_ip_address": fmt.Sprintf("192.168.1.%d", r.Intn(255)),
		"customer_user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
		"customer_note":       "Por favor dejar en portería",
		"billing": map[string]interface{}{
			"first_name": customer.Name,
			"last_name":  "",
			"company":    "",
			"address_1":  address.Street,
			"address_2":  address.Street2,
			"city":       address.City,
			"state":      address.State,
			"postcode":   address.PostalCode,
			"country":    address.Country,
			"email":      customer.Email,
			"phone":      customer.Phone,
		},
		"shipping": map[string]interface{}{
			"first_name": customer.Name,
			"last_name":  "",
			"company":    "",
			"address_1":  address.Street,
			"address_2":  address.Street2,
			"city":       address.City,
			"state":      address.State,
			"postcode":   address.PostalCode,
			"country":    address.Country,
		},
		"payment_method":       paymentMethod["id"],
		"payment_method_title": paymentMethod["title"],
		"transaction_id":       "",
		"date_paid":            datePaid,
		"date_completed":       nil,
		"cart_hash":            fmt.Sprintf("%x", r.Int63()),
		"meta_data": []map[string]interface{}{
			{
				"id":    131,
				"key":   "_billing_cedula",
				"value": customer.DNI,
			},
			{
				"id":    132,
				"key":   "is_vat_exempt",
				"value": "no",
			},
		},
		"line_items": lineItems,
		"tax_lines":  []interface{}{},
		"shipping_lines": []map[string]interface{}{
			{
				"id":          317,
				"method_title": "Envío Rápido",
				"method_id":    "flat_rate",
				"instance_id":  "2",
				"total":       fmt.Sprintf("%.2f", shippingTotal),
				"total_tax":   "0.00",
				"taxes":       []interface{}{},
				"meta_data":   []interface{}{},
			},
		},
		"fee_lines":   []interface{}{},
		"coupon_lines": []interface{}{},
		"refunds":      []interface{}{},
	}

	return wooOrder, nil
}

// MapWooCommerceJSONToCanonical mapea un JSON de WooCommerce al formato canónico
func (g *OrderGenerator) MapWooCommerceJSONToCanonical(wooJSON map[string]interface{}, integrationID uint, businessID *uint) (*domain.CanonicalOrderDTO, error) {
	// Helper para obtener valores
	getString := func(key string) string {
		if v, ok := wooJSON[key].(string); ok {
			return v
		}
		return ""
	}
	getFloat := func(key string) float64 {
		if v, ok := wooJSON[key].(string); ok {
			if f, err := strconv.ParseFloat(v, 64); err == nil {
				return f
			}
		}
		if v, ok := wooJSON[key].(float64); ok {
			return v
		}
		return 0
	}
	getInt := func(key string) int {
		if v, ok := wooJSON[key].(float64); ok {
			return int(v)
		}
		if v, ok := wooJSON[key].(int); ok {
			return v
		}
		return 0
	}

	// Parsear fechas
	dateCreatedStr := getString("date_created")
	dateCreated, _ := time.Parse("2006-01-02T15:04:05", dateCreatedStr)
	if dateCreated.IsZero() {
		dateCreated = time.Now()
	}

	// Obtener billing y shipping
	billingData, _ := wooJSON["billing"].(map[string]interface{})
	shippingData, _ := wooJSON["shipping"].(map[string]interface{})

	// Mapear line items
	lineItemsData, _ := wooJSON["line_items"].([]interface{})
	orderItems := make([]domain.CanonicalOrderItemDTO, 0, len(lineItemsData))
	
	for _, itemData := range lineItemsData {
		item, _ := itemData.(map[string]interface{})
		
		price := 0.0
		if p, ok := item["price"].(float64); ok {
			price = p
		}
		
		quantity := 0
		if q, ok := item["quantity"].(float64); ok {
			quantity = int(q)
		} else if q, ok := item["quantity"].(int); ok {
			quantity = q
		}
		
		totalPrice := 0.0
		if tp, ok := item["total"].(string); ok {
			if f, err := strconv.ParseFloat(tp, 64); err == nil {
				totalPrice = f
			}
		}
		
		tax := 0.0
		if tt, ok := item["total_tax"].(string); ok {
			if f, err := strconv.ParseFloat(tt, 64); err == nil {
				tax = f
			}
		}
		
		var sku string
		if s, ok := item["sku"].(string); ok {
			sku = s
		}
		
		var productName string
		if n, ok := item["name"].(string); ok {
			productName = n
		}

		orderItems = append(orderItems, domain.CanonicalOrderItemDTO{
			ProductSKU:   sku,
			ProductName:  productName,
			ProductTitle: productName,
			Quantity:     quantity,
			UnitPrice:    price,
			TotalPrice:   totalPrice,
			Currency:     "COP",
			Tax:          tax,
			TaxRate:      floatPtr(0.19),
		})
	}

	// Calcular totales
	subtotal := getFloat("total") - getFloat("shipping_total") - getFloat("cart_tax")
	tax := getFloat("cart_tax")
	shippingCost := getFloat("shipping_total")
	totalAmount := getFloat("total")

	// Obtener información del cliente
	customerName := ""
	customerEmail := ""
	customerPhone := ""
	if billingData != nil {
		if firstName, ok := billingData["first_name"].(string); ok {
			customerName = firstName
		}
		if lastName, ok := billingData["last_name"].(string); ok && lastName != "" {
			customerName += " " + lastName
		}
		if email, ok := billingData["email"].(string); ok {
			customerEmail = email
		}
		if phone, ok := billingData["phone"].(string); ok {
			customerPhone = phone
		}
	}

	// Obtener direcciones
	addresses := []domain.CanonicalAddressDTO{}
	if shippingData != nil {
		addresses = append(addresses, domain.CanonicalAddressDTO{
			Type:       "shipping",
			FirstName:  addressString(shippingData, "first_name"),
			LastName:   addressString(shippingData, "last_name"),
			Street:     addressString(shippingData, "address_1"),
			Street2:    addressString(shippingData, "address_2"),
			City:       addressString(shippingData, "city"),
			State:      addressString(shippingData, "state"),
			Country:    addressString(shippingData, "country"),
			PostalCode: addressString(shippingData, "postcode"),
		})
	}
	if billingData != nil {
		addresses = append(addresses, domain.CanonicalAddressDTO{
			Type:       "billing",
			FirstName:  addressString(billingData, "first_name"),
			LastName:   addressString(billingData, "last_name"),
			Street:     addressString(billingData, "address_1"),
			Street2:    addressString(billingData, "address_2"),
			City:       addressString(billingData, "city"),
			State:      addressString(billingData, "state"),
			Country:    addressString(billingData, "country"),
			PostalCode: addressString(billingData, "postcode"),
			Phone:      addressString(billingData, "phone"),
		})
	}

	// Mapear pagos
	payments := []domain.CanonicalPaymentDTO{}
	paymentMethod := getString("payment_method")
	datePaidStr := getString("date_paid")
	status := getString("status")
	
	paymentStatus := "pending"
	if status == "processing" || status == "completed" {
		paymentStatus = "completed"
	} else if status == "failed" {
		paymentStatus = "failed"
	} else if status == "refunded" {
		paymentStatus = "refunded"
	}
	
	paidAt := &dateCreated
	if datePaidStr != "" {
		if dt, err := time.Parse("2006-01-02T15:04:05", datePaidStr); err == nil {
			paidAt = &dt
		}
	}
	if paymentStatus == "pending" {
		paidAt = nil
	}
	
	paymentMethodTitle := getString("payment_method_title")
	if paymentMethodTitle == "" {
		paymentMethodTitle = paymentMethod
	}
	
	payments = append(payments, domain.CanonicalPaymentDTO{
		PaymentMethodID: 1, // Default
		Amount:          totalAmount,
		Currency:        "COP",
		Status:          paymentStatus,
		PaidAt:          paidAt,
		Gateway:         stringPtr(paymentMethod),
		PaymentReference: stringPtr(paymentMethodTitle),
	})

	// Serializar JSON de WooCommerce para metadata
	rawDataJSON, _ := json.Marshal(wooJSON)

	// Crear orden canónica
	orderNumber := getString("number")
	if orderNumber == "" {
		orderNumber = fmt.Sprintf("%d", getInt("id"))
	}
	externalID := orderNumber

	return &domain.CanonicalOrderDTO{
		BusinessID:      businessID,
		IntegrationID:   integrationID,
		IntegrationType: "woocommerce",
		Platform:        "woocommerce",
		ExternalID:      externalID,
		OrderNumber:     orderNumber,
		Subtotal:        subtotal,
		Tax:             tax,
		Discount:        getFloat("discount_total"),
		ShippingCost:    shippingCost,
		TotalAmount:     totalAmount,
		Currency:        "COP",
		CustomerName:    customerName,
		CustomerEmail:   customerEmail,
		CustomerPhone:   customerPhone,
		OrderTypeName:   "delivery",
		Status:          mapWooCommerceStatus(status),
		OriginalStatus:  status,
		OrderItems:      orderItems,
		Addresses:       addresses,
		Payments:        payments,
		OccurredAt:      dateCreated,
		ImportedAt:      time.Now(),
		ChannelMetadata: &domain.CanonicalChannelMetadataDTO{
			ChannelSource: "woocommerce",
			RawData:       datatypes.JSON(rawDataJSON),
			Version:       "1.0",
			ReceivedAt:    time.Now(),
			IsLatest:      true,
			SyncStatus:    "pending",
		},
	}, nil
}

func mapWooCommerceStatus(status string) string {
	switch status {
	case "processing":
		return "processing"
	case "pending":
		return "pending"
	case "on-hold":
		return "on_hold"
	case "completed":
		return "completed"
	case "cancelled":
		return "cancelled"
	case "refunded":
		return "refunded"
	case "failed":
		return "failed"
	default:
		return "pending"
	}
}

