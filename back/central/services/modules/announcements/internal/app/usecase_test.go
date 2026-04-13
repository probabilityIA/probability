package app

import (
	"context"
	"errors"
	"testing"
	"time"

	domainerrors "github.com/secamc93/probability/back/central/services/modules/announcements/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/mocks"
)

func newUseCase(repo *mocks.RepositoryMock, storage *mocks.StorageMock) IUseCase {
	return New(repo, storage, &mocks.LoggerMock{})
}

func TestCreateAnnouncement_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mocks.RepositoryMock{
		CreateFn: func(_ context.Context, a *entities.Announcement) (*entities.Announcement, error) {
			a.ID = 1
			return a, nil
		},
	}
	uc := newUseCase(repo, &mocks.StorageMock{})

	dto := dtos.CreateAnnouncementDTO{
		Title:         "Test",
		Message:       "Hello",
		DisplayType:   entities.DisplayTypeModalText,
		FrequencyType: entities.FrequencyAlways,
		IsGlobal:      true,
	}

	result, err := uc.CreateAnnouncement(ctx, dto)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected announcement, got nil")
	}
	if result.ID != 1 {
		t.Errorf("expected ID 1, got %d", result.ID)
	}
	if result.Status != entities.StatusActive {
		t.Errorf("expected status active, got %s", result.Status)
	}
}

func TestCreateAnnouncement_ScheduledWhenFutureStartsAt(t *testing.T) {
	ctx := context.Background()
	future := time.Now().Add(24 * time.Hour)
	repo := &mocks.RepositoryMock{
		CreateFn: func(_ context.Context, a *entities.Announcement) (*entities.Announcement, error) {
			a.ID = 2
			return a, nil
		},
	}
	uc := newUseCase(repo, &mocks.StorageMock{})

	dto := dtos.CreateAnnouncementDTO{
		Title:         "Scheduled",
		DisplayType:   entities.DisplayTypeModalText,
		FrequencyType: entities.FrequencyOnce,
		IsGlobal:      true,
		StartsAt:      &future,
	}

	result, err := uc.CreateAnnouncement(ctx, dto)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Status != entities.StatusScheduled {
		t.Errorf("expected status scheduled, got %s", result.Status)
	}
}

func TestCreateAnnouncement_InvalidDisplayType(t *testing.T) {
	ctx := context.Background()
	uc := newUseCase(&mocks.RepositoryMock{}, &mocks.StorageMock{})

	dto := dtos.CreateAnnouncementDTO{
		DisplayType:   "invalid_type",
		FrequencyType: entities.FrequencyAlways,
		IsGlobal:      true,
	}

	_, err := uc.CreateAnnouncement(ctx, dto)
	if !errors.Is(err, domainerrors.ErrInvalidDisplayType) {
		t.Errorf("expected ErrInvalidDisplayType, got %v", err)
	}
}

func TestCreateAnnouncement_InvalidFrequencyType(t *testing.T) {
	ctx := context.Background()
	uc := newUseCase(&mocks.RepositoryMock{}, &mocks.StorageMock{})

	dto := dtos.CreateAnnouncementDTO{
		DisplayType:   entities.DisplayTypeModalText,
		FrequencyType: "invalid_freq",
		IsGlobal:      true,
	}

	_, err := uc.CreateAnnouncement(ctx, dto)
	if !errors.Is(err, domainerrors.ErrInvalidFrequencyType) {
		t.Errorf("expected ErrInvalidFrequencyType, got %v", err)
	}
}

func TestCreateAnnouncement_InvalidDateRange(t *testing.T) {
	ctx := context.Background()
	uc := newUseCase(&mocks.RepositoryMock{}, &mocks.StorageMock{})

	starts := time.Now().Add(48 * time.Hour)
	ends := time.Now().Add(24 * time.Hour)

	dto := dtos.CreateAnnouncementDTO{
		DisplayType:   entities.DisplayTypeModalText,
		FrequencyType: entities.FrequencyAlways,
		IsGlobal:      true,
		StartsAt:      &starts,
		EndsAt:        &ends,
	}

	_, err := uc.CreateAnnouncement(ctx, dto)
	if !errors.Is(err, domainerrors.ErrInvalidDateRange) {
		t.Errorf("expected ErrInvalidDateRange, got %v", err)
	}
}

