package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/testing/integrations/shipit/internal/domain"
)

func shipitError(message string) gin.H {
	return gin.H{"state": "error", "message": message}
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/v")
	api.Use(h.authMiddleware())

	api.GET("/communes", h.handleCommunes)
	api.POST("/rates", h.handleRates)
	api.POST("/shipments", h.handleCreateShipment)
	api.GET("/trackings/number/:number", h.handleTrack)
	api.GET("/webhooks/package", h.handleWebhookConfig)

	router.POST("/simulate/status", h.handleSimulateStatus)
	router.GET("/labels/:id", h.handleLabel)
	router.GET("/health", h.handleHealth)
}

func (h *Handler) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		email := c.GetHeader("X-Shipit-Email")
		token := c.GetHeader("X-Shipit-Access-Token")
		if email == "" || token == "" {
			c.JSON(401, shipitError("Missing X-Shipit-Email or X-Shipit-Access-Token headers"))
			c.Abort()
			return
		}
		c.Next()
	}
}

func (h *Handler) handleHealth(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "ok",
		"service": "shipit-mock",
	})
}

func (h *Handler) handleCommunes(c *gin.Context) {
	c.JSON(200, h.apiSimulator.GetCommunes())
}

func (h *Handler) handleRates(c *gin.Context) {
	var req domain.RatesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, shipitError("Invalid request format"))
		return
	}

	resp, err := h.apiSimulator.HandleRates(req)
	if err != nil {
		c.JSON(400, shipitError(err.Error()))
		return
	}

	c.JSON(201, resp)
}

func (h *Handler) handleCreateShipment(c *gin.Context) {
	var req domain.ShipmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, shipitError("Invalid request format"))
		return
	}

	resp, err := h.apiSimulator.HandleCreateShipment(req)
	if err != nil {
		c.JSON(400, shipitError(err.Error()))
		return
	}

	c.JSON(200, resp)
}

func (h *Handler) handleTrack(c *gin.Context) {
	number := c.Param("number")

	resp, err := h.apiSimulator.HandleTrack(number)
	if err != nil {
		c.JSON(400, shipitError(err.Error()))
		return
	}

	c.JSON(200, resp)
}

func (h *Handler) handleWebhookConfig(c *gin.Context) {
	c.JSON(200, gin.H{
		"webhook": gin.H{
			"url": "",
			"options": gin.H{
				"sign_body":     gin.H{"required": false, "token": ""},
				"authorization": gin.H{"required": false, "kind": "", "token": ""},
			},
		},
	})
}

func (h *Handler) handleSimulateStatus(c *gin.Context) {
	var req struct {
		TrackingNumber string `json:"tracking_number"`
		Status         string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.TrackingNumber == "" || req.Status == "" {
		c.JSON(400, shipitError("tracking_number y status son requeridos"))
		return
	}

	stored, err := h.apiSimulator.SimulateStatusChange(req.TrackingNumber, req.Status)
	if err != nil {
		c.JSON(404, shipitError(err.Error()))
		return
	}

	c.JSON(200, gin.H{
		"success":         true,
		"tracking_number": stored.TrackingNumber,
		"status":          stored.Status,
		"statuses":        len(stored.Statuses),
	})
}

func (h *Handler) handleLabel(c *gin.Context) {
	id := c.Param("id")
	content := fmt.Sprintf("%%PDF-1.4\nShipit mock label %s\n%%%%EOF\n", id)
	c.Data(200, "application/pdf", []byte(content))
}
