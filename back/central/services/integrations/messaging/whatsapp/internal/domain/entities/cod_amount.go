package entities

func CodAmountToCollect(codTotal float64, codCarrierFee float64, codIncludesShipping bool) float64 {
	if codTotal <= 0 {
		return codTotal
	}
	if codIncludesShipping || codCarrierFee <= 0 {
		return codTotal
	}
	return codTotal + codCarrierFee
}
