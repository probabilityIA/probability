package consumer

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/constants"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/mocks"
)

func TestRetryBackoffCreceYSeTopea(t *testing.T) {
	cases := []struct {
		attempt int
		want    time.Duration
	}{
		{1, 15 * time.Minute},
		{2, 30 * time.Minute},
		{3, time.Hour},
		{4, 2 * time.Hour},
		{5, 4 * time.Hour},
		{6, retryMaxDelay},
		{50, retryMaxDelay},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("intento %d", tc.attempt), func(t *testing.T) {
			if got := retryBackoff(tc.attempt); got != tc.want {
				t.Fatalf("intento %d: esperado %s, obtenido %s", tc.attempt, tc.want, got)
			}
		})
	}
}

func TestRetryBackoffNoDecrece(t *testing.T) {
	var previo time.Duration
	for attempt := 1; attempt <= 100; attempt++ {
		actual := retryBackoff(attempt)
		if actual < previo {
			t.Fatalf("el backoff bajo de %s a %s en el intento %d", previo, actual, attempt)
		}
		if actual > retryMaxDelay {
			t.Fatalf("el backoff %s supera el tope %s", actual, retryMaxDelay)
		}
		previo = actual
	}
}

func TestErroresPermanentesSeReconocen(t *testing.T) {
	permanentes := []error{
		domainerrors.ErrRetryNotAllowed,
		domainerrors.ErrMaxRetriesExceeded,
		domainerrors.ErrInvoiceNotFound,
		domainerrors.ErrProviderNotConfigured,
		domainerrors.ErrSyncLogNotFound,
	}

	for _, err := range permanentes {
		if !isPermanentRetryError(err) {
			t.Fatalf("%v deberia ser permanente", err)
		}
		envuelto := fmt.Errorf("contexto adicional: %w", err)
		if !isPermanentRetryError(envuelto) {
			t.Fatalf("%v envuelto deberia seguir siendo permanente", err)
		}
	}
}

func TestErroresTransitoriosNoSeRetiran(t *testing.T) {
	transitorios := []error{
		domainerrors.ErrProviderTimeout,
		domainerrors.ErrProviderAPIError,
		domainerrors.ErrProviderRateLimitExceeded,
		domainerrors.ErrSyncFailed,
		fmt.Errorf("failed to get order: connection refused"),
	}

	for _, err := range transitorios {
		if isPermanentRetryError(err) {
			t.Fatalf("%v no deberia tratarse como permanente", err)
		}
	}
}

type fakeRetryRepo struct {
	ports.IRepository
	invoices map[uint]*entities.Invoice
	logs     map[uint]*entities.InvoiceSyncLog
	batch    []*entities.InvoiceSyncLog
	updates  int
}

func (r *fakeRetryRepo) GetInvoiceByID(_ context.Context, id uint) (*entities.Invoice, error) {
	invoice, ok := r.invoices[id]
	if !ok {
		return nil, fmt.Errorf("no existe")
	}
	return invoice, nil
}

func (r *fakeRetryRepo) GetPendingSyncLogRetries(_ context.Context, _ int) ([]*entities.InvoiceSyncLog, error) {
	return r.batch, nil
}

func (r *fakeRetryRepo) GetSyncLogsByInvoiceID(_ context.Context, invoiceID uint) ([]*entities.InvoiceSyncLog, error) {
	out := make([]*entities.InvoiceSyncLog, 0, 1)
	for _, l := range r.logs {
		if l.InvoiceID == invoiceID {
			out = append(out, l)
		}
	}
	return out, nil
}

func (r *fakeRetryRepo) UpdateInvoiceSyncLog(_ context.Context, l *entities.InvoiceSyncLog) error {
	r.updates++
	r.logs[l.ID] = l
	return nil
}

type fakeRetryUseCase struct {
	ports.IUseCase
	retried    []uint
	checked    []uint
	retryErr   error
	checkErr   error
	onDispatch func(invoiceID uint)
}

func (u *fakeRetryUseCase) RetryInvoice(_ context.Context, invoiceID uint, _ bool) error {
	u.retried = append(u.retried, invoiceID)
	if u.onDispatch != nil {
		u.onDispatch(invoiceID)
	}
	return u.retryErr
}

