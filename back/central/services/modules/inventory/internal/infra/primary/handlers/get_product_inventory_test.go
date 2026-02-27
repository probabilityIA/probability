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

// setupGetProductInventoryContext crea un contexto de Gin para el endpoint GetProductInventory
func setupGetProductInventoryContext(w *httptest.ResponseRecorder, productID string, businessID uint) *gin.Context {
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodGet, "/inventory/product/"+productID, nil)
	c.Request = req
	c.Params = gin.Params{{Key: "productId", Value: productID}}
	if businessID > 0 {
		c.Set("business_id", businessID)
	}
	return c
}

func TestGetProductInventory_SinBusinessID_RetornaBadRequest(t *testing.T) {
	// Arrange
	uc := &mocks.UseCaseMock{}
	h := New(uc)

	w := httptest.NewRecorder()
	c := setupGetProductInventoryContext(w, "prod-001", 0) // sin business_id

	// Act
	h.GetProductInventory(c)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("status esperado %d, se obtuvo %d", http.StatusBadRequest, w.Code)
	}

	var body map[string]string
	json.Unmarshal(w.Body.Bytes(), &body)
	if body["error"] == "" {
		t.Error("respuesta de error debe incluir campo 'error'")
	}
}

func TestGetProductInventory_ProductoNoEncontrado_RetornaNotFound(t *testing.T) {
	// Arrange
	uc := &mocks.UseCaseMock{
		GetProductInventoryFn: func(ctx context.Context, params dtos.GetProductInventoryParams) ([]entities.InventoryLevel, error) {
			return nil, domainerrors.ErrProductNotFound
		},
	}
	h := New(uc)

	w := httptest.NewRecorder()
	c := setupGetProductInventoryContext(w, "prod-inexistente", 10)

	// Act
	h.GetProductInventory(c)

	// Assert
	if w.Code != http.StatusNotFound {
		t.Errorf("status esperado %d, se obtuvo %d", http.StatusNotFound, w.Code)
	}
}

func TestGetProductInventory_ErrorInterno_RetornaInternalServerError(t *testing.T) {
	// Arrange
	uc := &mocks.UseCaseMock{
		GetProductInventoryFn: func(ctx context.Context, params dtos.GetProductInventoryParams) ([]entities.InventoryLevel, error) {
			return nil, errors.New("error de base de datos")
		},
	}
	h := New(uc)

	w := httptest.NewRecorder()
	c := setupGetProductInventoryContext(w, "prod-001", 10)

	// Act
	h.GetProductInventory(c)

	// Assert
	if w.Code != http.StatusInternalServerError {
		t.Errorf("status esperado %d, se obtuvo %d", http.StatusInternalServerError, w.Code)
	}
}

func TestGetProductInventory_Exitoso_RetornaListaInventario(t *testing.T) {
	// Arrange
	levelsEsperados := []entities.InventoryLevel{
		{
			ID:           1,
			ProductID:    "prod-001",
			WarehouseID:  1,
			BusinessID:   10,
			Quantity:     100,
			AvailableQty: 90,
			ReservedQty:  10,
			ProductName:  "Camisa Azul",
			ProductSKU:   "CAM-001",
		},
	}

	paramsCapturados := dtos.GetProductInventoryParams{}
	uc := &mocks.UseCaseMock{
		GetProductInventoryFn: func(ctx context.Context, params dtos.GetProductInventoryParams) ([]entities.InventoryLevel, error) {
			paramsCapturados = params
			return levelsEsperados, nil
		},
	}
	h := New(uc)

	w := httptest.NewRecorder()
	c := setupGetProductInventoryContext(w, "prod-001", 10)

	// Act
	h.GetProductInventory(c)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("status esperado %d, se obtuvo %d. Body: %s", http.StatusOK, w.Code, w.Body.String())
	}

	// Verificar params pasados al use case
	if paramsCapturados.ProductID != "prod-001" {
		t.Errorf("product_id esperado 'prod-001', se obtuvo '%s'", paramsCapturados.ProductID)
	}
	if paramsCapturados.BusinessID != 10 {
		t.Errorf("business_id esperado 10, se obtuvo %d", paramsCapturados.BusinessID)
	}

	// Verificar respuesta JSON es un array
	var responseBody []map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &responseBody); err != nil {
		t.Fatalf("error al parsear respuesta JSON: %v", err)
	}
	if len(responseBody) != 1 {
		t.Errorf("cantidad de items esperada 1, se obtuvo %d", len(responseBody))
	}
}

func TestGetProductInventory_InventarioVacio_RetornaArrayVacio(t *testing.T) {
	// Arrange
	uc := &mocks.UseCaseMock{
		GetProductInventoryFn: func(ctx context.Context, params dtos.GetProductInventoryParams) ([]entities.InventoryLevel, error) {
			return []entities.InventoryLevel{}, nil
		},
	}
	h := New(uc)

	w := httptest.NewRecorder()
	c := setupGetProductInventoryContext(w, "prod-001", 10)

	// Act
	h.GetProductInventory(c)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("status esperado %d, se obtuvo %d", http.StatusOK, w.Code)
	}
}

// -----------------------------------------------------------------------
// ListWarehouseInventory
// -----------------------------------------------------------------------

func setupListWarehouseContext(w *httptest.ResponseRecorder, warehouseID string, businessID uint, queryParams string) *gin.Context {
	c, _ := gin.CreateTestContext(w)
	url := "/inventory/warehouse/" + warehouseID
	if queryParams != "" {
		url += "?" + queryParams
	}
	req := httptest.NewRequest(http.MethodGet, url, nil)
	c.Request = req
	c.Params = gin.Params{{Key: "warehouseId", Value: warehouseID}}
	if businessID > 0 {
		c.Set("business_id", businessID)
	}
	return c
}

