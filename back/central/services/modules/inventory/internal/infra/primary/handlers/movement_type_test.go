package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/mocks"
)

// setupMovementTypeContext crea un contexto de Gin para endpoints de movement types
func setupMovementTypeContext(w *httptest.ResponseRecorder, method, path string, body []byte, businessID uint, params gin.Params) *gin.Context {
	c, _ := gin.CreateTestContext(w)
	var req *http.Request
	if body != nil {
		req = httptest.NewRequest(method, path, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	c.Request = req
	c.Params = params
	if businessID > 0 {
		c.Set("business_id", businessID)
	}
	return c
}

// -----------------------------------------------------------------------
// CreateMovementType
// -----------------------------------------------------------------------

func TestCreateMovementType_SinBusinessID_RetornaBadRequest(t *testing.T) {
	// Arrange
	uc := &mocks.UseCaseMock{}
	h := New(uc)

	w := httptest.NewRecorder()
	body := []byte(`{"code":"test","name":"Test","direction":"in"}`)
	c := setupMovementTypeContext(w, http.MethodPost, "/inventory/movement-types", body, 0, nil)

	// Act
	h.CreateMovementType(c)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("status esperado %d, se obtuvo %d", http.StatusBadRequest, w.Code)
	}
}

func TestCreateMovementType_BodyInvalido_RetornaBadRequest(t *testing.T) {
	// Arrange
	uc := &mocks.UseCaseMock{}
	h := New(uc)

	w := httptest.NewRecorder()
	// Falta "direction" (required con oneof=in out neutral)
	body := []byte(`{"code":"test","name":"Test"}`)
	c := setupMovementTypeContext(w, http.MethodPost, "/inventory/movement-types", body, 10, nil)

	// Act
	h.CreateMovementType(c)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("status esperado %d, se obtuvo %d", http.StatusBadRequest, w.Code)
	}
}

func TestCreateMovementType_CodigoYaExiste_RetornaConflict(t *testing.T) {
	// Arrange
	uc := &mocks.UseCaseMock{
		CreateMovementTypeFn: func(ctx context.Context, dto dtos.CreateStockMovementTypeDTO) (*entities.StockMovementType, error) {
			return nil, domainerrors.ErrMovementTypeCodeExists
		},
	}
	h := New(uc)

	w := httptest.NewRecorder()
	body := []byte(`{"code":"inbound","name":"Entrada","direction":"in"}`)
	c := setupMovementTypeContext(w, http.MethodPost, "/inventory/movement-types", body, 10, nil)

	// Act
	h.CreateMovementType(c)

	// Assert
	if w.Code != http.StatusConflict {
		t.Errorf("status esperado %d, se obtuvo %d", http.StatusConflict, w.Code)
	}
}

func TestCreateMovementType_ErrorInterno_RetornaInternalServerError(t *testing.T) {
	// Arrange
	uc := &mocks.UseCaseMock{
		CreateMovementTypeFn: func(ctx context.Context, dto dtos.CreateStockMovementTypeDTO) (*entities.StockMovementType, error) {
			return nil, errors.New("error de base de datos")
		},
	}
	h := New(uc)

	w := httptest.NewRecorder()
	body := []byte(`{"code":"nuevo","name":"Nuevo Tipo","direction":"neutral"}`)
	c := setupMovementTypeContext(w, http.MethodPost, "/inventory/movement-types", body, 10, nil)

	// Act
	h.CreateMovementType(c)

	// Assert
	if w.Code != http.StatusInternalServerError {
		t.Errorf("status esperado %d, se obtuvo %d", http.StatusInternalServerError, w.Code)
	}
}

