package announcements

import (
	"context"
	"fmt"
	"sync"

	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/domain/entities"
)

var (
	alertCategoryOnce sync.Once
	alertCategoryID   uint
	alertCategoryErr  error
)

func (b *Bundle) resolveAlertCategoryID(ctx context.Context) (uint, error) {
	alertCategoryOnce.Do(func() {
		categories, err := b.UseCase.ListCategories(ctx)
		if err != nil {
			alertCategoryErr = err
			return
		}
		for _, c := range categories {
			if c.Code == "alert" {
				alertCategoryID = c.ID
				return
			}
		}
		alertCategoryErr = fmt.Errorf("alert category not found")
	})
	return alertCategoryID, alertCategoryErr
}

func (b *Bundle) CreateBusinessAlert(ctx context.Context, businessID uint, title, message string, createdByID uint, daily bool) (uint, error) {
	categoryID, err := b.resolveAlertCategoryID(ctx)
	if err != nil {
		return 0, err
	}

	frequency := entities.FrequencyOnce
	if daily {
		frequency = entities.FrequencyDaily
	}

	announcement, err := b.UseCase.CreateAnnouncement(ctx, dtos.CreateAnnouncementDTO{
		CategoryID:    categoryID,
		Title:         title,
		Message:       message,
		DisplayType:   entities.DisplayTypeModalText,
		FrequencyType: frequency,
		IsGlobal:      false,
		TargetIDs:     []uint{businessID},
		CreatedByID:   createdByID,
	})
	if err != nil {
		return 0, err
	}

	return announcement.ID, nil
}

func (b *Bundle) FindActiveBusinessAlert(ctx context.Context, businessID uint, title string) (*uint, error) {
	active, err := b.UseCase.GetActiveAnnouncements(ctx, dtos.ActiveAnnouncementsParams{BusinessID: businessID})
	if err != nil {
		return nil, err
	}

	for _, a := range active {
		if a.Title == title {
			id := a.ID
			return &id, nil
		}
	}

	return nil, nil
}

func (b *Bundle) DeactivateAnnouncement(ctx context.Context, id uint) error {
	return b.UseCase.ChangeStatus(ctx, dtos.ChangeStatusDTO{ID: id, Status: entities.StatusInactive})
}
