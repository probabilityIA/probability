package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

type ProductHandler struct {
	useCase ports.IInvoiceUseCase
	log     log.ILogger
}

func NewProductHandler(useCase ports.IInvoiceUseCase, logger log.ILogger) *ProductHandler {
	return &ProductHandler{
		useCase: useCase,
		log:     logger.WithModule("siigo.product_handler"),
	}
}

func (h *ProductHandler) RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/siigo")
	{
		group.POST("/products/reconcile", middleware.JWT(), h.ReconcileProducts)
		group.POST("/products/apply", middleware.JWT(), h.ApplyProducts)
	}
}

func (h *ProductHandler) resolveBusinessID(c *gin.Context, bodyBusinessID *uint) (uint, bool) {
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

type productReconcileRequest struct {
	IntegrationID uint  `json:"integration_id" binding:"required"`
	BusinessID    *uint `json:"business_id"`
}

func productBriefsToResponse(items []dtos.ProductBrief) []gin.H {
	out := make([]gin.H, 0, len(items))
	for _, b := range items {
		out = append(out, gin.H{"sku": b.SKU, "name": b.Name})
	}
	return out
}

func (h *ProductHandler) ReconcileProducts(c *gin.Context) {
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
		"only_in_siigo":          productBriefsToResponse(result.OnlyInSiigo),
		"probability_no_sku":     result.ProbabilityNoSKU,
		"siigo_no_sku":           result.SiigoNoSKU,
	})
}

type productApplyRequest struct {
	IntegrationID uint     `json:"integration_id" binding:"required"`
	BusinessID    *uint    `json:"business_id"`
	Skus          []string `json:"skus"`
}

func (h *ProductHandler) ApplyProducts(c *gin.Context) {
	var req productApplyRequest
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
		if err := h.useCase.ApplyProductsToProbability(ctx, integrationID, businessID, correlationID, skus); err != nil {
			h.log.Error(ctx).Err(err).Msg("Error aplicando sincronizacion de productos Siigo")
		}
	}()

	c.JSON(http.StatusAccepted, gin.H{
		"success":        true,
		"correlation_id": correlationID,
		"message":        "Sincronizacion iniciada",
	})
}
