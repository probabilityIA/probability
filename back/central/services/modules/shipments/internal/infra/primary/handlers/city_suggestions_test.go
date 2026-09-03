package handlers

import (
	"context"
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/app/usecases"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/mocks"
)

func newTestUseCase(repo domain.IRepository) *usecases.UseCases {
	return usecases.New(repo, nil, nil)
}

func TestCitySuggestionsCaeAlaBusquedaNacionalSiElDepartamentoNoDaNada(t *testing.T) {
	llamadas := make([]string, 0, 2)
	repo := &mocks.RepositoryMock{
		SearchDaneCitiesFn: func(_ context.Context, stateCode, term string, limit int) ([]domain.DaneItem, error) {
			llamadas = append(llamadas, stateCode)
			if stateCode != "" {
				return nil, nil
			}
			return []domain.DaneItem{{Code: "13001", Name: "CARTAGENA DE INDIAS", StateName: "BOLIVAR"}}, nil
		},
	}

	h := &Handlers{uc: newTestUseCase(repo)}
	got := h.citySuggestions(context.Background(), "Cartagena", "ANT")

	if len(llamadas) != 2 || llamadas[0] != "ANT" || llamadas[1] != "" {
		t.Fatalf("se esperaba buscar primero en el departamento y luego en todo el pais, llamadas: %v", llamadas)
	}
	if len(got) != 1 || got[0].Code != "13001" || got[0].State != "BOLIVAR" {
		t.Fatalf("sugerencia inesperada: %+v", got)
	}
}

func TestCitySuggestionsSinCiudadNoConsulta(t *testing.T) {
	repo := &mocks.RepositoryMock{
		SearchDaneCitiesFn: func(_ context.Context, _, _ string, _ int) ([]domain.DaneItem, error) {
			t.Fatal("sin ciudad no se debe consultar la base")
			return nil, nil
		},
	}

	h := &Handlers{uc: newTestUseCase(repo)}
	if got := h.citySuggestions(context.Background(), "   ", "ANT"); len(got) != 0 {
		t.Fatalf("se esperaban cero sugerencias, llegaron %d", len(got))
	}
}
