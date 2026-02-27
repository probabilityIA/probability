package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/customers/internal/infra/primary/handlers/response"
)

// GetClient godoc
// @Summary      Obtener cliente por ID
// @Description  Obtiene un cliente específico con estadísticas de órdenes
// @Tags         Clients
// @Produce      json
// @Param        id  path  int  true  "ID del cliente"
// @Security     BearerAuth
// @Success      200  {object}  response.ClientDetailResponse
// @Router       /clients/{id} [get]
func (h *Handlers) GetClient(c *gin.Context) {
	businessID := c.GetUint("business_id")
	if businessID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "business_id not found in token"})
		return
	}

	clientID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || clientID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid client id"})
		return
	}

	client, err := h.uc.GetClient(c.Request.Context(), businessID, uint(clientID))
	if err != nil {
		if errors.Is(err, domainerrors.ErrClientNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.DetailFromEntity(client))
}
