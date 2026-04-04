package app

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"

	domain "github.com/secamc93/probability/back/central/services/modules/ai_sales/internal/domain"
)

var thinkingTagRegex = regexp.MustCompile(`(?s)<thinking>.*?</thinking>`)

const (
	maxHistoryMessages = 20
	defaultBusinessName = "Probability Demo"
)

func (uc *useCase) HandleIncoming(ctx context.Context, dto domain.IncomingMessageDTO) error {
	// Verificar si un humano tomó control de esta conversación
	if uc.pauseChecker != nil && uc.pauseChecker.IsAIPaused(ctx, dto.PhoneNumber) {
		uc.log.Info(ctx).
			Str("phone", dto.PhoneNumber).
			Msg("[AI Sales] - AI pausado por atención humana, ignorando mensaje")
		return nil
	}

	// 1. Obtener config
	config, err := uc.configProvider.GetAIConfig(ctx)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Error obteniendo AI config")
		return uc.sendErrorResponse(ctx, dto.PhoneNumber, config, "Lo siento, no puedo procesar tu mensaje en este momento. Intenta de nuevo mas tarde.")
	}

	if !config.Enabled {
		uc.log.Info(ctx).Msg("AI Sales deshabilitado, ignorando mensaje")
		return nil
	}

	businessID := config.DemoBusinessID
	if dto.BusinessID > 0 {
		businessID = dto.BusinessID
	}

	uc.log.Info(ctx).
		Str("phone", dto.PhoneNumber).
		Uint("business_id", businessID).
		Msg("Procesando mensaje AI incoming")

	// 2. Get/create session
	isNewSession := false
	session, err := uc.sessionCache.Get(ctx, dto.PhoneNumber)
	if err != nil || session == nil {
		session = &domain.AISession{
			ID:          uuid.New().String(),
			PhoneNumber: dto.PhoneNumber,
			BusinessID:  businessID,
			Messages:    []domain.AIMessage{},
			CreatedAt:   time.Now(),
		}
		isNewSession = true
	}

	// 3. Append user message
	session.Messages = append(session.Messages, domain.AIMessage{
		Role: domain.RoleUser,
		Content: []domain.ContentBlock{
			{Type: domain.ContentTypeText, Text: dto.MessageText},
		},
	})

	// Persistir mensaje entrante y conversación
	uc.persistMessage(ctx, session, "inbound", dto.MessageText, isNewSession)

	// 4. Trim history
	if len(session.Messages) > maxHistoryMessages {
		session.Messages = session.Messages[len(session.Messages)-maxHistoryMessages:]
	}

	// 5. Tool use loop
	systemPrompt := BuildSystemPrompt(defaultBusinessName)
	tools := GetToolDefinitions()
	maxIterations := config.MaxToolIterations
	if maxIterations <= 0 {
		maxIterations = 5
	}

	deps := &toolDeps{
		productRepo:    uc.productRepo,
		orderPublisher: uc.orderPublisher,
		businessID:     businessID,
	}

	for i := 0; i < maxIterations; i++ {
		resp, err := uc.aiProvider.Converse(ctx, session.Messages, systemPrompt, tools)
		if err != nil {
			uc.log.Error(ctx).Err(err).Msg("Error en Bedrock Converse")
			return uc.sendErrorResponse(ctx, dto.PhoneNumber, config, "Lo siento, tuve un problema procesando tu mensaje. Intenta de nuevo.")
		}

		// Append assistant response to history
		session.Messages = append(session.Messages, domain.AIMessage{
			Role:    domain.RoleAssistant,
			Content: resp.Content,
		})

		// Si end_turn, extraer texto y responder
		if resp.StopReason == domain.StopReasonEndTurn || resp.StopReason == domain.StopReasonMaxToken {
			responseText := extractTextFromContent(resp.Content)
			if responseText == "" {
				responseText = "No pude generar una respuesta. Intenta reformular tu pregunta."
			}

			// Save session y publicar respuesta
			uc.saveSession(ctx, session, config)
			uc.persistMessage(ctx, session, "outbound", responseText, false)
			return uc.responsePublisher.PublishResponse(ctx, dto.PhoneNumber, businessID, responseText)
		}

		// Si tool_use, ejecutar tools y continuar loop
		if resp.StopReason == domain.StopReasonToolUse {
			toolResults := uc.executeTools(ctx, resp.Content, deps)
			session.Messages = append(session.Messages, domain.AIMessage{
				Role:    domain.RoleUser,
				Content: toolResults,
			})
			continue
		}
	}

	// Max iterations reached
	uc.log.Warn(ctx).Int("max_iterations", maxIterations).Msg("Max tool iterations reached")
	uc.saveSession(ctx, session, config)
	return uc.responsePublisher.PublishResponse(ctx, dto.PhoneNumber, businessID,
		"Lo siento, no pude completar tu solicitud. Intenta reformular tu pregunta de forma mas simple.")
}