func TestListWarehouseInventory_SinBusinessID_RetornaBadRequest(t *testing.T) {
	// Arrange
	uc := &mocks.UseCaseMock{}
	h := New(uc)

	w := httptest.NewRecorder()
	c := setupListWarehouseContext(w, "1", 0, "")

	// Act
	h.ListWarehouseInventory(c)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("status esperado %d, se obtuvo %d", http.StatusBadRequest, w.Code)
	}
}

func TestListWarehouseInventory_WarehouseIDInvalido_RetornaBadRequest(t *testing.T) {
	// Arrange
	uc := &mocks.UseCaseMock{}
	h := New(uc)

	w := httptest.NewRecorder()
	c := setupListWarehouseContext(w, "no-es-numero", 10, "")

	// Act
	h.ListWarehouseInventory(c)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("status esperado %d, se obtuvo %d", http.StatusBadRequest, w.Code)
	}
}

func TestListWarehouseInventory_BodegaNoEncontrada_RetornaNotFound(t *testing.T) {
	// Arrange
	uc := &mocks.UseCaseMock{
		ListWarehouseInventoryFn: func(ctx context.Context, params dtos.ListWarehouseInventoryParams) ([]entities.InventoryLevel, int64, error) {
			return nil, 0, domainerrors.ErrWarehouseNotFound
		},
	}
	h := New(uc)

	w := httptest.NewRecorder()
	c := setupListWarehouseContext(w, "99", 10, "")

	// Act
	h.ListWarehouseInventory(c)

	// Assert
	if w.Code != http.StatusNotFound {
		t.Errorf("status esperado %d, se obtuvo %d", http.StatusNotFound, w.Code)
	}
}

func TestListWarehouseInventory_Exitoso_RetornaRespuestaPaginada(t *testing.T) {
	// Arrange
	uc := &mocks.UseCaseMock{
		ListWarehouseInventoryFn: func(ctx context.Context, params dtos.ListWarehouseInventoryParams) ([]entities.InventoryLevel, int64, error) {
			return []entities.InventoryLevel{
				{ID: 1, ProductID: "prod-001", WarehouseID: 1, Quantity: 50},
				{ID: 2, ProductID: "prod-002", WarehouseID: 1, Quantity: 30},
			}, 45, nil
		},
	}
	h := New(uc)

	w := httptest.NewRecorder()
	c := setupListWarehouseContext(w, "1", 10, "page=1&page_size=2")

	// Act
	h.ListWarehouseInventory(c)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("status esperado %d, se obtuvo %d. Body: %s", http.StatusOK, w.Code, w.Body.String())
	}

	var responseBody map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &responseBody); err != nil {
		t.Fatalf("error al parsear respuesta JSON: %v", err)
	}

	// Verificar estructura paginada
	if responseBody["data"] == nil {
		t.Error("respuesta debe incluir campo 'data'")
	}
	if responseBody["total"] == nil {
		t.Error("respuesta debe incluir campo 'total'")
	}
	if responseBody["page"] == nil {
		t.Error("respuesta debe incluir campo 'page'")
	}
	if responseBody["total_pages"] == nil {
		t.Error("respuesta debe incluir campo 'total_pages'")
	}
}

// -----------------------------------------------------------------------
// TransferStock
// -----------------------------------------------------------------------

func setupTransferStockContext(w *httptest.ResponseRecorder, body []byte, businessID uint) *gin.Context {
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodPost, "/inventory/transfer", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	if businessID > 0 {
		c.Set("business_id", businessID)
	}
	return c
}

func TestTransferStock_SinBusinessID_RetornaBadRequest(t *testing.T) {
	// Arrange
	uc := &mocks.UseCaseMock{}
	h := New(uc)

	w := httptest.NewRecorder()
	body := []byte(`{"product_id":"prod-001","from_warehouse_id":1,"to_warehouse_id":2,"quantity":5}`)
	c := setupTransferStockContext(w, body, 0)

	// Act
	h.TransferStock(c)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("status esperado %d, se obtuvo %d", http.StatusBadRequest, w.Code)
	}
}

func TestTransferStock_MismasBodegas_RetornaBadRequest(t *testing.T) {
	// Arrange
	uc := &mocks.UseCaseMock{
		TransferStockFn: func(ctx context.Context, dto dtos.TransferStockDTO) error {
			return domainerrors.ErrSameWarehouse
		},
	}
	h := New(uc)

	w := httptest.NewRecorder()
	body := []byte(`{"product_id":"prod-001","from_warehouse_id":1,"to_warehouse_id":1,"quantity":5}`)
	c := setupTransferStockContext(w, body, 10)

	// Act
	h.TransferStock(c)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("status esperado %d, se obtuvo %d", http.StatusBadRequest, w.Code)
	}
}

func TestTransferStock_Exitoso_RetornaOKConMensaje(t *testing.T) {
	// Arrange
	uc := &mocks.UseCaseMock{
		TransferStockFn: func(ctx context.Context, dto dtos.TransferStockDTO) error {
			return nil
		},
	}
	h := New(uc)

	w := httptest.NewRecorder()
	body := []byte(`{"product_id":"prod-001","from_warehouse_id":1,"to_warehouse_id":2,"quantity":5}`)
	c := setupTransferStockContext(w, body, 10)

	// Act
	h.TransferStock(c)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("status esperado %d, se obtuvo %d. Body: %s", http.StatusOK, w.Code, w.Body.String())
	}

	var responseBody map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &responseBody); err != nil {
		t.Fatalf("error al parsear respuesta JSON: %v", err)
	}
	if responseBody["message"] == nil {
		t.Error("respuesta exitosa debe incluir campo 'message'")
	}
}
