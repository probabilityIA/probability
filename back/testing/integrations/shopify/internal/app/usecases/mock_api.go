package usecases

import (
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"

	"github.com/secamc93/probability/back/testing/integrations/shopify/internal/domain"
	"github.com/secamc93/probability/back/testing/shared/log"
)

// MockAPIServer simula el API REST de Shopify (GET /admin/api/2024-10/orders.json).
// Pre-genera órdenes en memoria y las sirve con filtrado por fechas, estado y paginación via Link header.
type MockAPIServer struct {
	orders         []*domain.Order
	mu             sync.RWMutex
	dataGenerator  *RandomDataGenerator
	businessConfig *domain.BusinessConfig
	logger         log.ILogger
	orderSeq       int
}

// NewMockAPIServer crea un nuevo servidor mock de Shopify API.
func NewMockAPIServer(logger log.ILogger, businessConfig *domain.BusinessConfig) *MockAPIServer {
	return &MockAPIServer{
		orders:         make([]*domain.Order, 0),
		dataGenerator:  NewRandomDataGenerator(),
		businessConfig: businessConfig,
		logger:         logger,
		orderSeq:       1000,
	}
}

// GetBusinessConfig retorna la configuración del business para uso en handlers
func (m *MockAPIServer) GetBusinessConfig() *domain.BusinessConfig {
	return m.businessConfig
}

// GenerateOrders pre-genera N órdenes distribuidas en un rango de fechas.
// Las órdenes se crean con created_at distribuido uniformemente entre dateFrom y dateTo.
func (m *MockAPIServer) GenerateOrders(count int, dateFrom, dateTo time.Time) {
	m.mu.Lock()
	defer m.mu.Unlock()

	totalRange := dateTo.Sub(dateFrom)
	if totalRange <= 0 {
		totalRange = 24 * time.Hour
	}

	statuses := []string{"open", "closed", "cancelled"}
	financialStatuses := []string{"pending", "paid", "refunded", "partially_refunded"}
	fulfillmentStatuses := []*string{
		nil,
		strPtr("fulfilled"),
		strPtr("partial"),
	}

	for i := 0; i < count; i++ {
		m.orderSeq++
		orderNumber := fmt.Sprintf("#%d", m.orderSeq)

		// Distribuir created_at uniformemente en el rango
		offset := time.Duration(rand.Int63n(int64(totalRange)))
		createdAt := dateFrom.Add(offset)

		order := m.generateOrderAtDate(orderNumber, createdAt)

		// Asignar estados aleatorios
		statusIdx := rand.Intn(len(statuses))
		if statuses[statusIdx] == "cancelled" {
			reason := "customer"
			order.CancelReason = &reason
			order.CancelledAt = &createdAt
			order.FinancialStatus = "refunded"
		} else {
			order.FinancialStatus = financialStatuses[rand.Intn(len(financialStatuses))]
		}
		order.FulfillmentStatus = fulfillmentStatuses[rand.Intn(len(fulfillmentStatuses))]

		m.orders = append(m.orders, order)
	}

	// Ordenar por created_at descendente (como Shopify)
	sort.Slice(m.orders, func(i, j int) bool {
		return m.orders[i].CreatedAt.After(m.orders[j].CreatedAt)
	})

	m.logger.Info().
		Int("count", count).
		Str("from", dateFrom.Format(time.RFC3339)).
		Str("to", dateTo.Format(time.RFC3339)).
		Int("total_orders", len(m.orders)).
		Msg("📦 Órdenes generadas para mock Shopify API")
}

// QueryOrders filtra y pagina órdenes según los parámetros de Shopify.
// Retorna las órdenes de la página actual y si hay más páginas.
func (m *MockAPIServer) QueryOrders(params OrderQueryParams) (orders []*domain.Order, hasNextPage bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Filtrar
	filtered := make([]*domain.Order, 0)
	for _, o := range m.orders {
		if !m.matchesFilters(o, params) {
			continue
		}
		filtered = append(filtered, o)
	}

	// Paginación por since_id
	if params.SinceID > 0 {
		idx := -1
		for i, o := range filtered {
			if o.ID == params.SinceID {
				idx = i
				break
			}
		}
		if idx >= 0 && idx+1 < len(filtered) {
			filtered = filtered[idx+1:]
		} else if idx == -1 {
			// since_id no encontrado, retornar todo
		}
	}

	// Paginación por page_info (usamos offset simple)
	if params.PageOffset > 0 && params.PageOffset < len(filtered) {
		filtered = filtered[params.PageOffset:]
	}

	// Limitar
	limit := params.Limit
	if limit <= 0 || limit > 250 {
		limit = 250
	}

	if len(filtered) > limit {
		return filtered[:limit], true
	}

	return filtered, false
}

