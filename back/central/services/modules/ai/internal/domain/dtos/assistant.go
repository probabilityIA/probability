package dtos

import "github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/entities"

type AccessScope struct {
	UserID              uint
	TokenBusinessID     uint
	RequestedBusinessID uint
}

type ChatInput struct {
	Scope    AccessScope
	Messages []entities.ChatMessage
}

type ModelRequest struct {
	SystemPrompt    string
	Messages        []entities.ChatMessage
	DestinationKeys []string
}

type ModelReply struct {
	Message        string
	DestinationKey string
}
