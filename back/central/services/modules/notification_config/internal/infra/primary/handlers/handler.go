package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers/mappers"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers/request"
)

// Create godoc
// @Summary Crear configuración de notificación
// @Description Crea una nueva configuración de notificación para una integración
// @Tags notification-config
// @Accept json
// @Produce json
// @Param body body request.CreateNotificationConfig true "Datos de la configuración"
// @Success 201 {object} response.NotificationConfig
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/integrations/notification-configs [post]
func (h *handler) Create(c *gin.Context) {
	var req request.CreateNotificationConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error().Err(err).Msg("Invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convertir request HTTP a DTO de dominio usando mapper
	dto := mappers.CreateRequestToDomain(&req)

	result, err := h.useCase.Create(c.Request.Context(), dto)
	if err != nil {
		if err == errors.ErrDuplicateConfig {
			c.JSON(http.StatusConflict, gin.H{"error": "A similar notification config already exists"})
			return
		}
		h.logger.Error().Err(err).Msg("Error creating notification config")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Convertir DTO de dominio a response HTTP usando mapper
	response := mappers.DomainToResponse(*result)
	c.JSON(http.StatusCreated, response)
}

// List godoc
// @Summary Listar configuraciones de notificación
// @Description Obtiene una lista de configuraciones con filtros opcionales
// @Tags notification-config
// @Accept json
// @Produce json
// @Param integration_id query uint false "ID de la integración"
// @Param notification_type query string false "Tipo de notificación (whatsapp, email, sms)"
// @Param is_active query bool false "Filtrar por activas/inactivas"
// @Param trigger query string false "Filtrar por trigger"
// @Success 200 {array} response.NotificationConfig
// @Failure 500 {object} map[string]interface{}
// @Router /api/integrations/notification-configs [get]
func (h *handler) List(c *gin.Context) {
	var query request.FilterNotificationConfig
	if err := c.ShouldBindQuery(&query); err != nil {
		h.logger.Error().Err(err).Msg("Invalid query parameters")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convertir query params a DTO de dominio usando mapper
	filters := mappers.FilterRequestToDomain(&query)

	result, err := h.useCase.List(c.Request.Context(), filters)
	if err != nil {
		h.logger.Error().Err(err).Msg("Error listing notification configs")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Convertir lista de DTOs de dominio a responses HTTP usando mapper
	responses := mappers.DomainListToResponse(result)
	c.JSON(http.StatusOK, responses)
}

// GetByID godoc
// @Summary Obtener configuración por ID
// @Description Obtiene una configuración de notificación por su ID
// @Tags notification-config
// @Accept json
// @Produce json
// @Param id path uint true "ID de la configuración"
// @Success 200 {object} response.NotificationConfig
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/integrations/notification-configs/{id} [get]
func (h *handler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	result, err := h.useCase.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		if err == errors.ErrNotificationConfigNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Notification config not found"})
			return
		}
		h.logger.Error().Err(err).Msg("Error getting notification config")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Convertir DTO de dominio a response HTTP usando mapper
	response := mappers.DomainToResponse(*result)
	c.JSON(http.StatusOK, response)
}

// Update godoc
// @Summary Actualizar configuración
// @Description Actualiza una configuración de notificación existente
// @Tags notification-config
// @Accept json
// @Produce json
// @Param id path uint true "ID de la configuración"
// @Param body body request.UpdateNotificationConfig true "Datos a actualizar"
// @Success 200 {object} response.NotificationConfig
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/integrations/notification-configs/{id} [put]
func (h *handler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req request.UpdateNotificationConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error().Err(err).Msg("Invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convertir request HTTP a DTO de dominio usando mapper
	dto := mappers.UpdateRequestToDomain(&req)

	result, err := h.useCase.Update(c.Request.Context(), uint(id), dto)
	if err != nil {
		if err == errors.ErrNotificationConfigNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Notification config not found"})
			return
		}
		h.logger.Error().Err(err).Msg("Error updating notification config")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Convertir DTO de dominio a response HTTP usando mapper
	response := mappers.DomainToResponse(*result)
	c.JSON(http.StatusOK, response)
}

// Delete godoc
// @Summary Eliminar configuración
// @Description Elimina una configuración de notificación
// @Tags notification-config
// @Accept json
// @Produce json
// @Param id path uint true "ID de la configuración"
// @Success 204
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/integrations/notification-configs/{id} [delete]
func (h *handler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.useCase.Delete(c.Request.Context(), uint(id)); err != nil {
		if err == errors.ErrNotificationConfigNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Notification config not found"})
			return
		}
		h.logger.Error().Err(err).Msg("Error deleting notification config")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Status(http.StatusNoContent)
}
