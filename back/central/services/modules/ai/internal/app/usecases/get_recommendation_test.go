package usecases

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain"
)

type aiServiceMock struct {
	GetRecommendationFn func(req domain.RecommendationRequest) (*domain.AIRecommendation, error)

	Requests []domain.RecommendationRequest
}

var _ domain.AIService = (*aiServiceMock)(nil)

func (m *aiServiceMock) GetRecommendation(req domain.RecommendationRequest) (*domain.AIRecommendation, error) {
	m.Requests = append(m.Requests, req)
	if m.GetRecommendationFn != nil {
		return m.GetRecommendationFn(req)
	}
	return &domain.AIRecommendation{RecommendedCarrier: "SERVIENTREGA"}, nil
}

func TestExecute_TrasladaOrigenYDestinoAlServicioDeIA(t *testing.T) {
	svc := &aiServiceMock{}

	rec, err := NewGetRecommendationUseCase(svc).Execute("BOGOTA", "MEDELLIN")

	require.NoError(t, err)
	require.Len(t, svc.Requests, 1)
	assert.Equal(t, "BOGOTA", svc.Requests[0].Origin)
	assert.Equal(t, "MEDELLIN", svc.Requests[0].Destination)
	assert.Equal(t, "SERVIENTREGA", rec.RecommendedCarrier)
}

func TestExecute_NoRestringeLasTransportadorasCandidatas(t *testing.T) {
	svc := &aiServiceMock{}

	_, err := NewGetRecommendationUseCase(svc).Execute("BOGOTA", "CALI")

	require.NoError(t, err)
	assert.Empty(t, svc.Requests[0].Carriers,
		"el caso de uso nunca llena Carriers, asi que el modelo puede recomendar una transportadora que el negocio no tiene contratada")
}

func TestExecute_OrigenODestinoVacioIgualLlegaAlModelo(t *testing.T) {
	casos := []struct {
		nombre          string
		origen, destino string
	}{
		{"sin origen", "", "MEDELLIN"},
		{"sin destino", "BOGOTA", ""},
		{"ambos vacios", "", ""},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			svc := &aiServiceMock{}

			_, err := NewGetRecommendationUseCase(svc).Execute(caso.origen, caso.destino)

			require.NoError(t, err)
			require.Len(t, svc.Requests, 1,
				"no hay validacion previa: una ciudad vacia se gasta una llamada al modelo y devuelve una recomendacion sin base")
		})
	}
}

func TestExecute_PropagaElErrorDelServicioDeIA(t *testing.T) {
	fallo := errors.New("openrouter 429: rate limit")
	svc := &aiServiceMock{
		GetRecommendationFn: func(req domain.RecommendationRequest) (*domain.AIRecommendation, error) {
			return nil, fallo
		},
	}

	rec, err := NewGetRecommendationUseCase(svc).Execute("BOGOTA", "MEDELLIN")

	assert.ErrorIs(t, err, fallo)
	assert.Nil(t, rec)
}

func TestExecute_DevuelveLasCotizacionesTalComoLlegan(t *testing.T) {
	svc := &aiServiceMock{
		GetRecommendationFn: func(req domain.RecommendationRequest) (*domain.AIRecommendation, error) {
			return &domain.AIRecommendation{
				RecommendedCarrier: "INTERRAPIDISIMO",
				Reasoning:          "amplia cobertura en el eje cafetero",
				Alternatives:       []string{"SERVIENTREGA"},
				Quotations: []domain.Quotation{
					{Carrier: "INTERRAPIDISIMO", EstimatedCost: 15000, EstimatedDeliveryDays: 3},
				},
			}, nil
		},
	}

	rec, err := NewGetRecommendationUseCase(svc).Execute("PEREIRA", "ARMENIA")

	require.NoError(t, err)
	require.Len(t, rec.Quotations, 1)
	assert.EqualValues(t, 15000, rec.Quotations[0].EstimatedCost,
		"los costos son estimados por el modelo, no cotizaciones reales de la transportadora")
	assert.Equal(t, []string{"SERVIENTREGA"}, rec.Alternatives)
}
