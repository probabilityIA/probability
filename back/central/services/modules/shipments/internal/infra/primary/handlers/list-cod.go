package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

func (h *Handlers) ListCODShipments(c *gin.Context) {
	businessID, err := h.resolveBusinessIDForCOD(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"error":   "invalid_business_id",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	filter := domain.CODFilter{
		BusinessID: businessID,
		Status:     strings.TrimSpace(c.Query("status")),
		Page:       page,
		PageSize:   pageSize,
	}

	if isPaidStr := c.Query("is_paid"); isPaidStr != "" {
		if v, err := strconv.ParseBool(isPaidStr); err == nil {
			filter.IsPaid = &v
		}
	}

	resp, err := h.uc.ListCODShipments(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al obtener envios contra entrega",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"message":     "Envios contra entrega obtenidos exitosamente",
		"data":        resp.Data,
		"total":       resp.Total,
		"page":        resp.Page,
		"page_size":   resp.PageSize,
		"total_pages": resp.TotalPages,
	})
}

func (h *Handlers) resolveBusinessIDForCOD(c *gin.Context) (uint, error) {
	businessID, exists := middleware.GetBusinessID(c)
	if !exists {
		return 0, errors.New("no se pudo identificar la empresa")
	}
	if !middleware.IsSuperAdmin(c) {
		return businessID, nil
	}
	param := c.Query("business_id")
	if param == "" {
		return 0, errors.New("super admin: business_id es requerido como query param")
	}
	id, err := strconv.ParseUint(param, 10, 64)
	if err != nil || id == 0 {
		return 0, errors.New("super admin: business_id invalido")
	}
	return uint(id), nil
}
