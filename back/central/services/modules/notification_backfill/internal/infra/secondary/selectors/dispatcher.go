package selectors

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/notification_backfill/internal/domain/entities"
)

type GuideDispatcher interface {
	SendGuideNotification(ctx context.Context, orderID string) error
}

type ConfirmationDispatcher interface {
	RequestConfirmation(ctx context.Context, orderID string) error
}

type guideDispatchAdapter struct {
	inner GuideDispatcher
}

func NewGuideDispatchAdapter(inner GuideDispatcher) func(context.Context, entities.Candidate) error {
	a := &guideDispatchAdapter{inner: inner}
	return a.dispatch
}

func (a *guideDispatchAdapter) dispatch(ctx context.Context, c entities.Candidate) error {
	return a.inner.SendGuideNotification(ctx, c.OrderID)
}

type confirmationDispatchAdapter struct {
	inner ConfirmationDispatcher
}

func NewConfirmationDispatchAdapter(inner ConfirmationDispatcher) func(context.Context, entities.Candidate) error {
	a := &confirmationDispatchAdapter{inner: inner}
	return a.dispatch
}

func (a *confirmationDispatchAdapter) dispatch(ctx context.Context, c entities.Candidate) error {
	return a.inner.RequestConfirmation(ctx, c.OrderID)
}
