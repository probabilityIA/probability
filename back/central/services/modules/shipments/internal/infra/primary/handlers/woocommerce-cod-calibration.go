package handlers

import (
	"context"
	"math"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

const (
	codCalibrationMinGap   = 10000.0
	codCalibrationHeadroom = 1.3
	codCalibrationMinSlack = 0.05
	codCalibrationMaxRatio = 3.0
)

var codCalibrationFields = []string{
	"flete",
	"minimumInsurance",
	"extraInsurance",
	"codCarrierFee",
	"codProbabilityMargin",
	"codCarrierFeeInsured",
}

func codRateKey(rate map[string]interface{}) string {
	if id := toFloat(rate["idRate"]); id > 0 {
		return "id:" + strconv.FormatFloat(id, 'f', -1, 64)
	}
	return "carrier:" + toStr(rate["carrier"]) + "|" + toStr(rate["product"])
}

func codTargetCost(rate map[string]interface{}, alwaysInsure bool) float64 {
	cost := toFloat(rate["flete"]) + toFloat(rate["minimumInsurance"]) + toFloat(rate["codProbabilityMargin"])
	if alwaysInsure {
		cost += toFloat(rate["extraInsurance"])
	}
	return cost
}

func codProbeValue(ratesList []map[string]interface{}, contentValue float64, alwaysInsure bool) float64 {
	worst := 0.0
	for _, rate := range ratesList {
		total := codTargetCost(rate, alwaysInsure) + toFloat(rate["codCarrierFee"])
		if total > worst {
			worst = total
		}
	}
	gap := worst * codCalibrationHeadroom
	if gap < codCalibrationMinGap {
		gap = codCalibrationMinGap
	}
	return math.Ceil(contentValue + gap)
}

func calibrateCODRate(base, probe map[string]interface{}, contentValue, baseDeclared, probeDeclared float64, alwaysInsure bool) (map[string]interface{}, bool) {
	span := probeDeclared - baseDeclared
	if span <= 0 {
		return nil, false
	}

	baseFee := toFloat(base["codCarrierFee"])
	probeFee := toFloat(probe["codCarrierFee"])
	baseCost := codTargetCost(base, alwaysInsure)
	probeCost := codTargetCost(probe, alwaysInsure)

	feeSlope := (probeFee - baseFee) / span
	costSlope := (probeCost - baseCost) / span
	if feeSlope < 0 || costSlope < 0 {
		return nil, false
	}

	denom := 1 - feeSlope - costSlope
	if denom < codCalibrationMinSlack {
		return nil, false
	}

	feeFixed := baseFee - feeSlope*baseDeclared
	costFixed := baseCost - costSlope*baseDeclared
	if feeFixed < 0 || costFixed < 0 {
		return nil, false
	}

	declared := math.Ceil((contentValue + feeFixed + costFixed) / denom)
	if declared <= baseDeclared || declared > contentValue*codCalibrationMaxRatio+probeDeclared {
		return nil, false
	}

	ratio := (declared - baseDeclared) / span

	out := make(map[string]interface{}, len(base))
	for k, v := range base {
		out[k] = v
	}
	for _, field := range codCalibrationFields {
		from, hasFrom := base[field]
		to, hasTo := probe[field]
		if !hasFrom || !hasTo {
			continue
		}
		value := toFloat(from) + (toFloat(to)-toFloat(from))*ratio
		if value < 0 {
			return nil, false
		}
		out[field] = math.Ceil(value)
	}

	return out, true
}

func calibrateCODRates(base, probe []map[string]interface{}, contentValue, baseDeclared, probeDeclared float64, alwaysInsure bool) []map[string]interface{} {
	probeByKey := make(map[string]map[string]interface{}, len(probe))
	for _, rate := range probe {
		probeByKey[codRateKey(rate)] = rate
	}

	out := make([]map[string]interface{}, 0, len(base))
	for _, rate := range base {
		match, ok := probeByKey[codRateKey(rate)]
		if !ok {
			out = append(out, rate)
			continue
		}
		calibrated, ok := calibrateCODRate(rate, match, contentValue, baseDeclared, probeDeclared, alwaysInsure)
		if !ok {
			out = append(out, rate)
			continue
		}
		out = append(out, calibrated)
	}

	return out
}

func (h *Handlers) calibrateWooCODRates(
	ctx context.Context,
	carrier *domain.CarrierInfo,
	businessID uint,
	payload map[string]interface{},
	ratesList []map[string]interface{},
	contentValue float64,
	alwaysInsure bool,
) []map[string]interface{} {
	if len(ratesList) == 0 || contentValue <= 0 {
		return ratesList
	}

	probeDeclared := codProbeValue(ratesList, contentValue, alwaysInsure)
	if probeDeclared <= contentValue {
		return ratesList
	}

	probePayload := make(map[string]interface{}, len(payload))
	for k, v := range payload {
		probePayload[k] = v
	}
	probePayload["contentValue"] = probeDeclared
	probePayload["codValue"] = probeDeclared

	result, err := h.runQuote(ctx, carrier, businessID, probePayload, uuid.New().String(), 12*time.Second)
	if err != nil || result == nil || result.Status != quoteStatusSuccess {
		h.logCot().Warn(ctx).Err(err).
			Uint("business_id", businessID).
			Float64("sondeo", probeDeclared).
			Msg("Sin calibrar: el sondeo de comision COD no respondio, se cotiza sobre el valor del producto")
		return ratesList
	}

	probeRates := toRatesList(getRatesFromData(result.Data))
	if len(probeRates) == 0 {
		h.logCot().Warn(ctx).
			Uint("business_id", businessID).
			Float64("sondeo", probeDeclared).
			Msg("Sin calibrar: el sondeo no devolvio tarifas")
		return ratesList
	}

	calibradas := calibrateCODRates(ratesList, probeRates, contentValue, contentValue, probeDeclared, alwaysInsure)

	for i := range calibradas {
		antes := toFloat(ratesList[i]["codCarrierFee"])
		ahora := toFloat(calibradas[i]["codCarrierFee"])
		if antes == ahora {
			continue
		}
		h.logCot().Info(ctx).
			Uint("business_id", businessID).
			Str("transportadora", toStr(calibradas[i]["carrier"])).
			Float64("productos", contentValue).
			Float64("comision_sobre_producto", antes).
			Float64("comision_calibrada", ahora).
			Float64("envio", codTargetCost(calibradas[i], alwaysInsure)+ahora).
			Msg("Comision COD calibrada sobre el total a recaudar")
	}

	return calibradas
}
