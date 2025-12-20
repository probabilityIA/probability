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
	ListOrderStatuses(ctx context.Context, isActive *bool) ([]domain.OrderStatusInfo, error)
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
	exists, err := uc.repo.Exists(ctx, mapping.IntegrationTypeID, mapping.OriginalStatus)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("mapping already exists for this integration type and original status")
	}

	model := &models.OrderStatusMapping{
		IntegrationTypeID: mapping.IntegrationTypeID,
		OriginalStatus:    mapping.OriginalStatus,
		OrderStatusID:     mapping.OrderStatusID,
		Priority:          mapping.Priority,
		Description:       mapping.Description,
		IsActive:          true,
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
	model.OrderStatusID = mapping.OrderStatusID
	model.Priority = mapping.Priority
	model.Description = mapping.Description

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
	result := &domain.OrderStatusMapping{
		ID:                m.ID,
		IntegrationTypeID: m.IntegrationTypeID,
		OriginalStatus:    m.OriginalStatus,
		OrderStatusID:     m.OrderStatusID,
		IsActive:          m.IsActive,
		Priority:          m.Priority,
		Description:       m.Description,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
	}

	// Incluir información del IntegrationType si está cargado
	if m.IntegrationType.ID != 0 {
		result.IntegrationType = &domain.IntegrationTypeInfo{
			ID:       m.IntegrationType.ID,
			Code:     m.IntegrationType.Code,
			Name:     m.IntegrationType.Name,
			ImageURL: m.IntegrationType.ImageURL,
		}
	}

	// Incluir información del OrderStatus si está cargado
	if m.OrderStatus.ID != 0 {
		result.OrderStatus = &domain.OrderStatusInfo{
			ID:          m.OrderStatus.ID,
			Code:        m.OrderStatus.Code,
			Name:        m.OrderStatus.Name,
			Description: m.OrderStatus.Description,
			Category:    m.OrderStatus.Category,
			Color:       m.OrderStatus.Color,
		}
	}

	return result
}

// ListOrderStatuses lista todos los estados de órdenes de Probability
func (uc *UseCase) ListOrderStatuses(ctx context.Context, isActive *bool) ([]domain.OrderStatusInfo, error) {
	statuses, err := uc.repo.ListOrderStatuses(ctx, isActive)
	if err != nil {
		return nil, err
	}

	result := make([]domain.OrderStatusInfo, len(statuses))
	for i, status := range statuses {
		result[i] = domain.OrderStatusInfo{
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
