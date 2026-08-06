package app

import (
	"bytes"
	"context"
	stderrors "errors"
	"mime/multipart"
	"net/textproto"
	"strings"
	"testing"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/entities"
	dom "github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTicketsUseCase(repo *mocks.RepositoryMock, storage *mocks.StorageServiceMock) IUseCase {
	if repo == nil {
		repo = &mocks.RepositoryMock{}
	}
	if storage == nil {
		storage = &mocks.StorageServiceMock{}
	}
	return New(repo, storage, mocks.NewSilentLogger())
}

func uintPtr(v uint) *uint    { return &v }
func strPtr(v string) *string { return &v }

func ticketValido() dtos.CreateTicketDTO {
	return dtos.CreateTicketDTO{
		Title: "No carga el listado", Description: "Al entrar a ordenes queda en blanco",
		CreatedByID: 42,
	}
}

func TestCreate_TituloYDescripcionObligatorios(t *testing.T) {
	casos := []struct {
		nombre  string
		dto     dtos.CreateTicketDTO
		wantErr error
	}{
		{"sin titulo", dtos.CreateTicketDTO{Description: "d"}, dom.ErrTitleRequired},
		{"titulo en blanco", dtos.CreateTicketDTO{Title: "   ", Description: "d"}, dom.ErrTitleRequired},
		{"sin descripcion", dtos.CreateTicketDTO{Title: "t"}, dom.ErrDescriptionRequired},
		{"descripcion en blanco", dtos.CreateTicketDTO{Title: "t", Description: "  "}, dom.ErrDescriptionRequired},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			repo := &mocks.RepositoryMock{}

			got, err := newTicketsUseCase(repo, nil).Create(context.Background(), tc.dto)

			assert.ErrorIs(t, err, tc.wantErr)
			assert.Nil(t, got)
			assert.Nil(t, repo.CreatedTicket)
		})
	}
}

func TestCreate_ValoresPorDefecto(t *testing.T) {
	repo := &mocks.RepositoryMock{}

	_, err := newTicketsUseCase(repo, nil).Create(context.Background(), ticketValido())

	require.NoError(t, err)
	require.NotNil(t, repo.CreatedTicket)
	assert.Equal(t, "support", repo.CreatedTicket.Type)
	assert.Equal(t, "medium", repo.CreatedTicket.Priority)
	assert.Equal(t, "internal", repo.CreatedTicket.Source)
	assert.Equal(t, "soporte", repo.CreatedTicket.Area)
	assert.Equal(t, "open", repo.CreatedTicket.Status, "todo ticket nace abierto")
	assert.Empty(t, repo.CreatedTicket.Severity)
}

func TestCreate_NormalizaAMinusculasYRecorta(t *testing.T) {
	repo := &mocks.RepositoryMock{}
	dto := ticketValido()
	dto.Type = "  BUG  "
	dto.Priority = "  CRITICAL  "
	dto.Severity = "  HIGH  "
	dto.Area = "  DESARROLLO  "
	dto.Source = "  BUSINESS  "
	dto.Title = "  titulo  "
	dto.Description = "  descripcion  "

	_, err := newTicketsUseCase(repo, nil).Create(context.Background(), dto)

	require.NoError(t, err)
	assert.Equal(t, "bug", repo.CreatedTicket.Type)
	assert.Equal(t, "critical", repo.CreatedTicket.Priority)
	assert.Equal(t, "high", repo.CreatedTicket.Severity)
	assert.Equal(t, "desarrollo", repo.CreatedTicket.Area)
	assert.Equal(t, "business", repo.CreatedTicket.Source)
	assert.Equal(t, "titulo", repo.CreatedTicket.Title)
	assert.Equal(t, "descripcion", repo.CreatedTicket.Description)
}

func TestCreate_ValoresInvalidos_Rechaza(t *testing.T) {
	casos := []struct {
		nombre  string
		ajuste  func(*dtos.CreateTicketDTO)
		wantErr error
	}{
		{"tipo invalido", func(d *dtos.CreateTicketDTO) { d.Type = "inventado" }, dom.ErrInvalidType},
		{"prioridad invalida", func(d *dtos.CreateTicketDTO) { d.Priority = "urgentisimo" }, dom.ErrInvalidPriority},
		{"severidad invalida", func(d *dtos.CreateTicketDTO) { d.Severity = "catastrofica" }, dom.ErrInvalidSeverity},
		{"area invalida", func(d *dtos.CreateTicketDTO) { d.Area = "marketing" }, dom.ErrInvalidArea},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			repo := &mocks.RepositoryMock{}
			dto := ticketValido()
			tc.ajuste(&dto)

			_, err := newTicketsUseCase(repo, nil).Create(context.Background(), dto)

			assert.ErrorIs(t, err, tc.wantErr)
			assert.Nil(t, repo.CreatedTicket)
		})
	}
}

func TestCreate_OrigenInvalidoCaeAInternoSinFallar(t *testing.T) {
	repo := &mocks.RepositoryMock{}
	dto := ticketValido()
	dto.Source = "telepatia"

	_, err := newTicketsUseCase(repo, nil).Create(context.Background(), dto)

	require.NoError(t, err, "a diferencia del resto, un origen desconocido no rechaza el ticket")
	assert.Equal(t, "internal", repo.CreatedTicket.Source)
}

func TestCreate_AceptaTodosLosValoresDelCatalogo(t *testing.T) {
	for tipo := range validTypes {
		dto := ticketValido()
		dto.Type = tipo
		_, err := newTicketsUseCase(nil, nil).Create(context.Background(), dto)
		assert.NoError(t, err, "tipo %q del catalogo debe aceptarse", tipo)
	}

	for prioridad := range validPriorities {
		dto := ticketValido()
		dto.Priority = prioridad
		_, err := newTicketsUseCase(nil, nil).Create(context.Background(), dto)
		assert.NoError(t, err, "prioridad %q del catalogo debe aceptarse", prioridad)
	}

	for area := range validAreas {
		dto := ticketValido()
		dto.Area = area
		_, err := newTicketsUseCase(nil, nil).Create(context.Background(), dto)
		assert.NoError(t, err, "area %q del catalogo debe aceptarse", area)
	}

	for severidad := range validSeverities {
		dto := ticketValido()
		dto.Severity = severidad
		_, err := newTicketsUseCase(nil, nil).Create(context.Background(), dto)
		assert.NoError(t, err, "severidad %q del catalogo debe aceptarse", severidad)
	}
}

