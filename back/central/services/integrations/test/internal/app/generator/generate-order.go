package generator

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/test/internal/domain"
	"gorm.io/datatypes"
)

// GenerateRandomOrder genera una orden canónica aleatoria
// Si la plataforma es "shopify", genera primero el JSON de Shopify y luego lo mapea
func (g *OrderGenerator) GenerateRandomOrder(req *domain.GenerateOrderRequest) *domain.CanonicalOrderDTO {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Si es Shopify, generar JSON de Shopify y mapearlo
	if req.Platform == "shopify" {
		shopifyJSON, err := g.GenerateShopifyOrderJSON(r)
		if err != nil {
			// Si falla, usar el método genérico
			return g.generateGenericOrder(req, r)
		}
		canonicalOrder, err := g.MapShopifyJSONToCanonical(shopifyJSON, req.IntegrationID, req.BusinessID)
		if err != nil {
			// Si falla el mapeo, usar el método genérico
			return g.generateGenericOrder(req, r)
		}
		return canonicalOrder
	}

	// Si es Mercado Libre, generar JSON de Meli y mapearlo
	if req.Platform == "meli" || req.Platform == "mercado_libre" {
		meliJSON, err := g.GenerateMeliOrderJSON(r)
		if err != nil {
			// Si falla, usar el método genérico
			return g.generateGenericOrder(req, r)
		}
		canonicalOrder, err := g.MapMeliJSONToCanonical(meliJSON, req.IntegrationID, req.BusinessID)
		if err != nil {
			// Si falla el mapeo, usar el método genérico
			return g.generateGenericOrder(req, r)
		}
		return canonicalOrder
	}

	// Si es WooCommerce, generar JSON de WooCommerce y mapearlo
	if req.Platform == "woocommerce" || req.Platform == "woo" {
		wooJSON, err := g.GenerateWooCommerceOrderJSON(r)
		if err != nil {
			// Si falla, usar el método genérico
			return g.generateGenericOrder(req, r)
		}
		canonicalOrder, err := g.MapWooCommerceJSONToCanonical(wooJSON, req.IntegrationID, req.BusinessID)
		if err != nil {
			// Si falla el mapeo, usar el método genérico
			return g.generateGenericOrder(req, r)
		}
		return canonicalOrder
	}

	// Para otras plataformas, usar el método genérico
	return g.generateGenericOrder(req, r)
}

