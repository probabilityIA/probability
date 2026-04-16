package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/drivers/internal/domain/dtos"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/drivers/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/drivers/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/modules/drivers/internal/infra/primary/handlers/response"
)

func (h *Handlers) CreateDriver(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}

	var req request.CreateDriverRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dto := dtos.CreateDriverDTO{
		BusinessID:     businessID,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Email:          req.Email,
		Phone:          req.Phone,
		Identification: req.Identification,
		Status:         req.Status,
		PhotoURL:       req.PhotoURL,
		LicenseType:    req.LicenseType,
		LicenseExpiry:  req.LicenseExpiry,
		WarehouseID:    req.WarehouseID,
		Notes:          req.Notes,
	}

	driver, err := h.uc.CreateDriver(c.Request.Context(), dto)
	if err != nil {
		if errors.Is(err, domainerrors.ErrDuplicateIdentification) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response.FromEntity(driver))
}
