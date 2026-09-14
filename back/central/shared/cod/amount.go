package cod

import (
	"math"
	"strings"
)

type Order struct {
	TotalAmount        float64
	CodTotal           float64
	IncludesShipping   bool
	CheckoutCarrierFee float64
	QuotedCarrier      string
	GuideCarrier       string
	CollectAmount      float64
}

const (
	FeeSourceCarrier  = "carrier"
	FeeSourceCheckout = "checkout"
)

var accentReplacer = strings.NewReplacer(
	"\u00c1", "A", "\u00c9", "E", "\u00cd", "I", "\u00d3", "O", "\u00da", "U", "\u00dc", "U", "\u00d1", "N",
)

func NormalizeCarrier(name string) string {
	if idx := strings.Index(name, " - "); idx >= 0 {
		name = name[:idx]
	}
	name = accentReplacer.Replace(strings.ToUpper(strings.TrimSpace(name)))
	var b strings.Builder
	for _, r := range name {
		if (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func CarrierChanged(o Order) bool {
	if o.CheckoutCarrierFee <= 0 {
		return false
	}
	quoted := NormalizeCarrier(o.QuotedCarrier)
	guide := NormalizeCarrier(o.GuideCarrier)
	return quoted != "" && guide != "" && quoted != guide
}

func CheckoutTotal(o Order) float64 {
	if o.CodTotal <= 0 {
		return 0
	}
	return o.CodTotal + o.CheckoutCarrierFee
}

func CustomerCharge(o Order, carrierFee float64) float64 {
	if o.CodTotal <= 0 {
		return 0
	}
	if o.CollectAmount > 0 {
		return o.CollectAmount
	}
	if o.CheckoutCarrierFee > 0 {
		if CarrierChanged(o) && carrierFee > 0 {
			return o.CodTotal + carrierFee
		}
		return o.CodTotal + o.CheckoutCarrierFee
	}
	if o.IncludesShipping || carrierFee <= 0 {
		return o.CodTotal
	}
	return o.CodTotal + carrierFee
}

type Breakdown struct {
	CustomerCharge    float64
	CheckoutTotal     float64
	CarrierFee        float64
	CarrierFeeSource  string
	ChargedCarrierFee float64
	BusinessNet       float64
	CarrierChanged    bool
}

func Summarize(o Order, carrierFee float64) Breakdown {
	if o.CodTotal <= 0 {
		return Breakdown{}
	}
	charge := CustomerCharge(o, carrierFee)
	changed := CarrierChanged(o)
	b := Breakdown{
		CustomerCharge:    charge,
		CheckoutTotal:     CheckoutTotal(o),
		ChargedCarrierFee: charge - o.CodTotal,
		CarrierChanged:    changed,
	}
	switch {
	case carrierFee > 0:
		b.CarrierFee = carrierFee
		b.CarrierFeeSource = FeeSourceCarrier
	case o.CheckoutCarrierFee > 0:
		b.CarrierFee = o.CheckoutCarrierFee
		b.CarrierFeeSource = FeeSourceCheckout
	}
	netFee := b.CarrierFee
	if o.CheckoutCarrierFee > 0 && !changed {
		netFee = o.CheckoutCarrierFee
	}
	b.BusinessNet = charge - netFee
	return b
}

func NetTarget(o Order, guideTotalCost float64, carrierFee float64) float64 {
	if o.CodTotal <= 0 {
		return 0
	}
	if CarrierChanged(o) {
		return o.CodTotal
	}
	if !o.IncludesShipping && guideTotalCost > 0 {
		return o.TotalAmount + guideTotalCost
	}
	if o.IncludesShipping {
		return CheckoutTotal(o) - carrierFee
	}
	return o.CodTotal
}

func AmountToCollect(o Order, guideTotalCost float64, carrierFee float64) float64 {
	base := NetTarget(o, guideTotalCost, carrierFee)
	if base <= 0 {
		return 0
	}
	return math.Ceil(base + carrierFee)
}
