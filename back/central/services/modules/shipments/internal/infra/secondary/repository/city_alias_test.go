package repository

import "testing"

func TestAliasCityDane(t *testing.T) {
	casos := map[string]string{
		"Cartagena":  "13001",
		"CARTAGENA":  "13001",
		"cartagena ": "13001",
		"Buga":       "76111",
		"BUGA":       "76111",
		"Suba":       "",
		"Bo":         "",
		"casa":       "",
		"":           "",
	}

	for entrada, esperado := range casos {
		if got := aliasCityDane(entrada); got != esperado {
			t.Errorf("para %q se esperaba %q, llego %q", entrada, esperado, got)
		}
	}
}

func TestNormalizeCityKeyQuitaTildesYEspacios(t *testing.T) {
	casos := map[string]string{
		"  Medellín ":  "medellin",
		"BOGOTÁ, D.C.": "bogota dc",
		"Facatativá":   "facatativa",
		"Peñol":        "penol",
	}

	for entrada, esperado := range casos {
		if got := normalizeCityKey(entrada); got != esperado {
			t.Errorf("para %q se esperaba %q, llego %q", entrada, esperado, got)
		}
	}
}
