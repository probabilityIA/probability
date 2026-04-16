package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/monitoring/internal/infra/primary/handlers/mappers"
)

func (h *handler) GetLogs(c *gin.Context) {
	id := c.Param("id")

	tail := 100
	if t := c.Query("tail"); t != "" {
		if parsed, err := strconv.Atoi(t); err == nil && parsed > 0 {
			tail = parsed
		}
	}

	logs, err := h.useCase.GetContainerLogs(c.Request.Context(), id, tail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, mappers.LogsToResponse(logs))
}
