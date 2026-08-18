package usecaseshipment

import (
	"context"
	"errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type marginReaderStub struct {
	margin      domain.ShippingMargin
	err         error
	calls       int
	lastBiz     uint
	lastCarrier string
}

func (s *marginReaderStub) Get(ctx context.Context, businessID uint, carrierCode string) (domain.ShippingMargin, error) {
	s.calls++
	s.lastBiz = businessID
	s.lastCarrier = carrierCode
	return s.margin, s.err
}

func f64(v float64) *float64 {
	return &v
}

func s(v string) *string {
	return &v
}

func newUseCaseForCost(repo domain.IRepository, reader domain.IShippingMarginReader) *UseCaseShipment {
	return &UseCaseShipment{repo: repo, marginReader: reader}
}

func repoWithBusiness(businessID uint, err error) *mocks.RepositoryMock {
	return &mocks.RepositoryMock{
		GetOrderBusinessIDFn: func(ctx context.Context, orderUUID string) (uint, error) {
			return businessID, err
		},
	}
}

func baseShipment() *domain.Shipment {
	return &domain.Shipment{
		OrderID:     s("order-uuid"),
		TotalCost:   f64(20000),
		CarrierCode: s("servientrega"),
	}
}

func TestApplyCarrierCost_DescuentaMargenDelCostoTotal(t *testing.T) {
	reader := &marginReaderStub{margin: domain.ShippingMargin{MarginAmount: 3000, InsuranceMargin: 500}}
	uc := newUseCaseForCost(repoWithBusiness(36, nil), reader)
	shipment := baseShipment()

	uc.applyCarrierCost(context.Background(), shipment)

	require.NotNil(t, shipment.CarrierCost)
	require.NotNil(t, shipment.AppliedMargin)
	assert.InDelta(t, 16500.0, *shipment.CarrierCost, 0.001)
	assert.InDelta(t, 3500.0, *shipment.AppliedMargin, 0.001)
	assert.Equal(t, uint(36), reader.lastBiz)
	assert.Equal(t, "servientrega", reader.lastCarrier)
}

func TestApplyCarrierCost_ConCodAplicaMargenPorcentual(t *testing.T) {
	reader := &marginReaderStub{margin: domain.ShippingMargin{
		MarginAmount:     2000,
		InsuranceMargin:  0,
		CODMarginPercent: 10,
	}}
	uc := newUseCaseForCost(repoWithBusiness(36, nil), reader)
	shipment := baseShipment()
	shipment.CodCarrierFee = f64(5000)

	uc.applyCarrierCost(context.Background(), shipment)

	require.NotNil(t, shipment.CarrierCost)
	require.NotNil(t, shipment.CodProbabilityMargin)
	assert.InDelta(t, 500.0, *shipment.CodProbabilityMargin, 0.001)
	assert.InDelta(t, 2500.0, *shipment.AppliedMargin, 0.001)
	assert.InDelta(t, 17500.0, *shipment.CarrierCost, 0.001)
}

func TestApplyCarrierCost_SinCodNoAsignaMargenCod(t *testing.T) {
	reader := &marginReaderStub{margin: domain.ShippingMargin{MarginAmount: 1000, CODMarginPercent: 10}}
	uc := newUseCaseForCost(repoWithBusiness(36, nil), reader)
	shipment := baseShipment()

	uc.applyCarrierCost(context.Background(), shipment)

	assert.Nil(t, shipment.CodProbabilityMargin)
	assert.InDelta(t, 1000.0, *shipment.AppliedMargin, 0.001)
}

func TestApplyCarrierCost_MargenMayorAlCosto_NoDejaNegativo(t *testing.T) {
	reader := &marginReaderStub{margin: domain.ShippingMargin{MarginAmount: 50000}}
	uc := newUseCaseForCost(repoWithBusiness(36, nil), reader)
	shipment := baseShipment()

	uc.applyCarrierCost(context.Background(), shipment)

	require.NotNil(t, shipment.CarrierCost)
	assert.InDelta(t, 0.0, *shipment.CarrierCost, 0.001)
	assert.GreaterOrEqual(t, *shipment.CarrierCost, 0.0)
}

func TestApplyCarrierCost_UsaCarrierCuandoNoHayCarrierCode(t *testing.T) {
	reader := &marginReaderStub{margin: domain.ShippingMargin{MarginAmount: 1000}}
	uc := newUseCaseForCost(repoWithBusiness(36, nil), reader)
	shipment := baseShipment()
	shipment.CarrierCode = nil
	shipment.Carrier = s("  Interrapidisimo  ")

	uc.applyCarrierCost(context.Background(), shipment)

	assert.Equal(t, "interrapidisimo", reader.lastCarrier)
	require.NotNil(t, shipment.CarrierCost)
}