func TestCreateAnnouncement_TargetsRequiredWhenNotGlobal(t *testing.T) {
	ctx := context.Background()
	uc := newUseCase(&mocks.RepositoryMock{}, &mocks.StorageMock{})

	dto := dtos.CreateAnnouncementDTO{
		DisplayType:   entities.DisplayTypeModalText,
		FrequencyType: entities.FrequencyAlways,
		IsGlobal:      false,
		TargetIDs:     []uint{},
	}

	_, err := uc.CreateAnnouncement(ctx, dto)
	if !errors.Is(err, domainerrors.ErrTargetsRequired) {
		t.Errorf("expected ErrTargetsRequired, got %v", err)
	}
}

func TestCreateAnnouncement_NonGlobalWithTargetsSetsTargets(t *testing.T) {
	ctx := context.Background()
	repo := &mocks.RepositoryMock{
		CreateFn: func(_ context.Context, a *entities.Announcement) (*entities.Announcement, error) {
			a.ID = 5
			return a, nil
		},
	}
	uc := newUseCase(repo, &mocks.StorageMock{})

	dto := dtos.CreateAnnouncementDTO{
		DisplayType:   entities.DisplayTypeModalText,
		FrequencyType: entities.FrequencyAlways,
		IsGlobal:      false,
		TargetIDs:     []uint{10, 20},
	}

	result, err := uc.CreateAnnouncement(ctx, dto)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result.Targets) != 2 {
		t.Errorf("expected 2 targets, got %d", len(result.Targets))
	}
}

func TestCreateAnnouncement_RepositoryError(t *testing.T) {
	ctx := context.Background()
	repoErr := errors.New("db error")
	repo := &mocks.RepositoryMock{
		CreateFn: func(_ context.Context, _ *entities.Announcement) (*entities.Announcement, error) {
			return nil, repoErr
		},
	}
	uc := newUseCase(repo, &mocks.StorageMock{})

	dto := dtos.CreateAnnouncementDTO{
		DisplayType:   entities.DisplayTypeModalText,
		FrequencyType: entities.FrequencyAlways,
		IsGlobal:      true,
	}

	_, err := uc.CreateAnnouncement(ctx, dto)
	if !errors.Is(err, repoErr) {
		t.Errorf("expected repo error, got %v", err)
	}
}

func TestUpdateAnnouncement_Success(t *testing.T) {
	ctx := context.Background()
	existing := &entities.Announcement{ID: 1, Title: "Old", DisplayType: entities.DisplayTypeModalText, FrequencyType: entities.FrequencyAlways}
	updated := &entities.Announcement{ID: 1, Title: "New"}

	getByIDCallCount := 0
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(_ context.Context, id uint) (*entities.Announcement, error) {
			getByIDCallCount++
			if getByIDCallCount == 1 {
				return existing, nil
			}
			return updated, nil
		},
		UpdateFn: func(_ context.Context, a *entities.Announcement) (*entities.Announcement, error) {
			return a, nil
		},
		ReplaceLinksFn: func(_ context.Context, _ uint, _ []entities.AnnouncementLink) error {
			return nil
		},
		ReplaceTargetsFn: func(_ context.Context, _ uint, _ []entities.AnnouncementTarget) error {
			return nil
		},
	}
	uc := newUseCase(repo, &mocks.StorageMock{})

	dto := dtos.UpdateAnnouncementDTO{
		ID:            1,
		Title:         "New",
		DisplayType:   entities.DisplayTypeModalText,
		FrequencyType: entities.FrequencyAlways,
		IsGlobal:      true,
	}

	result, err := uc.UpdateAnnouncement(ctx, dto)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected updated announcement, got nil")
	}
}

func TestUpdateAnnouncement_NotFound(t *testing.T) {
	ctx := context.Background()
	repoErr := errors.New("not found")
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(_ context.Context, _ uint) (*entities.Announcement, error) {
			return nil, repoErr
		},
	}
	uc := newUseCase(repo, &mocks.StorageMock{})

	dto := dtos.UpdateAnnouncementDTO{
		ID:            99,
		DisplayType:   entities.DisplayTypeModalText,
		FrequencyType: entities.FrequencyAlways,
		IsGlobal:      true,
	}

	_, err := uc.UpdateAnnouncement(ctx, dto)
	if !errors.Is(err, repoErr) {
		t.Errorf("expected not found error, got %v", err)
	}
}