func TestCreate_ResponsableInexistente_NoCrea(t *testing.T) {
	repo := &mocks.RepositoryMock{
		UserExistsFn: func(ctx context.Context, userID uint) (bool, error) { return false, nil },
	}
	dto := ticketValido()
	dto.AssignedToID = uintPtr(404)

	_, err := newTicketsUseCase(repo, nil).Create(context.Background(), dto)

	assert.ErrorIs(t, err, dom.ErrAssigneeNotFound)
	assert.Nil(t, repo.CreatedTicket)
}

func TestCreate_SinResponsable_NoValidaUsuario(t *testing.T) {
	validado := false
	repo := &mocks.RepositoryMock{
		UserExistsFn: func(ctx context.Context, userID uint) (bool, error) {
			validado = true
			return true, nil
		},
	}

	_, err := newTicketsUseCase(repo, nil).Create(context.Background(), ticketValido())

	require.NoError(t, err)
	assert.False(t, validado)
}

func TestCreate_FechaLimiteInvalida_SeIgnoraSinFallar(t *testing.T) {
	casos := []string{"", "31-12-2026", "2026/12/31", "manana", "2026-13-45"}

	for _, fecha := range casos {
		t.Run(fecha, func(t *testing.T) {
			repo := &mocks.RepositoryMock{}
			dto := ticketValido()
			dto.DueDate = strPtr(fecha)

			_, err := newTicketsUseCase(repo, nil).Create(context.Background(), dto)

			require.NoError(t, err, "una fecha mal formada no debe tumbar la creacion del ticket")
			assert.Nil(t, repo.CreatedTicket.DueDate)
		})
	}
}

func TestCreate_FechaLimiteValida_SeGuarda(t *testing.T) {
	repo := &mocks.RepositoryMock{}
	dto := ticketValido()
	dto.DueDate = strPtr("2026-12-31")

	_, err := newTicketsUseCase(repo, nil).Create(context.Background(), dto)

	require.NoError(t, err)
	require.NotNil(t, repo.CreatedTicket.DueDate)
	assert.Equal(t, time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC), *repo.CreatedTicket.DueDate)
}

func TestCreate_AsignaElCodigoDelRepositorio(t *testing.T) {
	repo := &mocks.RepositoryMock{
		NextCodeFn: func(ctx context.Context) (string, error) { return "TKT-0099", nil },
	}

	_, err := newTicketsUseCase(repo, nil).Create(context.Background(), ticketValido())

	require.NoError(t, err)
	assert.Equal(t, "TKT-0099", repo.CreatedTicket.Code)
}

func TestCreate_RegistraElHistorialInicial(t *testing.T) {
	repo := &mocks.RepositoryMock{}

	_, err := newTicketsUseCase(repo, nil).Create(context.Background(), ticketValido())

	require.NoError(t, err)
	require.Len(t, repo.History, 1)
	assert.Empty(t, repo.History[0].From, "no hay estado anterior en la creacion")
	assert.Equal(t, "open", repo.History[0].To)
	assert.Equal(t, uint(42), repo.History[0].ChangedByID)
}

func TestCreate_ErroresDelRepo_SePropagan(t *testing.T) {
	dbErr := stderrors.New("db caida")

	t.Run("al pedir el codigo", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			NextCodeFn: func(ctx context.Context) (string, error) { return "", dbErr },
		}
		_, err := newTicketsUseCase(repo, nil).Create(context.Background(), ticketValido())
		assert.ErrorIs(t, err, dbErr)
		assert.Nil(t, repo.CreatedTicket)
	})

	t.Run("al validar el responsable", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			UserExistsFn: func(ctx context.Context, userID uint) (bool, error) { return false, dbErr },
		}
		dto := ticketValido()
		dto.AssignedToID = uintPtr(7)
		_, err := newTicketsUseCase(repo, nil).Create(context.Background(), dto)
		assert.ErrorIs(t, err, dbErr)
	})

	t.Run("al crear", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			CreateFn: func(ctx context.Context, tk *entities.Ticket) (*entities.Ticket, error) {
				return nil, dbErr
			},
		}
		_, err := newTicketsUseCase(repo, nil).Create(context.Background(), ticketValido())
		assert.ErrorIs(t, err, dbErr)
		assert.Empty(t, repo.History)
	})
}

func TestChangeStatus_EstadoInvalido_Rechaza(t *testing.T) {
	repo := &mocks.RepositoryMock{}

	_, err := newTicketsUseCase(repo, nil).ChangeStatus(context.Background(),
		dtos.ChangeStatusDTO{TicketID: 1, NewStatus: "inventado"})

	assert.ErrorIs(t, err, dom.ErrInvalidStatus)
	assert.Empty(t, repo.Updates)
}

func TestChangeStatus_AceptaTodosLosEstadosDelCatalogo(t *testing.T) {
	for estado := range validStatuses {
		repo := &mocks.RepositoryMock{
			GetByIDFn: func(ctx context.Context, id uint) (*entities.Ticket, error) {
				return &entities.Ticket{ID: id, Status: "open"}, nil
			},
		}
		_, err := newTicketsUseCase(repo, nil).ChangeStatus(context.Background(),
			dtos.ChangeStatusDTO{TicketID: 1, NewStatus: estado})
		assert.NoError(t, err, "estado %q del catalogo debe aceptarse", estado)
	}
}

func TestChangeStatus_MismoEstado_NoActualizaNiRegistraHistorial(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, id uint) (*entities.Ticket, error) {
			return &entities.Ticket{ID: id, Status: "in_review"}, nil
		},
	}

	got, err := newTicketsUseCase(repo, nil).ChangeStatus(context.Background(),
		dtos.ChangeStatusDTO{TicketID: 1, NewStatus: "IN_REVIEW"})

	require.NoError(t, err)
	assert.Equal(t, "in_review", got.Status)
	assert.Empty(t, repo.Updates, "no hay cambio real, no hay que escribir")
	assert.Empty(t, repo.History, "ni ensuciar el historial con un cambio que no ocurrio")
}

func TestChangeStatus_Resuelto_MarcaLaFechaDeResolucion(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, id uint) (*entities.Ticket, error) {
			return &entities.Ticket{ID: id, Status: "open"}, nil
		},
	}

	_, err := newTicketsUseCase(repo, nil).ChangeStatus(context.Background(),
		dtos.ChangeStatusDTO{TicketID: 1, NewStatus: "resolved"})

	require.NoError(t, err)
	require.Len(t, repo.Updates, 1)
	assert.Contains(t, repo.Updates[0], "resolved_at")
	assert.NotContains(t, repo.Updates[0], "closed_at")
}