func TestApplyCarrierCost_NoAplicaCuando(t *testing.T) {
	cases := []struct {
		name     string
		mutate   func(*domain.Shipment)
		repo     domain.IRepository
		reader   domain.IShippingMarginReader
		wantCall bool
	}{
		{
			name:   "total cost nulo",
			mutate: func(sh *domain.Shipment) { sh.TotalCost = nil },
		},
		{
			name:   "total cost cero",
			mutate: func(sh *domain.Shipment) { sh.TotalCost = f64(0) },
		},
		{
			name:   "total cost negativo",
			mutate: func(sh *domain.Shipment) { sh.TotalCost = f64(-100) },
		},
		{
			name:   "carrier cost ya calculado",
			mutate: func(sh *domain.Shipment) { sh.CarrierCost = f64(999) },
		},
		{
			name:   "sin order id",
			mutate: func(sh *domain.Shipment) { sh.OrderID = nil },
		},
		{
			name:   "order id vacio",
			mutate: func(sh *domain.Shipment) { sh.OrderID = s("") },
		},
		{
			name: "sin carrier ni carrier code",
			mutate: func(sh *domain.Shipment) {
				sh.CarrierCode = nil
				sh.Carrier = nil
			},
		},
		{
			name:   "carrier code solo espacios",
			mutate: func(sh *domain.Shipment) { sh.CarrierCode = s("   ") },
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			reader := &marginReaderStub{margin: domain.ShippingMargin{MarginAmount: 1000}}
			uc := newUseCaseForCost(repoWithBusiness(36, nil), reader)
			shipment := baseShipment()
			original := shipment.CarrierCost
			tc.mutate(shipment)

			uc.applyCarrierCost(context.Background(), shipment)

			if tc.name == "carrier cost ya calculado" {
				assert.InDelta(t, 999.0, *shipment.CarrierCost, 0.001)
			} else {
				assert.Equal(t, original, shipment.CarrierCost)
			}
			assert.Equal(t, 0, reader.calls)
			assert.Nil(t, shipment.AppliedMargin)
		})
	}
}

func TestApplyCarrierCost_SinMarginReader_NoHaceNada(t *testing.T) {
	uc := newUseCaseForCost(repoWithBusiness(36, nil), nil)
	shipment := baseShipment()

	uc.applyCarrierCost(context.Background(), shipment)

	assert.Nil(t, shipment.CarrierCost)
	assert.Nil(t, shipment.AppliedMargin)
}

func TestApplyCarrierCost_NegocioNoResuelto_NoAplica(t *testing.T) {
	cases := []struct {
		name       string
		businessID uint
		err        error
	}{
		{name: "error consultando negocio", businessID: 0, err: errors.New("db down")},
		{name: "negocio cero", businessID: 0, err: nil},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			reader := &marginReaderStub{margin: domain.ShippingMargin{MarginAmount: 1000}}
			uc := newUseCaseForCost(repoWithBusiness(tc.businessID, tc.err), reader)
			shipment := baseShipment()

			uc.applyCarrierCost(context.Background(), shipment)

			assert.Equal(t, 0, reader.calls)
			assert.Nil(t, shipment.CarrierCost)
		})
	}
}

func TestApplyCarrierCost_FallaLeerMargen_NoAplica(t *testing.T) {
	reader := &marginReaderStub{err: errors.New("margin service down")}
	uc := newUseCaseForCost(repoWithBusiness(36, nil), reader)
	shipment := baseShipment()

	uc.applyCarrierCost(context.Background(), shipment)

	assert.Equal(t, 1, reader.calls)
	assert.Nil(t, shipment.CarrierCost)
	assert.Nil(t, shipment.AppliedMargin)
}

func TestApplyCarrierCost_SinMargen_CostoCarrierIgualAlTotal(t *testing.T) {
	reader := &marginReaderStub{margin: domain.ShippingMargin{}}
	uc := newUseCaseForCost(repoWithBusiness(36, nil), reader)
	shipment := baseShipment()

	uc.applyCarrierCost(context.Background(), shipment)

	require.NotNil(t, shipment.CarrierCost)
	assert.InDelta(t, 20000.0, *shipment.CarrierCost, 0.001)
	assert.InDelta(t, 0.0, *shipment.AppliedMargin, 0.001)
}

func TestCarrierLogoKey(t *testing.T) {
	cases := map[string]string{
		"servientrega":     "SERVIENTREGA",
		"Inter Rapidisimo": "INTERRAPIDISIMO",
		"dhl-express":      "DHLEXPRESS",
		"mi_paquete":       "MIPAQUETE",
		"  ":               "",
		"":                 "",
	}

	for input, want := range cases {
		t.Run(input, func(t *testing.T) {
			assert.Equal(t, want, carrierLogoKey(input))
		})
	}
}

func TestGetCarrierLogoBytes_CarrierDesconocido_RetornaNil(t *testing.T) {
	assert.Nil(t, getCarrierLogoBytes("transportadora-inventada"))
	assert.Nil(t, getCarrierLogoBytes(""))
	assert.Nil(t, getCarrierLogoBytes("   "))
}
