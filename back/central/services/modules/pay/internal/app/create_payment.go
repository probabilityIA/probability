package app

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/constants"
	paydtos "github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/entities"
	payerrs "github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/errors"
)

// CreatePayment inicia una nueva transacción de pago
func (uc *useCase) CreatePayment(ctx context.Context, dto *paydtos.CreatePaymentDTO) (*entities.PaymentTransaction, error) {
	uc.log.Info(ctx).
		Uint("business_id", dto.BusinessID).
		Float64("amount", dto.Amount).
		Str("gateway", dto.GatewayCode).
		Msg("Creating payment transaction")

	// Validar monto
	if dto.Amount <= 0 {
		return nil, payerrs.ErrInvalidAmount
	}

	// Validar gateway
	if dto.GatewayCode != constants.GatewayNequi {
		return nil, fmt.Errorf("%w: %s", payerrs.ErrInvalidGateway, dto.GatewayCode)
	}

	// Generar referencia única
	reference := generateReference()

	// Valores default
	currency := dto.Currency
	if currency == "" {
		currency = "COP"
	}
	paymentMethod := dto.PaymentMethod
	if paymentMethod == "" {
		paymentMethod = constants.PaymentMethodQRCode
	}

	// Crear entidad de dominio
	tx := &entities.PaymentTransaction{
		BusinessID:    dto.BusinessID,
		Amount:        dto.Amount,
		Currency:      currency,
		Status:        entities.PaymentStatusPending,
		GatewayCode:   dto.GatewayCode,
		Reference:     reference,
		PaymentMethod: paymentMethod,
		Description:   dto.Description,
		CallbackURL:   dto.CallbackURL,
		Metadata:      dto.Metadata,
	}

	// Persistir transacción
	if err := uc.repo.CreatePaymentTransaction(ctx, tx); err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to create payment transaction")
		return nil, fmt.Errorf("failed to create payment transaction: %w", err)
	}

	// Crear sync log inicial
	syncLog := &entities.PaymentSyncLog{
		PaymentTransactionID: tx.ID,
		Status:               constants.StatusProcessing,
		RetryCount:           0,
	}
	if err := uc.repo.CreateSyncLog(ctx, syncLog); err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to create initial sync log")
		// No fallar - la transacción ya fue creada
	}

	// Publicar solicitud a la cola
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
		uc.log.Error(ctx).Err(err).Msg("Failed to publish payment request")
		// Actualizar status a failed
		tx.Status = entities.PaymentStatusFailed
		_ = uc.repo.UpdatePaymentTransaction(ctx, tx)
		return nil, fmt.Errorf("failed to publish payment request: %w", err)
	}

	uc.log.Info(ctx).
		Uint("transaction_id", tx.ID).
		Str("reference", tx.Reference).
		Msg("Payment transaction created and request published")

	return tx, nil
}

// generateReference genera una referencia única para la transacción
func generateReference() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