func TestUpdateAnnouncement_InvalidDisplayType(t *testing.T) {
	ctx := context.Background()
	uc := newUseCase(&mocks.RepositoryMock{}, &mocks.StorageMock{})

	dto := dtos.UpdateAnnouncementDTO{
		ID:            1,
		DisplayType:   "bad_type",
		FrequencyType: entities.FrequencyAlways,
		IsGlobal:      true,
	}

	_, err := uc.UpdateAnnouncement(ctx, dto)
	if !errors.Is(err, domainerrors.ErrInvalidDisplayType) {
		t.Errorf("expected ErrInvalidDisplayType, got %v", err)
	}
}

func TestUpdateAnnouncement_TargetsRequiredWhenNotGlobal(t *testing.T) {
	ctx := context.Background()
	uc := newUseCase(&mocks.RepositoryMock{}, &mocks.StorageMock{})

	dto := dtos.UpdateAnnouncementDTO{
		ID:            1,
		DisplayType:   entities.DisplayTypeModalText,
		FrequencyType: entities.FrequencyAlways,
		IsGlobal:      false,
		TargetIDs:     nil,
	}

	_, err := uc.UpdateAnnouncement(ctx, dto)
	if !errors.Is(err, domainerrors.ErrTargetsRequired) {
		t.Errorf("expected ErrTargetsRequired, got %v", err)
	}
}

func TestDeleteAnnouncement_Success(t *testing.T) {
	ctx := context.Background()
	imageURL := "https://s3.example.com/img.png"
	existing := &entities.Announcement{
		ID:     1,
		Images: []entities.AnnouncementImage{{ID: 1, ImageURL: imageURL}},
	}

	deletedURLs := []string{}
	storage := &mocks.StorageMock{
		DeleteFileFn: func(_ context.Context, url string) error {
			deletedURLs = append(deletedURLs, url)
			return nil
		},
	}
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(_ context.Context, _ uint) (*entities.Announcement, error) {
			return existing, nil
		},
		DeleteFn: func(_ context.Context, _ uint) error {
			return nil
		},
	}
	uc := newUseCase(repo, storage)

	err := uc.DeleteAnnouncement(ctx, 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(deletedURLs) != 1 || deletedURLs[0] != imageURL {
		t.Errorf("expected S3 delete for %s, got %v", imageURL, deletedURLs)
	}
}

func TestDeleteAnnouncement_NotFound(t *testing.T) {
	ctx := context.Background()
	repoErr := errors.New("not found")
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(_ context.Context, _ uint) (*entities.Announcement, error) {
			return nil, repoErr
		},
	}
	uc := newUseCase(repo, &mocks.StorageMock{})

	err := uc.DeleteAnnouncement(ctx, 99)
	if !errors.Is(err, repoErr) {
		t.Errorf("expected not found error, got %v", err)
	}
}

func TestDeleteAnnouncement_S3ErrorDoesNotStopDeletion(t *testing.T) {
	ctx := context.Background()
	existing := &entities.Announcement{
		ID:     1,
		Images: []entities.AnnouncementImage{{ID: 1, ImageURL: "https://s3.example.com/img.png"}},
	}
	storage := &mocks.StorageMock{
		DeleteFileFn: func(_ context.Context, _ string) error {
			return errors.New("s3 error")
		},
	}
	repoDeleteCalled := false
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(_ context.Context, _ uint) (*entities.Announcement, error) {
			return existing, nil
		},
		DeleteFn: func(_ context.Context, _ uint) error {
			repoDeleteCalled = true
			return nil
		},
	}
	uc := newUseCase(repo, storage)

	err := uc.DeleteAnnouncement(ctx, 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !repoDeleteCalled {
		t.Error("expected repo.Delete to be called despite S3 error")
	}
}

