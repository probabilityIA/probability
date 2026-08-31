package domain

import "testing"

func TestOrderCodBasisAmountToCollect(t *testing.T) {
	tests := []struct {
		name       string
		basis      OrderCodBasis
		totalCost  float64
		carrierFee float64
		want       float64
	}{
		{
			name:       "no cod",
			basis:      OrderCodBasis{TotalAmount: 157500, CodTotal: 0, CodIncludesShipping: false},
			totalCost:  18993.40,
			carrierFee: 6116,
			want:       0,
		},
		{
			name:       "flete aparte suma flete final y comision",
			basis:      OrderCodBasis{TotalAmount: 157500, CodTotal: 174918.40, CodIncludesShipping: false},
			totalCost:  18993.40,
			carrierFee: 6116,
			want:       182610,
		},
		{
			name:       "total ya incluye flete y comision: no se vuelve a sumar",
			basis:      OrderCodBasis{TotalAmount: 157500, CodTotal: 176493.40, CodIncludesShipping: true},
			totalCost:  18993.40,
			carrierFee: 6116,
			want:       176494,
		},
		{
			name:       "sin flete calculado conserva el cod_total",
			basis:      OrderCodBasis{TotalAmount: 157500, CodTotal: 157500, CodIncludesShipping: false},
			totalCost:  0,
			carrierFee: 6116,
			want:       163616,
		},
		{
			name:       "sin comision de recaudo, redondea al peso",
			basis:      OrderCodBasis{TotalAmount: 157500, CodTotal: 157500, CodIncludesShipping: false},
			totalCost:  18993.40,
			carrierFee: 0,
			want:       176494,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.basis.AmountToCollect(tt.totalCost, tt.carrierFee)
			if got != tt.want {
				t.Fatalf("AmountToCollect() = %v, want %v", got, tt.want)
			}
		})
	}
}
