package generator

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/test/internal/domain"
	"gorm.io/datatypes"
)

// GenerateMeliOrderJSON genera un JSON de orden de Mercado Libre simulado
func (g *OrderGenerator) GenerateMeliOrderJSON(r *rand.Rand) (map[string]interface{}, error) {
	// Seleccionar cliente aleatorio
	customer := domain.FakeCustomers[r.Intn(len(domain.FakeCustomers))]

	// Seleccionar dirección aleatoria
	address := domain.FakeAddresses[r.Intn(len(domain.FakeAddresses))]

	// Generar 1-3 items aleatorios
	numItems := r.Intn(3) + 1
	orderItems := make([]map[string]interface{}, numItems)
	subtotal := 0.0

	for i := 0; i < numItems; i++ {
		product := domain.FakeProducts[r.Intn(len(domain.FakeProducts))]
		quantity := r.Intn(2) + 1
		unitPrice := product.Price
		fullUnitPrice := unitPrice * 1.1 // Precio con margen
		saleFee := unitPrice * 0.1       // Comisión de Mercado Libre (10%)
		itemTotal := unitPrice * float64(quantity)
		subtotal += itemTotal

		orderItems[i] = map[string]interface{}{
			"item": map[string]interface{}{
				"id":                  fmt.Sprintf("MCO%d", 123456789+r.Intn(1000000)),
				"title":               product.Name,
				"category_id":         "MCO1714",
				"variation_id":        87654321098 + int64(r.Intn(1000)),
				"seller_custom_field": nil,
				"variation_attributes": []map[string]interface{}{
					{
						"id":         "COLOR",
						"name":       "Color",
						"value_id":   "52049",
						"value_name": []string{"Negro", "Blanco", "Plateado", "Rojo"}[r.Intn(4)],
					},
				},
				"warranty":   "Garantía de fábrica: 3 meses",
				"condition":  "new",
				"seller_sku": product.SKU,
			},
			"quantity":           quantity,
			"unit_price":         unitPrice,
			"full_unit_price":    fullUnitPrice,
			"currency_id":        "COP",
			"manufacturing_days": nil,
			"sale_fee":           saleFee,
		}
	}

	// Calcular shipping
	shippingCost := float64(r.Intn(5000) + 3000) // Entre 3000 y 8000 COP
	totalAmountWithShipping := subtotal + shippingCost
	paidAmount := totalAmountWithShipping

	// Generar IDs
	orderID := 2000003508971234 + int64(r.Intn(1000000))
	buyerID := 12345678 + int64(r.Intn(1000))
	sellerID := 98765432

	now := time.Now()
	dateCreated := now.Add(-time.Duration(r.Intn(24)) * time.Hour).Format(time.RFC3339)
	dateClosed := now.Add(-time.Duration(r.Intn(5)) * time.Hour).Format(time.RFC3339)
	dateApproved := dateClosed

	// Estado aleatorio
	statuses := []string{"paid", "pending", "confirmed", "cancelled"}
	status := statuses[r.Intn(len(statuses))]

	// Construir JSON de Mercado Libre
	meliOrder := map[string]interface{}{
		"id":                         orderID,
		"status":                     status,
		"status_detail":              nil,
		"date_created":               dateCreated,
		"date_closed":                dateClosed,
		"order_items":                orderItems,
		"total_amount":               subtotal,
		"total_amount_with_shipping": totalAmountWithShipping,
		"paid_amount":                paidAmount,
		"currency_id":                "COP",
		"shipping": map[string]interface{}{
			"id":            4000000099887766 + int64(r.Intn(1000)),
			"shipment_type": "shipping",
			"status":        "ready_to_ship",
			"date_created":  dateCreated,
			"receiver_address": map[string]interface{}{
				"id":            123456789 + int64(r.Intn(1000)),
				"address_line":  address.Street,
				"street_name":   address.Street,
				"street_number": address.Street2,
				"comment":       "Edificio",
				"zip_code":      address.PostalCode,
				"city": map[string]interface{}{
					"id":   "TUNPX0JPRzgyNjk",
					"name": address.City,
				},
				"state": map[string]interface{}{
					"id":   "CO-DC",
					"name": address.State,
				},
				"country": map[string]interface{}{
					"id":   "CO",
					"name": address.Country,
				},
				"latitude":  address.Latitude,
				"longitude": address.Longitude,
			},
		},
		"buyer": map[string]interface{}{
			"id":         buyerID,
			"nickname":   fmt.Sprintf("COMPRADOR_%d", buyerID),
			"first_name": customer.Name,
			"last_name":  "",
			"email":      customer.Email,
			"phone": map[string]interface{}{
				"area_code": "57",
				"number":    customer.Phone,
				"extension": "",
			},
			"billing_info": map[string]interface{}{
				"doc_type":   "CC",
				"doc_number": customer.DNI,
			},
		},
		"seller": map[string]interface{}{
			"id":       int64(sellerID),
			"nickname": "MI_TIENDA_OFICIAL",
			"email":    "ventas@mitienda.com",
		},
		"payments": []map[string]interface{}{
			{
				"id":       12345678900 + int64(r.Intn(1000)),
				"order_id": orderID,
				"payer_id": buyerID,
				"collector": map[string]interface{}{
					"id": int64(sellerID),
				},
				"card_id":           123456 + r.Intn(1000),
				"site_id":           "MCO",
				"reason":            orderItems[0]["item"].(map[string]interface{})["title"].(string),
				"payment_method_id": "master",
				"currency_id":       "COP",
				"installments":      1,
				"issuer_id":         "123",
				"atm_transfer_reference": map[string]interface{}{
					"company_id":     nil,
					"transaction_id": nil,
				},
				"coupon_id":          nil,
				"activation_uri":     nil,
				"operation_type":     "regular_payment",
				"payment_type":       "credit_card",
				"available_actions":  []string{"refund"},
				"status":             "approved",
				"status_code":        nil,
				"status_detail":      "accredited",
				"transaction_amount": totalAmountWithShipping,
				"taxes_amount":       0,
				"shipping_cost":      shippingCost,
				"total_paid_amount":  paidAmount,
				"installment_amount": paidAmount,
				"date_approved":      dateApproved,
				"date_created":       dateCreated,
			},
		},
		"pack_id": nil,
		"coupon": map[string]interface{}{
			"amount": 0,
			"id":     nil,
		},
		"tags": []string{"paid", "not_delivered"},
	}

	return meliOrder, nil
}