func TestCreateMovementType_Exitoso_RetornaCreatedConTipo(t *testing.T) {
	// Arrange
	tipoEsperado := &entities.StockMovementType{
		ID:        10,
		Code:      "devolucion",
		Name:      "Devolución de cliente",
		Direction: "in",
		IsActive:  true,
	}
	dtoCapturado := dtos.CreateStockMovementTypeDTO{}
	uc := &mocks.UseCaseMock{
		CreateMovementTypeFn: func(ctx context.Context, dto dtos.CreateStockMovementTypeDTO) (*entities.StockMovementType, error) {
			dtoCapturado = dto
			return tipoEsperado, nil
		},
	}
	h := New(uc)

	w := httptest.NewRecorder()
	body := []byte(`{"code":"devolucion","name":"Devolución de cliente","description":"Devolucion por garantia","direction":"in"}`)
	c := setupMovementTypeContext(w, http.MethodPost, "/inventory/movement-types", body, 10, nil)

	// Act
	h.CreateMovementType(c)

	// Assert
	if w.Code != http.StatusCreated {
		t.Errorf("status esperado %d, se obtuvo %d. Body: %s", http.StatusCreated, w.Code, w.Body.String())
	}

	if dtoCapturado.Code != "devolucion" {
		t.Errorf("code esperado 'devolucion', se obtuvo '%s'", dtoCapturado.Code)
	}
	if dtoCapturado.Direction != "in" {
		t.Errorf("direction esperada 'in', se obtuvo '%s'", dtoCapturado.Direction)
	}

	var responseBody map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &responseBody); err != nil {
		t.Fatalf("error al parsear respuesta JSON: %v", err)
	}
	if responseBody["id"] == nil {
		t.Error("respuesta debe incluir campo 'id'")
	}
}

// -----------------------------------------------------------------------
// UpdateMovementType
// -----------------------------------------------------------------------

func TestUpdateMovementType_SinBusinessID_RetornaBadRequest(t *testing.T) {
	// Arrange
	uc := &mocks.UseCaseMock{}
	h := New(uc)

	w := httptest.NewRecorder()
	body := []byte(`{"name":"Nuevo nombre"}`)
	params := gin.Params{{Key: "id", Value: "5"}}
	c := setupMovementTypeContext(w, http.MethodPut, "/inventory/movement-types/5", body, 0, params)

	// Act
	h.UpdateMovementType(c)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("status esperado %d, se obtuvo %d", http.StatusBadRequest, w.Code)
	}
}

func TestUpdateMovementType_IDInvalido_RetornaBadRequest(t *testing.T) {
	// Arrange
	uc := &mocks.UseCaseMock{}
	h := New(uc)

	w := httptest.NewRecorder()
	body := []byte(`{"name":"Nuevo nombre"}`)
	params := gin.Params{{Key: "id", Value: "no-es-numero"}}
	c := setupMovementTypeContext(w, http.MethodPut, "/inventory/movement-types/no-es-numero", body, 10, params)

	// Act
	h.UpdateMovementType(c)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("status esperado %d, se obtuvo %d", http.StatusBadRequest, w.Code)
	}
}

func TestUpdateMovementType_TipoNoExiste_RetornaNotFound(t *testing.T) {
	// Arrange
	uc := &mocks.UseCaseMock{
		UpdateMovementTypeFn: func(ctx context.Context, dto dtos.UpdateStockMovementTypeDTO) (*entities.StockMovementType, error) {
			return nil, domainerrors.ErrMovementTypeNotFound
		},
	}
	h := New(uc)

	w := httptest.NewRecorder()
	body := []byte(`{"name":"Actualizado"}`)
	params := gin.Params{{Key: "id", Value: "999"}}
	c := setupMovementTypeContext(w, http.MethodPut, "/inventory/movement-types/999", body, 10, params)

	// Act
	h.UpdateMovementType(c)

	// Assert
	if w.Code != http.StatusNotFound {
		t.Errorf("status esperado %d, se obtuvo %d", http.StatusNotFound, w.Code)
	}
}

