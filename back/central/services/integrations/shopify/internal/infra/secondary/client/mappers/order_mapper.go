package mappers

import (
	"strconv"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/shopify/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/shopify/internal/infra/secondary/client/response"
)

// MapOrderResponseToShopifyOrder mapea una Order de respuesta de Shopify a ShopifyOrder del dominio
func MapOrderResponseToShopifyOrder(orderResp response.Order, rawOrder []byte, businessID *uint, integrationID uint, integrationType string) domain.ShopifyOrder {
	// Convertir precios de string a float64 (shop_money - USD)
	totalAmount, _ := strconv.ParseFloat(orderResp.TotalPrice, 64)

	// Extraer precios en presentment_money (moneda local) si están disponibles
	var totalAmountPresentment float64
	var currencyPresentment string
	var subtotalPresentment, taxPresentment, discountPresentment, shippingCostPresentment float64

	if orderResp.TotalPriceSet != nil && orderResp.TotalPriceSet.PresentmentMoney.Amount != "" {
		totalAmountPresentment, _ = strconv.ParseFloat(orderResp.TotalPriceSet.PresentmentMoney.Amount, 64)
		currencyPresentment = orderResp.TotalPriceSet.PresentmentMoney.CurrencyCode
	}

	if orderResp.PresentmentCurrency != "" {
		currencyPresentment = orderResp.PresentmentCurrency
	}

	// Extraer subtotal en moneda local
	if orderResp.SubtotalPriceSet != nil && orderResp.SubtotalPriceSet.PresentmentMoney.Amount != "" {
		subtotalPresentment, _ = strconv.ParseFloat(orderResp.SubtotalPriceSet.PresentmentMoney.Amount, 64)
	}

	// Extraer tax en moneda local
	if orderResp.TotalTaxSet != nil && orderResp.TotalTaxSet.PresentmentMoney.Amount != "" {
		taxPresentment, _ = strconv.ParseFloat(orderResp.TotalTaxSet.PresentmentMoney.Amount, 64)
	}

	// Extraer discount en moneda local
	if orderResp.TotalDiscountsSet != nil && orderResp.TotalDiscountsSet.PresentmentMoney.Amount != "" {
		discountPresentment, _ = strconv.ParseFloat(orderResp.TotalDiscountsSet.PresentmentMoney.Amount, 64)
	}

	// Extraer shipping cost en moneda local (sumar todos los shipping lines)
	if orderResp.TotalShippingPriceSet != nil && orderResp.TotalShippingPriceSet.PresentmentMoney.Amount != "" {
		shippingCostPresentment, _ = strconv.ParseFloat(orderResp.TotalShippingPriceSet.PresentmentMoney.Amount, 64)
	}

	// Mapear customer
	customer := domain.ShopifyCustomer{
		Email: orderResp.Email,
		Phone: "",
	}
	if orderResp.Phone != nil {
		customer.Phone = *orderResp.Phone
	}
	if orderResp.Customer != nil {
		customer.Name = orderResp.Customer.FirstName + " " + orderResp.Customer.LastName
		customer.Email = orderResp.Customer.Email
		if orderResp.Customer.Phone != nil {
			customer.Phone = *orderResp.Customer.Phone
		}
	} else {
		// Si no hay customer, usar email de la orden
		customer.Name = orderResp.Email
		customer.Email = orderResp.Email
	}

	// Mapear shipping address
	shippingAddress := domain.ShopifyAddress{
		Street:     "",
		Address2:   "",
		City:       "",
		State:      "",
		Country:    "",
		PostalCode: "",
	}
	if orderResp.ShippingAddress != nil {
		shippingAddress.Street = orderResp.ShippingAddress.Address1
		if orderResp.ShippingAddress.Address2 != nil {
			shippingAddress.Address2 = *orderResp.ShippingAddress.Address2
		}
		shippingAddress.City = orderResp.ShippingAddress.City
		shippingAddress.State = orderResp.ShippingAddress.Province
		shippingAddress.Country = orderResp.ShippingAddress.Country
		shippingAddress.PostalCode = orderResp.ShippingAddress.Zip

		// Mapear coordenadas si existen
		if orderResp.ShippingAddress.Latitude != nil && orderResp.ShippingAddress.Longitude != nil {
			shippingAddress.Coordinates = &struct {
				Lat float64
				Lng float64
			}{
				Lat: *orderResp.ShippingAddress.Latitude,
				Lng: *orderResp.ShippingAddress.Longitude,
			}
		}
	}

	// Mapear items
	items := make([]domain.ShopifyOrderItem, len(orderResp.LineItems))
	for i, item := range orderResp.LineItems {
		unitPrice, _ := strconv.ParseFloat(item.Price, 64)
		totalDiscount, _ := strconv.ParseFloat(item.TotalDiscount, 64)

		// Calcular impuesto total de tax_lines
		var totalTax float64
		for _, taxLine := range item.TaxLines {
			taxPrice, _ := strconv.ParseFloat(taxLine.Price, 64)
			totalTax += taxPrice
		}

		// Extraer precios en moneda local del item (presentment_money)
		var unitPricePresentment, discountPresentment, taxPresentment float64
		if item.PriceSet != nil && item.PriceSet.PresentmentMoney.Amount != "" {
			unitPricePresentment, _ = strconv.ParseFloat(item.PriceSet.PresentmentMoney.Amount, 64)
		}
		if item.TotalDiscountSet != nil && item.TotalDiscountSet.PresentmentMoney.Amount != "" {
			discountPresentment, _ = strconv.ParseFloat(item.TotalDiscountSet.PresentmentMoney.Amount, 64)
		}
		// Calcular tax en moneda local sumando las tax_lines
		for _, taxLine := range item.TaxLines {
			if taxLine.PriceSet != nil && taxLine.PriceSet.PresentmentMoney.Amount != "" {
				taxPricePresentment, _ := strconv.ParseFloat(taxLine.PriceSet.PresentmentMoney.Amount, 64)
				taxPresentment += taxPricePresentment
			}
		}

		// Convertir gramos a float64 para peso
		var weight *float64
		if item.Grams > 0 {
			weightVal := float64(item.Grams) / 1000.0 // Convertir a kg
			weight = &weightVal
		}

		productID := item.ProductID
		variantID := item.VariantID

		items[i] = domain.ShopifyOrderItem{
			ExternalID:   strconv.FormatInt(item.VariantID, 10),
			Name:         item.Name,
			SKU:          item.SKU,
			Quantity:     item.Quantity,
			UnitPrice:    unitPrice,
			ProductID:    &productID,
			VariantID:    &variantID,
			Title:        item.Title,
			VariantTitle: item.VariantTitle,
			Discount:     totalDiscount,
			Tax:          totalTax,
			Weight:       weight,
			// Precios en moneda local
			UnitPricePresentment: unitPricePresentment,
			DiscountPresentment:  discountPresentment,
			TaxPresentment:       taxPresentment,
		}
	}

	// Determinar status
	status := orderResp.FinancialStatus
	if orderResp.FulfillmentStatus != nil {
		status = *orderResp.FulfillmentStatus
	}

	// Mapear metadata
	metadata := make(map[string]interface{})
	metadata["shopify_id"] = orderResp.ID
	metadata["shopify_name"] = orderResp.Name
	metadata["shopify_token"] = orderResp.Token
	if orderResp.Note != nil {
		metadata["note"] = *orderResp.Note
	}
	metadata["tags"] = orderResp.Tags
	metadata["source_name"] = orderResp.SourceName
	metadata["payment_gateway_names"] = orderResp.PaymentGatewayNames

	// Determinar OrderStatusURL
	orderStatusURL := ""
	if orderResp.OrderStatusURL != nil {
		orderStatusURL = *orderResp.OrderStatusURL
	}

	// Usar el Name de Shopify (ej: "#1001") como OrderNumber, o el OrderNumber numérico si Name está vacío
	orderNumber := orderResp.Name
	if orderNumber == "" {
		orderNumber = strconv.Itoa(orderResp.OrderNumber)
	}

	return domain.ShopifyOrder{
		BusinessID:      businessID,
		IntegrationID:   integrationID,
		IntegrationType: integrationType,
		Platform:        "shopify",
		ExternalID:      strconv.FormatInt(orderResp.ID, 10),
		OrderNumber:     orderNumber,
		TotalAmount:     totalAmount,
		Currency:        orderResp.Currency,
		Customer:        customer,
		ShippingAddress: shippingAddress,
		Status:          status,
		OriginalStatus:  orderResp.FinancialStatus,
		Items:           items,
		Metadata:        metadata,
		OccurredAt:      orderResp.CreatedAt,
		ImportedAt:      time.Now(),
		OrderStatusURL:  orderStatusURL,
		RawData:         rawOrder,
		// Precios en moneda local
		SubtotalPresentment:     subtotalPresentment,
		TaxPresentment:          taxPresentment,
		DiscountPresentment:     discountPresentment,
		ShippingCostPresentment: shippingCostPresentment,
		TotalAmountPresentment:  totalAmountPresentment,
		CurrencyPresentment:     currencyPresentment,
	}
}

