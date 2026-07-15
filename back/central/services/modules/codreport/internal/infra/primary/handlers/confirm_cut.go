package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/modules/codreport/internal/domain/dtos"
)

type draftCutRequest struct {
	PeriodStart string   `json:"period_start"`
	PeriodEnd   string   `json:"period_end"`
	OrderIDs    []string `json:"order_ids"`
}

func (h *Handlers) CreateDraftCut(c *gin.Context) {
	if !isAdminUser(c) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Solo un administrador puede crear cortes de pago"})
		return
	}

	businessID, err := resolveBusinessID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	var req draftCutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Cuerpo de la solicitud invalido"})
		return
	}

	start, err1 := time.Parse("2006-01-02", req.PeriodStart)
	end, err2 := time.Parse("2006-01-02", req.PeriodEnd)
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Fechas del periodo invalidas"})
		return
	}
	end = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 0, time.UTC)

	if len(req.OrderIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Debes seleccionar al menos una orden"})
		return
	}

	userID, _ := middleware.GetUserID(c)

	cut, err := h.uc.CreateDraft(c.Request.Context(), dtos.ConfirmCutDTO{
		BusinessID:  businessID,
		PeriodStart: start,
		PeriodEnd:   end,
		OrderIDs:    req.OrderIDs,
		UserID:      userID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al crear el borrador del corte",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Borrador de corte creado exitosamente",
		"data":    mapCut(cut),
	})
}

func (h *Handlers) ConfirmCut(c *gin.Context) {
	if !isAdminUser(c) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Solo un administrador puede confirmar cortes de pago"})
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

	userID, _ := middleware.GetUserID(c)

	if err := h.uc.ConfirmCut(c.Request.Context(), businessID, uint(cutID), userID, ""); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al confirmar el corte de pago",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Corte de pago confirmado exitosamente"})
}
