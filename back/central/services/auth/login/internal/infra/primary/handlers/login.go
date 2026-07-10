package authhandler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/login/internal/domain"
	"github.com/secamc93/probability/back/central/services/auth/login/internal/infra/primary/handlers/mapper"
	"github.com/secamc93/probability/back/central/services/auth/login/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/auth/login/internal/infra/primary/handlers/response"
	"github.com/secamc93/probability/back/central/shared/log"
)

func (h *AuthHandler) LoginHandler(c *gin.Context) {
	ctx := log.WithFunctionCtx(c.Request.Context(), "LoginHandler")

	var loginRequest request.LoginRequest

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		h.logger.Error(ctx).Err(err).Msg("Error al validar request de login")
		c.JSON(http.StatusBadRequest, response.LoginBadRequestResponse{
			Error:   "Datos de entrada inválidos",
			Details: err.Error(),
		})
		return
	}

	domainRequest := domain.LoginRequest{
		Email:    loginRequest.Email,
		Password: loginRequest.Password,
	}

	domainResponse, err := h.usecase.Login(ctx, domainRequest)
	if err != nil {
		h.logger.Error(ctx).Err(err).Str("email", loginRequest.Email).Msg("Error en proceso de login")

		statusCode := http.StatusInternalServerError
		errorMessage := "Error interno del servidor"

		if errors.Is(err, domain.ErrUserPendingVerification) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": domain.ErrUserPendingVerification.Error(),
				"code":  "USER_PENDING_VERIFICATION",
				"email": loginRequest.Email,
			})
			return
		}

		switch {
		case errors.Is(err, domain.ErrInvalidCredentials):
			statusCode = http.StatusUnauthorized
			errorMessage = domain.ErrInvalidCredentials.Error()
		case errors.Is(err, domain.ErrUserNotFound):
			statusCode = http.StatusUnauthorized
			errorMessage = domain.ErrInvalidCredentials.Error()
		case errors.Is(err, domain.ErrUserInactive):
			statusCode = http.StatusForbidden
			errorMessage = domain.ErrUserInactive.Error()
		case errors.Is(err, domain.ErrEmailPasswordRequired):
			statusCode = http.StatusBadRequest
			errorMessage = domain.ErrEmailPasswordRequired.Error()
		}

		c.JSON(statusCode, response.LoginErrorResponse{
			Error: errorMessage,
		})
		return
	}

	loginResponse := mapper.ToLoginResponse(domainResponse)

	clientType := c.GetHeader("X-Client-Type")
	isMobileClient := clientType == "mobile" || clientType == "api"

	if !isMobileClient {
		cookieValue := fmt.Sprintf(
			"%s=%s; Max-Age=%d; Path=%s; Domain=%s; Secure; HttpOnly; SameSite=None; Partitioned",
			"session_token",
			domainResponse.Token,
			7*24*60*60,
			"/",
			".probabilityia.com.co",
		)
		c.Header("Set-Cookie", cookieValue)

		loginResponse.Token = ""
	}

	h.logger.Info(ctx).
		Str("email", loginRequest.Email).
		Uint("user_id", domainResponse.User.ID).
		Str("scope", domainResponse.Scope).
		Bool("is_super_admin", domainResponse.IsSuperAdmin).
		Str("client_type", clientType).
		Msg("Login exitoso")

	c.JSON(http.StatusOK, response.LoginSuccessResponse{
		Success: true,
		Data:    *loginResponse,
	})
}
