package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/drivers/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/drivers/internal/infra/primary/handlers/response"
)

func (h *Handlers) GetDriver(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}

	driverID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || driverID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid driver id"})
		return
	}

	driver, err := h.uc.GetDriver(c.Request.Context(), businessID, uint(driverID))
	if err != nil {
		if errors.Is(err, domainerrors.ErrDriverNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.FromEntity(driver))
}
