package app

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/secamc93/probability/back/testing/modules/orders/internal/domain/dtos"
	"github.com/secamc93/probability/back/testing/modules/orders/internal/domain/entities"
	sharedtypes "github.com/secamc93/probability/back/testing/shared/types"
)

var (
	firstNames = []string{"Juan", "Maria", "Carlos", "Ana", "Pedro", "Laura", "Diego", "Sofia", "Andres", "Valentina"}
	lastNames  = []string{"Garcia", "Rodriguez", "Martinez", "Lopez", "Gonzalez", "Hernandez", "Perez", "Sanchez", "Ramirez", "Torres"}
	cities     = []string{"Bogota", "Medellin", "Cali", "Barranquilla", "Cartagena", "Bucaramanga", "Pereira", "Manizales", "Cucuta", "Ibague"}
	states     = []string{"Cundinamarca", "Antioquia", "Valle del Cauca", "Atlantico", "Bolivar", "Santander", "Risaralda", "Caldas", "Norte de Santander", "Tolima"}
	streets    = []string{"Calle 10 #5-20", "Carrera 15 #30-45", "Avenida 80 #12-33", "Calle 50 #25-10", "Carrera 7 #40-60", "Diagonal 35 #8-15"}
)

const (
	categoryEcommerce = 1
	categoryPlatform  = 6
)

func (uc *useCase) GenerateOrders(ctx context.Context, businessID uint, dto *dtos.GenerateOrdersDTO, token string) (*entities.GenerateResult, error) {
	dto.ApplyDefaults()

	if dto.IntegrationID == 0 {
		return nil, fmt.Errorf("integration_id is required")
	}

	// Determine if this is a platform or ecommerce integration
	categoryID, err := uc.repo.GetIntegrationCategoryID(ctx, dto.IntegrationID)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve integration category: %w", err)
	}

	switch categoryID {
	case categoryPlatform:
		return uc.buildPlatformPayloads(ctx, businessID, dto, token)
	case categoryEcommerce:
		return uc.buildEcommercePayloads(ctx, dto)
	default:
		return nil, fmt.Errorf("unsupported integration category: %d", categoryID)
	}
}

// buildPlatformPayloads builds payloads for native/platform order creation
func (uc *useCase) buildPlatformPayloads(ctx context.Context, businessID uint, dto *dtos.GenerateOrdersDTO, token string) (*entities.GenerateResult, error) {
	products, err := uc.repo.GetProducts(ctx, businessID)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}
	if len(products) == 0 {
		return nil, fmt.Errorf("no products found for business %d", businessID)
	}

	paymentMethods, err := uc.repo.GetPaymentMethods(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment methods: %w", err)
	}
	if len(paymentMethods) == 0 {
		return nil, fmt.Errorf("no payment methods found")
	}

	centralURL := uc.centralClient.GetBaseURL() + "/api/v1/orders"
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	result := &entities.GenerateResult{
		Total: dto.Count,
	}

	for i := 0; i < dto.Count; i++ {
		orderBody := buildRandomOrder(rng, businessID, dto.IntegrationID, products, paymentMethods, dto)

		payload := sharedtypes.WebhookPayload{
			URL:    centralURL,
			Method: "POST",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
			Body: orderBody,
		}
		result.Payloads = append(result.Payloads, payload)
	}

	uc.log.Info().
		Uint("business_id", businessID).
		Int("total", result.Total).
		Msg("Platform payloads built")

	return result, nil
}

// buildEcommercePayloads builds payloads using the webhook simulator for the integration type
func (uc *useCase) buildEcommercePayloads(ctx context.Context, dto *dtos.GenerateOrdersDTO) (*entities.GenerateResult, error) {
	typeCode, err := uc.repo.GetIntegrationTypeCode(ctx, dto.IntegrationID)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve integration type: %w", err)
	}

	simulator, ok := uc.webhookSimulators[typeCode]
	if !ok {
		return nil, fmt.Errorf("webhook simulator not available for integration type: %s", typeCode)
	}

	topic := dto.Topic
	if topic == "" {
		topic = "orders/create"
	}

	result := &entities.GenerateResult{
		Total: dto.Count,
	}

	for i := 0; i < dto.Count; i++ {
		payload, err := simulator.BuildWebhookPayload(topic, uc.centralClient.GetBaseURL())
		if err != nil {
			uc.log.Warn().
				Int("index", i).
				Err(err).
				Str("type_code", typeCode).
				Str("topic", topic).
				Msg("Failed to build webhook payload")
			result.Errors = append(result.Errors, entities.OrderError{
				Index:   i,
				Message: err.Error(),
			})
			continue
		}
		result.Payloads = append(result.Payloads, *payload)
	}

	uc.log.Info().
		Str("type_code", typeCode).
		Str("topic", topic).
		Int("total", result.Total).
		Int("built", len(result.Payloads)).
		Msg("Ecommerce webhook payloads built")

	return result, nil
}

