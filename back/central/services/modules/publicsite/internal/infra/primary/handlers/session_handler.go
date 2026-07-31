package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/publicsite/internal/domain/errors"
)

func (h *Handlers) GetSession(c *gin.Context) {
	slug := c.Param("slug")
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "sesion no valida"})
		return
	}

	session, err := h.uc.GetSession(c.Request.Context(), slug, userID)
	if err != nil {
		if errors.Is(err, domainerrors.ErrBusinessNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "sesion no asociada a este negocio"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "error consultando la sesion"})
		return
	}

	var address gin.H
	if session.Address != nil {
		address = gin.H{
			"street":      session.Address.Street,
			"city":        session.Address.City,
			"state":       session.Address.State,
			"country":     session.Address.Country,
			"postal_code": session.Address.PostalCode,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"name":        session.Name,
			"avatar_url":  session.AvatarURL,
			"email":       session.Email,
			"phone":       session.Phone,
			"dni":         session.Dni,
			"customer_id": session.CustomerID,
			"address":     address,
		},
	})
}
