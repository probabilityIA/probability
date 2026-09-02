package app

import (
	"context"
	stderrors "errors"
	"testing"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/entities"
	dom "github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdate_LimpiarSprint_NoConsultaElSprint(t *testing.T) {
	consultado := false
	repo := &mocks.RepositoryMock{
		FindSprintNameFn: func(ctx context.Context, sprintID uint) (string, bool, error) {
			consultado = true
			return "Sprint 1", true, nil
		},
	}

	_, err := newTicketsUseCase(repo, nil).Update(context.Background(),
		dtos.UpdateTicketDTO{ID: 1, ClearSprint: true})

	require.NoError(t, err)
	assert.False(t, consultado)
	require.Len(t, repo.Updates, 1)
	assert.Contains(t, repo.Updates[0], "sprint_id")
	assert.Nil(t, repo.Updates[0]["sprint_id"])
}

func TestUpdate_LimpiarSprintGanaSobreEnviarlo(t *testing.T) {
	repo := &mocks.RepositoryMock{}

	_, err := newTicketsUseCase(repo, nil).Update(context.Background(),
		dtos.UpdateTicketDTO{ID: 1, ClearSprint: true, SprintID: uintPtr(9)})

	require.NoError(t, err)
	require.Len(t, repo.Updates, 1)
	assert.Nil(t, repo.Updates[0]["sprint_id"],
		"si se pide limpiar el sprint se limpia aunque venga un sprint en el mismo request")
}

func TestUpdate_SprintExistente_SeGuarda(t *testing.T) {
	var consultado uint
	repo := &mocks.RepositoryMock{
		FindSprintNameFn: func(ctx context.Context, sprintID uint) (string, bool, error) {
			consultado = sprintID
			return "Sprint 1", true, nil
		},
	}

	_, err := newTicketsUseCase(repo, nil).Update(context.Background(),
		dtos.UpdateTicketDTO{ID: 1, SprintID: uintPtr(9)})

	require.NoError(t, err)
	assert.Equal(t, uint(9), consultado)
	assert.Equal(t, uint(9), repo.Updates[0]["sprint_id"])
}

func TestUpdate_SprintInexistente_Rechaza(t *testing.T) {
	repo := &mocks.RepositoryMock{
		FindSprintNameFn: func(ctx context.Context, sprintID uint) (string, bool, error) {
			return "", false, nil
		},
	}

	_, err := newTicketsUseCase(repo, nil).Update(context.Background(),
		dtos.UpdateTicketDTO{ID: 1, SprintID: uintPtr(404)})

	assert.ErrorIs(t, err, dom.ErrSprintNotFound)
	assert.Empty(t, repo.Updates)
}

func TestUpdate_ErrorAlBuscarSprint_SePropaga(t *testing.T) {
	dbErr := stderrors.New("db caida")
	repo := &mocks.RepositoryMock{
		FindSprintNameFn: func(ctx context.Context, sprintID uint) (string, bool, error) {
			return "", false, dbErr
		},
	}

	_, err := newTicketsUseCase(repo, nil).Update(context.Background(),
		dtos.UpdateTicketDTO{ID: 1, SprintID: uintPtr(9)})

	assert.ErrorIs(t, err, dbErr)
	assert.Empty(t, repo.Updates)
}

func TestUpdate_ErrorAlValidarResponsable_SePropaga(t *testing.T) {
	dbErr := stderrors.New("db caida")
	repo := &mocks.RepositoryMock{
		UserExistsFn: func(ctx context.Context, userID uint) (bool, error) { return false, dbErr },
	}

	_, err := newTicketsUseCase(repo, nil).Update(context.Background(),
		dtos.UpdateTicketDTO{ID: 1, AssignedToID: uintPtr(7)})

	assert.ErrorIs(t, err, dbErr)
	assert.Empty(t, repo.Updates)
}

func TestUpdate_ResponsableExistente_SeGuarda(t *testing.T) {
	repo := &mocks.RepositoryMock{}

	_, err := newTicketsUseCase(repo, nil).Update(context.Background(),
		dtos.UpdateTicketDTO{ID: 1, AssignedToID: uintPtr(7)})

	require.NoError(t, err)
	require.Len(t, repo.Updates, 1)
	assert.Equal(t, uint(7), repo.Updates[0]["assigned_to_id"])
}

func TestUpdate_DescripcionYFechaLimiteValidas_SeGuardan(t *testing.T) {
	repo := &mocks.RepositoryMock{}

	_, err := newTicketsUseCase(repo, nil).Update(context.Background(), dtos.UpdateTicketDTO{
		ID: 1, Description: strPtr("  el listado sigue vacio  "), DueDate: strPtr("2026-12-31"),
	})

	require.NoError(t, err)
	require.Len(t, repo.Updates, 1)
	assert.Equal(t, "el listado sigue vacio", repo.Updates[0]["description"])
	assert.Equal(t, time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC), repo.Updates[0]["due_date"])
}

func TestUpdate_FechaLimiteVacia_NoSeEscribe(t *testing.T) {
	repo := &mocks.RepositoryMock{}

	_, err := newTicketsUseCase(repo, nil).Update(context.Background(),
		dtos.UpdateTicketDTO{ID: 1, DueDate: strPtr("")})

	require.NoError(t, err)
	assert.Empty(t, repo.Updates, "una fecha vacia no es una orden de limpiarla: para eso esta ClearDueDate")
}

func TestUpdate_SinCambios_DevuelveElTicketTalCual(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, id uint) (*entities.Ticket, error) {
			return &entities.Ticket{ID: id, Title: "sin tocar"}, nil
		},
	}

	got, err := newTicketsUseCase(repo, nil).Update(context.Background(), dtos.UpdateTicketDTO{ID: 5})

	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "sin tocar", got.Title)
	assert.Empty(t, repo.Updates)
}

func TestUpdate_SinCambiosYTicketInexistente_PropagaElError(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, id uint) (*entities.Ticket, error) {
			return nil, dom.ErrTicketNotFound
		},
	}

	got, err := newTicketsUseCase(repo, nil).Update(context.Background(), dtos.UpdateTicketDTO{ID: 404})

	assert.ErrorIs(t, err, dom.ErrTicketNotFound)
	assert.Nil(t, got)
}
