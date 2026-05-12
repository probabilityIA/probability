package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"

	"github.com/secamc93/probability/back/central/services/integrations/pay/bold/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
)

type rawLogRow struct {
	ID                uuid.UUID      `gorm:"column:id;primaryKey"`
	CreatedAt         time.Time      `gorm:"column:created_at"`
	UpdatedAt         time.Time      `gorm:"column:updated_at"`
	Endpoint          string         `gorm:"column:endpoint"`
	HTTPStatus        int            `gorm:"column:http_status"`
	Status            string         `gorm:"column:status"`
	SignatureHeader   string         `gorm:"column:signature_header"`
	BoldEventID       string         `gorm:"column:bold_event_id"`
	EventType         string         `gorm:"column:event_type"`
	MerchantReference string         `gorm:"column:merchant_reference"`
	PaymentID         string         `gorm:"column:payment_id"`
	BodySize          int            `gorm:"column:body_size"`
	BodyJSON          datatypes.JSON `gorm:"column:body_json;type:jsonb"`
	BodyText          *string        `gorm:"column:body_text"`
	ErrorDetail       *string        `gorm:"column:error_detail"`
	ExpectedHash      string         `gorm:"column:expected_hash"`
}

func (rawLogRow) TableName() string {
	return "bold_webhook_raw_logs"
}

type RawWebhookLogRepository struct {
	db  db.IDatabase
	log log.ILogger
}

func NewRawWebhookLogRepository(database db.IDatabase, logger log.ILogger) ports.IRawWebhookLogger {
	return &RawWebhookLogRepository{
		db:  database,
		log: logger.WithModule("bold.raw_webhook_log_repository"),
	}
}

func (r *RawWebhookLogRepository) LogIncoming(ctx context.Context, raw *ports.RawWebhookLog) error {
	if raw == nil {
		return fmt.Errorf("raw webhook log is nil")
	}
	row := rawLogRow{
		Endpoint:          raw.Endpoint,
		Status:            "received",
		SignatureHeader:   raw.SignatureHeader,
		BoldEventID:       raw.BoldEventID,
		EventType:         raw.EventType,
		MerchantReference: raw.MerchantReference,
		PaymentID:         raw.PaymentID,
		BodySize:          raw.BodySize,
	}
	if id, err := uuid.Parse(raw.ID); err == nil {
		row.ID = id
	} else {
		row.ID = uuid.New()
		raw.ID = row.ID.String()
	}
	if json.Valid(raw.Body) {
		row.BodyJSON = datatypes.JSON(raw.Body)
	} else {
		text := string(raw.Body)
		row.BodyText = &text
	}
	if err := r.db.Conn(ctx).Create(&row).Error; err != nil {
		return fmt.Errorf("create raw webhook log: %w", err)
	}
	return nil
}

func (r *RawWebhookLogRepository) UpdateResult(ctx context.Context, result *ports.RawWebhookResult) error {
	if result == nil {
		return fmt.Errorf("result is nil")
	}
	id, err := uuid.Parse(result.ID)
	if err != nil {
		return fmt.Errorf("invalid raw webhook log id: %w", err)
	}
	updates := map[string]interface{}{
		"status":      result.Status,
		"http_status": result.HTTPStatus,
		"updated_at":  time.Now(),
	}
	if result.ErrorDetail != "" {
		updates["error_detail"] = result.ErrorDetail
	}
	if result.ExpectedHash != "" {
		updates["expected_hash"] = result.ExpectedHash
	}
	if err := r.db.Conn(ctx).Table("bold_webhook_raw_logs").Where("id = ?", id).Updates(updates).Error; err != nil {
		return fmt.Errorf("update raw webhook log: %w", err)
	}
	return nil
}

func (r *RawWebhookLogRepository) DeleteOlderThan(ctx context.Context, days int) (int64, error) {
	if days <= 0 {
		return 0, fmt.Errorf("invalid retention days: %d", days)
	}
	cutoff := time.Now().Add(-time.Duration(days) * 24 * time.Hour)
	result := r.db.Conn(ctx).Where("created_at < ?", cutoff).Delete(&rawLogRow{})
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}
