package mocks

import (
	"context"
	"sync"

	"github.com/rs/zerolog"
	domain "github.com/secamc93/probability/back/central/services/modules/ai_sales/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

type ConverseCall struct {
	Messages     []domain.AIMessage
	SystemPrompt string
	Tools        []domain.ToolDefinition
}

type AIProviderMock struct {
	ConverseFn func(ctx context.Context, messages []domain.AIMessage, systemPrompt string, tools []domain.ToolDefinition) (*domain.AIResponse, error)

	mu    sync.Mutex
	calls []ConverseCall
}

var _ domain.IAIProvider = (*AIProviderMock)(nil)

func (m *AIProviderMock) Converse(ctx context.Context, messages []domain.AIMessage, systemPrompt string, tools []domain.ToolDefinition) (*domain.AIResponse, error) {
	m.mu.Lock()
	cp := make([]domain.AIMessage, len(messages))
	copy(cp, messages)
	m.calls = append(m.calls, ConverseCall{Messages: cp, SystemPrompt: systemPrompt, Tools: tools})
	n := len(m.calls)
	m.mu.Unlock()
	if m.ConverseFn != nil {
		return m.ConverseFn(ctx, messages, systemPrompt, tools)
	}
	_ = n
	return &domain.AIResponse{
		StopReason: domain.StopReasonEndTurn,
		Content:    []domain.ContentBlock{{Type: domain.ContentTypeText, Text: "Hola"}},
	}, nil
}

func (m *AIProviderMock) Calls() []ConverseCall {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]ConverseCall, len(m.calls))
	copy(out, m.calls)
	return out
}

type SessionCacheMock struct {
	GetFn    func(ctx context.Context, phoneNumber string) (*domain.AISession, error)
	SaveFn   func(ctx context.Context, session *domain.AISession) error
	DeleteFn func(ctx context.Context, phoneNumber string) error

	Saved []domain.AISession
}

var _ domain.ISessionCache = (*SessionCacheMock)(nil)

func (m *SessionCacheMock) Get(ctx context.Context, phoneNumber string) (*domain.AISession, error) {
	if m.GetFn != nil {
		return m.GetFn(ctx, phoneNumber)
	}
	return nil, nil
}

func (m *SessionCacheMock) Save(ctx context.Context, session *domain.AISession) error {
	m.Saved = append(m.Saved, *session)
	if m.SaveFn != nil {
		return m.SaveFn(ctx, session)
	}
	return nil
}

func (m *SessionCacheMock) Delete(ctx context.Context, phoneNumber string) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, phoneNumber)
	}
	return nil
}

type ProductRepositoryMock struct {
	SearchProductsFn  func(ctx context.Context, businessID uint, query string, limit int) ([]domain.ProductSearchResult, error)
	GetProductBySKUFn func(ctx context.Context, businessID uint, sku string) (*domain.ProductSearchResult, error)

	SearchLimits []int
	SearchQuery  []string
}

var _ domain.IProductRepository = (*ProductRepositoryMock)(nil)

func (m *ProductRepositoryMock) SearchProducts(ctx context.Context, businessID uint, query string, limit int) ([]domain.ProductSearchResult, error) {
	m.SearchLimits = append(m.SearchLimits, limit)
	m.SearchQuery = append(m.SearchQuery, query)
	if m.SearchProductsFn != nil {
		return m.SearchProductsFn(ctx, businessID, query, limit)
	}
	return nil, nil
}

func (m *ProductRepositoryMock) GetProductBySKU(ctx context.Context, businessID uint, sku string) (*domain.ProductSearchResult, error) {
	if m.GetProductBySKUFn != nil {
		return m.GetProductBySKUFn(ctx, businessID, sku)
	}
	return &domain.ProductSearchResult{ID: "P-1", SKU: sku, Name: "Producto", Price: 1000, Currency: "COP"}, nil
}

type CustomerRepositoryMock struct {
	SearchCustomersFn          func(ctx context.Context, businessID uint, query string) ([]domain.CustomerSearchResult, error)
	GetCustomerLastAddressFn   func(ctx context.Context, businessID uint, customerID uint) (*domain.CustomerLastAddress, error)
	GetWhatsAppIntegrationIDFn func(ctx context.Context, businessID uint) (uint, error)
}

var _ domain.ICustomerRepository = (*CustomerRepositoryMock)(nil)

func (m *CustomerRepositoryMock) SearchCustomers(ctx context.Context, businessID uint, query string) ([]domain.CustomerSearchResult, error) {
	if m.SearchCustomersFn != nil {
		return m.SearchCustomersFn(ctx, businessID, query)
	}
	return nil, nil
}

func (m *CustomerRepositoryMock) GetCustomerLastAddress(ctx context.Context, businessID uint, customerID uint) (*domain.CustomerLastAddress, error) {
	if m.GetCustomerLastAddressFn != nil {
		return m.GetCustomerLastAddressFn(ctx, businessID, customerID)
	}
	return nil, nil
}

