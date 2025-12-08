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

// ShopifyOrderJSON representa el JSON de una orden de Shopify
type ShopifyOrderJSON struct {
	ID                    int64                 `json:"id"`
	Email                 string                `json:"email"`
	ClosedAt              *string               `json:"closed_at"`
	CreatedAt             string                `json:"created_at"`
	UpdatedAt             string                `json:"updated_at"`
	Number                int                   `json:"number"`
	Note                  *string               `json:"note"`
	Token                 string                `json:"token"`
	Gateway               string                `json:"gateway"`
	Test                  bool                  `json:"test"`
	TotalPrice            string                `json:"total_price"`
	SubtotalPrice         string                `json:"subtotal_price"`
	TotalWeight           int                   `json:"total_weight"`
	TotalTax              string                `json:"total_tax"`
	TaxesIncluded         bool                  `json:"taxes_included"`
	Currency              string                `json:"currency"`
	FinancialStatus       string                `json:"financial_status"`
	Confirmed             bool                  `json:"confirmed"`
	TotalDiscounts        string                `json:"total_discounts"`
	TotalLineItemsPrice   string                `json:"total_line_items_price"`
	CartToken             string                `json:"cart_token"`
	BuyerAcceptsMarketing bool                  `json:"buyer_accepts_marketing"`
	Name                  string                `json:"name"`
	ReferringSite         string                `json:"referring_site"`
	LandingSite           string                `json:"landing_site"`
	CancelledAt           *string               `json:"cancelled_at"`
	CancelReason          *string               `json:"cancel_reason"`
	TotalPriceUSD         string                `json:"total_price_usd"`
	CheckoutToken         string                `json:"checkout_token"`
	Reference             *string               `json:"reference"`
	UserID                *int64                `json:"user_id"`
	LocationID            *int64                `json:"location_id"`
	SourceIdentifier      *string               `json:"source_identifier"`
	SourceURL             *string               `json:"source_url"`
	ProcessedAt           string                `json:"processed_at"`
	DeviceID              *int64                `json:"device_id"`
	Phone                 *string               `json:"phone"`
	CustomerLocale        string                `json:"customer_locale"`
	AppID                 int                   `json:"app_id"`
	BrowserIP             string                `json:"browser_ip"`
	LandingSiteRef        *string               `json:"landing_site_ref"`
	OrderNumber           int                   `json:"order_number"`
	DiscountApplications  []interface{}         `json:"discount_applications"`
	DiscountCodes         []interface{}         `json:"discount_codes"`
	NoteAttributes        []interface{}         `json:"note_attributes"`
	PaymentGatewayNames   []string              `json:"payment_gateway_names"`
	ProcessingMethod      string                `json:"processing_method"`
	CheckoutID            int64                 `json:"checkout_id"`
	SourceName            string                `json:"source_name"`
	FulfillmentStatus     *string               `json:"fulfillment_status"`
	TaxLines              []ShopifyTaxLine      `json:"tax_lines"`
	Tags                  string                `json:"tags"`
	ContactEmail          string                `json:"contact_email"`
	OrderStatusURL        string                `json:"order_status_url"`
	PresentmentCurrency   string                `json:"presentment_currency"`
	LineItems             []ShopifyLineItem     `json:"line_items"`
	ShippingLines         []ShopifyShippingLine `json:"shipping_lines"`
	BillingAddress        ShopifyAddress        `json:"billing_address"`
	ShippingAddress       ShopifyAddress        `json:"shipping_address"`
	Customer              ShopifyCustomer       `json:"customer"`
}

type ShopifyTaxLine struct {
	Price string  `json:"price"`
	Rate  float64 `json:"rate"`
	Title string  `json:"title"`
}

