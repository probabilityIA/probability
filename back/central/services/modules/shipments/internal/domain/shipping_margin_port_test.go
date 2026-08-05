package domain

import "testing"

func TestShippingMarginCODMargin(t *testing.T) {
	tests := []struct {
		name   string
		margin ShippingMargin
		fee    float64
		want   float64
	}{
		{
			name:   "monto fijo tiene prioridad sobre el porcentaje",
			margin: ShippingMargin{CODMarginPercent: 15, CODMarginAmount: 2800},
			fee:    10326,
			want:   2800,
		},
		{
			name:   "monto fijo aplica aunque no haya comision del carrier",
			margin: ShippingMargin{CODMarginAmount: 2800},
			fee:    0,
			want:   2800,
		},
		{
			name:   "sin monto fijo se conserva el porcentaje",
			margin: ShippingMargin{CODMarginPercent: 15},
			fee:    10326,
			want:   1548.9,
		},
		{
			name:   "sin comision del carrier el porcentaje no aplica",
			margin: ShippingMargin{CODMarginPercent: 15},
			fee:    0,
			want:   0,
		},
		{
			name:   "sin configuracion no hay margen",
			margin: ShippingMargin{},
			fee:    10326,
			want:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.margin.CODMargin(tt.fee); got != tt.want {
				t.Fatalf("CODMargin(%v) = %v, se esperaba %v", tt.fee, got, tt.want)
			}
		})
	}
}
