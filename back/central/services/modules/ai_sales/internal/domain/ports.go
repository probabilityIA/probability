package domain

import "context"

// IAIProvider define la interfaz del proveedor de IA (adapta shared/bedrock)
type IAIProvider interface {
	Converse(ctx context.Context, messages []AIMessage, systemPrompt string, tools []ToolDefinition) (*AIResponse, error)
}

// ISessionCache define la interfaz del cache de sesiones en Redis
type ISessionCache interface {
	Get(ctx context.Context, phoneNumber string) (*AISession, error)
	Save(ctx context.Context, session *AISession) error
	Delete(ctx context.Context, phoneNumber string) error
}

// IProductRepository define la interfaz del repositorio de productos (replicado localmente)
type IProductRepository interface {
	SearchProducts(ctx context.Context, businessID uint, query string, limit int) ([]ProductSearchResult, error)
	GetProductBySKU(ctx context.Context, businessID uint, sku string) (*ProductSearchResult, error)
}

// IAIResponsePublisher publica respuestas de AI a la cola de WhatsApp
type IAIResponsePublisher interface {
	PublishResponse(ctx context.Context, phoneNumber string, businessID uint, text string) error
}

// IAIOrderPublisher publica ordenes al queue canonical de ordenes
type IAIOrderPublisher interface {
	PublishOrder(ctx context.Context, orderPayload []byte) error
}

// IConfigProvider obtiene la configuracion del AI desde platform_creds en Redis
type IConfigProvider interface {
	GetAIConfig(ctx context.Context) (*AIConfig, error)
}

// IAIPersistencePublisher persiste conversaciones y mensajes AI en BD via RabbitMQ
type IAIPersistencePublisher interface {
	PublishConversationUpsert(ctx context.Context, session *AISession) error
	PublishMessageLog(ctx context.Context, conversationID, phoneNumber, direction, content string) error
}

// IAIPauseChecker verifica si el AI está pausado para un número (humano tomó control).
// Implementado en infra/secondary/cache usando la misma clave Redis que el módulo WhatsApp.
type IAIPauseChecker interface {
	IsAIPaused(ctx context.Context, phoneNumber string) bool
}
