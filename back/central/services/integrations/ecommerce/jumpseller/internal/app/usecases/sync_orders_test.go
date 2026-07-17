package usecases

import (
	"testing"
	"time"
)

func TestBuildQueryParamsAceptaFechaDelFormulario(t *testing.T) {
	qp, err := buildQueryParams(map[string]interface{}{
		"created_at_min": "2026-06-16",
		"created_at_max": "2026-07-16",
	})
	if err != nil {
		t.Fatalf("el formulario manda YYYY-MM-DD y debe aceptarse: %v", err)
	}
	if qp.After == nil || qp.Before == nil {
		t.Fatal("el rango de fechas se perdio: la sincronizacion traeria TODO en vez del periodo elegido")
	}
	if qp.After.Year() != 2026 || qp.After.Month() != time.June || qp.After.Day() != 16 {
		t.Fatalf("created_at_min mal parseado: %v", qp.After)
	}
	if qp.Before.Day() != 16 || qp.Before.Month() != time.July {
		t.Fatalf("created_at_max mal parseado: %v", qp.Before)
	}
}

func TestBuildQueryParamsAceptaRFC3339(t *testing.T) {
	qp, err := buildQueryParams(map[string]interface{}{
		"created_at_min": "2026-06-16T00:00:00Z",
	})
	if err != nil {
		t.Fatalf("RFC3339 debe seguir funcionando: %v", err)
	}
	if qp.After == nil {
		t.Fatal("created_at_min en RFC3339 se perdio")
	}
	if qp.Before == nil {
		t.Fatal("con solo fecha inicial, la final debe completarse con ahora")
	}
}

func TestBuildQueryParamsFallaConFechaInvalida(t *testing.T) {
	_, err := buildQueryParams(map[string]interface{}{
		"created_at_min": "16/06/2026",
	})
	if err == nil {
		t.Fatal("una fecha invalida debe fallar visible, nunca ignorarse en silencio")
	}
}

func TestBuildQueryParamsSinFechasNoFalla(t *testing.T) {
	qp, err := buildQueryParams(map[string]interface{}{})
	if err != nil {
		t.Fatalf("sin fechas no debe fallar: %v", err)
	}
	if qp.After != nil || qp.Before != nil {
		t.Fatal("sin fechas no se debe inventar un rango")
	}
	if len(qp.Statuses) == 0 {
		t.Fatal("deben quedar los estados por defecto")
	}
}
