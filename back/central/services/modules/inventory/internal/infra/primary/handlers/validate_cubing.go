package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apprequest "github.com/secamc93/probability/back/central/services/modules/inventory/internal/app/request"
)

type validateCubingBody struct {
	ProductID  string `json:"product_id" binding:"required"`
	LocationID uint   `json:"location_id" binding:"required,min=1"`
	Quantity   int    `json:"quantity" binding:"required,min=1"`
}

func (h *handlers) ValidateCubing(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}

	var body validateCubingBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.uc.ValidateCubing(c.Request.Context(), apprequest.ValidateCubingDTO{
		ProductID:  body.ProductID,
		LocationID: body.LocationID,
		BusinessID: businessID,
		Quantity:   body.Quantity,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
