package businesshandler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/bussines/internal/infra/primary/controllers/businesshandler/response"
)

const (
	businessesSimpleDefaultPageSize = 1000
	businessesSimpleMaxPageSize     = 1000
)

func (h *BusinessHandler) GetBusinessesSimple(c *gin.Context) {
	page := 1
	if raw := c.Query("page"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			page = parsed
		}
	}

	pageSize := businessesSimpleDefaultPageSize
	if raw := c.Query("page_size"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			if parsed > businessesSimpleMaxPageSize {
				parsed = businessesSimpleMaxPageSize
			}
			pageSize = parsed
		}
	}

	search := c.Query("search")

	isActive := true
	businesses, total, err := h.usecase.GetBusinesses(c.Request.Context(), page, pageSize, search, nil, &isActive)
	if err != nil {
		h.logger.Error().Err(err).Msg("Error getting businesses for simple list")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al obtener negocios",
			"error":   err.Error(),
		})
		return
	}

	simpleBusinesses := make([]response.BusinessSimpleResponse, 0, len(businesses))
	for _, business := range businesses {
		simpleBusinesses = append(simpleBusinesses, response.BusinessSimpleResponse{
			ID:              business.ID,
			Name:            business.Name,
			Code:            business.Code,
			LogoURL:         business.LogoURL,
			PrimaryColor:    business.PrimaryColor,
			SecondaryColor:  business.SecondaryColor,
			TertiaryColor:   business.TertiaryColor,
			QuaternaryColor: business.QuaternaryColor,
		})
	}

	totalPages := 0
	if pageSize > 0 {
		totalPages = int((total + int64(pageSize) - 1) / int64(pageSize))
	}

	c.JSON(http.StatusOK, response.GetBusinessesSimpleResponse{
		Success:    true,
		Message:    "Negocios obtenidos exitosamente",
		Data:       simpleBusinesses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}
