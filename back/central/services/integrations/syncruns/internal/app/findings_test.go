package app

import (
	"context"
	"errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/syncruns/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/syncruns/internal/mocks"
	"github.com/secamc93/probability/back/central/shared/log"
)

func nuevoCasoHallazgos(channels []domain.ChannelSummary, cross domain.CrossChannelCounts) domain.IUseCase {
	repo := &mocks.RepositoryMock{
		ChannelSummariesFn: func(context.Context, uint) ([]domain.ChannelSummary, error) { return channels, nil },
		CrossChannelFn:     func(context.Context, uint) (domain.CrossChannelCounts, error) { return cross, nil },
	}
	return New(repo, log.New())
}

func nuevoCasoMatriz() (domain.IUseCase, *mocks.RepositoryMock) {
	repo := &mocks.RepositoryMock{}
	return New(repo, log.New()), repo
}

func buscar(report *domain.FindingsReport, code string) *domain.Finding {
	for i := range report.Findings {
		if report.Findings[i].Code == code {
			return &report.Findings[i]
		}
	}
	return nil
}

func TestSinProblemasNoReportaHallazgos(t *testing.T) {
	uc := nuevoCasoHallazgos([]domain.ChannelSummary{{IntegrationID: 1, IntegrationName: "Meli", Matched: 100}}, domain.CrossChannelCounts{})

	r, err := uc.Findings(context.Background(), 46)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if len(r.Findings) != 0 || r.Total != 0 {
		t.Fatalf("un negocio sano no debe mostrar nada: %+v", r.Findings)
	}
}

func TestSumaElMismoProblemaEntreCanales(t *testing.T) {
	uc := nuevoCasoHallazgos([]domain.ChannelSummary{
		{IntegrationID: 1, IntegrationName: "Mercado Libre", ChannelNoSKU: 320},
		{IntegrationID: 2, IntegrationName: "Woocommerce", ChannelNoSKU: 5},
	}, domain.CrossChannelCounts{})

	r, _ := uc.Findings(context.Background(), 46)
	f := buscar(r, domain.FindingChannelNoSKU)
	if f == nil || f.Count != 325 {
		t.Fatalf("debia sumar los dos canales, obtuve %+v", f)
	}
	if len(f.Channels) != 2 {
		t.Fatalf("debia decir en que canales pasa, obtuve %v", f.Channels)
	}
}

func TestNoNombraElCanalQueNoTieneEseProblema(t *testing.T) {
	uc := nuevoCasoHallazgos([]domain.ChannelSummary{
		{IntegrationID: 1, IntegrationName: "Mercado Libre", SKUTypo: 34},
		{IntegrationID: 2, IntegrationName: "Woocommerce"},
	}, domain.CrossChannelCounts{})

	f := buscar(mustFindings(t, uc), domain.FindingSKUTypo)
	if f == nil || len(f.Channels) != 1 || f.Channels[0] != "Mercado Libre" {
		t.Fatalf("solo debe nombrar al canal afectado, obtuve %+v", f)
	}
}

func TestLosGravesVanPrimero(t *testing.T) {
	uc := nuevoCasoHallazgos([]domain.ChannelSummary{
		{IntegrationID: 1, IntegrationName: "Meli", SKUTypo: 500, ChannelNoSKU: 1},
	}, domain.CrossChannelCounts{})

	r := mustFindings(t, uc)
	if r.Findings[0].Severity != domain.SeverityError {
		t.Fatalf("un error con 1 caso pesa mas que un aviso con 500, obtuve %+v", r.Findings[0])
	}
}

func TestEntreIgualesMandaLaCantidad(t *testing.T) {
	uc := nuevoCasoHallazgos([]domain.ChannelSummary{
		{IntegrationID: 1, IntegrationName: "Meli", ChannelNoSKU: 10, SKUChanged: 99},
	}, domain.CrossChannelCounts{})

	r := mustFindings(t, uc)
	if r.Findings[0].Code != domain.FindingSKUChanged {
		t.Fatalf("a igual gravedad manda el volumen, obtuve %q", r.Findings[0].Code)
	}
}

func TestReportaLoQueSeVendeSinExistir(t *testing.T) {
	uc := nuevoCasoHallazgos(nil, domain.CrossChannelCounts{SoldNotOwned: 7})

	f := buscar(mustFindings(t, uc), domain.FindingSoldNotOwned)
	if f == nil || f.Count != 7 || f.Severity != domain.SeverityError {
		t.Fatalf("vender algo que no existe es un error, no un aviso: %+v", f)
	}
}