func TestChangeStatus_CerradoSinResolverAntes_MarcaAmbasFechas(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, id uint) (*entities.Ticket, error) {
			return &entities.Ticket{ID: id, Status: "open", ResolvedAt: nil}, nil
		},
	}

	_, err := newTicketsUseCase(repo, nil).ChangeStatus(context.Background(),
		dtos.ChangeStatusDTO{TicketID: 1, NewStatus: "closed"})

	require.NoError(t, err)
	require.Len(t, repo.Updates, 1)
	assert.Contains(t, repo.Updates[0], "closed_at")
	assert.Contains(t, repo.Updates[0], "resolved_at",
		"cerrar sin haber pasado por resuelto tambien fija la fecha de resolucion")
}

func TestChangeStatus_CerradoYaResuelto_ConservaLaFechaOriginal(t *testing.T) {
	resuelto := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, id uint) (*entities.Ticket, error) {
			return &entities.Ticket{ID: id, Status: "resolved", ResolvedAt: &resuelto}, nil
		},
	}

	_, err := newTicketsUseCase(repo, nil).ChangeStatus(context.Background(),
		dtos.ChangeStatusDTO{TicketID: 1, NewStatus: "closed"})

	require.NoError(t, err)
	require.Len(t, repo.Updates, 1)
	assert.Contains(t, repo.Updates[0], "closed_at")
	assert.NotContains(t, repo.Updates[0], "resolved_at",
		"si ya estaba resuelto no se pisa la fecha original de resolucion")
}

func TestChangeStatus_RegistraElHistorialConAmbosEstados(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, id uint) (*entities.Ticket, error) {
			return &entities.Ticket{ID: id, Status: "open"}, nil
		},
	}

	_, err := newTicketsUseCase(repo, nil).ChangeStatus(context.Background(),
		dtos.ChangeStatusDTO{TicketID: 1, NewStatus: "in_development", ChangedByID: 42, Note: "empezamos"})

	require.NoError(t, err)
	require.Len(t, repo.History, 1)
	assert.Equal(t, "open", repo.History[0].From)
	assert.Equal(t, "in_development", repo.History[0].To)
	assert.Equal(t, uint(42), repo.History[0].ChangedByID)
	assert.Equal(t, "empezamos", repo.History[0].Note)
}

func TestChangeStatus_ErroresDelRepo_SePropagan(t *testing.T) {
	dbErr := stderrors.New("db caida")

	t.Run("al leer", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			GetByIDFn: func(ctx context.Context, id uint) (*entities.Ticket, error) { return nil, dbErr },
		}
		_, err := newTicketsUseCase(repo, nil).ChangeStatus(context.Background(),
			dtos.ChangeStatusDTO{TicketID: 1, NewStatus: "closed"})
		assert.ErrorIs(t, err, dbErr)
	})

	t.Run("al actualizar", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			UpdateFn: func(ctx context.Context, id uint, u map[string]any) (*entities.Ticket, error) {
				return nil, dbErr
			},
		}
		_, err := newTicketsUseCase(repo, nil).ChangeStatus(context.Background(),
			dtos.ChangeStatusDTO{TicketID: 1, NewStatus: "closed"})
		assert.ErrorIs(t, err, dbErr)
		assert.Empty(t, repo.History)
	})
}

func TestChangeArea_AreaInvalida_Rechaza(t *testing.T) {
	repo := &mocks.RepositoryMock{}

	_, err := newTicketsUseCase(repo, nil).ChangeArea(context.Background(),
		dtos.ChangeAreaDTO{TicketID: 1, NewArea: "marketing"})

	assert.ErrorIs(t, err, dom.ErrInvalidArea)
	assert.Empty(t, repo.Updates)
}

func TestChangeArea_MismaArea_NoActualiza(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, id uint) (*entities.Ticket, error) {
			return &entities.Ticket{ID: id, Area: "soporte"}, nil
		},
	}

	got, err := newTicketsUseCase(repo, nil).ChangeArea(context.Background(),
		dtos.ChangeAreaDTO{TicketID: 1, NewArea: "  SOPORTE  "})

	require.NoError(t, err)
	assert.Equal(t, "soporte", got.Area)
	assert.Empty(t, repo.Updates)
	assert.Empty(t, repo.AreaHistory)
}

func TestChangeArea_RegistraHistorialDeArea(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, id uint) (*entities.Ticket, error) {
			return &entities.Ticket{ID: id, Area: "soporte"}, nil
		},
	}

	_, err := newTicketsUseCase(repo, nil).ChangeArea(context.Background(),
		dtos.ChangeAreaDTO{TicketID: 1, NewArea: "desarrollo", ChangedByID: 42, Note: "es un bug"})

	require.NoError(t, err)
	require.Len(t, repo.AreaHistory, 1)
	assert.Equal(t, "soporte", repo.AreaHistory[0].From)
	assert.Equal(t, "desarrollo", repo.AreaHistory[0].To)
	assert.Equal(t, "es un bug", repo.AreaHistory[0].Note)
	assert.Empty(t, repo.History, "cambiar de area no toca el historial de estados")
}

func TestChangeArea_FalloAlRegistrarHistorial_NoRompeElCambio(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, id uint) (*entities.Ticket, error) {
			return &entities.Ticket{ID: id, Area: "soporte"}, nil
		},
		AddAreaHistoryFn: func(ctx context.Context, ticketID uint, from, to string, by uint, note string) error {
			return stderrors.New("db caida")
		},
	}

	got, err := newTicketsUseCase(repo, nil).ChangeArea(context.Background(),
		dtos.ChangeAreaDTO{TicketID: 1, NewArea: "desarrollo"})

	require.NoError(t, err, "el historial es auditoria, no debe bloquear el cambio de area")
	assert.NotNil(t, got)
}

func TestChangeArea_TicketNulo_Rechaza(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, id uint) (*entities.Ticket, error) { return nil, nil },
	}

	_, err := newTicketsUseCase(repo, nil).ChangeArea(context.Background(),
		dtos.ChangeAreaDTO{TicketID: 404, NewArea: "desarrollo"})

	assert.ErrorIs(t, err, dom.ErrTicketNotFound)
}

