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

const (
	boldWebhookSource           = "bold"
	boldIntegrationTypeCodeRaw  = "bold_pay"
	boldWebhookRetentionDays    = 15
)

type webhookLogRow struct {
	ID                  uuid.UUID      `gorm:"column:id;primaryKey"`
	CreatedAt           time.Time      `gorm:"column:created_at"`
	UpdatedAt           time.Time      `gorm:"column:updated_at"`
	Source              string         `gorm:"column:source"`
	EventType           string         `gorm:"column:event_type"`
	URL                 string         `gorm:"column:url"`
	Method              string         `gorm:"column:method"`
	Headers             datatypes.JSON `gorm:"column:headers"`
	RequestBody         datatypes.JSON `gorm:"column:request_body"`
	RemoteIP            string         `gorm:"column:remote_ip"`
	Status              string         `gorm:"column:status"`
	ResponseCode        int            `gorm:"column:response_code"`
	ProcessedAt         *time.Time     `gorm:"column:processed_at"`
	ErrorMessage        *string        `gorm:"column:error_message"`
	SignatureValid      *bool          `gorm:"column:signature_valid"`
	IntegrationTypeID   *uint          `gorm:"column:integration_type_id"`
	IntegrationTypeCode *string        `gorm:"column:integration_type_code"`
	IntegrationID       *uint          `gorm:"column:integration_id"`
	CorrelationID       *string        `gorm:"column:correlation_id"`
	RetentionUntil      *time.Time     `gorm:"column:retention_until"`
}

func (webhookLogRow) TableName() string {
	return "webhook_logs"
}

type RawWebhookLogRepository struct {
	db                db.IDatabase
	log               log.ILogger
	integrationTypeID *uint
}

func NewRawWebhookLogRepository(database db.IDatabase, logger log.ILogger) ports.IRawWebhookLogger {
	return &RawWebhookLogRepository{
		db:  database,
		log: logger.WithModule("bold.raw_webhook_log_repository"),
	}
}

func (r *RawWebhookLogRepository) SetIntegrationTypeID(id uint) {
	if id > 0 {
		copy := id
		r.integrationTypeID = &copy
	}
}

func (r *RawWebhookLogRepository) LogIncoming(ctx context.Context, raw *ports.RawWebhookLog) error {
	if raw == nil {
		return fmt.Errorf("raw webhook log is nil")
	}
	id := uuid.New()
	if parsed, err := uuid.Parse(raw.ID); err == nil {
		id = parsed
	}
	raw.ID = id.String()

	path := "/api/v1/webhooks/bold"
	if raw.Endpoint == "test" {
		path = "/api/v1/webhooks/bold/test"
	}

	headers := map[string]string{}
	if raw.SignatureHeader != "" {
		headers["x-bold-signature"] = raw.SignatureHeader
	}
	headersJSON, _ := json.Marshal(headers)

	var bodyJSON datatypes.JSON
	if len(raw.Body) > 0 && json.Valid(raw.Body) {
		bodyJSON = datatypes.JSON(raw.Body)
	} else if len(raw.Body) > 0 {
		wrap, _ := json.Marshal(map[string]string{"raw": string(raw.Body)})
		bodyJSON = datatypes.JSON(wrap)
	} else {
		bodyJSON = datatypes.JSON([]byte("{}"))
	}

	code := boldIntegrationTypeCodeRaw
	retentionUntil := time.Now().Add(boldWebhookRetentionDays * 24 * time.Hour)

	row := webhookLogRow{
		ID:                  id,
		Source:              boldWebhookSource,
		EventType:           coalesce(raw.EventType, "unknown"),
		URL:                 path,
		Method:              "POST",
		Headers:             datatypes.JSON(headersJSON),
		RequestBody:         bodyJSON,
		Status:              "received",
		ResponseCode:        0,
		IntegrationTypeID:   r.integrationTypeID,
		IntegrationTypeCode: &code,
		RetentionUntil:      &retentionUntil,
	}
	if raw.MerchantReference != "" {
		ref := raw.MerchantReference
		row.CorrelationID = &ref
	}

	if err := r.db.Conn(ctx).Create(&row).Error; err != nil {
		return fmt.Errorf("create webhook_logs row: %w", err)
	}
	return nil
}

func (r *RawWebhookLogRepository) UpdateResult(ctx context.Context, result *ports.RawWebhookResult) error {
	if result == nil {
		return fmt.Errorf("result is nil")
	}
	id, err := uuid.Parse(result.ID)
	if err != nil {
		return fmt.Errorf("invalid webhook log id: %w", err)
	}
	now := time.Now()
	updates := map[string]interface{}{
		"status":        result.Status,
		"response_code": result.HTTPStatus,
		"processed_at":  now,
		"updated_at":    now,
	}
	switch result.Status {
	case "ok":
		t := true
		updates["signature_valid"] = &t
	case "invalid_signature":
		f := false
		updates["signature_valid"] = &f
	}
	if result.ErrorDetail != "" {
		updates["error_message"] = result.ErrorDetail
	}
	if err := r.db.Conn(ctx).Table("webhook_logs").Where("id = ?", id).Updates(updates).Error; err != nil {
		return fmt.Errorf("update webhook_logs row: %w", err)
	}
	return nil
}

func (r *RawWebhookLogRepository) DeleteOlderThan(ctx context.Context, days int) (int64, error) {
	if days <= 0 {
		return 0, fmt.Errorf("invalid retention days: %d", days)
	}
	result := r.db.Conn(ctx).
		Where("source = ?", boldWebhookSource).
		Where("retention_until IS NOT NULL AND retention_until < ?", time.Now()).
		Delete(&webhookLogRow{})
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func coalesce(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
