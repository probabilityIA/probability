package mappers

import (
	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/ai/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/modules/ai/internal/infra/primary/handlers/response"
)

func ToChatInput(scope dtos.AccessScope, req request.AssistantChat) dtos.ChatInput {
	messages := make([]entities.ChatMessage, 0, len(req.Messages))
	for _, m := range req.Messages {
		messages = append(messages, entities.ChatMessage{Role: entities.MessageRole(m.Role), Text: m.Text})
	}
	return dtos.ChatInput{
		Scope:          scope,
		ConversationID: req.ConversationID,
		Pathname:       req.Pathname,
		Messages:       messages,
	}
}

func FromReply(reply *entities.AssistantReply) response.AssistantReply {
	out := response.AssistantReply{
		MessageID:      reply.MessageID,
		ConversationID: reply.ConversationID,
		Message:        reply.Message,
	}
	if reply.Destination != nil {
		out.Destination = &response.Destination{
			Key:         reply.Destination.Key,
			Label:       reply.Destination.Label,
			Route:       reply.Destination.Route,
			Description: reply.Destination.Description,
		}
		if h := reply.Destination.Highlight; h != nil {
			out.Destination.Highlight = &response.Highlight{Target: h.Target, Title: h.Title, Hint: h.Hint}
		}
	}
	return out
}

func FromState(state *entities.AssistantState) response.AssistantState {
	return response.AssistantState{
		IntroSeen: state.IntroSeen,
		Limit:     state.Limit,
		Remaining: state.Remaining,
		ResetAt:   state.ResetAt,
	}
}

func FromRecommendation(rec *entities.Recommendation) response.Recommendation {
	quotations := make([]response.Quotation, 0, len(rec.Quotations))
	for _, q := range rec.Quotations {
		quotations = append(quotations, response.Quotation{
			Carrier:               q.Carrier,
			EstimatedCost:         q.EstimatedCost,
			EstimatedDeliveryDays: q.EstimatedDeliveryDays,
		})
	}
	return response.Recommendation{
		RecommendedCarrier: rec.RecommendedCarrier,
		Reasoning:          rec.Reasoning,
		Alternatives:       rec.Alternatives,
		Quotations:         quotations,
	}
}

func FromReviewMessages(page *dtos.PaginatedResponse[entities.ReviewMessage]) response.ReviewMessages {
	data := make([]response.ReviewMessage, 0, len(page.Data))
	for _, m := range page.Data {
		data = append(data, response.ReviewMessage{
			ID:               m.ID,
			ConversationID:   m.ConversationID,
			BusinessID:       m.BusinessID,
			BusinessName:     m.BusinessName,
			UserID:           m.UserID,
			UserName:         m.UserName,
			UserEmail:        m.UserEmail,
			Pathname:         m.Pathname,
			Question:         m.Question,
			Answer:           m.Answer,
			DestinationKey:   m.DestinationKey,
			DestinationRoute: m.DestinationRoute,
			ErrorCode:        m.ErrorCode,
			Model:            m.Model,
			InputTokens:      m.InputTokens,
			OutputTokens:     m.OutputTokens,
			LatencyMs:        m.LatencyMs,
			Feedback:         m.Feedback,
			FeedbackAt:       m.FeedbackAt,
			ClickedAt:        m.ClickedAt,
			CreatedAt:        m.CreatedAt,
		})
	}
	return response.ReviewMessages{
		Data:       data,
		Total:      page.Total,
		Page:       page.Page,
		PageSize:   page.PageSize,
		TotalPages: page.TotalPages,
	}
}

func FromReviewSummary(summary *entities.ReviewSummary) response.ReviewSummary {
	destinations := make([]response.DestinationCount, 0, len(summary.TopDestinations))
	for _, d := range summary.TopDestinations {
		destinations = append(destinations, response.DestinationCount{Key: d.Key, Count: d.Count})
	}
	return response.ReviewSummary{
		Messages:         summary.Messages,
		Conversations:    summary.Conversations,
		Users:            summary.Users,
		InputTokens:      summary.InputTokens,
		OutputTokens:     summary.OutputTokens,
		NoDestination:    summary.NoDestination,
		Errors:           summary.Errors,
		Positive:         summary.Positive,
		Negative:         summary.Negative,
		WithDestination:  summary.WithDestination,
		Clicked:          summary.Clicked,
		EstimatedCostUSD: summary.EstimatedCostUSD,
		TopDestinations:  destinations,
	}
}
