package mappers

import (
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
		CreatedAt:            tx.CreatedAt,
	}
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