func (uc *useCase) executeTools(ctx context.Context, content []domain.ContentBlock, deps *toolDeps) []domain.ContentBlock {
	var results []domain.ContentBlock

	for _, block := range content {
		if block.Type != domain.ContentTypeToolUse {
			continue
		}

		uc.log.Info(ctx).
			Str("tool", block.ToolName).
			Str("tool_use_id", block.ToolUseID).
			Msg("Ejecutando tool")

		result, err := DispatchTool(ctx, block.ToolName, block.Input, deps)
		if err != nil {
			uc.log.Error(ctx).Err(err).Str("tool", block.ToolName).Msg("Error ejecutando tool")
			result = fmt.Sprintf(`{"error": "Error ejecutando %s: %s"}`, block.ToolName, err.Error())
		}

		results = append(results, domain.ContentBlock{
			Type:      domain.ContentTypeToolResult,
			ToolUseID: block.ToolUseID,
			Content:   result,
		})
	}

	return results
}

func extractTextFromContent(content []domain.ContentBlock) string {
	for _, block := range content {
		if block.Type == domain.ContentTypeText && block.Text != "" {
			// Strip <thinking>...</thinking> blocks that some models include in output
			clean := thinkingTagRegex.ReplaceAllString(block.Text, "")
			clean = strings.TrimSpace(clean)
			if clean != "" {
				return clean
			}
		}
	}
	return ""
}

func (uc *useCase) saveSession(ctx context.Context, session *domain.AISession, config *domain.AIConfig) {
	ttl := config.SessionTTLMinutes
	if ttl <= 0 {
		ttl = 20
	}
	session.UpdatedAt = time.Now()
	session.ExpiresAt = time.Now().Add(time.Duration(ttl) * time.Minute)

	if err := uc.sessionCache.Save(ctx, session); err != nil {
		uc.log.Error(ctx).Err(err).Msg("Error guardando session en Redis")
	}
}

func (uc *useCase) sendErrorResponse(ctx context.Context, phone string, config *domain.AIConfig, msg string) error {
	businessID := uint(1)
	if config != nil && config.DemoBusinessID > 0 {
		businessID = config.DemoBusinessID
	}
	return uc.responsePublisher.PublishResponse(ctx, phone, businessID, msg)
}

func (uc *useCase) persistMessage(ctx context.Context, session *domain.AISession, direction, content string, isNewSession bool) {
	if uc.persistencePublisher == nil {
		return
	}
	if isNewSession {
		if err := uc.persistencePublisher.PublishConversationUpsert(ctx, session); err != nil {
			uc.log.Error(ctx).Err(err).Msg("Error persistiendo conversación AI")
		}
	}
	if err := uc.persistencePublisher.PublishMessageLog(ctx, session.ID, session.PhoneNumber, direction, content); err != nil {
		uc.log.Error(ctx).Err(err).Str("direction", direction).Msg("Error persistiendo mensaje AI")
	}
}
