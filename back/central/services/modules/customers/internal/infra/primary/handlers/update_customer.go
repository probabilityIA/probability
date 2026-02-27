package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/customers/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/modules/customers/internal/infra/primary/handlers/response"
)

func (h *Handlers) UpdateClient(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}

	clientID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || clientID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid client id"})
		return
	}

	var req request.UpdateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dto := dtos.UpdateClientDTO{
		ID:         uint(clientID),
		BusinessID: businessID,
		Name:       req.Name,
		Email:      req.Email,
		Phone:      req.Phone,
		Dni:        req.Dni,
	}

	client, err := h.uc.UpdateClient(c.Request.Context(), dto)
	if err != nil {
		if errors.Is(err, domainerrors.ErrClientNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, domainerrors.ErrDuplicateEmail) || errors.Is(err, domainerrors.ErrDuplicateDni) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.FromEntity(client))
}
