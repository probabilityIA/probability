package app

import (
	"context"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/constants"
	paydtos "github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/entities"
)

// ProcessPaymentResponse procesa la respuesta del gateway de pago
func (uc *useCase) ProcessPaymentResponse(ctx context.Context, msg *paydtos.PaymentResponseMessage) error {
	uc.log.Info(ctx).
		Uint("transaction_id", msg.PaymentTransactionID).
		Str("gateway", msg.GatewayCode).
		Str("status", msg.Status).
		Str("correlation_id", msg.CorrelationID).
		Msg("Processing payment response")

	// Obtener transacción
	tx, err := uc.repo.GetPaymentTransactionByID(ctx, msg.PaymentTransactionID)
	if err != nil {
		uc.log.Error(ctx).Err(err).Uint("id", msg.PaymentTransactionID).Msg("Transaction not found")
		return fmt.Errorf("failed to get payment transaction: %w", err)
	}

	// Obtener sync logs activos
	syncLogs, err := uc.repo.GetSyncLogsByTransactionID(ctx, tx.ID)
	if err != nil {
		uc.log.Warn(ctx).Err(err).Msg("Failed to get sync logs")
	}

	// Encontrar el sync log más reciente (processing)
	var currentSyncLog *entities.PaymentSyncLog
	for _, sl := range syncLogs {
		if sl.Status == constants.StatusProcessing {
			currentSyncLog = sl
			break
		}
	}

	if msg.Status == "success" {
		// Éxito: actualizar transacción y sync log
		tx.Status = entities.PaymentStatusCompleted
		tx.ExternalID = msg.ExternalID

		if currentSyncLog != nil {
			currentSyncLog.Status = constants.StatusCompleted
			_ = uc.repo.UpdateSyncLog(ctx, currentSyncLog)
		}

		if err := uc.repo.UpdatePaymentTransaction(ctx, tx); err != nil {
			uc.log.Error(ctx).Err(err).Msg("Failed to update transaction to completed")
			return err
		}

		// Publicar evento SSE
		if uc.ssePublisher != nil {
			_ = uc.ssePublisher.PublishPaymentCompleted(ctx, tx)
		}

		uc.log.Info(ctx).
			Uint("transaction_id", tx.ID).
			Str("external_id", ptrStr(msg.ExternalID)).
			Msg("Payment completed successfully")

	} else {
		// Error: actualizar sync log y decidir si reintentar
		retryCount := 0
		if currentSyncLog != nil {
			retryCount = currentSyncLog.RetryCount
			currentSyncLog.Status = constants.StatusFailed
			errMsg := msg.Error
			currentSyncLog.ErrorMessage = &errMsg

			// Calcular próximo reintento
			if retryCount < constants.MaxRetries {
				nextRetry := time.Now().Add(5 * time.Minute)
				currentSyncLog.NextRetryAt = &nextRetry
			}
			_ = uc.repo.UpdateSyncLog(ctx, currentSyncLog)
		}

		// Actualizar transacción
		if retryCount >= constants.MaxRetries {
			tx.Status = entities.PaymentStatusFailed
		} else {
			// Dejar en pending para que el retry consumer lo reintente
			tx.Status = entities.PaymentStatusPending
		}

		if err := uc.repo.UpdatePaymentTransaction(ctx, tx); err != nil {
			uc.log.Error(ctx).Err(err).Msg("Failed to update transaction after error")
			return err
		}

		// Publicar evento SSE de fallo
		if tx.Status == entities.PaymentStatusFailed && uc.ssePublisher != nil {
			_ = uc.ssePublisher.PublishPaymentFailed(ctx, tx, msg.Error)
		}

		uc.log.Warn(ctx).
			Uint("transaction_id", tx.ID).
			Str("error", msg.Error).
			Int("retry_count", retryCount).
			Msg("Payment failed")
	}

	return nil
}

func ptrStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
