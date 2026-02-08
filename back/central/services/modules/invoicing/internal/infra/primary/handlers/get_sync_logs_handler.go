package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/primary/handlers/mappers"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/primary/handlers/response"
)

// GetInvoiceSyncLogs obtiene los logs de sincronización de una factura
func (h *handler) GetInvoiceSyncLogs(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error{
			Error:   "invalid_id",
			Message: "Invalid invoice ID",
		})
		return
	}

	logs, err := h.useCase.GetInvoiceSyncLogs(ctx, uint(id))
	if err != nil {
		h.log.Error(ctx).Err(err).Uint("invoice_id", uint(id)).Msg("Failed to get sync logs")
		c.JSON(http.StatusNotFound, response.Error{
			Error:   "sync_logs_not_found",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"sync_logs": mappers.SyncLogsToResponse(logs),
	})
}
