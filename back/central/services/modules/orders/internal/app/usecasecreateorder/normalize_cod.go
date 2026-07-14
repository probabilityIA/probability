package usecasecreateorder

import (
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
)

func normalizeCodOnDTO(dto *dtos.ProbabilityOrderDTO) error {
	var paymentMethodID uint
	if len(dto.Payments) > 0 {
		paymentMethodID = dto.Payments[0].PaymentMethodID
	}

	var isCod *bool
	if dto.IsCod {
		isCod = &dto.IsCod
	}

	cod, amount, method, err := domain.NormalizeCod(isCod, dto.CodTotal, paymentMethodID)
	if err != nil {
		return err
	}

	dto.IsCod = cod
	dto.CodTotal = amount
	if len(dto.Payments) > 0 {
		dto.Payments[0].PaymentMethodID = method
	}
	return nil
}
