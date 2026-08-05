package usecasecreateorder

import (
	"context"
	"errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func codDTO(isCod bool, codTotal *float64, paymentMethodID uint) *dtos.ProbabilityOrderDTO {
	dto := &dtos.ProbabilityOrderDTO{
		IsCod:    isCod,
		CodTotal: codTotal,
	}
	if paymentMethodID > 0 {
		dto.Payments = []dtos.ProbabilityPaymentDTO{{PaymentMethodID: paymentMethodID}}
	}
	return dto
}

func TestNormalizeCodOnDTO_SinFlagNiTotal_NoEsCod(t *testing.T) {
	dto := codDTO(false, nil, 0)

	err := normalizeCodOnDTO(dto)

	require.NoError(t, err)
	assert.False(t, dto.IsCod)
	require.NotNil(t, dto.CodTotal)
	assert.InDelta(t, 0.0, *dto.CodTotal, 0.001)
}

func TestNormalizeCodOnDTO_FlagFalsoConTotal_HoyQuedaComoCod(t *testing.T) {
	total := 150000.0
	dto := codDTO(false, &total, 0)

	err := normalizeCodOnDTO(dto)

	require.NoError(t, err)
	assert.True(t, dto.IsCod)
	require.NotNil(t, dto.CodTotal)
	assert.InDelta(t, 150000.0, *dto.CodTotal, 0.001)
}

func TestNormalizeCodOnDTO_EsCodConTotal_SeConserva(t *testing.T) {
	total := 150000.0
	dto := codDTO(true, &total, 0)

	err := normalizeCodOnDTO(dto)

	require.NoError(t, err)
	assert.True(t, dto.IsCod)
	require.NotNil(t, dto.CodTotal)
	assert.InDelta(t, 150000.0, *dto.CodTotal, 0.001)
}

func TestNormalizeCodOnDTO_EsCodSinTotal_RetornaError(t *testing.T) {
	dto := codDTO(true, nil, 0)

	err := normalizeCodOnDTO(dto)

	assert.ErrorIs(t, err, domain.ErrCodTotalRequired)
}

func TestNormalizeCodOnDTO_EsCodConTotalCero_RetornaError(t *testing.T) {
	cero := 0.0
	dto := codDTO(true, &cero, 0)

	err := normalizeCodOnDTO(dto)

	assert.ErrorIs(t, err, domain.ErrCodTotalRequired)
}

func TestNormalizeCodOnDTO_MetodoCOD_ConvierteACashYMarcaCod(t *testing.T) {
	total := 150000.0
	dto := codDTO(false, &total, domain.PaymentMethodCOD)

	err := normalizeCodOnDTO(dto)

	require.NoError(t, err)
	assert.True(t, dto.IsCod, "el metodo de pago COD fuerza la orden a contraentrega")
	require.Len(t, dto.Payments, 1)
	assert.Equal(t, domain.PaymentMethodCash, dto.Payments[0].PaymentMethodID)
	assert.InDelta(t, 150000.0, *dto.CodTotal, 0.001)
}

func TestNormalizeCodOnDTO_MetodoCODSinTotal_RetornaError(t *testing.T) {
	dto := codDTO(false, nil, domain.PaymentMethodCOD)

	err := normalizeCodOnDTO(dto)

	assert.ErrorIs(t, err, domain.ErrCodTotalRequired)
}

func TestNormalizeCodOnDTO_SinPagos_NoRompe(t *testing.T) {
	total := 50000.0
	dto := &dtos.ProbabilityOrderDTO{IsCod: true, CodTotal: &total}

	err := normalizeCodOnDTO(dto)

	require.NoError(t, err)
	assert.True(t, dto.IsCod)
	assert.Empty(t, dto.Payments)
}

