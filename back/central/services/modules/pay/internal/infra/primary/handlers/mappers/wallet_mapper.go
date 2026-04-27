package mappers

import (
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/infra/primary/handlers/response"
)

// WalletToResponse convierte una entidad Wallet a respuesta HTTP
func WalletToResponse(w *entities.Wallet) *response.WalletResponse {
	return &response.WalletResponse{
		ID:         w.ID.String(),
		BusinessID: w.BusinessID,
		Balance:    w.Balance,
		CreatedAt:  w.CreatedAt,
		UpdatedAt:  w.UpdatedAt,
	}
}

// WalletTxToResponse convierte una entidad WalletTransaction a respuesta HTTP
func WalletTxToResponse(tx *entities.WalletTransaction) *response.WalletTransactionResponse {
	return &response.WalletTransactionResponse{
		ID:                   tx.ID.String(),
		WalletID:             tx.WalletID.String(),
		Amount:               tx.Amount,
		Type:                 tx.Type,
		Status:               tx.Status,
		Reference:            tx.Reference,
		QrCode:               tx.QrCode,
		PaymentTransactionID: tx.PaymentTransactionID,
		IntegrationTypeID:    tx.IntegrationTypeID,
		IntegrationID:        tx.IntegrationID,
		IntegrationName:      integrationNameFromReference(tx.Reference, tx.IntegrationTypeID),
		IntegrationImageURL:  tx.IntegrationImageURL,
		GatewayRequest:       jsonOrNil(tx.GatewayRequest),
		GatewayResponse:      jsonOrNil(tx.GatewayResponse),
		CreatedAt:            tx.CreatedAt,
	}
}

func jsonOrNil(b []byte) any {
	if len(b) == 0 {
		return nil
	}
	var v any
	if err := json.Unmarshal(b, &v); err != nil {
		return string(b)
	}
	return v
}

func integrationNameFromReference(reference string, integrationTypeID *uint) string {
	if integrationTypeID != nil {
		if *integrationTypeID == 23 {
			return "Bold"
		}
	}
	switch {
	case startsWith(reference, "BOLD_SANDBOX_") || startsWith(reference, "WLT"):
		return "Bold"
	case startsWith(reference, "MAN_DEB_"):
		return "Debito manual"
	case startsWith(reference, "MANUAL_"):
		return "Nequi"
	default:
		return ""
	}
}

func startsWith(s, prefix string) bool {
	if len(s) < len(prefix) {
		return false
	}
	return s[:len(prefix)] == prefix
}

// WalletListToResponse convierte una lista de Wallet a respuesta HTTP
func WalletListToResponse(wallets []*entities.Wallet) []*response.WalletResponse {
	result := make([]*response.WalletResponse, len(wallets))
	for i, w := range wallets {
		result[i] = WalletToResponse(w)
	}
	return result
}

// WalletTxListToResponse convierte una lista de WalletTransaction a respuesta HTTP
func WalletTxListToResponse(txs []*entities.WalletTransaction) []*response.WalletTransactionResponse {
	result := make([]*response.WalletTransactionResponse, len(txs))
	for i, tx := range txs {
		result[i] = WalletTxToResponse(tx)
	}
	return result
}
