package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/dtos"
	"gorm.io/datatypes"
)

type bancolombiaWebhookEventRow struct {
	ID                   uuid.UUID
	EventID              string
	Type                 string
	TransferState        string
	Reference            string
	OccurredAt           *time.Time
	Payload              datatypes.JSON
	SignatureValid       bool
	ProcessedAt          *time.Time
	ProcessedError       *string
	PaymentTransactionID *uint
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

func (bancolombiaWebhookEventRow) TableName() string {
	return "bancolombia_webhook_events"
}

func (r *Repository) RecordBancolombiaWebhookEvent(ctx context.Context, event *dtos.BancolombiaWebhookEvent) (bool, error) {
	if event.EventID == "" {
		return false, fmt.Errorf("bancolombia event id required")
	}

	row := bancolombiaWebhookEventRow{
		ID:                   uuid.New(),
		EventID:              event.EventID,
		Type:                 event.Type,
		TransferState:        event.TransferState,
		Reference:            event.Reference,
		OccurredAt:           event.OccurredAt,
		Payload:              datatypes.JSON(event.Payload),
		SignatureValid:       event.SignatureValid,
		ProcessedAt:          event.ProcessedAt,
		ProcessedError:       event.ProcessedError,
		PaymentTransactionID: event.PaymentTransactionID,
	}

	res := r.db.Conn(ctx).
		Where("event_id = ?", event.EventID).
		Attrs(row).
		FirstOrCreate(&row)
	if res.Error != nil {
		if isUniqueViolation(res.Error) {
			return false, nil
		}
		return false, fmt.Errorf("upsert bancolombia webhook event: %w", res.Error)
	}

	event.ID = row.ID
	created := res.RowsAffected == 1
	return created, nil
}

func (r *Repository) MarkBancolombiaWebhookProcessed(ctx context.Context, id uuid.UUID, paymentTransactionID *uint, processErr error) error {
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
		Table("bancolombia_webhook_events").
		Where("id = ?", id).
		Updates(updates).Error
}
