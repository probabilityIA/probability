package usecases

import (
	"math"
	"testing"
)

func casiIgual(a, b float64) bool {
	return math.Abs(a-b) < 0.0001
}

func TestNormalizeWeightToKg(t *testing.T) {
	casos := []struct {
		nombre   string
		peso     float64
		unidad   string
		esperado float64
		ok       bool
	}{
		{"kg se queda igual", 1.5, "kg", 1.5, true},
		{"gramos a kg", 1500, "g", 1.5, true},
		{"gramos abreviado gr", 500, "gr", 0.5, true},
		{"libras a kg", 1, "lb", 0.45359237, true},
		{"onzas a kg", 16, "oz", 0.45359237, true},
		{"mayusculas y espacios", 1500, "  G  ", 1.5, true},
		{"unidad vacia NO asume kg", 1.5, "", 0, false},
		{"unidad desconocida no adivina", 1.5, "piedras", 0, false},
		{"peso cero se ignora", 0, "kg", 0, false},
		{"peso negativo se ignora", -3, "kg", 0, false},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			got, ok := normalizeWeightToKg(c.peso, c.unidad)
			if ok != c.ok {
				t.Fatalf("ok = %v, se esperaba %v (peso=%v unidad=%q)", ok, c.ok, c.peso, c.unidad)
			}
			if ok && !casiIgual(got, c.esperado) {
				t.Fatalf("peso = %v, se esperaba %v", got, c.esperado)
			}
		})
	}
}

func TestConvertKgToStoreUnit(t *testing.T) {
	got, ok := convertKgToStoreUnit(1.5, "g")
	if !ok || !casiIgual(got, 1500) {
		t.Fatalf("1.5 kg a gramos = %v (ok=%v), se esperaba 1500", got, ok)
	}

	got, ok = convertKgToStoreUnit(0.45359237, "lb")
	if !ok || !casiIgual(got, 1) {
		t.Fatalf("0.4536 kg a libras = %v (ok=%v), se esperaba 1", got, ok)
	}

	if _, ok := convertKgToStoreUnit(1.5, ""); ok {
		t.Fatal("sin unidad de tienda no debe convertir: seria adivinar")
	}

	if _, ok := convertKgToStoreUnit(1.5, "piedras"); ok {
		t.Fatal("unidad desconocida no debe convertir")
	}
}

func TestRoundTripPesoNoSeDeforma(t *testing.T) {
	for _, unidad := range []string{"kg", "g", "lb", "oz"} {
		enTienda, ok := convertKgToStoreUnit(2.5, unidad)
		if !ok {
			t.Fatalf("no convirtio a %s", unidad)
		}
		devuelta, ok := normalizeWeightToKg(enTienda, unidad)
		if !ok || !casiIgual(devuelta, 2.5) {
			t.Fatalf("round trip en %s: 2.5 -> %v -> %v", unidad, enTienda, devuelta)
		}
	}
}

func TestPositive(t *testing.T) {
	if positive(0) != nil {
		t.Fatal("0 debe ser nil para no escribir dimensiones vacias")
	}
	if positive(-1) != nil {
		t.Fatal("negativo debe ser nil")
	}
	v := positive(12.5)
	if v == nil || *v != 12.5 {
		t.Fatal("valor positivo debe pasar")
	}
}
