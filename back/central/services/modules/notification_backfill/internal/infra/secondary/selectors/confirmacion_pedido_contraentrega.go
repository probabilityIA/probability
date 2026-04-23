package selectors

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/notification_backfill/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/notification_backfill/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/notification_backfill/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
)

var confirmationEligibleStatuses = []string{
	"pending",
	"confirmed",
	"processing",
	"picking",
	"packing",
}

type confirmationSelector struct {
	db       db.IDatabase
	log      log.ILogger
	dispatch func(context.Context, entities.Candidate) error
}

func NewConfirmationSelector(database db.IDatabase, logger log.ILogger, dispatcher ConfirmationDispatcher) ports.IEligibilitySelector {
	return &confirmationSelector{
		db:       database,
		log:      logger.WithModule("notification_backfill.selector.confirmacion"),
		dispatch: NewConfirmationDispatchAdapter(dispatcher),
	}
}

func (s *confirmationSelector) EventCode() string { return "confirmacion_pedido_contraentrega" }
func (s *confirmationSelector) EventName() string { return "Confirmación de pedido contra entrega" }
func (s *confirmationSelector) Channel() string   { return "whatsapp" }

func (s *confirmationSelector) Preview(ctx context.Context, filter dtos.BackfillFilter) ([]entities.Candidate, error) {
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
		Select("orders.id, orders.order_number, orders.business_id, orders.customer_phone, COALESCE(orders.tracking_number,'') AS tracking_number, orders.status").
		Where("orders.deleted_at IS NULL").
		Where("COALESCE(orders.is_test, false) = false").
		Where("COALESCE(orders.is_paid, false) = false").
		Where("COALESCE(orders.is_confirmed, false) = false").
		Where("orders.customer_phone IS NOT NULL AND TRIM(orders.customer_phone) <> ''").
		Where("orders.status IN ?", confirmationEligibleStatuses).
		Where("orders.created_at >= NOW() - (? * INTERVAL '1 day')", days).
		Where(`NOT EXISTS (
			SELECT 1
			FROM whatsapp_message_logs ml
			JOIN whatsapp_conversations wc ON wc.id = ml.conversation_id
			WHERE ml.template_name = 'confirmacion_pedido_contraentrega'
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

func (s *confirmationSelector) Dispatch(ctx context.Context, c entities.Candidate) error {
	return s.dispatch(ctx, c)
}
