package mappers

import (
	"encoding/json"
	"math"
	"testing"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/infra/secondary/client/response"
)

// almostEqual compara dos float64 con tolerancia de 0.01
func almostEqual(a, b float64) bool {
	return math.Abs(a-b) < 0.01
}

// buildMinimalOrder construye una orden base con los campos mínimos para que el mapper no falle
func buildMinimalOrder(lineItems []response.LineItem) response.Order {
	return response.Order{
		ID:              1234567890,
		Name:            "#67269",
		OrderNumber:     67269,
		Email:           "test@example.com",
		Currency:        "USD",
		TotalPrice:      "44.99",
		SubtotalPrice:   "41.80",
		TotalTax:        "6.67",
		TotalDiscounts:  "4.70",
		FinancialStatus: "paid",
		SourceName:      "web",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		ProcessedAt:     time.Now(),
		PresentmentCurrency: "COP",
		TotalPriceSet: &response.MoneySet{
			ShopMoney:        response.Money{Amount: "44.99", CurrencyCode: "USD"},
			PresentmentMoney: response.Money{Amount: "169233.35", CurrencyCode: "COP"},
		},
		SubtotalPriceSet: &response.MoneySet{
			ShopMoney:        response.Money{Amount: "41.80", CurrencyCode: "USD"},
			PresentmentMoney: response.Money{Amount: "157233.35", CurrencyCode: "COP"},
		},
		TotalTaxSet: &response.MoneySet{
			ShopMoney:        response.Money{Amount: "6.67", CurrencyCode: "USD"},
			PresentmentMoney: response.Money{Amount: "25086.54", CurrencyCode: "COP"},
		},
		TotalDiscountsSet: &response.MoneySet{
			ShopMoney:        response.Money{Amount: "4.70", CurrencyCode: "USD"},
			PresentmentMoney: response.Money{Amount: "17666.65", CurrencyCode: "COP"},
		},
		TotalShippingPriceSet: &response.MoneySet{
			ShopMoney:        response.Money{Amount: "3.19", CurrencyCode: "USD"},
			PresentmentMoney: response.Money{Amount: "12000.00", CurrencyCode: "COP"},
		},
		LineItems: lineItems,
	}
}

