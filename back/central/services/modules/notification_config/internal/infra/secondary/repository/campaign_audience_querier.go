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

const campaignAudienceBase = `
	WITH base AS (
		SELECT
			c.id AS client_id,
			c.business_id,
			c.name,
			regexp_replace(COALESCE(c.phone, ''), '[^0-9]', '', 'g') AS phone_key,
			COALESCE(c.phone, '') AS phone,
			c.accepts_marketing,
			COALESCE(addr.city, '') AS city
		FROM client c
		LEFT JOIN LATERAL (
			SELECT city
			FROM customer_address
			WHERE customer_id = c.id AND deleted_at IS NULL
			ORDER BY is_primary DESC, times_used DESC, id DESC
			LIMIT 1
		) addr ON true
		WHERE c.business_id = @business_id
		  AND c.deleted_at IS NULL
	)
`

func (q *campaignAudienceQuerier) buildFilters(params entities.CampaignAudienceParams, audienceType string) (string, map[string]any) {
	conditions := []string{}
	args := map[string]any{}

	if audienceType == entities.CampaignAudienceFiltered {
		if strings.TrimSpace(params.City) != "" {
			conditions = append(conditions, "city ILIKE @city")
			args["city"] = "%" + strings.TrimSpace(params.City) + "%"
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
