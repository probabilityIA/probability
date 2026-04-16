package permissionhandler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/permissions/internal/domain"
	"github.com/secamc93/probability/back/central/services/auth/permissions/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTestRouter(handler IPermissionHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/permissions/bulk", handler.BulkCreatePermissionsHandler)
	return router
}

func TestBulkCreatePermissionsHandler_Success(t *testing.T) {
	// Arrange
	mockUseCase := new(mocks.UseCaseMock)
	mockLogger := new(mocks.LoggerMock)

	handler := New(mockUseCase, mockLogger)
	router := setupTestRouter(handler)

	mockLogger.On("Info", mock.Anything).Maybe()
	mockLogger.On("Error", mock.Anything).Maybe()

	expectedResults := []domain.BulkCreateResult{
		{Name: "Crear Usuarios", Success: true, Message: "permiso creado exitosamente"},
		{Name: "Editar Usuarios", Success: true, Message: "permiso creado exitosamente"},
	}

	mockUseCase.On("BulkCreatePermissions", mock.Anything, mock.AnythingOfType("[]domain.CreatePermissionDTO")).
		Return(expectedResults, nil)

	requestBody := map[string]interface{}{
		"permissions": []map[string]interface{}{
			{
				"name":        "Crear Usuarios",
				"resource_id": 1,
				"action_id":   1,
				"scope_id":    1,
			},
			{
				"name":        "Editar Usuarios",
				"resource_id": 1,
				"action_id":   2,
				"scope_id":    1,
			},
		},
	}

	body, err := json.Marshal(requestBody)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/permissions/bulk", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, true, response["success"])
	assert.Contains(t, response["message"], "2 de 2")

	results, ok := response["results"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, results, 2)

	mockUseCase.AssertExpectations(t)
}

func TestBulkCreatePermissionsHandler_InvalidJSON(t *testing.T) {
	// Arrange
	mockUseCase := new(mocks.UseCaseMock)
	mockLogger := new(mocks.LoggerMock)

	handler := New(mockUseCase, mockLogger)
	router := setupTestRouter(handler)

	mockLogger.On("Info", mock.Anything).Maybe()
	mockLogger.On("Error", mock.Anything).Maybe()

	// Body malformado
	malformedBody := []byte(`{"permissions": [{"name": "sin cerrar}`)

	req := httptest.NewRequest(http.MethodPost, "/permissions/bulk", bytes.NewBuffer(malformedBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, false, response["success"])
	assert.Contains(t, response["message"], "inválidos")

	// El use case no debe ser llamado ante un JSON invalido
	mockUseCase.AssertNotCalled(t, "BulkCreatePermissions", mock.Anything, mock.Anything)
}

func TestBulkCreatePermissionsHandler_PartialSuccess(t *testing.T) {
	// Arrange
	mockUseCase := new(mocks.UseCaseMock)
	mockLogger := new(mocks.LoggerMock)

	handler := New(mockUseCase, mockLogger)
	router := setupTestRouter(handler)

	mockLogger.On("Info", mock.Anything).Maybe()
	mockLogger.On("Error", mock.Anything).Maybe()

	// Resultados mixtos: uno exitoso, uno fallido
	expectedResults := []domain.BulkCreateResult{
		{Name: "Crear Usuarios", Success: true, Message: "permiso creado exitosamente"},
		{Name: "Permiso Duplicado", Success: false, Error: "ya existe un permiso con el nombre 'Permiso Duplicado'"},
	}

	mockUseCase.On("BulkCreatePermissions", mock.Anything, mock.AnythingOfType("[]domain.CreatePermissionDTO")).
		Return(expectedResults, nil)

	requestBody := map[string]interface{}{
		"permissions": []map[string]interface{}{
			{
				"name":        "Crear Usuarios",
				"resource_id": 1,
				"action_id":   1,
				"scope_id":    1,
			},
			{
				"name":        "Permiso Duplicado",
				"resource_id": 1,
				"action_id":   2,
				"scope_id":    1,
			},
		},
	}

	body, err := json.Marshal(requestBody)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/permissions/bulk", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, true, response["success"])
	// Solo 1 de 2 fue exitoso
	assert.Contains(t, response["message"], "1 de 2")

	results, ok := response["results"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, results, 2)

	// Verificar resultado exitoso
	firstResult := results[0].(map[string]interface{})
	assert.Equal(t, true, firstResult["success"])
	assert.Equal(t, "Crear Usuarios", firstResult["name"])

	// Verificar resultado fallido
	secondResult := results[1].(map[string]interface{})
	assert.Equal(t, false, secondResult["success"])
	assert.Equal(t, "Permiso Duplicado", secondResult["name"])

	mockUseCase.AssertExpectations(t)
}
