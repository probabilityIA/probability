package repository

import (
	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
)

func announcementFromEntity(e *entities.Announcement) *models.Announcement {
	m := &models.Announcement{
		BusinessID:     e.BusinessID,
		CategoryID:     e.CategoryID,
		Title:          e.Title,
		Message:        e.Message,
		DisplayType:    string(e.DisplayType),
		FrequencyType:  string(e.FrequencyType),
		Priority:       e.Priority,
		IsGlobal:       e.IsGlobal,
		Status:         string(e.Status),
		StartsAt:       e.StartsAt,
		EndsAt:         e.EndsAt,
		ForceRedisplay: e.ForceRedisplay,
		CreatedByID:    e.CreatedByID,
	}

	for _, l := range e.Links {
		m.Links = append(m.Links, models.AnnouncementLink{
			Label:     l.Label,
			URL:       l.URL,
			SortOrder: l.SortOrder,
		})
	}

	for _, t := range e.Targets {
		m.Targets = append(m.Targets, models.AnnouncementTarget{
			BusinessID: t.BusinessID,
		})
	}

	return m
}

func announcementToEntity(m *models.Announcement) *entities.Announcement {
	e := &entities.Announcement{
		ID:             m.ID,
		BusinessID:     m.BusinessID,
		CategoryID:     m.CategoryID,
		Title:          m.Title,
		Message:        m.Message,
		DisplayType:    entities.DisplayType(m.DisplayType),
		FrequencyType:  entities.FrequencyType(m.FrequencyType),
		Priority:       m.Priority,
		IsGlobal:       m.IsGlobal,
		Status:         entities.AnnouncementStatus(m.Status),
		StartsAt:       m.StartsAt,
		EndsAt:         m.EndsAt,
		ForceRedisplay: m.ForceRedisplay,
		CreatedByID:    m.CreatedByID,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}

	if m.Category.ID != 0 {
		e.Category = &entities.AnnouncementCategory{
			ID:    m.Category.ID,
			Code:  m.Category.Code,
			Name:  m.Category.Name,
			Icon:  m.Category.Icon,
			Color: m.Category.Color,
		}
	}

	for _, img := range m.Images {
		e.Images = append(e.Images, entities.AnnouncementImage{
			ID:             img.ID,
			AnnouncementID: img.AnnouncementID,
			ImageURL:       img.ImageURL,
			SortOrder:      img.SortOrder,
		})
	}

	for _, l := range m.Links {
		e.Links = append(e.Links, entities.AnnouncementLink{
			ID:             l.ID,
			AnnouncementID: l.AnnouncementID,
			Label:          l.Label,
			URL:            l.URL,
			SortOrder:      l.SortOrder,
		})
	}

	for _, t := range m.Targets {
		e.Targets = append(e.Targets, entities.AnnouncementTarget{
			ID:             t.ID,
			AnnouncementID: t.AnnouncementID,
			BusinessID:     t.BusinessID,
		})
	}

	return e
}

func imageToEntity(m *models.AnnouncementImage) entities.AnnouncementImage {
	return entities.AnnouncementImage{
		ID:             m.ID,
		AnnouncementID: m.AnnouncementID,
		ImageURL:       m.ImageURL,
		SortOrder:      m.SortOrder,
	}
}

func viewToEntity(m *models.AnnouncementView) entities.AnnouncementView {
	return entities.AnnouncementView{
		ID:             m.ID,
		AnnouncementID: m.AnnouncementID,
		UserID:         m.UserID,
		BusinessID:     m.BusinessID,
		Action:         entities.ViewAction(m.Action),
		LinkID:         m.LinkID,
		ViewedAt:       m.ViewedAt,
		CreatedAt:      m.CreatedAt,
	}
}

func categoryToEntity(m *models.AnnouncementCategory) entities.AnnouncementCategory {
	return entities.AnnouncementCategory{
		ID:    m.ID,
		Code:  m.Code,
		Name:  m.Name,
		Icon:  m.Icon,
		Color: m.Color,
	}
}
