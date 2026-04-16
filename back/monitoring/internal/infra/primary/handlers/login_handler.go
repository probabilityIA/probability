package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	domainErrors "github.com/secamc93/probability/back/monitoring/internal/domain/errors"
	"github.com/secamc93/probability/back/monitoring/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/monitoring/internal/infra/primary/handlers/response"
)

func (h *handler) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email and password are required"})
		return
	}

	user, err := h.useCase.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, domainErrors.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	token, err := h.useCase.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, response.LoginResponse{
		Token: token,
		Name:  user.Name,
		Email: user.Email,
	})
}

func (h *handler) Verify(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"valid": true,
		"email": c.GetString("email"),
		"name":  c.GetString("name"),
	})
}
