package handlers

import (
	"math"
	"testing"
)

func interraRate(declared float64) map[string]interface{} {
	return map[string]interface{}{
		"idRate":               float64(25780956),
		"carrier":              "INTERRAPIDISIMO",
		"product":              "Normal",
		"cod":                  true,
		"flete":                10418.0,
		"minimumInsurance":     math.Round(declared * 0.015),
		"extraInsurance":       0.0,
		"codCarrierFee":        math.Floor((840 + 0.05*declared) * 1.19),
		"codProbabilityMargin": 1500.0,
	}
}

func coordinadoraRate(declared float64) map[string]interface{} {
	fee := math.Floor((840 + 0.02*declared) * 1.19)
	if fee < 6116 {
		fee = 6116
	}
	return map[string]interface{}{
		"idRate":               float64(25693427),
		"carrier":              "COORDINADORA",
		"product":              "Normal",
		"cod":                  true,
		"flete":                10291.0,
		"minimumInsurance":     650.0,
		"extraInsurance":       0.0,
		"codCarrierFee":        fee,
		"codProbabilityMargin": 1500.0,
	}
}

func assertCODCuadra(t *testing.T, name string, build func(float64) map[string]interface{}, contentValue float64) {
	t.Helper()

	base := []map[string]interface{}{build(contentValue)}
	probeDeclared := codProbeValue(base, contentValue, false)
	probe := []map[string]interface{}{build(probeDeclared)}

	out := calibrateCODRates(base, probe, contentValue, contentValue, probeDeclared, false)
	if len(out) != 1 {
		t.Fatalf("%s: se esperaba 1 tarifa, llegaron %d", name, len(out))
	}

	declared := contentValue + codTargetCost(out[0], false) + toFloat(out[0]["codCarrierFee"])
	real := build(declared)
	neto := declared - toFloat(real["codCarrierFee"])
	objetivo := contentValue + codTargetCost(real, false)

	if neto < objetivo {
		t.Fatalf("%s: el negocio recibe %.0f y necesita %.0f (faltan %.0f)", name, neto, objetivo, objetivo-neto)
	}
	if neto-objetivo > 10 {
		t.Fatalf("%s: se cobra de mas, neto %.0f contra objetivo %.0f", name, neto, objetivo)
	}
}

func TestCalibracionCODInterrapidisimo(t *testing.T) {
	for _, contentValue := range []float64{45000, 94900, 176500, 353000} {
		assertCODCuadra(t, "interrapidisimo", interraRate, contentValue)
	}
}

func TestCalibracionCODCoordinadora(t *testing.T) {
	for _, contentValue := range []float64{45000, 176500, 353000} {
		assertCODCuadra(t, "coordinadora", coordinadoraRate, contentValue)
	}
}

func TestCalibracionCODOrden14687(t *testing.T) {
	contentValue := 45000.0
	base := []map[string]interface{}{interraRate(contentValue)}

	if fee := toFloat(base[0]["codCarrierFee"]); fee != 3677 {
		t.Fatalf("el modelo no reproduce la comision del checkout: %.0f", fee)
	}

	probeDeclared := codProbeValue(base, contentValue, false)
	out := calibrateCODRates(base, []map[string]interface{}{interraRate(probeDeclared)}, contentValue, contentValue, probeDeclared, false)

	fee := toFloat(out[0]["codCarrierFee"])
	envio := codTargetCost(out[0], false) + fee
	if fee <= 3677 {
		t.Fatalf("la comision calibrada quedo igual o menor que la del bug: %.0f", fee)
	}
	if envio <= 16270 {
		t.Fatalf("el envio cobrado al comprador no subio: %.0f", envio)
	}
	t.Logf("checkout viejo: envio 16270 comision 3677 | calibrado: envio %.0f comision %.0f", envio, fee)
}

func TestCalibracionSinTarifaEnElSondeoNoRompe(t *testing.T) {
	base := []map[string]interface{}{interraRate(45000)}
	out := calibrateCODRates(base, []map[string]interface{}{}, 45000, 45000, 70000, false)

	if len(out) != 1 || toFloat(out[0]["codCarrierFee"]) != 3677 {
		t.Fatal("sin sondeo la tarifa debe quedar tal cual llego")
	}
}
