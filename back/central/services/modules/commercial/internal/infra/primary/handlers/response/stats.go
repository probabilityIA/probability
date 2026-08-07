package response

import "github.com/secamc93/probability/back/central/services/modules/commercial/internal/domain/entities"

type CountByKeyResponse struct {
	Key   string `json:"key"`
	Count int64  `json:"count"`
}

type ScoreBucketResponse struct {
	RangeLabel string `json:"range_label"`
	Count      int64  `json:"count"`
}

type StatsResponse struct {
	Total          int64                 `json:"total"`
	ByPlatform     []CountByKeyResponse  `json:"by_platform"`
	ByStatus       []CountByKeyResponse  `json:"by_status"`
	ByCity         []CountByKeyResponse  `json:"by_city"`
	ScoreBuckets   []ScoreBucketResponse `json:"score_buckets"`
	HasCODCount    int64                 `json:"has_cod_count"`
	MultiChannel   int64                 `json:"multi_channel_count"`
	AverageScore   float64               `json:"average_score"`
	TopScoreDomain string                `json:"top_score_domain"`
	TopScore       float64               `json:"top_score"`
}

func mapCounts(items []entities.CountByKey) []CountByKeyResponse {
	out := make([]CountByKeyResponse, 0, len(items))
	for _, item := range items {
		out = append(out, CountByKeyResponse{Key: item.Key, Count: item.Count})
	}
	return out
}

func FromStats(s *entities.ProspectStats) StatsResponse {
	buckets := make([]ScoreBucketResponse, 0, len(s.ScoreBuckets))
	for _, b := range s.ScoreBuckets {
		buckets = append(buckets, ScoreBucketResponse{RangeLabel: b.RangeLabel, Count: b.Count})
	}

	return StatsResponse{
		Total:          s.Total,
		ByPlatform:     mapCounts(s.ByPlatform),
		ByStatus:       mapCounts(s.ByStatus),
		ByCity:         mapCounts(s.ByCity),
		ScoreBuckets:   buckets,
		HasCODCount:    s.HasCODCount,
		MultiChannel:   s.MultiChannel,
		AverageScore:   s.AverageScore,
		TopScoreDomain: s.TopScoreDomain,
		TopScore:       s.TopScore,
	}
}