func TestDeleteAnnouncement_RepositoryDeleteError(t *testing.T) {
	ctx := context.Background()
	repoErr := errors.New("db error")
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(_ context.Context, _ uint) (*entities.Announcement, error) {
			return &entities.Announcement{ID: 1}, nil
		},
		DeleteFn: func(_ context.Context, _ uint) error {
			return repoErr
		},
	}
	uc := newUseCase(repo, &mocks.StorageMock{})

	err := uc.DeleteAnnouncement(ctx, 1)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetAnnouncement_Success(t *testing.T) {
	ctx := context.Background()
	expected := &entities.Announcement{ID: 1, Title: "Hello"}
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(_ context.Context, id uint) (*entities.Announcement, error) {
			return expected, nil
		},
	}
	uc := newUseCase(repo, &mocks.StorageMock{})

	result, err := uc.GetAnnouncement(ctx, 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.ID != expected.ID {
		t.Errorf("expected ID %d, got %d", expected.ID, result.ID)
	}
}

func TestListAnnouncements_NormalizesPage(t *testing.T) {
	ctx := context.Background()
	capturedParams := dtos.ListAnnouncementsParams{}
	repo := &mocks.RepositoryMock{
		ListFn: func(_ context.Context, params dtos.ListAnnouncementsParams) ([]entities.Announcement, int64, error) {
			capturedParams = params
			return nil, 0, nil
		},
	}
	uc := newUseCase(repo, &mocks.StorageMock{})

	_, _, err := uc.ListAnnouncements(ctx, dtos.ListAnnouncementsParams{Page: 0, PageSize: 0})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if capturedParams.Page != 1 {
		t.Errorf("expected page normalized to 1, got %d", capturedParams.Page)
	}
	if capturedParams.PageSize != 20 {
		t.Errorf("expected pageSize normalized to 20, got %d", capturedParams.PageSize)
	}
}

func TestListAnnouncements_CapsBigPageSize(t *testing.T) {
	ctx := context.Background()
	capturedParams := dtos.ListAnnouncementsParams{}
	repo := &mocks.RepositoryMock{
		ListFn: func(_ context.Context, params dtos.ListAnnouncementsParams) ([]entities.Announcement, int64, error) {
			capturedParams = params
			return nil, 0, nil
		},
	}
	uc := newUseCase(repo, &mocks.StorageMock{})

	_, _, err := uc.ListAnnouncements(ctx, dtos.ListAnnouncementsParams{Page: 1, PageSize: 999})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if capturedParams.PageSize != 20 {
		t.Errorf("expected pageSize normalized to 20, got %d", capturedParams.PageSize)
	}
}

func TestGetActiveAnnouncements_FrequencyAlwaysVisible(t *testing.T) {
	ctx := context.Background()
	announcements := []entities.Announcement{
		{ID: 1, FrequencyType: entities.FrequencyAlways},
		{ID: 2, FrequencyType: entities.FrequencyAlways},
	}
	repo := &mocks.RepositoryMock{
		GetActiveAnnouncementsFn: func(_ context.Context, _ dtos.ActiveAnnouncementsParams) ([]entities.Announcement, error) {
			return announcements, nil
		},
	}
	uc := newUseCase(repo, &mocks.StorageMock{})

	result, err := uc.GetActiveAnnouncements(ctx, dtos.ActiveAnnouncementsParams{BusinessID: 1, UserID: 10})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 announcements, got %d", len(result))
	}
}

func TestGetActiveAnnouncements_FrequencyOnce_HiddenAfterView(t *testing.T) {
	ctx := context.Background()
	announcements := []entities.Announcement{
		{ID: 1, FrequencyType: entities.FrequencyOnce},
	}
	repo := &mocks.RepositoryMock{
		GetActiveAnnouncementsFn: func(_ context.Context, _ dtos.ActiveAnnouncementsParams) ([]entities.Announcement, error) {
			return announcements, nil
		},
		GetUserViewsFn: func(_ context.Context, _, _ uint) ([]entities.AnnouncementView, error) {
			return []entities.AnnouncementView{{ID: 1, Action: entities.ViewActionViewed}}, nil
		},
	}
	uc := newUseCase(repo, &mocks.StorageMock{})

	result, err := uc.GetActiveAnnouncements(ctx, dtos.ActiveAnnouncementsParams{BusinessID: 1, UserID: 10})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected 0 announcements (once already viewed), got %d", len(result))
	}
}