// MapMeliJSONToCanonical mapea un JSON de Mercado Libre al formato canónico
func (g *OrderGenerator) MapMeliJSONToCanonical(meliJSON map[string]interface{}, integrationID uint, businessID *uint) (*domain.CanonicalOrderDTO, error) {
	// Helper para obtener valores
	getString := func(key string) string {
		if v, ok := meliJSON[key].(string); ok {
			return v
		}
		return ""
	}
	getFloat := func(key string) float64 {
		if v, ok := meliJSON[key].(float64); ok {
			return v
		}
		return 0
	}
	getInt64 := func(key string) int64 {
		if v, ok := meliJSON[key].(float64); ok {
			return int64(v)
		}
		if v, ok := meliJSON[key].(int64); ok {
			return v
		}
		return 0
	}

	// Parsear fechas
	dateCreatedStr := getString("date_created")
	dateCreated, _ := time.Parse(time.RFC3339, dateCreatedStr)
	if dateCreated.IsZero() {
		dateCreated = time.Now()
	}

	// Obtener shipping y buyer
	shippingData, _ := meliJSON["shipping"].(map[string]interface{})
	buyerData, _ := meliJSON["buyer"].(map[string]interface{})

	// Mapear order items
	orderItemsData, _ := meliJSON["order_items"].([]interface{})
	orderItems := make([]domain.CanonicalOrderItemDTO, 0, len(orderItemsData))

	for _, itemData := range orderItemsData {
		item, _ := itemData.(map[string]interface{})
		itemInfo, _ := item["item"].(map[string]interface{})

		unitPrice := 0.0
		if up, ok := item["unit_price"].(float64); ok {
			unitPrice = up
		}

		quantity := 0
		if q, ok := item["quantity"].(float64); ok {
			quantity = int(q)
		} else if q, ok := item["quantity"].(int); ok {
			quantity = q
		}

		totalPrice := unitPrice * float64(quantity)
		tax := totalPrice * 0.19

		var sku string
		if s, ok := itemInfo["seller_sku"].(string); ok {
			sku = s
		}

		var productName string
		if n, ok := itemInfo["title"].(string); ok {
			productName = n
		}

		orderItems = append(orderItems, domain.CanonicalOrderItemDTO{
			ProductSKU:   sku,
			ProductName:  productName,
			ProductTitle: productName,
			Quantity:     quantity,
			UnitPrice:    unitPrice,
			TotalPrice:   totalPrice,
			Currency:     "COP",
			Tax:          tax,
			TaxRate:      floatPtr(0.19),
		})
	}

	// Calcular totales
	subtotal := getFloat("total_amount")
	shippingCost := getFloat("total_amount_with_shipping") - subtotal
	totalAmount := getFloat("paid_amount")

	// Obtener información del cliente
	customerName := ""
	customerEmail := ""
	customerPhone := ""
	if buyerData != nil {
		if firstName, ok := buyerData["first_name"].(string); ok {
			customerName = firstName
		}
		if lastName, ok := buyerData["last_name"].(string); ok && lastName != "" {
			customerName += " " + lastName
		}
		if email, ok := buyerData["email"].(string); ok {
			customerEmail = email
		}
		if phoneData, ok := buyerData["phone"].(map[string]interface{}); ok {
			if areaCode, ok := phoneData["area_code"].(string); ok {
				if number, ok := phoneData["number"].(string); ok {
					customerPhone = fmt.Sprintf("+%s%s", areaCode, number)
				}
			}
		}
	}

	// Obtener direcciones desde shipping
	addresses := []domain.CanonicalAddressDTO{}
	if shippingData != nil {
		if receiverAddr, ok := shippingData["receiver_address"].(map[string]interface{}); ok {
			lat := addressFloat(receiverAddr, "latitude")
			lng := addressFloat(receiverAddr, "longitude")

			cityName := ""
			if city, ok := receiverAddr["city"].(map[string]interface{}); ok {
				if name, ok := city["name"].(string); ok {
					cityName = name
				}
			}

			stateName := ""
			if state, ok := receiverAddr["state"].(map[string]interface{}); ok {
				if name, ok := state["name"].(string); ok {
					stateName = name
				}
			}

			countryName := ""
			if country, ok := receiverAddr["country"].(map[string]interface{}); ok {
				if name, ok := country["name"].(string); ok {
					countryName = name
				}
			}

			addresses = append(addresses, domain.CanonicalAddressDTO{
				Type:       "shipping",
				FirstName:  customerName,
				LastName:   "",
				Street:     addressString(receiverAddr, "address_line"),
				Street2:    addressString(receiverAddr, "street_number"),
				City:       cityName,
				State:      stateName,
				Country:    countryName,
				PostalCode: addressString(receiverAddr, "zip_code"),
				Latitude:   lat,
				Longitude:  lng,
			})

			// También crear billing address (usar la misma)
			addresses = append(addresses, domain.CanonicalAddressDTO{
				Type:       "billing",
				FirstName:  customerName,
				LastName:   "",
				Street:     addressString(receiverAddr, "address_line"),
				Street2:    addressString(receiverAddr, "street_number"),
				City:       cityName,
				State:      stateName,
				Country:    countryName,
				PostalCode: addressString(receiverAddr, "zip_code"),
			})
		}
	}

	// Mapear pagos
	payments := []domain.CanonicalPaymentDTO{}
	paymentsData, _ := meliJSON["payments"].([]interface{})
	if len(paymentsData) > 0 {
		if paymentData, ok := paymentsData[0].(map[string]interface{}); ok {
			status := "pending"
			if s, ok := paymentData["status"].(string); ok {
				if s == "approved" {
					status = "completed"
				}
			}

			paidAt := &dateCreated
			if dateApproved, ok := paymentData["date_approved"].(string); ok {
				if dt, err := time.Parse(time.RFC3339, dateApproved); err == nil {
					paidAt = &dt
				}
			}
			if status == "pending" {
				paidAt = nil
			}

			transactionAmount := 0.0
			if ta, ok := paymentData["transaction_amount"].(float64); ok {
				transactionAmount = ta
			}

			payments = append(payments, domain.CanonicalPaymentDTO{
				PaymentMethodID: 1, // Default
				Amount:          transactionAmount,
				Currency:        "COP",
				Status:          status,
				PaidAt:          paidAt,
				Gateway:         stringPtr("mercadopago"),
				TransactionID:   stringPtr(fmt.Sprintf("%.0f", paymentData["id"])),
			})
		}
	}

	// Serializar JSON de Mercado Libre para metadata
	rawDataJSON, _ := json.Marshal(meliJSON)

	// Crear orden canónica
	orderNumber := fmt.Sprintf("%d", getInt64("id"))
	externalID := orderNumber
	status := getString("status")

	return &domain.CanonicalOrderDTO{
		BusinessID:      businessID,
		IntegrationID:   integrationID,
		IntegrationType: "meli",
		Platform:        "meli",
		ExternalID:      externalID,
		OrderNumber:     orderNumber,
		Subtotal:        subtotal,
		Tax:             subtotal * 0.19, // Calcular impuesto
		Discount:        0,
		ShippingCost:    shippingCost,
		TotalAmount:     totalAmount,
		Currency:        "COP",
		CustomerName:    customerName,
		CustomerEmail:   customerEmail,
		CustomerPhone:   customerPhone,
		OrderTypeName:   "delivery",
		Status:          mapMeliStatus(status),
		OriginalStatus:  status,
		OrderItems:      orderItems,
		Addresses:       addresses,
		Payments:        payments,
		OccurredAt:      dateCreated,
		ImportedAt:      time.Now(),
		ChannelMetadata: &domain.CanonicalChannelMetadataDTO{
			ChannelSource: "meli",
			RawData:       datatypes.JSON(rawDataJSON),
			Version:       "1.0",
			ReceivedAt:    time.Now(),
			IsLatest:      true,
			SyncStatus:    "pending",
		},
	}, nil
}

func mapMeliStatus(status string) string {
	switch status {
	case "paid":
		return "paid"
	case "pending":
		return "pending"
	case "confirmed":
		return "confirmed"
	case "cancelled":
		return "cancelled"
	default:
		return "pending"
	}
}
