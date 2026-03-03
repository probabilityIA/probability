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

func init() {
	gin.SetMode(gin.TestMode)
}

// crearContextoGin crea un contexto de Gin con el business_id configurado como key
func crearContextoGin(w *httptest.ResponseRecorder, method, path string, body []byte) (*gin.Context, *gin.Engine) {
	router := gin.New()
	var capturedCtx *gin.Context

	router.Any(path, func(c *gin.Context) {
		capturedCtx = c
	})

	req := httptest.NewRequest(method, path, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	return capturedCtx, router
}

// setupTestContext crea un contexto de Gin de prueba con business_id
func setupTestContext(w *httptest.ResponseRecorder, body []byte, businessID uint) *gin.Context {
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodPost, "/inventory/adjust", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	if businessID > 0 {
		c.Set("business_id", businessID)
	}
	return c
}

func TestAdjustStock_SinBusinessID_RetornaBadRequest(t *testing.T) {
	// Arrange
	uc := &mocks.UseCaseMock{}
	h := New(uc)

	w := httptest.NewRecorder()
	body := []byte(`{"product_id":"prod-001","warehouse_id":1,"quantity":5,"reason":"ajuste"}`)
	c := setupTestContext(w, body, 0) // sin business_id

	// Act
	h.AdjustStock(c)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("status esperado %d, se obtuvo %d", http.StatusBadRequest, w.Code)
	}
}

func TestAdjustStock_BodyInvalido_RetornaBadRequest(t *testing.T) {
	// Arrange
	uc := &mocks.UseCaseMock{}
	h := New(uc)

	w := httptest.NewRecorder()
	// Falta el campo "reason" (required)
	body := []byte(`{"product_id":"prod-001","warehouse_id":1,"quantity":5}`)
	c := setupTestContext(w, body, 10)

	// Act
	h.AdjustStock(c)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("status esperado %d, se obtuvo %d", http.StatusBadRequest, w.Code)
	}
}

func TestAdjustStock_ProductoNoEncontrado_RetornaNotFound(t *testing.T) {
	// Arrange
	uc := &mocks.UseCaseMock{
		AdjustStockFn: func(ctx context.Context, dto dtos.AdjustStockDTO) (*entities.StockMovement, error) {
			return nil, domainerrors.ErrProductNotFound
		},
	}
	h := New(uc)

	w := httptest.NewRecorder()
	body := []byte(`{"product_id":"prod-inexistente","warehouse_id":1,"quantity":5,"reason":"ajuste test"}`)
	c := setupTestContext(w, body, 10)

	// Act
	h.AdjustStock(c)

	// Assert
	if w.Code != http.StatusNotFound {
		t.Errorf("status esperado %d, se obtuvo %d", http.StatusNotFound, w.Code)
	}
}

func TestAdjustStock_BodegaNoEncontrada_RetornaNotFound(t *testing.T) {
	// Arrange
	uc := &mocks.UseCaseMock{
		AdjustStockFn: func(ctx context.Context, dto dtos.AdjustStockDTO) (*entities.StockMovement, error) {
			return nil, domainerrors.ErrWarehouseNotFound
		},
	}
	h := New(uc)

	w := httptest.NewRecorder()
	body := []byte(`{"product_id":"prod-001","warehouse_id":99,"quantity":5,"reason":"ajuste test"}`)
	c := setupTestContext(w, body, 10)

	// Act
	h.AdjustStock(c)

	// Assert
	if w.Code != http.StatusNotFound {
		t.Errorf("status esperado %d, se obtuvo %d", http.StatusNotFound, w.Code)
	}
}

