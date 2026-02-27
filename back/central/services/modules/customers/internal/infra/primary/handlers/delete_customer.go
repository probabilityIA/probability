package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/errors"
)

// DeleteClient godoc
// @Summary      Eliminar cliente
// @Description  Elimina un cliente. Retorna 409 si tiene órdenes asociadas.
// @Tags         Clients
// @Produce      json
// @Param        id  path  int  true  "ID del cliente"
// @Security     BearerAuth
// @Success      200  {object}  map[string]string
// @Failure      409  {object}  map[string]string
// @Router       /clients/{id} [delete]
func (h *Handlers) DeleteClient(c *gin.Context) {
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

	if err := h.uc.DeleteClient(c.Request.Context(), businessID, uint(clientID)); err != nil {
		if errors.Is(err, domainerrors.ErrClientNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, domainerrors.ErrClientHasOrders) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "client deleted successfully"})
}