func (u *fakeRetryUseCase) CheckPendingInvoice(_ context.Context, invoiceID uint) error {
	u.checked = append(u.checked, invoiceID)
	if u.onDispatch != nil {
		u.onDispatch(invoiceID)
	}
	return u.checkErr
}

func nuevoSyncLog(id, invoiceID uint, status string, retryCount int) *entities.InvoiceSyncLog {
	next := time.Now().Add(-time.Hour)
	return &entities.InvoiceSyncLog{
		ID:          id,
		InvoiceID:   invoiceID,
		Status:      status,
		RetryCount:  retryCount,
		MaxRetries:  constants.MaxRetries,
		NextRetryAt: &next,
	}
}

func armar(logs []*entities.InvoiceSyncLog, invoices map[uint]*entities.Invoice) *fakeRetryRepo {
	repo := &fakeRetryRepo{
		invoices: invoices,
		logs:     make(map[uint]*entities.InvoiceSyncLog, len(logs)),
		batch:    logs,
	}
	for _, l := range logs {
		repo.logs[l.ID] = l
	}
	return repo
}

func nuevoConsumer(repo ports.IRepository, uc ports.IUseCase) *RetryConsumer {
	return NewRetryConsumer(repo, uc, mocks.NewSilentLogger())
}

func TestFacturaPendingUsaCheckNoRetry(t *testing.T) {
	syncLog := nuevoSyncLog(1, 100, constants.SyncStatusFailed, 0)
	repo := armar([]*entities.InvoiceSyncLog{syncLog}, map[uint]*entities.Invoice{
		100: {ID: 100, Status: constants.InvoiceStatusPending},
	})
	uc := &fakeRetryUseCase{}

	nuevoConsumer(repo, uc).processRetries(context.Background())

	if len(uc.retried) != 0 {
		t.Fatalf("no debia reenviar el POST de una factura pending, se llamo con %v", uc.retried)
	}
	if len(uc.checked) != 1 || uc.checked[0] != 100 {
		t.Fatalf("debia consultar el estado de la factura 100, se consulto %v", uc.checked)
	}
}

func TestErrorPermanenteRetiraElSyncLog(t *testing.T) {
	syncLog := nuevoSyncLog(1, 100, constants.SyncStatusFailed, 0)
	repo := armar([]*entities.InvoiceSyncLog{syncLog}, map[uint]*entities.Invoice{
		100: {ID: 100, Status: constants.InvoiceStatusFailed},
	})
	uc := &fakeRetryUseCase{retryErr: domainerrors.ErrRetryNotAllowed}

	nuevoConsumer(repo, uc).processRetries(context.Background())

	guardado := repo.logs[1]
	if guardado.Status != constants.SyncStatusCancelled {
		t.Fatalf("esperado status cancelled, obtenido %s", guardado.Status)
	}
	if guardado.NextRetryAt != nil {
		t.Fatalf("next_retry_at debia quedar en nil, quedo en %v", guardado.NextRetryAt)
	}
}

func TestErrorTransitorioAplazaConBackoff(t *testing.T) {
	syncLog := nuevoSyncLog(1, 100, constants.SyncStatusFailed, 0)
	repo := armar([]*entities.InvoiceSyncLog{syncLog}, map[uint]*entities.Invoice{
		100: {ID: 100, Status: constants.InvoiceStatusFailed},
	})
	uc := &fakeRetryUseCase{retryErr: domainerrors.ErrProviderTimeout}

	antes := time.Now()
	nuevoConsumer(repo, uc).processRetries(context.Background())

	guardado := repo.logs[1]
	if guardado.RetryCount != 1 {
		t.Fatalf("retry_count debia avanzar a 1, quedo en %d", guardado.RetryCount)
	}
	if guardado.NextRetryAt == nil || !guardado.NextRetryAt.After(antes.Add(10*time.Minute)) {
		t.Fatalf("next_retry_at debia moverse al futuro, quedo en %v", guardado.NextRetryAt)
	}
	if guardado.Status == constants.SyncStatusCancelled {
		t.Fatal("un error transitorio no debe cancelar el reintento")
	}
}