// MapOrdersResponseToShopifyOrders mapea múltiples órdenes de respuesta a ShopifyOrder del dominio
func MapOrdersResponseToShopifyOrders(ordersResp []response.Order, rawOrders [][]byte, businessID *uint, integrationID uint, integrationType string) []domain.ShopifyOrder {
	orders := make([]domain.ShopifyOrder, len(ordersResp))
	for i, orderResp := range ordersResp {
		var rawOrder []byte
		if i < len(rawOrders) {
			rawOrder = rawOrders[i]
		}
		orders[i] = MapOrderResponseToShopifyOrder(orderResp, rawOrder, businessID, integrationID, integrationType)
	}
	return orders
}

// MapOrderResponseToDomain mapea una Order de respuesta de Shopify a map[string]interface{} (legacy)
// DEPRECATED: Usar MapOrderResponseToShopifyOrder en su lugar
func MapOrderResponseToDomain(orderResp response.Order) map[string]interface{} {
	// Convertir la estructura tipada a map[string]interface{} para mantener compatibilidad
	// mientras se migra el código existente
	orderMap := make(map[string]interface{})

	orderMap["id"] = orderResp.ID
	orderMap["name"] = orderResp.Name
	orderMap["order_number"] = orderResp.OrderNumber
	orderMap["email"] = orderResp.Email
	if orderResp.Phone != nil {
		orderMap["phone"] = *orderResp.Phone
	} else {
		orderMap["phone"] = ""
	}
	orderMap["created_at"] = orderResp.CreatedAt.Format(time.RFC3339)
	orderMap["updated_at"] = orderResp.UpdatedAt.Format(time.RFC3339)
	orderMap["processed_at"] = orderResp.ProcessedAt.Format(time.RFC3339)
	orderMap["currency"] = orderResp.Currency
	orderMap["total_price"] = orderResp.TotalPrice
	orderMap["subtotal_price"] = orderResp.SubtotalPrice
	orderMap["total_tax"] = orderResp.TotalTax
	orderMap["total_discounts"] = orderResp.TotalDiscounts
	orderMap["financial_status"] = orderResp.FinancialStatus
	orderMap["fulfillment_status"] = orderResp.FulfillmentStatus
	orderMap["source_name"] = orderResp.SourceName
	orderMap["payment_gateway_names"] = orderResp.PaymentGatewayNames
	orderMap["tags"] = orderResp.Tags
	orderMap["total_weight"] = orderResp.TotalWeight
	if orderResp.Note != nil {
		orderMap["note"] = *orderResp.Note
	}
	if orderResp.LocationID != nil {
		orderMap["location_id"] = *orderResp.LocationID
	}

	// Mapear customer
	if orderResp.Customer != nil {
		orderMap["customer"] = mapCustomerToDomain(orderResp.Customer)
	}

	// Mapear shipping_address
	if orderResp.ShippingAddress != nil {
		orderMap["shipping_address"] = mapAddressToDomain(orderResp.ShippingAddress)
	}

	// Mapear billing_address
	if orderResp.BillingAddress != nil {
		orderMap["billing_address"] = mapAddressToDomain(orderResp.BillingAddress)
	}

	// Mapear line_items
	lineItems := make([]map[string]interface{}, len(orderResp.LineItems))
	for i, item := range orderResp.LineItems {
		lineItems[i] = mapLineItemToDomain(item)
	}
	orderMap["line_items"] = lineItems

	// Mapear shipping_lines
	shippingLines := make([]map[string]interface{}, len(orderResp.ShippingLines))
	for i, line := range orderResp.ShippingLines {
		shippingLines[i] = mapShippingLineToDomain(line)
	}
	orderMap["shipping_lines"] = shippingLines

	// Mapear fulfillments
	fulfillments := make([]map[string]interface{}, len(orderResp.Fulfillments))
	for i, fulfillment := range orderResp.Fulfillments {
		fulfillments[i] = mapFulfillmentToDomain(fulfillment)
	}
	orderMap["fulfillments"] = fulfillments

	return orderMap
}

