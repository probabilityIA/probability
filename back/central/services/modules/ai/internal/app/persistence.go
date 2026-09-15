package app

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/errors"
)

func (uc *UseCase) PersistMessage(ctx context.Context, record entities.MessageRecord) error {
	if uc.conversations == nil {
		return nil
	}
	return uc.conversations.SaveMessage(ctx, record)
}

func (uc *UseCase) PersistFeedback(ctx context.Context, userID uint, messageID string, value int, at time.Time) error {
	if uc.conversations == nil {
		return nil
	}
	found, err := uc.conversations.ApplyFeedback(ctx, userID, messageID, value, at)
	if err != nil {
		return err
	}
	if !found {
		return domainerrors.ErrMessageNotFound
	}
	return nil
}

func (uc *UseCase) PersistClick(ctx context.Context, userID uint, messageID string, at time.Time) error {
	if uc.conversations == nil {
		return nil
	}
	found, err := uc.conversations.ApplyClick(ctx, userID, messageID, at)
	if err != nil {
		return err
	}
	if !found {
		return domainerrors.ErrMessageNotFound
	}
	return nil
}

func (uc *UseCase) PurgeExpiredConversations(ctx context.Context) (int64, error) {
	if uc.conversations == nil {
		return 0, nil
	}
	return uc.conversations.DeleteOlderThan(ctx, time.Now().Add(-ConversationRetention))
}
