package usecases

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/payments/internal/domain/dtos"
)

func (uc *UseCase) ListChannelPaymentMethods(ctx context.Context, integrationType *string, isActive *bool) ([]dtos.ChannelPaymentMethodInfo, error) {
	methods, err := uc.repo.ListChannelPaymentMethods(ctx, integrationType, isActive)
	if err != nil {
		return nil, err
	}

	result := make([]dtos.ChannelPaymentMethodInfo, len(methods))
	for i, m := range methods {
		result[i] = dtos.ChannelPaymentMethodInfo{
			ID:              m.ID,
			IntegrationType: m.IntegrationType,
			Code:            m.Code,
			Name:            m.Name,
			Description:     m.Description,
			IsActive:        m.IsActive,
			DisplayOrder:    m.DisplayOrder,
		}
	}
	return result, nil
}
