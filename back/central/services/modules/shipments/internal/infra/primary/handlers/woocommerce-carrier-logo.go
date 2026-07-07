package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/app/usecaseshipment"
)

func (h *Handlers) WooCommerceCarrierLogo(c *gin.Context) {
	carrier := c.Param("carrier")
	data := usecaseshipment.CarrierLogoBytes(carrier)
	if len(data) == 0 {
		c.Status(http.StatusNotFound)
		return
	}
	c.Header("Cache-Control", "public, max-age=86400")
	c.Data(http.StatusOK, "image/png", data)
}
