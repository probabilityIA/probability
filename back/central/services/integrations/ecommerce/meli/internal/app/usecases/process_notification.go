package usecases

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/canonical"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/app/usecases/mapper"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
)

func (uc *meliUseCase) ProcessNotification(ctx context.Context, notification *domain.MeliNotification) error {
	switch notification.Topic {
	case "orders_v2", "orders":
		return uc.processOrderNotification(ctx, notification)
	case "shipments":
		return uc.processShipmentNotification(ctx, notification)
	case "payments":
		return uc.processPaymentNotification(ctx, notification)
	case "items", "items_prices", "stock-locations":
		return uc.processItemNotification(ctx, notification)
	case "claims", "messages", "post_purchase":
		return uc.processClaimNotification(ctx, notification)
	default:
		uc.logger.Info(ctx).
			Str("topic", notification.Topic).
			Str("resource", notification.Resource).
			Msg("Ignoring unsupported MercadoLibre notification topic")
		return nil
	}
}

func (uc *meliUseCase) resolveIntegrationAndToken(ctx context.Context, userID int64) (*domain.Integration, string, error) {
	sellerID := fmt.Sprintf("%d", userID)
	integration, err := uc.service.GetIntegrationByStoreID(ctx, sellerID)
	if err != nil {
		return nil, "", fmt.Errorf("finding integration: %w", err)
	}
	if integration == nil {
		uc.logger.Warn(ctx).Str("seller_id", sellerID).Msg("No integration found for MeLi seller_id")
		return nil, "", domain.ErrIntegrationNotFound
	}
	integrationID := fmt.Sprintf("%d", integration.ID)
	accessToken, err := uc.EnsureValidToken(ctx, integrationID)
	if err != nil {
		return nil, "", fmt.Errorf("ensuring valid token: %w", err)
	}
	return integration, accessToken, nil
}

func (uc *meliUseCase) fetchOrderDTO(ctx context.Context, integration *domain.Integration, accessToken string, orderID int64) (*canonical.ProbabilityOrderDTO, error) {
	order, rawJSON, err := uc.client.GetOrder(ctx, accessToken, orderID)
	if err != nil {
		return nil, fmt.Errorf("fetching order: %w", err)
	}

	uc.enrichBillingInfo(ctx, accessToken, order)

	if order.PackID != nil && *order.PackID > 0 {
		merged, perr := uc.consolidatePack(ctx, accessToken, *order.PackID, order)
		if perr != nil {
			uc.logger.Warn(ctx).Err(perr).Int64("pack_id", *order.PackID).Msg("Failed to consolidate pack, using single order")
		} else if merged != nil {
			order = merged
		}
	}

	var shippingDetail *domain.MeliShippingDetail
	if order.Shipping != nil && order.Shipping.ID > 0 {
		shippingDetail, err = uc.client.GetShipmentDetail(ctx, accessToken, order.Shipping.ID)
		if err != nil {
			uc.logger.Warn(ctx).Err(err).
				Int64("shipment_id", order.Shipping.ID).
				Msg("Failed to fetch shipment detail, continuing without shipping data")
			shippingDetail = nil
		}
	}

	dto := mapper.MapMeliOrderToProbability(order, shippingDetail, rawJSON)
	dto.IntegrationID = integration.ID
	dto.BusinessID = integration.BusinessID
	return dto, nil
}

func (uc *meliUseCase) publishOrder(ctx context.Context, integration *domain.Integration, accessToken string, orderID int64) error {
	dto, err := uc.fetchOrderDTO(ctx, integration, accessToken, orderID)
	if err != nil {
		return err
	}
	if err := uc.publisher.Publish(ctx, dto); err != nil {
		return fmt.Errorf("publishing order: %w", err)
	}
	uc.logger.Info(ctx).
		Int64("order_id", orderID).
		Uint("integration_id", integration.ID).
		Msg("MercadoLibre order published")
	return nil
}

func (uc *meliUseCase) processOrderNotification(ctx context.Context, notification *domain.MeliNotification) error {
	orderID, err := extractResourceIntID(notification.Resource)
	if err != nil {
		return fmt.Errorf("extracting order_id: %w", err)
	}
	integration, accessToken, err := uc.resolveIntegrationAndToken(ctx, notification.UserID)
	if err != nil {
		return err
	}
	return uc.publishOrder(ctx, integration, accessToken, orderID)
}

