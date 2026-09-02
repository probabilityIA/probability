package app

import (
	"context"
	stderrors "errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/entities"
	dom "github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChangeSprint_SinSprint_RetiraDelSprintSinConsultarlo(t *testing.T) {
	consultado := false
	repo := &mocks.RepositoryMock{
		FindSprintNameFn: func(ctx context.Context, sprintID uint) (string, bool, error) {
			consultado = true
			return "Sprint 1", true, nil
		},
	}

	_, err := newTicketsUseCase(repo, nil).ChangeSprint(context.Background(),
		dtos.ChangeSprintDTO{TicketID: 1, SprintID: nil, ChangedByID: 42})

	require.NoError(t, err)
	assert.False(t, consultado, "quitar el sprint no requiere validar ningun sprint")
	require.Len(t, repo.Updates, 1)
	assert.Contains(t, repo.Updates[0], "sprint_id")
	assert.Nil(t, repo.Updates[0]["sprint_id"])
	require.Len(t, repo.History, 1)
	assert.Contains(t, repo.History[0].Note, "retirado del sprint")
	assert.Equal(t, uint(42), repo.History[0].ChangedByID)
}

func TestChangeSprint_ConSprintExistente_GuardaElIDYNombraElSprintEnLaNota(t *testing.T) {
	var consultado uint
	repo := &mocks.RepositoryMock{
		FindSprintNameFn: func(ctx context.Context, sprintID uint) (string, bool, error) {
			consultado = sprintID
			return "Sprint de octubre", true, nil
		},
	}

	_, err := newTicketsUseCase(repo, nil).ChangeSprint(context.Background(),
		dtos.ChangeSprintDTO{TicketID: 1, SprintID: uintPtr(9), ChangedByID: 42})

	require.NoError(t, err)
	assert.Equal(t, uint(9), consultado)
	require.Len(t, repo.Updates, 1)
	assert.Equal(t, uint(9), repo.Updates[0]["sprint_id"])
	require.Len(t, repo.History, 1)
	assert.Contains(t, repo.History[0].Note, "Sprint de octubre")
}

func TestChangeSprint_SprintInexistente_NoMueveElTicket(t *testing.T) {
	repo := &mocks.RepositoryMock{
		FindSprintNameFn: func(ctx context.Context, sprintID uint) (string, bool, error) {
			return "", false, nil
		},
	}

	_, err := newTicketsUseCase(repo, nil).ChangeSprint(context.Background(),
		dtos.ChangeSprintDTO{TicketID: 1, SprintID: uintPtr(404)})

	assert.ErrorIs(t, err, dom.ErrSprintNotFound)
	assert.Empty(t, repo.Updates, "no se mueve un ticket a un sprint que no existe")
	assert.Empty(t, repo.History)
}

func TestChangeSprint_NoCambiaElEstadoDelTicket(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, id uint) (*entities.Ticket, error) {
			return &entities.Ticket{ID: id, Status: "in_development"}, nil
		},
	}

	_, err := newTicketsUseCase(repo, nil).ChangeSprint(context.Background(),
		dtos.ChangeSprintDTO{TicketID: 1, SprintID: uintPtr(9)})

	require.NoError(t, err)
	assert.NotContains(t, repo.Updates[0], "status")
	require.Len(t, repo.History, 1)
	assert.Equal(t, "in_development", repo.History[0].From)
	assert.Equal(t, repo.History[0].From, repo.History[0].To,
		"mover de sprint no es una transicion de estado")
}

func TestChangeSprint_ErroresDelRepo_SePropagan(t *testing.T) {
	dbErr := stderrors.New("db caida")

	t.Run("al leer el ticket", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			GetByIDFn: func(ctx context.Context, id uint) (*entities.Ticket, error) { return nil, dbErr },
		}
		_, err := newTicketsUseCase(repo, nil).ChangeSprint(context.Background(),
			dtos.ChangeSprintDTO{TicketID: 1})
		assert.ErrorIs(t, err, dbErr)
		assert.Empty(t, repo.Updates)
	})

	t.Run("al buscar el sprint", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			FindSprintNameFn: func(ctx context.Context, sprintID uint) (string, bool, error) {
				return "", false, dbErr
			},
		}
		_, err := newTicketsUseCase(repo, nil).ChangeSprint(context.Background(),
			dtos.ChangeSprintDTO{TicketID: 1, SprintID: uintPtr(9)})
		assert.ErrorIs(t, err, dbErr)
		assert.Empty(t, repo.Updates)
	})

	t.Run("al actualizar", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			UpdateFn: func(ctx context.Context, id uint, u map[string]any) (*entities.Ticket, error) {
				return nil, dbErr
			},
		}
		_, err := newTicketsUseCase(repo, nil).ChangeSprint(context.Background(),
			dtos.ChangeSprintDTO{TicketID: 1})
		assert.ErrorIs(t, err, dbErr)
		assert.Empty(t, repo.History, "sin escritura no se registra historial")
	})
}
