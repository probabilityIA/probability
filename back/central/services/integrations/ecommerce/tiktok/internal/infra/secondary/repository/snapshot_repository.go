package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiktok/internal/domain"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/migration/shared/models"
)

type SnapshotRepository struct {
	db     db.IDatabase
	logger log.ILogger
}

func New(database db.IDatabase, logger log.ILogger) domain.ISnapshotRepository {
	return &SnapshotRepository{
		db:     database,
		logger: logger.WithModule("tiktok.snapshot_repository"),
	}
}

func (r *SnapshotRepository) Save(ctx context.Context, snapshot *domain.ShopSnapshot) error {
	record := &models.TikTokShopSnapshot{
		BusinessID:           snapshot.BusinessID,
		IntegrationID:        snapshot.IntegrationID,
		ShopID:               snapshot.ShopID,
		ShopName:             snapshot.ShopName,
		ShopRating:           snapshot.ShopRating,
		ShopReviewCount:      snapshot.ShopReviewCount,
		ItemSoldCount:        snapshot.ItemSoldCount,
		ShopPerformanceValue: snapshot.ShopPerformanceValue,
	}

	if err := r.db.Conn(ctx).Create(record).Error; err != nil {
		return err
	}

	snapshot.ID = record.ID
	snapshot.CreatedAt = record.CreatedAt
	return nil
}

func (r *SnapshotRepository) ListByIntegration(ctx context.Context, businessID, integrationID uint, page, pageSize int) ([]domain.ShopSnapshot, int64, error) {
	var records []models.TikTokShopSnapshot
	var total int64

	query := r.db.Conn(ctx).Model(&models.TikTokShopSnapshot{}).
		Where("business_id = ? AND integration_id = ?", businessID, integrationID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&records).Error; err != nil {
		return nil, 0, err
	}

	snapshots := make([]domain.ShopSnapshot, 0, len(records))
	for _, rec := range records {
		snapshots = append(snapshots, domain.ShopSnapshot{
			ID:                   rec.ID,
			BusinessID:           rec.BusinessID,
			IntegrationID:        rec.IntegrationID,
			ShopID:               rec.ShopID,
			ShopName:             rec.ShopName,
			ShopRating:           rec.ShopRating,
			ShopReviewCount:      rec.ShopReviewCount,
			ItemSoldCount:        rec.ItemSoldCount,
			ShopPerformanceValue: rec.ShopPerformanceValue,
			CreatedAt:            rec.CreatedAt,
		})
	}

	return snapshots, total, nil
}
