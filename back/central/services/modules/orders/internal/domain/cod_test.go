package domain_test

import (
	"errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain"
)

func f(v float64) *float64 { return &v }
func b(v bool) *bool       { return &v }

func TestNormalizeCod(t *testing.T) {
	cases := []struct {
		name       string
		isCod      *bool
		codTotal   *float64
		method     uint
		wantCod    bool
		wantAmount float64
		wantMethod uint
		wantErr    error
	}{
		{
			name:       "contra entrega con efectivo",
			isCod:      b(true),
			codTotal:   f(198385),
			method:     domain.PaymentMethodCash,
			wantCod:    true,
			wantAmount: 198385,
			wantMethod: domain.PaymentMethodCash,
		},
		{
			name:       "contra entrega con tarjeta de debito",
			isCod:      b(true),
			codTotal:   f(50000),
			method:     2,
			wantCod:    true,
			wantAmount: 50000,
			wantMethod: 2,
		},
		{
			name:       "prepagada con tarjeta de credito",
			isCod:      b(false),
			codTotal:   f(0),
			method:     1,
			wantCod:    false,
			wantAmount: 0,
			wantMethod: 1,
		},
		{
			name:       "prepagada ignora cod_total residual",
			isCod:      b(false),
			codTotal:   f(120000),
			method:     1,
			wantCod:    false,
			wantAmount: 0,
			wantMethod: 1,
		},
		{
			name:       "sin is_cod se infiere del monto",
			codTotal:   f(75000),
			method:     1,
			wantCod:    true,
			wantAmount: 75000,
			wantMethod: 1,
		},
		{
			name:       "metodo legado contra entrega se normaliza a efectivo",
			codTotal:   f(90000),
			method:     domain.PaymentMethodCOD,
			wantCod:    true,
			wantAmount: 90000,
			wantMethod: domain.PaymentMethodCash,
		},
		{
			name:     "contra entrega sin monto es invalida",
			isCod:    b(true),
			codTotal: f(0),
			method:   domain.PaymentMethodCash,
			wantErr:  domain.ErrCodTotalRequired,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cod, amount, method, err := domain.NormalizeCod(tc.isCod, tc.codTotal, tc.method)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("error: esperado %v, obtenido %v", tc.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("error inesperado: %v", err)
			}
			if cod != tc.wantCod {
				t.Errorf("isCod: esperado %v, obtenido %v", tc.wantCod, cod)
			}
			if amount == nil || *amount != tc.wantAmount {
				t.Errorf("codTotal: esperado %v, obtenido %v", tc.wantAmount, amount)
			}
			if method != tc.wantMethod {
				t.Errorf("paymentMethodID: esperado %d, obtenido %d", tc.wantMethod, method)
			}
		})
	}
}
