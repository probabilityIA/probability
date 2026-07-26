package worker

import (
	"testing"
	"time"
)

func TestNextRun(t *testing.T) {
	location, err := time.LoadLocation(syncTimeZone)
	if err != nil {
		t.Skipf("zona horaria no disponible: %v", err)
	}

	tests := []struct {
		name         string
		ahora        time.Time
		esperadoHora int
		esperadoDia  int
	}{
		{
			name:         "madrugada corre a las 7",
			ahora:        time.Date(2026, 7, 26, 2, 30, 0, 0, location),
			esperadoHora: 7,
			esperadoDia:  26,
		},
		{
			name:         "media manana corre a las 13",
			ahora:        time.Date(2026, 7, 26, 9, 0, 0, 0, location),
			esperadoHora: 13,
			esperadoDia:  26,
		},
		{
			name:         "tarde corre a las 19",
			ahora:        time.Date(2026, 7, 26, 15, 45, 0, 0, location),
			esperadoHora: 19,
			esperadoDia:  26,
		},
		{
			name:         "noche pasa al dia siguiente",
			ahora:        time.Date(2026, 7, 26, 22, 10, 0, 0, location),
			esperadoHora: 7,
			esperadoDia:  27,
		},
		{
			name:         "justo en la hora programada espera la siguiente",
			ahora:        time.Date(2026, 7, 26, 13, 0, 0, 0, location),
			esperadoHora: 19,
			esperadoDia:  26,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := nextRun(tt.ahora, location)
			if got.Hour() != tt.esperadoHora || got.Day() != tt.esperadoDia {
				t.Errorf("nextRun() = %s, se esperaba dia %d a las %d:00", got.Format(time.RFC3339), tt.esperadoDia, tt.esperadoHora)
			}
			if !got.After(tt.ahora) {
				t.Errorf("nextRun() = %s debe ser posterior a %s", got, tt.ahora)
			}
		})
	}
}

func TestSyncHoursSonTresAlDia(t *testing.T) {
	if len(syncHours) != 3 {
		t.Fatalf("se esperaban 3 corridas diarias, hay %d", len(syncHours))
	}
	for i := 1; i < len(syncHours); i++ {
		if syncHours[i] <= syncHours[i-1] {
			t.Errorf("las horas deben ir en orden ascendente: %v", syncHours)
		}
	}
}
