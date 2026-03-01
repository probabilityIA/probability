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

func TestCreateHandler_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)

	expectedResponse := &dtos.NotificationConfigResponseDTO{
		ID:                      1,
		IntegrationID:           100,
		NotificationTypeID:      1,
		NotificationEventTypeID: 1,
		Enabled:                 true,
		Description:             "Test config",
	}

	mockUseCase := &mocks.UseCaseMock{
		CreateFn: func(ctx context.Context, dto dtos.CreateNotificationConfigDTO) (*dtos.NotificationConfigResponseDTO, error) {
			return expectedResponse, nil
		},
	}
	mockLogger := mocks.NewLoggerMock()

	handler := New(mockUseCase, mockLogger)

	requestBody := map[string]interface{}{
		"integration_id":             100,
		"notification_type_id":       1,
		"notification_event_type_id": 1,
		"enabled":                    true,
		"description":                "Test config",
		"order_status_ids":           []uint{1},
	}

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/notification-configs", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("business_id", uint(10)) // Simular middleware de autenticación

	// Act
	handler.Create(c)

	// Assert
	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["id"].(float64) != float64(expectedResponse.ID) {
		t.Errorf("expected ID %d, got %v", expectedResponse.ID, response["id"])
	}

	if response["integration_id"].(float64) != float64(expectedResponse.IntegrationID) {
		t.Errorf("expected IntegrationID %d, got %v", expectedResponse.IntegrationID, response["integration_id"])
	}
}

func TestCreateHandler_InvalidRequest(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)

	mockUseCase := &mocks.UseCaseMock{}
	mockLogger := mocks.NewLoggerMock()

	handler := New(mockUseCase, mockLogger)

	// Request body inválido (JSON malformado)
	invalidBody := []byte(`{"integration_id": "invalid"}`)
	req := httptest.NewRequest(http.MethodPost, "/notification-configs", bytes.NewBuffer(invalidBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Act
	handler.Create(c)

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

func TestCreateHandler_DuplicateConfig(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)

	mockUseCase := &mocks.UseCaseMock{
		CreateFn: func(ctx context.Context, dto dtos.CreateNotificationConfigDTO) (*dtos.NotificationConfigResponseDTO, error) {
			return nil, domainErrors.ErrDuplicateConfig
		},
	}
	mockLogger := mocks.NewLoggerMock()

	handler := New(mockUseCase, mockLogger)

	requestBody := map[string]interface{}{
		"integration_id":             100,
		"notification_type_id":       1,
		"notification_event_type_id": 1,
		"enabled":                    true,
		"description":                "Test config",
	}

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/notification-configs", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("business_id", uint(10)) // Simular middleware de autenticación

	// Act
	handler.Create(c)

	// Assert
	if w.Code != http.StatusConflict {
		t.Errorf("expected status %d, got %d", http.StatusConflict, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if _, ok := response["error"]; !ok {
		t.Error("expected error field in response")
	}
}

func TestCreateHandler_UseCaseError(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)

	mockUseCase := &mocks.UseCaseMock{
		CreateFn: func(ctx context.Context, dto dtos.CreateNotificationConfigDTO) (*dtos.NotificationConfigResponseDTO, error) {
			return nil, errors.New("database connection failed")
		},
	}
	mockLogger := mocks.NewLoggerMock()

	handler := New(mockUseCase, mockLogger)

	requestBody := map[string]interface{}{
		"integration_id":             100,
		"notification_type_id":       1,
		"notification_event_type_id": 1,
		"enabled":                    true,
		"description":                "Test config",
	}

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/notification-configs", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("business_id", uint(10)) // Simular middleware de autenticación

	// Act
	handler.Create(c)

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
