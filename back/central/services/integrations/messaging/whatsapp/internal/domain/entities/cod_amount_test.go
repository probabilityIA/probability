package entities

import "testing"

func TestCodAmountToCollect(t *testing.T) {
	tests := []struct {
		name                  string
		codTotal              float64
		codCarrierFee         float64
		codIncludesShipping   bool
		codCheckoutCarrierFee float64
		want                  float64
	}{
		{"orden manual: cod_total neto, se suma la comision", 135929, 9451, false, 0, 145380},
		{"checkout del plugin: cod_total neto + comision del checkout (orden 15789)", 68002, 5365, true, 5365, 73367},
		{"checkout del plugin con comision real distinta: manda la del checkout (orden 15788)", 65990, 6116, true, 5238, 71228},
		{"canal que cobra el total en su checkout sin plugin", 191322, 6116, true, 0, 191322},
		{"sin comision conocida todavia", 73367, 0, false, 0, 73367},
		{"sin cod_total no se inventa un monto", 0, 5365, false, 5365, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CodAmountToCollect(tt.codTotal, tt.codCarrierFee, tt.codIncludesShipping, tt.codCheckoutCarrierFee)
			if got != tt.want {
				t.Errorf("CodAmountToCollect() = %v, se esperaba %v", got, tt.want)
			}
		})
	}
}
