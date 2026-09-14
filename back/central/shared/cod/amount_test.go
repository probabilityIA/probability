package cod

import "testing"

func TestCustomerCharge(t *testing.T) {
	tests := []struct {
		name       string
		order      Order
		carrierFee float64
		want       float64
	}{
		{"orden manual: cod_total neto, se suma la comision real (VIG-0069)", Order{CodTotal: 135929}, 9451, 145380},
		{"orden manual sin guia todavia", Order{CodTotal: 73367}, 0, 73367},
		{"checkout del plugin: cod_total neto + comision del checkout (15789)", Order{CodTotal: 68002, IncludesShipping: true, CheckoutCarrierFee: 5365}, 5365, 73367},
		{"checkout del plugin con comision real distinta: manda la del checkout (15788)", Order{CodTotal: 65990, IncludesShipping: true, CheckoutCarrierFee: 5238}, 6116, 71228},
		{"canal que cobra el total en su checkout sin plugin", Order{CodTotal: 191322, IncludesShipping: true}, 6116, 191322},
		{"sin cod_total no se inventa un monto", Order{CheckoutCarrierFee: 5365}, 5365, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CustomerCharge(tt.order, tt.carrierFee); got != tt.want {
				t.Errorf("CustomerCharge() = %v, se esperaba %v", got, tt.want)
			}
		})
	}
}

func TestCheckoutTotal(t *testing.T) {
	if got := CheckoutTotal(Order{CodTotal: 68002, CheckoutCarrierFee: 5365}); got != 73367 {
		t.Errorf("plugin: CheckoutTotal() = %v, se esperaba 73367", got)
	}
	if got := CheckoutTotal(Order{CodTotal: 135929}); got != 135929 {
		t.Errorf("manual: CheckoutTotal() = %v, se esperaba 135929", got)
	}
	if got := CheckoutTotal(Order{CheckoutCarrierFee: 5365}); got != 0 {
		t.Errorf("sin cod_total: CheckoutTotal() = %v, se esperaba 0", got)
	}
}

func TestSummarize(t *testing.T) {
	tests := []struct {
		name       string
		order      Order
		carrierFee float64
		want       Breakdown
	}{
		{"orden manual con guia (VIG-0161)", Order{CodTotal: 106565}, 7805, Breakdown{CustomerCharge: 114370, CheckoutTotal: 106565, CarrierFee: 7805, CarrierFeeSource: FeeSourceCarrier, ChargedCarrierFee: 7805, BusinessNet: 106565}},
		{"orden manual sin guia", Order{CodTotal: 106565}, 0, Breakdown{CustomerCharge: 106565, CheckoutTotal: 106565, BusinessNet: 106565}},
		{"plugin sin guia: comision del checkout (15791)", Order{CodTotal: 65990, IncludesShipping: true, CheckoutCarrierFee: 5238}, 0, Breakdown{CustomerCharge: 71228, CheckoutTotal: 71228, CarrierFee: 5238, CarrierFeeSource: FeeSourceCheckout, ChargedCarrierFee: 5238, BusinessNet: 65990}},
		{"plugin con comision real mayor: el negocio recibe su neto (15788)", Order{CodTotal: 65990, IncludesShipping: true, CheckoutCarrierFee: 5238}, 6116, Breakdown{CustomerCharge: 71228, CheckoutTotal: 71228, CarrierFee: 6116, CarrierFeeSource: FeeSourceCarrier, ChargedCarrierFee: 5238, BusinessNet: 65990}},
		{"canal que cobra el total sin plugin", Order{CodTotal: 191322, IncludesShipping: true}, 6116, Breakdown{CustomerCharge: 191322, CheckoutTotal: 191322, CarrierFee: 6116, CarrierFeeSource: FeeSourceCarrier, BusinessNet: 185206}},
		{"sin cod_total", Order{CheckoutCarrierFee: 5365}, 5365, Breakdown{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Summarize(tt.order, tt.carrierFee); got != tt.want {
				t.Errorf("Summarize() = %+v, se esperaba %+v", got, tt.want)
			}
		})
	}
}

func TestAmountToCollect(t *testing.T) {
	tests := []struct {
		name       string
		order      Order
		totalCost  float64
		carrierFee float64
		want       float64
	}{
		{"no cod", Order{TotalAmount: 157500}, 18993.40, 6116, 0},
		{"flete aparte suma flete final y comision", Order{TotalAmount: 157500, CodTotal: 174918.40}, 18993.40, 6116, 182610},
		{"total ya incluye flete y comision: no se vuelve a sumar", Order{TotalAmount: 157500, CodTotal: 176493.40, IncludesShipping: true}, 18993.40, 6116, 176494},
		{"checkout del plugin: declara el total de la tienda (15789)", Order{TotalAmount: 45000, CodTotal: 68002, IncludesShipping: true, CheckoutCarrierFee: 5365}, 23001, 5365, 73367},
		{"checkout del plugin con comision real distinta: declara el total de la tienda (15788)", Order{TotalAmount: 45000, CodTotal: 65990, IncludesShipping: true, CheckoutCarrierFee: 5238}, 20225, 6116, 71228},
		{"sin flete calculado conserva el cod_total", Order{TotalAmount: 157500, CodTotal: 157500}, 0, 6116, 163616},
		{"sin comision de recaudo, redondea al peso", Order{TotalAmount: 157500, CodTotal: 157500}, 18993.40, 0, 176494},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AmountToCollect(tt.order, tt.totalCost, tt.carrierFee); got != tt.want {
				t.Fatalf("AmountToCollect() = %v, se esperaba %v", got, tt.want)
			}
		})
	}
}
