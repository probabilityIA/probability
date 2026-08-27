package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handlers) WooCommerceCheckoutConfig(c *gin.Context) {
	integrationID, ok := h.authWooPublic(c)
	if !ok {
		return
	}

	resolved, err := h.resolveWoo(c.Request.Context(), integrationID)
	if err != nil || resolved == nil {
		c.JSON(http.StatusOK, gin.H{
			"show_map": true,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"show_map": !resolved.DisableAddressMap,
	})
}
