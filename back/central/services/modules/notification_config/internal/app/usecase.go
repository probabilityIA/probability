package app

import (
	"context"
	"errors"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain"
)

type UseCase struct {
	repo domain.IRepository
}

func New(repo domain.IRepository) domain.IUseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) CreateConfig(ctx context.Context, dto domain.CreateConfigDTO) (*domain.NotificationConfig, error) {
	// Check if config already exists for this business and event type
	existing, err := uc.repo.GetByBusinessAndEventType(ctx, dto.BusinessID, dto.EventType)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("configuration already exists for this business and event type")
	}

	config := &domain.NotificationConfig{
		BusinessID:  dto.BusinessID,
		EventType:   dto.EventType,
		Enabled:     dto.Enabled,
		Channels:    dto.Channels,
		Filters:     dto.Filters,
		Description: dto.Description,
	}

	if err := uc.repo.Create(ctx, config); err != nil {
		return nil, err
	}

	return config, nil
}

func (uc *UseCase) GetConfig(ctx context.Context, id uint) (*domain.NotificationConfig, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *UseCase) UpdateConfig(ctx context.Context, id uint, dto domain.UpdateConfigDTO) (*domain.NotificationConfig, error) {
	config, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, errors.New("configuration not found")
	}

	if dto.Enabled != nil {
		config.Enabled = *dto.Enabled
	}
	if dto.Channels != nil {
		config.Channels = dto.Channels
	}
	if dto.Filters != nil {
		config.Filters = dto.Filters
	}
	if dto.Description != nil {
		config.Description = *dto.Description
	}

	if err := uc.repo.Update(ctx, id, config); err != nil {
		return nil, err
	}

	return config, nil
}

func (uc *UseCase) DeleteConfig(ctx context.Context, id uint) error {
	return uc.repo.Delete(ctx, id)
}

func (uc *UseCase) ListConfigs(ctx context.Context, filter domain.ConfigFilter) ([]*domain.NotificationConfig, error) {
	return uc.repo.List(ctx, filter)
}
