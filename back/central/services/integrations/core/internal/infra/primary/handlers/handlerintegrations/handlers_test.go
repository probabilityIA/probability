package handlerintegrations

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ============================================
// Mock del Use Case (domain.IIntegrationUseCase)
// ============================================

type mockIntegrationUseCase struct {
	mock.Mock
}

func (m *mockIntegrationUseCase) CreateIntegration(ctx context.Context, dto domain.CreateIntegrationDTO) (*domain.Integration, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Integration), args.Error(1)
}

func (m *mockIntegrationUseCase) UpdateIntegration(ctx context.Context, id uint, dto domain.UpdateIntegrationDTO) (*domain.Integration, error) {
	args := m.Called(ctx, id, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Integration), args.Error(1)
}

func (m *mockIntegrationUseCase) GetIntegrationByID(ctx context.Context, id uint) (*domain.Integration, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Integration), args.Error(1)
}

func (m *mockIntegrationUseCase) GetIntegrationByIDWithCredentials(ctx context.Context, id uint) (*domain.IntegrationWithCredentials, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.IntegrationWithCredentials), args.Error(1)
}

func (m *mockIntegrationUseCase) GetIntegrationByType(ctx context.Context, code string, businessID *uint) (*domain.IntegrationWithCredentials, error) {
	args := m.Called(ctx, code, businessID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.IntegrationWithCredentials), args.Error(1)
}

func (m *mockIntegrationUseCase) GetPublicIntegrationByID(ctx context.Context, integrationID string) (*domain.PublicIntegration, error) {
	args := m.Called(ctx, integrationID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PublicIntegration), args.Error(1)
}

func (m *mockIntegrationUseCase) GetIntegrationConfig(ctx context.Context, integrationType string, businessID *uint) (map[string]interface{}, error) {
	args := m.Called(ctx, integrationType, businessID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *mockIntegrationUseCase) DecryptCredentialField(ctx context.Context, integrationID string, fieldName string) (string, error) {
	args := m.Called(ctx, integrationID, fieldName)
	return args.String(0), args.Error(1)
}

func (m *mockIntegrationUseCase) DeleteIntegration(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockIntegrationUseCase) ListIntegrations(ctx context.Context, filters domain.IntegrationFilters) ([]*domain.Integration, int64, error) {
	args := m.Called(ctx, filters)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*domain.Integration), args.Get(1).(int64), args.Error(2)
}

func (m *mockIntegrationUseCase) ActivateIntegration(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockIntegrationUseCase) DeactivateIntegration(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockIntegrationUseCase) SetAsDefault(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockIntegrationUseCase) UpdateLastSync(ctx context.Context, integrationID string) error {
	args := m.Called(ctx, integrationID)
	return args.Error(0)
}

func (m *mockIntegrationUseCase) TestIntegration(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockIntegrationUseCase) TestConnectionRaw(ctx context.Context, integrationTypeCode string, config map[string]interface{}, credentials map[string]interface{}) error {
	args := m.Called(ctx, integrationTypeCode, config, credentials)
	return args.Error(0)
}

func (m *mockIntegrationUseCase) RegisterObserver(observer domain.IntegrationCreatedObserver) {
	m.Called(observer)
}

func (m *mockIntegrationUseCase) WarmCache(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *mockIntegrationUseCase) RegisterProvider(integrationType int, provider domain.IIntegrationContract) {
	m.Called(integrationType, provider)
}

func (m *mockIntegrationUseCase) GetProvider(integrationType int) (domain.IIntegrationContract, bool) {
	args := m.Called(integrationType)
	if args.Get(0) == nil {
		return nil, args.Bool(1)
	}
	return args.Get(0).(domain.IIntegrationContract), args.Bool(1)
}

func (m *mockIntegrationUseCase) SyncOrdersByIntegrationID(ctx context.Context, integrationID string) error {
	args := m.Called(ctx, integrationID)
	return args.Error(0)
}

func (m *mockIntegrationUseCase) SyncOrdersByIntegrationIDWithParams(ctx context.Context, integrationID string, params interface{}) error {
	args := m.Called(ctx, integrationID, params)
	return args.Error(0)
}

func (m *mockIntegrationUseCase) SyncOrdersByBusiness(ctx context.Context, businessID uint) error {
	args := m.Called(ctx, businessID)
	return args.Error(0)
}

