package consumer

import (
	"testing"
	"time"
)

func TestProviderUnavailableRetryDelay(t *testing.T) {
	cases := []struct {
		name    string
		elapsed time.Duration
		want    time.Duration
	}{
		{"primer fallo", 0, 15 * time.Minute},
		{"dentro de la primera hora", 59 * time.Minute, 15 * time.Minute},
		{"segunda hora", time.Hour, 30 * time.Minute},
		{"antes de las dos horas", 119 * time.Minute, 30 * time.Minute},
		{"tercera hora", 2 * time.Hour, time.Hour},
		{"antes de las cuatro horas", 239 * time.Minute, time.Hour},
		{"despues de cuatro horas", 4 * time.Hour, 2 * time.Hour},
		{"muy vieja", 30 * time.Hour, 2 * time.Hour},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := providerUnavailableRetryDelay(tc.elapsed); got != tc.want {
				t.Fatalf("elapsed %s: esperado %s, obtenido %s", tc.elapsed, tc.want, got)
			}
		})
	}
}

func TestProviderUnavailableRetryDelayNoDecrece(t *testing.T) {
	var previo time.Duration
	for elapsed := time.Duration(0); elapsed <= 8*time.Hour; elapsed += 5 * time.Minute {
		actual := providerUnavailableRetryDelay(elapsed)
		if actual < previo {
			t.Fatalf("el backoff bajo de %s a %s en elapsed %s", previo, actual, elapsed)
		}
		previo = actual
	}
}

func TestProviderUnavailableMaxWindowCortaAntesDeUnDia(t *testing.T) {
	if providerUnavailableMaxWindow <= 0 {
		t.Fatal("la ventana debe ser positiva")
	}
	if providerUnavailableMaxWindow > 24*time.Hour {
		t.Fatalf("ventana demasiado larga: %s", providerUnavailableMaxWindow)
	}
	if maxDelay := providerUnavailableRetryDelay(providerUnavailableMaxWindow); maxDelay >= providerUnavailableMaxWindow {
		t.Fatalf("el delay %s no permite ningun reintento dentro de la ventana %s", maxDelay, providerUnavailableMaxWindow)
	}
}
