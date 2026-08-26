package repository

import (
	"context"
	"strconv"

	"github.com/secamc93/probability/back/central/services/modules/codreport/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
)

const (
	emailLogModule        = "cod_report"
	emailLogReferenceType = "cod_payment_cut"
	emailLogEventType     = "cod_cut_report"
	emailLogProvider      = "resend"
)

func (r *Repository) SaveCutEmailLogs(ctx context.Context, logs []entities.CutEmailLog) error {
	if len(logs) == 0 {
		return nil
	}
	rows := make([]models.EmailLog, len(logs))
	for i := range logs {
		row := models.EmailLog{
			CreatedAt:     logs[i].SentAt,
			BusinessID:    logs[i].BusinessID,
			Module:        emailLogModule,
			ReferenceType: emailLogReferenceType,
			ReferenceID:   strconv.FormatUint(uint64(logs[i].CutID), 10),
			To:            logs[i].Recipient,
			Subject:       logs[i].Subject,
			EventType:     emailLogEventType,
			Status:        logs[i].Status,
			Provider:      emailLogProvider,
			SentBy:        logs[i].SentBy,
			SentByName:    logs[i].SentByName,
		}
		if logs[i].ErrorMessage != "" {
			msg := logs[i].ErrorMessage
			row.ErrorMessage = &msg
		}
		rows[i] = row
	}
	return r.db.Conn(ctx).Create(&rows).Error
}

func (r *Repository) CutEmailLogs(ctx context.Context, businessID uint, cutID uint) ([]entities.CutEmailLog, error) {
	var rows []models.EmailLog
	err := r.db.Conn(ctx).
		Where("business_id = ? AND reference_type = ? AND reference_id = ?", businessID, emailLogReferenceType, strconv.FormatUint(uint64(cutID), 10)).
		Order("created_at DESC").
		Limit(100).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make([]entities.CutEmailLog, len(rows))
	for i := range rows {
		errMsg := ""
		if rows[i].ErrorMessage != nil {
			errMsg = *rows[i].ErrorMessage
		}
		out[i] = entities.CutEmailLog{
			ID:           rows[i].ID.String(),
			CutID:        cutID,
			BusinessID:   rows[i].BusinessID,
			Recipient:    rows[i].To,
			Subject:      rows[i].Subject,
			Status:       rows[i].Status,
			ErrorMessage: errMsg,
			SentBy:       rows[i].SentBy,
			SentByName:   rows[i].SentByName,
			SentAt:       rows[i].CreatedAt,
		}
	}
	return out, nil
}
