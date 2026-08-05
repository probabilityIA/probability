package domain

import (
	"math"
	"testing"
)

func interFee(declared float64) float64 {
	return math.Floor((840 + 0.05*declared) * 1.19)
}

func coordFee(declared float64) float64 {
	fee := math.Floor((840 + 0.02*declared) * 1.19)
	if fee < 6116 {
		return 6116
	}
	return fee
}

func TestSolveCODDeclaredValueDejaElNetoExacto(t *testing.T) {
	tests := []struct {
		name      string
		netTarget float64
		feeOf     func(float64) float64
	}{
		{"interrapidisimo VIG-0068", 159300.90, interFee},
		{"interrapidisimo VIG-0060", 198291.25, interFee},
		{"interrapidisimo VIG-0066", 65050.75, interFee},
		{"coordinadora VIG-0062", 767108.00, coordFee},
		{"coordinadora VIG-0067", 255307.70, coordFee},
		{"coordinadora en el minimo VIG-0065", 171948.40, coordFee},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fee0 := tt.feeOf(tt.netTarget)
			declared1 := tt.netTarget + fee0
			fee1 := tt.feeOf(declared1)

			declared, ok := SolveCODDeclaredValue(
				CODQuotePoint{Declared: tt.netTarget, Fee: fee0},
				CODQuotePoint{Declared: declared1, Fee: fee1},
				tt.netTarget,
			)
			if !ok {
				t.Fatalf("SolveCODDeclaredValue() no resolvio")
			}

			recibe := declared - tt.feeOf(declared)
			if recibe < tt.netTarget {
				t.Fatalf("el negocio recibe %.2f y esperaba al menos %.2f (declarado %.2f)", recibe, tt.netTarget, declared)
			}
			if recibe-math.Ceil(tt.netTarget) > 2 {
				t.Fatalf("el negocio recibe %.2f, se paso de %.2f por mas de 2 pesos (declarado %.2f)", recibe, math.Ceil(tt.netTarget), declared)
			}
		})
	}
}

func TestSolveCODDeclaredValueComisionFijaSumaDirecto(t *testing.T) {
	declared, ok := SolveCODDeclaredValue(
		CODQuotePoint{Declared: 100000, Fee: 6116},
		CODQuotePoint{Declared: 106116, Fee: 6116},
		100000,
	)
	if !ok {
		t.Fatalf("SolveCODDeclaredValue() no resolvio con comision fija")
	}
	if declared != 106116 {
		t.Fatalf("declared = %v, se esperaba 106116", declared)
	}
}

func TestSolveCODDeclaredValueRechazaTarifasImposibles(t *testing.T) {
	if _, ok := SolveCODDeclaredValue(
		CODQuotePoint{Declared: 100000, Fee: 1000},
		CODQuotePoint{Declared: 200000, Fee: 150000},
		100000,
	); ok {
		t.Fatal("una tarifa mayor al 100% debio rechazarse")
	}

	if _, ok := SolveCODDeclaredValue(
		CODQuotePoint{Declared: 100000, Fee: 5000},
		CODQuotePoint{Declared: 200000, Fee: 1000},
		100000,
	); ok {
		t.Fatal("una tarifa negativa debio rechazarse")
	}

	if _, ok := SolveCODDeclaredValue(
		CODQuotePoint{Declared: 100000, Fee: 1000},
		CODQuotePoint{Declared: 200000, Fee: 2000},
		0,
	); ok {
		t.Fatal("un neto objetivo de cero debio rechazarse")
	}
}

func TestFindCODFeePrefiereElCarrierDelEnvio(t *testing.T) {
	rates := []Rate{
		{Carrier: "COORDINADORA", CODDetails: &CODDetails{CODCost: 6116}},
		{Carrier: "INTERRAPIDISIMO", CODDetails: &CODDetails{CODCost: 11092}},
	}

	fee, ok := FindCODFee(rates, "interrapidisimo")
	if !ok || fee != 11092 {
		t.Fatalf("FindCODFee() = %v, %v; se esperaba 11092, true", fee, ok)
	}

	if _, ok := FindCODFee([]Rate{{Carrier: "ENVIA"}}, "ENVIA"); ok {
		t.Fatal("un rate sin codDetails no debe reportar comision")
	}
}
