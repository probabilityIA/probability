package handlers

import (
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/app/usecasenumbers"
)

const maxFotoBytes = 5 << 20

type businessProfileRequest struct {
	About       *string   `json:"about"`
	Address     *string   `json:"address"`
	Description *string   `json:"description"`
	Email       *string   `json:"email"`
	Vertical    *string   `json:"vertical"`
	Websites    *[]string `json:"websites"`
}

func (h *handler) GetBusinessProfile(c *gin.Context) {
	if !h.numbersReady(c) {
		return
	}

	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		return
	}

	perfil, err := h.numbersUseCase.GetProfile(c.Request.Context(), businessID)
	h.respondProfile(c, perfil, err, "consultando el perfil de WhatsApp")
}

func (h *handler) UpdateBusinessProfile(c *gin.Context) {
	if !h.numbersReady(c) {
		return
	}

	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		return
	}

	var req businessProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "no se pudo leer el formulario"})
		return
	}

	perfil, err := h.numbersUseCase.UpdateProfile(c.Request.Context(), businessID, usecasenumbers.BusinessProfileInput{
		About:       req.About,
		Address:     req.Address,
		Description: req.Description,
		Email:       req.Email,
		Vertical:    req.Vertical,
		Websites:    req.Websites,
	})
	h.respondProfile(c, perfil, err, "actualizando el perfil de WhatsApp")
}

func (h *handler) UpdateBusinessProfilePhoto(c *gin.Context) {
	if !h.numbersReady(c) {
		return
	}

	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		return
	}

	archivo, err := c.FormFile("photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "adjunta la imagen en el campo 'photo'"})
		return
	}

	if archivo.Size > maxFotoBytes {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "la imagen no puede pesar mas de 5 MB"})
		return
	}

	contentType := archivo.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "el archivo debe ser una imagen"})
		return
	}

	abierto, err := archivo.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "no se pudo leer la imagen"})
		return
	}
	defer abierto.Close()

	data, err := io.ReadAll(abierto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "no se pudo leer la imagen"})
		return
	}

	perfil, err := h.numbersUseCase.UpdateProfilePhoto(c.Request.Context(), businessID, contentType, data)
	h.respondProfile(c, perfil, err, "actualizando la foto del perfil de WhatsApp")
}

func (h *handler) respondProfile(c *gin.Context, perfil *usecasenumbers.BusinessProfile, err error, accion string) {
	if err != nil {
		h.log.Warn(c.Request.Context()).Err(err).Msgf("[WhatsApp Perfil] - error %s", accion)
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	if perfil == nil {
		c.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{
		"about":               perfil.About,
		"address":             perfil.Address,
		"description":         perfil.Description,
		"email":               perfil.Email,
		"profile_picture_url": perfil.ProfilePictureURL,
		"websites":            perfil.Websites,
		"vertical":            perfil.Vertical,
	}})
}
