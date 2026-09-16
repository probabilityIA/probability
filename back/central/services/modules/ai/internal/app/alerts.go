package app

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/entities"
)

const (
	DefaultAlertsPageSize = 20
	MaxAlertsPageSize     = 50
	RecentAlertsForModel  = 8
	AlertRetention        = 90 * 24 * time.Hour
)

func (uc *UseCase) IngestAlertEvent(ctx context.Context, event dtos.AlertEvent) error {
	if uc.alerts == nil || event.ID == "" || event.BusinessID == 0 {
		return nil
	}
	alert, ok := buildAlert(event)
	if !ok {
		return nil
	}
	if alert.CreatedAt.IsZero() {
		alert.CreatedAt = uc.now()
	}
	alert.ID = uuid.NewString()

	created, err := uc.alerts.SaveAlert(ctx, *alert)
	if err != nil {
		return err
	}
	if created {
		uc.log.Info(ctx).Uint("business_id", alert.BusinessID).Str("event_type", alert.EventType).Msg("[ai.assistant] alerta registrada")
	}
	return nil
}

func (uc *UseCase) ListAlerts(ctx context.Context, query dtos.AlertQuery) (*dtos.PaginatedResponse[entities.Alert], error) {
	if uc.alerts == nil {
		return &dtos.PaginatedResponse[entities.Alert]{Data: []entities.Alert{}, Page: 1, PageSize: DefaultAlertsPageSize}, nil
	}
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 {
		query.PageSize = DefaultAlertsPageSize
	}
	if query.PageSize > MaxAlertsPageSize {
		query.PageSize = MaxAlertsPageSize
	}
	items, total, err := uc.alerts.ListAlerts(ctx, query)
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []entities.Alert{}
	}
	pages := int((total + int64(query.PageSize) - 1) / int64(query.PageSize))
	return &dtos.PaginatedResponse[entities.Alert]{Data: items, Total: total, Page: query.Page, PageSize: query.PageSize, TotalPages: pages}, nil
}

func (uc *UseCase) GetAlertsUnread(ctx context.Context, businessID, userID uint) (*entities.AlertsUnread, error) {
	if uc.alerts == nil {
		return &entities.AlertsUnread{}, nil
	}
	count, latest, err := uc.alerts.CountUnread(ctx, businessID, userID)
	if err != nil {
		return nil, err
	}
	return &entities.AlertsUnread{Count: count, Latest: latest}, nil
}

func (uc *UseCase) MarkAlertsSeen(ctx context.Context, businessID, userID uint) error {
	if uc.alerts == nil {
		return nil
	}
	return uc.alerts.MarkSeen(ctx, businessID, userID, uc.now())
}

func (uc *UseCase) PurgeExpiredAlerts(ctx context.Context) (int64, error) {
	if uc.alerts == nil {
		return 0, nil
	}
	return uc.alerts.DeleteAlertsOlderThan(ctx, uc.now().Add(-AlertRetention))
}
