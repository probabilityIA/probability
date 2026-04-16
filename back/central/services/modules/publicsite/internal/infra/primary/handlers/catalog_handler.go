package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/publicsite/internal/domain/dtos"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/publicsite/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/publicsite/internal/infra/primary/handlers/response"
)

func (h *Handlers) ListCatalog(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "slug es requerido"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "12"))

	filters := dtos.CatalogFilters{
		Search:   c.Query("search"),
		Category: c.Query("category"),
		Page:     page,
		PageSize: pageSize,
	}

	products, total, err := h.uc.ListCatalog(c.Request.Context(), slug, filters)
	if err != nil {
		if errors.Is(err, domainerrors.ErrBusinessNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, domainerrors.ErrPublicSiteNotActive) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	filters.Normalize()
	totalPages := int(total) / filters.PageSize
	if int(total)%filters.PageSize != 0 {
		totalPages++
	}

	imageURLBase := h.getImageURLBase()

	c.JSON(http.StatusOK, response.CatalogListResponse{
		Data:       response.ProductsFromEntities(products, imageURLBase),
		Total:      total,
		Page:       filters.Page,
		PageSize:   filters.PageSize,
		TotalPages: totalPages,
	})
}
