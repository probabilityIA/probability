package handlerintegrationtype

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ============================================
// Mock del UseCase (usecaseintegrationtype.IIntegrationTypeUseCase)
// ============================================

type mockIntegrationTypeUseCase struct {
	mock.Mock
}

func (m *mockIntegrationTypeUseCase) CreateIntegrationType(ctx context.Context, dto domain.CreateIntegrationTypeDTO) (*domain.IntegrationType, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.IntegrationType), args.Error(1)
}

func (m *mockIntegrationTypeUseCase) UpdateIntegrationType(ctx context.Context, id uint, dto domain.UpdateIntegrationTypeDTO) (*domain.IntegrationType, error) {
	args := m.Called(ctx, id, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.IntegrationType), args.Error(1)
}

func (m *mockIntegrationTypeUseCase) GetIntegrationTypeByID(ctx context.Context, id uint) (*domain.IntegrationType, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.IntegrationType), args.Error(1)
}

func (m *mockIntegrationTypeUseCase) GetIntegrationTypeByCode(ctx context.Context, code string) (*domain.IntegrationType, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.IntegrationType), args.Error(1)
}

func (m *mockIntegrationTypeUseCase) DeleteIntegrationType(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockIntegrationTypeUseCase) ListIntegrationTypes(ctx context.Context, categoryID *uint) ([]*domain.IntegrationType, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.IntegrationType), args.Error(1)
}

func (m *mockIntegrationTypeUseCase) GetPlatformCredentials(ctx context.Context, id uint) (map[string]interface{}, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *mockIntegrationTypeUseCase) ListActiveIntegrationTypes(ctx context.Context) ([]*domain.IntegrationType, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.IntegrationType), args.Error(1)
}

func (m *mockIntegrationTypeUseCase) ListIntegrationCategories(ctx context.Context) ([]*domain.IntegrationCategory, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.IntegrationCategory), args.Error(1)
}

// ============================================
// Mock de ILogger para handlers de tipos
// ============================================

type mockTypeHandlerLogger struct {
	l zerolog.Logger
}

func newMockTypeHandlerLogger() *mockTypeHandlerLogger {
	return &mockTypeHandlerLogger{l: zerolog.New(io.Discard)}
}

func (m *mockTypeHandlerLogger) Info(ctx ...context.Context) *zerolog.Event  { return m.l.Info() }
func (m *mockTypeHandlerLogger) Error(ctx ...context.Context) *zerolog.Event { return m.l.Error() }
func (m *mockTypeHandlerLogger) Debug(ctx ...context.Context) *zerolog.Event { return m.l.Debug() }
func (m *mockTypeHandlerLogger) Warn(ctx ...context.Context) *zerolog.Event  { return m.l.Warn() }
func (m *mockTypeHandlerLogger) Fatal(ctx ...context.Context) *zerolog.Event { return m.l.Fatal() }
func (m *mockTypeHandlerLogger) Panic(ctx ...context.Context) *zerolog.Event { return m.l.Panic() }
func (m *mockTypeHandlerLogger) With() zerolog.Context                       { return m.l.With() }
func (m *mockTypeHandlerLogger) WithService(service string) log.ILogger      { return m }
func (m *mockTypeHandlerLogger) WithModule(module string) log.ILogger        { return m }
func (m *mockTypeHandlerLogger) WithBusinessID(businessID uint) log.ILogger  { return m }

// Mock de IConfig
type mockTypeHandlerConfig struct {
	mock.Mock
}

func (m *mockTypeHandlerConfig) Get(key string) string {
	args := m.Called(key)
	return args.String(0)
}

// typeHandlerSetup construye un IntegrationTypeHandler listo para tests
func typeHandlerSetup(uc *mockIntegrationTypeUseCase) (*IntegrationTypeHandler, *mockTypeHandlerConfig) {
	logger := newMockTypeHandlerLogger()
	cfg := new(mockTypeHandlerConfig)
	cfg.On("Get", mock.Anything).Return("").Maybe()

	h := &IntegrationTypeHandler{
		usecase: uc,
		logger:  logger,
		env:     cfg,
	}
	return h, cfg
}

func init() {
	gin.SetMode(gin.TestMode)
}

// ============================================
// ListIntegrationTypesHandler
// ============================================

