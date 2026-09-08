package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
)

const whatsappIntegrationTypeID = 2

type campaignSenderQuerier struct {
	db     db.IDatabase
	logger log.ILogger
}

func (q *campaignSenderQuerier) GetOwnWhatsappSender(ctx context.Context, businessID uint) (uint, string, error) {
	var row struct {
		ID          uint
		PhoneNumber string
	}

	err := q.db.Conn(ctx).
		Table("integrations").
		Select(`id, COALESCE(config->>'display_phone_number', config->>'phone_number_id', '') AS phone_number`).
		Where("business_id = ? AND integration_type_id = ? AND is_active = true AND deleted_at IS NULL", businessID, whatsappIntegrationTypeID).
		Where(`COALESCE(config->>'waba_id', '') <> ''`).
		Where(`COALESCE(config->>'phone_number_id', '') <> ''`).
		Limit(1).
		Scan(&row).Error

	if err != nil {
		q.logger.Error().Err(err).Uint("business_id", businessID).Msg("Error resolving own whatsapp sender")
		return 0, "", err
	}

	return row.ID, row.PhoneNumber, nil
}

func newCampaignSenderQuerier(database db.IDatabase, logger log.ILogger) *campaignSenderQuerier {
	return &campaignSenderQuerier{db: database, logger: logger.WithModule("campaign_sender_querier")}
}