func TestGetActiveAnnouncements_FrequencyOnce_VisibleWithoutViews(t *testing.T) {
	ctx := context.Background()
	announcements := []entities.Announcement{
		{ID: 1, FrequencyType: entities.FrequencyOnce},
	}
	repo := &mocks.RepositoryMock{
		GetActiveAnnouncementsFn: func(_ context.Context, _ dtos.ActiveAnnouncementsParams) ([]entities.Announcement, error) {
			return announcements, nil
		},
		GetUserViewsFn: func(_ context.Context, _, _ uint) ([]entities.AnnouncementView, error) {
			return []entities.AnnouncementView{}, nil
		},
	}
	uc := newUseCase(repo, &mocks.StorageMock{})

	result, err := uc.GetActiveAnnouncements(ctx, dtos.ActiveAnnouncementsParams{BusinessID: 1, UserID: 10})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 announcement, got %d", len(result))
	}
}

func TestGetActiveAnnouncements_FrequencyRequiresAcceptance_HiddenAfterAccept(t *testing.T) {
	ctx := context.Background()
	announcements := []entities.Announcement{
		{ID: 1, FrequencyType: entities.FrequencyRequiresAcceptance},
	}
	repo := &mocks.RepositoryMock{
		GetActiveAnnouncementsFn: func(_ context.Context, _ dtos.ActiveAnnouncementsParams) ([]entities.Announcement, error) {
			return announcements, nil
		},
		GetUserViewsFn: func(_ context.Context, _, _ uint) ([]entities.AnnouncementView, error) {
			return []entities.AnnouncementView{{ID: 1, Action: entities.ViewActionAccepted}}, nil
		},
	}
	uc := newUseCase(repo, &mocks.StorageMock{})

	result, err := uc.GetActiveAnnouncements(ctx, dtos.ActiveAnnouncementsParams{BusinessID: 1, UserID: 10})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected 0 (accepted), got %d", len(result))
	}
}

func TestGetActiveAnnouncements_FrequencyRequiresAcceptance_VisibleWithoutAccept(t *testing.T) {
	ctx := context.Background()
	announcements := []entities.Announcement{
		{ID: 1, FrequencyType: entities.FrequencyRequiresAcceptance},
	}
	repo := &mocks.RepositoryMock{
		GetActiveAnnouncementsFn: func(_ context.Context, _ dtos.ActiveAnnouncementsParams) ([]entities.Announcement, error) {
			return announcements, nil
		},
		GetUserViewsFn: func(_ context.Context, _, _ uint) ([]entities.AnnouncementView, error) {
			return []entities.AnnouncementView{{ID: 1, Action: entities.ViewActionViewed}}, nil
		},
	}
	uc := newUseCase(repo, &mocks.StorageMock{})

	result, err := uc.GetActiveAnnouncements(ctx, dtos.ActiveAnnouncementsParams{BusinessID: 1, UserID: 10})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 (not yet accepted), got %d", len(result))
	}
}

func TestGetActiveAnnouncements_FrequencyDaily_HiddenIfViewedToday(t *testing.T) {
	ctx := context.Background()
	announcements := []entities.Announcement{
		{ID: 1, FrequencyType: entities.FrequencyDaily},
	}
	repo := &mocks.RepositoryMock{
		GetActiveAnnouncementsFn: func(_ context.Context, _ dtos.ActiveAnnouncementsParams) ([]entities.Announcement, error) {
			return announcements, nil
		},
		GetUserViewsFn: func(_ context.Context, _, _ uint) ([]entities.AnnouncementView, error) {
			return []entities.AnnouncementView{{ID: 1, ViewedAt: time.Now()}}, nil
		},
	}
	uc := newUseCase(repo, &mocks.StorageMock{})

	result, err := uc.GetActiveAnnouncements(ctx, dtos.ActiveAnnouncementsParams{BusinessID: 1, UserID: 10})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected 0 (viewed today), got %d", len(result))
	}
}

func TestGetActiveAnnouncements_FrequencyDaily_VisibleIfViewedYesterday(t *testing.T) {
	ctx := context.Background()
	announcements := []entities.Announcement{
		{ID: 1, FrequencyType: entities.FrequencyDaily},
	}
	yesterday := time.Now().Add(-25 * time.Hour)
	repo := &mocks.RepositoryMock{
		GetActiveAnnouncementsFn: func(_ context.Context, _ dtos.ActiveAnnouncementsParams) ([]entities.Announcement, error) {
			return announcements, nil
		},
		GetUserViewsFn: func(_ context.Context, _, _ uint) ([]entities.AnnouncementView, error) {
			return []entities.AnnouncementView{{ID: 1, ViewedAt: yesterday}}, nil
		},
	}
	uc := newUseCase(repo, &mocks.StorageMock{})

	result, err := uc.GetActiveAnnouncements(ctx, dtos.ActiveAnnouncementsParams{BusinessID: 1, UserID: 10})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 (viewed yesterday), got %d", len(result))
	}
}