func TestUpdateMovementType_Exitoso_RetornaOKConTipoActualizado(t *testing.T) {
	// Arrange
	isActive := false
	tipoActualizado := &entities.StockMovementType{
		ID:       5,
		Code:     "ajuste",
		Name:     "Ajuste actualizado",
		IsActive: false,
	}
	dtoCapturado := dtos.UpdateStockMovementTypeDTO{}
	uc := &mocks.UseCaseMock{
		UpdateMovementTypeFn: func(ctx context.Context, dto dtos.UpdateStockMovementTypeDTO) (*entities.StockMovementType, error) {
			dtoCapturado = dto
			return tipoActualizado, nil
		},
	}
	h := New(uc)

	w := httptest.NewRecorder()
	body := []byte(`{"name":"Ajuste actualizado","is_active":false}`)
	params := gin.Params{{Key: "id", Value: "5"}}
	c := setupMovementTypeContext(w, http.MethodPut, "/inventory/movement-types/5", body, 10, params)

	// Act
	h.UpdateMovementType(c)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("status esperado %d, se obtuvo %d. Body: %s", http.StatusOK, w.Code, w.Body.String())
	}
	if dtoCapturado.ID != 5 {
		t.Errorf("ID en DTO esperado 5, se obtuvo %d", dtoCapturado.ID)
	}
	if dtoCapturado.IsActive == nil || *dtoCapturado.IsActive != isActive {
		t.Errorf("IsActive en DTO esperado false")
	}
}

// -----------------------------------------------------------------------
// DeleteMovementType
// -----------------------------------------------------------------------

func TestDeleteMovementType_SinBusinessID_RetornaBadRequest(t *testing.T) {
	// Arrange
	uc := &mocks.UseCaseMock{}
	h := New(uc)

	w := httptest.NewRecorder()
	params := gin.Params{{Key: "id", Value: "5"}}
	c := setupMovementTypeContext(w, http.MethodDelete, "/inventory/movement-types/5", nil, 0, params)

	// Act
	h.DeleteMovementType(c)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("status esperado %d, se obtuvo %d", http.StatusBadRequest, w.Code)
	}
}

func TestDeleteMovementType_IDInvalido_RetornaBadRequest(t *testing.T) {
	// Arrange
	uc := &mocks.UseCaseMock{}
	h := New(uc)

	w := httptest.NewRecorder()
	params := gin.Params{{Key: "id", Value: "abc"}}
	c := setupMovementTypeContext(w, http.MethodDelete, "/inventory/movement-types/abc", nil, 10, params)

	// Act
	h.DeleteMovementType(c)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("status esperado %d, se obtuvo %d", http.StatusBadRequest, w.Code)
	}
}

func TestDeleteMovementType_ErrorInterno_RetornaInternalServerError(t *testing.T) {
	// Arrange
	uc := &mocks.UseCaseMock{
		DeleteMovementTypeFn: func(ctx context.Context, id uint) error {
			return errors.New("no se puede eliminar: tiene movimientos asociados")
		},
	}
	h := New(uc)

	w := httptest.NewRecorder()
	params := gin.Params{{Key: "id", Value: "1"}}
	c := setupMovementTypeContext(w, http.MethodDelete, "/inventory/movement-types/1", nil, 10, params)

	// Act
	h.DeleteMovementType(c)

	// Assert
	if w.Code != http.StatusInternalServerError {
		t.Errorf("status esperado %d, se obtuvo %d", http.StatusInternalServerError, w.Code)
	}
}

func TestDeleteMovementType_Exitoso_RetornaOKConMensaje(t *testing.T) {
	// Arrange
	idCapturado := uint(0)
	uc := &mocks.UseCaseMock{
		DeleteMovementTypeFn: func(ctx context.Context, id uint) error {
			idCapturado = id
			return nil
		},
	}
	h := New(uc)

	w := httptest.NewRecorder()
	params := gin.Params{{Key: "id", Value: "7"}}
	c := setupMovementTypeContext(w, http.MethodDelete, "/inventory/movement-types/7", nil, 10, params)

	// Act
	h.DeleteMovementType(c)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("status esperado %d, se obtuvo %d. Body: %s", http.StatusOK, w.Code, w.Body.String())
	}
	if idCapturado != 7 {
		t.Errorf("ID capturado esperado 7, se obtuvo %d", idCapturado)
	}

	var responseBody map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &responseBody); err != nil {
		t.Fatalf("error al parsear respuesta JSON: %v", err)
	}
	if responseBody["message"] == nil {
		t.Error("respuesta exitosa debe incluir campo 'message'")
	}
}

// -----------------------------------------------------------------------
// ListMovementTypes
// -----------------------------------------------------------------------