func TestChangeArea_ErroresDelRepo_SePropagan(t *testing.T) {
	dbErr := stderrors.New("db caida")

	t.Run("al leer", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			GetByIDFn: func(ctx context.Context, id uint) (*entities.Ticket, error) { return nil, dbErr },
		}
		_, err := newTicketsUseCase(repo, nil).ChangeArea(context.Background(),
			dtos.ChangeAreaDTO{TicketID: 1, NewArea: "desarrollo"})
		assert.ErrorIs(t, err, dbErr)
	})

	t.Run("al actualizar", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			GetByIDFn: func(ctx context.Context, id uint) (*entities.Ticket, error) {
				return &entities.Ticket{ID: id, Area: "soporte"}, nil
			},
			UpdateFn: func(ctx context.Context, id uint, u map[string]any) (*entities.Ticket, error) {
				return nil, dbErr
			},
		}
		_, err := newTicketsUseCase(repo, nil).ChangeArea(context.Background(),
			dtos.ChangeAreaDTO{TicketID: 1, NewArea: "desarrollo"})
		assert.ErrorIs(t, err, dbErr)
	})
}

func TestAssign_ResponsableInexistente_NoAsigna(t *testing.T) {
	repo := &mocks.RepositoryMock{
		UserExistsFn: func(ctx context.Context, userID uint) (bool, error) { return false, nil },
	}

	_, err := newTicketsUseCase(repo, nil).Assign(context.Background(),
		dtos.AssignTicketDTO{TicketID: 1, AssignedToID: uintPtr(404)})

	assert.ErrorIs(t, err, dom.ErrAssigneeNotFound)
	assert.Empty(t, repo.Updates)
}

func TestAssign_ConResponsable_GuardaElIDYDejaNota(t *testing.T) {
	repo := &mocks.RepositoryMock{}

	_, err := newTicketsUseCase(repo, nil).Assign(context.Background(),
		dtos.AssignTicketDTO{TicketID: 1, AssignedToID: uintPtr(7), ChangedByID: 42})

	require.NoError(t, err)
	require.Len(t, repo.Updates, 1)
	assert.Equal(t, uint(7), repo.Updates[0]["assigned_to_id"])
	require.Len(t, repo.History, 1)
	assert.Contains(t, repo.History[0].Note, "7")
}

func TestAssign_SinResponsable_DesasignaSinValidarUsuario(t *testing.T) {
	validado := false
	repo := &mocks.RepositoryMock{
		UserExistsFn: func(ctx context.Context, userID uint) (bool, error) {
			validado = true
			return true, nil
		},
	}

	_, err := newTicketsUseCase(repo, nil).Assign(context.Background(),
		dtos.AssignTicketDTO{TicketID: 1, AssignedToID: nil})

	require.NoError(t, err)
	assert.False(t, validado)
	require.Len(t, repo.Updates, 1)
	assert.Nil(t, repo.Updates[0]["assigned_to_id"])
	assert.Contains(t, repo.History[0].Note, "sin asignar")
}

func TestAssign_NoCambiaElEstadoDelTicket(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, id uint) (*entities.Ticket, error) {
			return &entities.Ticket{ID: id, Status: "in_development"}, nil
		},
	}

	_, err := newTicketsUseCase(repo, nil).Assign(context.Background(),
		dtos.AssignTicketDTO{TicketID: 1, AssignedToID: uintPtr(7)})

	require.NoError(t, err)
	assert.NotContains(t, repo.Updates[0], "status")
	require.Len(t, repo.History, 1)
	assert.Equal(t, repo.History[0].From, repo.History[0].To,
		"asignar no es una transicion de estado")
	assert.Equal(t, "in_development", repo.History[0].From)
}

func TestAssign_ErroresDelRepo_SePropagan(t *testing.T) {
	dbErr := stderrors.New("db caida")

	t.Run("al leer", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			GetByIDFn: func(ctx context.Context, id uint) (*entities.Ticket, error) { return nil, dbErr },
		}
		_, err := newTicketsUseCase(repo, nil).Assign(context.Background(), dtos.AssignTicketDTO{TicketID: 1})
		assert.ErrorIs(t, err, dbErr)
	})

	t.Run("al validar usuario", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			UserExistsFn: func(ctx context.Context, userID uint) (bool, error) { return false, dbErr },
		}
		_, err := newTicketsUseCase(repo, nil).Assign(context.Background(),
			dtos.AssignTicketDTO{TicketID: 1, AssignedToID: uintPtr(7)})
		assert.ErrorIs(t, err, dbErr)
	})

	t.Run("al actualizar", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			UpdateFn: func(ctx context.Context, id uint, u map[string]any) (*entities.Ticket, error) {
				return nil, dbErr
			},
		}
		_, err := newTicketsUseCase(repo, nil).Assign(context.Background(), dtos.AssignTicketDTO{TicketID: 1})
		assert.ErrorIs(t, err, dbErr)
		assert.Empty(t, repo.History)
	})
}

func TestEscalate_MarcaEscaladoConFecha(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, id uint) (*entities.Ticket, error) {
			return &entities.Ticket{ID: id, Status: "in_review"}, nil
		},
	}

	_, err := newTicketsUseCase(repo, nil).Escalate(context.Background(),
		dtos.EscalateTicketDTO{TicketID: 1, ChangedByID: 42})

	require.NoError(t, err)
	require.Len(t, repo.Updates, 1)
	assert.Equal(t, true, repo.Updates[0]["escalated_to_dev"])
	assert.Contains(t, repo.Updates[0], "escalated_at")
}

func TestEscalate_DesdeAbierto_PasaAEnRevision(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, id uint) (*entities.Ticket, error) {
			return &entities.Ticket{ID: id, Status: "open"}, nil
		},
	}

	_, err := newTicketsUseCase(repo, nil).Escalate(context.Background(),
		dtos.EscalateTicketDTO{TicketID: 1})

	require.NoError(t, err)
	assert.Equal(t, "in_review", repo.Updates[0]["status"],
		"escalar un ticket abierto lo mueve a revision automaticamente")
}

func TestEscalate_DesdeOtroEstado_NoLoCambia(t *testing.T) {
	for _, estado := range []string{"in_review", "in_development", "blocked", "testing"} {
		t.Run(estado, func(t *testing.T) {
			repo := &mocks.RepositoryMock{
				GetByIDFn: func(ctx context.Context, id uint) (*entities.Ticket, error) {
					return &entities.Ticket{ID: id, Status: estado}, nil
				},
			}

			_, err := newTicketsUseCase(repo, nil).Escalate(context.Background(),
				dtos.EscalateTicketDTO{TicketID: 1})

			require.NoError(t, err)
			assert.NotContains(t, repo.Updates[0], "status",
				"escalar no debe retroceder un ticket que ya avanzo")
		})
	}
}

func TestEscalate_SinNota_UsaLaPorDefecto(t *testing.T) {
	repo := &mocks.RepositoryMock{}

	_, err := newTicketsUseCase(repo, nil).Escalate(context.Background(),
		dtos.EscalateTicketDTO{TicketID: 1})

	require.NoError(t, err)
	require.Len(t, repo.History, 1)
	assert.Equal(t, "Escalado a desarrollo", repo.History[0].Note)
}