type ShopifyLineItem struct {
	ID                         int64            `json:"id"`
	VariantID                  int64            `json:"variant_id"`
	Title                      string           `json:"title"`
	Quantity                   int              `json:"quantity"`
	SKU                        string           `json:"sku"`
	VariantTitle               string           `json:"variant_title"`
	Vendor                     string           `json:"vendor"`
	FulfillmentService         string           `json:"fulfillment_service"`
	ProductID                  int64            `json:"product_id"`
	RequiresShipping           bool             `json:"requires_shipping"`
	Taxable                    bool             `json:"taxable"`
	GiftCard                   bool             `json:"gift_card"`
	Name                       string           `json:"name"`
	VariantInventoryManagement string           `json:"variant_inventory_management"`
	Properties                 []interface{}    `json:"properties"`
	ProductExists              bool             `json:"product_exists"`
	FulfillableQuantity        int              `json:"fulfillable_quantity"`
	Grams                      int              `json:"grams"`
	Price                      string           `json:"price"`
	TotalDiscount              string           `json:"total_discount"`
	FulfillmentStatus          *string          `json:"fulfillment_status"`
	PriceSet                   ShopifyPriceSet  `json:"price_set"`
	TaxLines                   []ShopifyTaxLine `json:"tax_lines"`
}

type ShopifyPriceSet struct {
	ShopMoney        ShopifyMoney `json:"shop_money"`
	PresentmentMoney ShopifyMoney `json:"presentment_money"`
}

type ShopifyMoney struct {
	Amount       string `json:"amount"`
	CurrencyCode string `json:"currency_code"`
}

type ShopifyShippingLine struct {
	ID                            int64   `json:"id"`
	Title                         string  `json:"title"`
	Price                         string  `json:"price"`
	Code                          string  `json:"code"`
	Source                        string  `json:"source"`
	Phone                         *string `json:"phone"`
	RequestedFulfillmentServiceID *int64  `json:"requested_fulfillment_service_id"`
	DeliveryCategory              *string `json:"delivery_category"`
	CarrierIdentifier             *string `json:"carrier_identifier"`
	DiscountedPrice               string  `json:"discounted_price"`
}

type ShopifyAddress struct {
	FirstName    string   `json:"first_name"`
	Address1     string   `json:"address1"`
	Phone        *string  `json:"phone"`
	City         string   `json:"city"`
	Zip          string   `json:"zip"`
	Province     string   `json:"province"`
	Country      string   `json:"country"`
	LastName     string   `json:"last_name"`
	Address2     *string  `json:"address2"`
	Company      *string  `json:"company"`
	Latitude     *float64 `json:"latitude"`
	Longitude    *float64 `json:"longitude"`
	Name         string   `json:"name"`
	CountryCode  string   `json:"country_code"`
	ProvinceCode string   `json:"province_code"`
}

type ShopifyCustomer struct {
	ID                  int64          `json:"id"`
	Email               string         `json:"email"`
	AcceptsMarketing    bool           `json:"accepts_marketing"`
	CreatedAt           string         `json:"created_at"`
	UpdatedAt           string         `json:"updated_at"`
	FirstName           string         `json:"first_name"`
	LastName            string         `json:"last_name"`
	OrdersCount         int            `json:"orders_count"`
	State               string         `json:"state"`
	TotalSpent          string         `json:"total_spent"`
	LastOrderID         int64          `json:"last_order_id"`
	Note                *string        `json:"note"`
	VerifiedEmail       bool           `json:"verified_email"`
	MultipassIdentifier *string        `json:"multipass_identifier"`
	TaxExempt           bool           `json:"tax_exempt"`
	Phone               *string        `json:"phone"`
	Tags                string         `json:"tags"`
	LastOrderName       string         `json:"last_order_name"`
	Currency            string         `json:"currency"`
	DefaultAddress      ShopifyAddress `json:"default_address"`
}