// TestDiscountFromDiscountAllocations verifica que los descuentos automaticos
// de Shopify (discount_allocations) se extraen correctamente cuando total_discount es "0.00".
// Caso real: orden #67269 con "Descuento Miembro" 10.101% automatico.
func TestDiscountFromDiscountAllocations(t *testing.T) {
	lineItems := []response.LineItem{
		{
			ID:            12841714417873,
			Name:          "Colageno Hidrolizado - 300g",
			SKU:           "PT01015",
			Title:         "Colageno Hidrolizado",
			Quantity:      1,
			Price:         "13.16",
			Grams:         300,
			ProductID:     7235290136785,
			VariantID:     42749434233041,
			TotalDiscount: "0.00",
			TotalDiscountSet: &response.MoneySet{
				ShopMoney:        response.Money{Amount: "0.00", CurrencyCode: "USD"},
				PresentmentMoney: response.Money{Amount: "0.00", CurrencyCode: "COP"},
			},
			PriceSet: &response.MoneySet{
				ShopMoney:        response.Money{Amount: "13.16", CurrencyCode: "USD"},
				PresentmentMoney: response.Money{Amount: "49500.00", CurrencyCode: "COP"},
			},
			TaxLines: []response.TaxLine{
				{
					Rate:  0.19,
					Price: "1.88",
					Title: "VAT",
					PriceSet: &response.MoneySet{
						ShopMoney:        response.Money{Amount: "1.88", CurrencyCode: "USD"},
						PresentmentMoney: response.Money{Amount: "7070.87", CurrencyCode: "COP"},
					},
				},
			},
			DiscountAllocations: []response.DiscountAllocation{
				{
					Amount: "1.33",
					AmountSet: &response.MoneySet{
						ShopMoney:        response.Money{Amount: "1.33", CurrencyCode: "USD"},
						PresentmentMoney: response.Money{Amount: "4999.99", CurrencyCode: "COP"},
					},
					DiscountApplicationIndex: 1,
				},
			},
		},
		{
			ID:            12841714450641,
			Name:          "Creatina Monohidrato - 300g",
			SKU:           "PT01004",
			Title:         "Creatina Monohidrato",
			Quantity:      2,
			Price:         "16.67",
			Grams:         300,
			ProductID:     6866571296977,
			VariantID:     42749438034129,
			TotalDiscount: "0.00",
			TotalDiscountSet: &response.MoneySet{
				ShopMoney:        response.Money{Amount: "0.00", CurrencyCode: "USD"},
				PresentmentMoney: response.Money{Amount: "0.00", CurrencyCode: "COP"},
			},
			PriceSet: &response.MoneySet{
				ShopMoney:        response.Money{Amount: "16.67", CurrencyCode: "USD"},
				PresentmentMoney: response.Money{Amount: "62700.00", CurrencyCode: "COP"},
			},
			TaxLines: []response.TaxLine{
				{
					Rate:  0.19,
					Price: "4.79",
					Title: "VAT",
					PriceSet: &response.MoneySet{
						ShopMoney:        response.Money{Amount: "4.79", CurrencyCode: "USD"},
						PresentmentMoney: response.Money{Amount: "18015.67", CurrencyCode: "COP"},
					},
				},
			},
			DiscountAllocations: []response.DiscountAllocation{
				{
					Amount: "3.37",
					AmountSet: &response.MoneySet{
						ShopMoney:        response.Money{Amount: "3.37", CurrencyCode: "USD"},
						PresentmentMoney: response.Money{Amount: "12666.66", CurrencyCode: "COP"},
					},
					DiscountApplicationIndex: 0,
				},
			},
		},
	}

	order := buildMinimalOrder(lineItems)
	businessID := uint(34)
	rawOrder, _ := json.Marshal(order)

	result := MapOrderResponseToShopifyOrder(order, rawOrder, &businessID, 1, "shopify")

	// Item 0: Colageno - debe tener discount 1.33 USD / 4999.99 COP
	if !almostEqual(result.Items[0].Discount, 1.33) {
		t.Errorf("Item[0].Discount: esperado 1.33, obtenido %.2f", result.Items[0].Discount)
	}
	if !almostEqual(result.Items[0].DiscountPresentment, 4999.99) {
		t.Errorf("Item[0].DiscountPresentment: esperado 4999.99, obtenido %.2f", result.Items[0].DiscountPresentment)
	}

	// Item 1: Creatina x2 - debe tener discount 3.37 USD / 12666.66 COP
	if !almostEqual(result.Items[1].Discount, 3.37) {
		t.Errorf("Item[1].Discount: esperado 3.37, obtenido %.2f", result.Items[1].Discount)
	}
	if !almostEqual(result.Items[1].DiscountPresentment, 12666.66) {
		t.Errorf("Item[1].DiscountPresentment: esperado 12666.66, obtenido %.2f", result.Items[1].DiscountPresentment)
	}

	// Discount total de la orden
	if !almostEqual(result.Discount, 4.70) {
		t.Errorf("Order.Discount: esperado 4.70, obtenido %.2f", result.Discount)
	}
	if !almostEqual(result.DiscountPresentment, 17666.65) {
		t.Errorf("Order.DiscountPresentment: esperado 17666.65, obtenido %.2f", result.DiscountPresentment)
	}
}

