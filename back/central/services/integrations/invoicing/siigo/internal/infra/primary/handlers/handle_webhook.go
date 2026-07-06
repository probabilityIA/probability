package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type siigoWebhookPayload struct {
	Topic      string `json:"topic"`
	CompanyKey string `json:"company_key"`
}

type inventorySyncRequest struct {
	InvoiceID     uint              `json:"invoice_id"`
	Provider      string            `json:"provider"`
	Operation     string            `json:"operation"`
	InvoiceData   inventorySyncData `json:"invoice_data"`
	CorrelationID string            `json:"correlation_id"`
	Timestamp     time.Time         `json:"timestamp"`
}

type inventorySyncData struct {
	IntegrationID uint                   `json:"integration_id"`
	Config        map[string]interface{} `json:"config"`
}

func (h *Handler) HandleWebhook(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot read body"})
		return
	}

	var payload siigoWebhookPayload
	_ = json.Unmarshal(body, &payload)

	var integrationID *uint
	if raw := c.Query("integration_id"); raw != "" {
		if id, perr := strconv.ParseUint(raw, 10, 64); perr == nil && id > 0 {
			v := uint(id)
			integrationID = &v
		}
	}

	logID, logErr := h.webhookLog.LogIncoming(c.Request.Context(), ports.WebhookLogEntry{
		Source:              "siigo",
		EventType:           payload.Topic,
		URL:                 c.Request.URL.Path,
		Body:                body,
		RemoteIP:            c.ClientIP(),
		IntegrationID:       integrationID,
		IntegrationTypeCode: "siigo",
	})
	if logErr != nil {
		h.log.Error(c.Request.Context()).Err(logErr).Msg("Error al registrar webhook Siigo en historico")
	}

	c.JSON(http.StatusOK, gin.H{"received": true})

	go h.process(logID, integrationID, payload.Topic)
}

func (h *Handler) process(logID string, integrationID *uint, topic string) {
	ctx := context.Background()

	if integrationID == nil {
		h.finish(ctx, logID, "failed", http.StatusBadRequest, "missing integration_id")
		return
	}

	integration, err := h.coreIntegration.GetIntegrationByID(ctx, strconv.FormatUint(uint64(*integrationID), 10))
	if err != nil || integration == nil {
		h.finish(ctx, logID, "failed", http.StatusNotFound, "integration not found")
		return
	}

	var businessID uint
	if integration.BusinessID != nil {
		businessID = *integration.BusinessID
	}

	if err := h.publishInventorySync(ctx, *integrationID, businessID); err != nil {
		h.log.Error(ctx).Err(err).Uint("integration_id", *integrationID).Msg("Error al disparar sync de inventario desde webhook Siigo")
		h.finish(ctx, logID, "failed", http.StatusInternalServerError, err.Error())
		return
	}

	h.log.Info(ctx).
		Uint("integration_id", *integrationID).
		Uint("business_id", businessID).
		Str("topic", topic).
		Msg("Webhook Siigo procesado, sync de inventario disparado")

	h.finish(ctx, logID, "ok", http.StatusOK, "")
}

func (h *Handler) publishInventorySync(ctx context.Context, integrationID, businessID uint) error {
	if h.rabbit == nil {
		return fmt.Errorf("rabbitmq no disponible")
	}

	msg := inventorySyncRequest{
		InvoiceID: 0,
		Provider:  "siigo",
		Operation: "inventory_sync",
		InvoiceData: inventorySyncData{
			IntegrationID: integrationID,
			Config: map[string]interface{}{
				"business_id": businessID,
			},
		},
		CorrelationID: uuid.New().String(),
		Timestamp:     time.Now(),
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	if err := h.rabbit.DeclareQueue(rabbitmq.QueueInvoicingSiigoRequests, true); err != nil {
		return err
	}

	return h.rabbit.Publish(ctx, rabbitmq.QueueInvoicingSiigoRequests, data)
}

func (h *Handler) finish(ctx context.Context, logID, status string, httpStatus int, errMessage string) {
	if logID == "" {
		return
	}
	if err := h.webhookLog.UpdateResult(ctx, logID, status, httpStatus, errMessage); err != nil {
		h.log.Error(ctx).Err(err).Msg("Error al actualizar resultado de webhook Siigo")
	}
}
