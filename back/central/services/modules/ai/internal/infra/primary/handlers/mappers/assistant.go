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
	return dtos.ChatInput{Scope: scope, Messages: messages}
}

func FromReply(reply *entities.AssistantReply) response.AssistantReply {
	out := response.AssistantReply{Message: reply.Message}
	if reply.Destination != nil {
		out.Destination = &response.Destination{
			Key:         reply.Destination.Key,
			Label:       reply.Destination.Label,
			Route:       reply.Destination.Route,
			Description: reply.Destination.Description,
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
