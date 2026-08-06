package app

import (
	"context"
	stderrors "errors"
	"strings"
	"testing"

	"github.com/secamc93/probability/back/central/services/auth/actions/internal/domain"
	"github.com/secamc93/probability/back/central/services/auth/actions/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newActionsUseCase(repo *mocks.RepositoryMock) IUseCaseAction {
	if repo == nil {
		repo = &mocks.RepositoryMock{}
	}
	return New(repo, mocks.NewSilentLogger())
}

func TestCreateAction_NombreObligatorio(t *testing.T) {
	for _, nombre := range []string{"", "   ", "\t"} {
		repo := &mocks.RepositoryMock{}

		got, err := newActionsUseCase(repo).CreateAction(context.Background(),
			domain.CreateActionDTO{Name: nombre})

		require.Error(t, err, "nombre=%q", nombre)
		assert.Contains(t, err.Error(), "nombre del action es obligatorio")
		assert.Nil(t, got)
		assert.Empty(t, repo.CreatedActions)
	}
}

func TestCreateAction_LimitesDeLongitud(t *testing.T) {
	t.Run("nombre en el limite de 20 se acepta", func(t *testing.T) {
		_, err := newActionsUseCase(nil).CreateAction(context.Background(),
			domain.CreateActionDTO{Name: strings.Repeat("a", 20)})
		require.NoError(t, err)
	})

	t.Run("nombre de 21 se rechaza", func(t *testing.T) {
		repo := &mocks.RepositoryMock{}
		_, err := newActionsUseCase(repo).CreateAction(context.Background(),
			domain.CreateActionDTO{Name: strings.Repeat("a", 21)})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no puede exceder 20")
		assert.Empty(t, repo.CreatedActions)
	})

	t.Run("descripcion en el limite de 255 se acepta", func(t *testing.T) {
		_, err := newActionsUseCase(nil).CreateAction(context.Background(),
			domain.CreateActionDTO{Name: "read", Description: strings.Repeat("d", 255)})
		require.NoError(t, err)
	})

	t.Run("descripcion de 256 se rechaza", func(t *testing.T) {
		repo := &mocks.RepositoryMock{}
		_, err := newActionsUseCase(repo).CreateAction(context.Background(),
			domain.CreateActionDTO{Name: "read", Description: strings.Repeat("d", 256)})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no puede exceder 255")
		assert.Empty(t, repo.CreatedActions)
	})
}

func TestCreateAction_NombreDuplicado_Rechaza(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetActionByNameFn: func(ctx context.Context, name string) (*domain.Action, error) {
			return &domain.Action{ID: 1, Name: name}, nil
		},
	}

	got, err := newActionsUseCase(repo).CreateAction(context.Background(),
		domain.CreateActionDTO{Name: "read"})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "ya existe un action con el nombre 'read'")
	assert.Nil(t, got)
	assert.Empty(t, repo.CreatedActions,
		"dos acciones con el mismo nombre harian ambiguo el permiso")
}

func TestCreateAction_ErrorAlBuscarDuplicado_NoBloquea(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetActionByNameFn: func(ctx context.Context, name string) (*domain.Action, error) {
			return nil, stderrors.New("no encontrado")
		},
	}

	_, err := newActionsUseCase(repo).CreateAction(context.Background(),
		domain.CreateActionDTO{Name: "read"})

	require.NoError(t, err, "que la busqueda falle se interpreta como que no existe")
	assert.Len(t, repo.CreatedActions, 1)
}

func TestCreateAction_RecortaEspacios(t *testing.T) {
	repo := &mocks.RepositoryMock{}

	_, err := newActionsUseCase(repo).CreateAction(context.Background(),
		domain.CreateActionDTO{Name: "  read  ", Description: "  leer  "})

	require.NoError(t, err)
	require.Len(t, repo.CreatedActions, 1)
	assert.Equal(t, "read", repo.CreatedActions[0].Name)
	assert.Equal(t, "leer", repo.CreatedActions[0].Description)
}

func TestCreateAction_ErroresDelRepo_SePropagan(t *testing.T) {
	dbErr := stderrors.New("constraint violation")

	t.Run("al crear", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			CreateActionFn: func(ctx context.Context, a domain.Action) (uint, error) { return 0, dbErr },
		}
		_, err := newActionsUseCase(repo).CreateAction(context.Background(),
			domain.CreateActionDTO{Name: "read"})
		assert.ErrorIs(t, err, dbErr)
	})

	t.Run("al recargar", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			GetActionByIDFn: func(ctx context.Context, id uint) (*domain.Action, error) { return nil, dbErr },
		}
		_, err := newActionsUseCase(repo).CreateAction(context.Background(),
			domain.CreateActionDTO{Name: "read"})
		assert.ErrorIs(t, err, dbErr)
	})
}

