package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apprequest "github.com/secamc93/probability/back/central/services/modules/inventory/internal/app/request"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/infra/primary/handlers/response"
)

func (h *handlers) ListInventoryStates(c *gin.Context) {
	states, err := h.uc.ListInventoryStates(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	data := make([]response.InventoryStateResponse, len(states))
	for i := range states {
		data[i] = response.StateFromEntity(&states[i])
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *handlers) ChangeInventoryState(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}
	var body request.ChangeInventoryStateBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("user_id")
	var createdByID *uint
	if userID > 0 {
		createdByID = &userID
	}

	movement, err := h.uc.ChangeInventoryState(c.Request.Context(), apprequest.ChangeInventoryStateDTO{
		BusinessID:    businessID,
		LevelID:       body.LevelID,
		FromStateCode: body.FromStateCode,
		ToStateCode:   body.ToStateCode,
		Quantity:      body.Quantity,
		Reason:        body.Reason,
		CreatedByID:   createdByID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.StockMovementFromEntity(movement))
}
