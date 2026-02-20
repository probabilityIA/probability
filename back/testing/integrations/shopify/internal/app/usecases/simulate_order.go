package usecases

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/secamc93/probability/back/testing/integrations/shopify/internal/domain"
	"github.com/secamc93/probability/back/testing/shared/env"
	"github.com/secamc93/probability/back/testing/shared/log"
)

// OrderSimulator simula órdenes de Shopify
type OrderSimulator struct {
	webhookClient   domain.IWebhookClient
	config          env.IConfig
	logger          log.ILogger
	orderRepository *domain.OrderRepository
	dataGenerator   *RandomDataGenerator
	orderNumberSeq  int
	businessConfig  *domain.BusinessConfig // Configuración del business de prueba
}

// SimulateOrder simula una orden según el topic especificado
func (s *OrderSimulator) SimulateOrder(topic string) error {
	// Usar el shop domain del business config
	shopDomain := s.businessConfig.ShopDomain
	if shopDomain == "" {
		return fmt.Errorf("Shop domain no configurado en business config")
	}

	var order *domain.Order
	var err error

	switch topic {
	case "orders/create":
		order, err = s.CreateRandomOrder()
		if err != nil {
			return err
		}
	case "orders/paid":
		order, err = s.GetOrCreateRandomOrder()
		if err != nil {
			return err
		}
		order = s.markAsPaid(order)
	case "orders/updated":
		order, err = s.GetOrCreateRandomOrder()
		if err != nil {
			return err
		}
		order = s.updateOrder(order)
	case "orders/cancelled":
		order, err = s.GetOrCreateRandomOrder()
		if err != nil {
			return err
		}
		order = s.cancelOrder(order)
	case "orders/fulfilled":
		order, err = s.GetOrCreateRandomOrder()
		if err != nil {
			return err
		}
		order = s.fulfillOrder(order)
	case "orders/partially_fulfilled":
		order, err = s.GetOrCreateRandomOrder()
		if err != nil {
			return err
		}
		order = s.partiallyFulfillOrder(order)
	default:
		return fmt.Errorf("topic no soportado: %s", topic)
	}

	// Guardar la orden actualizada en el repositorio
	s.orderRepository.Save(order)

	return s.webhookClient.SendWebhook(topic, shopDomain, *order)
}

// CreateRandomOrder crea una nueva orden aleatoria
func (s *OrderSimulator) CreateRandomOrder() (*domain.Order, error) {
	orderNumber := s.generateUniqueOrderNumber()
	order := s.generateRandomOrder(orderNumber)
	s.orderRepository.Save(order)
	s.logger.Info().
		Str("order_number", orderNumber).
		Str("business", s.businessConfig.BusinessName).
		Uint("business_id", s.businessConfig.BusinessID).
		Uint("integration_id", s.businessConfig.IntegrationID).
		Str("shop_domain", s.businessConfig.ShopDomain).
		Msg("✅ Orden creada con configuración de business real")
	return order, nil
}

// GetOrCreateRandomOrder obtiene una orden existente o crea una nueva
func (s *OrderSimulator) GetOrCreateRandomOrder() (*domain.Order, error) {
	allOrders := s.orderRepository.GetAll()
	if len(allOrders) > 0 && rand.Float32() < 0.7 {
		// 70% de probabilidad de usar una orden existente
		selected := allOrders[rand.Intn(len(allOrders))]
		s.logger.Info().Str("order_number", selected.Name).Msg("Usando orden existente")
		return selected, nil
	}
	// Crear una nueva orden
	return s.CreateRandomOrder()
}

// GetAllOrders retorna todas las órdenes almacenadas
func (s *OrderSimulator) GetAllOrders() []*domain.Order {
	return s.orderRepository.GetAll()
}

// GetOrderByNumber obtiene una orden por su número
func (s *OrderSimulator) GetOrderByNumber(orderNumber string) (*domain.Order, bool) {
	return s.orderRepository.Get(orderNumber)
}

// generateUniqueOrderNumber genera un número de orden único
func (s *OrderSimulator) generateUniqueOrderNumber() string {
	for {
		s.orderNumberSeq++
		orderNumber := fmt.Sprintf("#%d", s.orderNumberSeq)
		if !s.orderRepository.Exists(orderNumber) {
			return orderNumber
		}
	}
}