// TestDiscountFromTotalDiscount verifica que cuando total_discount tiene valor,
// se usa directamente (sin sumar discount_allocations).
// Caso: descuento aplicado directamente al line item (no automatico).
func TestDiscountFromTotalDiscount(t *testing.T) {
	lineItems := []response.LineItem{
		{
			ID:            100,
			Name:          "Producto con descuento directo",
			SKU:           "PROD001",
			Title:         "Producto Test",
			Quantity:      2,
			Price:         "100.00",
			ProductID:     1,
			VariantID:     1,
			TotalDiscount: "20.00",
			TotalDiscountSet: &response.MoneySet{
				ShopMoney:        response.Money{Amount: "20.00", CurrencyCode: "USD"},
				PresentmentMoney: response.Money{Amount: "75000.00", CurrencyCode: "COP"},
			},
			PriceSet: &response.MoneySet{
				ShopMoney:        response.Money{Amount: "100.00", CurrencyCode: "USD"},
				PresentmentMoney: response.Money{Amount: "376200.00", CurrencyCode: "COP"},
			},
			// Tiene discount_allocations pero total_discount > 0, asi que se ignoran
			DiscountAllocations: []response.DiscountAllocation{
				{
					Amount: "20.00",
					AmountSet: &response.MoneySet{
						ShopMoney:        response.Money{Amount: "20.00", CurrencyCode: "USD"},
						PresentmentMoney: response.Money{Amount: "75000.00", CurrencyCode: "COP"},
					},
					DiscountApplicationIndex: 0,
				},
			},
		},
	}

	order := buildMinimalOrder(lineItems)
	businessID := uint(1)
	rawOrder, _ := json.Marshal(order)

	result := MapOrderResponseToShopifyOrder(order, rawOrder, &businessID, 1, "shopify")

	// Debe usar total_discount (20.00), NO sumar discount_allocations
	if !almostEqual(result.Items[0].Discount, 20.00) {
		t.Errorf("Item[0].Discount: esperado 20.00, obtenido %.2f", result.Items[0].Discount)
	}
	if !almostEqual(result.Items[0].DiscountPresentment, 75000.00) {
		t.Errorf("Item[0].DiscountPresentment: esperado 75000.00, obtenido %.2f", result.Items[0].DiscountPresentment)
	}
}

// TestNoDiscountAtAll verifica que items sin descuento quedan en 0.
func TestNoDiscountAtAll(t *testing.T) {
	lineItems := []response.LineItem{
		{
			ID:            200,
			Name:          "Producto sin descuento",
			SKU:           "PROD002",
			Title:         "Producto Normal",
			Quantity:      1,
			Price:         "50.00",
			ProductID:     2,
			VariantID:     2,
			TotalDiscount: "0.00",
			TotalDiscountSet: &response.MoneySet{
				ShopMoney:        response.Money{Amount: "0.00", CurrencyCode: "USD"},
				PresentmentMoney: response.Money{Amount: "0.00", CurrencyCode: "COP"},
			},
			PriceSet: &response.MoneySet{
				ShopMoney:        response.Money{Amount: "50.00", CurrencyCode: "USD"},
				PresentmentMoney: response.Money{Amount: "188100.00", CurrencyCode: "COP"},
			},
			DiscountAllocations: []response.DiscountAllocation{},
		},
	}

	order := buildMinimalOrder(lineItems)
	businessID := uint(1)
	rawOrder, _ := json.Marshal(order)

	result := MapOrderResponseToShopifyOrder(order, rawOrder, &businessID, 1, "shopify")

	if result.Items[0].Discount != 0 {
		t.Errorf("Item[0].Discount: esperado 0, obtenido %.2f", result.Items[0].Discount)
	}
	if result.Items[0].DiscountPresentment != 0 {
		t.Errorf("Item[0].DiscountPresentment: esperado 0, obtenido %.2f", result.Items[0].DiscountPresentment)
	}
}

