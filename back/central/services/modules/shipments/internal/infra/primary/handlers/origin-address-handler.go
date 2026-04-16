package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

// resolveBusinessID obtiene el business_id efectivo.
// Para usuarios normales usa el del JWT.
// Para super admins (business_id=0 en JWT) lee el query param ?business_id=X.
func (h *Handlers) resolveBusinessID(c *gin.Context) (uint, bool) {
	businessID, exists := middleware.GetBusinessID(c)
	if !exists {
		return 0, false
	}
	if businessID > 0 {
		return businessID, true
	}
	// Super admin: leer de query param
	if param := c.Query("business_id"); param != "" {
		if id, err := strconv.ParseUint(param, 10, 64); err == nil && id > 0 {
			return uint(id), true
		}
	}
	return 0, false
}

// ListOriginAddresses lista las direcciones de origen del comercio
func (h *Handlers) ListOriginAddresses(c *gin.Context) {
	businessID, exists := h.resolveBusinessID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No se pudo identificar la empresa"})
		return
	}

	addresses, err := h.uc.OriginAddress.List(c.Request.Context(), businessID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, addresses)
}

// CreateOriginAddress crea una nueva dirección de origen
func (h *Handlers) CreateOriginAddress(c *gin.Context) {
	businessID, exists := h.resolveBusinessID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No se pudo identificar la empresa"})
		return
	}

	var req domain.CreateOriginAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	address, err := h.uc.OriginAddress.Create(c.Request.Context(), businessID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, address)
}

// UpdateOriginAddress actualiza una dirección de origen
func (h *Handlers) UpdateOriginAddress(c *gin.Context) {
	businessID, exists := h.resolveBusinessID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No se pudo identificar la empresa"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var req domain.UpdateOriginAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	address, err := h.uc.OriginAddress.Update(c.Request.Context(), uint(id), businessID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, address)
}

// DeleteOriginAddress elimina una dirección de origen
func (h *Handlers) DeleteOriginAddress(c *gin.Context) {
	businessID, exists := h.resolveBusinessID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No se pudo identificar la empresa"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := h.uc.OriginAddress.Delete(c.Request.Context(), uint(id), businessID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Dirección eliminada correctamente"})
}
