package entities

import "testing"

func TestCodAmountToCollect(t *testing.T) {
	tests := []struct {
		name                string
		codTotal            float64
		codCarrierFee       float64
		codIncludesShipping bool
		want                float64
	}{
		{"orden manual: cod_total neto, se suma la comision", 135929, 9451, false, 145380},
		{"checkout woocommerce: cod_total ya trae la comision (orden 15789)", 73367, 5365, true, 73367},
		{"sin comision conocida todavia", 73367, 0, false, 73367},
		{"sin cod_total no se inventa un monto", 0, 5365, false, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CodAmountToCollect(tt.codTotal, tt.codCarrierFee, tt.codIncludesShipping)
			if got != tt.want {
				t.Errorf("CodAmountToCollect(%v, %v, %v) = %v, se esperaba %v",
					tt.codTotal, tt.codCarrierFee, tt.codIncludesShipping, got, tt.want)
			}
		})
	}
}
