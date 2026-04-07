package domain

import "context"

type IAIProvider interface {
	Converse(ctx context.Context, messages []AIMessage, systemPrompt string, tools []ToolDefinition) (*AIResponse, error)
}

type ISessionCache interface {
	Get(ctx context.Context, phoneNumber string) (*AISession, error)
	Save(ctx context.Context, session *AISession) error
	Delete(ctx context.Context, phoneNumber string) error
}

type IProductRepository interface {
	SearchProducts(ctx context.Context, businessID uint, query string, limit int) ([]ProductSearchResult, error)
	GetProductBySKU(ctx context.Context, businessID uint, sku string) (*ProductSearchResult, error)
}

type ICustomerRepository interface {
	SearchCustomers(ctx context.Context, businessID uint, query string) ([]CustomerSearchResult, error)
	GetCustomerLastAddress(ctx context.Context, businessID uint, customerID uint) (*CustomerLastAddress, error)
	GetWhatsAppIntegrationID(ctx context.Context, businessID uint) (uint, error)
}

type IAIResponsePublisher interface {
	PublishResponse(ctx context.Context, phoneNumber string, businessID uint, text string) error
}

type IAIOrderPublisher interface {
	PublishOrder(ctx context.Context, orderPayload []byte) error
}

type IConfigProvider interface {
	GetAIConfig(ctx context.Context) (*AIConfig, error)
}

type IAIPersistencePublisher interface {
	PublishConversationUpsert(ctx context.Context, session *AISession) error
	PublishMessageLog(ctx context.Context, conversationID, phoneNumber, direction, content string) error
}

type IAIPauseChecker interface {
	IsAIPaused(ctx context.Context, phoneNumber string) bool
}
