package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/paymentstatus/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

// IUseCase define la interfaz para la lógica de negocio de estados de pago
type IUseCase interface {
	GetPaymentStatusIDByCode(ctx context.Context, code string) (*uint, error)
	ListPaymentStatuses(ctx context.Context, isActive *bool) ([]domain.PaymentStatusInfo, error)
}

type UseCase struct {
	repo   domain.IRepository
	logger log.ILogger
}

func New(repo domain.IRepository, logger log.ILogger) IUseCase {
	return &UseCase{
		repo:   repo,
		logger: logger,
	}
}

func (uc *UseCase) GetPaymentStatusIDByCode(ctx context.Context, code string) (*uint, error) {
	return uc.repo.GetPaymentStatusIDByCode(ctx, code)
}

func (uc *UseCase) ListPaymentStatuses(ctx context.Context, isActive *bool) ([]domain.PaymentStatusInfo, error) {
	statuses, err := uc.repo.ListPaymentStatuses(ctx, isActive)
	if err != nil {
		return nil, err
	}

	result := make([]domain.PaymentStatusInfo, len(statuses))
	for i, status := range statuses {
		result[i] = domain.PaymentStatusInfo{
			ID:          status.ID,
			Code:        status.Code,
			Name:        status.Name,
			Description: status.Description,
			Category:    status.Category,
			Color:       status.Color,
		}
	}

	return result, nil
}
