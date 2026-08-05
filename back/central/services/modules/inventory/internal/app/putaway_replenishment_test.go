package app

import (
	"context"
	"errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/app/request"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func repoPutawayFeliz() *mocks.RepositoryMock {
	return &mocks.RepositoryMock{
		FindApplicableRuleFn: func(ctx context.Context, businessID uint, productID string) (*entities.PutawayRule, error) {
			return &entities.PutawayRule{ID: 1, TargetZoneID: 3, Strategy: "nearest_empty"}, nil
		},
		PickLocationInZoneFn: func(ctx context.Context, zoneID uint) (uint, error) {
			return 42, nil
		},
		CreatePutawaySuggestionFn: func(ctx context.Context, s *entities.PutawaySuggestion) (*entities.PutawaySuggestion, error) {
			return s, nil
		},
	}
}

func TestSuggestPutaway_SinItems_ResultadoVacio(t *testing.T) {
	uc := buildUseCase(repoPutawayFeliz(), nil, nil)

	got, err := uc.SuggestPutaway(context.Background(), request.PutawaySuggestDTO{BusinessID: 10})

	require.NoError(t, err)
	assert.Empty(t, got.Suggestions)
	assert.Empty(t, got.UnresolvedItems)
}

func TestSuggestPutaway_GeneraSugerencia(t *testing.T) {
	var creada *entities.PutawaySuggestion
	repo := repoPutawayFeliz()
	repo.CreatePutawaySuggestionFn = func(ctx context.Context, s *entities.PutawaySuggestion) (*entities.PutawaySuggestion, error) {
		creada = s
		return s, nil
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.SuggestPutaway(context.Background(), request.PutawaySuggestDTO{
		BusinessID: 10,
		Items:      []request.PutawaySuggestItem{{ProductID: "prod-1", Quantity: 20}},
	})

	require.NoError(t, err)
	require.Len(t, got.Suggestions, 1)
	require.NotNil(t, creada)
	assert.Equal(t, uint(42), creada.RecommendedLocationID)
	assert.Equal(t, "pending", creada.Status)
	assert.Equal(t, "nearest_empty", creada.Reason)
	assert.Equal(t, 20, creada.Quantity)
	require.NotNil(t, creada.RuleID)
	assert.Equal(t, uint(1), *creada.RuleID)
}

func TestSuggestPutaway_SinReglaAplicable_QuedaSinResolver(t *testing.T) {
	repo := repoPutawayFeliz()
	repo.FindApplicableRuleFn = func(ctx context.Context, businessID uint, productID string) (*entities.PutawayRule, error) {
		return nil, domainerrors.ErrNoPutawayRule
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.SuggestPutaway(context.Background(), request.PutawaySuggestDTO{
		BusinessID: 10,
		Items:      []request.PutawaySuggestItem{{ProductID: "prod-1", Quantity: 20}},
	})

	require.NoError(t, err)
	assert.Empty(t, got.Suggestions)
	assert.Equal(t, []string{"prod-1"}, got.UnresolvedItems)
}

func TestSuggestPutaway_SinUbicacionLibre_QuedaSinResolver(t *testing.T) {
	repo := repoPutawayFeliz()
	repo.PickLocationInZoneFn = func(ctx context.Context, zoneID uint) (uint, error) {
		return 0, errors.New("zona llena")
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.SuggestPutaway(context.Background(), request.PutawaySuggestDTO{
		BusinessID: 10,
		Items:      []request.PutawaySuggestItem{{ProductID: "prod-1", Quantity: 20}},
	})

	require.NoError(t, err)
	assert.Equal(t, []string{"prod-1"}, got.UnresolvedItems)
}

func TestSuggestPutaway_FallaGuardar_QuedaSinResolver(t *testing.T) {
	repo := repoPutawayFeliz()
	repo.CreatePutawaySuggestionFn = func(ctx context.Context, s *entities.PutawaySuggestion) (*entities.PutawaySuggestion, error) {
		return nil, errors.New("db down")
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.SuggestPutaway(context.Background(), request.PutawaySuggestDTO{
		BusinessID: 10,
		Items:      []request.PutawaySuggestItem{{ProductID: "prod-1", Quantity: 20}},
	})

	require.NoError(t, err)
	assert.Equal(t, []string{"prod-1"}, got.UnresolvedItems)
}

func TestSuggestPutaway_UnItemFallaYOtroNo(t *testing.T) {
	repo := repoPutawayFeliz()
	repo.FindApplicableRuleFn = func(ctx context.Context, businessID uint, productID string) (*entities.PutawayRule, error) {
		if productID == "prod-malo" {
			return nil, domainerrors.ErrNoPutawayRule
		}
		return &entities.PutawayRule{ID: 1, TargetZoneID: 3, Strategy: "abc"}, nil
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.SuggestPutaway(context.Background(), request.PutawaySuggestDTO{
		BusinessID: 10,
		Items: []request.PutawaySuggestItem{
			{ProductID: "prod-bueno", Quantity: 5},
			{ProductID: "prod-malo", Quantity: 5},
		},
	})

	require.NoError(t, err)
	assert.Len(t, got.Suggestions, 1)
	assert.Equal(t, []string{"prod-malo"}, got.UnresolvedItems)
}

func TestConfirmPutaway_YaConfirmada_Rechaza(t *testing.T) {
	actualizada := false
	repo := &mocks.RepositoryMock{
		GetPutawaySuggestionByIDFn: func(ctx context.Context, businessID, id uint) (*entities.PutawaySuggestion, error) {
			return &entities.PutawaySuggestion{ID: id, Status: "confirmed"}, nil
		},
		UpdatePutawaySuggestionFn: func(ctx context.Context, s *entities.PutawaySuggestion) (*entities.PutawaySuggestion, error) {
			actualizada = true
			return s, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.ConfirmPutaway(context.Background(), request.ConfirmPutawayDTO{BusinessID: 10, SuggestionID: 1})

	assert.ErrorIs(t, err, domainerrors.ErrPutawayAlreadyConfirmed)
	assert.Nil(t, got)
	assert.False(t, actualizada)
}

func TestConfirmPutaway_Pendiente_SeConfirma(t *testing.T) {
	var guardada *entities.PutawaySuggestion
	repo := &mocks.RepositoryMock{
		GetPutawaySuggestionByIDFn: func(ctx context.Context, businessID, id uint) (*entities.PutawaySuggestion, error) {
			return &entities.PutawaySuggestion{ID: id, Status: "pending"}, nil
		},
		UpdatePutawaySuggestionFn: func(ctx context.Context, s *entities.PutawaySuggestion) (*entities.PutawaySuggestion, error) {
			guardada = s
			return s, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	_, err := uc.ConfirmPutaway(context.Background(), request.ConfirmPutawayDTO{
		BusinessID:       10,
		SuggestionID:     1,
		ActualLocationID: 77,
	})

	require.NoError(t, err)
	assert.Equal(t, "confirmed", guardada.Status)
	require.NotNil(t, guardada.ActualLocationID)
	assert.Equal(t, uint(77), *guardada.ActualLocationID)
	assert.NotNil(t, guardada.ConfirmedAt)
}

func TestConfirmPutaway_Inexistente_PropagaError(t *testing.T) {
	dbErr := errors.New("not found")
	repo := &mocks.RepositoryMock{
		GetPutawaySuggestionByIDFn: func(ctx context.Context, businessID, id uint) (*entities.PutawaySuggestion, error) {
			return nil, dbErr
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.ConfirmPutaway(context.Background(), request.ConfirmPutawayDTO{BusinessID: 10, SuggestionID: 1})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestListPutawaySuggestions_NormalizaPaginacion(t *testing.T) {
	var recibido dtos.ListPutawaySuggestionsParams
	repo := &mocks.RepositoryMock{
		ListPutawaySuggestionsFn: func(ctx context.Context, params dtos.ListPutawaySuggestionsParams) ([]entities.PutawaySuggestion, int64, error) {
			recibido = params
			return nil, 0, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	_, _, err := uc.ListPutawaySuggestions(context.Background(), dtos.ListPutawaySuggestionsParams{Page: 0, PageSize: 700})

	require.NoError(t, err)
	assert.Equal(t, 1, recibido.Page)
	assert.Equal(t, 10, recibido.PageSize)
}

func TestNoPutawayRuleFound(t *testing.T) {
	uc := &useCase{}

	assert.True(t, uc.NoPutawayRuleFound(domainerrors.ErrNoPutawayRule))
	assert.False(t, uc.NoPutawayRuleFound(errors.New("otro error")))
	assert.False(t, uc.NoPutawayRuleFound(nil))
}

func TestCreateReplenishmentTask_CantidadInvalida_RetornaError(t *testing.T) {
	uc := buildUseCase(&mocks.RepositoryMock{}, nil, nil)

	for _, qty := range []int{0, -7} {
		got, err := uc.CreateReplenishmentTask(context.Background(), request.CreateReplenishmentTaskDTO{
			BusinessID: 10,
			Quantity:   qty,
		})

		assert.ErrorIs(t, err, domainerrors.ErrInvalidQuantity)
		assert.Nil(t, got)
	}
}

func TestCreateReplenishmentTask_OrigenPorDefectoEsManual(t *testing.T) {
	var guardada *entities.ReplenishmentTask
	repo := &mocks.RepositoryMock{
		CreateReplenishmentTaskFn: func(ctx context.Context, t *entities.ReplenishmentTask) (*entities.ReplenishmentTask, error) {
			guardada = t
			return t, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	_, err := uc.CreateReplenishmentTask(context.Background(), request.CreateReplenishmentTaskDTO{
		BusinessID: 10,
		Quantity:   30,
	})

	require.NoError(t, err)
	assert.Equal(t, "manual", guardada.TriggeredBy)
	assert.Equal(t, "pending", guardada.Status)
	assert.Equal(t, 30, guardada.Quantity)
}

func TestCreateReplenishmentTask_OrigenExplicito(t *testing.T) {
	var guardada *entities.ReplenishmentTask
	repo := &mocks.RepositoryMock{
		CreateReplenishmentTaskFn: func(ctx context.Context, t *entities.ReplenishmentTask) (*entities.ReplenishmentTask, error) {
			guardada = t
			return t, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	_, err := uc.CreateReplenishmentTask(context.Background(), request.CreateReplenishmentTaskDTO{
		BusinessID:  10,
		Quantity:    30,
		TriggeredBy: "min_max",
	})

	require.NoError(t, err)
	assert.Equal(t, "min_max", guardada.TriggeredBy)
}

func TestAssignReplenishment_TareaCerrada_Rechaza(t *testing.T) {
	for _, estado := range []string{"completed", "cancelled"} {
		t.Run(estado, func(t *testing.T) {
			repo := &mocks.RepositoryMock{
				GetReplenishmentTaskByIDFn: func(ctx context.Context, businessID, id uint) (*entities.ReplenishmentTask, error) {
					return &entities.ReplenishmentTask{ID: id, Status: estado}, nil
				},
			}
			uc := buildUseCase(repo, nil, nil)

			got, err := uc.AssignReplenishment(context.Background(), request.AssignReplenishmentDTO{
				BusinessID: 10,
				TaskID:     1,
				UserID:     7,
			})

			assert.ErrorIs(t, err, domainerrors.ErrReplenishmentClosed)
			assert.Nil(t, got)
		})
	}
}

func TestAssignReplenishment_Pendiente_SeAsigna(t *testing.T) {
	var guardada *entities.ReplenishmentTask
	repo := &mocks.RepositoryMock{
		GetReplenishmentTaskByIDFn: func(ctx context.Context, businessID, id uint) (*entities.ReplenishmentTask, error) {
			return &entities.ReplenishmentTask{ID: id, Status: "pending"}, nil
		},
		UpdateReplenishmentTaskFn: func(ctx context.Context, t *entities.ReplenishmentTask) (*entities.ReplenishmentTask, error) {
			guardada = t
			return t, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	_, err := uc.AssignReplenishment(context.Background(), request.AssignReplenishmentDTO{
		BusinessID: 10,
		TaskID:     1,
		UserID:     7,
	})

	require.NoError(t, err)
	assert.Equal(t, "in_progress", guardada.Status)
	require.NotNil(t, guardada.AssignedToID)
	assert.Equal(t, uint(7), *guardada.AssignedToID)
	assert.NotNil(t, guardada.AssignedAt)
}

func TestCompleteReplenishment_TareaCerrada_Rechaza(t *testing.T) {
	for _, estado := range []string{"completed", "cancelled"} {
		repo := &mocks.RepositoryMock{
			GetReplenishmentTaskByIDFn: func(ctx context.Context, businessID, id uint) (*entities.ReplenishmentTask, error) {
				return &entities.ReplenishmentTask{ID: id, Status: estado}, nil
			},
		}
		uc := buildUseCase(repo, nil, nil)

		got, err := uc.CompleteReplenishment(context.Background(), request.CompleteReplenishmentDTO{
			BusinessID: 10,
			TaskID:     1,
		})

		assert.ErrorIs(t, err, domainerrors.ErrReplenishmentClosed)
		assert.Nil(t, got)
	}
}

func TestCompleteReplenishment_EnProgreso_SeCompleta(t *testing.T) {
	var guardada *entities.ReplenishmentTask
	repo := &mocks.RepositoryMock{
		GetReplenishmentTaskByIDFn: func(ctx context.Context, businessID, id uint) (*entities.ReplenishmentTask, error) {
			return &entities.ReplenishmentTask{ID: id, Status: "in_progress", Notes: "original"}, nil
		},
		UpdateReplenishmentTaskFn: func(ctx context.Context, t *entities.ReplenishmentTask) (*entities.ReplenishmentTask, error) {
			guardada = t
			return t, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	_, err := uc.CompleteReplenishment(context.Background(), request.CompleteReplenishmentDTO{
		BusinessID: 10,
		TaskID:     1,
		Notes:      "movido a picking",
	})

	require.NoError(t, err)
	assert.Equal(t, "completed", guardada.Status)
	assert.NotNil(t, guardada.CompletedAt)
	assert.Equal(t, "movido a picking", guardada.Notes)
}

func TestCompleteReplenishment_SinNotas_ConservaLasExistentes(t *testing.T) {
	var guardada *entities.ReplenishmentTask
	repo := &mocks.RepositoryMock{
		GetReplenishmentTaskByIDFn: func(ctx context.Context, businessID, id uint) (*entities.ReplenishmentTask, error) {
			return &entities.ReplenishmentTask{ID: id, Status: "pending", Notes: "original"}, nil
		},
		UpdateReplenishmentTaskFn: func(ctx context.Context, t *entities.ReplenishmentTask) (*entities.ReplenishmentTask, error) {
			guardada = t
			return t, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	_, err := uc.CompleteReplenishment(context.Background(), request.CompleteReplenishmentDTO{BusinessID: 10, TaskID: 1})

	require.NoError(t, err)
	assert.Equal(t, "original", guardada.Notes)
}

func TestCancelReplenishment_YaCompletada_Rechaza(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetReplenishmentTaskByIDFn: func(ctx context.Context, businessID, id uint) (*entities.ReplenishmentTask, error) {
			return &entities.ReplenishmentTask{ID: id, Status: "completed"}, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.CancelReplenishment(context.Background(), 10, 1, "ya no aplica")

	assert.ErrorIs(t, err, domainerrors.ErrReplenishmentClosed)
	assert.Nil(t, got)
}

func TestCancelReplenishment_Pendiente_SeCancela(t *testing.T) {
	var guardada *entities.ReplenishmentTask
	repo := &mocks.RepositoryMock{
		GetReplenishmentTaskByIDFn: func(ctx context.Context, businessID, id uint) (*entities.ReplenishmentTask, error) {
			return &entities.ReplenishmentTask{ID: id, Status: "pending"}, nil
		},
		UpdateReplenishmentTaskFn: func(ctx context.Context, t *entities.ReplenishmentTask) (*entities.ReplenishmentTask, error) {
			guardada = t
			return t, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	_, err := uc.CancelReplenishment(context.Background(), 10, 1, "ya no aplica")

	require.NoError(t, err)
	assert.Equal(t, "cancelled", guardada.Status)
	assert.Equal(t, "ya no aplica", guardada.Notes)
}

func TestCancelReplenishment_YaCancelada_SePuedeRecancelar(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetReplenishmentTaskByIDFn: func(ctx context.Context, businessID, id uint) (*entities.ReplenishmentTask, error) {
			return &entities.ReplenishmentTask{ID: id, Status: "cancelled"}, nil
		},
		UpdateReplenishmentTaskFn: func(ctx context.Context, t *entities.ReplenishmentTask) (*entities.ReplenishmentTask, error) {
			return t, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.CancelReplenishment(context.Background(), 10, 1, "")

	require.NoError(t, err, "cancelar solo bloquea sobre completadas, no sobre canceladas")
	assert.Equal(t, "cancelled", got.Status)
}

func TestDetectReplenishmentNeeds_CreaTareasParaCadaCandidato(t *testing.T) {
	repo := &mocks.RepositoryMock{
		DetectReplenishmentCandidatesFn: func(ctx context.Context, businessID uint) ([]entities.ReplenishmentTask, error) {
			return []entities.ReplenishmentTask{{ProductID: "prod-1"}, {ProductID: "prod-2"}}, nil
		},
		CreateReplenishmentTaskFn: func(ctx context.Context, t *entities.ReplenishmentTask) (*entities.ReplenishmentTask, error) {
			return t, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.DetectReplenishmentNeeds(context.Background(), 10)

	require.NoError(t, err)
	assert.Equal(t, 2, got.Created)
	assert.Len(t, got.Tasks, 2)
}

func TestDetectReplenishmentNeeds_UnFalloNoAbortaElResto(t *testing.T) {
	repo := &mocks.RepositoryMock{
		DetectReplenishmentCandidatesFn: func(ctx context.Context, businessID uint) ([]entities.ReplenishmentTask, error) {
			return []entities.ReplenishmentTask{{ProductID: "prod-1"}, {ProductID: "prod-2"}}, nil
		},
		CreateReplenishmentTaskFn: func(ctx context.Context, t *entities.ReplenishmentTask) (*entities.ReplenishmentTask, error) {
			if t.ProductID == "prod-1" {
				return nil, errors.New("db down")
			}
			return t, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.DetectReplenishmentNeeds(context.Background(), 10)

	require.NoError(t, err)
	assert.Equal(t, 1, got.Created)
}

func TestDetectReplenishmentNeeds_FallaDeteccion_PropagaError(t *testing.T) {
	dbErr := errors.New("db down")
	repo := &mocks.RepositoryMock{
		DetectReplenishmentCandidatesFn: func(ctx context.Context, businessID uint) ([]entities.ReplenishmentTask, error) {
			return nil, dbErr
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.DetectReplenishmentNeeds(context.Background(), 10)

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestRunSlotting_PeriodoPorDefectoEs30d(t *testing.T) {
	var periodoUsado string
	repo := &mocks.RepositoryMock{
		ComputeVelocitiesFn: func(ctx context.Context, businessID, warehouseID uint, period string) error {
			periodoUsado = period
			return nil
		},
		ListVelocitiesFn: func(ctx context.Context, params dtos.ListVelocityParams) ([]entities.ProductVelocity, error) {
			return []entities.ProductVelocity{{ProductID: "prod-1"}}, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.RunSlotting(context.Background(), request.RunSlottingDTO{BusinessID: 10, WarehouseID: 5})

	require.NoError(t, err)
	assert.Equal(t, "30d", periodoUsado)
	assert.Equal(t, "30d", got.Period)
	assert.Equal(t, 1, got.TotalScanned)
}

func TestRunSlotting_PeriodoExplicito(t *testing.T) {
	var periodoUsado string
	repo := &mocks.RepositoryMock{
		ComputeVelocitiesFn: func(ctx context.Context, businessID, warehouseID uint, period string) error {
			periodoUsado = period
			return nil
		},
		ListVelocitiesFn: func(ctx context.Context, params dtos.ListVelocityParams) ([]entities.ProductVelocity, error) {
			return nil, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	_, err := uc.RunSlotting(context.Background(), request.RunSlottingDTO{
		BusinessID: 10,
		Period:     "90d",
	})

	require.NoError(t, err)
	assert.Equal(t, "90d", periodoUsado)
}

func TestRunSlotting_FallaComputo_PropagaError(t *testing.T) {
	dbErr := errors.New("db down")
	repo := &mocks.RepositoryMock{
		ComputeVelocitiesFn: func(ctx context.Context, businessID, warehouseID uint, period string) error {
			return dbErr
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.RunSlotting(context.Background(), request.RunSlottingDTO{BusinessID: 10})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestRunSlotting_FallaListar_PropagaError(t *testing.T) {
	dbErr := errors.New("db down")
	repo := &mocks.RepositoryMock{
		ComputeVelocitiesFn: func(ctx context.Context, businessID, warehouseID uint, period string) error {
			return nil
		},
		ListVelocitiesFn: func(ctx context.Context, params dtos.ListVelocityParams) ([]entities.ProductVelocity, error) {
			return nil, dbErr
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.RunSlotting(context.Background(), request.RunSlottingDTO{BusinessID: 10})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestListVelocities_DelegaEnRepositorio(t *testing.T) {
	repo := &mocks.RepositoryMock{
		ListVelocitiesFn: func(ctx context.Context, params dtos.ListVelocityParams) ([]entities.ProductVelocity, error) {
			return []entities.ProductVelocity{{ProductID: "prod-1"}, {ProductID: "prod-2"}}, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.ListVelocities(context.Background(), dtos.ListVelocityParams{BusinessID: 10})

	require.NoError(t, err)
	assert.Len(t, got, 2)
}
