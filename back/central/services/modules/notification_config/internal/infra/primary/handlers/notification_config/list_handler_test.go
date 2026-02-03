package notification_config

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/mocks"
)

func TestListHandler_Success_NoFilters(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)

	expectedConfigs := []dtos.NotificationConfigResponseDTO{
		{
			ID:                      1,
			IntegrationID:           100,
			NotificationTypeID:      1,
			NotificationEventTypeID: 1,
			Enabled:                 true,
		},
		{
			ID:                      2,
			IntegrationID:           200,
			NotificationTypeID:      2,
			NotificationEventTypeID: 2,
			Enabled:                 false,
		},
	}

	mockUseCase := &mocks.UseCaseMock{
		ListFn: func(ctx context.Context, filters dtos.FilterNotificationConfigDTO) ([]dtos.NotificationConfigResponseDTO, error) {
			return expectedConfigs, nil
		},
	}
	mockLogger := mocks.NewLoggerMock()

	handler := New(mockUseCase, mockLogger)

	req := httptest.NewRequest(http.MethodGet, "/notification-configs", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Act
	handler.List(c)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if len(response) != len(expectedConfigs) {
		t.Errorf("expected %d configs, got %d", len(expectedConfigs), len(response))
	}
}

func TestListHandler_Success_WithFilters(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)

	expectedConfigs := []dtos.NotificationConfigResponseDTO{
		{
			ID:                      1,
			IntegrationID:           100,
			NotificationTypeID:      1,
			NotificationEventTypeID: 1,
			Enabled:                 true,
		},
	}

	mockUseCase := &mocks.UseCaseMock{
		ListFn: func(ctx context.Context, filters dtos.FilterNotificationConfigDTO) ([]dtos.NotificationConfigResponseDTO, error) {
			// Verificar que los filtros se pasaron correctamente
			if filters.IntegrationID != nil && *filters.IntegrationID != 100 {
				t.Errorf("expected IntegrationID 100, got %d", *filters.IntegrationID)
			}
			if filters.NotificationTypeID != nil && *filters.NotificationTypeID != 1 {
				t.Errorf("expected NotificationTypeID 1, got %d", *filters.NotificationTypeID)
			}
			return expectedConfigs, nil
		},
	}
	mockLogger := mocks.NewLoggerMock()

	handler := New(mockUseCase, mockLogger)

	req := httptest.NewRequest(http.MethodGet, "/notification-configs?integration_id=100&notification_type_id=1&enabled=true", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Act
	handler.List(c)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if len(response) != 1 {
		t.Errorf("expected 1 config, got %d", len(response))
	}
}

func TestListHandler_EmptyResult(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)

	mockUseCase := &mocks.UseCaseMock{
		ListFn: func(ctx context.Context, filters dtos.FilterNotificationConfigDTO) ([]dtos.NotificationConfigResponseDTO, error) {
			return []dtos.NotificationConfigResponseDTO{}, nil
		},
	}
	mockLogger := mocks.NewLoggerMock()

	handler := New(mockUseCase, mockLogger)

	req := httptest.NewRequest(http.MethodGet, "/notification-configs", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Act
	handler.List(c)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if len(response) != 0 {
		t.Errorf("expected 0 configs, got %d", len(response))
	}
}

func TestListHandler_UseCaseError(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)

	mockUseCase := &mocks.UseCaseMock{
		ListFn: func(ctx context.Context, filters dtos.FilterNotificationConfigDTO) ([]dtos.NotificationConfigResponseDTO, error) {
			return nil, errors.New("database connection failed")
		},
	}
	mockLogger := mocks.NewLoggerMock()

	handler := New(mockUseCase, mockLogger)

	req := httptest.NewRequest(http.MethodGet, "/notification-configs", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Act
	handler.List(c)

	// Assert
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if _, ok := response["error"]; !ok {
		t.Error("expected error field in response")
	}
}