func buildRandomOrder(rng *rand.Rand, businessID, integrationID uint, products []entities.Product, paymentMethods []entities.PaymentMethod, dto *dtos.GenerateOrdersDTO) map[string]interface{} {
	firstName := firstNames[rng.Intn(len(firstNames))]
	lastName := lastNames[rng.Intn(len(lastNames))]
	if dto.CustomerName != "" {
		parts := splitName(dto.CustomerName)
		firstName = parts[0]
		lastName = parts[1]
	}
	cityIdx := rng.Intn(len(cities))

	numItems := rng.Intn(dto.MaxItemsPerOrder) + 1
	if numItems > len(products) {
		numItems = len(products)
	}

	perm := rng.Perm(len(products))
	var items []map[string]interface{}
	var subtotal float64

	for j := 0; j < numItems; j++ {
		p := products[perm[j]]
		qty := rng.Intn(3) + 1
		price := p.Price
		if price <= 0 {
			price = float64(rng.Intn(200000) + 10000)
		}
		itemTotal := price * float64(qty)
		subtotal += itemTotal

		items = append(items, map[string]interface{}{
			"sku":      p.SKU,
			"name":     p.Name,
			"quantity": qty,
			"price":    price,
		})
	}

	tax := subtotal * 0.19
	shippingCost := float64(rng.Intn(15000) + 5000)
	total := subtotal + tax + shippingCost

	pm := paymentMethods[rng.Intn(len(paymentMethods))]

	email := fmt.Sprintf("%s.%s%d@test.probability.com", firstName, lastName, rng.Intn(100))
	phone := fmt.Sprintf("+5730%08d", rng.Intn(100000000))
	if dto.CustomerPhone != "" {
		phone = dto.CustomerPhone
	}
	dni := fmt.Sprintf("%d", 1000000000+rng.Intn(999999999))

	timestamp := time.Now().UnixNano()
	externalID := fmt.Sprintf("test-%d-%d", timestamp, rng.Intn(10000))

	now := time.Now().Format(time.RFC3339)

	return map[string]interface{}{
		"business_id":           businessID,
		"integration_id":        integrationID,
		"integration_type":      "platform",
		"platform":              "manual",
		"external_id":           externalID,
		"order_number":          "AUTO",
		"subtotal":              subtotal,
		"tax":                   tax,
		"discount":              0,
		"shipping_cost":         shippingCost,
		"total_amount":          total,
		"currency":              "COP",
		"customer_name":         fmt.Sprintf("%s %s", firstName, lastName),
		"customer_first_name":   firstName,
		"customer_last_name":    lastName,
		"customer_email":        email,
		"customer_phone":        phone,
		"customer_dni":          dni,
		"shipping_street":       streets[rng.Intn(len(streets))],
		"shipping_city":         cities[cityIdx],
		"shipping_state":        states[cityIdx],
		"shipping_country":      "Colombia",
		"shipping_postal_code":  fmt.Sprintf("0%d0001", cityIdx+1),
		"payment_method_id":     pm.ID,
		"is_paid":               rng.Intn(2) == 1,
		"status":                "pending",
		"invoiceable":           true,
		"items":                 buildItemsJSON(items),
		"occurred_at":           now,
		"imported_at":           now,
	}
}

func splitName(fullName string) [2]string {
	parts := strings.SplitN(strings.TrimSpace(fullName), " ", 2)
	if len(parts) == 1 {
		return [2]string{parts[0], ""}
	}
	return [2]string{parts[0], parts[1]}
}

func buildItemsJSON(items []map[string]interface{}) json.RawMessage {
	data, err := json.Marshal(items)
	if err != nil {
		return json.RawMessage("[]")
	}
	return json.RawMessage(data)
}
