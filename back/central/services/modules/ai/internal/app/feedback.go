package app

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/errors"
)

func (uc *UseCase) SubmitFeedback(ctx context.Context, userID uint, messageID string, value int) error {
	id, err := parseMessageID(messageID)
	if err != nil {
		return err
	}
	if value != entities.FeedbackNone && value != entities.FeedbackPositive && value != entities.FeedbackNegative {
		return domainerrors.ErrInvalidFeedback
	}

	now := time.Now()
	if uc.recorder != nil {
		return uc.recorder.RecordFeedback(ctx, userID, id, value, now)
	}
	return uc.PersistFeedback(ctx, userID, id, value, now)
}

func (uc *UseCase) MarkDestinationClicked(ctx context.Context, userID uint, messageID string) error {
	id, err := parseMessageID(messageID)
	if err != nil {
		return err
	}

	now := time.Now()
	if uc.recorder != nil {
		return uc.recorder.RecordClick(ctx, userID, id, now)
	}
	return uc.PersistClick(ctx, userID, id, now)
}

func parseMessageID(raw string) (string, error) {
	parsed, err := uuid.Parse(strings.TrimSpace(raw))
	if err != nil {
		return "", domainerrors.ErrInvalidMessageID
	}
	return parsed.String(), nil
}
