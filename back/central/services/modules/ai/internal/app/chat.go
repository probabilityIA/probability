package app

import (
	"context"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/errors"
)

const noDestination = "none"

func (uc *UseCase) Chat(ctx context.Context, input dtos.ChatInput) (*entities.AssistantReply, error) {
	messages, err := normalizeConversation(input.Messages)
	if err != nil {
		return nil, err
	}

	catalog, err := uc.navigation.ForUser(ctx, input.Scope)
	if err != nil {
		return nil, err
	}

	if err := uc.consumeQuota(ctx, input.Scope.UserID); err != nil {
		return nil, err
	}

	reply, err := uc.model.Reply(ctx, dtos.ModelRequest{
		SystemPrompt:    buildSystemPrompt(catalog),
		Messages:        messages,
		DestinationKeys: catalog.Keys(),
	})
	if err != nil {
		uc.log.Error(ctx).Err(err).Uint("user_id", input.Scope.UserID).Msg("[ai.assistant] el modelo no respondio")
		return nil, domainerrors.ErrModelUnavailable
	}

	return composeReply(reply, catalog)
}

func (uc *UseCase) consumeQuota(ctx context.Context, userID uint) error {
	if uc.store == nil {
		return nil
	}
	usage, err := uc.store.ConsumeMessage(ctx, userID, AssistantWindow)
	if err != nil {
		uc.log.Warn(ctx).Err(err).Uint("user_id", userID).Msg("[ai.assistant] no se pudo registrar el consumo; se atiende sin limite")
		return nil
	}
	if usage.Count > AssistantMessageLimit {
		return &domainerrors.RateLimitedError{Limit: AssistantMessageLimit, ResetAt: usage.ResetAt}
	}
	return nil
}

func normalizeConversation(messages []entities.ChatMessage) ([]entities.ChatMessage, error) {
	cleaned := make([]entities.ChatMessage, 0, len(messages))
	for _, m := range messages {
		if m.Role != entities.RoleUser && m.Role != entities.RoleAssistant {
			continue
		}
		text := strings.TrimSpace(m.Text)
		if text == "" {
			continue
		}
		if m.Role == entities.RoleUser && utf8.RuneCountInString(text) > MaxUserMessageRunes {
			return nil, domainerrors.ErrMessageTooLong
		}
		if m.Role == entities.RoleAssistant {
			text = truncateRunes(text, MaxAssistantRunes)
		}
		cleaned = append(cleaned, entities.ChatMessage{Role: m.Role, Text: text})
	}

	if len(cleaned) > MaxHistoryMessages {
		cleaned = cleaned[len(cleaned)-MaxHistoryMessages:]
	}
	for len(cleaned) > 0 && cleaned[0].Role != entities.RoleUser {
		cleaned = cleaned[1:]
	}

	merged := make([]entities.ChatMessage, 0, len(cleaned))
	for _, m := range cleaned {
		if n := len(merged); n > 0 && merged[n-1].Role == m.Role {
			merged[n-1].Text += "\n\n" + m.Text
			continue
		}
		merged = append(merged, m)
	}

	if len(merged) == 0 || merged[len(merged)-1].Role != entities.RoleUser {
		return nil, domainerrors.ErrEmptyConversation
	}
	return merged, nil
}

func composeReply(reply *dtos.ModelReply, catalog *entities.NavigationCatalog) (*entities.AssistantReply, error) {
	if reply == nil {
		return nil, domainerrors.ErrModelUnavailable
	}

	message, inlineKey := extractInlineDestination(cleanModelText(reply.Message))

	var destination *entities.Destination
	key := strings.TrimSpace(reply.DestinationKey)
	if key == "" || key == noDestination {
		key = inlineKey
	}
	if key != "" && key != noDestination {
		if d, ok := catalog.Find(key); ok {
			destination = d
		}
	}
	if message == "" {
		if destination == nil {
			return nil, domainerrors.ErrModelUnavailable
		}
		message = "Esto lo encuentras en " + destination.Label + "."
	}

	return &entities.AssistantReply{Message: message, Destination: destination}, nil
}

var inlineDestinationPattern = regexp.MustCompile(`(?im)^[\s*_]*destination[\s*_]*[:=][\s*_]*([a-z_.]+)[\s*_.]*$`)

func extractInlineDestination(text string) (string, string) {
	match := inlineDestinationPattern.FindStringSubmatchIndex(text)
	if match == nil {
		return text, ""
	}
	key := text[match[2]:match[3]]
	cleaned := strings.TrimSpace(text[:match[0]] + text[match[1]:])
	return cleaned, key
}

func cleanModelText(text string) string {
	for {
		start := strings.Index(text, "<think>")
		if start < 0 {
			break
		}
		end := strings.Index(text[start:], "</think>")
		if end < 0 {
			text = text[:start]
			break
		}
		text = text[:start] + text[start+end+len("</think>"):]
	}
	return strings.TrimSpace(text)
}

func truncateRunes(text string, max int) string {
	if utf8.RuneCountInString(text) <= max {
		return text
	}
	return string([]rune(text)[:max])
}