func TestEscalate_ConNota_LaRespeta(t *testing.T) {
	repo := &mocks.RepositoryMock{}

	_, err := newTicketsUseCase(repo, nil).Escalate(context.Background(),
		dtos.EscalateTicketDTO{TicketID: 1, Note: "bloquea a 3 clientes"})

	require.NoError(t, err)
	assert.Equal(t, "bloquea a 3 clientes", repo.History[0].Note)
}

func TestEscalate_ErroresDelRepo_SePropagan(t *testing.T) {
	dbErr := stderrors.New("db caida")

	t.Run("al leer", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			GetByIDFn: func(ctx context.Context, id uint) (*entities.Ticket, error) { return nil, dbErr },
		}
		_, err := newTicketsUseCase(repo, nil).Escalate(context.Background(), dtos.EscalateTicketDTO{TicketID: 1})
		assert.ErrorIs(t, err, dbErr)
	})

	t.Run("al actualizar", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			UpdateFn: func(ctx context.Context, id uint, u map[string]any) (*entities.Ticket, error) {
				return nil, dbErr
			},
		}
		_, err := newTicketsUseCase(repo, nil).Escalate(context.Background(), dtos.EscalateTicketDTO{TicketID: 1})
		assert.ErrorIs(t, err, dbErr)
		assert.Empty(t, repo.History)
	})
}

func TestGet_SuperAdminVeCualquierTicket(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, id uint) (*entities.Ticket, error) {
			return &entities.Ticket{ID: id, BusinessID: uintPtr(99)}, nil
		},
	}

	got, err := newTicketsUseCase(repo, nil).Get(context.Background(), 1, 42, uintPtr(26), true)

	require.NoError(t, err)
	assert.Equal(t, uint(1), got.ID)
}

func TestGet_UsuarioDeOtroNegocio_Rechaza(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, id uint) (*entities.Ticket, error) {
			return &entities.Ticket{ID: id, BusinessID: uintPtr(99)}, nil
		},
	}

	got, err := newTicketsUseCase(repo, nil).Get(context.Background(), 1, 42, uintPtr(26), false)

	assert.ErrorIs(t, err, dom.ErrForbidden,
		"un negocio no puede ver los tickets de otro aunque adivine el ID")
	assert.Nil(t, got)
}

func TestGet_UsuarioDeSuPropioNegocio_Permite(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, id uint) (*entities.Ticket, error) {
			return &entities.Ticket{ID: id, BusinessID: uintPtr(26)}, nil
		},
	}

	got, err := newTicketsUseCase(repo, nil).Get(context.Background(), 1, 42, uintPtr(26), false)

	require.NoError(t, err)
	assert.Equal(t, uint(1), got.ID)
}

func TestGet_TicketSinNegocioOSolicitanteSinNegocio_Rechaza(t *testing.T) {
	casos := []struct {
		nombre           string
		ticketBusinessID *uint
		requesterID      *uint
	}{
		{"ticket interno sin negocio", nil, uintPtr(26)},
		{"solicitante sin negocio", uintPtr(26), nil},
		{"ambos sin negocio", nil, nil},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			repo := &mocks.RepositoryMock{
				GetByIDFn: func(ctx context.Context, id uint) (*entities.Ticket, error) {
					return &entities.Ticket{ID: id, BusinessID: tc.ticketBusinessID}, nil
				},
			}

			_, err := newTicketsUseCase(repo, nil).Get(context.Background(), 1, 42, tc.requesterID, false)

			assert.ErrorIs(t, err, dom.ErrForbidden,
				"sin negocio comparable no se concede acceso")
		})
	}
}

func TestGet_ErrorDelRepo_SePropaga(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, id uint) (*entities.Ticket, error) {
			return nil, dom.ErrTicketNotFound
		},
	}

	_, err := newTicketsUseCase(repo, nil).Get(context.Background(), 404, 42, uintPtr(26), true)

	assert.ErrorIs(t, err, dom.ErrTicketNotFound)
}

func TestList_NormalizaPaginacion(t *testing.T) {
	casos := []struct {
		page, pageSize, wantPage, wantSize int
	}{
		{0, 10, 1, 10},
		{-5, 10, 1, 10},
		{1, 0, 1, 10},
		{1, -1, 1, 10},
		{1, 101, 1, 10},
		{1, 100, 1, 100},
		{3, 50, 3, 50},
	}

	for _, tc := range casos {
		var visto dtos.ListTicketsParams
		repo := &mocks.RepositoryMock{
			ListFn: func(ctx context.Context, p dtos.ListTicketsParams) ([]entities.Ticket, int64, error) {
				visto = p
				return nil, 0, nil
			},
		}

		_, _, err := newTicketsUseCase(repo, nil).List(context.Background(),
			dtos.ListTicketsParams{Page: tc.page, PageSize: tc.pageSize})

		require.NoError(t, err)
		assert.Equal(t, tc.wantPage, visto.Page)
		assert.Equal(t, tc.wantSize, visto.PageSize)
	}
}

func TestList_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := stderrors.New("timeout")
	repo := &mocks.RepositoryMock{
		ListFn: func(ctx context.Context, p dtos.ListTicketsParams) ([]entities.Ticket, int64, error) {
			return nil, 0, dbErr
		},
	}

	_, _, err := newTicketsUseCase(repo, nil).List(context.Background(), dtos.ListTicketsParams{})

	assert.ErrorIs(t, err, dbErr)
}

func TestUpdate_CamposNilNoSeTocan(t *testing.T) {
	leido := false
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, id uint) (*entities.Ticket, error) {
			leido = true
			return &entities.Ticket{ID: id}, nil
		},
	}

	_, err := newTicketsUseCase(repo, nil).Update(context.Background(), dtos.UpdateTicketDTO{ID: 1})

	require.NoError(t, err)
	assert.Empty(t, repo.Updates, "sin cambios no se escribe nada")
	assert.True(t, leido, "se devuelve el ticket tal como esta")
}

