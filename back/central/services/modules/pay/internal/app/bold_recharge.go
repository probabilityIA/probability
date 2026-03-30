package app

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/dtos"
)

// BoldGenerateSignature genera la firma de integridad SHA256 para Bold.co
// Hash = SHA256(Identificador + Monto + Divisa + LlaveSecreta)
func (uc *walletUseCase) BoldGenerateSignature(ctx context.Context, amount float64, currency string) (*dtos.BoldSignatureResponse, error) {
	orderID := "WLT" + uuid.New().String()
	// Bold requiere el monto sin decimales (como entero en string)
	amountInt := int64(amount)
	
	identityKey := uc.config.Get("BOLD_IDENTITY_KEY")
	secretKey := uc.config.Get("BOLD_SECRET_KEY")

	// Formato: {Identificador}{Monto}{Divisa}{LlaveSecreta}
	rawString := fmt.Sprintf("%s%d%s%s", orderID, amountInt, currency, secretKey)
	
	hash := sha256.Sum256([]byte(rawString))
	signature := hex.EncodeToString(hash[:])

	uc.log.Info(ctx).
		Str("order_id", orderID).
		Float64("amount", amount).
		Msg("Generated Bold.co integrity signature")

	return &dtos.BoldSignatureResponse{
		OrderID:            orderID,
		IntegritySignature: signature,
		Amount:             amount,
		Currency:           currency,
		IdentityKey:        identityKey,
	}, nil
}

// GetBoldStatus consulta el estado de una orden directamente con la API de Bold.co
func (uc *walletUseCase) GetBoldStatus(ctx context.Context, boldOrderID string) (*dtos.BoldStatusResponse, error) {
	apiKey := uc.config.Get("BOLD_SECRET_KEY") // Para el API de consulta se suele usar la Secret Key
	
	url := fmt.Sprintf("https://api.bold.co/v2/payment-orders/%s", boldOrderID)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Add("Authorization", "Bearer "+apiKey)
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		uc.log.Error(ctx).Err(err).Str("bold_order_id", boldOrderID).Msg("Failed to call Bold API")
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Bold API error: status %d", resp.StatusCode)
	}

	var boldData struct {
		ID            string  `json:"id"`
		PaymentStatus string  `json:"payment_status"`
		Amount        float64 `json:"amount"`
		Currency      string  `json:"currency"`
		PaymentMethod string  `json:"payment_method"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&boldData); err != nil {
		return nil, err
	}

	uc.log.Info(ctx).
		Str("bold_order_id", boldOrderID).
		Str("status", boldData.PaymentStatus).
		Msg("Fetched Bold.co transaction status")

	return &dtos.BoldStatusResponse{
		BoldOrderID:   boldData.ID,
		Status:        strings.ToUpper(boldData.PaymentStatus),
		Amount:        boldData.Amount,
		Currency:      boldData.Currency,
		PaymentMethod: boldData.PaymentMethod,
	}, nil
}
