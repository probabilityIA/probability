package consumer

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/constants"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/mocks"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

const localDSN = "host=127.0.0.1 port=5434 user=postgres password=postgres dbname=probability sslmode=disable"

type localDB struct {
	conn *gorm.DB
}

func (d *localDB) Connect(context.Context) error        { return nil }
func (d *localDB) Close() error                         { return nil }
func (d *localDB) Conn(context.Context) *gorm.DB        { return d.conn }
func (d *localDB) WithContext(context.Context) *gorm.DB { return d.conn }
func (d *localDB) DebugConn(context.Context) *gorm.DB   { return d.conn }

type rechazaTodoUseCase struct {
	ports.IUseCase
	retried []uint
	checked []uint
}

func (u *rechazaTodoUseCase) RetryInvoice(_ context.Context, invoiceID uint, _ bool) error {
	u.retried = append(u.retried, invoiceID)
	return domainerrors.ErrRetryNotAllowed
}

func (u *rechazaTodoUseCase) CheckPendingInvoice(_ context.Context, invoiceID uint) error {
	u.checked = append(u.checked, invoiceID)
	return domainerrors.ErrRetryNotAllowed
}

type filaElegible struct {
	ID          uint
	InvoiceID   uint
	Status      string
	RetryCount  int
	MaxRetries  int
	NextRetryAt *time.Time
}

func elegibles(t *testing.T, conn *gorm.DB) []filaElegible {
	t.Helper()
	var filas []filaElegible
	err := conn.Table("invoice_sync_logs").
		Select("id, invoice_id, status, retry_count, max_retries, next_retry_at").
		Where("status IN (?, ?) AND next_retry_at IS NOT NULL AND next_retry_at <= ? AND retry_count < max_retries",
			"failed", "pending", time.Now()).
		Order("next_retry_at ASC").
		Scan(&filas).Error
	if err != nil {
		t.Fatalf("no se pudo leer el conjunto elegible: %v", err)
	}
	return filas
}

func TestBucleDeReintentosContraBaseLocal(t *testing.T) {
	if os.Getenv("INVOICING_LOCAL_DB_TEST") != "1" {
		t.Skip("requiere la base local en 5434: INVOICING_LOCAL_DB_TEST=1")
	}

	conn, err := gorm.Open(postgres.Open(localDSN), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		t.Fatalf("no se pudo conectar a la base local: %v", err)
	}

	antes := elegibles(t, conn)
	if len(antes) == 0 {
		t.Skip("la base local no tiene reintentos pendientes que probar")
	}

	snapshot := make([]filaElegible, len(antes))
	copy(snapshot, antes)
	t.Cleanup(func() {
		for _, fila := range snapshot {
			conn.Exec(
				"UPDATE invoice_sync_logs SET status = ?, retry_count = ?, next_retry_at = ?, error_message = NULL WHERE id = ?",
				fila.Status, fila.RetryCount, fila.NextRetryAt, fila.ID,
			)
		}
	})

	t.Logf("conjunto elegible inicial: %d filas", len(antes))
	for _, fila := range antes {
		var invoiceStatus string
		conn.Raw("SELECT status FROM invoices WHERE id = ?", fila.InvoiceID).Scan(&invoiceStatus)
		t.Logf("  sync_log=%d invoice=%d sync_status=%s invoice_status=%s retry_count=%d",
			fila.ID, fila.InvoiceID, fila.Status, invoiceStatus, fila.RetryCount)
	}

	repo := repository.New(&localDB{conn: conn}, nil, mocks.NewSilentLogger())
	useCase := &rechazaTodoUseCase{}
	retryConsumer := NewRetryConsumer(repo, useCase, mocks.NewSilentLogger())

	ctx := context.Background()
	retryConsumer.processRetries(ctx)

	llamadasPrimeraVuelta := len(useCase.retried) + len(useCase.checked)
	t.Logf("primera vuelta: %d retry, %d check", len(useCase.retried), len(useCase.checked))

	despues := elegibles(t, conn)
	if len(despues) != 0 {
		t.Fatalf("quedaron %d filas elegibles despues de la primera vuelta: el bucle seguiria vivo", len(despues))
	}
	t.Log("conjunto elegible despues de la primera vuelta: 0 filas")

	retryConsumer.processRetries(ctx)

	llamadasSegundaVuelta := len(useCase.retried) + len(useCase.checked) - llamadasPrimeraVuelta
	if llamadasSegundaVuelta != 0 {
		t.Fatalf("la segunda vuelta reproceso %d facturas: el bucle no se corto", llamadasSegundaVuelta)
	}
	t.Log("segunda vuelta: 0 llamadas, el bucle quedo cortado")

	var cancelados int64
	ids := make([]uint, 0, len(snapshot))
	for _, fila := range snapshot {
		ids = append(ids, fila.ID)
	}
	conn.Raw("SELECT count(*) FROM invoice_sync_logs WHERE id IN ? AND status = ? AND next_retry_at IS NULL",
		ids, constants.SyncStatusCancelled).Scan(&cancelados)
	if cancelados != int64(len(snapshot)) {
		t.Fatalf("esperaba %d sync logs retirados, hay %d", len(snapshot), cancelados)
	}
	t.Logf("%d sync logs retirados (cancelled + next_retry_at NULL)", cancelados)
}
