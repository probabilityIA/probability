package dtos

import (
	"time"

	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/domain/entities"
)

type CreateAnnouncementDTO struct {
	BusinessID    *uint
	CategoryID    uint
	Title         string
	Message       string
	DisplayType   entities.DisplayType
	FrequencyType entities.FrequencyType
	Priority      int
	IsGlobal      bool
	StartsAt      *time.Time
	EndsAt        *time.Time
	CreatedByID   uint
	Links         []CreateLinkDTO
	TargetIDs     []uint
}

type CreateLinkDTO struct {
	Label     string
	URL       string
	SortOrder int
}

type UpdateAnnouncementDTO struct {
	ID            uint
	BusinessID    *uint
	CategoryID    uint
	Title         string
	Message       string
	DisplayType   entities.DisplayType
	FrequencyType entities.FrequencyType
	Priority      int
	IsGlobal      bool
	StartsAt      *time.Time
	EndsAt        *time.Time
	Links         []CreateLinkDTO
	TargetIDs     []uint
}

type ListAnnouncementsParams struct {
	Status     string
	CategoryID *uint
	Search     string
	Page       int
	PageSize   int
}

func (p ListAnnouncementsParams) Offset() int {
	if p.Page < 1 {
		p.Page = 1
	}
	return (p.Page - 1) * p.PageSize
}

type ActiveAnnouncementsParams struct {
	BusinessID uint
	UserID     uint
}

type RegisterViewDTO struct {
	AnnouncementID uint
	UserID         uint
	BusinessID     uint
	Action         entities.ViewAction
	LinkID         *uint
}

type ChangeStatusDTO struct {
	ID     uint
	Status entities.AnnouncementStatus
}
