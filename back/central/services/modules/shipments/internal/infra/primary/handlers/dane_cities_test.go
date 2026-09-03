package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/mocks"
)

func TestBuscadorDeCiudadesPasaTerminoDepartamentoYLimite(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var vistoState, vistoTerm string
	var vistoLimit int
	repo := &mocks.RepositoryMock{
		SearchDaneCitiesFn: func(_ context.Context, stateCode, term string, limit int) ([]domain.DaneItem, error) {
			vistoState, vistoTerm, vistoLimit = stateCode, term, limit
			return []domain.DaneItem{{Code: "13001", Name: "CARTAGENA DE INDIAS", StateName: "BOLIVAR"}}, nil
		},
	}

	h := &Handlers{uc: newTestUseCase(repo)}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/dane/1/cities?state=BOL&q=carta&limit=8", nil)

	h.buscarCiudadesDane(c)

	if vistoState != "BOL" || vistoTerm != "carta" || vistoLimit != 8 {
		t.Fatalf("parametros mal propagados: state=%q term=%q limit=%d", vistoState, vistoTerm, vistoLimit)
	}

	var body struct {
		Cities []daneItemResponse `json:"cities"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	if len(body.Cities) != 1 || body.Cities[0].Code != "13001" || body.Cities[0].State != "BOLIVAR" {
		t.Fatalf("respuesta inesperada: %+v", body.Cities)
	}
}

func TestBuscadorSinDepartamentoNiTerminoNoConsulta(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &mocks.RepositoryMock{
		SearchDaneCitiesFn: func(context.Context, string, string, int) ([]domain.DaneItem, error) {
			t.Fatal("sin state ni q no se debe consultar")
			return nil, nil
		},
	}

	h := &Handlers{uc: newTestUseCase(repo)}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/dane/1/cities", nil)

	h.buscarCiudadesDane(c)

	var body struct {
		Cities []daneItemResponse `json:"cities"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	if len(body.Cities) != 0 {
		t.Fatalf("se esperaba lista vacia, llego %+v", body.Cities)
	}
}
