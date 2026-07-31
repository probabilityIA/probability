package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/publicsite/internal/domain/errors"
)

func (h *Handlers) GetMyOrders(c *gin.Context) {
	slug := c.Param("slug")
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "sesion no valida"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	orders, total, err := h.uc.GetMyOrders(c.Request.Context(), slug, userID, page, pageSize)
	if err != nil {
		if errors.Is(err, domainerrors.ErrBusinessNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "negocio no encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "error consultando pedidos"})
		return
	}

	items := make([]gin.H, 0, len(orders))
	for _, o := range orders {
		items = append(items, gin.H{
			"id":           o.ID,
			"order_number": o.OrderNumber,
			"total":        o.Total,
			"currency":     o.Currency,
			"status":       o.Status,
			"created_at":   o.CreatedAt,
			"items_count":  o.ItemsCount,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    items,
		"total":   total,
		"page":    page,
	})
}