// generateRandomOrder genera una orden completamente aleatoria
func (s *OrderSimulator) generateRandomOrder(orderNumber string) *domain.Order {
	now := time.Now()
	orderID := int64(rand.Intn(9999999999) + 1000000000)
	currency := s.dataGenerator.randomChoice([]string{"COP", "USD", "EUR"})

	customer := s.dataGenerator.GenerateCustomer()
	lineItems := s.dataGenerator.GenerateLineItems(rand.Intn(3) + 1)

	// Calcular totales
	subtotal := 0.0
	for _, item := range lineItems {
		var price float64
		fmt.Sscanf(item.Price, "%f", &price)
		subtotal += price * float64(item.Quantity)
	}

	taxRate := 0.19
	tax := subtotal * taxRate

	// Generar shipping lines para calcular el costo de envío
	shippingLines := s.generateShippingLines(currency)
	var shippingCost float64
	if len(shippingLines) > 0 {
		fmt.Sscanf(shippingLines[0].Price, "%f", &shippingCost)
	}

	total := subtotal + tax + shippingCost

	subtotalStr := fmt.Sprintf("%.2f", subtotal)
	taxStr := fmt.Sprintf("%.2f", tax)
	totalStr := fmt.Sprintf("%.2f", total)
	shippingCostStr := fmt.Sprintf("%.2f", shippingCost)

	order := &domain.Order{
		ID:                       orderID,
		AdminGraphQLAPIID:        fmt.Sprintf("gid://shopify/Order/%d", orderID),
		AppID:                    int64Ptr(int64(rand.Intn(999999) + 100000)),
		BrowserIP:                s.dataGenerator.stringPtr(fmt.Sprintf("%d.%d.%d.%d", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255))),
		BuyerAcceptsMarketing:    rand.Float32() < 0.3,
		CancelReason:             nil,
		CancelledAt:              nil,
		CartToken:                s.dataGenerator.stringPtr(fmt.Sprintf("%x", rand.Int63())),
		CheckoutID:               int64Ptr(int64(rand.Intn(9999999999) + 1000000000)),
		CheckoutToken:            s.dataGenerator.stringPtr(fmt.Sprintf("ct_%x", rand.Int63())),
		ClientDetails:            s.generateClientDetails(),
		ClosedAt:                 nil,
		Confirmed:                true,
		ContactEmail:             customer.Email,
		CreatedAt:                now.Add(-time.Duration(rand.Intn(30)) * time.Minute),
		Currency:                 currency,
		CurrentSubtotalPrice:     subtotalStr,
		CurrentSubtotalPriceSet:  s.generateMoneySet(subtotalStr, currency),
		CurrentTotalDiscounts:    "0.00",
		CurrentTotalDiscountsSet: s.generateMoneySet("0.00", currency),
		CurrentTotalDutiesSet:    nil,
		CurrentTotalPrice:        totalStr,
		CurrentTotalPriceSet:     s.generateMoneySet(totalStr, currency),
		CurrentTotalTax:          taxStr,
		CurrentTotalTaxSet:       s.generateMoneySet(taxStr, currency),
		CustomerLocale:           s.dataGenerator.stringPtr("es-CO"),
		DeviceID:                 nil,
		DiscountCodes:            []domain.DiscountCode{},
		Email:                    customer.Email,
		EstimatedTaxes:           false,
		FinancialStatus:          "pending",
		FulfillmentStatus:        nil,
		Gateway:                  s.dataGenerator.stringPtr("shopify_payments"),
		LandingSite:              s.dataGenerator.stringPtr("/"),
		LandingSiteRef:           nil,
		LocationID:               nil,
		MerchantOfRecordAppID:    nil,
		Name:                     orderNumber,
		Note:                     nil,
		NoteAttributes:           s.generateNoteAttributes(),
		Number:                   s.orderNumberSeq,
		OrderNumber:              s.orderNumberSeq,
		OrderStatusURL:           s.dataGenerator.stringPtr(fmt.Sprintf("https://%s/orders/%s/authenticate?key=abc", s.businessConfig.ShopDomain, orderNumber)),
		OriginalTotalDutiesSet:   nil,
		PaymentGatewayNames:      []string{"shopify_payments"},
		Phone:                    customer.Phone,
		PresentmentCurrency:      currency,
		ProcessedAt:              now.Add(-time.Duration(rand.Intn(30)) * time.Minute),
		ProcessingMethod:         s.dataGenerator.stringPtr("direct"),
		Reference:                nil,
		ReferringSite:            nil,
		SourceIdentifier:         nil,
		SourceName:               s.dataGenerator.randomChoice([]string{"web", "pos", "mobile", "api"}),
		SourceURL:                nil,
		SubtotalPrice:            subtotalStr,
		SubtotalPriceSet:         s.generateMoneySet(subtotalStr, currency),
		Tags:                     "",
		TaxLines:                 s.generateTaxLines(taxStr, currency, taxRate),
		TaxesIncluded:            false,
		Test:                     false,
		Token:                    fmt.Sprintf("ct_%x", rand.Int63()),
		TotalDiscounts:           "0.00",
		TotalDiscountsSet:        s.generateMoneySet("0.00", currency),
		TotalLineItemsPrice:      subtotalStr,
		TotalLineItemsPriceSet:   s.generateMoneySet(subtotalStr, currency),
		TotalOutstanding:         totalStr,
		TotalPrice:               totalStr,
		TotalPriceSet:            s.generateMoneySet(totalStr, currency),
		TotalPriceUSD:            fmt.Sprintf("%.2f", total*0.00025), // Aproximación
		TotalShippingPriceSet:    s.generateMoneySet(shippingCostStr, currency),
		TotalTax:                 taxStr,
		TotalTaxSet:              s.generateMoneySet(taxStr, currency),
		TotalTipReceived:         "0.00",
		TotalWeight:              rand.Intn(5000) + 100,
		UpdatedAt:                now,
		UserID:                   nil,
		BillingAddress:           s.dataGenerator.GenerateAddress(),
		Customer:                 customer,
		DiscountApplications:     []domain.DiscountApplication{},
		Fulfillments:             []domain.Fulfillment{},
		LineItems:                lineItems,
		PaymentTerms:             nil,
		Refunds:                  []domain.Refund{},
		ShippingAddress:          s.dataGenerator.GenerateAddress(),
		ShippingLines:            shippingLines,
	}

	return order
}