// GetTotalOrders retorna el total de órdenes en el mock.
func (m *MockAPIServer) GetTotalOrders() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.orders)
}

// OrderQueryParams representa los parámetros de consulta de órdenes de Shopify.
type OrderQueryParams struct {
	CreatedAtMin      *time.Time
	CreatedAtMax      *time.Time
	Status            string // any, open, closed, cancelled
	FinancialStatus   string
	FulfillmentStatus string
	Limit             int
	SinceID           int64
	PageOffset        int // Para paginación interna vía page_info
}

// matchesFilters verifica si una orden cumple con los filtros.
func (m *MockAPIServer) matchesFilters(order *domain.Order, params OrderQueryParams) bool {
	// Filtro por created_at_min
	if params.CreatedAtMin != nil && order.CreatedAt.Before(*params.CreatedAtMin) {
		return false
	}

	// Filtro por created_at_max
	if params.CreatedAtMax != nil && order.CreatedAt.After(*params.CreatedAtMax) {
		return false
	}

	// Filtro por status
	if params.Status != "" && params.Status != "any" {
		orderStatus := "open"
		if order.CancelledAt != nil {
			orderStatus = "cancelled"
		} else if order.ClosedAt != nil {
			orderStatus = "closed"
		}
		if orderStatus != params.Status {
			return false
		}
	}

	// Filtro por financial_status
	if params.FinancialStatus != "" && params.FinancialStatus != "any" {
		if order.FinancialStatus != params.FinancialStatus {
			return false
		}
	}

	// Filtro por fulfillment_status
	if params.FulfillmentStatus != "" && params.FulfillmentStatus != "any" {
		orderFS := ""
		if order.FulfillmentStatus != nil {
			orderFS = *order.FulfillmentStatus
		}
		if params.FulfillmentStatus == "unshipped" || params.FulfillmentStatus == "unfulfilled" {
			if orderFS != "" {
				return false
			}
		} else if orderFS != params.FulfillmentStatus {
			return false
		}
	}

	return true
}

// generateOrderAtDate genera una orden con un created_at específico.
// Delega a OrderSimulator para reutilizar la lógica dual-currency.
func (m *MockAPIServer) generateOrderAtDate(orderNumber string, createdAt time.Time) *domain.Order {
	if m.businessConfig.IsDualCurrency() {
		return m.generateDualCurrencyOrderAtDate(orderNumber, createdAt)
	}
	return m.generateSingleCurrencyOrderAtDate(orderNumber, createdAt)
}

