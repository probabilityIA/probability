package handlers

import (
	"encoding/json"
	stderrors "errors"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	bcErrors "github.com/secamc93/probability/back/central/services/integrations/pay/bancolombia/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/integrations/pay/bancolombia/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

var signatureHeaderCandidates = []string{
	"x-bancolombia-signature",
	"x-signature",
	"signature",
}

type WebhookHandlers struct {
	useCase ports.IWebhookUseCase
	rawLog  ports.IRawWebhookLogger
	log     log.ILogger
}

func NewWebhookHandlers(useCase ports.IWebhookUseCase, rawLog ports.IRawWebhookLogger, logger log.ILogger) *WebhookHandlers {
	return &WebhookHandlers{
		useCase: useCase,
		rawLog:  rawLog,
		log:     logger.WithModule("bancolombia.webhook_handler"),
	}
}

func (h *WebhookHandlers) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/webhooks/bancolombia", h.handle(false))
	router.POST("/webhooks/bancolombia/test", h.handle(true))
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
		h.log.Error(ctx).Err(err).Msg("bancolombia webhook read body failed")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	signature := ""
	for _, header := range signatureHeaderCandidates {
		if v := c.GetHeader(header); v != "" {
			signature = v
			break
		}
	}

	endpoint := "prod"
	if isTest {
		endpoint = "test"
	}

	rawLog := h.buildRawLog(endpoint, signature, body)
	if h.rawLog != nil {
		if err := h.rawLog.LogIncoming(ctx, rawLog); err != nil {
			h.log.Warn(ctx).Err(err).Msg("bancolombia webhook: persist raw log failed (non-blocking)")
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
			h.log.Error(ctx).Err(useCaseErr).Msg("bancolombia webhook processing failed")
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
		var probe map[string]interface{}
		if err := json.Unmarshal(body, &probe); err == nil {
			raw.EventID = probeString(probe, "eventId", "event_id", "messageId", "id")
			raw.EventType = probeString(probe, "eventType", "event_type", "type")
			raw.Reference = probeString(probe, "transferReference", "paymentReference", "reference")
			raw.TransferCode = probeString(probe, "transferCode", "voucher", "cus")
		}
	}
	return raw
}

func probeString(m map[string]interface{}, keys ...string) string {
	for _, k := range keys {
		if v, ok := m[k].(string); ok && v != "" {
			return v
		}
	}
	for _, nested := range []string{"data", "transaction", "transfer"} {
		switch typed := m[nested].(type) {
		case map[string]interface{}:
			for _, k := range keys {
				if v, ok := typed[k].(string); ok && v != "" {
					return v
				}
			}
		case []interface{}:
			if len(typed) > 0 {
				if first, ok := typed[0].(map[string]interface{}); ok {
					for _, k := range keys {
						if v, ok := first[k].(string); ok && v != "" {
							return v
						}
					}
				}
			}
		}
	}
	return ""
}

func classifyResult(err error) (status, code string, httpStatus int, detail string) {
	if err == nil {
		return "ok", "ok", http.StatusOK, ""
	}
	switch {
	case stderrors.Is(err, bcErrors.ErrInvalidSignature):
		return "invalid_signature", "invalid signature", http.StatusUnauthorized, err.Error()
	case stderrors.Is(err, bcErrors.ErrConfigNotFound),
		stderrors.Is(err, bcErrors.ErrInvalidCredentials):
		return "config_missing", "bancolombia not configured", http.StatusServiceUnavailable, err.Error()
	default:
		detail = err.Error()
		if isParseError(err) {
			return "parse_error", "invalid payload", http.StatusBadRequest, detail
		}
		return "process_error", "internal", http.StatusInternalServerError, detail
	}
}

func isParseError(err error) bool {
	msg := strings.ToLower(err.Error())
	for _, fragment := range []string{"decode bancolombia webhook envelope", "empty body", "unmarshal"} {
		if strings.Contains(msg, fragment) {
			return true
		}
	}
	return false
}

func codeUpper(s string) string {
	return "BANCOLOMBIA_" + strings.ToUpper(strings.ReplaceAll(s, " ", "_"))
}
