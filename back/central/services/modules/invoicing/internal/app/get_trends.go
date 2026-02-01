package app

import (
	"context"
	"errors"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
)

// GetTrends obtiene datos de tendencias temporales para gráficos
func (uc *useCase) GetTrends(ctx context.Context, businessID uint, startDate, endDate, granularity, metric string) (*entities.TrendData, error) {
	uc.log.Info(ctx).
		Uint("business_id", businessID).
		Str("start_date", startDate).
		Str("end_date", endDate).
		Str("granularity", granularity).
		Str("metric", metric).
		Msg("Getting trends")

	// Validar y parsear fechas
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Invalid start_date format")
		return nil, errors.New("formato de start_date inválido, usar YYYY-MM-DD")
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Invalid end_date format")
		return nil, errors.New("formato de end_date inválido, usar YYYY-MM-DD")
	}

	// Validar que start < end
	if start.After(end) {
		return nil, errors.New("start_date debe ser anterior a end_date")
	}

	// Validar granularidad
	if granularity != "day" && granularity != "week" && granularity != "month" {
		granularity = "day" // Default
	}

	// Validar métrica
	if metric != "count" && metric != "amount" && metric != "success_rate" {
		metric = "count" // Default
	}

	// Obtener tendencias desde el repositorio
	trends, err := uc.invoiceRepo.GetTrends(ctx, businessID, start, end, granularity, metric)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to get trends")
		return nil, err
	}

	uc.log.Info(ctx).
		Int("data_points", len(trends.DataPoints)).
		Str("trend_direction", trends.Trend.Direction).
		Msg("Trends retrieved successfully")

	return trends, nil
}
