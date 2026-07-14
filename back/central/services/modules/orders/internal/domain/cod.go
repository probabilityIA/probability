package domain

import "errors"

const (
	PaymentMethodCash uint = 5
	PaymentMethodCOD  uint = 6
)

var ErrCodTotalRequired = errors.New("una orden contra entrega requiere cod_total mayor a cero")

func NormalizeCod(isCod *bool, codTotal *float64, paymentMethodID uint) (bool, *float64, uint, error) {
	amount := 0.0
	if codTotal != nil {
		amount = *codTotal
	}

	cod := amount > 0
	if isCod != nil {
		cod = *isCod
	}

	if paymentMethodID == PaymentMethodCOD {
		cod = true
		paymentMethodID = PaymentMethodCash
	}

	if !cod {
		zero := 0.0
		return false, &zero, paymentMethodID, nil
	}

	if amount <= 0 {
		return false, codTotal, paymentMethodID, ErrCodTotalRequired
	}

	return true, &amount, paymentMethodID, nil
}
