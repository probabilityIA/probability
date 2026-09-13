package repository

import (
	"context"
	"strings"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
)

type campaignAudienceQuerier struct {
	db     db.IDatabase
	logger log.ILogger
}

type campaignCandidateRow struct {
	ClientID   uint
	BusinessID uint
	Name       string
	Phone      string
	City       string
}

type audienceLocationRow struct {
	CityName  string
	StateName string
	Clients   uint
}

func locationKey(expression string) string {
	normalized := "lower(unaccent(trim(" + expression + ")))"
	return "btrim(replace(replace(" + normalized + ", ', d.c.', ''), ' d.c.', ''))"
}

var (
	rawCityExpression = "COALESCE(NULLIF(geo.geo_city, ''), COALESCE(addr.city, ''))"
	stateExpression   = "COALESCE(NULLIF(keyed.geo_state, ''), cat.state_name, '')"

	campaignAudienceBase = `
	WITH city_catalog AS (
		SELECT norm_name, min(city_name) AS city_name, min(state_name) AS state_name
		FROM (
			SELECT
				` + locationKey("gc.name") + ` AS norm_name,
				gc.name AS city_name,
				gs.name AS state_name
			FROM geozones gc
			JOIN geozones gs ON gs.id = gc.parent_id AND gs.type = 'state'
			WHERE gc.type = 'city'
			  AND gc.deleted_at IS NULL
			  AND gc.business_id = 0
		) catalog
		GROUP BY norm_name
		HAVING count(*) = 1
	),
	client_geo AS (
		SELECT DISTINCT ON (o.customer_id)
			o.customer_id,
			COALESCE(gc.name, '') AS geo_city,
			COALESCE(gs.name, '') AS geo_state
		FROM orders o
		LEFT JOIN geozones gc ON gc.id = o.geozone_city_id
		LEFT JOIN geozones gs ON gs.id = o.geozone_state_id
		WHERE o.business_id = @business_id
		  AND o.deleted_at IS NULL
		  AND o.geozone_city_id IS NOT NULL
		ORDER BY o.customer_id, o.id DESC
	),
	raw AS (
		SELECT
			c.id AS client_id,
			c.business_id,
			c.name,
			regexp_replace(COALESCE(c.phone, ''), '[^0-9]', '', 'g') AS phone_key,
			COALESCE(c.phone, '') AS phone,
			c.accepts_marketing,
			` + rawCityExpression + ` AS city,
			COALESCE(geo.geo_state, '') AS geo_state
		FROM client c
		LEFT JOIN LATERAL (
			SELECT city
			FROM customer_address
			WHERE customer_id = c.id AND deleted_at IS NULL
			ORDER BY is_primary DESC, times_used DESC, id DESC
			LIMIT 1
		) addr ON true
		LEFT JOIN client_geo geo ON geo.customer_id = c.id
		WHERE c.business_id = @business_id
		  AND c.deleted_at IS NULL
	),
	keyed AS (
		SELECT raw.*, ` + locationKey("raw.city") + ` AS city_key
		FROM raw
	),
	base AS (
		SELECT
			keyed.client_id,
			keyed.business_id,
			keyed.name,
			keyed.phone_key,
			keyed.phone,
			keyed.accepts_marketing,
			keyed.city_key,
			COALESCE(NULLIF(cat.city_name, ''), keyed.city) AS city,
			` + stateExpression + ` AS state,
			` + locationKey(stateExpression) + ` AS state_key
		FROM keyed
		LEFT JOIN city_catalog cat ON cat.norm_name = keyed.city_key
	)
`
)