// MapOrdersResponseToDomain mapea múltiples órdenes de respuesta a maps del dominio
func MapOrdersResponseToDomain(ordersResp []response.Order) []map[string]interface{} {
	orders := make([]map[string]interface{}, len(ordersResp))
	for i, orderResp := range ordersResp {
		orders[i] = MapOrderResponseToDomain(orderResp)
	}
	return orders
}

// mapCustomerToDomain mapea un Customer de respuesta a map del dominio
func mapCustomerToDomain(customer *response.Customer) map[string]interface{} {
	customerMap := make(map[string]interface{})
	customerMap["id"] = customer.ID
	customerMap["email"] = customer.Email
	customerMap["first_name"] = customer.FirstName
	customerMap["last_name"] = customer.LastName
	if customer.Phone != nil {
		customerMap["phone"] = *customer.Phone
	}
	customerMap["verified_email"] = customer.VerifiedEmail
	customerMap["created_at"] = customer.CreatedAt.Format(time.RFC3339)
	customerMap["updated_at"] = customer.UpdatedAt.Format(time.RFC3339)
	customerMap["state"] = customer.State
	if customer.Note != nil {
		customerMap["note"] = *customer.Note
	}
	customerMap["tags"] = customer.Tags
	customerMap["currency"] = customer.Currency
	customerMap["tax_exempt"] = customer.TaxExempt
	customerMap["accepts_marketing"] = customer.AcceptsMarketing

	if customer.DefaultAddress != nil {
		customerMap["default_address"] = mapAddressToDomain(customer.DefaultAddress)
	}

	return customerMap
}

