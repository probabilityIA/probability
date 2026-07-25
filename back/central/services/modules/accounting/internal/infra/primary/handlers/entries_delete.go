package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/errors"
)

func (h *Handlers) DeleteEntry(c *gin.Context) {
	if !h.requireSuperAdmin(c) {
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "id invalido"})
		return
	}
	if err := h.uc.DeleteManualEntry(c.Request.Context(), uint(id)); err != nil {
		switch {
		case errors.Is(err, domainerrors.ErrEntryNotFound):
			c.JSON(http.StatusNotFound, gin.H{"success": false, "error": err.Error()})
		case errors.Is(err, domainerrors.ErrEntryNotManual):
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
