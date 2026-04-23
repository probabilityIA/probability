package selectors

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/notification_backfill/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/notification_backfill/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/notification_backfill/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
)

var guideEligibleStatuses = []string{
	"shipped",
	"ready_to_ship",
	"in_transit",
	"out_for_delivery",
	"picked_up",
	"assigned_to_driver",
	"delivery_novelty",
	"delivery_failed",
}

type guideSelector struct {
	db       db.IDatabase
	log      log.ILogger
	dispatch func(context.Context, entities.Candidate) error
}

func NewGuideSelector(database db.IDatabase, logger log.ILogger, dispatcher GuideDispatcher) ports.IEligibilitySelector {
	return &guideSelector{
		db:       database,
		log:      logger.WithModule("notification_backfill.selector.guia"),
		dispatch: NewGuideDispatchAdapter(dispatcher),
	}
}

func (s *guideSelector) EventCode() string { return "guia_envio_generada" }
func (s *guideSelector) EventName() string { return "Guía de envío generada" }
func (s *guideSelector) Channel() string   { return "whatsapp" }

func (s *guideSelector) Preview(ctx context.Context, filter dtos.BackfillFilter) ([]entities.Candidate, error) {
	days := filter.Days
	if days <= 0 {
		days = 4
	}
	limit := filter.Limit
	if limit <= 0 || limit > 2000 {
		limit = 500
	}

	q := s.db.Conn(ctx).
		Table("orders").
		Select("orders.id, orders.order_number, orders.business_id, orders.customer_phone, orders.tracking_number, orders.status").
		Where("orders.deleted_at IS NULL").
		Where("COALESCE(orders.is_test, false) = false").
		Where("orders.tracking_number IS NOT NULL AND TRIM(orders.tracking_number) <> ''").
		Where("orders.customer_phone IS NOT NULL AND TRIM(orders.customer_phone) <> ''").
		Where("orders.delivered_at IS NULL").
		Where("orders.status IN ?", guideEligibleStatuses).
		Where("orders.created_at >= NOW() - (? * INTERVAL '1 day')", days).
		Where(`NOT EXISTS (
			SELECT 1
			FROM whatsapp_message_logs ml
			JOIN whatsapp_conversations wc ON wc.id = ml.conversation_id
			WHERE ml.template_name = 'guia_envio_generada'
			  AND ml.status IN ('sent','delivered','read')
			  AND wc.business_id = orders.business_id
			  AND wc.order_number = orders.order_number
		)`)

	if filter.BusinessID != nil && *filter.BusinessID > 0 {
		q = q.Where("orders.business_id = ?", *filter.BusinessID)
	}

	type row struct {
		ID             string
		OrderNumber    string
		BusinessID     uint
		CustomerPhone  string
		TrackingNumber string
		Status         string
	}

	var rows []row
	if err := q.Order("orders.created_at ASC").Limit(limit).Scan(&rows).Error; err != nil {
		return nil, err
	}

	out := make([]entities.Candidate, 0, len(rows))
	for _, r := range rows {
		out = append(out, entities.Candidate{
			OrderID:        r.ID,
			OrderNumber:    r.OrderNumber,
			BusinessID:     r.BusinessID,
			CustomerPhone:  r.CustomerPhone,
			TrackingNumber: r.TrackingNumber,
			Status:         r.Status,
		})
	}
	return out, nil
}

func (s *guideSelector) Dispatch(ctx context.Context, c entities.Candidate) error {
	return s.dispatch(ctx, c)
}
