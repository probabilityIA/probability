package usecases

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/app/usecases/mapper"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
)

// ProcessNotification procesa una notificación IPN de MercadoLibre.
// MeLi solo envía topic + resource URL; debemos obtener la orden completa de la API.
func (uc *meliUseCase) ProcessNotification(ctx context.Context, notification *domain.MeliNotification) error {
	// 1. Filtrar por topic — solo procesar "orders_v2"
	if notification.Topic != "orders_v2" {
		uc.logger.Info(ctx).
			Str("topic", notification.Topic).
			Str("resource", notification.Resource).
			Msg("Ignoring non-order notification")
		return nil
	}

	// 2. Extraer order_id del resource path ("/orders/123456789")
	orderID, err := extractOrderIDFromResource(notification.Resource)
	if err != nil {
		uc.logger.Error(ctx).Err(err).
			Str("resource", notification.Resource).
			Msg("Failed to extract order ID from notification resource")
		return fmt.Errorf("extracting order_id: %w", err)
	}

	uc.logger.Info(ctx).
		Int64("order_id", orderID).
		Int64("user_id", notification.UserID).
		Msg("Processing MercadoLibre order notification")

	// 3. Buscar integración por seller_id (notification.UserID = seller_id de MeLi)
	sellerID := fmt.Sprintf("%d", notification.UserID)
	integration, err := uc.service.GetIntegrationByStoreID(ctx, sellerID)
	if err != nil {
		uc.logger.Error(ctx).Err(err).
			Str("seller_id", sellerID).
			Msg("Failed to find integration by seller_id")
		return fmt.Errorf("finding integration: %w", err)
	}
	if integration == nil {
		uc.logger.Warn(ctx).
			Str("seller_id", sellerID).
			Msg("No integration found for MeLi seller_id")
		return domain.ErrIntegrationNotFound
	}

	integrationID := fmt.Sprintf("%d", integration.ID)

	// 4. Obtener access_token vigente
	accessToken, err := uc.EnsureValidToken(ctx, integrationID)
	if err != nil {
		uc.logger.Error(ctx).Err(err).
			Uint("integration_id", integration.ID).
			Msg("Failed to ensure valid token")
		return fmt.Errorf("ensuring valid token: %w", err)
	}

	// 5. Obtener orden completa de la API de MeLi
	order, rawJSON, err := uc.client.GetOrder(ctx, accessToken, orderID)
	if err != nil {
		uc.logger.Error(ctx).Err(err).
			Int64("order_id", orderID).
			Msg("Failed to fetch order from MeLi API")
		return fmt.Errorf("fetching order: %w", err)
	}

	// 6. Obtener detalles de envío si hay shipping
	var shippingDetail *domain.MeliShippingDetail
	if order.Shipping != nil && order.Shipping.ID > 0 {
		shippingDetail, err = uc.client.GetShipmentDetail(ctx, accessToken, order.Shipping.ID)
		if err != nil {
			// No es fatal — logear y continuar sin datos de envío
			uc.logger.Warn(ctx).Err(err).
				Int64("shipment_id", order.Shipping.ID).
				Msg("Failed to fetch shipment detail, continuing without shipping data")
		}
	}

	// 7. Mapear a DTO canónico
	dto := mapper.MapMeliOrderToProbability(order, shippingDetail, rawJSON)
	dto.IntegrationID = integration.ID
	dto.BusinessID = integration.BusinessID

	// 8. Publicar a la cola
	if err := uc.publisher.Publish(ctx, dto); err != nil {
		uc.logger.Error(ctx).Err(err).
			Int64("order_id", orderID).
			Msg("Failed to publish MeLi order to queue")
		return fmt.Errorf("publishing order: %w", err)
	}

	uc.logger.Info(ctx).
		Int64("order_id", orderID).
		Str("status", order.Status).
		Uint("integration_id", integration.ID).
		Msg("MercadoLibre order published successfully")

	return nil
}

// extractOrderIDFromResource extrae el order_id numérico del resource path.
// Ejemplo: "/orders/123456789" → 123456789
func extractOrderIDFromResource(resource string) (int64, error) {
	// El resource puede ser "/orders/123456789" o "orders/123456789"
	parts := strings.Split(strings.Trim(resource, "/"), "/")
	if len(parts) < 2 {
		return 0, fmt.Errorf("invalid resource format: %s", resource)
	}
	// El último segmento es el ID
	idStr := parts[len(parts)-1]
	return strconv.ParseInt(idStr, 10, 64)
}