// GenerateShopifyOrderJSON genera un JSON de orden de Shopify simulado
func (g *OrderGenerator) GenerateShopifyOrderJSON(r *rand.Rand) (map[string]interface{}, error) {
	// Seleccionar cliente aleatorio
	customer := domain.FakeCustomers[r.Intn(len(domain.FakeCustomers))]

	// Seleccionar dirección aleatoria
	address := domain.FakeAddresses[r.Intn(len(domain.FakeAddresses))]

	// Generar 1-3 items aleatorios
	numItems := r.Intn(3) + 1
	lineItems := make([]map[string]interface{}, numItems)
	subtotal := 0.0
	totalTax := 0.0
	totalWeight := 0

	for i := 0; i < numItems; i++ {
		product := domain.FakeProducts[r.Intn(len(domain.FakeProducts))]
		quantity := r.Intn(2) + 1
		price, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", product.Price), 64)
		itemTotal := price * float64(quantity)
		tax := itemTotal * 0.19
		totalTax += tax
		subtotal += itemTotal
		totalWeight += int(product.Weight * 1000 * float64(quantity)) // Convertir a gramos

		lineItems[i] = map[string]interface{}{
			"id":                           450789469 + int64(i),
			"variant_id":                   808950810 + int64(i),
			"title":                        product.Name,
			"quantity":                     quantity,
			"sku":                          product.SKU,
			"variant_title":                fmt.Sprintf("Talla %s", []string{"S", "M", "L", "XL"}[r.Intn(4)]),
			"vendor":                       "Mi Marca",
			"fulfillment_service":          "manual",
			"product_id":                   632910392 + int64(i),
			"requires_shipping":            true,
			"taxable":                      true,
			"gift_card":                    false,
			"name":                         fmt.Sprintf("%s - %s", product.Name, product.Title),
			"variant_inventory_management": "shopify",
			"properties":                   []interface{}{},
			"product_exists":               true,
			"fulfillable_quantity":         quantity,
			"grams":                        int(product.Weight * 1000),
			"price":                        fmt.Sprintf("%.2f", price),
			"total_discount":               "0.00",
			"fulfillment_status":           nil,
			"price_set": map[string]interface{}{
				"shop_money": map[string]interface{}{
					"amount":        fmt.Sprintf("%.2f", price),
					"currency_code": "COP",
				},
				"presentment_money": map[string]interface{}{
					"amount":        fmt.Sprintf("%.2f", price),
					"currency_code": "COP",
				},
			},
			"tax_lines": []map[string]interface{}{
				{
					"title": "IVA",
					"price": fmt.Sprintf("%.2f", tax),
					"rate":  0.19,
					"price_set": map[string]interface{}{
						"shop_money": map[string]interface{}{
							"amount":        fmt.Sprintf("%.2f", tax),
							"currency_code": "COP",
						},
						"presentment_money": map[string]interface{}{
							"amount":        fmt.Sprintf("%.2f", tax),
							"currency_code": "COP",
						},
					},
				},
			},
		}
	}

	// Calcular shipping
	shippingCost := 4500.0
	if r.Float64() < 0.5 {
		shippingCost = float64(r.Intn(3000) + 3000)
	}
	totalAmount := subtotal + totalTax + shippingCost

	// Generar IDs
	orderID := 450789469 + int64(r.Intn(1000000))
	orderNumber := r.Intn(9999) + 1000
	token := fmt.Sprintf("%x", r.Int63())

	now := time.Now()
	createdAt := now.Add(-time.Duration(r.Intn(24)) * time.Hour).Format(time.RFC3339)
	updatedAt := now.Add(-time.Duration(r.Intn(5)) * time.Hour).Format(time.RFC3339)
	processedAt := createdAt

	// Estado financiero aleatorio
	financialStatuses := []string{"paid", "pending", "authorized", "partially_paid", "refunded"}
	financialStatus := financialStatuses[r.Intn(len(financialStatuses))]

	// Construir JSON de Shopify
	shopifyOrder := map[string]interface{}{
		"id":                      orderID,
		"email":                   customer.Email,
		"closed_at":               nil,
		"created_at":              createdAt,
		"updated_at":              updatedAt,
		"number":                  orderNumber,
		"note":                    nil,
		"token":                   token,
		"gateway":                 "shopify_payments",
		"test":                    false,
		"total_price":             fmt.Sprintf("%.2f", totalAmount),
		"subtotal_price":          fmt.Sprintf("%.2f", subtotal),
		"total_weight":            totalWeight,
		"total_tax":               fmt.Sprintf("%.2f", totalTax),
		"taxes_included":          true,
		"currency":                "COP",
		"financial_status":        financialStatus,
		"confirmed":               true,
		"total_discounts":         "0.00",
		"total_line_items_price":  fmt.Sprintf("%.2f", subtotal),
		"cart_token":              fmt.Sprintf("%x", r.Int63()),
		"buyer_accepts_marketing": false,
		"name":                    fmt.Sprintf("#%d", orderNumber),
		"referring_site":          "",
		"landing_site":            "/",
		"cancelled_at":            nil,
		"cancel_reason":           nil,
		"total_price_usd":         fmt.Sprintf("%.2f", totalAmount/3900), // Aproximado
		"checkout_token":          fmt.Sprintf("%x", r.Int63()),
		"reference":               nil,
		"user_id":                 nil,
		"location_id":             nil,
		"source_identifier":       nil,
		"source_url":              nil,
		"processed_at":            processedAt,
		"device_id":               nil,
		"phone":                   customer.Phone,
		"customer_locale":         "es",
		"app_id":                  580111,
		"browser_ip":              fmt.Sprintf("192.168.1.%d", r.Intn(255)),
		"landing_site_ref":        nil,
		"order_number":            orderNumber,
		"discount_applications":   []interface{}{},
		"discount_codes":          []interface{}{},
		"note_attributes":         []interface{}{},
		"payment_gateway_names":   []string{"shopify_payments"},
		"processing_method":       "direct",
		"checkout_id":             901414060 + int64(r.Intn(1000)),
		"source_name":             "web",
		"fulfillment_status":      nil,
		"tax_lines": []map[string]interface{}{
			{
				"price": fmt.Sprintf("%.2f", totalTax),
				"rate":  0.19,
				"title": "IVA",
			},
		},
		"tags":                 "",
		"contact_email":        customer.Email,
		"order_status_url":     fmt.Sprintf("https://mi-tienda.myshopify.com/548380009/orders/%s/authenticate?key=...", token),
		"presentment_currency": "COP",
		"line_items":           lineItems,
		"shipping_lines": []map[string]interface{}{
			{
				"id":                               123456789 + int64(r.Intn(1000)),
				"title":                            "Envío Estándar",
				"price":                            fmt.Sprintf("%.2f", shippingCost),
				"code":                             "Standard",
				"source":                           "shopify",
				"phone":                            nil,
				"requested_fulfillment_service_id": nil,
				"delivery_category":                nil,
				"carrier_identifier":               nil,
				"discounted_price":                 fmt.Sprintf("%.2f", shippingCost),
			},
		},
		"billing_address": map[string]interface{}{
			"first_name":    customer.Name,
			"address1":      address.Street,
			"phone":         customer.Phone,
			"city":          address.City,
			"zip":           address.PostalCode,
			"province":      address.State,
			"country":       address.Country,
			"last_name":     "",
			"address2":      address.Street2,
			"company":       nil,
			"latitude":      address.Latitude,
			"longitude":     address.Longitude,
			"name":          customer.Name,
			"country_code":  "CO",
			"province_code": "DC",
		},
		"shipping_address": map[string]interface{}{
			"first_name":    customer.Name,
			"address1":      address.Street,
			"phone":         customer.Phone,
			"city":          address.City,
			"zip":           address.PostalCode,
			"province":      address.State,
			"country":       address.Country,
			"last_name":     "",
			"address2":      address.Street2,
			"company":       nil,
			"latitude":      address.Latitude,
			"longitude":     address.Longitude,
			"name":          customer.Name,
			"country_code":  "CO",
			"province_code": "DC",
		},
		"customer": map[string]interface{}{
			"id":                   207119551 + int64(r.Intn(1000)),
			"email":                customer.Email,
			"accepts_marketing":    false,
			"created_at":           createdAt,
			"updated_at":           updatedAt,
			"first_name":           customer.Name,
			"last_name":            "",
			"orders_count":         r.Intn(10) + 1,
			"state":                "enabled",
			"total_spent":          fmt.Sprintf("%.2f", totalAmount),
			"last_order_id":        orderID,
			"note":                 nil,
			"verified_email":       true,
			"multipass_identifier": nil,
			"tax_exempt":           false,
			"phone":                customer.Phone,
			"tags":                 "VIP",
			"last_order_name":      fmt.Sprintf("#%d", orderNumber),
			"currency":             "COP",
			"default_address": map[string]interface{}{
				"id":            207119551 + int64(r.Intn(1000)),
				"customer_id":   207119551 + int64(r.Intn(1000)),
				"first_name":    customer.Name,
				"last_name":     "",
				"company":       nil,
				"address1":      address.Street,
				"address2":      address.Street2,
				"city":          address.City,
				"province":      address.State,
				"country":       address.Country,
				"zip":           address.PostalCode,
				"phone":         customer.Phone,
				"name":          customer.Name,
				"province_code": "DC",
				"country_code":  "CO",
				"country_name":  "Colombia",
				"default":       true,
			},
		},
	}

	return shopifyOrder, nil
}

