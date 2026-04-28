package app

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	stderrors "errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/errors"
)

const boldStatusEndpoint = "/online/link/v1/%s"

func (uc *walletUseCase) BoldGenerateSignature(ctx context.Context, businessID uint, amount float64, currency string) (*dtos.BoldSignatureResponse, error) {
	creds, err := uc.repo.GetBoldCredentialsForBusiness(ctx, businessID)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to load Bold credentials")
		return nil, err
	}

	orderID := "WLT" + strings.ReplaceAll(uuid.New().String(), "-", "")[:25]
	amountInt := int64(amount)

	rawString := fmt.Sprintf("%s%d%s%s", orderID, amountInt, currency, creds.SecretKey)
	hash := sha256.Sum256([]byte(rawString))
	signature := hex.EncodeToString(hash[:])

	isSandbox := creds.Environment == "sandbox"

	wallet, err := uc.GetWallet(ctx, businessID)
	if err != nil {
		return nil, err
	}

	bizIntegration, _ := uc.repo.GetBoldIntegrationForBusiness(ctx, businessID)
	var integrationTypeID, integrationID *uint
	if bizIntegration != nil {
		if bizIntegration.IntegrationTypeID != 0 {
			id := bizIntegration.IntegrationTypeID
			integrationTypeID = &id
		}
		if bizIntegration.IntegrationID != 0 {
			id := bizIntegration.IntegrationID
			integrationID = &id
		}
	}
	if integrationTypeID == nil && creds.IntegrationTypeID != 0 {
		id := creds.IntegrationTypeID
		integrationTypeID = &id
	}

	gatewayRequest, _ := json.Marshal(map[string]any{
		"order_id":     orderID,
		"amount":       amount,
		"currency":     currency,
		"public_key":   creds.APIKey,
		"hash":         signature,
		"environment":  creds.Environment,
		"is_sandbox":   isSandbox,
		"generated_at": time.Now().Format(time.RFC3339),
	})

	pendingTx := &entities.WalletTransaction{
		WalletID:          wallet.ID,
		Amount:            amount,
		Type:              entities.WalletTxTypeRecharge,
		Status:            entities.WalletTxStatusPending,
		Reference:         orderID,
		IntegrationTypeID: integrationTypeID,
		IntegrationID:     integrationID,
		GatewayRequest:    gatewayRequest,
	}
	if err := uc.repo.CreateWalletTransaction(ctx, pendingTx); err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to create pending Bold wallet transaction")
		return nil, err
	}

	redirectionURL := uc.config.Get("WEBHOOK_BASE_URL")
	if redirectionURL == "" {
		redirectionURL = uc.config.Get("URL_BASE_SWAGGER")
	}
	if redirectionURL != "" {
		redirectionURL = strings.TrimRight(redirectionURL, "/") + "/wallet?bold_order_id=" + orderID
	}

	uc.log.Info(ctx).
		Str("order_id", orderID).
		Str("wallet_tx_id", pendingTx.ID.String()).
		Float64("amount", amount).
		Int64("amount_int", amountInt).
		Str("currency", currency).
		Str("raw_string_preview", fmt.Sprintf("%s%d%s***", orderID, amountInt, currency)).
		Str("hash_preview", signature[:16]+"...").
		Str("api_key_preview", creds.APIKey[:8]+"...").
		Str("redirection_url", redirectionURL).
		Str("environment", creds.Environment).
		Bool("is_sandbox", isSandbox).
		Msg("Generated Bold integrity signature and created pending tx")

	return &dtos.BoldSignatureResponse{
		OrderID:        orderID,
		Hash:           signature,
		Amount:         amount,
		Currency:       currency,
		PublicKey:      creds.APIKey,
		RedirectionURL: redirectionURL,
		IsSandbox:      isSandbox,
	}, nil
}

func (uc *walletUseCase) GetBoldStatus(ctx context.Context, boldOrderID string) (*dtos.BoldStatusResponse, error) {
	if boldOrderID == "" {
		return nil, fmt.Errorf("bold order id is required")
	}

	creds, err := uc.repo.GetBoldCredentials(ctx)
	if err != nil {
		return nil, err
	}

	client := resty.New().
		SetTimeout(10 * time.Second).
		SetRetryCount(3).
		SetRetryWaitTime(500 * time.Millisecond).
		SetRetryMaxWaitTime(3 * time.Second).
		AddRetryCondition(func(r *resty.Response, err error) bool {
			if err != nil {
				return true
			}
			return r.StatusCode() >= 500 || r.StatusCode() == 429
		})

	resp, err := client.R().
		SetContext(ctx).
		SetHeader("Authorization", "x-api-key "+creds.APIKey).
		SetHeader("Accept", "application/json").
		Get(creds.BaseURL + fmt.Sprintf(boldStatusEndpoint, boldOrderID))

	if err != nil {
		uc.log.Error(ctx).Err(err).Str("bold_order_id", boldOrderID).Msg("Bold status request failed")
		return nil, stderrors.Join(domainerrors.ErrBoldUpstreamUnavailable, err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
	case http.StatusNotFound:
		return nil, domainerrors.ErrBoldOrderNotFound
	case http.StatusUnauthorized, http.StatusForbidden:
		return nil, domainerrors.ErrBoldUnauthorized
	default:
		uc.log.Error(ctx).
			Int("status", resp.StatusCode()).
			Str("body", resp.String()).
			Str("bold_order_id", boldOrderID).
			Msg("Bold status returned non-OK")
		return nil, fmt.Errorf("%w: status %d", domainerrors.ErrBoldUpstreamUnavailable, resp.StatusCode())
	}

	var data struct {
		ID            string  `json:"id"`
		PaymentStatus string  `json:"payment_status"`
		Status        string  `json:"status"`
		Amount        float64 `json:"amount"`
		Currency      string  `json:"currency"`
		PaymentMethod string  `json:"payment_method"`
	}
	if err := json.Unmarshal(resp.Body(), &data); err != nil {
		return nil, fmt.Errorf("decode bold response: %w", err)
	}

	status := data.PaymentStatus
	if status == "" {
		status = data.Status
	}

	uc.log.Info(ctx).
		Str("bold_order_id", boldOrderID).
		Str("status", status).
		Msg("Fetched Bold transaction status")

	return &dtos.BoldStatusResponse{
		BoldOrderID:   data.ID,
		Status:        strings.ToUpper(status),
		Amount:        data.Amount,
		Currency:      data.Currency,
		PaymentMethod: data.PaymentMethod,
	}, nil
}
