package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

func (h *Handlers) requireSuperAdmin(c *gin.Context) bool {
	if !middleware.IsSuperAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "solo super admin puede administrar la tienda"})
		return false
	}
	return true
}

func (h *Handlers) Status(c *gin.Context) {
	if !h.requireSuperAdmin(c) {
		return
	}
	st, err := h.uc.Status(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, st)
}

func (h *Handlers) Start(c *gin.Context) {
	if !h.requireSuperAdmin(c) {
		return
	}
	st, err := h.uc.Start(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, st)
}

func (h *Handlers) Stop(c *gin.Context) {
	if !h.requireSuperAdmin(c) {
		return
	}
	st, err := h.uc.Stop(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, st)
}