func (m *mockIntegrationUseCase) GetWebhookURL(ctx context.Context, integrationID uint) (*domain.WebhookInfo, error) {
	args := m.Called(ctx, integrationID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.WebhookInfo), args.Error(1)
}

func (m *mockIntegrationUseCase) ListWebhooks(ctx context.Context, integrationID string) ([]interface{}, error) {
	args := m.Called(ctx, integrationID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]interface{}), args.Error(1)
}

func (m *mockIntegrationUseCase) DeleteWebhook(ctx context.Context, integrationID, webhookID string) error {
	args := m.Called(ctx, integrationID, webhookID)
	return args.Error(0)
}

func (m *mockIntegrationUseCase) VerifyWebhooksByURL(ctx context.Context, integrationID string) ([]interface{}, error) {
	args := m.Called(ctx, integrationID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]interface{}), args.Error(1)
}

func (m *mockIntegrationUseCase) CreateWebhookForIntegration(ctx context.Context, integrationID string) (interface{}, error) {
	args := m.Called(ctx, integrationID)
	return args.Get(0), args.Error(1)
}

func (m *mockIntegrationUseCase) GetIntegrationByExternalID(ctx context.Context, externalID string, integrationType int) (*domain.PublicIntegration, error) {
	args := m.Called(ctx, externalID, integrationType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PublicIntegration), args.Error(1)
}

func (m *mockIntegrationUseCase) UpdateIntegrationConfig(ctx context.Context, integrationID string, newConfig map[string]interface{}) error {
	args := m.Called(ctx, integrationID, newConfig)
	return args.Error(0)
}