func TestCreateAction_DevuelveElRecargado(t *testing.T) {
	repo := &mocks.RepositoryMock{
		CreateActionFn: func(ctx context.Context, a domain.Action) (uint, error) { return 9, nil },
		GetActionByIDFn: func(ctx context.Context, id uint) (*domain.Action, error) {
			return &domain.Action{ID: id, Name: "read", Description: "leer"}, nil
		},
	}

	got, err := newActionsUseCase(repo).CreateAction(context.Background(),
		domain.CreateActionDTO{Name: "read"})

	require.NoError(t, err)
	assert.Equal(t, uint(9), got.ID)
	assert.Equal(t, "leer", got.Description)
}

func TestGetActions_NormalizaPaginacion(t *testing.T) {
	casos := []struct {
		page, pageSize, wantPage, wantSize int
	}{
		{0, 10, 1, 10},
		{-3, 10, 1, 10},
		{1, 0, 1, 10},
		{1, -1, 1, 10},
		{2, 50, 2, 50},
	}

	for _, tc := range casos {
		var vistoPage, vistoSize int
		repo := &mocks.RepositoryMock{
			GetActionsFn: func(ctx context.Context, page, pageSize int, name string) ([]domain.Action, int64, error) {
				vistoPage, vistoSize = page, pageSize
				return nil, 0, nil
			},
		}

		got, err := newActionsUseCase(repo).GetActions(context.Background(), tc.page, tc.pageSize, "")

		require.NoError(t, err)
		assert.Equal(t, tc.wantPage, vistoPage)
		assert.Equal(t, tc.wantSize, vistoSize)
		assert.Equal(t, tc.wantPage, got.Page)
		assert.Equal(t, tc.wantSize, got.PageSize)
	}
}

func TestGetActions_CalculaTotalPagesConRedondeoHaciaArriba(t *testing.T) {
	casos := []struct {
		total    int64
		pageSize int
		want     int
	}{
		{0, 10, 0},
		{1, 10, 1},
		{10, 10, 1},
		{11, 10, 2},
		{95, 10, 10},
	}

	for _, tc := range casos {
		repo := &mocks.RepositoryMock{
			GetActionsFn: func(ctx context.Context, page, pageSize int, name string) ([]domain.Action, int64, error) {
				return nil, tc.total, nil
			},
		}

		got, err := newActionsUseCase(repo).GetActions(context.Background(), 1, tc.pageSize, "")

		require.NoError(t, err)
		assert.Equal(t, tc.want, got.TotalPages, "total=%d pageSize=%d", tc.total, tc.pageSize)
	}
}

func TestGetActions_PasaElFiltroDeNombre(t *testing.T) {
	var visto string
	repo := &mocks.RepositoryMock{
		GetActionsFn: func(ctx context.Context, page, pageSize int, name string) ([]domain.Action, int64, error) {
			visto = name
			return nil, 0, nil
		},
	}

	_, err := newActionsUseCase(repo).GetActions(context.Background(), 1, 10, "rea")

	require.NoError(t, err)
	assert.Equal(t, "rea", visto)
}

func TestGetActions_MapeaLosCampos(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetActionsFn: func(ctx context.Context, page, pageSize int, name string) ([]domain.Action, int64, error) {
			return []domain.Action{
				{ID: 1, Name: "read", Description: "leer"},
				{ID: 2, Name: "write", Description: "escribir"},
			}, 2, nil
		},
	}

	got, err := newActionsUseCase(repo).GetActions(context.Background(), 1, 10, "")

	require.NoError(t, err)
	require.Len(t, got.Actions, 2)
	assert.Equal(t, "read", got.Actions[0].Name)
	assert.Equal(t, "escribir", got.Actions[1].Description)
	assert.Equal(t, int64(2), got.Total)
}

func TestGetActions_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := stderrors.New("timeout")
	repo := &mocks.RepositoryMock{
		GetActionsFn: func(ctx context.Context, page, pageSize int, name string) ([]domain.Action, int64, error) {
			return nil, 0, dbErr
		},
	}

	got, err := newActionsUseCase(repo).GetActions(context.Background(), 1, 10, "")

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestGetActionByID_Inexistente_RetornaError(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetActionByIDFn: func(ctx context.Context, id uint) (*domain.Action, error) { return nil, nil },
	}

	got, err := newActionsUseCase(repo).GetActionByID(context.Background(), 404)

	require.Error(t, err)
	assert.Nil(t, got)
	assert.Contains(t, err.Error(), "404 no encontrado")
}

