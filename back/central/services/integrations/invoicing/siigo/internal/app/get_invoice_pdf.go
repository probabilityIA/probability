package app

import (
	"context"
	"fmt"
	"strings"
)

func (uc *invoicingUseCase) GetInvoicePDF(ctx context.Context, integrationID uint, siigoInvoiceID string) ([]byte, string, error) {
	if integrationID == 0 {
		return nil, "", fmt.Errorf("integration_id es requerido")
	}
	if strings.TrimSpace(siigoInvoiceID) == "" {
		return nil, "", fmt.Errorf("la factura no tiene identificador de Siigo: no se puede descargar el PDF")
	}

	creds, err := uc.resolveWebhookCredentials(ctx, fmt.Sprintf("%d", integrationID))
	if err != nil {
		return nil, "", err
	}

	return uc.siigoClient.GetInvoicePDF(ctx, creds, siigoInvoiceID)
}
