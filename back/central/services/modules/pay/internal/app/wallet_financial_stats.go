package app

import (
	"context"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/dtos"
)

// GetFinancialStats obtiene estadísticas financieras (membresías + guías) para super admins
func (uc *walletUseCase) GetFinancialStats(ctx context.Context, dto *dtos.FinancialStatsDTO) (*dtos.FinancialStatsResponse, error) {
	uc.log.Info(ctx).
		Interface("business_id", dto.BusinessID).
		Str("start_date", dto.StartDate).
		Str("end_date", dto.EndDate).
		Msg("Fetching financial stats")

	// Validar fechas
	startDate, err := time.Parse("2006-01-02", dto.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start_date: %w", err)
	}

	endDate, err := time.Parse("2006-01-02", dto.EndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end_date: %w", err)
	}

	if endDate.Before(startDate) {
		return nil, fmt.Errorf("end_date cannot be before start_date")
	}

	// Obtener stats del repositorio
	response, err := uc.repo.GetFinancialStats(ctx, dto)
	if err != nil {
		uc.log.Error(ctx).
			Err(err).
			Msg("Failed to get financial stats from repository")
		return nil, err
	}

	uc.log.Info(ctx).
		Float64("total_income", response.TotalIncome).
		Int("businesses_count", len(response.Businesses)).
		Msg("Financial stats retrieved successfully")

	return response, nil
}