func TestUpdate_ValoresInvalidos_Rechaza(t *testing.T) {
	casos := []struct {
		nombre  string
		dto     dtos.UpdateTicketDTO
		wantErr error
	}{
		{"titulo vacio", dtos.UpdateTicketDTO{ID: 1, Title: strPtr("  ")}, dom.ErrTitleRequired},
		{"descripcion vacia", dtos.UpdateTicketDTO{ID: 1, Description: strPtr("  ")}, dom.ErrDescriptionRequired},
		{"tipo invalido", dtos.UpdateTicketDTO{ID: 1, Type: strPtr("inventado")}, dom.ErrInvalidType},
		{"prioridad invalida", dtos.UpdateTicketDTO{ID: 1, Priority: strPtr("urgentisimo")}, dom.ErrInvalidPriority},
		{"severidad invalida", dtos.UpdateTicketDTO{ID: 1, Severity: strPtr("catastrofica")}, dom.ErrInvalidSeverity},
		{"area invalida", dtos.UpdateTicketDTO{ID: 1, Area: strPtr("marketing")}, dom.ErrInvalidArea},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			repo := &mocks.RepositoryMock{}

			_, err := newTicketsUseCase(repo, nil).Update(context.Background(), tc.dto)

			assert.ErrorIs(t, err, tc.wantErr)
			assert.Empty(t, repo.Updates)
		})
	}
}

func TestUpdate_SeveridadVaciaSePermiteParaLimpiarla(t *testing.T) {
	repo := &mocks.RepositoryMock{}

	_, err := newTicketsUseCase(repo, nil).Update(context.Background(),
		dtos.UpdateTicketDTO{ID: 1, Severity: strPtr("")})

	require.NoError(t, err, "una severidad vacia es valida: significa quitarla")
	require.Len(t, repo.Updates, 1)
	assert.Equal(t, "", repo.Updates[0]["severity"])
}

func TestUpdate_NormalizaYRecortaLosCampos(t *testing.T) {
	repo := &mocks.RepositoryMock{}

	_, err := newTicketsUseCase(repo, nil).Update(context.Background(), dtos.UpdateTicketDTO{
		ID: 1, Title: strPtr("  nuevo  "), Type: strPtr("  BUG  "),
		Priority: strPtr("  HIGH  "), Area: strPtr("  DESARROLLO  "),
		Category: strPtr("  facturacion  "),
	})

	require.NoError(t, err)
	require.Len(t, repo.Updates, 1)
	u := repo.Updates[0]
	assert.Equal(t, "nuevo", u["title"])
	assert.Equal(t, "bug", u["type"])
	assert.Equal(t, "high", u["priority"])
	assert.Equal(t, "desarrollo", u["area"])
	assert.Equal(t, "facturacion", u["category"])
}

func TestUpdate_ResponsableInexistente_NoActualiza(t *testing.T) {
	repo := &mocks.RepositoryMock{
		UserExistsFn: func(ctx context.Context, userID uint) (bool, error) { return false, nil },
	}

	_, err := newTicketsUseCase(repo, nil).Update(context.Background(),
		dtos.UpdateTicketDTO{ID: 1, AssignedToID: uintPtr(404)})

	assert.ErrorIs(t, err, dom.ErrAssigneeNotFound)
	assert.Empty(t, repo.Updates)
}

func TestUpdate_LimpiarFechaLimiteGanaSobreEnviarla(t *testing.T) {
	repo := &mocks.RepositoryMock{}

	_, err := newTicketsUseCase(repo, nil).Update(context.Background(), dtos.UpdateTicketDTO{
		ID: 1, ClearDueDate: true, DueDate: strPtr("2026-12-31"),
	})

	require.NoError(t, err)
	require.Len(t, repo.Updates, 1)
	assert.Nil(t, repo.Updates[0]["due_date"],
		"si se pide limpiar la fecha, se limpia aunque venga una fecha en el mismo request")
}

func TestUpdate_FechaLimiteInvalida_NoSeEscribe(t *testing.T) {
	repo := &mocks.RepositoryMock{}

	_, err := newTicketsUseCase(repo, nil).Update(context.Background(),
		dtos.UpdateTicketDTO{ID: 1, DueDate: strPtr("31/12/2026")})

	require.NoError(t, err)
	assert.Empty(t, repo.Updates, "una fecha mal formada se ignora, no se escribe basura")
}

func TestUpdate_NoPermiteCambiarElEstado(t *testing.T) {
	repo := &mocks.RepositoryMock{}

	_, err := newTicketsUseCase(repo, nil).Update(context.Background(),
		dtos.UpdateTicketDTO{ID: 1, Title: strPtr("nuevo")})

	require.NoError(t, err)
	assert.NotContains(t, repo.Updates[0], "status",
		"el estado se cambia solo por ChangeStatus, que registra historial")
}

func TestUpdate_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := stderrors.New("db caida")
	repo := &mocks.RepositoryMock{
		UpdateFn: func(ctx context.Context, id uint, u map[string]any) (*entities.Ticket, error) {
			return nil, dbErr
		},
	}

	_, err := newTicketsUseCase(repo, nil).Update(context.Background(),
		dtos.UpdateTicketDTO{ID: 1, Title: strPtr("nuevo")})

	assert.ErrorIs(t, err, dbErr)
}

func TestAddComment_CuerpoVacio_Rechaza(t *testing.T) {
	for _, body := range []string{"", "   ", "\t\n"} {
		repo := &mocks.RepositoryMock{}

		_, err := newTicketsUseCase(repo, nil).AddComment(context.Background(),
			dtos.CreateCommentDTO{TicketID: 1, Body: body})

		assert.ErrorIs(t, err, dom.ErrDescriptionRequired)
		assert.Empty(t, repo.AddedComments)
	}
}

func TestAddComment_RecortaElCuerpo(t *testing.T) {
	repo := &mocks.RepositoryMock{}

	_, err := newTicketsUseCase(repo, nil).AddComment(context.Background(),
		dtos.CreateCommentDTO{TicketID: 1, Body: "  ya lo revise  "})

	require.NoError(t, err)
	require.Len(t, repo.AddedComments, 1)
	assert.Equal(t, "ya lo revise", repo.AddedComments[0].Body)
}

func TestAddComment_TicketInexistente_NoComenta(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, id uint) (*entities.Ticket, error) {
			return nil, dom.ErrTicketNotFound
		},
	}

	_, err := newTicketsUseCase(repo, nil).AddComment(context.Background(),
		dtos.CreateCommentDTO{TicketID: 404, Body: "hola"})

	assert.ErrorIs(t, err, dom.ErrTicketNotFound)
	assert.Empty(t, repo.AddedComments)
}

func TestListComments_PasaElFlagDeInternos(t *testing.T) {
	for _, incluir := range []bool{true, false} {
		var visto bool
		repo := &mocks.RepositoryMock{
			ListCommentsFn: func(ctx context.Context, ticketID uint, includeInternal bool) ([]entities.TicketComment, error) {
				visto = includeInternal
				return nil, nil
			},
		}

		_, err := newTicketsUseCase(repo, nil).ListComments(context.Background(), 1, incluir)

		require.NoError(t, err)
		assert.Equal(t, incluir, visto,
			"el flag decide si el cliente ve los comentarios internos del equipo")
	}
}

