package app

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
)

// GetDetailedStats obtiene estadísticas detalladas para dashboards
func (uc *useCase) GetDetailedStats(ctx context.Context, businessID uint, filters map[string]interface{}) (*entities.DetailedStats, error) {
	uc.log.Info(ctx).Uint("business_id", businessID).Msg("Getting detailed stats")

	// Parsear fechas de filtros
	processedFilters := processStatsFilters(filters)

	// Obtener estadísticas desde el repositorio
	stats, err := uc.repo.GetInvoiceDetailedStats(ctx, businessID, processedFilters)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to get detailed stats")
		return nil, err
	}

	uc.log.Info(ctx).
		Int("total_invoices", stats.Summary.TotalInvoices).
		Float64("success_rate", stats.Summary.SuccessRate).
		Msg("Detailed stats retrieved successfully")

	return stats, nil
}

// processStatsFilters procesa y valida los filtros de estadísticas
func processStatsFilters(filters map[string]interface{}) map[string]interface{} {
	processed := make(map[string]interface{})

	// Procesar start_date
	if startStr, ok := filters["start_date"].(string); ok {
		if start, err := time.Parse("2006-01-02", startStr); err == nil {
			processed["start_date"] = start
		}
	}

	// Procesar end_date
	if endStr, ok := filters["end_date"].(string); ok {
		if end, err := time.Parse("2006-01-02", endStr); err == nil {
			processed["end_date"] = end
		}
	}

	// Copiar otros filtros
	if integrationID, ok := filters["integration_id"]; ok {
		processed["integration_id"] = integrationID
	}
	if invoicingIntegrationID, ok := filters["invoicing_integration_id"]; ok {
		processed["invoicing_integration_id"] = invoicingIntegrationID
	}

	return processed
}
