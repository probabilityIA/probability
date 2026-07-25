package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/dtos"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/errors"
)

func (h *Handlers) Report(c *gin.Context) {
	if !h.requireSuperAdmin(c) {
		return
	}
	from, err := time.Parse("2006-01-02", c.Query("from"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "from es requerido (YYYY-MM-DD)"})
		return
	}
	to, err := time.Parse("2006-01-02", c.Query("to"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "to es requerido (YYYY-MM-DD)"})
		return
	}
	report, err := h.uc.Report(c.Request.Context(), dtos.ReportParams{From: from, To: to})
	if err != nil {
		if errors.Is(err, domainerrors.ErrInvalidPeriod) {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": report})
}
