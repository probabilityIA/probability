package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	apprequest "github.com/secamc93/probability/back/central/services/modules/inventory/internal/app/request"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/infra/primary/handlers/response"
)

func (h *handlers) CreateLot(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}
	var body request.CreateLotBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	lot, err := h.uc.CreateLot(c.Request.Context(), apprequest.CreateLotDTO{
		BusinessID:      businessID,
		ProductID:       body.ProductID,
		LotCode:         body.LotCode,
		ManufactureDate: body.ManufactureDate,
		ExpirationDate:  body.ExpirationDate,
		ReceivedAt:      body.ReceivedAt,
		SupplierID:      body.SupplierID,
		Status:          body.Status,
	})
	if err != nil {
		switch {
		case errors.Is(err, domainerrors.ErrProductNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, domainerrors.ErrDuplicateLotCode):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusCreated, response.LotFromEntity(lot))
}

func (h *handlers) GetLot(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}
	lotID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || lotID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lot id"})
		return
	}
	lot, err := h.uc.GetLot(c.Request.Context(), businessID, uint(lotID))
	if err != nil {
		if errors.Is(err, domainerrors.ErrLotNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.LotFromEntity(lot))
}

func (h *handlers) ListLots(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	expiringInDays, _ := strconv.Atoi(c.DefaultQuery("expiring_in_days", "0"))

	lots, total, err := h.uc.ListLots(c.Request.Context(), dtos.ListLotsParams{
		BusinessID:     businessID,
		ProductID:      c.Query("product_id"),
		Status:         c.Query("status"),
		ExpiringInDays: expiringInDays,
		Page:           page,
		PageSize:       pageSize,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	data := make([]response.LotResponse, len(lots))
	for i := range lots {
		data[i] = response.LotFromEntity(&lots[i])
	}
	if pageSize < 1 {
		pageSize = 10
	}
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	c.JSON(http.StatusOK, response.LotListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}

func (h *handlers) UpdateLot(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}
	lotID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || lotID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lot id"})
		return
	}
	var body request.UpdateLotBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	lot, err := h.uc.UpdateLot(c.Request.Context(), apprequest.UpdateLotDTO{
		ID:              uint(lotID),
		BusinessID:      businessID,
		LotCode:         body.LotCode,
		ManufactureDate: body.ManufactureDate,
		ExpirationDate:  body.ExpirationDate,
		ReceivedAt:      body.ReceivedAt,
		SupplierID:      body.SupplierID,
		Status:          body.Status,
	})
	if err != nil {
		switch {
		case errors.Is(err, domainerrors.ErrLotNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, domainerrors.ErrDuplicateLotCode):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, response.LotFromEntity(lot))
}

func (h *handlers) DeleteLot(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}
	lotID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || lotID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lot id"})
		return
	}
	if err := h.uc.DeleteLot(c.Request.Context(), businessID, uint(lotID)); err != nil {
		if errors.Is(err, domainerrors.ErrLotNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "lot deleted"})
}
