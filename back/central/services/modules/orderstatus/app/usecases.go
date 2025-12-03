package app

import (
	"context"
	"errors"

	"github.com/secamc93/probability/back/central/services/modules/orderstatus/domain"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/migration/shared/models"
)

// IUseCase define la interfaz para la lógica de negocio de mapeos de estado
type IUseCase interface {
	CreateOrderStatusMapping(ctx context.Context, mapping *domain.OrderStatusMapping) (*domain.OrderStatusMapping, error)
	GetOrderStatusMapping(ctx context.Context, id uint) (*domain.OrderStatusMapping, error)
	ListOrderStatusMappings(ctx context.Context, filters map[string]interface{}) ([]domain.OrderStatusMapping, int64, error)
	UpdateOrderStatusMapping(ctx context.Context, id uint, mapping *domain.OrderStatusMapping) (*domain.OrderStatusMapping, error)
	DeleteOrderStatusMapping(ctx context.Context, id uint) error
	ToggleOrderStatusMappingActive(ctx context.Context, id uint) (*domain.OrderStatusMapping, error)
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

func (uc *UseCase) CreateOrderStatusMapping(ctx context.Context, mapping *domain.OrderStatusMapping) (*domain.OrderStatusMapping, error) {
	// Verificar si ya existe
	exists, err := uc.repo.Exists(ctx, mapping.IntegrationType, mapping.OriginalStatus)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("mapping already exists for this integration type and original status")
	}

	model := &models.OrderStatusMapping{
		IntegrationType: mapping.IntegrationType,
		OriginalStatus:  mapping.OriginalStatus,
		MappedStatus:    mapping.MappedStatus,
		Priority:        mapping.Priority,
		Description:     mapping.Description,
		IsActive:        true,
	}

	if !model.IsValidMappedStatus() {
		return nil, errors.New("invalid mapped status")
	}

	if err := uc.repo.Create(ctx, model); err != nil {
		return nil, err
	}

	return toDomain(model), nil
}

func (uc *UseCase) GetOrderStatusMapping(ctx context.Context, id uint) (*domain.OrderStatusMapping, error) {
	model, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toDomain(model), nil
}

func (uc *UseCase) ListOrderStatusMappings(ctx context.Context, filters map[string]interface{}) ([]domain.OrderStatusMapping, int64, error) {
	modelsList, total, err := uc.repo.List(ctx, filters)
	if err != nil {
		return nil, 0, err
	}

	var response []domain.OrderStatusMapping
	for _, m := range modelsList {
		response = append(response, *toDomain(&m))
	}

	return response, total, nil
}

func (uc *UseCase) UpdateOrderStatusMapping(ctx context.Context, id uint, mapping *domain.OrderStatusMapping) (*domain.OrderStatusMapping, error) {
	model, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	model.OriginalStatus = mapping.OriginalStatus
	model.MappedStatus = mapping.MappedStatus
	model.Priority = mapping.Priority
	model.Description = mapping.Description

	if !model.IsValidMappedStatus() {
		return nil, errors.New("invalid mapped status")
	}

	if err := uc.repo.Update(ctx, model); err != nil {
		return nil, err
	}

	return toDomain(model), nil
}

func (uc *UseCase) DeleteOrderStatusMapping(ctx context.Context, id uint) error {
	return uc.repo.Delete(ctx, id)
}

func (uc *UseCase) ToggleOrderStatusMappingActive(ctx context.Context, id uint) (*domain.OrderStatusMapping, error) {
	model, err := uc.repo.ToggleActive(ctx, id)
	if err != nil {
		return nil, err
	}
	return toDomain(model), nil
}

func toDomain(m *models.OrderStatusMapping) *domain.OrderStatusMapping {
	return &domain.OrderStatusMapping{
		ID:              m.ID,
		IntegrationType: m.IntegrationType,
		OriginalStatus:  m.OriginalStatus,
		MappedStatus:    m.MappedStatus,
		IsActive:        m.IsActive,
		Priority:        m.Priority,
		Description:     m.Description,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}
