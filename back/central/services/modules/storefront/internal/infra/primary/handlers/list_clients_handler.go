package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/infra/primary/handlers/response"
)

func (h *Handlers) ListClients(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id es requerido"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	clients, total, err := h.uc.ListClients(c.Request.Context(), businessID, page, pageSize)
	if err != nil {
		if errors.Is(err, domainerrors.ErrStorefrontNotActive) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if pageSize < 1 {
		pageSize = 20
	}
	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, response.ClientListResponse{
		Data:       response.ClientsFromEntities(clients),
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}
