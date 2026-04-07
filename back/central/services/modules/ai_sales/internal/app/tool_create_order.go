package app

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	domain "github.com/secamc93/probability/back/central/services/modules/ai_sales/internal/domain"
)

type createOrderInput struct {
	CustomerName  string            `json:"customer_name"`
	CustomerPhone string            `json:"customer_phone"`
	Items         []createOrderItem `json:"items"`
}

type createOrderItem struct {
	ProductSKU string `json:"product_sku"`
	Quantity   int    `json:"quantity"`
}

// executeCreateOrder valida productos, construye la orden canonica y la publica
func executeCreateOrder(ctx context.Context, inputJSON string, deps *toolDeps) (string, error) {
	var input createOrderInput
	if err := parseToolInput(inputJSON, &input); err != nil {
		return "", fmt.Errorf("error parsing CreateOrder input: %w", err)
	}

	if len(input.Items) == 0 {
		return `{"error": "No se proporcionaron productos para el pedido"}`, nil
	}

	// Validar cada producto y calcular totales
	type validatedItem struct {
		Product  *domain.ProductSearchResult
		Quantity int
	}

	var items []validatedItem
	var totalAmount float64
	currency := "COP"

	for _, item := range input.Items {
		if item.Quantity <= 0 {
			return fmt.Sprintf(`{"error": "Cantidad invalida para producto %s"}`, item.ProductSKU), nil
		}

		product, err := deps.productRepo.GetProductBySKU(ctx, deps.businessID, item.ProductSKU)
		if err != nil {
			return fmt.Sprintf(`{"error": "Producto no encontrado: %s"}`, item.ProductSKU), nil
		}

		// Solo validar stock si el producto trackea inventario
		if product.TrackInventory && product.StockQuantity < item.Quantity {
			return fmt.Sprintf(`{"error": "Stock insuficiente para %s. Disponible: %d, solicitado: %d"}`,
				product.Name, product.StockQuantity, item.Quantity), nil
		}

		items = append(items, validatedItem{Product: product, Quantity: item.Quantity})
		totalAmount += product.Price * float64(item.Quantity)
		currency = product.Currency
	}

	// Construir orden canonica serializable
	externalID := uuid.New().String()
	orderNumber := fmt.Sprintf("AI-%s", externalID[:8])
	now := time.Now().Format(time.RFC3339)
	businessID := deps.businessID

	orderItems := make([]serializableOrderItem, 0, len(items))
	for _, item := range items {
		itemTotal := item.Product.Price * float64(item.Quantity)
		productID := item.Product.ID
		orderItems = append(orderItems, serializableOrderItem{
			ProductID:   &productID,
			ProductSKU:  item.Product.SKU,
			ProductName: item.Product.Name,
			Quantity:    item.Quantity,
			UnitPrice:   item.Product.Price,
			TotalPrice:  itemTotal,
			Currency:    item.Product.Currency,
		})
	}

	order := serializableOrder{
		BusinessID:      &businessID,
		IntegrationType: "whatsapp_ai",
		Platform:        "whatsapp_ai",
		ExternalID:      externalID,
		OrderNumber:     orderNumber,
		Subtotal:        totalAmount,
		TotalAmount:     totalAmount,
		Currency:        currency,
		CurrencyPresentment:    currency,
		SubtotalPresentment:    totalAmount,
		TotalAmountPresentment: totalAmount,
		CustomerName:    input.CustomerName,
		CustomerPhone:   input.CustomerPhone,
		Status:          "pending",
		OriginalStatus:  "pending",
		OccurredAt:      now,
		ImportedAt:      now,
		OrderItems:      orderItems,
	}

	payload, err := json.Marshal(order)
	if err != nil {
		return "", fmt.Errorf("error marshaling order: %w", err)
	}

	if err := deps.orderPublisher.PublishOrder(ctx, payload); err != nil {
		return "", fmt.Errorf("error publishing order: %w", err)
	}

	// Retornar confirmacion al AI
	response, _ := json.Marshal(map[string]any{
		"success":      true,
		"order_number": orderNumber,
		"total":        totalAmount,
		"currency":     currency,
		"items_count":  len(items),
		"message":      fmt.Sprintf("Pedido %s creado exitosamente", orderNumber),
	})

	return string(response), nil
}

// serializableOrder estructura con tags JSON para publicar a la cola canonical
type serializableOrder struct {
	BusinessID              *uint                 `json:"business_id"`
	IntegrationID           uint                  `json:"integration_id"`
	IntegrationType         string                `json:"integration_type"`
	Platform                string                `json:"platform"`
	ExternalID              string                `json:"external_id"`
	OrderNumber             string                `json:"order_number"`
	Subtotal                float64               `json:"subtotal"`
	TotalAmount             float64               `json:"total_amount"`
	Currency                string                `json:"currency"`
	CurrencyPresentment     string                `json:"currency_presentment"`
	SubtotalPresentment     float64               `json:"subtotal_presentment"`
	TotalAmountPresentment  float64               `json:"total_amount_presentment"`
	CustomerName            string                `json:"customer_name"`
	CustomerPhone           string                `json:"customer_phone"`
	Status                  string                `json:"status"`
	OriginalStatus          string                `json:"original_status"`
	OccurredAt              string                `json:"occurred_at"`
	ImportedAt              string                `json:"imported_at"`
	OrderItems              []serializableOrderItem `json:"order_items"`
}

type serializableOrderItem struct {
	ProductID   *string `json:"product_id"`
	ProductSKU  string  `json:"product_sku"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	TotalPrice  float64 `json:"total_price"`
	Currency    string  `json:"currency"`
}
