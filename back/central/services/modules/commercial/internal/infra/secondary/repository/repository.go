package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/commercial/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/commercial/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) ListProspects(ctx context.Context, filters dtos.ListProspectsFilters) ([]entities.Prospect, int64, error) {
	base := r.db.Conn(ctx).Model(&models.CommercialProspect{})

	if filters.Platform != "" {
		base = base.Where("platform = ?", filters.Platform)
	}
	if filters.Status != "" {
		base = base.Where("status = ?", filters.Status)
	}
	if filters.City != "" {
		base = base.Where("initcap(unaccent(TRIM(city))) = ?", filters.City)
	}
	if filters.HasCOD != nil {
		base = base.Where("has_cod = ?", *filters.HasCOD)
	}
	if filters.MinScore != nil {
		base = base.Where("score >= ?", *filters.MinScore)
	}
	if filters.MaxScore != nil {
		base = base.Where("score <= ?", *filters.MaxScore)
	}
	if filters.Search != "" {
		like := "%" + filters.Search + "%"
		base = base.Where("domain ILIKE ? OR business_name ILIKE ? OR instagram_handle ILIKE ?", like, like, like)
	}

	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (filters.Page - 1) * filters.PageSize

	var rows []models.CommercialProspect
	if err := base.
		Order("score DESC, id DESC").
		Offset(offset).
		Limit(filters.PageSize).
		Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	prospects := make([]entities.Prospect, 0, len(rows))
	for i := range rows {
		prospects = append(prospects, modelToProspect(&rows[i]))
	}

	return prospects, total, nil
}

func (r *Repository) UpdateSeen(ctx context.Context, id uint, seen bool) error {
	return r.db.Conn(ctx).Model(&models.CommercialProspect{}).Where("id = ?", id).Update("seen", seen).Error
}

func (r *Repository) GetStats(ctx context.Context) (*entities.ProspectStats, error) {
	conn := r.db.Conn(ctx)
	stats := &entities.ProspectStats{}

	if err := conn.Model(&models.CommercialProspect{}).Count(&stats.Total).Error; err != nil {
		return nil, err
	}

	var platformRows []struct {
		Platform string
		Count    int64
	}
	if err := conn.Model(&models.CommercialProspect{}).
		Select("platform, COUNT(*) as count").
		Group("platform").
		Order("count DESC").
		Scan(&platformRows).Error; err != nil {
		return nil, err
	}
	for _, row := range platformRows {
		stats.ByPlatform = append(stats.ByPlatform, entities.CountByKey{Key: row.Platform, Count: row.Count})
	}

	var statusRows []struct {
		Status string
		Count  int64
	}
	if err := conn.Model(&models.CommercialProspect{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Order("count DESC").
		Scan(&statusRows).Error; err != nil {
		return nil, err
	}
	for _, row := range statusRows {
		stats.ByStatus = append(stats.ByStatus, entities.CountByKey{Key: row.Status, Count: row.Count})
	}

	var cityRows []struct {
		City  string
		Count int64
	}
	if err := conn.Model(&models.CommercialProspect{}).
		Select("initcap(unaccent(TRIM(city))) as city, COUNT(*) as count").
		Where("city <> ''").
		Group("initcap(unaccent(TRIM(city)))").
		Order("count DESC").
		Limit(10).
		Scan(&cityRows).Error; err != nil {
		return nil, err
	}
	for _, row := range cityRows {
		stats.ByCity = append(stats.ByCity, entities.CountByKey{Key: row.City, Count: row.Count})
	}

	buckets := []entities.ScoreBucket{
		{RangeLabel: "0-20", MinScore: 0, MaxScore: 20},
		{RangeLabel: "20-40", MinScore: 20, MaxScore: 40},
		{RangeLabel: "40-60", MinScore: 40, MaxScore: 60},
		{RangeLabel: "60-80", MinScore: 60, MaxScore: 80},
		{RangeLabel: "80-100", MinScore: 80, MaxScore: 100},
	}
	for i := range buckets {
		var count int64
		q := conn.Model(&models.CommercialProspect{}).Where("score >= ? AND score < ?", buckets[i].MinScore, buckets[i].MaxScore)
		if buckets[i].MaxScore == 100 {
			q = conn.Model(&models.CommercialProspect{}).Where("score >= ? AND score <= ?", buckets[i].MinScore, buckets[i].MaxScore)
		}
		if err := q.Count(&count).Error; err != nil {
			return nil, err
		}
		buckets[i].Count = count
	}
	stats.ScoreBuckets = buckets

	if err := conn.Model(&models.CommercialProspect{}).Where("has_cod = ?", true).Count(&stats.HasCODCount).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(&models.CommercialProspect{}).Where("uses_probability_channel = ?", true).Count(&stats.MultiChannel).Error; err != nil {
		return nil, err
	}

	var avg struct{ Avg float64 }
	if err := conn.Model(&models.CommercialProspect{}).Select("COALESCE(AVG(score), 0) as avg").Scan(&avg).Error; err != nil {
		return nil, err
	}
	stats.AverageScore = avg.Avg

	var top models.CommercialProspect
	if err := conn.Model(&models.CommercialProspect{}).Order("score DESC").Limit(1).Find(&top).Error; err != nil {
		return nil, err
	}
	stats.TopScoreDomain = top.Domain
	stats.TopScore = top.Score

	return stats, nil
}
