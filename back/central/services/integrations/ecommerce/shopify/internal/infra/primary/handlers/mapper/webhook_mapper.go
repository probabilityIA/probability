package mapper

import (
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/infra/secondary/client/mappers"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/infra/secondary/client/response"
)

// MapWebhookPayloadToOrderResponse convierte el payload del webhook (que viene sin wrapper)
// a la estructura OrderResponse que espera el mapper existente
func MapWebhookPayloadToOrderResponse(payload map[string]interface{}) (response.Order, error) {
	// El webhook envía los datos directamente, así que convertimos el map a JSON
	// y luego deserializamos a response.Order
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return response.Order{}, err
	}

	var orderResp response.Order
	if err := json.Unmarshal(jsonData, &orderResp); err != nil {
		return response.Order{}, err
	}

	return orderResp, nil
}

// MapWebhookPayloadToShopifyOrder mapea directamente el payload del webhook a ShopifyOrder
// usando el mapper existente después de convertir a response.Order
func MapWebhookPayloadToShopifyOrder(payload map[string]interface{}, businessID *uint, integrationID uint, integrationType string) (domain.ShopifyOrder, error) {
	// Convertir el payload completo a response.Order usando JSON
	orderResp, err := MapWebhookPayloadToOrderResponse(payload)
	if err != nil {
		return domain.ShopifyOrder{}, err
	}

	// Usar el mapper existente

	// Marshal payload again to get raw bytes for storage
	rawOrder, _ := json.Marshal(payload)

	// Usar el mapper existente
	shopifyOrder := mappers.MapOrderResponseToShopifyOrder(orderResp, rawOrder, businessID, integrationID, integrationType)

	return shopifyOrder, nil
}