// generateGenericOrder genera una orden canónica genérica (método original)
func (g *OrderGenerator) generateGenericOrder(req *domain.GenerateOrderRequest, r *rand.Rand) *domain.CanonicalOrderDTO {

	// Seleccionar cliente aleatorio
	customer := domain.FakeCustomers[r.Intn(len(domain.FakeCustomers))]

	// Seleccionar dirección aleatoria
	address := domain.FakeAddresses[r.Intn(len(domain.FakeAddresses))]

	// Generar 1-5 items aleatorios
	numItems := r.Intn(5) + 1
	orderItems := make([]domain.CanonicalOrderItemDTO, numItems)
	subtotal := 0.0

	for i := 0; i < numItems; i++ {
		product := domain.FakeProducts[r.Intn(len(domain.FakeProducts))]
		quantity := r.Intn(3) + 1
		unitPrice := product.Price
		discount := 0.0
		if r.Float64() < 0.3 { // 30% de probabilidad de descuento
			discount = unitPrice * 0.1 * float64(r.Intn(3)+1) // 10%, 20% o 30%
		}
		totalPrice := (unitPrice - discount) * float64(quantity)
		tax := totalPrice * 0.19 // IVA 19%
		taxRate := 0.19

		orderItems[i] = domain.CanonicalOrderItemDTO{
			ProductID:    &product.ID,
			ProductSKU:   product.SKU,
			ProductName:  product.Name,
			ProductTitle: product.Title,
			Quantity:     quantity,
			UnitPrice:    unitPrice,
			TotalPrice:   totalPrice,
			Currency:     "CLP",
			Discount:     discount * float64(quantity),
			Tax:          tax,
			TaxRate:      &taxRate,
			ImageURL:     &product.ImageURL,
			ProductURL:   &product.ProductURL,
			Weight:       &product.Weight,
		}

		subtotal += totalPrice
	}

	// Calcular totales
	tax := subtotal * 0.19
	shippingCost := 0.0
	if r.Float64() < 0.7 { // 70% de probabilidad de tener envío
		shippingCost = float64(r.Intn(5000) + 3000) // Entre 3000 y 8000 CLP
	}
	totalAmount := subtotal + tax + shippingCost

	// Generar IDs
	externalID := fmt.Sprintf("TEST-%d-%d", time.Now().Unix(), r.Intn(10000))
	orderNumber := fmt.Sprintf("ORD-%d", r.Intn(999999)+100000)

	// Plataforma y tipo
	platform := req.Platform
	if platform == "" {
		platform = "test"
	}
	status := req.Status
	if status == "" {
		status = "pending"
	}

	// Crear orden canónica
	order := &domain.CanonicalOrderDTO{
		BusinessID:      req.BusinessID,
		IntegrationID:   req.IntegrationID,
		IntegrationType: "test",
		Platform:        platform,
		ExternalID:      externalID,
		OrderNumber:     orderNumber,
		Subtotal:        subtotal,
		Tax:             tax,
		Discount:        0,
		ShippingCost:    shippingCost,
		TotalAmount:     totalAmount,
		Currency:        "CLP",
		CustomerName:    customer.Name,
		CustomerEmail:   customer.Email,
		CustomerPhone:   customer.Phone,
		CustomerDNI:     customer.DNI,
		OrderTypeName:   "delivery",
		Status:          status,
		OriginalStatus:  status,
		OrderItems:      orderItems,
		OccurredAt:      time.Now(),
		ImportedAt:      time.Now(),
	}

	// Agregar direcciones
	order.Addresses = []domain.CanonicalAddressDTO{
		{
			Type:       "shipping",
			FirstName:  customer.Name,
			LastName:   "",
			Street:     address.Street,
			Street2:    address.Street2,
			City:       address.City,
			State:      address.State,
			Country:    address.Country,
			PostalCode: address.PostalCode,
			Latitude:   &address.Latitude,
			Longitude:  &address.Longitude,
		},
		{
			Type:       "billing",
			FirstName:  customer.Name,
			LastName:   "",
			Street:     address.Street,
			Street2:    address.Street2,
			City:       address.City,
			State:      address.State,
			Country:    address.Country,
			PostalCode: address.PostalCode,
		},
	}

	// Agregar pago si se solicita
	if req.IncludePayment {
		paymentMethod := domain.FakePaymentMethods[r.Intn(len(domain.FakePaymentMethods))]
		paidAt := time.Now()
		if r.Float64() < 0.7 { // 70% pagadas
			order.Payments = []domain.CanonicalPaymentDTO{
				{
					PaymentMethodID: paymentMethod.ID,
					Amount:          totalAmount,
					Currency:        "CLP",
					Status:          "completed",
					PaidAt:          &paidAt,
					TransactionID:   stringPtr(fmt.Sprintf("TXN-%d", r.Intn(999999)+100000)),
					Gateway:         stringPtr("test_gateway"),
				},
			}
		} else {
			order.Payments = []domain.CanonicalPaymentDTO{
				{
					PaymentMethodID: paymentMethod.ID,
					Amount:          totalAmount,
					Currency:        "CLP",
					Status:          "pending",
					Gateway:         stringPtr("test_gateway"),
				},
			}
		}
	}

	// Agregar envío si se solicita
	if req.IncludeShipment {
		carrier := domain.FakeCarriers[r.Intn(len(domain.FakeCarriers))]
		trackingNumber := fmt.Sprintf("TRK-%d", r.Intn(999999999)+100000000)
		trackingURL := fmt.Sprintf("https://tracking.example.com/%s", trackingNumber)
		shippedAt := time.Now().Add(-24 * time.Hour)
		estimatedDelivery := time.Now().Add(3 * 24 * time.Hour)

		// Calcular peso total
		totalWeight := 0.0
		for _, item := range orderItems {
			if item.Weight != nil {
				totalWeight += *item.Weight * float64(item.Quantity)
			}
		}

		order.Shipments = []domain.CanonicalShipmentDTO{
			{
				TrackingNumber:    &trackingNumber,
				TrackingURL:       &trackingURL,
				Carrier:           &carrier.Name,
				CarrierCode:       &carrier.Code,
				Status:            "in_transit",
				ShippedAt:         &shippedAt,
				EstimatedDelivery: &estimatedDelivery,
				ShippingCost:      &shippingCost,
				Weight:            &totalWeight,
				WarehouseName:     "Almacén Central",
			},
		}
	}

	// Agregar metadata del canal
	rawDataJSON, _ := datatypes.JSON(fmt.Sprintf(`{"test": true, "generated_at": "%s", "external_id": "%s"}`, time.Now().Format(time.RFC3339), externalID)).MarshalJSON()
	order.ChannelMetadata = &domain.CanonicalChannelMetadataDTO{
		ChannelSource: "test",
		RawData:       datatypes.JSON(rawDataJSON),
		Version:       "1.0",
		ReceivedAt:    time.Now(),
		IsLatest:      true,
		SyncStatus:    "pending",
	}

	return order
}

// stringPtr retorna un puntero a string
func stringPtr(s string) *string {
	return &s
}
