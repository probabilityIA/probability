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

func repoSerialBase() *mocks.RepositoryMock {
	return &mocks.RepositoryMock{
		GetProductByIDFn: func(ctx context.Context, productID string, businessID uint) (string, string, bool, error) {
			return "SKU-1", "Producto", true, nil
		},
		SerialExistsFn: func(ctx context.Context, businessID uint, productID, serial string, excludeID *uint) (bool, error) {
			return false, nil
		},
		GetInventoryStateByCodeFn: func(ctx context.Context, code string) (*entities.InventoryState, error) {
			return &entities.InventoryState{ID: 3, Code: code}, nil
		},
		CreateSerialFn: func(ctx context.Context, serial *entities.InventorySerial) (*entities.InventorySerial, error) {
			return serial, nil
		},
	}
}

func TestCreateSerial_ProductoInexistente_RetornaError(t *testing.T) {
	repo := repoSerialBase()
	repo.GetProductByIDFn = func(ctx context.Context, productID string, businessID uint) (string, string, bool, error) {
		return "", "", false, errors.New("not found")
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.CreateSerial(context.Background(), request.CreateSerialDTO{ProductID: "prod-1", BusinessID: 10})

	assert.ErrorIs(t, err, domainerrors.ErrProductNotFound)
	assert.Nil(t, got)
}

func TestCreateSerial_SerialDuplicado_RetornaError(t *testing.T) {
	repo := repoSerialBase()
	repo.SerialExistsFn = func(ctx context.Context, businessID uint, productID, serial string, excludeID *uint) (bool, error) {
		return true, nil
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.CreateSerial(context.Background(), request.CreateSerialDTO{
		ProductID:    "prod-1",
		BusinessID:   10,
		SerialNumber: "SN-001",
	})

	assert.ErrorIs(t, err, domainerrors.ErrDuplicateSerial)
	assert.Nil(t, got)
}

func TestCreateSerial_FallaVerificarDuplicado_PropagaError(t *testing.T) {
	dbErr := errors.New("db down")
	repo := repoSerialBase()
	repo.SerialExistsFn = func(ctx context.Context, businessID uint, productID, serial string, excludeID *uint) (bool, error) {
		return false, dbErr
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.CreateSerial(context.Background(), request.CreateSerialDTO{ProductID: "prod-1", BusinessID: 10})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestCreateSerial_SinEstado_NoAsignaEstado(t *testing.T) {
	consultoEstado := false
	var guardado *entities.InventorySerial
	repo := repoSerialBase()
	repo.GetInventoryStateByCodeFn = func(ctx context.Context, code string) (*entities.InventoryState, error) {
		consultoEstado = true
		return &entities.InventoryState{ID: 3}, nil
	}
	repo.CreateSerialFn = func(ctx context.Context, serial *entities.InventorySerial) (*entities.InventorySerial, error) {
		guardado = serial
		return serial, nil
	}
	uc := buildUseCase(repo, nil, nil)

	_, err := uc.CreateSerial(context.Background(), request.CreateSerialDTO{
		ProductID:    "prod-1",
		BusinessID:   10,
		SerialNumber: "SN-001",
	})

	require.NoError(t, err)
	assert.False(t, consultoEstado)
	assert.Nil(t, guardado.CurrentStateID)
	assert.NotNil(t, guardado.ReceivedAt, "la fecha de recepcion se sella al crear")
}

func TestCreateSerial_ConEstado_LoResuelvePorCodigo(t *testing.T) {
	var codigoPedido string
	var guardado *entities.InventorySerial
	repo := repoSerialBase()
	repo.GetInventoryStateByCodeFn = func(ctx context.Context, code string) (*entities.InventoryState, error) {
		codigoPedido = code
		return &entities.InventoryState{ID: 3, Code: code}, nil
	}
	repo.CreateSerialFn = func(ctx context.Context, serial *entities.InventorySerial) (*entities.InventorySerial, error) {
		guardado = serial
		return serial, nil
	}
	uc := buildUseCase(repo, nil, nil)

	_, err := uc.CreateSerial(context.Background(), request.CreateSerialDTO{
		ProductID:    "prod-1",
		BusinessID:   10,
		SerialNumber: "SN-001",
		StateCode:    "quarantine",
	})

	require.NoError(t, err)
	assert.Equal(t, "quarantine", codigoPedido)
	require.NotNil(t, guardado.CurrentStateID)
	assert.Equal(t, uint(3), *guardado.CurrentStateID)
}

func TestCreateSerial_EstadoInexistente_PropagaError(t *testing.T) {
	dbErr := errors.New("estado no existe")
	repo := repoSerialBase()
	repo.GetInventoryStateByCodeFn = func(ctx context.Context, code string) (*entities.InventoryState, error) {
		return nil, dbErr
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.CreateSerial(context.Background(), request.CreateSerialDTO{
		ProductID:  "prod-1",
		BusinessID: 10,
		StateCode:  "inventado",
	})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestGetSerial_DelegaEnRepositorio(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetSerialByIDFn: func(ctx context.Context, businessID, serialID uint) (*entities.InventorySerial, error) {
			return &entities.InventorySerial{ID: serialID}, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.GetSerial(context.Background(), 10, 5)

	require.NoError(t, err)
	assert.Equal(t, uint(5), got.ID)
}

func TestListSerials_NormalizaPaginacion(t *testing.T) {
	cases := []struct {
		page         int
		pageSize     int
		wantPage     int
		wantPageSize int
	}{
		{page: 0, pageSize: 10, wantPage: 1, wantPageSize: 10},
		{page: 1, pageSize: 0, wantPage: 1, wantPageSize: 10},
		{page: 1, pageSize: 600, wantPage: 1, wantPageSize: 10},
		{page: 2, pageSize: 60, wantPage: 2, wantPageSize: 60},
	}

	for _, tc := range cases {
		var recibido dtos.ListSerialsParams
		repo := &mocks.RepositoryMock{
			ListSerialsFn: func(ctx context.Context, params dtos.ListSerialsParams) ([]entities.InventorySerial, int64, error) {
				recibido = params
				return nil, 0, nil
			},
		}
		uc := buildUseCase(repo, nil, nil)

		_, _, err := uc.ListSerials(context.Background(), dtos.ListSerialsParams{Page: tc.page, PageSize: tc.pageSize})

		require.NoError(t, err)
		assert.Equal(t, tc.wantPage, recibido.Page)
		assert.Equal(t, tc.wantPageSize, recibido.PageSize)
	}
}

func TestUpdateSerial_ActualizaLoteUbicacionYEstado(t *testing.T) {
	lote := uint(7)
	ubicacion := uint(42)

	var guardado *entities.InventorySerial
	repo := &mocks.RepositoryMock{
		GetSerialByIDFn: func(ctx context.Context, businessID, serialID uint) (*entities.InventorySerial, error) {
			return &entities.InventorySerial{ID: serialID}, nil
		},
		GetInventoryStateByCodeFn: func(ctx context.Context, code string) (*entities.InventoryState, error) {
			return &entities.InventoryState{ID: 9, Code: code}, nil
		},
		UpdateSerialFn: func(ctx context.Context, serial *entities.InventorySerial) (*entities.InventorySerial, error) {
			guardado = serial
			return serial, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	_, err := uc.UpdateSerial(context.Background(), request.UpdateSerialDTO{
		ID:         1,
		BusinessID: 10,
		LotID:      &lote,
		LocationID: &ubicacion,
		StateCode:  "damaged",
	})

	require.NoError(t, err)
	assert.Equal(t, uint(7), *guardado.LotID)
	assert.Equal(t, uint(42), *guardado.CurrentLocationID)
	assert.Equal(t, uint(9), *guardado.CurrentStateID)
}

func TestUpdateSerial_SinCambios_ConservaLoExistente(t *testing.T) {
	lote := uint(7)
	estado := uint(3)

	var guardado *entities.InventorySerial
	repo := &mocks.RepositoryMock{
		GetSerialByIDFn: func(ctx context.Context, businessID, serialID uint) (*entities.InventorySerial, error) {
			return &entities.InventorySerial{ID: serialID, LotID: &lote, CurrentStateID: &estado}, nil
		},
		UpdateSerialFn: func(ctx context.Context, serial *entities.InventorySerial) (*entities.InventorySerial, error) {
			guardado = serial
			return serial, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	_, err := uc.UpdateSerial(context.Background(), request.UpdateSerialDTO{ID: 1, BusinessID: 10})

	require.NoError(t, err)
	assert.Equal(t, uint(7), *guardado.LotID)
	assert.Equal(t, uint(3), *guardado.CurrentStateID)
}

func TestUpdateSerial_SerialInexistente_PropagaError(t *testing.T) {
	dbErr := errors.New("not found")
	repo := &mocks.RepositoryMock{
		GetSerialByIDFn: func(ctx context.Context, businessID, serialID uint) (*entities.InventorySerial, error) {
			return nil, dbErr
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.UpdateSerial(context.Background(), request.UpdateSerialDTO{ID: 1, BusinessID: 10})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestUpdateSerial_EstadoInexistente_PropagaError(t *testing.T) {
	dbErr := errors.New("estado no existe")
	repo := &mocks.RepositoryMock{
		GetSerialByIDFn: func(ctx context.Context, businessID, serialID uint) (*entities.InventorySerial, error) {
			return &entities.InventorySerial{ID: serialID}, nil
		},
		GetInventoryStateByCodeFn: func(ctx context.Context, code string) (*entities.InventoryState, error) {
			return nil, dbErr
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.UpdateSerial(context.Background(), request.UpdateSerialDTO{
		ID:         1,
		BusinessID: 10,
		StateCode:  "inventado",
	})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestListInventoryStates_DelegaEnRepositorio(t *testing.T) {
	repo := &mocks.RepositoryMock{
		ListInventoryStatesFn: func(ctx context.Context) ([]entities.InventoryState, error) {
			return []entities.InventoryState{{Code: "available"}, {Code: "quarantine"}}, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.ListInventoryStates(context.Background())

	require.NoError(t, err)
	assert.Len(t, got, 2)
}

func TestChangeInventoryState_CantidadInvalida_RetornaError(t *testing.T) {
	uc := buildUseCase(&mocks.RepositoryMock{}, nil, nil)

	for _, qty := range []int{0, -4} {
		got, err := uc.ChangeInventoryState(context.Background(), request.ChangeInventoryStateDTO{
			BusinessID:    10,
			Quantity:      qty,
			FromStateCode: "available",
			ToStateCode:   "quarantine",
		})

		assert.ErrorIs(t, err, domainerrors.ErrInvalidQuantity)
		assert.Nil(t, got)
	}
}

func TestChangeInventoryState_MismoEstado_Rechaza(t *testing.T) {
	ejecutado := false
	repo := &mocks.RepositoryMock{
		ChangeStateTxFn: func(ctx context.Context, params dtos.ChangeInventoryStateTxParams) (*entities.StockMovement, error) {
			ejecutado = true
			return nil, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.ChangeInventoryState(context.Background(), request.ChangeInventoryStateDTO{
		BusinessID:    10,
		Quantity:      5,
		FromStateCode: "available",
		ToStateCode:   "available",
	})

	assert.ErrorIs(t, err, domainerrors.ErrStateTransition)
	assert.Nil(t, got)
	assert.False(t, ejecutado, "una transicion al mismo estado no debe generar movimiento")
}

func TestChangeInventoryState_Exitoso_PropagaLosParametros(t *testing.T) {
	var recibido dtos.ChangeInventoryStateTxParams
	repo := &mocks.RepositoryMock{
		ChangeStateTxFn: func(ctx context.Context, params dtos.ChangeInventoryStateTxParams) (*entities.StockMovement, error) {
			recibido = params
			return &entities.StockMovement{ID: 100}, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.ChangeInventoryState(context.Background(), request.ChangeInventoryStateDTO{
		BusinessID:    10,
		LevelID:       55,
		Quantity:      5,
		FromStateCode: "available",
		ToStateCode:   "damaged",
		Reason:        "producto averiado en recepcion",
	})

	require.NoError(t, err)
	assert.Equal(t, uint(100), got.ID)
	assert.Equal(t, uint(55), recibido.LevelID)
	assert.Equal(t, "available", recibido.FromStateCode)
	assert.Equal(t, "damaged", recibido.ToStateCode)
	assert.Equal(t, 5, recibido.Quantity)
	assert.Equal(t, "producto averiado en recepcion", recibido.Reason)
}

func TestChangeInventoryState_FallaTransaccion_PropagaError(t *testing.T) {
	dbErr := errors.New("tx failed")
	repo := &mocks.RepositoryMock{
		ChangeStateTxFn: func(ctx context.Context, params dtos.ChangeInventoryStateTxParams) (*entities.StockMovement, error) {
			return nil, dbErr
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.ChangeInventoryState(context.Background(), request.ChangeInventoryStateDTO{
		BusinessID:    10,
		Quantity:      5,
		FromStateCode: "available",
		ToStateCode:   "damaged",
	})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestScan_CodigoResuelto(t *testing.T) {
	var evento *entities.ScanEvent
	repo := &mocks.RepositoryMock{
		ResolveScanCodeFn: func(ctx context.Context, businessID uint, code string) (*entities.ScanResolution, error) {
			return &entities.ScanResolution{CodeType: "lpn"}, nil
		},
		RecordScanEventFn: func(ctx context.Context, e *entities.ScanEvent) (*entities.ScanEvent, error) {
			evento = e
			return e, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.Scan(context.Background(), request.ScanDTO{
		BusinessID: 10,
		Code:       "LPN-001",
		Action:     "pick",
		DeviceID:   "scanner-1",
	})

	require.NoError(t, err)
	assert.True(t, got.Resolved)
	require.NotNil(t, evento)
	assert.Equal(t, "lpn", evento.CodeType)
	assert.Equal(t, "LPN-001", evento.ScannedCode)
	assert.Equal(t, "pick", evento.Action)
	assert.Equal(t, "scanner-1", evento.DeviceID)
}

func TestScan_CodigoNoResuelto_QuedaComoUnknown(t *testing.T) {
	var evento *entities.ScanEvent
	repo := &mocks.RepositoryMock{
		ResolveScanCodeFn: func(ctx context.Context, businessID uint, code string) (*entities.ScanResolution, error) {
			return nil, nil
		},
		RecordScanEventFn: func(ctx context.Context, e *entities.ScanEvent) (*entities.ScanEvent, error) {
			evento = e
			return e, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.Scan(context.Background(), request.ScanDTO{BusinessID: 10, Code: "XYZ"})

	require.NoError(t, err)
	assert.False(t, got.Resolved)
	assert.Equal(t, "unknown", evento.CodeType)
	assert.Nil(t, got.Resolution)
}

func TestScan_FallaResolver_PropagaError(t *testing.T) {
	dbErr := errors.New("db down")
	repo := &mocks.RepositoryMock{
		ResolveScanCodeFn: func(ctx context.Context, businessID uint, code string) (*entities.ScanResolution, error) {
			return nil, dbErr
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.Scan(context.Background(), request.ScanDTO{BusinessID: 10, Code: "XYZ"})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestScan_FallaRegistrarEvento_NoRompeElEscaneo(t *testing.T) {
	repo := &mocks.RepositoryMock{
		ResolveScanCodeFn: func(ctx context.Context, businessID uint, code string) (*entities.ScanResolution, error) {
			return &entities.ScanResolution{CodeType: "serial"}, nil
		},
		RecordScanEventFn: func(ctx context.Context, e *entities.ScanEvent) (*entities.ScanEvent, error) {
			return nil, errors.New("db down")
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.Scan(context.Background(), request.ScanDTO{BusinessID: 10, Code: "SN-001"})

	require.NoError(t, err, "no poder auditar el escaneo no debe bloquear al operario")
	assert.True(t, got.Resolved)
	require.NotNil(t, got.Event)
	assert.Equal(t, "serial", got.Event.CodeType)
}

func TestListReplenishmentTasks_NormalizaPaginacion(t *testing.T) {
	var recibido dtos.ListReplenishmentTasksParams
	repo := &mocks.RepositoryMock{
		ListReplenishmentTasksFn: func(ctx context.Context, params dtos.ListReplenishmentTasksParams) ([]entities.ReplenishmentTask, int64, error) {
			recibido = params
			return nil, 0, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	_, _, err := uc.ListReplenishmentTasks(context.Background(), dtos.ListReplenishmentTasksParams{Page: 0, PageSize: 900})

	require.NoError(t, err)
	assert.Equal(t, 1, recibido.Page)
	assert.Equal(t, 10, recibido.PageSize)
}