// markAsPaid marca una orden como pagada
func (s *OrderSimulator) markAsPaid(order *domain.Order) *domain.Order {
	order.FinancialStatus = "paid"
	order.ProcessedAt = time.Now()
	order.UpdatedAt = time.Now()
	return order
}

// updateOrder actualiza una orden existente
func (s *OrderSimulator) updateOrder(order *domain.Order) *domain.Order {
	order.UpdatedAt = time.Now()
	// Simular cambios aleatorios
	if rand.Float32() < 0.5 {
		order.Tags = s.dataGenerator.randomChoice([]string{"urgente", "vip", "regalo", "envio_express", ""})
	}
	if rand.Float32() < 0.3 {
		note := fmt.Sprintf("Nota actualizada: %d", rand.Intn(9999))
		order.Note = &note
	}
	return order
}

// cancelOrder cancela una orden
func (s *OrderSimulator) cancelOrder(order *domain.Order) *domain.Order {
	order.FinancialStatus = "refunded"
	reasons := []string{"customer", "fraud", "inventory", "other"}
	reason := s.dataGenerator.randomChoice(reasons)
	order.CancelReason = &reason
	now := time.Now()
	order.CancelledAt = &now
	order.ClosedAt = &now
	order.UpdatedAt = now
	return order
}

// fulfillOrder marca una orden como cumplida
func (s *OrderSimulator) fulfillOrder(order *domain.Order) *domain.Order {
	order.FulfillmentStatus = s.dataGenerator.stringPtr("fulfilled")
	order.Fulfillments = s.generateFulfillments(order)
	order.UpdatedAt = time.Now()
	return order
}

// partiallyFulfillOrder marca una orden como parcialmente cumplida
func (s *OrderSimulator) partiallyFulfillOrder(order *domain.Order) *domain.Order {
	order.FulfillmentStatus = s.dataGenerator.stringPtr("partial")
	order.Fulfillments = s.generateFulfillments(order)
	order.UpdatedAt = time.Now()
	return order
}

// Funciones auxiliares
func (s *OrderSimulator) generateClientDetails() *domain.ClientDetails {
	return &domain.ClientDetails{
		AcceptLanguage: s.dataGenerator.stringPtr("es"),
		BrowserHeight:  nil,
		BrowserIP:      s.dataGenerator.stringPtr(fmt.Sprintf("%d.%d.%d.%d", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255))),
		BrowserWidth:   nil,
		SessionHash:    nil,
		UserAgent: s.dataGenerator.stringPtr(s.dataGenerator.randomChoice([]string{
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)",
			"Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X)",
			"Mozilla/5.0 (Android; Mobile)",
		})),
	}
}

func (s *OrderSimulator) generateMoneySet(amount, currency string) *domain.MoneySet {
	return &domain.MoneySet{
		ShopMoney: domain.Money{
			Amount:       amount,
			CurrencyCode: currency,
		},
		PresentmentMoney: domain.Money{
			Amount:       amount,
			CurrencyCode: currency,
		},
	}
}

