package cod

import "math"

type Order struct {
	TotalAmount        float64
	CodTotal           float64
	IncludesShipping   bool
	CheckoutCarrierFee float64
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
	if o.CheckoutCarrierFee > 0 {
		return o.CodTotal + o.CheckoutCarrierFee
	}
	if o.IncludesShipping || carrierFee <= 0 {
		return o.CodTotal
	}
	return o.CodTotal + carrierFee
}

const (
	FeeSourceCarrier  = "carrier"
	FeeSourceCheckout = "checkout"
)

type Breakdown struct {
	CustomerCharge    float64
	CheckoutTotal     float64
	CarrierFee        float64
	CarrierFeeSource  string
	ChargedCarrierFee float64
	BusinessNet       float64
}

func Summarize(o Order, carrierFee float64) Breakdown {
	if o.CodTotal <= 0 {
		return Breakdown{}
	}
	charge := CustomerCharge(o, carrierFee)
	b := Breakdown{
		CustomerCharge:    charge,
		CheckoutTotal:     CheckoutTotal(o),
		ChargedCarrierFee: charge - o.CodTotal,
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
	if o.CheckoutCarrierFee > 0 {
		netFee = o.CheckoutCarrierFee
	}
	b.BusinessNet = charge - netFee
	return b
}

func NetTarget(o Order, guideTotalCost float64, carrierFee float64) float64 {
	if o.CodTotal <= 0 {
		return 0
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
