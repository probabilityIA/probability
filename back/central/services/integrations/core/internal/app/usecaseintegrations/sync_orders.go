package usecaseintegrations

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// SyncOrdersByIntegrationID sincroniza órdenes para una integración específica.
func (uc *IntegrationUseCase) SyncOrdersByIntegrationID(ctx context.Context, integrationID string) error {
	ctx = log.WithFunctionCtx(ctx, "SyncOrdersByIntegrationID")

	provider, err := uc.getProviderForIntegration(ctx, integrationID)
	if err != nil {
		return err
	}

	return provider.SyncOrdersByIntegrationID(ctx, integrationID)
}

// SyncOrdersByIntegrationIDWithParams sincroniza órdenes con parámetros de filtrado.
// Si la integración no soporta params, hace fallback a SyncOrdersByIntegrationID.
func (uc *IntegrationUseCase) SyncOrdersByIntegrationIDWithParams(ctx context.Context, integrationID string, params interface{}) error {
	ctx = log.WithFunctionCtx(ctx, "SyncOrdersByIntegrationIDWithParams")

	provider, err := uc.getProviderForIntegration(ctx, integrationID)
	if err != nil {
		return err
	}

	// Intentar con parámetros; si no está soportado, fallback sin params
	err = provider.SyncOrdersByIntegrationIDWithParams(ctx, integrationID, params)
	if errors.Is(err, domain.ErrNotSupported) {
		return provider.SyncOrdersByIntegrationID(ctx, integrationID)
	}
	return err
}

// SyncOrdersByIntegrationIDWithBatches divide un rango de fechas grande en chunks de 7 días,
// publica cada chunk como mensaje a la cola de lotes, y retorna inmediatamente.
// Si el rango resulta en un solo chunk, usa el flujo directo (sin overhead de cola).
func (uc *IntegrationUseCase) SyncOrdersByIntegrationIDWithBatches(ctx context.Context, integrationID string, params *domain.SyncBatchParams) error {
	ctx = log.WithFunctionCtx(ctx, "SyncOrdersByIntegrationIDWithBatches")

	// Resolver integración para obtener type_id y business_id
	integration, err := uc.GetPublicIntegrationByID(ctx, integrationID)
	if err != nil {
		return fmt.Errorf("error al obtener integración: %w", err)
	}

	// Determinar rango de fechas
	now := time.Now()
	start := now.AddDate(0, 0, -30)
	end := now

	if params.CreatedAtMin != nil {
		start = *params.CreatedAtMin
	}
	if params.CreatedAtMax != nil {
		end = *params.CreatedAtMax
	}

	// Dividir en chunks de 7 días
	chunks := SplitDateRange(start, end, 7)

	// Si solo 1 chunk → flujo directo (sin overhead de cola)
	if len(chunks) <= 1 {
		genericParams := params.ToGenericMap()
		return uc.SyncOrdersByIntegrationIDWithParams(ctx, integrationID, genericParams)
	}

	jobID := uuid.New().String()
	uc.log.Info(ctx).
		Str("job_id", jobID).
		Str("integration_id", integrationID).
		Int("total_batches", len(chunks)).
		Time("date_from", start).
		Time("date_to", end).
		Msg("📦 Creando job de sincronización por lotes")

	// Publicar cada chunk como mensaje
	for i, chunk := range chunks {
		msg := domain.SyncBatchMessage{
			JobID:             jobID,
			IntegrationID:     integrationID,
			IntegrationTypeID: integration.IntegrationType,
			BusinessID:        integration.BusinessID,
			BatchIndex:        i,
			TotalBatches:      len(chunks),
			CreatedAtMin:      chunk.Start,
			CreatedAtMax:      chunk.End,
			Status:            params.Status,
			FinancialStatus:   params.FinancialStatus,
			FulfillmentStatus: params.FulfillmentStatus,
			EnqueuedAt:        time.Now(),
		}

		msgBytes, err := json.Marshal(msg)
		if err != nil {
			uc.log.Error(ctx).Err(err).Int("batch_index", i).Msg("Error al serializar mensaje de lote")
			return fmt.Errorf("error serializando lote %d: %w", i, err)
		}

		if err := uc.queue.Publish(ctx, rabbitmq.QueueSyncBatches, msgBytes); err != nil {
			uc.log.Error(ctx).Err(err).Int("batch_index", i).Msg("Error al publicar mensaje de lote")
			return fmt.Errorf("error publicando lote %d: %w", i, err)
		}
	}

	// Publicar evento SSE de inicio
	var businessID uint
	if integration.BusinessID != nil {
		businessID = *integration.BusinessID
	}
	rabbitmq.PublishEvent(ctx, uc.queue, rabbitmq.EventEnvelope{ //nolint:errcheck
		Type:          "integration.sync.batched.started",
		Category:      "integration",
		BusinessID:    businessID,
		IntegrationID: integration.ID,
		Data: map[string]interface{}{
			"job_id":         jobID,
			"integration_id": integrationID,
			"total_batches":  len(chunks),
			"date_from":      start.Format(time.RFC3339),
			"date_to":        end.Format(time.RFC3339),
			"chunk_days":     7,
		},
	})

	uc.log.Info(ctx).
		Str("job_id", jobID).
		Int("total_batches", len(chunks)).
		Msg("✅ Todos los mensajes de lote publicados exitosamente")

	return nil
}

// SyncOrdersByBusiness sincroniza órdenes de todas las integraciones activas de un negocio
// que soporten sincronización.
func (uc *IntegrationUseCase) SyncOrdersByBusiness(ctx context.Context, businessID uint) error {
	ctx = log.WithFunctionCtx(ctx, "SyncOrdersByBusiness")

	businessIDPtr := &businessID
	filters := domain.IntegrationFilters{
		BusinessID: businessIDPtr,
		IsActive:   boolPtr(true),
	}

	integrations, _, err := uc.repo.ListIntegrations(ctx, filters)
	if err != nil {
		return fmt.Errorf("error al obtener integraciones: %w", err)
	}

	for _, integration := range integrations {
		if integration.IntegrationType == nil {
			continue
		}

		integrationTypeID := int(integration.IntegrationTypeID)
		impl, registered := uc.providerReg.Get(integrationTypeID)
		if !registered {
			continue
		}

		integrationIDStr := fmt.Sprintf("%d", integration.ID)
		err := impl.SyncOrdersByIntegrationID(ctx, integrationIDStr)
		if err != nil {
			// Skip silently — ErrNotSupported means this provider doesn't sync orders
			continue
		}
	}

	return nil
}
