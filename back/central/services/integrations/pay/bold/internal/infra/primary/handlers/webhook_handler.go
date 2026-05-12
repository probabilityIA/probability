package handlers

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	boldErrors "github.com/secamc93/probability/back/central/services/integrations/pay/bold/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/integrations/pay/bold/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

const boldSignatureHeader = "x-bold-signature"

type WebhookHandlers struct {
	useCase ports.IWebhookUseCase
	rawLog  ports.IRawWebhookLogger
	log     log.ILogger
}

func NewWebhookHandlers(useCase ports.IWebhookUseCase, rawLog ports.IRawWebhookLogger, logger log.ILogger) *WebhookHandlers {
	return &WebhookHandlers{
		useCase: useCase,
		rawLog:  rawLog,
		log:     logger.WithModule("bold.webhook_handler"),
	}
}

func (h *WebhookHandlers) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/webhooks/bold", h.handle(false))
	router.POST("/webhooks/bold/test", h.handle(true))
}

func (h *WebhookHandlers) handle(isTest bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		h.HandleWebhook(c, isTest)
	}
}

func (h *WebhookHandlers) HandleWebhook(c *gin.Context, isTest bool) {
	ctx := c.Request.Context()
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		h.log.Error(ctx).Err(err).Msg("bold webhook read body failed")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	signature := c.GetHeader(boldSignatureHeader)
	endpoint := "prod"
	if isTest {
		endpoint = "test"
	}

	rawLog := h.buildRawLog(endpoint, signature, body)
	if h.rawLog != nil {
		if err := h.rawLog.LogIncoming(ctx, rawLog); err != nil {
			h.log.Warn(ctx).Err(err).Msg("bold webhook: persist raw log failed (non-blocking)")
		}
	}

	useCaseErr := h.useCase.HandleIncomingWebhook(ctx, signature, body, isTest)

	status, code, httpStatus, errDetail := classifyResult(useCaseErr)
	if h.rawLog != nil && rawLog.ID != "" {
		_ = h.rawLog.UpdateResult(ctx, &ports.RawWebhookResult{
			ID:          rawLog.ID,
			Status:      status,
			HTTPStatus:  httpStatus,
			ErrorDetail: errDetail,
		})
	}

	if useCaseErr != nil {
		if status == "process_error" {
			h.log.Error(ctx).Err(useCaseErr).Msg("bold webhook processing failed")
		}
		c.JSON(httpStatus, gin.H{"error": code, "code": codeUpper(code)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"received": true})
}

func (h *WebhookHandlers) buildRawLog(endpoint, signature string, body []byte) *ports.RawWebhookLog {
	raw := &ports.RawWebhookLog{
		ID:              uuid.New().String(),
		Endpoint:        endpoint,
		SignatureHeader: signature,
		BodySize:        len(body),
		Body:            body,
	}
	if len(body) > 0 && json.Valid(body) {
		var probe struct {
			ID      string `json:"id"`
			Type    string `json:"type"`
			Subject string `json:"subject"`
			Data    struct {
				PaymentID         string `json:"payment_id"`
				MerchantReference string `json:"merchant_reference"`
			} `json:"data"`
		}
		if err := json.Unmarshal(body, &probe); err == nil {
			raw.BoldEventID = probe.ID
			raw.EventType = probe.Type
			raw.MerchantReference = probe.Data.MerchantReference
			raw.PaymentID = probe.Data.PaymentID
		}
	}
	return raw
}

func classifyResult(err error) (status, code string, httpStatus int, detail string) {
	if err == nil {
		return "ok", "ok", http.StatusOK, ""
	}
	switch {
	case stderrors.Is(err, boldErrors.ErrInvalidSignature):
		return "invalid_signature", "invalid signature", http.StatusUnauthorized, err.Error()
	case stderrors.Is(err, boldErrors.ErrBoldConfigNotFound),
		stderrors.Is(err, boldErrors.ErrInvalidCredentials):
		return "config_missing", "bold not configured", http.StatusServiceUnavailable, err.Error()
	default:
		detail = err.Error()
		if isParseError(err) {
			return "parse_error", "invalid payload", http.StatusBadRequest, detail
		}
		return "process_error", "internal", http.StatusInternalServerError, detail
	}
}

func isParseError(err error) bool {
	msg := err.Error()
	for _, fragment := range []string{"decode bold webhook envelope", "missing event id", "unmarshal"} {
		if containsFold(msg, fragment) {
			return true
		}
	}
	return false
}

func containsFold(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	lo := []byte(toLower(s))
	sub := []byte(toLower(substr))
	for i := 0; i+len(sub) <= len(lo); i++ {
		match := true
		for j := 0; j < len(sub); j++ {
			if lo[i+j] != sub[j] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

func toLower(s string) string {
	b := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 32
		}
		b[i] = c
	}
	return string(b)
}

func codeUpper(s string) string {
	b := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c >= 'a' && c <= 'z':
			b[i] = c - 32
		case c == ' ':
			b[i] = '_'
		default:
			b[i] = c
		}
	}
	return "BOLD_" + string(b)
}

var _ = context.Background