func (m *mockIntegrationUseCase) TestConnectionFromConfig(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error {
	args := m.Called(ctx, config, credentials)
	return args.Error(0)
}

func (m *mockIntegrationUseCase) OnIntegrationCreated(integrationType int, observer func(context.Context, *domain.PublicIntegration)) {
	m.Called(integrationType, observer)
}

func (m *mockIntegrationUseCase) GetPlatformCredentialByIntegrationID(ctx context.Context, integrationID string, fieldName string) (string, error) {
	args := m.Called(ctx, integrationID, fieldName)
	return args.String(0), args.Error(1)
}

// ============================================
// Mock de IConfig y ILogger para handlers
// ============================================

type mockHandlerConfig struct {
	mock.Mock
}

func (m *mockHandlerConfig) Get(key string) string {
	args := m.Called(key)
	return args.String(0)
}

type mockHandlerLogger struct {
	l zerolog.Logger
}

func newMockHandlerLogger() *mockHandlerLogger {
	return &mockHandlerLogger{l: zerolog.New(io.Discard)}
}

func (m *mockHandlerLogger) Info(ctx ...context.Context) *zerolog.Event  { return m.l.Info() }
func (m *mockHandlerLogger) Error(ctx ...context.Context) *zerolog.Event { return m.l.Error() }
func (m *mockHandlerLogger) Debug(ctx ...context.Context) *zerolog.Event { return m.l.Debug() }
func (m *mockHandlerLogger) Warn(ctx ...context.Context) *zerolog.Event  { return m.l.Warn() }
func (m *mockHandlerLogger) Fatal(ctx ...context.Context) *zerolog.Event { return m.l.Fatal() }
func (m *mockHandlerLogger) Panic(ctx ...context.Context) *zerolog.Event { return m.l.Panic() }
func (m *mockHandlerLogger) With() zerolog.Context                       { return m.l.With() }

// Implementar WithService, WithModule, WithBusinessID retornando log.ILogger
func (m *mockHandlerLogger) WithService(service string) log.ILogger    { return m }
func (m *mockHandlerLogger) WithModule(module string) log.ILogger      { return m }
func (m *mockHandlerLogger) WithBusinessID(businessID uint) log.ILogger { return m }

// handlerSetup construye un handler listo para tests
func handlerSetup(uc *mockIntegrationUseCase) (*IntegrationHandler, *mockHandlerConfig) {
	logger := newMockHandlerLogger()
	cfg := new(mockHandlerConfig)
	cfg.On("Get", mock.Anything).Return("").Maybe()

	h := &IntegrationHandler{
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
// GetIntegrationsHandler
// ============================================

func TestGetIntegrationsHandler_RetornaLista(t *testing.T) {
	// Arrange
	uc := new(mockIntegrationUseCase)
	h, _ := handlerSetup(uc)

	integraciones := []*domain.Integration{
		{ID: 1, Name: "Shopify"},
		{ID: 2, Name: "Factus"},
	}
	uc.On("ListIntegrations", mock.Anything, mock.Anything).Return(integraciones, int64(2), nil)

	r := gin.New()
	r.GET("/integrations", h.GetIntegrationsHandler)

	req := httptest.NewRequest(http.MethodGet, "/integrations?page=1&page_size=10", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var body map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &body)
	assert.Equal(t, true, body["success"])
	assert.Equal(t, float64(2), body["total"])
}

func TestGetIntegrationsHandler_ErrorDeUseCase(t *testing.T) {
	// Arrange
	uc := new(mockIntegrationUseCase)
	h, _ := handlerSetup(uc)

	uc.On("ListIntegrations", mock.Anything, mock.Anything).Return(nil, int64(0), errors.New("db error"))

	r := gin.New()
	r.GET("/integrations", h.GetIntegrationsHandler)

	req := httptest.NewRequest(http.MethodGet, "/integrations", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var body map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &body)
	assert.Equal(t, false, body["success"])
}

// ============================================
// GetIntegrationByIDHandler
// ============================================

func TestGetIntegrationByIDHandler_IDValido(t *testing.T) {
	// Arrange
	uc := new(mockIntegrationUseCase)
	h, _ := handlerSetup(uc)

	integracion := &domain.Integration{ID: 5, Name: "Shopify Store"}
	uc.On("GetIntegrationByID", mock.Anything, uint(5)).Return(integracion, nil)

	r := gin.New()
	r.GET("/integrations/:id", h.GetIntegrationByIDHandler)

	req := httptest.NewRequest(http.MethodGet, "/integrations/5", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var body map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &body)
	assert.Equal(t, true, body["success"])
}

func TestGetIntegrationByIDHandler_IDInvalido(t *testing.T) {
	// Arrange
	uc := new(mockIntegrationUseCase)
	h, _ := handlerSetup(uc)

	r := gin.New()
	r.GET("/integrations/:id", h.GetIntegrationByIDHandler)

	req := httptest.NewRequest(http.MethodGet, "/integrations/abc", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetIntegrationByIDHandler_NoEncontrado(t *testing.T) {
	// Arrange
	uc := new(mockIntegrationUseCase)
	h, _ := handlerSetup(uc)

	uc.On("GetIntegrationByID", mock.Anything, uint(99)).Return(nil, errors.New("not found"))

	r := gin.New()
	r.GET("/integrations/:id", h.GetIntegrationByIDHandler)

	req := httptest.NewRequest(http.MethodGet, "/integrations/99", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert — Sin ErrIntegrationNotFound explícito retorna 500
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ============================================
// DeleteIntegrationHandler
// ============================================

func TestDeleteIntegrationHandler_EliminaConPermisos(t *testing.T) {
	// Arrange
	uc := new(mockIntegrationUseCase)
	h, _ := handlerSetup(uc)

	uc.On("DeleteIntegration", mock.Anything, uint(5)).Return(nil)

	r := gin.New()
	// IsSuperAdmin verifica business_id == 0
	r.DELETE("/integrations/:id", func(c *gin.Context) {
		c.Set("business_id", uint(0)) // super admin = business_id 0
		h.DeleteIntegrationHandler(c)
	})

	req := httptest.NewRequest(http.MethodDelete, "/integrations/5", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteIntegrationHandler_SinPermisos_Retorna403(t *testing.T) {
	// Arrange
	uc := new(mockIntegrationUseCase)
	h, _ := handlerSetup(uc)

	r := gin.New()
	// Sin super admin
	r.DELETE("/integrations/:id", h.DeleteIntegrationHandler)

	req := httptest.NewRequest(http.MethodDelete, "/integrations/5", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestDeleteIntegrationHandler_IDInvalido(t *testing.T) {
	// Arrange
	uc := new(mockIntegrationUseCase)
	h, _ := handlerSetup(uc)

	r := gin.New()
	r.DELETE("/integrations/:id", func(c *gin.Context) {
		c.Set("business_id", uint(0)) // super admin
		h.DeleteIntegrationHandler(c)
	})

	req := httptest.NewRequest(http.MethodDelete, "/integrations/abc", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteIntegrationHandler_ErrorNoEncontrado(t *testing.T) {
	// Arrange
	uc := new(mockIntegrationUseCase)
	h, _ := handlerSetup(uc)

	uc.On("DeleteIntegration", mock.Anything, uint(99)).Return(domain.ErrIntegrationNotFound)

	r := gin.New()
	r.DELETE("/integrations/:id", func(c *gin.Context) {
		c.Set("business_id", uint(0)) // super admin
		h.DeleteIntegrationHandler(c)
	})

	req := httptest.NewRequest(http.MethodDelete, "/integrations/99", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ============================================
// SyncOrdersByIntegrationIDHandler
// ============================================

func TestSyncOrdersByIntegrationIDHandler_SinBody(t *testing.T) {
	// Arrange
	uc := new(mockIntegrationUseCase)
	h, _ := handlerSetup(uc)

	uc.On("SyncOrdersByIntegrationID", mock.Anything, "5").Return(nil)

	r := gin.New()
	r.POST("/integrations/:id/sync", h.SyncOrdersByIntegrationIDHandler)

	req := httptest.NewRequest(http.MethodPost, "/integrations/5/sync", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusAccepted, w.Code)
}

func TestSyncOrdersByIntegrationIDHandler_ErrorDeUseCase(t *testing.T) {
	// Arrange
	uc := new(mockIntegrationUseCase)
	h, _ := handlerSetup(uc)

	uc.On("SyncOrdersByIntegrationID", mock.Anything, "5").Return(errors.New("sync error"))

	r := gin.New()
	r.POST("/integrations/:id/sync", h.SyncOrdersByIntegrationIDHandler)

	req := httptest.NewRequest(http.MethodPost, "/integrations/5/sync", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestSyncOrdersByIntegrationIDHandler_ConFiltros(t *testing.T) {
	// Arrange
	uc := new(mockIntegrationUseCase)
	h, _ := handlerSetup(uc)

	uc.On("SyncOrdersByIntegrationIDWithParams", mock.Anything, "5", mock.Anything).Return(nil)

	r := gin.New()
	r.POST("/integrations/:id/sync", h.SyncOrdersByIntegrationIDHandler)

	body := bytes.NewBufferString(`{"status":"paid","financial_status":"paid"}`)
	req := httptest.NewRequest(http.MethodPost, "/integrations/5/sync", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusAccepted, w.Code)
}

// ============================================
// GetWebhookURLHandler
// ============================================

func TestGetWebhookURLHandler_IDValido(t *testing.T) {
	// Arrange
	uc := new(mockIntegrationUseCase)
	h, _ := handlerSetup(uc)

	webhookInfo := &domain.WebhookInfo{
		URL:    "https://api.example.com/webhook",
		Method: "POST",
	}
	uc.On("GetWebhookURL", mock.Anything, uint(5)).Return(webhookInfo, nil)

	r := gin.New()
	r.GET("/integrations/:id/webhook", h.GetWebhookURLHandler)

	req := httptest.NewRequest(http.MethodGet, "/integrations/5/webhook", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var body map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &body)
	assert.Equal(t, true, body["success"])
}

func TestGetWebhookURLHandler_IDInvalido(t *testing.T) {
	// Arrange
	uc := new(mockIntegrationUseCase)
	h, _ := handlerSetup(uc)

	r := gin.New()
	r.GET("/integrations/:id/webhook", h.GetWebhookURLHandler)

	req := httptest.NewRequest(http.MethodGet, "/integrations/abc/webhook", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetWebhookURLHandler_ErrorDeUseCase(t *testing.T) {
	// Arrange
	uc := new(mockIntegrationUseCase)
	h, _ := handlerSetup(uc)

	uc.On("GetWebhookURL", mock.Anything, uint(5)).Return(nil, errors.New("config missing"))

	r := gin.New()
	r.GET("/integrations/:id/webhook", h.GetWebhookURLHandler)

	req := httptest.NewRequest(http.MethodGet, "/integrations/5/webhook", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ============================================
// ListWebhooksHandler
// ============================================

func TestListWebhooksHandler_Exitoso(t *testing.T) {
	// Arrange
	uc := new(mockIntegrationUseCase)
	h, _ := handlerSetup(uc)

	webhooks := []interface{}{map[string]interface{}{"id": "wh_1"}}
	uc.On("ListWebhooks", mock.Anything, "5").Return(webhooks, nil)

	r := gin.New()
	r.GET("/integrations/:id/webhooks", h.ListWebhooksHandler)

	req := httptest.NewRequest(http.MethodGet, "/integrations/5/webhooks", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestListWebhooksHandler_ErrorDeUseCase(t *testing.T) {
	// Arrange
	uc := new(mockIntegrationUseCase)
	h, _ := handlerSetup(uc)

	uc.On("ListWebhooks", mock.Anything, "5").Return(nil, errors.New("api error"))

	r := gin.New()
	r.GET("/integrations/:id/webhooks", h.ListWebhooksHandler)

	req := httptest.NewRequest(http.MethodGet, "/integrations/5/webhooks", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ============================================
// DeleteWebhookHandler
// ============================================

func TestDeleteWebhookHandler_Exitoso(t *testing.T) {
	// Arrange
	uc := new(mockIntegrationUseCase)
	h, _ := handlerSetup(uc)

	uc.On("DeleteWebhook", mock.Anything, "5", "wh_123").Return(nil)

	r := gin.New()
	r.DELETE("/integrations/:id/webhooks/:webhook_id", h.DeleteWebhookHandler)

	req := httptest.NewRequest(http.MethodDelete, "/integrations/5/webhooks/wh_123", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteWebhookHandler_ErrorDeUseCase(t *testing.T) {
	// Arrange
	uc := new(mockIntegrationUseCase)
	h, _ := handlerSetup(uc)

	uc.On("DeleteWebhook", mock.Anything, "5", "wh_999").Return(errors.New("webhook not found"))

	r := gin.New()
	r.DELETE("/integrations/:id/webhooks/:webhook_id", h.DeleteWebhookHandler)

	req := httptest.NewRequest(http.MethodDelete, "/integrations/5/webhooks/wh_999", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ============================================
// VerifyWebhooksHandler
// ============================================

func TestVerifyWebhooksHandler_Exitoso(t *testing.T) {
	// Arrange
	uc := new(mockIntegrationUseCase)
	h, _ := handlerSetup(uc)

	webhooks := []interface{}{map[string]interface{}{"id": "wh_1"}}
	uc.On("VerifyWebhooksByURL", mock.Anything, "5").Return(webhooks, nil)

	r := gin.New()
	r.GET("/integrations/:id/webhooks/verify", h.VerifyWebhooksHandler)

	req := httptest.NewRequest(http.MethodGet, "/integrations/5/webhooks/verify", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var body map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &body)
	assert.Equal(t, true, body["success"])
}

func TestVerifyWebhooksHandler_ErrorDeUseCase(t *testing.T) {
	// Arrange
	uc := new(mockIntegrationUseCase)
	h, _ := handlerSetup(uc)

	uc.On("VerifyWebhooksByURL", mock.Anything, "5").Return(nil, errors.New("base url missing"))

	r := gin.New()
	r.GET("/integrations/:id/webhooks/verify", h.VerifyWebhooksHandler)

	req := httptest.NewRequest(http.MethodGet, "/integrations/5/webhooks/verify", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ============================================
// CreateWebhookHandler
// ============================================

func TestCreateWebhookHandler_Exitoso(t *testing.T) {
	// Arrange
	uc := new(mockIntegrationUseCase)
	h, _ := handlerSetup(uc)

	webhookResult := map[string]interface{}{"id": "new_wh"}
	uc.On("CreateWebhookForIntegration", mock.Anything, "5").Return(webhookResult, nil)

	r := gin.New()
	r.POST("/integrations/:id/webhooks", h.CreateWebhookHandler)

	req := httptest.NewRequest(http.MethodPost, "/integrations/5/webhooks", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateWebhookHandler_ErrorDeUseCase(t *testing.T) {
	// Arrange
	uc := new(mockIntegrationUseCase)
	h, _ := handlerSetup(uc)

	uc.On("CreateWebhookForIntegration", mock.Anything, "5").Return(nil, errors.New("base url not configured"))

	r := gin.New()
	r.POST("/integrations/:id/webhooks", h.CreateWebhookHandler)

	req := httptest.NewRequest(http.MethodPost, "/integrations/5/webhooks", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ============================================
// ActivateIntegrationHandler
// ============================================

func TestActivateIntegrationHandler_Exitoso(t *testing.T) {
	// Arrange
	uc := new(mockIntegrationUseCase)
	h, _ := handlerSetup(uc)

	uc.On("ActivateIntegration", mock.Anything, uint(5)).Return(nil)

	r := gin.New()
	r.POST("/integrations/:id/activate", func(c *gin.Context) {
		c.Set("business_id", uint(0)) // super admin
		h.ActivateIntegrationHandler(c)
	})

	req := httptest.NewRequest(http.MethodPost, "/integrations/5/activate", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestActivateIntegrationHandler_IDInvalido(t *testing.T) {
	// Arrange
	uc := new(mockIntegrationUseCase)
	h, _ := handlerSetup(uc)

	r := gin.New()
	r.POST("/integrations/:id/activate", func(c *gin.Context) {
		c.Set("business_id", uint(0)) // super admin
		h.ActivateIntegrationHandler(c)
	})

	req := httptest.NewRequest(http.MethodPost, "/integrations/abc/activate", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ============================================
// DeactivateIntegrationHandler
// ============================================

func TestDeactivateIntegrationHandler_Exitoso(t *testing.T) {
	// Arrange
	uc := new(mockIntegrationUseCase)
	h, _ := handlerSetup(uc)

	uc.On("DeactivateIntegration", mock.Anything, uint(5)).Return(nil)

	r := gin.New()
	r.POST("/integrations/:id/deactivate", func(c *gin.Context) {
		c.Set("business_id", uint(0)) // super admin
		h.DeactivateIntegrationHandler(c)
	})

	req := httptest.NewRequest(http.MethodPost, "/integrations/5/deactivate", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeactivateIntegrationHandler_ErrorDeUseCase(t *testing.T) {
	// Arrange
	uc := new(mockIntegrationUseCase)
	h, _ := handlerSetup(uc)

	uc.On("DeactivateIntegration", mock.Anything, uint(5)).Return(errors.New("error deactivating"))

	r := gin.New()
	r.POST("/integrations/:id/deactivate", func(c *gin.Context) {
		c.Set("business_id", uint(0)) // super admin
		h.DeactivateIntegrationHandler(c)
	})

	req := httptest.NewRequest(http.MethodPost, "/integrations/5/deactivate", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ============================================
// TestIntegrationHandler
// ============================================

func TestTestIntegrationHandler_Exitoso(t *testing.T) {
	// Arrange
	uc := new(mockIntegrationUseCase)
	h, _ := handlerSetup(uc)

	uc.On("TestIntegration", mock.Anything, uint(5)).Return(nil)

	r := gin.New()
	// IsSuperAdmin verifica business_id == 0
	r.POST("/integrations/:id/test", func(c *gin.Context) {
		c.Set("business_id", uint(0)) // super admin
		h.TestIntegrationHandler(c)
	})

	req := httptest.NewRequest(http.MethodPost, "/integrations/5/test", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestTestIntegrationHandler_SinPermisos_Retorna403(t *testing.T) {
	// Arrange
	uc := new(mockIntegrationUseCase)
	h, _ := handlerSetup(uc)

	r := gin.New()
	// Sin business_id configurado → IsSuperAdmin retorna false → 403
	r.POST("/integrations/:id/test", h.TestIntegrationHandler)

	req := httptest.NewRequest(http.MethodPost, "/integrations/5/test", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestTestIntegrationHandler_IDInvalido(t *testing.T) {
	// Arrange
	uc := new(mockIntegrationUseCase)
	h, _ := handlerSetup(uc)

	r := gin.New()
	r.POST("/integrations/:id/test", func(c *gin.Context) {
		c.Set("business_id", uint(0)) // super admin
		h.TestIntegrationHandler(c)
	})

	req := httptest.NewRequest(http.MethodPost, "/integrations/abc/test", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTestIntegrationHandler_ErrorDeUseCase(t *testing.T) {
	// Arrange
	uc := new(mockIntegrationUseCase)
	h, _ := handlerSetup(uc)

	uc.On("TestIntegration", mock.Anything, uint(5)).Return(errors.New("connection failed"))

	r := gin.New()
	r.POST("/integrations/:id/test", func(c *gin.Context) {
		c.Set("business_id", uint(0)) // super admin
		h.TestIntegrationHandler(c)
	})

	req := httptest.NewRequest(http.MethodPost, "/integrations/5/test", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ============================================
// GetIntegrationByTypeHandler
// ============================================

func TestGetIntegrationByTypeHandler_Exitoso(t *testing.T) {
	// Arrange
	uc := new(mockIntegrationUseCase)
	h, _ := handlerSetup(uc)

	integracion := &domain.IntegrationWithCredentials{
		Integration: domain.Integration{ID: 1, Name: "Shopify", Code: "shopify"},
	}
	uc.On("GetIntegrationByType", mock.Anything, "shopify", mock.Anything).Return(integracion, nil)

	r := gin.New()
	r.GET("/integrations/type/:type", h.GetIntegrationByTypeHandler)

	req := httptest.NewRequest(http.MethodGet, "/integrations/type/shopify", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetIntegrationByTypeHandler_NoEncontrado(t *testing.T) {
	// Arrange
	uc := new(mockIntegrationUseCase)
	h, _ := handlerSetup(uc)

	uc.On("GetIntegrationByType", mock.Anything, "unknown", mock.Anything).Return(nil, domain.ErrIntegrationNotFound)

	r := gin.New()
	r.GET("/integrations/type/:type", h.GetIntegrationByTypeHandler)

	req := httptest.NewRequest(http.MethodGet, "/integrations/type/unknown", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetIntegrationByTypeHandler_ErrorInterno(t *testing.T) {
	// Arrange
	uc := new(mockIntegrationUseCase)
	h, _ := handlerSetup(uc)

	uc.On("GetIntegrationByType", mock.Anything, "shopify", mock.Anything).Return(nil, errors.New("db error"))

	r := gin.New()
	r.GET("/integrations/type/:type", h.GetIntegrationByTypeHandler)

	req := httptest.NewRequest(http.MethodGet, "/integrations/type/shopify", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ============================================
// SetAsDefaultHandler
// ============================================

func TestSetAsDefaultHandler_Exitoso(t *testing.T) {
	// Arrange
	uc := new(mockIntegrationUseCase)
	h, _ := handlerSetup(uc)

	uc.On("SetAsDefault", mock.Anything, uint(5)).Return(nil)

	r := gin.New()
	r.PUT("/integrations/:id/set-default", func(c *gin.Context) {
		c.Set("business_id", uint(0)) // super admin
		h.SetAsDefaultHandler(c)
	})

	req := httptest.NewRequest(http.MethodPut, "/integrations/5/set-default", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSetAsDefaultHandler_SinPermisos_Retorna403(t *testing.T) {
	// Arrange
	uc := new(mockIntegrationUseCase)
	h, _ := handlerSetup(uc)

	r := gin.New()
	// Sin business_id — IsSuperAdmin retorna false
	r.PUT("/integrations/:id/set-default", h.SetAsDefaultHandler)

	req := httptest.NewRequest(http.MethodPut, "/integrations/5/set-default", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestSetAsDefaultHandler_IDInvalido(t *testing.T) {
	// Arrange
	uc := new(mockIntegrationUseCase)
	h, _ := handlerSetup(uc)

	r := gin.New()
	r.PUT("/integrations/:id/set-default", func(c *gin.Context) {
		c.Set("business_id", uint(0)) // super admin
		h.SetAsDefaultHandler(c)
	})

	req := httptest.NewRequest(http.MethodPut, "/integrations/abc/set-default", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSetAsDefaultHandler_ErrorDeUseCase(t *testing.T) {
	// Arrange
	uc := new(mockIntegrationUseCase)
	h, _ := handlerSetup(uc)

	uc.On("SetAsDefault", mock.Anything, uint(5)).Return(errors.New("db error"))

	r := gin.New()
	r.PUT("/integrations/:id/set-default", func(c *gin.Context) {
		c.Set("business_id", uint(0)) // super admin
		h.SetAsDefaultHandler(c)
	})

	req := httptest.NewRequest(http.MethodPut, "/integrations/5/set-default", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