func (q *campaignAudienceQuerier) buildFilters(params entities.CampaignAudienceParams, audienceType string) (string, map[string]any) {
	conditions := []string{}
	args := map[string]any{}

	if audienceType == entities.CampaignAudienceFiltered {
		if strings.TrimSpace(params.City) != "" {
			conditions = append(conditions, "city_key = "+locationKey("@city"))
			args["city"] = strings.TrimSpace(params.City)
		}
		if strings.TrimSpace(params.State) != "" {
			conditions = append(conditions, "state_key = "+locationKey("@state"))
			args["state"] = strings.TrimSpace(params.State)
		}
		if params.CreatedFromDays > 0 {
			conditions = append(conditions, "client_id IN (SELECT id FROM client WHERE business_id = @business_id AND created_at >= NOW() - make_interval(days => @created_from_days))")
			args["created_from_days"] = params.CreatedFromDays
		}
		if params.OnlyWithoutOrder {
			conditions = append(conditions, "NOT EXISTS (SELECT 1 FROM orders o WHERE o.customer_id = base.client_id AND o.deleted_at IS NULL)")
		}
		if len(params.ClientIDs) > 0 {
			conditions = append(conditions, "client_id IN @client_ids")
			args["client_ids"] = params.ClientIDs
		}
		if params.RegisteredBeforeDays > 0 {
			conditions = append(conditions, "client_id IN (SELECT id FROM client WHERE business_id = @business_id AND created_at <= NOW() - make_interval(days => @registered_before_days))")
			args["registered_before_days"] = params.RegisteredBeforeDays
		}
		if params.MinOrders > 0 {
			conditions = append(conditions, `EXISTS (
				SELECT 1 FROM customer_summary cs
				WHERE cs.customer_id = base.client_id
				  AND cs.business_id = @business_id
				  AND cs.deleted_at IS NULL
				  AND cs.total_orders >= @min_orders
			)`)
			args["min_orders"] = params.MinOrders
		}
		if params.MinSpent > 0 {
			conditions = append(conditions, `EXISTS (
				SELECT 1 FROM customer_summary cs
				WHERE cs.customer_id = base.client_id
				  AND cs.business_id = @business_id
				  AND cs.deleted_at IS NULL
				  AND cs.total_spent >= @min_spent
			)`)
			args["min_spent"] = params.MinSpent
		}
		if params.LastPurchaseBeforeDays > 0 {
			conditions = append(conditions, `EXISTS (
				SELECT 1 FROM customer_summary cs
				WHERE cs.customer_id = base.client_id
				  AND cs.business_id = @business_id
				  AND cs.deleted_at IS NULL
				  AND cs.last_order_at IS NOT NULL
				  AND cs.last_order_at <= NOW() - make_interval(days => @last_purchase_before_days)
			)`)
			args["last_purchase_before_days"] = params.LastPurchaseBeforeDays
		}
	}

	if params.ExcludeRecentDays > 0 {
		maxMessages := params.ExcludeRecentMax
		if maxMessages <= 0 {
			maxMessages = 1
		}

		conditions = append(conditions, `(
			SELECT count(*)
			FROM whatsapp_campaign_sends s
			WHERE s.client_id = base.client_id
			  AND s.business_id = @business_id
			  AND s.deleted_at IS NULL
			  AND s.status IN ('sent', 'delivered', 'read', 'replied')
			  AND s.sent_at >= NOW() - make_interval(days => @exclude_recent_days)
		) < @exclude_recent_max`)

		args["exclude_recent_days"] = params.ExcludeRecentDays
		args["exclude_recent_max"] = maxMessages
	}

	if len(conditions) == 0 {
		return "", args
	}

	return " AND " + strings.Join(conditions, " AND "), args
}