func TestGetActionByID_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := stderrors.New("timeout")
	repo := &mocks.RepositoryMock{
		GetActionByIDFn: func(ctx context.Context, id uint) (*domain.Action, error) { return nil, dbErr },
	}

	_, err := newActionsUseCase(repo).GetActionByID(context.Background(), 1)

	require.Error(t, err)
}

func TestGetActionByID_MapeaLosCampos(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetActionByIDFn: func(ctx context.Context, id uint) (*domain.Action, error) {
			return &domain.Action{ID: id, Name: "read", Description: "leer"}, nil
		},
	}

	got, err := newActionsUseCase(repo).GetActionByID(context.Background(), 5)

	require.NoError(t, err)
	assert.Equal(t, uint(5), got.ID)
	assert.Equal(t, "read", got.Name)
	assert.Equal(t, "leer", got.Description)
}

func TestUpdateAction_ValidaAntesDeLeer(t *testing.T) {
	leido := false
	repo := &mocks.RepositoryMock{
		GetActionByIDFn: func(ctx context.Context, id uint) (*domain.Action, error) {
			leido = true
			return &domain.Action{ID: id}, nil
		},
	}

	_, err := newActionsUseCase(repo).UpdateAction(context.Background(), 1,
		domain.UpdateActionDTO{Name: "  "})

	require.Error(t, err)
	assert.False(t, leido)
}

func TestUpdateAction_Inexistente_NoActualiza(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetActionByIDFn: func(ctx context.Context, id uint) (*domain.Action, error) { return nil, nil },
	}

	_, err := newActionsUseCase(repo).UpdateAction(context.Background(), 404,
		domain.UpdateActionDTO{Name: "read"})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "404 no encontrado")
	assert.Empty(t, repo.UpdatedActions)
}

func TestUpdateAction_NombreDeOtroAction_Rechaza(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetActionByIDFn: func(ctx context.Context, id uint) (*domain.Action, error) {
			return &domain.Action{ID: id, Name: "viejo"}, nil
		},
		GetActionByNameFn: func(ctx context.Context, name string) (*domain.Action, error) {
			return &domain.Action{ID: 99, Name: name}, nil
		},
	}

	_, err := newActionsUseCase(repo).UpdateAction(context.Background(), 1,
		domain.UpdateActionDTO{Name: "ocupado"})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "ya existe otro action con el nombre 'ocupado'")
	assert.Empty(t, repo.UpdatedActions)
}

func TestUpdateAction_SuPropioNombreNoCuentaComoDuplicado(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetActionByIDFn: func(ctx context.Context, id uint) (*domain.Action, error) {
			return &domain.Action{ID: id, Name: "viejo"}, nil
		},
		GetActionByNameFn: func(ctx context.Context, name string) (*domain.Action, error) {
			return &domain.Action{ID: 1, Name: name}, nil
		},
	}

	_, err := newActionsUseCase(repo).UpdateAction(context.Background(), 1,
		domain.UpdateActionDTO{Name: "nuevo"})

	require.NoError(t, err)
	assert.Len(t, repo.UpdatedActions, 1)
}

func TestUpdateAction_MismoNombreNoDisparaChequeo(t *testing.T) {
	chequeado := false
	repo := &mocks.RepositoryMock{
		GetActionByIDFn: func(ctx context.Context, id uint) (*domain.Action, error) {
			return &domain.Action{ID: id, Name: "read"}, nil
		},
		GetActionByNameFn: func(ctx context.Context, name string) (*domain.Action, error) {
			chequeado = true
			return nil, nil
		},
	}

	_, err := newActionsUseCase(repo).UpdateAction(context.Background(), 1,
		domain.UpdateActionDTO{Name: "  read  "})

	require.NoError(t, err)
	assert.False(t, chequeado)
}

func TestUpdateAction_RecortaEspacios(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetActionByIDFn: func(ctx context.Context, id uint) (*domain.Action, error) {
			return &domain.Action{ID: id, Name: "viejo"}, nil
		},
	}

	_, err := newActionsUseCase(repo).UpdateAction(context.Background(), 1,
		domain.UpdateActionDTO{Name: "  nuevo  ", Description: "  desc  "})

	require.NoError(t, err)
	require.Len(t, repo.UpdatedActions, 1)
	assert.Equal(t, "nuevo", repo.UpdatedActions[0].Name)
	assert.Equal(t, "desc", repo.UpdatedActions[0].Description)
	assert.Equal(t, uint(1), repo.UpdatedActions[0].ID)
}

