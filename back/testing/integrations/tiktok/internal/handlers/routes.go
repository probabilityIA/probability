package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	router.GET("/health", h.handleHealth)
	router.POST("/v2/oauth/token/", h.handleToken)
	router.POST("/v2/research/tts/shop/", h.handleShopQuery)
}

func (h *Handler) handleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "mock": "tiktok"})
}

func (h *Handler) handleToken(c *gin.Context) {
	clientKey := c.PostForm("client_key")
	clientSecret := c.PostForm("client_secret")
	grantType := c.PostForm("grant_type")

	if clientKey == "" || clientSecret == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":             "invalid_request",
			"error_description": "client_key and client_secret are required",
		})
		return
	}
	if grantType != "client_credentials" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":             "unsupported_grant_type",
			"error_description": "grant_type must be client_credentials",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": MockAccessToken,
		"expires_in":   7200,
		"token_type":   "Bearer",
	})
}

type shopQueryRequest struct {
	ShopName string `json:"shop_name"`
	Fields   string `json:"fields"`
	Limit    int    `json:"limit"`
}

func (h *Handler) handleShopQuery(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") || strings.TrimPrefix(auth, "Bearer ") != MockAccessToken {
		c.JSON(http.StatusUnauthorized, gin.H{
			"data": gin.H{},
			"error": gin.H{
				"code":    "access_token_invalid",
				"message": "The access token is invalid or expired",
			},
		})
		return
	}

	var req shopQueryRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.ShopName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"data": gin.H{},
			"error": gin.H{
				"code":    "invalid_params",
				"message": "shop_name is required",
			},
		})
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	results := make([]shopData, 0, 1)
	if shop, ok := h.shops[req.ShopName]; ok {
		results = append(results, *shop)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"shop_data": results,
		},
		"error": gin.H{
			"code":    "ok",
			"message": "",
		},
	})
}