func setupListMovementTypesContext(w *httptest.ResponseRecorder, businessID uint, queryParams string) *gin.Context {
	c, _ := gin.CreateTestContext(w)
	url := "/inventory/movement-types"
	if queryParams != "" {
		url += "?" + queryParams
	}
	req := httptest.NewRequest(http.MethodGet, url, nil)
	c.Request = req
	if businessID > 0 {
		c.Set("business_id", businessID)
	}
	return c
}

func TestListMovementTypes_Exitoso_RetornaRespuestaPaginada(t *testing.T) {
	// Arrange
	uc := &mocks.UseCaseMock{
		ListMovementTypesFn: func(ctx context.Context, params dtos.ListStockMovementTypesParams) ([]entities.StockMovementType, int64, error) {
			return []entities.StockMovementType{
				{ID: 1, Code: "inbound", Name: "Entrada", IsActive: true, Direction: "in"},
				{ID: 2, Code: "outbound", Name: "Salida", IsActive: true, Direction: "out"},
			}, 2, nil
		},
	}
	h := New(uc)

	w := httptest.NewRecorder()
	c := setupListMovementTypesContext(w, 10, "page=1&page_size=20")

	// Act
	h.ListMovementTypes(c)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("status esperado %d, se obtuvo %d. Body: %s", http.StatusOK, w.Code, w.Body.String())
	}

	var responseBody map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &responseBody); err != nil {
		t.Fatalf("error al parsear respuesta JSON: %v", err)
	}
	if responseBody["data"] == nil {
		t.Error("respuesta debe incluir campo 'data'")
	}
	if responseBody["total"] == nil {
		t.Error("respuesta debe incluir campo 'total'")
	}
}

func TestListMovementTypes_SinBusinessID_RetornaBadRequest(t *testing.T) {
	// Arrange
	uc := &mocks.UseCaseMock{}
	h := New(uc)

	w := httptest.NewRecorder()
	c := setupListMovementTypesContext(w, 0, "")

	// Act
	h.ListMovementTypes(c)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("status esperado %d, se obtuvo %d", http.StatusBadRequest, w.Code)
	}
}

// -----------------------------------------------------------------------
// ListMovements
// -----------------------------------------------------------------------

func setupListMovementsContext(w *httptest.ResponseRecorder, businessID uint, queryParams string) *gin.Context {
	c, _ := gin.CreateTestContext(w)
	url := "/inventory/movements"
	if queryParams != "" {
		url += "?" + queryParams
	}
	req := httptest.NewRequest(http.MethodGet, url, nil)
	c.Request = req
	if businessID > 0 {
		c.Set("business_id", businessID)
	}
	return c
}

func TestListMovements_SinBusinessID_RetornaBadRequest(t *testing.T) {
	// Arrange
	uc := &mocks.UseCaseMock{}
	h := New(uc)

	w := httptest.NewRecorder()
	c := setupListMovementsContext(w, 0, "")

	// Act
	h.ListMovements(c)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("status esperado %d, se obtuvo %d", http.StatusBadRequest, w.Code)
	}
}

func TestListMovements_Exitoso_RetornaRespuestaPaginada(t *testing.T) {
	// Arrange
	uc := &mocks.UseCaseMock{
		ListMovementsFn: func(ctx context.Context, params dtos.ListMovementsParams) ([]entities.StockMovement, int64, error) {
			return []entities.StockMovement{
				{ID: 1, ProductID: "prod-001", Quantity: 5},
			}, 1, nil
		},
	}
	h := New(uc)

	w := httptest.NewRecorder()
	c := setupListMovementsContext(w, 10, "product_id=prod-001&page=1&page_size=10")

	// Act
	h.ListMovements(c)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("status esperado %d, se obtuvo %d. Body: %s", http.StatusOK, w.Code, w.Body.String())
	}

	var responseBody map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &responseBody); err != nil {
		t.Fatalf("error al parsear respuesta JSON: %v", err)
	}
	if responseBody["data"] == nil {
		t.Error("respuesta debe incluir campo 'data'")
	}
	if responseBody["total_pages"] == nil {
		t.Error("respuesta debe incluir campo 'total_pages'")
	}
}
