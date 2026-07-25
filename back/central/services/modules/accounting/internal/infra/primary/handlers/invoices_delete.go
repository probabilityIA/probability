package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handlers) DeleteInvoice(c *gin.Context) {
	if !h.requireSuperAdmin(c) {
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "id invalido"})
		return
	}
	if err := h.uc.DeleteInvoice(c.Request.Context(), uint(id)); err != nil {
		h.invoiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
