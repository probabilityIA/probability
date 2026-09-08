package whatsapp_template

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *handler) Variables(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "data": h.useCase.VariableCatalog()})
}
