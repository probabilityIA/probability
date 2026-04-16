package domain

import "time"

type MessageRole string

const (
	RoleUser      MessageRole = "user"
	RoleAssistant MessageRole = "assistant"
)

type ContentType string

const (
	ContentTypeText       ContentType = "text"
	ContentTypeToolUse    ContentType = "toolUse"
	ContentTypeToolResult ContentType = "toolResult"
)

type StopReason string

const (
	StopReasonEndTurn  StopReason = "end_turn"
	StopReasonToolUse  StopReason = "tool_use"
	StopReasonMaxToken StopReason = "max_tokens"
)

type AISession struct {
	ID          string
	PhoneNumber string
	BusinessID  uint
	Messages    []AIMessage
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ExpiresAt   time.Time
}

type AIMessage struct {
	Role    MessageRole
	Content []ContentBlock
}

type ContentBlock struct {
	Type      ContentType
	Text      string
	ToolUseID string
	ToolName  string
	Input     string // JSON string para toolUse
	Content   string // Para toolResult
}

type AIResponse struct {
	Content    []ContentBlock
	StopReason StopReason
}

type ToolDefinition struct {
	Name        string
	Description string
	InputSchema string // JSON Schema como string
}

type ProductSearchResult struct {
	ID               string
	SKU              string
	Name             string
	Description      string
	ShortDescription string
	Price            float64
	Currency         string
	StockQuantity    int
	TrackInventory   bool
	Category         string
	Brand            string
	ImageURL         string
	IsActive         bool
}

type CustomerSearchResult struct {
	ID    uint
	Name  string
	Email string
	Phone string
	DNI   string
}

type CustomerLastAddress struct {
	Street     string
	City       string
	State      string
	Country    string
	PostalCode string
	OrderDate  string
}

type OrderResultDTO struct {
	ExternalID   string `json:"external_id"`
	PhoneNumber  string `json:"phone_number"`
	BusinessID   uint   `json:"business_id"`
	Success      bool   `json:"success"`
	OrderID      string `json:"order_id"`
	OrderNumber  string `json:"order_number"`
	ErrorMessage string `json:"error_message"`
}

type AIConfig struct {
	Enabled           bool
	ModelID           string
	SessionTTLMinutes int
	MaxToolIterations int
	DemoBusinessID    uint
}
