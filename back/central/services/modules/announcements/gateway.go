package announcements

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/domain/entities"
)

// resolveAlertCategoryID cachea el ID de la categoria "alert" solo cuando se
// resuelve con exito. Un fallo transitorio (ej. DB no lista al arrancar) NUNCA
// queda cacheado: la proxima llamada vuelve a intentar la consulta.
func (b *Bundle) resolveAlertCategoryID(ctx context.Context) (uint, error) {
	b.alertCategoryMu.RLock()
	if b.alertCategoryResolved {
		id := b.alertCategoryID
		b.alertCategoryMu.RUnlock()
		return id, nil
	}
	b.alertCategoryMu.RUnlock()

	categories, err := b.UseCase.ListCategories(ctx)
	if err != nil {
		return 0, err
	}
	for _, c := range categories {
		if c.Code == "alert" {
			b.alertCategoryMu.Lock()
			b.alertCategoryID = c.ID
			b.alertCategoryResolved = true
			b.alertCategoryMu.Unlock()
			return c.ID, nil
		}
	}
	return 0, fmt.Errorf("alert category not found")
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
