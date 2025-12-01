package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/auth/users/internal/domain"
	"github.com/secamc93/probability/back/central/services/auth/users/internal/infra/primary/handlers/mapper"
	"github.com/secamc93/probability/back/central/services/auth/users/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/auth/users/internal/infra/primary/handlers/response"
	"github.com/secamc93/probability/back/central/shared/log"
)

// Updatehandlers maneja la solicitud de actualizar un usuario
//
//	@Summary		Actualizar usuario
//	@Description	Actualiza un usuario existente.
//	@Description	ENVÍO: SOLO multipart/form-data (no JSON body)
//	@Tags			Users
//	@Accept			mpfd
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id				path		int		true	"ID del usuario"
//	@Param			name			formData	string	false	"Nombre (2-100)"
//	@Param			email			formData	string	false	"Email válido"
//	@Param			phone			formData	string	false	"Teléfono (exactamente 10 dígitos)"
//	@Param			is_active		formData	boolean	false	"¿Activo?"
//	@Param			remove_avatar	formData	boolean	false	"Eliminar avatar actual (true/false)"
//	@Param			avatarFile		formData	file	false	"Imagen de avatar (sube a S3)"
//	@Param			business_ids	formData	string	false	"IDs de negocios separados por comas. Ej: '16,21'"
//	@Success		200				{object}	response.UserSuccessResponse	"Usuario actualizado (data con usuario)"
//	@Failure		400				{object}	response.UserErrorResponse	"Datos inválidos o IDs inexistentes"
//	@Failure		401				{object}	response.UserErrorResponse	"Token de acceso requerido"
//	@Failure		404				{object}	response.UserErrorResponse	"Usuario no encontrado"
//	@Failure		409				{object}	response.UserErrorResponse	"Email ya existe"
//	@Failure		500				{object}	response.UserErrorResponse	"Error interno del servidor"
//	@Router			/users/{id} [put]
func (h *handlers) Updatehandlers(c *gin.Context) {
	ctx := log.WithFunctionCtx(c.Request.Context(), "Updatehandlers")
	var uriReq request.GetUserByIDRequest
	var bodyReq request.UpdateUserRequest

	// Binding automático para parámetros de URL
	if err := c.ShouldBindUri(&uriReq); err != nil {
		h.logger.Error(ctx).Err(err).Msg("Error al validar ID del usuario")
		c.JSON(http.StatusBadRequest, response.UserErrorResponse{
			Error: "ID inválido: " + err.Error(),
		})
		return
	}

	// Binding automático para el body (multipart/form-data o JSON)
	if c.ContentType() == "application/json" {
		if err := c.ShouldBindJSON(&bodyReq); err != nil {
			h.logger.Error(ctx).Err(err).Msg("Error al validar datos de la solicitud (JSON)")
			c.JSON(http.StatusBadRequest, response.UserErrorResponse{
				Error: buildUpdateValidationMessage(err),
			})
			return
		}
	} else if err := c.ShouldBind(&bodyReq); err != nil {
		h.logger.Error(ctx).Err(err).Msg("Error al validar datos de la solicitud")
		c.JSON(http.StatusBadRequest, response.UserErrorResponse{
			Error: buildUpdateValidationMessage(err),
		})
		return
	}

	// Si viene multipart/form-data, deserializar business_ids manualmente (JSON en string o separado por comas)
	if raw := bodyReq.BusinessIDsRaw; raw != "" {
		var ids []uint
		if err := json.Unmarshal([]byte(raw), &ids); err != nil {
			// intentar coma separada
			parts := strings.Split(raw, ",")
			for _, p := range parts {
				p = strings.TrimSpace(p)
				if p == "" {
					continue
				}
				if idv, err := strconv.ParseUint(p, 10, 32); err == nil {
					ids = append(ids, uint(idv))
				}
			}
		}
		bodyReq.BusinessIDs = ids
	}

	// Si no es super admin, validar que solo se asignen businesses del token
	isSuperAdmin := middleware.IsSuperAdmin(c)
	if !isSuperAdmin && len(bodyReq.BusinessIDs) > 0 {
		tokenBusinessID, ok := middleware.GetBusinessID(c)
		if !ok {
			h.logger.Error(ctx).Msg("Business ID no disponible en token")
			c.JSON(http.StatusUnauthorized, response.UserErrorResponse{
				Error: "Token inválido: business_id no disponible",
			})
			return
		}
		filtered := make([]uint, 0)
		for _, bid := range bodyReq.BusinessIDs {
			if bid == tokenBusinessID {
				filtered = append(filtered, bid)
			}
		}
		bodyReq.BusinessIDs = filtered
		h.logger.Info(ctx).Uint("business_id", tokenBusinessID).Int("business_ids_count", len(filtered)).Msg("Filtrando businesses al del token para usuario normal")
	}

	h.logger.Info(ctx).Uint("id", uriReq.ID).Str("email", bodyReq.Email).Bool("is_super_admin", isSuperAdmin).Msg("Iniciando solicitud para actualizar usuario")

	// Log para verificar si el archivo está llegando
	if bodyReq.AvatarFile != nil {
		h.logger.Info(ctx).Uint("id", uriReq.ID).Str("filename", bodyReq.AvatarFile.Filename).Int64("size", bodyReq.AvatarFile.Size).Msg("Archivo de avatar recibido")
	} else {
		h.logger.Info(ctx).Uint("id", uriReq.ID).Msg("No se recibió archivo de avatar")
	}

	h.logger.Info(ctx).Int("business_ids_count", len(bodyReq.BusinessIDs)).Any("business_ids", bodyReq.BusinessIDs).Msg("Businesses recibidos para actualización de usuario")
	userDTO := mapper.ToUpdateUserDTO(bodyReq)
	message, err := h.usecase.UpdateUser(ctx, uriReq.ID, userDTO)
	if err != nil {
		h.logger.Error(ctx).Err(err).Uint("id", uriReq.ID).Msg("Error al actualizar usuario desde el caso de uso")

		statusCode := http.StatusInternalServerError
		errorMessage := "Error interno del servidor"

		// Manejar errores de dominio usando errors.Is()
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			statusCode = http.StatusNotFound
			errorMessage = domain.ErrUserNotFound.Error()
		case errors.Is(err, domain.ErrUserEmailExists):
			statusCode = http.StatusConflict
			errorMessage = domain.ErrUserEmailExists.Error()
		case errors.Is(err, domain.ErrBusinessesNotFound):
			statusCode = http.StatusBadRequest
			errorMessage = domain.ErrBusinessesNotFound.Error()
		case errors.Is(err, domain.ErrRolesNotFound):
			statusCode = http.StatusBadRequest
			errorMessage = domain.ErrRolesNotFound.Error()
		case errors.Is(err, domain.ErrUserAvatarUploadFailed):
			statusCode = http.StatusInternalServerError
			errorMessage = "Error al subir la imagen de avatar"
		default:
			// Para errores wrapped, verificar el mensaje
			if strings.Contains(err.Error(), "algunos businesses no existen") {
				statusCode = http.StatusBadRequest
				errorMessage = "Algunos businesses no existen"
			} else if strings.Contains(err.Error(), "algunos roles no existen") {
				statusCode = http.StatusBadRequest
				errorMessage = "Algunos roles no existen"
			}
		}

		c.JSON(statusCode, response.UserErrorResponse{
			Error: errorMessage,
		})
		return
	}

	h.logger.Info(ctx).Uint("id", uriReq.ID).Msg("Usuario actualizado exitosamente")
	// Devolver el usuario actualizado completo en data para compatibilidad con el front
	updatedUser, err := h.usecase.GetUserByID(ctx, uriReq.ID)
	if err != nil || updatedUser == nil {
		c.JSON(http.StatusOK, response.UserMessageResponse{
			Success: true,
			Message: message,
		})
		return
	}

	userResp := mapper.ToUserResponse(*updatedUser)
	c.JSON(http.StatusOK, response.UserSuccessResponse{
		Success: true,
		Data:    userResp,
	})
}

// buildUpdateValidationMessage construye mensajes claros a partir de errores de validación
func buildUpdateValidationMessage(err error) string {
	var verrs validator.ValidationErrors
	if errors.As(err, &verrs) {
		fe := verrs[0]
		field := strings.ToLower(fe.Field())
		switch field {
		case "email":
			if fe.Tag() == "email" {
				return "El email no tiene un formato válido"
			}
		case "phone":
			if fe.Tag() == "len" {
				return "El teléfono debe tener exactamente 10 dígitos"
			}
		case "password":
			if fe.Tag() == "min" {
				return "La contraseña debe tener al menos 6 caracteres"
			}
		}
		return "Datos de entrada inválidos"
	}
	return "Datos de entrada inválidos: " + err.Error()
}