// TestMultipleDiscountAllocationsPerItem verifica que se suman correctamente
// multiples discount_allocations en un solo item (ej: descuento miembro + cupon).
func TestMultipleDiscountAllocationsPerItem(t *testing.T) {
	lineItems := []response.LineItem{
		{
			ID:            300,
			Name:          "Producto con 2 descuentos",
			SKU:           "PROD003",
			Title:         "Producto Multi-Descuento",
			Quantity:      1,
			Price:         "100.00",
			ProductID:     3,
			VariantID:     3,
			TotalDiscount: "0.00",
			TotalDiscountSet: &response.MoneySet{
				ShopMoney:        response.Money{Amount: "0.00", CurrencyCode: "USD"},
				PresentmentMoney: response.Money{Amount: "0.00", CurrencyCode: "COP"},
			},
			PriceSet: &response.MoneySet{
				ShopMoney:        response.Money{Amount: "100.00", CurrencyCode: "USD"},
				PresentmentMoney: response.Money{Amount: "376200.00", CurrencyCode: "COP"},
			},
			DiscountAllocations: []response.DiscountAllocation{
				{
					Amount: "5.00",
					AmountSet: &response.MoneySet{
						ShopMoney:        response.Money{Amount: "5.00", CurrencyCode: "USD"},
						PresentmentMoney: response.Money{Amount: "18810.00", CurrencyCode: "COP"},
					},
					DiscountApplicationIndex: 0,
				},
				{
					Amount: "3.00",
					AmountSet: &response.MoneySet{
						ShopMoney:        response.Money{Amount: "3.00", CurrencyCode: "USD"},
						PresentmentMoney: response.Money{Amount: "11286.00", CurrencyCode: "COP"},
					},
					DiscountApplicationIndex: 1,
				},
			},
		},
	}

	order := buildMinimalOrder(lineItems)
	businessID := uint(1)
	rawOrder, _ := json.Marshal(order)

	result := MapOrderResponseToShopifyOrder(order, rawOrder, &businessID, 1, "shopify")

	// 5.00 + 3.00 = 8.00
	if !almostEqual(result.Items[0].Discount, 8.00) {
		t.Errorf("Item[0].Discount: esperado 8.00, obtenido %.2f", result.Items[0].Discount)
	}
	// 18810.00 + 11286.00 = 30096.00
	if !almostEqual(result.Items[0].DiscountPresentment, 30096.00) {
		t.Errorf("Item[0].DiscountPresentment: esperado 30096.00, obtenido %.2f", result.Items[0].DiscountPresentment)
	}
}

// TestDiscountAllocationWithoutAmountSet verifica que si un discount_allocation
// no tiene amount_set (presentment_money), solo se suma el shop_money.
func TestDiscountAllocationWithoutAmountSet(t *testing.T) {
	lineItems := []response.LineItem{
		{
			ID:            400,
			Name:          "Producto allocation sin AmountSet",
			SKU:           "PROD004",
			Title:         "Producto Test",
			Quantity:      1,
			Price:         "50.00",
			ProductID:     4,
			VariantID:     4,
			TotalDiscount: "0.00",
			PriceSet: &response.MoneySet{
				ShopMoney:        response.Money{Amount: "50.00", CurrencyCode: "USD"},
				PresentmentMoney: response.Money{Amount: "188100.00", CurrencyCode: "COP"},
			},
			DiscountAllocations: []response.DiscountAllocation{
				{
					Amount:    "2.50",
					AmountSet: nil, // sin presentment
					DiscountApplicationIndex: 0,
				},
			},
		},
	}

	order := buildMinimalOrder(lineItems)
	businessID := uint(1)
	rawOrder, _ := json.Marshal(order)

	result := MapOrderResponseToShopifyOrder(order, rawOrder, &businessID, 1, "shopify")

	// shop_money discount debe funcionar
	if !almostEqual(result.Items[0].Discount, 2.50) {
		t.Errorf("Item[0].Discount: esperado 2.50, obtenido %.2f", result.Items[0].Discount)
	}
	// presentment debe quedar en 0 porque no hay AmountSet
	if result.Items[0].DiscountPresentment != 0 {
		t.Errorf("Item[0].DiscountPresentment: esperado 0, obtenido %.2f", result.Items[0].DiscountPresentment)
	}
}

