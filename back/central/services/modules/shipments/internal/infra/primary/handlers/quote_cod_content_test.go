package handlers

import "testing"

func TestAlignCODContentValueIgualaElDeclarado(t *testing.T) {
	raw := map[string]interface{}{"codValue": 60261.0, "contentValue": 45000.0}
	alignCODContentValue(raw)

	if raw["contentValue"] != 60261.0 {
		t.Fatalf("se esperaba contentValue 60261, quedo %v", raw["contentValue"])
	}
}

func TestAlignCODContentValueNoTocaLoQueNoEsCOD(t *testing.T) {
	raw := map[string]interface{}{"contentValue": 45000.0}
	alignCODContentValue(raw)

	if raw["contentValue"] != 45000.0 {
		t.Fatalf("sin codValue no se debe tocar contentValue, quedo %v", raw["contentValue"])
	}

	raw = map[string]interface{}{"codValue": 0.0, "contentValue": 45000.0}
	alignCODContentValue(raw)

	if raw["contentValue"] != 45000.0 {
		t.Fatalf("con codValue en cero no se debe tocar contentValue, quedo %v", raw["contentValue"])
	}
}