func (q *campaignAudienceQuerier) FindCampaignCandidates(ctx context.Context, businessID uint, params entities.CampaignAudienceParams, audienceType string, limit int) ([]entities.CampaignCandidate, error) {
	filters, args := q.buildFilters(params, audienceType)
	args["business_id"] = businessID
	args["limit"] = limit

	query := campaignAudienceBase + `
		, eligible AS (
			SELECT * FROM base
			WHERE accepts_marketing = true
			  AND phone_key <> ''
			  AND length(phone_key) >= 10
	` + filters + `
		),
		deduped AS (
			SELECT DISTINCT ON (phone_key) *
			FROM eligible
			ORDER BY phone_key, client_id DESC
		)
		SELECT client_id, business_id, name, phone, city
		FROM deduped
		ORDER BY client_id
		LIMIT @limit`

	var rows []campaignCandidateRow
	if err := q.db.Conn(ctx).Raw(query, args).Scan(&rows).Error; err != nil {
		q.logger.Error().Err(err).Uint("business_id", businessID).Msg("Error finding campaign candidates")
		return nil, err
	}

	out := make([]entities.CampaignCandidate, 0, len(rows))
	for _, row := range rows {
		out = append(out, entities.CampaignCandidate{
			ClientID:   row.ClientID,
			BusinessID: row.BusinessID,
			Name:       row.Name,
			FirstName:  firstName(row.Name),
			Phone:      row.Phone,
			City:       row.City,
		})
	}

	return out, nil
}

func (q *campaignAudienceQuerier) CountCampaignAudience(ctx context.Context, businessID uint, params entities.CampaignAudienceParams, audienceType string) (uint, uint, uint, uint, error) {
	filters, args := q.buildFilters(params, audienceType)
	args["business_id"] = businessID

	query := campaignAudienceBase + `
		, filtered AS (
			SELECT * FROM base WHERE true
	` + filters + `
		)
		SELECT
			COUNT(*) AS total,
			COUNT(*) FILTER (WHERE accepts_marketing = false) AS opted_out,
			COUNT(*) FILTER (WHERE phone_key = '' OR length(phone_key) < 10) AS no_phone,
			COUNT(DISTINCT phone_key) FILTER (
				WHERE accepts_marketing = true AND phone_key <> '' AND length(phone_key) >= 10
			) AS reachable
		FROM filtered`

	var result struct {
		Total     uint
		OptedOut  uint
		NoPhone   uint
		Reachable uint
	}

	if err := q.db.Conn(ctx).Raw(query, args).Scan(&result).Error; err != nil {
		q.logger.Error().Err(err).Uint("business_id", businessID).Msg("Error counting campaign audience")
		return 0, 0, 0, 0, err
	}

	return result.Total, result.OptedOut, result.NoPhone, result.Reachable, nil
}

func (q *campaignAudienceQuerier) ListAudienceLocations(ctx context.Context, businessID uint) ([]entities.AudienceLocation, error) {
	args := map[string]any{"business_id": businessID}

	query := campaignAudienceBase + `
		SELECT
			initcap(lower(max(city))) AS city_name,
			initcap(lower(max(state))) AS state_name,
			COUNT(DISTINCT phone_key) AS clients
		FROM base
		WHERE accepts_marketing = true
		  AND phone_key <> ''
		  AND length(phone_key) >= 10
		  AND city_key <> ''
		GROUP BY city_key, state_key
		ORDER BY clients DESC, city_name`

	var rows []audienceLocationRow
	if err := q.db.Conn(ctx).Raw(query, args).Scan(&rows).Error; err != nil {
		q.logger.Error().Err(err).Uint("business_id", businessID).Msg("Error listing audience locations")
		return nil, err
	}

	out := make([]entities.AudienceLocation, 0, len(rows))
	for _, row := range rows {
		out = append(out, entities.AudienceLocation{
			City:    row.CityName,
			State:   row.StateName,
			Clients: row.Clients,
		})
	}

	return out, nil
}

func firstName(fullName string) string {
	trimmed := strings.TrimSpace(fullName)
	if trimmed == "" {
		return ""
	}
	parts := strings.Fields(trimmed)
	return parts[0]
}

func newCampaignAudienceQuerier(database db.IDatabase, logger log.ILogger) *campaignAudienceQuerier {
	return &campaignAudienceQuerier{db: database, logger: logger.WithModule("campaign_audience_querier")}
}
