package usecaseupdateorder

import (
	"context"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
)

// UpdateOrder actualiza una orden existente con los datos del DTO
func (uc *UseCaseUpdateOrder) UpdateOrder(ctx context.Context, existingOrder *entities.ProbabilityOrder, dto *dtos.ProbabilityOrderDTO) (*dtos.OrderResponse, error) {
	// Guardar el estado anterior antes de actualizar (para detectar cambios de estado)
	previousStatus := existingOrder.Status

	uc.saveChannelMetadata(ctx, existingOrder, dto)

	// 1. Validar y actualizar todos los campos de la orden
	hasChanges := uc.updateOrderFields(ctx, existingOrder, dto)

	// 2. Si no hay cambios, notificar como omitida (para el sync) y retornar sin actualizar
	if !hasChanges {
		if !dto.IsManualOrder {
			uc.publishSyncOrderSkipped(ctx, existingOrder)
		}
		return uc.mapOrderToResponse(existingOrder), nil
	}

	statusChanged := previousStatus != existingOrder.Status
	if statusChanged && !dto.IsManualOrder {
		existingOrder.StatusSource = entities.StatusSourceSalesChannel
		existingOrder.StatusChangedBy = entities.SalesChannelChangedBy(existingOrder.Platform)
		statusChangedAt := time.Now()
		existingOrder.StatusChangedAt = &statusChangedAt
	}

	if err := uc.repo.UpdateOrder(ctx, existingOrder); err != nil {
		return nil, fmt.Errorf("error updating order: %w", err)
	}

	if statusChanged {
		if err := uc.repo.CreateOrderHistory(ctx, &entities.OrderHistory{
			OrderID:        existingOrder.ID,
			PreviousStatus: previousStatus,
			NewStatus:      existingOrder.Status,
			ChangedByName:  existingOrder.StatusChangedBy,
			Source:         existingOrder.StatusSource,
		}); err != nil {
			uc.logger.Warn(ctx).Err(err).Str("order_id", existingOrder.ID).Msg("No se pudo registrar el historial del cambio de estado")
		}
	}

	if existingOrder.BusinessID != nil {
		if err := uc.repo.ResolveOrderGeozone(ctx, existingOrder.ID, *existingOrder.BusinessID); err != nil {
			uc.logger.Warn(ctx).Err(err).Str("order_id", existingOrder.ID).Msg("Failed to resolve order geozone")
		}
	}

	uc.publishUpdateEvents(ctx, existingOrder, previousStatus, dto.IsManualOrder)

	// 5. Retornar la respuesta actualizada
	return uc.mapOrderToResponse(existingOrder), nil
}

// updateOrderFields actualiza todos los campos de la orden y retorna si hubo cambios
func (uc *UseCaseUpdateOrder) updateOrderFields(ctx context.Context, order *entities.ProbabilityOrder, dto *dtos.ProbabilityOrderDTO) bool {
	hasChanges := false

	// Actualizar estados
	if uc.updateOrderStatuses(ctx, order, dto) {
		hasChanges = true
	}

	// Actualizar información financiera
	if uc.updateFinancialFields(order, dto) {
		hasChanges = true
	}

	// Actualizar información del cliente
	if uc.updateCustomerFields(order, dto) {
		hasChanges = true
	}

	// Actualizar dirección de envío
	if uc.updateShippingFields(order, dto) {
		hasChanges = true
	}

	// Actualizar información de pago
	if uc.updatePaymentFields(ctx, order, dto) {
		hasChanges = true
	}

	// Actualizar información de fulfillment
	if uc.updateFulfillmentFields(order, dto) {
		hasChanges = true
	}

	// Actualizar campos adicionales
	if uc.updateAdditionalFields(order, dto) {
		hasChanges = true
	}

	// Actualizar datos estructurados (JSONB)
	if uc.updateStructuredData(order, dto) {
		hasChanges = true
	}

	return hasChanges
}

func (uc *UseCaseUpdateOrder) saveChannelMetadata(ctx context.Context, order *entities.ProbabilityOrder, dto *dtos.ProbabilityOrderDTO) {
	if dto.ChannelMetadata == nil || len(dto.ChannelMetadata.RawData) == 0 {
		return
	}

	if err := uc.repo.MarkChannelMetadataNotLatest(ctx, order.ID); err != nil {
		uc.logger.Warn(ctx).Err(err).Str("order_id", order.ID).Msg("No se pudo marcar el raw anterior como no vigente")
	}

	metadata := &entities.ProbabilityOrderChannelMetadata{
		OrderID:       order.ID,
		ChannelSource: dto.ChannelMetadata.ChannelSource,
		IntegrationID: dto.IntegrationID,
		RawData:       dto.ChannelMetadata.RawData,
		Version:       dto.ChannelMetadata.Version,
		ReceivedAt:    dto.ChannelMetadata.ReceivedAt,
		ProcessedAt:   dto.ChannelMetadata.ProcessedAt,
		IsLatest:      true,
		LastSyncedAt:  dto.ChannelMetadata.LastSyncedAt,
		SyncStatus:    dto.ChannelMetadata.SyncStatus,
	}
	if metadata.ChannelSource == "" {
		metadata.ChannelSource = order.Platform
	}
	if metadata.IntegrationID == 0 {
		metadata.IntegrationID = order.IntegrationID
	}
	if metadata.ReceivedAt.IsZero() {
		metadata.ReceivedAt = time.Now()
	}
	if metadata.SyncStatus == "" {
		metadata.SyncStatus = "synced"
	}

	if err := uc.repo.CreateChannelMetadata(ctx, metadata); err != nil {
		uc.logger.Warn(ctx).Err(err).Str("order_id", order.ID).Msg("No se pudo guardar el JSON crudo del canal")
	}
}
