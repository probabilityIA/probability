package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/entities"
)

// GetSpendSummary total gastado por el negocio, agrupado por concepto. Nunca
// incluye margen/ganancia de Probability: solo movimientos de la propia
// billetera del negocio (type=USAGE).
func (uc *walletUseCase) GetSpendSummary(ctx context.Context, dto *dtos.SpendSummaryDTO) ([]dtos.ConceptTotal, error) {
	return uc.repo.GetSpendSummary(ctx, dto)
}

// ListSpendTransactions detalle paginado de los debitos de billetera de un negocio.
func (uc *walletUseCase) ListSpendTransactions(ctx context.Context, dto *dtos.TransactionFilterDTO) ([]*entities.WalletTransaction, int64, error) {
	return uc.repo.ListTransactionsFiltered(ctx, dto)
}