func TestDeleteAttachment_DeOtroUsuario_Rechaza(t *testing.T) {
	storage := &mocks.StorageServiceMock{}
	repo := &mocks.RepositoryMock{
		GetAttachmentFn: func(ctx context.Context, id uint) (*entities.TicketAttachment, error) {
			return &entities.TicketAttachment{ID: id, UploadedByID: 99, FileURL: "s3://x.png"}, nil
		},
	}

	err := newTicketsUseCase(repo, storage).DeleteAttachment(context.Background(), 1, 42, false)

	assert.ErrorIs(t, err, dom.ErrForbidden)
	assert.Empty(t, storage.DeletedFiles, "no se borra el archivo de alguien mas")
}

func TestDeleteAttachment_SuPropioAdjunto_Permite(t *testing.T) {
	storage := &mocks.StorageServiceMock{}
	repo := &mocks.RepositoryMock{
		GetAttachmentFn: func(ctx context.Context, id uint) (*entities.TicketAttachment, error) {
			return &entities.TicketAttachment{ID: id, UploadedByID: 42, FileURL: "s3://x.png"}, nil
		},
	}

	err := newTicketsUseCase(repo, storage).DeleteAttachment(context.Background(), 1, 42, false)

	require.NoError(t, err)
	assert.Equal(t, []string{"s3://x.png"}, storage.DeletedFiles)
}

func TestDeleteAttachment_SuperAdminBorraCualquiera(t *testing.T) {
	storage := &mocks.StorageServiceMock{}
	repo := &mocks.RepositoryMock{
		GetAttachmentFn: func(ctx context.Context, id uint) (*entities.TicketAttachment, error) {
			return &entities.TicketAttachment{ID: id, UploadedByID: 99, FileURL: "s3://x.png"}, nil
		},
	}

	err := newTicketsUseCase(repo, storage).DeleteAttachment(context.Background(), 1, 42, true)

	require.NoError(t, err)
	assert.Len(t, storage.DeletedFiles, 1)
}

func TestDeleteAttachment_FalloEnS3_IgualBorraElRegistro(t *testing.T) {
	borrado := false
	storage := &mocks.StorageServiceMock{
		DeleteFileFn: func(ctx context.Context, fileURL string) error {
			return stderrors.New("s3 caido")
		},
	}
	repo := &mocks.RepositoryMock{
		GetAttachmentFn: func(ctx context.Context, id uint) (*entities.TicketAttachment, error) {
			return &entities.TicketAttachment{ID: id, UploadedByID: 42}, nil
		},
		DeleteAttachmentFn: func(ctx context.Context, id uint) error {
			borrado = true
			return nil
		},
	}

	err := newTicketsUseCase(repo, storage).DeleteAttachment(context.Background(), 1, 42, false)

	require.NoError(t, err, "un archivo huerfano en S3 es mejor que un registro que apunta a la nada")
	assert.True(t, borrado)
}

func TestDeleteAttachment_AdjuntoInexistente_SePropaga(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetAttachmentFn: func(ctx context.Context, id uint) (*entities.TicketAttachment, error) {
			return nil, dom.ErrAttachmentNotFound
		},
	}

	err := newTicketsUseCase(repo, nil).DeleteAttachment(context.Background(), 404, 42, true)

	assert.ErrorIs(t, err, dom.ErrAttachmentNotFound)
}

func TestDelete_BorraLosArchivosDeS3AntesDelTicket(t *testing.T) {
	storage := &mocks.StorageServiceMock{}
	repo := &mocks.RepositoryMock{
		ListAttachmentsFn: func(ctx context.Context, ticketID uint) ([]entities.TicketAttachment, error) {
			return []entities.TicketAttachment{
				{ID: 1, FileURL: "s3://a.png"},
				{ID: 2, FileURL: ""},
				{ID: 3, FileURL: "s3://b.pdf"},
			}, nil
		},
	}

	err := newTicketsUseCase(repo, storage).Delete(context.Background(), 1)

	require.NoError(t, err)
	assert.Equal(t, []string{"s3://a.png", "s3://b.pdf"}, storage.DeletedFiles,
		"los adjuntos sin URL se saltan")
	assert.Equal(t, []uint{1}, repo.DeletedTickets)
}

func TestDelete_FalloAlListarAdjuntos_IgualBorraElTicket(t *testing.T) {
	repo := &mocks.RepositoryMock{
		ListAttachmentsFn: func(ctx context.Context, ticketID uint) ([]entities.TicketAttachment, error) {
			return nil, stderrors.New("db caida")
		},
	}

	err := newTicketsUseCase(repo, nil).Delete(context.Background(), 1)

	require.NoError(t, err)
	assert.Equal(t, []uint{1}, repo.DeletedTickets)
}

func TestDelete_FalloEnS3_IgualBorraElTicket(t *testing.T) {
	storage := &mocks.StorageServiceMock{
		DeleteFileFn: func(ctx context.Context, fileURL string) error {
			return stderrors.New("s3 caido")
		},
	}
	repo := &mocks.RepositoryMock{
		ListAttachmentsFn: func(ctx context.Context, ticketID uint) ([]entities.TicketAttachment, error) {
			return []entities.TicketAttachment{{ID: 1, FileURL: "s3://a.png"}}, nil
		},
	}

	err := newTicketsUseCase(repo, storage).Delete(context.Background(), 1)

	require.NoError(t, err)
	assert.Equal(t, []uint{1}, repo.DeletedTickets)
}

func TestDelete_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := stderrors.New("fk violation")
	repo := &mocks.RepositoryMock{
		DeleteFn: func(ctx context.Context, id uint) error { return dbErr },
	}

	err := newTicketsUseCase(repo, nil).Delete(context.Background(), 1)

	assert.ErrorIs(t, err, dbErr)
}

func TestListAttachmentsYListHistory_DeleganAlRepositorio(t *testing.T) {
	var vistoAdj, vistoHist uint
	repo := &mocks.RepositoryMock{
		ListAttachmentsFn: func(ctx context.Context, ticketID uint) ([]entities.TicketAttachment, error) {
			vistoAdj = ticketID
			return []entities.TicketAttachment{{ID: 1}}, nil
		},
		ListHistoryFn: func(ctx context.Context, ticketID uint) ([]entities.TicketStatusHistory, error) {
			vistoHist = ticketID
			return []entities.TicketStatusHistory{{ID: 1}}, nil
		},
	}
	uc := newTicketsUseCase(repo, nil)

	adj, err := uc.ListAttachments(context.Background(), 5)
	require.NoError(t, err)
	assert.Equal(t, uint(5), vistoAdj)
	assert.Len(t, adj, 1)

	hist, err := uc.ListHistory(context.Background(), 5)
	require.NoError(t, err)
	assert.Equal(t, uint(5), vistoHist)
	assert.Len(t, hist, 1)
}