// generateSingleCurrencyOrderAtDate genera una orden single-currency con created_at específico
func (m *MockAPIServer) generateSingleCurrencyOrderAtDate(orderNumber string, createdAt time.Time) *domain.Order {
	orderID := int64(rand.Intn(9999999999) + 1000000000)
	currency := "COP"

	customer := m.dataGenerator.GenerateCustomer()
	lineItems := m.dataGenerator.GenerateLineItems(rand.Intn(3) + 1)

	// Calcular totales
	subtotal := 0.0
	for _, item := range lineItems {
		var price float64
		fmt.Sscanf(item.Price, "%f", &price)
		subtotal += price * float64(item.Quantity)
	}

	taxRate := 0.19
	tax := subtotal * taxRate

	shippingLines := m.generateShippingLines(currency)
	var shippingCost float64
	if len(shippingLines) > 0 {
		fmt.Sscanf(shippingLines[0].Price, "%f", &shippingCost)
	}

	total := subtotal + tax + shippingCost

	subtotalStr := fmt.Sprintf("%.2f", subtotal)
	taxStr := fmt.Sprintf("%.2f", tax)
	totalStr := fmt.Sprintf("%.2f", total)
	shippingCostStr := fmt.Sprintf("%.2f", shippingCost)

	updatedAt := createdAt.Add(time.Duration(rand.Intn(3600)) * time.Second)

	order := &domain.Order{
		ID:                       orderID,
		AdminGraphQLAPIID:        fmt.Sprintf("gid://shopify/Order/%d", orderID),
		AppID:                    int64Ptr(int64(rand.Intn(999999) + 100000)),
		BrowserIP:                m.dataGenerator.stringPtr(fmt.Sprintf("%d.%d.%d.%d", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255))),
		BuyerAcceptsMarketing:    rand.Float32() < 0.3,
		CancelReason:             nil,
		CancelledAt:              nil,
		CartToken:                m.dataGenerator.stringPtr(fmt.Sprintf("%x", rand.Int63())),
		CheckoutID:               int64Ptr(int64(rand.Intn(9999999999) + 1000000000)),
		CheckoutToken:            m.dataGenerator.stringPtr(fmt.Sprintf("ct_%x", rand.Int63())),
		ClientDetails:            m.generateClientDetails(),
		ClosedAt:                 nil,
		Confirmed:                true,
		ContactEmail:             customer.Email,
		CreatedAt:                createdAt,
		Currency:                 currency,
		CurrentSubtotalPrice:     subtotalStr,
		CurrentSubtotalPriceSet:  m.generateMoneySet(subtotalStr, currency),
		CurrentTotalDiscounts:    "0.00",
		CurrentTotalDiscountsSet: m.generateMoneySet("0.00", currency),
		CurrentTotalDutiesSet:    nil,
		CurrentTotalPrice:        totalStr,
		CurrentTotalPriceSet:     m.generateMoneySet(totalStr, currency),
		CurrentTotalTax:          taxStr,
		CurrentTotalTaxSet:       m.generateMoneySet(taxStr, currency),
		CustomerLocale:           m.dataGenerator.stringPtr("es-CO"),
		DeviceID:                 nil,
		DiscountCodes:            []domain.DiscountCode{},
		Email:                    customer.Email,
		EstimatedTaxes:           false,
		FinancialStatus:          "pending",
		FulfillmentStatus:        nil,
		Gateway:                  m.dataGenerator.stringPtr("shopify_payments"),
		LandingSite:              m.dataGenerator.stringPtr("/"),
		LandingSiteRef:           nil,
		LocationID:               nil,
		MerchantOfRecordAppID:    nil,
		Name:                     orderNumber,
		Note:                     nil,
		NoteAttributes:           m.generateNoteAttributes(),
		Number:                   m.orderSeq,
		OrderNumber:              m.orderSeq,
		OrderStatusURL:           m.dataGenerator.stringPtr(fmt.Sprintf("https://%s/orders/%s/authenticate?key=abc", m.businessConfig.ShopDomain, orderNumber)),
		OriginalTotalDutiesSet:   nil,
		PaymentGatewayNames:      []string{"shopify_payments"},
		Phone:                    customer.Phone,
		PresentmentCurrency:      currency,
		ProcessedAt:              createdAt,
		ProcessingMethod:         m.dataGenerator.stringPtr("direct"),
		Reference:                nil,
		ReferringSite:            nil,
		SourceIdentifier:         nil,
		SourceName:               m.dataGenerator.randomChoice([]string{"web", "pos", "mobile", "api"}),
		SourceURL:                nil,
		SubtotalPrice:            subtotalStr,
		SubtotalPriceSet:         m.generateMoneySet(subtotalStr, currency),
		Tags:                     "",
		TaxLines:                 m.generateTaxLines(taxStr, currency, taxRate),
		TaxesIncluded:            false,
		Test:                     false,
		Token:                    fmt.Sprintf("ct_%x", rand.Int63()),
		TotalDiscounts:           "0.00",
		TotalDiscountsSet:        m.generateMoneySet("0.00", currency),
		TotalLineItemsPrice:      subtotalStr,
		TotalLineItemsPriceSet:   m.generateMoneySet(subtotalStr, currency),
		TotalOutstanding:         totalStr,
		TotalPrice:               totalStr,
		TotalPriceSet:            m.generateMoneySet(totalStr, currency),
		TotalPriceUSD:            fmt.Sprintf("%.2f", total*0.00025),
		TotalShippingPriceSet:    m.generateMoneySet(shippingCostStr, currency),
		TotalTax:                 taxStr,
		TotalTaxSet:              m.generateMoneySet(taxStr, currency),
		TotalTipReceived:         "0.00",
		TotalWeight:              rand.Intn(5000) + 100,
		UpdatedAt:                updatedAt,
		UserID:                   nil,
		BillingAddress:           m.dataGenerator.GenerateAddress(),
		Customer:                 customer,
		DiscountApplications:     []domain.DiscountApplication{},
		Fulfillments:             []domain.Fulfillment{},
		LineItems:                lineItems,
		PaymentTerms:             nil,
		Refunds:                  []domain.Refund{},
		ShippingAddress:          m.dataGenerator.GenerateAddress(),
		ShippingLines:            shippingLines,
	}

	return order
}