// MapShopifyJSONToCanonical mapea un JSON de Shopify al formato canónico
func (g *OrderGenerator) MapShopifyJSONToCanonical(shopifyJSON map[string]interface{}, integrationID uint, businessID *uint) (*domain.CanonicalOrderDTO, error) {
	// Helper para obtener valores
	getString := func(key string) string {
		if v, ok := shopifyJSON[key].(string); ok {
			return v
		}
		return ""
	}
	getFloat := func(key string) float64 {
		if v, ok := shopifyJSON[key].(string); ok {
			if f, err := strconv.ParseFloat(v, 64); err == nil {
				return f
			}
		}
		return 0
	}
	getInt := func(key string) int {
		if v, ok := shopifyJSON[key].(float64); ok {
			return int(v)
		}
		if v, ok := shopifyJSON[key].(int64); ok {
			return int(v)
		}
		return 0
	}

	// Parsear fechas
	createdAtStr := getString("created_at")
	createdAt, _ := time.Parse(time.RFC3339, createdAtStr)
	if createdAt.IsZero() {
		createdAt = time.Now()
	}

	// Obtener direcciones
	billingAddr, _ := shopifyJSON["billing_address"].(map[string]interface{})
	shippingAddr, _ := shopifyJSON["shipping_address"].(map[string]interface{})
	customerData, _ := shopifyJSON["customer"].(map[string]interface{})

	// Mapear line items
	lineItemsData, _ := shopifyJSON["line_items"].([]interface{})
	orderItems := make([]domain.CanonicalOrderItemDTO, 0, len(lineItemsData))

	for _, itemData := range lineItemsData {
		item, _ := itemData.(map[string]interface{})
		price := getFloat("price")
		quantity := 0
		if q, ok := item["quantity"].(float64); ok {
			quantity = int(q)
		} else if q, ok := item["quantity"].(int); ok {
			quantity = q
		}
		totalPrice := price * float64(quantity)
		tax := totalPrice * 0.19

		var sku string
		if s, ok := item["sku"].(string); ok {
			sku = s
		}

		var productName string
		if n, ok := item["title"].(string); ok {
			productName = n
		}

		var variantTitle string
		if vt, ok := item["variant_title"].(string); ok {
			variantTitle = vt
		}

		orderItems = append(orderItems, domain.CanonicalOrderItemDTO{
			ProductSKU:   sku,
			ProductName:  productName,
			ProductTitle: fmt.Sprintf("%s - %s", productName, variantTitle),
			Quantity:     quantity,
			UnitPrice:    price,
			TotalPrice:   totalPrice,
			Currency:     "COP",
			Tax:          tax,
			TaxRate:      floatPtr(0.19),
		})
	}

	// Calcular totales
	subtotal := getFloat("subtotal_price")
	tax := getFloat("total_tax")
	shippingCost := 0.0
	if shippingLines, ok := shopifyJSON["shipping_lines"].([]interface{}); ok && len(shippingLines) > 0 {
		if shippingLine, ok := shippingLines[0].(map[string]interface{}); ok {
			if price, ok := shippingLine["price"].(string); ok {
				if f, err := strconv.ParseFloat(price, 64); err == nil {
					shippingCost = f
				}
			}
		}
	}
	totalAmount := getFloat("total_price")

	// Obtener información del cliente
	customerName := ""
	customerEmail := getString("email")
	customerPhone := ""
	if customerData != nil {
		if name, ok := customerData["first_name"].(string); ok {
			customerName = name
		}
		if email, ok := customerData["email"].(string); ok && customerEmail == "" {
			customerEmail = email
		}
		if phone, ok := customerData["phone"].(string); ok {
			customerPhone = phone
		}
	}

	// Obtener direcciones
	addresses := []domain.CanonicalAddressDTO{}
	if shippingAddr != nil {
		lat := addressFloat(shippingAddr, "latitude")
		lng := addressFloat(shippingAddr, "longitude")
		addresses = append(addresses, domain.CanonicalAddressDTO{
			Type:       "shipping",
			FirstName:  addressString(shippingAddr, "first_name"),
			LastName:   addressString(shippingAddr, "last_name"),
			Street:     addressString(shippingAddr, "address1"),
			Street2:    addressString(shippingAddr, "address2"),
			City:       addressString(shippingAddr, "city"),
			State:      addressString(shippingAddr, "province"),
			Country:    addressString(shippingAddr, "country"),
			PostalCode: addressString(shippingAddr, "zip"),
			Phone:      addressString(shippingAddr, "phone"),
			Latitude:   lat,
			Longitude:  lng,
		})
	}
	if billingAddr != nil {
		addresses = append(addresses, domain.CanonicalAddressDTO{
			Type:       "billing",
			FirstName:  addressString(billingAddr, "first_name"),
			LastName:   addressString(billingAddr, "last_name"),
			Street:     addressString(billingAddr, "address1"),
			Street2:    addressString(billingAddr, "address2"),
			City:       addressString(billingAddr, "city"),
			State:      addressString(billingAddr, "province"),
			Country:    addressString(billingAddr, "country"),
			PostalCode: addressString(billingAddr, "zip"),
			Phone:      addressString(billingAddr, "phone"),
		})
	}

	// Mapear pagos
	payments := []domain.CanonicalPaymentDTO{}
	financialStatus := getString("financial_status")
	paymentStatus := "pending"
	if financialStatus == "paid" {
		paymentStatus = "completed"
	}
	paidAt := &createdAt
	if paymentStatus == "pending" {
		paidAt = nil
	}
	payments = append(payments, domain.CanonicalPaymentDTO{
		PaymentMethodID: 1, // Default
		Amount:          totalAmount,
		Currency:        "COP",
		Status:          paymentStatus,
		PaidAt:          paidAt,
		Gateway:         stringPtr("shopify_payments"),
	})

	// Serializar JSON de Shopify para metadata
	rawDataJSON, _ := json.Marshal(shopifyJSON)

	// Crear orden canónica
	orderNumber := fmt.Sprintf("%d", getInt("order_number"))
	var externalID string
	if id, ok := shopifyJSON["id"].(float64); ok {
		externalID = fmt.Sprintf("%.0f", id)
	} else if id, ok := shopifyJSON["id"].(int64); ok {
		externalID = fmt.Sprintf("%d", id)
	} else {
		externalID = fmt.Sprintf("%v", shopifyJSON["id"])
	}

	return &domain.CanonicalOrderDTO{
		BusinessID:      businessID,
		IntegrationID:   integrationID,
		IntegrationType: "shopify",
		Platform:        "shopify",
		ExternalID:      externalID,
		OrderNumber:     orderNumber,
		Subtotal:        subtotal,
		Tax:             tax,
		Discount:        0,
		ShippingCost:    shippingCost,
		TotalAmount:     totalAmount,
		Currency:        "COP",
		CustomerName:    customerName,
		CustomerEmail:   customerEmail,
		CustomerPhone:   customerPhone,
		OrderTypeName:   "delivery",
		Status:          mapShopifyStatus(financialStatus),
		OriginalStatus:  financialStatus,
		OrderItems:      orderItems,
		Addresses:       addresses,
		Payments:        payments,
		OccurredAt:      createdAt,
		ImportedAt:      time.Now(),
		ChannelMetadata: &domain.CanonicalChannelMetadataDTO{
			ChannelSource: "shopify",
			RawData:       datatypes.JSON(rawDataJSON),
			Version:       "1.0",
			ReceivedAt:    time.Now(),
			IsLatest:      true,
			SyncStatus:    "pending",
		},
	}, nil
}

func addressString(addr map[string]interface{}, key string) string {
	if v, ok := addr[key].(string); ok {
		return v
	}
	return ""
}

func addressFloat(addr map[string]interface{}, key string) *float64 {
	if v, ok := addr[key].(float64); ok {
		return &v
	}
	return nil
}

func floatPtr(f float64) *float64 {
	return &f
}

func mapShopifyStatus(status string) string {
	switch status {
	case "paid":
		return "paid"
	case "pending":
		return "pending"
	case "authorized":
		return "authorized"
	case "partially_paid":
		return "partially_paid"
	case "refunded":
		return "refunded"
	default:
		return "pending"
	}
}
