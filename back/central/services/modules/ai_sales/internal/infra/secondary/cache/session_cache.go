package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	domain "github.com/secamc93/probability/back/central/services/modules/ai_sales/internal/domain"
)

const (
	sessionPrefix    = "ai_sales:session:"
	defaultSessionTTL = 20 * time.Minute
)

// cachedSession estructura con JSON tags para serializar en Redis
type cachedSession struct {
	ID          string           `json:"id"`
	PhoneNumber string           `json:"phone_number"`
	BusinessID  uint             `json:"business_id"`
	Messages    []cachedMessage  `json:"messages"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	ExpiresAt   time.Time        `json:"expires_at"`
}

type cachedMessage struct {
	Role    string               `json:"role"`
	Content []cachedContentBlock `json:"content"`
}

type cachedContentBlock struct {
	Type      string `json:"type"`
	Text      string `json:"text,omitempty"`
	ToolUseID string `json:"tool_use_id,omitempty"`
	ToolName  string `json:"tool_name,omitempty"`
	Input     string `json:"input,omitempty"`
	Content   string `json:"content,omitempty"`
}

func sessionKey(phoneNumber string) string {
	return fmt.Sprintf("%s%s", sessionPrefix, phoneNumber)
}

func (c *sessionCache) Get(ctx context.Context, phoneNumber string) (*domain.AISession, error) {
	data, err := c.redis.Get(ctx, sessionKey(phoneNumber))
	if err != nil || data == "" {
		return nil, &domain.ErrSessionNotFound{PhoneNumber: phoneNumber}
	}

	var cached cachedSession
	if err := json.Unmarshal([]byte(data), &cached); err != nil {
		return nil, fmt.Errorf("error deserializing session: %w", err)
	}

	return toDomainSession(&cached), nil
}

func (c *sessionCache) Save(ctx context.Context, session *domain.AISession) error {
	cached := toCachedSession(session)
	data, err := json.Marshal(cached)
	if err != nil {
		return fmt.Errorf("error serializing session: %w", err)
	}

	ttl := time.Until(session.ExpiresAt)
	if ttl <= 0 {
		ttl = defaultSessionTTL
	}

	return c.redis.Set(ctx, sessionKey(session.PhoneNumber), string(data), ttl)
}

func (c *sessionCache) Delete(ctx context.Context, phoneNumber string) error {
	return c.redis.Delete(ctx, sessionKey(phoneNumber))
}

func toDomainSession(cached *cachedSession) *domain.AISession {
	session := &domain.AISession{
		ID:          cached.ID,
		PhoneNumber: cached.PhoneNumber,
		BusinessID:  cached.BusinessID,
		CreatedAt:   cached.CreatedAt,
		UpdatedAt:   cached.UpdatedAt,
		ExpiresAt:   cached.ExpiresAt,
	}

	for _, msg := range cached.Messages {
		domainMsg := domain.AIMessage{
			Role: domain.MessageRole(msg.Role),
		}
		for _, block := range msg.Content {
			domainMsg.Content = append(domainMsg.Content, domain.ContentBlock{
				Type:      domain.ContentType(block.Type),
				Text:      block.Text,
				ToolUseID: block.ToolUseID,
				ToolName:  block.ToolName,
				Input:     block.Input,
				Content:   block.Content,
			})
		}
		session.Messages = append(session.Messages, domainMsg)
	}

	return session
}

func toCachedSession(session *domain.AISession) *cachedSession {
	cached := &cachedSession{
		ID:          session.ID,
		PhoneNumber: session.PhoneNumber,
		BusinessID:  session.BusinessID,
		CreatedAt:   session.CreatedAt,
		UpdatedAt:   session.UpdatedAt,
		ExpiresAt:   session.ExpiresAt,
	}

	for _, msg := range session.Messages {
		cachedMsg := cachedMessage{
			Role: string(msg.Role),
		}
		for _, block := range msg.Content {
			cachedMsg.Content = append(cachedMsg.Content, cachedContentBlock{
				Type:      string(block.Type),
				Text:      block.Text,
				ToolUseID: block.ToolUseID,
				ToolName:  block.ToolName,
				Input:     block.Input,
				Content:   block.Content,
			})
		}
		cached.Messages = append(cached.Messages, cachedMsg)
	}

	return cached
}
