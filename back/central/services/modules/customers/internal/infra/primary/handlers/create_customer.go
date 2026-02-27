package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/customers/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/modules/customers/internal/infra/primary/handlers/response"
)

// CreateClient godoc
// @Summary      Crear cliente
// @Description  Crea un nuevo cliente en el negocio
// @Tags         Clients
// @Accept       json
// @Produce      json
// @Param        request  body  request.CreateClientRequest  true  "Datos del cliente"
// @Security     BearerAuth
// @Success      201  {object}  response.ClientResponse
// @Router       /clients [post]
func (h *Handlers) CreateClient(c *gin.Context) {
	businessID := c.GetUint("business_id")
	if businessID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "business_id not found in token"})
		return
	}

	var req request.CreateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dto := dtos.CreateClientDTO{
		BusinessID: businessID,
		Name:       req.Name,
		Email:      req.Email,
		Phone:      req.Phone,
		Dni:        req.Dni,
	}

	client, err := h.uc.CreateClient(c.Request.Context(), dto)
	if err != nil {
		if errors.Is(err, domainerrors.ErrDuplicateEmail) || errors.Is(err, domainerrors.ErrDuplicateDni) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response.FromEntity(client))
}
