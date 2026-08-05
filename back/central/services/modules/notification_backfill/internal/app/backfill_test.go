package app

import (
	"context"
	stderrors "errors"
	"testing"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/notification_backfill/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/notification_backfill/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/notification_backfill/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/modules/notification_backfill/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newBackfillUseCase(
	registry *mocks.RegistryMock,
	store *mocks.JobStoreMock,
	pub *mocks.ProgressPublisherMock,
	resolver *mocks.BusinessNameResolverMock,
) ports.IUseCase {
	if registry == nil {
		registry = &mocks.RegistryMock{}
	}
	if store == nil {
		store = mocks.NewJobStore()
	}
	if pub == nil {
		pub = mocks.NewProgressPublisher()
	}
	if resolver == nil {
		return New(registry, store, pub, nil, mocks.NewSilentLogger())
	}
	return New(registry, store, pub, resolver, mocks.NewSilentLogger())
}

func uintPtr(v uint) *uint { return &v }

func registryCon(sel *mocks.SelectorMock) *mocks.RegistryMock {
	return &mocks.RegistryMock{
		GetFn: func(eventCode string) (ports.IEligibilitySelector, bool) { return sel, true },
	}
}

func registrySinEvento() *mocks.RegistryMock {
	return &mocks.RegistryMock{
		GetFn: func(eventCode string) (ports.IEligibilitySelector, bool) { return nil, false },
	}
}

func esperarFin(t *testing.T, pub *mocks.ProgressPublisherMock) {
	t.Helper()
	select {
	case <-pub.Done():
	case <-time.After(3 * time.Second):
		t.Fatal("el job no termino a tiempo")
	}
}

func TestListEvents_MapeaElRegistro(t *testing.T) {
	registry := &mocks.RegistryMock{
		ListFn: func() []ports.IEligibilitySelector {
			return []ports.IEligibilitySelector{
				&mocks.SelectorMock{
					CodeFn:    func() string { return "guide_created" },
					NameFn:    func() string { return "Guia creada" },
					ChannelFn: func() string { return "whatsapp" },
				},
				&mocks.SelectorMock{
					CodeFn:    func() string { return "order_delivered" },
					NameFn:    func() string { return "Orden entregada" },
					ChannelFn: func() string { return "email" },
				},
			}
		},
	}

	got := newBackfillUseCase(registry, nil, nil, nil).ListEvents(context.Background())

	require.Len(t, got, 2)
	assert.Equal(t, "guide_created", got[0].EventCode)
	assert.Equal(t, "Guia creada", got[0].EventName)
	assert.Equal(t, "whatsapp", got[0].Channel)
	assert.Equal(t, "email", got[1].Channel)
}

func TestListEvents_RegistroVacio_DevuelveListaVaciaNoNil(t *testing.T) {
	got := newBackfillUseCase(nil, nil, nil, nil).ListEvents(context.Background())

	require.NotNil(t, got)
	assert.Empty(t, got)
}

func TestPreview_EventoNoSoportado_Rechaza(t *testing.T) {
	got, err := newBackfillUseCase(registrySinEvento(), nil, nil, nil).
		Preview(context.Background(), dtos.BackfillFilter{EventCode: "inventado"})

	require.Error(t, err)
	assert.Nil(t, got)
	assert.Contains(t, err.Error(), "event_code no soportado: inventado")
}

func TestPreview_PasaElFiltroCompletoAlSelector(t *testing.T) {
	sel := &mocks.SelectorMock{}
	filtro := dtos.BackfillFilter{
		EventCode: "guide_created", BusinessID: uintPtr(26), Days: 7, Limit: 100,
	}

	_, err := newBackfillUseCase(registryCon(sel), nil, nil, nil).
		Preview(context.Background(), filtro)

	require.NoError(t, err)
	require.Len(t, sel.PreviewFilters(), 1)
	assert.Equal(t, filtro, sel.PreviewFilters()[0])
}

func TestPreview_ErrorDelSelector_SePropaga(t *testing.T) {
	selErr := stderrors.New("query fallida")
	sel := &mocks.SelectorMock{
		PreviewFn: func(ctx context.Context, f dtos.BackfillFilter) ([]entities.Candidate, error) {
			return nil, selErr
		},
	}

	_, err := newBackfillUseCase(registryCon(sel), nil, nil, nil).
		Preview(context.Background(), dtos.BackfillFilter{EventCode: "guide_created"})

	assert.ErrorIs(t, err, selErr)
}

func TestPreview_AgrupaPorNegocio(t *testing.T) {
	sel := &mocks.SelectorMock{
		PreviewFn: func(ctx context.Context, f dtos.BackfillFilter) ([]entities.Candidate, error) {
			return []entities.Candidate{
				{OrderID: "O-1", BusinessID: 26},
				{OrderID: "O-2", BusinessID: 27},
				{OrderID: "O-3", BusinessID: 26},
			}, nil
		},
	}

	got, err := newBackfillUseCase(registryCon(sel), nil, nil, nil).
		Preview(context.Background(), dtos.BackfillFilter{EventCode: "guide_created"})

	require.NoError(t, err)
	assert.Equal(t, 3, got.TotalEligible)
	require.Len(t, got.Businesses, 2)
	assert.Equal(t, uint(26), got.Businesses[0].BusinessID, "el de mas candidatos va primero")
	assert.Equal(t, 2, got.Businesses[0].Count)
	assert.Len(t, got.Businesses[0].Orders, 2)
	assert.Equal(t, uint(27), got.Businesses[1].BusinessID)
	assert.Equal(t, 1, got.Businesses[1].Count)
}

func TestPreview_OrdenaPorCantidadYLuegoPorID(t *testing.T) {
	sel := &mocks.SelectorMock{
		PreviewFn: func(ctx context.Context, f dtos.BackfillFilter) ([]entities.Candidate, error) {
			return []entities.Candidate{
				{OrderID: "a", BusinessID: 30},
				{OrderID: "b", BusinessID: 10},
				{OrderID: "c", BusinessID: 20},
				{OrderID: "d", BusinessID: 20},
			}, nil
		},
	}

	got, err := newBackfillUseCase(registryCon(sel), nil, nil, nil).
		Preview(context.Background(), dtos.BackfillFilter{EventCode: "guide_created"})

	require.NoError(t, err)
	require.Len(t, got.Businesses, 3)
	assert.Equal(t, uint(20), got.Businesses[0].BusinessID, "2 candidatos, va primero")
	assert.Equal(t, uint(10), got.Businesses[1].BusinessID, "empate en 1: desempata el ID menor")
	assert.Equal(t, uint(30), got.Businesses[2].BusinessID)
}

func TestPreview_ResuelveLosNombresDeNegocio(t *testing.T) {
	resolver := &mocks.BusinessNameResolverMock{
		ResolveNamesFn: func(ctx context.Context, ids []uint) (map[uint]string, error) {
			return map[uint]string{26: "Demo", 27: "Otro"}, nil
		},
	}
	sel := &mocks.SelectorMock{
		PreviewFn: func(ctx context.Context, f dtos.BackfillFilter) ([]entities.Candidate, error) {
			return []entities.Candidate{{BusinessID: 26}, {BusinessID: 27}}, nil
		},
	}

	got, err := newBackfillUseCase(registryCon(sel), nil, nil, resolver).
		Preview(context.Background(), dtos.BackfillFilter{EventCode: "guide_created"})

	require.NoError(t, err)
	nombres := map[uint]string{}
	for _, b := range got.Businesses {
		nombres[b.BusinessID] = b.BusinessName
	}
	assert.Equal(t, "Demo", nombres[26])
	assert.Equal(t, "Otro", nombres[27])
	require.Len(t, resolver.ReceivedIDs, 1)
	assert.ElementsMatch(t, []uint{26, 27}, resolver.ReceivedIDs[0])
}

func TestPreview_FalloAlResolverNombres_NoAborta(t *testing.T) {
	resolver := &mocks.BusinessNameResolverMock{
		ResolveNamesFn: func(ctx context.Context, ids []uint) (map[uint]string, error) {
			return nil, stderrors.New("timeout")
		},
	}
	sel := &mocks.SelectorMock{
		PreviewFn: func(ctx context.Context, f dtos.BackfillFilter) ([]entities.Candidate, error) {
			return []entities.Candidate{{BusinessID: 26}}, nil
		},
	}

	got, err := newBackfillUseCase(registryCon(sel), nil, nil, resolver).
		Preview(context.Background(), dtos.BackfillFilter{EventCode: "guide_created"})

	require.NoError(t, err, "el nombre es cosmetico, no debe tumbar el preview")
	require.Len(t, got.Businesses, 1)
	assert.Empty(t, got.Businesses[0].BusinessName)
}

func TestPreview_SinResolver_NoRompe(t *testing.T) {
	sel := &mocks.SelectorMock{
		PreviewFn: func(ctx context.Context, f dtos.BackfillFilter) ([]entities.Candidate, error) {
			return []entities.Candidate{{BusinessID: 26}}, nil
		},
	}

	got, err := newBackfillUseCase(registryCon(sel), nil, nil, nil).
		Preview(context.Background(), dtos.BackfillFilter{EventCode: "guide_created"})

	require.NoError(t, err)
	assert.Empty(t, got.Businesses[0].BusinessName)
}

func TestPreview_SinCandidatos_NoConsultaNombres(t *testing.T) {
	resolver := &mocks.BusinessNameResolverMock{}
	sel := &mocks.SelectorMock{}

	got, err := newBackfillUseCase(registryCon(sel), nil, nil, resolver).
		Preview(context.Background(), dtos.BackfillFilter{EventCode: "guide_created"})

	require.NoError(t, err)
	assert.Zero(t, got.TotalEligible)
	require.NotNil(t, got.Businesses)
	assert.Empty(t, got.Businesses)
	assert.Empty(t, resolver.ReceivedIDs)
}

func TestPreview_MapeaTodosLosCamposDelCandidato(t *testing.T) {
	sel := &mocks.SelectorMock{
		PreviewFn: func(ctx context.Context, f dtos.BackfillFilter) ([]entities.Candidate, error) {
			return []entities.Candidate{{
				OrderID: "O-1", OrderNumber: "ORD-001", BusinessID: 26,
				CustomerPhone: "3001234567", TrackingNumber: "034057376067",
				Status: "in_transit", Carrier: "Servientrega",
				CarrierLogoURL: "https://cdn/servientrega.png",
			}}, nil
		},
	}

	got, err := newBackfillUseCase(registryCon(sel), nil, nil, nil).
		Preview(context.Background(), dtos.BackfillFilter{EventCode: "guide_created"})

	require.NoError(t, err)
	require.Len(t, got.Businesses, 1)
	require.Len(t, got.Businesses[0].Orders, 1)
	o := got.Businesses[0].Orders[0]
	assert.Equal(t, "O-1", o.OrderID)
	assert.Equal(t, "ORD-001", o.OrderNumber)
	assert.Equal(t, "3001234567", o.CustomerPhone)
	assert.Equal(t, "034057376067", o.TrackingNumber)
	assert.Equal(t, "in_transit", o.Status)
	assert.Equal(t, "Servientrega", o.Carrier)
	assert.Equal(t, "https://cdn/servientrega.png", o.CarrierLogoURL)
}

func TestPreview_DevuelveElEventCodePedido(t *testing.T) {
	got, err := newBackfillUseCase(registryCon(&mocks.SelectorMock{}), nil, nil, nil).
		Preview(context.Background(), dtos.BackfillFilter{EventCode: "order_delivered"})

	require.NoError(t, err)
	assert.Equal(t, "order_delivered", got.EventCode)
}

func TestPreview_NoDespachaNada(t *testing.T) {
	sel := &mocks.SelectorMock{
		PreviewFn: func(ctx context.Context, f dtos.BackfillFilter) ([]entities.Candidate, error) {
			return []entities.Candidate{{OrderID: "O-1", BusinessID: 26}}, nil
		},
	}

	_, err := newBackfillUseCase(registryCon(sel), nil, nil, nil).
		Preview(context.Background(), dtos.BackfillFilter{EventCode: "guide_created"})

	require.NoError(t, err)
	assert.Empty(t, sel.Dispatched(),
		"el preview es solo lectura: no debe enviar ni un mensaje")
}

func TestRun_EventoNoSoportado_Rechaza(t *testing.T) {
	store := mocks.NewJobStore()

	got, err := newBackfillUseCase(registrySinEvento(), store, nil, nil).
		Run(context.Background(), dtos.RunRequest{EventCode: "inventado"}, 42)

	require.Error(t, err)
	assert.Nil(t, got)
	assert.Contains(t, err.Error(), "no soportado")
	assert.Empty(t, store.List(nil), "no se crea job para un evento que no existe")
}

func TestRun_ErrorAlPreviewear_NoCreaJob(t *testing.T) {
	selErr := stderrors.New("query fallida")
	store := mocks.NewJobStore()
	sel := &mocks.SelectorMock{
		PreviewFn: func(ctx context.Context, f dtos.BackfillFilter) ([]entities.Candidate, error) {
			return nil, selErr
		},
	}

	_, err := newBackfillUseCase(registryCon(sel), store, nil, nil).
		Run(context.Background(), dtos.RunRequest{EventCode: "guide_created"}, 42)

	assert.ErrorIs(t, err, selErr)
	assert.Empty(t, store.List(nil))
	assert.Empty(t, sel.Dispatched())
}

func TestRun_ArmaElFiltroDesdeElRequest(t *testing.T) {
	sel := &mocks.SelectorMock{}
	pub := mocks.NewProgressPublisher()

	_, err := newBackfillUseCase(registryCon(sel), nil, pub, nil).Run(context.Background(),
		dtos.RunRequest{EventCode: "guide_created", BusinessID: uintPtr(26), Days: 7, Limit: 50}, 42)

	require.NoError(t, err)
	esperarFin(t, pub)
	require.Len(t, sel.PreviewFilters(), 1)
	f := sel.PreviewFilters()[0]
	assert.Equal(t, "guide_created", f.EventCode)
	require.NotNil(t, f.BusinessID)
	assert.Equal(t, uint(26), *f.BusinessID)
	assert.Equal(t, 7, f.Days)
	assert.Equal(t, 50, f.Limit)
}

func TestRun_CreaElJobEnEjecucionConElTotal(t *testing.T) {
	sel := &mocks.SelectorMock{
		PreviewFn: func(ctx context.Context, f dtos.BackfillFilter) ([]entities.Candidate, error) {
			return []entities.Candidate{{OrderID: "O-1"}, {OrderID: "O-2"}}, nil
		},
	}
	pub := mocks.NewProgressPublisher()

	got, err := newBackfillUseCase(registryCon(sel), nil, pub, nil).Run(context.Background(),
		dtos.RunRequest{EventCode: "guide_created", BusinessID: uintPtr(26)}, 42)

	require.NoError(t, err)
	assert.NotEmpty(t, got.ID)
	assert.Equal(t, entities.JobStatusRunning, got.Status)
	assert.Equal(t, 2, got.TotalEligible)
	assert.Equal(t, "guide_created", got.EventCode)
	assert.Equal(t, uint(42), got.CreatedBy)
	require.NotNil(t, got.BusinessID)
	assert.Equal(t, uint(26), *got.BusinessID)
	assert.False(t, got.StartedAt.IsZero())
	esperarFin(t, pub)
}

func TestRun_JobIDsUnicos(t *testing.T) {
	vistos := map[string]bool{}
	for i := 0; i < 10; i++ {
		pub := mocks.NewProgressPublisher()
		got, err := newBackfillUseCase(registryCon(&mocks.SelectorMock{}), nil, pub, nil).
			Run(context.Background(), dtos.RunRequest{EventCode: "guide_created"}, 42)
		require.NoError(t, err)
		require.False(t, vistos[got.ID], "job_id repetido: %q", got.ID)
		vistos[got.ID] = true
		esperarFin(t, pub)
	}
	assert.Len(t, vistos, 10)
}

func TestRun_DespachaTodosLosCandidatos(t *testing.T) {
	sel := &mocks.SelectorMock{
		PreviewFn: func(ctx context.Context, f dtos.BackfillFilter) ([]entities.Candidate, error) {
			return []entities.Candidate{
				{OrderID: "O-1", BusinessID: 26},
				{OrderID: "O-2", BusinessID: 26},
			}, nil
		},
	}
	pub := mocks.NewProgressPublisher()

	_, err := newBackfillUseCase(registryCon(sel), nil, pub, nil).
		Run(context.Background(), dtos.RunRequest{EventCode: "guide_created"}, 42)

	require.NoError(t, err)
	esperarFin(t, pub)

	despachados := sel.Dispatched()
	require.Len(t, despachados, 2)
	assert.Equal(t, "O-1", despachados[0].OrderID)
	assert.Equal(t, "O-2", despachados[1].OrderID)
}

func TestRun_CuentaEnviadosYFallidos(t *testing.T) {
	store := mocks.NewJobStore()
	pub := mocks.NewProgressPublisher()
	sel := &mocks.SelectorMock{
		PreviewFn: func(ctx context.Context, f dtos.BackfillFilter) ([]entities.Candidate, error) {
			return []entities.Candidate{
				{OrderID: "O-1"}, {OrderID: "O-2"}, {OrderID: "O-3"},
			}, nil
		},
		DispatchFn: func(ctx context.Context, c entities.Candidate) error {
			if c.OrderID == "O-2" {
				return stderrors.New("whatsapp rechazo el envio")
			}
			return nil
		},
	}

	job, err := newBackfillUseCase(registryCon(sel), store, pub, nil).
		Run(context.Background(), dtos.RunRequest{EventCode: "guide_created"}, 42)

	require.NoError(t, err)
	esperarFin(t, pub)

	final, ok := store.Get(job.ID)
	require.True(t, ok)
	assert.Equal(t, 2, final.Sent)
	assert.Equal(t, 1, final.Failed)
	assert.Equal(t, 3, final.TotalEligible)
	assert.Len(t, sel.Dispatched(), 3,
		"un envio fallido no debe abortar los siguientes")
}

func TestRun_AlTerminarQuedaCompletadoConFecha(t *testing.T) {
	store := mocks.NewJobStore()
	pub := mocks.NewProgressPublisher()

	job, err := newBackfillUseCase(registryCon(&mocks.SelectorMock{}), store, pub, nil).
		Run(context.Background(), dtos.RunRequest{EventCode: "guide_created"}, 42)

	require.NoError(t, err)
	esperarFin(t, pub)

	final, ok := store.Get(job.ID)
	require.True(t, ok)
	assert.Equal(t, entities.JobStatusCompleted, final.Status)
	require.NotNil(t, final.FinishedAt)
	assert.False(t, final.FinishedAt.Before(final.StartedAt))
}

func TestRun_TodosFallan_IgualQuedaCompletado(t *testing.T) {
	store := mocks.NewJobStore()
	pub := mocks.NewProgressPublisher()
	sel := &mocks.SelectorMock{
		PreviewFn: func(ctx context.Context, f dtos.BackfillFilter) ([]entities.Candidate, error) {
			return []entities.Candidate{{OrderID: "O-1"}, {OrderID: "O-2"}}, nil
		},
		DispatchFn: func(ctx context.Context, c entities.Candidate) error {
			return stderrors.New("caido")
		},
	}

	job, err := newBackfillUseCase(registryCon(sel), store, pub, nil).
		Run(context.Background(), dtos.RunRequest{EventCode: "guide_created"}, 42)

	require.NoError(t, err)
	esperarFin(t, pub)

	final, _ := store.Get(job.ID)
	assert.Equal(t, entities.JobStatusCompleted, final.Status,
		"el job termina aunque no se haya enviado nada; el detalle esta en Failed")
	assert.Equal(t, 2, final.Failed)
	assert.Zero(t, final.Sent)
}

func TestRun_PublicaProgresoAlIniciarPorCadaEnvioYAlCerrar(t *testing.T) {
	pub := mocks.NewProgressPublisher()
	sel := &mocks.SelectorMock{
		PreviewFn: func(ctx context.Context, f dtos.BackfillFilter) ([]entities.Candidate, error) {
			return []entities.Candidate{{OrderID: "O-1"}, {OrderID: "O-2"}}, nil
		},
	}

	_, err := newBackfillUseCase(registryCon(sel), nil, pub, nil).
		Run(context.Background(), dtos.RunRequest{EventCode: "guide_created"}, 42)

	require.NoError(t, err)
	esperarFin(t, pub)

	snaps := pub.Snapshots()
	require.GreaterOrEqual(t, len(snaps), 4,
		"al menos: inicio, uno por cada candidato y el cierre")
	assert.Equal(t, entities.JobStatusRunning, snaps[0].Status)
	assert.Equal(t, entities.JobStatusCompleted, snaps[len(snaps)-1].Status)
}

func TestRun_ElProgresoAvanzaMonotonamente(t *testing.T) {
	pub := mocks.NewProgressPublisher()
	sel := &mocks.SelectorMock{
		PreviewFn: func(ctx context.Context, f dtos.BackfillFilter) ([]entities.Candidate, error) {
			return []entities.Candidate{{OrderID: "O-1"}, {OrderID: "O-2"}, {OrderID: "O-3"}}, nil
		},
	}

	_, err := newBackfillUseCase(registryCon(sel), nil, pub, nil).
		Run(context.Background(), dtos.RunRequest{EventCode: "guide_created"}, 42)

	require.NoError(t, err)
	esperarFin(t, pub)

	previo := -1
	for i, s := range pub.Snapshots() {
		procesados := s.Sent + s.Failed + s.Skipped
		assert.GreaterOrEqual(t, procesados, previo,
			"el contador retrocedio en el snapshot %d", i)
		assert.LessOrEqual(t, procesados, s.TotalEligible,
			"se procesaron mas de los elegibles en el snapshot %d", i)
		previo = procesados
	}
}

func TestRun_SinCandidatos_TerminaDeInmediato(t *testing.T) {
	store := mocks.NewJobStore()
	pub := mocks.NewProgressPublisher()
	sel := &mocks.SelectorMock{}

	job, err := newBackfillUseCase(registryCon(sel), store, pub, nil).
		Run(context.Background(), dtos.RunRequest{EventCode: "guide_created"}, 42)

	require.NoError(t, err)
	esperarFin(t, pub)

	final, _ := store.Get(job.ID)
	assert.Equal(t, entities.JobStatusCompleted, final.Status)
	assert.Zero(t, final.TotalEligible)
	assert.Empty(t, sel.Dispatched())
}

func TestRun_RespetaElThrottleEntreEnvios(t *testing.T) {
	pub := mocks.NewProgressPublisher()
	sel := &mocks.SelectorMock{
		PreviewFn: func(ctx context.Context, f dtos.BackfillFilter) ([]entities.Candidate, error) {
			return []entities.Candidate{{OrderID: "O-1"}, {OrderID: "O-2"}, {OrderID: "O-3"}}, nil
		},
	}

	inicio := time.Now()
	_, err := newBackfillUseCase(registryCon(sel), nil, pub, nil).
		Run(context.Background(), dtos.RunRequest{EventCode: "guide_created"}, 42)
	require.NoError(t, err)
	esperarFin(t, pub)
	transcurrido := time.Since(inicio)

	esperado := time.Duration(2) * (time.Second / defaultThrottlePerSec)
	assert.GreaterOrEqual(t, transcurrido, esperado,
		"3 envios deben esperar 2 intervalos: sin throttle se satura el proveedor de WhatsApp")
}

func TestRun_DevuelveInmediatamenteSinEsperarElDespacho(t *testing.T) {
	pub := mocks.NewProgressPublisher()
	sel := &mocks.SelectorMock{
		PreviewFn: func(ctx context.Context, f dtos.BackfillFilter) ([]entities.Candidate, error) {
			return []entities.Candidate{{OrderID: "O-1"}, {OrderID: "O-2"}, {OrderID: "O-3"}}, nil
		},
	}

	inicio := time.Now()
	got, err := newBackfillUseCase(registryCon(sel), nil, pub, nil).
		Run(context.Background(), dtos.RunRequest{EventCode: "guide_created"}, 42)
	respuesta := time.Since(inicio)

	require.NoError(t, err)
	assert.Equal(t, entities.JobStatusRunning, got.Status)
	assert.Less(t, respuesta, 100*time.Millisecond,
		"el endpoint responde ya con el job_id; el envio corre en background")
	esperarFin(t, pub)
}

func TestGetJob_DevuelveElJobGuardado(t *testing.T) {
	store := mocks.NewJobStore()
	pub := mocks.NewProgressPublisher()

	job, err := newBackfillUseCase(registryCon(&mocks.SelectorMock{}), store, pub, nil).
		Run(context.Background(), dtos.RunRequest{EventCode: "guide_created"}, 42)
	require.NoError(t, err)
	esperarFin(t, pub)

	uc := newBackfillUseCase(registryCon(&mocks.SelectorMock{}), store, pub, nil)
	got, ok := uc.GetJob(job.ID)

	require.True(t, ok)
	assert.Equal(t, job.ID, got.ID)
	assert.Equal(t, "guide_created", got.EventCode)
}

func TestGetJob_Inexistente_DevuelveFalse(t *testing.T) {
	got, ok := newBackfillUseCase(nil, mocks.NewJobStore(), nil, nil).GetJob("no-existe")

	assert.False(t, ok)
	assert.Nil(t, got)
}

func TestThrottle_NoBajaDelMinimo(t *testing.T) {
	intervalo := time.Second / defaultThrottlePerSec
	assert.GreaterOrEqual(t, intervalo, minThrottleInterval,
		"el intervalo configurado no debe quedar por debajo del piso de seguridad")
}
