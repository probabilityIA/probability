package app

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/entities"
	"github.com/secamc93/probability/back/central/shared/log"
)

func TestGetRecommendation_TrasladaOrigenYDestinoAlProveedor(t *testing.T) {
	provider := &recommendationFake{}
	uc := New(provider, nil, nil, nil, nil, nil, nil, log.New())

	rec, err := uc.GetRecommendation(context.Background(), "BOGOTA", "MEDELLIN")

	require.NoError(t, err)
	require.Len(t, provider.calls, 1)
	assert.Equal(t, [2]string{"BOGOTA", "MEDELLIN"}, provider.calls[0])
	assert.Equal(t, "SERVIENTREGA", rec.RecommendedCarrier)
}

func TestGetRecommendation_PropagaElErrorDelProveedor(t *testing.T) {
	fallo := errors.New("openrouter 429: rate limit")
	provider := &recommendationFake{fn: func(string, string) (*entities.Recommendation, error) { return nil, fallo }}
	uc := New(provider, nil, nil, nil, nil, nil, nil, log.New())

	rec, err := uc.GetRecommendation(context.Background(), "BOGOTA", "MEDELLIN")

	assert.ErrorIs(t, err, fallo)
	assert.Nil(t, rec)
}

func TestGetRecommendation_DevuelveLasCotizacionesTalComoLlegan(t *testing.T) {
	provider := &recommendationFake{fn: func(string, string) (*entities.Recommendation, error) {
		return &entities.Recommendation{
			RecommendedCarrier: "INTERRAPIDISIMO",
			Alternatives:       []string{"SERVIENTREGA"},
			Quotations:         []entities.Quotation{{Carrier: "INTERRAPIDISIMO", EstimatedCost: 15000, EstimatedDeliveryDays: 3}},
		}, nil
	}}
	uc := New(provider, nil, nil, nil, nil, nil, nil, log.New())

	rec, err := uc.GetRecommendation(context.Background(), "PEREIRA", "ARMENIA")

	require.NoError(t, err)
	require.Len(t, rec.Quotations, 1)
	assert.EqualValues(t, 15000, rec.Quotations[0].EstimatedCost)
	assert.Equal(t, []string{"SERVIENTREGA"}, rec.Alternatives)
}
