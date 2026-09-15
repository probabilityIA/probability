package app

import (
	"context"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
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

	started := time.Now()
	record := entities.MessageRecord{
		ID:             uuid.NewString(),
		ConversationID: conversationIDOrNew(input.ConversationID),
		BusinessID:     businessOf(input.Scope),
		UserID:         input.Scope.UserID,
		Pathname:       normalizePathname(input.Pathname),
		Question:       messages[len(messages)-1].Text,
		CreatedAt:      started,
	}

	catalog, err := uc.navigation.ForUser(ctx, input.Scope)
	if err != nil {
		return nil, err
	}

	if err := uc.consumeQuota(ctx, input.Scope.UserID); err != nil {
		record.ErrorCode = entities.ErrorCodeRateLimited
		uc.recordMessage(ctx, record, started)
		return nil, err
	}

	access := uc.resolveDataAccess(catalog, input.Scope)
	request := dtos.ModelRequest{
		SystemPrompt:    buildSystemPrompt(catalog, access, uc.now()),
		Messages:        toModelMessages(messages),
		DestinationKeys: catalog.Keys(),
		Tools:           access.Tools,
	}

	var reply *dtos.ModelReply
	for round := 0; ; round++ {
		request.ForceReply = len(request.Tools) == 0 || round >= MaxToolRounds
		reply, err = uc.model.Reply(ctx, request)
		if err != nil {
			uc.log.Error(ctx).Err(err).Uint("user_id", input.Scope.UserID).Int("round", round).Msg("[ai.assistant] el modelo no respondio")
			record.ErrorCode = entities.ErrorCodeUnavailable
			uc.recordMessage(ctx, record, started)
			return nil, domainerrors.ErrModelUnavailable
		}

		record.Model = reply.Model
		record.InputTokens += reply.InputTokens
		record.OutputTokens += reply.OutputTokens

		if request.ForceReply || len(reply.ToolCalls) == 0 {
			break
		}
		request.Messages = append(request.Messages,
			dtos.ModelMessage{Role: entities.RoleAssistant, Text: reply.Message, ToolCalls: reply.ToolCalls},
			dtos.ModelMessage{Role: entities.RoleUser, ToolResults: uc.runTools(ctx, access, reply.ToolCalls)},
		)
	}

	result, err := composeReply(reply, catalog)
	if err != nil {
		record.ErrorCode = entities.ErrorCodeUnavailable
		uc.recordMessage(ctx, record, started)
		return nil, err
	}

	record.Answer = result.Message
	if result.Destination != nil {
		record.DestinationKey = result.Destination.Key
		record.DestinationRoute = result.Destination.Route
	}
	uc.recordMessage(ctx, record, started)

	result.MessageID = record.ID
	result.ConversationID = record.ConversationID
	return result, nil
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

func (uc *UseCase) recordMessage(ctx context.Context, record entities.MessageRecord, started time.Time) {
	record.LatencyMs = int(time.Since(started).Milliseconds())

	var err error
	switch {
	case uc.recorder != nil:
		err = uc.recorder.RecordMessage(ctx, record)
	case uc.conversations != nil:
		err = uc.conversations.SaveMessage(ctx, record)
	default:
		return
	}
	if err != nil {
		uc.log.Warn(ctx).Err(err).Str("message_id", record.ID).Msg("[ai.assistant] no se pudo guardar el mensaje")
	}
}

func conversationIDOrNew(raw string) string {
	if parsed, err := uuid.Parse(strings.TrimSpace(raw)); err == nil {
		return parsed.String()
	}
	return uuid.NewString()
}

func businessOf(scope dtos.AccessScope) *uint {
	if scope.TokenBusinessID > 0 {
		id := scope.TokenBusinessID
		return &id
	}
	if scope.RequestedBusinessID > 0 {
		id := scope.RequestedBusinessID
		return &id
	}
	return nil
}

func toModelMessages(messages []entities.ChatMessage) []dtos.ModelMessage {
	out := make([]dtos.ModelMessage, 0, len(messages))
	for _, m := range messages {
		out = append(out, dtos.ModelMessage{Role: m.Role, Text: m.Text})
	}
	return out
}

func normalizePathname(raw string) string {
	path := strings.TrimSpace(raw)
	if !strings.HasPrefix(path, "/") {
		return ""
	}
	return truncateRunes(path, MaxPathnameRunes)
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

var inlineDestinationPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?is)\{\s*"?destination"?\s*:\s*"([a-z_.]+)"\s*\}`),
	regexp.MustCompile(`(?im)^[\s*_]*destination[\s*_]*[:=][\s*_]*([a-z_.]+)[\s*_.]*$`),
}

func extractInlineDestination(text string) (string, string) {
	for _, pattern := range inlineDestinationPatterns {
		match := pattern.FindStringSubmatchIndex(text)
		if match == nil {
			continue
		}
		key := text[match[2]:match[3]]
		cleaned := strings.TrimSpace(text[:match[0]] + text[match[1]:])
		return cleaned, key
	}
	return text, ""
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
