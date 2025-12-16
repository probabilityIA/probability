package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/secamc93/probability/back/central/services/auth/users/internal/domain"
	"github.com/secamc93/probability/back/central/services/auth/users/internal/infra/primary/handlers/mapper"
	"github.com/secamc93/probability/back/central/services/auth/users/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/auth/users/internal/infra/primary/handlers/response"
)

// Createhandlers maneja la solicitud de crear un usuario
//
//	@Summary		Crear usuario
//	@Description	Crea un nuevo usuario.
//	@Description	ENVÍO: SOLO multipart/form-data (no JSON body)
//	@Tags			Users
//	@Accept			mpfd
//	@Produce		json
//	@Security		BearerAuth
//	@Param			name			formData	string	true	"Nombre (2-100)"
//	@Param			email			formData	string	true	"Email válido"
//	@Param			phone			formData	string	false	"Teléfono (exactamente 10 dígitos)"
//	@Param			is_active		formData	boolean	false	"¿Activo? (true/false)"
//	@Param			avatarFile		formData	file	false	"Imagen de avatar (sube a S3)"
//	@Param			business_ids	formData	string	false	"IDs de negocios separados por comas. Ej: '16,21'"
//	@Success		201				{object}	response.UserCreatedResponse	"Usuario creado exitosamente (incluye contraseña generada)"
//	@Failure		400				{object}	response.UserErrorResponse	"Datos inválidos (mensajes claros de validación)"
//	@Failure		401				{object}	response.UserErrorResponse	"Token de acceso requerido"
//	@Failure		409				{object}	response.UserErrorResponse	"Email ya existe"
//	@Failure		500				{object}	response.UserErrorResponse	"Error interno del servidor"
//	@Router			/users [post]
func (h *handlers) Createhandlers(c *gin.Context) {
	var req request.CreateUserRequest
	// Seleccionar binder según Content-Type
	if c.ContentType() == "application/json" {
		if err := c.ShouldBindJSON(&req); err != nil {
			h.logger.Error().Err(err).Msg("Error al validar datos de entrada (JSON)")
			c.JSON(http.StatusBadRequest, response.UserErrorResponse{
				Error: buildValidationMessage(err),
			})
			return
		}
	} else if err := c.ShouldBind(&req); err != nil {
		h.logger.Error().Err(err).Msg("Error al validar datos de entrada")
		c.JSON(http.StatusBadRequest, response.UserErrorResponse{
			Error: buildValidationMessage(err),
		})
		return
	}

	// Si viene multipart/form-data, deserializar business_ids manualmente (JSON en string o separado por comas)
	if raw := req.BusinessIDsRaw; raw != "" {
		var ids []uint
		if err := json.Unmarshal([]byte(raw), &ids); err != nil {
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
		req.BusinessIDs = ids
	}

	// Parsear scope_id manualmente si viene como string desde form-data
	if req.ScopeID == nil {
		if scopeStr := c.PostForm("scope_id"); scopeStr != "" {
			if scopeIDVal, err := strconv.ParseUint(scopeStr, 10, 32); err == nil {
				scopeID := uint(scopeIDVal)
				req.ScopeID = &scopeID
			}
		}
	}

	h.logger.Info().
		Str("email", req.Email).
		Int("business_ids_count", len(req.BusinessIDs)).
		Any("business_ids", req.BusinessIDs).
		Any("scope_id", req.ScopeID).
		Msg("Iniciando solicitud para crear usuario")

	// Log simple para confirmar recepción de archivo
	if req.AvatarFile != nil {
		h.logger.Info().Str("email", req.Email).Str("filename", req.AvatarFile.Filename).Int64("size", req.AvatarFile.Size).Msg("Archivo de avatar recibido en creación")
	}

	userDTO := mapper.ToCreateUserDTO(req)
	email, password, message, err := h.usecase.CreateUser(c.Request.Context(), userDTO)
	if err != nil {
		h.logger.Error().Err(err).Msg("Error al crear usuario desde el caso de uso")

		statusCode := http.StatusInternalServerError
		errorMessage := "Error interno del servidor"

		// Manejar errores de dominio usando errors.Is()
		switch {
		case errors.Is(err, domain.ErrUserEmailExists):
			statusCode = http.StatusConflict
			errorMessage = domain.ErrUserEmailExists.Error()
		case errors.Is(err, domain.ErrBusinessesNotFound):
			statusCode = http.StatusBadRequest
			errorMessage = domain.ErrBusinessesNotFound.Error()
		case errors.Is(err, domain.ErrRolesNotFound):
			statusCode = http.StatusBadRequest
			errorMessage = domain.ErrRolesNotFound.Error()
		case errors.Is(err, domain.ErrUserPasswordError):
			statusCode = http.StatusInternalServerError
			errorMessage = "Error al generar la contraseña del usuario"
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

	h.logger.Info().
		Str("email", email).
		Bool("password_received", password != "").
		Msg("Usuario creado exitosamente - enviando respuesta con contraseña generada")
	c.JSON(http.StatusCreated, response.UserCreatedResponse{
		Success:  true,
		Email:    email,
		Password: password,
		Message:  message,
	})
}

// buildValidationMessage construye mensajes claros a partir de errores de validación
func buildValidationMessage(err error) string {
	var verrs validator.ValidationErrors
	if errors.As(err, &verrs) {
		// Tomar el primer error significativo
		fe := verrs[0]
		field := strings.ToLower(fe.Field())
		switch field {
		case "name":
			if fe.Tag() == "min" || fe.Tag() == "max" {
				return "El nombre debe tener entre 2 y 100 caracteres"
			}
			return "El nombre es inválido"
		case "email":
			if fe.Tag() == "email" {
				return "El email no tiene un formato válido"
			}
			return "El email es inválido"
		case "phone":
			if fe.Tag() == "len" {
				return "El teléfono debe tener exactamente 10 dígitos"
			}
			return "El teléfono es inválido"
		default:
			return "Datos de entrada inválidos"
		}
	}
	return "Datos de entrada inválidos: " + err.Error()
}
