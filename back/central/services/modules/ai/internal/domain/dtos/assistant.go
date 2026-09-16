package dtos

import "github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/entities"

type AccessScope struct {
	UserID              uint
	TokenBusinessID     uint
	RequestedBusinessID uint
}

type ChatInput struct {
	Scope          AccessScope
	ConversationID string
	Pathname       string
	Messages       []entities.ChatMessage
}

type ToolDefinition struct {
	Name        string
	Description string
	Schema      map[string]any
}

type ToolCall struct {
	ID    string
	Name  string
	Input map[string]any
}

type ToolResult struct {
	ToolCallID string
	Content    map[string]any
}

type ModelMessage struct {
	Role        entities.MessageRole
	Text        string
	ToolCalls   []ToolCall
	ToolResults []ToolResult
}

type ModelRequest struct {
	SystemPrompt    string
	Messages        []ModelMessage
	DestinationKeys []string
	Tools           []ToolDefinition
	ForceReply      bool
}

type ModelReply struct {
	Message        string
	DestinationKey string
	ToolCalls      []ToolCall
	Model          string
	InputTokens    int
	OutputTokens   int
}
