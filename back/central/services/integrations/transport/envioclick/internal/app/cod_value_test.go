package app

import (
	"context"
	"math"
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/transport/envioclick/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/transport/envioclick/internal/mocks"
)

const testCODCarrier = "INTERRAPIDISIMO"

func quoteMockConComisionProporcional(porcentaje float64, llamadas *[]float64) *mocks.EnvioClickClientMock {
	return &mocks.EnvioClickClientMock{
		QuoteFn: func(baseURL, apiKey string, req domain.QuoteRequest) (*domain.QuoteResponse, error) {
			if llamadas != nil {
				*llamadas = append(*llamadas, req.CODValue)
			}
			fee := math.Ceil(req.CODValue * porcentaje)
			return &domain.QuoteResponse{
				Status: "success",
				Data: domain.QuoteData{
					Rates: []domain.Rate{
						{IDRate: 1, Carrier: testCODCarrier, COD: false},
						{
							IDRate:     2,
							Carrier:    testCODCarrier,
							COD:        true,
							CODDetails: &domain.CODDetails{CODCost: fee, CODCostWithInsurance: fee},
						},
					},
				},
			}, nil
		},
	}
}

func TestResolveCODValue_DevuelveLaComisionDelValorFinalmenteDeclarado(t *testing.T) {
	ctx := context.Background()

	var declarados []float64
	uc := New(quoteMockConComisionProporcional(0.07, &declarados), &mocks.LoggerMock{})

	netTarget := 91241.0

	declared, fee, ok := uc.ResolveCODValue(ctx, testBaseURL, testAPIKey, domain.QuoteRequest{CODValue: netTarget}, testCODCarrier, 0, netTarget, nil)

	if !ok {
		t.Fatal("se esperaba que la calibracion resolviera")
	}
	if declared <= netTarget {
		t.Fatalf("el valor declarado debe superar el objetivo, se obtuvo: %v", declared)
	}
	if declared-fee < netTarget {
		t.Errorf("declarado %v menos comision %v deja %v, por debajo del objetivo %v", declared, fee, declared-fee, netTarget)
	}

	esperada := math.Ceil(declared * 0.07)
	if fee != esperada {
		t.Errorf("se esperaba la comision medida sobre %v (%v), se obtuvo: %v", declared, esperada, fee)
	}

	comisionDelObjetivo := math.Ceil(netTarget * 0.07)
	if fee == comisionDelObjetivo {
		t.Errorf("la comision devuelta es la del valor base (%v), no la del valor declarado", comisionDelObjetivo)
	}

	if len(declarados) < 3 {
		t.Errorf("se esperaban al menos 3 sondeos, se obtuvieron: %d", len(declarados))
	}
	if len(declarados) > 0 && declarados[len(declarados)-1] != declared {
		t.Errorf("el ultimo sondeo debe ser sobre el valor declarado %v, fue sobre %v", declared, declarados[len(declarados)-1])
	}
}

func TestResolveCODValue_ComisionFijaDevuelveLaMismaComision(t *testing.T) {
	ctx := context.Background()

	const comisionFija = 6116.0

	client := &mocks.EnvioClickClientMock{
		QuoteFn: func(baseURL, apiKey string, req domain.QuoteRequest) (*domain.QuoteResponse, error) {
			return &domain.QuoteResponse{
				Status: "success",
				Data: domain.QuoteData{
					Rates: []domain.Rate{{
						IDRate:     2,
						Carrier:    testCODCarrier,
						COD:        true,
						CODDetails: &domain.CODDetails{CODCost: comisionFija, CODCostWithInsurance: comisionFija},
					}},
				},
			}, nil
		},
	}

	uc := New(client, &mocks.LoggerMock{})

	netTarget := 112267.0

	declared, fee, ok := uc.ResolveCODValue(ctx, testBaseURL, testAPIKey, domain.QuoteRequest{CODValue: netTarget}, testCODCarrier, 0, netTarget, nil)

	if !ok {
		t.Fatal("se esperaba que la calibracion resolviera")
	}
	if fee != comisionFija {
		t.Errorf("se esperaba comision %v, se obtuvo: %v", comisionFija, fee)
	}
	if declared != netTarget+comisionFija {
		t.Errorf("se esperaba declarado %v, se obtuvo: %v", netTarget+comisionFija, declared)
	}
}

func TestResolveCODValue_SinComisionEnLaCotizacionNoDevuelveNada(t *testing.T) {
	ctx := context.Background()

	client := &mocks.EnvioClickClientMock{
		QuoteFn: func(baseURL, apiKey string, req domain.QuoteRequest) (*domain.QuoteResponse, error) {
			return &domain.QuoteResponse{
				Status: "success",
				Data: domain.QuoteData{
					Rates: []domain.Rate{{IDRate: 1, Carrier: testCODCarrier, COD: false}},
				},
			}, nil
		},
	}

	uc := New(client, &mocks.LoggerMock{})

	declared, fee, ok := uc.ResolveCODValue(ctx, testBaseURL, testAPIKey, domain.QuoteRequest{CODValue: 90000}, testCODCarrier, 0, 90000, nil)

	if ok {
		t.Error("sin comision en la cotizacion no se debe calibrar")
	}
	if declared != 0 || fee != 0 {
		t.Errorf("se esperaban ceros, se obtuvo declarado=%v comision=%v", declared, fee)
	}
}