// TestRealOrder67269FullMapping test de integracion con datos reales de la orden #67269.
// Verifica el mapeo completo incluyendo montos, descuentos, impuestos y shipping.
func TestRealOrder67269FullMapping(t *testing.T) {
	rawJSON := `{
		"id": 5178312229073,
		"name": "#67269",
		"email": "krito2107@hotmail.com",
		"currency": "USD",
		"total_price": "44.99",
		"subtotal_price": "41.80",
		"total_tax": "6.67",
		"total_discounts": "4.70",
		"financial_status": "paid",
		"source_name": "web",
		"order_number": 67269,
		"created_at": "2026-01-05T23:45:14-05:00",
		"updated_at": "2026-01-06T13:42:19-05:00",
		"processed_at": "2026-01-05T23:45:12-05:00",
		"presentment_currency": "COP",
		"tags": "Melonn, Melonn-Entregado",
		"payment_gateway_names": ["shopify_payments"],
		"total_price_set": {
			"shop_money": {"amount": "44.99", "currency_code": "USD"},
			"presentment_money": {"amount": "169233.35", "currency_code": "COP"}
		},
		"subtotal_price_set": {
			"shop_money": {"amount": "41.80", "currency_code": "USD"},
			"presentment_money": {"amount": "157233.35", "currency_code": "COP"}
		},
		"total_tax_set": {
			"shop_money": {"amount": "6.67", "currency_code": "USD"},
			"presentment_money": {"amount": "25086.54", "currency_code": "COP"}
		},
		"total_discounts_set": {
			"shop_money": {"amount": "4.70", "currency_code": "USD"},
			"presentment_money": {"amount": "17666.65", "currency_code": "COP"}
		},
		"total_shipping_price_set": {
			"shop_money": {"amount": "3.19", "currency_code": "USD"},
			"presentment_money": {"amount": "12000.00", "currency_code": "COP"}
		},
		"customer": {
			"id": 6268602450129,
			"email": "krito2107@hotmail.com",
			"first_name": "Carolina",
			"last_name": "Ramirez",
			"phone": "+573178949340",
			"state": "enabled",
			"currency": "COP",
			"verified_email": true,
			"created_at": "2024-06-25T21:06:37-05:00",
			"updated_at": "2026-01-22T14:59:42-05:00",
			"default_address": {
				"address1": "Cra 50 n 144-61",
				"address2": "Apt 311",
				"city": "Bogota",
				"province": "Bogota, D.C.",
				"country": "Colombia",
				"zip": "111156",
				"company": "1019061902",
				"country_code": "CO",
				"province_code": "DC",
				"first_name": "Carolina",
				"last_name": "Ramirez",
				"name": "Carolina Ramirez"
			}
		},
		"shipping_address": {
			"address1": "Cra 50 n 144-61",
			"address2": "Apt 311",
			"city": "Bogota",
			"province": "Bogota, D.C.",
			"country": "Colombia",
			"zip": "111156",
			"company": "1019061902",
			"phone": "3178949340",
			"first_name": "Carolina",
			"last_name": "Ramirez",
			"country_code": "CO",
			"province_code": "DC",
			"latitude": 4.7278517,
			"longitude": -74.05411819999999,
			"name": "Carolina Ramirez"
		},
		"billing_address": {
			"address1": "Cra 50 n 144-61",
			"address2": "Apt 311",
			"city": "Bogota",
			"province": "Bogota, D.C.",
			"country": "Colombia",
			"zip": "111156",
			"company": "1019061902",
			"phone": "3178949340",
			"first_name": "Carolina",
			"last_name": "Ramirez",
			"country_code": "CO",
			"province_code": "DC",
			"latitude": 4.7278517,
			"longitude": -74.05411819999999,
			"name": "Carolina Ramirez"
		},
		"line_items": [
			{
				"id": 12841714417873,
				"sku": "PT01015",
				"name": "Colageno Hidrolizado - 300g",
				"grams": 300,
				"price": "13.16",
				"title": "Colageno Hidrolizado",
				"vendor": "Fitness Food SaS",
				"taxable": true,
				"quantity": 1,
				"gift_card": false,
				"price_set": {
					"shop_money": {"amount": "13.16", "currency_code": "USD"},
					"presentment_money": {"amount": "49500.00", "currency_code": "COP"}
				},
				"tax_lines": [{
					"rate": 0.19,
					"price": "1.88",
					"title": "VAT",
					"price_set": {
						"shop_money": {"amount": "1.88", "currency_code": "USD"},
						"presentment_money": {"amount": "7070.87", "currency_code": "COP"}
					}
				}],
				"product_id": 7235290136785,
				"variant_id": 42749434233041,
				"variant_title": "300g",
				"total_discount": "0.00",
				"requires_shipping": true,
				"fulfillment_status": "fulfilled",
				"total_discount_set": {
					"shop_money": {"amount": "0.00", "currency_code": "USD"},
					"presentment_money": {"amount": "0.00", "currency_code": "COP"}
				},
				"discount_allocations": [{
					"amount": "1.33",
					"amount_set": {
						"shop_money": {"amount": "1.33", "currency_code": "USD"},
						"presentment_money": {"amount": "4999.99", "currency_code": "COP"}
					},
					"discount_application_index": 1
				}]
			},
			{
				"id": 12841714450641,
				"sku": "PT01004",
				"name": "Creatina Monohidrato - 300g",
				"grams": 300,
				"price": "16.67",
				"title": "Creatina Monohidrato",
				"vendor": "Fitness Food SaS",
				"taxable": true,
				"quantity": 2,
				"gift_card": false,
				"price_set": {
					"shop_money": {"amount": "16.67", "currency_code": "USD"},
					"presentment_money": {"amount": "62700.00", "currency_code": "COP"}
				},
				"tax_lines": [{
					"rate": 0.19,
					"price": "4.79",
					"title": "VAT",
					"price_set": {
						"shop_money": {"amount": "4.79", "currency_code": "USD"},
						"presentment_money": {"amount": "18015.67", "currency_code": "COP"}
					}
				}],
				"product_id": 6866571296977,
				"variant_id": 42749438034129,
				"variant_title": "300g",
				"total_discount": "0.00",
				"requires_shipping": true,
				"fulfillment_status": "fulfilled",
				"total_discount_set": {
					"shop_money": {"amount": "0.00", "currency_code": "USD"},
					"presentment_money": {"amount": "0.00", "currency_code": "COP"}
				},
				"discount_allocations": [{
					"amount": "3.37",
					"amount_set": {
						"shop_money": {"amount": "3.37", "currency_code": "USD"},
						"presentment_money": {"amount": "12666.66", "currency_code": "COP"}
					},
					"discount_application_index": 0
				}]
			}
		],
		"shipping_lines": [{
			"id": 4331493228753,
			"title": "Entrega Estandar CUNDINAMARCA",
			"price": "3.19",
			"price_set": {
				"shop_money": {"amount": "3.19", "currency_code": "USD"},
				"presentment_money": {"amount": "12000.00", "currency_code": "COP"}
			},
			"discounted_price": "3.19"
		}],
		"discount_applications": [
			{"type": "automatic", "title": "Descuento Miembro", "value": "10.101", "value_type": "percentage", "target_type": "line_item", "target_selection": "entitled", "allocation_method": "across"},
			{"type": "automatic", "title": "Descuento Miembro", "value": "10.101", "value_type": "percentage", "target_type": "line_item", "target_selection": "entitled", "allocation_method": "across"}
		],
		"fulfillments": []
	}`

	var order response.Order
	if err := json.Unmarshal([]byte(rawJSON), &order); err != nil {
		t.Fatalf("Error al parsear JSON de orden: %v", err)
	}

	businessID := uint(34)
	result := MapOrderResponseToShopifyOrder(order, []byte(rawJSON), &businessID, 1, "shopify")

	// Verificar datos generales de la orden
	if result.OrderNumber != "#67269" {
		t.Errorf("OrderNumber: esperado #67269, obtenido %s", result.OrderNumber)
	}
	if result.Currency != "USD" {
		t.Errorf("Currency: esperado USD, obtenido %s", result.Currency)
	}
	if result.CurrencyPresentment != "COP" {
		t.Errorf("CurrencyPresentment: esperado COP, obtenido %s", result.CurrencyPresentment)
	}

	// Verificar montos de orden
	if !almostEqual(result.TotalAmount, 44.99) {
		t.Errorf("TotalAmount: esperado 44.99, obtenido %.2f", result.TotalAmount)
	}
	if !almostEqual(result.TotalAmountPresentment, 169233.35) {
		t.Errorf("TotalAmountPresentment: esperado 169233.35, obtenido %.2f", result.TotalAmountPresentment)
	}
	if !almostEqual(result.Discount, 4.70) {
		t.Errorf("Order.Discount: esperado 4.70, obtenido %.2f", result.Discount)
	}
	if !almostEqual(result.DiscountPresentment, 17666.65) {
		t.Errorf("Order.DiscountPresentment: esperado 17666.65, obtenido %.2f", result.DiscountPresentment)
	}

	// Verificar items
	if len(result.Items) != 2 {
		t.Fatalf("Items count: esperado 2, obtenido %d", len(result.Items))
	}

	// Item 0: Colageno
	colageno := result.Items[0]
	if colageno.SKU != "PT01015" {
		t.Errorf("Item[0].SKU: esperado PT01015, obtenido %s", colageno.SKU)
	}
	if colageno.Quantity != 1 {
		t.Errorf("Item[0].Quantity: esperado 1, obtenido %d", colageno.Quantity)
	}
	if !almostEqual(colageno.UnitPrice, 13.16) {
		t.Errorf("Item[0].UnitPrice: esperado 13.16, obtenido %.2f", colageno.UnitPrice)
	}
	if !almostEqual(colageno.UnitPricePresentment, 49500.00) {
		t.Errorf("Item[0].UnitPricePresentment: esperado 49500.00, obtenido %.2f", colageno.UnitPricePresentment)
	}
	if !almostEqual(colageno.Discount, 1.33) {
		t.Errorf("Item[0].Discount: esperado 1.33 (de discount_allocations), obtenido %.2f", colageno.Discount)
	}
	if !almostEqual(colageno.DiscountPresentment, 4999.99) {
		t.Errorf("Item[0].DiscountPresentment: esperado 4999.99, obtenido %.2f", colageno.DiscountPresentment)
	}
	if !almostEqual(colageno.Tax, 1.88) {
		t.Errorf("Item[0].Tax: esperado 1.88, obtenido %.2f", colageno.Tax)
	}
	if !almostEqual(colageno.TaxPresentment, 7070.87) {
		t.Errorf("Item[0].TaxPresentment: esperado 7070.87, obtenido %.2f", colageno.TaxPresentment)
	}

	// Item 1: Creatina x2
	creatina := result.Items[1]
	if creatina.SKU != "PT01004" {
		t.Errorf("Item[1].SKU: esperado PT01004, obtenido %s", creatina.SKU)
	}
	if creatina.Quantity != 2 {
		t.Errorf("Item[1].Quantity: esperado 2, obtenido %d", creatina.Quantity)
	}
	if !almostEqual(creatina.UnitPrice, 16.67) {
		t.Errorf("Item[1].UnitPrice: esperado 16.67, obtenido %.2f", creatina.UnitPrice)
	}
	if !almostEqual(creatina.Discount, 3.37) {
		t.Errorf("Item[1].Discount: esperado 3.37 (de discount_allocations), obtenido %.2f", creatina.Discount)
	}
	if !almostEqual(creatina.DiscountPresentment, 12666.66) {
		t.Errorf("Item[1].DiscountPresentment: esperado 12666.66, obtenido %.2f", creatina.DiscountPresentment)
	}

	// Verificar suma de descuentos de items = descuento total de orden
	totalItemDiscountUSD := colageno.Discount + creatina.Discount
	if !almostEqual(totalItemDiscountUSD, 4.70) {
		t.Errorf("Suma descuentos items USD: esperado 4.70, obtenido %.2f", totalItemDiscountUSD)
	}
	totalItemDiscountCOP := colageno.DiscountPresentment + creatina.DiscountPresentment
	if !almostEqual(totalItemDiscountCOP, 17666.65) {
		t.Errorf("Suma descuentos items COP: esperado 17666.65, obtenido %.2f", totalItemDiscountCOP)
	}

	// Verificar customer
	if result.Customer.Name != "Carolina Ramirez" {
		t.Errorf("Customer.Name: esperado 'Carolina Ramirez', obtenido '%s'", result.Customer.Name)
	}
	if result.Customer.Phone != "+573178949340" {
		t.Errorf("Customer.Phone: esperado '+573178949340', obtenido '%s'", result.Customer.Phone)
	}
}
