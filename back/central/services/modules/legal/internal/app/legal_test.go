package app

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/legal/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/legal/internal/domain/entities"
	legalerrors "github.com/secamc93/probability/back/central/services/modules/legal/internal/domain/errors"
	"github.com/secamc93/probability/back/central/shared/log"
)

type repoFake struct {
	activos     []entities.LegalDocument
	porID       []entities.LegalDocument
	aceptados   map[uint]bool
	esPlatforma bool
	guardadas   []entities.LegalAcceptance
	errGuardar  error
}

func (r *repoFake) GetActiveDocuments(_ context.Context) ([]entities.LegalDocument, error) {
	return r.activos, nil
}

func (r *repoFake) GetDocumentsByIDs(_ context.Context, ids []uint) ([]entities.LegalDocument, error) {
	if r.porID != nil {
		return r.porID, nil
	}
	encontrados := make([]entities.LegalDocument, 0, len(ids))
	for _, id := range ids {
		for _, doc := range r.activos {
			if doc.ID == id {
				encontrados = append(encontrados, doc)
			}
		}
	}
	return encontrados, nil
}

func (r *repoFake) GetAcceptedDocumentIDs(_ context.Context, _ uint) (map[uint]bool, error) {
	if r.aceptados == nil {
		return map[uint]bool{}, nil
	}
	return r.aceptados, nil
}

func (r *repoFake) SaveAcceptances(_ context.Context, aceptaciones []entities.LegalAcceptance) error {
	if r.errGuardar != nil {
		return r.errGuardar
	}
	r.guardadas = append(r.guardadas, aceptaciones...)
	return nil
}

func (r *repoFake) IsPlatformUser(_ context.Context, _ uint) (bool, error) {
	return r.esPlatforma, nil
}

func documentosVigentes() []entities.LegalDocument {
	return []entities.LegalDocument{
		{ID: 1, Code: "terms_of_service", Version: "1.0", Title: "Terminos", SHA256: "aaa", IsActive: true, EffectiveFrom: time.Now().Add(-time.Hour)},
		{ID: 2, Code: "privacy_policy", Version: "1.0", Title: "Politica", SHA256: "bbb", IsActive: true, EffectiveFrom: time.Now().Add(-time.Hour)},
	}
}

func nuevoUseCase(repo *repoFake) IUseCase {
	return New(repo, log.New())
}

func TestPendientesDevuelveLosNoAceptados(t *testing.T) {
	repo := &repoFake{activos: documentosVigentes(), aceptados: map[uint]bool{1: true}}
	uc := nuevoUseCase(repo)

	res, err := uc.GetPendingDocuments(context.Background(), 10)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if !res.RequiresAcceptance {
		t.Fatal("deberia requerir aceptacion")
	}
	if len(res.Documents) != 1 || res.Documents[0].ID != 2 {
		t.Fatalf("documentos pendientes = %v", res.Documents)
	}
}

func TestPendientesVacioCuandoYaAceptoTodo(t *testing.T) {
	repo := &repoFake{activos: documentosVigentes(), aceptados: map[uint]bool{1: true, 2: true}}
	res, err := nuevoUseCase(repo).GetPendingDocuments(context.Background(), 10)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if res.RequiresAcceptance || len(res.Documents) != 0 {
		t.Fatalf("no deberia requerir aceptacion: %+v", res)
	}
}

func TestSuperAdminNoAcepta(t *testing.T) {
	repo := &repoFake{activos: documentosVigentes(), esPlatforma: true}
	res, err := nuevoUseCase(repo).GetPendingDocuments(context.Background(), 1)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if res.RequiresAcceptance || len(res.Documents) != 0 {
		t.Fatalf("el super admin no debe ver el modal: %+v", res)
	}
}

func TestPendientesSinUsuario(t *testing.T) {
	_, err := nuevoUseCase(&repoFake{}).GetPendingDocuments(context.Background(), 0)
	if !errors.Is(err, legalerrors.ErrUserRequired) {
		t.Fatalf("error = %v, se esperaba ErrUserRequired", err)
	}
}

func TestAceptarGuardaVersionYHash(t *testing.T) {
	repo := &repoFake{activos: documentosVigentes()}
	businessID := uint(46)

	res, err := nuevoUseCase(repo).AcceptDocuments(context.Background(), dtos.AcceptDocumentsInput{
		UserID:      10,
		BusinessID:  &businessID,
		DocumentIDs: []uint{1, 2},
		IPAddress:   "190.1.2.3",
		UserAgent:   "Mozilla/5.0",
	})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if len(res.DocumentIDs) != 2 {
		t.Fatalf("document_ids = %v", res.DocumentIDs)
	}
	if len(repo.guardadas) != 2 {
		t.Fatalf("aceptaciones guardadas = %d", len(repo.guardadas))
	}

	primera := repo.guardadas[0]
	if primera.DocumentSHA256 != "aaa" || primera.DocumentVersion != "1.0" || primera.DocumentCode != "terms_of_service" {
		t.Errorf("la aceptacion no guardo la identidad del documento: %+v", primera)
	}
	if primera.IPAddress != "190.1.2.3" || primera.UserAgent != "Mozilla/5.0" {
		t.Errorf("falta la evidencia de origen: %+v", primera)
	}
	if primera.Method != entities.AcceptanceMethodWebModal {
		t.Errorf("method = %s", primera.Method)
	}
	if primera.BusinessID == nil || *primera.BusinessID != 46 {
		t.Errorf("business_id = %v", primera.BusinessID)
	}
	if primera.AcceptedAt.IsZero() {
		t.Error("accepted_at vacio")
	}
}

func TestAceptarRechazaDocumentoInexistente(t *testing.T) {
	repo := &repoFake{activos: documentosVigentes()}
	_, err := nuevoUseCase(repo).AcceptDocuments(context.Background(), dtos.AcceptDocumentsInput{
		UserID:      10,
		DocumentIDs: []uint{1, 99},
	})
	if !errors.Is(err, legalerrors.ErrDocumentNotAvailable) {
		t.Fatalf("error = %v, se esperaba ErrDocumentNotAvailable", err)
	}
	if len(repo.guardadas) != 0 {
		t.Error("no debio guardar nada")
	}
}

func TestAceptarRechazaDocumentoInactivo(t *testing.T) {
	inactivo := []entities.LegalDocument{{ID: 3, Code: "terms_of_service", Version: "0.9", IsActive: false}}
	repo := &repoFake{activos: documentosVigentes(), porID: inactivo}

	_, err := nuevoUseCase(repo).AcceptDocuments(context.Background(), dtos.AcceptDocumentsInput{
		UserID:      10,
		DocumentIDs: []uint{3},
	})
	if !errors.Is(err, legalerrors.ErrDocumentNotAvailable) {
		t.Fatalf("error = %v, se esperaba ErrDocumentNotAvailable", err)
	}
}

func TestAceptarSinDocumentos(t *testing.T) {
	_, err := nuevoUseCase(&repoFake{}).AcceptDocuments(context.Background(), dtos.AcceptDocumentsInput{UserID: 10})
	if !errors.Is(err, legalerrors.ErrNoDocumentsSelected) {
		t.Fatalf("error = %v, se esperaba ErrNoDocumentsSelected", err)
	}
}
