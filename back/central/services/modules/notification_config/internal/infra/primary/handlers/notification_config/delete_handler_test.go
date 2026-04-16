package notification_config

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	domainErrors "github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/mocks"
)

func TestDeleteHandler_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)

	mockUseCase := &mocks.UseCaseMock{
		DeleteFn: func(ctx context.Context, id uint) error {
			if id == 1 {
				return nil
			}
			return domainErrors.ErrNotificationConfigNotFound
		},
	}
	mockLogger := mocks.NewLoggerMock()

	handler := New(mockUseCase, mockLogger)

	req := httptest.NewRequest(http.MethodDelete, "/notification-configs/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	// Act
	handler.Delete(c)

	// Assert
	// Gin en modo test puede retornar 200 en lugar de 204
	if w.Code != http.StatusNoContent && w.Code != http.StatusOK {
		t.Errorf("expected status %d or %d, got %d", http.StatusNoContent, http.StatusOK, w.Code)
	}

	// Si es 204, debe tener body vacío
	if w.Code == http.StatusNoContent && w.Body.Len() != 0 {
		t.Error("expected empty body for 204 response")
	}
}

func TestDeleteHandler_InvalidID(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)

	mockUseCase := &mocks.UseCaseMock{}
	mockLogger := mocks.NewLoggerMock()

	handler := New(mockUseCase, mockLogger)

	req := httptest.NewRequest(http.MethodDelete, "/notification-configs/invalid", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "invalid"}}

	// Act
	handler.Delete(c)

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

func TestDeleteHandler_NotFound(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)

	mockUseCase := &mocks.UseCaseMock{
		DeleteFn: func(ctx context.Context, id uint) error {
			return domainErrors.ErrNotificationConfigNotFound
		},
	}
	mockLogger := mocks.NewLoggerMock()

	handler := New(mockUseCase, mockLogger)

	req := httptest.NewRequest(http.MethodDelete, "/notification-configs/999", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "999"}}

	// Act
	handler.Delete(c)

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

func TestDeleteHandler_UseCaseError(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)

	mockUseCase := &mocks.UseCaseMock{
		DeleteFn: func(ctx context.Context, id uint) error {
			return errors.New("database delete failed")
		},
	}
	mockLogger := mocks.NewLoggerMock()

	handler := New(mockUseCase, mockLogger)

	req := httptest.NewRequest(http.MethodDelete, "/notification-configs/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	// Act
	handler.Delete(c)

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
