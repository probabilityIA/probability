package cache

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	goredis "github.com/redis/go-redis/v9"
)

// ============================================
// Mock de redis.IRedis
// ============================================

type mockRedis struct {
	mock.Mock
}

func (m *mockRedis) Connect(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}

func (m *mockRedis) Close() error {
	return m.Called().Error(0)
}

func (m *mockRedis) Client(ctx context.Context) *goredis.Client {
	return nil
}

func (m *mockRedis) Ping(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}

func (m *mockRedis) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *mockRedis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return m.Called(ctx, key, value, expiration).Error(0)
}

func (m *mockRedis) Delete(ctx context.Context, keys ...string) error {
	args := []interface{}{ctx}
	for _, k := range keys {
		args = append(args, k)
	}
	return m.Called(args...).Error(0)
}

func (m *mockRedis) Exists(ctx context.Context, keys ...string) (int64, error) {
	args := m.Called(ctx, keys)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockRedis) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return m.Called(ctx, key, expiration).Error(0)
}

func (m *mockRedis) TTL(ctx context.Context, key string) (time.Duration, error) {
	args := m.Called(ctx, key)
	return args.Get(0).(time.Duration), args.Error(1)
}

func (m *mockRedis) Keys(ctx context.Context, pattern string) ([]string, error) {
	args := m.Called(ctx, pattern)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *mockRedis) Incr(ctx context.Context, key string) (int64, error) {
	args := m.Called(ctx, key)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockRedis) Decr(ctx context.Context, key string) (int64, error) {
	args := m.Called(ctx, key)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockRedis) HGet(ctx context.Context, key, field string) (string, error) {
	args := m.Called(ctx, key, field)
	return args.String(0), args.Error(1)
}

func (m *mockRedis) HSet(ctx context.Context, key string, values ...interface{}) error {
	return m.Called(ctx, key, values).Error(0)
}

func (m *mockRedis) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]string), args.Error(1)
}

func (m *mockRedis) HDel(ctx context.Context, key string, fields ...string) error {
	return m.Called(ctx, key, fields).Error(0)
}

func (m *mockRedis) RegisterCachePrefix(prefix string) {
	m.Called(prefix)
}

func (m *mockRedis) RegisterChannel(channel string) {
	m.Called(channel)
}

// ============================================
// Mock de log.ILogger para cache
// ============================================

type mockCacheLogger struct {
	l zerolog.Logger
}

func newMockCacheLogger() *mockCacheLogger {
	return &mockCacheLogger{l: zerolog.New(io.Discard)}
}

func (m *mockCacheLogger) Info(ctx ...context.Context) *zerolog.Event  { return m.l.Info() }
func (m *mockCacheLogger) Error(ctx ...context.Context) *zerolog.Event { return m.l.Error() }
func (m *mockCacheLogger) Debug(ctx ...context.Context) *zerolog.Event { return m.l.Debug() }
func (m *mockCacheLogger) Warn(ctx ...context.Context) *zerolog.Event  { return m.l.Warn() }
func (m *mockCacheLogger) Fatal(ctx ...context.Context) *zerolog.Event { return m.l.Fatal() }
func (m *mockCacheLogger) Panic(ctx ...context.Context) *zerolog.Event { return m.l.Panic() }
func (m *mockCacheLogger) With() zerolog.Context                       { return m.l.With() }
func (m *mockCacheLogger) WithService(s string) log.ILogger            { return m }
func (m *mockCacheLogger) WithModule(s string) log.ILogger             { return m }
func (m *mockCacheLogger) WithBusinessID(id uint) log.ILogger          { return m }

// cacheSetup crea una instancia del cache con mocks
func cacheSetup(t *testing.T) (*IntegrationCache, *mockRedis) {
	t.Helper()
	r := new(mockRedis)
	logger := newMockCacheLogger()
	c := &IntegrationCache{redis: r, log: logger}
	return c, r
}

// ============================================
// SetIntegration
// ============================================

func TestSetIntegration_Exitoso(t *testing.T) {
	// Arrange
	c, r := cacheSetup(t)
	ctx := context.Background()

	bizID := uint(10)
	integracion := &domain.CachedIntegration{
		ID:                1,
		Code:              "shopify_store",
		BusinessID:        &bizID,
		IntegrationTypeID: 1,
	}

	// El mock responde OK a los tres Set (metadata, code index, biz+type index)
	r.On("Set", mock.Anything, "integration:meta:1", mock.Anything, ttlMetadata).Return(nil)
	r.On("Set", mock.Anything, "integration:code:shopify_store", mock.Anything, ttlMetadata).Return(nil)
	r.On("Set", mock.Anything, "integration:idx:biz:10:type:1", mock.Anything, ttlMetadata).Return(nil)

	// Act
	err := c.SetIntegration(ctx, integracion)

	// Assert
	assert.NoError(t, err)
	r.AssertCalled(t, "Set", mock.Anything, "integration:meta:1", mock.Anything, ttlMetadata)
}