func TestReportaElStockQueNadieOfrece(t *testing.T) {
	uc := nuevoCasoHallazgos(nil, domain.CrossChannelCounts{NotPublished: 120})

	f := buscar(mustFindings(t, uc), domain.FindingNotPublished)
	if f == nil || f.Count != 120 {
		t.Fatalf("obtuve %+v", f)
	}
}

func TestElDesbalanceEsInformativoNoUnError(t *testing.T) {
	uc := nuevoCasoHallazgos(nil, domain.CrossChannelCounts{Imbalance: 40})

	f := buscar(mustFindings(t, uc), domain.FindingImbalance)
	if f == nil || f.Severity != domain.SeverityInfo {
		t.Fatalf("no publicar en todos los canales puede ser a proposito: %+v", f)
	}
}

func TestElTotalSumaTodosLosHallazgos(t *testing.T) {
	uc := nuevoCasoHallazgos([]domain.ChannelSummary{
		{IntegrationID: 1, IntegrationName: "Meli", ChannelNoSKU: 320, SKUTypo: 34, NotAssociated: 645},
	}, domain.CrossChannelCounts{SoldNotOwned: 7, NotPublished: 120, Imbalance: 40})

	r := mustFindings(t, uc)
	if r.Total != 320+34+645+7+120+40 {
		t.Fatalf("obtuve %d", r.Total)
	}
}

func TestDevuelveElResumenPorCanal(t *testing.T) {
	uc := nuevoCasoHallazgos([]domain.ChannelSummary{
		{IntegrationID: 254, IntegrationName: "Mercadolibre", ChannelCode: "Mercado Libre", Matched: 274},
	}, domain.CrossChannelCounts{})

	r := mustFindings(t, uc)
	if len(r.Channels) != 1 || r.Channels[0].IntegrationID != 254 {
		t.Fatalf("el panel necesita el detalle por canal: %+v", r.Channels)
	}
}

func TestPropagaElErrorDelResumen(t *testing.T) {
	repo := &mocks.RepositoryMock{
		ChannelSummariesFn: func(context.Context, uint) ([]domain.ChannelSummary, error) {
			return nil, errors.New("base caida")
		},
	}
	if _, err := New(repo, log.New()).Findings(context.Background(), 46); err == nil {
		t.Fatal("debia propagar")
	}
}

func TestPropagaElErrorDelCruce(t *testing.T) {
	repo := &mocks.RepositoryMock{
		CrossChannelFn: func(context.Context, uint) (domain.CrossChannelCounts, error) {
			return domain.CrossChannelCounts{}, errors.New("base caida")
		},
	}
	if _, err := New(repo, log.New()).Findings(context.Background(), 46); err == nil {
		t.Fatal("debia propagar")
	}
}

func mustFindings(t *testing.T, uc domain.IUseCase) *domain.FindingsReport {
	t.Helper()
	r, err := uc.Findings(context.Background(), 46)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	return r
}

func TestMatchMatrixExigeNegocio(t *testing.T) {
	uc, _ := nuevoCasoMatriz()
	if _, err := uc.MatchMatrix(context.Background(), domain.MatrixQuery{}); err == nil {
		t.Fatal("sin business_id deberia fallar")
	}
}

func TestMatchMatrixNormalizaPaginacion(t *testing.T) {
	uc, repo := nuevoCasoMatriz()
	var visto domain.MatrixQuery
	repo.MatchMatrixFn = func(_ context.Context, q domain.MatrixQuery) (*domain.MatrixPage, error) {
		visto = q
		return &domain.MatrixPage{Page: q.Page, PageSize: q.PageSize}, nil
	}

	if _, err := uc.MatchMatrix(context.Background(), domain.MatrixQuery{BusinessID: 46, Page: 0, PageSize: 9999}); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if visto.Page != 1 || visto.PageSize != 50 {
		t.Fatalf("no normalizo la paginacion: %+v", visto)
	}
}

func TestMatchMatrixPropagaError(t *testing.T) {
	uc, repo := nuevoCasoMatriz()
	repo.MatchMatrixFn = func(_ context.Context, _ domain.MatrixQuery) (*domain.MatrixPage, error) {
		return nil, errors.New("boom")
	}
	if _, err := uc.MatchMatrix(context.Background(), domain.MatrixQuery{BusinessID: 46}); err == nil {
		t.Fatal("deberia propagar el error del repositorio")
	}
}
