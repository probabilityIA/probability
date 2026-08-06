package testkit

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"regexp"
	"sync"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
	redisclient "github.com/secamc93/probability/back/central/shared/redis"
)

type SilentLogger struct{}

func NewSilentLogger() log.ILogger {
	return &SilentLogger{}
}

func (l *SilentLogger) nop() zerolog.Logger {
	return zerolog.Nop()
}

func (l *SilentLogger) Info(ctx ...context.Context) *zerolog.Event {
	n := l.nop()
	return n.Info()
}

func (l *SilentLogger) Error(ctx ...context.Context) *zerolog.Event {
	n := l.nop()
	return n.Error()
}

func (l *SilentLogger) Warn(ctx ...context.Context) *zerolog.Event {
	n := l.nop()
	return n.Warn()
}

func (l *SilentLogger) Debug(ctx ...context.Context) *zerolog.Event {
	n := l.nop()
	return n.Debug()
}

func (l *SilentLogger) Fatal(ctx ...context.Context) *zerolog.Event {
	n := l.nop()
	return n.Fatal()
}

func (l *SilentLogger) Panic(ctx ...context.Context) *zerolog.Event {
	n := l.nop()
	return n.Panic()
}

func (l *SilentLogger) With() zerolog.Context {
	n := l.nop()
	return n.With()
}

func (l *SilentLogger) WithService(service string) log.ILogger { return l }

func (l *SilentLogger) WithModule(module string) log.ILogger { return l }

func (l *SilentLogger) WithBusinessID(businessID uint) log.ILogger { return l }

type ConfigMock struct {
	Values map[string]string
}

func NewConfig(pares ...string) *ConfigMock {
	valores := map[string]string{}
	for i := 0; i+1 < len(pares); i += 2 {
		valores[pares[i]] = pares[i+1]
	}
	return &ConfigMock{Values: valores}
}

func (m *ConfigMock) Get(key string) string { return m.Values[key] }

type QueueMock struct {
	mu     sync.Mutex
	espera chan struct{}

	PublishFn           func(ctx context.Context, queueName string, message []byte) error
	PublishToExchangeFn func(ctx context.Context, exchangeName, routingKey string, message []byte) error
	DeclareQueueFn      func(queueName string, durable bool) error
	DeclareExchangeFn   func(exchangeName, exchangeType string, durable bool) error
	BindQueueFn         func(queueName, exchangeName, routingKey string) error

	Publicados  []Publicacion
	Declaradas  []Declaracion
	Exchanges   []DeclaracionExchange
	Bindings    []Binding
	Consumos    []string
	PingLlamado bool
	Cerrado     bool
}

type Publicacion struct {
	Queue      string
	Exchange   string
	RoutingKey string
	Body       []byte
}

type Declaracion struct {
	Name    string
	Durable bool
}

type DeclaracionExchange struct {
	Name    string
	Type    string
	Durable bool
}

type Binding struct {
	Queue      string
	Exchange   string
	RoutingKey string
}

func (m *QueueMock) registrar(p Publicacion) {
	m.mu.Lock()
	m.Publicados = append(m.Publicados, p)
	if m.espera != nil {
		select {
		case m.espera <- struct{}{}:
		default:
		}
	}
	m.mu.Unlock()
}

func (m *QueueMock) Publicaciones() []Publicacion {
	m.mu.Lock()
	defer m.mu.Unlock()
	copia := make([]Publicacion, len(m.Publicados))
	copy(copia, m.Publicados)
	return copia
}

func (m *QueueMock) EsperarPublicacion(t *testing.T) Publicacion {
	t.Helper()
	m.mu.Lock()
	if m.espera == nil {
		m.espera = make(chan struct{}, 64)
	}
	yaHay := len(m.Publicados)
	espera := m.espera
	m.mu.Unlock()

	if yaHay > 0 {
		return m.Publicaciones()[0]
	}
	select {
	case <-espera:
		pubs := m.Publicaciones()
		if len(pubs) == 0 {
			t.Fatal("se desperto la espera pero no hay publicaciones")
		}
		return pubs[0]
	case <-time.After(3 * time.Second):
		t.Fatal("no se publico nada en 3 segundos")
		return Publicacion{}
	}
}

func (m *QueueMock) Publish(ctx context.Context, queueName string, message []byte) error {
	m.registrar(Publicacion{Queue: queueName, Body: message})
	if m.PublishFn != nil {
		return m.PublishFn(ctx, queueName, message)
	}
	return nil
}

func (m *QueueMock) PublishToExchange(ctx context.Context, exchangeName, routingKey string, message []byte) error {
	m.registrar(Publicacion{Exchange: exchangeName, RoutingKey: routingKey, Body: message})
	if m.PublishToExchangeFn != nil {
		return m.PublishToExchangeFn(ctx, exchangeName, routingKey, message)
	}
	return nil
}

func (m *QueueMock) Consume(ctx context.Context, queueName string, handler func([]byte) error) error {
	m.Consumos = append(m.Consumos, queueName)
	return nil
}

func (m *QueueMock) ConsumeConcurrent(ctx context.Context, queueName string, handler func([]byte) error, workers int) error {
	m.Consumos = append(m.Consumos, queueName)
	return nil
}

func (m *QueueMock) Close() error {
	m.Cerrado = true
	return nil
}

func (m *QueueMock) DeclareQueue(queueName string, durable bool) error {
	m.Declaradas = append(m.Declaradas, Declaracion{Name: queueName, Durable: durable})
	if m.DeclareQueueFn != nil {
		return m.DeclareQueueFn(queueName, durable)
	}
	return nil
}

