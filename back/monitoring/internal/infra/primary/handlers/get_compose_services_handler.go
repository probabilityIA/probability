package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/monitoring/internal/infra/primary/handlers/mappers"
)

func (h *handler) GetComposeServices(c *gin.Context) {
	services, err := h.useCase.GetComposeServices(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, mappers.ComposeServicesToResponse(services))
}
