package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/infra/primary/handlers/mappers"
	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/infra/primary/handlers/response"
)

func (h *Handlers) CreateClient(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id es requerido"})
		return
	}

	requesterUserID := c.GetUint("user_id")

	var req request.CreateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dto := mappers.RequestToCreateClientDTO(&req)

	client, tempPassword, err := h.uc.CreateClient(c.Request.Context(), businessID, requesterUserID, dto)
	if err != nil {
		switch {
		case errors.Is(err, domainerrors.ErrStorefrontNotActive):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case errors.Is(err, domainerrors.ErrRoleNotAllowed):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case errors.Is(err, domainerrors.ErrEmailAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case errors.Is(err, domainerrors.ErrRoleNotFound):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "configuracion del sistema incompleta"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, response.CreateClientResponse{
		Client:       response.ClientFromEntity(client),
		TempPassword: tempPassword,
	})
}
