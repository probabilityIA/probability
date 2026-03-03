package app

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/secamc93/probability/back/testing/internal/modules/orders/internal/domain/dtos"
	"github.com/secamc93/probability/back/testing/internal/modules/orders/internal/domain/entities"
)

var (
	firstNames = []string{"Juan", "Maria", "Carlos", "Ana", "Pedro", "Laura", "Diego", "Sofia", "Andres", "Valentina"}
	lastNames  = []string{"Garcia", "Rodriguez", "Martinez", "Lopez", "Gonzalez", "Hernandez", "Perez", "Sanchez", "Ramirez", "Torres"}
	cities     = []string{"Bogota", "Medellin", "Cali", "Barranquilla", "Cartagena", "Bucaramanga", "Pereira", "Manizales", "Cucuta", "Ibague"}
	states     = []string{"Cundinamarca", "Antioquia", "Valle del Cauca", "Atlantico", "Bolivar", "Santander", "Risaralda", "Caldas", "Norte de Santander", "Tolima"}
	streets    = []string{"Calle 10 #5-20", "Carrera 15 #30-45", "Avenida 80 #12-33", "Calle 50 #25-10", "Carrera 7 #40-60", "Diagonal 35 #8-15"}
)

func (uc *useCase) GenerateOrders(ctx context.Context, businessID uint, dto *dtos.GenerateOrdersDTO, token string) (*entities.GenerateResult, error) {
	dto.ApplyDefaults()

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

	// Resolve integration ID
	integrationID := dto.IntegrationID
	if integrationID == 0 {
		integrations, err := uc.repo.GetIntegrations(ctx, businessID)
		if err != nil {
			return nil, fmt.Errorf("failed to get integrations: %w", err)
		}
		// Prefer platform integration
		for _, intg := range integrations {
			if intg.Category == "platform" {
				integrationID = intg.ID
				break
			}
		}
		if integrationID == 0 && len(integrations) > 0 {
			integrationID = integrations[0].ID
		}
		if integrationID == 0 {
			return nil, fmt.Errorf("no integrations found for business %d", businessID)
		}
	}

	result := &entities.GenerateResult{
		Total: dto.Count,
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < dto.Count; i++ {
		order := buildRandomOrder(rng, businessID, integrationID, products, paymentMethods, dto)

		created, apiLog, err := uc.centralClient.CreateOrder(ctx, token, order)
		if apiLog != nil {
			apiLog.Index = i
			result.APILogs = append(result.APILogs, *apiLog)
		}
		if err != nil {
			uc.log.Warn().
				Int("index", i).
				Err(err).
				Msg("Failed to create test order")
			result.Failed++
			result.Errors = append(result.Errors, entities.OrderError{
				Index:   i,
				Message: err.Error(),
			})
			continue
		}

		result.Created++
		result.Orders = append(result.Orders, *created)
	}

	uc.log.Info().
		Uint("business_id", businessID).
		Int("total", result.Total).
		Int("created", result.Created).
		Int("failed", result.Failed).
		Msg("Order generation completed")

	return result, nil
}

func buildRandomOrder(rng *rand.Rand, businessID, integrationID uint, products []entities.Product, paymentMethods []entities.PaymentMethod, dto *dtos.GenerateOrdersDTO) map[string]interface{} {
	// Pick random customer
	firstName := firstNames[rng.Intn(len(firstNames))]
	lastName := lastNames[rng.Intn(len(lastNames))]
	cityIdx := rng.Intn(len(cities))

	// Pick random products
	numItems := rng.Intn(dto.MaxItemsPerOrder) + 1
	if numItems > len(products) {
		numItems = len(products)
	}

	// Shuffle and pick products
	perm := rng.Perm(len(products))
	var items []map[string]interface{}
	var subtotal float64

	for j := 0; j < numItems; j++ {
		p := products[perm[j]]
		qty := rng.Intn(3) + 1
		price := p.Price
		if price <= 0 {
			price = float64(rng.Intn(200000)+10000) // 10k-210k COP
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
	shippingCost := float64(rng.Intn(15000) + 5000) // 5k-20k
	total := subtotal + tax + shippingCost

	// Pick random payment method
	pm := paymentMethods[rng.Intn(len(paymentMethods))]

	// Random email
	email := fmt.Sprintf("%s.%s%d@test.probability.com", firstName, lastName, rng.Intn(100))
	phone := fmt.Sprintf("+5730%08d", rng.Intn(100000000))
	dni := fmt.Sprintf("%d", 1000000000+rng.Intn(999999999))

	timestamp := time.Now().UnixNano()
	externalID := fmt.Sprintf("test-%d-%d", timestamp, rng.Intn(10000))

	now := time.Now().Format(time.RFC3339)

	return map[string]interface{}{
		"business_id":      businessID,
		"integration_id":   integrationID,
		"integration_type": "platform",
		"platform":         "manual",
		"external_id":      externalID,
		"order_number":     "AUTO",
		"subtotal":         subtotal,
		"tax":              tax,
		"discount":         0,
		"shipping_cost":    shippingCost,
		"total_amount":     total,
		"currency":         "COP",
		"customer_name":       fmt.Sprintf("%s %s", firstName, lastName),
		"customer_first_name": firstName,
		"customer_last_name":  lastName,
		"customer_email":      email,
		"customer_phone":      phone,
		"customer_dni":        dni,
		"shipping_street":      streets[rng.Intn(len(streets))],
		"shipping_city":        cities[cityIdx],
		"shipping_state":       states[cityIdx],
		"shipping_country":     "Colombia",
		"shipping_postal_code": fmt.Sprintf("0%d0001", cityIdx+1),
		"payment_method_id": pm.ID,
		"is_paid":           rng.Intn(2) == 1,
		"status":            "pending",
		"invoiceable":       true,
		"items":             buildItemsJSON(items),
		"occurred_at":       now,
		"imported_at":       now,
	}
}

func buildItemsJSON(items []map[string]interface{}) json.RawMessage {
	data, err := json.Marshal(items)
	if err != nil {
		return json.RawMessage("[]")
	}
	return json.RawMessage(data)
}
