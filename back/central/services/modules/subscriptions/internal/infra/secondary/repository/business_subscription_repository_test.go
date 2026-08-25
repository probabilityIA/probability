package repository

import "testing"

func TestPlanIncludesModule(t *testing.T) {
	tests := []struct {
		name        string
		moduleCodes []string
		code        string
		want        bool
	}{
		{"incluido", []string{"orders", "inventory", "shipments"}, "inventory", true},
		{"no incluido", []string{"orders", "shipments"}, "inventory", false},
		{"lista vacia", []string{}, "inventory", false},
		{"lista nil", nil, "inventory", false},
		{"case sensitive, no matchea", []string{"Inventory"}, "inventory", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := planIncludesModule(tt.moduleCodes, tt.code)
			if got != tt.want {
				t.Errorf("planIncludesModule(%v, %q) = %v, want %v", tt.moduleCodes, tt.code, got, tt.want)
			}
		})
	}
}

func TestPlanIncludesModule_ConFeaturesSerializadas(t *testing.T) {
	features := marshalModuleCodes([]string{"orders", "inventory"})
	codes := unmarshalModuleCodes(features)

	if !planIncludesModule(codes, inventoryModuleCode) {
		t.Fatal("esperaba que el plan con inventory serializado lo detectara como incluido")
	}

	featuresSinInventario := marshalModuleCodes([]string{"orders", "customers"})
	codesSinInventario := unmarshalModuleCodes(featuresSinInventario)

	if planIncludesModule(codesSinInventario, inventoryModuleCode) {
		t.Fatal("esperaba que el plan sin inventory NO lo detectara como incluido")
	}
}