func TestListIntegrationTypesHandler_RetornaLista(t *testing.T) {
	// Arrange
	uc := new(mockIntegrationTypeUseCase)
	h, _ := typeHandlerSetup(uc)

	tipos := []*domain.IntegrationType{
		{ID: 1, Name: "Shopify", Code: "shopify", IsActive: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 2, Name: "Factus", Code: "factus", IsActive: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	uc.On("ListIntegrationTypes", mock.Anything).Return(tipos, nil)

	r := gin.New()
	r.GET("/integration-types", h.ListIntegrationTypesHandler)

	req := httptest.NewRequest(http.MethodGet, "/integration-types", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestListIntegrationTypesHandler_ErrorDeUseCase(t *testing.T) {
	// Arrange
	uc := new(mockIntegrationTypeUseCase)
	h, _ := typeHandlerSetup(uc)

	uc.On("ListIntegrationTypes", mock.Anything).Return(nil, errors.New("db error"))

	r := gin.New()
	r.GET("/integration-types", h.ListIntegrationTypesHandler)

	req := httptest.NewRequest(http.MethodGet, "/integration-types", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ============================================
// ListActiveIntegrationTypesHandler
// ============================================

func TestListActiveIntegrationTypesHandler_RetornaActivos(t *testing.T) {
	// Arrange
	uc := new(mockIntegrationTypeUseCase)
	h, _ := typeHandlerSetup(uc)

	tipos := []*domain.IntegrationType{
		{ID: 1, Name: "Shopify", Code: "shopify", IsActive: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	uc.On("ListActiveIntegrationTypes", mock.Anything).Return(tipos, nil)

	r := gin.New()
	r.GET("/integration-types/active", h.ListActiveIntegrationTypesHandler)

	req := httptest.NewRequest(http.MethodGet, "/integration-types/active", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestListActiveIntegrationTypesHandler_ErrorDeUseCase(t *testing.T) {
	// Arrange
	uc := new(mockIntegrationTypeUseCase)
	h, _ := typeHandlerSetup(uc)

	uc.On("ListActiveIntegrationTypes", mock.Anything).Return(nil, errors.New("db error"))

	r := gin.New()
	r.GET("/integration-types/active", h.ListActiveIntegrationTypesHandler)

	req := httptest.NewRequest(http.MethodGet, "/integration-types/active", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ============================================
// ListIntegrationCategoriesHandler
// ============================================

func TestListIntegrationCategoriesHandler_RetornaCategorias(t *testing.T) {
	// Arrange
	uc := new(mockIntegrationTypeUseCase)
	h, _ := typeHandlerSetup(uc)

	categorias := []*domain.IntegrationCategory{
		{ID: 1, Code: "ecommerce", Name: "Ecommerce", IsActive: true, IsVisible: true},
		{ID: 2, Code: "invoicing", Name: "Facturación", IsActive: true, IsVisible: true},
	}
	uc.On("ListIntegrationCategories", mock.Anything).Return(categorias, nil)

	r := gin.New()
	r.GET("/integration-categories", h.ListIntegrationCategoriesHandler)

	req := httptest.NewRequest(http.MethodGet, "/integration-categories", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestListIntegrationCategoriesHandler_ErrorDeUseCase(t *testing.T) {
	// Arrange
	uc := new(mockIntegrationTypeUseCase)
	h, _ := typeHandlerSetup(uc)

	uc.On("ListIntegrationCategories", mock.Anything).Return(nil, errors.New("db error"))

	r := gin.New()
	r.GET("/integration-categories", h.ListIntegrationCategoriesHandler)

	req := httptest.NewRequest(http.MethodGet, "/integration-categories", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ============================================
// GetIntegrationTypeByIDHandler
// ============================================

func TestGetIntegrationTypeByIDHandler_IDValido(t *testing.T) {
	// Arrange
	uc := new(mockIntegrationTypeUseCase)
	h, _ := typeHandlerSetup(uc)

	tipo := &domain.IntegrationType{ID: 1, Name: "Shopify", Code: "shopify", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	uc.On("GetIntegrationTypeByID", mock.Anything, uint(1)).Return(tipo, nil)

	r := gin.New()
	r.GET("/integration-types/:id", h.GetIntegrationTypeByIDHandler)

	req := httptest.NewRequest(http.MethodGet, "/integration-types/1", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetIntegrationTypeByIDHandler_IDInvalido(t *testing.T) {
	// Arrange
	uc := new(mockIntegrationTypeUseCase)
	h, _ := typeHandlerSetup(uc)

	r := gin.New()
	r.GET("/integration-types/:id", h.GetIntegrationTypeByIDHandler)

	req := httptest.NewRequest(http.MethodGet, "/integration-types/abc", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetIntegrationTypeByIDHandler_NoEncontrado(t *testing.T) {
	// Arrange
	uc := new(mockIntegrationTypeUseCase)
	h, _ := typeHandlerSetup(uc)

	uc.On("GetIntegrationTypeByID", mock.Anything, uint(99)).Return(nil, domain.ErrIntegrationTypeNotFound)

	r := gin.New()
	r.GET("/integration-types/:id", h.GetIntegrationTypeByIDHandler)

	req := httptest.NewRequest(http.MethodGet, "/integration-types/99", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetIntegrationTypeByIDHandler_ErrorInterno(t *testing.T) {
	// Arrange
	uc := new(mockIntegrationTypeUseCase)
	h, _ := typeHandlerSetup(uc)

	uc.On("GetIntegrationTypeByID", mock.Anything, uint(5)).Return(nil, errors.New("db error"))

	r := gin.New()
	r.GET("/integration-types/:id", h.GetIntegrationTypeByIDHandler)

	req := httptest.NewRequest(http.MethodGet, "/integration-types/5", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ============================================
// GetIntegrationTypeByCodeHandler
// ============================================

func TestGetIntegrationTypeByCodeHandler_CodigoValido(t *testing.T) {
	// Arrange
	uc := new(mockIntegrationTypeUseCase)
	h, _ := typeHandlerSetup(uc)

	tipo := &domain.IntegrationType{ID: 1, Name: "Shopify", Code: "shopify", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	uc.On("GetIntegrationTypeByCode", mock.Anything, "shopify").Return(tipo, nil)

	r := gin.New()
	r.GET("/integration-types/code/:code", h.GetIntegrationTypeByCodeHandler)

	req := httptest.NewRequest(http.MethodGet, "/integration-types/code/shopify", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetIntegrationTypeByCodeHandler_NoEncontrado(t *testing.T) {
	// Arrange
	uc := new(mockIntegrationTypeUseCase)
	h, _ := typeHandlerSetup(uc)

	uc.On("GetIntegrationTypeByCode", mock.Anything, "unknown").Return(nil, domain.ErrIntegrationTypeNotFound)

	r := gin.New()
	r.GET("/integration-types/code/:code", h.GetIntegrationTypeByCodeHandler)

	req := httptest.NewRequest(http.MethodGet, "/integration-types/code/unknown", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
}