func (m *CustomerRepositoryMock) GetWhatsAppIntegrationID(ctx context.Context, businessID uint) (uint, error) {
	if m.GetWhatsAppIntegrationIDFn != nil {
		return m.GetWhatsAppIntegrationIDFn(ctx, businessID)
	}
	return 9, nil
}

type ResponseCall struct {
	PhoneNumber string
	BusinessID  uint
	Text        string
}

type ResponsePublisherMock struct {
	PublishResponseFn func(ctx context.Context, phoneNumber string, businessID uint, text string) error
	Published         []ResponseCall
}

var _ domain.IAIResponsePublisher = (*ResponsePublisherMock)(nil)

func (m *ResponsePublisherMock) PublishResponse(ctx context.Context, phoneNumber string, businessID uint, text string) error {
	m.Published = append(m.Published, ResponseCall{phoneNumber, businessID, text})
	if m.PublishResponseFn != nil {
		return m.PublishResponseFn(ctx, phoneNumber, businessID, text)
	}
	return nil
}

type OrderPublisherMock struct {
	PublishOrderFn func(ctx context.Context, orderPayload []byte) error
	Published      [][]byte
}

var _ domain.IAIOrderPublisher = (*OrderPublisherMock)(nil)

func (m *OrderPublisherMock) PublishOrder(ctx context.Context, orderPayload []byte) error {
	m.Published = append(m.Published, orderPayload)
	if m.PublishOrderFn != nil {
		return m.PublishOrderFn(ctx, orderPayload)
	}
	return nil
}

type ConfigProviderMock struct {
	GetAIConfigFn func(ctx context.Context) (*domain.AIConfig, error)
}

var _ domain.IConfigProvider = (*ConfigProviderMock)(nil)

func (m *ConfigProviderMock) GetAIConfig(ctx context.Context) (*domain.AIConfig, error) {
	if m.GetAIConfigFn != nil {
		return m.GetAIConfigFn(ctx)
	}
	return &domain.AIConfig{Enabled: true, DemoBusinessID: 1, MaxToolIterations: 5, SessionTTLMinutes: 20}, nil
}

type MessageLogCall struct {
	ConversationID string
	PhoneNumber    string
	Direction      string
	Content        string
}

type PersistencePublisherMock struct {
	PublishConversationUpsertFn func(ctx context.Context, session *domain.AISession) error
	PublishMessageLogFn         func(ctx context.Context, conversationID, phoneNumber, direction, content string) error

	Upserts     []domain.AISession
	MessageLogs []MessageLogCall
}

var _ domain.IAIPersistencePublisher = (*PersistencePublisherMock)(nil)

func (m *PersistencePublisherMock) PublishConversationUpsert(ctx context.Context, session *domain.AISession) error {
	m.Upserts = append(m.Upserts, *session)
	if m.PublishConversationUpsertFn != nil {
		return m.PublishConversationUpsertFn(ctx, session)
	}
	return nil
}

func (m *PersistencePublisherMock) PublishMessageLog(ctx context.Context, conversationID, phoneNumber, direction, content string) error {
	m.MessageLogs = append(m.MessageLogs, MessageLogCall{conversationID, phoneNumber, direction, content})
	if m.PublishMessageLogFn != nil {
		return m.PublishMessageLogFn(ctx, conversationID, phoneNumber, direction, content)
	}
	return nil
}

type PauseCheckerMock struct {
	IsAIPausedFn func(ctx context.Context, phoneNumber string) bool
}

var _ domain.IAIPauseChecker = (*PauseCheckerMock)(nil)

func (m *PauseCheckerMock) IsAIPaused(ctx context.Context, phoneNumber string) bool {
	if m.IsAIPausedFn != nil {
		return m.IsAIPausedFn(ctx, phoneNumber)
	}
	return false
}

type SilentLogger struct{}

func NewSilentLogger() log.ILogger { return &SilentLogger{} }

func (l *SilentLogger) nop() zerolog.Logger { return zerolog.Nop() }

func (l *SilentLogger) Info(ctx ...context.Context) *zerolog.Event  { n := l.nop(); return n.Info() }
func (l *SilentLogger) Error(ctx ...context.Context) *zerolog.Event { n := l.nop(); return n.Error() }
func (l *SilentLogger) Warn(ctx ...context.Context) *zerolog.Event  { n := l.nop(); return n.Warn() }
func (l *SilentLogger) Debug(ctx ...context.Context) *zerolog.Event { n := l.nop(); return n.Debug() }
func (l *SilentLogger) Fatal(ctx ...context.Context) *zerolog.Event { n := l.nop(); return n.Fatal() }
func (l *SilentLogger) Panic(ctx ...context.Context) *zerolog.Event { n := l.nop(); return n.Panic() }
func (l *SilentLogger) With() zerolog.Context                       { n := l.nop(); return n.With() }
func (l *SilentLogger) WithService(service string) log.ILogger      { return l }
func (l *SilentLogger) WithModule(module string) log.ILogger        { return l }
func (l *SilentLogger) WithBusinessID(businessID uint) log.ILogger  { return l }
