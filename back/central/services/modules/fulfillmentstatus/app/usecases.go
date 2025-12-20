package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/fulfillmentstatus/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

// IUseCase define la interfaz para la lógica de negocio de estados de fulfillment
type IUseCase interface {
	GetFulfillmentStatusIDByCode(ctx context.Context, code string) (*uint, error)
	ListFulfillmentStatuses(ctx context.Context, isActive *bool) ([]domain.FulfillmentStatusInfo, error)
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

func (uc *UseCase) GetFulfillmentStatusIDByCode(ctx context.Context, code string) (*uint, error) {
	return uc.repo.GetFulfillmentStatusIDByCode(ctx, code)
}

func (uc *UseCase) ListFulfillmentStatuses(ctx context.Context, isActive *bool) ([]domain.FulfillmentStatusInfo, error) {
	statuses, err := uc.repo.ListFulfillmentStatuses(ctx, isActive)
	if err != nil {
		return nil, err
	}

	result := make([]domain.FulfillmentStatusInfo, len(statuses))
	for i, status := range statuses {
		result[i] = domain.FulfillmentStatusInfo{
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
