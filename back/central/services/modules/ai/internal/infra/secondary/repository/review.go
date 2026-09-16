package repository

import (
	"context"
	"strings"

	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

const topDestinationsLimit = 8

var likeEscaper = strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)

func (r *Repository) scoped(ctx context.Context, filter dtos.ReviewFilter, withKind bool) *gorm.DB {
	query := r.db.Conn(ctx).Model(&models.AssistantMessage{})
	if filter.BusinessID != nil {
		query = query.Where("business_id = ?", *filter.BusinessID)
	}
	if filter.From != nil {
		query = query.Where("created_at >= ?", *filter.From)
	}
	if filter.To != nil {
		query = query.Where("created_at < ?", *filter.To)
	}
	if filter.Search != "" {
		pattern := "%" + likeEscaper.Replace(filter.Search) + "%"
		query = query.Where("(question ILIKE ? OR answer ILIKE ?)", pattern, pattern)
	}
	if !withKind {
		return query
	}

	switch filter.Kind {
	case dtos.ReviewNoDestination:
		query = query.Where("destination_key = '' AND error_code = ''")
	case dtos.ReviewNegative:
		query = query.Where("feedback = -1")
	case dtos.ReviewPositive:
		query = query.Where("feedback = 1")
	case dtos.ReviewErrors:
		query = query.Where("error_code <> ''")
	case dtos.ReviewNotClicked:
		query = query.Where("destination_key <> '' AND clicked_at IS NULL")
	}
	return query
}

func (r *Repository) ListMessages(ctx context.Context, filter dtos.ReviewFilter) ([]entities.ReviewMessage, int64, error) {
	var total int64
	if err := r.scoped(ctx, filter, true).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []models.AssistantMessage
	err := r.scoped(ctx, filter, true).
		Preload("Business", func(db *gorm.DB) *gorm.DB { return db.Select("id", "name") }).
		Preload("User", func(db *gorm.DB) *gorm.DB { return db.Select("id", "name", "email") }).
		Order("created_at DESC").
		Offset((filter.Page - 1) * filter.PageSize).
		Limit(filter.PageSize).
		Find(&rows).Error
	if err != nil {
		return nil, 0, err
	}

	items := make([]entities.ReviewMessage, 0, len(rows))
	for _, row := range rows {
		items = append(items, toReviewMessage(row))
	}
	return items, total, nil
}

type summaryRow struct {
	Messages        int64
	Conversations   int64
	Users           int64
	InputTokens     int64
	OutputTokens    int64
	NoDestination   int64
	Errors          int64
	Positive        int64
	Negative        int64
	WithDestination int64
	Clicked         int64
}

type destinationRow struct {
	DestKey string
	Total   int64
}

func (r *Repository) Summary(ctx context.Context, filter dtos.ReviewFilter) (*entities.ReviewSummary, error) {
	var row summaryRow
	err := r.scoped(ctx, filter, false).Select(
		"COUNT(*) AS messages, " +
			"COUNT(DISTINCT conversation_id) AS conversations, " +
			"COUNT(DISTINCT user_id) AS users, " +
			"COALESCE(SUM(input_tokens), 0) AS input_tokens, " +
			"COALESCE(SUM(output_tokens), 0) AS output_tokens, " +
			"COUNT(*) FILTER (WHERE destination_key = '' AND error_code = '') AS no_destination, " +
			"COUNT(*) FILTER (WHERE error_code <> '') AS errors, " +
			"COUNT(*) FILTER (WHERE feedback = 1) AS positive, " +
			"COUNT(*) FILTER (WHERE feedback = -1) AS negative, " +
			"COUNT(*) FILTER (WHERE destination_key <> '') AS with_destination, " +
			"COUNT(*) FILTER (WHERE clicked_at IS NOT NULL) AS clicked",
	).Scan(&row).Error
	if err != nil {
		return nil, err
	}

	var destinations []destinationRow
	err = r.scoped(ctx, filter, false).
		Select("destination_key AS dest_key, COUNT(*) AS total").
		Where("destination_key <> ''").
		Group("destination_key").
		Order("total DESC").
		Limit(topDestinationsLimit).
		Scan(&destinations).Error
	if err != nil {
		return nil, err
	}

	summary := &entities.ReviewSummary{
		Messages:        row.Messages,
		Conversations:   row.Conversations,
		Users:           row.Users,
		InputTokens:     row.InputTokens,
		OutputTokens:    row.OutputTokens,
		NoDestination:   row.NoDestination,
		Errors:          row.Errors,
		Positive:        row.Positive,
		Negative:        row.Negative,
		WithDestination: row.WithDestination,
		Clicked:         row.Clicked,
		TopDestinations: make([]entities.DestinationCount, 0, len(destinations)),
	}
	for _, d := range destinations {
		summary.TopDestinations = append(summary.TopDestinations, entities.DestinationCount{Key: d.DestKey, Count: d.Total})
	}
	return summary, nil
}

func toReviewMessage(row models.AssistantMessage) entities.ReviewMessage {
	item := entities.ReviewMessage{
		MessageRecord: entities.MessageRecord{
			ID:               row.ID,
			ConversationID:   row.ConversationID,
			BusinessID:       row.BusinessID,
			UserID:           row.UserID,
			Pathname:         row.Pathname,
			Question:         row.Question,
			Answer:           row.Answer,
			DestinationKey:   row.DestinationKey,
			DestinationRoute: row.DestinationRoute,
			ErrorCode:        row.ErrorCode,
			Model:            row.Model,
			InputTokens:      row.InputTokens,
			OutputTokens:     row.OutputTokens,
			LatencyMs:        row.LatencyMs,
			CreatedAt:        row.CreatedAt,
		},
		UserName:   row.User.Name,
		UserEmail:  row.User.Email,
		Feedback:   int(row.Feedback),
		FeedbackAt: row.FeedbackAt,
		ClickedAt:  row.ClickedAt,
	}
	if row.Business != nil {
		item.BusinessName = row.Business.Name
	}
	return item
}
