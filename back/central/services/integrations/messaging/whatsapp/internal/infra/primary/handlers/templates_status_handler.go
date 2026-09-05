package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/ports"
)

func (h *handler) resolveBusinessID(c *gin.Context) (uint, bool) {
	businessID, ok := middleware.GetBusinessIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "contexto de negocio no encontrado",
		})
		return 0, false
	}

	if businessID > 0 {
		return businessID, true
	}

	param := c.Query("business_id")
	if param == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "business_id es requerido para super admin",
		})
		return 0, false
	}

	parsed, err := strconv.ParseUint(param, 10, 64)
	if err != nil || parsed == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "business_id inválido",
		})
		return 0, false
	}

	return uint(parsed), true
}

func (h *handler) GetTemplatesStatus(c *gin.Context) {
	ctx := c.Request.Context()

	if h.templatesUseCase == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"success": false,
			"error":   "el módulo de plantillas no está disponible",
		})
		return
	}

	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		return
	}

	refresh := c.Query("refresh") == "true"

	snapshot, err := h.templatesUseCase.GetStatus(ctx, businessID, refresh)
	if err != nil {
		h.log.Error(ctx).Err(err).
			Uint("business_id", businessID).
			Msg("[WhatsApp Templates] - error consultando estado de plantillas")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    mapSnapshotToResponse(snapshot),
	})
}

func (h *handler) ProvisionTemplates(c *gin.Context) {
	ctx := c.Request.Context()

	if h.templatesUseCase == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"success": false,
			"error":   "el módulo de plantillas no está disponible",
		})
		return
	}

	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		return
	}

	result, err := h.templatesUseCase.Provision(ctx, businessID)
	if err != nil {
		h.log.Error(ctx).Err(err).
			Uint("business_id", businessID).
			Msg("[WhatsApp Templates] - error aprovisionando plantillas")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"integration_id": result.IntegrationID,
			"business_id":    result.BusinessID,
			"waba_id":        result.WABAID,
			"created":        result.Created,
			"already_exists": result.AlreadyExists,
			"failed":         result.Failed,
			"templates":      mapTemplatesToResponse(result.Templates),
		},
	})
}

func mapSnapshotToResponse(snapshot *ports.WABATemplatesSnapshot) gin.H {
	if snapshot == nil {
		return gin.H{"templates": []gin.H{}}
	}

	return gin.H{
		"integration_id": snapshot.IntegrationID,
		"business_id":    snapshot.BusinessID,
		"waba_id":        snapshot.WABAID,
		"refreshed_at":   snapshot.RefreshedAt,
		"templates":      mapTemplatesToResponse(snapshot.Templates),
	}
}

func mapTemplatesToResponse(templates []ports.TemplateStatus) []gin.H {
	out := make([]gin.H, 0, len(templates))
	for _, tpl := range templates {
		out = append(out, gin.H{
			"name":        tpl.Name,
			"language":    tpl.Language,
			"status":      tpl.Status,
			"category":    tpl.Category,
			"meta_id":     tpl.MetaID,
			"reason":      tpl.Reason,
			"provisioned": tpl.Provisioned,
			"updated_at":  tpl.UpdatedAt,
		})
	}
	return out
}