// generateDualCurrencyOrderAtDate genera una orden dual-currency USD/COP con created_at específico
func (m *MockAPIServer) generateDualCurrencyOrderAtDate(orderNumber string, createdAt time.Time) *domain.Order {
	orderID := int64(rand.Intn(9999999999) + 1000000000)
	exchangeRate := m.businessConfig.ExchangeRate

	customer := m.dataGenerator.GenerateCustomer()
	customer.Currency = "COP"

	lineItems := m.dataGenerator.GenerateDualCurrencyLineItems(rand.Intn(3)+1, exchangeRate)

	// Calcular totales en COP (precios incluyen IVA)
	totalCOP := 0.0
	totalTaxCOP := 0.0
	for _, item := range lineItems {
		var copPrice float64
		fmt.Sscanf(item.PriceSet.PresentmentMoney.Amount, "%f", &copPrice)
		totalCOP += copPrice * float64(item.Quantity)
		for _, tax := range item.TaxLines {
			var taxCOP float64
			fmt.Sscanf(tax.PriceSet.PresentmentMoney.Amount, "%f", &taxCOP)
			totalTaxCOP += taxCOP
		}
	}

	subtotalCOP := totalCOP - totalTaxCOP

	shippingLines := m.dataGenerator.GenerateDualCurrencyShippingLines(exchangeRate)
	var shippingCOP float64
	if len(shippingLines) > 0 {
		fmt.Sscanf(shippingLines[0].PriceSet.PresentmentMoney.Amount, "%f", &shippingCOP)
	}

	grandTotalCOP := totalCOP + shippingCOP

	subtotalUSD := subtotalCOP / exchangeRate
	taxUSD := totalTaxCOP / exchangeRate
	grandTotalUSD := grandTotalCOP / exchangeRate

	subtotalUSDStr := fmt.Sprintf("%.2f", subtotalUSD)
	taxUSDStr := fmt.Sprintf("%.2f", taxUSD)
	totalUSDStr := fmt.Sprintf("%.2f", grandTotalUSD)

	updatedAt := createdAt.Add(time.Duration(rand.Intn(3600)) * time.Second)

	billingAddress := m.dataGenerator.GenerateAddress()
	if rand.Float32() < 0.4 {
		dni := fmt.Sprintf("%d", rand.Intn(90000000)+10000000)
		billingAddress.Company = &dni
	}

	noteAttrs := m.generateNoteAttributes()
	noteAttrs = append(noteAttrs, domain.NoteAttribute{
		Name:  "_shipping_cost_cop",
		Value: fmt.Sprintf("%.0f", shippingCOP),
	})

	order := &domain.Order{
		ID:                    orderID,
		AdminGraphQLAPIID:     fmt.Sprintf("gid://shopify/Order/%d", orderID),
		AppID:                 int64Ptr(int64(rand.Intn(999999) + 100000)),
		BrowserIP:             m.dataGenerator.stringPtr(fmt.Sprintf("%d.%d.%d.%d", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255))),
		BuyerAcceptsMarketing: rand.Float32() < 0.3,
		CancelReason:          nil,
		CancelledAt:           nil,
		CartToken:             m.dataGenerator.stringPtr(fmt.Sprintf("%x", rand.Int63())),
		CheckoutID:            int64Ptr(int64(rand.Intn(9999999999) + 1000000000)),
		CheckoutToken:         m.dataGenerator.stringPtr(fmt.Sprintf("ct_%x", rand.Int63())),
		ClientDetails:         m.generateClientDetails(),
		ClosedAt:              nil,
		Confirmed:             true,
		ContactEmail:          customer.Email,
		CreatedAt:             createdAt,
		Currency:              "USD",
		CurrentSubtotalPrice:  subtotalUSDStr,
		CurrentSubtotalPriceSet:  m.dataGenerator.GenerateDualCurrencyMoneySet(subtotalCOP, exchangeRate),
		CurrentTotalDiscounts:    "0.00",
		CurrentTotalDiscountsSet: m.dataGenerator.GenerateDualCurrencyMoneySet(0, exchangeRate),
		CurrentTotalDutiesSet:    nil,
		CurrentTotalPrice:        totalUSDStr,
		CurrentTotalPriceSet:     m.dataGenerator.GenerateDualCurrencyMoneySet(grandTotalCOP, exchangeRate),
		CurrentTotalTax:          taxUSDStr,
		CurrentTotalTaxSet:       m.dataGenerator.GenerateDualCurrencyMoneySet(totalTaxCOP, exchangeRate),
		CustomerLocale:           m.dataGenerator.stringPtr("es-CO"),
		DeviceID:                 nil,
		DiscountCodes:            []domain.DiscountCode{},
		Email:                    customer.Email,
		EstimatedTaxes:           false,
		FinancialStatus:          "pending",
		FulfillmentStatus:        nil,
		Gateway:                  m.dataGenerator.stringPtr("shopify_payments"),
		LandingSite:              m.dataGenerator.stringPtr("/"),
		LandingSiteRef:           nil,
		LocationID:               nil,
		MerchantOfRecordAppID:    nil,
		Name:                     orderNumber,
		Note:                     nil,
		NoteAttributes:           noteAttrs,
		Number:                   m.orderSeq,
		OrderNumber:              m.orderSeq,
		OrderStatusURL:           m.dataGenerator.stringPtr(fmt.Sprintf("https://%s/orders/%s/authenticate?key=abc", m.businessConfig.ShopDomain, orderNumber)),
		OriginalTotalDutiesSet:   nil,
		PaymentGatewayNames:      []string{"shopify_payments"},
		Phone:                    customer.Phone,
		PresentmentCurrency:      "COP",
		ProcessedAt:              createdAt,
		ProcessingMethod:         m.dataGenerator.stringPtr("direct"),
		Reference:                nil,
		ReferringSite:            nil,
		SourceIdentifier:         nil,
		SourceName:               "web",
		SourceURL:                nil,
		SubtotalPrice:            subtotalUSDStr,
		SubtotalPriceSet:         m.dataGenerator.GenerateDualCurrencyMoneySet(subtotalCOP, exchangeRate),
		Tags:                     "",
		TaxLines: []domain.TaxLine{
			{
				Price:         taxUSDStr,
				Rate:          0.19,
				Title:         "IVA",
				PriceSet:      m.dataGenerator.GenerateDualCurrencyMoneySet(totalTaxCOP, exchangeRate),
				ChannelLiable: false,
			},
		},
		TaxesIncluded:          true,
		Test:                   false,
		Token:                  fmt.Sprintf("ct_%x", rand.Int63()),
		TotalDiscounts:         "0.00",
		TotalDiscountsSet:      m.dataGenerator.GenerateDualCurrencyMoneySet(0, exchangeRate),
		TotalLineItemsPrice:    fmt.Sprintf("%.2f", totalCOP/exchangeRate),
		TotalLineItemsPriceSet: m.dataGenerator.GenerateDualCurrencyMoneySet(totalCOP, exchangeRate),
		TotalOutstanding:       totalUSDStr,
		TotalPrice:             totalUSDStr,
		TotalPriceSet:          m.dataGenerator.GenerateDualCurrencyMoneySet(grandTotalCOP, exchangeRate),
		TotalPriceUSD:          totalUSDStr,
		TotalShippingPriceSet:  m.dataGenerator.GenerateDualCurrencyMoneySet(shippingCOP, exchangeRate),
		TotalTax:               taxUSDStr,
		TotalTaxSet:            m.dataGenerator.GenerateDualCurrencyMoneySet(totalTaxCOP, exchangeRate),
		TotalTipReceived:       "0.00",
		TotalWeight:            rand.Intn(5000) + 100,
		UpdatedAt:              updatedAt,
		UserID:                 nil,
		BillingAddress:         billingAddress,
		Customer:               customer,
		DiscountApplications:   []domain.DiscountApplication{},
		Fulfillments:           []domain.Fulfillment{},
		LineItems:              lineItems,
		PaymentTerms:           nil,
		Refunds:                []domain.Refund{},
		ShippingAddress:        m.dataGenerator.GenerateAddress(),
		ShippingLines:          shippingLines,
	}

	return order
}