func (uc *meliUseCase) processShipmentNotification(ctx context.Context, notification *domain.MeliNotification) error {
	shipmentID, err := extractResourceIntID(notification.Resource)
	if err != nil {
		return fmt.Errorf("extracting shipment_id: %w", err)
	}
	integration, accessToken, err := uc.resolveIntegrationAndToken(ctx, notification.UserID)
	if err != nil {
		return err
	}
	orderIDs, err := uc.client.GetShipmentOrderIDs(ctx, accessToken, shipmentID)
	if err != nil {
		return fmt.Errorf("resolving shipment orders: %w", err)
	}
	if len(orderIDs) == 0 {
		uc.logger.Warn(ctx).Int64("shipment_id", shipmentID).Msg("No orders found for shipment")
		return nil
	}
	for _, orderID := range orderIDs {
		if perr := uc.publishOrder(ctx, integration, accessToken, orderID); perr != nil {
			uc.logger.Error(ctx).Err(perr).Int64("order_id", orderID).Msg("Failed to publish order from shipment notification")
		}
	}
	return nil
}

func (uc *meliUseCase) processPaymentNotification(ctx context.Context, notification *domain.MeliNotification) error {
	paymentID, err := extractResourceIntID(notification.Resource)
	if err != nil {
		return fmt.Errorf("extracting payment_id: %w", err)
	}
	integration, accessToken, err := uc.resolveIntegrationAndToken(ctx, notification.UserID)
	if err != nil {
		return err
	}
	orderID, err := uc.client.GetPaymentOrderID(ctx, accessToken, paymentID)
	if err != nil {
		if err == domain.ErrOrderNotFound {
			uc.logger.Info(ctx).Int64("payment_id", paymentID).Msg("Payment has no associated order, skipping")
			return nil
		}
		return fmt.Errorf("resolving payment order: %w", err)
	}
	return uc.publishOrder(ctx, integration, accessToken, orderID)
}

func (uc *meliUseCase) processClaimNotification(ctx context.Context, notification *domain.MeliNotification) error {
	claimID, err := extractResourceIntID(notification.Resource)
	if err != nil {
		return fmt.Errorf("extracting claim_id: %w", err)
	}
	integration, accessToken, err := uc.resolveIntegrationAndToken(ctx, notification.UserID)
	if err != nil {
		return err
	}
	claim, err := uc.client.GetClaim(ctx, accessToken, claimID)
	if err != nil {
		if err == domain.ErrOrderNotFound {
			return nil
		}
		return fmt.Errorf("fetching claim: %w", err)
	}
	if claim.ResourceType != "order" || claim.ResourceID <= 0 {
		uc.logger.Info(ctx).Int64("claim_id", claimID).Msg("Claim not attached to an order, skipping")
		return nil
	}
	dto, err := uc.fetchOrderDTO(ctx, integration, accessToken, claim.ResourceID)
	if err != nil {
		return err
	}
	note := buildClaimNote(claim)
	dto.Notes = &note
	if err := uc.publisher.Publish(ctx, dto); err != nil {
		return fmt.Errorf("publishing order with claim note: %w", err)
	}
	return nil
}

func buildClaimNote(claim *domain.MeliClaim) string {
	parts := []string{fmt.Sprintf("Reclamo MercadoLibre (%s)", claim.Status)}
	if claim.Reason != "" {
		parts = append(parts, claim.Reason)
	}
	if len(claim.Messages) > 0 {
		parts = append(parts, strings.Join(claim.Messages, " | "))
	}
	return strings.Join(parts, ": ")
}

func extractResourceIntID(resource string) (int64, error) {
	segment := extractResourceStringID(resource)
	if segment == "" {
		return 0, fmt.Errorf("invalid resource format: %s", resource)
	}
	return strconv.ParseInt(segment, 10, 64)
}

func extractResourceStringID(resource string) string {
	trimmed := strings.Trim(resource, "/")
	if trimmed == "" {
		return ""
	}
	if idx := strings.IndexByte(trimmed, '?'); idx >= 0 {
		trimmed = trimmed[:idx]
	}
	parts := strings.Split(trimmed, "/")
	return parts[len(parts)-1]
}
