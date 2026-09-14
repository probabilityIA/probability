package entities

func CodAmountToCollect(codTotal float64, codCarrierFee float64, codIncludesShipping bool, codCheckoutCarrierFee float64) float64 {
	if codTotal <= 0 {
		return codTotal
	}
	if codCheckoutCarrierFee > 0 {
		return codTotal + codCheckoutCarrierFee
	}
	if codIncludesShipping || codCarrierFee <= 0 {
		return codTotal
	}
	return codTotal + codCarrierFee
}
