package usecasebusinesstype

import (
	"context"

	"github.com/secamc93/probability/back/central/services/auth/bussines/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

type IUseCaseBusinessType interface {
	GetBusinessTypes(ctx context.Context) ([]domain.BusinessTypeResponse, error)
	GetBusinessTypeByID(ctx context.Context, id uint) (*domain.BusinessTypeResponse, error)
	CreateBusinessType(ctx context.Context, request domain.BusinessTypeRequest) (*domain.BusinessTypeResponse, error)
	UpdateBusinessType(ctx context.Context, id uint, request domain.BusinessTypeRequest) (*domain.BusinessTypeResponse, error)
	DeleteBusinessType(ctx context.Context, id uint) error
}

type BusinessTypeUseCase struct {
	repository domain.IBusinessRepository
	log        log.ILogger
}

func New(repository domain.IBusinessRepository, log log.ILogger) IUseCaseBusinessType {
	return &BusinessTypeUseCase{
		repository: repository,
		log:        log,
	}
}
