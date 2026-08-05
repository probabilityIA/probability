package app

import (
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/dtos"
	invoicingerrors "github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCurrencyCOPValidator(t *testing.T) {
	cases := []struct {
		name     string
		currency string
		wantErr  bool
	}{
		{name: "COP pasa", currency: "COP", wantErr: false},
		{name: "cop minuscula pasa", currency: "cop", wantErr: false},
		{name: "con espacios pasa", currency: "  COP  ", wantErr: false},
		{name: "vacia pasa por defecto", currency: "", wantErr: false},
		{name: "solo espacios pasa", currency: "   ", wantErr: false},
		{name: "USD se rechaza", currency: "USD", wantErr: true},
		{name: "EUR se rechaza", currency: "EUR", wantErr: true},
		{name: "usd minuscula se rechaza", currency: "usd", wantErr: true},
	}

	v := &CurrencyCOPValidator{}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := v.Validate(&dtos.OrderData{Currency: tc.currency})

			if tc.wantErr {
				assert.ErrorIs(t, err, invoicingerrors.ErrCurrencyNotAllowed)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestPaymentStatusValidator_CODConPermiso_SeFactura(t *testing.T) {
	v := &PaymentStatusValidator{RequiredStatus: "paid", AllowCOD: true}

	err := v.Validate(&dtos.OrderData{IsPaid: false, IsCOD: true})

	assert.NoError(t, err)
}

func TestPaymentStatusValidator_CODSinPermiso_SeRechaza(t *testing.T) {
	v := &PaymentStatusValidator{RequiredStatus: "paid", AllowCOD: false}

	err := v.Validate(&dtos.OrderData{IsPaid: false, IsCOD: true})

	assert.ErrorIs(t, err, invoicingerrors.ErrOrderNotPaid)
}

func TestPaymentStatusValidator_NoPagadaNiCOD_SeRechaza(t *testing.T) {
	v := &PaymentStatusValidator{RequiredStatus: "paid", AllowCOD: true}

	err := v.Validate(&dtos.OrderData{IsPaid: false, IsCOD: false})

	assert.ErrorIs(t, err, invoicingerrors.ErrOrderNotPaid)
}

func TestPaymentStatusValidator_Pagada_SiemprePasa(t *testing.T) {
	for _, allowCOD := range []bool{true, false} {
		v := &PaymentStatusValidator{RequiredStatus: "paid", AllowCOD: allowCOD}

		assert.NoError(t, v.Validate(&dtos.OrderData{IsPaid: true}))
	}
}

func TestPaymentStatusValidator_EstadoDistintoDePaid_NoValidaNada(t *testing.T) {
	v := &PaymentStatusValidator{RequiredStatus: "pending"}

	err := v.Validate(&dtos.OrderData{IsPaid: false, IsCOD: false})

	assert.NoError(t, err)
}

func TestDateRangeValidator_SinLimites_Pasa(t *testing.T) {
	v := &DateRangeValidator{}

	assert.NoError(t, v.Validate(&dtos.OrderData{}))
}

func TestDateRangeValidator_ConLimites_HoyPasaSiempre(t *testing.T) {
	inicio := "2026-01-01"
	fin := "2026-12-31"

	cases := []*DateRangeValidator{
		{StartDate: &inicio},
		{EndDate: &fin},
		{StartDate: &inicio, EndDate: &fin},
	}

	for _, v := range cases {
		err := v.Validate(&dtos.OrderData{})
		require.NoError(t, err, "la validacion de fechas esta pendiente de implementar y hoy siempre pasa")
	}
}
