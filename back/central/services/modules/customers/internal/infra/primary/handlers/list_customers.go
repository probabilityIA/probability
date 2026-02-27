package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/customers/internal/infra/primary/handlers/response"
)

// ListClients godoc
// @Summary      Listar clientes
// @Description  Obtiene una lista paginada de clientes del negocio
// @Tags         Clients
// @Produce      json
// @Param        page       query  int     false  "Página (default: 1)"
// @Param        page_size  query  int     false  "Tamaño de página (default: 20, max: 100)"
// @Param        search     query  string  false  "Buscar por nombre, email o teléfono"
// @Security     BearerAuth
// @Success      200  {object}  response.ClientsListResponse
// @Router       /clients [get]
func (h *Handlers) ListClients(c *gin.Context) {
	businessID := c.GetUint("business_id")
	if businessID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "business_id not found in token"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	params := dtos.ListClientsParams{
		BusinessID: businessID,
		Search:     c.Query("search"),
		Page:       page,
		PageSize:   pageSize,
	}

	clients, total, err := h.uc.ListClients(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	data := make([]response.ClientResponse, len(clients))
	for i, cl := range clients {
		data[i] = response.FromEntity(&cl)
	}

	totalPages := int(total) / params.PageSize
	if int(total)%params.PageSize != 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, response.ClientsListResponse{
		Data:       data,
		Total:      total,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalPages: totalPages,
	})
}
