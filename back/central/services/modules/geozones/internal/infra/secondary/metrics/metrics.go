package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	ProbabilityCacheOperations = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "geozones_probability_cache_operations_total",
		Help: "Operaciones sobre la cache Redis de probabilidad geozones (hit/miss/set/invalidate)",
	}, []string{"result"})

	ProbabilityQueryDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "geozones_probability_query_duration_seconds",
		Help:    "Duracion del SELECT a geozone_carrier_stats para una orden",
		Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2, 5},
	})

	AggregateRefreshDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "geozones_aggregate_refresh_duration_seconds",
		Help:    "Duracion del refresh completo de geozone_carrier_stats",
		Buckets: []float64{0.5, 1, 2, 5, 10, 30, 60, 120, 300},
	})

	AggregateRefreshLastSuccess = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "geozones_aggregate_refresh_last_success_timestamp_seconds",
		Help: "Unix timestamp del ultimo refresh exitoso de geozone_carrier_stats",
	})
)
