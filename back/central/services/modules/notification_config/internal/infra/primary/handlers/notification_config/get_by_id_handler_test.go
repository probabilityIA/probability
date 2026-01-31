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
	domainErrors "github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/mocks"
)

func TestGetByIDHandler_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)

	expectedConfig := &dtos.NotificationConfigResponseDTO{
		ID:               1,
		IntegrationID:    100,
		NotificationType: "whatsapp",
		IsActive:         true,
		Description:      "Test config",
		Priority:         1,
	}

	mockUseCase := &mocks.UseCaseMock{
		GetByIDFn: func(ctx context.Context, id uint) (*dtos.NotificationConfigResponseDTO, error) {
			if id == 1 {
				return expectedConfig, nil
			}
			return nil, domainErrors.ErrNotificationConfigNotFound
		},
	}
	mockLogger := mocks.NewLoggerMock()

	handler := New(mockUseCase, mockLogger)

	req := httptest.NewRequest(http.MethodGet, "/notification-configs/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	// Act
	handler.GetByID(c)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["id"].(float64) != float64(expectedConfig.ID) {
		t.Errorf("expected ID %d, got %v", expectedConfig.ID, response["id"])
	}

	if response["integration_id"].(float64) != float64(expectedConfig.IntegrationID) {
		t.Errorf("expected IntegrationID %d, got %v", expectedConfig.IntegrationID, response["integration_id"])
	}
}

func TestGetByIDHandler_InvalidID(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)

	mockUseCase := &mocks.UseCaseMock{}
	mockLogger := mocks.NewLoggerMock()

	handler := New(mockUseCase, mockLogger)

	req := httptest.NewRequest(http.MethodGet, "/notification-configs/invalid", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "invalid"}}

	// Act
	handler.GetByID(c)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if _, ok := response["error"]; !ok {
		t.Error("expected error field in response")
	}
}

func TestGetByIDHandler_NotFound(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)

	mockUseCase := &mocks.UseCaseMock{
		GetByIDFn: func(ctx context.Context, id uint) (*dtos.NotificationConfigResponseDTO, error) {
			return nil, domainErrors.ErrNotificationConfigNotFound
		},
	}
	mockLogger := mocks.NewLoggerMock()

	handler := New(mockUseCase, mockLogger)

	req := httptest.NewRequest(http.MethodGet, "/notification-configs/999", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "999"}}

	// Act
	handler.GetByID(c)

	// Assert
	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if _, ok := response["error"]; !ok {
		t.Error("expected error field in response")
	}
}

func TestGetByIDHandler_UseCaseError(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)

	mockUseCase := &mocks.UseCaseMock{
		GetByIDFn: func(ctx context.Context, id uint) (*dtos.NotificationConfigResponseDTO, error) {
			return nil, errors.New("database connection failed")
		},
	}
	mockLogger := mocks.NewLoggerMock()

	handler := New(mockUseCase, mockLogger)

	req := httptest.NewRequest(http.MethodGet, "/notification-configs/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	// Act
	handler.GetByID(c)

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
