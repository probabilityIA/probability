package app

import (
	"context"
	"errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/mocks"
)

func TestDelete_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()

	mockRepo := &mocks.RepositoryMock{
		DeleteFn: func(ctx context.Context, id uint) error {
			if id == 1 {
				return nil
			}
			return errors.New("not found")
		},
	}
	mockNotificationTypeRepo := &mocks.NotificationTypeRepositoryMock{}
	mockEventTypeRepo := &mocks.NotificationEventTypeRepositoryMock{}
	mockCacheManager := &mocks.CacheManagerMock{}
	mockLogger := mocks.NewLoggerMock()

	useCase := New(mockRepo, mockNotificationTypeRepo, mockEventTypeRepo, mockCacheManager, &mocks.MessageAuditQuerierMock{}, mockLogger)

	// Act
	err := useCase.Delete(ctx, 1)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestDelete_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	expectedErr := errors.New("record not found")

	mockRepo := &mocks.RepositoryMock{
		DeleteFn: func(ctx context.Context, id uint) error {
			return expectedErr
		},
	}
	mockNotificationTypeRepo := &mocks.NotificationTypeRepositoryMock{}
	mockEventTypeRepo := &mocks.NotificationEventTypeRepositoryMock{}
	mockCacheManager := &mocks.CacheManagerMock{}
	mockLogger := mocks.NewLoggerMock()

	useCase := New(mockRepo, mockNotificationTypeRepo, mockEventTypeRepo, mockCacheManager, &mocks.MessageAuditQuerierMock{}, mockLogger)

	// Act
	err := useCase.Delete(ctx, 999)

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, expectedErr) {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}
}

func TestDelete_RepositoryError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	expectedErr := errors.New("database delete failed")

	mockRepo := &mocks.RepositoryMock{
		DeleteFn: func(ctx context.Context, id uint) error {
			return expectedErr
		},
	}
	mockNotificationTypeRepo := &mocks.NotificationTypeRepositoryMock{}
	mockEventTypeRepo := &mocks.NotificationEventTypeRepositoryMock{}
	mockCacheManager := &mocks.CacheManagerMock{}
	mockLogger := mocks.NewLoggerMock()

	useCase := New(mockRepo, mockNotificationTypeRepo, mockEventTypeRepo, mockCacheManager, &mocks.MessageAuditQuerierMock{}, mockLogger)

	// Act
	err := useCase.Delete(ctx, 1)

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, expectedErr) {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}
}
