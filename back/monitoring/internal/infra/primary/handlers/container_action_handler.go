package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	domainErrors "github.com/secamc93/probability/back/monitoring/internal/domain/errors"
)

func (h *handler) ContainerAction(action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		err := h.useCase.ContainerAction(c.Request.Context(), id, action)
		if err != nil {
			if errors.Is(err, domainErrors.ErrContainerNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			if errors.Is(err, domainErrors.ErrInvalidAction) {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("container %s: %s successful", id[:12], action)})
	}
}
