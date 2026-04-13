package app

import (
	"context"
	"time"

	domainerrors "github.com/secamc93/probability/back/central/services/modules/announcements/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/domain/entities"
)

func (uc *UseCase) CreateAnnouncement(ctx context.Context, dto dtos.CreateAnnouncementDTO) (*entities.Announcement, error) {
	if err := validateDisplayType(dto.DisplayType); err != nil {
		return nil, err
	}
	if err := validateFrequencyType(dto.FrequencyType); err != nil {
		return nil, err
	}
	if dto.StartsAt != nil && dto.EndsAt != nil && dto.StartsAt.After(*dto.EndsAt) {
		return nil, domainerrors.ErrInvalidDateRange
	}
	if !dto.IsGlobal && len(dto.TargetIDs) == 0 {
		return nil, domainerrors.ErrTargetsRequired
	}

	status := resolveInitialStatus(dto.StartsAt)

	announcement := &entities.Announcement{
		BusinessID:    dto.BusinessID,
		CategoryID:    dto.CategoryID,
		Title:         dto.Title,
		Message:       dto.Message,
		DisplayType:   dto.DisplayType,
		FrequencyType: dto.FrequencyType,
		Priority:      dto.Priority,
		IsGlobal:      dto.IsGlobal,
		Status:        status,
		StartsAt:      dto.StartsAt,
		EndsAt:        dto.EndsAt,
		CreatedByID:   dto.CreatedByID,
	}

	for _, l := range dto.Links {
		announcement.Links = append(announcement.Links, entities.AnnouncementLink{
			Label:     l.Label,
			URL:       l.URL,
			SortOrder: l.SortOrder,
		})
	}

	if !dto.IsGlobal {
		for _, bid := range dto.TargetIDs {
			announcement.Targets = append(announcement.Targets, entities.AnnouncementTarget{
				BusinessID: bid,
			})
		}
	}

	created, err := uc.repo.Create(ctx, announcement)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("failed to create announcement")
		return nil, err
	}

	return created, nil
}

func resolveInitialStatus(startsAt *time.Time) entities.AnnouncementStatus {
	if startsAt != nil && startsAt.After(time.Now()) {
		return entities.StatusScheduled
	}
	return entities.StatusActive
}