func TestGetActiveAnnouncements_SkipsOnViewError(t *testing.T) {
	ctx := context.Background()
	announcements := []entities.Announcement{
		{ID: 1, FrequencyType: entities.FrequencyOnce},
	}
	repo := &mocks.RepositoryMock{
		GetActiveAnnouncementsFn: func(_ context.Context, _ dtos.ActiveAnnouncementsParams) ([]entities.Announcement, error) {
			return announcements, nil
		},
		GetUserViewsFn: func(_ context.Context, _, _ uint) ([]entities.AnnouncementView, error) {
			return nil, errors.New("db error")
		},
	}
	uc := newUseCase(repo, &mocks.StorageMock{})

	result, err := uc.GetActiveAnnouncements(ctx, dtos.ActiveAnnouncementsParams{BusinessID: 1, UserID: 10})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected 0 (skipped due to error), got %d", len(result))
	}
}

func TestChangeStatus_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(_ context.Context, _ uint) (*entities.Announcement, error) {
			return &entities.Announcement{ID: 1}, nil
		},
		ChangeStatusFn: func(_ context.Context, _ uint, _ entities.AnnouncementStatus) error {
			return nil
		},
	}
	uc := newUseCase(repo, &mocks.StorageMock{})

	err := uc.ChangeStatus(ctx, dtos.ChangeStatusDTO{ID: 1, Status: entities.StatusInactive})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestChangeStatus_InvalidStatus(t *testing.T) {
	ctx := context.Background()
	uc := newUseCase(&mocks.RepositoryMock{}, &mocks.StorageMock{})

	err := uc.ChangeStatus(ctx, dtos.ChangeStatusDTO{ID: 1, Status: "nonexistent"})
	if !errors.Is(err, domainerrors.ErrInvalidStatus) {
		t.Errorf("expected ErrInvalidStatus, got %v", err)
	}
}

func TestChangeStatus_AllValidStatuses(t *testing.T) {
	statuses := []entities.AnnouncementStatus{
		entities.StatusDraft,
		entities.StatusScheduled,
		entities.StatusActive,
		entities.StatusInactive,
	}

	for _, status := range statuses {
		t.Run(string(status), func(t *testing.T) {
			ctx := context.Background()
			repo := &mocks.RepositoryMock{
				GetByIDFn: func(_ context.Context, _ uint) (*entities.Announcement, error) {
					return &entities.Announcement{ID: 1}, nil
				},
				ChangeStatusFn: func(_ context.Context, _ uint, _ entities.AnnouncementStatus) error {
					return nil
				},
			}
			uc := newUseCase(repo, &mocks.StorageMock{})

			err := uc.ChangeStatus(ctx, dtos.ChangeStatusDTO{ID: 1, Status: status})
			if err != nil {
				t.Errorf("expected no error for status %s, got %v", status, err)
			}
		})
	}
}

func TestChangeStatus_NotFound(t *testing.T) {
	ctx := context.Background()
	repoErr := errors.New("not found")
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(_ context.Context, _ uint) (*entities.Announcement, error) {
			return nil, repoErr
		},
	}
	uc := newUseCase(repo, &mocks.StorageMock{})

	err := uc.ChangeStatus(ctx, dtos.ChangeStatusDTO{ID: 99, Status: entities.StatusActive})
	if !errors.Is(err, repoErr) {
		t.Errorf("expected not found error, got %v", err)
	}
}

