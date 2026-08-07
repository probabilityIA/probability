package domain

import "math"

type CODQuotePoint struct {
	Declared float64
	Fee      float64
}

func SolveCODDeclaredValue(p0, p1 CODQuotePoint, netTarget float64) (float64, bool) {
	if netTarget <= 0 {
		return 0, false
	}

	target := math.Ceil(netTarget)

	deltaDeclared := p1.Declared - p0.Declared
	if deltaDeclared == 0 {
		return math.Ceil(target + p0.Fee), true
	}

	rate := (p1.Fee - p0.Fee) / deltaDeclared
	if rate < 0 || rate >= 1 {
		return 0, false
	}

	fixed := p0.Fee - rate*p0.Declared
	if fixed < 0 {
		return 0, false
	}

	declared := (target + fixed) / (1 - rate)
	if declared < target {
		return 0, false
	}

	return math.Ceil(declared), true
}

func FindCODFee(rates []Rate, carrier string, rateID int64) (float64, bool) {
	if rateID > 0 {
		for i := range rates {
			if rates[i].CODDetails != nil && rates[i].IDRate == rateID {
				return rates[i].CODDetails.CODCost, true
			}
		}
	}

	for i := range rates {
		if rates[i].CODDetails != nil && carrier != "" && equalCarrier(rates[i].Carrier, carrier) {
			return rates[i].CODDetails.CODCost, true
		}
	}

	return 0, false
}

func equalCarrier(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		ca, cb := a[i], b[i]
		if ca >= 'a' && ca <= 'z' {
			ca -= 32
		}
		if cb >= 'a' && cb <= 'z' {
			cb -= 32
		}
		if ca != cb {
			return false
		}
	}
	return true
}
