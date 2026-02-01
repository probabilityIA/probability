package app

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/dtos"
)

// CancelInvoice cancela una factura emitida
// DEPRECATED: Esta funcionalidad fue migrada a integrations/invoicing/softpymes/
// TODO: Re-implementar usando integrationCore si es necesario
func (uc *useCase) CancelInvoice(ctx context.Context, dto *dtos.CancelInvoiceDTO) error {
	return fmt.Errorf("CancelInvoice is deprecated and was moved to softpymes integration")
}

/* IMPLEMENTACIÓN ORIGINAL (requiere providerRepo y providerClient que ya no existen):

func (uc *useCase) CancelInvoice(ctx context.Context, dto *dtos.CancelInvoiceDTO) error {
	uc.log.Info(ctx).Uint("invoice_id", dto.InvoiceID).Msg("Cancelling invoice")

	// 1. Obtener factura
	invoice, err := uc.invoiceRepo.GetByID(ctx, dto.InvoiceID)
	if err != nil {
		return errors.ErrInvoiceNotFound
	}

	// 2. Validar que la factura esté emitida
	if invoice.Status != constants.InvoiceStatusIssued {
		return errors.ErrInvoiceCannotBeCancelled
	}

	// 3. Obtener proveedor
	provider, err := uc.providerRepo.GetByID(ctx, invoice.InvoicingProviderID)
	if err != nil {
		return errors.ErrProviderNotFound
	}

	// 4. Desencriptar credenciales
	credentials, err := uc.encryption.Decrypt(provider.Credentials)
	if err != nil {
		return errors.ErrDecryptionFailed
	}

	// 5. Autenticar
	token, err := uc.providerClient.Authenticate(ctx, credentials)
	if err != nil {
		return errors.ErrAuthenticationFailed
	}

	// 6. Cancelar en el proveedor
	if invoice.ExternalID == nil {
		return fmt.Errorf("invoice has no external ID")
	}

	if err := uc.providerClient.CancelInvoice(ctx, token, *invoice.ExternalID, dto.Reason); err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to cancel invoice with provider")
		return errors.ErrProviderAPIError
	}

	// 7. Actualizar factura
	now := time.Now()
	invoice.Status = constants.InvoiceStatusCancelled
	invoice.CancelledAt = &now

	if err := uc.invoiceRepo.Update(ctx, invoice); err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to update invoice")
		return err
	}

	// 8. Publicar evento
	if err := uc.eventPublisher.PublishInvoiceCancelled(ctx, invoice); err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to publish event")
	}

	uc.log.Info(ctx).Uint("invoice_id", dto.InvoiceID).Msg("Invoice cancelled successfully")
	return nil
}
*/
