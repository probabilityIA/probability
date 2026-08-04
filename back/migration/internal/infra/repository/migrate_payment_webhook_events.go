package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

const backfillBoldWebhookEventsSQL = `
INSERT INTO payment_webhook_events (
    id, created_at, updated_at, gateway, event_id, type, external_id,
    external_state, reference, source, occurred_at, payload, signature_valid,
    processed_at, processed_error, payment_transaction_id, wallet_transaction_id
)
SELECT
    id, created_at, updated_at, 'bold', bold_event_id, type, COALESCE(subject, ''),
    '', '', COALESCE(source, ''), occurred_at, payload, signature_valid,
    processed_at, processed_error, payment_transaction_id, wallet_transaction_id
FROM bold_webhook_events
ON CONFLICT (gateway, event_id) DO NOTHING`

func (r *Repository) migratePaymentWebhookEvents(ctx context.Context) error {
	if err := r.db.Conn(ctx).AutoMigrate(&models.PaymentWebhookEvent{}); err != nil {
		return fmt.Errorf("migratePaymentWebhookEvents: creando tabla: %w", err)
	}

	var legacy *string
	if err := r.db.Conn(ctx).
		Raw(`SELECT to_regclass('public.bold_webhook_events')::text`).
		Scan(&legacy).Error; err != nil {
		return fmt.Errorf("migratePaymentWebhookEvents: verificando tabla legacy: %w", err)
	}
	if legacy == nil {
		return nil
	}

	res := r.db.Conn(ctx).Exec(backfillBoldWebhookEventsSQL)
	if res.Error != nil {
		return fmt.Errorf("migratePaymentWebhookEvents: backfill bold: %w", res.Error)
	}

	var pending int64
	if err := r.db.Conn(ctx).
		Raw(`SELECT count(*) FROM bold_webhook_events b
             WHERE NOT EXISTS (
                 SELECT 1 FROM payment_webhook_events p
                 WHERE p.gateway = 'bold' AND p.event_id = b.bold_event_id
             )`).
		Scan(&pending).Error; err != nil {
		return fmt.Errorf("migratePaymentWebhookEvents: verificando backfill: %w", err)
	}
	if pending > 0 {
		return fmt.Errorf("migratePaymentWebhookEvents: %d eventos de bold_webhook_events no migraron", pending)
	}

	return nil
}