func (s *OrderSimulator) generateTaxLines(taxAmount, currency string, rate float64) []domain.TaxLine {
	return []domain.TaxLine{
		{
			Price:         taxAmount,
			Rate:          rate,
			Title:         "IVA",
			PriceSet:      s.generateMoneySet(taxAmount, currency),
			ChannelLiable: false,
		},
	}
}

func (s *OrderSimulator) generateFulfillments(order *domain.Order) []domain.Fulfillment {
	now := time.Now()
	trackingNumber := fmt.Sprintf("%d", rand.Intn(9999999999)+1000000000)
	shipmentStatus := "confirmed"

	return []domain.Fulfillment{
		{
			ID:                int64(rand.Intn(999999999) + 100000000),
			OrderID:           order.ID,
			Status:            "success",
			CreatedAt:         now,
			Service:           s.dataGenerator.stringPtr("manual"),
			UpdatedAt:         now,
			TrackingCompany:   s.dataGenerator.stringPtr(s.dataGenerator.randomChoice([]string{"FedEx", "DHL", "Servientrega", "Coordinadora", "TCC"})),
			ShipmentStatus:    &shipmentStatus,
			LocationID:        int64Ptr(int64(rand.Intn(99999999) + 10000000)),
			OriginAddress:     nil,
			Receipt:           nil,
			Name:              fmt.Sprintf("%s.1", order.Name),
			AdminGraphQLAPIID: fmt.Sprintf("gid://shopify/Fulfillment/%d", rand.Intn(999999999)+100000000),
			TrackingNumbers:   []string{trackingNumber},
			TrackingUrls:      []string{fmt.Sprintf("https://tracking.example.com/%s", trackingNumber)},
			TrackingNumber:    s.dataGenerator.stringPtr(trackingNumber),
			TrackingURL:       s.dataGenerator.stringPtr(fmt.Sprintf("https://tracking.example.com/%s", trackingNumber)),
			UpdatedAtCustom:   nil,
			LineItems:         order.LineItems, // Incluir los line items en el fulfillment
		},
	}
}

// generateShippingLines genera líneas de envío para la orden
func (s *OrderSimulator) generateShippingLines(currency string) []domain.ShippingLine {
	shippingMethods := []struct {
		title string
		code  string
		price float64
	}{
		{"Entrega Estándar CUNDINAMARCA (3 a 6 días hábiles municipios principales en Cundinamarca - 3 o más días a otros municipios)", "standard_cundinamarca", 3.15},
		{"Envío Express", "express", 5.00},
		{"Envío Gratis", "free_shipping", 0.00},
		{"Recogida en Tienda", "pickup", 0.00},
	}

	method := shippingMethods[rand.Intn(len(shippingMethods))]
	priceStr := fmt.Sprintf("%.2f", method.price)

	return []domain.ShippingLine{
		{
			ID:                            int64(rand.Intn(999999999) + 100000000),
			Title:                         method.title,
			Code:                          &method.code,
			Price:                         priceStr,
			PriceSet:                      s.generateMoneySet(priceStr, currency),
			DiscountedPrice:               priceStr,
			DiscountedPriceSet:            s.generateMoneySet(priceStr, currency),
			Source:                        s.dataGenerator.stringPtr("shopify"),
			CarrierIdentifier:             nil,
			DeliveryCategory:              nil,
			Phone:                         nil,
			RequestedFulfillmentServiceID: nil,
			TaxLines:                      []domain.TaxLine{},
			DiscountAllocations:           []domain.DiscountAllocation{},
		},
	}
}

// generateNoteAttributes genera atributos de nota con metadatos del business
// Estos atributos ayudan al backend a identificar de qué business e integración proviene la orden
func (s *OrderSimulator) generateNoteAttributes() []domain.NoteAttribute {
	return []domain.NoteAttribute{
		{
			Name:  "_business_id",
			Value: fmt.Sprintf("%d", s.businessConfig.BusinessID),
		},
		{
			Name:  "_business_code",
			Value: s.businessConfig.BusinessCode,
		},
		{
			Name:  "_integration_id",
			Value: fmt.Sprintf("%d", s.businessConfig.IntegrationID),
		},
		{
			Name:  "_integration_code",
			Value: s.businessConfig.IntegrationCode,
		},
		{
			Name:  "_integration_type_id",
			Value: fmt.Sprintf("%d", s.businessConfig.IntegrationTypeID),
		},
		{
			Name:  "_customer_dni",
			Value: "1098720627",
		},
	}
}

// Helpers
func int64Ptr(i int64) *int64 {
	return &i
}
