package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/monitoring/internal/infra/primary/handlers/mappers"
)

func (h *handler) ListContainers(c *gin.Context) {
	containers, err := h.useCase.ListContainers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, mappers.ContainersToResponse(containers))
}