func (m *MockAPIServer) generateNoteAttributes() []domain.NoteAttribute {
	return []domain.NoteAttribute{
		{Name: "_business_id", Value: fmt.Sprintf("%d", m.businessConfig.BusinessID)},
		{Name: "_business_code", Value: m.businessConfig.BusinessCode},
		{Name: "_integration_id", Value: fmt.Sprintf("%d", m.businessConfig.IntegrationID)},
		{Name: "_integration_code", Value: m.businessConfig.IntegrationCode},
		{Name: "_integration_type_id", Value: fmt.Sprintf("%d", m.businessConfig.IntegrationTypeID)},
		{Name: "_customer_dni", Value: fmt.Sprintf("%d", rand.Intn(90000000)+10000000)},
	}
}

func (m *MockAPIServer) generateClientDetails() *domain.ClientDetails {
	return &domain.ClientDetails{
		AcceptLanguage: m.dataGenerator.stringPtr("es"),
		BrowserIP:      m.dataGenerator.stringPtr(fmt.Sprintf("%d.%d.%d.%d", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255))),
		UserAgent:      m.dataGenerator.stringPtr("Mozilla/5.0 (Windows NT 10.0; Win64; x64)"),
	}
}

func (m *MockAPIServer) generateMoneySet(amount, currency string) *domain.MoneySet {
	return &domain.MoneySet{
		ShopMoney:        domain.Money{Amount: amount, CurrencyCode: currency},
		PresentmentMoney: domain.Money{Amount: amount, CurrencyCode: currency},
	}
}

