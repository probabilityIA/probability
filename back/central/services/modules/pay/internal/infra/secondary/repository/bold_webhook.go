package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/dtos"
	"gorm.io/datatypes"
)

type boldWebhookEventRow struct {
	ID                   uuid.UUID
	BoldEventID          string
	Type                 string
	Subject              string
	Source               string
	OccurredAt           *time.Time
	Payload              datatypes.JSON
	SignatureValid       bool
	ProcessedAt          *time.Time
	ProcessedError       *string
	PaymentTransactionID *uint
	WalletTransactionID  *uuid.UUID
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

func (boldWebhookEventRow) TableName() string {
	return "bold_webhook_events"
}

func (r *Repository) RecordBoldWebhookEvent(ctx context.Context, event *dtos.BoldWebhookEvent) (bool, error) {
	if event.BoldEventID == "" {
		return false, fmt.Errorf("bold event id required")
	}

	newID := uuid.New()
	row := boldWebhookEventRow{
		ID:                   newID,
		BoldEventID:          event.BoldEventID,
		Type:                 event.Type,
		Subject:              event.Subject,
		Source:               event.Source,
		OccurredAt:           event.OccurredAt,
		Payload:              datatypes.JSON(event.Payload),
		SignatureValid:       event.SignatureValid,
		ProcessedAt:          event.ProcessedAt,
		ProcessedError:       event.ProcessedError,
		PaymentTransactionID: event.PaymentTransactionID,
	}

	res := r.db.Conn(ctx).
		Where("bold_event_id = ?", event.BoldEventID).
		Attrs(row).
		FirstOrCreate(&row)
	if res.Error != nil {
		if isUniqueViolation(res.Error) {
			return false, nil
		}
		return false, fmt.Errorf("upsert bold webhook event: %w", res.Error)
	}

	event.ID = row.ID
	created := res.RowsAffected == 1
	return created, nil
}

func (r *Repository) MarkBoldWebhookProcessed(ctx context.Context, id uuid.UUID, paymentTransactionID *uint, processErr error) error {
	now := time.Now()
	updates := map[string]any{
		"processed_at":           &now,
		"payment_transaction_id": paymentTransactionID,
	}
	if processErr != nil {
		msg := processErr.Error()
		updates["processed_error"] = &msg
	} else {
		updates["processed_error"] = nil
	}
	return r.db.Conn(ctx).
		Table("bold_webhook_events").
		Where("id = ?", id).
		Updates(updates).Error
}

func (r *Repository) LinkBoldWebhookToWalletTransaction(ctx context.Context, eventID, walletTransactionID uuid.UUID) error {
	return r.db.Conn(ctx).
		Table("bold_webhook_events").
		Where("id = ?", eventID).
		Update("wallet_transaction_id", walletTransactionID).Error
}

func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "duplicate key") || strings.Contains(msg, "unique constraint")
}
