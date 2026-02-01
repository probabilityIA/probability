package usecases

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/paymentstatus/internal/app/usecases/mappers"
	"github.com/secamc93/probability/back/central/services/modules/paymentstatus/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/paymentstatus/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

// IUseCase define la interfaz para la lógica de negocio de estados de pago
type IUseCase interface {
	GetPaymentStatusIDByCode(ctx context.Context, code string) (*uint, error)
	ListPaymentStatuses(ctx context.Context, isActive *bool) ([]dtos.PaymentStatusInfo, error)
}

type UseCase struct {
	repo   ports.IRepository
	logger log.ILogger
}

func (uc *UseCase) GetPaymentStatusIDByCode(ctx context.Context, code string) (*uint, error) {
	return uc.repo.GetPaymentStatusIDByCode(ctx, code)
}

func (uc *UseCase) ListPaymentStatuses(ctx context.Context, isActive *bool) ([]dtos.PaymentStatusInfo, error) {
	// ✅ Obtener entidades de dominio desde repositorio
	statuses, err := uc.repo.ListPaymentStatuses(ctx, isActive)
	if err != nil {
		return nil, err
	}

	// ✅ Delegar mapeo a mapper
	return mappers.ToPaymentStatusInfoList(statuses), nil
}