func (m *MockAPIServer) generateTaxLines(taxAmount, currency string, rate float64) []domain.TaxLine {
	return []domain.TaxLine{
		{
			Price:         taxAmount,
			Rate:          rate,
			Title:         "IVA",
			PriceSet:      m.generateMoneySet(taxAmount, currency),
			ChannelLiable: false,
		},
	}
}

func (m *MockAPIServer) generateShippingLines(currency string) []domain.ShippingLine {
	methods := []struct {
		title string
		code  string
		price float64
	}{
		{"Entrega Estándar CUNDINAMARCA", "standard_cundinamarca", 3.15},
		{"Envío Express", "express", 5.00},
		{"Envío Gratis", "free_shipping", 0.00},
		{"Recogida en Tienda", "pickup", 0.00},
	}

	method := methods[rand.Intn(len(methods))]
	priceStr := fmt.Sprintf("%.2f", method.price)

	return []domain.ShippingLine{
		{
			ID:                 int64(rand.Intn(999999999) + 100000000),
			Title:              method.title,
			Code:               &method.code,
			Price:              priceStr,
			PriceSet:           m.generateMoneySet(priceStr, currency),
			DiscountedPrice:    priceStr,
			DiscountedPriceSet: m.generateMoneySet(priceStr, currency),
			Source:             m.dataGenerator.stringPtr("shopify"),
			TaxLines:           []domain.TaxLine{},
			DiscountAllocations: []domain.DiscountAllocation{},
		},
	}
}

func strPtr(s string) *string {
	return &s
}
