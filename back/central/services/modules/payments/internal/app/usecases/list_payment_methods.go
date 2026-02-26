package usecases

import (
	"context"
	"math"

	"github.com/secamc93/probability/back/central/services/modules/payments/internal/app/usecases/mappers"
	"github.com/secamc93/probability/back/central/services/modules/payments/internal/domain/dtos"
)

// ListPaymentMethods obtiene una lista paginada de métodos de pago
func (uc *UseCase) ListPaymentMethods(ctx context.Context, page, pageSize int, filters map[string]interface{}) (*dtos.PaymentMethodsListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	methods, total, err := uc.repo.ListPaymentMethods(ctx, page, pageSize, filters)
	if err != nil {
		return nil, err
	}

	data := mappers.EntitiesToResponses(methods)

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return &dtos.PaymentMethodsListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}
