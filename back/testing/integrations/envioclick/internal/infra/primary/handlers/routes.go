package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/testing/integrations/envioclick/internal/domain"
)

// RegisterRoutes registra todas las rutas del simulador de EnvioClick
func (h *Handler) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api/v2")
	api.Use(h.authMiddleware())

	api.POST("/quotation", h.handleQuote)
	api.POST("/shipment", h.handleGenerate)
	api.POST("/track", h.handleTrack)
	api.DELETE("/shipment/:id", h.handleCancel)

	// Health check (no auth required)
	router.GET("/health", h.handleHealth)
}

// authMiddleware validates the Authorization header (accepts any non-empty value)
func (h *Handler) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.JSON(401, envioClickError("Missing authorization token"))
			c.Abort()
			return
		}
		c.Next()
	}
}

// handleHealth maneja el health check
func (h *Handler) handleHealth(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "ok",
		"service": "envioclick-mock",
	})
}

// handleQuote handles POST /api/v2/quotation
func (h *Handler) handleQuote(c *gin.Context) {
	var req domain.QuoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, envioClickError("Invalid request format"))
		return
	}

	resp, err := h.apiSimulator.HandleQuote(req)
	if err != nil {
		c.JSON(422, envioClickError(err.Error()))
		return
	}

	c.JSON(200, resp)
}

// handleGenerate handles POST /api/v2/shipment
func (h *Handler) handleGenerate(c *gin.Context) {
	var req domain.QuoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, envioClickError("Invalid request format"))
		return
	}

	resp, err := h.apiSimulator.HandleGenerate(req)
	if err != nil {
		c.JSON(422, envioClickError(err.Error()))
		return
	}

	c.JSON(200, resp)
}

// handleTrack handles POST /api/v2/track
func (h *Handler) handleTrack(c *gin.Context) {
	var req domain.TrackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, envioClickError("Invalid request format"))
		return
	}

	resp, err := h.apiSimulator.HandleTrack(req.TrackingCode)
	if err != nil {
		c.JSON(404, envioClickError(err.Error()))
		return
	}

	c.JSON(200, resp)
}

// handleCancel handles DELETE /api/v2/shipment/:id
func (h *Handler) handleCancel(c *gin.Context) {
	id := c.Param("id")

	resp, err := h.apiSimulator.HandleCancel(id)
	if err != nil {
		c.JSON(404, envioClickError(err.Error()))
		return
	}

	c.JSON(200, resp)
}

// envioClickError formats errors in the same structure the real API uses
// The real client parses: {"status_messages": [{"error": ["message"]}]}
func envioClickError(msg string) gin.H {
	return gin.H{
		"status_messages": []gin.H{
			{
				"error": []string{msg},
			},
		},
	}
}
