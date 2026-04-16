package app

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/domain/entities"
)

func (uc *UseCase) GetActiveAnnouncements(ctx context.Context, params dtos.ActiveAnnouncementsParams) ([]entities.Announcement, error) {
	announcements, err := uc.repo.GetActiveAnnouncements(ctx, params)
	if err != nil {
		return nil, err
	}

	var filtered []entities.Announcement
	for _, a := range announcements {
		visible, err := uc.shouldShow(ctx, a, params.UserID)
		if err != nil {
			uc.log.Warn(ctx).Err(err).Uint("announcement_id", a.ID).Msg("error checking visibility")
			continue
		}
		if visible {
			filtered = append(filtered, a)
		}
	}

	return filtered, nil
}

func (uc *UseCase) shouldShow(ctx context.Context, a entities.Announcement, userID uint) (bool, error) {
	if a.FrequencyType == entities.FrequencyAlways {
		return true, nil
	}

	views, err := uc.repo.GetUserViews(ctx, userID, a.ID)
	if err != nil {
		return false, err
	}

	if len(views) == 0 {
		return true, nil
	}

	switch a.FrequencyType {
	case entities.FrequencyOnce:
		return false, nil
	case entities.FrequencyRequiresAcceptance:
		for _, v := range views {
			if v.Action == entities.ViewActionAccepted {
				return false, nil
			}
		}
		return true, nil
	case entities.FrequencyDaily:
		today := time.Now().Truncate(24 * time.Hour)
		for _, v := range views {
			if v.ViewedAt.After(today) {
				return false, nil
			}
		}
		return true, nil
	}

	return true, nil
}
