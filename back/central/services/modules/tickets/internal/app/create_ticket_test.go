package app

import (
	"context"
	stderrors "errors"
	"testing"

	dom "github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreate_SprintInexistente_NoCrea(t *testing.T) {
	repo := &mocks.RepositoryMock{
		FindSprintNameFn: func(ctx context.Context, sprintID uint) (string, bool, error) {
			return "", false, nil
		},
	}
	dto := ticketValido()
	dto.SprintID = uintPtr(404)

	_, err := newTicketsUseCase(repo, nil).Create(context.Background(), dto)

	assert.ErrorIs(t, err, dom.ErrSprintNotFound)
	assert.Nil(t, repo.CreatedTicket)
}

func TestCreate_ErrorAlBuscarSprint_SePropaga(t *testing.T) {
	dbErr := stderrors.New("db caida")
	repo := &mocks.RepositoryMock{
		FindSprintNameFn: func(ctx context.Context, sprintID uint) (string, bool, error) {
			return "", false, dbErr
		},
	}
	dto := ticketValido()
	dto.SprintID = uintPtr(9)

	_, err := newTicketsUseCase(repo, nil).Create(context.Background(), dto)

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, repo.CreatedTicket)
}

func TestCreate_SprintExistente_SeGuardaEnElTicket(t *testing.T) {
	var consultado uint
	repo := &mocks.RepositoryMock{
		FindSprintNameFn: func(ctx context.Context, sprintID uint) (string, bool, error) {
			consultado = sprintID
			return "Sprint 1", true, nil
		},
	}
	dto := ticketValido()
	dto.SprintID = uintPtr(9)

	_, err := newTicketsUseCase(repo, nil).Create(context.Background(), dto)

	require.NoError(t, err)
	assert.Equal(t, uint(9), consultado)
	require.NotNil(t, repo.CreatedTicket)
	require.NotNil(t, repo.CreatedTicket.SprintID)
	assert.Equal(t, uint(9), *repo.CreatedTicket.SprintID)
}

func TestCreate_SinSprint_NoConsultaSprints(t *testing.T) {
	consultado := false
	repo := &mocks.RepositoryMock{
		FindSprintNameFn: func(ctx context.Context, sprintID uint) (string, bool, error) {
			consultado = true
			return "Sprint 1", true, nil
		},
	}

	_, err := newTicketsUseCase(repo, nil).Create(context.Background(), ticketValido())

	require.NoError(t, err)
	assert.False(t, consultado)
}

func TestCreate_DevuelveElTicketReleidoDelRepositorio(t *testing.T) {
	repo := &mocks.RepositoryMock{}

	got, err := newTicketsUseCase(repo, nil).Create(context.Background(), ticketValido())

	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, uint(1), got.ID,
		"tras crear se relee el ticket para devolverlo con sus relaciones cargadas")
}
