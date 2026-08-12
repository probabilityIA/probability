package app

import (
	"context"
	"errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/syncruns/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/syncruns/internal/mocks"
	"github.com/secamc93/probability/back/central/shared/log"
)

func TestFindingItemsPasaLaConsultaAlRepositorio(t *testing.T) {
	var recibida domain.FindingItemsQuery
	repo := &mocks.RepositoryMock{
		FindingItemsFn: func(_ context.Context, q domain.FindingItemsQuery) (*domain.FindingItemsPage, error) {
			recibida = q
			return &domain.FindingItemsPage{Items: []domain.FindingItem{{SKU: "A-1"}}}, nil
		},
	}

	_, err := New(repo, log.New()).FindingItems(context.Background(), domain.FindingItemsQuery{
		BusinessID: 46, Code: domain.FindingSKUTypo, Search: "SDH", Page: 2, PageSize: 20,
	})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if recibida.Code != domain.FindingSKUTypo || recibida.Search != "SDH" || recibida.Page != 2 {
		t.Fatalf("la consulta llego alterada: %+v", recibida)
	}
}

func TestFindingItemsSinNegocioNoConsulta(t *testing.T) {
	llamado := false
	repo := &mocks.RepositoryMock{
		FindingItemsFn: func(context.Context, domain.FindingItemsQuery) (*domain.FindingItemsPage, error) {
			llamado = true
			return nil, nil
		},
	}

	page, err := New(repo, log.New()).FindingItems(context.Background(), domain.FindingItemsQuery{Code: domain.FindingSKUTypo})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if llamado {
		t.Fatal("sin negocio no hay nada que consultar")
	}
	if page == nil || len(page.Items) != 0 {
		t.Fatal("debia devolver una pagina vacia, no nil")
	}
}

func TestFindingItemsSinCodigoNoConsulta(t *testing.T) {
	llamado := false
	repo := &mocks.RepositoryMock{
		FindingItemsFn: func(context.Context, domain.FindingItemsQuery) (*domain.FindingItemsPage, error) {
			llamado = true
			return nil, nil
		},
	}

	if _, err := New(repo, log.New()).FindingItems(context.Background(), domain.FindingItemsQuery{BusinessID: 46}); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if llamado {
		t.Fatal("sin codigo de hallazgo no se sabe que listar")
	}
}

func TestFindingItemsNormalizaLaPaginacion(t *testing.T) {
	var recibida domain.FindingItemsQuery
	repo := &mocks.RepositoryMock{
		FindingItemsFn: func(_ context.Context, q domain.FindingItemsQuery) (*domain.FindingItemsPage, error) {
			recibida = q
			return &domain.FindingItemsPage{}, nil
		},
	}

	_, _ = New(repo, log.New()).FindingItems(context.Background(), domain.FindingItemsQuery{
		BusinessID: 46, Code: domain.FindingSKUTypo, Page: 0, PageSize: 9999,
	})
	if recibida.Page != 1 {
		t.Fatalf("la pagina debia caer en 1, obtuve %d", recibida.Page)
	}
	if recibida.PageSize != 50 {
		t.Fatalf("un tamano absurdo debia caer al por defecto, obtuve %d", recibida.PageSize)
	}
}

func TestFindingItemsPropagaElError(t *testing.T) {
	repo := &mocks.RepositoryMock{
		FindingItemsFn: func(context.Context, domain.FindingItemsQuery) (*domain.FindingItemsPage, error) {
			return nil, errors.New("hallazgo desconocido")
		},
	}
	if _, err := New(repo, log.New()).FindingItems(context.Background(), domain.FindingItemsQuery{
		BusinessID: 46, Code: "inventado",
	}); err == nil {
		t.Fatal("debia propagar")
	}
}

func TestOffsetSaltaLasPaginasAnteriores(t *testing.T) {
	q := domain.FindingItemsQuery{Page: 3, PageSize: 50}
	if q.Offset() != 100 {
		t.Fatalf("obtuve %d", q.Offset())
	}
	if (domain.FindingItemsQuery{Page: 0, PageSize: 50}).Offset() != 0 {
		t.Fatal("la primera pagina no salta nada")
	}
}
