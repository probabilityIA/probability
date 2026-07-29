package handlerintegrations

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/shared/productmatch"
)

type productMatchRuleRequest struct {
	Probability string `json:"probability" binding:"required"`
	Channel     string `json:"channel" binding:"required"`
}

type updateProductMatchRequest struct {
	BusinessID *uint                     `json:"business_id"`
	Rules      []productMatchRuleRequest `json:"rules" binding:"required"`
}

func (h *IntegrationHandler) resolveProductMatchBusinessID(c *gin.Context, bodyBusinessID *uint) (uint, bool) {
	businessID, ok := middleware.GetBusinessIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "contexto de negocio no encontrado"})
		return 0, false
	}
	if businessID > 0 {
		return businessID, true
	}
	if bodyBusinessID != nil && *bodyBusinessID > 0 {
		return *bodyBusinessID, true
	}
	if param := c.Query("business_id"); param != "" {
		if id, err := strconv.ParseUint(param, 10, 64); err == nil && id > 0 {
			return uint(id), true
		}
	}
	c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "business_id es requerido para super admin"})
	return 0, false
}

func productMatchResponse(cfg *domain.ProductMatchConfig) gin.H {
	return gin.H{
		"integration_id": cfg.IntegrationID,
		"rules":          cfg.Rules,
		"default_rules":  cfg.DefaultRules,
		"options":        cfg.Options,
		"is_override":    cfg.IsOverride,
	}
}

func (h *IntegrationHandler) GetProductMatchConfigHandler(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "id de integracion invalido"})
		return
	}
	businessID, ok := h.resolveProductMatchBusinessID(c, nil)
	if !ok {
		return
	}

	cfg, err := h.usecase.GetProductMatchConfig(c.Request.Context(), uint(id), businessID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": productMatchResponse(cfg)})
}

func (h *IntegrationHandler) UpdateProductMatchConfigHandler(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "id de integracion invalido"})
		return
	}

	var req updateProductMatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "rules es requerido"})
		return
	}
	if len(req.Rules) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "debe seleccionar al menos una regla de match"})
		return
	}

	businessID, ok := h.resolveProductMatchBusinessID(c, req.BusinessID)
	if !ok {
		return
	}

	rules := make([]productmatch.Rule, 0, len(req.Rules))
	for _, r := range req.Rules {
		rules = append(rules, productmatch.Rule{Probability: r.Probability, Channel: r.Channel})
	}

	cfg, err := h.usecase.UpdateProductMatchConfig(c.Request.Context(), uint(id), businessID, rules)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": productMatchResponse(cfg)})
}
