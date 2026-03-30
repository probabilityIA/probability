package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	domainErrors "github.com/secamc93/probability/back/monitoring/internal/domain/errors"
	"github.com/secamc93/probability/back/monitoring/internal/infra/primary/handlers/mappers"
)

func (h *handler) GetContainer(c *gin.Context) {
	id := c.Param("id")

	container, err := h.useCase.GetContainer(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domainErrors.ErrContainerNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, mappers.ContainerToResponse(container))
}