func (m *QueueMock) DeclareExchange(exchangeName, exchangeType string, durable bool) error {
	m.Exchanges = append(m.Exchanges, DeclaracionExchange{Name: exchangeName, Type: exchangeType, Durable: durable})
	if m.DeclareExchangeFn != nil {
		return m.DeclareExchangeFn(exchangeName, exchangeType, durable)
	}
	return nil
}

func (m *QueueMock) BindQueue(queueName, exchangeName, routingKey string) error {
	m.Bindings = append(m.Bindings, Binding{Queue: queueName, Exchange: exchangeName, RoutingKey: routingKey})
	if m.BindQueueFn != nil {
		return m.BindQueueFn(queueName, exchangeName, routingKey)
	}
	return nil
}

func (m *QueueMock) Ping() error {
	m.PingLlamado = true
	return nil
}

type DBMock struct {
	gorm *gorm.DB
	Mock sqlmock.Sqlmock
}

func (d *DBMock) Connect(ctx context.Context) error { return nil }

func (d *DBMock) Close() error { return nil }

func (d *DBMock) Conn(ctx context.Context) *gorm.DB { return d.gorm.WithContext(ctx) }

func (d *DBMock) WithContext(ctx context.Context) *gorm.DB { return d.gorm.WithContext(ctx) }

func (d *DBMock) DebugConn(ctx context.Context) *gorm.DB { return d.gorm.WithContext(ctx) }

func NewDB(t *testing.T) *DBMock {
	t.Helper()

	sqlDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if err != nil {
		t.Fatalf("no se pudo crear sqlmock: %v", err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn:                 sqlDB,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger:                 gormlogger.Discard,
		SkipDefaultTransaction: true,
	})
	if err != nil {
		t.Fatalf("no se pudo abrir gorm sobre sqlmock: %v", err)
	}

	return &DBMock{gorm: gormDB, Mock: mock}
}

func (d *DBMock) SinPendientes(t *testing.T) {
	t.Helper()
	if err := d.Mock.ExpectationsWereMet(); err != nil {
		t.Errorf("quedaron consultas esperadas sin ejecutar: %v", err)
	}
}

var _ db.IDatabase = (*DBMock)(nil)

func Rx(query string) string {
	return regexp.QuoteMeta(query)
}

var _ rabbitmq.IQueue = (*QueueMock)(nil)

func Peticion(t *testing.T, metodo, ruta string, cuerpo []byte) (*gin.Context, *httptest.ResponseRecorder) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)

	var lector io.Reader
	if cuerpo != nil {
		lector = bytes.NewReader(cuerpo)
	}
	req := httptest.NewRequest(metodo, ruta, lector)
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	return c, rec
}

func CuerpoJSON(t *testing.T, rec *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var out map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("la respuesta no es JSON valido (%q): %v", rec.Body.String(), err)
	}
	return out
}

type RedisMock struct {
	mu     sync.Mutex
	Values map[string]string

	GetFn func(ctx context.Context, key string) (string, error)
	SetFn func(ctx context.Context, key string, value interface{}, expiration time.Duration) error

	Leidas   []string
	Escritas []string
}

var _ redisclient.IRedis = (*RedisMock)(nil)

func NewRedis(pares ...string) *RedisMock {
	valores := map[string]string{}
	for i := 0; i+1 < len(pares); i += 2 {
		valores[pares[i]] = pares[i+1]
	}
	return &RedisMock{Values: valores}
}

func (m *RedisMock) Connect(ctx context.Context) error { return nil }

func (m *RedisMock) Close() error { return nil }

func (m *RedisMock) Client(ctx context.Context) *redis.Client { return nil }

func (m *RedisMock) Ping(ctx context.Context) error { return nil }

func (m *RedisMock) Get(ctx context.Context, key string) (string, error) {
	m.mu.Lock()
	m.Leidas = append(m.Leidas, key)
	valor, ok := m.Values[key]
	m.mu.Unlock()
	if m.GetFn != nil {
		return m.GetFn(ctx, key)
	}
	if !ok {
		return "", errors.New("redis: nil")
	}
	return valor, nil
}

func (m *RedisMock) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	m.mu.Lock()
	m.Escritas = append(m.Escritas, key)
	if s, ok := value.(string); ok {
		if m.Values == nil {
			m.Values = map[string]string{}
		}
		m.Values[key] = s
	}
	m.mu.Unlock()
	if m.SetFn != nil {
		return m.SetFn(ctx, key, value, expiration)
	}
	return nil
}

func (m *RedisMock) Delete(ctx context.Context, keys ...string) error { return nil }

func (m *RedisMock) Exists(ctx context.Context, keys ...string) (int64, error) { return 0, nil }

func (m *RedisMock) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return nil
}

func (m *RedisMock) TTL(ctx context.Context, key string) (time.Duration, error) { return 0, nil }

func (m *RedisMock) Keys(ctx context.Context, pattern string) ([]string, error) { return nil, nil }

func (m *RedisMock) Incr(ctx context.Context, key string) (int64, error) { return 0, nil }

func (m *RedisMock) Decr(ctx context.Context, key string) (int64, error) { return 0, nil }

func (m *RedisMock) HGet(ctx context.Context, key, field string) (string, error) { return "", nil }

func (m *RedisMock) HSet(ctx context.Context, key string, values ...interface{}) error { return nil }

func (m *RedisMock) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return nil, nil
}

func (m *RedisMock) HDel(ctx context.Context, key string, fields ...string) error { return nil }

func (m *RedisMock) RegisterCachePrefix(prefix string) {}

func (m *RedisMock) RegisterChannel(channel string) {}
