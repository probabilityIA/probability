package app

import (
	"context"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/constants"
	paydtos "github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/entities"
	payerrs "github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/errors"
)

// RetryPayment reintenta una transacción fallida
func (uc *useCase) RetryPayment(ctx context.Context, transactionID uint) error {
	uc.log.Info(ctx).Uint("transaction_id", transactionID).Msg("Retrying payment")

	// Obtener transacción
	tx, err := uc.repo.GetPaymentTransactionByID(ctx, transactionID)
	if err != nil {
		return fmt.Errorf("failed to get payment transaction: %w", err)
	}

	// Verificar que no esté ya completada
	if tx.Status == entities.PaymentStatusCompleted {
		return payerrs.ErrPaymentAlreadyProcessed
	}

	// Obtener sync logs para contar reintentos
	syncLogs, err := uc.repo.GetSyncLogsByTransactionID(ctx, tx.ID)
	if err != nil {
		return fmt.Errorf("failed to get sync logs: %w", err)
	}

	// Contar intentos totales
	totalAttempts := len(syncLogs)
	if totalAttempts >= constants.MaxRetries {
		tx.Status = entities.PaymentStatusFailed
		_ = uc.repo.UpdatePaymentTransaction(ctx, tx)
		return payerrs.ErrMaxRetriesReached
	}

	// Cancelar sync logs pendientes
	if err := uc.repo.CancelPendingSyncLogs(ctx, tx.ID); err != nil {
		uc.log.Warn(ctx).Err(err).Msg("Failed to cancel pending sync logs")
	}

	// Crear nuevo sync log
	newLog := &entities.PaymentSyncLog{
		PaymentTransactionID: tx.ID,
		Status:               constants.StatusProcessing,
		RetryCount:           totalAttempts,
	}
	if err := uc.repo.CreateSyncLog(ctx, newLog); err != nil {
		return fmt.Errorf("failed to create retry sync log: %w", err)
	}

	// Publicar solicitud de reintento
	correlationID := generateReference()
	msg := &paydtos.PaymentRequestMessage{
		PaymentTransactionID: tx.ID,
		BusinessID:           tx.BusinessID,
		GatewayCode:          tx.GatewayCode,
		Amount:               tx.Amount,
		Currency:             tx.Currency,
		Reference:            tx.Reference,
		PaymentMethod:        tx.PaymentMethod,
		Description:          tx.Description,
		Metadata:             tx.Metadata,
		CorrelationID:        correlationID,
		Timestamp:            time.Now(),
	}

	if err := uc.requestPublisher.PublishPaymentRequest(ctx, msg); err != nil {
		return fmt.Errorf("failed to publish retry request: %w", err)
	}

	uc.log.Info(ctx).
		Uint("transaction_id", tx.ID).
		Int("attempt", totalAttempts+1).
		Msg("Payment retry published")

	return nil
}