func TestNoTocaElSyncLogSiElCasoDeUsoYaLoManejo(t *testing.T) {
	syncLog := nuevoSyncLog(1, 100, constants.SyncStatusFailed, 0)
	repo := armar([]*entities.InvoiceSyncLog{syncLog}, map[uint]*entities.Invoice{
		100: {ID: 100, Status: constants.InvoiceStatusFailed},
	})
	uc := &fakeRetryUseCase{retryErr: domainerrors.ErrProviderTimeout}
	uc.onDispatch = func(uint) {
		syncLog.Status = constants.SyncStatusCancelled
		syncLog.NextRetryAt = nil
	}

	nuevoConsumer(repo, uc).processRetries(context.Background())

	if repo.updates != 0 {
		t.Fatalf("el consumidor no debia reescribir un sync log que el caso de uso ya cerro (%d escrituras)", repo.updates)
	}
}

func TestFacturaEmitidaRetiraElReintento(t *testing.T) {
	syncLog := nuevoSyncLog(1, 100, constants.SyncStatusFailed, 0)
	repo := armar([]*entities.InvoiceSyncLog{syncLog}, map[uint]*entities.Invoice{
		100: {ID: 100, Status: constants.InvoiceStatusIssued},
	})
	uc := &fakeRetryUseCase{}

	nuevoConsumer(repo, uc).processRetries(context.Background())

	if len(uc.retried) != 0 || len(uc.checked) != 0 {
		t.Fatal("una factura ya emitida no debe disparar ningun reintento")
	}
	if repo.logs[1].Status != constants.SyncStatusCancelled {
		t.Fatalf("el sync log debia retirarse, quedo en %s", repo.logs[1].Status)
	}
}

func TestElBucleNoSeRepiteEnLaSiguienteVuelta(t *testing.T) {
	syncLog := nuevoSyncLog(1, 100, constants.SyncStatusFailed, 0)
	repo := armar([]*entities.InvoiceSyncLog{syncLog}, map[uint]*entities.Invoice{
		100: {ID: 100, Status: constants.InvoiceStatusPending},
	})
	uc := &fakeRetryUseCase{checkErr: domainerrors.ErrRetryNotAllowed}
	consumer := nuevoConsumer(repo, uc)

	consumer.processRetries(context.Background())

	elegibles := make([]*entities.InvoiceSyncLog, 0, 1)
	for _, l := range repo.logs {
		if l.NextRetryAt != nil && l.RetryCount < l.MaxRetries &&
			(l.Status == constants.SyncStatusFailed || l.Status == constants.SyncStatusPending) {
			elegibles = append(elegibles, l)
		}
	}
	repo.batch = elegibles

	if len(repo.batch) != 0 {
		t.Fatalf("la fila envenenada seguiria elegible: %d filas", len(repo.batch))
	}

	consumer.processRetries(context.Background())

	if len(uc.checked) != 1 {
		t.Fatalf("la segunda vuelta no debia reprocesar nada, total de consultas: %d", len(uc.checked))
	}
}

func TestDespachoExitosoDejaDeDisparar(t *testing.T) {
	syncLog := nuevoSyncLog(1, 100, constants.SyncStatusFailed, 0)
	repo := armar([]*entities.InvoiceSyncLog{syncLog}, map[uint]*entities.Invoice{
		100: {ID: 100, Status: constants.InvoiceStatusPending},
	})
	uc := &fakeRetryUseCase{}
	consumidor := nuevoConsumer(repo, uc)

	consumidor.processRetries(context.Background())

	if len(uc.checked) != 1 {
		t.Fatalf("la primera vuelta debia despachar una consulta, despacho %d", len(uc.checked))
	}
	if repo.logs[1].NextRetryAt != nil {
		t.Fatal("la fila despachada debia dejar de estar elegible: next_retry_at sigue puesto")
	}

	elegibles := make([]*entities.InvoiceSyncLog, 0, 1)
	for _, l := range repo.logs {
		if l.NextRetryAt != nil && l.RetryCount < l.MaxRetries &&
			(l.Status == constants.SyncStatusFailed || l.Status == constants.SyncStatusPending) {
			elegibles = append(elegibles, l)
		}
	}
	repo.batch = elegibles

	consumidor.processRetries(context.Background())

	if len(uc.checked) != 1 {
		t.Fatalf("la segunda vuelta volvio a despachar: %d consultas en total", len(uc.checked))
	}
}
