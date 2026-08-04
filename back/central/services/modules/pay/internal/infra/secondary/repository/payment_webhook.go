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

type paymentWebhookEventRow struct {
	ID                   uuid.UUID
	Gateway              string
	EventID              string
	Type                 string
	ExternalID           string
	ExternalState        string
	Reference            string
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

func (paymentWebhookEventRow) TableName() string {
	return "payment_webhook_events"
}

func (r *Repository) RecordPaymentWebhookEvent(ctx context.Context, event *dtos.PaymentWebhookEvent) (bool, error) {
	if event.Gateway == "" {
		return false, fmt.Errorf("payment webhook gateway required")
	}
	if event.EventID == "" {
		return false, fmt.Errorf("payment webhook event id required")
	}

	row := paymentWebhookEventRow{
		ID:                   uuid.New(),
		Gateway:              event.Gateway,
		EventID:              event.EventID,
		Type:                 event.Type,
		ExternalID:           event.ExternalID,
		ExternalState:        event.ExternalState,
		Reference:            event.Reference,
		Source:               event.Source,
		OccurredAt:           event.OccurredAt,
		Payload:              datatypes.JSON(event.Payload),
		SignatureValid:       event.SignatureValid,
		ProcessedAt:          event.ProcessedAt,
		ProcessedError:       event.ProcessedError,
		PaymentTransactionID: event.PaymentTransactionID,
	}

	res := r.db.Conn(ctx).
		Where("gateway = ? AND event_id = ?", event.Gateway, event.EventID).
		Attrs(row).
		FirstOrCreate(&row)
	if res.Error != nil {
		if isUniqueViolation(res.Error) {
			return false, nil
		}
		return false, fmt.Errorf("upsert payment webhook event: %w", res.Error)
	}

	event.ID = row.ID
	created := res.RowsAffected == 1
	return created, nil
}

func (r *Repository) MarkPaymentWebhookProcessed(ctx context.Context, id uuid.UUID, paymentTransactionID *uint, processErr error) error {
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
		Table("payment_webhook_events").
		Where("id = ?", id).
		Updates(updates).Error
}

func (r *Repository) LinkPaymentWebhookToWalletTransaction(ctx context.Context, eventID, walletTransactionID uuid.UUID) error {
	return r.db.Conn(ctx).
		Table("payment_webhook_events").
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