// mapAddressToDomain mapea una Address de respuesta a map del dominio
func mapAddressToDomain(address *response.Address) map[string]interface{} {
	addressMap := make(map[string]interface{})
	addressMap["first_name"] = address.FirstName
	addressMap["last_name"] = address.LastName
	if address.Company != nil {
		addressMap["company"] = *address.Company
	}
	addressMap["address1"] = address.Address1
	if address.Address2 != nil {
		addressMap["address2"] = *address.Address2
	}
	addressMap["city"] = address.City
	addressMap["province"] = address.Province
	addressMap["country"] = address.Country
	addressMap["zip"] = address.Zip
	addressMap["country_code"] = address.CountryCode
	if address.Phone != nil {
		addressMap["phone"] = *address.Phone
	}
	if address.ProvinceCode != nil {
		addressMap["province_code"] = *address.ProvinceCode
	}
	if address.Latitude != nil {
		addressMap["latitude"] = *address.Latitude
	}
	if address.Longitude != nil {
		addressMap["longitude"] = *address.Longitude
	}
	addressMap["name"] = address.Name

	return addressMap
}

// mapLineItemToDomain mapea un LineItem de respuesta a map del dominio
func mapLineItemToDomain(item response.LineItem) map[string]interface{} {
	itemMap := make(map[string]interface{})
	itemMap["id"] = item.ID
	itemMap["product_id"] = item.ProductID
	itemMap["variant_id"] = item.VariantID
	itemMap["title"] = item.Title
	if item.VariantTitle != nil {
		itemMap["variant_title"] = *item.VariantTitle
	}
	itemMap["sku"] = item.SKU
	itemMap["quantity"] = item.Quantity
	itemMap["price"] = item.Price
	itemMap["grams"] = item.Grams
	itemMap["total_discount"] = item.TotalDiscount
	if item.FulfillmentStatus != nil {
		itemMap["fulfillment_status"] = *item.FulfillmentStatus
	}
	itemMap["name"] = item.Name
	itemMap["taxable"] = item.Taxable
	itemMap["requires_shipping"] = item.RequiresShipping
	itemMap["gift_card"] = item.GiftCard
	itemMap["fulfillable_quantity"] = item.FulfillableQuantity
	if item.Vendor != nil {
		itemMap["vendor"] = *item.Vendor
	}

	// Mapear tax_lines
	if len(item.TaxLines) > 0 {
		taxLines := make([]map[string]interface{}, len(item.TaxLines))
		for i, taxLine := range item.TaxLines {
			taxLines[i] = map[string]interface{}{
				"title": taxLine.Title,
				"price": taxLine.Price,
				"rate":  taxLine.Rate,
			}
		}
		itemMap["tax_lines"] = taxLines
	}

	return itemMap
}

