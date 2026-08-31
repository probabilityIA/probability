package dtos

import "testing"

func TestInvoiceShippingCost(t *testing.T) {
	tests := []struct {
		name  string
		order *OrderData
		want  float64
	}{
		{"nil", nil, 0},
		{"prepago usa el costo de la guia", &OrderData{IsCOD: false, ShippingCost: 17799}, 17799},
		{"cod sin comision conocida no cambia", &OrderData{IsCOD: true, ShippingCost: 19031}, 19031},
		{"prepago ignora la comision", &OrderData{IsCOD: false, ShippingCost: 17799, CodCarrierFee: 6116}, 17799},
		{"orden 14665", &OrderData{IsCOD: true, ShippingCost: 12006, CodCarrierFee: 6116}, 18122},
		{"orden 14666", &OrderData{IsCOD: true, ShippingCost: 23132, CodCarrierFee: 6116}, 29248},
		{"orden 14667 interrapidisimo", &OrderData{IsCOD: true, ShippingCost: 20789, CodCarrierFee: 4445}, 25234},
		{"orden 14668", &OrderData{IsCOD: true, ShippingCost: 18901, CodCarrierFee: 6116}, 25017},
		{"orden 14670", &OrderData{IsCOD: true, ShippingCost: 15645, CodCarrierFee: 3677}, 19322},
		{"orden 14672", &OrderData{IsCOD: true, ShippingCost: 28382, CodCarrierFee: 6116}, 34498},
		{"orden 14673", &OrderData{IsCOD: true, ShippingCost: 19031, CodCarrierFee: 6116}, 25147},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InvoiceShippingCost(tt.order); got != tt.want {
				t.Errorf("InvoiceShippingCost() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInvoiceTotalAmount(t *testing.T) {
	tests := []struct {
		name  string
		order *OrderData
		want  float64
	}{
		{"nil", nil, 0},
		{"prepago no cambia", &OrderData{IsCOD: false, TotalAmount: 60000, ShippingCost: 17799}, 60000},
		{"cod sin comision no cambia", &OrderData{IsCOD: true, TotalAmount: 45000, ShippingCost: 19031}, 45000},
		{"orden 14673 suma la comision", &OrderData{IsCOD: true, TotalAmount: 45000, ShippingCost: 19031, CodCarrierFee: 6116}, 51116},
		{"shopify con iva incluido no cambia", &OrderData{IsCOD: true, TotalAmount: 211780, Subtotal: 197780, Tax: 31580.80, ShippingCost: 14000}, 211780},
		{"shopify con descuento no cambia", &OrderData{IsCOD: true, TotalAmount: 493158.94, Subtotal: 493158.94, Discount: 76911.06, ShippingCost: 11000}, 493158.94},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InvoiceTotalAmount(tt.order); got != tt.want {
				t.Errorf("InvoiceTotalAmount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrdenesSinComisionNoSeTocan(t *testing.T) {
	casos := []*OrderData{
		{IsCOD: true, TotalAmount: 211780, Subtotal: 197780, Tax: 31580.80, ShippingCost: 14000},
		{IsCOD: true, TotalAmount: 493158.94, Subtotal: 493158.94, Discount: 76911.06, ShippingCost: 11000},
		{IsCOD: true, TotalAmount: 164780, Subtotal: 164780, Discount: 15000, ShippingCost: 15000},
		{IsCOD: false, TotalAmount: 296112, Subtotal: 270000, ShippingCost: 26112},
	}

	for _, o := range casos {
		if got := InvoiceShippingCost(o); got != o.ShippingCost {
			t.Errorf("envio cambio: %v, esperado %v", got, o.ShippingCost)
		}
		if got := InvoiceTotalAmount(o); got != o.TotalAmount {
			t.Errorf("total cambio: %v, esperado %v", got, o.TotalAmount)
		}
	}
}
