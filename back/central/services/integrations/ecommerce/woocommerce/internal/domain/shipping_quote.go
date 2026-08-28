package domain

type ShippingQuoteRate struct {
	Flete                float64
	MinimumInsurance     float64
	ExtraInsurance       float64
	COD                  bool
	CODCarrierFee        float64
	CODProbabilityMargin float64
}
