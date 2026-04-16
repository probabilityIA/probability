package app

import (
	"context"

	domainerrors "github.com/secamc93/probability/back/central/services/modules/announcements/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/domain/entities"
)

func (uc *UseCase) UpdateAnnouncement(ctx context.Context, dto dtos.UpdateAnnouncementDTO) (*entities.Announcement, error) {
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

	existing, err := uc.repo.GetByID(ctx, dto.ID)
	if err != nil {
		return nil, err
	}

	existing.BusinessID = dto.BusinessID
	existing.CategoryID = dto.CategoryID
	existing.Title = dto.Title
	existing.Message = dto.Message
	existing.DisplayType = dto.DisplayType
	existing.FrequencyType = dto.FrequencyType
	existing.Priority = dto.Priority
	existing.IsGlobal = dto.IsGlobal
	existing.StartsAt = dto.StartsAt
	existing.EndsAt = dto.EndsAt

	updated, err := uc.repo.Update(ctx, existing)
	if err != nil {
		uc.log.Error(ctx).Err(err).Uint("announcement_id", dto.ID).Msg("failed to update announcement")
		return nil, err
	}

	var links []entities.AnnouncementLink
	for _, l := range dto.Links {
		links = append(links, entities.AnnouncementLink{
			AnnouncementID: dto.ID,
			Label:          l.Label,
			URL:            l.URL,
			SortOrder:      l.SortOrder,
		})
	}
	if err := uc.repo.ReplaceLinks(ctx, dto.ID, links); err != nil {
		return nil, err
	}

	if dto.IsGlobal {
		if err := uc.repo.ReplaceTargets(ctx, dto.ID, nil); err != nil {
			return nil, err
		}
	} else {
		var targets []entities.AnnouncementTarget
		for _, bid := range dto.TargetIDs {
			targets = append(targets, entities.AnnouncementTarget{
				AnnouncementID: dto.ID,
				BusinessID:     bid,
			})
		}
		if err := uc.repo.ReplaceTargets(ctx, dto.ID, targets); err != nil {
			return nil, err
		}
	}

	return uc.repo.GetByID(ctx, updated.ID)
}