func TestForceRedisplay_Success(t *testing.T) {
	ctx := context.Background()
	existing := &entities.Announcement{ID: 1, ForceRedisplay: false}

	updatedAnnouncement := (*entities.Announcement)(nil)
	deleteViewsCalled := false
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(_ context.Context, _ uint) (*entities.Announcement, error) {
			return existing, nil
		},
		UpdateFn: func(_ context.Context, a *entities.Announcement) (*entities.Announcement, error) {
			updatedAnnouncement = a
			return a, nil
		},
		DeleteViewsByAnnouncementIDFn: func(_ context.Context, _ uint) error {
			deleteViewsCalled = true
			return nil
		},
	}
	uc := newUseCase(repo, &mocks.StorageMock{})

	err := uc.ForceRedisplay(ctx, 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if updatedAnnouncement == nil || !updatedAnnouncement.ForceRedisplay {
		t.Error("expected ForceRedisplay to be set to true")
	}
	if !deleteViewsCalled {
		t.Error("expected DeleteViewsByAnnouncementID to be called")
	}
}

func TestForceRedisplay_NotFound(t *testing.T) {
	ctx := context.Background()
	repoErr := errors.New("not found")
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(_ context.Context, _ uint) (*entities.Announcement, error) {
			return nil, repoErr
		},
	}
	uc := newUseCase(repo, &mocks.StorageMock{})

	err := uc.ForceRedisplay(ctx, 99)
	if !errors.Is(err, repoErr) {
		t.Errorf("expected not found error, got %v", err)
	}
}

func TestForceRedisplay_UpdateError(t *testing.T) {
	ctx := context.Background()
	updateErr := errors.New("update error")
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(_ context.Context, _ uint) (*entities.Announcement, error) {
			return &entities.Announcement{ID: 1}, nil
		},
		UpdateFn: func(_ context.Context, _ *entities.Announcement) (*entities.Announcement, error) {
			return nil, updateErr
		},
	}
	uc := newUseCase(repo, &mocks.StorageMock{})

	err := uc.ForceRedisplay(ctx, 1)
	if !errors.Is(err, updateErr) {
		t.Errorf("expected update error, got %v", err)
	}
}

func TestRegisterView_Success(t *testing.T) {
	ctx := context.Background()
	var capturedView *entities.AnnouncementView
	repo := &mocks.RepositoryMock{
		RegisterViewFn: func(_ context.Context, view *entities.AnnouncementView) error {
			capturedView = view
			return nil
		},
	}
	uc := newUseCase(repo, &mocks.StorageMock{})

	dto := dtos.RegisterViewDTO{
		AnnouncementID: 5,
		UserID:         10,
		BusinessID:     20,
		Action:         entities.ViewActionViewed,
	}

	err := uc.RegisterView(ctx, dto)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if capturedView == nil {
		t.Fatal("expected view to be registered")
	}
	if capturedView.AnnouncementID != dto.AnnouncementID {
		t.Errorf("expected AnnouncementID %d, got %d", dto.AnnouncementID, capturedView.AnnouncementID)
	}
	if capturedView.UserID != dto.UserID {
		t.Errorf("expected UserID %d, got %d", dto.UserID, capturedView.UserID)
	}
	if capturedView.ViewedAt.IsZero() {
		t.Error("expected ViewedAt to be set")
	}
}

func TestGetAnnouncementStats_Success(t *testing.T) {
	ctx := context.Background()
	expected := &entities.AnnouncementStats{TotalViews: 100, UniqueUsers: 50}
	repo := &mocks.RepositoryMock{
		GetStatsFn: func(_ context.Context, id uint) (*entities.AnnouncementStats, error) {
			return expected, nil
		},
	}
	uc := newUseCase(repo, &mocks.StorageMock{})

	result, err := uc.GetAnnouncementStats(ctx, 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.TotalViews != expected.TotalViews {
		t.Errorf("expected TotalViews %d, got %d", expected.TotalViews, result.TotalViews)
	}
}

func TestListCategories_Success(t *testing.T) {
	ctx := context.Background()
	expected := []entities.AnnouncementCategory{
		{ID: 1, Code: "promo", Name: "Promo"},
		{ID: 2, Code: "alert", Name: "Alert"},
	}
	repo := &mocks.RepositoryMock{
		ListCategoriesFn: func(_ context.Context) ([]entities.AnnouncementCategory, error) {
			return expected, nil
		},
	}
	uc := newUseCase(repo, &mocks.StorageMock{})

	result, err := uc.ListCategories(ctx)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != len(expected) {
		t.Errorf("expected %d categories, got %d", len(expected), len(result))
	}
}
