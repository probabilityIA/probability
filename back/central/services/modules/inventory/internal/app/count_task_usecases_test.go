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

func planDeConteo() *entities.CycleCountPlan {
	return &entities.CycleCountPlan{
		ID:          1,
		BusinessID:  10,
		WarehouseID: 5,
		Strategy:    "abc",
	}
}

func TestGenerateCountTask_UsaLaEstrategiaDelPlanSiNoHayScope(t *testing.T) {
	var creada *entities.CycleCountTask
	repo := &mocks.RepositoryMock{
		GetCountPlanByIDFn: func(ctx context.Context, businessID, id uint) (*entities.CycleCountPlan, error) {
			return planDeConteo(), nil
		},
		CreateCountTaskFn: func(ctx context.Context, t *entities.CycleCountTask) (*entities.CycleCountTask, error) {
			creada = t
			t.ID = 100
			return t, nil
		},
		GenerateCountLinesForTaskFn: func(ctx context.Context, task *entities.CycleCountTask, strategy string) ([]entities.CycleCountLine, error) {
			return []entities.CycleCountLine{{ID: 1}, {ID: 2}}, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.GenerateCountTask(context.Background(), request.GenerateCountTaskDTO{
		BusinessID: 10,
		PlanID:     1,
	})

	require.NoError(t, err)
	require.NotNil(t, got)
	require.NotNil(t, creada)
	assert.Equal(t, "abc", creada.ScopeType)
	assert.Equal(t, "pending", creada.Status)
	assert.Equal(t, uint(5), creada.WarehouseID)
	assert.Len(t, got.Lines, 2)
}

func TestGenerateCountTask_ScopeExplicitoGana(t *testing.T) {
	var creada *entities.CycleCountTask
	repo := &mocks.RepositoryMock{
		GetCountPlanByIDFn: func(ctx context.Context, businessID, id uint) (*entities.CycleCountPlan, error) {
			return planDeConteo(), nil
		},
		CreateCountTaskFn: func(ctx context.Context, t *entities.CycleCountTask) (*entities.CycleCountTask, error) {
			creada = t
			return t, nil
		},
		GenerateCountLinesForTaskFn: func(ctx context.Context, task *entities.CycleCountTask, strategy string) ([]entities.CycleCountLine, error) {
			return nil, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	_, err := uc.GenerateCountTask(context.Background(), request.GenerateCountTaskDTO{
		BusinessID: 10,
		PlanID:     1,
		ScopeType:  "location",
	})

	require.NoError(t, err)
	assert.Equal(t, "location", creada.ScopeType)
}

func TestGenerateCountTask_PlanInexistente_PropagaError(t *testing.T) {
	dbErr := errors.New("not found")
	repo := &mocks.RepositoryMock{
		GetCountPlanByIDFn: func(ctx context.Context, businessID, id uint) (*entities.CycleCountPlan, error) {
			return nil, dbErr
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.GenerateCountTask(context.Background(), request.GenerateCountTaskDTO{BusinessID: 10, PlanID: 1})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestGenerateCountTask_FallaGenerarLineas_PropagaError(t *testing.T) {
	dbErr := errors.New("db down")
	repo := &mocks.RepositoryMock{
		GetCountPlanByIDFn: func(ctx context.Context, businessID, id uint) (*entities.CycleCountPlan, error) {
			return planDeConteo(), nil
		},
		CreateCountTaskFn: func(ctx context.Context, t *entities.CycleCountTask) (*entities.CycleCountTask, error) {
			return t, nil
		},
		GenerateCountLinesForTaskFn: func(ctx context.Context, task *entities.CycleCountTask, strategy string) ([]entities.CycleCountLine, error) {
			return nil, dbErr
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.GenerateCountTask(context.Background(), request.GenerateCountTaskDTO{BusinessID: 10, PlanID: 1})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestStartCountTask_TareaCerrada_Rechaza(t *testing.T) {
	for _, estado := range []string{"completed", "cancelled"} {
		t.Run(estado, func(t *testing.T) {
			actualizada := false
			repo := &mocks.RepositoryMock{
				GetCountTaskByIDFn: func(ctx context.Context, businessID, id uint) (*entities.CycleCountTask, error) {
					return &entities.CycleCountTask{ID: id, Status: estado}, nil
				},
				UpdateCountTaskFn: func(ctx context.Context, t *entities.CycleCountTask) (*entities.CycleCountTask, error) {
					actualizada = true
					return t, nil
				},
			}
			uc := buildUseCase(repo, nil, nil)

			got, err := uc.StartCountTask(context.Background(), request.StartCountTaskDTO{
				BusinessID: 10,
				TaskID:     1,
				UserID:     7,
			})

			assert.ErrorIs(t, err, domainerrors.ErrCountTaskClosed)
			assert.Nil(t, got)
			assert.False(t, actualizada)
		})
	}
}

func TestStartCountTask_Exitoso_AsignaYMarcaInicio(t *testing.T) {
	var guardada *entities.CycleCountTask
	repo := &mocks.RepositoryMock{
		GetCountTaskByIDFn: func(ctx context.Context, businessID, id uint) (*entities.CycleCountTask, error) {
			return &entities.CycleCountTask{ID: id, Status: "pending"}, nil
		},
		UpdateCountTaskFn: func(ctx context.Context, t *entities.CycleCountTask) (*entities.CycleCountTask, error) {
			guardada = t
			return t, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	_, err := uc.StartCountTask(context.Background(), request.StartCountTaskDTO{
		BusinessID: 10,
		TaskID:     1,
		UserID:     7,
	})

	require.NoError(t, err)
	assert.Equal(t, "in_progress", guardada.Status)
	require.NotNil(t, guardada.AssignedToID)
	assert.Equal(t, uint(7), *guardada.AssignedToID)
	assert.NotNil(t, guardada.StartedAt)
}

func TestFinishCountTask_YaCompletada_Rechaza(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetCountTaskByIDFn: func(ctx context.Context, businessID, id uint) (*entities.CycleCountTask, error) {
			return &entities.CycleCountTask{ID: id, Status: "completed"}, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.FinishCountTask(context.Background(), 10, 1)

	assert.ErrorIs(t, err, domainerrors.ErrCountTaskClosed)
	assert.Nil(t, got)
}

func TestFinishCountTask_Exitoso_MarcaFin(t *testing.T) {
	var guardada *entities.CycleCountTask
	repo := &mocks.RepositoryMock{
		GetCountTaskByIDFn: func(ctx context.Context, businessID, id uint) (*entities.CycleCountTask, error) {
			return &entities.CycleCountTask{ID: id, Status: "in_progress"}, nil
		},
		UpdateCountTaskFn: func(ctx context.Context, t *entities.CycleCountTask) (*entities.CycleCountTask, error) {
			guardada = t
			return t, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	_, err := uc.FinishCountTask(context.Background(), 10, 1)

	require.NoError(t, err)
	assert.Equal(t, "completed", guardada.Status)
	assert.NotNil(t, guardada.FinishedAt)
}

func TestSubmitCountLine_LineaYaEnviada_Rechaza(t *testing.T) {
	for _, estado := range []string{"submitted", "resolved"} {
		t.Run(estado, func(t *testing.T) {
			repo := &mocks.RepositoryMock{
				GetCountLineByIDFn: func(ctx context.Context, businessID, id uint) (*entities.CycleCountLine, error) {
					return &entities.CycleCountLine{ID: id, Status: estado}, nil
				},
			}
			uc := buildUseCase(repo, nil, nil)

			got, err := uc.SubmitCountLine(context.Background(), request.SubmitCountLineDTO{
				BusinessID: 10,
				LineID:     1,
				CountedQty: 5,
			})

			assert.ErrorIs(t, err, domainerrors.ErrCountLineSubmitted)
			assert.Nil(t, got)
		})
	}
}

func TestSubmitCountLine_SinDiferencia_QuedaResuelta(t *testing.T) {
	discrepanciaCreada := false
	var guardada *entities.CycleCountLine
	repo := &mocks.RepositoryMock{
		GetCountLineByIDFn: func(ctx context.Context, businessID, id uint) (*entities.CycleCountLine, error) {
			return &entities.CycleCountLine{ID: id, Status: "pending", ExpectedQty: 10}, nil
		},
		UpdateCountLineFn: func(ctx context.Context, line *entities.CycleCountLine) (*entities.CycleCountLine, error) {
			guardada = line
			return line, nil
		},
		CreateDiscrepancyFn: func(ctx context.Context, d *entities.InventoryDiscrepancy) (*entities.InventoryDiscrepancy, error) {
			discrepanciaCreada = true
			return d, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.SubmitCountLine(context.Background(), request.SubmitCountLineDTO{
		BusinessID: 10,
		LineID:     1,
		CountedQty: 10,
	})

	require.NoError(t, err)
	assert.Equal(t, "resolved", guardada.Status)
	assert.Equal(t, 0, guardada.Variance)
	assert.False(t, discrepanciaCreada, "sin diferencia no se abre discrepancia")
	assert.Nil(t, got.Discrepancy)
}

func TestSubmitCountLine_ConDiferencia_AbreDiscrepancia(t *testing.T) {
	var creada *entities.InventoryDiscrepancy
	var guardada *entities.CycleCountLine
	repo := &mocks.RepositoryMock{
		GetCountLineByIDFn: func(ctx context.Context, businessID, id uint) (*entities.CycleCountLine, error) {
			return &entities.CycleCountLine{ID: id, TaskID: 55, BusinessID: 10, Status: "pending", ExpectedQty: 10}, nil
		},
		UpdateCountLineFn: func(ctx context.Context, line *entities.CycleCountLine) (*entities.CycleCountLine, error) {
			guardada = line
			return line, nil
		},
		CreateDiscrepancyFn: func(ctx context.Context, d *entities.InventoryDiscrepancy) (*entities.InventoryDiscrepancy, error) {
			creada = d
			return d, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.SubmitCountLine(context.Background(), request.SubmitCountLineDTO{
		BusinessID: 10,
		LineID:     1,
		CountedQty: 3,
	})

	require.NoError(t, err)
	assert.Equal(t, "submitted", guardada.Status)
	assert.Equal(t, -7, guardada.Variance, "contado menos esperado")
	require.NotNil(t, creada)
	assert.Equal(t, uint(55), creada.TaskID)
	assert.Equal(t, "open", creada.Status)
	assert.NotNil(t, got.Discrepancy)
}

func TestSubmitCountLine_SobranteTambienAbreDiscrepancia(t *testing.T) {
	var guardada *entities.CycleCountLine
	repo := &mocks.RepositoryMock{
		GetCountLineByIDFn: func(ctx context.Context, businessID, id uint) (*entities.CycleCountLine, error) {
			return &entities.CycleCountLine{ID: id, Status: "pending", ExpectedQty: 10}, nil
		},
		UpdateCountLineFn: func(ctx context.Context, line *entities.CycleCountLine) (*entities.CycleCountLine, error) {
			guardada = line
			return line, nil
		},
		CreateDiscrepancyFn: func(ctx context.Context, d *entities.InventoryDiscrepancy) (*entities.InventoryDiscrepancy, error) {
			return d, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	_, err := uc.SubmitCountLine(context.Background(), request.SubmitCountLineDTO{
		BusinessID: 10,
		LineID:     1,
		CountedQty: 15,
	})

	require.NoError(t, err)
	assert.Equal(t, 5, guardada.Variance)
	assert.Equal(t, "submitted", guardada.Status)
}

func TestSubmitCountLine_FallaCrearDiscrepancia_NoRompeElConteo(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetCountLineByIDFn: func(ctx context.Context, businessID, id uint) (*entities.CycleCountLine, error) {
			return &entities.CycleCountLine{ID: id, Status: "pending", ExpectedQty: 10}, nil
		},
		UpdateCountLineFn: func(ctx context.Context, line *entities.CycleCountLine) (*entities.CycleCountLine, error) {
			return line, nil
		},
		CreateDiscrepancyFn: func(ctx context.Context, d *entities.InventoryDiscrepancy) (*entities.InventoryDiscrepancy, error) {
			return nil, errors.New("db down")
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.SubmitCountLine(context.Background(), request.SubmitCountLineDTO{
		BusinessID: 10,
		LineID:     1,
		CountedQty: 3,
	})

	require.NoError(t, err, "el conteo se guarda aunque falle abrir la discrepancia")
	assert.Nil(t, got.Discrepancy)
}

func TestSubmitCountLine_FallaGuardarLinea_PropagaError(t *testing.T) {
	dbErr := errors.New("db down")
	repo := &mocks.RepositoryMock{
		GetCountLineByIDFn: func(ctx context.Context, businessID, id uint) (*entities.CycleCountLine, error) {
			return &entities.CycleCountLine{ID: id, Status: "pending", ExpectedQty: 10}, nil
		},
		UpdateCountLineFn: func(ctx context.Context, line *entities.CycleCountLine) (*entities.CycleCountLine, error) {
			return nil, dbErr
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.SubmitCountLine(context.Background(), request.SubmitCountLineDTO{BusinessID: 10, LineID: 1})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestListCountTasksYLineas_NormalizanPaginacion(t *testing.T) {
	var tareas dtos.ListCycleCountTasksParams
	var lineas dtos.ListCycleCountLinesParams
	repo := &mocks.RepositoryMock{
		ListCountTasksFn: func(ctx context.Context, params dtos.ListCycleCountTasksParams) ([]entities.CycleCountTask, int64, error) {
			tareas = params
			return nil, 0, nil
		},
		ListCountLinesFn: func(ctx context.Context, params dtos.ListCycleCountLinesParams) ([]entities.CycleCountLine, int64, error) {
			lineas = params
			return nil, 0, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	_, _, err := uc.ListCountTasks(context.Background(), dtos.ListCycleCountTasksParams{Page: 0, PageSize: 500})
	require.NoError(t, err)
	assert.Equal(t, 1, tareas.Page)
	assert.Equal(t, 10, tareas.PageSize)

	_, _, err = uc.ListCountLines(context.Background(), dtos.ListCycleCountLinesParams{Page: -1, PageSize: 0})
	require.NoError(t, err)
	assert.Equal(t, 1, lineas.Page)
	assert.Equal(t, 10, lineas.PageSize)
}

func TestGetCountTask_DelegaEnRepositorio(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetCountTaskByIDFn: func(ctx context.Context, businessID, id uint) (*entities.CycleCountTask, error) {
			return &entities.CycleCountTask{ID: id}, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.GetCountTask(context.Background(), 10, 3)

	require.NoError(t, err)
	assert.Equal(t, uint(3), got.ID)
}
