package domain

import "time"

// MessageRole representa el rol de un mensaje en la conversacion
type MessageRole string

const (
	RoleUser      MessageRole = "user"
	RoleAssistant MessageRole = "assistant"
)

// ContentType representa el tipo de bloque de contenido
type ContentType string

const (
	ContentTypeText       ContentType = "text"
	ContentTypeToolUse    ContentType = "toolUse"
	ContentTypeToolResult ContentType = "toolResult"
)

// StopReason indica por que Bedrock dejo de generar
type StopReason string

const (
	StopReasonEndTurn  StopReason = "end_turn"
	StopReasonToolUse  StopReason = "tool_use"
	StopReasonMaxToken StopReason = "max_tokens"
)

// AISession representa una sesion de chat activa en Redis
type AISession struct {
	ID          string
	PhoneNumber string
	BusinessID  uint
	Messages    []AIMessage
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ExpiresAt   time.Time
}

// AIMessage representa un mensaje individual en el historial
type AIMessage struct {
	Role    MessageRole
	Content []ContentBlock
}

// ContentBlock representa un bloque de contenido dentro de un mensaje
type ContentBlock struct {
	Type      ContentType
	Text      string
	ToolUseID string
	ToolName  string
	Input     string // JSON string para toolUse
	Content   string // Para toolResult
}

// AIResponse representa la respuesta de Bedrock
type AIResponse struct {
	Content    []ContentBlock
	StopReason StopReason
}

// ToolDefinition define una herramienta para Bedrock
type ToolDefinition struct {
	Name        string
	Description string
	InputSchema string // JSON Schema como string
}

// ProductSearchResult resultado de busqueda de producto
type ProductSearchResult struct {
	ID               string
	SKU              string
	Name             string
	Description      string
	ShortDescription string
	Price            float64
	Currency         string
	StockQuantity    int
	Category         string
	Brand            string
	ImageURL         string
	IsActive         bool
}

// AIConfig configuracion del agente de ventas IA leida de platform_creds
type AIConfig struct {
	Enabled           bool
	ModelID           string
	SessionTTLMinutes int
	MaxToolIterations int
	DemoBusinessID    uint
}
