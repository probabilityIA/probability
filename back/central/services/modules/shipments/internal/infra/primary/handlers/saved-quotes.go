package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

func (h *Handlers) ListSavedQuotes(c *gin.Context) {
	businessID, ok := h.resolveBusinessIDParam(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no se pudo identificar la empresa"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	filter := domain.SavedQuoteFilter{
		BusinessID: businessID,
		Source:     c.Query("source"),
		Status:     c.Query("status"),
		OrderUUID:  c.Query("order_uuid"),
		Page:       page,
		PageSize:   pageSize,
	}

	resp, err := h.uc.Quotes.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handlers) GetSavedQuote(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id invalido"})
		return
	}
	resp, err := h.uc.Quotes.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if resp == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cotizacion no encontrada"})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handlers) resolveBusinessIDParam(c *gin.Context) (uint, bool) {
	businessID, exists := middleware.GetBusinessID(c)
	if !exists {
		return 0, false
	}
	if businessID > 0 {
		return businessID, true
	}
	if param := c.Query("business_id"); param != "" {
		if id, err := strconv.ParseUint(param, 10, 64); err == nil && id > 0 {
			return uint(id), true
		}
	}
	return 0, false
}
