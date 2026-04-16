package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/drivers/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/drivers/internal/infra/primary/handlers/response"
)

func (h *Handlers) ListDrivers(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	params := dtos.ListDriversParams{
		BusinessID: businessID,
		Search:     c.Query("search"),
		Status:     c.Query("status"),
		Page:       page,
		PageSize:   pageSize,
	}

	drivers, total, err := h.uc.ListDrivers(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	data := make([]response.DriverResponse, len(drivers))
	for i, d := range drivers {
		data[i] = response.FromEntity(&d)
	}

	totalPages := int(total) / params.PageSize
	if int(total)%params.PageSize != 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, response.DriversListResponse{
		Data:       data,
		Total:      total,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalPages: totalPages,
	})
}
