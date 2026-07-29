package pay

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/dtos"
)

func (b *Bundle) GetBalance(ctx context.Context, businessID uint) (float64, error) {
	wallet, err := b.WalletUseCase.GetWallet(ctx, businessID)
	if err != nil {
		return 0, err
	}
	return wallet.Balance, nil
}

func (b *Bundle) Debit(ctx context.Context, businessID uint, amount float64, reference, concept string, userID uint) error {
	return b.WalletUseCase.ManualDebit(ctx, &dtos.ManualDebitDTO{
		BusinessID: businessID,
		Amount:     amount,
		Reference:  reference,
		Concept:    concept,
		UserID:     &userID,
	})
}

// BoldSignature es la respuesta publica (cruza limites de modulo) para iniciar un pago
// con Bold que no es recarga de billetera.
type BoldSignature struct {
	OrderID        string
	Hash           string
	Amount         float64
	Currency       string
	PublicKey      string
	RedirectionURL string
	IsSandbox      bool
	PollingEnabled bool
}

// BoldGenerateSignatureForReference permite a otros modulos (ej. publicsite, para el
// checkout de la tienda publica) iniciar un cobro con Bold usando su propio prefijo de
// referencia, sin crear un WalletTransaction.
func (b *Bundle) BoldGenerateSignatureForReference(ctx context.Context, businessID uint, amount float64, currency, referencePrefix string) (*BoldSignature, error) {
	resp, err := b.WalletUseCase.BoldGenerateSignatureForReference(ctx, businessID, amount, currency, referencePrefix)
	if err != nil {
		return nil, err
	}
	return &BoldSignature{
		OrderID:        resp.OrderID,
		Hash:           resp.Hash,
		Amount:         resp.Amount,
		Currency:       resp.Currency,
		PublicKey:      resp.PublicKey,
		RedirectionURL: resp.RedirectionURL,
		IsSandbox:      resp.IsSandbox,
		PollingEnabled: resp.PollingEnabled,
	}, nil
}