func TestCatalogos_NoSeSolapanEntreSi(t *testing.T) {
	assert.Len(t, validStatuses, 8)
	assert.Len(t, validPriorities, 4)
	assert.Len(t, validTypes, 9)
	assert.Len(t, validSeverities, 3)
	assert.Len(t, validSources, 2)
	assert.Len(t, validAreas, 3)

	for estado := range closedStatuses {
		assert.True(t, validStatuses[estado],
			"el estado cerrado %q debe existir en el catalogo de estados", estado)
	}
}

func archivoDePrueba(t *testing.T, nombre, contenido, contentType string) *multipart.FileHeader {
	t.Helper()
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="file"; filename="`+nombre+`"`)
	h.Set("Content-Type", contentType)
	part, err := w.CreatePart(h)
	require.NoError(t, err)
	_, err = part.Write([]byte(contenido))
	require.NoError(t, err)
	require.NoError(t, w.Close())

	r := multipart.NewReader(&buf, w.Boundary())
	form, err := r.ReadForm(1 << 20)
	require.NoError(t, err)
	return form.File["file"][0]
}

func TestUploadAttachment_TicketInexistente_NoSube(t *testing.T) {
	storage := &mocks.StorageServiceMock{}
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, id uint) (*entities.Ticket, error) {
			return nil, dom.ErrTicketNotFound
		},
	}

	_, err := newTicketsUseCase(repo, storage).UploadAttachment(context.Background(),
		404, nil, 42, archivoDePrueba(t, "captura.png", "datos", "image/png"))

	assert.ErrorIs(t, err, dom.ErrTicketNotFound)
	assert.Empty(t, storage.UploadedFolders)
}

func TestUploadAttachment_GuardaEnLaCarpetaDelTicket(t *testing.T) {
	storage := &mocks.StorageServiceMock{}
	repo := &mocks.RepositoryMock{}

	_, err := newTicketsUseCase(repo, storage).UploadAttachment(context.Background(),
		7, nil, 42, archivoDePrueba(t, "captura.png", "datos", "image/png"))

	require.NoError(t, err)
	assert.Equal(t, []string{"tickets/7"}, storage.UploadedFolders,
		"cada ticket tiene su carpeta, no se mezclan adjuntos entre tickets")
}

func TestUploadAttachment_ConservaElNombreOriginalPeroNoLoUsaComoRuta(t *testing.T) {
	var nombreEnS3 string
	storage := &mocks.StorageServiceMock{
		UploadFileFn: func(ctx context.Context, folder, filename string, data []byte, contentType string) (string, error) {
			nombreEnS3 = filename
			return "s3://" + folder + "/" + filename, nil
		},
	}
	repo := &mocks.RepositoryMock{}

	_, err := newTicketsUseCase(repo, storage).UploadAttachment(context.Background(),
		7, nil, 42, archivoDePrueba(t, "../../etc/passwd.png", "datos", "image/png"))

	require.NoError(t, err)
	require.NotNil(t, repo.AddedAttachment)
	assert.Equal(t, "passwd.png", repo.AddedAttachment.FileName,
		"el parser multipart de Go ya recorta la ruta y deja solo el nombre base")
	assert.NotContains(t, nombreEnS3, "..",
		"y ademas el nombre real en S3 es generado, no permite escapar de la carpeta")
	assert.NotContains(t, nombreEnS3, "/")
	assert.True(t, strings.HasSuffix(nombreEnS3, ".png"), "conserva la extension: %q", nombreEnS3)
}

func TestUploadAttachment_NombresIgualesNoSePisan(t *testing.T) {
	vistos := map[string]bool{}
	storage := &mocks.StorageServiceMock{
		UploadFileFn: func(ctx context.Context, folder, filename string, data []byte, contentType string) (string, error) {
			require.False(t, vistos[filename], "nombre repetido en S3: %q", filename)
			vistos[filename] = true
			return "s3://" + filename, nil
		},
	}
	uc := newTicketsUseCase(&mocks.RepositoryMock{}, storage)

	for i := 0; i < 20; i++ {
		_, err := uc.UploadAttachment(context.Background(), 7, nil, 42,
			archivoDePrueba(t, "captura.png", "datos", "image/png"))
		require.NoError(t, err)
	}
	assert.Len(t, vistos, 20)
}

func TestUploadAttachment_GuardaMetadatosDelArchivo(t *testing.T) {
	storage := &mocks.StorageServiceMock{
		UploadFileFn: func(ctx context.Context, folder, filename string, data []byte, contentType string) (string, error) {
			return "s3://bucket/x.png", nil
		},
	}
	repo := &mocks.RepositoryMock{}
	comentario := uint(3)

	_, err := newTicketsUseCase(repo, storage).UploadAttachment(context.Background(),
		7, &comentario, 42, archivoDePrueba(t, "captura.png", "contenido", "image/png"))

	require.NoError(t, err)
	require.NotNil(t, repo.AddedAttachment)
	a := repo.AddedAttachment
	assert.Equal(t, uint(7), a.TicketID)
	require.NotNil(t, a.CommentID)
	assert.Equal(t, uint(3), *a.CommentID)
	assert.Equal(t, uint(42), a.UploadedByID)
	assert.Equal(t, "s3://bucket/x.png", a.FileURL)
	assert.Equal(t, "image/png", a.MimeType)
	assert.Equal(t, int64(len("contenido")), a.Size)
}

func TestUploadAttachment_FalloEnS3_NoGuardaElRegistro(t *testing.T) {
	storage := &mocks.StorageServiceMock{
		UploadFileFn: func(ctx context.Context, folder, filename string, data []byte, contentType string) (string, error) {
			return "", stderrors.New("s3 caido")
		},
	}
	repo := &mocks.RepositoryMock{}

	_, err := newTicketsUseCase(repo, storage).UploadAttachment(context.Background(),
		7, nil, 42, archivoDePrueba(t, "captura.png", "datos", "image/png"))

	require.Error(t, err)
	assert.Contains(t, err.Error(), "upload")
	assert.Nil(t, repo.AddedAttachment,
		"un registro sin archivo real apuntaria a la nada")
}
