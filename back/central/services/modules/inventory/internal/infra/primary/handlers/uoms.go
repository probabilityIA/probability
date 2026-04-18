package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	apprequest "github.com/secamc93/probability/back/central/services/modules/inventory/internal/app/request"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/infra/primary/handlers/response"
)

func (h *handlers) ListUoMs(c *gin.Context) {
	uoms, err := h.uc.ListUoMs(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	data := make([]response.UoMResponse, len(uoms))
	for i := range uoms {
		data[i] = response.UoMFromEntity(&uoms[i])
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *handlers) ListProductUoMs(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}
	productID := c.Param("productId")
	if productID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product_id is required"})
		return
	}
	list, err := h.uc.ListProductUoMs(c.Request.Context(), businessID, productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	data := make([]response.ProductUoMResponse, len(list))
	for i := range list {
		data[i] = response.ProductUoMFromEntity(&list[i])
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *handlers) CreateProductUoM(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}
	productID := c.Param("productId")
	if productID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product_id is required"})
		return
	}
	var body request.CreateProductUoMBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := h.uc.CreateProductUoM(c.Request.Context(), apprequest.CreateProductUoMDTO{
		BusinessID:       businessID,
		ProductID:        productID,
		UomCode:          body.UomCode,
		ConversionFactor: body.ConversionFactor,
		IsBase:           body.IsBase,
		Barcode:          body.Barcode,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, response.ProductUoMFromEntity(result))
}

func (h *handlers) DeleteProductUoM(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.uc.DeleteProductUoM(c.Request.Context(), businessID, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "product uom deleted"})
}

func (h *handlers) ConvertUoM(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}
	var body request.ConvertUoMBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := h.uc.ConvertUoM(c.Request.Context(), apprequest.ConvertUoMDTO{
		BusinessID:  businessID,
		ProductID:   body.ProductID,
		FromUomCode: body.FromUomCode,
		ToUomCode:   body.ToUomCode,
		Quantity:    body.Quantity,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
