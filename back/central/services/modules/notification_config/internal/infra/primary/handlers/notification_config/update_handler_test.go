package notification_config

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
	domainErrors "github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/mocks"
)

func TestUpdateHandler_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)

	expectedResponse := &dtos.NotificationConfigResponseDTO{
		ID:                      1,
		IntegrationID:           100,
		NotificationTypeID:      1,
		NotificationEventTypeID: 1,
		Enabled:                 false,
		Description:             "Updated config",
	}

	mockUseCase := &mocks.UseCaseMock{
		UpdateFn: func(ctx context.Context, id uint, dto dtos.UpdateNotificationConfigDTO) (*dtos.NotificationConfigResponseDTO, error) {
			if id == 1 {
				return expectedResponse, nil
			}
			return nil, domainErrors.ErrNotificationConfigNotFound
		},
	}
	mockLogger := mocks.NewLoggerMock()

	handler := New(mockUseCase, mockLogger)

	requestBody := map[string]interface{}{
		"notification_type": "email",
		"is_active":         false,
		"description":       "Updated config",
		"priority":          2,
	}

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPut, "/notification-configs/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	// Act
	handler.Update(c)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["id"].(float64) != float64(expectedResponse.ID) {
		t.Errorf("expected ID %d, got %v", expectedResponse.ID, response["id"])
	}

	if response["enabled"].(bool) != expectedResponse.Enabled {
		t.Errorf("expected Enabled %v, got %v", expectedResponse.Enabled, response["enabled"])
	}
}

func TestUpdateHandler_InvalidID(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)

	mockUseCase := &mocks.UseCaseMock{}
	mockLogger := mocks.NewLoggerMock()

	handler := New(mockUseCase, mockLogger)

	requestBody := map[string]interface{}{
		"is_active": false,
	}

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPut, "/notification-configs/invalid", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "invalid"}}

	// Act
	handler.Update(c)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestUpdateHandler_InvalidRequest(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)

	mockUseCase := &mocks.UseCaseMock{}
	mockLogger := mocks.NewLoggerMock()

	handler := New(mockUseCase, mockLogger)

	invalidBody := []byte(`{"is_active": "invalid"}`)
	req := httptest.NewRequest(http.MethodPut, "/notification-configs/1", bytes.NewBuffer(invalidBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	// Act
	handler.Update(c)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestUpdateHandler_NotFound(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)

	mockUseCase := &mocks.UseCaseMock{
		UpdateFn: func(ctx context.Context, id uint, dto dtos.UpdateNotificationConfigDTO) (*dtos.NotificationConfigResponseDTO, error) {
			return nil, domainErrors.ErrNotificationConfigNotFound
		},
	}
	mockLogger := mocks.NewLoggerMock()

	handler := New(mockUseCase, mockLogger)

	requestBody := map[string]interface{}{
		"is_active": false,
	}

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPut, "/notification-configs/999", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "999"}}

	// Act
	handler.Update(c)

	// Assert
	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestUpdateHandler_UseCaseError(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)

	mockUseCase := &mocks.UseCaseMock{
		UpdateFn: func(ctx context.Context, id uint, dto dtos.UpdateNotificationConfigDTO) (*dtos.NotificationConfigResponseDTO, error) {
			return nil, errors.New("database update failed")
		},
	}
	mockLogger := mocks.NewLoggerMock()

	handler := New(mockUseCase, mockLogger)

	requestBody := map[string]interface{}{
		"is_active": false,
	}

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPut, "/notification-configs/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	// Act
	handler.Update(c)

	// Assert
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}
