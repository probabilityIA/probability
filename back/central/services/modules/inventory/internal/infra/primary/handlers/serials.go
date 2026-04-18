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

func parseOptionalUintQuery(c *gin.Context, key string) *uint {
	val := c.Query(key)
	if val == "" {
		return nil
	}
	n, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return nil
	}
	u := uint(n)
	return &u
}

func (h *handlers) CreateSerial(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}
	var body request.CreateSerialBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	serial, err := h.uc.CreateSerial(c.Request.Context(), apprequest.CreateSerialDTO{
		BusinessID:   businessID,
		ProductID:    body.ProductID,
		SerialNumber: body.SerialNumber,
		LotID:        body.LotID,
		LocationID:   body.LocationID,
		StateCode:    body.StateCode,
	})
	if err != nil {
		switch {
		case errors.Is(err, domainerrors.ErrProductNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, domainerrors.ErrDuplicateSerial):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusCreated, response.SerialFromEntity(serial))
}

func (h *handlers) GetSerial(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}
	serialID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || serialID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid serial id"})
		return
	}
	serial, err := h.uc.GetSerial(c.Request.Context(), businessID, uint(serialID))
	if err != nil {
		if errors.Is(err, domainerrors.ErrSerialNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.SerialFromEntity(serial))
}

func (h *handlers) ListSerials(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	serials, total, err := h.uc.ListSerials(c.Request.Context(), dtos.ListSerialsParams{
		BusinessID: businessID,
		ProductID:  c.Query("product_id"),
		LotID:      parseOptionalUintQuery(c, "lot_id"),
		StateID:    parseOptionalUintQuery(c, "state_id"),
		LocationID: parseOptionalUintQuery(c, "location_id"),
		Page:       page,
		PageSize:   pageSize,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	data := make([]response.SerialResponse, len(serials))
	for i := range serials {
		data[i] = response.SerialFromEntity(&serials[i])
	}
	if pageSize < 1 {
		pageSize = 10
	}
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	c.JSON(http.StatusOK, response.SerialListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}

func (h *handlers) UpdateSerial(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}
	serialID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || serialID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid serial id"})
		return
	}
	var body request.UpdateSerialBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	serial, err := h.uc.UpdateSerial(c.Request.Context(), apprequest.UpdateSerialDTO{
		ID:         uint(serialID),
		BusinessID: businessID,
		LotID:      body.LotID,
		LocationID: body.LocationID,
		StateCode:  body.StateCode,
	})
	if err != nil {
		if errors.Is(err, domainerrors.ErrSerialNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.SerialFromEntity(serial))
}