func TestUpdateAction_LimitesDeLongitud(t *testing.T) {
	repo := &mocks.RepositoryMock{}

	_, err := newActionsUseCase(repo).UpdateAction(context.Background(), 1,
		domain.UpdateActionDTO{Name: strings.Repeat("a", 21)})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no puede exceder 20")

	_, err = newActionsUseCase(repo).UpdateAction(context.Background(), 1,
		domain.UpdateActionDTO{Name: "read", Description: strings.Repeat("d", 256)})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no puede exceder 255")

	assert.Empty(t, repo.UpdatedActions)
}

func TestUpdateAction_ErroresDelRepo_SePropagan(t *testing.T) {
	dbErr := stderrors.New("deadlock")

	t.Run("al actualizar", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			GetActionByIDFn: func(ctx context.Context, id uint) (*domain.Action, error) {
				return &domain.Action{ID: id, Name: "read"}, nil
			},
			UpdateActionFn: func(ctx context.Context, id uint, a domain.Action) (string, error) {
				return "", dbErr
			},
		}
		_, err := newActionsUseCase(repo).UpdateAction(context.Background(), 1,
			domain.UpdateActionDTO{Name: "read"})
		assert.ErrorIs(t, err, dbErr)
	})

	t.Run("al leer para actualizar", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			GetActionByIDFn: func(ctx context.Context, id uint) (*domain.Action, error) { return nil, dbErr },
		}
		_, err := newActionsUseCase(repo).UpdateAction(context.Background(), 1,
			domain.UpdateActionDTO{Name: "read"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no encontrado")
	})
}

func TestUpdateAction_ErrorAlRecargar_SePropaga(t *testing.T) {
	dbErr := stderrors.New("timeout")
	llamadas := 0
	repo := &mocks.RepositoryMock{
		GetActionByIDFn: func(ctx context.Context, id uint) (*domain.Action, error) {
			llamadas++
			if llamadas == 1 {
				return &domain.Action{ID: id, Name: "read"}, nil
			}
			return nil, dbErr
		},
	}

	_, err := newActionsUseCase(repo).UpdateAction(context.Background(), 1,
		domain.UpdateActionDTO{Name: "read"})

	require.Error(t, err)
}

func TestDeleteAction_Inexistente_NoBorra(t *testing.T) {
	borrado := false
	repo := &mocks.RepositoryMock{
		GetActionByIDFn: func(ctx context.Context, id uint) (*domain.Action, error) { return nil, nil },
		DeleteActionFn: func(ctx context.Context, id uint) (string, error) {
			borrado = true
			return "", nil
		},
	}

	got, err := newActionsUseCase(repo).DeleteAction(context.Background(), 404)

	require.Error(t, err)
	assert.Empty(t, got)
	assert.False(t, borrado)
}

func TestDeleteAction_ErrorAlVerificar_NoBorra(t *testing.T) {
	borrado := false
	repo := &mocks.RepositoryMock{
		GetActionByIDFn: func(ctx context.Context, id uint) (*domain.Action, error) {
			return nil, stderrors.New("timeout")
		},
		DeleteActionFn: func(ctx context.Context, id uint) (string, error) {
			borrado = true
			return "", nil
		},
	}

	_, err := newActionsUseCase(repo).DeleteAction(context.Background(), 1)

	require.Error(t, err)
	assert.False(t, borrado)
}

func TestDeleteAction_Exitoso_DevuelveElMensajeDelRepo(t *testing.T) {
	var borradoID uint
	repo := &mocks.RepositoryMock{
		GetActionByIDFn: func(ctx context.Context, id uint) (*domain.Action, error) {
			return &domain.Action{ID: id, Name: "read"}, nil
		},
		DeleteActionFn: func(ctx context.Context, id uint) (string, error) {
			borradoID = id
			return "action eliminado", nil
		},
	}

	got, err := newActionsUseCase(repo).DeleteAction(context.Background(), 5)

	require.NoError(t, err)
	assert.Equal(t, uint(5), borradoID)
	assert.Equal(t, "action eliminado", got)
}

func TestDeleteAction_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := stderrors.New("fk violation: hay permisos usando la accion")
	repo := &mocks.RepositoryMock{
		GetActionByIDFn: func(ctx context.Context, id uint) (*domain.Action, error) {
			return &domain.Action{ID: id, Name: "read"}, nil
		},
		DeleteActionFn: func(ctx context.Context, id uint) (string, error) { return "", dbErr },
	}

	got, err := newActionsUseCase(repo).DeleteAction(context.Background(), 5)

	assert.ErrorIs(t, err, dbErr)
	assert.Empty(t, got)
}
