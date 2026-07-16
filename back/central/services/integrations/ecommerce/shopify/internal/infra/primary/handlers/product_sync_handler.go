package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/domain"
)

func (h *ShopifyHandler) resolveBusinessID(c *gin.Context, bodyBusinessID *uint) (uint, bool) {
	businessID, ok := middleware.GetBusinessIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "contexto de negocio no encontrado"})
		return 0, false
	}
	if businessID == 0 {
		if bodyBusinessID == nil || *bodyBusinessID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "business_id es requerido para super admin"})
			return 0, false
		}
		businessID = *bodyBusinessID
	}
	return businessID, true
}

func productBriefsToResponse(items []domain.ProductBrief) []gin.H {
	out := make([]gin.H, 0, len(items))
	for _, b := range items {
		out = append(out, gin.H{"sku": b.SKU, "name": b.Name})
	}
	return out
}

type productReconcileRequest struct {
	IntegrationID uint  `json:"integration_id" binding:"required"`
	BusinessID    *uint `json:"business_id"`
}

func (h *ShopifyHandler) ReconcileProducts(c *gin.Context) {
	var req productReconcileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "integration_id es requerido"})
		return
	}
	businessID, ok := h.resolveBusinessID(c, req.BusinessID)
	if !ok {
		return
	}

	integrationID := strconv.FormatUint(uint64(req.IntegrationID), 10)
	result, err := h.useCase.ReconcileProducts(c.Request.Context(), integrationID, businessID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":                true,
		"matched":                result.Matched,
		"matched_not_associated": productBriefsToResponse(result.MatchedNotAssociated),
		"only_in_probability":    productBriefsToResponse(result.OnlyInProbability),
		"only_in_shopify":        productBriefsToResponse(result.OnlyInShopify),
		"probability_no_sku":     result.ProbabilityNoSKU,
		"shopify_no_sku":         result.ShopifyNoSKU,
	})
}

type productApplyRequest struct {
	IntegrationID uint   `json:"integration_id" binding:"required"`
	BusinessID    *uint  `json:"business_id"`
	Direction     string `json:"direction" binding:"required"`
}

func (h *ShopifyHandler) ApplyProducts(c *gin.Context) {
	var req productApplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "integration_id y direction son requeridos"})
		return
	}
	if req.Direction != "to_shopify" && req.Direction != "to_probability" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "direction debe ser to_shopify o to_probability"})
		return
	}
	businessID, ok := h.resolveBusinessID(c, req.BusinessID)
	if !ok {
		return
	}

	integrationID := strconv.FormatUint(uint64(req.IntegrationID), 10)
	correlationID := uuid.New().String()
	direction := req.Direction

	go func() {
		ctx := context.Background()
		var err error
		if direction == "to_shopify" {
			err = h.useCase.ApplyProductsToShopify(ctx, integrationID, businessID, correlationID)
		} else {
			err = h.useCase.ApplyProductsToProbability(ctx, integrationID, businessID, correlationID)
		}
		if err != nil {
			h.logger.Error(ctx).Err(err).Str("direction", direction).Msg("Error aplicando reconciliacion de productos Shopify")
		}
	}()

	c.JSON(http.StatusAccepted, gin.H{
		"success":        true,
		"correlation_id": correlationID,
		"message":        "Reconciliacion iniciada",
	})
}

type productAssociateRequest struct {
	IntegrationID uint     `json:"integration_id" binding:"required"`
	BusinessID    *uint    `json:"business_id"`
	Skus          []string `json:"skus"`
}

func (h *ShopifyHandler) AssociateProducts(c *gin.Context) {
	var req productAssociateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "integration_id es requerido"})
		return
	}
	businessID, ok := h.resolveBusinessID(c, req.BusinessID)
	if !ok {
		return
	}

	integrationID := strconv.FormatUint(uint64(req.IntegrationID), 10)
	correlationID := uuid.New().String()
	skus := req.Skus

	go func() {
		ctx := context.Background()
		if err := h.useCase.AssociateProducts(ctx, integrationID, businessID, correlationID, skus); err != nil {
			h.logger.Error(ctx).Err(err).Msg("Error asociando productos a Shopify")
		}
	}()

	c.JSON(http.StatusAccepted, gin.H{
		"success":        true,
		"correlation_id": correlationID,
		"message":        "Asociacion iniciada",
	})
}
