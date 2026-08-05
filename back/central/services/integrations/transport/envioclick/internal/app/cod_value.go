package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/transport/envioclick/internal/domain"
)

const codProbeGap = 10000.0

func (u *useCase) ResolveCODValue(ctx context.Context, baseURL, apiKey string, req domain.QuoteRequest, carrier string, netTarget float64, metas *[]domain.SyncMeta) (float64, bool) {
	if netTarget <= 0 {
		return 0, false
	}

	firstFee, ok := u.probeCODFee(ctx, baseURL, apiKey, req, carrier, netTarget, metas)
	if !ok {
		return 0, false
	}

	secondDeclared := netTarget + firstFee
	if secondDeclared <= netTarget {
		secondDeclared = netTarget + codProbeGap
	}

	secondFee, ok := u.probeCODFee(ctx, baseURL, apiKey, req, carrier, secondDeclared, metas)
	if !ok {
		return 0, false
	}

	declared, ok := domain.SolveCODDeclaredValue(
		domain.CODQuotePoint{Declared: netTarget, Fee: firstFee},
		domain.CODQuotePoint{Declared: secondDeclared, Fee: secondFee},
		netTarget,
	)
	if !ok {
		return 0, false
	}

	u.log.Info(ctx).
		Str("carrier", carrier).
		Float64("net_target", netTarget).
		Float64("fee_probe_1", firstFee).
		Float64("fee_probe_2", secondFee).
		Float64("declared", declared).
		Msg("COD declarado recalculado con calibracion de comision")

	return declared, true
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
