package app

import (
	"context"
	"fmt"
	"strings"

	"github.com/secamc93/probability/back/central/services/modules/orderscompare/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/shared/orderscompare"
)

const maxApplyBatch = 200

func (uc *useCase) Apply(ctx context.Context, cmd dtos.ApplyCommand) (*dtos.ApplyResult, error) {
	integration, err := uc.resolveIntegration(ctx, cmd.IntegrationID, cmd.BusinessID)
	if err != nil {
		return nil, err
	}

	if !uc.channels.Supports(integration.IntegrationType) {
		return nil, fmt.Errorf("el canal %s todavia no permite importar ordenes desde el comparativo", integration.Name)
	}

	wanted := dedupe(cmd.ExternalIDs)
	if len(wanted) == 0 {
		return nil, fmt.Errorf("no se recibio ninguna orden para crear")
	}
	if len(wanted) > maxApplyBatch {
		return nil, fmt.Errorf("maximo %d ordenes por lote, se recibieron %d", maxApplyBatch, len(wanted))
	}

	locals, err := uc.repo.ListLocalOrders(ctx, integration.ID, cmd.BusinessID, nil, nil, 0)
	if err != nil {
		return nil, fmt.Errorf("leyendo las ordenes de Probability: %w", err)
	}
	existing := make(map[string]bool, len(locals))
	for _, local := range locals {
		existing[strings.ToLower(strings.TrimSpace(local.ExternalID))] = true
	}

	pending := make([]string, 0, len(wanted))
	skipped := make([]string, 0)
	for _, id := range wanted {
		if existing[strings.ToLower(id)] {
			skipped = append(skipped, id)
			continue
		}
		pending = append(pending, id)
	}

	result := &dtos.ApplyResult{
		Queued:           []string{},
		Skipped:          skipped,
		WithoutInventory: []string{},
	}

	if len(pending) == 0 {
		result.Note = "todas las ordenes seleccionadas ya existian en Probability"
		return result, nil
	}

	channelOrders, err := uc.channels.ListOrders(ctx, integration.IntegrationType, fmt.Sprintf("%d", integration.ID), orderscompare.ChannelFilters{Limit: maxLimit})
	if err != nil {
		uc.logger.Warn(ctx).Err(err).Uint("integration_id", integration.ID).
			Msg("No se pudo releer el canal para clasificar el inventario de las ordenes a importar")
	}
	statusByExternal := make(map[string]string, len(channelOrders))
	for _, order := range channelOrders {
		statusByExternal[strings.ToLower(strings.TrimSpace(order.ExternalID))] = order.Status
	}

	for _, id := range pending {
		if skips, _ := orderscompare.SkipsInventory(statusByExternal[strings.ToLower(id)]); skips {
			result.WithoutInventory = append(result.WithoutInventory, id)
		}
	}

	imported, err := uc.channels.ImportOrders(ctx, integration.IntegrationType, fmt.Sprintf("%d", integration.ID), pending)
	if err != nil {
		return nil, fmt.Errorf("importando las ordenes del canal: %w", err)
	}

	result.Queued = imported.Queued
	result.Failed = imported.Failed
	if len(imported.NotFound) > 0 {
		result.Failed = mergeFailures(result.Failed, imported.NotFound, "la orden ya no existe en el canal")
	}

	if len(result.WithoutInventory) > 0 {
		result.Note = fmt.Sprintf(
			"%d de %d ordenes entran como historicas: no mueven inventario porque el canal ya las dio por entregadas, despachadas, canceladas o devueltas",
			len(result.WithoutInventory), len(pending),
		)
	}

	uc.logger.Info(ctx).
		Uint("integration_id", integration.ID).
		Uint("business_id", cmd.BusinessID).
		Int("queued", len(result.Queued)).
		Int("skipped", len(result.Skipped)).
		Int("without_inventory", len(result.WithoutInventory)).
		Msg("Ordenes encoladas desde el comparativo")

	return result, nil
}

func dedupe(ids []string) []string {
	seen := make(map[string]bool, len(ids))
	out := make([]string, 0, len(ids))
	for _, id := range ids {
		clean := strings.TrimSpace(id)
		if clean == "" {
			continue
		}
		key := strings.ToLower(clean)
		if seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, clean)
	}
	return out
}

func mergeFailures(failed map[string]string, ids []string, reason string) map[string]string {
	if failed == nil {
		failed = make(map[string]string, len(ids))
	}
	for _, id := range ids {
		failed[id] = reason
	}
	return failed
}
