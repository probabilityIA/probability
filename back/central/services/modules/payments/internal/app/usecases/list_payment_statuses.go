package usecases

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/payments/internal/domain/dtos"
)

// ListPaymentStatuses obtiene el catálogo de estados de pago
func (uc *UseCase) ListPaymentStatuses(ctx context.Context, isActive *bool) ([]dtos.PaymentStatusInfo, error) {
	statuses, err := uc.repo.ListPaymentStatuses(ctx, isActive)
	if err != nil {
		return nil, err
	}

	result := make([]dtos.PaymentStatusInfo, len(statuses))
	for i, s := range statuses {
		result[i] = dtos.PaymentStatusInfo{
			ID:          s.ID,
			Code:        s.Code,
			Name:        s.Name,
			Description: s.Description,
			Category:    s.Category,
			Color:       s.Color,
			IsActive:    s.IsActive,
		}
	}

	return result, nil
}