func TestAdjustStock_CantidadInvalida_RetornaBadRequest(t *testing.T) {
	// Arrange
	uc := &mocks.UseCaseMock{
		AdjustStockFn: func(ctx context.Context, dto dtos.AdjustStockDTO) (*entities.StockMovement, error) {
			return nil, domainerrors.ErrInvalidQuantity
		},
	}
	h := New(uc)

	w := httptest.NewRecorder()
	body := []byte(`{"product_id":"prod-001","warehouse_id":1,"quantity":0,"reason":"ajuste test"}`)
	c := setupTestContext(w, body, 10)

	// Act
	h.AdjustStock(c)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("status esperado %d, se obtuvo %d", http.StatusBadRequest, w.Code)
	}
}

func TestAdjustStock_ErrorInterno_RetornaInternalServerError(t *testing.T) {
	// Arrange
	uc := &mocks.UseCaseMock{
		AdjustStockFn: func(ctx context.Context, dto dtos.AdjustStockDTO) (*entities.StockMovement, error) {
			return nil, errors.New("error de base de datos")
		},
	}
	h := New(uc)

	w := httptest.NewRecorder()
	body := []byte(`{"product_id":"prod-001","warehouse_id":1,"quantity":5,"reason":"ajuste test"}`)
	c := setupTestContext(w, body, 10)

	// Act
	h.AdjustStock(c)

	// Assert
	if w.Code != http.StatusInternalServerError {
		t.Errorf("status esperado %d, se obtuvo %d", http.StatusInternalServerError, w.Code)
	}
}

func TestAdjustStock_Exitoso_RetornaCreatedConMovimiento(t *testing.T) {
	// Arrange
	movimientoEsperado := &entities.StockMovement{
		ID:          1,
		ProductID:   "prod-001",
		WarehouseID: 1,
		BusinessID:  10,
		Quantity:    5,
		NewQty:      55,
	}
	dtoCapturado := dtos.AdjustStockDTO{}
	uc := &mocks.UseCaseMock{
		AdjustStockFn: func(ctx context.Context, dto dtos.AdjustStockDTO) (*entities.StockMovement, error) {
			dtoCapturado = dto
			return movimientoEsperado, nil
		},
	}
	h := New(uc)

	w := httptest.NewRecorder()
	body := []byte(`{"product_id":"prod-001","warehouse_id":1,"quantity":5,"reason":"reposicion de stock","notes":"notas opcionales"}`)
	c := setupTestContext(w, body, 10)

	// Act
	h.AdjustStock(c)

	// Assert
	if w.Code != http.StatusCreated {
		t.Errorf("status esperado %d, se obtuvo %d. Body: %s", http.StatusCreated, w.Code, w.Body.String())
	}

	// Verificar que el DTO fue construido correctamente
	if dtoCapturado.BusinessID != 10 {
		t.Errorf("business_id esperado 10, se obtuvo %d", dtoCapturado.BusinessID)
	}
	if dtoCapturado.ProductID != "prod-001" {
		t.Errorf("product_id esperado 'prod-001', se obtuvo '%s'", dtoCapturado.ProductID)
	}
	if dtoCapturado.WarehouseID != 1 {
		t.Errorf("warehouse_id esperado 1, se obtuvo %d", dtoCapturado.WarehouseID)
	}

	// Verificar respuesta JSON
	var responseBody map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &responseBody); err != nil {
		t.Fatalf("error al parsear respuesta JSON: %v", err)
	}
	if responseBody["id"] == nil {
		t.Error("respuesta debe incluir el campo 'id'")
	}
}

func TestAdjustStock_ProductoSinTracking_RetornaBadRequest(t *testing.T) {
	// Arrange
	uc := &mocks.UseCaseMock{
		AdjustStockFn: func(ctx context.Context, dto dtos.AdjustStockDTO) (*entities.StockMovement, error) {
			return nil, domainerrors.ErrProductNoTracking
		},
	}
	h := New(uc)

	w := httptest.NewRecorder()
	body := []byte(`{"product_id":"prod-servicio","warehouse_id":1,"quantity":5,"reason":"ajuste test"}`)
	c := setupTestContext(w, body, 10)

	// Act
	h.AdjustStock(c)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("status esperado %d, se obtuvo %d", http.StatusBadRequest, w.Code)
	}
}
