package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
)

type segmentQuerier struct {
	db     db.IDatabase
	logger log.ILogger
}

const inactiveCustomersQuery = `
WITH base AS (
	SELECT
		c.id AS client_id,
		c.business_id,
		c.name,
		c.phone,
		right(regexp_replace(c.phone, '\D', '', 'g'), 10) AS phone_key,
		(CURRENT_DATE - cs.last_order_at::date) AS days_inactive,
		cs.total_orders,
		cs.last_order_at,
		b.name AS business_name
	FROM client c
	JOIN customer_summary cs
		ON cs.customer_id = c.id AND cs.business_id = c.business_id
	JOIN business b ON b.id = c.business_id
	WHERE c.business_id = @business_id
		AND c.deleted_at IS NULL
		AND cs.deleted_at IS NULL
		AND coalesce(c.phone, '') <> ''
		AND length(regexp_replace(c.phone, '\D', '', 'g')) >= 10
		AND cs.last_order_at IS NOT NULL
		AND cs.last_order_at > DATE '2000-01-01'
		AND cs.last_order_at < NOW() - make_interval(days => @days_without_purchase)
		AND (@requires_opt_in = false OR c.accepts_marketing = true)
		AND (@min_orders = 0 OR cs.total_orders >= @min_orders)
		AND (@max_orders = 0 OR cs.total_orders <= @max_orders)
),
deduped AS (
	SELECT DISTINCT ON (phone_key) *
	FROM base
	ORDER BY phone_key, last_order_at DESC, client_id DESC
)
SELECT d.client_id, d.business_id, d.name, d.phone, d.days_inactive, d.total_orders, d.business_name
FROM deduped d
WHERE NOT EXISTS (
	SELECT 1
	FROM scheduled_notification_sends s
	WHERE s.rule_id = @rule_id
		AND s.status <> 'discarded'
		AND s.queued_at > NOW() - make_interval(days => @cooldown_days)
		AND (
			s.client_id = d.client_id
			OR right(regexp_replace(s.phone, '\D', '', 'g'), 10) = d.phone_key
		)
)
ORDER BY d.days_inactive DESC, d.client_id ASC
LIMIT @max_rows
`

func (q *segmentQuerier) FindInactiveCustomers(
	ctx context.Context,
	businessID uint,
	params entities.SegmentParams,
	requiresOptIn bool,
	cooldownDays uint,
	ruleID uint,
	limit int,
) ([]entities.SegmentCandidate, error) {
	if params.DaysWithoutPurchase <= 0 {
		return nil, fmt.Errorf("days_without_purchase debe ser mayor a cero")
	}
	if limit <= 0 {
		limit = 200
	}

	type row struct {
		ClientID     uint
		BusinessID   uint
		Name         string
		Phone        string
		DaysInactive int
		TotalOrders  int
		BusinessName string
	}

	var rows []row

	if err := q.db.Conn(ctx).Raw(inactiveCustomersQuery, map[string]any{
		"business_id":           businessID,
		"days_without_purchase": params.DaysWithoutPurchase,
		"requires_opt_in":       requiresOptIn,
		"min_orders":            params.MinOrders,
		"max_orders":            params.MaxOrders,
		"rule_id":               ruleID,
		"cooldown_days":         int(cooldownDays),
		"max_rows":              limit,
	}).Scan(&rows).Error; err != nil {
		q.logger.Error().Err(err).Uint("business_id", businessID).
			Msg("Error querying inactive customers segment")
		return nil, err
	}

	out := make([]entities.SegmentCandidate, 0, len(rows))
	for _, item := range rows {
		out = append(out, entities.SegmentCandidate{
			ClientID:     item.ClientID,
			BusinessID:   item.BusinessID,
			Name:         item.Name,
			Phone:        NormalizePhone(item.Phone),
			DaysInactive: item.DaysInactive,
			TotalOrders:  item.TotalOrders,
			BusinessName: item.BusinessName,
		})
	}

	return out, nil
}

func NormalizePhone(phone string) string {
	var digits strings.Builder
	for _, char := range phone {
		if char >= '0' && char <= '9' {
			digits.WriteRune(char)
		}
	}

	value := digits.String()
	if len(value) > 10 {
		value = value[len(value)-10:]
	}
	if value == "" {
		return ""
	}

	return "57" + value
}
