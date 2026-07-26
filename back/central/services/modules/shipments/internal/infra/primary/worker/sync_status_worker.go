package worker

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/app/usecaseshipment"
	domain "github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

var syncHours = []int{7, 13, 19}

const (
	syncTimeZone      = "America/Bogota"
	syncLookbackDays  = 30
	syncBusinessPause = 2 * time.Second
)

type SyncStatusWorker struct {
	repo domain.IRepository
	pub  domain.ITransportRequestPublisher
	log  log.ILogger
}

func NewSyncStatusWorker(repo domain.IRepository, pub domain.ITransportRequestPublisher, logger log.ILogger) *SyncStatusWorker {
	return &SyncStatusWorker{
		repo: repo,
		pub:  pub,
		log:  logger.WithModule("shipments.sync_status_worker"),
	}
}

func (w *SyncStatusWorker) Start(ctx context.Context) {
	location, err := time.LoadLocation(syncTimeZone)
	if err != nil {
		location = time.UTC
	}

	for {
		wait := time.Until(nextRun(time.Now().In(location), location))
		timer := time.NewTimer(wait)

		select {
		case <-ctx.Done():
			timer.Stop()
			return
		case <-timer.C:
			w.runAll(ctx)
		}
	}
}

func nextRun(now time.Time, location *time.Location) time.Time {
	for _, hour := range syncHours {
		candidate := time.Date(now.Year(), now.Month(), now.Day(), hour, 0, 0, 0, location)
		if candidate.After(now) {
			return candidate
		}
	}
	tomorrow := now.AddDate(0, 0, 1)
	return time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), syncHours[0], 0, 0, 0, location)
}

func (w *SyncStatusWorker) runAll(ctx context.Context) {
	businessIDs, err := w.repo.ListBusinessesWithActiveProvider(ctx, domain.SyncProviderEnvioclick)
	if err != nil {
		w.log.Error(ctx).Err(err).Msg("Error listando negocios con transportadora activa")
		return
	}

	if len(businessIDs) == 0 {
		w.log.Info(ctx).Msg("Sin negocios con envios activos para sincronizar")
		return
	}

	w.log.Info(ctx).
		Int("businesses", len(businessIDs)).
		Msg("Iniciando sincronizacion programada de estados de envio")

	from := time.Now().AddDate(0, 0, -syncLookbackDays)
	to := time.Now()
	syncUC := usecaseshipment.NewSyncShipments(w.repo, w.pub)

	var totalShipments int
	for _, businessID := range businessIDs {
		select {
		case <-ctx.Done():
			return
		default:
		}

		result, err := syncUC.SyncShipments(ctx, domain.SyncShipmentsFilter{
			BusinessID: businessID,
			Provider:   domain.SyncProviderEnvioclick,
			DateFrom:   &from,
			DateTo:     &to,
		})
		if err != nil {
			w.log.Warn(ctx).
				Err(err).
				Uint("business_id", businessID).
				Msg("No se pudo sincronizar los envios del negocio")
			continue
		}

		totalShipments += result.TotalShipments
		w.log.Info(ctx).
			Uint("business_id", businessID).
			Int("shipments", result.TotalShipments).
			Str("correlation_id", result.CorrelationID).
			Msg("Sincronizacion encolada para el negocio")

		time.Sleep(syncBusinessPause)
	}

	w.log.Info(ctx).
		Int("businesses", len(businessIDs)).
		Int("shipments", totalShipments).
		Msg("Sincronizacion programada de estados de envio completada")
}
