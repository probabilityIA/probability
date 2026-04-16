package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/dtos"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/mocks"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// ---- helpers ----

// newEngine crea un router de gin limpio en modo test con el business_id inyectado en el contexto.
func newEngine(businessID uint) *gin.Engine {
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("business_id", businessID)
		c.Next()
	})
	return r
}

// newEngineNoBusinessID crea un router sin business_id en el contexto ni query param,
// para simular la ausencia de contexto de negocio.
func newEngineNoBusinessID() *gin.Engine {
	return gin.New()
}

func mustMarshal(t *testing.T, v any) *bytes.Buffer {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("error al serializar body: %v", err)
	}
	return bytes.NewBuffer(b)
}

func decodeJSON(t *testing.T, body *bytes.Buffer, dst any) {
	t.Helper()
	if err := json.NewDecoder(body).Decode(dst); err != nil {
		t.Fatalf("error al deserializar respuesta: %v", err)
	}
}

// warehouseFixture retorna una entidad Warehouse de prueba.
func warehouseFixture() *entities.Warehouse {
	return &entities.Warehouse{
		ID:            1,
		BusinessID:    10,
		Name:          "Bodega Central",
		Code:          "WH-001",
		Address:       "Calle 123",
		City:          "Bogotá",
		State:         "Cundinamarca",
		Country:       "Colombia",
		ZipCode:       "110111",
		Phone:         "3001234567",
		ContactName:   "Juan Pérez",
		ContactEmail:  "juan@test.com",
		IsActive:      true,
		IsDefault:     false,
		IsFulfillment: false,
		Locations:     []entities.WarehouseLocation{},
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

// locationFixture retorna una entidad WarehouseLocation de prueba.
func locationFixture() *entities.WarehouseLocation {
	cap := 100
	return &entities.WarehouseLocation{
		ID:            5,
		WarehouseID:   1,
		Name:          "Zona A",
		Code:          "LOC-A",
		Type:          "storage",
		IsActive:      true,
		IsFulfillment: false,
		Capacity:      &cap,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

// ========================================================================
// CREATE WAREHOUSE
// ========================================================================

func TestCreateWarehouse_Handler_Success(t *testing.T) {
	// Arrange
	expected := warehouseFixture()
	mock := &mocks.MockUseCase{
		CreateWarehouseFn: func(_ context.Context, dto dtos.CreateWarehouseDTO) (*entities.Warehouse, error) {
			return expected, nil
		},
	}

	r := newEngine(10)
	h := handlers.New(mock)
	r.POST("/warehouses", h.CreateWarehouse)

	body := mustMarshal(t, map[string]any{
		"name": "Bodega Central",
		"code": "WH-001",
		"city": "Bogotá",
	})
	req := httptest.NewRequest(http.MethodPost, "/warehouses", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusCreated {
		t.Errorf("esperado 201, obtenido %d — body: %s", w.Code, w.Body.String())
	}

	var resp map[string]any
	decodeJSON(t, w.Body, &resp)

	if resp["id"] == nil {
		t.Error("se esperaba un campo 'id' en la respuesta")
	}
	if resp["name"] != "Bodega Central" {
		t.Errorf("name incorrecto: %v", resp["name"])
	}
	if resp["code"] != "WH-001" {
		t.Errorf("code incorrecto: %v", resp["code"])
	}
}

func TestCreateWarehouse_Handler_MissingBusinessID(t *testing.T) {
	// Arrange: sin business_id en contexto ni query param
	mock := &mocks.MockUseCase{}
	r := newEngineNoBusinessID()
	h := handlers.New(mock)
	r.POST("/warehouses", h.CreateWarehouse)

	body := mustMarshal(t, map[string]any{"name": "Bodega X", "code": "WH-X"})
	req := httptest.NewRequest(http.MethodPost, "/warehouses", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado 400, obtenido %d", w.Code)
	}
}

func TestCreateWarehouse_Handler_InvalidBody(t *testing.T) {
	// Arrange: body sin campos obligatorios (name y code faltan)
	mock := &mocks.MockUseCase{}
	r := newEngine(10)
	h := handlers.New(mock)
	r.POST("/warehouses", h.CreateWarehouse)

	body := mustMarshal(t, map[string]any{"address": "Solo dirección"})
	req := httptest.NewRequest(http.MethodPost, "/warehouses", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado 400, obtenido %d — body: %s", w.Code, w.Body.String())
	}
}

func TestCreateWarehouse_Handler_DuplicateCode(t *testing.T) {
	// Arrange
	mock := &mocks.MockUseCase{
		CreateWarehouseFn: func(_ context.Context, dto dtos.CreateWarehouseDTO) (*entities.Warehouse, error) {
			return nil, domainerrors.ErrDuplicateCode
		},
	}
	r := newEngine(10)
	h := handlers.New(mock)
	r.POST("/warehouses", h.CreateWarehouse)

	body := mustMarshal(t, map[string]any{"name": "Bodega Dup", "code": "WH-001"})
	req := httptest.NewRequest(http.MethodPost, "/warehouses", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusConflict {
		t.Errorf("esperado 409, obtenido %d", w.Code)
	}
}

func TestCreateWarehouse_Handler_InternalError(t *testing.T) {
	// Arrange
	mock := &mocks.MockUseCase{
		CreateWarehouseFn: func(_ context.Context, dto dtos.CreateWarehouseDTO) (*entities.Warehouse, error) {
			return nil, errors.New("error de base de datos")
		},
	}
	r := newEngine(10)
	h := handlers.New(mock)
	r.POST("/warehouses", h.CreateWarehouse)

	body := mustMarshal(t, map[string]any{"name": "Bodega Error", "code": "WH-ERR"})
	req := httptest.NewRequest(http.MethodPost, "/warehouses", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusInternalServerError {
		t.Errorf("esperado 500, obtenido %d", w.Code)
	}
}

// ========================================================================
// LIST WAREHOUSES
// ========================================================================

func TestListWarehouses_Handler_Success(t *testing.T) {
	// Arrange
	warehouses := []entities.Warehouse{*warehouseFixture()}
	mock := &mocks.MockUseCase{
		ListWarehousesFn: func(_ context.Context, params dtos.ListWarehousesParams) ([]entities.Warehouse, int64, error) {
			return warehouses, 1, nil
		},
	}

	r := newEngine(10)
	h := handlers.New(mock)
	r.GET("/warehouses", h.ListWarehouses)

	req := httptest.NewRequest(http.MethodGet, "/warehouses?page=1&page_size=20", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("esperado 200, obtenido %d — body: %s", w.Code, w.Body.String())
	}

	var resp map[string]any
	decodeJSON(t, w.Body, &resp)

	if resp["data"] == nil {
		t.Error("se esperaba el campo 'data' en la respuesta")
	}
	if resp["total"] == nil {
		t.Error("se esperaba el campo 'total' en la respuesta")
	}
	if resp["total_pages"] == nil {
		t.Error("se esperaba el campo 'total_pages' en la respuesta")
	}
}

func TestListWarehouses_Handler_EmptyList(t *testing.T) {
	// Arrange: lista vacía
	mock := &mocks.MockUseCase{
		ListWarehousesFn: func(_ context.Context, params dtos.ListWarehousesParams) ([]entities.Warehouse, int64, error) {
			return []entities.Warehouse{}, 0, nil
		},
	}

	r := newEngine(10)
	h := handlers.New(mock)
	r.GET("/warehouses", h.ListWarehouses)

	req := httptest.NewRequest(http.MethodGet, "/warehouses", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("esperado 200, obtenido %d", w.Code)
	}

	var resp map[string]any
	decodeJSON(t, w.Body, &resp)

	total, _ := resp["total"].(float64)
	if total != 0 {
		t.Errorf("se esperaba total=0, obtenido %v", total)
	}
}

func TestListWarehouses_Handler_MissingBusinessID(t *testing.T) {
	// Arrange
	mock := &mocks.MockUseCase{}
	r := newEngineNoBusinessID()
	h := handlers.New(mock)
	r.GET("/warehouses", h.ListWarehouses)

	req := httptest.NewRequest(http.MethodGet, "/warehouses", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado 400, obtenido %d", w.Code)
	}
}

func TestListWarehouses_Handler_SuperAdminQueryParam(t *testing.T) {
	// Arrange: super admin pasa business_id como query param (business_id=0 en contexto)
	mock := &mocks.MockUseCase{
		ListWarehousesFn: func(_ context.Context, params dtos.ListWarehousesParams) ([]entities.Warehouse, int64, error) {
			return []entities.Warehouse{}, 0, nil
		},
	}

	r := newEngine(0) // super admin con business_id=0 en JWT
	h := handlers.New(mock)
	r.GET("/warehouses", h.ListWarehouses)

	req := httptest.NewRequest(http.MethodGet, "/warehouses?business_id=5", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("esperado 200, obtenido %d — body: %s", w.Code, w.Body.String())
	}
}

func TestListWarehouses_Handler_InternalError(t *testing.T) {
	// Arrange
	mock := &mocks.MockUseCase{
		ListWarehousesFn: func(_ context.Context, params dtos.ListWarehousesParams) ([]entities.Warehouse, int64, error) {
			return nil, 0, errors.New("fallo en repositorio")
		},
	}

	r := newEngine(10)
	h := handlers.New(mock)
	r.GET("/warehouses", h.ListWarehouses)

	req := httptest.NewRequest(http.MethodGet, "/warehouses", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusInternalServerError {
		t.Errorf("esperado 500, obtenido %d", w.Code)
	}
}

// ========================================================================
// GET WAREHOUSE
// ========================================================================

func TestGetWarehouse_Handler_Success(t *testing.T) {
	// Arrange
	expected := warehouseFixture()
	mock := &mocks.MockUseCase{
		GetWarehouseFn: func(_ context.Context, businessID, warehouseID uint) (*entities.Warehouse, error) {
			return expected, nil
		},
	}

	r := newEngine(10)
	h := handlers.New(mock)
	r.GET("/warehouses/:id", h.GetWarehouse)

	req := httptest.NewRequest(http.MethodGet, "/warehouses/1", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("esperado 200, obtenido %d — body: %s", w.Code, w.Body.String())
	}

	var resp map[string]any
	decodeJSON(t, w.Body, &resp)

	if resp["id"] == nil {
		t.Error("se esperaba el campo 'id' en la respuesta")
	}
	if resp["locations"] == nil {
		t.Error("se esperaba el campo 'locations' (detail) en la respuesta")
	}
}

func TestGetWarehouse_Handler_NotFound(t *testing.T) {
	// Arrange
	mock := &mocks.MockUseCase{
		GetWarehouseFn: func(_ context.Context, businessID, warehouseID uint) (*entities.Warehouse, error) {
			return nil, domainerrors.ErrWarehouseNotFound
		},
	}

	r := newEngine(10)
	h := handlers.New(mock)
	r.GET("/warehouses/:id", h.GetWarehouse)

	req := httptest.NewRequest(http.MethodGet, "/warehouses/999", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusNotFound {
		t.Errorf("esperado 404, obtenido %d", w.Code)
	}
}

func TestGetWarehouse_Handler_InvalidID(t *testing.T) {
	// Arrange: ID no numérico
	mock := &mocks.MockUseCase{}
	r := newEngine(10)
	h := handlers.New(mock)
	r.GET("/warehouses/:id", h.GetWarehouse)

	req := httptest.NewRequest(http.MethodGet, "/warehouses/abc", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado 400, obtenido %d", w.Code)
	}
}

func TestGetWarehouse_Handler_MissingBusinessID(t *testing.T) {
	// Arrange
	mock := &mocks.MockUseCase{}
	r := newEngineNoBusinessID()
	h := handlers.New(mock)
	r.GET("/warehouses/:id", h.GetWarehouse)

	req := httptest.NewRequest(http.MethodGet, "/warehouses/1", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado 400, obtenido %d", w.Code)
	}
}

func TestGetWarehouse_Handler_InternalError(t *testing.T) {
	// Arrange
	mock := &mocks.MockUseCase{
		GetWarehouseFn: func(_ context.Context, businessID, warehouseID uint) (*entities.Warehouse, error) {
			return nil, errors.New("error de conexión")
		},
	}

	r := newEngine(10)
	h := handlers.New(mock)
	r.GET("/warehouses/:id", h.GetWarehouse)

	req := httptest.NewRequest(http.MethodGet, "/warehouses/1", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusInternalServerError {
		t.Errorf("esperado 500, obtenido %d", w.Code)
	}
}

// ========================================================================
// UPDATE WAREHOUSE
// ========================================================================

func TestUpdateWarehouse_Handler_Success(t *testing.T) {
	// Arrange
	updated := warehouseFixture()
	updated.Name = "Bodega Actualizada"

	mock := &mocks.MockUseCase{
		UpdateWarehouseFn: func(_ context.Context, dto dtos.UpdateWarehouseDTO) (*entities.Warehouse, error) {
			return updated, nil
		},
	}

	r := newEngine(10)
	h := handlers.New(mock)
	r.PUT("/warehouses/:id", h.UpdateWarehouse)

	body := mustMarshal(t, map[string]any{
		"name": "Bodega Actualizada",
		"code": "WH-001",
	})
	req := httptest.NewRequest(http.MethodPut, "/warehouses/1", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("esperado 200, obtenido %d — body: %s", w.Code, w.Body.String())
	}

	var resp map[string]any
	decodeJSON(t, w.Body, &resp)
	if resp["name"] != "Bodega Actualizada" {
		t.Errorf("nombre incorrecto en respuesta: %v", resp["name"])
	}
}

func TestUpdateWarehouse_Handler_InvalidID(t *testing.T) {
	// Arrange
	mock := &mocks.MockUseCase{}
	r := newEngine(10)
	h := handlers.New(mock)
	r.PUT("/warehouses/:id", h.UpdateWarehouse)

	body := mustMarshal(t, map[string]any{"name": "X", "code": "Y"})
	req := httptest.NewRequest(http.MethodPut, "/warehouses/abc", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado 400, obtenido %d", w.Code)
	}
}

func TestUpdateWarehouse_Handler_InvalidBody(t *testing.T) {
	// Arrange: falta "name" y "code" requeridos
	mock := &mocks.MockUseCase{}
	r := newEngine(10)
	h := handlers.New(mock)
	r.PUT("/warehouses/:id", h.UpdateWarehouse)

	body := mustMarshal(t, map[string]any{"city": "Cali"})
	req := httptest.NewRequest(http.MethodPut, "/warehouses/1", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado 400, obtenido %d — body: %s", w.Code, w.Body.String())
	}
}

func TestUpdateWarehouse_Handler_NotFound(t *testing.T) {
	// Arrange
	mock := &mocks.MockUseCase{
		UpdateWarehouseFn: func(_ context.Context, dto dtos.UpdateWarehouseDTO) (*entities.Warehouse, error) {
			return nil, domainerrors.ErrWarehouseNotFound
		},
	}

	r := newEngine(10)
	h := handlers.New(mock)
	r.PUT("/warehouses/:id", h.UpdateWarehouse)

	body := mustMarshal(t, map[string]any{"name": "No existe", "code": "WH-X"})
	req := httptest.NewRequest(http.MethodPut, "/warehouses/999", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusNotFound {
		t.Errorf("esperado 404, obtenido %d", w.Code)
	}
}

func TestUpdateWarehouse_Handler_DuplicateCode(t *testing.T) {
	// Arrange
	mock := &mocks.MockUseCase{
		UpdateWarehouseFn: func(_ context.Context, dto dtos.UpdateWarehouseDTO) (*entities.Warehouse, error) {
			return nil, domainerrors.ErrDuplicateCode
		},
	}

	r := newEngine(10)
	h := handlers.New(mock)
	r.PUT("/warehouses/:id", h.UpdateWarehouse)

	body := mustMarshal(t, map[string]any{"name": "Bodega Dup", "code": "WH-DUP"})
	req := httptest.NewRequest(http.MethodPut, "/warehouses/1", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusConflict {
		t.Errorf("esperado 409, obtenido %d", w.Code)
	}
}

func TestUpdateWarehouse_Handler_InternalError(t *testing.T) {
	// Arrange
	mock := &mocks.MockUseCase{
		UpdateWarehouseFn: func(_ context.Context, dto dtos.UpdateWarehouseDTO) (*entities.Warehouse, error) {
			return nil, errors.New("timeout de BD")
		},
	}

	r := newEngine(10)
	h := handlers.New(mock)
	r.PUT("/warehouses/:id", h.UpdateWarehouse)

	body := mustMarshal(t, map[string]any{"name": "Bodega Error", "code": "WH-ERR"})
	req := httptest.NewRequest(http.MethodPut, "/warehouses/1", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusInternalServerError {
		t.Errorf("esperado 500, obtenido %d", w.Code)
	}
}

// ========================================================================
// DELETE WAREHOUSE
// ========================================================================

func TestDeleteWarehouse_Handler_Success(t *testing.T) {
	// Arrange
	mock := &mocks.MockUseCase{
		DeleteWarehouseFn: func(_ context.Context, businessID, warehouseID uint) error {
			return nil
		},
	}

	r := newEngine(10)
	h := handlers.New(mock)
	r.DELETE("/warehouses/:id", h.DeleteWarehouse)

	req := httptest.NewRequest(http.MethodDelete, "/warehouses/1", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("esperado 200, obtenido %d — body: %s", w.Code, w.Body.String())
	}

	var resp map[string]any
	decodeJSON(t, w.Body, &resp)
	if resp["message"] == nil {
		t.Error("se esperaba campo 'message' en respuesta de eliminación")
	}
}

func TestDeleteWarehouse_Handler_NotFound(t *testing.T) {
	// Arrange
	mock := &mocks.MockUseCase{
		DeleteWarehouseFn: func(_ context.Context, businessID, warehouseID uint) error {
			return domainerrors.ErrWarehouseNotFound
		},
	}

	r := newEngine(10)
	h := handlers.New(mock)
	r.DELETE("/warehouses/:id", h.DeleteWarehouse)

	req := httptest.NewRequest(http.MethodDelete, "/warehouses/999", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusNotFound {
		t.Errorf("esperado 404, obtenido %d", w.Code)
	}
}

func TestDeleteWarehouse_Handler_InvalidID(t *testing.T) {
	// Arrange
	mock := &mocks.MockUseCase{}
	r := newEngine(10)
	h := handlers.New(mock)
	r.DELETE("/warehouses/:id", h.DeleteWarehouse)

	req := httptest.NewRequest(http.MethodDelete, "/warehouses/abc", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado 400, obtenido %d", w.Code)
	}
}

func TestDeleteWarehouse_Handler_MissingBusinessID(t *testing.T) {
	// Arrange
	mock := &mocks.MockUseCase{}
	r := newEngineNoBusinessID()
	h := handlers.New(mock)
	r.DELETE("/warehouses/:id", h.DeleteWarehouse)

	req := httptest.NewRequest(http.MethodDelete, "/warehouses/1", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado 400, obtenido %d", w.Code)
	}
}

func TestDeleteWarehouse_Handler_InternalError(t *testing.T) {
	// Arrange
	mock := &mocks.MockUseCase{
		DeleteWarehouseFn: func(_ context.Context, businessID, warehouseID uint) error {
			return errors.New("error inesperado")
		},
	}

	r := newEngine(10)
	h := handlers.New(mock)
	r.DELETE("/warehouses/:id", h.DeleteWarehouse)

	req := httptest.NewRequest(http.MethodDelete, "/warehouses/1", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusInternalServerError {
		t.Errorf("esperado 500, obtenido %d", w.Code)
	}
}

// ========================================================================
// CREATE LOCATION
// ========================================================================

func TestCreateLocation_Handler_Success(t *testing.T) {
	// Arrange
	expected := locationFixture()
	mock := &mocks.MockUseCase{
		CreateLocationFn: func(_ context.Context, dto dtos.CreateLocationDTO) (*entities.WarehouseLocation, error) {
			return expected, nil
		},
	}

	r := newEngine(10)
	h := handlers.New(mock)
	r.POST("/warehouses/:id/locations", h.CreateLocation)

	body := mustMarshal(t, map[string]any{
		"name": "Zona A",
		"code": "LOC-A",
		"type": "storage",
	})
	req := httptest.NewRequest(http.MethodPost, "/warehouses/1/locations", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusCreated {
		t.Errorf("esperado 201, obtenido %d — body: %s", w.Code, w.Body.String())
	}

	var resp map[string]any
	decodeJSON(t, w.Body, &resp)
	if resp["id"] == nil {
		t.Error("se esperaba campo 'id' en la respuesta de ubicación")
	}
	if resp["code"] != "LOC-A" {
		t.Errorf("code incorrecto: %v", resp["code"])
	}
}

func TestCreateLocation_Handler_MissingBusinessID(t *testing.T) {
	// Arrange
	mock := &mocks.MockUseCase{}
	r := newEngineNoBusinessID()
	h := handlers.New(mock)
	r.POST("/warehouses/:id/locations", h.CreateLocation)

	body := mustMarshal(t, map[string]any{"name": "Zona X", "code": "LOC-X"})
	req := httptest.NewRequest(http.MethodPost, "/warehouses/1/locations", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado 400, obtenido %d", w.Code)
	}
}

func TestCreateLocation_Handler_InvalidWarehouseID(t *testing.T) {
	// Arrange
	mock := &mocks.MockUseCase{}
	r := newEngine(10)
	h := handlers.New(mock)
	r.POST("/warehouses/:id/locations", h.CreateLocation)

	body := mustMarshal(t, map[string]any{"name": "Zona X", "code": "LOC-X"})
	req := httptest.NewRequest(http.MethodPost, "/warehouses/abc/locations", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado 400, obtenido %d", w.Code)
	}
}

func TestCreateLocation_Handler_InvalidBody(t *testing.T) {
	// Arrange: falta name y code requeridos
	mock := &mocks.MockUseCase{}
	r := newEngine(10)
	h := handlers.New(mock)
	r.POST("/warehouses/:id/locations", h.CreateLocation)

	body := mustMarshal(t, map[string]any{"type": "storage"})
	req := httptest.NewRequest(http.MethodPost, "/warehouses/1/locations", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado 400, obtenido %d — body: %s", w.Code, w.Body.String())
	}
}

func TestCreateLocation_Handler_WarehouseNotFound(t *testing.T) {
	// Arrange
	mock := &mocks.MockUseCase{
		CreateLocationFn: func(_ context.Context, dto dtos.CreateLocationDTO) (*entities.WarehouseLocation, error) {
			return nil, domainerrors.ErrWarehouseNotFound
		},
	}

	r := newEngine(10)
	h := handlers.New(mock)
	r.POST("/warehouses/:id/locations", h.CreateLocation)

	body := mustMarshal(t, map[string]any{"name": "Zona X", "code": "LOC-X"})
	req := httptest.NewRequest(http.MethodPost, "/warehouses/999/locations", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusNotFound {
		t.Errorf("esperado 404, obtenido %d", w.Code)
	}
}

func TestCreateLocation_Handler_DuplicateLocCode(t *testing.T) {
	// Arrange
	mock := &mocks.MockUseCase{
		CreateLocationFn: func(_ context.Context, dto dtos.CreateLocationDTO) (*entities.WarehouseLocation, error) {
			return nil, domainerrors.ErrDuplicateLocCode
		},
	}

	r := newEngine(10)
	h := handlers.New(mock)
	r.POST("/warehouses/:id/locations", h.CreateLocation)

	body := mustMarshal(t, map[string]any{"name": "Zona Dup", "code": "LOC-DUP"})
	req := httptest.NewRequest(http.MethodPost, "/warehouses/1/locations", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusConflict {
		t.Errorf("esperado 409, obtenido %d", w.Code)
	}
}

func TestCreateLocation_Handler_InternalError(t *testing.T) {
	// Arrange
	mock := &mocks.MockUseCase{
		CreateLocationFn: func(_ context.Context, dto dtos.CreateLocationDTO) (*entities.WarehouseLocation, error) {
			return nil, errors.New("error interno")
		},
	}

	r := newEngine(10)
	h := handlers.New(mock)
	r.POST("/warehouses/:id/locations", h.CreateLocation)

	body := mustMarshal(t, map[string]any{"name": "Zona Err", "code": "LOC-ERR"})
	req := httptest.NewRequest(http.MethodPost, "/warehouses/1/locations", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusInternalServerError {
		t.Errorf("esperado 500, obtenido %d", w.Code)
	}
}

// ========================================================================
// LIST LOCATIONS
// ========================================================================

func TestListLocations_Handler_Success(t *testing.T) {
	// Arrange
	locations := []entities.WarehouseLocation{*locationFixture()}
	mock := &mocks.MockUseCase{
		ListLocationsFn: func(_ context.Context, params dtos.ListLocationsParams) ([]entities.WarehouseLocation, error) {
			return locations, nil
		},
	}

	r := newEngine(10)
	h := handlers.New(mock)
	r.GET("/warehouses/:id/locations", h.ListLocations)

	req := httptest.NewRequest(http.MethodGet, "/warehouses/1/locations", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("esperado 200, obtenido %d — body: %s", w.Code, w.Body.String())
	}

	var resp []map[string]any
	decodeJSON(t, w.Body, &resp)
	if len(resp) != 1 {
		t.Errorf("se esperaba 1 ubicación, obtenidas %d", len(resp))
	}
}

func TestListLocations_Handler_MissingBusinessID(t *testing.T) {
	// Arrange
	mock := &mocks.MockUseCase{}
	r := newEngineNoBusinessID()
	h := handlers.New(mock)
	r.GET("/warehouses/:id/locations", h.ListLocations)

	req := httptest.NewRequest(http.MethodGet, "/warehouses/1/locations", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado 400, obtenido %d", w.Code)
	}
}

func TestListLocations_Handler_InvalidWarehouseID(t *testing.T) {
	// Arrange
	mock := &mocks.MockUseCase{}
	r := newEngine(10)
	h := handlers.New(mock)
	r.GET("/warehouses/:id/locations", h.ListLocations)

	req := httptest.NewRequest(http.MethodGet, "/warehouses/abc/locations", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado 400, obtenido %d", w.Code)
	}
}

func TestListLocations_Handler_WarehouseNotFound(t *testing.T) {
	// Arrange
	mock := &mocks.MockUseCase{
		ListLocationsFn: func(_ context.Context, params dtos.ListLocationsParams) ([]entities.WarehouseLocation, error) {
			return nil, domainerrors.ErrWarehouseNotFound
		},
	}

	r := newEngine(10)
	h := handlers.New(mock)
	r.GET("/warehouses/:id/locations", h.ListLocations)

	req := httptest.NewRequest(http.MethodGet, "/warehouses/999/locations", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusNotFound {
		t.Errorf("esperado 404, obtenido %d", w.Code)
	}
}

func TestListLocations_Handler_InternalError(t *testing.T) {
	// Arrange
	mock := &mocks.MockUseCase{
		ListLocationsFn: func(_ context.Context, params dtos.ListLocationsParams) ([]entities.WarehouseLocation, error) {
			return nil, errors.New("error de repositorio")
		},
	}

	r := newEngine(10)
	h := handlers.New(mock)
	r.GET("/warehouses/:id/locations", h.ListLocations)

	req := httptest.NewRequest(http.MethodGet, "/warehouses/1/locations", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusInternalServerError {
		t.Errorf("esperado 500, obtenido %d", w.Code)
	}
}

// ========================================================================
// UPDATE LOCATION
// ========================================================================

func TestUpdateLocation_Handler_Success(t *testing.T) {
	// Arrange
	updated := locationFixture()
	updated.Name = "Zona Actualizada"

	mock := &mocks.MockUseCase{
		UpdateLocationFn: func(_ context.Context, dto dtos.UpdateLocationDTO) (*entities.WarehouseLocation, error) {
			return updated, nil
		},
	}

	r := newEngine(10)
	h := handlers.New(mock)
	r.PUT("/warehouses/:id/locations/:locationId", h.UpdateLocation)

	body := mustMarshal(t, map[string]any{
		"name": "Zona Actualizada",
		"code": "LOC-A",
	})
	req := httptest.NewRequest(http.MethodPut, "/warehouses/1/locations/5", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("esperado 200, obtenido %d — body: %s", w.Code, w.Body.String())
	}

	var resp map[string]any
	decodeJSON(t, w.Body, &resp)
	if resp["name"] != "Zona Actualizada" {
		t.Errorf("nombre incorrecto: %v", resp["name"])
	}
}

func TestUpdateLocation_Handler_InvalidWarehouseID(t *testing.T) {
	// Arrange
	mock := &mocks.MockUseCase{}
	r := newEngine(10)
	h := handlers.New(mock)
	r.PUT("/warehouses/:id/locations/:locationId", h.UpdateLocation)

	body := mustMarshal(t, map[string]any{"name": "X", "code": "Y"})
	req := httptest.NewRequest(http.MethodPut, "/warehouses/abc/locations/5", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado 400, obtenido %d", w.Code)
	}
}

func TestUpdateLocation_Handler_InvalidLocationID(t *testing.T) {
	// Arrange
	mock := &mocks.MockUseCase{}
	r := newEngine(10)
	h := handlers.New(mock)
	r.PUT("/warehouses/:id/locations/:locationId", h.UpdateLocation)

	body := mustMarshal(t, map[string]any{"name": "X", "code": "Y"})
	req := httptest.NewRequest(http.MethodPut, "/warehouses/1/locations/abc", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado 400, obtenido %d", w.Code)
	}
}

func TestUpdateLocation_Handler_InvalidBody(t *testing.T) {
	// Arrange: falta name y code
	mock := &mocks.MockUseCase{}
	r := newEngine(10)
	h := handlers.New(mock)
	r.PUT("/warehouses/:id/locations/:locationId", h.UpdateLocation)

	body := mustMarshal(t, map[string]any{"capacity": 50})
	req := httptest.NewRequest(http.MethodPut, "/warehouses/1/locations/5", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado 400, obtenido %d — body: %s", w.Code, w.Body.String())
	}
}

func TestUpdateLocation_Handler_NotFound(t *testing.T) {
	// Arrange
	mock := &mocks.MockUseCase{
		UpdateLocationFn: func(_ context.Context, dto dtos.UpdateLocationDTO) (*entities.WarehouseLocation, error) {
			return nil, domainerrors.ErrLocationNotFound
		},
	}

	r := newEngine(10)
	h := handlers.New(mock)
	r.PUT("/warehouses/:id/locations/:locationId", h.UpdateLocation)

	body := mustMarshal(t, map[string]any{"name": "Zona X", "code": "LOC-X"})
	req := httptest.NewRequest(http.MethodPut, "/warehouses/1/locations/999", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusNotFound {
		t.Errorf("esperado 404, obtenido %d", w.Code)
	}
}

func TestUpdateLocation_Handler_DuplicateCode(t *testing.T) {
	// Arrange
	mock := &mocks.MockUseCase{
		UpdateLocationFn: func(_ context.Context, dto dtos.UpdateLocationDTO) (*entities.WarehouseLocation, error) {
			return nil, domainerrors.ErrDuplicateLocCode
		},
	}

	r := newEngine(10)
	h := handlers.New(mock)
	r.PUT("/warehouses/:id/locations/:locationId", h.UpdateLocation)

	body := mustMarshal(t, map[string]any{"name": "Zona Dup", "code": "LOC-DUP"})
	req := httptest.NewRequest(http.MethodPut, "/warehouses/1/locations/5", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusConflict {
		t.Errorf("esperado 409, obtenido %d", w.Code)
	}
}

func TestUpdateLocation_Handler_InternalError(t *testing.T) {
	// Arrange
	mock := &mocks.MockUseCase{
		UpdateLocationFn: func(_ context.Context, dto dtos.UpdateLocationDTO) (*entities.WarehouseLocation, error) {
			return nil, errors.New("error de BD")
		},
	}

	r := newEngine(10)
	h := handlers.New(mock)
	r.PUT("/warehouses/:id/locations/:locationId", h.UpdateLocation)

	body := mustMarshal(t, map[string]any{"name": "Zona Err", "code": "LOC-ERR"})
	req := httptest.NewRequest(http.MethodPut, "/warehouses/1/locations/5", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusInternalServerError {
		t.Errorf("esperado 500, obtenido %d", w.Code)
	}
}

// ========================================================================
// DELETE LOCATION
// ========================================================================

func TestDeleteLocation_Handler_Success(t *testing.T) {
	// Arrange
	mock := &mocks.MockUseCase{
		DeleteLocationFn: func(_ context.Context, warehouseID, locationID uint, businessID uint) error {
			return nil
		},
	}

	r := newEngine(10)
	h := handlers.New(mock)
	r.DELETE("/warehouses/:id/locations/:locationId", h.DeleteLocation)

	req := httptest.NewRequest(http.MethodDelete, "/warehouses/1/locations/5", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("esperado 200, obtenido %d — body: %s", w.Code, w.Body.String())
	}

	var resp map[string]any
	decodeJSON(t, w.Body, &resp)
	if resp["message"] == nil {
		t.Error("se esperaba campo 'message' en respuesta de eliminación de ubicación")
	}
}

func TestDeleteLocation_Handler_MissingBusinessID(t *testing.T) {
	// Arrange
	mock := &mocks.MockUseCase{}
	r := newEngineNoBusinessID()
	h := handlers.New(mock)
	r.DELETE("/warehouses/:id/locations/:locationId", h.DeleteLocation)

	req := httptest.NewRequest(http.MethodDelete, "/warehouses/1/locations/5", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado 400, obtenido %d", w.Code)
	}
}

func TestDeleteLocation_Handler_InvalidWarehouseID(t *testing.T) {
	// Arrange
	mock := &mocks.MockUseCase{}
	r := newEngine(10)
	h := handlers.New(mock)
	r.DELETE("/warehouses/:id/locations/:locationId", h.DeleteLocation)

	req := httptest.NewRequest(http.MethodDelete, "/warehouses/abc/locations/5", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado 400, obtenido %d", w.Code)
	}
}

func TestDeleteLocation_Handler_InvalidLocationID(t *testing.T) {
	// Arrange
	mock := &mocks.MockUseCase{}
	r := newEngine(10)
	h := handlers.New(mock)
	r.DELETE("/warehouses/:id/locations/:locationId", h.DeleteLocation)

	req := httptest.NewRequest(http.MethodDelete, "/warehouses/1/locations/abc", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado 400, obtenido %d", w.Code)
	}
}

func TestDeleteLocation_Handler_NotFound(t *testing.T) {
	// Arrange
	mock := &mocks.MockUseCase{
		DeleteLocationFn: func(_ context.Context, warehouseID, locationID uint, businessID uint) error {
			return domainerrors.ErrLocationNotFound
		},
	}

	r := newEngine(10)
	h := handlers.New(mock)
	r.DELETE("/warehouses/:id/locations/:locationId", h.DeleteLocation)

	req := httptest.NewRequest(http.MethodDelete, "/warehouses/1/locations/999", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusNotFound {
		t.Errorf("esperado 404, obtenido %d", w.Code)
	}
}

func TestDeleteLocation_Handler_InternalError(t *testing.T) {
	// Arrange
	mock := &mocks.MockUseCase{
		DeleteLocationFn: func(_ context.Context, warehouseID, locationID uint, businessID uint) error {
			return errors.New("error inesperado en repositorio")
		},
	}

	r := newEngine(10)
	h := handlers.New(mock)
	r.DELETE("/warehouses/:id/locations/:locationId", h.DeleteLocation)

	req := httptest.NewRequest(http.MethodDelete, "/warehouses/1/locations/5", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusInternalServerError {
		t.Errorf("esperado 500, obtenido %d", w.Code)
	}
}
