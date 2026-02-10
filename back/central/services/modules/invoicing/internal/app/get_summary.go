package app

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
)

// GetSummary obtiene un resumen general de facturas con KPIs principales
func (uc *useCase) GetSummary(ctx context.Context, businessID uint, period string) (*entities.InvoiceSummary, error) {
	uc.log.Info(ctx).Uint("business_id", businessID).Str("period", period).Msg("Getting invoice summary")

	// Calcular fechas según el período solicitado
	start, end := calculatePeriod(period)

	// Obtener resumen desde el repositorio
<<<<<<< HEAD
	summary, err := uc.invoiceRepo.GetSummary(ctx, businessID, start, end)
=======
	summary, err := uc.repo.GetInvoiceSummary(ctx, businessID, start, end)
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to get invoice summary")
		return nil, err
	}

	uc.log.Info(ctx).
		Int("total_invoices", summary.Totals.TotalInvoices).
		Float64("total_amount", summary.Totals.TotalAmount).
		Msg("Invoice summary retrieved successfully")

	return summary, nil
}

// calculatePeriod convierte un string de período a fechas de inicio y fin
func calculatePeriod(period string) (time.Time, time.Time) {
	now := time.Now()

	switch period {
	case "today":
		return startOfDay(now), endOfDay(now)
	case "week":
		return startOfWeek(now), endOfWeek(now)
	case "year":
		return startOfYear(now), endOfYear(now)
	case "all":
		// Últimos 10 años
		return now.AddDate(-10, 0, 0), now
	default: // month (por defecto)
		return startOfMonth(now), endOfMonth(now)
	}
}

// startOfDay retorna el inicio del día (00:00:00)
func startOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

// endOfDay retorna el fin del día (23:59:59)
func endOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 23, 59, 59, 999999999, t.Location())
}

// startOfWeek retorna el inicio de la semana (lunes 00:00:00)
func startOfWeek(t time.Time) time.Time {
	// Ajustar al lunes más reciente
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7 // Domingo = 7
	}
	daysToSubtract := weekday - 1
	monday := t.AddDate(0, 0, -daysToSubtract)
	return startOfDay(monday)
}

// endOfWeek retorna el fin de la semana (domingo 23:59:59)
func endOfWeek(t time.Time) time.Time {
	start := startOfWeek(t)
	sunday := start.AddDate(0, 0, 6)
	return endOfDay(sunday)
}

// startOfMonth retorna el inicio del mes (día 1 a las 00:00:00)
func startOfMonth(t time.Time) time.Time {
	year, month, _ := t.Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, t.Location())
}

// endOfMonth retorna el fin del mes (último día a las 23:59:59)
func endOfMonth(t time.Time) time.Time {
	start := startOfMonth(t)
	nextMonth := start.AddDate(0, 1, 0)
	lastDay := nextMonth.Add(-time.Second)
	return lastDay
}

// startOfYear retorna el inicio del año (1 de enero a las 00:00:00)
func startOfYear(t time.Time) time.Time {
	return time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())
}

// endOfYear retorna el fin del año (31 de diciembre a las 23:59:59)
func endOfYear(t time.Time) time.Time {
	return time.Date(t.Year(), 12, 31, 23, 59, 59, 999999999, t.Location())
}
