package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/entities"
)

func (uc *UseCase) GetAssistantState(ctx context.Context, userID uint) (*entities.AssistantState, error) {
	state := &entities.AssistantState{Limit: AssistantMessageLimit, Remaining: AssistantMessageLimit}
	if uc.store == nil {
		return state, nil
	}

	seen, err := uc.store.IsIntroSeen(ctx, userID)
	if err != nil {
		return nil, err
	}
	state.IntroSeen = seen

	usage, err := uc.store.GetUsage(ctx, userID)
	if err != nil {
		return nil, err
	}
	state.ResetAt = usage.ResetAt
	if remaining := AssistantMessageLimit - usage.Count; remaining > 0 {
		state.Remaining = remaining
	} else {
		state.Remaining = 0
	}
	return state, nil
}

func (uc *UseCase) MarkIntroSeen(ctx context.Context, userID uint) error {
	if uc.store == nil {
		return nil
	}
	return uc.store.MarkIntroSeen(ctx, userID)
}
