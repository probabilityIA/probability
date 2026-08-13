package domain

import (
	"context"
	"time"
)

const (
	ModeFillEmpty = "fill_empty"
	ModeOverwrite = "overwrite"
)

const MaxProductsPerApply = 5000

type DataApplyRequest struct {
	BusinessID    uint
	IntegrationID uint
	Field         string
	Mode          string
	SKUs          []string
	UserID        uint
}

type DataApplyResult struct {
	BatchID  string    `json:"batch_id"`
	Field    string    `json:"field"`
	Applied  int       `json:"applied"`
	AppliedAt time.Time `json:"applied_at"`
}

type DataUndoResult struct {
	BatchID  string `json:"batch_id"`
	Reverted int    `json:"reverted"`
}

type DataBatch struct {
	BatchID       string    `json:"batch_id"`
	Field         string    `json:"field"`
	IntegrationID *uint     `json:"integration_id,omitempty"`
	Products      int       `json:"products"`
	CreatedAt     time.Time `json:"created_at"`
}

type ProgressFunc func(processed, total int)

type IDataApplyRepository interface {
	ApplyChannelField(ctx context.Context, req DataApplyRequest, batchID string, at time.Time) (int, error)
	ApplyChannelVariant(ctx context.Context, req DataApplyRequest, batchID string, at time.Time, progress ProgressFunc) (int, error)
	UndoBatch(ctx context.Context, businessID uint, batchID string, userID uint, at time.Time) (int, error)
	ListDataBatches(ctx context.Context, businessID uint, limit int) ([]DataBatch, error)
}
