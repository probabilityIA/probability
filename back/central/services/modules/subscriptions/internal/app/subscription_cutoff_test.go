package app

import (
	"testing"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/entities"
)

func TestCutoffReached_NoCutoffConfigured(t *testing.T) {
	end := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	business := entities.ExpiringBusiness{EndDate: end, CutoffDay: nil}

	if cutoffReached(business, end) {
		t.Fatal("sin dia de corte configurado, debe dar 6 dias de gracia antes de expirar")
	}
	if cutoffReached(business, end.AddDate(0, 0, 5)) {
		t.Fatal("no deberia expirar antes de cumplirse los 6 dias de gracia")
	}
	if !cutoffReached(business, end.AddDate(0, 0, 6)) {
		t.Fatal("deberia expirar exactamente a los 6 dias de gracia")
	}
}

func TestCutoffReached_SameMonthFuture(t *testing.T) {
	end := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	cutoff := 5
	business := entities.ExpiringBusiness{EndDate: end, CutoffDay: &cutoff}

	if cutoffReached(business, time.Date(2026, 9, 3, 0, 0, 0, 0, time.UTC)) {
		t.Fatal("no deberia expirar antes del dia de corte")
	}
	if !cutoffReached(business, time.Date(2026, 9, 5, 0, 0, 0, 0, time.UTC)) {
		t.Fatal("deberia expirar en el dia de corte")
	}
	if !cutoffReached(business, time.Date(2026, 9, 20, 0, 0, 0, 0, time.UTC)) {
		t.Fatal("deberia seguir expirado despues del dia de corte")
	}
}

func TestCutoffReached_CutoffAlreadyPassedRollsToNextMonth(t *testing.T) {
	end := time.Date(2026, 9, 20, 0, 0, 0, 0, time.UTC)
	cutoff := 5
	business := entities.ExpiringBusiness{EndDate: end, CutoffDay: &cutoff}

	if cutoffReached(business, time.Date(2026, 9, 25, 0, 0, 0, 0, time.UTC)) {
		t.Fatal("el corte del 5 ya paso en septiembre, no deberia expirar hasta el 5 de octubre")
	}
	if !cutoffReached(business, time.Date(2026, 10, 5, 0, 0, 0, 0, time.UTC)) {
		t.Fatal("deberia expirar el 5 de octubre")
	}
}

func TestCutoffReached_ClampsToLastDayOfMonth(t *testing.T) {
	end := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)
	cutoff := 31
	business := entities.ExpiringBusiness{EndDate: end, CutoffDay: &cutoff}

	if cutoffReached(business, time.Date(2026, 2, 27, 0, 0, 0, 0, time.UTC)) {
		t.Fatal("febrero 2026 tiene 28 dias, no deberia expirar antes del 28")
	}
	if !cutoffReached(business, time.Date(2026, 2, 28, 0, 0, 0, 0, time.UTC)) {
		t.Fatal("deberia expirar el ultimo dia del mes cuando el corte configurado no existe")
	}
}

func TestCutoffReached_YearRollover(t *testing.T) {
	end := time.Date(2026, 12, 20, 0, 0, 0, 0, time.UTC)
	cutoff := 5
	business := entities.ExpiringBusiness{EndDate: end, CutoffDay: &cutoff}

	if cutoffReached(business, time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)) {
		t.Fatal("el corte del 5 ya paso en diciembre, no deberia expirar hasta el 5 de enero")
	}
	if !cutoffReached(business, time.Date(2027, 1, 5, 0, 0, 0, 0, time.UTC)) {
		t.Fatal("deberia expirar el 5 de enero del anio siguiente")
	}
}
