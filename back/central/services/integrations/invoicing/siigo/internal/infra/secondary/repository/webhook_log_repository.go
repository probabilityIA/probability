package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
)

const (
	siigoWebhookSource        = "siigo"
	siigoWebhookRetentionDays = 30
)

type webhookLogRow struct {
	ID                  uuid.UUID      `gorm:"column:id;primaryKey"`
	CreatedAt           time.Time      `gorm:"column:created_at"`
	UpdatedAt           time.Time      `gorm:"column:updated_at"`
	Source              string         `gorm:"column:source"`
	EventType           string         `gorm:"column:event_type"`
	URL                 string         `gorm:"column:url"`
	Method              string         `gorm:"column:method"`
	RequestBody         datatypes.JSON `gorm:"column:request_body"`
	RemoteIP            string         `gorm:"column:remote_ip"`
	Status              string         `gorm:"column:status"`
	ResponseCode        int            `gorm:"column:response_code"`
	ProcessedAt         *time.Time     `gorm:"column:processed_at"`
	ErrorMessage        *string        `gorm:"column:error_message"`
	IntegrationTypeCode *string        `gorm:"column:integration_type_code"`
	IntegrationID       *uint          `gorm:"column:integration_id"`
	RetentionUntil      *time.Time     `gorm:"column:retention_until"`
}

func (webhookLogRow) TableName() string {
	return "webhook_logs"
}

type WebhookLogRepository struct {
	db  db.IDatabase
	log log.ILogger
}

func New(database db.IDatabase, logger log.ILogger) ports.IWebhookLogRepository {
	return &WebhookLogRepository{
		db:  database,
		log: logger.WithModule("siigo.webhook_log_repository"),
	}
}

func (r *WebhookLogRepository) LogIncoming(ctx context.Context, entry ports.WebhookLogEntry) (string, error) {
	id := uuid.New()

	var bodyJSON datatypes.JSON
	if len(entry.Body) > 0 && json.Valid(entry.Body) {
		bodyJSON = datatypes.JSON(entry.Body)
	} else if len(entry.Body) > 0 {
		wrap, _ := json.Marshal(map[string]string{"raw": string(entry.Body)})
		bodyJSON = datatypes.JSON(wrap)
	} else {
		bodyJSON = datatypes.JSON([]byte("{}"))
	}

	eventType := entry.EventType
	if eventType == "" {
		eventType = "unknown"
	}

	source := entry.Source
	if source == "" {
		source = siigoWebhookSource
	}

	retentionUntil := time.Now().Add(siigoWebhookRetentionDays * 24 * time.Hour)

	row := webhookLogRow{
		ID:             id,
		Source:         source,
		EventType:      eventType,
		URL:            entry.URL,
		Method:         "POST",
		RequestBody:    bodyJSON,
		RemoteIP:       entry.RemoteIP,
		Status:         "received",
		ResponseCode:   0,
		IntegrationID:  entry.IntegrationID,
		RetentionUntil: &retentionUntil,
	}
	if entry.IntegrationTypeCode != "" {
		code := entry.IntegrationTypeCode
		row.IntegrationTypeCode = &code
	}

	if err := r.db.Conn(ctx).Create(&row).Error; err != nil {
		return "", fmt.Errorf("create webhook_logs row: %w", err)
	}

	return id.String(), nil
}

func (r *WebhookLogRepository) UpdateResult(ctx context.Context, id string, status string, httpStatus int, errMessage string) error {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid webhook log id: %w", err)
	}
	now := time.Now()
	updates := map[string]interface{}{
		"status":        status,
		"response_code": httpStatus,
		"processed_at":  now,
		"updated_at":    now,
	}
	if errMessage != "" {
		updates["error_message"] = errMessage
	}
	if err := r.db.Conn(ctx).Table("webhook_logs").Where("id = ?", parsed).Updates(updates).Error; err != nil {
		return fmt.Errorf("update webhook_logs row: %w", err)
	}
	return nil
}
