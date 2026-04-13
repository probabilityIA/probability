package mocks

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/domain/entities"
)

type RepositoryMock struct {
	CreateFn                      func(ctx context.Context, announcement *entities.Announcement) (*entities.Announcement, error)
	GetByIDFn                     func(ctx context.Context, id uint) (*entities.Announcement, error)
	ListFn                        func(ctx context.Context, params dtos.ListAnnouncementsParams) ([]entities.Announcement, int64, error)
	UpdateFn                      func(ctx context.Context, announcement *entities.Announcement) (*entities.Announcement, error)
	DeleteFn                      func(ctx context.Context, id uint) error
	CreateImagesFn                func(ctx context.Context, images []entities.AnnouncementImage) error
	DeleteImagesByAnnouncementIDFn func(ctx context.Context, announcementID uint) error
	GetImagesByAnnouncementIDFn   func(ctx context.Context, announcementID uint) ([]entities.AnnouncementImage, error)
	GetImageByIDFn                func(ctx context.Context, id uint) (*entities.AnnouncementImage, error)
	DeleteImageByIDFn             func(ctx context.Context, id uint) error
	ReplaceLinksFn                func(ctx context.Context, announcementID uint, links []entities.AnnouncementLink) error
	ReplaceTargetsFn              func(ctx context.Context, announcementID uint, targets []entities.AnnouncementTarget) error
	RegisterViewFn                func(ctx context.Context, view *entities.AnnouncementView) error
	GetUserViewsFn                func(ctx context.Context, userID, announcementID uint) ([]entities.AnnouncementView, error)
	DeleteViewsByAnnouncementIDFn func(ctx context.Context, announcementID uint) error
	GetStatsFn                    func(ctx context.Context, announcementID uint) (*entities.AnnouncementStats, error)
	GetActiveAnnouncementsFn      func(ctx context.Context, params dtos.ActiveAnnouncementsParams) ([]entities.Announcement, error)
	ListCategoriesFn              func(ctx context.Context) ([]entities.AnnouncementCategory, error)
	ChangeStatusFn                func(ctx context.Context, id uint, status entities.AnnouncementStatus) error
}

func (m *RepositoryMock) Create(ctx context.Context, announcement *entities.Announcement) (*entities.Announcement, error) {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, announcement)
	}
	return announcement, nil
}

func (m *RepositoryMock) GetByID(ctx context.Context, id uint) (*entities.Announcement, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *RepositoryMock) List(ctx context.Context, params dtos.ListAnnouncementsParams) ([]entities.Announcement, int64, error) {
	if m.ListFn != nil {
		return m.ListFn(ctx, params)
	}
	return nil, 0, nil
}

func (m *RepositoryMock) Update(ctx context.Context, announcement *entities.Announcement) (*entities.Announcement, error) {
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, announcement)
	}
	return announcement, nil
}

func (m *RepositoryMock) Delete(ctx context.Context, id uint) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, id)
	}
	return nil
}

func (m *RepositoryMock) CreateImages(ctx context.Context, images []entities.AnnouncementImage) error {
	if m.CreateImagesFn != nil {
		return m.CreateImagesFn(ctx, images)
	}
	return nil
}

func (m *RepositoryMock) DeleteImagesByAnnouncementID(ctx context.Context, announcementID uint) error {
	if m.DeleteImagesByAnnouncementIDFn != nil {
		return m.DeleteImagesByAnnouncementIDFn(ctx, announcementID)
	}
	return nil
}

func (m *RepositoryMock) GetImagesByAnnouncementID(ctx context.Context, announcementID uint) ([]entities.AnnouncementImage, error) {
	if m.GetImagesByAnnouncementIDFn != nil {
		return m.GetImagesByAnnouncementIDFn(ctx, announcementID)
	}
	return nil, nil
}

func (m *RepositoryMock) GetImageByID(ctx context.Context, id uint) (*entities.AnnouncementImage, error) {
	if m.GetImageByIDFn != nil {
		return m.GetImageByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *RepositoryMock) DeleteImageByID(ctx context.Context, id uint) error {
	if m.DeleteImageByIDFn != nil {
		return m.DeleteImageByIDFn(ctx, id)
	}
	return nil
}

func (m *RepositoryMock) ReplaceLinks(ctx context.Context, announcementID uint, links []entities.AnnouncementLink) error {
	if m.ReplaceLinksFn != nil {
		return m.ReplaceLinksFn(ctx, announcementID, links)
	}
	return nil
}

func (m *RepositoryMock) ReplaceTargets(ctx context.Context, announcementID uint, targets []entities.AnnouncementTarget) error {
	if m.ReplaceTargetsFn != nil {
		return m.ReplaceTargetsFn(ctx, announcementID, targets)
	}
	return nil
}

func (m *RepositoryMock) RegisterView(ctx context.Context, view *entities.AnnouncementView) error {
	if m.RegisterViewFn != nil {
		return m.RegisterViewFn(ctx, view)
	}
	return nil
}

func (m *RepositoryMock) GetUserViews(ctx context.Context, userID, announcementID uint) ([]entities.AnnouncementView, error) {
	if m.GetUserViewsFn != nil {
		return m.GetUserViewsFn(ctx, userID, announcementID)
	}
	return nil, nil
}

func (m *RepositoryMock) DeleteViewsByAnnouncementID(ctx context.Context, announcementID uint) error {
	if m.DeleteViewsByAnnouncementIDFn != nil {
		return m.DeleteViewsByAnnouncementIDFn(ctx, announcementID)
	}
	return nil
}

func (m *RepositoryMock) GetStats(ctx context.Context, announcementID uint) (*entities.AnnouncementStats, error) {
	if m.GetStatsFn != nil {
		return m.GetStatsFn(ctx, announcementID)
	}
	return nil, nil
}

func (m *RepositoryMock) GetActiveAnnouncements(ctx context.Context, params dtos.ActiveAnnouncementsParams) ([]entities.Announcement, error) {
	if m.GetActiveAnnouncementsFn != nil {
		return m.GetActiveAnnouncementsFn(ctx, params)
	}
	return nil, nil
}

func (m *RepositoryMock) ListCategories(ctx context.Context) ([]entities.AnnouncementCategory, error) {
	if m.ListCategoriesFn != nil {
		return m.ListCategoriesFn(ctx)
	}
	return nil, nil
}

func (m *RepositoryMock) ChangeStatus(ctx context.Context, id uint, status entities.AnnouncementStatus) error {
	if m.ChangeStatusFn != nil {
		return m.ChangeStatusFn(ctx, id, status)
	}
	return nil
}
