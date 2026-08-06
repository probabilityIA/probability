package app

import (
	"context"
	"math"

	"github.com/secamc93/probability/back/central/services/integrations/transport/envioclick/internal/domain"
)

const (
	codMaxRefinements = 3
	codProbeGap       = 10000.0
)

func (u *useCase) ResolveCODValue(ctx context.Context, baseURL, apiKey string, req domain.QuoteRequest, carrier string, netTarget float64, metas *[]domain.SyncMeta) (float64, bool) {
	if netTarget <= 0 {
		return 0, false
	}

	target := math.Ceil(netTarget)

	firstFee, ok := u.probeCODFee(ctx, baseURL, apiKey, req, carrier, target, metas)
	if !ok {
		return 0, false
	}

	secondDeclared := target + firstFee
	if secondDeclared <= target {
		secondDeclared = target + codProbeGap
	}

	secondFee, ok := u.probeCODFee(ctx, baseURL, apiKey, req, carrier, secondDeclared, metas)
	if !ok {
		return 0, false
	}

	declared, ok := domain.SolveCODDeclaredValue(
		domain.CODQuotePoint{Declared: target, Fee: firstFee},
		domain.CODQuotePoint{Declared: secondDeclared, Fee: secondFee},
		target,
	)
	if !ok {
		declared = secondDeclared
	}

	for i := 0; i < codMaxRefinements; i++ {
		fee, ok := u.probeCODFee(ctx, baseURL, apiKey, req, carrier, declared, metas)
		if !ok {
			return 0, false
		}

		if declared-fee >= target {
			u.log.Info(ctx).
				Str("carrier", carrier).
				Float64("net_target", target).
				Float64("declared", declared).
				Float64("carrier_fee", fee).
				Float64("net_recibido", declared-fee).
				Int("refinamientos", i).
				Msg("Valor COD declarado resuelto contra la comision real del carrier")
			return declared, true
		}

		next := math.Ceil(target + fee)
		if next <= declared {
			next = declared + 1
		}
		declared = next
	}

	u.log.Warn(ctx).
		Str("carrier", carrier).
		Float64("net_target", target).
		Float64("declared", declared).
		Msg("No se logro converger al valor COD que deja el neto esperado, se conserva el valor original")

	return 0, false
}

func (u *useCase) probeCODFee(ctx context.Context, baseURL, apiKey string, req domain.QuoteRequest, carrier string, declared float64, metas *[]domain.SyncMeta) (float64, bool) {
	probe := req
	probe.CODValue = declared
	probe.IDRate = 0
	probe.RequestPickup = false

	resp, err := u.Quote(ctx, baseURL, apiKey, probe, metas)
	if err != nil {
		u.log.Warn(ctx).
			Err(err).
			Float64("declared", declared).
			Msg("No se pudo cotizar la comision COD, se conserva el valor original")
		return 0, false
	}
	if resp == nil {
		return 0, false
	}

	fee, found := domain.FindCODFee(resp.Data.Rates, carrier)
	if !found {
		u.log.Warn(ctx).
			Str("carrier", carrier).
			Float64("declared", declared).
			Msg("La cotizacion no devolvio codDetails, se conserva el valor original")
		return 0, false
	}

	return fee, true
}
