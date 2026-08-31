package app

import (
	"context"
	"errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/tours/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/tours/internal/domain/entities"
	tourerrors "github.com/secamc93/probability/back/central/services/modules/tours/internal/domain/errors"
	"github.com/secamc93/probability/back/central/shared/log"
)

type repoFake struct {
	listado    []entities.TourProgress
	guardado   *entities.TourProgress
	masivos    []entities.TourProgress
	borrados   []string
	borroTodos bool
	errGuardar error
}

func (r *repoFake) ListByUser(_ context.Context, _, _ uint) ([]entities.TourProgress, error) {
	return r.listado, nil
}

func (r *repoFake) Upsert(_ context.Context, progress entities.TourProgress) (*entities.TourProgress, error) {
	if r.errGuardar != nil {
		return nil, r.errGuardar
	}
	r.guardado = &progress
	return &progress, nil
}

func (r *repoFake) UpsertMany(_ context.Context, items []entities.TourProgress) error {
	if r.errGuardar != nil {
		return r.errGuardar
	}
	r.masivos = append(r.masivos, items...)
	return nil
}

func (r *repoFake) Delete(_ context.Context, _, _ uint, tourKey string) error {
	r.borrados = append(r.borrados, tourKey)
	return nil
}

func (r *repoFake) DeleteAll(_ context.Context, _, _ uint) error {
	r.borroTodos = true
	return nil
}

func nuevoUseCase(repo *repoFake) IUseCase {
	return New(repo, log.New())
}

func TestSaveProgressRechazaEntradasInvalidas(t *testing.T) {
	casos := []struct {
		nombre   string
		entrada  dtos.SaveProgressInput
		esperado error
	}{
		{
			nombre:   "sin usuario",
			entrada:  dtos.SaveProgressInput{TourKey: "orders", Version: 1, Status: entities.StatusCompleted},
			esperado: tourerrors.ErrUserRequired,
		},
		{
			nombre:   "sin tour_key",
			entrada:  dtos.SaveProgressInput{UserID: 1, Version: 1, Status: entities.StatusCompleted},
			esperado: tourerrors.ErrTourKeyRequired,
		},
		{
			nombre:   "status invalido",
			entrada:  dtos.SaveProgressInput{UserID: 1, TourKey: "orders", Version: 1, Status: "raro"},
			esperado: tourerrors.ErrInvalidStatus,
		},
		{
			nombre:   "version cero",
			entrada:  dtos.SaveProgressInput{UserID: 1, TourKey: "orders", Version: 0, Status: entities.StatusCompleted},
			esperado: tourerrors.ErrInvalidVersion,
		},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			repo := &repoFake{}
			_, err := nuevoUseCase(repo).SaveProgress(context.Background(), caso.entrada)

			if !errors.Is(err, caso.esperado) {
				t.Fatalf("esperaba %v, obtuvo %v", caso.esperado, err)
			}
			if repo.guardado != nil {
				t.Fatal("no debio guardar nada con una entrada invalida")
			}
		})
	}
}

func TestSaveProgressMarcaCompletedAtSoloAlCerrar(t *testing.T) {
	casos := []struct {
		status      string
		esperaFecha bool
	}{
		{entities.StatusInProgress, false},
		{entities.StatusPending, false},
		{entities.StatusCompleted, true},
		{entities.StatusSkipped, true},
	}

	for _, caso := range casos {
		t.Run(caso.status, func(t *testing.T) {
			repo := &repoFake{}
			resultado, err := nuevoUseCase(repo).SaveProgress(context.Background(), dtos.SaveProgressInput{
				UserID:  7,
				TourKey: "orders",
				Version: 2,
				Status:  caso.status,
			})
			if err != nil {
				t.Fatalf("error inesperado: %v", err)
			}

			tieneFecha := resultado.CompletedAt != nil
			if tieneFecha != caso.esperaFecha {
				t.Fatalf("status %s: esperaba completed_at presente=%v, obtuvo %v", caso.status, caso.esperaFecha, tieneFecha)
			}
		})
	}
}

