package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/infra/primary/handlers/response"
)

func (h *Handlers) ListOverrides(c *gin.Context) {
	if !h.requireSuperAdmin(c) {
		return
	}

	businessID, err := strconv.ParseUint(c.Param("businessId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid business id"})
		return
	}

	overrides, err := h.uc.ListOverrides(c.Request.Context(), uint(businessID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list overrides"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": response.FromOverrides(overrides)})
}

func (h *Handlers) GrantOverride(c *gin.Context) {
	if !h.requireSuperAdmin(c) {
		return
	}

	userID, _ := middleware.GetUserID(c)

	var req request.GrantOverrideRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var expiresAt *time.Time
	if req.ExpiresAt != nil && *req.ExpiresAt != "" {
		parsed, err := time.Parse("2006-01-02", *req.ExpiresAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "vencimiento invalido, formato esperado YYYY-MM-DD"})
			return
		}
		expiresAt = &parsed
	}

	err := h.uc.GrantOverride(c.Request.Context(), dtos.GrantOverrideDTO{
		BusinessID:      req.BusinessID,
		ModuleCode:      req.ModuleCode,
		Notes:           req.Notes,
		ExpiresAt:       expiresAt,
		GrantedByUserID: userID,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "module override granted"})
}

func (h *Handlers) RevokeOverride(c *gin.Context) {
	if !h.requireSuperAdmin(c) {
		return
	}

	businessID, err := strconv.ParseUint(c.Param("businessId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid business id"})
		return
	}

	moduleCode := c.Param("moduleCode")
	if moduleCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "module code is required"})
		return
	}

	userID, _ := middleware.GetUserID(c)
	if err := h.uc.RevokeOverride(c.Request.Context(), uint(businessID), moduleCode, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to revoke override"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "module override revoked"})
}
