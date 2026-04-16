package ports

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/domain/entities"
)

type IRepository interface {
	Create(ctx context.Context, announcement *entities.Announcement) (*entities.Announcement, error)
	GetByID(ctx context.Context, id uint) (*entities.Announcement, error)
	List(ctx context.Context, params dtos.ListAnnouncementsParams) ([]entities.Announcement, int64, error)
	Update(ctx context.Context, announcement *entities.Announcement) (*entities.Announcement, error)
	Delete(ctx context.Context, id uint) error

	CreateImages(ctx context.Context, images []entities.AnnouncementImage) error
	DeleteImagesByAnnouncementID(ctx context.Context, announcementID uint) error
	GetImagesByAnnouncementID(ctx context.Context, announcementID uint) ([]entities.AnnouncementImage, error)
	GetImageByID(ctx context.Context, id uint) (*entities.AnnouncementImage, error)
	DeleteImageByID(ctx context.Context, id uint) error

	ReplaceLinks(ctx context.Context, announcementID uint, links []entities.AnnouncementLink) error

	ReplaceTargets(ctx context.Context, announcementID uint, targets []entities.AnnouncementTarget) error

	RegisterView(ctx context.Context, view *entities.AnnouncementView) error
	GetUserViews(ctx context.Context, userID, announcementID uint) ([]entities.AnnouncementView, error)
	DeleteViewsByAnnouncementID(ctx context.Context, announcementID uint) error

	GetStats(ctx context.Context, announcementID uint) (*entities.AnnouncementStats, error)

	GetActiveAnnouncements(ctx context.Context, params dtos.ActiveAnnouncementsParams) ([]entities.Announcement, error)

	ListCategories(ctx context.Context) ([]entities.AnnouncementCategory, error)

	ChangeStatus(ctx context.Context, id uint, status entities.AnnouncementStatus) error
}

type IStorageService interface {
	UploadFile(ctx context.Context, folder string, filename string, data []byte, contentType string) (string, error)
	DeleteFile(ctx context.Context, fileURL string) error
}
