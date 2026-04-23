package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/notification_backfill/internal/domain/dtos"
)

func (uc *useCase) ListEvents(ctx context.Context) []dtos.RegisteredEventResponse {
	selectors := uc.registry.List()
	out := make([]dtos.RegisteredEventResponse, 0, len(selectors))
	for _, s := range selectors {
		out = append(out, dtos.RegisteredEventResponse{
			EventCode: s.EventCode(),
			EventName: s.EventName(),
			Channel:   s.Channel(),
		})
	}
	return out
}