// mapShippingLineToDomain mapea una ShippingLine de respuesta a map del dominio
func mapShippingLineToDomain(line response.ShippingLine) map[string]interface{} {
	lineMap := make(map[string]interface{})
	lineMap["id"] = line.ID
	lineMap["title"] = line.Title
	lineMap["price"] = line.Price
	if line.Code != nil {
		lineMap["code"] = *line.Code
	}
	if line.Source != nil {
		lineMap["source"] = *line.Source
	}
	if line.Phone != nil {
		lineMap["phone"] = *line.Phone
	}
	if line.CarrierIdentifier != nil {
		lineMap["carrier_identifier"] = *line.CarrierIdentifier
	}
	if line.DeliveryCategory != nil {
		lineMap["delivery_category"] = *line.DeliveryCategory
	}
	lineMap["discounted_price"] = line.DiscountedPrice

	return lineMap
}

// mapFulfillmentToDomain mapea un Fulfillment de respuesta a map del dominio
func mapFulfillmentToDomain(fulfillment response.Fulfillment) map[string]interface{} {
	fulfillmentMap := make(map[string]interface{})
	fulfillmentMap["id"] = fulfillment.ID
	fulfillmentMap["order_id"] = fulfillment.OrderID
	fulfillmentMap["status"] = fulfillment.Status
	fulfillmentMap["created_at"] = fulfillment.CreatedAt.Format(time.RFC3339)
	fulfillmentMap["updated_at"] = fulfillment.UpdatedAt.Format(time.RFC3339)
	if fulfillment.Service != nil {
		fulfillmentMap["service"] = *fulfillment.Service
	}
	if fulfillment.TrackingCompany != nil {
		fulfillmentMap["tracking_company"] = *fulfillment.TrackingCompany
	}
	if fulfillment.ShipmentStatus != nil {
		fulfillmentMap["shipment_status"] = *fulfillment.ShipmentStatus
	}
	if fulfillment.TrackingNumber != nil {
		fulfillmentMap["tracking_number"] = *fulfillment.TrackingNumber
	}
	if fulfillment.TrackingURL != nil {
		fulfillmentMap["tracking_url"] = *fulfillment.TrackingURL
	}
	if len(fulfillment.TrackingNumbers) > 0 {
		fulfillmentMap["tracking_numbers"] = fulfillment.TrackingNumbers
	}
	if len(fulfillment.TrackingUrls) > 0 {
		fulfillmentMap["tracking_urls"] = fulfillment.TrackingUrls
	}
	if fulfillment.LocationID != nil {
		fulfillmentMap["location_id"] = *fulfillment.LocationID
	}
	fulfillmentMap["name"] = fulfillment.Name

	return fulfillmentMap
}

