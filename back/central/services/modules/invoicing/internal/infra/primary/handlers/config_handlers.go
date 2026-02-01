package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/primary/handlers/mappers"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/primary/handlers/response"
)

// CreateConfig crea una nueva configuración de facturación
func (h *handler) CreateConfig(c *gin.Context) {
	ctx := c.Request.Context()

	var req request.CreateConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error(ctx).Err(err).Msg("Invalid request body")
		c.JSON(http.StatusBadRequest, response.Error{
			Error:   "invalid_request",
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	logInfo := h.log.Info(ctx).
		Uint("business_id", req.BusinessID).
		Uint("integration_id", req.IntegrationID).
		Uint("invoicing_integration_id", req.InvoicingIntegrationID)

	if req.InvoicingProviderID != nil {
		logInfo = logInfo.Uint("provider_id", *req.InvoicingProviderID)
	}

	logInfo.Msg("Creating invoicing config")

	// Obtener user ID del contexto (JWT)
	userID, exists := c.Get("user_id")
	if !exists {
		h.log.Error(ctx).Msg("User ID not found in context")
		c.JSON(http.StatusUnauthorized, response.Error{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	dto := mappers.CreateConfigRequestToDTO(&req, userID.(uint))

	config, err := h.useCase.CreateConfig(ctx, dto)
	if err != nil {
		h.log.Error(ctx).Err(err).Msg("Failed to create config")
		c.JSON(http.StatusInternalServerError, response.Error{
			Error:   "config_creation_failed",
			Message: err.Error(),
		})
		return
	}

	resp := mappers.ConfigToResponse(config)

	h.log.Info(ctx).
		Uint("config_id", config.ID).
		Msg("Config created successfully")

	c.JSON(http.StatusCreated, resp)
}

// ListConfigs lista configuraciones de facturación
func (h *handler) ListConfigs(c *gin.Context) {
	ctx := c.Request.Context()

	filters := make(map[string]interface{})

	var businessID uint
	if businessIDStr := c.Query("business_id"); businessIDStr != "" {
		id, _ := strconv.ParseUint(businessIDStr, 10, 32)
		businessID = uint(id)
		filters["business_id"] = businessID
	}

	if integrationID := c.Query("integration_id"); integrationID != "" {
		id, _ := strconv.ParseUint(integrationID, 10, 32)
		filters["integration_id"] = uint(id)
	}

	if providerID := c.Query("provider_id"); providerID != "" {
		id, _ := strconv.ParseUint(providerID, 10, 32)
		filters["invoicing_provider_id"] = uint(id)
	}

	if enabled := c.Query("enabled"); enabled != "" {
		filters["enabled"] = enabled == "true"
	}

	if autoInvoice := c.Query("auto_invoice"); autoInvoice != "" {
		filters["auto_invoice"] = autoInvoice == "true"
	}

	page := 1
	if p := c.Query("page"); p != "" {
		page, _ = strconv.Atoi(p)
	}

	pageSize := 20
	if ps := c.Query("page_size"); ps != "" {
		pageSize, _ = strconv.Atoi(ps)
	}

	if pageSize > 100 {
		pageSize = 100
	}

	h.log.Debug(ctx).Uint("business_id", businessID).Msg("Listing configs")

	configs, err := h.useCase.ListConfigs(ctx, businessID)
	if err != nil {
		h.log.Error(ctx).Err(err).Msg("Failed to list configs")
		c.JSON(http.StatusInternalServerError, response.Error{
			Error:   "list_configs_failed",
			Message: err.Error(),
		})
		return
	}

	resp := mappers.ConfigsToResponse(configs, int64(len(configs)), page, pageSize)
	c.JSON(http.StatusOK, resp)
}

// GetConfig obtiene una configuración por ID
func (h *handler) GetConfig(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error{
			Error:   "invalid_id",
			Message: "Invalid config ID",
		})
		return
	}

	h.log.Debug(ctx).Uint("config_id", uint(id)).Msg("Getting config")

	config, err := h.useCase.GetConfig(ctx, uint(id))
	if err != nil {
		h.log.Error(ctx).Err(err).Uint("config_id", uint(id)).Msg("Failed to get config")
		c.JSON(http.StatusNotFound, response.Error{
			Error:   "config_not_found",
			Message: err.Error(),
		})
		return
	}

	resp := mappers.ConfigToResponse(config)
	c.JSON(http.StatusOK, resp)
}

// UpdateConfig actualiza una configuración
func (h *handler) UpdateConfig(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error{
			Error:   "invalid_id",
			Message: "Invalid config ID",
		})
		return
	}

	var req request.UpdateConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error(ctx).Err(err).Msg("Invalid request body")
		c.JSON(http.StatusBadRequest, response.Error{
			Error:   "invalid_request",
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	h.log.Info(ctx).Uint("config_id", uint(id)).Msg("Updating config")

	dto := mappers.UpdateConfigRequestToDTO(&req)

	config, err := h.useCase.UpdateConfig(ctx, uint(id), dto)
	if err != nil {
		h.log.Error(ctx).Err(err).Uint("config_id", uint(id)).Msg("Failed to update config")
		c.JSON(http.StatusInternalServerError, response.Error{
			Error:   "config_update_failed",
			Message: err.Error(),
		})
		return
	}

	resp := mappers.ConfigToResponse(config)

	h.log.Info(ctx).Uint("config_id", config.ID).Msg("Config updated successfully")

	c.JSON(http.StatusOK, resp)
}

// DeleteConfig elimina una configuración
func (h *handler) DeleteConfig(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error{
			Error:   "invalid_id",
			Message: "Invalid config ID",
		})
		return
	}

	h.log.Info(ctx).Uint("config_id", uint(id)).Msg("Deleting config")

	err = h.useCase.DeleteConfig(ctx, uint(id))
	if err != nil {
		h.log.Error(ctx).Err(err).Uint("config_id", uint(id)).Msg("Failed to delete config")
		c.JSON(http.StatusInternalServerError, response.Error{
			Error:   "config_deletion_failed",
			Message: err.Error(),
		})
		return
	}

	h.log.Info(ctx).Uint("config_id", uint(id)).Msg("Config deleted successfully")

	c.JSON(http.StatusOK, response.Success{
		Success: true,
		Message: "Config deleted successfully",
	})
}