func TestNormalizeCodOnDTO_MetodoDistintoNoSeToca(t *testing.T) {
	total := 50000.0
	dto := codDTO(true, &total, 99)

	err := normalizeCodOnDTO(dto)

	require.NoError(t, err)
	require.Len(t, dto.Payments, 1)
	assert.Equal(t, uint(99), dto.Payments[0].PaymentMethodID)
}

func TestResolveCodIncludesShipping_PorDefectoIncluyeEnvio(t *testing.T) {
	uc := newTestCreateUseCase(&mockRepository{}, nil, nil, nil)
	dto := &dtos.ProbabilityOrderDTO{IntegrationID: 0}

	uc.resolveCodIncludesShipping(context.Background(), dto)

	assert.True(t, dto.CodIncludesShipping, "sin integracion se asume que el total ya incluye el envio")
}

func TestResolveCodIncludesShipping_LeeLaConfiguracionDeLaIntegracion(t *testing.T) {
	var consultada uint
	repo := &mockRepository{
		GetIntegrationCodIncludesShippingFn: func(ctx context.Context, integrationID uint) (bool, error) {
			consultada = integrationID
			return false, nil
		},
	}
	uc := newTestCreateUseCase(repo, nil, nil, nil)
	dto := &dtos.ProbabilityOrderDTO{IntegrationID: 197}

	uc.resolveCodIncludesShipping(context.Background(), dto)

	assert.False(t, dto.CodIncludesShipping)
	assert.Equal(t, uint(197), consultada)
}

func TestResolveCodIncludesShipping_FallaLectura_UsaElDefaultSeguro(t *testing.T) {
	repo := &mockRepository{
		GetIntegrationCodIncludesShippingFn: func(ctx context.Context, integrationID uint) (bool, error) {
			return false, errors.New("db down")
		},
	}
	uc := newTestCreateUseCase(repo, nil, nil, nil)
	dto := &dtos.ProbabilityOrderDTO{IntegrationID: 197}

	uc.resolveCodIncludesShipping(context.Background(), dto)

	assert.True(t, dto.CodIncludesShipping, "ante un fallo se asume que el total ya incluye el envio")
}

func TestNormalizeCod_TablaDeCasos(t *testing.T) {
	f := func(v float64) *float64 { return &v }
	b := func(v bool) *bool { return &v }

	cases := []struct {
		name       string
		isCod      *bool
		codTotal   *float64
		methodID   uint
		wantCod    bool
		wantAmount float64
		wantMethod uint
		wantErr    error
	}{
		{name: "sin datos no es cod", isCod: nil, codTotal: nil, methodID: 0, wantCod: false, wantAmount: 0, wantMethod: 0},
		{name: "total positivo infiere cod", isCod: nil, codTotal: f(1000), methodID: 0, wantCod: true, wantAmount: 1000, wantMethod: 0},
		{name: "flag falso gana sobre el total", isCod: b(false), codTotal: f(1000), methodID: 0, wantCod: false, wantAmount: 0, wantMethod: 0},
		{name: "flag verdadero sin total falla", isCod: b(true), codTotal: nil, methodID: 0, wantErr: domain.ErrCodTotalRequired},
		{name: "total negativo no es cod", isCod: nil, codTotal: f(-500), methodID: 0, wantCod: false, wantAmount: 0, wantMethod: 0},
		{name: "metodo cod se convierte a efectivo", isCod: nil, codTotal: f(1000), methodID: domain.PaymentMethodCOD, wantCod: true, wantAmount: 1000, wantMethod: domain.PaymentMethodCash},
		{name: "metodo cod vence al flag falso", isCod: b(false), codTotal: f(1000), methodID: domain.PaymentMethodCOD, wantCod: true, wantAmount: 1000, wantMethod: domain.PaymentMethodCash},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cod, amount, method, err := domain.NormalizeCod(tc.isCod, tc.codTotal, tc.methodID)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.wantCod, cod)
			assert.Equal(t, tc.wantMethod, method)
			require.NotNil(t, amount)
			assert.InDelta(t, tc.wantAmount, *amount, 0.001)
		})
	}
}