func TestSetIntegration_ErrorAlCachearMetadata(t *testing.T) {
	// Arrange
	c, r := cacheSetup(t)
	ctx := context.Background()

	integracion := &domain.CachedIntegration{ID: 2}
	r.On("Set", mock.Anything, "integration:meta:2", mock.Anything, ttlMetadata).Return(errors.New("redis down"))

	// Act
	err := c.SetIntegration(ctx, integracion)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "redis down")
}

func TestSetIntegration_SinCodigo_NoCreaMacIndice(t *testing.T) {
	// Arrange
	c, r := cacheSetup(t)
	ctx := context.Background()

	integracion := &domain.CachedIntegration{
		ID:   3,
		Code: "", // Sin código → no debe crear índice de código
	}
	r.On("Set", mock.Anything, "integration:meta:3", mock.Anything, ttlMetadata).Return(nil)

	// Act
	err := c.SetIntegration(ctx, integracion)

	// Assert
	assert.NoError(t, err)
	// El índice por código NO debe llamarse
	r.AssertNotCalled(t, "Set", mock.Anything, mock.MatchedBy(func(key string) bool {
		return key == "integration:code:"
	}), mock.Anything, mock.Anything)
}

// ============================================
// GetIntegration
// ============================================

func TestGetIntegration_CacheHit(t *testing.T) {
	// Arrange
	c, r := cacheSetup(t)
	ctx := context.Background()

	integracion := &domain.CachedIntegration{
		ID:   5,
		Code: "factus_main",
	}
	data, _ := json.Marshal(integracion)

	r.On("Get", mock.Anything, "integration:meta:5").Return(string(data), nil)

	// Act
	resultado, err := c.GetIntegration(ctx, 5)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resultado)
	assert.Equal(t, uint(5), resultado.ID)
	assert.Equal(t, "factus_main", resultado.Code)
}

func TestGetIntegration_CacheMiss(t *testing.T) {
	// Arrange
	c, r := cacheSetup(t)
	ctx := context.Background()

	r.On("Get", mock.Anything, "integration:meta:99").Return("", errors.New("cache miss"))

	// Act
	resultado, err := c.GetIntegration(ctx, 99)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resultado)
}

func TestGetIntegration_JSONInvalidoRetornaError(t *testing.T) {
	// Arrange
	c, r := cacheSetup(t)
	ctx := context.Background()

	r.On("Get", mock.Anything, "integration:meta:10").Return("invalid json", nil)

	// Act
	resultado, err := c.GetIntegration(ctx, 10)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resultado)
}

// ============================================
// SetCredentials / GetCredentials
// ============================================

func TestSetCredentials_Exitoso(t *testing.T) {
	// Arrange
	c, r := cacheSetup(t)
	ctx := context.Background()

	creds := &domain.CachedCredentials{
		IntegrationID: 1,
		Credentials: map[string]interface{}{
			"api_key": "key_123",
		},
	}
	r.On("Set", mock.Anything, "integration:creds:1", mock.Anything, ttlCredentials).Return(nil)

	// Act
	err := c.SetCredentials(ctx, creds)

	// Assert
	assert.NoError(t, err)
	// CachedAt debe haberse seteado
	assert.False(t, creds.CachedAt.IsZero())
}

func TestGetCredentials_CacheHit(t *testing.T) {
	// Arrange
	c, r := cacheSetup(t)
	ctx := context.Background()

	creds := &domain.CachedCredentials{
		IntegrationID: 1,
		Credentials: map[string]interface{}{
			"api_key": "key_abc",
		},
		CachedAt: time.Now(),
	}
	data, _ := json.Marshal(creds)
	r.On("Get", mock.Anything, "integration:creds:1").Return(string(data), nil)

	// Act
	resultado, err := c.GetCredentials(ctx, 1)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resultado)
	assert.Equal(t, uint(1), resultado.IntegrationID)
}

func TestGetCredentials_CacheMiss(t *testing.T) {
	// Arrange
	c, r := cacheSetup(t)
	ctx := context.Background()

	r.On("Get", mock.Anything, "integration:creds:50").Return("", errors.New("key not found"))

	// Act
	resultado, err := c.GetCredentials(ctx, 50)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resultado)
}

// ============================================
// GetCredentialField
// ============================================

func TestGetCredentialField_CampoExistente(t *testing.T) {
	// Arrange
	c, r := cacheSetup(t)
	ctx := context.Background()

	creds := &domain.CachedCredentials{
		IntegrationID: 1,
		Credentials: map[string]interface{}{
			"shop_domain": "mi-tienda.myshopify.com",
		},
		CachedAt: time.Now(),
	}
	data, _ := json.Marshal(creds)
	r.On("Get", mock.Anything, "integration:creds:1").Return(string(data), nil)

	// Act
	valor, err := c.GetCredentialField(ctx, 1, "shop_domain")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "mi-tienda.myshopify.com", valor)
}

