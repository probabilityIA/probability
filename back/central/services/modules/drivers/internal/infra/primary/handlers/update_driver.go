package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/drivers/internal/domain/dtos"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/drivers/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/drivers/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/modules/drivers/internal/infra/primary/handlers/response"
)

func (h *Handlers) UpdateDriver(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}

	driverID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || driverID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid driver id"})
		return
	}

	var req request.UpdateDriverRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dto := dtos.UpdateDriverDTO{
		ID:             uint(driverID),
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

	driver, err := h.uc.UpdateDriver(c.Request.Context(), dto)
	if err != nil {
		if errors.Is(err, domainerrors.ErrDriverNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, domainerrors.ErrDuplicateIdentification) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.FromEntity(driver))
}
