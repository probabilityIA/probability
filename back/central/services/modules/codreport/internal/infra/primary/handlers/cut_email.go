package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/modules/codreport/internal/domain/dtos"
)

type sendCutEmailRequest struct {
	Emails []string `json:"emails"`
}

func (h *Handlers) SendCutEmail(c *gin.Context) {
	if !isAdminUser(c) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Solo un administrador puede enviar el corte por correo"})
		return
	}

	businessID, err := resolveBusinessID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	cutID, err := strconv.ParseUint(c.Query("cut_id"), 10, 64)
	if err != nil || cutID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "cut_id invalido"})
		return
	}

	var req sendCutEmailRequest
	_ = c.ShouldBindJSON(&req)

	recipients := make([]string, 0, len(req.Emails))
	seen := map[string]struct{}{}
	for _, e := range req.Emails {
		v := strings.ToLower(strings.TrimSpace(e))
		if v == "" {
			continue
		}
		if !strings.Contains(v, "@") || strings.ContainsAny(v, " ,;") {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Correo invalido: " + v})
			return
		}
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		recipients = append(recipients, v)
	}
	if len(recipients) == 0 {
		if own, ok := middleware.GetUserEmail(c); ok && own != "" {
			recipients = append(recipients, strings.ToLower(strings.TrimSpace(own)))
		}
	}
	if len(recipients) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Debes indicar al menos un correo"})
		return
	}
	if len(recipients) > 10 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Maximo 10 destinatarios por envio"})
		return
	}

	userID, _ := middleware.GetUserID(c)

	if err := h.uc.SendCutEmail(c.Request.Context(), dtos.SendCutEmailDTO{
		BusinessID: businessID,
		CutID:      uint(cutID),
		Recipients: recipients,
		UserID:     userID,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al enviar el corte por correo",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Corte enviado por correo exitosamente",
		"data":    gin.H{"sent_to": recipients},
	})
}
