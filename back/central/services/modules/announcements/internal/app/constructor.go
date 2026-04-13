package app

import (
	"context"
	"mime/multipart"

	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

type IUseCase interface {
	CreateAnnouncement(ctx context.Context, dto dtos.CreateAnnouncementDTO) (*entities.Announcement, error)
	UpdateAnnouncement(ctx context.Context, dto dtos.UpdateAnnouncementDTO) (*entities.Announcement, error)
	DeleteAnnouncement(ctx context.Context, id uint) error
	UploadImage(ctx context.Context, announcementID uint, file *multipart.FileHeader, sortOrder int) (*entities.AnnouncementImage, error)
	DeleteImage(ctx context.Context, announcementID, imageID uint) error
	GetAnnouncement(ctx context.Context, id uint) (*entities.Announcement, error)
	ListAnnouncements(ctx context.Context, params dtos.ListAnnouncementsParams) ([]entities.Announcement, int64, error)
	GetActiveAnnouncements(ctx context.Context, params dtos.ActiveAnnouncementsParams) ([]entities.Announcement, error)
	RegisterView(ctx context.Context, dto dtos.RegisterViewDTO) error
	GetAnnouncementStats(ctx context.Context, announcementID uint) (*entities.AnnouncementStats, error)
	ListCategories(ctx context.Context) ([]entities.AnnouncementCategory, error)
	ChangeStatus(ctx context.Context, dto dtos.ChangeStatusDTO) error
	ForceRedisplay(ctx context.Context, id uint) error
}

type UseCase struct {
	repo    ports.IRepository
	storage ports.IStorageService
	log     log.ILogger
}

func New(repo ports.IRepository, storage ports.IStorageService, logger log.ILogger) IUseCase {
	return &UseCase{repo: repo, storage: storage, log: logger}
}
