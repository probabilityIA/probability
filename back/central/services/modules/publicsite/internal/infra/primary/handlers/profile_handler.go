package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/publicsite/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/publicsite/internal/domain/errors"
)

type updateProfileRequest struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Dni   string `json:"dni"`
}

type addressRequest struct {
	Street     string `json:"street" binding:"required"`
	City       string `json:"city"`
	State      string `json:"state"`
	Country    string `json:"country"`
	PostalCode string `json:"postal_code"`
	IsPrimary  bool   `json:"is_primary"`
}

func (h *Handlers) respondProfileError(c *gin.Context, err error) {
	if errors.Is(err, domainerrors.ErrBusinessNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "sesion no asociada a este negocio"})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
}

func (h *Handlers) UpdateProfile(c *gin.Context) {
	slug := c.Param("slug")
	userID := c.GetUint("user_id")

	var req updateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "datos invalidos"})
		return
	}

	if err := h.uc.UpdateProfile(c.Request.Context(), slug, userID, strings.TrimSpace(req.Name), strings.TrimSpace(req.Phone), strings.TrimSpace(req.Dni)); err != nil {
		h.respondProfileError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Perfil actualizado"})
}

func (h *Handlers) UploadAvatar(c *gin.Context) {
	slug := c.Param("slug")
	userID := c.GetUint("user_id")

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "imagen es requerida"})
		return
	}

	folder := fmt.Sprintf("avatars/%d", userID)
	relativePath, err := h.s3.UploadImage(c.Request.Context(), file, folder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "error subiendo imagen"})
		return
	}

	avatarURL := relativePath
	if base := h.env.Get("URL_BASE_DOMAIN_S3"); base != "" && !strings.HasPrefix(avatarURL, "http") {
		avatarURL = base + avatarURL
	}

	if err := h.uc.SaveAvatar(c.Request.Context(), slug, userID, avatarURL); err != nil {
		h.respondProfileError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "avatar_url": avatarURL})
}

func (h *Handlers) ListAddresses(c *gin.Context) {
	slug := c.Param("slug")
	userID := c.GetUint("user_id")

	addresses, err := h.uc.ListAddresses(c.Request.Context(), slug, userID)
	if err != nil {
		h.respondProfileError(c, err)
		return
	}

	items := make([]gin.H, 0, len(addresses))
	for _, a := range addresses {
		items = append(items, gin.H{
			"id":          a.ID,
			"street":      a.Street,
			"city":        a.City,
			"state":       a.State,
			"country":     a.Country,
			"postal_code": a.PostalCode,
			"is_primary":  a.IsPrimary,
		})
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": items})
}

func (h *Handlers) AddAddress(c *gin.Context) {
	slug := c.Param("slug")
	userID := c.GetUint("user_id")

	var req addressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "direccion invalida"})
		return
	}

	addr := &entities.CustomerAddress{
		Street:     strings.TrimSpace(req.Street),
		City:       strings.TrimSpace(req.City),
		State:      strings.TrimSpace(req.State),
		Country:    strings.TrimSpace(req.Country),
		PostalCode: strings.TrimSpace(req.PostalCode),
		IsPrimary:  req.IsPrimary,
	}
	if addr.Country == "" {
		addr.Country = "Colombia"
	}

	if err := h.uc.AddAddress(c.Request.Context(), slug, userID, addr); err != nil {
		h.respondProfileError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Direccion agregada"})
}

func (h *Handlers) SetPrimaryAddress(c *gin.Context) {
	slug := c.Param("slug")
	userID := c.GetUint("user_id")
	addressID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || addressID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "id invalido"})
		return
	}

	if err := h.uc.SetPrimaryAddress(c.Request.Context(), slug, userID, uint(addressID)); err != nil {
		h.respondProfileError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Direccion principal actualizada"})
}

func (h *Handlers) DeleteAddress(c *gin.Context) {
	slug := c.Param("slug")
	userID := c.GetUint("user_id")
	addressID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || addressID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "id invalido"})
		return
	}

	if err := h.uc.DeleteAddress(c.Request.Context(), slug, userID, uint(addressID)); err != nil {
		h.respondProfileError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Direccion eliminada"})
}