func TestGetCredentialField_CampoNoExistente(t *testing.T) {
	// Arrange
	c, r := cacheSetup(t)
	ctx := context.Background()

	creds := &domain.CachedCredentials{
		IntegrationID: 1,
		Credentials:   map[string]interface{}{},
		CachedAt:      time.Now(),
	}
	data, _ := json.Marshal(creds)
	r.On("Get", mock.Anything, "integration:creds:1").Return(string(data), nil)

	// Act
	valor, err := c.GetCredentialField(ctx, 1, "campo_inexistente")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "campo_inexistente")
	assert.Equal(t, "", valor)
}

func TestGetCredentialField_CacheMiss(t *testing.T) {
	// Arrange
	c, r := cacheSetup(t)
	ctx := context.Background()

	r.On("Get", mock.Anything, "integration:creds:99").Return("", errors.New("miss"))

	// Act
	valor, err := c.GetCredentialField(ctx, 99, "api_key")

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "", valor)
}

// ============================================
// InvalidateIntegration
// ============================================

func TestInvalidateIntegration_EliminaMetadataYCredenciales(t *testing.T) {
	// Arrange
	c, r := cacheSetup(t)
	ctx := context.Background()

	r.On("Delete", mock.Anything, "integration:meta:1").Return(nil)
	r.On("Delete", mock.Anything, "integration:creds:1").Return(nil)

	// Act
	err := c.InvalidateIntegration(ctx, 1)

	// Assert
	assert.NoError(t, err)
	r.AssertCalled(t, "Delete", mock.Anything, "integration:meta:1")
	r.AssertCalled(t, "Delete", mock.Anything, "integration:creds:1")
}

func TestInvalidateIntegration_ErrorNoBloquea(t *testing.T) {
	// Arrange — incluso si Delete falla, InvalidateIntegration no retorna error
	c, r := cacheSetup(t)
	ctx := context.Background()

	r.On("Delete", mock.Anything, "integration:meta:5").Return(errors.New("redis error"))
	r.On("Delete", mock.Anything, "integration:creds:5").Return(errors.New("redis error"))

	// Act
	err := c.InvalidateIntegration(ctx, 5)

	// Assert — no retorna error al caller
	assert.NoError(t, err)
}

// ============================================
// GetByCode
// ============================================

func TestGetByCode_ExisteEnCache(t *testing.T) {
	// Arrange
	c, r := cacheSetup(t)
	ctx := context.Background()

	integracion := &domain.CachedIntegration{ID: 7, Code: "shopify_main"}
	data, _ := json.Marshal(integracion)

	// Índice de código devuelve el ID "7"
	r.On("Get", mock.Anything, "integration:code:shopify_main").Return("7", nil)
	// Metadata de ID 7
	r.On("Get", mock.Anything, "integration:meta:7").Return(string(data), nil)

	// Act
	resultado, err := c.GetByCode(ctx, "shopify_main")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resultado)
	assert.Equal(t, uint(7), resultado.ID)
}

func TestGetByCode_CacheMissEnIndice(t *testing.T) {
	// Arrange
	c, r := cacheSetup(t)
	ctx := context.Background()

	r.On("Get", mock.Anything, "integration:code:unknown").Return("", errors.New("miss"))

	// Act
	resultado, err := c.GetByCode(ctx, "unknown")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resultado)
}

func TestGetByCode_IDInvalidoEnIndice(t *testing.T) {
	// Arrange
	c, r := cacheSetup(t)
	ctx := context.Background()

	// El índice devuelve un valor que no es un número
	r.On("Get", mock.Anything, "integration:code:bad_code").Return("not_a_number", nil)

	// Act
	resultado, err := c.GetByCode(ctx, "bad_code")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resultado)
}

// ============================================
// GetByBusinessAndType
// ============================================

func TestGetByBusinessAndType_ExisteEnCache(t *testing.T) {
	// Arrange
	c, r := cacheSetup(t)
	ctx := context.Background()

	bizID := uint(10)
	integracion := &domain.CachedIntegration{ID: 8, IntegrationTypeID: 1}
	data, _ := json.Marshal(integracion)

	// Índice biz+type devuelve el ID "8"
	r.On("Get", mock.Anything, "integration:idx:biz:10:type:1").Return("8", nil)
	r.On("Get", mock.Anything, "integration:meta:8").Return(string(data), nil)

	// Act
	resultado, err := c.GetByBusinessAndType(ctx, bizID, 1)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resultado)
	assert.Equal(t, uint(8), resultado.ID)
}

func TestGetByBusinessAndType_CacheMiss(t *testing.T) {
	// Arrange
	c, r := cacheSetup(t)
	ctx := context.Background()

	r.On("Get", mock.Anything, "integration:idx:biz:99:type:5").Return("", errors.New("miss"))

	// Act
	resultado, err := c.GetByBusinessAndType(ctx, 99, 5)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resultado)
}