func TestSaveProgressNormalizaStepIndexNegativo(t *testing.T) {
	repo := &repoFake{}
	resultado, err := nuevoUseCase(repo).SaveProgress(context.Background(), dtos.SaveProgressInput{
		UserID:    7,
		TourKey:   "orders",
		Version:   1,
		Status:    entities.StatusInProgress,
		StepIndex: -4,
	})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if resultado.StepIndex != 0 {
		t.Fatalf("esperaba step_index 0, obtuvo %d", resultado.StepIndex)
	}
}

func TestSkipAllMarcaTodosComoSkipped(t *testing.T) {
	repo := &repoFake{}
	err := nuevoUseCase(repo).SkipAll(context.Background(), dtos.SkipAllInput{
		UserID:     7,
		BusinessID: 26,
		Tours: []dtos.SkipAllTour{
			{TourKey: "home", Version: 1},
			{TourKey: "orders", Version: 1},
			{TourKey: "products", Version: 2},
		},
	})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}

	if len(repo.masivos) != 3 {
		t.Fatalf("esperaba 3 registros, obtuvo %d", len(repo.masivos))
	}
	for _, item := range repo.masivos {
		if item.Status != entities.StatusSkipped {
			t.Fatalf("tour %s quedo con status %s", item.TourKey, item.Status)
		}
		if item.CompletedAt == nil {
			t.Fatalf("tour %s quedo sin completed_at", item.TourKey)
		}
		if item.BusinessID != 26 {
			t.Fatalf("tour %s quedo con business_id %d", item.TourKey, item.BusinessID)
		}
	}
}

func TestSkipAllRechazaVersionInvalida(t *testing.T) {
	repo := &repoFake{}
	err := nuevoUseCase(repo).SkipAll(context.Background(), dtos.SkipAllInput{
		UserID: 7,
		Tours:  []dtos.SkipAllTour{{TourKey: "home", Version: 0}},
	})

	if !errors.Is(err, tourerrors.ErrInvalidVersion) {
		t.Fatalf("esperaba ErrInvalidVersion, obtuvo %v", err)
	}
	if len(repo.masivos) != 0 {
		t.Fatal("no debio guardar nada")
	}
}

func TestSkipAllSinToursFalla(t *testing.T) {
	repo := &repoFake{}
	err := nuevoUseCase(repo).SkipAll(context.Background(), dtos.SkipAllInput{UserID: 7})

	if !errors.Is(err, tourerrors.ErrTourKeyRequired) {
		t.Fatalf("esperaba ErrTourKeyRequired, obtuvo %v", err)
	}
}

func TestResetTourExigeUsuarioYClave(t *testing.T) {
	repo := &repoFake{}
	uc := nuevoUseCase(repo)

	if err := uc.ResetTour(context.Background(), dtos.ResetInput{TourKey: "orders"}); !errors.Is(err, tourerrors.ErrUserRequired) {
		t.Fatalf("esperaba ErrUserRequired, obtuvo %v", err)
	}
	if err := uc.ResetTour(context.Background(), dtos.ResetInput{UserID: 7}); !errors.Is(err, tourerrors.ErrTourKeyRequired) {
		t.Fatalf("esperaba ErrTourKeyRequired, obtuvo %v", err)
	}
	if len(repo.borrados) != 0 {
		t.Fatal("no debio borrar nada")
	}

	if err := uc.ResetTour(context.Background(), dtos.ResetInput{UserID: 7, TourKey: "orders"}); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if len(repo.borrados) != 1 || repo.borrados[0] != "orders" {
		t.Fatalf("esperaba borrar orders, obtuvo %v", repo.borrados)
	}
}

func TestListProgressExigeUsuario(t *testing.T) {
	repo := &repoFake{}
	_, err := nuevoUseCase(repo).ListProgress(context.Background(), dtos.ListProgressInput{})

	if !errors.Is(err, tourerrors.ErrUserRequired) {
		t.Fatalf("esperaba ErrUserRequired, obtuvo %v", err)
	}
}

func TestListProgressDevuelveLoDelRepositorio(t *testing.T) {
	repo := &repoFake{listado: []entities.TourProgress{
		{TourKey: "home", Status: entities.StatusCompleted},
		{TourKey: "orders", Status: entities.StatusSkipped},
	}}

	resultado, err := nuevoUseCase(repo).ListProgress(context.Background(), dtos.ListProgressInput{UserID: 7})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if len(resultado.Items) != 2 {
		t.Fatalf("esperaba 2 items, obtuvo %d", len(resultado.Items))
	}
}