// MapToShopifyAPIOrder mapea una Order de respuesta a ShopifyAPIOrder del dominio
func MapToShopifyAPIOrder(orderResp response.Order) domain.ShopifyAPIOrder {
	order := domain.ShopifyAPIOrder{
		ID:                  orderResp.ID,
		Name:                orderResp.Name,
		OrderNumber:         orderResp.OrderNumber,
		Email:               orderResp.Email,
		CreatedAt:           orderResp.CreatedAt.Format(time.RFC3339),
		UpdatedAt:           orderResp.UpdatedAt.Format(time.RFC3339),
		ProcessedAt:         orderResp.ProcessedAt.Format(time.RFC3339),
		Currency:            orderResp.Currency,
		TotalPrice:          orderResp.TotalPrice,
		SubtotalPrice:       orderResp.SubtotalPrice,
		TotalTax:            orderResp.TotalTax,
		TotalDiscounts:      orderResp.TotalDiscounts,
		FinancialStatus:     orderResp.FinancialStatus,
		FulfillmentStatus:   orderResp.FulfillmentStatus,
		SourceName:          orderResp.SourceName,
		PaymentGatewayNames: orderResp.PaymentGatewayNames,
		Tags:                orderResp.Tags,
		TotalWeight:         orderResp.TotalWeight,
		Note:                orderResp.Note,
		LocationID:          orderResp.LocationID,
	}

	if orderResp.Phone != nil {
		order.Phone = *orderResp.Phone
	}

	// Mapear customer
	if orderResp.Customer != nil {
		order.Customer = &domain.ShopifyAPICustomer{
			ID:            orderResp.Customer.ID,
			Email:         orderResp.Customer.Email,
			FirstName:     orderResp.Customer.FirstName,
			LastName:      orderResp.Customer.LastName,
			Phone:         orderResp.Customer.Phone,
			VerifiedEmail: orderResp.Customer.VerifiedEmail,
			State:         orderResp.Customer.State,
			CreatedAt:     orderResp.Customer.CreatedAt.Format(time.RFC3339),
			UpdatedAt:     orderResp.Customer.UpdatedAt.Format(time.RFC3339),
		}
	}

	// Mapear shipping_address
	if orderResp.ShippingAddress != nil {
		order.ShippingAddress = &domain.ShopifyAPIAddress{
			FirstName:   orderResp.ShippingAddress.FirstName,
			LastName:    orderResp.ShippingAddress.LastName,
			Company:     orderResp.ShippingAddress.Company,
			Address1:    orderResp.ShippingAddress.Address1,
			Address2:    orderResp.ShippingAddress.Address2,
			City:        orderResp.ShippingAddress.City,
			Province:    orderResp.ShippingAddress.Province,
			Country:     orderResp.ShippingAddress.Country,
			Zip:         orderResp.ShippingAddress.Zip,
			Phone:       orderResp.ShippingAddress.Phone,
			CountryCode: orderResp.ShippingAddress.CountryCode,
			Latitude:    orderResp.ShippingAddress.Latitude,
			Longitude:   orderResp.ShippingAddress.Longitude,
		}
		if orderResp.ShippingAddress.ProvinceCode != nil {
			order.ShippingAddress.ProvinceCode = *orderResp.ShippingAddress.ProvinceCode
		}
	}

	// Mapear billing_address
	if orderResp.BillingAddress != nil {
		order.BillingAddress = &domain.ShopifyAPIAddress{
			FirstName:   orderResp.BillingAddress.FirstName,
			LastName:    orderResp.BillingAddress.LastName,
			Company:     orderResp.BillingAddress.Company,
			Address1:    orderResp.BillingAddress.Address1,
			Address2:    orderResp.BillingAddress.Address2,
			City:        orderResp.BillingAddress.City,
			Province:    orderResp.BillingAddress.Province,
			Country:     orderResp.BillingAddress.Country,
			Zip:         orderResp.BillingAddress.Zip,
			Phone:       orderResp.BillingAddress.Phone,
			CountryCode: orderResp.BillingAddress.CountryCode,
			Latitude:    orderResp.BillingAddress.Latitude,
			Longitude:   orderResp.BillingAddress.Longitude,
		}
		if orderResp.BillingAddress.ProvinceCode != nil {
			order.BillingAddress.ProvinceCode = *orderResp.BillingAddress.ProvinceCode
		}
	}

	// Mapear line_items
	order.LineItems = make([]domain.ShopifyLineItem, len(orderResp.LineItems))
	for i, item := range orderResp.LineItems {
		order.LineItems[i] = domain.ShopifyLineItem{
			ID:                item.ID,
			ProductID:         &item.ProductID,
			VariantID:         &item.VariantID,
			Title:             item.Title,
			VariantTitle:      item.VariantTitle,
			SKU:               item.SKU,
			Quantity:          item.Quantity,
			Price:             item.Price,
			Grams:             item.Grams,
			TotalDiscount:     item.TotalDiscount,
			FulfillmentStatus: item.FulfillmentStatus,
			Name:              item.Name,
		}
	}

	// Mapear shipping_lines
	order.ShippingLines = make([]domain.ShopifyShippingLine, len(orderResp.ShippingLines))
	for i, line := range orderResp.ShippingLines {
		code := ""
		if line.Code != nil {
			code = *line.Code
		}
		source := ""
		if line.Source != nil {
			source = *line.Source
		}
		order.ShippingLines[i] = domain.ShopifyShippingLine{
			ID:                line.ID,
			Title:             line.Title,
			Price:             line.Price,
			Code:              code,
			Source:            source,
			Phone:             line.Phone,
			CarrierIdentifier: line.CarrierIdentifier,
			DeliveryCategory:  line.DeliveryCategory,
		}
		if line.RequestedFulfillmentServiceID != nil {
			serviceIDStr := strconv.FormatInt(*line.RequestedFulfillmentServiceID, 10)
			order.ShippingLines[i].RequestedFulfillmentServiceID = &serviceIDStr
		}
	}

	// Mapear fulfillments
	order.Fulfillments = make([]domain.ShopifyFulfillment, len(orderResp.Fulfillments))
	for i, fulfillment := range orderResp.Fulfillments {
		order.Fulfillments[i] = domain.ShopifyFulfillment{
			ID:              fulfillment.ID,
			OrderID:         fulfillment.OrderID,
			Status:          fulfillment.Status,
			CreatedAt:       fulfillment.CreatedAt.Format(time.RFC3339),
			UpdatedAt:       fulfillment.UpdatedAt.Format(time.RFC3339),
			TrackingCompany: fulfillment.TrackingCompany,
			TrackingNumber:  fulfillment.TrackingNumber,
			TrackingNumbers: fulfillment.TrackingNumbers,
			TrackingURL:     fulfillment.TrackingURL,
			TrackingURLs:    fulfillment.TrackingUrls,
			ShipmentStatus:  fulfillment.ShipmentStatus,
		}
	}

	return order
}
